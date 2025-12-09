# PowerShell script to build likhis.exe to build folder

$ExeName = "likhis.exe"
$BuildDir = "build"
$SourceFile = "main.go"

Write-Host "Building $ExeName..." -ForegroundColor Cyan
Write-Host ""

# Create build directory if it doesn't exist
if (-not (Test-Path $BuildDir)) {
    Write-Host "Creating build directory: $BuildDir" -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
}

# Build the executable
Write-Host "Building executable..." -ForegroundColor Cyan
$buildPath = Join-Path $BuildDir $ExeName
go build -o $buildPath $SourceFile

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Success! Executable built to: $buildPath" -ForegroundColor Green
    Write-Host ""
    Write-Host "You can now run the link script to create a symbolic link:" -ForegroundColor Yellow
    Write-Host "  scripts\link.ps1" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "Error: Build failed!" -ForegroundColor Red
}

Write-Host ""
Read-Host "Press Enter to exit"

