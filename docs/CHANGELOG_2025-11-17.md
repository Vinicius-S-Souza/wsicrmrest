# Changelog - 2025-11-17

**Data:** 2025-11-17
**Autor:** Claude Code
**Versão:** 1.26.4.27

## Resumo das Alterações

Esta atualização remove o suporte multi-banco de dados e padroniza o sistema para usar **exclusivamente Oracle**, além de adicionar suporte completo para compilação em **Windows 32 e 64 bits**.

## Mudanças Principais

### 1. Remoção do Parâmetro `driver`

**Motivo:** O sistema foi projetado exclusivamente para Oracle Database. O parâmetro `driver` que permitia escolher entre SQL Server (1) e Oracle (2) foi removido para simplificar a configuração.

**Arquivos modificados:**

- **`dbinit.ini.example`**
  - ❌ Removido: `driver = 2`
  - ✅ Mantido: Apenas configurações Oracle (tns_name, username, password)

- **`internal/config/config.go`**
  - ❌ Removido: Campo `Driver int` da struct `DatabaseConfig`
  - ❌ Removido: Linha `Driver: cfg.Section("database").Key("driver").MustInt(2)`
  - ✅ Atualizado: Comentário da struct para "configurações do banco de dados Oracle"

- **`internal/database/database.go`**
  - ❌ Removido: Validação `if cfg.Database.Driver != 2`
  - ✅ Simplificado: Conexão direta com Oracle via godror

**Impacto:**
- **Configuração mais simples:** Usuários não precisam mais especificar o driver
- **Código mais limpo:** Menos validações e menos complexidade condicional
- **Sem perda de funcionalidade:** Oracle sempre foi o único banco suportado na prática

### 2. Suporte para Windows 32 e 64 bits

**Motivo:** Permitir que o sistema seja compilado e executado em ambientes Windows, tanto em arquiteturas de 32 quanto 64 bits.

#### Novos Targets no Makefile

**Arquivo:** `Makefile`

Adicionados os seguintes comandos:

```makefile
# Compila ambas as versões Windows
make build-windows

# Compila apenas Windows 32 bits
make build-windows-32

# Compila apenas Windows 64 bits
make build-windows-64
```

**Detalhes técnicos:**
- Usa cross-compilation com `GOOS=windows` e `GOARCH=386|amd64`
- Requer MinGW (i686-w64-mingw32-gcc ou x86_64-w64-mingw32-gcc)
- CGO habilitado (necessário para o driver Oracle godror)
- Binários gerados:
  - `build/wsicrmrest_win32.exe`
  - `build/wsicrmrest_win64.exe`

#### Scripts Windows

**Arquivo:** `scripts/build_windows.bat`
- Script interativo para compilação no Windows
- Menu de opções: 32 bits, 64 bits ou ambos
- Verifica instalação do Go
- Mostra mensagens de erro claras

**Arquivo:** `scripts/run_windows.bat`
- Detecta automaticamente arquitetura do Windows (32/64 bits)
- Verifica existência do `dbinit.ini`
- Cria diretório `log/` automaticamente
- Executa o binário correspondente

#### Documentação

**Arquivo:** `docs/setup/BUILD_WINDOWS.md`

Documentação completa incluindo:

1. **Pré-requisitos:**
   - Go 1.21+
   - MinGW (GCC para Windows)
   - Oracle Instant Client (32 ou 64 bits)

2. **Métodos de compilação:**
   - Script Batch (Windows)
   - Makefile (Linux/WSL cross-compilation)
   - Linha de comando manual

3. **Execução e distribuição:**
   - Como executar no Windows
   - Arquivos necessários para distribuição
   - Configuração no cliente

4. **Troubleshooting:**
   - Erros comuns e soluções
   - Problemas com Oracle Client
   - Problemas de compilação CGO

5. **Tabela comparativa:** 32 vs 64 bits

### 3. Alteração da Pasta de Build

**Motivo:** Organizar melhor os artefatos de compilação em uma pasta dedicada.

**Mudanças:**

- **Antes:** Binário compilado na raiz do projeto (`./wsicrmrest`)
- **Depois:** Binário compilado em `build/` (`build/wsicrmrest`)

**Arquivos modificados:**

- `Makefile`:
  - `BUILD_DIR` alterado de `.` para `build`
  - Comando `clean` agora remove `build/` inteiro
  - Comando `run` atualizado para `./build/wsicrmrest`

- `.gitignore`:
  - Adicionado `build/` e `dist/` (para pacotes de distribuição)

**Benefícios:**
- Projeto mais organizado
- Fácil de limpar (apenas deletar `build/`)
- Suporta múltiplos binários (Linux, Win32, Win64) na mesma pasta

### 4. Atualizações na Documentação

**Arquivo:** `CLAUDE.md`

Atualizações na seção "Build and Run Commands":
- Dividido em seções "Linux/WSL" e "Windows"
- Adicionados novos comandos de build para Windows
- Nota sobre localização dos binários em `build/`

Atualizações na seção "Configuration Loading Priority":
- Clarificado que database connection é "Oracle only"

Adicionada referência à nova documentação:
- `docs/setup/BUILD_WINDOWS.md`: Windows compilation guide

## Migração para Usuários Existentes

### Se você já tem o projeto configurado:

1. **Atualizar dbinit.ini (OPCIONAL):**
   ```bash
   # A linha "driver = 2" pode ser removida, mas não é obrigatório
   # O sistema simplesmente ignora esse parâmetro agora
   ```

2. **Recompilar:**
   ```bash
   # Linux
   make clean
   make build

   # Windows
   scripts\build_windows.bat
   ```

3. **Novo caminho do binário:**
   ```bash
   # Antes:
   ./wsicrmrest

   # Depois:
   ./build/wsicrmrest  # ou build\wsicrmrest_win64.exe no Windows
   ```

4. **Usar o Makefile normalmente:**
   ```bash
   make run  # Funciona automaticamente com o novo caminho
   ```

### Se você está iniciando um novo projeto:

1. **Copie o exemplo de configuração:**
   ```bash
   cp dbinit.ini.example dbinit.ini
   ```

2. **Edite dbinit.ini** (NÃO precisa mais do parâmetro `driver`):
   ```ini
   [database]
   tns_name = ORCL
   username = wsuser
   password = wspass
   ```

3. **Compile e execute:**
   ```bash
   # Linux
   make build
   make run

   # Windows
   scripts\build_windows.bat
   scripts\run_windows.bat
   ```

## Compatibilidade

### ✅ Mantido (100% compatível)

- Todas as APIs REST existentes
- Formato de resposta JSON
- Autenticação JWT
- Logs (arquivo e banco de dados)
- CORS
- Estrutura de tabelas Oracle
- Variáveis globais hardcoded

### ⚠️ Alterado (não quebra funcionalidade)

- Parâmetro `driver` no `dbinit.ini` é **ignorado** (pode ser removido mas não obrigatório)
- Caminho do binário mudou de `./wsicrmrest` para `build/wsicrmrest`

### ❌ Removido (sem impacto prático)

- Suporte teórico a SQL Server (nunca foi implementado de fato)
- Validação do campo `Driver` na configuração

## Testes Realizados

- ✅ Compilação Linux bem-sucedida
- ✅ Binário gerado em `build/wsicrmrest` (43MB)
- ✅ Configuração sem parâmetro `driver` funciona corretamente
- ⚠️ Compilação Windows requer MinGW instalado (documentado)

## Próximos Passos Recomendados

1. **Testar em ambiente Windows real:**
   - Compilar com `scripts\build_windows.bat`
   - Validar com Oracle Instant Client 32 e 64 bits

2. **Atualizar CHANGELOG.md principal:**
   - Incorporar essas mudanças no changelog oficial

3. **Criar release:**
   - Tag Git: `v1.26.4.27-windows-support`
   - Incluir binários pré-compilados (Linux + Win32 + Win64)

4. **CI/CD (opcional):**
   - Automatizar builds multi-plataforma
   - GitHub Actions ou GitLab CI

## Referências

- [BUILD_WINDOWS.md](setup/BUILD_WINDOWS.md) - Guia completo de compilação Windows
- [CLAUDE.md](../CLAUDE.md) - Documentação atualizada do projeto
- [Makefile](../Makefile) - Novos targets de build
