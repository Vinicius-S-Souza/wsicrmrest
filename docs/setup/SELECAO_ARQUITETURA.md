# Seleção de Arquitetura - 32 vs 64 bits

**Data de criação:** 2025-11-24
**Última atualização:** 2025-11-24

Este documento explica como funciona a detecção e seleção de arquitetura (32/64 bits) nos scripts de instalação do WSICRMREST.

---

## 🎯 Detecção Automática

O script `install_service_windows.bat` detecta automaticamente a arquitetura do Windows usando variáveis de ambiente:

### Como funciona?

```batch
if defined PROCESSOR_ARCHITEW6432 (
    REM Windows 64 bits rodando script 32 bits
    set ARCH=64
) else if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    REM Windows 64 bits nativo
    set ARCH=64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    REM Windows 32 bits
    set ARCH=32
) else (
    REM Fallback para 64 bits
    set ARCH=64
)
```

### Cenários:

| Sistema Windows | Detecção | Executável Recomendado |
|----------------|----------|------------------------|
| Windows 11 (64 bits) | ✅ 64 bits | wsicrmrest_win64.exe |
| Windows 10 (64 bits) | ✅ 64 bits | wsicrmrest_win64.exe |
| Windows Server 2016+ | ✅ 64 bits | wsicrmrest_win64.exe |
| Windows 10 (32 bits) | ✅ 32 bits | wsicrmrest_win32.exe |
| Windows 7/8 (32 bits) | ✅ 32 bits | wsicrmrest_win32.exe |

---

## 📋 Menu de Seleção

Quando você executa o instalador, ele mostra:

```
Detectando arquitetura do Windows...
Sistema detectado: Windows 64 bits

Executáveis disponíveis:
  [1] wsicrmrest_win32.exe (32 bits)
  [2] wsicrmrest_win64.exe (64 bits)
  [A] Detectar automaticamente (recomendado)

Escolha o executável [1/2/A - padrão A]:
```

### Opções:

- **[A] Automático (Padrão):** Usa detecção automática
  - Windows 64 bits → seleciona win64.exe
  - Windows 32 bits → seleciona win32.exe
  - Se o executável ideal não existir, usa o disponível com aviso

- **[1] Manual 32 bits:** Força uso do executável 32 bits
  - Útil em sistemas 64 bits por razões de compatibilidade
  - Funciona, mas pode ter desempenho reduzido

- **[2] Manual 64 bits:** Força uso do executável 64 bits
  - Recomendado para sistemas 64 bits
  - **Não funciona em sistemas 32 bits!**

---

## ⚠️ Avisos e Validações

### Aviso 1: Sistema 32 bits + Executável 64 bits

```
AVISO: Sistema 32 bits, mas apenas executável 64 bits disponível.
Usando wsicrmrest_win64.exe (pode não funcionar em Windows 32 bits)
```

**O que fazer:**
- Compilar versão 32 bits:
  ```bash
  make build-windows-32
  ```
- Ou aceitar que pode não funcionar

### Aviso 2: Sistema 64 bits + Executável 32 bits

```
AVISO: Usando executável 32 bits em sistema 64 bits.
Recomenda-se compilar versão 64 bits para melhor desempenho.
```

**O que fazer:**
- Compilar versão 64 bits:
  ```bash
  make build-windows-64
  ```
- Ou aceitar desempenho reduzido (funciona normalmente)

### Erro: Nenhum executável disponível

```
ERRO: Nenhum executável encontrado!

Por favor, compile o projeto primeiro:
  1. Execute: scripts\build_windows.bat
  2. Ou execute: make build-windows
```

**Solução:**
Compilar pelo menos uma das versões.

---

## 🔧 Compilação

### Compilar Ambas Versões (Recomendado)

**Linux/WSL:**
```bash
make build-windows
```

Ou separadamente:
```bash
make build-windows-32
make build-windows-64
```

**Windows:**
```batch
scripts\build_windows.bat
```

Menu interativo que compila ambas versões.

### Compilação Manual

**32 bits:**
```bash
GOOS=windows GOARCH=386 go build -o build/wsicrmrest_win32.exe ./cmd/server
```

**64 bits:**
```bash
GOOS=windows GOARCH=amd64 go build -o build/wsicrmrest_win64.exe ./cmd/server
```

---

## 📊 Comparação de Desempenho

| Característica | 32 bits | 64 bits |
|----------------|---------|---------|
| **Memória Máxima** | ~3.5 GB | Praticamente ilimitada |
| **Desempenho** | Normal | +10-30% mais rápido |
| **Compatibilidade** | Windows 32/64 bits | Apenas Windows 64 bits |
| **Tamanho do Executável** | Menor (~8 MB) | Maior (~10 MB) |
| **Uso de Memória** | Menor | Ligeiramente maior |

### Quando usar 32 bits?

- ✅ Sistema Windows 32 bits (obrigatório)
- ✅ Servidor muito antigo com recursos limitados
- ✅ Compatibilidade com sistemas legados
- ⚠️ Aplicação usa < 3 GB de memória

### Quando usar 64 bits?

- ✅ Sistema Windows 64 bits moderno (recomendado)
- ✅ Aplicação pode usar > 3 GB de memória
- ✅ Melhor desempenho é importante
- ✅ Servidor em produção

**Recomendação geral:** Use 64 bits se o sistema suportar.

---

## 🔍 Verificar Arquitetura Instalada

### Via Script de Gerenciamento

```batch
scripts\manage_service_windows.bat
```

Mostra:
```
Status do Serviço: RUNNING
Arquitetura: 64 bits
```

### Via Script de Desinstalação

```batch
scripts\uninstall_service_windows.bat
```

Mostra:
```
Serviço encontrado: WSICRMREST
Status atual: RUNNING
Executável: C:\CRM\WSICRMREST\wsicrmrest_win64.exe
Arquitetura: 64 bits
```

### Manualmente (sc qc)

```batch
sc qc WSICRMREST | findstr BINARY_PATH_NAME
```

Saída:
```
BINARY_PATH_NAME   : "C:\CRM\WSICRMREST\wsicrmrest_win64.exe"
```

---

## 🔄 Trocar de Arquitetura

Se você instalou a versão errada e quer trocar:

### Passo 1: Desinstalar serviço atual

```batch
scripts\uninstall_service_windows.bat
```

### Passo 2: Compilar versão desejada

```bash
# Para 64 bits
make build-windows-64

# Para 32 bits
make build-windows-32
```

### Passo 3: Copiar executável

```batch
copy /Y build\wsicrmrest_win64.exe C:\CRM\WSICRMREST\wsicrmrest_win64.exe
```

### Passo 4: Reinstalar com nova versão

```batch
scripts\install_service_windows.bat
```

Escolha a arquitetura desejada no menu.

---

## 🐛 Solução de Problemas

### Problema: "O sistema não pode executar o programa especificado"

**Causa:** Executável 64 bits em sistema 32 bits.

**Solução:**
1. Desinstalar serviço
2. Compilar versão 32 bits
3. Reinstalar com versão correta

### Problema: Serviço usa muita memória

**Causa:** Versão 64 bits usa mais memória que 32 bits.

**Solução:**
- Normal, versão 64 bits usa ~20-30% mais memória
- Se servidor tem < 2 GB RAM, considere versão 32 bits
- Ou aumente RAM do servidor

### Problema: Desempenho ruim

**Causa:** Versão 32 bits em sistema 64 bits potente.

**Solução:**
- Trocar para versão 64 bits
- Ganho de 10-30% em desempenho

---

## 📚 Referências Técnicas

### Variáveis de Ambiente Windows

| Variável | Significado |
|----------|-------------|
| `PROCESSOR_ARCHITECTURE` | Arquitetura do processo atual |
| `PROCESSOR_ARCHITEW6432` | Arquitetura real (se diferente) |

**Valores:**
- `AMD64` = Windows 64 bits
- `x86` = Windows 32 bits (ou processo 32 bits em 64 bits)

### Detalhes de Compilação Go

**GOARCH valores:**
- `386` = 32 bits (Intel 80386+)
- `amd64` = 64 bits (AMD64/Intel 64/x86-64)

**Cross-compilation:**
```bash
GOOS=windows GOARCH=386   # Windows 32 bits
GOOS=windows GOARCH=amd64 # Windows 64 bits
```

---

## ✅ Checklist de Instalação

Antes de instalar, verifique:

- [ ] Sistema Windows é 32 ou 64 bits?
  ```batch
  systeminfo | findstr /C:"Tipo de Sistema"
  ```

- [ ] Executável correspondente compilado?
  ```batch
  dir build\*.exe
  ```

- [ ] Espaço em disco suficiente? (mínimo 100 MB)
  ```batch
  dir C:\CRM\WSICRMREST
  ```

- [ ] Permissões de administrador?
  ```batch
  net session
  ```

Se todos OK, prossiga com instalação!

---

**Última atualização:** 2025-11-24
