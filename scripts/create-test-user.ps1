# Create Test User Script
# This script creates a test user in the database for login testing

param(
    [string]$Email = "admin@bairesdev.com",
    [string]$Password = "Test123!",
    [string]$Role = "logistics",
    [string]$DBName = "laptop_tracking_dev",
    [string]$DBUser = "postgres",
    [string]$DBHost = "localhost",
    [string]$DBPort = "5432"
)

Write-Host "Creating test user..." -ForegroundColor Cyan
Write-Host "Email: $Email" -ForegroundColor Yellow
Write-Host "Password: $Password" -ForegroundColor Yellow
Write-Host "Role: $Role" -ForegroundColor Yellow
Write-Host ""

# Note: This uses the bcrypt hash for "Test123!" generated with cost 12
# In production, you should generate your own hash using the Go bcrypt package
$passwordHash = '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW'  # Hash for "Test123!"

# SQL to insert test user
$sql = @"
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES ('$Email', '$passwordHash', '$Role', NOW(), NOW())
ON CONFLICT (email) DO UPDATE
SET password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW()
RETURNING id, email, role;
"@

Write-Host "Executing SQL..." -ForegroundColor Cyan

# Set PostgreSQL password environment variable
$env:PGPASSWORD = $env:DB_PASSWORD
if (-not $env:PGPASSWORD) {
    $env:PGPASSWORD = "password"
}

# Execute SQL using psql
try {
    $result = & psql -h $DBHost -p $DBPort -U $DBUser -d $DBName -c $sql 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "Success! Test user created successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "==================================" -ForegroundColor Cyan
        Write-Host "LOGIN CREDENTIALS" -ForegroundColor Cyan
        Write-Host "==================================" -ForegroundColor Cyan
        Write-Host "Email:    $Email" -ForegroundColor White
        Write-Host "Password: $Password" -ForegroundColor White
        Write-Host "Role:     $Role" -ForegroundColor White
        Write-Host "==================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "You can now log in at: http://localhost:8080/login" -ForegroundColor Yellow
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "Failed to create user" -ForegroundColor Red
        Write-Host "Error: $result" -ForegroundColor Red
        Write-Host ""
        Write-Host "Make sure PostgreSQL is running and database exists:" -ForegroundColor Yellow
        Write-Host "  psql -h $DBHost -p $DBPort -U $DBUser -l" -ForegroundColor Gray
    }
} catch {
    Write-Host ""
    Write-Host "Error executing psql command" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Make sure:" -ForegroundColor Yellow
    Write-Host "  1. PostgreSQL is installed and running" -ForegroundColor Gray
    Write-Host "  2. psql is in your PATH" -ForegroundColor Gray
    Write-Host "  3. Database exists: $DBName" -ForegroundColor Gray
    Write-Host "  4. Migrations have been run: make migrate-up" -ForegroundColor Gray
}
