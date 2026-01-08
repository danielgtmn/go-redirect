#!/bin/bash

set -e

echo "Running Go Redirect Tests"
echo "========================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    local status=$1
    local message=$2
    if [ "$status" -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $message"
    else
        echo -e "${RED}✗${NC} $message"
        return 1
    fi
}

cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    docker stop redirect-test 2>/dev/null || true
    docker rm redirect-test 2>/dev/null || true
}

trap cleanup EXIT

echo -e "\n${YELLOW}1. Testing Docker build...${NC}"
if docker build -t go-redirect:test .; then
    print_status 0 "Docker build successful"
else
    print_status 1 "Docker build failed"
    exit 1
fi

echo -e "\n${YELLOW}2. Testing container startup...${NC}"
if docker run -d --name redirect-test \
    -p 8082:8080 \
    -e REDIRECT_TARGET=https://example.com \
    -e PRESERVE_PATH=true \
    go-redirect:test; then
    print_status 0 "Container started successfully"
else
    print_status 1 "Container failed to start"
    exit 1
fi

sleep 2

echo -e "\n${YELLOW}3. Testing health endpoint...${NC}"
if curl -f -s http://localhost:8082/health | grep -q "OK"; then
    print_status 0 "Health endpoint works"
else
    print_status 1 "Health endpoint failed"
    docker logs redirect-test
    exit 1
fi

echo -e "\n${YELLOW}4. Testing redirect (301)...${NC}"
LOCATION=$(curl -s -I http://localhost:8082/test/path 2>&1 | grep -i "^location:" | tr -d '\r')
if echo "$LOCATION" | grep -q "https://example.com/test/path"; then
    print_status 0 "Redirect with path preservation works"
else
    print_status 1 "Redirect failed. Location: $LOCATION"
    exit 1
fi

echo -e "\n${YELLOW}5. Testing redirect code...${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8082/)
if [ "$HTTP_CODE" = "301" ]; then
    print_status 0 "Correct redirect code (301)"
else
    print_status 1 "Wrong redirect code: $HTTP_CODE"
    exit 1
fi

echo -e "\n${GREEN}All tests passed!${NC}"
