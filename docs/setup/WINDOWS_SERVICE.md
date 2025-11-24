# Instalação como Serviço Windows

**Data de criação:** 2025-11-24
**Última atualização:** 2025-11-24

Este documento descreve como instalar, configurar e gerenciar o WSICRMREST como um serviço do Windows (Windows Service).

## ⚠️ Versão Atual

**Versão:** 3.0.0.2 (2025-11-24)

**Correções Incluídas:**
- ✅ Windows Service API completa (resolve erro 1053)
- ✅ Mudança automática de diretório de trabalho (resolve problema de `dbinit.ini` não encontrado)
- ✅ Event Log integration
- ✅ Graceful shutdown

## 📋 Índice

- [Requisitos](#requisitos)
- [Arquitetura do Serviço](#arquitetura-do-serviço)
- [Instalação](#instalação)
- [Gerenciamento](#gerenciamento)
- [Solução de Problemas](#solução-de-problemas)
- [Desinstalação](#desinstalação)

---

## Requisitos

### Sistema Operacional
- Windows Server 2016 ou superior
- Windows 10/11 (para desenvolvimento)

### Permissões
- **Administrador** é obrigatório para:
  - Instalar/desinstalar serviços
  - Iniciar/parar serviços
  - Registrar Event Log sources

### Pré-requisitos
1. Aplicação compilada para Windows (64 bits):
   ```batch
   scripts\build_windows.bat
   ```
   Ou via Make:
   ```bash
   make build-windows-64
   ```

2. Arquivo `dbinit.ini` configurado no diretório da aplicação

3. Acesso ao banco de dados Oracle configurado

4. Tabela `ORGANIZADOR` populada no banco de dados

---

## Arquitetura do Serviço

### Detecção Automática de Modo

O executável `wsicrmrest_win64.exe` funciona em **dois modos**:

1. **Modo Console** (execução direta):
   - Inicia servidor HTTP/HTTPS normalmente
   - Logs aparecem no terminal
   - Pode ser interrompido com Ctrl+C

2. **Modo Serviço** (registrado no Windows):
   - Detecta automaticamente que está rodando como serviço
   - Implementa Windows Service API
   - Responde a comandos do Service Control Manager
   - Registra eventos no Event Log do Windows

### Componentes Principais

```
wsicrmrest/
├── cmd/server/
│   ├── main.go                    # Ponto de entrada, detecta modo de execução
│   ├── service_windows.go         # Funções específicas Windows (build tag)
│   └── service_other.go           # Stubs para Linux/Mac (build tag)
├── internal/service/
│   └── windows_service.go         # Implementação Windows Service API
└── scripts/
    ├── install_service_windows.bat    # Instalador
    ├── uninstall_service_windows.bat  # Desinstalador
    └── manage_service_windows.bat     # Gerenciador
```

### Event Log

O serviço registra eventos importantes no **Windows Event Log**:

- **Local**: Applications and Services Logs → Application
- **Source**: WSICRMREST
- **Tipos de Eventos**:
  - ✅ **Informação**: Inicialização, parada normal, operações bem-sucedidas
  - ⚠️ **Aviso**: Configurações não recomendadas (ex: TLS desabilitado)
  - ❌ **Erro**: Falhas de inicialização, erros de banco de dados, crashes

---

## Instalação

### Método 1: Script Automatizado (Recomendado)

1. **Abra PowerShell ou CMD como Administrador**

2. **Navegue até o diretório do projeto:**
   ```batch
   cd C:\CRM\WSICRMREST
   ```

3. **Execute o instalador:**
   ```batch
   scripts\install_service_windows.bat
   ```

4. **O script irá:**
   - ✅ Detectar automaticamente a arquitetura do Windows (32 ou 64 bits)
   - ✅ Listar executáveis disponíveis (win32.exe e/ou win64.exe)
   - ✅ Permitir seleção manual ou detecção automática
   - ✅ Verificar permissões de administrador
   - ✅ Validar existência do executável e `dbinit.ini`
   - ✅ Registrar Event Log source
   - ✅ Criar o serviço Windows
   - ✅ Configurar início automático (com delay)
   - ✅ Configurar recuperação automática em caso de falha
   - ✅ Perguntar se deseja iniciar o serviço

**Exemplo de saída:**
```
Detectando arquitetura do Windows...
Sistema detectado: Windows 64 bits

Executáveis disponíveis:
  [1] wsicrmrest_win32.exe (32 bits)
  [2] wsicrmrest_win64.exe (64 bits)
  [A] Detectar automaticamente (recomendado)

Escolha o executável [1/2/A - padrão A]:

============================================
Configurações do Serviço:
============================================
  Nome do Serviço: WSICRMREST
  Nome de Exibição: WSICRMREST API Service do Sistema ICRM
  Executável: C:\CRM\WSICRMREST\wsicrmrest_win64.exe
  Arquitetura: 64 bits (auto)
  Diretório de Trabalho: C:\CRM\WSICRMREST
```

### Método 2: Manual via `sc` Command

```batch
REM 1. Registrar Event Log (como Admin)
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\WSICRMREST" /v EventMessageFile /t REG_EXPAND_SZ /d "%SystemRoot%\System32\EventCreate.exe" /f
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\WSICRMREST" /v TypesSupported /t REG_DWORD /d 7 /f

REM 2. Criar serviço
sc create WSICRMREST binPath= "C:\CRM\WSICRMREST\wsicrmrest_win64.exe" start= auto DisplayName= "WSICRMREST API Service do Sistema ICRM"

REM 3. Configurar descrição
sc description WSICRMREST "Web Service REST API para Integração com Sistema ICRM"

REM 4. Configurar recuperação automática
sc failure WSICRMREST reset= 86400 actions= restart/60000/restart/60000/restart/60000

REM 5. Configurar início atrasado
sc config WSICRMREST start= delayed-auto
```

### Método 3: Via Interface Gráfica (services.msc)

1. Primeiro, crie o serviço via linha de comando (método 2 acima)
2. Abra `services.msc`
3. Localize "WSICRMREST API Service do Sistema ICRM"
4. Configure manualmente propriedades adicionais:
   - Recovery options
   - Log On account (se necessário)
   - Dependencies

---

## Gerenciamento

### Iniciar o Serviço

**Linha de comando:**
```batch
sc start WSICRMREST
```

**PowerShell:**
```powershell
Start-Service WSICRMREST
```

**Interface Gráfica:**
1. Abra `services.msc`
2. Localize "WSICRMREST API Service do Sistema ICRM"
3. Clique com botão direito → **Iniciar**

### Parar o Serviço

**Linha de comando:**
```batch
sc stop WSICRMREST
```

**PowerShell:**
```powershell
Stop-Service WSICRMREST
```

### Verificar Status

**Linha de comando:**
```batch
sc query WSICRMREST
```

**PowerShell:**
```powershell
Get-Service WSICRMREST
```

**Script de Gerenciamento:**
```batch
scripts\manage_service_windows.bat
```
Menu interativo com opções de start/stop/restart/status.

### Reiniciar o Serviço

**PowerShell:**
```powershell
Restart-Service WSICRMREST
```

**Linha de comando:**
```batch
sc stop WSICRMREST && timeout /t 3 /nobreak && sc start WSICRMREST
```

---

## Solução de Problemas

### Erro 1053: "O serviço não respondeu à requisição de início ou controle em tempo hábil"

**Causa:** Este é o erro que você estava enfrentando. Ocorre quando o executável não implementa a Windows Service API corretamente.

**Solução:** A implementação agora inclui:
- ✅ Detecção automática de modo serviço (`svc.IsWindowsService()`)
- ✅ Implementação da interface `svc.Handler`
- ✅ Resposta adequada a comandos do SCM (Start, Stop, Shutdown)
- ✅ Event Log integration

**Após atualizar o código, recompile:**
```batch
make build-windows-64
```

### Serviço não inicia: "dbinit.ini não encontrado"

**Causa:** Versões antigas não mudavam o diretório de trabalho automaticamente.

**Solução:** Atualize para versão 3.0.0.2 ou superior que inclui mudança automática de diretório.

**Verificação:**
```batch
cd C:\CRM\WSICRMREST
dir dbinit.ini
dir wsicrmrest_win64.exe
```

Ambos os arquivos devem estar no **mesmo diretório**.

### Serviço não inicia: Erro de conexão com banco de dados

**Verificações:**

1. **Arquivo `dbinit.ini` existe e está no diretório correto?**
   ```batch
   cd C:\CRM\WSICRMREST
   dir dbinit.ini
   ```

2. **Configurações do Oracle estão corretas?**
   ```batch
   notepad dbinit.ini
   ```
   Verifique:
   - TNS Name
   - Username/Password
   - LD_LIBRARY_PATH (se aplicável)

3. **Tabela ORGANIZADOR tem dados?**
   ```sql
   SELECT * FROM ORGANIZADOR WHERE ORGCODIGO > 0;
   ```

4. **Verifique os logs:**
   ```batch
   type log\wsicrmrest_YYYY-MM-DD.log
   ```

### Serviço inicia mas para sozinho

**Verifique Event Log do Windows:**

1. Abra **Event Viewer** (`eventvwr.msc`)
2. Navegue: Windows Logs → **Application**
3. Filtre por Source: **WSICRMREST**
4. Procure por erros (ícone vermelho ❌)

**Causas comuns:**
- Porta já em uso (8080 ou configurada)
- Certificado TLS inválido/inexistente
- Conexão com banco caiu após inicialização

### Permissão negada ao instalar

**Erro:**
```
ERRO: Este script precisa ser executado como Administrador!
```

**Solução:**
1. Feche o prompt atual
2. Clique com **botão direito** no script → **Executar como administrador**
3. Ou abra CMD/PowerShell como Admin primeiro

### Event Log não aparece

**Reinstalar Event Log source manualmente:**

```batch
REM Como Administrador
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\WSICRMREST" /v EventMessageFile /t REG_EXPAND_SZ /d "%SystemRoot%\System32\EventCreate.exe" /f
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\WSICRMREST" /v TypesSupported /t REG_DWORD /d 7 /f
```

Depois reinicie o serviço Event Log:
```batch
net stop eventlog
net start eventlog
```

---

## Desinstalação

### Método 1: Script Automatizado (Recomendado)

```batch
REM Como Administrador
scripts\uninstall_service_windows.bat
```

O script irá:
1. ✅ Verificar se serviço existe
2. ✅ Parar o serviço se estiver rodando
3. ✅ Remover o serviço
4. ✅ Limpar registros do Event Log

### Método 2: Manual

```batch
REM 1. Parar o serviço
sc stop WSICRMREST

REM 2. Aguardar parada completa
timeout /t 5 /nobreak

REM 3. Remover serviço
sc delete WSICRMREST

REM 4. Remover Event Log source
reg delete "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\WSICRMREST" /f
```

### Limpeza Completa

Após desinstalar o serviço:

```batch
REM Remover diretório completo (cuidado!)
rmdir /s /q C:\CRM\WSICRMREST

REM Ou manter logs e configurações
del /q C:\CRM\WSICRMREST\*.exe
```

---

## Configuração Avançada

### Configurar conta de execução específica

Por padrão, o serviço roda como **Local System**. Para usar conta específica:

```batch
sc config WSICRMREST obj= "DOMAIN\Username" password= "Password"
```

**Atenção:** A conta precisa de:
- Permissão "Log on as a service"
- Acesso de leitura ao diretório de instalação
- Acesso ao banco de dados Oracle

### Ajustar timeout de inicialização

Se o banco de dados é lento para conectar:

```batch
REM Aumentar timeout para 120 segundos
sc config WSICRMREST start= delayed-auto
```

### Configurar dependências

Se precisa de outros serviços iniciados primeiro:

```batch
sc config WSICRMREST depend= "OracleServiceXE/Tcpip"
```

---

## Logs e Monitoramento

### Locais de Log

1. **Logs da Aplicação:**
   - `C:\CRM\WSICRMREST\log\wsicrmrest_YYYY-MM-DD.log`
   - Formato JSON estruturado
   - Rotação diária automática

2. **Event Log do Windows:**
   - Event Viewer → Application → WSICRMREST
   - Eventos críticos de serviço

### Monitoramento via PowerShell

```powershell
# Status em tempo real
while ($true) {
    Clear-Host
    Get-Service WSICRMREST | Format-List
    Start-Sleep -Seconds 5
}

# Últimos eventos do Event Log
Get-EventLog -LogName Application -Source WSICRMREST -Newest 10
```

---

## Referências

- [Windows Service API Documentation](https://learn.microsoft.com/en-us/windows/win32/services/services)
- [sc.exe Command Reference](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/sc-create)
- [Event Log Management](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-logging)
- [golang.org/x/sys/windows/svc Package](https://pkg.go.dev/golang.org/x/sys/windows/svc)

---

## Changelog

| Data       | Versão | Alteração                                      |
|------------|--------|------------------------------------------------|
| 2025-11-24 | 1.0.0  | Criação da documentação e implementação inicial do Windows Service |
