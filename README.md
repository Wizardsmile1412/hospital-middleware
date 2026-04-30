# 🏥 Hospital Middleware

A backend middleware system that enables hospital staff to search for patient information across Hospital Information Systems (HIS).

Built with **Go**, **Gin**, **PostgreSQL**, **Docker**, and **Nginx** — following **SOLID principles**.

---

## 📋 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)
- [Production Notes](#production-notes)

---

## Overview

This system acts as a middleware layer between hospital staff and multiple Hospital Information Systems (HIS). It allows staff members to securely search for patient records, with strict access control ensuring staff can only view patients from their own hospital.

**Key design decisions:**
- Patient data is stored in the system's own PostgreSQL database for fast, reliable queries
- Hospital staff are scoped to their hospital via JWT tokens — enforced at the service layer
- Searching by `national_id` or `passport_id` triggers a real-time lookup against the Hospital's external API, with automatic fallback to the local database if the external API is unavailable or returns no result
- All other searches query the local database directly

---

## Tech Stack

| Technology | Purpose |
|---|---|
| Go | Primary language |
| Gin | HTTP web framework |
| PostgreSQL | Relational database |
| golang-migrate | Database migrations |
| Docker + Docker Compose | Containerization |
| Nginx | Reverse proxy |
| JWT | Access token authentication |
| HttpOnly Cookie | Refresh token storage |
| bcrypt | Password hashing |
| testify | Unit testing |

---

## Architecture

```
Client
  │
  ▼
Nginx (reverse proxy :80)
  │
  ▼
Go + Gin HTTP Server (:8080)
  ├── Handler Layer    → Parse HTTP request, return response
  ├── Middleware Layer → JWT auth, logging
  ├── Service Layer    → Business logic + hospital isolation
  ├── Repository Layer → PostgreSQL queries
  └── Client Layer     → External Hospital API calls
  │
  ▼
PostgreSQL
```

### Data Strategy

```
Patient Search Flow:

Search by national_id / passport_id
  └──► Call Hospital External API (real-time)
         ├── Result found   → return immediately
         └── Error / no result → fallback: query local DB (error logged silently)

Search by name / DOB / phone / email
  └──► Query patients table in our PostgreSQL DB
         (fast, reliable, supports complex filters)
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose installed
- Git

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/hospital-middleware.git
cd hospital-middleware
```

### 2. Set up environment variables

```bash
cp .env.example .env
# Edit .env with your values
```

### 3. Start the services

```bash
docker compose up --build
```

This will start:
- **Nginx** on port `80`
- **Go API server** on port `8080` (internal)
- **PostgreSQL** on port `5432` (internal)

Database migrations and seed data run automatically on startup.

### 4. Verify it's running

```bash
curl -X POST http://localhost/api/v1/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","hospital_slug":"hospital-a"}'
# Expected: 201 Created
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values:

```env
# App
APP_ENV=development          # Set to "production" for secure cookies

# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hospital_middleware

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_ACCESS_EXPIRY=900        # 15 minutes in seconds
JWT_REFRESH_EXPIRY=604800    # 7 days in seconds

# Hospital A External API
HOSPITAL_A_API_URL=https://hospital-a.api.co.th
```

---

## API Endpoints

Base URL: `http://localhost/api/v1`

### Auth

| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| POST | `/staff/create` | Register new staff member | No |
| POST | `/staff/login` | Login and get tokens | No |
| POST | `/staff/refresh` | Refresh access token | Cookie |

### Patient

| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| GET | `/patient/search` | Search patients | Yes (Bearer token) |

---

### POST /staff/create

```bash
curl -X POST http://localhost/api/v1/staff/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "nurse_somjai",
    "password": "secret123",
    "hospital_slug": "hospital-a"
  }'
```

### POST /staff/login

```bash
curl -X POST http://localhost/api/v1/staff/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "nurse_somjai",
    "password": "secret123",
    "hospital_slug": "hospital-a"
  }'
```

### GET /patient/search

```bash
curl -X GET "http://localhost/api/v1/patient/search?first_name=Somchai&last_name=Jaidee" \
  -H "Authorization: Bearer <access_token>"
```

Search by national ID (triggers real-time Hospital API call):

```bash
curl -X GET "http://localhost/api/v1/patient/search?national_id=1234567890123" \
  -H "Authorization: Bearer <access_token>"
```

---

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific package tests
go test -v ./internal/service/...
go test -v ./internal/handler/...
```

---

## Project Structure

```
hospital-middleware/
├── cmd/
│   └── main.go                     # Entry point
├── internal/
│   ├── handler/                    # HTTP layer
│   │   ├── staff_handler.go
│   │   └── patient_handler.go
│   ├── routes/                     # Route registration
│   │   ├── routes.go               # Register() — wires all route groups
│   │   ├── staff_routes.go         # /staff group
│   │   └── patient_routes.go       # /patient group (auth middleware applied)
│   ├── service/                    # Business logic
│   │   ├── staff_service.go
│   │   └── patient_service.go
│   ├── repository/                 # Database access
│   │   ├── staff_repo.go
│   │   └── patient_repo.go
│   ├── middleware/
│   │   └── auth.go                 # JWT middleware
│   ├── token/
│   │   └── jwt.go                  # JWT generation and validation
│   ├── model/                      # Data structs
│   │   ├── staff.go
│   │   ├── patient.go
│   │   └── hospital.go
│   ├── client/                     # External API client
│   │   ├── hospital_client.go      # Interface
│   │   └── mock_hospital_client.go # Mock for tests
│   └── config/
│       └── config.go
├── migrations/                     # SQL migration files
│   ├── 000001_create_hospitals.{up,down}.sql
│   ├── 000002_create_staff.{up,down}.sql
│   ├── 000003_create_patients.{up,down}.sql
│   └── 000004_seed_data.{up,down}.sql
├── docker-compose.yml
├── Dockerfile
├── nginx.conf
├── Makefile
├── .env.example
└── README.md
```

---

## Production Notes

> ⚠️ **Sync Strategy**
>
> In this assignment, patient data is pre-seeded into the database for demonstration purposes.
>
> In a **real production system**, a scheduled background sync job should:
> - Periodically pull patient data from each Hospital's external API
> - Update the `patients` table with new or changed records
> - Handle failures gracefully with retry logic
>
> This keeps search fast and available even if a Hospital API goes down temporarily.
> The current architecture is intentionally designed to support this — the `HospitalClient` interface
> makes it easy to add a sync service without changing any existing code.

---

## SOLID Principles

This project is designed with SOLID principles throughout:

| Principle | Implementation |
|---|---|
| **S**ingle Responsibility | Each layer (handler/service/repository) has exactly one job |
| **O**pen/Closed | New hospitals only require a new `HospitalClient` implementation |
| **L**iskov Substitution | `MockHospitalClient` and real client are interchangeable |
| **I**nterface Segregation | `StaffRepository` and `PatientRepository` are separate interfaces |
| **D**ependency Inversion | All dependencies injected via constructors using interfaces |

---

## License

This project is created as a candidate assignment for Agnos Health.
