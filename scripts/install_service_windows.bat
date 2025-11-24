@echo off
REM Script de instalação do WSICRMREST como serviço Windows
REM Data de criação: 2025-11-17
REM
REM IMPORTANTE: Execute este script como Administrador

echo ========================================
echo   WSICRMREST - Instalação do Serviço
echo ========================================
echo.

REM Verificar se está rodando como Administrador
net session >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERRO: Este script precisa ser executado como Administrador!
    echo.
    echo Clique com botão direito no arquivo e selecione "Executar como administrador"
    echo.
    pause
    exit /b 1
)

REM Configurações do serviço
set SERVICE_NAME=WSICRMREST
set SERVICE_DISPLAY_NAME=WSICRMREST API Service do Sistema ICRM
set SERVICE_DESCRIPTION=Web Service REST API para Integração com Sistema ICRM

REM Converter para path absoluto
pushd %~dp0
set WORK_DIR=%CD%
popd

REM ============================================
REM Detectar arquitetura do Windows
REM ============================================
echo Detectando arquitetura do Windows...
echo.

REM Verificar se é 64 bits
if defined PROCESSOR_ARCHITEW6432 (
    set ARCH=64
    set ARCH_NAME=64 bits
) else if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=64
    set ARCH_NAME=64 bits
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set ARCH=32
    set ARCH_NAME=32 bits
) else (
    set ARCH=64
    set ARCH_NAME=64 bits (detectado por padrão)
)

echo Sistema detectado: Windows %ARCH_NAME%
echo.

REM ============================================
REM Verificar executáveis disponíveis
REM ============================================
set EXE_32=%WORK_DIR%\wsicrmrest_win32.exe
set EXE_64=%WORK_DIR%\wsicrmrest_win64.exe
set HAS_32=0
set HAS_64=0

if exist "%EXE_32%" set HAS_32=1
if exist "%EXE_64%" set HAS_64=1

REM ============================================
REM Selecionar executável
REM ============================================
echo Executáveis disponíveis:
if %HAS_32%==1 echo   [1] wsicrmrest_win32.exe (32 bits)
if %HAS_64%==1 echo   [2] wsicrmrest_win64.exe (64 bits)
echo   [A] Detectar automaticamente (recomendado)
echo.

REM Se nenhum executável existe
if %HAS_32%==0 if %HAS_64%==0 (
    echo ERRO: Nenhum executável encontrado!
    echo.
    echo Por favor, compile o projeto primeiro:
    echo   1. Execute: scripts\build_windows.bat
    echo   2. Ou execute: make build-windows
    echo.
    pause
    exit /b 1
)

REM Perguntar ao usuário
set /p CHOICE="Escolha o executável [1/2/A - padrão A]: "

REM Se não escolheu nada, usar auto-detect
if "%CHOICE%"=="" set CHOICE=A

REM Processar escolha
if /i "%CHOICE%"=="A" (
    REM Auto-detect baseado na arquitetura
    if %ARCH%==32 (
        if %HAS_32%==1 (
            set BINARY_PATH=%EXE_32%
            set SELECTED_ARCH=32 bits (auto)
        ) else if %HAS_64%==1 (
            echo.
            echo AVISO: Sistema 32 bits, mas apenas executável 64 bits disponível.
            echo Usando wsicrmrest_win64.exe (pode não funcionar em Windows 32 bits)
            set BINARY_PATH=%EXE_64%
            set SELECTED_ARCH=64 bits (fallback)
        )
    ) else (
        if %HAS_64%==1 (
            set BINARY_PATH=%EXE_64%
            set SELECTED_ARCH=64 bits (auto)
        ) else if %HAS_32%==1 (
            echo.
            echo AVISO: Usando executável 32 bits em sistema 64 bits.
            echo Recomenda-se compilar versão 64 bits para melhor desempenho.
            set BINARY_PATH=%EXE_32%
            set SELECTED_ARCH=32 bits (fallback)
        )
    )
) else if "%CHOICE%"=="1" (
    if %HAS_32%==1 (
        set BINARY_PATH=%EXE_32%
        set SELECTED_ARCH=32 bits (manual)
    ) else (
        echo ERRO: wsicrmrest_win32.exe não encontrado!
        pause
        exit /b 1
    )
) else if "%CHOICE%"=="2" (
    if %HAS_64%==1 (
        set BINARY_PATH=%EXE_64%
        set SELECTED_ARCH=64 bits (manual)
    ) else (
        echo ERRO: wsicrmrest_win64.exe não encontrado!
        pause
        exit /b 1
    )
) else (
    echo ERRO: Escolha inválida!
    pause
    exit /b 1
)

echo.
echo ============================================
echo Configurações do Serviço:
echo ============================================
echo   Nome do Serviço: %SERVICE_NAME%
echo   Nome de Exibição: %SERVICE_DISPLAY_NAME%
echo   Executável: %BINARY_PATH%
echo   Arquitetura: %SELECTED_ARCH%
echo   Diretório de Trabalho: %WORK_DIR%
echo.

REM Verificar se o executável existe (dupla verificação)
if not exist "%BINARY_PATH%" (
    echo ERRO: Executável não encontrado: %BINARY_PATH%
    echo.
    pause
    exit /b 1
)

REM Verificar se dbinit.ini existe
if not exist "%WORK_DIR%\dbinit.ini" (
    echo ERRO: Arquivo dbinit.ini não encontrado!
    echo.
    echo Copie dbinit.ini.example para dbinit.ini e configure suas credenciais.
    echo.
    pause
    exit /b 1
)

REM Verificar se o serviço já existe
sc query "%SERVICE_NAME%" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo.
    echo O serviço %SERVICE_NAME% já está instalado.
    echo.
    set /p REINSTALL="Deseja reinstalar? (S/N): "
    if /i "%REINSTALL%" NEQ "S" (
        echo Instalação cancelada.
        pause
        exit /b 0
    )

    echo.
    echo Parando serviço existente...
    sc stop "%SERVICE_NAME%" >nul 2>&1
    timeout /t 3 /nobreak >nul

    echo Removendo serviço existente...
    sc delete "%SERVICE_NAME%"
    if %ERRORLEVEL% NEQ 0 (
        echo ERRO ao remover serviço existente!
        pause
        exit /b 1
    )

    echo Aguardando remoção...
    timeout /t 2 /nobreak >nul
)

echo.
echo Instalando serviço...

REM Registrar Event Log Source (silencioso, ignora se já existir)
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /v EventMessageFile /t REG_EXPAND_SZ /d "%SystemRoot%\System32\EventCreate.exe" /f >nul 2>&1
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /v TypesSupported /t REG_DWORD /d 7 /f >nul 2>&1

REM Criar serviço usando sc create
sc create "%SERVICE_NAME%" ^
    binPath= "\"%BINARY_PATH%\"" ^
    start= auto ^
    DisplayName= "%SERVICE_DISPLAY_NAME%"

if %ERRORLEVEL% NEQ 0 (
    echo ERRO ao criar o serviço!
    pause
    exit /b 1
)

echo Serviço criado com sucesso!

REM Configurar descrição
sc description "%SERVICE_NAME%" "%SERVICE_DESCRIPTION%"

REM Configurar recuperação em caso de falha
echo Configurando recuperação automática...
sc failure "%SERVICE_NAME%" reset= 86400 actions= restart/60000/restart/60000/restart/60000

REM Configurar delay para início automático (aguardar outros serviços)
sc config "%SERVICE_NAME%" start= delayed-auto

echo.
echo ========================================
echo   Instalação Concluída!
echo ========================================
echo.
echo O serviço foi instalado com as seguintes configurações:
echo   - Início: Automático (Atrasado)
echo   - Recuperação: Reiniciar automaticamente após falha
echo   - Diretório: %WORK_DIR%
echo.

set /p START_NOW="Deseja iniciar o serviço agora? (S/N): "
if /i "%START_NOW%" EQU "S" (
    echo.
    echo Iniciando serviço...
    sc start "%SERVICE_NAME%"

    if %ERRORLEVEL% EQU 0 (
        echo.
        echo Serviço iniciado com sucesso!
        echo.
        echo Verifique os logs em: %WORK_DIR%\log\
    ) else (
        echo.
        echo ERRO ao iniciar o serviço!
        echo Verifique:
        echo   1. Se o arquivo dbinit.ini está configurado corretamente
        echo   2. Se a conexão com o banco de dados está funcionando
        echo   3. Os logs em: %WORK_DIR%\log\
    )
) else (
    echo.
    echo Para iniciar o serviço manualmente:
    echo   - Via Services.msc: Procure por "%SERVICE_DISPLAY_NAME%"
    echo   - Via linha de comando: sc start %SERVICE_NAME%
    echo   - Via PowerShell: Start-Service %SERVICE_NAME%
)

echo.
echo Comandos úteis:
echo   Iniciar:    sc start %SERVICE_NAME%
echo   Parar:      sc stop %SERVICE_NAME%
echo   Status:     sc query %SERVICE_NAME%
echo   Remover:    sc delete %SERVICE_NAME% (após parar)
echo.
pause
