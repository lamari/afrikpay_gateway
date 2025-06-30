🎯 MISSION WINDSURF : Développement Complet Afrikpay Gateway

Tu es maintenant l'architecte principal du projet Afrikpay Gateway. Ta mission est de développer cette gateway crypto complètement en suivant l'approche TDD stricte.

📋 DOCUMENTS DE RÉFÉRENCE (à lire en premier) :
1. .windsurfrules - Toutes les règles et standards du projet
2. Instructions Windsurf - Workflow de développement détaillé  
3. taches.md - Fichier de suivi de progression (OBLIGATOIRE à mettre à jour)
4. description.md - Fichier de description du projet
5. structure.md - Fichier de structure du projet

🔄 TON WORKFLOW :
1. Lis et comprends les documents de référence (*.md)
2. Consulte taches.md pour voir la progression actuelle
3. Identifie la première tâche non cochée [ ]
4. Développe avec TDD strict : RED → GREEN → REFACTOR
5. Teste et valide selon les critères définis
6. Marque la tâche terminée dans taches.md avec [x] et la date/heure
7. Commit et push les modifications
8. Passe automatiquement à la tâche suivante
9. Répète jusqu'à completion du projet


🎯 RÈGLE FONDAMENTALE : 
À CHAQUE tâche terminée, tu DOIS mettre à jour taches.md en cochant [x] et notant la date/heure.

DÉMARRE MAINTENANT !


je veux que tu crée les endpoints dans api, de temporal et les workflows associé :
- "GET /api/exchange/binance/v1/orders (get all orders)"- "GET /api/exchange/binance/v1/quotes (get all crypto quotes)"

modifie l'api le worker et ajoute les workflows correspandante, normalement le activities sont deja crée, tu peux les utiliser, si tu ne tourve pas d'activité il faut crée dabord la fonction dans le clients et ensuite tu crée l'activité 

ensuite redemarre les service et fait un test E2E et regarde si tout est ok

si tu as des problemes a seconnecter à l'api de biance consulte la collextion postman/Afrikpay Gateway - API Clients E2E Complete .postman_collection et regarde dans les pre-request et test elle peuvent contenir de la logique qu'on doit suivre pour appeler les differents apis, utilise cette collection pour avoir la bonne methode a suivre

