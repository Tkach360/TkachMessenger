@echo off
::chcp 1252 > nul

if "%1"=="" (
  call :help
  goto :eof
)

if "%1"=="runcs" (
  call :runcs
  goto :eof
)

if "%1"=="run2cs" (
  call :run2cs
  goto :eof
)

if "%1"=="runcli" (
  call :runcli
  goto :eof
)



:help
echo Functions:
echo  runcs     - run client and server
echo  runui     - run user interface
goto :eof

:runcs
call :RunInNewWindow "go run ./cmd/server/main.go" "Server"
call :RunInNewWindow "go run ./cmd/client/main.go" "Client"
goto :eof

:run2cs
call :RunInNewWindow "go run ./cmd/server/main.go" "Server"
call :RunInNewWindow "go run ./cmd/client/main.go" "Client1"
call :RunInNewWindow "go run ./cmd/client/main.go" "Client2"
goto :eof

:runcli
call :RunInNewWindow "go run ./cmd/client/main.go" "Client"
goto :eof

:: Функция для запуска в новом окне
:RunInNewWindow
setlocal
set CMD=%~1
set TITLE=%~2
start "%TITLE%" cmd /c "%CMD% & pause"
endlocal
goto :eof
