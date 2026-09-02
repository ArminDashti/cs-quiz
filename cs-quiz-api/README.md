# cs-quiz-api

Go (Gin) + PostgreSQL API for computer science quizzes.

## Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/health` | — | Liveness |
| POST | `/api/v1/auth/register` | invite code | Sign up |
| POST | `/api/v1/auth/login` | — | Login |
| GET | `/api/v1/auth/me` | JWT | Current user |
| GET | `/api/v1/quizzes` | — | List quizzes |
| GET | `/api/v1/quizzes/:slug/play` | JWT | Questions (no answers) |
| POST | `/api/v1/quizzes/:slug/submit` | JWT | Grade attempt |
| * | `/api/v1/admin/...` | JWT + admin | Manage quizzes/questions |

## Run

```bash
# Postgres (Docker)
docker compose up -d

cp .env.example .env
go mod tidy
go run ./cmd/server
```

Default listen: `:8090` (avoids clashing with Geoquiz on `:8080`). Invite code: `csquiz`.
Default login (seeded on startup): username `armin` / password `dopadopa123` (admin).
`ADMIN_EMAIL` defaults to `armin@local` (matches the seeded account).

Seeded quiz topics: `csharp`, `dotnet`, `python`, `js`, `sql-server` (migrations `002_seed.sql`, `003_seed_js.sql`).
