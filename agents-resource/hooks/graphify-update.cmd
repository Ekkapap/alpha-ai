@echo off
if "%PROJECT_ROOT%"=="" (
  echo Error: PROJECT_ROOT is not set.
  exit /b 1
)
bash "%PROJECT_ROOT%\.agents\scripts\graphify.sh" . --update %*
