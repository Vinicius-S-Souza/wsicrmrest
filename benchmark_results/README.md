# Resultados de Benchmarks

**Data de criação:** 26/11/2025 15:00
**Versão:** 3.0.0.6

Este diretório armazena automaticamente os resultados dos benchmarks executados.

## Estrutura dos Arquivos

```
benchmark_results/
├── benchmark_YYYYMMDD_HHMMSS.txt     - Resultados completos
├── benchmark_mem_YYYYMMDD_HHMMSS.txt - Relatório de memória
└── README.md                         - Este arquivo
```

## Formato dos Arquivos

### Arquivos de Resultados

Cada execução do script `run_benchmarks.sh` cria automaticamente:

- **benchmark_YYYYMMDD_HHMMSS.txt**: Resultados completos com todas as métricas
- **benchmark_mem_YYYYMMDD_HHMMSS.txt**: Foco em uso de memória

### Exemplo de Conteúdo

```
==========================================
WSICRMREST - Benchmark Results
Data: 2025-11-26 15:00:00
Go Version: go version go1.22.1 linux/amd64
==========================================

BenchmarkZenviaEmailWebhook_JSONParsing-8    967306    2357 ns/op    792 B/op    11 allocs/op
BenchmarkZenviaSMSWebhook_JSONParsing-8      946827    2357 ns/op    792 B/op    11 allocs/op
...
```

## Métricas Explicadas

- **ns/op**: Tempo médio por operação em nanosegundos
- **B/op**: Bytes alocados por operação
- **allocs/op**: Número de alocações de memória por operação

## Analisando Resultados

### Performance Atual (Baseline)

**Parsing JSON** (~2357 ns/op = ~0.002ms):
- ✅ Excelente performance
- ✅ Baixo uso de memória (792 B)
- ✅ Poucas alocações (11 allocs)

**Requisição HTTP Completa** (~5131 ns/op = ~0.005ms):
- ✅ Muito bom
- ✅ Memória moderada (6962 B)
- ✅ Alocações controladas (31 allocs)

**JSON Inválido** (~405 ns/op = ~0.0004ms):
- ✅ Falha rápida (desejável)
- ✅ Baixíssimo overhead

**Payload Grande** (~8646 ns/op = ~0.008ms):
- ✅ Performance ainda boa mesmo com 100x de dados
- ✅ Memória: apenas 2464 B (eficiente)

**Execução Paralela** (~447 ns/op):
- ✅ Excelente escalabilidade
- ✅ Performance similar ao sequencial

### Metas de Performance

| Métrica | Atual | Meta | Status |
|---------|-------|------|--------|
| JSON Parsing | 2.3 µs | < 5 µs | ✅ Excelente |
| HTTP Request | 5.1 µs | < 10 µs | ✅ Excelente |
| Memória | 792 B | < 2 KB | ✅ Excelente |
| Alocações | 11 | < 20 | ✅ Excelente |

## Comparando Resultados

### Usando benchstat

Instale a ferramenta:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Compare dois arquivos:

```bash
benchstat benchmark_20251126_100000.txt benchmark_20251126_150000.txt
```

### Exemplo de Comparação

```
name                            old time/op    new time/op    delta
ZenviaEmailWebhook_JSONParsing  2.35µs ± 2%    2.30µs ± 1%   -2.13%
ZenviaEmailWebhook_HTTPRequest  5.13µs ± 3%    5.01µs ± 2%   -2.34%

name                            old alloc/op   new alloc/op   delta
ZenviaEmailWebhook_JSONParsing  792B ± 0%      792B ± 0%      0.00%
```

## Histórico de Melhorias

### Versão 3.0.0.6 (26/11/2025)

- ✅ Implementação inicial dos benchmarks
- ✅ Baseline estabelecido:
  - JSON Parsing: 2.36 µs/op
  - HTTP Request: 5.13 µs/op
  - Memória: 792 B/op (parsing), 6962 B/op (HTTP)
  - Alocações: 11 (parsing), 31 (HTTP)

## Executando Novos Benchmarks

### Opção 1: Script Interativo

```bash
cd /home/vinicius/projetos/wsicrmrest
./scripts/run_benchmarks.sh
```

### Opção 2: Comando Direto

```bash
# Todos os benchmarks
go test -bench=. -benchmem ./internal/handlers/ > benchmark_results/resultado_$(date +%Y%m%d_%H%M%S).txt

# Apenas email
go test -bench=Email -benchmem ./internal/handlers/ > benchmark_results/email_$(date +%Y%m%d_%H%M%S).txt

# Apenas SMS
go test -bench=SMS -benchmem ./internal/handlers/ > benchmark_results/sms_$(date +%Y%m%d_%H%M%S).txt
```

## Manutenção

### Limpeza de Arquivos Antigos

```bash
# Manter apenas os últimos 10 arquivos
cd benchmark_results
ls -t benchmark_*.txt | tail -n +11 | xargs rm -f
```

### Backup

```bash
# Criar backup compactado
tar -czf benchmark_backup_$(date +%Y%m%d).tar.gz benchmark_*.txt
```

## Alertas de Regressão

⚠️ **Investigar se:**

1. **Tempo aumentar > 10%**: Possível regressão de performance
2. **Memória aumentar > 20%**: Possível vazamento ou ineficiência
3. **Alocações aumentarem significativamente**: Possível problema de GC

## Referências

- [Documentação Completa](/home/vinicius/projetos/wsicrmrest/docs/BENCHMARKS.md)
- [Go Benchmark Guide](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Benchstat Tool](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)

---

**Última atualização:** 26/11/2025 15:00
**Versão:** 3.0.0.6
