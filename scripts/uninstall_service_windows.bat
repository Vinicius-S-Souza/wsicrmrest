@echo off
REM Script de desinstalação do WSICRMREST do Windows Services
REM Data de criação: 2025-11-17
REM
REM IMPORTANTE: Execute este script como Administrador

echo ========================================
echo   WSICRMREST - Desinstalar Serviço
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

set SERVICE_NAME=WSICRMREST

REM Verificar se o serviço existe
sc query "%SERVICE_NAME%" >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo O serviço %SERVICE_NAME% não está instalado.
    echo.
    pause
    exit /b 0
)

echo Serviço encontrado: %SERVICE_NAME%
echo.

REM Obter informações do serviço
for /f "tokens=3" %%i in ('sc query "%SERVICE_NAME%" ^| findstr "STATE"') do set SERVICE_STATE=%%i

REM Obter caminho do executável
for /f "tokens=2*" %%i in ('sc qc "%SERVICE_NAME%" ^| findstr "BINARY_PATH_NAME"') do set BINARY_PATH_NAME=%%j

REM Remover aspas do path
set BINARY_PATH_NAME=%BINARY_PATH_NAME:"=%

REM Detectar arquitetura do executável instalado
set INSTALLED_ARCH=desconhecida
if not "%BINARY_PATH_NAME%"=="" (
    echo %BINARY_PATH_NAME% | findstr /i "win32.exe" >nul
    if %ERRORLEVEL%==0 set INSTALLED_ARCH=32 bits

    echo %BINARY_PATH_NAME% | findstr /i "win64.exe" >nul
    if %ERRORLEVEL%==0 set INSTALLED_ARCH=64 bits
)

echo Status atual: %SERVICE_STATE%
echo Executável: %BINARY_PATH_NAME%
echo Arquitetura: %INSTALLED_ARCH%
echo.

set /p CONFIRM="Tem certeza que deseja remover o serviço? (S/N): "
if /i "%CONFIRM%" NEQ "S" (
    echo Desinstalação cancelada.
    pause
    exit /b 0
)

REM Parar o serviço se estiver rodando
if /i "%SERVICE_STATE%" NEQ "STOPPED" (
    echo.
    echo Parando o serviço...
    sc stop "%SERVICE_NAME%"

    if %ERRORLEVEL% EQU 0 (
        echo Serviço parado com sucesso.
        echo Aguardando finalização...
        timeout /t 5 /nobreak >nul
    ) else (
        echo Aviso: Não foi possível parar o serviço gracefully.
        echo Tentando forçar parada...
        taskkill /F /FI "SERVICES eq %SERVICE_NAME%" >nul 2>&1
        timeout /t 2 /nobreak >nul
    )
)

echo.
echo Removendo o serviço...
sc delete "%SERVICE_NAME%"

if %ERRORLEVEL% EQU 0 (
    echo Serviço removido com sucesso!

    REM Remover Event Log Source (silencioso)
    echo Removendo registros do Event Log...
    reg delete "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /f >nul 2>&1

    echo.
    echo ========================================
    echo   Serviço removido com sucesso!
    echo ========================================
    echo.
    echo O serviço %SERVICE_NAME% foi completamente removido.
    echo.
    echo Nota: Os arquivos do aplicativo e logs não foram removidos.
    echo Para remover completamente:
    echo   1. Exclua a pasta do projeto manualmente
    echo   2. Ou mantenha para reinstalar depois
) else (
    echo.
    echo ERRO ao remover o serviço!
    echo.
    echo Possíveis causas:
    echo   - O serviço ainda está em execução
    echo   - Permissões insuficientes
    echo   - O serviço está sendo usado por outro processo
    echo.
    echo Tente:
    echo   1. Reiniciar o computador
    echo   2. Executar este script novamente
)

echo.
pause
