@echo off
title Car Specs Launcher
echo ==========================================
echo       Car Specs Platform Launcher
echo ==========================================
echo.

echo 1. Normal Start (Backend + Frontend)
echo 2. Sync Data THEN Start (API Ninjas -> DB -> Start)
echo.
set /p choice="Choose mode (1 or 2): "

if "%choice%"=="2" (
    echo.
    echo [SYNC] Starting Data Sync...
    echo PLEASE WAIT until sync finishes...
    cd backend
    go run cmd/api/main.go -sync-ninjas
    cd ..
    echo [SYNC] Finished!
    echo.
)

echo [START] Starting Backend...
start "Car Specs Backend" cmd /k "cd backend && go run cmd/api/main.go"

echo [START] Starting Frontend...
start "Car Specs Frontend" cmd /k "cd frontend && npm run dev"

echo.
echo [INFO] Waiting for services to initialize...
timeout /t 5 >nul
echo [INFO] Opening Browser...
start http://localhost:5173

echo.
echo ==========================================
echo    System Running!
echo    Backend: http://localhost:8080
echo    Frontend: http://localhost:5173
echo ==========================================
echo.
echo Press any key to exit this launcher (servers will keep running)...
pause
