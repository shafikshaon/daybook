#!/bin/bash

echo "========================================"
echo "Testing /api/v1/auth/me Fix"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test credentials - you can modify these
read -p "Enter username or email: " USERNAME
read -s -p "Enter password: " PASSWORD
echo ""
echo ""

# Test 1: Login to get token
echo "1. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

# Extract token and user ID
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
USER_ID=$(echo "$LOGIN_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Login failed${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo "  Token: ${TOKEN:0:30}..."
echo "  User ID from token: $USER_ID"
echo ""

# Test 2: Get profile with /auth/me
echo "2. Testing /auth/me endpoint..."
ME_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me)

HTTP_CODE=$(echo "$ME_RESPONSE" | grep HTTP_STATUS | cut -d':' -f2)
RESPONSE_BODY=$(echo "$ME_RESPONSE" | grep -v HTTP_STATUS)

echo "  HTTP Status: $HTTP_CODE"

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ /auth/me endpoint working!${NC}"
    echo ""
    echo "Response:"
    echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}SUCCESS! The fix is working!${NC}"
    echo -e "${GREEN}========================================${NC}"
elif [ "$HTTP_CODE" = "404" ]; then
    echo -e "${RED}✗ Still getting 404 - user not found${NC}"
    echo "Response: $RESPONSE_BODY"
    echo ""
    echo -e "${YELLOW}This means the user ID from the token doesn't exist in the database${NC}"
else
    echo -e "${RED}✗ Unexpected status code: $HTTP_CODE${NC}"
    echo "Response: $RESPONSE_BODY"
fi
echo ""

# Test 3: Test production if user wants
read -p "Test production endpoint? (y/n): " TEST_PROD

if [ "$TEST_PROD" = "y" ] || [ "$TEST_PROD" = "Y" ]; then
    echo ""
    echo "3. Testing Production /auth/me..."

    PROD_ME=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
      -H "Authorization: Bearer $TOKEN" \
      https://api.daybook.shafik.xyz/api/v1/auth/me)

    PROD_HTTP_CODE=$(echo "$PROD_ME" | grep HTTP_STATUS | cut -d':' -f2)
    PROD_RESPONSE=$(echo "$PROD_ME" | grep -v HTTP_STATUS)

    echo "  HTTP Status: $PROD_HTTP_CODE"

    if [ "$PROD_HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓ Production endpoint working!${NC}"
        echo "Response:"
        echo "$PROD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$PROD_RESPONSE"
    elif [ "$PROD_HTTP_CODE" = "404" ]; then
        echo -e "${RED}✗ Production still returning 404${NC}"
        echo -e "${YELLOW}→ You need to deploy the fix to production${NC}"
        echo "  1. Build backend: go build -o daybook-backend"
        echo "  2. Copy to production server"
        echo "  3. Restart backend service"
    elif [ "$PROD_HTTP_CODE" = "401" ]; then
        echo -e "${YELLOW}! Token not valid for production (different database)${NC}"
        echo "  This is expected if local and production use different databases"
    else
        echo -e "${RED}✗ Unexpected status code: $PROD_HTTP_CODE${NC}"
        echo "Response: $PROD_RESPONSE"
    fi
fi

echo ""
echo "========================================"
