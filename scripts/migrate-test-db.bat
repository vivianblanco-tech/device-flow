@echo off
echo Applying migrations to test database...

for %%f in (migrations\*up.sql) do (
    echo Applying: %%~nxf
    docker cp "%%f" laptop-tracking-db:/tmp/migration.sql >nul 2>&1
    docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -f /tmp/migration.sql >nul 2>&1
    if errorlevel 1 (
        echo   Warning: %%~nxf
    ) else (
        echo   Success: %%~nxf
    )
)

echo.
echo Migrations complete!

