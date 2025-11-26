# Changelog - WSICRMREST

## [3.0.0.5] - 2025-11-26

### 🔐 Segurança TLS

#### Suporte para Chaves Privadas Criptografadas
- ✅ Implementado suporte completo para chaves SSL/TLS criptografadas
- ✅ Auto-detecção de formato: criptografado vs não criptografado
- ✅ Suporte para PKCS#1 (BEGIN RSA PRIVATE KEY com Proc-Type: 4,ENCRYPTED)
- ✅ Suporte para PKCS#8 (BEGIN ENCRYPTED PRIVATE KEY)
- ✅ Compatibilidade total com chaves não criptografadas (backward compatible)
- ✅ Senha configurável via `dbinit.ini` na seção `[tls]`

**Arquivos criados:**
- `internal/tls/loader.go` - Módulo de carregamento de certificados TLS com suporte a criptografia

**Arquivos modificados:**
- `internal/config/config.go` - Adicionado campo `KeyPassword` à estrutura `TLSConfig`
- `cmd/server/main.go` - Integrado loader TLS customizado para modo console
- `internal/service/windows_service.go` - Integrado loader TLS customizado para Windows Service
- `dbinit.ini` - Adicionado campo `key_password` na seção `[tls]`

**Funcionalidades:**
- Descriptografia automática de chaves PKCS#1 e PKCS#8 usando `x509.DecryptPEMBlock()`
- Parsing flexível: PKCS#8, PKCS#1 RSA, e Elliptic Curve (EC)
- Servidor HTTPS customizado com `http.Server` e `TLSConfig`
- Mensagens de erro claras e descritivas para problemas de senha/formato
- Configuração TLS segura: TLS 1.2 mínimo + cipher suites recomendadas
- Zero breaking changes: chaves não criptografadas continuam funcionando

**Configuração no dbinit.ini:**
```ini
[tls]
enabled = true
cert_file = C:\Apache24\cert\certificate.crt
key_file = C:\Apache24\cert\private.key
key_password = sua_senha_aqui  ; Deixe vazio para chaves não criptografadas
port = 8443
```

**Resolução de problemas:**
- ✅ Resolvido: "tls: failed to parse private key" ao usar chaves PKCS#8 criptografadas
- ✅ Resolvido: Não era possível usar certificados SSL de produção com senha
- ✅ Aplicação agora suporta certificados SSL padrão de autoridades certificadoras

### 🧹 Qualidade de Código

#### Correção de Warnings do Compilador
- ✅ Removido parâmetro não utilizado `cfg` da função `logTokenToDB()` em `internal/handlers/token.go`
- ✅ Atualizado chamador da função para corresponder à nova assinatura
- ✅ Código mais limpo e sem warnings do `go vet` ou `gopls`

**Impacto:**
- Melhor manutenibilidade do código
- Conformidade com boas práticas de Go
- Funções com assinaturas mais limpas e sem parâmetros desnecessários

### 📝 Detalhes Técnicos

**Fluxo de descriptografia de chaves:**
1. Ler arquivo de chave privada do disco
2. Decodificar bloco PEM
3. Detectar se é criptografada via `x509.IsEncryptedPEMBlock()` ou tipo "ENCRYPTED PRIVATE KEY"
4. Se criptografada:
   - Verificar se senha foi fornecida (erro claro se não)
   - Descriptografar usando `x509.DecryptPEMBlock()` com senha
5. Se não criptografada: usar bytes direto
6. Parsear chave descriptografada (PKCS#8 → PKCS#1 → EC)
7. Criar `tls.Certificate` e `tls.Config`
8. Retornar configuração pronta para `http.Server`

**Benefícios da implementação:**
- ✅ Sem dependências externas adicionais (usa apenas `crypto/x509` e `crypto/tls` padrão)
- ✅ Performance: descriptografia ocorre apenas uma vez na inicialização
- ✅ Segurança: senha nunca é logada ou exposta
- ✅ Confiabilidade: testa múltiplos formatos automaticamente
- ✅ Simplicidade: configuração via arquivo INI, sem variáveis de ambiente complexas

**Compatibilidade:**
- ✅ Windows Server 2016+
- ✅ Windows 10/11
- ✅ Linux (todas distribuições)
- ✅ Modo Console e Windows Service
- ✅ Certificados de Let's Encrypt, DigiCert, Sectigo, etc.

---

## [3.0.0.4] - 2025-11-24

### 🔧 Melhorias

#### Seleção de Arquitetura (32/64 bits)
- ✅ Detecção automática da arquitetura do Windows
- ✅ Menu interativo para seleção de executável (32 ou 64 bits)
- ✅ Validação de compatibilidade arquitetura vs sistema
- ✅ Avisos quando usa executável incompatível
- ✅ Fallback inteligente quando executável ideal não existe

**Scripts atualizados:**
- `install_service_windows.bat` - Seleção de arquitetura na instalação
- `uninstall_service_windows.bat` - Mostra arquitetura instalada
- `manage_service_windows.bat` - Exibe arquitetura no menu

**Documentação:**
- `docs/setup/SELECAO_ARQUITETURA.md` - Guia completo sobre 32/64 bits
- `docs/setup/WINDOWS_SERVICE.md` - Atualizado com novo fluxo

**Funcionalidades:**
- Detecta PROCESSOR_ARCHITEW6432 e PROCESSOR_ARCHITECTURE
- Lista executáveis disponíveis (win32.exe e win64.exe)
- Permite seleção manual ou automática (padrão)
- Valida existência do executável escolhido
- Exibe aviso se arquitetura não é ideal

**Opções de seleção:**
- [A] Automático (recomendado) - Detecta e usa executável correto
- [1] Manual 32 bits - Força uso de win32.exe
- [2] Manual 64 bits - Força uso de win64.exe

---

## [3.0.0.3] - 2025-11-24

### 🛡️ Segurança

#### Fail2Ban Middleware
- ✅ Implementado middleware de proteção contra ataques de força bruta e scanning
- ✅ Proteção contra scanning (404s): 10 tentativas em 5min = ban de 1h
- ✅ Proteção contra brute force (401s): 5 tentativas em 5min = ban de 2h
- ✅ Rastreamento em memória com limpeza automática
- ✅ Thread-safe usando sync.RWMutex
- ✅ Logs detalhados de IPs banidos e tentativas bloqueadas

#### Scripts de Monitoramento
- ✅ `monitor_security.sh` (Linux) - Análise completa de segurança
- ✅ `monitor_security.ps1` (Windows) - Análise completa de segurança
- ✅ `watch_security.sh` (Linux) - Monitoramento em tempo real
- ✅ Detecta IPs suspeitos com múltiplos 404s
- ✅ Lista IPs banidos pelo Fail2Ban
- ✅ Identifica falhas de autenticação
- ✅ Mostra paths mais atacados
- ✅ Calcula estatísticas e fornece recomendações

#### Documentação de Segurança
- ✅ `docs/ANALISE_SEGURANCA.md` - Análise completa dos ataques detectados
- ✅ `docs/setup/MONITORAMENTO_SEGURANCA.md` - Guia completo de monitoramento

### 📝 Detalhes Técnicos

**Middleware Fail2Ban (`internal/middleware/fail2ban.go`):**
- Estrutura `IPTracker` com rastreamento de tentativas falhas
- Dois trackers independentes: um para 404s e outro para 401s
- Método `IsBanned()` para verificar se IP está banido
- Método `RecordFailure()` para registrar tentativas e aplicar ban
- Cleanup automático a cada 5 minutos via goroutine
- Resposta 403 com mensagem clara ao usuário banido

**Integração:**
- Aplicado em `cmd/server/main.go`
- Aplicado em `internal/service/windows_service.go`
- Logs de configuração ao iniciar servidor

---

## [1.26.4.28] - 2025-11-24

### 🔧 Corrigido

#### Windows Service Support
- ✅ **CRÍTICO**: Implementado suporte adequado para Windows Service API
- ✅ Resolvido erro 1053 ("O serviço não respondeu à requisição de início ou controle em tempo hábil")
- ✅ Detecção automática de modo de execução (Console vs Service)
- ✅ Implementação da interface `svc.Handler` para responder ao Service Control Manager
- ✅ Integração com Windows Event Log para registro de eventos do serviço
- ✅ Graceful shutdown quando recebe comandos STOP/SHUTDOWN do Windows
- ✅ Scripts de instalação/desinstalação atualizados com registro de Event Log
- ✅ **Mudança automática de diretório de trabalho** para o diretório do executável
  - Corrige problema de `dbinit.ini` não encontrado
  - Serviços Windows iniciam em `C:\Windows\System32` por padrão
  - Código agora usa `os.Executable()` e `os.Chdir()` para definir diretório correto

#### Novos Componentes
- ✅ `internal/service/windows_service.go` - Implementação completa Windows Service API
- ✅ `cmd/server/service_windows.go` - Funções específicas Windows (build tag)
- ✅ `cmd/server/service_other.go` - Stubs para Linux/Mac (build tag)

#### Documentação
- ✅ `docs/setup/WINDOWS_SERVICE.md` - Guia completo de instalação e gerenciamento
- ✅ `docs/WINDOWS_SERVICE_UPDATE.md` - Guia de atualização com antes/depois

### 📝 Detalhes Técnicos

**Problema:** Executável Go comum não pode ser simplesmente registrado como serviço Windows com `sc create`. É necessário implementar a Windows Service API para responder aos comandos do Service Control Manager (SCM).

**Solução:**
- Uso de `golang.org/x/sys/windows/svc` para implementar interface Windows Service
- Detecção automática via `svc.IsWindowsService()` no `main()`
- Servidor HTTP executa em goroutine enquanto serviço monitora comandos do SCM
- Event Log integration via `golang.org/x/sys/windows/svc/eventlog`

**Compatibilidade:**
- ✅ Windows Server 2016+
- ✅ Windows 10/11
- ✅ Mantém compatibilidade com execução console (desenvolvimento)
- ✅ Build tags garantem que código Windows não afeta Linux/Mac

---

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
