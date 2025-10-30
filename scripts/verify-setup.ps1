# Phase 0 Setup Verification Script
# This script validates that all Phase 0 components are in place and working

Write-Host "=== Phase 0 Setup Verification ===" -ForegroundColor Cyan
Write-Host ""

$errors = 0
$warnings = 0

# Check Go installation
Write-Host "[1/10] Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "  ✓ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Go is not installed or not in PATH" -ForegroundColor Red
    $errors++
}

# Check directory structure
Write-Host "`n[2/10] Checking directory structure..." -ForegroundColor Yellow
$requiredDirs = @(
    "cmd\web",
    "internal\config",
    "internal\database",
    "internal\models",
    "internal\handlers",
    "internal\middleware",
    "internal\auth",
    "internal\email",
    "internal\jira",
    "internal\validator",
    "migrations",
    "templates\layouts",
    "templates\pages",
    "templates\components",
    "static\css",
    "static\js",
    "static\images",
    "tests\unit",
    "tests\integration",
    "tests\e2e",
    "docs",
    "uploads"
)

foreach ($dir in $requiredDirs) {
    if (Test-Path $dir) {
        Write-Host "  ✓ $dir" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Missing: $dir" -ForegroundColor Red
        $errors++
    }
}

# Check required files
Write-Host "`n[3/10] Checking required files..." -ForegroundColor Yellow
$requiredFiles = @(
    "go.mod",
    "go.sum",
    "cmd\web\main.go",
    "internal\config\config.go",
    "internal\database\database.go",
    ".gitignore",
    "README.md",
    "Makefile",
    ".env.example",
    "Dockerfile",
    "docker-compose.yml"
)

foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Host "  ✓ $file" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Missing: $file" -ForegroundColor Red
        $errors++
    }
}

# Check .env file
Write-Host "`n[4/10] Checking environment configuration..." -ForegroundColor Yellow
if (Test-Path ".env") {
    Write-Host "  OK .env file exists" -ForegroundColor Green
} else {
    Write-Host "  WARNING .env file not found (copy from .env.example)" -ForegroundColor Yellow
    $warnings++
}

# Check Go modules
Write-Host "`n[5/10] Checking Go modules..." -ForegroundColor Yellow
try {
    $modCheck = go list -m all 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Go modules are valid" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Go modules have errors" -ForegroundColor Red
        $errors++
    }
} catch {
    Write-Host "  ✗ Failed to check Go modules" -ForegroundColor Red
    $errors++
}

# Test build
Write-Host "`n[6/10] Testing application build..." -ForegroundColor Yellow
try {
    $buildOutput = go build -o bin\test-build.exe cmd\web\main.go 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Application builds successfully" -ForegroundColor Green
        Remove-Item "bin\test-build.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "  ✗ Build failed: $buildOutput" -ForegroundColor Red
        $errors++
    }
} catch {
    Write-Host "  ✗ Build failed: $_" -ForegroundColor Red
    $errors++
}

# Run tests
Write-Host "`n[7/10] Running tests..." -ForegroundColor Yellow
try {
    $testOutput = go test ./... 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ All tests pass" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Tests failed: $testOutput" -ForegroundColor Red
        $errors++
    }
} catch {
    Write-Host "  ✗ Failed to run tests: $_" -ForegroundColor Red
    $errors++
}

# Check go vet
Write-Host "`n[8/10] Running go vet..." -ForegroundColor Yellow
try {
    $vetOutput = go vet ./... 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ No issues found by go vet" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Go vet found issues: $vetOutput" -ForegroundColor Red
        $errors++
    }
} catch {
    Write-Host "  ✗ Failed to run go vet: $_" -ForegroundColor Red
    $errors++
}

# Check git repository
Write-Host "`n[9/10] Checking git repository..." -ForegroundColor Yellow
if (Test-Path ".git") {
    Write-Host "  ✓ Git repository initialized" -ForegroundColor Green
    try {
        $commitCount = git rev-list --count HEAD 2>&1
        Write-Host "  ✓ Commits: $commitCount" -ForegroundColor Green
    } catch {
        Write-Host "  ⚠ Could not count commits" -ForegroundColor Yellow
        $warnings++
    }
} else {
    Write-Host "  ✗ Git repository not initialized" -ForegroundColor Red
    $errors++
}

# Check documentation
Write-Host "`n[10/10] Checking documentation..." -ForegroundColor Yellow
$docFiles = @("README.md", "docs\SETUP.md", "CONTRIBUTING.md")
$docCount = 0
foreach ($doc in $docFiles) {
    if (Test-Path $doc) {
        $docCount++
    }
}
if ($docCount -eq $docFiles.Count) {
    Write-Host "  ✓ All documentation files present ($docCount/$($docFiles.Count))" -ForegroundColor Green
} else {
    Write-Host "  ⚠ Some documentation missing ($docCount/$($docFiles.Count))" -ForegroundColor Yellow
    $warnings++
}

# Summary
Write-Host "`n=== Verification Summary ===" -ForegroundColor Cyan
Write-Host ""

if ($errors -eq 0 -and $warnings -eq 0) {
    Write-Host "SUCCESS Phase 0 setup is PERFECT!" -ForegroundColor Green
    Write-Host "All components are in place and working correctly." -ForegroundColor Green
    exit 0
} elseif ($errors -eq 0) {
    Write-Host "SUCCESS Phase 0 setup is COMPLETE with $warnings warning(s)" -ForegroundColor Yellow
    Write-Host "The setup is functional but has minor issues to address." -ForegroundColor Yellow
    exit 0
} else {
    Write-Host "ERROR Phase 0 setup has $errors error(s) and $warnings warning(s)" -ForegroundColor Red
    Write-Host "Please fix the errors before proceeding." -ForegroundColor Red
    exit 1
}

