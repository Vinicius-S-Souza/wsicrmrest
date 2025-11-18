# Scripts do WSICRMREST

**Data de criação:** 2025-11-17

## 📁 Conteúdo da Pasta

Esta pasta contém scripts úteis para compilação, execução e gerenciamento do WSICRMREST.

---

## 🪟 Scripts Windows

### 🔨 Compilação

#### `build_windows.bat`
Compila o projeto para Windows (32 e/ou 64 bits).

**Como usar:**
```batch
build_windows.bat
```

**Menu interativo:**
- Opção 1: Windows 32 bits → `build\wsicrmrest_win32.exe`
- Opção 2: Windows 64 bits → `build\wsicrmrest_win64.exe`
- Opção 3: Ambas as versões

**Pré-requisitos:**
- Go instalado
- MinGW/CGO configurado (para cross-compilation no Linux)

---

### ▶️ Execução

#### `run_windows.bat`
Executa o WSICRMREST no Windows.

**Como usar:**
```batch
run_windows.bat
```

**O que faz:**
- Detecta automaticamente arquitetura (32 ou 64 bits)
- Verifica se `dbinit.ini` existe
- Executa o binário correspondente

---

### 🔧 Gerenciamento de Serviço Windows

#### `install_service_windows.bat` ⭐
Instala o WSICRMREST como serviço do Windows.

**Como usar:**
```batch
REM Clique com botão direito → Executar como administrador
install_service_windows.bat
```

**O que faz:**
- ✅ Verifica permissões de administrador
- ✅ Verifica se executável e dbinit.ini existem
- ✅ Cria serviço Windows
- ✅ Configura início automático (atrasado)
- ✅ Configura recuperação automática
- ✅ Oferece opção de iniciar imediatamente

**Após instalação:**
- Nome do serviço: `WSICRMREST`
- Gerenciável via `services.msc`
- Inicia automaticamente no boot

---

#### `uninstall_service_windows.bat`
Remove o WSICRMREST dos serviços do Windows.

**Como usar:**
```batch
REM Clique com botão direito → Executar como administrador
uninstall_service_windows.bat
```

**O que faz:**
- ✅ Verifica se serviço existe
- ✅ Para o serviço (se estiver rodando)
- ✅ Remove o serviço do sistema
- ⚠️ Não remove arquivos do projeto (apenas o serviço)

---

#### `manage_service_windows.bat` ⭐
Menu interativo para gerenciar o serviço.

**Como usar:**
```batch
manage_service_windows.bat
```

**Funcionalidades:**
- 1️⃣ Iniciar serviço
- 2️⃣ Parar serviço
- 3️⃣ Reiniciar serviço
- 4️⃣ Ver status detalhado
- 5️⃣ Ver logs (últimas 50 linhas)
- 6️⃣ Abrir pasta de logs
- 7️⃣ Testar API
- 0️⃣ Sair

**Não requer privilégios de administrador** (exceto para iniciar/parar serviço)

---

## 🐧 Scripts Linux

### `test_apis.sh`
Testa os endpoints da API.

**Como usar:**
```bash
chmod +x scripts/test_apis.sh
./scripts/test_apis.sh
```

**Pré-requisitos:**
- Servidor deve estar rodando
- `curl` ou `httpie` instalado

---

## 📋 Workflow Típico

### Desenvolvimento Local (Windows)

```batch
REM 1. Compilar
scripts\build_windows.bat

REM 2. Configurar (primeira vez)
copy dbinit.ini.example dbinit.ini
notepad dbinit.ini

REM 3. Testar manualmente
scripts\run_windows.bat

REM 4. Testar API
curl http://localhost:8080/connect/v1/wsteste
```

### Instalação como Serviço (Windows)

```batch
REM 1. Compilar
scripts\build_windows.bat
REM Selecione opção 2 (64 bits)

REM 2. Instalar serviço (como Administrador)
scripts\install_service_windows.bat

REM 3. Gerenciar serviço
scripts\manage_service_windows.bat
```

### Atualização do Serviço (Windows)

```batch
REM 1. Parar serviço
sc stop WSICRMREST

REM 2. Aguardar
timeout /t 3 /nobreak

REM 3. Compilar nova versão
scripts\build_windows.bat

REM 4. Iniciar serviço
sc start WSICRMREST

REM 5. Verificar logs
scripts\manage_service_windows.bat
REM Opção 5 ou 6
```

---

## 🔐 Permissões

### Windows

**Scripts de serviço requerem privilégios de Administrador:**
- `install_service_windows.bat` ✅ Administrador obrigatório
- `uninstall_service_windows.bat` ✅ Administrador obrigatório

**Scripts normais:**
- `build_windows.bat` ✅ Usuário normal
- `run_windows.bat` ✅ Usuário normal
- `manage_service_windows.bat` ⚠️ Algumas opções requerem Administrador

### Linux

```bash
# Dar permissão de execução
chmod +x scripts/*.sh
```

---

## 📚 Documentação Relacionada

- **Instalação Windows Service:** `docs/setup/INSTALACAO_WINDOWS_SERVICE.md`
- **Build Windows:** `docs/setup/BUILD_WINDOWS.md`
- **Configuração:** `dbinit.ini.example`

---

## 💡 Dicas

1. **Sempre compile antes de instalar como serviço**
2. **Use o script `manage_service_windows.bat` para operações diárias**
3. **Mantenha backups do `dbinit.ini` antes de atualizar**
4. **Verifique logs após instalar/atualizar o serviço**
5. **Teste manualmente com `run_windows.bat` antes de instalar como serviço**

---

## ⚠️ Problemas Comuns

### "Este script precisa ser executado como Administrador"

**Solução:**
```batch
REM Clique com botão direito no .bat → Executar como administrador
```

### "Go não encontrado no PATH"

**Solução:**
```batch
REM Instale Go: https://golang.org/dl/
REM Reinicie o terminal após instalação
```

### "Executável não encontrado"

**Solução:**
```batch
REM Compile primeiro
scripts\build_windows.bat
```

---

**Última atualização:** 2025-11-17
