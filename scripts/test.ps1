# PowerShell script to test likhis on all framework examples

$ExePath = "build\likhis.exe"
$ExpDir = "exp"

Write-Host "Testing likhis on all framework examples..." -ForegroundColor Cyan
Write-Host ""

# Check if executable exists
if (-not (Test-Path $ExePath)) {
    Write-Host "Error: $ExePath not found!" -ForegroundColor Red
    Write-Host "Please build the executable first: scripts\build.ps1" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Check if exp directory exists
if (-not (Test-Path $ExpDir)) {
    Write-Host "Error: $ExpDir directory not found!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Create output directories
$OutDir = "out"
$frameworkDirs = @("express", "flask", "django", "laravel", "spring")

if (-not (Test-Path $OutDir)) {
    New-Item -ItemType Directory -Path $OutDir | Out-Null
}

foreach ($dir in $frameworkDirs) {
    $frameworkPath = Join-Path $OutDir $dir
    if (-not (Test-Path $frameworkPath)) {
        New-Item -ItemType Directory -Path $frameworkPath | Out-Null
    }
}

$testResults = @()

function Test-Framework {
    param(
        [string]$Name,
        [string]$Path,
        [string]$Framework,
        [string]$OutputSubDir
    )
    
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Testing $Name" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    
    $fullPath = Join-Path $ExpDir $Path
    $outputPath = Join-Path $OutDir $OutputSubDir
    & $ExePath -p $fullPath -o postman -F $Framework -O $outputPath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "[PASSED] $Name test passed" -ForegroundColor Green
        $script:testResults += @{Name=$Name; Status="PASSED"}
    } else {
        Write-Host "[FAILED] $Name test failed" -ForegroundColor Red
        $script:testResults += @{Name=$Name; Status="FAILED"}
    }
    Write-Host ""
}

# Test each framework
Test-Framework "Express.js" "express" "express" "express"
Test-Framework "Flask" "flask" "flask" "flask"
Test-Framework "Django" "django" "django" "django"
Test-Framework "Laravel" "laravel" "laravel" "laravel"
Test-Framework "Spring Boot" "spring" "spring" "spring"

# Test auto-detect
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing Auto-detect (all frameworks)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
$expressPath = Join-Path $ExpDir "express"
$expressOutPath = Join-Path $OutDir "express"
& $ExePath -p $expressPath -o postman -F auto -O $expressOutPath
if ($LASTEXITCODE -eq 0) {
    Write-Host "[PASSED] Auto-detect test passed" -ForegroundColor Green
    $testResults += @{Name="Auto-detect"; Status="PASSED"}
} else {
    Write-Host "[FAILED] Auto-detect test failed" -ForegroundColor Red
    $testResults += @{Name="Auto-detect"; Status="FAILED"}
}
Write-Host ""

# Test --full flag
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing --full flag (dev, staging, prod)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
& $ExePath -p $expressPath -o postman -F express --full -O $expressOutPath
if ($LASTEXITCODE -eq 0) {
    Write-Host "[PASSED] Full export test passed" -ForegroundColor Green
    $testResults += @{Name="Full Export"; Status="PASSED"}
} else {
    Write-Host "[FAILED] Full export test failed" -ForegroundColor Red
    $testResults += @{Name="Full Export"; Status="FAILED"}
}
Write-Host ""

# Test different output formats
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing different output formats" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

Write-Host "Testing Insomnia export..." -ForegroundColor Yellow
& $ExePath -p $expressPath -o insomnia -F express -O $expressOutPath
if ($LASTEXITCODE -eq 0) {
    Write-Host "[PASSED] Insomnia export passed" -ForegroundColor Green
    $testResults += @{Name="Insomnia Export"; Status="PASSED"}
} else {
    Write-Host "[FAILED] Insomnia export failed" -ForegroundColor Red
    $testResults += @{Name="Insomnia Export"; Status="FAILED"}
}

Write-Host "Testing HTTPie export..." -ForegroundColor Yellow
& $ExePath -p $expressPath -o httpie -F express -O $expressOutPath
if ($LASTEXITCODE -eq 0) {
    Write-Host "[PASSED] HTTPie export passed" -ForegroundColor Green
    $testResults += @{Name="HTTPie Export"; Status="PASSED"}
} else {
    Write-Host "[FAILED] HTTPie export failed" -ForegroundColor Red
    $testResults += @{Name="HTTPie Export"; Status="FAILED"}
}

Write-Host "Testing CURL export..." -ForegroundColor Yellow
& $ExePath -p $expressPath -o curl -F express -O $expressOutPath
if ($LASTEXITCODE -eq 0) {
    Write-Host "[PASSED] CURL export passed" -ForegroundColor Green
    $testResults += @{Name="CURL Export"; Status="PASSED"}
} else {
    Write-Host "[FAILED] CURL export failed" -ForegroundColor Red
    $testResults += @{Name="CURL Export"; Status="FAILED"}
}
Write-Host ""

# Test Summary
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$passed = ($testResults | Where-Object { $_.Status -eq "PASSED" }).Count
$failed = ($testResults | Where-Object { $_.Status -eq "FAILED" }).Count
$total = $testResults.Count

Write-Host "Total Tests: $total" -ForegroundColor White
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -gt 0) { "Red" } else { "Green" })
Write-Host ""

if ($failed -eq 0) {
    Write-Host "All tests passed! ✓" -ForegroundColor Green
} else {
    Write-Host "Some tests failed. Check the output above." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Check the generated files organized by framework:" -ForegroundColor Gray
Write-Host "  - $OutDir\express\" -ForegroundColor Gray
Write-Host "  - $OutDir\flask\" -ForegroundColor Gray
Write-Host "  - $OutDir\django\" -ForegroundColor Gray
Write-Host "  - $OutDir\laravel\" -ForegroundColor Gray
Write-Host "  - $OutDir\spring\" -ForegroundColor Gray
Write-Host ""
Read-Host "Press Enter to exit"

