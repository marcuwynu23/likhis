# PowerShell script to create symbolic link for likhis.exe to C:\Bin\webserve

$ExeName = "likhis.exe"
$BuildDir = Join-Path $PSScriptRoot "..\build"
$SourcePath = Join-Path $BuildDir $ExeName
$TargetDir = "C:\Bin\webserve"
$TargetPath = Join-Path $TargetDir $ExeName

Write-Host "Creating symbolic link for $ExeName..." -ForegroundColor Cyan
Write-Host ""

# Check if source file exists
if (-not (Test-Path $SourcePath)) {
    Write-Host "Error: $ExeName not found in $BuildDir" -ForegroundColor Red
    Write-Host "Please build the executable first: scripts\build.ps1" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Create target directory if it doesn't exist
if (-not (Test-Path $TargetDir)) {
    Write-Host "Creating directory: $TargetDir" -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
}

# Remove existing link if it exists
if (Test-Path $TargetPath) {
    Write-Host "Removing existing link..." -ForegroundColor Yellow
    Remove-Item $TargetPath -Force
}

# Create symbolic link
Write-Host "Creating symbolic link..." -ForegroundColor Cyan
Write-Host "  Source: $SourcePath" -ForegroundColor Gray
Write-Host "  Target: $TargetPath" -ForegroundColor Gray

try {
    $SourcePath = (Resolve-Path $SourcePath).Path
    New-Item -ItemType SymbolicLink -Path $TargetPath -Target $SourcePath -Force | Out-Null
    
    Write-Host ""
    Write-Host "Success! Symbolic link created." -ForegroundColor Green
    Write-Host "You can now run: $TargetPath" -ForegroundColor Green
} catch {
    Write-Host ""
    Write-Host "Error: Failed to create symbolic link." -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host "Make sure you're running PowerShell as Administrator." -ForegroundColor Yellow
}

Write-Host ""
Read-Host "Press Enter to exit"

