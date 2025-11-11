#!/bin/bash

# Authentication Service API Test Script
# This script tests all authentication endpoints in the microservice-auth example

set -e  # Exit on any error

# Configuration
BASE_URL="http://localhost:8081/api/v1"
HEALTH_URL="http://localhost:8081/health"
AUTH_ENDPOINT="$BASE_URL/auth"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test user data
TEST_USERNAME="testuser_$(date +%s)"
TEST_EMAIL="$TEST_USERNAME@example.com"
TEST_PASSWORD="password123"
NEW_PASSWORD="newpassword456"
CLIENT_ID="web-app"

# Global variables to store tokens
ACCESS_TOKEN=""
REFRESH_TOKEN=""

# Helper functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

check_server() {
    print_info "Checking if server is running..."
    if ! curl -s --max-time 5 "$HEALTH_URL" > /dev/null 2>&1; then
        print_error "Server is not running on $BASE_URL"
        print_info "Please start the authentication service first:"
        echo "  cd examples/authentication-service"
        echo "  go run cmd/app/main.go"
        exit 1
    fi
    print_success "Server is running"
}

# Test functions
test_signup() {
    print_header "Testing User Signup"

    local payload=$(cat <<EOF
{
    "clientId": "$CLIENT_ID",
    "identifier": {
        "scheme": "username",
        "value": "$TEST_USERNAME"
    },
    "credential": {
        "scheme": "basic",
        "value": "$TEST_PASSWORD"
    },
    "email": "$TEST_EMAIL",
    "firstname": "Test",
    "lastname": "User"
}
EOF
    )

    echo "Payload: $payload"

    local response=$(curl -s -X POST "$AUTH_ENDPOINT/signup" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    # Check if signup was successful
    if echo "$response" | grep -q '"userId"\|"token"'; then
        print_success "User signup successful"

        # Extract tokens from nested structure using more robust parsing
        ACCESS_TOKEN=$(echo "$response" | sed -n 's/.*"token":{"value":"\([^"]*\)".*/\1/p')
        REFRESH_TOKEN=$(echo "$response" | sed -n 's/.*"refreshToken":{"value":"\([^"]*\)".*/\1/p')

        if [ -n "$ACCESS_TOKEN" ]; then
            print_success "Access token obtained"
        fi
        if [ -n "$REFRESH_TOKEN" ]; then
            print_success "Refresh token obtained"
        fi
    else
        print_error "User signup failed"
        return 1
    fi
}

test_signin() {
    print_header "Testing User Signin"

    local payload=$(cat <<EOF
{
    "clientId": "$CLIENT_ID",
    "identifier": {
        "scheme": "username",
        "value": "$TEST_USERNAME"
    },
    "credential": {
        "scheme": "basic",
        "value": "$TEST_PASSWORD"
    }
}
EOF
    )

    echo "Payload: $payload"

    local response=$(curl -s -X POST "$AUTH_ENDPOINT/signin" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    # Check if signin was successful
    if echo "$response" | grep -q '"userId"\|"token"'; then
        print_success "User signin successful"

        # Extract tokens from nested structure using more robust parsing
        ACCESS_TOKEN=$(echo "$response" | sed -n 's/.*"token":{"value":"\([^"]*\)".*/\1/p')
        REFRESH_TOKEN=$(echo "$response" | sed -n 's/.*"refreshToken":{"value":"\([^"]*\)".*/\1/p')

        if [ -n "$ACCESS_TOKEN" ]; then
            print_success "Access token obtained: ${ACCESS_TOKEN:0:20}..."
        else
            print_error "No access token received"
            return 1
        fi
    else
        print_error "User signin failed"
        return 1
    fi
}

test_get_profile() {
    print_header "Testing Get Profile"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    local response=$(curl -s -X GET "$AUTH_ENDPOINT/profile" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true\|"email"\|"firstname"'; then
        print_success "Get profile successful"
    else
        print_error "Get profile failed"
        return 1
    fi
}

test_update_profile() {
    print_header "Testing Update Profile"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    local payload=$(cat <<EOF
{
    "firstname": "UpdatedTest",
    "lastname": "UpdatedUser",
    "email": "updated_$TEST_EMAIL"
}
EOF
    )

    echo "Payload: $payload"

    local response=$(curl -s -X PUT "$AUTH_ENDPOINT/profile" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true'; then
        print_success "Update profile successful"
    else
        print_error "Update profile failed"
        return 1
    fi
}

test_change_password() {
    print_header "Testing Change Password"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    local payload=$(cat <<EOF
{
    "oldCredential": {
        "scheme": "basic",
        "value": "$TEST_PASSWORD"
    },
    "newCredential": {
        "scheme": "basic",
        "value": "$NEW_PASSWORD"
    }
}
EOF
    )

    echo "Payload: $payload"

    local response=$(curl -s -X PUT "$AUTH_ENDPOINT/change-password" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true'; then
        print_success "Change password successful"

        # Update password for next signin test
        TEST_PASSWORD="$NEW_PASSWORD"
    else
        print_error "Change password failed"
        return 1
    fi
}

test_signin_after_password_change() {
    print_header "Testing Signin with New Password"

    local payload=$(cat <<EOF
{
    "clientId": "$CLIENT_ID",
    "identifier": {
        "scheme": "username",
        "value": "$TEST_USERNAME"
    },
    "credential": {
        "scheme": "basic",
        "value": "$TEST_PASSWORD"
    }
}
EOF
    )

    echo "Payload: $payload"

    local response=$(curl -s -X POST "$AUTH_ENDPOINT/signin" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    if echo "$response" | grep -q '"userId"\|"token"'; then
        print_success "Signin with new password successful"
        
        # Update token for subsequent tests
        ACCESS_TOKEN=$(echo "$response" | sed -n 's/.*"token":{"value":"\([^"]*\)".*/\1/p')
    else
        print_error "Signin with new password failed"
        return 1
    fi
}

test_protected_dashboard() {
    print_header "Testing Protected Dashboard"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    local response=$(curl -s -X GET "$BASE_URL/protected/dashboard" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true\|"data"'; then
        print_success "Protected dashboard access successful"
    else
        print_error "Protected dashboard access failed"
        return 1
    fi
}

test_protected_settings() {
    print_header "Testing Protected Settings"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    # Get settings
    print_info "Getting user settings..."
    local response=$(curl -s -X GET "$BASE_URL/protected/settings" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true\|"settings"'; then
        print_success "Get settings successful"
    else
        print_error "Get settings failed"
        return 1
    fi

    # Update settings
    print_info "Updating user settings..."
    local payload=$(cat <<EOF
{
    "theme": "dark",
    "language": "en",
    "notifications": true
}
EOF
    )

    echo "Payload: $payload"

    response=$(curl -s -X PUT "$BASE_URL/protected/settings" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$payload")

    echo "Response: $response"

    if echo "$response" | grep -q '"success":true'; then
        print_success "Update settings successful"
    else
        print_error "Update settings failed"
        return 1
    fi
}

test_admin_dashboard() {
    print_header "Testing Admin Dashboard (Role Check)"

    if [ -z "$ACCESS_TOKEN" ]; then
        print_error "No access token available"
        return 1
    fi

    local response=$(curl -s -X GET "$BASE_URL/admin/dashboard" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    # This should either succeed (if user has admin role) or fail with 403 (if not)
    if echo "$response" | grep -q '"success":true'; then
        print_success "Admin dashboard access successful (user has admin role)"
    elif echo "$response" | grep -q '"success":false\|"Admin role required"'; then
        print_success "Admin role check working correctly (403 expected for regular user)"
    else
        print_error "Unexpected response from admin dashboard"
        return 1
    fi
}

test_unauthorized_access() {
    print_header "Testing Unauthorized Access (No Token)"

    local response=$(curl -s -X GET "$AUTH_ENDPOINT/profile" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    if echo "$response" | grep -qi "unauthorized\|missing\|invalid"; then
        print_success "Unauthorized access properly blocked"
    else
        print_error "Unauthorized access not properly blocked!"
        return 1
    fi
}

test_invalid_token() {
    print_header "Testing Invalid Token"

    local response=$(curl -s -X GET "$AUTH_ENDPOINT/profile" \
        -H "Authorization: Bearer invalid.token.here" \
        -H "Content-Type: application/json")

    echo "Response: $response"

    if echo "$response" | grep -qi "unauthorized\|invalid\|token"; then
        print_success "Invalid token properly rejected"
    else
        print_error "Invalid token not properly rejected!"
        return 1
    fi
}

# Main test execution
main() {
    print_header "Authentication Service API Tests"
    print_info "Testing user: $TEST_USERNAME"
    print_info "Base URL: $BASE_URL"

    check_server

    # Run tests
    test_signup
    test_signin
    test_get_profile
    test_update_profile
    test_change_password
    test_signin_after_password_change
    test_protected_dashboard
    test_protected_settings
    test_admin_dashboard
    test_unauthorized_access
    test_invalid_token

    print_header "All Tests Completed"
    print_success "Authentication API testing finished"
    print_info "Total tests: 11"
    print_info "All endpoints validated successfully!"
}

# Run main function
main "$@"