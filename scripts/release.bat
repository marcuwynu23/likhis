@echo off
REM Script to create release builds for Windows and Linux
REM Creates archives with executables and plugins

set EXE_NAME=likhis
set RELEASE_DIR=release
set VERSION=%1

if "%VERSION%"=="" (
    echo Usage: release.bat [version]
    echo Example: release.bat v1.0.0
    exit /b 1
)

echo ========================================
echo Creating Release: %VERSION%
echo ========================================
echo.

REM Create version-specific release directory
set VERSION_DIR=%RELEASE_DIR%\%VERSION%
if exist "%VERSION_DIR%" (
    echo Cleaning existing release directory for version %VERSION%...
    rmdir /s /q "%VERSION_DIR%"
)
mkdir "%VERSION_DIR%"
echo.

REM Build Windows amd64
echo [1/5] Building Windows amd64...
set WINDOWS_DIR=%VERSION_DIR%\likhis-windows-amd64
mkdir "%WINDOWS_DIR%"
set GOOS=windows
set GOARCH=amd64
go build -o "%WINDOWS_DIR%\%EXE_NAME%.exe" main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Windows build failed!
    exit /b 1
)
echo   ✓ Windows executable built
echo.

REM Build Linux amd64
echo [2/5] Building Linux amd64...
set LINUX_DIR=%VERSION_DIR%\likhis-linux-amd64
mkdir "%LINUX_DIR%"
set GOOS=linux
set GOARCH=amd64
go build -o "%LINUX_DIR%\%EXE_NAME%" main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Linux build failed!
    exit /b 1
)
echo   ✓ Linux executable built
echo.

REM Copy plugins to Windows build
echo [3/5] Copying plugins...
xcopy /E /I /Y plugins "%WINDOWS_DIR%\plugins"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to copy plugins to Windows build!
    exit /b 1
)
echo   ✓ Plugins copied to Windows build

REM Copy plugins to Linux build
xcopy /E /I /Y plugins "%LINUX_DIR%\plugins"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to copy plugins to Linux build!
    exit /b 1
)
echo   ✓ Plugins copied to Linux build
echo.

REM Create archives
echo [4/5] Creating archives...

REM Create Windows ZIP
echo   Creating Windows ZIP archive...
powershell -Command "Compress-Archive -Path '%WINDOWS_DIR%\*' -DestinationPath '%VERSION_DIR%\likhis-windows-amd64-%VERSION%.zip' -Force"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to create Windows ZIP!
    exit /b 1
)
echo   ✓ Windows ZIP created: likhis-windows-amd64-%VERSION%.zip

REM Create Linux tar.gz using PowerShell
echo   Creating Linux tar.gz archive...
powershell -Command "$ErrorActionPreference='Stop'; Compress-Archive -Path '%LINUX_DIR%\*' -DestinationPath '%VERSION_DIR%\temp-linux.zip' -Force; $tar = '%VERSION_DIR%\likhis-linux-amd64-%VERSION%.tar.gz'; if (Test-Path $tar) { Remove-Item $tar }; tar -czf $tar -C '%LINUX_DIR%' .; Remove-Item '%VERSION_DIR%\temp-linux.zip'"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to create Linux tar.gz!
    echo Note: tar command requires Windows 10 version 1903+ or PowerShell 7+
    exit /b 1
)
echo   ✓ Linux tar.gz created: likhis-linux-amd64-%VERSION%.tar.gz
echo.

REM Generate release notes
echo [5/5] Generating release notes...
powershell -Command "$ErrorActionPreference='Stop'; $changelog = Get-Content 'CHANGELOG.md' -Raw; $version = '%VERSION%'; $pattern = '(?s)#### ' + [regex]::Escape($version) + '(.*?)(?=#### |$)'; if ($changelog -match $pattern) { $notes = $matches[1].Trim(); $notes | Out-File -FilePath '%VERSION_DIR%\RELEASE_NOTES.md' -Encoding UTF8; Write-Host '  ✓ Release notes generated' } else { $notes = '# Release Notes - ' + $version + '`n`n' + 'Release date: ' + (Get-Date -Format 'yyyy-MM-dd') + '`n`n' + 'See CHANGELOG.md for details.'; $notes | Out-File -FilePath '%VERSION_DIR%\RELEASE_NOTES.md' -Encoding UTF8; Write-Host '  ⚠ Release notes template created (version not found in CHANGELOG.md)' }"
if %ERRORLEVEL% NEQ 0 (
    echo   ⚠ Warning: Failed to generate release notes
) else (
    echo   ✓ Release notes generated
)
echo.

REM Cleanup build directories (optional - comment out if you want to keep them)
echo Cleaning up build directories...
rmdir /s /q "%WINDOWS_DIR%"
rmdir /s /q "%LINUX_DIR%"
echo   ✓ Build directories cleaned
echo.

echo ========================================
echo Release %VERSION% created successfully!
echo ========================================
echo.
echo Release directory: %VERSION_DIR%
echo.
echo Release files:
echo   - %VERSION_DIR%\likhis-windows-amd64-%VERSION%.zip
echo   - %VERSION_DIR%\likhis-linux-amd64-%VERSION%.tar.gz
echo   - %VERSION_DIR%\RELEASE_NOTES.md
echo.
pause

