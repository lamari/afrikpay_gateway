#!/bin/bash

# Script pour lancer le service auth en local
# Usage: ./run_local.sh

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Démarrage du service Afrikpay Auth API en local ===${NC}"

# Vérifier que les clés JWT existent
if [ ! -f "../../config/keys/private.pem" ] || [ ! -f "../../config/keys/public.pem" ]; then
    echo "Création des clés JWT..."
    mkdir -p ../../config/keys
    openssl genrsa -out ../../config/keys/private.pem 2048
    openssl rsa -in ../../config/keys/private.pem -pubout -out ../../config/keys/public.pem
fi

# Configurer les variables d'environnement pour le développement local
export AUTH_PORT=8001
export AUTH_HOST="0.0.0.0"
export JWT_PRIVATE_KEY_PATH="../../config/keys/private.pem"
export JWT_PUBLIC_KEY_PATH="../../config/keys/public.pem"
export LOG_LEVEL="debug"

echo -e "${GREEN}Lancement du service auth sur http://localhost:8001${NC}"
echo "Utilisez Ctrl+C pour arrêter le service"

# Lancer le service
cd cmd && go run main.go
