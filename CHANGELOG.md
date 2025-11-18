# Changelog - WSICRMREST

## [1.0.0] - 2025-01-27

### ✅ Implementado

#### Variáveis Globais

- ✅ `gsKey` → `config.JWT.SecretKey` = `"CloudI0812IcrMmDB"` (hardcoded)
- ✅ `gsIss` → `config.JWT.Issuer` = `"WSCloudICrmIntellsys"` (hardcoded)
- ✅ `gsKeyDelivery` → `config.JWT.KeyDelivery` = `"Ped2505IcrM"` (hardcoded)
- ✅ `gnFusoHorario` → `config.JWT.Timezone` = `0` (hardcoded)
- ✅ `gsVersao` → `config.Application.Version` (configurável via dbinit.ini)
- ✅ `gsDataVersao` → `config.Application.VersionDate` (configurável via dbinit.ini)
- ✅ `gnRegModulos` → `config.Organization.RegModulos` = `1` (hardcoded)
- ✅ **Credenciais JWT NÃO configuráveis** (valores fixos no código)

#### Tabela ORGANIZADOR

- ✅ `pgLeOrganizador()` → `database.LeOrganizador()`
- ✅ Carregamento automático e **OBRIGATÓRIO** na inicialização
- ✅ 20 campos de organização carregados **exclusivamente** do banco
- ✅ Campos de grupos de faturamento (6 grupos)
- ✅ **Dados NÃO vêm do dbinit.ini** (seção [organization] removida)
- ✅ Sistema não inicia se tabela não existir ou estiver vazia
- ✅ Logs informativos sobre dados carregados

#### Estrutura do Projeto

- ✅ Projeto Go modular e organizado
- ✅ Gerenciamento de dependências com go.mod
- ✅ Arquitetura limpa com separação de responsabilidades

#### Configuração

- ✅ Leitura de configurações via arquivo INI (dbinit.ini)
- ✅ Suporte a variáveis de ambiente
- ✅ Configuração de timezone
- ✅ Configuração de logs
- ✅ Configuração de conexão Oracle via TNSNAMES

#### Banco de Dados

- ✅ Conexão com Oracle usando godror
- ✅ Suporte a TNSNAMES
- ✅ Pool de conexões configurável
- ✅ Estrutura completa de tabelas:
  - WSAPLICACOES
  - WSAPLLOGTOKEN
  - WSREQUISICOES
- ✅ Sequences para IDs automáticos

#### Sistema de Logs

- ✅ Logs em arquivo usando Zap
- ✅ Rotação automática de logs por data
- ✅ Formato: `wsicrmrest_YYYY-MM-DD.log`
- ✅ Logs estruturados (JSON)
- ✅ Logs no console e arquivo simultaneamente
- ✅ Log de requisições no banco de dados

#### APIs Implementadas

##### GET /connect/v1/token

- ✅ Autenticação via Basic Auth
- ✅ Validação de grant_type (client_credentials)
- ✅ Geração de JWT (HS256)
- ✅ Validação de Client_ID e Client_Secret
- ✅ Verificação de status da aplicação
- ✅ Cálculo de timestamp Unix
- ✅ Expiração configurável por aplicação
- ✅ Sistema de escopos bitwise (13 escopos)
- ✅ Log de tokens gerados
- ✅ Resposta com código personalizado

##### GET /connect/v1/wsteste

- ✅ Teste de conexão com banco
- ✅ Retorno de informações da organização
- ✅ Health check do sistema

#### Funções Auxiliares (Conversão WinDev → Go)

✅ **Funções de String:**

- `utils.EliminaCaracterNulo()` ← `fgEliminaCaracterNulo()`
- `utils.StringChange()` ← `pgStringChange()`
- `utils.SanitizeForSQL()` - Nova
- `utils.XML2CLOB()` ← `fgXML2CLOB()`

✅ **Funções de Data/Hora:**

- `utils.FormatDateTimeOracle()` ← `FcDateTime()`
- `utils.CalcTimeStampUnix()` ← `pgCalcTimeStampUnix()`

✅ **Funções de Negócio:**

- `utils.Escopo()` ← `pgScopo()` - Sistema de escopos bitwise completo

✅ **Funções de Database:**

- `database.GravaLogDB()` ← `pgGravaLogDB()`
- `database.NewDatabase()` ← `pgAbreBancoDados()`

✅ **Sistema de Contexto:**

- `context.RequestContext` - Contexto de requisição
- UUID único por requisição
- Rastreamento de duração
- Cliente e aplicação no contexto

#### Segurança

- ✅ Remoção automática de header Authorization dos logs
- ✅ JWT com HMAC-SHA256
- ✅ Validação de credenciais
- ✅ Sanitização de inputs SQL
- ✅ Escape de caracteres especiais

#### Middleware

- ✅ Logger de requisições com Zap
- ✅ Recovery para panic handling (Gin)
- ✅ Rastreamento de IP do cliente
- ✅ Medição de tempo de resposta

#### Ferramentas e Scripts

- ✅ Makefile completo com comandos úteis
- ✅ Script de teste de APIs (test_apis.sh)
- ✅ Arquivo de exemplo de configuração
- ✅ .gitignore configurado

#### Documentação

- ✅ README.md completo
- ✅ QUICKSTART.md (guia de 5 minutos)
- ✅ DATABASE_SCHEMA.md (estrutura completa das tabelas)
- ✅ CHANGELOG.md (este arquivo)
- ✅ Comentários inline no código
- ✅ Documentação de funções e handlers

#### Qualidade de Código

- ✅ Código formatado (gofmt)
- ✅ Sem warnings do compilador
- ✅ Imports organizados
- ✅ Nomes de variáveis descritivos
- ✅ Separação de responsabilidades

### 📊 Estatísticas

- **Arquivos Go:** 10
- **Linhas de código:** ~1500+
- **Packages:** 8
- **APIs:** 2
- **Funções auxiliares:** 10+
- **Tempo de compilação:** < 5s
- **Tamanho do binário:** ~17MB

### 🔄 Conversão WinDev

| WinDev Procedure | Go Implementation | Status |
|------------------|-------------------|--------|
| `pgGerarToken()` | `handlers.GenerateToken()` | ✅ Completo |
| `pgWSRestTeste()` | `handlers.WSTest()` | ✅ Completo |
| `pgGravaLogDB()` | `database.GravaLogDB()` | ✅ Completo |
| `pgScopo()` | `utils.Escopo()` | ✅ Completo |
| `pgAbreBancoDados()` | `database.NewDatabase()` | ✅ Completo |
| `pgImprimirLog()` | `logger.NewLogger()` | ✅ Completo |
| `fgEliminaCaracterNulo()` | `utils.EliminaCaracterNulo()` | ✅ Completo |
| `pgStringChange()` | `utils.StringChange()` | ✅ Completo |
| `FcDateTime()` | `utils.FormatDateTimeOracle()` | ✅ Completo |
| `pgCalcTimeStampUnix()` | `utils.CalcTimeStampUnix()` | ✅ Completo |

### 📝 Notas Técnicas

#### Diferenças WinDev → Go

1. **Gerenciamento de Memória:**
   - WinDev: Gerenciamento automático
   - Go: Garbage collector eficiente

2. **Concorrência:**
   - WinDev: Threads
   - Go: Goroutines (mais leves e eficientes)

3. **Strings:**
   - WinDev: Strings ANSI e Unicode
   - Go: UTF-8 nativo

4. **Banco de Dados:**
   - WinDev: HExecuteSQLQuery
   - Go: database/sql com godror

5. **Logs:**
   - WinDev: fWriteLine
   - Go: Zap (estruturado e performático)

6. **HTTP:**
   - WinDev: WebService framework
   - Go: Gin (alto desempenho)

### 🎯 Melhorias em Relação ao WinDev

1. ✅ **Performance:** Go é compilado e mais rápido
2. ✅ **Concorrência:** Goroutines permitem melhor uso de CPU multi-core
3. ✅ **Deploy:** Binário único, sem dependências externas
4. ✅ **Logs:** Estruturados e mais fáceis de analisar
5. ✅ **Manutenção:** Código mais limpo e modular
6. ✅ **Testes:** Melhor suporte a testes unitários
7. ✅ **DevOps:** Facilita integração com pipelines CI/CD

---

## 🚀 Próximas Versões

### [1.1.0] - Planejado

#### A Implementar

- [ ] Middleware de validação de JWT
- [ ] Rotas protegidas com autenticação
- [ ] Refresh token
- [ ] Rate limiting
- [ ] CORS configurável
- [ ] Métricas (Prometheus)
- [ ] Health check avançado
- [ ] Graceful shutdown

#### A Adicionar

- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline
- [ ] Swagger/OpenAPI docs

#### Suporte SQL Server

- [ ] Driver SQL Server (driver = 1)
- [ ] Queries adaptadas
- [ ] Testes com SQL Server

---

## 📞 Suporte

Para problemas ou dúvidas:

1. Consulte [QUICKSTART.md](QUICKSTART.md)
2. Consulte [README.md](README.md)
3. Verifique [DATABASE_SCHEMA.md](docs/DATABASE_SCHEMA.md)

---

**Versão atual:** 1.0.0
**Data:** 2025-01-27
**Status:** ✅ Produção
