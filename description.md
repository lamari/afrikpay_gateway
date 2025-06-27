# Afrikpay Gateway - Roadmap de Développement

## 🎯 Vue d'ensemble du projet

**Objectif** : Développer une gateway crypto qui permet aux utilisateurs d'acheter des cryptomonnaies et de recharger leur wallet Afrikpay via Mobile Money.

### Fonctionnalités principales
- **Achat de crypto** : USDT/BTC via APIs d'exchange (Binance, Bitget)
- **Recharge wallet** : Via Mobile Money (MTN, Orange)
- **Gestion sécurisée** : Authentification JWT, transactions fiables
- **Architecture microservices** : Services découplés et scalables

### Technologies confirmées
- **Backend** : Go + Temporal + Docker
- **Bases de données** : MongoDB (CRUD) + PostgreSQL (Temporal)
- **Authentification** : JWT avec microservice dédié
- **Approche** : TDD (Test-Driven Development)
- **Déploiement** : Docker Compose + Makefile

---

## 🏗️ Architecture Microservices

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   CRUD Service  │    │ Temporal Service│
│    (JWT)        │    │   (MongoDB)     │    │  (PostgreSQL)   │
│   Port: 8001    │    │   Port: 8002    │    │   Port: 8003    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  Client Module  │
                    │ (3rd Party APIs)│
                    └─────────────────┘
```

## Structure du Projet (TDD)
afrikpay/
├── services/
│   ├── auth/                 # Microservice Auth JWT
│   ├── crud/                 # Microservice CRUD MongoDB  
│   ├── temporal/             # Microservice Temporal Workflows
│   └── client/               # Module client APIs tierces
├── shared/
│   ├── config/               # Configuration centralisée
│   ├── models/               # Modèles partagés
│   └── utils/                # Utilitaires communs
├── tests/                    # Tests d'intégration
├── docker-compose.yml        # Environnement local complet
├── Makefile                  # Commandes simplifiées
└── docs/                     # Documentation API

### Services
1. **Auth Service** : Gestion JWT (remplaçable par Authentik/Keycloak)
2. **CRUD Service** : API REST pour User/Wallet/Transaction (MongoDB)
3. **Temporal Service** : Orchestration des workflows (PostgreSQL)
4. **Client Module** : Connexions APIs tierces (Binance, Mobile Money)

---

## 🎯 Critères de Validation par Phase

### Phase 1 (Auth Service)
- ✅ Tests passent avec coverage > 90%
- ✅ Service démarrable via `make start-auth`
- ✅ JWT valides générés et validés
- ✅ Documentation API disponible

### Phase 2 (CRUD Service) 
- ✅ CRUD opérationnel avec MongoDB
- ✅ Authentification JWT fonctionnelle
- ✅ Validation des données robuste
- ✅ Tests d'intégration passent

### Phase 3 (Client Module)
- ✅ Connexions APIs tierces simulées
- ✅ Resilience patterns implémentés
- ✅ Gestion d'erreurs complète
- ✅ Logs et monitoring en place

### Phase 4 (Temporal Service)
- ✅ Workflows fonctionnels avec compensation
- ✅ Activities robustes avec retry
- ✅ API REST opérationnelle
- ✅ Tests end-to-end passent

### Phase 5 (Intégration)
- ✅ Communication inter-services fluide
- ✅ Scénarios métier complets
- ✅ Performance acceptable
- ✅ Sécurité validée

### Phase 6 (Finalisation)
- ✅ Documentation complète
- ✅ Déploiement automatisé
- ✅ Monitoring opérationnel
- ✅ Prêt pour production

---

## 🚀 Commandes Rapides

```bash
# Démarrer l'environnement complet
make start

# Lancer tous les tests
make test

# Vérifier le coverage
make coverage

# Build toutes les images Docker
make build

# Nettoyer l'environnement
make clean
```

## 📊 Suivi de Progression

**Phases complétées** : 0/6
**Étapes complétées** : 0/XX
**Coverage global** : 0%
**Dernière mise à jour** : [Date]

---

*Ce roadmap sera mis à jour au fur et à mesure de l'avancement du projet.*