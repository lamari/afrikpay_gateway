# Guide de Validation des APIs Tierces - Afrikpay Gateway

Ce guide explique comment utiliser la collection Postman pour valider que nos implémentations de clients correspondent exactement aux comportements des APIs tierces réelles.

## 📋 Vue d'ensemble

La collection `third_party_apis.postman_collection.json` contient des tests directs pour :

- **Binance API** - Exchange crypto (prix, ordres)
- **Bitget API** - Exchange crypto (prix, ordres) 
- **MTN Mobile Money API** - Paiements mobile (initiation, statut)
- **Orange Money API** - Paiements mobile (initiation, statut)

## 🎯 Objectif

Valider que nos clients Go implémentent correctement :
- Les endpoints exacts utilisés par les APIs
- Les formats de requête et réponse attendus
- Les codes de statut HTTP appropriés
- Les headers d'authentification requis
- Les structures de données JSON

## 🚀 Installation et Configuration

### 1. Installer Newman (CLI Postman)

```bash
npm install -g newman
```

### 2. Configurer les Clés API

Éditez le fichier `docs/third_party_apis.postman_environment.json` et ajoutez vos clés API :

```json
{
  "key": "binance_api_key",
  "value": "VOTRE_CLE_BINANCE",
  "type": "secret"
}
```

**⚠️ Important** : Utilisez les environnements sandbox/testnet pour éviter les transactions réelles.

### 3. URLs des Environnements

- **Binance Testnet** : `https://testnet.binance.vision`
- **Bitget Sandbox** : `https://api.bitget.com` (avec clés de test)
- **MTN Sandbox** : `https://sandbox.momodeveloper.mtn.com`
- **Orange Sandbox** : `https://api.orange.com` (environnement de test)

## 🧪 Exécution des Tests

### Via Script Automatisé

```bash
# Tous les tests
./scripts/test_third_party_apis.sh

# Tests spécifiques
./scripts/test_third_party_apis.sh binance
./scripts/test_third_party_apis.sh bitget
./scripts/test_third_party_apis.sh mtn
./scripts/test_third_party_apis.sh orange

# Aide
./scripts/test_third_party_apis.sh help
```

### Via Newman Direct

```bash
# Tous les tests
newman run docs/third_party_apis.postman_collection.json \
  --environment docs/third_party_apis.postman_environment.json

# Tests d'un dossier spécifique
newman run docs/third_party_apis.postman_collection.json \
  --environment docs/third_party_apis.postman_environment.json \
  --folder "Binance API Tests"
```

### Via Postman GUI

1. Importer `third_party_apis.postman_collection.json`
2. Importer `third_party_apis.postman_environment.json`
3. Sélectionner l'environnement
4. Configurer les variables avec vos clés API
5. Exécuter les tests individuellement ou en collection

## 📊 Tests Inclus

### Binance API Tests

| Test | Endpoint | Validation |
|------|----------|------------|
| Get Price BTCUSDT | `/api/v3/ticker/price` | Format prix string, symbol correct |
| Get Price ETHUSDT | `/api/v3/ticker/price` | Format prix string, symbol correct |
| Get All Prices | `/api/v3/ticker/24hr` | Array de tickers, champs requis |

**Validations clés** :
- Prix retourné comme string (notre client utilise `strconv.ParseFloat`)
- Header `X-MBX-APIKEY` requis
- Structure de réponse : `{"symbol": "BTCUSDT", "price": "43250.50"}`

### Bitget API Tests

| Test | Endpoint | Validation |
|------|----------|------------|
| Get Ticker BTCUSDT | `/api/spot/v1/market/ticker` | Structure response avec code/data |
| Get All Tickers | `/api/spot/v1/market/tickers` | Array dans data, champs lastPr/baseVolume |

**Validations clés** :
- Réponse encapsulée : `{"code": "00000", "data": {...}}`
- Prix dans `lastPr` (notre client mappe vers `LastPrice`)
- Headers d'authentification complexes (ACCESS-KEY, ACCESS-SIGN, etc.)

### MTN Mobile Money API Tests

| Test | Endpoint | Validation |
|------|----------|------------|
| Get Balance | `/collection/v1_0/accountbalance` | Balance string, currency |
| Request to Pay | `/collection/v1_0/requesttopay` | Status 202, X-Reference-Id header |
| Get Payment Status | `/collection/v1_0/requesttopay/{id}` | Status PENDING/SUCCESSFUL/FAILED |

**Validations clés** :
- Balance comme string (notre client attend string)
- Status 202 pour initiation (pas 200)
- Champ `reason` pour messages d'erreur (pas `message`)
- Headers : `Authorization`, `X-Target-Environment`, `Ocp-Apim-Subscription-Key`

### Orange Money API Tests

| Test | Endpoint | Validation |
|------|----------|------------|
| Get Balance | `/omcoreapis/1.0.2/mp/balance` | Balance number, currency |
| Initiate Payment | `/omcoreapis/1.0.2/mp/pay` | Status 201, transactionId |
| Get Payment Status | `/omcoreapis/1.0.2/mp/status/{id}` | Status SUCCESS/PENDING/FAILED |

**Validations clés** :
- Balance comme number (notre client attend float64)
- Status 201 pour initiation (pas 202 comme MTN)
- Champ `message` pour messages (pas `reason` comme MTN)
- `transactionId` dans réponse (notre client mappe vers `ReferenceID`)

## 🔍 Analyse des Résultats

### Rapports Générés

Les tests génèrent des rapports dans `reports/api_tests/` :
- `*_report.html` - Rapport visuel détaillé
- `*_results.json` - Données brutes pour analyse

### Points de Validation Critiques

1. **Formats de données** :
   - Binance : prix en string
   - Bitget : prix en string dans `lastPr`
   - MTN : montants en string
   - Orange : montants en number

2. **Codes de statut** :
   - MTN : 202 pour initiation
   - Orange : 201 pour initiation
   - Tous : 200 pour consultation

3. **Champs de message** :
   - MTN : utilise `reason`
   - Orange : utilise `message`

4. **Headers d'authentification** :
   - Binance : `X-MBX-APIKEY`
   - Bitget : `ACCESS-KEY`, `ACCESS-SIGN`, `ACCESS-TIMESTAMP`, `ACCESS-PASSPHRASE`
   - MTN : `Authorization`, `X-Target-Environment`, `Ocp-Apim-Subscription-Key`
   - Orange : `Authorization`

## 🐛 Résolution des Problèmes

### Erreurs Communes

1. **401 Unauthorized** :
   - Vérifier les clés API dans l'environnement
   - S'assurer d'utiliser les bonnes URLs (sandbox vs production)

2. **403 Forbidden** :
   - Vérifier les permissions des clés API
   - Contrôler les restrictions IP si configurées

3. **Rate Limiting** :
   - Ajouter des délais entre requêtes (`--delay-request 1000`)
   - Utiliser des clés avec limites plus élevées

4. **Timeout** :
   - Augmenter `--timeout-request`
   - Vérifier la connectivité réseau

### Validation des Implémentations

Si un test échoue, comparer :

1. **URL utilisée** dans le test vs notre client Go
2. **Headers** envoyés vs ceux dans notre implémentation
3. **Structure de la réponse** attendue vs celle parsée par notre client
4. **Gestion des erreurs** (codes de statut, messages)

## 📝 Mise à Jour des Clients

Si des divergences sont trouvées :

1. **Corriger le client Go** pour correspondre à l'API réelle
2. **Mettre à jour les tests unitaires** avec les bons formats
3. **Valider avec les tests d'intégration** du service Client
4. **Re-exécuter cette collection** pour confirmer la correction

## 🔄 Intégration Continue

Intégrer ces tests dans le pipeline CI/CD :

```yaml
# .github/workflows/api-validation.yml
- name: Validate Third Party APIs
  run: |
    npm install -g newman
    ./scripts/test_third_party_apis.sh
```

## 📚 Références

- [Documentation Binance API](https://binance-docs.github.io/apidocs/spot/en/)
- [Documentation Bitget API](https://bitgetlimited.github.io/apidoc/en/spot/)
- [Documentation MTN Mobile Money](https://momodeveloper.mtn.com/docs/services/collection/)
- [Documentation Orange Money](https://developer.orange.com/apis/mobile-money/)
- [Newman Documentation](https://learning.postman.com/docs/running-collections/using-newman-cli/)

---

**Note** : Cette validation est essentielle avant l'intégration avec Temporal pour s'assurer que nos clients fonctionnent correctement avec les vraies APIs.
