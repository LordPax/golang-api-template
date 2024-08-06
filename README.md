# UrEsport

## Description

UrEsport est une plateforme de gestion de tournois de jeux vidéos. Elle permet aux joueurs de s'inscrire à des tournois, de consulter les résultats et de suivre l'actualité des tournois.

## Technologies utilisées

* [Go](https://golang.org/)
* [Swag](https://github.com/swaggo/swag)
* [docker](https://www.docker.com/)
* [docker-compose](https://docs.docker.com/compose/)

## Initialisation du projet

1. Clonez le dépôt :
```bash
git clone https://github.com/Pietrucci-Blacher/UrEsport.git
```

2. Accédez au répertoire du projet :
```bash
cd UrEsport
```

3. Lancer les conteneur docker
```bash
docker-compose up
```

## Installation back

1. Accédez au répertoire du projet :
```bash
cd back
```

2. Créez un fichier `.env` à la racine du répertoire `back` et ajoutez les variables d'environnement suivantes :
```bash
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
BREVO_SENDER=noreply@uresport.fr

CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
```

3. Intaller les dépendances :
```bash
go mod tidy
go mod vendor
swag init
```

4. Build le projet :
```bash
go build
```

5. Lancer les migrations et les fixtures:
```bash
./golang-api migrate
./golang-api fixtures
```

6. Lancer le projet
```bash
./golang-api
```
