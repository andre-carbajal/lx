# Uninstall script for lx on Windows PowerShell
# Usage: powershell -ExecutionPolicy Bypass -File uninstall.ps1

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\lx"
)

$ErrorActionPreference = "Stop"

Write-Host "lx - Linux-to-Windows Command Translator" -ForegroundColor Cyan
Write-Host "Uninstallation Script for PowerShell"
Write-Host ""

# Confirm uninstallation
$confirm = Read-Host "Are you sure you want to uninstall lx? (Y/n)"
if ($confirm -ne "Y" -and $confirm -ne "") {
    Write-Host "Uninstallation cancelled" -ForegroundColor Yellow
    exit 0
}
Write-Host ""

# Step 1: Remove from PATH
Write-Host "Step 1: Removing from PATH..."
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -like "*$InstallDir*") {
    $pathArray = $userPath -split ";"
    $newPathArray = $pathArray | Where-Object { $_ -ne $InstallDir }
    $newPath = $newPathArray -join ";"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Removed from PATH" -ForegroundColor Green
} else {
    Write-Host "Not found in PATH" -ForegroundColor Gray
}
Write-Host ""

# Step 2: Remove alias from PowerShell profile
Write-Host "Step 2: Removing PowerShell alias..."
$profilePath = $PROFILE.CurrentUserCurrentHost
if (Test-Path $profilePath) {
    $content = Get-Content $profilePath
    $newContent = $content | Where-Object { $_ -notlike "*Set-Alias*lx*" }
    if ($newContent) {
        Set-Content -Path $profilePath -Value $newContent
    } else {
        Remove-Item -Path $profilePath -Force
    }
    Write-Host "Alias removed" -ForegroundColor Green
} else {
    Write-Host "Profile not found" -ForegroundColor Gray
}
Write-Host ""

# Step 3: Remove installation directory
Write-Host "Step 3: Removing installation directory..."
if (Test-Path $InstallDir) {
    Remove-Item -Path $InstallDir -Recurse -Force
    Write-Host "Installation directory removed" -ForegroundColor Green
} else {
    Write-Host "Installation directory not found" -ForegroundColor Gray
}
Write-Host ""

Write-Host "Uninstallation complete!" -ForegroundColor Green
Write-Host "Close and reopen PowerShell to complete cleanup"
