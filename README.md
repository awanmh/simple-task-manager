
# Simple Task Manager (Go + React)

![Go Version](https://img.shields.io/badge/Go-1.24-blue) ![React](https://img.shields.io/badge/React-TypeScript-blue) ![Docker](https://img.shields.io/badge/Docker-Ready-2496ED) ![License](https://img.shields.io/badge/License-MIT-green)

A robust, full-stack task management application designed to demonstrate **Clean Architecture**, **Scalability**, and **Modern DevOps practices**. This project implements a RESTful API using Golang and a reactive frontend using React (Vite).

## 🌟 Key Features

* **🔐 Authentication:** Secure JWT-based authentication (Login/Register) with Bcrypt password hashing.
* **🏗 Clean Architecture:** Strict separation of concerns (Domain, Usecase, Repository, Delivery) for maintainability and testability.
* **✅ Task & Subtasks:** Complete CRUD with progress tracking (percentage bar) for subtasks.
* **⏰ Real-time Reminders:** Browser-based alarm system for due tasks.
* **🔄 Recurring Tasks (Automation):** Backend **Cron Job** scheduler to automatically generate recurring tasks (Daily/Weekly).
* **🐳 Dockerized:** Fully containerized environment (Backend, Frontend, Database) using Docker Compose & Nginx.
* **📄 API Documentation:** Auto-generated Swagger/OpenAPI documentation.
* **🧪 Unit Testing:** Comprehensive unit tests for business logic using `Testify` and `Mockery`.

## 🛠 Tech Stack

### Backend
* **Language:** Golang 1.24
* **Framework:** Gin Gonic
* **Database Driver:** pgx/v5 (High performance PostgreSQL driver)
* **Scheduler:** Robfig Cron v3
* **Testing:** Testify & Mockery
* **Docs:** Swaggo

### Frontend
* **Framework:** React 19 (Vite)
* **Language:** TypeScript
* **Styling:** CSS Modules / Custom UI
* **HTTP Client:** Axios with Interceptors

### Infrastructure
* **Database:** PostgreSQL 15 (Alpine)
* **Containerization:** Docker & Docker Compose (Multi-stage builds)
* **Web Server:** Nginx (Reverse Proxy for Frontend)

---

## 🚀 Getting Started (The Fast Way)

Prerequisites: **Docker** and **Docker Compose** installed.

1.  **Clone the repository**
    ```bash
    git clone [https://github.com/username/simple-task-manager.git](https://github.com/username/simple-task-manager.git)
    cd simple-task-manager
    ```

2.  **Run with Docker Compose**
    This command will build the Go binary, compile React assets, and setup the Database automatically.
    ```bash
    docker-compose up --build
    ```

3.  **Access the Application**
    * Frontend: [http://localhost:3000](http://localhost:3000)
    * Backend API: `http://localhost:8080`
    * Swagger Docs: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🔧 Manual Installation (Development)

If you want to run it without Docker:

### 1. Database Setup
Ensure PostgreSQL is running and create a database named `taskdb`.
```bash
# Run migration manually (or execute init.sql)
psql -U postgres -d taskdb -f init.sql
```

### 2\. Backend Setup

```bash
# Create .env file
cp .env.example .env
```

```bash
# Install dependencies
go mod tidy
```
```bash
# Run Server
go run cmd/main.go
```
### 3. Frontend Setup

```bash
cd client
npm install
npm run dev
```
### 4. Docker Setup

```bash
docker-compose build
```

```bash
docker-compose up
```

-----

## 📂 Project Structure (Clean Architecture)

```
simple-task-manager/
├── client/
│   ├── node_modules/
│   ├── public/
│   ├── src/
│   │   └── assets/
├── cmd/
├── docs/
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   └── usecase/
│   │       └── mocks/
│   ├── infra/
│   │   ├── delivery/
│   │   │   ├── http/
│   │   │   └── middleware/
│   │   ├── repository/ 
│   │   └── scheduler/
│   ├── db/
│   │   └── migrations/
├── node_modules/
├── pkg/
│   └── util/
├── .dockerignore
├── .gitignore
├── .mockery.yaml
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.mod
├── init.sql
├── package-lock.json
├── package.json
├── README.md
└── go.sum
```

-----

## 🧪 Testing

This project uses **Unit Testing** for the Usecase layer, mocking the repository layer to ensure isolation.

To run the tests:

```bash
go test ./internal/core/usecase/... -v
```

**Coverage:**

  * User Registration logic (Duplicate email checks, Hashing)
  * Task Creation logic
  * Subtask logic

-----

## 📝 API Documentation

Full API documentation is available via Swagger UI.
Once the server is running, visit:
**[http://localhost:8080/swagger/index.html](https://www.google.com/search?q=http://localhost:8080/swagger/index.html)**

-----

## 🛡 License

This project is licensed under the MIT License.
