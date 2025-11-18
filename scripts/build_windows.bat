@echo off
REM Script de compilação para Windows
REM Data de criação: 2025-11-17

echo ========================================
echo   WSICRMREST - Build para Windows
echo ========================================
echo.

REM Criar diretório de build se não existir
if not exist "build" mkdir build

REM Verificar se Go está instalado
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERRO: Go não encontrado no PATH!
    echo Por favor, instale o Go em https://golang.org/dl/
    pause
    exit /b 1
)

echo Versão do Go:
go version
echo.

REM Perguntar qual versão compilar
echo Escolha a versão para compilar:
echo 1 - Windows 32 bits
echo 2 - Windows 64 bits
echo 3 - Ambas as versões
echo.
set /p choice="Digite sua escolha (1-3): "

if "%choice%"=="1" goto build32
if "%choice%"=="2" goto build64
if "%choice%"=="3" goto buildall
echo Opção inválida!
pause
exit /b 1

:build32
echo.
echo Compilando para Windows 32 bits...
echo Obtendo data/hora de compilação...
for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'"') do set BUILD_TIME=%%i
echo Data/Hora da compilação: %BUILD_TIME%
set GOOS=windows
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags "-X 'wsicrmrest/internal/config.VersionDate=%BUILD_TIME%' -X 'wsicrmrest/internal/config.BuildTime=%BUILD_TIME%'" -o build\wsicrmrest_win32.exe cmd\server\main.go
if %ERRORLEVEL% EQU 0 (
    echo OK - Compilação concluída: build\wsicrmrest_win32.exe
) else (
    echo ERRO na compilação!
)
goto end

:build64
echo.
echo Compilando para Windows 64 bits...
echo Obtendo data/hora de compilação...
for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'"') do set BUILD_TIME=%%i
echo Data/Hora da compilação: %BUILD_TIME%
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -ldflags "-X 'wsicrmrest/internal/config.VersionDate=%BUILD_TIME%' -X 'wsicrmrest/internal/config.BuildTime=%BUILD_TIME%'" -o build\wsicrmrest_win64.exe cmd\server\main.go
if %ERRORLEVEL% EQU 0 (
    echo OK - Compilação concluída: build\wsicrmrest_win64.exe
) else (
    echo ERRO na compilação!
)
goto end

:buildall
echo.
echo Compilando para Windows 32 bits...
echo Obtendo data/hora de compilação...
for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'"') do set BUILD_TIME=%%i
echo Data/Hora da compilação: %BUILD_TIME%
set GOOS=windows
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags "-X 'wsicrmrest/internal/config.VersionDate=%BUILD_TIME%' -X 'wsicrmrest/internal/config.BuildTime=%BUILD_TIME%'" -o build\wsicrmrest_win32.exe cmd\server\main.go
if %ERRORLEVEL% EQU 0 (
    echo OK - Compilação 32 bits concluída: build\wsicrmrest_win32.exe
) else (
    echo ERRO na compilação 32 bits!
)

echo.
echo Compilando para Windows 64 bits...
echo Data/Hora da compilação: %BUILD_TIME%
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -ldflags "-X 'wsicrmrest/internal/config.VersionDate=%BUILD_TIME%' -X 'wsicrmrest/internal/config.BuildTime=%BUILD_TIME%'" -o build\wsicrmrest_win64.exe cmd\server\main.go
if %ERRORLEVEL% EQU 0 (
    echo OK - Compilação 64 bits concluída: build\wsicrmrest_win64.exe
) else (
    echo ERRO na compilação 64 bits!
)

:end
echo.
echo ========================================
echo   Build finalizado!
echo ========================================
echo.
echo Arquivos gerados na pasta build\
dir build\*.exe
echo.
pause
