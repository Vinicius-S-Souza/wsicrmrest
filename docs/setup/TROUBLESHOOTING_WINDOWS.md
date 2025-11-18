# Troubleshooting - Windows Service

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-17

## 🔍 Problemas Comuns e Soluções

---

## ⚠️ Erro: "Falha na ativação do aplicativo Microsoft.Windows.Cortana"

### Mensagem de Erro

```
Falha na ativação do aplicativo Microsoft.Windows.Cortana_cw5n1h2txyewy!CortanaUI
com o erro: Este aplicativo não pode ser ativado pelo Administrador Interno.
Veja o log Microsoft-Windows-TWinUI/Operational para obter informações adicionais.
```

### ✅ Status: FALSO POSITIVO - Pode ser Ignorado

**Este erro NÃO afeta o funcionamento do WSICRMREST!**

### Por que acontece?

Quando um serviço Windows é executado com a conta **Sistema Local** (padrão), o Windows tenta ativar alguns componentes do sistema, incluindo a Cortana. Como a Cortana não pode ser ativada pelo Administrador Interno, esse erro aparece nos logs do Event Viewer.

**Importante:**
- ✅ O serviço WSICRMREST **está funcionando corretamente**
- ✅ Este erro é do **Windows**, não do WSICRMREST
- ✅ Pode ser **completamente ignorado**
- ✅ Não afeta performance ou estabilidade

### Como Verificar se o Serviço Está Funcionando

#### Opção 1: Testar a API

```batch
REM Testar endpoint
curl http://localhost:8080/connect/v1/wsteste

REM Ou abrir no navegador
start http://localhost:8080/connect/v1/wsteste
```

**Se a API responder, o serviço está OK!**

#### Opção 2: Verificar Logs do WSICRMREST

```batch
REM Ver logs do serviço (não do Event Viewer)
notepad log\wsicrmrest_2025-11-17.log

REM Ou via PowerShell
Get-Content log\wsicrmrest_*.log -Tail 50
```

**Se houver logs recentes, o serviço está rodando!**

#### Opção 3: Verificar Status do Serviço

```batch
REM Status via sc
sc query WSICRMREST

REM Ou via PowerShell
Get-Service WSICRMREST

REM Ou via Services.msc
services.msc
REM Procure por "WSICRMREST API Service do Sistema ICRM"
```

**Se STATUS = RUNNING, está tudo OK!**

---

### Como Ocultar o Erro da Cortana (Opcional)

Se você quiser parar de ver esse erro nos logs do Windows:

#### Solução 1: Desabilitar Log da Cortana

```batch
REM Execute como Administrador (PowerShell)
wevtutil sl Microsoft-Windows-TWinUI/Operational /e:false
```

**Para reativar:**
```batch
wevtutil sl Microsoft-Windows-TWinUI/Operational /e:true
```

#### Solução 2: Criar Filtro no Event Viewer

1. Abra **Event Viewer** (eventvwr.msc)
2. Navegue até: **Applications and Services Logs → Microsoft → Windows → TWinUI → Operational**
3. Clique com botão direito → **Filter Current Log**
4. Em **Event sources**, desmarque **Cortana**
5. Clique **OK**

#### Solução 3: Executar com Conta Diferente (Avançado)

⚠️ **Não recomendado** - Requer mais configuração e pode causar problemas de permissões.

1. Criar conta de serviço dedicada
2. Dar permissões necessárias (rede, disco, Oracle)
3. Reconfigurar serviço para usar essa conta

```batch
REM Exemplo (NÃO RECOMENDADO para iniciantes)
sc config WSICRMREST obj= "DOMINIO\usuario_servico" password= "senha"
```

---

## 🔴 Erros Reais que Requerem Atenção

### Erro: Serviço Não Inicia (Erro 1053)

**Mensagem:**
```
O serviço não respondeu ao pedido de início ou controle em tempo hábil.
```

**Causas:**
1. Erro no `dbinit.ini`
2. Banco de dados Oracle inacessível
3. Porta já em uso
4. Tabela ORGANIZADOR vazia

**Solução:**

```batch
REM 1. Verificar logs do WSICRMREST (NÃO Event Viewer)
notepad log\wsicrmrest_2025-11-17.log

REM 2. Testar executável manualmente
cd build
wsicrmrest_win64.exe
REM Se houver erro, será exibido no console

REM 3. Verificar conexão Oracle
sqlplus usuario/senha@tns_name

REM 4. Verificar porta
netstat -ano | findstr :8080
```

---

### Erro: Access Denied

**Mensagem:**
```
Acesso negado / Access is denied
```

**Causa:** Script não foi executado como Administrador

**Solução:**
1. Clique com botão direito no `.bat`
2. Selecione **"Executar como administrador"**
3. Aceite o UAC (Controle de Conta de Usuário)

---

### Erro: Executável Não Encontrado

**Mensagem:**
```
ERRO: Executável não encontrado: C:\path\wsicrmrest_win64.exe
```

**Solução:**

```batch
REM Compilar o projeto
scripts\build_windows.bat
REM Selecione opção 2 (64 bits)

REM Verificar se foi criado
dir wsicrmrest_win64.exe
```

---

### Erro: dbinit.ini Não Encontrado

**Solução:**

```batch
REM Copiar exemplo
copy dbinit.ini.example dbinit.ini

REM Editar configurações
notepad dbinit.ini

REM Configurar minimamente:
REM [database]
REM tns_name = SEU_TNS
REM username = SEU_USUARIO
REM password = SUA_SENHA
```

---

### Erro: Porta 8080 em Uso

**Mensagem nos logs:**
```
bind: address already in use
```

**Verificar processo:**
```batch
netstat -ano | findstr :8080
REM Último número é o PID do processo
```

**Solução 1: Matar processo**
```batch
taskkill /PID 1234 /F
REM Substitua 1234 pelo PID encontrado
```

**Solução 2: Mudar porta**
```ini
[application]
port = 8081
```

---

### Erro: Conexão com Oracle Falhou

**Mensagens nos logs:**
```
Erro ao conectar ao banco de dados
ORA-12154: TNS:could not resolve the connect identifier
ORA-01017: invalid username/password
```

**Verificações:**

```batch
REM 1. Testar TNS
tnsping SEU_TNS_NAME

REM 2. Testar conexão
sqlplus usuario/senha@tns_name

REM 3. Verificar variável ORACLE_HOME
echo %ORACLE_HOME%

REM 4. Verificar tnsnames.ora
notepad %ORACLE_HOME%\network\admin\tnsnames.ora
```

---

### Erro: Organizador Não Cadastrado

**Mensagem nos logs:**
```
Organizador Não Cadastrado
Erro ao carregar dados do organizador
```

**Causa:** Tabela ORGANIZADOR está vazia ou não existe

**Solução:**

```sql
-- Verificar se tabela existe
SELECT COUNT(*) FROM ORGANIZADOR WHERE ORGCODIGO > 0;

-- Se COUNT = 0, inserir registro:
INSERT INTO ORGANIZADOR (
    ORGCODIGO,
    ORGNOME,
    ORGCNPJ,
    ORGCODLOJAMATRIZ,
    ORGCODISGA
) VALUES (
    1,
    'Minha Empresa',
    '12345678000190',
    1,
    123
);
COMMIT;
```

---

## 📋 Checklist de Diagnóstico

Quando o serviço não funciona, verifique **nesta ordem**:

- [ ] **1. O serviço está instalado?**
  ```batch
  sc query WSICRMREST
  ```

- [ ] **2. O serviço está rodando?**
  ```batch
  sc query WSICRMREST | findstr "RUNNING"
  ```

- [ ] **3. A API responde?**
  ```batch
  curl http://localhost:8080/connect/v1/wsteste
  ```

- [ ] **4. Há logs recentes?**
  ```batch
  dir log\wsicrmrest_*.log /O-D
  type log\wsicrmrest_2025-11-17.log
  ```

- [ ] **5. Há erros nos logs do WSICRMREST?**
  ```batch
  findstr /C:"ERROR" log\wsicrmrest_*.log
  ```

- [ ] **6. O dbinit.ini está correto?**
  ```batch
  type dbinit.ini
  ```

- [ ] **7. O Oracle está acessível?**
  ```batch
  tnsping SEU_TNS_NAME
  sqlplus usuario/senha@tns_name
  ```

- [ ] **8. A porta está livre?**
  ```batch
  netstat -ano | findstr :8080
  ```

---

## 🔧 Comandos Úteis de Diagnóstico

### Informações do Serviço

```batch
REM Status básico
sc query WSICRMREST

REM Configuração completa
sc qc WSICRMREST

REM Informações detalhadas (PowerShell)
Get-Service WSICRMREST | Format-List *

REM Dependências
sc enumdepend WSICRMREST
```

### Logs do Sistema Windows

```batch
REM Event Viewer
eventvwr.msc

REM Navegar até:
REM - Windows Logs → System
REM - Windows Logs → Application

REM Filtrar por "WSICRMREST"
```

### Verificar Processos

```batch
REM Listar processos do WSICRMREST
tasklist | findstr wsicrmrest

REM Detalhes do processo (PowerShell)
Get-Process wsicrmrest* | Format-List *

REM Ver portas abertas pelo processo
netstat -ano | findstr "PID_DO_PROCESSO"
```

---

## 💡 Dicas de Prevenção

1. **Sempre teste manualmente antes de instalar como serviço**
   ```batch
   wsicrmrest_win64.exe
   ```

2. **Mantenha backups do dbinit.ini**
   ```batch
   copy dbinit.ini dbinit.ini.backup
   ```

3. **Monitore logs regularmente**
   ```batch
   findstr /C:"ERROR" log\wsicrmrest_*.log
   ```

4. **Configure alertas para o serviço**
   - Use Task Scheduler para verificar se serviço está rodando
   - Envie alerta se parar

5. **Documente sua configuração**
   - Anote porta usada
   - Anote TNS name
   - Anote versão instalada

---

## 📞 Quando Pedir Ajuda

Se após verificar **todos os itens acima** ainda houver problemas:

### Informações para Fornecer

```batch
REM 1. Versão do serviço
type wsicrmrest_win64.exe | findstr "version"

REM 2. Status do serviço
sc query WSICRMREST

REM 3. Últimos 100 logs
powershell -Command "Get-Content log\wsicrmrest_*.log -Tail 100"

REM 4. Erros recentes
findstr /C:"ERROR" log\wsicrmrest_*.log | more

REM 5. Configuração (sem senha)
type dbinit.ini | findstr /V "password"

REM 6. Versão do Windows
ver
systeminfo | findstr /C:"OS"

REM 7. Versão do Oracle
sqlplus -v
```

---

## 🎯 Resumo: Erro da Cortana

| Pergunta | Resposta |
|----------|----------|
| **É um problema?** | ❌ Não |
| **Afeta o WSICRMREST?** | ❌ Não |
| **Preciso corrigir?** | ❌ Não |
| **Posso ignorar?** | ✅ Sim |
| **Como verificar se está OK?** | Testar a API: `curl http://localhost:8080/connect/v1/wsteste` |

---

**Documentação mantida por:** Equipe de Desenvolvimento
**Última revisão:** 2025-11-17
