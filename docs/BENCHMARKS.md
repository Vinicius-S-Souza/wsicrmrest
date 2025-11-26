# Benchmarks - Webhooks Zenvia

**Data de criação:** 26/11/2025 14:48
**Versão:** 3.0.0.6

## Visão Geral

Este documento descreve a infraestrutura de benchmarks para os webhooks Zenvia (email e SMS) do WSICRMREST. Os benchmarks permitem medir a performance das APIs sem conectar ao banco de dados Oracle real.

## Importante: Uso de Mocks

⚠️ **TODOS os benchmarks usam MOCKS** - não conectam ao banco de dados real:

- **MockDatabase**: Simula operações de banco de dados (GetEmailByAPIMessageID, InsereLogsAPI, etc.)
- **MockLogger**: Simula logging do Zap
- **Nenhuma alteração** é feita no banco de dados Oracle durante os testes

## Estrutura dos Arquivos

```
internal/
  testhelpers/
    mocks.go      - Mock database e logger (180 linhas)
    fixtures.go   - Payloads JSON de teste (281 linhas)
    setup.go      - Funções helper para setup (229 linhas)
  handlers/
    webhook_bench_test.go - Benchmarks (530+ linhas)

scripts/
  run_benchmarks.sh - Script para executar benchmarks

docs/
  BENCHMARKS.md - Esta documentação

benchmark_results/  - Resultados salvos automaticamente
```

## Benchmarks Disponíveis

### 1. Benchmarks de Email

| Benchmark | Descrição | Status Zenvia | Código Interno |
|-----------|-----------|---------------|----------------|
| `BenchmarkZenviaEmailWebhook_Sent` | Email enviado | SENT | 121 |
| `BenchmarkZenviaEmailWebhook_Delivered` | Email entregue | DELIVERED | 122 |
| `BenchmarkZenviaEmailWebhook_Read` | Email lido/aberto | READ | 123 |
| `BenchmarkZenviaEmailWebhook_Rejected` | Email rejeitado | NOT_SENT | 124 |
| `BenchmarkZenviaEmailWebhook_Bounce` | Email devolvido | BOUNCED | 125 |
| `BenchmarkZenviaEmailWebhook_AllStatuses` | Todos os status alternados | Vários | 121-125 |
| `BenchmarkZenviaEmailWebhook_InvalidJSON` | JSON inválido (erro) | N/A | N/A |
| `BenchmarkZenviaEmailWebhook_Parallel` | Execução paralela | SENT | 121 |
| `BenchmarkZenviaEmailWebhook_LargePayload` | Payload grande (100x) | SENT | 121 |

### 2. Benchmarks de SMS

| Benchmark | Descrição | Status Zenvia | Código Interno |
|-----------|-----------|---------------|----------------|
| `BenchmarkZenviaSMSWebhook_Sent` | SMS enviado | SENT | 121 |
| `BenchmarkZenviaSMSWebhook_Delivered` | SMS entregue | DELIVERED | 122 |
| `BenchmarkZenviaSMSWebhook_Read` | SMS lido | READ | 123 |
| `BenchmarkZenviaSMSWebhook_Rejected` | SMS rejeitado | NOT_SENT | 124 |
| `BenchmarkZenviaSMSWebhook_Bounce` | SMS devolvido | BOUNCED | 125 |
| `BenchmarkZenviaSMSWebhook_AllStatuses` | Todos os status alternados | Vários | 121-125 |
| `BenchmarkZenviaSMSWebhook_InvalidJSON` | JSON inválido (erro) | N/A | N/A |
| `BenchmarkZenviaSMSWebhook_Parallel` | Execução paralela | SENT | 121 |
| `BenchmarkZenviaSMSWebhook_LargePayload` | Payload grande (100x) | SENT | 121 |

## Como Executar

### Opção 1: Script Interativo (Recomendado)

```bash
cd /home/vinicius/projetos/wsicrmrest
./scripts/run_benchmarks.sh
```

**Menu de opções:**

1. Email - Todos os status (Sent, Delivered, Read, Rejected, Bounce)
2. SMS - Todos os status (Sent, Delivered, Read, Rejected, Bounce)
3. Email e SMS - Status individuais
4. Benchmarks paralelos (Email e SMS)
5. Benchmarks com payloads grandes
6. Benchmarks com JSON inválido
7. Todos os benchmarks (completo)
8. Benchmark rápido (apenas Sent)

### Opção 2: Comandos Go Diretos

```bash
# Executar todos os benchmarks
go test -bench=. -benchmem ./internal/handlers/

# Executar apenas benchmarks de email
go test -bench=BenchmarkZenviaEmailWebhook -benchmem ./internal/handlers/

# Executar apenas benchmarks de SMS
go test -bench=BenchmarkZenviaSMSWebhook -benchmem ./internal/handlers/

# Executar benchmark específico
go test -bench=BenchmarkZenviaEmailWebhook_Sent -benchmem ./internal/handlers/

# Executar com tempo personalizado (padrão: 1s, recomendado: 3-5s)
go test -bench=. -benchmem -benchtime=5s ./internal/handlers/

# Salvar resultados em arquivo
go test -bench=. -benchmem ./internal/handlers/ > benchmark_results.txt
```

### Opção 3: Benchmarks Paralelos

```bash
# Executar com mais CPUs (padrão: GOMAXPROCS)
go test -bench=Parallel -benchmem -cpu=1,2,4,8 ./internal/handlers/
```

## Interpretando Resultados

### Formato de Saída

```
BenchmarkZenviaEmailWebhook_Sent-8    15234    78563 ns/op    12456 B/op    142 allocs/op
```

**Explicação:**

- `BenchmarkZenviaEmailWebhook_Sent-8`: Nome do benchmark e número de CPUs
- `15234`: Número de iterações executadas
- `78563 ns/op`: Tempo médio por operação (nanosegundos)
- `12456 B/op`: Bytes alocados por operação
- `142 allocs/op`: Número de alocações de memória por operação

### Métricas Importantes

#### 1. Tempo por Operação (ns/op)

- **< 100,000 ns (< 0.1ms)**: Excelente
- **100,000 - 500,000 ns (0.1-0.5ms)**: Bom
- **500,000 - 1,000,000 ns (0.5-1ms)**: Aceitável
- **> 1,000,000 ns (> 1ms)**: Requer otimização

#### 2. Memória Alocada (B/op)

- **< 10KB**: Excelente
- **10KB - 50KB**: Bom
- **50KB - 100KB**: Aceitável
- **> 100KB**: Requer otimização

#### 3. Número de Alocações (allocs/op)

- **< 100**: Excelente
- **100 - 200**: Bom
- **200 - 500**: Aceitável
- **> 500**: Requer otimização

### Exemplo de Análise

```
BenchmarkZenviaEmailWebhook_Sent-8         15234    78563 ns/op    12456 B/op    142 allocs/op
BenchmarkZenviaEmailWebhook_Delivered-8    14892    81234 ns/op    12512 B/op    145 allocs/op
```

**Análise:**

- Tempo de resposta: **~0.08ms** (excelente)
- Memória: **~12KB** por requisição (muito bom)
- Alocações: **142-145** (bom)
- Diferença entre SENT e DELIVERED: **~3.4%** (normal)

## Comparando Resultados

### 1. Comparação Temporal

```bash
# Executar benchmark atual
go test -bench=. -benchmem ./internal/handlers/ > benchmark_2025_11_26.txt

# Comparar com versão anterior (requer benchstat)
go install golang.org/x/perf/cmd/benchstat@latest
benchstat benchmark_2025_11_20.txt benchmark_2025_11_26.txt
```

**Exemplo de saída do benchstat:**

```
name                            old time/op    new time/op    delta
ZenviaEmailWebhook_Sent-8       82.5µs ± 2%    78.6µs ± 1%   -4.73%
ZenviaEmailWebhook_Delivered-8  84.2µs ± 3%    81.2µs ± 2%   -3.56%

name                            old alloc/op   new alloc/op   delta
ZenviaEmailWebhook_Sent-8       13.2kB ± 0%    12.5kB ± 0%   -5.30%
```

### 2. Análise de Regressão

Se os novos benchmarks forem **>10% mais lentos**, investigar:

1. Mudanças recentes no código
2. Novos middlewares ou validações
3. Aumento de alocações de memória
4. Operações de I/O adicionadas

## Estrutura dos Testes

### MockDatabase

```go
type MockDatabase struct {
    EmailData *models.EmailData
    SMSData   *models.SMSData
    Error     error

    // Contadores
    GetEmailCalls    int
    GetSMSCalls      int
    InsereLogsCalls  int
    InsereOcorrCalls int
}
```

**Métodos implementados:**

- `GetEmailByAPIMessageID(messageID string)` - Retorna EmailData simulado
- `GetSMSByMessageID(messageID string)` - Retorna SMSData simulado
- `InsereLogsAPI(...)` - Simula inserção de log (não grava no banco)
- `InsereLogsAPISMS(...)` - Simula inserção de log SMS
- `InsereOcorrenciaEmailInconsistente(...)` - Simula criação de ocorrência
- `InsereOcorrenciaSmsInconsistente(...)` - Simula criação de ocorrência SMS
- `GravaLogDB(...)` - Simula gravação de log de requisição

### MockLogger

```go
type MockLogger struct {
    InfoCalls  int
    WarnCalls  int
    ErrorCalls int
    DebugCalls int
}
```

**Todos os métodos do zap.SugaredLogger** estão implementados.

### Fixtures de Teste

**Payloads de Email:**

- `ZenviaEmailPayloadSent` - Status SENT (121)
- `ZenviaEmailPayloadDelivered` - Status DELIVERED (122)
- `ZenviaEmailPayloadRead` - Status READ (123)
- `ZenviaEmailPayloadRejected` - Status NOT_SENT (124)
- `ZenviaEmailPayloadBounce` - Status BOUNCED (125)

**Payloads de SMS:**

- `ZenviaSMSPayloadSent` - Status SENT (121)
- `ZenviaSMSPayloadDelivered` - Status DELIVERED (122)
- `ZenviaSMSPayloadRead` - Status READ (123)
- `ZenviaSMSPayloadRejected` - Status NOT_SENT (124)
- `ZenviaSMSPayloadBounce` - Status BOUNCED (125)

**Payloads Inválidos:**

- `InvalidJSON` - JSON malformado
- `EmptyJSON` - JSON vazio
- `MissingMessageID` - Falta campo obrigatório

## Boas Práticas

### 1. Execução de Benchmarks

- Execute benchmarks em ambiente **estável** (sem outras aplicações rodando)
- Execute **múltiplas vezes** e use a mediana
- Use `-benchtime=3s` ou `-benchtime=5s` para resultados mais estáveis
- Execute benchmarks **antes e depois** de mudanças no código

### 2. Análise de Performance

- Compare **sempre com baseline** (versão anterior)
- Foque em **tendências**, não em valores absolutos
- Monitore **regressões** (> 10% de degradação)
- Documente **melhorias** significativas (> 20%)

### 3. Otimização

Prioridades de otimização:

1. **Reduzir alocações** (allocs/op) - Maior impacto
2. **Reduzir memória** (B/op) - Impacto moderado
3. **Reduzir tempo** (ns/op) - Consequência natural das otimizações acima

## Troubleshooting

### Problema: "No test files"

```bash
# Solução: Verificar se está no diretório correto
cd /home/vinicius/projetos/wsicrmrest
go test -bench=. ./internal/handlers/
```

### Problema: "Package not found"

```bash
# Solução: Executar go mod tidy
go mod tidy
go test -bench=. ./internal/handlers/
```

### Problema: Resultados inconsistentes

```bash
# Solução: Aumentar benchtime e repetir
go test -bench=. -benchtime=10s -count=5 ./internal/handlers/
```

### Problema: "Too slow"

Se os benchmarks demorarem muito:

1. Execute benchmarks específicos: `-bench=BenchmarkZenviaEmailWebhook_Sent`
2. Reduza benchtime: `-benchtime=1s`
3. Use o menu do script: opção 8 (benchmark rápido)

## Arquivos de Resultados

Os resultados são salvos automaticamente em `benchmark_results/`:

```
benchmark_results/
  benchmark_20251126_144800.txt     - Resultados completos
  benchmark_mem_20251126_144800.txt - Relatório de memória
```

**Formato do nome:** `benchmark_YYYYMMDD_HHMMSS.txt`

## Exemplos de Uso

### 1. Validar Performance Após Mudança

```bash
# Antes da mudança
./scripts/run_benchmarks.sh
# Escolher opção 7 (todos)
# Resultados salvos em benchmark_results/benchmark_20251126_140000.txt

# Fazer mudanças no código...

# Depois da mudança
./scripts/run_benchmarks.sh
# Escolher opção 7 (todos)
# Resultados salvos em benchmark_results/benchmark_20251126_150000.txt

# Comparar
benchstat benchmark_results/benchmark_20251126_140000.txt \
          benchmark_results/benchmark_20251126_150000.txt
```

### 2. Testar Cenário Específico

```bash
# Testar apenas email com JSON inválido
go test -bench=BenchmarkZenviaEmailWebhook_InvalidJSON -benchmem ./internal/handlers/
```

### 3. Teste de Carga Paralelo

```bash
# Simular múltiplos CPUs
go test -bench=Parallel -benchmem -cpu=1,2,4,8 ./internal/handlers/
```

## Referências

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Go Benchmark Guide](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Benchstat Tool](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [WSICRMREST - CLAUDE.md](/home/vinicius/projetos/wsicrmrest/CLAUDE.md)

## Histórico de Versões

| Versão | Data | Descrição |
|--------|------|-----------|
| 3.0.0.6 | 26/11/2025 | Implementação inicial da infraestrutura de benchmarks |

---

**Última atualização:** 26/11/2025 14:48
**Versão do documento:** 3.0.0.6
