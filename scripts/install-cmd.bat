@echo off
REM Install script for lx on Windows CMD
REM This script installs lx to a user directory and adds it to PATH

setlocal enabledelayedexpansion

set "INSTALL_DIR=%LOCALAPPDATA%\lx"
set "BINARY_NAME=lx.exe"

echo.
echo lx - Linux-to-Windows Command Translator
echo Installation Script for CMD
echo.

REM Step 1: Create installation directory
echo Step 1: Preparing installation directory...
if not exist "%INSTALL_DIR%" (
    echo Creating directory: %INSTALL_DIR%
    mkdir "%INSTALL_DIR%"
)
echo Installation directory: %INSTALL_DIR%
echo.

REM Step 2: Check if binary exists
echo Step 2: Checking if %BINARY_NAME% exists...
if exist "%INSTALL_DIR%\%BINARY_NAME%" (
    echo Note: %BINARY_NAME% already exists at: %INSTALL_DIR%\%BINARY_NAME%
    echo Please copy the new binary to overwrite the existing one.
) else (
    echo Note: Please copy %BINARY_NAME% to: %INSTALL_DIR%\%BINARY_NAME%
)
echo.

REM Step 3: Add to PATH
echo Step 3: Adding %INSTALL_DIR% to PATH...
setx PATH "!PATH!;%INSTALL_DIR%" > nul
if errorlevel 1 (
    echo Warning: Failed to add to PATH
) else (
    echo Added to PATH successfully
)
echo.

REM Step 4: Create DOSKEY alias
echo Step 4: Creating DOSKEY macro...
doskey lx=%INSTALL_DIR%\%BINARY_NAME% $*
echo DOSKEY alias created for this session
echo.

REM Step 5: Persist DOSKEY macro via registry
echo Step 5: Persisting DOSKEY alias...
reg add "HKCU\Software\Microsoft\Command Processor" /v AutoRun /t REG_SZ /d "doskey lx=%INSTALL_DIR%\%BINARY_NAME% $*" /f > nul
if errorlevel 1 (
    echo Warning: Failed to persist DOSKEY alias via registry
) else (
    echo DOSKEY alias persisted in registry
)
echo.

REM Step 6: Verify installation
echo Step 6: Verifying installation...
if exist "%INSTALL_DIR%\%BINARY_NAME%" (
    echo.
    echo ✓ Installation completed successfully!
    echo.
    echo Next steps:
    echo   1. Close and reopen CMD to reload PATH and DOSKEY
    echo   2. Run 'lx --help' to see available commands
    echo   3. Try 'lx ls' to test the translator
) else (
    echo Note: %INSTALL_DIR%\%BINARY_NAME% not found yet
    echo Please copy lx.exe to the installation directory
)
echo.

endlocal
pause
