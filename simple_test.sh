#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

BASE_URL="http://localhost:3000"
EMAIL="test$(date +%s)@example.com"
PASSWORD="SecurePass123"
NAME="Test User"

echo -e "${YELLOW}=== UserMes API Test ===${NC}\n"

# Test 1: Health Check
echo -e "${YELLOW}1. Testing Health Check...${NC}"
curl -s "$BASE_URL/health" | grep -q "ok" && echo -e "${GREEN}✓ Health check passed${NC}" || echo -e "${RED}✗ Failed${NC}"
echo ""

# Test 2: Register
echo -e "${YELLOW}2. Testing User Registration...${NC}"
echo "Email: $EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"name\":\"$NAME\"}")
echo "$REGISTER_RESPONSE" | grep -q "User registered successfully" && echo -e "${GREEN}✓ Registration successful${NC}" || echo -e "${RED}✗ Failed${NC}"
echo ""

# Test 3: Login
echo -e "${YELLOW}3. Testing Login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$LOGIN_RESPONSE" | grep -q "token" && echo -e "${GREEN}✓ Login successful${NC}" || echo -e "${RED}✗ Failed${NC}"

# Extract token using grep and sed
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"//; s/"$//')
echo "Token (first 40 chars): ${TOKEN:0:40}..."
echo ""

# Test 4: Get Profile
echo -e "${YELLOW}4. Testing Get Profile...${NC}"
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/me" \
    -H "Authorization: Bearer $TOKEN")
echo "$PROFILE_RESPONSE" | grep -q "$EMAIL" && echo -e "${GREEN}✓ Get profile successful${NC}" || echo -e "${RED}✗ Failed${NC}"
echo "Profile: $PROFILE_RESPONSE"
echo ""

# Test 5: Invalid Token
echo -e "${YELLOW}5. Testing Invalid Token (should fail)...${NC}"
INVALID_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/me" \
    -H "Authorization: Bearer invalid_token")
echo "$INVALID_RESPONSE" | grep -q "unauthorized" && echo -e "${GREEN}✓ Invalid token correctly rejected${NC}" || echo -e "${RED}✗ Failed${NC}"
echo ""

echo -e "${GREEN}=== Tests Complete ===${NC}"
echo -e "Test user: ${YELLOW}$EMAIL${NC}"
