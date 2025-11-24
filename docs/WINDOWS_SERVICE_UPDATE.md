# Atualização: Suporte a Windows Service

**Data:** 2025-11-24

## 📋 Problema Resolvido

**Erro anterior:**
```
[SC] StartService FALHA 1053:
O serviço não respondeu à requisição de início ou controle em tempo hábil.
```

**Causa:** O executável Go não implementava a **Windows Service API** necessária para rodar como serviço Windows. Um executável comum não pode ser simplesmente registrado com `sc create` - ele precisa responder aos comandos do Service Control Manager (SCM).

---

## ✅ Solução Implementada

### 1. **Nova Arquitetura de Serviço**

Foram criados os seguintes componentes:

```
Novos arquivos criados:
├── internal/service/windows_service.go      # Implementação Windows Service API
├── cmd/server/service_windows.go            # Funções específicas Windows
└── cmd/server/service_other.go              # Stubs para Linux/Mac
```

### 2. **Detecção Automática de Modo**

O executável agora detecta automaticamente se está rodando como:
- **Serviço Windows** → Ativa Windows Service API
- **Console/Terminal** → Execução normal (desenvolvimento)

### 3. **Event Log Integration**

O serviço agora registra eventos no Event Log do Windows:
- ✅ Inicialização e parada
- ⚠️ Avisos de configuração
- ❌ Erros críticos

### 4. **Scripts Atualizados**

- `scripts/install_service_windows.bat` → Agora registra Event Log
- `scripts/uninstall_service_windows.bat` → Limpa Event Log

---

## 🔄 Passos para Atualização

### **Passo 1: Parar e Remover Serviço Existente (se instalado)**

```batch
REM Como Administrador
sc stop WSICRMREST
timeout /t 5 /nobreak
sc delete WSICRMREST
```

Ou use o script:
```batch
scripts\uninstall_service_windows.bat
```

### **Passo 2: Recompilar com Novo Código**

No **Linux/WSL**:
```bash
make build-windows-64
```

Ou no **Windows**:
```batch
scripts\build_windows.bat
```

Isso irá gerar: `build/wsicrmrest_win64.exe` com suporte a Windows Service.

**IMPORTANTE:** Esta nova versão inclui:
- ✅ Windows Service API completa
- ✅ Mudança automática de diretório de trabalho (resolve problema de `dbinit.ini` não encontrado)

### **Passo 3: Copiar Novo Executável**

Substitua o executável antigo pelo novo:

```batch
REM Exemplo: copiar de build/ para pasta de produção
copy /Y build\wsicrmrest_win64.exe C:\CRM\WSICRMREST\wsicrmrest_win64.exe
```

### **Passo 4: Reinstalar o Serviço**

```batch
REM Como Administrador, no diretório C:\CRM\WSICRMREST
scripts\install_service_windows.bat
```

O script irá:
1. ✅ Registrar Event Log source
2. ✅ Criar o serviço
3. ✅ Configurar recuperação automática
4. ✅ Oferecer para iniciar o serviço

### **Passo 5: Iniciar e Verificar**

**Iniciar:**
```batch
sc start WSICRMREST
```

**Verificar status:**
```batch
sc query WSICRMREST
```

Você deve ver:
```
STATE              : 4  RUNNING
```

**Verificar logs:**
```batch
type C:\CRM\WSICRMREST\log\wsicrmrest_YYYY-MM-DD.log
```

**Verificar Event Log:**
1. Abra Event Viewer (`eventvwr.msc`)
2. Windows Logs → Application
3. Filtre por Source: **WSICRMREST**

---

## 🧪 Testando

### Teste 1: Verificar que API está respondendo

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/wsteste" -Method GET
```

Ou via `curl`:
```batch
curl http://localhost:8080/wsteste
```

### Teste 2: Verificar que serviço responde a comandos

```batch
REM Parar
sc stop WSICRMREST

REM Aguardar
timeout /t 3 /nobreak

REM Verificar que parou
sc query WSICRMREST
REM Deve mostrar: STATE: 1 STOPPED

REM Iniciar novamente
sc start WSICRMREST

REM Verificar que iniciou
sc query WSICRMREST
REM Deve mostrar: STATE: 4 RUNNING
```

### Teste 3: Verificar reinicialização automática

```batch
REM Simular crash (forçar parada)
taskkill /F /IM wsicrmrest_win64.exe

REM Aguardar 60 segundos (configurado no install script)
timeout /t 65 /nobreak

REM Verificar que reiniciou automaticamente
sc query WSICRMREST
REM Deve mostrar: STATE: 4 RUNNING
```

---

## 📊 Diferenças Técnicas

### **Antes (não funcionava)**

```go
// main.go simplesmente iniciava servidor HTTP
func main() {
    router := gin.Default()
    router.Run(":8080") // Não responde a comandos do SCM
}
```

Quando registrado como serviço:
- ❌ Não responde a `SERVICE_CONTROL_STOP`
- ❌ SCM aguarda resposta → timeout 30s → erro 1053
- ❌ Não registra eventos no Event Log

### **Depois (funciona)**

```go
// main.go detecta modo
func main() {
    if runtime.GOOS == "windows" {
        isService, _ := svc.IsWindowsService()
        if isService {
            runAsWindowsService() // ← Implementa Windows Service API
            return
        }
    }
    // Execução normal (console)
    startHTTPServer()
}
```

**Windows Service API implementa:**
```go
type WindowsService struct { ... }

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) {
    changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

    for {
        select {
        case c := <-r:
            switch c.Cmd {
            case svc.Stop, svc.Shutdown:
                // Responde ao SCM
                changes <- svc.Status{State: svc.StopPending}
                ws.stop()
                return
            }
        }
    }
}
```

Quando registrado como serviço:
- ✅ Responde a `SERVICE_CONTROL_STOP` em <1s
- ✅ SCM recebe confirmação → sem timeout
- ✅ Registra eventos no Event Log
- ✅ Suporta recuperação automática

---

## 🔍 Verificação de Implementação

Você pode verificar que o novo código está compilado:

```batch
REM Executar manualmente para ver logs
wsicrmrest_win64.exe
```

Se tiver a mensagem no log:
```
"Iniciando WSICRMREST", "version": "Ver 1.26.4.27"
```

E o servidor **não mencionar** "como Windows Service", significa que está rodando em modo console (correto quando executado manualmente).

Quando rodando via `sc start`, o log deve mostrar:
```json
{"level":"info","msg":"Iniciando WSICRMREST como Windows Service","version":"Ver 1.26.4.27"}
```

---

## 📚 Documentação Adicional

- **Guia Completo:** `docs/setup/WINDOWS_SERVICE.md`
- **Solução de Problemas:** Ver seção no guia completo
- **Event Log:** Como monitorar eventos do serviço

---

## ⚙️ Dependências Adicionadas

O código agora usa:
```go
import (
    "golang.org/x/sys/windows/svc"
    "golang.org/x/sys/windows/svc/eventlog"
)
```

Essas dependências são automaticamente baixadas ao rodar:
```bash
go mod tidy
```

---

## 🎯 Próximos Passos Recomendados

Após instalar com sucesso:

1. **Configurar Monitoramento:**
   - Configure alertas no Event Log para eventos de erro
   - Use Windows Performance Monitor se necessário

2. **Backup da Configuração:**
   ```batch
   copy C:\CRM\WSICRMREST\dbinit.ini C:\CRM\WSICRMREST\dbinit.ini.backup
   ```

3. **Documentar Ambiente:**
   - Anotar porta configurada
   - Anotar se TLS está habilitado
   - Anotar account do serviço (se alterado de Local System)

4. **Testar Recuperação:**
   - Simule falha e verifique que reinicia automaticamente
   - Teste reinicialização do servidor Windows

---

## 📞 Suporte

Se encontrar problemas:

1. **Verificar Event Log** primeiro
2. **Verificar logs da aplicação** em `log/`
3. **Testar conexão com banco** via `sqlplus`
4. **Verificar dbinit.ini** está no diretório correto

---

## ✅ Checklist de Validação

Após seguir os passos acima:

- [ ] Serviço instalado sem erros
- [ ] Serviço iniciado com sucesso (`sc query WSICRMREST` → RUNNING)
- [ ] API responde em `http://localhost:8080/wsteste`
- [ ] Eventos aparecem no Event Log do Windows
- [ ] Logs da aplicação sendo criados em `log/`
- [ ] Serviço para e inicia corretamente via `sc stop/start`
- [ ] Recuperação automática funciona (testar com `taskkill /F`)
- [ ] Após reboot do Windows, serviço inicia automaticamente

Se todos os itens estiverem ✅, a instalação foi bem-sucedida!
