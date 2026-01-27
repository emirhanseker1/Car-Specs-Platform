@echo off
echo Starting Car Specs API...
echo.
set DB_PATH=vehicles.db
set PORT=8080
cd cmd\api
go run main.go
