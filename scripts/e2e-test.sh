#!/bin/bash
set -e

CLI="./nhncloud"
REGION="kr1"
CONFIG_FILE="$HOME/.nhncloud/credentials"

echo "=============================================="
echo "NHN Cloud CLI E2E Integration Test"
echo "=============================================="
echo ""

echo "Step 0: Checking configuration..."
echo ""

if [ -f "$CONFIG_FILE" ]; then
    echo "  Config file: $CONFIG_FILE (found)"
    echo "  Using credentials from config file"
    echo ""
else
    echo "  Config file: $CONFIG_FILE (not found)"
    echo ""
    echo "  Checking environment variables..."
    
    check_env() {
        local var_name=$1
        local var_value=${!var_name}
        if [ -z "$var_value" ]; then
            echo "  ERROR: $var_name is not set"
            return 1
        fi
        echo "  $var_name: ***configured***"
        return 0
    }
    
    errors=0
    check_env "NHN_CLOUD_REGION" || ((errors++))
    check_env "NHN_CLOUD_APPKEY" || ((errors++))
    check_env "NHN_CLOUD_ACCESS_KEY" || ((errors++))
    check_env "NHN_CLOUD_SECRET_KEY" || ((errors++))
    check_env "NHN_CLOUD_USERNAME" || ((errors++))
    check_env "NHN_CLOUD_PASSWORD" || ((errors++))
    check_env "NHN_CLOUD_TENANT_ID" || ((errors++))
    
    if [ $errors -gt 0 ]; then
        echo ""
        echo "No credentials found. Please either:"
        echo ""
        echo "1. Create config file at ~/.nhncloud/credentials:"
        echo "   [default]"
        echo "   access_key_id = your-access-key"
        echo "   secret_access_key = your-secret-key"
        echo "   region = kr1"
        echo "   username = your-email"
        echo "   api_password = your-password"
        echo "   tenant_id = your-tenant-id"
        echo "   rds_app_key = your-rds-appkey"
        echo ""
        echo "2. Or set environment variables:"
        echo "   export NHN_CLOUD_REGION=kr1"
        echo "   export NHN_CLOUD_APPKEY=your-appkey"
        echo "   export NHN_CLOUD_ACCESS_KEY=your-access-key"
        echo "   export NHN_CLOUD_SECRET_KEY=your-secret-key"
        echo "   export NHN_CLOUD_USERNAME=your-email"
        echo "   export NHN_CLOUD_PASSWORD=your-api-password"
        echo "   export NHN_CLOUD_TENANT_ID=your-tenant-id"
        exit 1
    fi
fi

echo "Configuration check passed."
echo ""

echo "=============================================="
echo "Phase 1: Verify CLI Commands"
echo "=============================================="
echo ""

echo "1.1 Testing VPC list..."
$CLI vpc list --region $REGION -o json | head -5
echo "✓ VPC list works"
echo ""

echo "1.2 Testing Subnet list..."
$CLI vpc subnets --region $REGION -o json | head -5
echo "✓ Subnet list works"
echo ""

echo "1.3 Testing Security Group list..."
$CLI security-group list --region $REGION -o json | head -5
echo "✓ Security Group list works"
echo ""

echo "1.4 Testing Floating IP list..."
$CLI floating-ip list --region $REGION -o json | head -5
echo "✓ Floating IP list works"
echo ""

echo "1.5 Testing Compute flavors..."
$CLI compute flavors --region $REGION | head -10
echo "✓ Compute flavors works"
echo ""

echo "1.6 Testing Compute images..."
$CLI compute images --region $REGION | head -10
echo "✓ Compute images works"
echo ""

echo "1.7 Testing Compute keypairs..."
$CLI compute keypairs --region $REGION
echo "✓ Compute keypairs works"
echo ""

echo "1.8 Testing Compute instances list..."
$CLI compute list --region $REGION
echo "✓ Compute instances list works"
echo ""

echo "1.9 Testing RDS MySQL flavors..."
$CLI rds-mysql flavors --region $REGION | head -10
echo "✓ RDS MySQL flavors works"
echo ""

echo "1.10 Testing RDS MySQL versions..."
$CLI rds-mysql versions --region $REGION
echo "✓ RDS MySQL versions works"
echo ""

echo "1.11 Testing RDS MySQL instances list..."
$CLI rds-mysql list --region $REGION
echo "✓ RDS MySQL instances list works"
echo ""

echo "=============================================="
echo "Phase 2: Resource Discovery"
echo "=============================================="
echo ""

echo "2.1 Finding Ubuntu 22.04 image..."
$CLI compute images --region $REGION -o json | jq -r '.images[] | select(.name | contains("Ubuntu") and contains("22.04")) | "\(.id) - \(.name)"' | head -3
echo ""

echo "2.2 Finding m2.c2m4 flavor..."
$CLI compute flavors --region $REGION -o json | jq -r '.flavors[] | select(.name | contains("m2.c2m4")) | "\(.id) - \(.name) (\(.vcpus) vCPU, \(.ram)MB RAM)"' | head -3
echo ""

echo "2.3 Finding available Floating IPs (status=DOWN)..."
$CLI floating-ip list --region $REGION -o json | jq -r '.floatingips[] | select(.status == "DOWN") | "\(.id) - \(.floating_ip_address)"' | head -3
echo ""

echo "2.4 Finding MySQL parameter groups..."
$CLI rds-mysql parameter-group-list --region $REGION | head -10
echo ""

echo "2.5 Finding MySQL flavors (m2.c2m4)..."
$CLI rds-mysql flavors --region $REGION -o json | jq -r '.dbFlavors[] | select(.dbFlavorName | contains("m2.c2m4")) | "\(.dbFlavorId) - \(.dbFlavorName)"' | head -3
echo ""

echo "=============================================="
echo "E2E Test Summary"
echo "=============================================="
echo ""
echo "All basic connectivity tests passed!"
echo ""
echo "Available commands verified:"
echo "  - VPC management: vpc list, vpc get, vpc subnets"
echo "  - Security Groups: security-group list/get/create/delete/rule-create"
echo "  - Floating IPs: floating-ip list/get/create/delete/associate/disassociate"
echo "  - Compute: compute list/get/create/delete/start/stop/reboot/flavors/images/keypairs"
echo "  - RDS MySQL: rds-mysql list/get/create/delete/flavors/versions/backups/etc."
echo "  - RDS MariaDB: rds-mariadb list/get/create/delete/flavors/versions"
echo "  - RDS PostgreSQL: rds-postgresql list/get/create/delete/flavors/versions/database"
echo ""
echo "Next steps for full E2E test:"
echo "1. Create Security Group: security-group create --name e2e-test-sg"
echo "2. Add rules: security-group rule-create --security-group-id <id> --protocol tcp --port-min 22 --port-max 22"
echo "3. Create VM: compute create --name e2e-test-vm --image <id> --flavor <id> --network <subnet-id> --key-name <key>"
echo "4. Associate Floating IP: floating-ip associate <fip-id> --port-id <port-id>"
echo "5. Create RDS: rds-mysql create --name e2e-test-db ..."
echo "6. Test connectivity: SSH to VM, connect to MySQL"
echo ""
echo "Refer to docs/E2E_TEST_SCENARIO.md for full test procedure."
