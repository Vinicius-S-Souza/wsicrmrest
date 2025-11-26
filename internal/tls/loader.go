package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadEncryptedTLSConfig carrega certificado e chave privada (criptografada ou não)
// e retorna uma configuração TLS pronta para uso com http.Server
//
// Parâmetros:
//   - certFile: Caminho do arquivo de certificado (.crt ou .pem)
//   - keyFile: Caminho do arquivo de chave privada (.key ou .pem)
//   - password: Senha para descriptografar a chave (vazio se chave não for criptografada)
//
// Retorna:
//   - *tls.Config configurado e pronto para uso
//   - error se houver falha no carregamento ou descriptografia
//
// Suporta automaticamente:
//   - Chaves não criptografadas (BEGIN PRIVATE KEY, BEGIN RSA PRIVATE KEY)
//   - Chaves criptografadas PKCS#1 (BEGIN RSA PRIVATE KEY + Proc-Type: 4,ENCRYPTED)
//   - Chaves criptografadas PKCS#8 (BEGIN ENCRYPTED PRIVATE KEY)
//
// Data de criação: 2025-11-26
func LoadEncryptedTLSConfig(certFile, keyFile, password string) (*tls.Config, error) {
	// Ler arquivo de certificado
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler certificado %s: %w", certFile, err)
	}

	// Ler arquivo de chave privada
	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler chave privada %s: %w", keyFile, err)
	}

	// Decodificar bloco PEM da chave
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("falha ao decodificar bloco PEM da chave privada")
	}

	// Detectar se a chave está criptografada
	isEncrypted := x509.IsEncryptedPEMBlock(keyBlock) || keyBlock.Type == "ENCRYPTED PRIVATE KEY"

	var keyDER []byte

	if isEncrypted {
		// Chave criptografada - requer senha
		if password == "" {
			return nil, fmt.Errorf("chave privada está criptografada mas nenhuma senha foi fornecida. Configure key_password em [tls] no dbinit.ini")
		}

		// Descriptografar usando x509.DecryptPEMBlock (funciona para PKCS#1 e PKCS#8)
		decrypted, err := x509.DecryptPEMBlock(keyBlock, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("falha ao descriptografar chave privada (senha incorreta?): %w", err)
		}
		keyDER = decrypted
	} else {
		// Chave não criptografada
		if password != "" {
			// Aviso: senha fornecida mas chave não é criptografada
			// Não é erro fatal, apenas continuamos
		}
		keyDER = keyBlock.Bytes
	}

	// Tentar parsear a chave descriptografada como diferentes formatos
	var privateKey interface{}

	// Tentar PKCS#8 primeiro (formato mais moderno)
	privateKey, err = x509.ParsePKCS8PrivateKey(keyDER)
	if err != nil {
		// Se falhar, tentar PKCS#1 (RSA tradicional)
		privateKey, err = x509.ParsePKCS1PrivateKey(keyDER)
		if err != nil {
			// Se ainda falhar, tentar EC (Elliptic Curve)
			privateKey, err = x509.ParseECPrivateKey(keyDER)
			if err != nil {
				return nil, fmt.Errorf("falha ao parsear chave privada (formato não suportado): %w", err)
			}
		}
	}

	// Parsear certificado
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("falha ao decodificar bloco PEM do certificado")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("falha ao parsear certificado: %w", err)
	}

	// Criar tls.Certificate
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
	}

	// Criar configuração TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12, // TLS 1.2 mínimo (segurança)
		CipherSuites: []uint16{
			// Cipher suites recomendadas (seguras e compatíveis)
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	return tlsConfig, nil
}
