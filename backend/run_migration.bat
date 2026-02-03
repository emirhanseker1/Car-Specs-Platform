@echo off
REM Simple script to run SQL migration using compiled backend

cd /d %~dp0

echo Stopping backend...
taskkill /F /IM api.exe 2>nul
timeout /t 2 >nul

echo.
echo Creating migration runner...

REM Create a simple Go program that uses the compiled binary's database connection
echo package main > temp_migrate.go
echo. >> temp_migrate.go
echo import ( >> temp_migrate.go
echo     "database/sql" >> temp_migrate.go
echo     "fmt" >> temp_migrate.go
echo     "log" >> temp_migrate.go
echo     "os" >> temp_migrate.go
echo     _ "modernc.org/sqlite" >> temp_migrate.go
echo ) >> temp_migrate.go
echo. >> temp_migrate.go
echo func main() { >> temp_migrate.go
echo     db, err := sql.Open("sqlite", "./vehicles.db") >> temp_migrate.go
echo     if err != nil { log.Fatal(err) } >> temp_migrate.go
echo     defer db.Close() >> temp_migrate.go
echo. >> temp_migrate.go
echo     sqlBytes, err := os.ReadFile("./migrations/008_fix_audi_a3_data.sql") >> temp_migrate.go
echo     if err != nil { log.Fatal(err) } >> temp_migrate.go
echo. >> temp_migrate.go
echo     _, err = db.Exec(string(sqlBytes)) >> temp_migrate.go
echo     if err != nil { log.Fatal(err) } >> temp_migrate.go
echo. >> temp_migrate.go
echo     fmt.Println("Migration completed!") >> temp_migrate.go
echo } >> temp_migrate.go

echo Running migration with modernc.org/sqlite (no CGO needed)...
go run temp_migrate.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Migration successful!
    del temp_migrate.go
) else (
    echo.
    echo ❌ Migration failed!
    del temp_migrate.go
    exit /b 1
)

echo.
echo Restarting backend...
start /B api.exe

timeout /t 2 >nul
echo Backend restarted!
