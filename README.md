# Ticket System API Microservice

A backend-only RESTful microservice built in Golang for managing user authentication and support ticket lifecycle management, featuring an optional standalone Web Dashboard UI.

- **GitHub Repository**: [https://github.com/Shauryakant/eva-assignment.git](https://github.com/Shauryakant/eva-assignment.git)
- **Database**: MongoDB Atlas (`evabharat`)

---

## Technical Stack

- **Language**: Go 1.22+ (`net/http` + `github.com/go-chi/chi/v5`)
- **Database**: MongoDB (Official `go.mongodb.org/mongo-driver`)
- **Authentication**: JWT (HS256) via `github.com/golang-jwt/jwt/v5`
- **Password Hashing**: `bcrypt` (`golang.org/x/crypto/bcrypt`)
- **Configuration**: Environment variables (`.env` via `github.com/joho/godotenv`)
- **Frontend / Dashboard**: HTML5, CSS3, Vanilla JavaScript

---

## Data Model

### Users Collection (`users`)
```json
{
  "_id": "65f1a2b3c4d5e6f7a8b9c0d1",
  "email": "user@example.com",
  "password_hash": "$2a$10$...",
  "created_at": "2026-09-15T10:00:00Z"
}
```

### Tickets Collection (`tickets`)
```json
{
  "_id": "65f1a2b3c4d5e6f7a8b9c0d2",
  "user_id": "65f1a2b3c4d5e6f7a8b9c0d1",
  "title": "Database connection drop",
  "description": "App loses connection to MongoDB under high concurrency",
  "status": "open",
  "created_at": "2026-09-15T10:05:00Z",
  "updated_at": "2026-09-15T10:05:00Z"
}
```

---

## API Endpoint Contract

| Method | Endpoint | Auth Required | Purpose | Success Code |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | No | Public health check endpoint | `200 OK` |
| `POST` | `/auth/register` | No | Register new user (bcrypt password hashing) | `201 Created` |
| `POST` | `/auth/login` | No | Authenticate user & return JWT token | `200 OK` |
| `POST` | `/auth/logout` | No | Logout user & clear session | `200 OK` |
| `POST` | `/tickets` | Yes (`Bearer <token>`) | Create a new ticket owned by caller | `201 Created` |
| `GET` | `/tickets` | Yes (`Bearer <token>`) | List all tickets owned by caller | `200 OK` |
| `GET` | `/tickets/{id}` | Yes (`Bearer <token>`) | Fetch single ticket by ID | `200 OK` / `404` |
| `PATCH` | `/tickets/{id}/status` | Yes (`Bearer <token>`) | Update ticket status | `200 OK` / `400` / `404` |

---

## Status Transition Rules

The status of a ticket strictly follows a single-step forward progression:

`open` -> `in_progress` -> `closed`

- **`open` -> `in_progress`**: Allowed (`200 OK`).
- **`in_progress` -> `closed`**: Allowed (`200 OK`).
- **`open` -> `closed`**: Prohibited (`400 Bad Request`). Tickets cannot skip intermediate states.
- **`closed` -> `open` / `in_progress`**: Prohibited (`400 Bad Request`). Closed tickets are immutable and cannot be reopened.

---

## Security & Authorization Rules

1. **Bearer Token Authentication**: Protected endpoints require `Authorization: Bearer <token>`. Missing, expired, or invalid tokens return `401 Unauthorized`.
2. **Ownership Isolation**: Users can only read and modify tickets they created (`user_id == caller_id`).
3. **No Existence Leaks (Strict 404)**: Accessing non-owned ticket IDs returns `404 Not Found` (never `403`), avoiding resource enumeration.
4. **Password Security**: Passwords stored as `bcrypt` hashes (cost 10). Password hashes are marked `json:"-"` and never returned in API responses.

---

## Environment Variables

Configure `.env`:

```env
PORT=8080
MONGODB_URI=mongodb+srv://udemy:udemy123@cluster0.ywipqhb.mongodb.net/evabharat
DB_NAME=evabharat
JWT_SECRET=supersecretjwtkeychangeinproduction
```

---

## Local Setup & Run

### 1. Running with Go directly

```bash
git clone https://github.com/Shauryakant/eva-assignment.git
cd eva-assignment
go run main.go
```

Open your browser at **[http://localhost:8080](http://localhost:8080)** to access the dashboard.

### 2. Running Automated Terminal Demo Script

```bash
bash test_demo.sh
```

### 3. Running Unit Tests

```bash
go test -v ./...
```

---

## Assumptions & Design Decisions

1. **MongoDB Choice**: Aligns with JSON document structure, high read/write performance, and unique indexing (`users.email`).
2. **Sequential Status Enforcement**: Defaulted to strict step-by-step state transitions (`open` -> `in_progress` -> `closed`).
3. **HTTP 404 vs 403 on Non-Owned Tickets**: System intentionally returns `404 Not Found` to prevent resource enumeration.
4. **Minimal Dependencies**: Standard Go library with Chi router for lightweight binary size and clean code architecture.
