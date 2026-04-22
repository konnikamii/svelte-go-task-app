# Copilot Instructions

> **Keep this file up to date.** Whenever there is an important change — new domain, auth change, new convention, restructure — update this file as part of the same task.

## Project Overview

Full-stack web application:

- `/backend` — Go REST API (chi router, sqlc, pgx, PostgreSQL)
- `/frontend` — Svelte frontend (SPA, no SSR)

---

## Backend Architecture

### Domain structure

Each feature lives in its own domain folder under `internal/`:

```
<domain>.handler.go   — HTTP handlers (thin, delegate to service)
<domain>.routes.go    — route registration, receives *pgxpool.Pool
<domain>.service.go   — business logic, uses repo.Queries + *pgxpool.Pool
<domain>.types.go     — request/response DTOs for this domain (optional)
```

Current domains: `auth`, `seed`, `tasks`, `users`.

Current domains: `auth`, `contact`, `seed`, `tasks`, `users`.

### Dependency direction

```
handler → service → repo (sqlc-generated)
```

Handlers must not access the database directly.

### Database layer

- PostgreSQL via `pgx/v5`
- `sqlc` for type-safe query generation
- SQL queries live in `internal/adapters/postgresql/sqlc/` (one `.sql` file per domain)
- Generated Go code lives in `internal/adapters/postgresql/sqlc/out/` — **do not edit generated files by hand unless sqlc cannot model the query**
- Migrations live in `internal/adapters/postgresql/migrations/` (goose format)
- After changing any `.sql` query file, run `sqlc generate` from `backend/`

### Route registration pattern

Every `Routes(r chi.Router, db *pgxpool.Pool)` function must:

1. Construct the repo via `repo.New(db)`
2. Construct the service with `repo.Queries` (and `*pgxpool.Pool` when raw queries or transactions are needed)
3. Construct the handler with the service
4. Register routes under the domain prefix, split into open and protected groups

Use a connection pool for application-wide database access. Do not share a single `*pgx.Conn` across HTTP requests because concurrent requests will hit `conn busy` failures.

### Protected vs open routes

Use `middleware.RequireAuth` (from `internal/middleware`) to guard route groups:

```go
r.Route("/tasks", func(r chi.Router) {
    r.Use(middleware.RequireAuth)   // all routes in this block require a session
    r.Post("/list", handler.List)
    r.Get("/{id}", handler.GetByID)
})

r.Route("/users", func(r chi.Router) {
    r.Post("/", handler.CreateUser) // open — registration

    r.Group(func(r chi.Router) {
        r.Use(middleware.RequireAuth) // protected group within the same prefix
        r.Get("/{id}", handler.GetByID)
    })
})
```

### Authentication (session-based)

- `internal/middleware/session.go` — session setup and helpers for DB-backed sessions. Cookie stores an opaque token; session records are stored in `user_sessions` table with device fingerprinting.
  - `DeviceIDFromRequest(r)` derives stable device identity from User-Agent + Accept-Language + Accept-Encoding headers (supports multi-device logins).
  - On login: revokes all **active sessions for this user on this device**, then creates a new session. This allows one active session per device per user without cross-device interference.
  - Background cleanup deletes expired and revoked-30-days-old sessions opportunistically.
- `internal/middleware/context.go` — `ContextWithUserID` / `UserIDFromContext` helpers.
- `internal/middleware/authenticate.go` — `RequireAuth` HTTP middleware. Reads `sid` cookie token → validates active non-revoked non-expired session from DB → injects user ID into context.
- Auth endpoints (all under `/api`):
  - `POST /login` — open; derives device_id, revokes old sessions on this device for this user, creates new session, sets `sid` cookie
  - `POST /logout` — open; revokes DB session and clears `sid` cookie
  - `GET /me` — protected; returns current user info
- Contact endpoint:
  - `POST /contact/` — open; stores a contact request with required `email`, `title`, and `message`; after save, the backend attempts a non-blocking SMTP notification (for local compose: MailHog)
- Seed endpoint:
  - `POST /seed` — open only in the sense that it does not require auth, but it must return `409` when any user already exists; it seeds roles, permissions, role mappings, demo users, user-role mappings, and demo tasks in one transaction for an empty database only

Cookie settings (set in `auth.Handler`):

- `HttpOnly: true`, `SameSite: Lax`, `Secure` driven by `COOKIE_SECURE` env var.

To read the authenticated user ID inside any protected handler:

```go
userID := middleware.UserIDFromContext(r.Context())
```

### Authorization (RBAC + ABAC scope)

- RBAC tables:
  - `roles`
  - `permissions` (`resource`, `action`, `scope` where scope is `any` or `own`)
  - `users_roles` (composite PK: `user_id`, `role_id`)
  - `roles_permissions` (composite PK: `role_id`, `permission_id`)
- Every authorization table includes `created_at` and `updated_at`.
- `admin` role is mapped to `scope = any` permissions.
- `user` role is mapped to `scope = own` permissions.
- ABAC ownership checks (`scope = own`) must compare the authenticated user ID from context against the resource owner ID in service layer.
- `tasks` ownership is tracked by `tasks.owner_id`.
- `tasks` supports optional `due_date`.
- New registrations should be assigned the `user` role in the same transaction as user creation.

### CORS

Configured in `cmd/api/api.go` using `github.com/go-chi/cors`. Must be the first middleware.

- Allowed origin: `FRONTEND_URL` env var (default `http://localhost:5173`)
- `AllowCredentials: true` — required for cross-origin session cookies

### Error handling

- `internal/apperrors/errors.go` — centralized app error constructors: `BadRequest`, `Unauthorized`, `Unauthenticated`, `Forbidden`, `NotFound`, `Conflict`, `Internal`
- `apperrors.FromError(err)` — unwraps `*AppError` or falls back to `Internal`
- In handlers use `h.AppError(w, err)` to send any service error (maps to correct HTTP status + `errorCode` field)
- Never expose raw DB or internal error strings to the client

### Paginated list endpoints

Use `POST /<domain>/list` (not `GET /`) for paginated, filtered, sorted lists.
Service method: `Get<Domain>sPaginated(ctx, PaginatedParams) (PaginatedReturn, error)`
Params struct carries: `Page`, `PageSize`, `SortBy`, `SortType`, `Filters`.

### Response types

- `repo.*` types (sqlc-generated) must NOT cross the service → handler boundary
- Each domain defines its own response struct (e.g. `UserResponse`, `UserInfo`) that omits sensitive fields like `Password`
- Map `repo.*` → response struct inside the service using a mapper function (e.g. `userToResponse`)

### Handler helpers (via `handlers.BaseHandler`)

| Method                                  | Purpose                             |
| --------------------------------------- | ----------------------------------- |
| `h.JSON(w, status, data)`               | Write success response              |
| `h.Error(w, status, msg)`               | Write plain error (no code)         |
| `h.ErrorWithCode(w, status, msg, code)` | Write error with machine code       |
| `h.AppError(w, err)`                    | Map `*AppError` or fall back to 500 |
| `h.Read(r, &body)`                      | Decode JSON body (strict)           |

---

## Do

- Place new features inside a domain module
- Write SQL in the sqlc query file, regenerate, then implement the service
- Use transactions (`db.Begin`) for any write that spans multiple queries
- Whitelist sort columns explicitly (never interpolate user input directly into SQL)
- Cast non-text columns before using `ILIKE` in dynamic queries (e.g. `CAST(id AS TEXT)`)
- Check `rows.Err()` after iterating `pgx.Rows`
- Propagate `tx.Commit` errors instead of ignoring them
- Use `middleware.RequireAuth` on every route group that requires authentication
- Return `apperrors.*` values from services for all user-facing errors

## Do not

- Do not add a `Repository` interface layer — use `repo.Queries` directly
- Do not use mock repositories in production code paths
- Do not place business logic in handlers or in `cmd/api`
- Do not expose internal error details (stack traces, DB errors) in HTTP responses
- Do not add extra abstraction layers, helpers, or utilities unless the task clearly requires them
- Do not guess SQL column names — check the migration file first
- Do not edit `out/*.sql.go` or `out/querier.go` by hand when `sqlc generate` can produce the correct output
- Do not use `GET` for endpoints that accept a request body
- Do not duplicate auth middleware in multiple packages; keep a single `RequireAuth` implementation in `internal/middleware`

---

## Environment variables

| Variable                   | Default                                     | Description                           |
| -------------------------- | ------------------------------------------- | ------------------------------------- |
| `SERVER_HOST`              | `localhost`                                 | Bind host                             |
| `SERVER_PORT`              | `8000`                                      | Bind port                             |
| `GOOSE_DBSTRING`           | —                                           | PostgreSQL DSN                        |
| `GOOSE_DRIVER`             | `postgres`                                  | Goose driver                          |
| `GOOSE_MIGRATION_DIR`      | `./internal/adapters/postgresql/migrations` | Migrations path                       |
| `FRONTEND_URL`             | `http://localhost:5173`                     | Allowed CORS origin                   |
| `COOKIE_SECURE`            | `false`                                     | Set `true` in production (HTTPS only) |
| `SESSION_SECRET`           | —                                           | Cookie session signing/encryption key |
| `SESSION_DURATION_MINUTES` | `1440`                                      | Session cookie max age in minutes     |
| `SMTP_HOST`                | `mailhog`                                   | SMTP host for contact notifications   |
| `SMTP_PORT`                | `1025`                                      | SMTP port for contact notifications   |
| `SMTP_FROM`                | `no-reply@taskify.local`                    | Sender for contact notifications      |
| `SMTP_TO`                  | `contact@taskify.local`                     | Recipient for contact notifications   |

---

## Frontend

- Located in `/frontend`, built with Svelte
- Communicates with the backend over HTTP
- Must send `credentials: 'include'` on all fetch/axios calls for the session cookie to be sent
- No SSR — pure SPA

---

## Running the project

```bash
# Backend (with live reload)
cd backend && air

# Frontend
cd frontend && pnpm dev

# Full stack with Docker Compose
docker compose up --build

# Regenerate sqlc types
cd backend && sqlc generate

# New migration
cd backend && goose create <name> -s sql
```

Docker notes:

- `docker-compose.yaml` builds `frontend` and `backend` from their local Dockerfiles and runs PostgreSQL as `postgres`
- The frontend build arg `PUBLIC_BACKEND_URL` should point to the browser-visible backend URL (for local compose: `http://localhost:8000`)
- PostgreSQL initializes schema from `backend/internal/adapters/postgresql/migrations/` on first startup of a fresh `postgres-data` volume

You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

## Available Svelte MCP Tools:

### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.
