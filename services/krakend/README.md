# KrakenD API Gateway - Afrikpay Gateway

Ce service sert de passerelle API (API Gateway) pour tous les microservices d'Afrikpay Gateway.

## Description

KrakenD est une passerelle API haute performance qui agit comme un point d'entrée unique pour toutes les requêtes clients. Cette passerelle :

- Unifie l'accès à tous les microservices (Auth, CRUD, Temporal)
- Gère l'authentification JWT via le service Auth
- Effectue le routage vers les bons services en fonction de l'URL
- Offre des fonctionnalités de rate limiting, caching et monitoring

## Architecture

```
┌───────────┐     ┌──────────────────┐     ┌──────────────┐
│           │     │                  │     │              │
│  Client   │────▶│  KrakenD (8000)  │────▶│  Auth (8001) │
│           │     │                  │     │              │
└───────────┘     └──────────────────┘     └──────────────┘
                           │
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
    ┌──────────────┐           ┌──────────────────┐
    │              │           │                  │
    │ CRUD (8002)  │           │ Temporal (8003)  │
    │              │           │                  │
    └──────────────┘           └──────────────────┘
```

## Configuration

La configuration de KrakenD est définie dans le fichier `config/krakend.json`. Ce fichier contient :

- Les routes vers tous les services
- La configuration d'authentification JWT
- Les paramètres CORS
- La configuration de logging et monitoring

## Démarrage

Pour démarrer le service KrakenD:

```bash
cd services/krakend
docker-compose up -d
```

Le service sera accessible à l'adresse `http://localhost:8000`

## Endpoints

Le service expose les endpoints suivants:

### Auth Service
- POST `/api/auth/login` - Authentification utilisateur
- GET `/api/auth/verify` - Vérification de token
- POST `/api/auth/refresh` - Renouvellement de token

### CRUD Service
- POST `/api/users` - Création d'utilisateur
- GET/PUT/DELETE `/api/users/{id}` - Opérations sur un utilisateur
- POST `/api/wallets` - Création de portefeuille
- GET/PUT/DELETE `/api/wallets/{id}` - Opérations sur un portefeuille
- POST `/api/transactions` - Création de transaction
- GET/PUT/DELETE `/api/transactions/{id}` - Opérations sur une transaction

### Temporal Service
- POST `/api/v1/mtn/payments` - Initier un paiement MTN
- GET `/api/v1/mtn/payments/{id}` - Vérifier l'état d'un paiement MTN
- POST `/api/v1/orange/payments` - Initier un paiement Orange Money
- GET `/api/v1/orange/payments/{id}` - Vérifier l'état d'un paiement Orange Money
- POST `/api/v1/crypto/buy` - Initier un achat de cryptomonnaie
- GET `/api/v1/crypto/exchanges` - Lister les exchanges supportés

## Sécurité

Toutes les routes protégées nécessitent un token JWT valide dans le header `Authorization` avec le format `Bearer {token}`. Ce token est vérifié auprès du service Auth.
