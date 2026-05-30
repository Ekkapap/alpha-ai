@echo off
if "%PROJECT_ROOT%"=="" (
  echo Error: PROJECT_ROOT is not set.
  exit /b 1
)
start "" "%PROJECT_ROOT%\graphify-out\graph.html"
