@echo off
REM Script de execução para Windows
REM Data de criação: 2025-11-17

echo ========================================
echo   WSICRMREST - Servidor REST API
echo ========================================
echo.

REM Verificar se dbinit.ini existe
if not exist "dbinit.ini" (
    echo ERRO: Arquivo dbinit.ini não encontrado!
    echo.
    echo Por favor, copie dbinit.ini.example para dbinit.ini
    echo e configure suas credenciais do banco de dados.
    echo.
    pause
    exit /b 1
)

REM Detectar arquitetura do Windows
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set BINARY=build\wsicrmrest_win64.exe
    set ARCH=64 bits
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set BINARY=build\wsicrmrest_win32.exe
    set ARCH=32 bits
) else (
    echo Arquitetura não suportada: %PROCESSOR_ARCHITECTURE%
    pause
    exit /b 1
)

REM Verificar se o binário existe
if not exist "%BINARY%" (
    echo ERRO: Binário não encontrado: %BINARY%
    echo.
    echo Execute primeiro: scripts\build_windows.bat
    echo.
    pause
    exit /b 1
)

REM Criar diretório de logs se não existir
if not exist "log" mkdir log

echo Iniciando servidor (%ARCH%)...
echo Binário: %BINARY%
echo.
echo Pressione Ctrl+C para parar o servidor
echo.
echo ========================================
echo.

REM Executar o servidor
"%BINARY%"

REM Se o servidor parar, mostrar mensagem
echo.
echo Servidor encerrado.
pause
