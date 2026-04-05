# Install script for lx on Windows PowerShell
# Usage: powershell -ExecutionPolicy Bypass -File install.ps1

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\lx"
)

$ErrorActionPreference = "Stop"

Write-Host "lx - Linux-to-Windows Command Translator" -ForegroundColor Cyan
Write-Host "Installation Script for PowerShell"
Write-Host ""

# Step 1: Create installation directory
Write-Host "Step 1: Creating installation directory..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Write-Host "Installation directory: $InstallDir" -ForegroundColor Green
Write-Host ""

# Step 2: Check if binary exists
Write-Host "Step 2: Checking for lx.exe..."
$binaryPath = Join-Path $InstallDir "lx.exe"
if (Test-Path $binaryPath) {
    Write-Host "Found existing lx.exe at: $binaryPath" -ForegroundColor Yellow
} else {
    Write-Host "Note: Copy lx.exe to: $binaryPath" -ForegroundColor Yellow
}
Write-Host ""

# Step 3: Add to PATH
Write-Host "Step 3: Adding to PATH..."
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$InstallDir*") {
    $newPath = "$userPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Added to PATH" -ForegroundColor Green
} else {
    Write-Host "Already in PATH" -ForegroundColor Green
}
Write-Host ""

# Step 4: Create PowerShell profile alias
Write-Host "Step 4: Creating PowerShell alias..."
$profilePath = $PROFILE.CurrentUserCurrentHost
$profileDir = Split-Path -Parent $profilePath

if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
}

if (-not (Test-Path $profilePath)) {
    New-Item -ItemType File -Path $profilePath -Force | Out-Null
}

$aliasCommand = "Set-Alias -Name lx -Value '$InstallDir\lx.exe' -Force"
$content = Get-Content $profilePath -ErrorAction SilentlyContinue
if ($content -notlike "*Set-Alias*lx*") {
    Add-Content -Path $profilePath -Value "`n$aliasCommand"
    Write-Host "Alias added to profile" -ForegroundColor Green
} else {
    Write-Host "Alias already exists" -ForegroundColor Green
}
Write-Host ""

Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "Close and reopen PowerShell to load changes"
