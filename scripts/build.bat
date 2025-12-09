@echo off
REM Script to build likhis.exe to build folder

set EXE_NAME=likhis.exe
set BUILD_DIR=build
set SOURCE_FILE=main.go

echo Building %EXE_NAME%...
echo.

REM Create build directory if it doesn't exist
if not exist "%BUILD_DIR%" (
    echo Creating build directory: %BUILD_DIR%
    mkdir "%BUILD_DIR%"
)

REM Build the executable
echo Building executable...
go build -o "%BUILD_DIR%\%EXE_NAME%" %SOURCE_FILE%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Success! Executable built to: %BUILD_DIR%\%EXE_NAME%
    echo.
    echo You can now run the link script to create a symbolic link:
    echo   scripts\link.bat
) else (
    echo.
    echo Error: Build failed!
)

echo.
pause

