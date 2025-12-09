@echo off
REM Script to test likhis on all framework examples

set EXE_PATH=build\likhis.exe
set EXP_DIR=exp

echo Testing likhis on all framework examples...
echo.

REM Check if executable exists
if not exist "%EXE_PATH%" (
    echo Error: %EXE_PATH% not found!
    echo Please build the executable first: scripts\build.bat
    pause
    exit /b 1
)

REM Check if exp directory exists
if not exist "%EXP_DIR%" (
    echo Error: %EXP_DIR% directory not found!
    pause
    exit /b 1
)

REM Create output directories
if not exist "out" mkdir out
if not exist "out\express" mkdir out\express
if not exist "out\flask" mkdir out\flask
if not exist "out\django" mkdir out\django
if not exist "out\laravel" mkdir out\laravel
if not exist "out\spring" mkdir out\spring

echo ========================================
echo Testing Express.js
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\express" -o postman -F express -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Express test failed
) else (
    echo [PASSED] Express test passed
)
echo.

echo ========================================
echo Testing Flask
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\flask" -o postman -F flask -O out\flask
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Flask test failed
) else (
    echo [PASSED] Flask test passed
)
echo.

echo ========================================
echo Testing Django
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\django" -o postman -F django -O out\django
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Django test failed
) else (
    echo [PASSED] Django test passed
)
echo.

echo ========================================
echo Testing Laravel
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\laravel" -o postman -F laravel -O out\laravel
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Laravel test failed
) else (
    echo [PASSED] Laravel test passed
)
echo.

echo ========================================
echo Testing Spring Boot
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\spring" -o postman -F spring -O out\spring
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Spring test failed
) else (
    echo [PASSED] Spring test passed
)
echo.

echo ========================================
echo Testing Auto-detect (all frameworks)
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\express" -o postman -F auto -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Auto-detect test failed
) else (
    echo [PASSED] Auto-detect test passed
)
echo.

echo ========================================
echo Testing --full flag (dev, staging, prod)
echo ========================================
"%EXE_PATH%" -p "%EXP_DIR%\express" -o postman -F express --full -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Full export test failed
) else (
    echo [PASSED] Full export test passed
)
echo.

echo ========================================
echo Testing different output formats
echo ========================================
echo Testing Insomnia export...
"%EXE_PATH%" -p "%EXP_DIR%\express" -o insomnia -F express -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] Insomnia export failed
) else (
    echo [PASSED] Insomnia export passed
)

echo Testing HTTPie export...
"%EXE_PATH%" -p "%EXP_DIR%\express" -o httpie -F express -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] HTTPie export failed
) else (
    echo [PASSED] HTTPie export passed
)

echo Testing CURL export...
"%EXE_PATH%" -p "%EXP_DIR%\express" -o curl -F express -O out\express
if %ERRORLEVEL% NEQ 0 (
    echo [FAILED] CURL export failed
) else (
    echo [PASSED] CURL export passed
)
echo.

echo ========================================
echo Test Summary
echo ========================================
echo All tests completed!
echo Check the generated files organized by framework:
echo   - out\express\
echo   - out\flask\
echo   - out\django\
echo   - out\laravel\
echo   - out\spring\
echo.
pause

