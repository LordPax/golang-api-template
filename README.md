# Golang api template

## Description

Ce projet est un template d'API REST en Go.

## Technologies utilisées

* [Go](https://golang.org/)
* [Swag](https://github.com/swaggo/swag)
* [docker](https://www.docker.com/)
* [docker-compose](https://docs.docker.com/compose/)

## Initialisation du projet

1. Clonez le dépôt :
```bash
git clone https://github.com/LordPax/golang-api-template.git
```

2. Accédez au répertoire du projet :
```bash
cd golang-api-template
```

3. Lancer les conteneur docker
```bash
docker-compose up
```

## Installation back

1. Créez un fichier `.env` à la racine du répertoire `back` et ajoutez les variables d'environnement suivantes :
```bash
NAME=Golang Api
DOMAIN=localhost:8080
GIN_MODE=debug
ALLOWED_ORIGINS='*'

DB_HOST=localhost
DB_USER=root
DB_PASSWORD=root
DB_NAME=golang-app
DB_PORT=5432

COOKIE_SECURE=false
JWT_SECRET_KEY=secret
BREVO_API_KEY=
BREVO_SENDER=noreply@example.fr

CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
```

2. Intaller les dépendances :
```bash
go mod download
go mod vendor
swag init
```

3. Build le projet :
```bash
go build
```

4. Lancer les migrations et les fixtures:
```bash
./golang-api migrate
./golang-api fixtures
```

5. Lancer le projet
```bash
./golang-api
```
