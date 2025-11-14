#!/bin/bash

# Email Notification Testing Script
# This script tests all email notifications by sending them and verifying they arrive in Mailhog

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Email Notifications Testing Setup${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if Mailhog is running
echo -e "${YELLOW}Checking if Mailhog is running...${NC}"
if ! curl -s http://localhost:8025/api/v2/messages > /dev/null 2>&1; then
    echo -e "${RED}❌ Mailhog is not running${NC}"
    echo ""
    echo -e "${YELLOW}To install and run Mailhog:${NC}"
    echo ""
    echo "  macOS:     brew install mailhog && mailhog"
    echo "  Linux:     go install github.com/mailhog/MailHog@latest && MailHog"
    echo "  Windows:   Download from https://github.com/mailhog/MailHog/releases"
    echo ""
    echo "Then access the web UI at: http://localhost:8025"
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ Mailhog is running${NC}"
echo ""

# Set environment variables for testing
export DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable}"
export MAILHOG_URL="${MAILHOG_URL:-http://localhost:8025}"
export SMTP_HOST="${SMTP_HOST:-localhost}"
export SMTP_PORT="${SMTP_PORT:-1025}"

echo -e "${BLUE}Environment Configuration:${NC}"
echo "  DATABASE_URL: $DATABASE_URL"
echo "  MAILHOG_URL:  $MAILHOG_URL"
echo "  SMTP_HOST:    $SMTP_HOST"
echo "  SMTP_PORT:    $SMTP_PORT"
echo ""

# Check if database is accessible
echo -e "${YELLOW}Checking database connection...${NC}"
if ! psql "$DATABASE_URL" -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to database${NC}"
    echo ""
    echo "Make sure your database is running and the DATABASE_URL is correct."
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ Database is accessible${NC}"
echo ""

# Build and run the test script
echo -e "${YELLOW}Building test script...${NC}"
cd "$(dirname "$0")/.."
go build -o /tmp/test_email_notifications ./scripts/email-test

echo -e "${GREEN}✅ Test script built${NC}"
echo ""

# Run the tests
echo -e "${YELLOW}Running email notification tests...${NC}"
echo ""
/tmp/test_email_notifications

# Capture exit code
EXIT_CODE=$?

# Clean up
rm -f /tmp/test_email_notifications

# Open Mailhog UI
if [ $EXIT_CODE -eq 0 ]; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  All Tests Passed!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${BLUE}💡 You can view the emails in Mailhog:${NC}"
    echo "   http://localhost:8025"
    echo ""
    
    # Try to open browser (platform-specific)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        open http://localhost:8025 2>/dev/null || true
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        xdg-open http://localhost:8025 2>/dev/null || true
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        # Windows
        start http://localhost:8025 2>/dev/null || true
    fi
fi

exit $EXIT_CODE

