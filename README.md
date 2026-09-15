# Ticket System API Microservice

A backend-only RESTful microservice built in Golang for managing user authentication and support ticket lifecycle management. Designed and implemented for a technical take-home assignment following strict contract guidelines and endpoint specifications.

---

## Technical Stack

- **Language**: Go 1.22+ (`net/http` + `github.com/go-chi/chi/v5`)
- **Database**: MongoDB (Official `go.mongodb.org/mongo-driver`)
- **Authentication**: JWT (HS256) via `github.com/golang-jwt/jwt/v5`
- **Password Hashing**: `bcrypt` (`golang.org/x/crypto/bcrypt`)
- **Configuration**: Environment variables (`.env` via `github.com/joho/godotenv`)

---

## Data Model

### Users Collection (`users`)
```json
{
  "_id": "65f1a2b3c4d5e6f7a8b9c0d1",
  "email": "user@example.com",
  "password_hash": "$2a$10$...", // Omitted from all JSON responses
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
  "status": "open", // Supported: "open" | "in_progress" | "closed"
  "created_at": "2026-09-15T10:05:00Z",
  "updated_at": "2026-09-15T10:05:00Z"
}
```

---

## API Endpoint Contract

| Method | Endpoint | Auth Required | Purpose | Success Code |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | No | Health check endpoint | `200 OK` |
| `POST` | `/auth/register` | No | Register new user | `201 Created` |
| `POST` | `/auth/login` | No | Authenticate user & return JWT | `200 OK` |
| `POST` | `/tickets` | Yes (`Bearer <token>`) | Create a new ticket owned by caller | `201 Created` |
| `GET` | `/tickets` | Yes (`Bearer <token>`) | List all tickets owned by caller | `200 OK` |
| `GET` | `/tickets/{id}` | Yes (`Bearer <token>`) | Fetch single ticket by ID | `200 OK` / `404 Not Found` |
| `PATCH` | `/tickets/{id}/status` | Yes (`Bearer <token>`) | Update ticket status | `200 OK` / `400 Bad Request` / `404` |

---

## Status Transition Rules

The status of a ticket strictly follows a single-step forward progression:

$$\text{open} \xrightarrow{\quad} \text{in\_progress} \xrightarrow{\quad} \text{closed}$$

- **`open` $\rightarrow$ `in_progress`**: Allowed (`200 OK`).
- **`in_progress` $\rightarrow$ `closed`**: Allowed (`200 OK`).
- **`open` $\rightarrow$ `closed`**: Prohibited (`400 Bad Request`). Tickets cannot skip intermediate states.
- **`closed` $\rightarrow$ `open` / `in_progress`**: Prohibited (`400 Bad Request`). Closed tickets are immutable and cannot be reopened.

---

## Security & Authorization Rules

1. **Bearer Token Authentication**: Protected endpoints require `Authorization: Bearer <token>`. Missing, expired, or invalid tokens return `401 Unauthorized`.
2. **Ownership Isolation**: Users can only read and modify tickets they created (`user_id == caller_id`).
3. **No Existence Leaks (Strict 404)**: If a user attempts to access or update a ticket ID that exists but belongs to another user, the API responds with `404 Not Found` (identical to non-existent IDs), avoiding resource enumeration vectors.
4. **Password Security**: Passwords are saved strictly as `bcrypt` hashes with cost factor 10. Plain text passwords and password hashes are never logged or exposed in API responses.

---

## Environment Variables

Copy `.env.example` to `.env`:

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DB_NAME=ticket_system
JWT_SECRET=supersecretjwtkeychangeinproduction
```

---

## Local Setup & Run

### 1. Running with Go directly

```bash
# Clone the repository
git clone <repository-url>
cd ticket-system

# Install dependencies and run
go run main.go
```

### 2. Running with Docker

```bash
# Build Docker image
docker build -t ticket-system .

# Run container
docker run -p 8080:8080 -e MONGODB_URI="mongodb://host.docker.internal:27017" ticket-system
```

### 3. Running with Docker Compose (App + Local MongoDB)

```bash
docker-compose up --build
```

### 4. Verifying Health Check

```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

---

## Running Automated Tests

```bash
go test -v ./...
```

---

## Deployment Information

- **Deployed Application URL**: `https://ticket-system-golang.onrender.com`
- **Public Health Check URL**: `https://ticket-system-golang.onrender.com/health`

---

## Assumptions & Design Decisions

1. **MongoDB Choice**: MongoDB was selected for seamless alignment with JSON document structures, high read/write performance, and flexible indexing (such as unique constraints on `users.email`).
2. **Sequential Status Enforcement**: Defaulted to strict step-by-step state transitions (`open` -> `in_progress` -> `closed`). Jumping directly from `open` to `closed` is treated as an illegal state transition to prevent skipping resolution phases.
3. **HTTP 404 vs 403 on Non-Owned Tickets**: The system intentionally returns `404 Not Found` instead of `403 Forbidden` when accessing non-owned tickets to prevent attacker resource enumeration.
4. **Minimal Dependencies**: Standard Go library with Chi router was chosen over bloated frameworks to maintain fast startup, minimal binary size, and clean code layout.
