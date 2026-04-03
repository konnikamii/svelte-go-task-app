# Project Architecture Guide (for Copilot & Contributors)

## Project Overview

This repository contains a full-stack web application:

- `/backend` — Go API server
- `/frontend` — Svelte frontend application

Currently the backend is the primary focus.

---

# Backend Architecture (Go)

Backend follows a **domain-based modular structure**.

Root folder:

```
/backend
```

Important directories:

```
backend/
│
├── cmd/
│   └── api/
│       ├── main.go     # application entrypoint
│       └── api.go      # server setup and bootstrapping
│
├── internal/
│   ├── adapters/       # external integrations
│   │   └── postgresql/
│   │       └── sqlc/
│   │           ├── migrations/  # database migrations
│   │           ├── queries.sql  # SQL queries used by sqlc
│   │           └── out/         # generated sqlc code
│   │
│   ├── auth/           # authentication logic
│   │
│   ├── env/            # environment variable helpers
│   │
│   ├── tasks/          # tasks domain module
│   │   ├── tasks.handler.go
│   │   ├── tasks.routes.go
│   │   └── tasks.service.go
│   │
│   ├── users/          # users domain module
│   │   ├── users.handler.go
│   │   ├── users.routes.go
│   │   └── users.service.go
```

---

# Backend Design Principles

### Domain-based modules

Each domain (tasks, users, etc.) contains:

```
<domain>.handler.go   -> HTTP handlers
<domain>.routes.go    -> route registration
<domain>.service.go   -> business logic
```

Handlers should remain thin and delegate logic to services.

---

### Services

Services contain the core business logic and interact with repositories.

Repositories are generated using **sqlc**.

---

### Database Layer

Database access uses:

- PostgreSQL
- sqlc for type-safe queries

Location:

```
internal/adapters/postgresql/sqlc/
```

Structure:

```
migrations/    -> schema migrations
queries.sql    -> SQL queries
out/           -> generated Go code
```

Services should call the generated **sqlc repository interfaces**.

---

### Application Startup

Application entry point:

```
cmd/api/main.go
```

Server configuration and wiring:

```
cmd/api/api.go
```

These files are responsible for:

- loading environment variables
- connecting to the database
- initializing services
- registering routes
- starting the HTTP server

---

# Frontend Architecture

Frontend is located in:

```
/frontend
```

Technology:

- Svelte

The frontend communicates with the Go API via HTTP.

---

# Conventions

### Naming

Domain files follow this naming pattern:

```
tasks.handler.go
tasks.routes.go
tasks.service.go
```

### Dependency direction

```
handlers → services → repositories (sqlc)
```

Handlers must not directly access the database.

---

# Notes for Copilot

When generating code:

- new features should be placed inside a **domain module**
- business logic belongs in **services**
- HTTP request logic belongs in **handlers**
- routes are registered in `<domain>.routes.go`
- database queries should be written in `queries.sql` and generated using sqlc
- migrations go into `migrations/`

Avoid placing application logic in `cmd/api`.
