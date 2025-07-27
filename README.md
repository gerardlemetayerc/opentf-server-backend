# OpenTF Server (Go)

Ce projet expose une API REST pour gérer des modules OpenToFu/Terraform, leurs instances et les jobs associés.

## Démarrage rapide

1. Installez Go (https://golang.org/dl/)
2. Installez les dépendances :
   ```sh
   go mod tidy
   ```
3. Lancez le serveur :
   ```sh
   go run ./cmd/main.go
   ```

L’API sera disponible sur http://localhost:8080

## Présentation
OpenTF Server est une API backend écrite en Go, conçue pour la gestion centralisée de modules OpenToFu/Terraform, le suivi de leurs instances et l'orchestration des jobs associés. Elle facilite l'automatisation et la gouvernance de l'infrastructure as code dans des environnements multi-utilisateurs.

### Fonctionnalités principales
- Gestion des modules Terraform/OpenToFu (CRUD)
- Suivi des instances déployées
- Orchestration et suivi des jobs d'exécution
- Authentification locale, OIDC et par token API (JWT)
- Gestion des utilisateurs, groupes et méthodes d'authentification
- Système de valeurs suggérées dynamiques pour les propriétés de modules
- Stockage SQLite avec ORM GORM
- API sécurisée (Bearer JWT)

## État du projet
Le backend est fonctionnel et expose toutes les routes principales pour la gestion des modules, instances, jobs, utilisateurs et authentification. La gestion des sessions utilise des tokens JWT. La logique métier de base est en place, mais certaines fonctionnalités avancées (OIDC complet, gestion fine des dépendances, UI) restent à finaliser.

### Ce qui fonctionne
- CRUD complet sur les entités principales
- Authentification locale et par token (JWT)
- Sécurité de base via Bearer token
- Stockage et migration automatique de la base SQLite

### Points à améliorer / roadmap
- Finaliser l'intégration OIDC
- Ajouter des tests automatisés
- Améliorer la gestion des erreurs et la validation des entrées
- Ajouter une interface utilisateur (frontend)
- Documentation API détaillée (OpenAPI/Swagger)

## Structure
- `internal/models/` : Modèles GORM (Module, Instance, Job)
- `internal/api/` : Routes et handlers API
- `cmd/main.go` : Point d’entrée du serveur

## À faire
- Implémenter la logique métier dans les handlers
- Ajouter la gestion de la base SQLite et des migrations
- Sécuriser l’API
