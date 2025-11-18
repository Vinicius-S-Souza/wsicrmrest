@echo off
REM Script de gerenciamento do serviço WSICRMREST
REM Data de criação: 2025-11-17

echo ========================================
echo   WSICRMREST - Gerenciar Serviço
echo ========================================
echo.

set SERVICE_NAME=WSICRMREST

REM Verificar se o serviço existe
sc query "%SERVICE_NAME%" >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo O serviço %SERVICE_NAME% não está instalado.
    echo.
    echo Para instalar, execute: install_service_windows.bat
    echo.
    pause
    exit /b 0
)

:MENU
cls
echo ========================================
echo   WSICRMREST - Menu de Gerenciamento
echo ========================================
echo.

REM Obter status do serviço
for /f "tokens=3" %%i in ('sc query "%SERVICE_NAME%" ^| findstr "STATE"') do set SERVICE_STATE=%%i

echo Status do Serviço: %SERVICE_STATE%
echo.
echo Opções:
echo   1 - Iniciar serviço
echo   2 - Parar serviço
echo   3 - Reiniciar serviço
echo   4 - Ver status detalhado
echo   5 - Ver logs (últimas 50 linhas)
echo   6 - Abrir pasta de logs
echo   7 - Testar API (GET /connect/v1/wsteste)
echo   0 - Sair
echo.

set /p OPTION="Digite sua opção: "

if "%OPTION%"=="1" goto START_SERVICE
if "%OPTION%"=="2" goto STOP_SERVICE
if "%OPTION%"=="3" goto RESTART_SERVICE
if "%OPTION%"=="4" goto STATUS_SERVICE
if "%OPTION%"=="5" goto VIEW_LOGS
if "%OPTION%"=="6" goto OPEN_LOGS
if "%OPTION%"=="7" goto TEST_API
if "%OPTION%"=="0" goto END

echo Opção inválida!
timeout /t 2 /nobreak >nul
goto MENU

:START_SERVICE
echo.
echo Iniciando serviço...
sc start "%SERVICE_NAME%"
if %ERRORLEVEL% EQU 0 (
    echo Serviço iniciado com sucesso!
) else (
    echo Erro ao iniciar serviço. Código: %ERRORLEVEL%
)
echo.
pause
goto MENU

:STOP_SERVICE
echo.
echo Parando serviço...
sc stop "%SERVICE_NAME%"
if %ERRORLEVEL% EQU 0 (
    echo Serviço parado com sucesso!
) else (
    echo Erro ao parar serviço. Código: %ERRORLEVEL%
)
echo.
pause
goto MENU

:RESTART_SERVICE
echo.
echo Parando serviço...
sc stop "%SERVICE_NAME%"
timeout /t 3 /nobreak >nul

echo Iniciando serviço...
sc start "%SERVICE_NAME%"
if %ERRORLEVEL% EQU 0 (
    echo Serviço reiniciado com sucesso!
) else (
    echo Erro ao reiniciar serviço. Código: %ERRORLEVEL%
)
echo.
pause
goto MENU

:STATUS_SERVICE
echo.
echo Status detalhado:
echo ----------------------------------------
sc query "%SERVICE_NAME%"
echo ----------------------------------------
echo.

echo Configuração:
echo ----------------------------------------
sc qc "%SERVICE_NAME%"
echo ----------------------------------------
echo.
pause
goto MENU

:VIEW_LOGS
echo.
set LOG_DIR=%~dp0..\log
if not exist "%LOG_DIR%" (
    echo Pasta de logs não encontrada: %LOG_DIR%
    pause
    goto MENU
)

REM Encontrar o log mais recente
for /f "delims=" %%i in ('dir /b /o-d "%LOG_DIR%\wsicrmrest_*.log" 2^>nul') do (
    set LATEST_LOG=%%i
    goto SHOW_LOG
)

echo Nenhum arquivo de log encontrado.
pause
goto MENU

:SHOW_LOG
echo Exibindo últimas 50 linhas do log: %LATEST_LOG%
echo ----------------------------------------
powershell -Command "Get-Content '%LOG_DIR%\%LATEST_LOG%' -Tail 50"
echo ----------------------------------------
echo.
pause
goto MENU

:OPEN_LOGS
echo.
set LOG_DIR=%~dp0..\log
if exist "%LOG_DIR%" (
    echo Abrindo pasta de logs...
    explorer "%LOG_DIR%"
) else (
    echo Pasta de logs não encontrada: %LOG_DIR%
)
timeout /t 2 /nobreak >nul
goto MENU

:TEST_API
echo.
echo Testando API...
echo ----------------------------------------

REM Verificar se o serviço está rodando
sc query "%SERVICE_NAME%" | findstr "RUNNING" >nul
if %ERRORLEVEL% NEQ 0 (
    echo ERRO: O serviço não está rodando!
    echo Inicie o serviço primeiro.
    pause
    goto MENU
)

REM Testar endpoint /connect/v1/wsteste
echo Requisição: GET http://localhost:8080/connect/v1/wsteste
echo.

curl -s http://localhost:8080/connect/v1/wsteste 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERRO: Não foi possível conectar à API.
    echo.
    echo Verifique:
    echo   1. Se o serviço está rodando
    echo   2. Se a porta 8080 está configurada corretamente
    echo   3. Se há firewall bloqueando a conexão
    echo.
    echo Nota: curl deve estar instalado para este teste funcionar.
) else (
    echo.
    echo API respondeu com sucesso!
)

echo ----------------------------------------
echo.
pause
goto MENU

:END
echo.
echo Saindo...
exit /b 0
