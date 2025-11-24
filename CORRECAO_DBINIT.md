# Correção: Problema de Diretório de Trabalho

**Data:** 2025-11-24
**Problema:** `open dbinit.ini: The system cannot find the file specified`

---

## 🎉 Boa Notícia!

O erro **1053** foi resolvido! Agora o serviço está **tentando iniciar** corretamente, mas encontrou um problema secundário: não consegue localizar o arquivo `dbinit.ini`.

### Por que isso aconteceu?

Serviços Windows iniciam com o diretório de trabalho padrão em `C:\Windows\System32`. Quando o código tenta abrir `dbinit.ini`, ele procura nesse diretório em vez do diretório onde o executável está.

---

## ✅ Correção Aplicada

Foram feitas duas mudanças:

### 1. Código Go - Mudança de Diretório Automática

**Arquivo:** `internal/service/windows_service.go`

```go
// Mudar para o diretório do executável
// Serviços Windows iniciam em C:\Windows\System32 por padrão
exePath, err := os.Executable()
if err == nil {
    exeDir := filepath.Dir(exePath)
    os.Chdir(exeDir)
}
```

Agora, antes de carregar `dbinit.ini`, o serviço muda automaticamente para o diretório onde o executável está localizado.

### 2. Script de Instalação - Path Absoluto

**Arquivo:** `scripts/install_service_windows.bat`

```batch
REM Converter para path absoluto
pushd %~dp0
set WORK_DIR=%CD%
set BINARY_PATH=%CD%\wsicrmrest_win64.exe
popd
```

O script agora garante que sempre usa caminhos absolutos.

---

## 🔄 Passos para Aplicar a Correção

### No Windows (onde está instalado):

1. **Parar o serviço:**
   ```batch
   sc stop WSICRMREST
   ```

2. **Remover serviço (opcional, mas recomendado):**
   ```batch
   sc delete WSICRMREST
   ```

### No Linux/WSL (onde compila):

3. **Recompilar com a correção:**
   ```bash
   cd /home/vinicius/projetos/wsicrmrest
   make build-windows-64
   ```

   Ou:
   ```bash
   go mod tidy
   GOOS=windows GOARCH=amd64 go build -o build/wsicrmrest_win64.exe ./cmd/server
   ```

### No Windows (instalar nova versão):

4. **Copiar novo executável:**
   ```batch
   copy /Y \\wsl$\Ubuntu\home\vinicius\projetos\wsicrmrest\build\wsicrmrest_win64.exe C:\CRM\WSICRMREST\wsicrmrest_win64.exe
   ```

   Ou copie manualmente de `build/wsicrmrest_win64.exe` para o diretório de instalação.

5. **Verificar que dbinit.ini está no mesmo diretório:**
   ```batch
   cd C:\CRM\WSICRMREST
   dir dbinit.ini
   ```

   Deve aparecer o arquivo. Se não existir:
   ```batch
   copy dbinit.ini.example dbinit.ini
   notepad dbinit.ini
   ```

6. **Reinstalar o serviço (como Administrador):**
   ```batch
   cd C:\CRM\WSICRMREST
   scripts\install_service_windows.bat
   ```

7. **Verificar status:**
   ```batch
   sc query WSICRMREST
   ```

   Deve mostrar:
   ```
   STATE              : 4  RUNNING
   ```

8. **Testar API:**
   ```batch
   curl http://localhost:8080/wsteste
   ```

---

## 🔍 Como Verificar que Funcionou

### 1. Event Log do Windows

Abra Event Viewer (`eventvwr.msc`):
- Windows Logs → Application
- Filtre por Source: **WSICRMREST**

Você deve ver eventos como:
- ✅ "Serviço WSICRMREST iniciado com sucesso"
- ✅ "Conexão com banco de dados estabelecida com sucesso"
- ✅ "Dados do organizador carregados com sucesso"

### 2. Log da Aplicação

```batch
type C:\CRM\WSICRMREST\log\wsicrmrest_2025-11-24.log
```

Deve conter linhas como:
```json
{"level":"info","msg":"Iniciando WSICRMREST como Windows Service","version":"Versão 3.0.0.2 (GO)"}
{"level":"info","msg":"Conexão com banco de dados estabelecida com sucesso"}
{"level":"info","msg":"Dados do organizador carregados com sucesso","codigo":1,"nome":"Empresa Exemplo"}
{"level":"info","msg":"Servidor HTTP iniciado","port":"8080"}
```

### 3. Status do Serviço

```batch
sc query WSICRMREST
```

Saída esperada:
```
SERVICE_NAME: WSICRMREST
TYPE               : 10  WIN32_OWN_PROCESS
STATE              : 4  RUNNING
                        (STOPPABLE, NOT_PAUSABLE, ACCEPTS_SHUTDOWN)
WIN32_EXIT_CODE    : 0  (0x0)
SERVICE_EXIT_CODE  : 0  (0x0)
CHECKPOINT         : 0x0
WAIT_HINT          : 0x0
```

---

## ⚠️ Problemas Comuns

### "dbinit.ini ainda não encontrado"

**Causa:** Arquivo não está no mesmo diretório do executável.

**Solução:**
```batch
cd C:\CRM\WSICRMREST
dir
```

Certifique-se que ambos estão no mesmo diretório:
- ✅ `wsicrmrest_win64.exe`
- ✅ `dbinit.ini`

### "Erro ao conectar ao banco de dados"

**Causa:** Configurações do Oracle em `dbinit.ini` estão incorretas.

**Solução:**
```batch
notepad C:\CRM\WSICRMREST\dbinit.ini
```

Verifique:
```ini
[database]
tns_name = SEU_TNS_NAME
username = seu_usuario
password = sua_senha
```

Teste a conexão manualmente:
```batch
sqlplus usuario/senha@TNS_NAME
```

### "Tabela ORGANIZADOR não encontrada"

**Causa:** Banco de dados não tem a tabela ou está vazia.

**Solução:**
```sql
-- Conectar ao banco
sqlplus usuario/senha@TNS_NAME

-- Verificar se tabela existe
SELECT COUNT(*) FROM ORGANIZADOR WHERE ORGCODIGO > 0;
```

Se retornar 0, insira pelo menos um registro.

---

## 📊 Diferença Visual

### Antes (erro):
```
[SC] StartService FALHA 1053:
O serviço não respondeu à requisição de início ou controle em tempo hábil.
```

### Agora (tentando iniciar):
```
Erro ao carregar configurações: erro ao ler arquivo dbinit.ini:
open dbinit.ini: The system cannot find the file specified.
```

### Após correção (sucesso):
```
SERVICE_NAME: WSICRMREST
STATE              : 4  RUNNING
```

---

## 🎯 Resumo

| Problema | Status | Ação |
|----------|--------|------|
| Erro 1053 (Windows Service API) | ✅ Resolvido | Implementado em versão anterior |
| Diretório de trabalho incorreto | ✅ Resolvido | Aplicar esta atualização |
| dbinit.ini não encontrado | ⚠️ Em teste | Recompilar e reinstalar |

---

## 📞 Próximo Passo

**Recompile o executável no WSL/Linux e teste no Windows conforme os passos acima.**

Após aplicar a correção, informe se:
- ✅ Serviço iniciou com sucesso
- ✅ API está respondendo
- ❌ Ainda há algum erro (compartilhe o erro)
