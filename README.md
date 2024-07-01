# Projet_Forum
### README pour le lancement du projet Forum KOKO GOAT.

---

#### 1. **Introduction**

Ce projet vise à créer un forum de discussion dédié aux passionnés de cinéma, en utilisant une API en Go, un backend en Node.js avec des templates EJS et une base de données SQL.

---

#### 2. **Prérequis**

Assurez-vous d'avoir les logiciels suivants installés :

- **Go** (Golang)
- **Node.js** et **npm**
- **MySQL** ou une autre base de données SQL compatible
- **Git**

---

#### 3. **Configuration de l'Environnement**

1. **Clonez le dépôt du projet :**
   ```bash
   git clone <URL_DU_DEPOT>
   cd <REPERTOIRE_DU_DEPOT>
   ```

2. **Configuration de la base de données :**
   - Installez MySQL.
   - Créez une nouvelle base de données :
     ```sql
     CREATE DATABASE forum_cinema;
     ```
   - Créez un utilisateur et accordez-lui les droits nécessaires :
     ```sql
     CREATE USER 'forum_user'@'localhost' IDENTIFIED BY 'password';
     GRANT ALL PRIVILEGES ON forum_cinema.* TO 'forum_user'@'localhost';
     FLUSH PRIVILEGES;
     ```

---

#### 4. **API en Go**

1. **Installation des dépendances :**
   ```bash
   cd api
   go mod tidy
   ```


#### 5. **Backend en Node.js**

1. **Installation des dépendances :**
   ```bash
   cd Servweb_js
   npm install
   ```

2. **Configuration :**
   - Créez un fichier `.env` dans le répertoire `GO_FORUM_API` avec le contenu suivant :
     ```env
     DB_USER=your_db_user
     DB_PASSWORD=your_password
     DB_HOST=localhost
     DB_PORT=3306
     DB_NAME=your_db_name
     ```



#### 6. **Templates en EJS**

1. **Structure du répertoire :**
   - Les templates EJS sont placés dans le répertoire `templates` à l'intérieur du répertoire `Servweb_js`.
   -  structure:
     ```
 Servweb_js/
     ├── backend/
     │   ├── src/
     │   │     ├── routes/
     │   └── serv.js 
     ├── frontend/
         ├── assets/
         ├── templates/
        
     ```
---

#### 7. **Base de Données**

1. **Création de la base de  données  :**

---

#### 8. **Lancement Complet du Projet**

1. **Lancer l'API en Go :**
   ```bash
   cd GO_FORUM_API
   go run .
   ```

2. **Lancer le backend en Node.js :**
   ```bash
   cd Servweb_js
   npm start
   ```

3. **Accéder à l'application :**
   - Ouvrez votre navigateur et accédez à `http://localhost:3000`.

---
