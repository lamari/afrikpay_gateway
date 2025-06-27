# Structure Détaillée du Projet Afrikpay Gateway

## 📁 Arborescence Complète

```
afrikpay_gateway/
├── 📁 .github/                              # GitHub workflows et templates
│   ├── workflows/
│   │   ├── ci.yml                          # CI/CD pipeline
│   │   ├── test.yml                        # Tests automatisés
│   │   └── security.yml                    # Scan sécurité
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   └── pull_request_template.md
│
├── 📁 services/                             # Microservices principaux
│   ├── 📁 auth/                            # Service d'authentification JWT
│   │   ├── 📁 cmd/
│   │   │   └── main.go                     # Point d'entrée du service
│   │   ├── 📁 internal/
│   │   │   ├── 📁 config/
│   │   │   │   ├── config.go               # Configuration du service
│   │   │   │   └── config_test.go
│   │   │   ├── 📁 handlers/
│   │   │   │   ├── auth_handler.go         # Handlers HTTP
│   │   │   │   ├── auth_handler_test.go
│   │   │   │   ├── health_handler.go
│   │   │   │   └── health_handler_test.go
│   │   │   ├── 📁 middleware/
│   │   │   │   ├── auth_middleware.go      # Middleware JWT
│   │   │   │   ├── auth_middleware_test.go
│   │   │   │   ├── cors_middleware.go
│   │   │   │   └── logging_middleware.go
│   │   │   ├── 📁 models/
│   │   │   │   ├── auth.go                 # Modèles d'auth
│   │   │   │   ├── auth_test.go
│   │   │   │   ├── token.go
│   │   │   │   └── token_test.go
│   │   │   ├── 📁 services/
│   │   │   │   ├── jwt_service.go          # Logique métier JWT
│   │   │   │   ├── jwt_service_test.go
│   │   │   │   ├── auth_service.go
│   │   │   │   └── auth_service_test.go
│   │   │   └── 📁 utils/
│   │   │       ├── crypto.go               # Utilitaires crypto
│   │   │       ├── crypto_test.go
│   │   │       ├── validator.go
│   │   │       └── validator_test.go
│   │   ├── 📁 api/
│   │   │   └── openapi.yml                 # Spécification OpenAPI
│   │   ├── 📁 deployments/
│   │   │   ├── Dockerfile
│   │   │   └── docker-compose.yml
│   │   ├── 📁 scripts/
│   │   │   ├── generate-keys.sh            # Génération clés JWT
│   │   │   └── migrate.sh
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Makefile
│   │   └── README.md
│   │
│   ├── 📁 crud/                            # Service CRUD avec MongoDB
│   │   ├── 📁 cmd/
│   │   │   └── main.go
│   │   ├── 📁 internal/
│   │   │   ├── 📁 config/
│   │   │   │   ├── config.go
│   │   │   │   └── config_test.go
│   │   │   ├── 📁 database/
│   │   │   │   ├── mongodb.go              # Connexion MongoDB
│   │   │   │   ├── mongodb_test.go
│   │   │   │   ├── migrations.go
│   │   │   │   └── migrations_test.go
│   │   │   ├── 📁 handlers/
│   │   │   │   ├── user_handler.go         # CRUD Users
│   │   │   │   ├── user_handler_test.go
│   │   │   │   ├── wallet_handler.go       # CRUD Wallets
│   │   │   │   ├── wallet_handler_test.go
│   │   │   │   ├── transaction_handler.go  # CRUD Transactions
│   │   │   │   ├── transaction_handler_test.go
│   │   │   │   ├── health_handler.go
│   │   │   │   └── health_handler_test.go
│   │   │   ├── 📁 middleware/
│   │   │   │   ├── auth_middleware.go      # Validation JWT
│   │   │   │   ├── auth_middleware_test.go
│   │   │   │   ├── validation_middleware.go
│   │   │   │   ├── validation_middleware_test.go
│   │   │   │   └── logging_middleware.go
│   │   │   ├── 📁 models/
│   │   │   │   ├── user.go                 # Modèle User
│   │   │   │   ├── user_test.go
│   │   │   │   ├── wallet.go               # Modèle Wallet
│   │   │   │   ├── wallet_test.go
│   │   │   │   ├── transaction.go          # Modèle Transaction
│   │   │   │   ├── transaction_test.go
│   │   │   │   ├── base.go                 # Modèle de base
│   │   │   │   └── base_test.go
│   │   │   ├── 📁 repositories/
│   │   │   │   ├── user_repository.go      # Repository User
│   │   │   │   ├── user_repository_test.go
│   │   │   │   ├── wallet_repository.go    # Repository Wallet
│   │   │   │   ├── wallet_repository_test.go
│   │   │   │   ├── transaction_repository.go
│   │   │   │   ├── transaction_repository_test.go
│   │   │   │   ├── interfaces.go           # Interfaces repositories
│   │   │   │   └── base_repository.go
│   │   │   ├── 📁 services/
│   │   │   │   ├── user_service.go         # Logique métier User
│   │   │   │   ├── user_service_test.go
│   │   │   │   ├── wallet_service.go       # Logique métier Wallet
│   │   │   │   ├── wallet_service_test.go
│   │   │   │   ├── transaction_service.go
│   │   │   │   ├── transaction_service_test.go
│   │   │   │   └── interfaces.go
│   │   │   └── 📁 utils/
│   │   │       ├── pagination.go
│   │   │       ├── pagination_test.go
│   │   │       ├── validator.go
│   │   │       ├── validator_test.go
│   │   │       ├── response.go
│   │   │       └── response_test.go
│   │   ├── 📁 api/
│   │   │   └── openapi.yml
│   │   ├── 📁 deployments/
│   │   │   ├── Dockerfile
│   │   │   └── docker-compose.yml
│   │   ├── 📁 scripts/
│   │   │   ├── migrate.sh
│   │   │   └── seed.sh                     # Données de test
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Makefile
│   │   └── README.md
│   │
│   ├── 📁 temporal/                        # Service Temporal Workflows
│   │   ├── 📁 cmd/
│   │   │   ├── worker/
│   │   │   │   └── main.go                 # Worker Temporal
│   │   │   └── server/
│   │   │       └── main.go                 # API Server
│   │   ├── 📁 internal/
│   │   │   ├── 📁 config/
│   │   │   │   ├── config.go
│   │   │   │   └── config_test.go
│   │   │   ├── 📁 activities/
│   │   │   │   ├── binance_activity.go     # Activity Binance
│   │   │   │   ├── binance_activity_test.go
│   │   │   │   ├── crud_activity.go        # Activity CRUD calls
│   │   │   │   ├── crud_activity_test.go
│   │   │   │   ├── mobilemoney_activity.go
│   │   │   │   ├── mobilemoney_activity_test.go
│   │   │   │   ├── notification_activity.go
│   │   │   │   ├── notification_activity_test.go
│   │   │   │   └── interfaces.go
│   │   │   ├── 📁 workflows/
│   │   │   │   ├── crypto_buy_workflow.go  # Workflow achat crypto
│   │   │   │   ├── crypto_buy_workflow_test.go
│   │   │   │   ├── wallet_deposit_workflow.go
│   │   │   │   ├── wallet_deposit_workflow_test.go
│   │   │   │   ├── compensation_workflow.go # Pattern Saga
│   │   │   │   ├── compensation_workflow_test.go
│   │   │   │   └── interfaces.go
│   │   │   ├── 📁 handlers/
│   │   │   │   ├── crypto_handler.go       # /crypto/quote, /crypto/buy
│   │   │   │   ├── crypto_handler_test.go
│   │   │   │   ├── wallet_handler.go       # /wallet/deposit
│   │   │   │   ├── wallet_handler_test.go
│   │   │   │   ├── workflow_handler.go     # Status workflows
│   │   │   │   ├── workflow_handler_test.go
│   │   │   │   ├── health_handler.go
│   │   │   │   └── health_handler_test.go
│   │   │   ├── 📁 models/
│   │   │   │   ├── crypto.go               # Modèles crypto
│   │   │   │   ├── crypto_test.go
│   │   │   │   ├── workflow.go             # Modèles workflows
│   │   │   │   ├── workflow_test.go
│   │   │   │   ├── activity.go
│   │   │   │   └── activity_test.go
│   │   │   ├── 📁 clients/
│   │   │   │   ├── crud_client.go          # Client vers CRUD service
│   │   │   │   ├── crud_client_test.go
│   │   │   │   ├── auth_client.go          # Client vers Auth service
│   │   │   │   ├── auth_client_test.go
│   │   │   │   └── interfaces.go
│   │   │   └── 📁 utils/
│   │   │       ├── errors.go               # Gestion erreurs Temporal
│   │   │       ├── errors_test.go
│   │   │       ├── retry.go                # Politiques retry
│   │   │       ├── retry_test.go
│   │   │       ├── timeout.go
│   │   │       └── timeout_test.go
│   │   ├── 📁 api/
│   │   │   └── openapi.yml
│   │   ├── 📁 deployments/
│   │   │   ├── Dockerfile.worker
│   │   │   ├── Dockerfile.server
│   │   │   └── docker-compose.yml
│   │   ├── 📁 scripts/
│   │   │   ├── setup-temporal.sh           # Setup namespace Temporal
│   │   │   └── migrate.sh
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Makefile
│   │   └── README.md
│   │
│   └── 📁 client/                          # Module clients APIs tierces
│       ├── 📁 cmd/
│       │   └── main.go                     # CLI de test
│       ├── 📁 internal/
│       │   ├── 📁 config/
│       │   │   ├── config.go
│       │   │   └── config_test.go
│       │   ├── 📁 binance/
│       │   │   ├── client.go               # Client Binance
│       │   │   ├── client_test.go
│       │   │   ├── models.go               # Modèles Binance
│       │   │   ├── models_test.go
│       │   │   ├── auth.go                 # Auth Binance
│       │   │   └── auth_test.go
│       │   ├── 📁 coinbase/
│       │   │   ├── client.go               # Client Coinbase
│       │   │   ├── client_test.go
│       │   │   ├── models.go
│       │   │   ├── models_test.go
│       │   │   ├── auth.go
│       │   │   └── auth_test.go
│       │   ├── 📁 mobilemoney/
│       │   │   ├── 📁 mtn/
│       │   │   │   ├── client.go           # Client MTN MoMo
│       │   │   │   ├── client_test.go
│       │   │   │   ├── models.go
│       │   │   │   ├── models_test.go
│       │   │   │   ├── auth.go
│       │   │   │   └── auth_test.go
│       │   │   ├── 📁 orange/
│       │   │   │   ├── client.go           # Client Orange Money
│       │   │   │   ├── client_test.go
│       │   │   │   ├── models.go
│       │   │   │   ├── models_test.go
│       │   │   │   ├── auth.go
│       │   │   │   └── auth_test.go
│       │   │   ├── interfaces.go           # Interfaces communes
│       │   │   └── factory.go              # Factory pattern
│       │   ├── 📁 common/
│       │   │   ├── 📁 resilience/
│       │   │   │   ├── circuit_breaker.go  # Circuit breaker
│       │   │   │   ├── circuit_breaker_test.go
│       │   │   │   ├── retry.go            # Retry avec backoff
│       │   │   │   ├── retry_test.go
│       │   │   │   ├── timeout.go          # Gestion timeouts
│       │   │   │   └── timeout_test.go
│       │   │   ├── 📁 http/
│       │   │   │   ├── client.go           # Client HTTP commun
│       │   │   │   ├── client_test.go
│       │   │   │   ├── middleware.go       # Middleware HTTP
│       │   │   │   └── middleware_test.go
│       │   │   └── 📁 monitoring/
│       │   │       ├── metrics.go          # Métriques
│       │   │       ├── metrics_test.go
│       │   │       ├── logger.go           # Logging structuré
│       │   │       └── logger_test.go
│       │   └── 📁 mocks/
│       │       ├── binance_mock.go         # Mocks pour tests
│       │       ├── coinbase_mock.go
│       │       ├── mtn_mock.go
│       │       └── orange_mock.go
│       ├── 📁 examples/
│       │   ├── binance_example.go          # Exemples d'utilisation
│       │   ├── coinbase_example.go
│       │   ├── mtn_example.go
│       │   └── orange_example.go
│       ├── go.mod
│       ├── go.sum
│       ├── Makefile
│       └── README.md
│
├── 📁 shared/                              # Code partagé entre services
│   ├── 📁 config/
│   │   ├── config.go                       # Configuration globale
│   │   ├── config_test.go
│   │   ├── loader.go                       # Chargement config
│   │   └── loader_test.go
│   ├── 📁 models/
│   │   ├── common.go                       # Modèles communs
│   │   ├── common_test.go
│   │   ├── errors.go                       # Erreurs standardisées
│   │   ├── errors_test.go
│   │   ├── responses.go                    # Réponses API
│   │   └── responses_test.go
│   ├── 📁 utils/
│   │   ├── crypto.go                       # Utilitaires crypto
│   │   ├── crypto_test.go
│   │   ├── strings.go                      # Utilitaires strings
│   │   ├── strings_test.go
│   │   ├── time.go                         # Utilitaires temps
│   │   ├── time_test.go
│   │   ├── validator.go                    # Validation commune
│   │   └── validator_test.go
│   ├── 📁 middleware/
│   │   ├── cors.go                         # Middleware CORS
│   │   ├── cors_test.go
│   │   ├── logging.go                      # Middleware logging
│   │   ├── logging_test.go
│   │   ├── ratelimit.go                    # Rate limiting
│   │   └── ratelimit_test.go
│   ├── go.mod
│   ├── go.sum
│   └── README.md
│
├── 📁 config/                              # Configuration centralisée
│   ├── config.yml                          # Config principale
│   ├── config.dev.yml                      # Config développement
│   ├── config.prod.yml                     # Config production
│   ├── config.test.yml                     # Config tests
│   ├── 📁 keys/
│   │   ├── jwt_private.key                 # Clé privée JWT
│   │   ├── jwt_public.key                  # Clé publique JWT
│   │   ├── .gitkeep
│   │   └── README.md                       # Instructions génération clés
│   └── 📁 secrets/
│       ├── .env.example                    # Exemple variables d'env
│       ├── .gitkeep
│       └── README.md
│
├── 📁 deployments/                         # Configuration déploiement
│   ├── docker-compose.yml                  # Compose principal
│   ├── docker-compose.dev.yml              # Compose développement
│   ├── docker-compose.prod.yml             # Compose production
│   ├── docker-compose.test.yml             # Compose tests
│   ├── 📁 kubernetes/
│   │   ├── namespace.yml
│   │   ├── configmap.yml
│   │   ├── secrets.yml
│   │   ├── auth-service.yml
│   │   ├── crud-service.yml
│   │   ├── temporal-service.yml
│   │   └── ingress.yml
│   └── 📁 helm/
│       ├── Chart.yml
│       ├── values.yml
│       └── 📁 templates/
│           ├── deployment.yml
│           ├── service.yml
│           └── configmap.yml
│
├── 📁 scripts/                             # Scripts utilitaires
│   ├── setup.sh                           # Setup environnement complet
│   ├── build.sh                           # Build tous les services
│   ├── test.sh                            # Lancer tous les tests
│   ├── clean.sh                           # Nettoyage
│   ├── migrate.sh                         # Migrations DB
│   ├── seed.sh                            # Données de test
│   ├── generate-keys.sh                   # Génération clés JWT
│   ├── docker-build.sh                    # Build images Docker
│   └── coverage.sh                        # Rapport coverage global
│
├── 📁 tests/                              # Tests d'intégration globaux
│   ├── 📁 integration/
│   │   ├── auth_test.go                   # Tests intégration auth
│   │   ├── crud_test.go                   # Tests intégration CRUD
│   │   ├── temporal_test.go               # Tests intégration Temporal
│   │   ├── e2e_test.go                    # Tests end-to-end
│   │   └── setup_test.go                  # Setup tests
│   ├── 📁 load/
│   │   ├── crypto_buy_test.go             # Tests de charge
│   │   ├── wallet_deposit_test.go
│   │   └── concurrent_test.go
│   ├── 📁 fixtures/
│   │   ├── users.json                     # Données de test
│   │   ├── wallets.json
│   │   ├── transactions.json
│   │   └── crypto_prices.json
│   ├── 📁 mocks/
│   │   ├── auth_mock.go                   # Mocks services
│   │   ├── crud_mock.go
│   │   └── temporal_mock.go
│   ├── go.mod
│   ├── go.sum
│   └── README.md
│
├── 📁 docs/                               # Documentation
│   ├── README.md                          # README principal
│   ├── ROADMAP.md                         # Roadmap du projet
│   ├── ARCHITECTURE.md                    # Documentation architecture
│   ├── DEVELOPMENT.md                     # Guide de développement
│   ├── DEPLOYMENT.md                      # Guide de déploiement
│   ├── CONTRIBUTING.md                    # Guide contribution
│   ├── SECURITY.md                        # Documentation sécurité
│   ├── 📁 api/
│   │   ├── auth-api.md                    # Doc API Auth
│   │   ├── crud-api.md                    # Doc API CRUD
│   │   ├── temporal-api.md                # Doc API Temporal
│   │   └── postman/
│   │       ├── auth-collection.json
│   │       ├── crud-collection.json
│   │       └── temporal-collection.json
│   ├── 📁 adr/                           # Architecture Decision Records
│   │   ├── 001-microservices-architecture.md
│   │   ├── 002-temporal-workflow-engine.md
│   │   ├── 003-jwt-authentication.md
│   │   ├── 004-mongodb-vs-postgresql.md
│   │   └── 005-go-vs-other-languages.md
│   ├── 📁 diagrams/
│   │   ├── architecture-overview.png
│   │   ├── sequence-crypto-buy.png
│   │   ├── sequence-wallet-deposit.png
│   │   ├── database-schema.png
│   │   └── temporal-workflows.png
│   └── 📁 tutorials/
│       ├── getting-started.md
│       ├── adding-new-service.md
│       ├── testing-guide.md
│       └── troubleshooting.md
│
├── 📁 tools/                             # Outils de développement
│   ├── 📁 generators/
│   │   ├── service-generator.go           # Générateur de service
│   │   ├── model-generator.go             # Générateur de modèles
│   │   └── test-generator.go              # Générateur de tests
│   ├── 📁 linters/
│   │   ├── .golangci.yml                  # Config linter Go
│   │   └── custom-rules.yml
│   ├── 📁 monitoring/
│   │   ├── prometheus.yml                 # Config Prometheus
│   │   ├── grafana-dashboard.json         # Dashboard Grafana
│   │   └── alerting-rules.yml
│   └── 📁 security/
│       ├── security-scan.sh               # Scan sécurité
│       ├── dependency-check.sh            # Vérification dépendances
│       └── vulnerability-scan.sh
│
├── 📁 vendor/                            # Dépendances vendorisées (optionnel)
│
├── .gitignore                            # Git ignore
├── .gitattributes                        # Git attributes
├── .editorconfig                         # Config éditeur
├── .env.example                          # Exemple variables environnement
├── go.work                               # Go workspace
├── go.work.sum                           # Go workspace sum
├── Makefile                              # Makefile principal
├── docker-compose.yml                    # Compose principal
├── LICENSE                               # Licence
├── README.md                             # README principal
├── CHANGELOG.md                          # Journal des modifications
├── ROADMAP.md                            # Roadmap du projet
├── CONTRIBUTING.md                       # Guide de contribution
└── SECURITY.md                           # Politique de sécurité
```

## 📊 Statistiques de la Structure

### **Services Principaux**
- **4 microservices** : Auth, CRUD, Temporal, Client
- **Indépendants** : Chaque service a son propre go.mod
- **Testables** : Chaque fichier .go a son fichier _test.go
- **Documentés** : API docs + README par service

### **Organisation par Responsabilité**

#### 🔐 **Auth Service**
- **Responsabilité** : Authentification JWT uniquement
- **Technologies** : Go + JWT + clés RSA
- **Remplaçable** : Par Authentik/Keycloak plus tard
- **Tests** : 15+ fichiers de tests

#### 💾 **CRUD Service**
- **Responsabilité** : Gestion des données métier
- **Technologies** : Go + MongoDB + validation
- **Modèles** : User, Wallet, Transaction
- **Tests** : 20+ fichiers de tests

#### ⚡ **Temporal Service**
- **Responsabilité** : Orchestration des workflows
- **Technologies** : Go + Temporal + PostgreSQL
- **Patterns** : Saga, compensation, retry
- **Tests** : 15+ fichiers de tests

#### 🌐 **Client Module**
- **Responsabilité** : Connexions APIs tierces
- **Technologies** : HTTP clients + resilience patterns
- **APIs** : Binance, Coinbase, MTN, Orange
- **Tests** : 25+ fichiers de tests

### **Code Partagé**
- **shared/** : Code commun entre services
- **config/** : Configuration centralisée
- **tests/** : Tests d'intégration globaux
- **docs/** : Documentation complète

### **DevOps & Tooling**
- **Docker** : Compose multi-environnements
- **Kubernetes** : Manifests + Helm charts
- **CI/CD** : GitHub Actions
- **Monitoring** : Prometheus + Grafana
- **Sécurité** : Scans automatiques

## 🚀 **Avantages de cette Structure**

### **Modularité**
- Chaque service est indépendant
- Déploiement séparé possible
- Évolution indépendante

### **Testabilité**
- TDD respecté (test pour chaque fonction)
- Mocks et fixtures organisés
- Coverage trackable

### **Scalabilité**
- Services découplés
- Configuration par environnement
- Monitoring intégré

### **Maintenabilité**
- Documentation complète
- Standards de code
- Architecture Decision Records

## 📝 **Prochaines Étapes**

1. **Créer la structure de base** avec les dossiers principaux
2. **Initialiser les go.mod** pour chaque service
3. **Setup docker-compose** et Makefile
4. **Commencer par l'Auth Service** (Phase 1 TDD)

Cette structure respecte toutes vos exigences et est prête pour un développement professionnel et évolutif !