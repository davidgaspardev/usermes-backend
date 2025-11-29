#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:3000"

# Test variables
EMAIL="test$(date +%s)@example.com"
PASSWORD="SecurePass123"
NAME="Test User"
TOKEN=""
USER_ID=""

echo -e "${YELLOW}=== UserMes API Test Script ===${NC}\n"

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Health check passed"
    echo "Response: $BODY"
else
    print_result 1 "Health check failed (HTTP $HTTP_CODE)"
fi
echo ""

# Test 2: Register User
echo -e "${YELLOW}Test 2: Register User${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"name\":\"$NAME\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 201 ]; then
    print_result 0 "User registration successful"
    USER_ID=$(echo "$BODY" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "User ID: $USER_ID"
    echo "Email: $EMAIL"
else
    print_result 1 "User registration failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 3: Login
echo -e "${YELLOW}Test 3: Login${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Login successful"
    TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "Token: ${TOKEN:0:50}..."
else
    print_result 1 "Login failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 4: Get Profile
echo -e "${YELLOW}Test 4: Get User Profile (Me)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/users/me" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Get profile successful"
    echo "Response: $BODY"
else
    print_result 1 "Get profile failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 5: Get User by ID
echo -e "${YELLOW}Test 5: Get User by ID${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/users/$USER_ID" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Get user by ID successful"
    echo "Response: $BODY"
else
    print_result 1 "Get user by ID failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 6: Update User
echo -e "${YELLOW}Test 6: Update User Name${NC}"
NEW_NAME="Updated Test User"
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/users/$USER_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"name\":\"$NEW_NAME\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Update user successful"
    echo "New name: $NEW_NAME"
else
    print_result 1 "Update user failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 7: Change Password
echo -e "${YELLOW}Test 7: Change Password${NC}"
NEW_PASSWORD="NewSecurePass456"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/$USER_ID/change-password" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"old_password\":\"$PASSWORD\",\"new_password\":\"$NEW_PASSWORD\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Change password successful"
    PASSWORD=$NEW_PASSWORD
else
    print_result 1 "Change password failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 8: Login with new password
echo -e "${YELLOW}Test 8: Login with New Password${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$NEW_PASSWORD\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Login with new password successful"
else
    print_result 1 "Login with new password failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 9: Deactivate User
echo -e "${YELLOW}Test 9: Deactivate User${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/$USER_ID/deactivate" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Deactivate user successful"
else
    print_result 1 "Deactivate user failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 10: Try login with deactivated user (should fail)
echo -e "${YELLOW}Test 10: Login with Deactivated User (Should Fail)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$NEW_PASSWORD\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -eq 403 ] || [ "$HTTP_CODE" -eq 401 ]; then
    print_result 0 "Login correctly blocked for deactivated user"
else
    print_result 1 "Login should have failed for deactivated user (HTTP $HTTP_CODE)"
fi
echo ""

# Test 11: Activate User
echo -e "${YELLOW}Test 11: Activate User${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/$USER_ID/activate" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Activate user successful"
else
    print_result 1 "Activate user failed (HTTP $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 12: Login after reactivation
echo -e "${YELLOW}Test 12: Login After Reactivation${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$NEW_PASSWORD\"}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Login after reactivation successful"
else
    print_result 1 "Login after reactivation failed (HTTP $HTTP_CODE)"
fi
echo ""

# Test 13: Invalid token
echo -e "${YELLOW}Test 13: Access with Invalid Token (Should Fail)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/users/me" \
    -H "Authorization: Bearer invalid_token_123")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -eq 401 ]; then
    print_result 0 "Invalid token correctly rejected"
else
    print_result 1 "Invalid token should have been rejected (HTTP $HTTP_CODE)"
fi
echo ""

# Summary
echo -e "${GREEN}=== All Tests Passed! ===${NC}"
echo ""
echo "Summary:"
echo "- User registered: $EMAIL"
echo "- User ID: $USER_ID"
echo "- All endpoints working correctly"
echo "- Authentication and authorization working"
echo ""
echo -e "${YELLOW}Note: The user created during testing remains in memory${NC}"
echo -e "${YELLOW}Restart the server to clear all data${NC}"
