#!/bin/bash

# Script to test third-party APIs using Newman
# This validates that our client implementations match the real API behaviors

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COLLECTION_FILE="docs/third_party_apis.postman_collection.json"
ENVIRONMENT_FILE="docs/third_party_apis.postman_environment.json"
REPORTS_DIR="reports/api_tests"

echo -e "${BLUE}=== Afrikpay Third Party APIs Validation ===${NC}"
echo "Testing direct API calls to validate our client implementations"
echo

# Check if Newman is installed
if ! command -v newman &> /dev/null; then
    echo -e "${RED}Error: Newman is not installed${NC}"
    echo "Install Newman with: npm install -g newman"
    exit 1
fi

# Create reports directory
mkdir -p "$REPORTS_DIR"

# Check if collection file exists
if [ ! -f "$COLLECTION_FILE" ]; then
    echo -e "${RED}Error: Collection file not found: $COLLECTION_FILE${NC}"
    exit 1
fi

# Check if environment file exists
if [ ! -f "$ENVIRONMENT_FILE" ]; then
    echo -e "${RED}Error: Environment file not found: $ENVIRONMENT_FILE${NC}"
    exit 1
fi

# Function to run tests for a specific folder
run_folder_tests() {
    local folder_name="$1"
    local report_name="$2"
    
    echo -e "${YELLOW}Testing $folder_name APIs...${NC}"
    
    newman run "$COLLECTION_FILE" \
        --environment "$ENVIRONMENT_FILE" \
        --folder "$folder_name" \
        --reporters cli,json,html \
        --reporter-json-export "$REPORTS_DIR/${report_name}_results.json" \
        --reporter-html-export "$REPORTS_DIR/${report_name}_report.html" \
        --timeout-request 10000 \
        --delay-request 1000 \
        --bail \
        || echo -e "${RED}$folder_name tests failed${NC}"
    
    echo
}

# Function to run all tests
run_all_tests() {
    echo -e "${YELLOW}Running all API tests...${NC}"
    
    newman run "$COLLECTION_FILE" \
        --environment "$ENVIRONMENT_FILE" \
        --reporters cli,json,html \
        --reporter-json-export "$REPORTS_DIR/all_apis_results.json" \
        --reporter-html-export "$REPORTS_DIR/all_apis_report.html" \
        --timeout-request 10000 \
        --delay-request 1000 \
        || echo -e "${RED}Some tests failed - check reports for details${NC}"
    
    echo
}

# Check command line arguments
case "${1:-all}" in
    "binance")
        run_folder_tests "Binance API Tests" "binance"
        ;;
    "bitget")
        run_folder_tests "Bitget API Tests" "bitget"
        ;;
    "mtn")
        run_folder_tests "MTN Mobile Money API Tests" "mtn"
        ;;
    "orange")
        run_folder_tests "Orange Money API Tests" "orange"
        ;;
    "all")
        run_all_tests
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [binance|bitget|mtn|orange|all|help]"
        echo
        echo "Commands:"
        echo "  binance  - Test only Binance APIs"
        echo "  bitget   - Test only Bitget APIs"
        echo "  mtn      - Test only MTN Mobile Money APIs"
        echo "  orange   - Test only Orange Money APIs"
        echo "  all      - Test all APIs (default)"
        echo "  help     - Show this help message"
        echo
        echo "Environment Variables (set in environment file):"
        echo "  BINANCE_API_KEY      - Binance API key"
        echo "  BITGET_API_KEY       - Bitget API key"
        echo "  BITGET_SECRET_KEY    - Bitget secret key"
        echo "  BITGET_PASSPHRASE    - Bitget passphrase"
        echo "  MTN_API_KEY          - MTN Mobile Money API key"
        echo "  MTN_SUBSCRIPTION_KEY - MTN subscription key"
        echo "  ORANGE_API_KEY       - Orange Money API key"
        exit 0
        ;;
    *)
        echo -e "${RED}Error: Unknown command '$1'${NC}"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac

echo -e "${GREEN}=== API Testing Complete ===${NC}"
echo "Reports generated in: $REPORTS_DIR"
echo
echo "Next steps:"
echo "1. Review the HTML reports for detailed results"
echo "2. Compare API responses with our client implementations"
echo "3. Update client code if any discrepancies are found"
echo "4. Run integration tests with the Client service"
