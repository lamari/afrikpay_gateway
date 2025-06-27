#!/bin/bash

# Script pour synchroniser les variables d'environnement avec l'environnement Postman
# Lit les clés API depuis .env.api_keys et met à jour le fichier Postman environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENV_FILE="docs/.env.api_keys"
POSTMAN_ENV_FILE="docs/third_party_apis.postman_environment.json"
POSTMAN_ENV_TEMPLATE="docs/third_party_apis.postman_environment.json"

echo -e "${BLUE}=== Synchronisation Environnement Postman ===${NC}"
echo "Mise à jour des variables API depuis .env.api_keys"
echo

# Vérifier si le fichier .env.api_keys existe
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}Erreur: Fichier .env.api_keys non trouvé${NC}"
    echo "Créez le fichier depuis le template :"
    echo "  cp docs/.env.api_keys.example docs/.env.api_keys"
    echo "  # Puis éditez avec vos vraies clés API"
    exit 1
fi

# Vérifier si le fichier Postman environment existe
if [ ! -f "$POSTMAN_ENV_FILE" ]; then
    echo -e "${RED}Erreur: Fichier environnement Postman non trouvé${NC}"
    echo "Fichier attendu : $POSTMAN_ENV_FILE"
    exit 1
fi

# Charger les variables d'environnement
echo -e "${YELLOW}Chargement des variables depuis .env.api_keys...${NC}"
source "$ENV_FILE"

# Fonction pour mettre à jour une variable dans le JSON Postman
update_postman_var() {
    local var_name="$1"
    local var_value="$2"
    local temp_file=$(mktemp)
    
    if [ -n "$var_value" ]; then
        echo "  Mise à jour : $var_name"
        jq --arg name "$var_name" --arg value "$var_value" \
           '(.values[] | select(.key == $name) | .value) = $value' \
           "$POSTMAN_ENV_FILE" > "$temp_file" && mv "$temp_file" "$POSTMAN_ENV_FILE"
    else
        echo "  Ignoré (vide) : $var_name"
    fi
}

# Vérifier si jq est installé
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Erreur: jq n'est pas installé${NC}"
    echo "Installez jq avec : brew install jq (macOS) ou apt-get install jq (Ubuntu)"
    exit 1
fi

echo -e "${YELLOW}Mise à jour des variables Postman...${NC}"

# Mettre à jour les variables Binance
update_postman_var "binance_api_key" "$BINANCE_API_KEY"
update_postman_var "binance_base_url" "${BINANCE_BASE_URL:-https://testnet.binance.vision}"

# Mettre à jour les variables Bitget
update_postman_var "bitget_api_key" "$BITGET_API_KEY"
update_postman_var "bitget_secret_key" "$BITGET_SECRET_KEY"
update_postman_var "bitget_passphrase" "$BITGET_PASSPHRASE"
update_postman_var "bitget_base_url" "${BITGET_BASE_URL:-https://api.bitget.com}"

# Mettre à jour les variables MTN
update_postman_var "mtn_api_key" "$MTN_API_KEY"
update_postman_var "mtn_subscription_key" "$MTN_SUBSCRIPTION_KEY"
update_postman_var "mtn_sandbox_url" "${MTN_SANDBOX_URL:-https://sandbox.momodeveloper.mtn.com}"

# Mettre à jour les variables Orange
update_postman_var "orange_api_key" "$ORANGE_API_KEY"
update_postman_var "orange_sandbox_url" "${ORANGE_SANDBOX_URL:-https://api.orange.com}"

echo
echo -e "${GREEN}✅ Synchronisation terminée${NC}"
echo "Fichier mis à jour : $POSTMAN_ENV_FILE"
echo

# Vérifier les variables critiques
echo -e "${YELLOW}Vérification des variables critiques...${NC}"

missing_vars=()

[ -z "$BINANCE_API_KEY" ] && missing_vars+=("BINANCE_API_KEY")
[ -z "$BITGET_API_KEY" ] && missing_vars+=("BITGET_API_KEY")
[ -z "$BITGET_SECRET_KEY" ] && missing_vars+=("BITGET_SECRET_KEY")
[ -z "$BITGET_PASSPHRASE" ] && missing_vars+=("BITGET_PASSPHRASE")
[ -z "$MTN_API_KEY" ] && missing_vars+=("MTN_API_KEY")
[ -z "$MTN_SUBSCRIPTION_KEY" ] && missing_vars+=("MTN_SUBSCRIPTION_KEY")
[ -z "$ORANGE_API_KEY" ] && missing_vars+=("ORANGE_API_KEY")

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Variables manquantes :${NC}"
    for var in "${missing_vars[@]}"; do
        echo "  - $var"
    done
    echo
    echo "Les tests correspondants échoueront. Configurez ces variables dans .env.api_keys"
else
    echo -e "${GREEN}✅ Toutes les variables critiques sont configurées${NC}"
fi

echo
echo -e "${BLUE}Prochaines étapes :${NC}"
echo "1. Vérifiez le fichier Postman environment mis à jour"
echo "2. Importez-le dans Postman ou utilisez avec Newman"
echo "3. Exécutez les tests : ./test_third_party_apis.sh"
echo
echo -e "${YELLOW}Commandes utiles :${NC}"
echo "  # Tester avec Newman"
echo "  newman run docs/third_party_apis.postman_collection.json \\"
echo "    --environment docs/third_party_apis.postman_environment.json"
echo
echo "  # Tester un service spécifique"
echo "  ./test_third_party_apis.sh binance"
