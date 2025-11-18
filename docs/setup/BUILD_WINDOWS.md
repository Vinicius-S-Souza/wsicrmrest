# Compilação para Windows

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-17

Este guia explica como compilar o WSICRMREST para Windows 32 e 64 bits.

## Pré-requisitos

### No Windows

1. **Go 1.21 ou superior**
   - Download: https://golang.org/dl/
   - Adicione Go ao PATH durante a instalação

2. **GCC/MinGW (para CGO)**
   - Download MinGW-w64: https://www.mingw-w64.org/downloads/
   - Ou via Chocolatey: `choco install mingw`
   - Adicione MinGW ao PATH

3. **Oracle Instant Client**
   - Download: https://www.oracle.com/database/technologies/instant-client/downloads.html
   - Escolha a versão 32 ou 64 bits conforme necessário
   - Extraia e adicione ao PATH
   - Configure variável `LD_LIBRARY_PATH` (ou `PATH` no Windows)

### No Linux (Cross-Compilation)

Para compilar executáveis Windows a partir do Linux, você precisa:

1. **Go 1.21 ou superior**
   ```bash
   go version
   ```

2. **MinGW Cross-Compiler**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install gcc-mingw-w64-i686 gcc-mingw-w64-x86-64

   # Fedora/RHEL
   sudo dnf install mingw32-gcc mingw64-gcc
   ```

3. **Oracle Instant Client (Windows version)**
   - Baixe as versões Windows do Oracle Instant Client
   - Você precisará incluí-las no pacote de distribuição

## Métodos de Compilação

### Opção 1: Script Batch (Windows)

O método mais simples para usuários Windows:

```batch
cd wsicrmrest
scripts\build_windows.bat
```

O script irá:
1. Verificar se Go está instalado
2. Criar o diretório `build\` se não existir
3. Perguntar qual versão você deseja compilar:
   - 32 bits
   - 64 bits
   - Ambas

**Executáveis gerados:**
- `build\wsicrmrest_win32.exe` (32 bits)
- `build\wsicrmrest_win64.exe` (64 bits)

### Opção 2: Makefile (Linux/WSL)

Para compilar a partir de Linux ou WSL:

```bash
# Compilar ambas as versões Windows
make build-windows

# Ou compilar versões específicas
make build-windows-32
make build-windows-64
```

### Opção 3: Linha de Comando Manual

#### Windows 32 bits
```bash
set GOOS=windows
set GOARCH=386
set CGO_ENABLED=1
go build -o build\wsicrmrest_win32.exe cmd\server\main.go
```

#### Windows 64 bits
```bash
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -o build\wsicrmrest_win64.exe cmd\server\main.go
```

#### Cross-compilation no Linux (32 bits)
```bash
GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc \
  go build -o build/wsicrmrest_win32.exe cmd/server/main.go
```

#### Cross-compilation no Linux (64 bits)
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -o build/wsicrmrest_win64.exe cmd/server/main.go
```

## Executando no Windows

### Usando o Script
```batch
scripts\run_windows.bat
```

O script automaticamente:
- Detecta a arquitetura do Windows (32 ou 64 bits)
- Verifica se `dbinit.ini` existe
- Cria o diretório `log\` se necessário
- Executa o binário correto

### Execução Manual

1. **Configure o banco de dados:**
   ```batch
   copy dbinit.ini.example dbinit.ini
   notepad dbinit.ini
   ```

2. **Execute o servidor:**
   ```batch
   REM Para 64 bits
   build\wsicrmrest_win64.exe

   REM Para 32 bits
   build\wsicrmrest_win32.exe
   ```

## Distribuição

### Arquivos Necessários

Para distribuir o aplicativo, inclua:

```
wsicrmrest/
├── build/
│   ├── wsicrmrest_win32.exe  (ou wsicrmrest_win64.exe)
├── dbinit.ini.example         (renomear para dbinit.ini)
├── log/                       (diretório vazio)
├── oracle_instant_client/     (bibliotecas Oracle)
│   ├── oci.dll
│   ├── oraociei*.dll
│   └── ...
└── README.md
```

### Configuração no Cliente

1. **Instalar Oracle Instant Client:**
   - Extrair na pasta do aplicativo ou em `C:\oracle\instantclient`
   - Adicionar ao PATH ou colocar DLLs junto com o executável

2. **Configurar tnsnames.ora:**
   - Criar arquivo em `C:\oracle\instantclient\network\admin\tnsnames.ora`
   - Ou definir variável `TNS_ADMIN` apontando para o diretório

3. **Configurar dbinit.ini:**
   - Copiar `dbinit.ini.example` para `dbinit.ini`
   - Editar com credenciais corretas:
     ```ini
     [database]
     tns_name = ORCL
     username = wsuser
     password = wspass
     ```

4. **Executar:**
   ```batch
   wsicrmrest_win64.exe
   ```

## Troubleshooting

### Erro: "go: not found" ou "Go não encontrado no PATH"
**Solução:** Instale Go e adicione ao PATH:
```batch
set PATH=%PATH%;C:\Go\bin
```

### Erro: "gcc: not found" ou "MinGW não encontrado"
**Solução:** Instale MinGW e adicione ao PATH:
```batch
set PATH=%PATH%;C:\mingw64\bin
```

### Erro: "OCI.dll not found"
**Solução:**
1. Instale Oracle Instant Client
2. Adicione ao PATH:
   ```batch
   set PATH=%PATH%;C:\oracle\instantclient_19_x
   ```

### Erro: "ORA-12154: TNS:could not resolve the connect identifier"
**Solução:**
1. Verifique se `tnsnames.ora` está configurado
2. Defina `TNS_ADMIN`:
   ```batch
   set TNS_ADMIN=C:\oracle\instantclient\network\admin
   ```

### Erro CGO durante compilação
**Solução:** Certifique-se de que:
```batch
set CGO_ENABLED=1
set CC=gcc
```

## Diferenças entre 32 e 64 bits

| Aspecto | 32 bits | 64 bits |
|---------|---------|---------|
| **Memória máxima** | ~4 GB | Ilimitada (praticamente) |
| **Performance** | Menor | Maior |
| **Compatibilidade** | Roda em Windows 32 e 64 bits | Só roda em Windows 64 bits |
| **Oracle Client** | Requer Instant Client 32 bits | Requer Instant Client 64 bits |
| **Tamanho executável** | Menor | Maior |

**Recomendação:** Use a versão 64 bits para servidores modernos. Use 32 bits apenas se necessário para compatibilidade com sistemas antigos.

## Notas Importantes

- **CGO está habilitado:** O driver Oracle (`godror`) requer CGO, então a compilação precisa do GCC/MinGW
- **Arquitetura do Oracle Client:** A arquitetura do executável (32/64 bits) deve corresponder à do Oracle Instant Client instalado
- **Cross-compilation:** É possível compilar Windows a partir do Linux, mas requer MinGW cross-compiler instalado
- **Variáveis de ambiente:** Em produção, considere criar um arquivo `.bat` com as variáveis necessárias (PATH, TNS_ADMIN, etc.)

## Automação

### Script de Deploy Completo

Crie `scripts\deploy_windows.bat`:

```batch
@echo off
echo Criando pacote de distribuição...
mkdir dist
mkdir dist\log

REM Copiar executável
copy build\wsicrmrest_win64.exe dist\

REM Copiar configuração
copy dbinit.ini.example dist\dbinit.ini

REM Copiar Oracle Instant Client (ajuste o caminho)
xcopy /E /I C:\oracle\instantclient_19_x dist\oracle_instant_client

echo Pacote criado em: dist\
echo Configure dist\dbinit.ini antes de distribuir!
pause
```

## Referências

- Go Cross Compilation: https://go.dev/doc/install/source#environment
- MinGW-w64: https://www.mingw-w64.org/
- Oracle Instant Client: https://www.oracle.com/database/technologies/instant-client.html
- godror (Oracle driver): https://github.com/godror/godror
