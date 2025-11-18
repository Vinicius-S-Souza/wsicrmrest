# Sistema de Logs - WSICRMREST

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-17

## 📋 Visão Geral

O WSICRMREST utiliza um sistema de logging robusto com rotação automática por data, baseado no Zap (go.uber.org/zap).

## 🔄 Rotação Automática de Logs

### Como Funciona

O sistema implementa **rotação automática de arquivos de log por data**:

- ✅ **Verificação contínua:** A cada escrita no log, verifica se a data mudou
- ✅ **Mudança à meia-noite:** Quando a data muda, fecha o arquivo atual e cria um novo
- ✅ **Thread-safe:** Usa mutex para garantir operações seguras em ambientes concorrentes
- ✅ **Sem reinicialização:** Funciona automaticamente sem precisar reiniciar o serviço

### Formato do Nome do Arquivo

```
wsicrmrest_YYYY-MM-DD.log
```

**Exemplos:**
- `wsicrmrest_2025-11-17.log` - Logs do dia 17 de novembro de 2025
- `wsicrmrest_2025-11-18.log` - Logs do dia 18 de novembro de 2025
- `wsicrmrest_2025-12-01.log` - Logs do dia 1º de dezembro de 2025

### Localização

Todos os logs são gravados na pasta:
```
log/
├── wsicrmrest_2025-11-15.log
├── wsicrmrest_2025-11-16.log
├── wsicrmrest_2025-11-17.log  # ← Arquivo ativo hoje
└── ...
```

## 📝 Formato dos Logs

### Estrutura JSON

Cada linha de log é um objeto JSON com os seguintes campos:

```json
{
  "level": "INFO",
  "timestamp": "2025-11-17T18:12:28-0300",
  "caller": "gin@v1.11.0/context.go:192",
  "message": "Iniciando WSICRMREST",
  "version": "Versão 3.0.0.1 (GO)",
  "version_date": "2025-11-17T17:51:37",
  "build_time": "2025-11-17T17:51:37"
}
```

### Campos Padrão

| Campo | Descrição | Exemplo |
|-------|-----------|---------|
| `level` | Nível de log (INFO, WARN, ERROR) | `"INFO"` |
| `timestamp` | Data/hora no formato ISO 8601 | `"2025-11-17T18:12:28-0300"` |
| `caller` | Arquivo e linha que gerou o log | `"handlers/token.go:78"` |
| `message` | Mensagem descritiva | `"Token gerado com sucesso"` |
| `...` | Campos adicionais contextuais | Variável |

### Níveis de Log

- **INFO:** Informações gerais de operação
- **WARN:** Avisos que não impedem funcionamento
- **ERROR:** Erros que requerem atenção

## 🔧 Implementação Técnica

### DailyRotatingWriter

A rotação é implementada através da struct `DailyRotatingWriter`:

```go
type DailyRotatingWriter struct {
    logDir      string      // Diretório dos logs
    currentDate string      // Data atual (YYYY-MM-DD)
    file        *os.File    // Arquivo aberto
    mu          sync.Mutex  // Mutex para thread-safety
}
```

### Funcionamento Interno

1. **Na inicialização:**
   - Cria diretório `log/` se não existir
   - Abre arquivo com a data atual
   - Armazena a data no formato `YYYY-MM-DD`

2. **A cada escrita (`Write`):**
   - Obtém a data atual
   - Compara com a data armazenada
   - Se diferente, executa rotação:
     - Fecha arquivo atual
     - Cria novo arquivo com nova data
     - Atualiza data armazenada

3. **Thread-safety:**
   - Usa mutex para proteger operações
   - Garante que apenas uma goroutine acesse por vez

### Exemplo de Uso

O logger é criado automaticamente na inicialização:

```go
// cmd/server/main.go
log, err := logger.NewLogger()
if err != nil {
    fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
    os.Exit(1)
}
defer log.Sync()

log.Info("Sistema iniciado")
```

## 📊 Exemplos de Logs

### Inicialização do Servidor

```json
{"level":"INFO","timestamp":"2025-11-17T18:12:28-0300","caller":"server/main.go:78","message":"Iniciando WSICRMREST","version":"Versão 3.0.0.1 (GO)","version_date":"2025-11-17T17:51:37","build_time":"2025-11-17T17:51:37"}
{"level":"INFO","timestamp":"2025-11-17T18:12:28-0300","caller":"server/main.go:90","message":"Conexão com banco de dados estabelecida com sucesso"}
{"level":"INFO","timestamp":"2025-11-17T18:12:28-0300","caller":"server/main.go:99","message":"Dados do organizador carregados com sucesso","codigo":1,"nome":"Minha Empresa"}
```

### Webhook Recebido

```json
{"level":"INFO","timestamp":"2025-11-17T18:15:32-0300","caller":"handlers/webhook_zenvia_sms.go:100","message":"Payload recebido do webhook Zenvia SMS","payload":"{\"type\":\"message_status\",...}"}
{"level":"INFO","timestamp":"2025-11-17T18:15:32-0300","caller":"handlers/webhook_zenvia_sms.go:193","message":"Processando ocorrência","from":"5573988769791","evento":"sent","messageId":"MSG-123","externalId":"12345","description":"Message sent"}
```

### Erro de Autenticação

```json
{"level":"WARN","timestamp":"2025-11-17T18:20:15-0300","caller":"handlers/token.go:85","message":"Credenciais inválidas","authorization":"Basic xxxxx","grant_type":"client_credentials"}
```

## 🗂️ Gerenciamento de Logs

### Ver Logs em Tempo Real (Linux)

```bash
# Seguir logs em tempo real
tail -f log/wsicrmrest_$(date +%Y-%m-%d).log

# Últimas 100 linhas
tail -n 100 log/wsicrmrest_$(date +%Y-%m-%d).log

# Filtrar por nível
grep '"level":"ERROR"' log/wsicrmrest_*.log

# Filtrar por mensagem
grep "webhook" log/wsicrmrest_*.log
```

### Ver Logs em Tempo Real (Windows)

```batch
REM PowerShell - Seguir logs em tempo real
Get-Content log\wsicrmrest_2025-11-17.log -Wait -Tail 50

REM CMD - Últimas 50 linhas
powershell -Command "Get-Content log\wsicrmrest_2025-11-17.log -Tail 50"

REM Filtrar por nível
findstr "ERROR" log\wsicrmrest_*.log
```

### Limpeza de Logs Antigos

#### Script Linux

```bash
#!/bin/bash
# Remover logs com mais de 30 dias
find log/ -name "wsicrmrest_*.log" -mtime +30 -delete

# Comprimir logs com mais de 7 dias
find log/ -name "wsicrmrest_*.log" -mtime +7 -exec gzip {} \;
```

#### Script Windows (PowerShell)

```powershell
# Remover logs com mais de 30 dias
Get-ChildItem log\wsicrmrest_*.log | Where-Object {
    $_.LastWriteTime -lt (Get-Date).AddDays(-30)
} | Remove-Item

# Comprimir logs com mais de 7 dias
Get-ChildItem log\wsicrmrest_*.log | Where-Object {
    $_.LastWriteTime -lt (Get-Date).AddDays(-7)
} | ForEach-Object {
    Compress-Archive -Path $_.FullName -DestinationPath "$($_.FullName).zip"
    Remove-Item $_.FullName
}
```

### Tamanho Aproximado

- **Por dia:** ~1-10 MB dependendo do volume de requisições
- **Por mês:** ~30-300 MB
- **Recomendação:** Manter últimos 30-90 dias

## 📈 Monitoramento

### Verificar Saúde dos Logs

```bash
# Verificar se logs estão sendo gravados
ls -lht log/ | head -5

# Ver tamanho dos logs
du -sh log/

# Contar linhas por dia
wc -l log/wsicrmrest_*.log

# Encontrar erros recentes
grep '"level":"ERROR"' log/wsicrmrest_$(date +%Y-%m-%d).log | tail -10
```

### Alertas Recomendados

1. **Log não rotacionou:** Verificar se arquivo de ontem ainda está sendo escrito
2. **Tamanho excessivo:** Arquivo maior que 100 MB pode indicar problema
3. **Muitos erros:** Mais de 100 erros por hora requer investigação
4. **Disco cheio:** Logs podem encher disco se não limpos

## 🔍 Análise de Logs

### Ferramentas Úteis

#### jq (Linux/Mac)

```bash
# Logs de erro formatados
cat log/wsicrmrest_2025-11-17.log | grep ERROR | jq '.'

# Contar por nível
cat log/wsicrmrest_2025-11-17.log | jq -r '.level' | sort | uniq -c

# Extrair apenas mensagens
cat log/wsicrmrest_2025-11-17.log | jq -r '.message'
```

#### PowerShell (Windows)

```powershell
# Converter JSON e filtrar
Get-Content log\wsicrmrest_2025-11-17.log | ConvertFrom-Json | Where-Object level -eq "ERROR"

# Contar por nível
Get-Content log\wsicrmrest_2025-11-17.log | ConvertFrom-Json | Group-Object level | Select-Object Name, Count
```

## 🚨 Troubleshooting

### Problema: Logs não estão sendo gravados

**Verificações:**
```bash
# 1. Verificar permissões da pasta
ls -la log/

# 2. Verificar espaço em disco
df -h

# 3. Verificar se o serviço está rodando
ps aux | grep wsicrmrest  # Linux
sc query WSICRMREST       # Windows
```

### Problema: Arquivo não rotaciona à meia-noite

**Causa:** A rotação acontece na **primeira escrita após a meia-noite**, não exatamente à meia-noite.

**Solução:** Isso é normal. Se o sistema não tiver nenhuma requisição após meia-noite, continuará usando o arquivo do dia anterior até a primeira requisição do novo dia.

### Problema: Arquivo muito grande

**Causa:** Volume alto de requisições ou logging excessivo.

**Solução:**
```bash
# Comprimir arquivo atual
gzip log/wsicrmrest_2025-11-17.log

# Limpar logs antigos
find log/ -name "*.log" -mtime +7 -delete
```

## 📋 Boas Práticas

1. **Rotação regular:** Limpe logs com mais de 30-90 dias
2. **Monitoramento:** Configure alertas para erros frequentes
3. **Backup:** Faça backup de logs importantes antes de deletar
4. **Compressão:** Comprima logs antigos para economizar espaço
5. **Análise:** Revise logs periodicamente para identificar padrões

## 🔐 Segurança

### Informações Sensíveis

O sistema **remove automaticamente** informações sensíveis dos logs:

- ✅ Headers `Authorization` são removidos antes de logar
- ✅ Senhas não são logadas
- ✅ Tokens JWT não aparecem completos nos logs

### Controle de Acesso

Recomendações:
```bash
# Linux - Permissões recomendadas
chmod 755 log/           # Pasta
chmod 644 log/*.log      # Arquivos

# Dono: Usuário do serviço
chown wsicrmrest:wsicrmrest log/ -R
```

## 📚 Referências

- **Implementação:** `internal/logger/logger.go`
- **Configuração:** Zap (go.uber.org/zap)
- **Formato:** JSON Lines (http://jsonlines.org/)

---

**Documentação mantida por:** Equipe de Desenvolvimento
**Última revisão:** 2025-11-17
