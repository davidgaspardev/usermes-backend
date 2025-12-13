#!/bin/bash

BASE_URL="http://localhost:3001"

echo "==================================="
echo "UserMes API Test - New Structure"
echo "==================================="
echo ""

# Test 1: Health check
echo "1. Health Check"
curl -s ${BASE_URL}/health | jq .
echo ""
echo ""

# Test 2: Register user
echo "2. Register User"
REGISTER_RESPONSE=$(curl -s -X POST ${BASE_URL}/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "plant.owner@company.com",
    "password": "SecurePass123!",
    "username": "plantowner",
    "name": "Plant Owner"
  }')
echo $REGISTER_RESPONSE | jq .
echo ""
echo ""

# Test 3: Login
echo "3. Login User"
LOGIN_RESPONSE=$(curl -s -X POST ${BASE_URL}/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "plantowner",
    "password": "SecurePass123!"
  }')
echo $LOGIN_RESPONSE | jq .
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token // empty')
echo ""
echo "Token: ${TOKEN:0:50}..."
echo ""
echo ""

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token. Stopping tests."
  exit 1
fi

# Test 4: Create Plant (NEW!)
echo "4. Create Plant (NEW MODULE!)"
PLANT_RESPONSE=$(curl -s -X POST ${BASE_URL}/v1/plants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "code": "SP01",
    "name": "São Paulo Manufacturing Plant",
    "latitude": -23.5505,
    "longitude": -46.6333
  }')
echo $PLANT_RESPONSE | jq .
echo ""
echo ""

# Test 5: Create Resource in Plant (UPDATED!)
echo "5. Create Resource in Plant SP01 (UPDATED WITH PLANT SCOPE!)"
RESOURCE_RESPONSE=$(curl -s -X POST ${BASE_URL}/v1/plants/SP01/production/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "code": "MACHINE-001",
    "type": "CNC-Lathe",
    "stop_factor": 10,
    "tags": ["critical", "automated"]
  }')
echo $RESOURCE_RESPONSE | jq .
echo ""
echo ""

echo "==================================="
echo "✅ All tests completed!"
echo "==================================="
