# Collections Postman pour Afrikpay Gateway

Ce répertoire contient des collections Postman permettant de tester les différentes API des microservices Afrikpay Gateway.

## Structure

```
postman/
├── collections/         # Collections Postman pour chaque service
│   ├── auth_service.json
│   └── ... (autres services)
├── environments/        # Environnements Postman (dev, staging, prod)
└── README.md           # Ce fichier
```

## Collections disponibles

### Auth Service

La collection `auth_service.json` permet de tester le service d'authentification. Elle inclut :

- **Health Check** : Vérification de l'état du service
  - `GET /health` - État général du service
  - `GET /ready` - État de disponibilité du service

- **Authentication** : Endpoints d'authentification
  - `POST /auth/login` - Connexion utilisateur (obtention des tokens JWT)
  - `POST /auth/refresh` - Rafraîchissement d'un token JWT
  - `GET /auth/verify` - Vérification de la validité d'un token JWT

- **Protected Routes** : Endpoints protégés nécessitant une authentification
  - `GET /protected/profile` - Exemple d'endpoint protégé retournant le profil utilisateur

## Comment utiliser

### Avec Postman

1. Importer la collection dans Postman
2. Configurer la variable `baseUrl` selon votre environnement (par défaut: `http://localhost:8001` pour le service auth)
3. Exécuter les requêtes dans l'ordre logique (ex: login → verify → protected routes)

### Avec Newman (ligne de commande)

Pour exécuter les tests via Newman, utilisez la commande suivante :

```bash
# Installation de Newman (si nécessaire)
npm install -g newman

# Exécution d'une collection
newman run postman/collections/auth_service.json
```

## Authentification

La collection gère automatiquement les tokens JWT :

1. La requête `Login` stocke automatiquement les tokens obtenus dans les variables de collection
2. Les autres requêtes utilisent ces variables pour l'authentification
3. La requête `Refresh Token` met à jour ces variables avec les nouveaux tokens
