package utils

import (
	"strings"
	"time"
)

// EliminaCaracterNulo remove caracteres nulos de uma string
// Equivalente a fgEliminaCaracterNulo do WinDev
func EliminaCaracterNulo(s string) string {
	return strings.ReplaceAll(s, "\x00", "")
}

// StringChange substitui todas as ocorrências de old por new em s
// Equivalente a pgStringChange do WinDev
func StringChange(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// CalcTimeStampUnix calcula o timestamp Unix para uma data/hora específica
// Equivalente a pgCalcTimeStampUnix do WinDev
// timezone: fuso horário em horas (ex: -3 para Brasília)
func CalcTimeStampUnix(date time.Time, timezone int) int64 {
	// Ajustar para UTC baseado no fuso horário
	utcTime := date.Add(time.Duration(-timezone) * time.Hour)
	return utcTime.Unix()
}

// FormatDateTimeOracle formata data/hora para SQL Oracle
// Equivalente a FcDateTime do WinDev
// Retorna no formato: TO_TIMESTAMP('MM/DD/YYYY HH24:MI:SS.FF', 'MM/DD/YYYY HH24:MI:SS.FF')
func FormatDateTimeOracle(date time.Time, withMilliseconds bool) string {
	if withMilliseconds {
		dateStr := date.Format("01/02/2006 15:04:05.00")
		return "TO_TIMESTAMP('" + dateStr + "', 'MM/DD/YYYY HH24:MI:SS.FF')"
	}
	dateStr := date.Format("01/02/2006 15:04:05")
	return "TO_TIMESTAMP('" + dateStr + "', 'MM/DD/YYYY HH24:MI:SS')"
}

// Escopo retorna a string de escopo baseada no código bitwise
// Equivalente a pgScopo do WinDev
func Escopo(codigo int64) string {
	if codigo <= 0 {
		return ""
	}

	scopes := []string{}

	// Bit 1 (1) - clientes
	if codigo&1 != 0 {
		scopes = append(scopes, "clientes")
	}

	// Bit 2 (2) - lojas
	if codigo&2 != 0 {
		scopes = append(scopes, "lojas")
	}

	// Bit 3 (4) - ofertas
	if codigo&4 != 0 {
		scopes = append(scopes, "ofertas")
	}

	// Bit 4 (8) - produtos
	if codigo&8 != 0 {
		scopes = append(scopes, "produtos")
	}

	// Bit 5 (16) - pontos
	if codigo&16 != 0 {
		scopes = append(scopes, "pontos")
	}

	// Bit 6 (32) - private
	if codigo&32 != 0 {
		scopes = append(scopes, "private")
	}

	// Bit 7 (64) - convenio
	if codigo&64 != 0 {
		scopes = append(scopes, "convenio")
	}

	// Bit 8 (128) - giftcard
	if codigo&128 != 0 {
		scopes = append(scopes, "giftcard")
	}

	// Bit 9 (256) - cobranca
	if codigo&256 != 0 {
		scopes = append(scopes, "cobranca")
	}

	// Bit 10 (512) - basico
	if codigo&512 != 0 {
		scopes = append(scopes, "basico")
	}

	// Bit 11 (1024) - sistema
	if codigo&1024 != 0 {
		scopes = append(scopes, "sistema")
	}

	// Bit 12 (2048) - terceiros
	if codigo&2048 != 0 {
		scopes = append(scopes, "terceiros")
	}

	// Bit 13 (4096) - totem
	if codigo&4096 != 0 {
		scopes = append(scopes, "totem")
	}

	return strings.Join(scopes, " ")
}

// RemoveBase64Padding remove padding de strings base64
func RemoveBase64Padding(s string) string {
	return strings.TrimRight(s, "=")
}

// Base64URLSafe converte base64 padrão para URL-safe
func Base64URLSafe(s string) string {
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	return RemoveBase64Padding(s)
}

// XML2CLOB prepara string para inserção em CLOB Oracle
// Equivalente a fgXML2CLOB do WinDev
func XML2CLOB(s string) string {
	if s == "" {
		return "NULL"
	}
	// Escapar aspas simples para Oracle
	s = strings.ReplaceAll(s, "'", "''")
	return "'" + s + "'"
}

// SanitizeForSQL sanitiza string para uso em SQL
func SanitizeForSQL(s string) string {
	// Substituir aspas duplas por aspas duplas escapadas
	s = strings.ReplaceAll(s, "\"", "\"\"")
	// Substituir aspas simples por backtick
	s = strings.ReplaceAll(s, "'", "`")
	return s
}
