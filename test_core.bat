@echo off

echo 🧪 Testing AI Todo Assistant Core Functionality...

REM Build the application first
echo Building application...
go build -o todo.exe cmd/todo/main.go
if %errorlevel% neq 0 (
    echo ❌ Build failed
    exit /b 1
)

echo ✅ Build successful

REM Test 1: Add a simple task
echo Testing task addition...
echo y | todo.exe add "Buy milk tomorrow" > test_output.txt 2>&1
if %errorlevel% equ 0 (
    echo ✅ Task addition test passed
) else (
    echo ⚠️ Task addition test had issues
    type test_output.txt
)

REM Test 2: List tasks
echo Testing task listing...
todo.exe list > test_output.txt 2>&1
if %errorlevel% equ 0 (
    echo ✅ Task listing test passed
) else (
    echo ⚠️ Task listing test had issues
    type test_output.txt
)

REM Clean up
del todo.exe >nul 2>&1
del test_output.txt >nul 2>&1

echo.
echo 🧪 Core functionality tests complete!
echo Check for any issues above.