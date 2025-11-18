# Instalação do WSICRMREST como Serviço do Windows

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-18

## 📋 Índice

- [Pré-requisitos](#pré-requisitos)
- [Instalação Rápida](#instalação-rápida)
- [Instalação Detalhada](#instalação-detalhada)
- [Gerenciamento do Serviço](#gerenciamento-do-serviço)
- [Desinstalação](#desinstalação)
- [Troubleshooting](#troubleshooting)

---

## 🔧 Pré-requisitos

### 1. Compilar o Executável

Antes de instalar como serviço, você precisa compilar o executável Windows:

```batch
REM Opção 1: Usar script batch
cd scripts
build_windows.bat
REM Selecione opção 2 (Windows 64 bits) ou 3 (Ambas)

REM Opção 2: Usar Makefile (requer make no Windows)
make build-windows-64
```

**Resultado esperado:**
- Executável gerado em: `build\wsicrmrest_win64.exe`
- Tamanho aproximado: ~40-50 MB

### 2. Configurar dbinit.ini

```batch
REM Se dbinit.ini não existir, copie do exemplo
copy dbinit.ini.example dbinit.ini

REM Edite o arquivo e configure:
notepad dbinit.ini
```

**Configurações obrigatórias:**
```ini
[database]
tns_name = SEU_TNS_NAME
username = SEU_USUARIO
password = SUA_SENHA

[application]
port = 8080
environment = production
```

### 3. Verificar Oracle Client

O Oracle Instant Client deve estar instalado e configurado:

```batch
REM Verificar se ORACLE_HOME está definido
echo %ORACLE_HOME%

REM Verificar se tnsnames.ora existe
dir %ORACLE_HOME%\network\admin\tnsnames.ora
```

---

## ⚡ Instalação Rápida

### Passos Rápidos:

1. **Compilar** (se ainda não compilou):
   ```batch
   scripts\build_windows.bat
   ```

2. **Instalar como serviço**:
   ```batch
   REM Clique com botão direito e "Executar como administrador"
   scripts\install_service_windows.bat
   ```

3. **Pronto!** O serviço está instalado e pode ser iniciado.

---

## 📝 Instalação Detalhada

### Passo 1: Preparar Ambiente

```batch
REM 1. Criar estrutura de diretórios (se necessário)
mkdir log
mkdir build

REM 2. Verificar arquivos necessários
dir build\wsicrmrest_win64.exe
dir dbinit.ini
```

### Passo 2: Executar Script de Instalação

**IMPORTANTE:** Execute como Administrador!

1. Navegue até a pasta `scripts`
2. Clique com botão direito em `install_service_windows.bat`
3. Selecione **"Executar como administrador"**

**O script irá:**
- ✅ Verificar permissões de administrador
- ✅ Verificar se o executável existe
- ✅ Verificar se `dbinit.ini` existe
- ✅ Criar o serviço Windows
- ✅ Configurar início automático (atrasado)
- ✅ Configurar recuperação automática em caso de falha
- ✅ Perguntar se deseja iniciar imediatamente

### Passo 3: Verificar Instalação

```batch
REM Ver status do serviço
sc query WSICRMREST

REM Ou usando PowerShell
Get-Service WSICRMREST
```

**Saída esperada:**
```
SERVICE_NAME: WSICRMREST
DISPLAY_NAME: WSICRMREST API Service
TYPE               : 10  WIN32_OWN_PROCESS
STATE              : 4  RUNNING (ou STOPPED se não iniciou)
```

### Passo 4: Iniciar Serviço

**Opção 1: Via Services.msc**
1. Pressione `Win + R`
2. Digite: `services.msc`
3. Procure por "WSICRMREST API Service"
4. Clique com botão direito → Iniciar

**Opção 2: Via Linha de Comando (Administrador)**
```batch
sc start WSICRMREST
```

**Opção 3: Via PowerShell (Administrador)**
```powershell
Start-Service WSICRMREST
```

### Passo 5: Testar API

```batch
REM Teste básico
curl http://localhost:8080/connect/v1/wsteste

REM Ou abra no navegador
start http://localhost:8080/connect/v1/wsteste
```

---

## 🎛️ Gerenciamento do Serviço

### Script de Gerenciamento Interativo

Execute o script de gerenciamento:

```batch
scripts\manage_service_windows.bat
```

**Menu disponível:**
```
1 - Iniciar serviço
2 - Parar serviço
3 - Reiniciar serviço
4 - Ver status detalhado
5 - Ver logs (últimas 50 linhas)
6 - Abrir pasta de logs
7 - Testar API
0 - Sair
```

### Comandos Manuais

#### Iniciar Serviço
```batch
REM Linha de comando
sc start WSICRMREST

REM PowerShell
Start-Service WSICRMREST
```

#### Parar Serviço
```batch
REM Linha de comando
sc stop WSICRMREST

REM PowerShell
Stop-Service WSICRMREST
```

#### Reiniciar Serviço
```batch
REM Linha de comando
sc stop WSICRMREST
timeout /t 3 /nobreak
sc start WSICRMREST

REM PowerShell
Restart-Service WSICRMREST
```

#### Ver Status
```batch
REM Linha de comando
sc query WSICRMREST

REM PowerShell
Get-Service WSICRMREST

REM Status detalhado
sc qc WSICRMREST
```

#### Ver Logs
```batch
REM Abrir pasta de logs
explorer log\

REM Ver log mais recente (PowerShell)
Get-Content log\wsicrmrest_*.log -Tail 50
```

---

## 🗑️ Desinstalação

### Opção 1: Script Automático (Recomendado)

```batch
REM Clique com botão direito e "Executar como administrador"
scripts\uninstall_service_windows.bat
```

O script irá:
1. Verificar se o serviço existe
2. Parar o serviço (se estiver rodando)
3. Remover o serviço do sistema
4. Confirmar remoção

### Opção 2: Manual

```batch
REM 1. Parar o serviço
sc stop WSICRMREST

REM 2. Aguardar alguns segundos
timeout /t 5 /nobreak

REM 3. Remover o serviço
sc delete WSICRMREST
```

**Nota:** Os arquivos do projeto e logs **não são removidos** automaticamente. Para remover completamente, exclua a pasta do projeto manualmente.

---

## 🔍 Troubleshooting

### Problema: "Este script precisa ser executado como Administrador"

**Solução:**
1. Clique com botão direito no arquivo `.bat`
2. Selecione "Executar como administrador"
3. Aceite o UAC (Controle de Conta de Usuário)

---

### Problema: "Executável não encontrado"

**Erro:**
```
ERRO: Executável não encontrado: C:\path\build\wsicrmrest_win64.exe
```

**Solução:**
```batch
REM Compilar o projeto
cd scripts
build_windows.bat
REM Selecione opção 2 ou 3
```

---

### Problema: "Arquivo dbinit.ini não encontrado"

**Solução:**
```batch
REM Copiar arquivo de exemplo
copy dbinit.ini.example dbinit.ini

REM Editar configurações
notepad dbinit.ini
```

---

### Problema: Serviço não inicia (Erro 1053)

**Erro:**
```
O serviço não respondeu ao pedido de início ou controle em tempo hábil.
```

**Causas comuns:**
1. Erro no `dbinit.ini`
2. Banco de dados Oracle inacessível
3. Porta 8080 já em uso

**Solução:**
```batch
REM 1. Verificar logs
notepad log\wsicrmrest_2025-11-17.log

REM 2. Testar executável manualmente
cd build
wsicrmrest_win64.exe
REM Se houver erro, será exibido no console

REM 3. Verificar porta
netstat -ano | findstr :8080
```

---

### ⚠️ Aviso: Erro "Falha na ativação do aplicativo Microsoft.Windows.Cortana"

**Status: FALSO POSITIVO - Pode ser Ignorado**

Ao iniciar o serviço, você pode ver este erro no **Event Viewer**:

```
Falha na ativação do aplicativo Microsoft.Windows.Cortana_cw5n1h2txyewy!CortanaUI
com o erro: Este aplicativo não pode ser ativado pelo Administrador Interno.
```

**✅ Este erro NÃO afeta o funcionamento do WSICRMREST!**

**Por que acontece:**
- O serviço roda como "Sistema Local" (padrão do Windows)
- Windows tenta ativar componentes do sistema como Cortana
- Cortana não pode ser ativada pelo Administrador Interno
- Este é um comportamento normal do Windows

**Como verificar se o serviço está OK:**
```batch
REM 1. Testar a API
curl http://localhost:8080/connect/v1/wsteste

REM 2. Verificar status do serviço
sc query WSICRMREST

REM 3. Ver logs do WSICRMREST (não do Event Viewer)
notepad log\wsicrmrest_2025-11-17.log
```

**Se a API responder corretamente, ignore o erro da Cortana completamente.**

Para mais detalhes, consulte: `docs/setup/TROUBLESHOOTING_WINDOWS.md`

---

### Problema: "Access is denied" ao instalar

**Solução:**
1. Certifique-se de estar executando como **Administrador**
2. Desative temporariamente o antivírus (pode estar bloqueando)
3. Verifique se há outro serviço com o mesmo nome

---

### Problema: Serviço para sozinho após alguns segundos

**Causas:**
- Erro de conexão com banco de dados
- Configuração incorreta no `dbinit.ini`
- Tabela ORGANIZADOR vazia ou inexistente

**Solução:**
```batch
REM 1. Ver logs imediatamente após tentar iniciar
sc start WSICRMREST
timeout /t 2 /nobreak
notepad log\wsicrmrest_2025-11-17.log

REM 2. Procurar por erros como:
REM    - "Erro ao conectar ao banco de dados"
REM    - "Organizador Não Cadastrado"
REM    - "Erro ao carregar dados do organizador"
```

---

### Problema: Porta 8080 já está em uso

**Erro nos logs:**
```
bind: address already in use
```

**Solução 1: Mudar porta no dbinit.ini**
```ini
[application]
port = 8081  # Trocar para outra porta
```

**Solução 2: Identificar processo usando porta 8080**
```batch
REM Ver qual processo está usando a porta
netstat -ano | findstr :8080

REM Matar processo (substitua PID pelo número mostrado)
taskkill /PID 1234 /F
```

---

## 📊 Configurações do Serviço

### Propriedades Padrão

| Propriedade | Valor |
|------------|-------|
| **Nome do Serviço** | WSICRMREST |
| **Nome de Exibição** | WSICRMREST API Service |
| **Tipo de Início** | Automático (Atrasado) |
| **Conta** | Sistema Local |
| **Dependências** | Nenhuma |

### Configuração de Recuperação

O serviço é configurado para **reiniciar automaticamente** em caso de falha:

| Tentativa | Ação | Delay |
|-----------|------|-------|
| 1ª falha | Reiniciar serviço | 1 minuto |
| 2ª falha | Reiniciar serviço | 1 minuto |
| 3ª falha | Reiniciar serviço | 1 minuto |

**Resetar contador após:** 24 horas sem falhas

---

## 🔐 Permissões e Segurança

### Conta do Serviço

Por padrão, o serviço roda como **Sistema Local**, que tem:
- ✅ Acesso total ao sistema local
- ✅ Permissão para abrir portas de rede
- ✅ Acesso a arquivos locais
- ⚠️ Sem acesso a recursos de rede por padrão

### Firewall

Se precisar acessar a API de outras máquinas:

```batch
REM Adicionar regra de firewall (como Administrador)
netsh advfirewall firewall add rule ^
    name="WSICRMREST API" ^
    dir=in ^
    action=allow ^
    protocol=TCP ^
    localport=8080
```

---

## 📋 Checklist de Instalação

- [ ] Compilar executável Windows 64 bits
- [ ] Criar/configurar `dbinit.ini`
- [ ] Verificar conexão com Oracle (tnsping)
- [ ] Executar `install_service_windows.bat` como Administrador
- [ ] Verificar instalação (`sc query WSICRMREST`)
- [ ] Iniciar serviço
- [ ] Testar API (`http://localhost:8080/connect/v1/wsteste`)
- [ ] Verificar logs (`log\wsicrmrest_*.log`)
- [ ] Configurar firewall (se necessário)
- [ ] Documentar credenciais e configurações

---

## 🚀 Próximos Passos

Após instalar com sucesso:

1. **Monitoramento:** Configure alertas para quando o serviço parar
2. **Backup:** Faça backup regular do `dbinit.ini` e da pasta `log`
3. **Atualizações:** Para atualizar, pare o serviço, substitua o executável, inicie novamente
4. **Documentação:** Documente qualquer configuração customizada

---

## 📚 Referências

- **Scripts de Instalação:** `scripts/install_service_windows.bat`
- **Scripts de Desinstalação:** `scripts/uninstall_service_windows.bat`
- **Scripts de Gerenciamento:** `scripts/manage_service_windows.bat`
- **Configuração:** `dbinit.ini.example`
- **Logs:** `log/wsicrmrest_*.log`

---

## 💡 Dicas

1. **Sempre execute scripts de instalação/desinstalação como Administrador**
2. **Teste o executável manualmente antes de instalar como serviço**
3. **Mantenha backups do dbinit.ini**
4. **Configure rotação de logs para evitar enchimento de disco**
5. **Use o script `manage_service_windows.bat` para operações diárias**

---

**Documentação mantida por:** Equipe de Desenvolvimento
**Última revisão:** 2025-11-17
