@echo off
if "%PROJECT_ROOT%"=="" (
  echo Error: PROJECT_ROOT is not set.
  exit /b 1
)
if exist "%PROJECT_ROOT%\.agents\tools\bin\windows\my-understand.exe" (
  "%PROJECT_ROOT%\.agents\tools\bin\windows\my-understand.exe" %*
) else (
  echo Error: my-understand.exe not found under α/tools/bin/windows/
  exit /b 1
)
