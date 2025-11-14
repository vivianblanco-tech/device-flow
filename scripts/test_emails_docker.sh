#!/bin/bash

# Email Notification Testing Script (Docker Version)
# This script runs email notification tests inside Docker containers

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  Email Testing - Docker Environment${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Parse command line arguments
TEST_MODE="default"
CLEANUP=true

while [[ $# -gt 0 ]]; do
    case $1 in
        --isolated)
            TEST_MODE="isolated"
            shift
            ;;
        --no-cleanup)
            CLEANUP=false
            shift
            ;;
        --help)
            echo "Usage: ./scripts/test_emails_docker.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --isolated      Use separate test database (recommended)"
            echo "  --no-cleanup    Don't stop containers after test"
            echo "  --help          Show this help message"
            echo ""
            echo "Examples:"
            echo "  ./scripts/test_emails_docker.sh"
            echo "  ./scripts/test_emails_docker.sh --isolated"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Check if Docker is running
echo -e "${YELLOW}Checking Docker...${NC}"
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running${NC}"
    echo ""
    echo "Please start Docker and try again."
    exit 1
fi
echo -e "${GREEN}✅ Docker is running${NC}"
echo ""

# Check if docker-compose is available
echo -e "${YELLOW}Checking docker-compose...${NC}"
DOCKER_COMPOSE_CMD="docker-compose"
if ! command -v docker-compose &> /dev/null; then
    # Try docker compose (v2 syntax)
    if docker compose version &> /dev/null; then
        DOCKER_COMPOSE_CMD="docker compose"
    else
        echo -e "${RED}❌ docker-compose not found${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✅ docker-compose is available${NC}"
echo ""

# Get project root
cd "$(dirname "$0")/.."

# Start required services
echo -e "${YELLOW}Starting required services...${NC}"
$DOCKER_COMPOSE_CMD up -d postgres mailhog

echo -e "${GREEN}✅ Services started${NC}"
echo ""

# Wait for services to be healthy
echo -e "${YELLOW}Waiting for services to be ready...${NC}"

MAX_WAIT=30
WAITED=0
READY=false

while [ $READY = false ] && [ $WAITED -lt $MAX_WAIT ]; do
    sleep 2
    WAITED=$((WAITED + 2))
    
    # Check PostgreSQL
    PG_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' laptop-tracking-db 2>/dev/null || echo "unknown")
    
    # Check Mailhog
    MAILHOG_RUNNING=$(docker inspect --format='{{.State.Running}}' laptop-tracking-mailhog 2>/dev/null || echo "false")
    
    if [ "$PG_HEALTH" = "healthy" ] && [ "$MAILHOG_RUNNING" = "true" ]; then
        READY=true
    else
        echo -n "."
    fi
done

echo ""

if [ $READY = false ]; then
    echo -e "${YELLOW}⚠️  Services may not be fully ready, but continuing...${NC}"
else
    echo -e "${GREEN}✅ Services are ready${NC}"
fi
echo ""

# Build test container
echo -e "${YELLOW}Building test container...${NC}"

if [ "$TEST_MODE" = "isolated" ]; then
    echo -e "${CYAN}Using isolated test database${NC}"
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml build email-test-isolated
else
    echo -e "${CYAN}Using development database${NC}"
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml build email-test
fi

echo -e "${GREEN}✅ Test container built${NC}"
echo ""

# Run the tests
echo -e "${MAGENTA}========================================${NC}"
echo -e "${MAGENTA}  Running Email Notification Tests${NC}"
echo -e "${MAGENTA}========================================${NC}"
echo ""

if [ "$TEST_MODE" = "isolated" ]; then
    # Start isolated test database
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated up -d postgres-test
    sleep 5
    
    # Run tests
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated run --rm email-test-isolated
    EXIT_CODE=$?
else
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test
    EXIT_CODE=$?
fi

echo ""

# Show results
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  All Tests Passed!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${CYAN}💡 View emails in Mailhog:${NC}"
    echo "   http://localhost:8025"
    echo ""
    
    # Try to open browser
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open http://localhost:8025 2>/dev/null || true
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open http://localhost:8025 2>/dev/null || true
    fi
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}  Some Tests Failed${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
fi

# Cleanup
if [ "$CLEANUP" = true ]; then
    echo ""
    echo -e "${YELLOW}Cleaning up...${NC}"
    
    if [ "$TEST_MODE" = "isolated" ]; then
        $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated down -v
    fi
    
    echo -e "${GREEN}✅ Cleanup complete${NC}"
else
    echo ""
    echo -e "${YELLOW}⚠️  Containers left running (--no-cleanup flag used)${NC}"
    echo "To stop manually:"
    echo "  docker-compose down"
fi

exit $EXIT_CODE

