@echo off
echo Setting up Ka-Ping database...
echo.

REM Change to XAMPP MySQL bin directory
cd /d "C:\xampp2\mysql\bin"

REM Run the SQL setup script
mysql -u root -p < "C:\xampp2\htdocs\go\ping\ka-ping-server\setup_database.sql"

echo.
echo Database setup complete!
echo Press any key to continue...
pause > nul
