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

## Structure
- `internal/models/` : Modèles GORM (Module, Instance, Job)
- `internal/api/` : Routes et handlers API
- `cmd/main.go` : Point d’entrée du serveur

## À faire
- Implémenter la logique métier dans les handlers
- Ajouter la gestion de la base SQLite et des migrations
- Sécuriser l’API
