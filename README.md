# Broker Backend (Go)

A clean-architecture backend for a brokerage platform, written in Go 1.23.  
Features authentication with JWT, MySQL persistence, versioned migrations, and OpenAPI-described REST APIs for holdings & positions.

---

## ✨ Features

* Clean Architecture layers (handlers → services → domain ← infra).
* MySQL storage via `sqlx`.
* Database migrations with **golang-migrate** (auto-run on boot).
* Secure password hashing (bcrypt) & JWT auth (chi jwtauth).
* CRUD for:
  * Users (signup / login)
  * Holdings
  * Positions
* `.env` configuration (loaded with **godotenv**).  `.env` is in `.gitignore`.
* OpenAPI 3 specification (`openapi.yaml`).
* Extensive logging in service layer.

---

## 📂 Project structure (key parts)

```
cmd/            – entry-point (wires everything)
internal/
  domain/       – entities & repository interfaces
  infra/        – DB adapters (MySQL)
  services/     – application/business logic
  handlers/     – HTTP delivery (chi)
  migrations/   – SQL migration files
pkg/            – shared helpers (auth, ...)
```

---

## 🚀 Quick start

### 1. Prerequisites

* Go ≥ 1.23
* MySQL (local or AWS RDS)

### 2. Clone & setup

```bash
git clone <repo-url> broker-backend
cd broker-backend
cp env.example .env     # edit credentials & JWT_SECRET
```

Relevant `.env` keys:

| Key        | Example                        | Description                         |
|------------|--------------------------------|-------------------------------------|
| DB_HOST    | 127.0.0.1                      | MySQL host                          |
| DB_PORT    | 3306                           | MySQL port                          |
| DB_USER    | root                           | DB user                             |
| DB_PASS    | password                       | DB password                         |
| DB_NAME    | broker                         | Database name (auto-created)        |
| JWT_SECRET | super-secret-string            | HMAC key for JWT tokens             |

### 3. Run

```bash
go run ./cmd/server
```

* On first start the app will:
  1. Create the database (if missing).
  2. Apply SQL migrations under `internal/migrations`.
  3. Listen on `:8080`.

Logs look like:
```
2025/06/14 13:50:00 auth_service: user 1 signed up
2025/06/14 13:50:05 holding_service: created holding 3 for user 1 (AAPL x10.00)
```

---

## 🖥️ API overview

All routes are JSON; protected routes require
`Authorization: Bearer <token>`.

| Method | Path              | Description                            |
|--------|-------------------|----------------------------------------|
| POST   | `/signup`         | Register new user (email + password)   |
| POST   | `/login`          | Login, returns JWT token               |
| POST   | `/holdings/create`| Create holding                         |
| GET    | `/holdings`       | List holdings (DB)                     |
| POST   | `/positions`      | Create position                        |
| GET    | `/positions`      | List positions                         |

Full request/response schemas in **openapi.yaml**.

Render docs quickly:
```bash
npm i -g swagger-ui-watcher
swagger-ui-watcher openapi.yaml
# open http://localhost:3200
```

---

## 🛠️ Development

* `go vet ./...` – static analysis
* `go test ./...` – (tests TBD)
* Migrations: `migrate -path internal/migrations -database "mysql://user:pass@tcp(host:3306)/broker" up`

---

## 📄 License

MIT (see LICENSE) 