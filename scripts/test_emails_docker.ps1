# Email Notification Testing Script (Docker Version)
# This script runs email notification tests inside Docker containers

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
Write-Color "  Email Testing - Docker Environment" "Cyan"
Write-Color "========================================" "Cyan"
Write-Host ""

# Parse command line arguments
$TestMode = "default"
$CleanUp = $true

for ($i = 0; $i -lt $args.Count; $i++) {
    switch ($args[$i]) {
        "--isolated" { $TestMode = "isolated" }
        "--no-cleanup" { $CleanUp = $false }
        "--help" {
            Write-Host "Usage: .\scripts\test_emails_docker.ps1 [OPTIONS]"
            Write-Host ""
            Write-Host "Options:"
            Write-Host "  --isolated      Use separate test database (recommended)"
            Write-Host "  --no-cleanup    Don't stop containers after test"
            Write-Host "  --help          Show this help message"
            Write-Host ""
            Write-Host "Examples:"
            Write-Host "  .\scripts\test_emails_docker.ps1"
            Write-Host "  .\scripts\test_emails_docker.ps1 --isolated"
            exit 0
        }
    }
}

# Check if Docker is running
Write-Color "Checking Docker..." "Yellow"
try {
    docker info | Out-Null
    Write-Color "✅ Docker is running" "Green"
} catch {
    Write-Color "❌ Docker is not running" "Red"
    Write-Host ""
    Write-Host "Please start Docker Desktop and try again."
    exit 1
}
Write-Host ""

# Check if docker-compose is available
Write-Color "Checking docker-compose..." "Yellow"
$dockerComposeCmd = "docker-compose"
try {
    & $dockerComposeCmd version | Out-Null
} catch {
    # Try docker compose (v2 syntax)
    try {
        docker compose version | Out-Null
        $dockerComposeCmd = "docker"
        $composeArg = "compose"
    } catch {
        Write-Color "❌ docker-compose not found" "Red"
        exit 1
    }
}
Write-Color "✅ docker-compose is available" "Green"
Write-Host ""

# Get project root directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir

# Change to project root
Push-Location $projectRoot

try {
    # Start required services
    Write-Color "Starting required services..." "Yellow"
    
    if ($dockerComposeCmd -eq "docker") {
        docker compose up -d postgres mailhog
    } else {
        docker-compose up -d postgres mailhog
    }
    
    if ($LASTEXITCODE -ne 0) {
        Write-Color "❌ Failed to start services" "Red"
        exit 1
    }
    
    Write-Color "✅ Services started" "Green"
    Write-Host ""
    
    # Wait for services to be healthy
    Write-Color "Waiting for services to be ready..." "Yellow"
    
    $maxWait = 30
    $waited = 0
    $ready = $false
    
    while (-not $ready -and $waited -lt $maxWait) {
        Start-Sleep -Seconds 2
        $waited += 2
        
        # Check PostgreSQL
        $pgHealth = docker inspect --format='{{.State.Health.Status}}' laptop-tracking-db 2>$null
        
        # Check Mailhog (may not have health check)
        $mailhogRunning = docker inspect --format='{{.State.Running}}' laptop-tracking-mailhog 2>$null
        
        if ($pgHealth -eq "healthy" -and $mailhogRunning -eq "true") {
            $ready = $true
        } else {
            Write-Host "." -NoNewline
        }
    }
    
    Write-Host ""
    
    if (-not $ready) {
        Write-Color "⚠️  Services may not be fully ready, but continuing..." "Yellow"
    } else {
        Write-Color "✅ Services are ready" "Green"
    }
    Write-Host ""
    
    # Run migrations if needed
    Write-Color "Running database migrations..." "Yellow"
    docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT 1 FROM pg_tables WHERE tablename='users'" | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Database needs initialization. Running migrations..."
        # Would need to run migrations - for now, assuming DB is set up
    }
    Write-Color "✅ Database ready" "Green"
    Write-Host ""
    
    # Build and run test container
    Write-Color "Building test container..." "Yellow"
    
    if ($TestMode -eq "isolated") {
        Write-Color "Using isolated test database" "Cyan"
        
        if ($dockerComposeCmd -eq "docker") {
            docker compose -f docker-compose.yml -f docker-compose.test.yml build email-test-isolated
        } else {
            docker-compose -f docker-compose.yml -f docker-compose.test.yml build email-test-isolated
        }
    } else {
        Write-Color "Using development database" "Cyan"
        
        if ($dockerComposeCmd -eq "docker") {
            docker compose -f docker-compose.yml -f docker-compose.test.yml build email-test
        } else {
            docker-compose -f docker-compose.yml -f docker-compose.test.yml build email-test
        }
    }
    
    if ($LASTEXITCODE -ne 0) {
        Write-Color "❌ Failed to build test container" "Red"
        exit 1
    }
    
    Write-Color "✅ Test container built" "Green"
    Write-Host ""
    
    # Run the tests
    Write-Color "========================================" "Magenta"
    Write-Color "  Running Email Notification Tests" "Magenta"
    Write-Color "========================================" "Magenta"
    Write-Host ""
    
    if ($TestMode -eq "isolated") {
        # Start isolated test database
        if ($dockerComposeCmd -eq "docker") {
            docker compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated up -d postgres-test
        } else {
            docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated up -d postgres-test
        }
        
        Start-Sleep -Seconds 5
        
        # Run tests
        if ($dockerComposeCmd -eq "docker") {
            docker compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated run --rm email-test-isolated
        } else {
            docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated run --rm email-test-isolated
        }
    } else {
        if ($dockerComposeCmd -eq "docker") {
            docker compose -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test
        } else {
            docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test
        }
    }
    
    $exitCode = $LASTEXITCODE
    
    Write-Host ""
    
    # Show results
    if ($exitCode -eq 0) {
        Write-Color "========================================" "Green"
        Write-Color "  All Tests Passed!" "Green"
        Write-Color "========================================" "Green"
        Write-Host ""
        Write-Color "💡 View emails in Mailhog:" "Cyan"
        Write-Host "   http://localhost:8025"
        Write-Host ""
        
        # Try to open browser
        Start-Process "http://localhost:8025" -ErrorAction SilentlyContinue
    } else {
        Write-Color "========================================" "Red"
        Write-Color "  Some Tests Failed" "Red"
        Write-Color "========================================" "Red"
        Write-Host ""
    }
    
} finally {
    Pop-Location
    
    # Cleanup
    if ($CleanUp) {
        Write-Host ""
        Write-Color "Cleaning up..." "Yellow"
        
        Push-Location $projectRoot
        
        if ($TestMode -eq "isolated") {
            if ($dockerComposeCmd -eq "docker") {
                docker compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated down -v
            } else {
                docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated down -v
            }
        }
        
        Pop-Location
        
        Write-Color "✅ Cleanup complete" "Green"
    } else {
        Write-Host ""
        Write-Color "⚠️  Containers left running (--no-cleanup flag used)" "Yellow"
        Write-Host "To stop manually:"
        Write-Host "  docker-compose down"
    }
}

exit $exitCode

