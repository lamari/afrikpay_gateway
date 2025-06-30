
## Contexte du projet
- Projet Temporal existant avec workflows et activités déjà créé
- **IMPORTANT** : Tous les fichiers existent déjà - NE PAS créer de nouveaux fichiers sans validation explicite
- Objectif : Implémenter le workflow MTN pour créer des payment requests via l'API MTN

## Instructions de développement

### Étape 1 : Analyser la structure existante
- Examiner la structure du projet pour identifier les fichiers appropriés
- Localiser le client MTN existant pour y ajouter les nouvelles méthodes
- Identifier le fichier de workflow MTN à modifier

### Étape 2 : Implémenter les méthodes client MTN
Ajouter dans le client MTN existant les 4 méthodes suivantes :

**Base URL** : `https://sandbox.momodeveloper.mtn.com`

#### 1. CreateUser()
- **Endpoint** : `POST /v1_0/apiuser`
- **Headers** :
  - `Ocp-Apim-Subscription-Key` : primary key depuis config
  - `X-Reference-Id` : UUID généré (stocker dans `mtn_reference_id`)
- **Body** : `{"providerCallbackHost": "string"}`

#### 2. CreateApiKey()
- **Endpoint** : `POST /v1_0/apiuser/{mtn_reference_id}/apikey`
- **Headers** :
  - `Ocp-Apim-Subscription-Key` : primary key depuis config  
  - `X-Reference-Id` : `mtn_reference_id` généré précédemment
- **Body** : `{"providerCallbackHost": "string"}`
- **Response** : Stocker l'API key dans `mtn_api_key`

#### 3. GetAccessToken()
- **Endpoint** : `POST /v1_0/token`
- **Headers** :
  - `Ocp-Apim-Subscription-Key` : primary key depuis config
  - `Authorization` : `Basic ` + base64(`mtn_reference_id:mtn_api_key`)
- **Body** : vide
- **Response** : Stocker l'access_token dans `mtn_access_token`

#### 4. CreatePaymentRequest()
- **Endpoint** : `POST /v1_0/paymentrequest`
- **Headers** :
  - `Ocp-Apim-Subscription-Key` : primary key depuis config
  - `Authorization` : `Bearer {mtn_access_token}`
  - `X-Reference-Id` : `mtn_reference_id`
- **Body exemple** :
```json
{
  "amount": "100",
  "currency": "EUR", 
  "externalId": "test-{{timestamp}}",
  "payer": {
    "partyIdType": "MSISDN",
    "partyId": "256774290781"
  },
  "payerMessage": "Test payment from Afrikpay Gateway",
  "payeeNote": "Payment for crypto purchase"
}
```
- **Response attendue** : Status 202

### Étape 3 : Tests E2E
- Modifier le fichier existant : `services/temporal/internal/e2e/mtn_activities_e2e_test.go`
- Créer des tests pour chaque méthode client
- Tester le workflow complet end-to-end

### Étape 4 : Implémentation du workflow
- Créer les activités Temporal pour chaque méthode client
- Implémenter le workflow MTN qui orchestre ces activités dans l'ordre
- Gérer les erreurs et retry policies appropriés

**Reminder** : Modifier uniquement les fichiers existants, analyser d'abord la structure pour identifier les bons emplacements.