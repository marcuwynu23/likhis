#!/bin/bash
# Script to test likhis on all framework examples

EXE_PATH="build/likhis"
EXP_DIR="exp"

echo "Testing likhis on all framework examples..."
echo ""

# Check if executable exists
if [ ! -f "$EXE_PATH" ]; then
    echo "Error: $EXE_PATH not found!"
    echo "Please build the executable first: scripts/build.sh"
    exit 1
fi

# Check if exp directory exists
if [ ! -d "$EXP_DIR" ]; then
    echo "Error: $EXP_DIR directory not found!"
    exit 1
fi

# Create output directories
mkdir -p out/express
mkdir -p out/flask
mkdir -p out/django
mkdir -p out/laravel
mkdir -p out/spring

test_results=()

test_framework() {
    local name=$1
    local path=$2
    local framework=$3
    local output_dir=$4
    
    echo "========================================"
    echo "Testing $name"
    echo "========================================"
    
    "$EXE_PATH" -p "$EXP_DIR/$path" -o postman -F "$framework" -O "$output_dir"
    
    if [ $? -eq 0 ]; then
        echo "[PASSED] $name test passed"
        test_results+=("PASSED:$name")
    else
        echo "[FAILED] $name test failed"
        test_results+=("FAILED:$name")
    fi
    echo ""
}

# Test each framework
test_framework "Express.js" "express" "express" "out/express"
test_framework "Flask" "flask" "flask" "out/flask"
test_framework "Django" "django" "django" "out/django"
test_framework "Laravel" "laravel" "laravel" "out/laravel"
test_framework "Spring Boot" "spring" "spring" "out/spring"

# Test auto-detect
echo "========================================"
echo "Testing Auto-detect (all frameworks)"
echo "========================================"
"$EXE_PATH" -p "$EXP_DIR/express" -o postman -F auto -O out/express
if [ $? -eq 0 ]; then
    echo "[PASSED] Auto-detect test passed"
    test_results+=("PASSED:Auto-detect")
else
    echo "[FAILED] Auto-detect test failed"
    test_results+=("FAILED:Auto-detect")
fi
echo ""

# Test --full flag
echo "========================================"
echo "Testing --full flag (dev, staging, prod)"
echo "========================================"
"$EXE_PATH" -p "$EXP_DIR/express" -o postman -F express --full -O out/express
if [ $? -eq 0 ]; then
    echo "[PASSED] Full export test passed"
    test_results+=("PASSED:Full Export")
else
    echo "[FAILED] Full export test failed"
    test_results+=("FAILED:Full Export")
fi
echo ""

# Test different output formats
echo "========================================"
echo "Testing different output formats"
echo "========================================"

echo "Testing Insomnia export..."
"$EXE_PATH" -p "$EXP_DIR/express" -o insomnia -F express -O out/express
if [ $? -eq 0 ]; then
    echo "[PASSED] Insomnia export passed"
    test_results+=("PASSED:Insomnia Export")
else
    echo "[FAILED] Insomnia export failed"
    test_results+=("FAILED:Insomnia Export")
fi

echo "Testing HTTPie export..."
"$EXE_PATH" -p "$EXP_DIR/express" -o httpie -F express -O out/express
if [ $? -eq 0 ]; then
    echo "[PASSED] HTTPie export passed"
    test_results+=("PASSED:HTTPie Export")
else
    echo "[FAILED] HTTPie export failed"
    test_results+=("FAILED:HTTPie Export")
fi

echo "Testing CURL export..."
"$EXE_PATH" -p "$EXP_DIR/express" -o curl -F express -O out/express
if [ $? -eq 0 ]; then
    echo "[PASSED] CURL export passed"
    test_results+=("PASSED:CURL Export")
else
    echo "[FAILED] CURL export failed"
    test_results+=("FAILED:CURL Export")
fi
echo ""

# Test Summary
echo "========================================"
echo "Test Summary"
echo "========================================"

passed=0
failed=0
for result in "${test_results[@]}"; do
    if [[ $result == PASSED:* ]]; then
        ((passed++))
    else
        ((failed++))
    fi
done

total=$((passed + failed))
echo "Total Tests: $total"
echo "Passed: $passed"
echo "Failed: $failed"
echo ""

if [ $failed -eq 0 ]; then
    echo "All tests passed! ✓"
else
    echo "Some tests failed. Check the output above."
fi

echo ""
echo "Check the generated files organized by framework:"
echo "  - out/express/"
echo "  - out/flask/"
echo "  - out/django/"
echo "  - out/laravel/"
echo "  - out/spring/"
echo ""

