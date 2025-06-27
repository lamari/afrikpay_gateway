# Documentation Afrikpay Gateway

Ce dossier contient la documentation technique et les outils de validation pour le projet Afrikpay Gateway.

## 📁 Contenu

### Collections de Tests API

- **`third_party_apis.postman_collection.json`** - Collection Postman pour tester les APIs tierces
- **`third_party_apis.postman_environment.json`** - Variables d'environnement pour les tests
- **`API_VALIDATION_GUIDE.md`** - Guide complet d'utilisation des tests API

### Documentation Technique

- **Architecture des microservices** (à venir)
- **Workflows Temporal** (à venir)
- **Intégration des APIs tierces** (à venir)

## 🎯 Validation des APIs Tierces

### Démarrage Rapide

1. **Installer Newman** :
   ```bash
   npm install -g newman
   ```

2. **Configurer les clés API** dans `third_party_apis.postman_environment.json`

3. **Exécuter les tests** :
   ```bash
   cd /Users/lamari_alaa/Documents/Projects/afrikpay
   ./scripts/test_third_party_apis.sh
   ```

### APIs Testées

| Service | Endpoints | Authentification |
|---------|-----------|------------------|
| **Binance** | Prix, Tickers | API Key |
| **Bitget** | Prix, Tickers | API Key + Signature |
| **MTN Mobile Money** | Balance, Paiements | Bearer Token + Subscription Key |
| **Orange Money** | Balance, Paiements | Bearer Token |

### Objectifs de Validation

- ✅ **Endpoints corrects** - URLs et méthodes HTTP
- ✅ **Headers d'authentification** - Formats et noms exacts
- ✅ **Structures de requête** - JSON payloads conformes
- ✅ **Structures de réponse** - Parsing correct des données
- ✅ **Codes de statut** - Gestion appropriée des succès/erreurs
- ✅ **Formats de données** - Types (string vs number) cohérents

## 🔧 Utilisation

### Tests Complets

```bash
# Tous les services
./scripts/test_third_party_apis.sh

# Service spécifique
./scripts/test_third_party_apis.sh binance
./scripts/test_third_party_apis.sh bitget
./scripts/test_third_party_apis.sh mtn
./scripts/test_third_party_apis.sh orange
```

### Tests via Postman GUI

1. Importer la collection `third_party_apis.postman_collection.json`
2. Importer l'environnement `third_party_apis.postman_environment.json`
3. Configurer les variables avec vos clés API
4. Exécuter les tests individuellement ou en batch

### Tests via Newman CLI

```bash
newman run third_party_apis.postman_collection.json \
  --environment third_party_apis.postman_environment.json \
  --reporters cli,html \
  --reporter-html-export reports/api_validation.html
```

## 📊 Rapports

Les tests génèrent des rapports dans `../reports/api_tests/` :

- **HTML Reports** - Visualisation détaillée des résultats
- **JSON Reports** - Données brutes pour analyse automatisée
- **CLI Output** - Résultats en temps réel dans le terminal

## 🚨 Points d'Attention

### Différences Critiques Entre APIs

1. **Formats de Prix** :
   - Binance/Bitget : string (`"43250.50"`)
   - Nos clients : conversion vers float64

2. **Codes de Statut Paiement** :
   - MTN : 202 pour initiation
   - Orange : 201 pour initiation

3. **Champs de Message** :
   - MTN : `reason`
   - Orange : `message`

4. **Authentification** :
   - Binance : Simple API key
   - Bitget : Signature complexe (HMAC-SHA256)
   - MTN : Multiple headers (Bearer + Subscription)
   - Orange : Bearer token simple

## 🔄 Workflow de Validation

1. **Exécuter les tests** contre les APIs réelles
2. **Analyser les divergences** avec nos implémentations
3. **Corriger les clients Go** si nécessaire
4. **Mettre à jour les tests unitaires** pour refléter les corrections
5. **Valider l'intégration** avec le service Client
6. **Procéder aux tests Temporal** une fois validé

## 📚 Ressources

- [Guide de Validation API](./API_VALIDATION_GUIDE.md) - Documentation complète
- [Collection Postman](./third_party_apis.postman_collection.json) - Tests automatisés
- [Environnement Postman](./third_party_apis.postman_environment.json) - Configuration
- [Script de Test](../scripts/test_third_party_apis.sh) - Automatisation Newman

---

**Important** : Ces tests doivent être exécutés avant toute intégration avec Temporal pour garantir que nos clients correspondent exactement aux comportements des APIs tierces réelles.
