# Configuration des Clés API Binance pour Postman

## 📋 Prérequis

1. **Compte Binance** : Créez un compte sur [binance.com](https://binance.com) ou utilisez le testnet
2. **Clés API** : Générez vos clés API depuis les paramètres de votre compte

## 🔑 Génération des Clés API Binance

### Pour le Testnet (Recommandé pour les tests)
1. Allez sur [testnet.binance.vision](https://testnet.binance.vision)
2. Connectez-vous avec votre compte GitHub
3. Générez une nouvelle clé API
4. Copiez votre **API Key** et **Secret Key**

### Pour le Mainnet (Production)
1. Connectez-vous sur [binance.com](https://binance.com)
2. Allez dans **Profil** > **Sécurité API**
3. Créez une nouvelle clé API
4. Configurez les permissions nécessaires (lecture/trading)
5. Copiez votre **API Key** et **Secret Key**

## ⚙️ Configuration dans Postman

### Option 1 : Variables de Collection (Recommandé)
1. Ouvrez la collection `Afrikpay Gateway - Temporal Binance Workflows (Updated)`
2. Cliquez sur l'onglet **Variables**
3. Remplacez les valeurs :
   ```
   binance_api_key: "VOTRE_CLE_API_BINANCE"
   binance_secret_key: "VOTRE_CLE_SECRETE_BINANCE"
   binance_testnet: "true" (pour testnet) ou "false" (pour mainnet)
   ```

### Option 2 : Fichier d'Environnement
1. Importez le fichier `binance_env.json` dans Postman
2. Modifiez les variables d'environnement
3. Sélectionnez l'environnement avant d'exécuter les tests

## 🚀 Configuration du Service Temporal

Le service Temporal doit être configuré avec les mêmes clés API :

### Variables d'Environnement
```bash
export BINANCE_API_KEY="votre_cle_api"
export BINANCE_SECRET_KEY="votre_cle_secrete"
export BINANCE_TESTNET="true"  # ou "false" pour mainnet
```

### Fichier .env
```env
BINANCE_API_KEY=votre_cle_api
BINANCE_SECRET_KEY=votre_cle_secrete
BINANCE_TESTNET=true
```

## 🧪 Test de la Configuration

### 1. Démarrer le Service Temporal
```bash
cd services/temporal
go run cmd/worker/main.go &
go run cmd/api/main.go &
```

### 2. Tester avec Newman CLI
```bash
cd postman/collections
newman run temporal_binance_workflows.json -e ../environments/binance_env.json
```

### 3. Endpoints Disponibles

✅ **Endpoints qui fonctionnent SANS clés API :**
- `GET /api/exchange/binance/v1/quotes` - Données publiques
- `POST /api/workflow/v1/BinancePrice` - Prix public

🔑 **Endpoints qui nécessitent des clés API :**
- `POST /api/exchange/binance/v1/order` - Placer un ordre
- `GET /api/exchange/binance/v1/order/{orderId}` - Statut d'ordre
- `GET /api/exchange/binance/v1/orders` - Tous les ordres (si implémenté avec l'API réelle)

## ⚠️ Sécurité

1. **Ne partagez jamais vos clés API** dans du code public
2. **Utilisez le testnet** pour tous les tests de développement
3. **Limitez les permissions** de vos clés API (lecture seule si possible)
4. **Régénérez vos clés** régulièrement
5. **Utilisez des variables d'environnement** ou des secrets managers en production

## 🐛 Dépannage

### Erreur "context deadline exceeded"
- Vérifiez que vos clés API sont correctes
- Vérifiez que le service Temporal est démarré
- Vérifiez la connectivité réseau avec l'API Binance

### Erreur "Signature invalid"
- Vérifiez que vos clés API et secrètes sont correctes
- Vérifiez l'horodatage système (doit être synchronisé)
- Vérifiez que vous utilisez le bon environnement (testnet vs mainnet)

### Réponses avec des prix à zéro
- Normal sur le testnet Binance - les prix peuvent être fictifs
- Vérifiez que vous utilisez des symboles supportés (BTCUSDT, ETHUSDT, etc.)

## 📚 Documentation Supplémentaire

- [Documentation API Binance](https://binance-docs.github.io/apidocs/)
- [Binance Testnet](https://testnet.binance.vision/)
- [Temporal Documentation](https://docs.temporal.io/)
