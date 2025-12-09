# Script to create release builds for Windows and Linux
# Creates archives with executables and plugins

param(
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$ErrorActionPreference = "Stop"

$ExeName = "likhis"
$ReleaseDir = "release"
$WindowsDir = "$ReleaseDir\likhis-windows-amd64"
$LinuxDir = "$ReleaseDir\likhis-linux-amd64"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Creating Release: $Version" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Create release directory
if (Test-Path $ReleaseDir) {
    Write-Host "Cleaning existing release directory..." -ForegroundColor Yellow
    Remove-Item -Path $ReleaseDir -Recurse -Force
}
New-Item -ItemType Directory -Path $ReleaseDir -Force | Out-Null
Write-Host ""

# Build Windows amd64
Write-Host "[1/4] Building Windows amd64..." -ForegroundColor Green
New-Item -ItemType Directory -Path $WindowsDir -Force | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$WindowsDir\$ExeName.exe" main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Windows build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ Windows executable built" -ForegroundColor Green
Write-Host ""

# Build Linux amd64
Write-Host "[2/4] Building Linux amd64..." -ForegroundColor Green
New-Item -ItemType Directory -Path $LinuxDir -Force | Out-Null
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$LinuxDir\$ExeName" main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Linux build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ Linux executable built" -ForegroundColor Green
Write-Host ""

# Copy plugins to Windows build
Write-Host "[3/4] Copying plugins..." -ForegroundColor Green
Copy-Item -Path "plugins\*" -Destination "$WindowsDir\plugins\" -Recurse -Force
if (-not (Test-Path "$WindowsDir\plugins")) {
    Write-Host "Error: Failed to copy plugins to Windows build!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ Plugins copied to Windows build" -ForegroundColor Green

# Copy plugins to Linux build
Copy-Item -Path "plugins\*" -Destination "$LinuxDir\plugins\" -Recurse -Force
if (-not (Test-Path "$LinuxDir\plugins")) {
    Write-Host "Error: Failed to copy plugins to Linux build!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ Plugins copied to Linux build" -ForegroundColor Green
Write-Host ""

# Create archives
Write-Host "[4/4] Creating archives..." -ForegroundColor Green

# Create Windows ZIP
Write-Host "  Creating Windows ZIP archive..." -ForegroundColor Yellow
$WindowsZip = "$ReleaseDir\likhis-windows-amd64-$Version.zip"
if (Test-Path $WindowsZip) {
    Remove-Item $WindowsZip -Force
}
Compress-Archive -Path "$WindowsDir\*" -DestinationPath $WindowsZip -Force
Write-Host "  ✓ Windows ZIP created: likhis-windows-amd64-$Version.zip" -ForegroundColor Green

# Create Linux tar.gz
Write-Host "  Creating Linux tar.gz archive..." -ForegroundColor Yellow
$LinuxTarGz = "$ReleaseDir\likhis-linux-amd64-$Version.tar.gz"
if (Test-Path $LinuxTarGz) {
    Remove-Item $LinuxTarGz -Force
}

# Check if tar command is available (Windows 10 1903+ or PowerShell 7+)
try {
    # Use tar command if available
    Push-Location $LinuxDir
    tar -czf "..\likhis-linux-amd64-$Version.tar.gz" *
    Pop-Location
    Write-Host "  ✓ Linux tar.gz created: likhis-linux-amd64-$Version.tar.gz" -ForegroundColor Green
} catch {
    Write-Host "  Warning: tar command not available, creating ZIP instead..." -ForegroundColor Yellow
    $LinuxZip = "$ReleaseDir\likhis-linux-amd64-$Version.zip"
    Compress-Archive -Path "$LinuxDir\*" -DestinationPath $LinuxZip -Force
    Write-Host "  ✓ Linux ZIP created: likhis-linux-amd64-$Version.zip" -ForegroundColor Green
}
Write-Host ""

# Cleanup build directories (optional - comment out if you want to keep them)
Write-Host "Cleaning up build directories..." -ForegroundColor Yellow
Remove-Item -Path $WindowsDir -Recurse -Force
Remove-Item -Path $LinuxDir -Recurse -Force
Write-Host "  ✓ Build directories cleaned" -ForegroundColor Green
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Release $Version created successfully!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Release files:" -ForegroundColor Green
Get-ChildItem -Path $ReleaseDir -Filter "*.zip","*.tar.gz" | ForEach-Object {
    Write-Host "  - $($_.Name)" -ForegroundColor White
}
Write-Host ""

