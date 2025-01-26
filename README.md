# Project Setup Guide

This document provides step-by-step instructions to set up and run the project, including setting up the PostgreSQL database, backend, and frontend.

---

## Prerequisites

Ensure you have the following installed on your machine:

1. **Git**: For cloning the repository.
2. **Go**: To run the backend server (version 1.18+ recommended).
3. **Node.js & npm**: To run the frontend (Node.js 16+ recommended).
4. **PostgreSQL**: To set up the database.

---

## Steps to Set Up the Project

### 1. Clone the Repository

```bash
git clone https://github.com/Choy050823/CVWO-Travel-Web-Forum.git
cd CVWO-Travel-Web-Forum
```

### 2. Install Dependencies

#### For the Backend
Navigate to the backend folder:

```bash
cd go-backend
```

Install Go dependencies:

```bash
go mod tidy
```

#### For the Frontend
Navigate to the frontend folder:

```bash
cd ../reactTS-frontend
```

Install Node.js dependencies:

```bash
npm install
```

---

### 3. Set Up PostgreSQL Database

#### Create a Database
1. Open PostgreSQL in your terminal or preferred client.
2. Create a new database:

```sql
CREATE DATABASE CVWO_Travel_Web_Forum;
```

3. (Optional) Create a dedicated user:

```sql
CREATE USER project_user WITH PASSWORD 'securepassword';
GRANT ALL PRIVILEGES ON DATABASE project_db TO project_user;
```

#### Run SQL Scripts
Navigate to the folder containing the `schema.sql` file and execute it:

```bash
psql -U project_user -d project_db -f /path/to/schema.sql
```

This script will set up the necessary tables and relationships for the database.

---

### 4. Configure Environment Variables

#### Backend
Create a `.env` file in the `backend` directory and set the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=project_user
DB_PASSWORD=securepassword
DB_NAME=project_db
JWT_SECRET=your_jwt_secret_key
```

#### Frontend
If necessary, configure API endpoints or other variables in the `frontend/.env` file:

```
REACT_APP_API_URL=http://localhost:8080
```

---

### 5. Run the Project

#### Start the Backend
From the `backend` directory:

```bash
go run main.go
```

The backend server will start at `http://localhost:8080`.

#### Start the Frontend
From the `frontend` directory:

```bash
npm run dev
```

The frontend application will be available at `http://localhost:5173`.

---

### Additional Notes

1. **Frontend Dependencies**: The project uses tailwind CSS as styling, as listed in the `package.json`.
2. **Backend WebSocket**: The backend integrates Gorilla WebSocket for real-time communication alongside REST APIs.
3. **Database Schema**: Ensure the `schema.sql` file is executed properly to avoid runtime issues.

---

You're now ready to run and develop the project! For any issues, feel free to open an issue in the repository or contact the maintainers.

