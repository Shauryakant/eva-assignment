#!/usr/bin/env bash

# Terminal colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"
EMAIL="demo_user_$(date +%s)@example.com"
PASSWORD="password123"

echo -e "${CYAN}====================================================${NC}"
echo -e "${CYAN}    TICKET SYSTEM API - AUTOMATED DEMO TEST RUN      ${NC}"
echo -e "${CYAN}====================================================${NC}"
sleep 1.5

# 1. Health Check
echo -e "\n${YELLOW}[TEST 1/8] GET /health (Public Health Check)${NC}"
curl -s -X GET "$BASE_URL/health" | jq . || curl -s -X GET "$BASE_URL/health"
sleep 2

# 2. Register User
echo -e "\n${YELLOW}[TEST 2/8] POST /auth/register (Register New User)${NC}"
echo -e "${BLUE}Registering user: $EMAIL${NC}"
REGISTER_RESP=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$REGISTER_RESP"
sleep 2

# 3. Login
echo -e "\n${YELLOW}[TEST 3/8] POST /auth/login (Login & Obtain JWT)${NC}"
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$LOGIN_RESP"

TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo -e "${RED}Failed to obtain JWT token! Is the server running on port 8080?${NC}"
  exit 1
fi

echo -e "${GREEN}JWT Token obtained successfully!${NC}"
sleep 2

# 4. Create Ticket
echo -e "\n${YELLOW}[TEST 4/8] POST /tickets (Create Ticket Owned by Caller)${NC}"
CREATE_RESP=$(curl -s -X POST "$BASE_URL/tickets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Database Connection Spike","description":"High latency detected on Atlas cluster"}')
echo "$CREATE_RESP"

TICKET_ID=$(echo "$CREATE_RESP" | grep -o '"id":"[^"]*' | grep -o '[^"]*$')
sleep 2

# 5. List Tickets
echo -e "\n${YELLOW}[TEST 5/8] GET /tickets (List User's Tickets)${NC}"
curl -s -X GET "$BASE_URL/tickets" \
  -H "Authorization: Bearer $TOKEN"
sleep 2

# 6. Fetch Ticket by ID
echo -e "\n${YELLOW}[TEST 6/8] GET /tickets/$TICKET_ID (Fetch Single Ticket)${NC}"
curl -s -X GET "$BASE_URL/tickets/$TICKET_ID" \
  -H "Authorization: Bearer $TOKEN"
sleep 2

# 7. Valid Status Transition (open -> in_progress)
echo -e "\n${YELLOW}[TEST 7/8] PATCH /tickets/$TICKET_ID/status (Valid Transition: open -> in_progress)${NC}"
curl -s -X PATCH "$BASE_URL/tickets/$TICKET_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status":"in_progress"}'
sleep 2

# 8. Illegal Status Transition (in_progress -> open -> Expected 400 Bad Request)
echo -e "\n${YELLOW}[TEST 8/8] PATCH /tickets/$TICKET_ID/status (Invalid Transition: in_progress -> open)${NC}"
curl -s -X PATCH "$BASE_URL/tickets/$TICKET_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status":"open"}'
sleep 2

echo -e "\n${GREEN}====================================================${NC}"
echo -e "${GREEN}    DEMO COMPLETED SUCCESSFULLY!                     ${NC}"
echo -e "${GREEN}====================================================${NC}"
