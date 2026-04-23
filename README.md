# Svelte / Go Task Manager App

<p align="center" style="display: flex; justify-content: center; align-items: center;">
	<a href="https://svelte.dev/" rel="noopener noreferrer" target="_blank">
		<img src="https://upload.wikimedia.org/wikipedia/commons/1/1b/Svelte_Logo.svg" height="120" alt="Svelte Logo">
	</a>
	<a href="https://github.com/konnikamii/" rel="noopener noreferrer" target="_blank" style="margin: 0px 20px 0px 20px ">
		<img src="https://upload.wikimedia.org/wikipedia/commons/6/6a/JavaScript-logo.png" width="40" alt="plus">
	</a>
	<a href="https://go.dev/" rel="noopener noreferrer" target="_blank">
		<img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" width="260" alt="Go Logo">
	</a>
</p>

<table align="center">
	<tr>
		<td>
			<a href="https://kit.svelte.dev/" rel="noopener noreferrer" target="_blank">
				<img src="https://raw.githubusercontent.com/sveltejs/branding/master/svelte-logo.png" height="50" alt="SvelteKit Logo">
			</a>
		</td>
		<td>
			<a href="https://vite.dev/" rel="noopener noreferrer" target="_blank">
				<img src="https://de.vitejs.dev/logo-with-shadow.png" height="50" alt="Vite Logo">
			</a>
		</td>
		<td>
			<a href="https://www.postgresql.org/" rel="noopener noreferrer" target="_blank">
				<img src="https://www.postgresql.org/media/img/about/press/elephant.png" height="50" alt="PostgreSQL Logo">
			</a>
		</td>
		<td>
			<a href="https://github.com/mailhog/MailHog" rel="noopener noreferrer" target="_blank">
				<img src="https://avatars.githubusercontent.com/u/10258541?s=200&v=4" height="50" alt="Mailhog Logo">
			</a>
		</td>
		<td>
			<a href="https://www.docker.com/" rel="noopener noreferrer" target="_blank">
				<img src="https://www.logo.wine/a/logo/Docker_(software)/Docker_(software)-Logo.wine.svg" height="90" alt="Docker Logo">
			</a>
		</td>
		<td>
			<a href="https://marketingplatform.google.com/about/analytics/" rel="noopener noreferrer" target="_blank">
				<img src="https://img.icons8.com/color/512/google-analytics.png" height="50" alt="GA4 Logo">
			</a>
		</td>
	</tr>
</table>

## Description

This repository contains a full-stack task management platform application.

### Technologies Used

- **Backend**: Go (Chi, pgx, sqlc)
- **Frontend**: SvelteKit + Vite written in TypeScript
- **Database**: PostgreSQL

### Additional features include:

- **Mailhog** - email testing tool;
- **Session-based authentication** with DB-backed sessions;
- **RBAC + ABAC authorization** for task access control;
- **Cookie consent banner** with optional **GA4** and **Simple Analytics** integration;
- **Dockerized**: easy setup and deployment.

---

## Prerequisites

Ensure you have the following installed on your machine:

- [Node.js](https://nodejs.org/) (version **22** or higher) (if running locally)
- [pnpm](https://pnpm.io/installation) (latest version) (if running locally)
- [Go](https://go.dev/dl/) (version **1.26.1** or higher) (if running locally)
- [PostgreSQL](https://www.postgresql.org/download/) (version **17** recommended for parity with Docker)
- [Docker](https://www.docker.com/) (latest version)
- [Docker Compose](https://docs.docker.com/compose/) (latest version)

---

## Getting Started

#### 1. Make a folder where you will store the code:

```bash
mkdir svelte-go-task-app
```

#### 2. Clone the repository in the folder of your choice:

```bash
git clone https://github.com/konnikamii/svelte-go-task-app.git .
```

#### 3. Copy the example environment files and configure them:

```bash
cd backend && cp .env.example .env
cd ../frontend && cp .env.example .env
cd ..
```

## Backend Setup

#### 1. Navigate to the backend directory:

```bash
cd backend
```

#### 2. Copy the example environment file and configure it:

```bash
cp .env.example .env
```

#### 3. Install Go dependencies:

```bash
go mod download
```

#### 4. Run database migrations (requires PostgreSQL):

```bash
goose up
```

#### 5. (Optional) Regenerate sqlc types after editing SQL files:

```bash
sqlc generate
```

#### 6. Start the backend server with live reload:

```bash
air
```

If you do not use `air`, run:

```bash
go run ./cmd/api
```

## Frontend Setup

#### 1. Navigate to the frontend directory:

```bash
cd frontend
```

#### 2. Copy the example environment file and configure it:

```bash
cp .env.example .env
```

#### 3. Install dependencies:

```bash
pnpm install
```

#### 4. Start the SvelteKit development server:

```bash
pnpm dev
```

#### 5. Build frontend for production:

```bash
pnpm build
```

## Docker Setup

#### 1. Navigate to the root directory:

```bash
cd ..
```

#### 2. Build and start Docker containers:

```bash
docker compose up --build
```

### Access the application:

By default:

- The **frontend** is available at [http://localhost:3000](http://localhost:3000)
- The **backend** is available at [http://localhost:8000](http://localhost:8000)
- **Mailhog UI** is available at [http://localhost:8025](http://localhost:8025)

Try creating an account and logging in. Then create tasks and test protected routes.

---

### Additional Information:

- **Mailhog**:
  Mailhog is included in the Docker setup to catch outgoing emails.

- **Database**:
  PostgreSQL is included in Docker setup and initialized using migration scripts under `backend/internal/adapters/postgresql/migrations/`.

- **Auth**:
  Session cookies are HTTP-only and tied to entries in the `user_sessions` table.

If you run into startup issues, verify `.env` values for:

- PostgreSQL host/user/password/database/port;
- frontend/backend URLs and ports;
- SMTP host/port for Mailhog;
- `SESSION_SECRET` in backend environment.

---

### Running Locally Without Docker

You need the following:

- **PostgreSQL**:
  Install PostgreSQL, create a database, and update `GOOSE_DBSTRING` in `backend/.env`.
  - Default port: `5432`
  - Ensure your local DB hostname and credentials are correct.

- **Mailhog** (optional):
  Install and configure Mailhog to run on the following ports (or use Docker):
  - SMTP: `1025`
  - HTTP: `8025`

If Mailhog is not configured, the contact notification attempt will fail gracefully and not crash the app.

---

### Testing & Checks

Backend:

```bash
cd backend
go test ./...
```

Frontend checks:

```bash
cd frontend
pnpm check
pnpm test
```

If checks fail, verify your local environment variables and service connectivity first.

#### Analytics (GA4 + Simple Analytics)

To connect analytics, update `frontend/.env`:

```env
PUBLIC_GA_MEASUREMENT_ID=G-XXXXXXXXXX
PUBLIC_SIMPLE_ANALYTICS_DOMAIN=example.com
PUBLIC_SIMPLE_ANALYTICS_SCRIPT_URL=https://scripts.simpleanalyticscdn.com/latest.js
```

Analytics scripts are loaded only after user consent via the cookie consent banner.

#### Helpful commands for PostgreSQL

**Windows:**

```bash
pg_ctl.exe register -N "PostgreSQL" -U "NT AUTHORITY\NetworkService" -D "C:\Program Files\PostgreSQL\[version]\data" -w
																																	# creates a service to start on boot
pg_ctl status -D "C:\Program Files\PostgreSQL\[version]\data"     # checks PostgreSQL process status
pg_ctl restart -D "C:\Program Files\PostgreSQL\[version]\data"    # restart PostgreSQL process

psql -p <port> -d <database-name> -U <user>                          # login as user in postgres database

\l+                                                                  # size of DBs
\du                                                                  # show all users
SELECT * FROM pg_user;                                                # show all users
CREATE DATABASE <database-name> TEMPLATE template0;                   # create DB from clean template
DROP DATABASE IF EXISTS <database-name>;                              # remove DB
```

**Linux:**

```bash
pg_ctl status -D /var/lib/postgresql/[version]/main                  # checks PostgreSQL process status
pg_ctl restart -D /var/lib/postgresql/[version]/main                 # restart PostgreSQL process
pg_ctl start -D /var/lib/postgresql/[version]/main                   # start PostgreSQL process
pg_ctl stop -D /var/lib/postgresql/[version]/main                    # stop PostgreSQL process
```

#### Helpful commands for Docker

```bash
docker compose up                                                    # builds images and starts containers
docker compose down                                                  # removes containers
docker compose config                                                # troubleshoot compose setup

docker compose up --build                                            # force image rebuilds
docker compose --project-name "my-app" up                           # set project name
docker compose -p "my-app" up                                       # shorthand for project name
docker compose -f <filename.yaml> up                                 # run a specific compose file

docker ps                                                            # list running containers
docker logs <container_name_or_id>                                   # check container logs
docker stats                                                         # live container resource utilization

docker exec -it <container_name_or_id> /bin/sh                       # enter container using sh
docker exec -it <container_name_or_id> bash                          # enter container using bash (if installed)
```
