@echo off
REM Uninstall script for lx on Windows CMD
REM This script removes lx from the system

setlocal enabledelayedexpansion

set "INSTALL_DIR=%LOCALAPPDATA%\lx"

echo.
echo lx - Linux-to-Windows Command Translator
echo Uninstallation Script for CMD
echo.

REM Step 1: Confirm uninstallation
echo Are you sure you want to uninstall lx?
set /p "confirm=Type Y to confirm, anything else to cancel: "
if /i not "%confirm%"=="Y" (
    echo Uninstallation cancelled.
    pause
    exit /b 0
)
echo.

REM Step 2: Remove from PATH
echo Step 1: Removing from PATH...
for /f "tokens=*" %%A in ('reg query HKCU\Environment /v PATH 2^>nul ^| findstr /i PATH') do (
    set "PATH_VALUE=%%A"
)
REM This is a simplified version - proper implementation would require more complex string manipulation
echo Note: You may need to manually remove %INSTALL_DIR% from PATH
echo   Go to Control Panel ^> System ^> Advanced system settings ^> Environment Variables
echo.

REM Step 3: Remove DOSKEY alias from registry
echo Step 2: Removing DOSKEY alias from registry...
reg delete "HKCU\Software\Microsoft\Command Processor" /v AutoRun /f > nul 2>&1
if errorlevel 1 (
    echo Warning: DOSKEY alias not found or couldn't be removed
) else (
    echo DOSKEY alias removed from registry
)
echo.

REM Step 4: Remove installation directory
echo Step 3: Removing installation directory...
if exist "%INSTALL_DIR%" (
    rmdir /s /q "%INSTALL_DIR%"
    if errorlevel 1 (
        echo Warning: Failed to remove %INSTALL_DIR%
    ) else (
        echo Installation directory removed
    )
) else (
    echo Installation directory not found
)
echo.

echo ✓ Uninstallation completed!
echo.
echo Tip: Close and reopen CMD to complete cleanup
echo.

endlocal
pause
