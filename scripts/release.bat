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

REM Create release directory
if exist "%RELEASE_DIR%" (
    echo Cleaning existing release directory...
    rmdir /s /q "%RELEASE_DIR%"
)
mkdir "%RELEASE_DIR%"
echo.

REM Build Windows amd64
echo [1/4] Building Windows amd64...
set WINDOWS_DIR=%RELEASE_DIR%\likhis-windows-amd64
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
echo [2/4] Building Linux amd64...
set LINUX_DIR=%RELEASE_DIR%\likhis-linux-amd64
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
echo [3/4] Copying plugins...
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
echo [4/4] Creating archives...

REM Create Windows ZIP
echo   Creating Windows ZIP archive...
powershell -Command "Compress-Archive -Path '%WINDOWS_DIR%\*' -DestinationPath '%RELEASE_DIR%\likhis-windows-amd64-%VERSION%.zip' -Force"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to create Windows ZIP!
    exit /b 1
)
echo   ✓ Windows ZIP created: likhis-windows-amd64-%VERSION%.zip

REM Create Linux tar.gz using PowerShell
echo   Creating Linux tar.gz archive...
powershell -Command "$ErrorActionPreference='Stop'; Compress-Archive -Path '%LINUX_DIR%\*' -DestinationPath '%RELEASE_DIR%\temp-linux.zip' -Force; $tar = '%RELEASE_DIR%\likhis-linux-amd64-%VERSION%.tar.gz'; if (Test-Path $tar) { Remove-Item $tar }; tar -czf $tar -C '%LINUX_DIR%' .; Remove-Item '%RELEASE_DIR%\temp-linux.zip'"
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to create Linux tar.gz!
    echo Note: tar command requires Windows 10 version 1903+ or PowerShell 7+
    exit /b 1
)
echo   ✓ Linux tar.gz created: likhis-linux-amd64-%VERSION%.tar.gz
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
echo Release files:
echo   - %RELEASE_DIR%\likhis-windows-amd64-%VERSION%.zip
echo   - %RELEASE_DIR%\likhis-linux-amd64-%VERSION%.tar.gz
echo.
pause

