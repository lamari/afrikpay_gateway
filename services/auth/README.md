# Afrikpay Gateway – Auth Service

Service d’authentification JWT pour Afrikpay Gateway.

## Sommaire
1. Présentation
2. Démarrage rapide
3. Variables d’environnement
4. Endpoints
5. Exemples d’utilisation (curl)
6. Collection Postman
7. Tests

---

## 1. Présentation
Ce micro-service fournit :
* Authentification utilisateur (`/auth/login`)
* Vérification de token (`/auth/verify`)
* Renouvellement de token (`/auth/refresh`)
* Endpoint protégé exemple (`/protected/profile`)
* Endpoints de santé (`/health`, `/ready`, `/live`)

Les tokens sont signés avec des clés RSA 2048 bits.

---

## 2. Démarrage rapide
```bash
# Lancer uniquement le service Auth
make run-auth          # via Makefile
# ou via Docker Compose (tous les services)
docker compose up auth # nécessite Docker
```
Le service écoute par défaut sur `http://localhost:8001`.

---

## 3. Variables d’environnement essentielles
| Variable | Par défaut | Description |
|----------|-----------|-------------|
| `AUTH_HOST` | 0.0.0.0 | Adresse d’écoute |
| `AUTH_PORT` | 8001 | Port HTTP |
| `JWT_PRIVATE_KEY` | ./config/keys/private.pem | Chemin clé privée RSA |
| `JWT_PUBLIC_KEY` | ./config/keys/public.pem | Chemin clé publique RSA |
| `ACCESS_TOKEN_TTL` | 900 | Durée de vie access token (s) |
| `REFRESH_TOKEN_TTL` | 604800 | Durée de vie refresh token (s) |
| `ALLOWED_ORIGINS` | * | Origines CORS autorisées |

---

## 4. Endpoints
| Méthode | Chemin | Description |
|---------|--------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |
| GET | `/live` | Liveness check |
| POST | `/auth/login` | Authentification utilisateur |
| GET | `/auth/verify` | Vérification de token (header `Authorization: Bearer <token>`) |
| POST | `/auth/refresh` | Renouvellement de token |
| GET | `/protected/profile` | Endpoint protégé, nécessite un token valide |

La spécification complète OpenAPI est disponible dans `api/openapi.yml`.

---

## 5. Exemples d’utilisation (curl)
```bash
# Login
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Vérification (remplacez $ACCESS_TOKEN)
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
     http://localhost:8001/auth/verify

# Renouvellement (remplacez $REFRESH_TOKEN)
curl -X POST http://localhost:8001/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"'$REFRESH_TOKEN'"}'
```

---

## 6. Collection Postman
Une collection Postman prête à l’emploi est fournie : `api/postman_collection.json`.
1. Importez la collection dans Postman.
2. Définissez la variable d’environnement `base_url` (`http://localhost:8001`).
3. Exécutez la requête **Login** pour initialiser `access_token` et `refresh_token`.
4. Les autres requêtes utiliseront ces variables automatiquement.

---

## 7. Tests
```bash
make test        # Tests unitaires
make integration  # Tests d’intégration
```
Les rapports de couverture sont générés dans `coverage.html`.
