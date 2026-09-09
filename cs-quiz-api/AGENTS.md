# AGENTS.md

## Cursor Cloud specific instructions

`cs-quiz-api` is a single Go (Gin) + PostgreSQL REST API (no frontend in this repo). The only
runtime dependencies are the Go toolchain (installed) and a reachable PostgreSQL instance.

### Services

| Service | How to run | Port | Notes |
|---------|-----------|------|-------|
| PostgreSQL 16 | `sudo pg_ctlcluster 16 main start` | 5432 | Installed via apt (no Docker in this VM). Role/DB `csquiz`/`csquiz` already created. |
| cs-quiz-api (Go server) | `go run ./cmd/server` | 8090 | Reads `.env`; runs SQL migrations + seeds on startup. |

### Startup notes (non-obvious)

- **Postgres is native, not Docker.** The repo's `docker-compose.yml` is not used here (Docker is
  not installed, and its network is declared `external`). Start Postgres with
  `sudo pg_ctlcluster 16 main start` and confirm with `sudo pg_lsclusters`. It does not auto-start
  on VM boot.
- **`pgcrypto` extension is pre-created** in the `csquiz` database. The app's migration
  (`migrations/001_init.sql`) runs `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`, which needs
  superuser; the DB owner role `csquiz` is not superuser, so the extension was created once as the
  `postgres` superuser. If the database is ever recreated from scratch, re-run
  `sudo -u postgres psql -d csquiz -c 'CREATE EXTENSION IF NOT EXISTS "pgcrypto";'` before starting
  the server, otherwise migrations fail.
- **`.env` is required.** Copy it once with `cp .env.example .env`. Config is loaded from real
  environment variables first and `.env` only fills gaps (shell/exported vars win over `.env`).
- **Migrations are idempotent** (all `CREATE TABLE IF NOT EXISTS` / seeds use `ON CONFLICT`), so the
  server can be restarted freely against the same database.

### Auth / test flow

- Default local accounts (seeded on every startup):
  - username `armin` / password `dopadopa123` (admin)
  - username `admin` / password `dopadopa123` (admin)
- External ingest: `POST /api/v1/external/quizzes` with header `X-API-Key: <EXTERNAL_API_KEY>`.
- Registration is gated by an invite code (default `csquiz`, from `INVITE_CODE`).
- Whoever registers/logs in with `ADMIN_EMAIL` (default `armin@local`) is granted admin.
- Quick end-to-end check: `GET /health` → login as `armin` → `GET /api/v1/quizzes/csharp/play`
  (JWT) → `POST /api/v1/quizzes/csharp/submit`.

## Learned User Preferences

- Always use default login username `armin` and password `dopadopa123` for local/demo api-webui auth seeding and sign-in.
### Lint / build / test

- Lint: `go vet ./...`
- Build: `go build ./cmd/server`
- Tests: none exist in this repo (`go test ./...` is a no-op).
