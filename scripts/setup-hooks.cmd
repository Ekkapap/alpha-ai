@echo off
:: α/scripts/setup-hooks.cmd
:: Windows Native Setup — self-contained shims, no project-hook dependency
setlocal enabledelayedexpansion

echo Running Windows Alpha Toolchain Setup...

:: 1. Resolve INSTALL_ROOT (walk up from this script's location to find α\)
set "CURRENT_DIR=%~dp0"
:find_root
if exist "!CURRENT_DIR!α" (
    set "INSTALL_ROOT=!CURRENT_DIR:~0,-1!"
    goto :found_root
)
for %%I in ("!CURRENT_DIR!..") do set "PARENT=%%~fI"
if "!PARENT!"=="!CURRENT_DIR:~0,-1!" goto :no_root
set "CURRENT_DIR=!PARENT!\"
goto :find_root
:no_root
echo Error: Could not find Alpha project root (α/ directory).
exit /b 1
:found_root

set "ALPHA=%INSTALL_ROOT%\α"
set "BIN_DIR=%ALPHA%\hooks\bin"
set "SCRIPT_DIR=%ALPHA%\scripts"

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%ALPHA%\tools\bin\windows" mkdir "%ALPHA%\tools\bin\windows"
if not exist "%ALPHA%\memories" mkdir "%ALPHA%\memories"

echo Install Root: %INSTALL_ROOT%

:: 2. Locate Go and compile binaries
set "GO_PATH="
where go >nul 2>&1
if %errorlevel% equ 0 (
    set "GO_PATH=go"
) else if exist "C:\Program Files\Go\bin\go.exe" (
    set "GO_PATH=C:\Program Files\Go\bin\go.exe"
)

if not "%GO_PATH%"=="" (
    echo Go Compiler found: %GO_PATH%
    if exist "%ALPHA%\tools\my-graphify" (
        echo    Compiling my-graphify.exe...
        cd /d "%ALPHA%\tools\my-graphify"
        "%GO_PATH%" build -o "%ALPHA%\tools\bin\windows\my-graphify.exe" main.go
        if %errorlevel% equ 0 (echo    OK: my-graphify.exe) else (echo    FAILED: my-graphify.exe)
    )
    if exist "%ALPHA%\tools\my-understand" (
        echo    Compiling my-understand.exe...
        cd /d "%ALPHA%\tools\my-understand"
        "%GO_PATH%" build -o "%ALPHA%\tools\bin\windows\my-understand.exe" main.go
        if %errorlevel% equ 0 (echo    OK: my-understand.exe) else (echo    FAILED: my-understand.exe)
    )
    cd /d "%INSTALL_ROOT%"
) else (
    echo WARNING: Go not found. Using existing binaries.
)

:: 3. Helper — write a self-contained root-walking CMD shim
::    Usage: call :write_shim <name> <body_lines_file>
::    Each shim detects PROJECT_ROOT by walking up to find α\

set "TOOLS=awake sync focus forget map graphify graphify-update graphify-cluster-only overview sketch detail understand"

for %%T in (%TOOLS%) do (
    set "SHIM=%BIN_DIR%\%%T.cmd"
    call :write_shim_header "!SHIM!"
    call :write_shim_body_%%T "!SHIM!"
    echo    Shim: %%T.cmd
)

:: 4. Add α\hooks\bin to User PATH (persistent, no ~/.local/bin needed)
echo.
echo Adding α\hooks\bin to user PATH...
powershell -NoProfile -Command ^
  "$path = [Environment]::GetEnvironmentVariable('PATH','User');" ^
  "$alpha = '%ALPHA%\hooks\bin';" ^
  "if ($path -notlike ('*' + $alpha + '*')) {" ^
  "  [Environment]::SetEnvironmentVariable('PATH', $path + ';' + $alpha, 'User');" ^
  "  Write-Host 'Added to PATH: ' + $alpha" ^
  "} else { Write-Host 'Already in PATH.' }"

:: 5. PowerShell profile — add doskey-style functions
set "PS_PROFILE=%USERPROFILE%\Documents\PowerShell\Microsoft.PowerShell_profile.ps1"
if not exist "%USERPROFILE%\Documents\PowerShell" mkdir "%USERPROFILE%\Documents\PowerShell"
powershell -NoProfile -Command ^
  "$marker = '# === ALPHA_HOOKS_START ===';" ^
  "$end = '# === ALPHA_HOOKS_END ===';" ^
  "$profile = '%PS_PROFILE%';" ^
  "$content = if (Test-Path $profile) { Get-Content $profile -Raw } else { '' };" ^
  "$content = $content -replace ('(?s)' + [regex]::Escape($marker) + '.*?' + [regex]::Escape($end) + '\r?\n?'), '';" ^
  "$hooks = @(" ^
  "  $marker," ^
  "  'function Find-AlphaRoot { $d=(Get-Location).Path; while ($d -ne (Split-Path $d -Parent)) { if (Test-Path (Join-Path $d \"α\")) { return $d } ; $d=Split-Path $d -Parent }; return $null }'," ^
  "  'function Invoke-AlphaTool($tool, $rest) { $r=Find-AlphaRoot; if (!$r) { Write-Error \"No α/ project found\" ; return }; & \"$r\α\tools\bin\windows\my-graphify.exe\" $tool @rest }'," ^
  "  'function awake { Invoke-AlphaTool awake $args }'," ^
  "  'function sync { Invoke-AlphaTool sync $args }'," ^
  "  'function overview { Invoke-AlphaTool overview $args }'," ^
  "  'function sketch { Invoke-AlphaTool sketch $args }'," ^
  "  'function detail { Invoke-AlphaTool detail $args }'," ^
  "  'function focus { Invoke-AlphaTool focus $args }'," ^
  "  'function forget { Invoke-AlphaTool forget $args }'," ^
  "  $end" ^
  ") -join \"`n\";" ^
  "Set-Content $profile ($content.TrimEnd() + \"`n`n\" + $hooks)" ^
  "Write-Host 'PowerShell profile updated: ' + $profile"

echo.
echo Setup complete.
echo   - CMD:        awake, sync, overview, sketch, detail, focus, forget
echo   - PowerShell: awake, sync, overview, sketch, detail, focus, forget
echo   - Restart terminal to activate PATH changes.
exit /b 0

:: ── Subroutines ──────────────────────────────────────────────────────────────

:write_shim_header
set "f=%~1"
(
echo @echo off
echo :: Self-contained — walks up directory tree to find α\, no project-hook needed
echo setlocal enabledelayedexpansion
echo set "SEARCH=%CD%"
echo :_find
echo if exist "!SEARCH!\α" ^( set "PROJECT_ROOT=!SEARCH!" ^& goto :_found ^)
echo set "PREV=!SEARCH!"
echo for %%%%I in ^("!SEARCH!\.."^) do set "SEARCH=%%%%~fI"
echo if "!SEARCH!"=="!PREV!" goto :_nofound
echo goto :_find
echo :_nofound
echo echo Error: No α/ project found in directory tree.
echo exit /b 1
echo :_found
) > "%f%"
exit /b 0

:write_shim_body_awake
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" awake %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_sync
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" sync %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_focus
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" focus %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_forget
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" forget %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_overview
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" overview %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_sketch
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" sketch %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_detail
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-graphify.exe" detail %%* >> %~1
echo ^) else ^( echo Error: my-graphify.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_understand
echo if exist "%%PROJECT_ROOT%%\α\tools\bin\windows\my-understand.exe" ^( >> %~1
echo   "%%PROJECT_ROOT%%\α\tools\bin\windows\my-understand.exe" %%* >> %~1
echo ^) else ^( echo Error: my-understand.exe not found. ^& exit /b 1 ^) >> %~1
exit /b 0

:write_shim_body_map
echo start "" "%%PROJECT_ROOT%%\graphify-out\graph.html" >> %~1
exit /b 0

:write_shim_body_graphify
echo bash "%%PROJECT_ROOT%%\α\scripts\graphify.sh" . %%* >> %~1
exit /b 0

:write_shim_body_graphify-update
echo bash "%%PROJECT_ROOT%%\α\scripts\graphify.sh" . --update %%* >> %~1
exit /b 0

:write_shim_body_graphify-cluster-only
echo bash "%%PROJECT_ROOT%%\α\scripts\graphify.sh" . --cluster-only %%* >> %~1
exit /b 0
