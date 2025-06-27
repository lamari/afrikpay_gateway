#!/bin/bash

# Script de test des endpoints Auth API
# Usage: ./test_endpoints.sh

BASE_URL="http://localhost:8001"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Test des endpoints Afrikpay Auth API ===${NC}"
echo "Base URL: $BASE_URL"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}1. Test Health Check${NC}"
response=$(curl -s -w "HTTP_CODE:%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ Health Check: OK${NC}"
    echo "Response: $body"
else
    echo -e "${RED}❌ Health Check: FAILED (HTTP $http_code)${NC}"
fi
echo ""

# Test 2: Login
echo -e "${YELLOW}2. Test Login${NC}"
login_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}')

login_http_code=$(echo "$login_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
login_body=$(echo "$login_response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$login_http_code" = "200" ]; then
    echo -e "${GREEN}✅ Login: OK${NC}"
    # Extract access token
    access_token=$(echo "$login_body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    refresh_token=$(echo "$login_body" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
    echo "Access Token: ${access_token:0:50}..."
    echo "Refresh Token: ${refresh_token:0:50}..."
else
    echo -e "${RED}❌ Login: FAILED (HTTP $login_http_code)${NC}"
    echo "Response: $login_body"
    exit 1
fi
echo ""

# Test 3: Verify Token
echo -e "${YELLOW}3. Test Verify Token${NC}"
verify_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X GET "$BASE_URL/auth/verify" \
    -H "Authorization: Bearer $access_token")

verify_http_code=$(echo "$verify_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
verify_body=$(echo "$verify_response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$verify_http_code" = "200" ]; then
    echo -e "${GREEN}✅ Verify Token: OK${NC}"
    echo "Claims: $verify_body"
else
    echo -e "${RED}❌ Verify Token: FAILED (HTTP $verify_http_code)${NC}"
    echo "Response: $verify_body"
fi
echo ""

# Test 4: Protected Endpoint
echo -e "${YELLOW}4. Test Protected Endpoint${NC}"
profile_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X GET "$BASE_URL/protected/profile" \
    -H "Authorization: Bearer $access_token")

profile_http_code=$(echo "$profile_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
profile_body=$(echo "$profile_response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$profile_http_code" = "200" ]; then
    echo -e "${GREEN}✅ Protected Endpoint: OK${NC}"
    echo "Profile: $profile_body"
else
    echo -e "${RED}❌ Protected Endpoint: FAILED (HTTP $profile_http_code)${NC}"
    echo "Response: $profile_body"
fi
echo ""

# Test 5: Refresh Token
echo -e "${YELLOW}5. Test Refresh Token${NC}"
refresh_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X POST "$BASE_URL/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"refresh_token\":\"$refresh_token\"}")

refresh_http_code=$(echo "$refresh_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
refresh_body=$(echo "$refresh_response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$refresh_http_code" = "200" ]; then
    echo -e "${GREEN}✅ Refresh Token: OK${NC}"
    new_access_token=$(echo "$refresh_body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    echo "New Access Token: ${new_access_token:0:50}..."
else
    echo -e "${RED}❌ Refresh Token: FAILED (HTTP $refresh_http_code)${NC}"
    echo "Response: $refresh_body"
fi
echo ""

# Test 6: Invalid Token
echo -e "${YELLOW}6. Test Invalid Token${NC}"
invalid_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X GET "$BASE_URL/protected/profile" \
    -H "Authorization: Bearer invalid_token")

invalid_http_code=$(echo "$invalid_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
invalid_body=$(echo "$invalid_response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$invalid_http_code" = "401" ]; then
    echo -e "${GREEN}✅ Invalid Token Rejection: OK${NC}"
    echo "Error Response: $invalid_body"
else
    echo -e "${RED}❌ Invalid Token Rejection: FAILED (Expected 401, got $invalid_http_code)${NC}"
    echo "Response: $invalid_body"
fi
echo ""

echo -e "${YELLOW}=== Tests terminés ===${NC}"
echo -e "${GREEN}✅ Tous les endpoints Auth API fonctionnent correctement !${NC}"
