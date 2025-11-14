# Email Notification Testing Script (PowerShell)
# This script tests all email notifications by sending them and verifying they arrive in Mailhog

$ErrorActionPreference = "Stop"

# Colors
function Write-Color {
    param(
        [string]$Text,
        [string]$Color = "White"
    )
    Write-Host $Text -ForegroundColor $Color
}

Write-Color "========================================" "Cyan"
Write-Color "  Email Notifications Testing Setup" "Cyan"
Write-Color "========================================" "Cyan"
Write-Host ""

# Check if Mailhog is running
Write-Color "Checking if Mailhog is running..." "Yellow"
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8025/api/v2/messages" -Method Get -TimeoutSec 2 -UseBasicParsing
    Write-Color "✅ Mailhog is running" "Green"
} catch {
    Write-Color "❌ Mailhog is not running" "Red"
    Write-Host ""
    Write-Color "To install and run Mailhog:" "Yellow"
    Write-Host ""
    Write-Host "  1. Download from: https://github.com/mailhog/MailHog/releases"
    Write-Host "  2. Extract MailHog.exe"
    Write-Host "  3. Run: MailHog.exe"
    Write-Host ""
    Write-Host "Then access the web UI at: http://localhost:8025"
    Write-Host ""
    exit 1
}
Write-Host ""

# Set environment variables for testing
if (-not $env:DATABASE_URL) {
    $env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"
}
if (-not $env:MAILHOG_URL) {
    $env:MAILHOG_URL = "http://localhost:8025"
}
if (-not $env:SMTP_HOST) {
    $env:SMTP_HOST = "localhost"
}
if (-not $env:SMTP_PORT) {
    $env:SMTP_PORT = "1025"
}

Write-Color "Environment Configuration:" "Cyan"
Write-Host "  DATABASE_URL: $env:DATABASE_URL"
Write-Host "  MAILHOG_URL:  $env:MAILHOG_URL"
Write-Host "  SMTP_HOST:    $env:SMTP_HOST"
Write-Host "  SMTP_PORT:    $env:SMTP_PORT"
Write-Host ""

# Check if database is accessible
Write-Color "Checking database connection..." "Yellow"
$dbCheck = & psql $env:DATABASE_URL -c "SELECT 1" 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Color "❌ Cannot connect to database" "Red"
    Write-Host ""
    Write-Host "Make sure your database is running and the DATABASE_URL is correct."
    Write-Host ""
    exit 1
}

Write-Color "✅ Database is accessible" "Green"
Write-Host ""

# Get script directory and project root
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir

# Build the test script
Write-Color "Building test script..." "Yellow"
Push-Location $projectRoot
go build -o "$env:TEMP\test_email_notifications.exe" .\scripts\email-test
if ($LASTEXITCODE -ne 0) {
    Write-Color "❌ Failed to build test script" "Red"
    Pop-Location
    exit 1
}
Pop-Location

Write-Color "✅ Test script built" "Green"
Write-Host ""

# Run the tests
Write-Color "Running email notification tests..." "Yellow"
Write-Host ""
& "$env:TEMP\test_email_notifications.exe"
$exitCode = $LASTEXITCODE

# Clean up
Remove-Item "$env:TEMP\test_email_notifications.exe" -ErrorAction SilentlyContinue

# Open Mailhog UI
if ($exitCode -eq 0) {
    Write-Host ""
    Write-Color "========================================" "Green"
    Write-Color "  All Tests Passed!" "Green"
    Write-Color "========================================" "Green"
    Write-Host ""
    Write-Color "💡 You can view the emails in Mailhog:" "Cyan"
    Write-Host "   http://localhost:8025"
    Write-Host ""
    
    # Try to open browser
    Start-Process "http://localhost:8025"
}

exit $exitCode

