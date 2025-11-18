# Guia de Segurança - WSICRMREST

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-17

Este documento descreve as implementações de segurança do WSICRMREST e como configurá-las adequadamente.

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [SQL Injection Protection](#sql-injection-protection)
3. [HTTPS/TLS](#httpstls)
4. [Rate Limiting](#rate-limiting)
5. [CORS](#cors)
6. [Request Limits](#request-limits)
7. [Checklist de Produção](#checklist-de-produção)
8. [Troubleshooting](#troubleshooting)

---

## Visão Geral

O WSICRMREST implementa múltiplas camadas de segurança:

| Camada | Proteção | Status |
|--------|----------|--------|
| **SQL Injection** | Bind variables parametrizadas | ✅ Implementado |
| **HTTPS/TLS** | Criptografia de transporte | ✅ Implementado |
| **Rate Limiting** | Proteção contra brute-force e DoS | ✅ Implementado |
| **CORS** | Controle de origens permitidas | ✅ Implementado |
| **Request Limits** | Timeout e tamanho máximo | ✅ Implementado |
| **Security Headers** | X-Frame-Options, HSTS, etc. | ✅ Implementado |
| **Timing Attack Protection** | Comparação de secrets em tempo constante | ✅ Implementado |

---

## SQL Injection Protection

### O que foi implementado

Todas as queries SQL agora usam **bind variables** (`:1`, `:2`, etc.) em vez de concatenação de strings.

**Antes (VULNERÁVEL):**
```go
query := fmt.Sprintf("UPDATE clientes SET clicelular = '' WHERE clicodigo = %d", cliCodigo)
```

**Depois (SEGURO):**
```go
query := "UPDATE clientes SET clicelular = '' WHERE clicodigo = :1"
db.Exec(query, cliCodigo)
```

### Arquivos corrigidos

- `internal/database/log.go` - Todas as queries de logging
- `internal/database/webhook.go` - Todas as queries de webhooks (10 funções)
- `internal/handlers/token.go` - Query de log de tokens

### Proteção adicional

- **Validação de nomes de colunas**: Em `LimpaEmailInconsistente()`, apenas `CliExtEmail2` e `CliExtEmail3` são permitidos
- **Sanitização removida**: Não é mais necessário `utils.SanitizeForSQL()` pois bind variables são seguros

### Timing Attack Protection

Comparação de `client_secret` agora usa `crypto/subtle.ConstantTimeCompare()`:

```go
func constantTimeCompare(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
```

Isso previne ataques que tentam descobrir secrets medindo o tempo de resposta.

---

## HTTPS/TLS

### Configuração

Edite `dbinit.ini`:

```ini
[tls]
; Habilitar HTTPS/TLS
enabled = true

; Caminhos dos certificados
cert_file = certs/server.crt
key_file = certs/server.key

; Porta HTTPS (padrão: 8443)
port = 8443
```

### Gerar certificados auto-assinados (desenvolvimento)

```bash
# Criar diretório
mkdir -p certs

# Gerar certificado auto-assinado válido por 365 dias
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt \
  -days 365 -nodes \
  -subj "/C=BR/ST=SP/L=SaoPaulo/O=MyCompany/CN=localhost"
```

### Obter certificado válido (produção)

**Opção 1: Let's Encrypt (gratuito)**

```bash
# Instalar Certbot
sudo apt-get install certbot

# Obter certificado (requer domínio válido)
sudo certbot certonly --standalone -d seu-dominio.com

# Certificados estarão em:
# /etc/letsencrypt/live/seu-dominio.com/fullchain.pem
# /etc/letsencrypt/live/seu-dominio.com/privkey.pem
```

Atualizar `dbinit.ini`:
```ini
cert_file = /etc/letsencrypt/live/seu-dominio.com/fullchain.pem
key_file = /etc/letsencrypt/live/seu-dominio.com/privkey.pem
```

**Opção 2: Certificado comercial**

Compre certificado de uma CA (DigiCert, GlobalSign, etc.) e configure os caminhos.

### Headers de segurança HTTPS

Quando TLS está habilitado, o servidor adiciona automaticamente:

- `Strict-Transport-Security: max-age=31536000; includeSubDomains` (HSTS)
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: no-referrer`

### Testar HTTPS

```bash
# Com certificado auto-assinado (aceitar warning)
curl -k https://localhost:8443/health

# Com certificado válido
curl https://seu-dominio.com:8443/health
```

---

## Rate Limiting

### Configuração

Edite `dbinit.ini`:

```ini
[security]
; Habilitar rate limiting
rate_limit_enabled = true

; Limite de requests por minuto (por IP)
rate_limit_per_min = 60

; Limite de requests por hora (por IP)
rate_limit_per_hour = 1000
```

### Como funciona

- **Por IP**: Cada endereço IP tem seu próprio contador
- **Janela deslizante**: Contadores resetam após 1 minuto/hora
- **Em memória**: Não requer Redis (limpar automaticamente)

### Headers de resposta

Toda resposta inclui headers informativos:

```
X-RateLimit-Limit-Minute: 60
X-RateLimit-Limit-Hour: 1000
X-RateLimit-Remaining-Minute: 45
X-RateLimit-Remaining-Hour: 980
```

### Resposta quando limite excedido

**Status:** `429 Too Many Requests`

```json
{
  "code": "429",
  "message": "Rate limit exceeded. Please try again later."
}
```

### Valores recomendados

| Ambiente | Per Minute | Per Hour | Uso |
|----------|------------|----------|-----|
| **Desenvolvimento** | 0 (desabilitado) | 0 | Testes sem limites |
| **Staging** | 100 | 2000 | Testes realistas |
| **Produção** | 60 | 1000 | Balanceado |
| **Alta carga** | 120 | 5000 | APIs públicas |
| **Restrito** | 10 | 100 | Endpoints sensíveis |

### Desabilitar rate limiting

```ini
[security]
rate_limit_enabled = false
```

---

## CORS

### Configuração

**Desenvolvimento (permite todas as origens):**

```ini
[CORS]
AllowedOrigins=
```

**Produção (origens específicas):**

```ini
[CORS]
AllowedOrigins=https://app.example.com,https://admin.example.com
```

### Avisos de segurança

⚠️ **AVISO CRÍTICO**: Se `AllowedOrigins` estiver vazio em **production**, você verá:

```
WARN ⚠️  CORS configurado para permitir TODAS as origens (*) em PRODUÇÃO!
         Isso é um risco de segurança.
WARN Configure AllowedOrigins no dbinit.ini para restringir as origens permitidas.
```

### Bloqueio de origens

Quando uma origem não permitida tenta acessar:

```
WARN CORS: Origem não permitida bloqueada
     origin=https://malicious-site.com ip=192.168.1.100
```

A resposta **não inclui** header `Access-Control-Allow-Origin`, então o navegador bloqueia.

### Configuração completa

```ini
[CORS]
; Origens permitidas (separadas por vírgula)
AllowedOrigins=https://app.example.com,https://admin.example.com

; Métodos HTTP permitidos
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS

; Headers permitidos
AllowedHeaders=Origin,Content-Type,Content-Length,Accept-Encoding,Authorization,Grant_type,X-CSRF-Token

; Permitir credenciais (cookies, auth headers)
AllowCredentials=true

; Tempo de cache do preflight (12 horas)
MaxAge=43200
```

### Testar CORS

```bash
# Simular requisição de origem específica
curl -H "Origin: https://app.example.com" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/token

# Deve retornar 204 com headers CORS
```

---

## Request Limits

### Configuração

```ini
[application]
; Timeout de requisição em segundos
request_timeout = 30

[security]
; Tamanho máximo do body em bytes (1MB = 1048576)
max_body_size = 1048576
```

### Limite de tamanho do body

**Por que é importante:**
- Previne DoS via payloads gigantes
- Protege memória do servidor
- Padrão: 1MB (suficiente para APIs REST)

**Resposta quando excedido:**

**Status:** `413 Request Entity Too Large`

```json
{
  "code": "413",
  "message": "Request body too large. Max size: 1048576 bytes"
}
```

### Timeout de requisição

**Por que é importante:**
- Previne conexões travadas
- Libera recursos rapidamente
- Padrão: 30 segundos

**Resposta quando excedido:**

**Status:** `408 Request Timeout`

```json
{
  "code": "408",
  "message": "Request timeout"
}
```

### Valores recomendados

| Tipo de API | Max Body Size | Timeout |
|-------------|---------------|---------|
| **REST puro** | 1MB | 30s |
| **Com uploads** | 10MB | 60s |
| **Webhooks** | 5MB | 45s |
| **Interno/Admin** | 50MB | 120s |

---

## Checklist de Produção

Use este checklist antes de implantar em produção:

### Configurações obrigatórias

- [ ] **HTTPS/TLS habilitado** (`tls.enabled = true`)
- [ ] **Certificado válido** (não auto-assinado)
- [ ] **CORS restrito** (`AllowedOrigins` com domínios específicos)
- [ ] **Rate limiting habilitado** (`rate_limit_enabled = true`)
- [ ] **Environment = production** (`environment = production`)

### Configurações recomendadas

- [ ] **Request timeout configurado** (30-60s)
- [ ] **Max body size ajustado** (1-5MB conforme necessidade)
- [ ] **Logs de banco habilitados** (`ws_grava_log_db = true`)
- [ ] **Permissões do dbinit.ini** (`chmod 600 dbinit.ini`)

### Validações de segurança

- [ ] **Testar SQL injection** (tentativas de bypass devem falhar)
- [ ] **Testar rate limiting** (ultrapassar limite retorna 429)
- [ ] **Testar CORS** (origem não permitida é bloqueada)
- [ ] **Testar HTTPS** (conexão criptografada, HSTS ativo)
- [ ] **Testar timeout** (requisição longa é interrompida)

### Monitoramento

- [ ] **Logs de CORS bloqueados** (verificar origens suspeitas)
- [ ] **Logs de rate limit** (identificar possíveis ataques)
- [ ] **Logs de SQL errors** (não devem ocorrer SQL injection)
- [ ] **Certificado TLS** (configurar renovação automática)

---

## Troubleshooting

### Erro: "Certificado TLS não encontrado"

```
ERROR Certificado TLS não encontrado cert_file=certs/server.crt
```

**Solução:**
```bash
# Verificar se arquivo existe
ls -l certs/server.crt

# Se não existir, gerar certificado auto-assinado
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \
  -out certs/server.crt -days 365 -nodes
```

### Erro: "Rate limit exceeded" mesmo com poucos requests

**Causa:** Múltiplas instâncias compartilham o mesmo IP

**Solução:** Aumentar limites ou usar IP real em vez de proxy:

```go
// Em main.go, antes de iniciar o router
router.ForwardedByClientIP = true
```

### Aviso: "CORS configurado para permitir TODAS as origens"

**Causa:** `AllowedOrigins` está vazio em produção

**Solução:** Adicionar origens permitidas em `dbinit.ini`:

```ini
AllowedOrigins=https://seu-dominio.com
```

### Erro: "Request timeout" em operações longas

**Causa:** Timeout padrão é 30s

**Solução:** Aumentar timeout no `dbinit.ini`:

```ini
[application]
request_timeout = 120
```

### HTTPS funciona mas navegador avisa "Não seguro"

**Causa:** Certificado auto-assinado

**Solução:** Usar Let's Encrypt ou certificado comercial

---

## Recursos Adicionais

### Ferramentas de teste de segurança

```bash
# Verificar vulnerabilidades conhecidas
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Análise estática de segurança
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

### Referências

- **OWASP API Security Top 10**: https://owasp.org/www-project-api-security/
- **Go Security Best Practices**: https://github.com/OWASP/Go-SCP
- **JWT Best Practices**: https://tools.ietf.org/html/rfc8725
- **Let's Encrypt**: https://letsencrypt.org/
- **TLS Best Practices**: https://wiki.mozilla.org/Security/Server_Side_TLS

---

## Resumo de Comandos Rápidos

```bash
# Gerar certificado auto-assinado
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \
  -out certs/server.crt -days 365 -nodes

# Proteger dbinit.ini
chmod 600 dbinit.ini

# Verificar segurança
gosec ./...
govulncheck ./...

# Testar HTTPS
curl -k https://localhost:8443/health

# Testar rate limiting
for i in {1..70}; do curl http://localhost:8080/health; done

# Ver logs de segurança
tail -f log/wsicrmrest_$(date +%Y-%m-%d).log | grep -E "CORS|rate|TLS"
```

---

**Documentação mantida por:** Claude Code
**Versão do sistema:** 1.26.4.27
**Data:** 2025-11-17
