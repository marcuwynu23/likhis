@echo off
REM Script to create symbolic link for likhis.exe to C:\Bin\webserve

set EXE_NAME=likhis.exe
set BUILD_DIR=%~dp0..\build
set SOURCE_PATH=%BUILD_DIR%\%EXE_NAME%
set TARGET_DIR=C:\Bin\webserve
set TARGET_PATH=%TARGET_DIR%\%EXE_NAME%

echo Creating symbolic link for %EXE_NAME%...
echo.

REM Check if source file exists
if not exist "%SOURCE_PATH%" (
    echo Error: %EXE_NAME% not found in %BUILD_DIR%
    echo Please build the executable first: scripts\build.bat
    pause
    exit /b 1
)

REM Create target directory if it doesn't exist
if not exist "%TARGET_DIR%" (
    echo Creating directory: %TARGET_DIR%
    mkdir "%TARGET_DIR%"
)

REM Remove existing link if it exists
if exist "%TARGET_PATH%" (
    echo Removing existing link...
    del "%TARGET_PATH%"
)

REM Create symbolic link
echo Creating symbolic link...
echo   Source: %SOURCE_PATH%
echo   Target: %TARGET_PATH%
mklink "%TARGET_PATH%" "%SOURCE_PATH%"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Success! Symbolic link created.
    echo You can now run: %TARGET_PATH%
) else (
    echo.
    echo Error: Failed to create symbolic link.
    echo Make sure you're running as Administrator.
)

echo.
pause

