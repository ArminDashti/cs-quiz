# cs-quiz-webui

Vue 3 + Vite + Tailwind (shadcn-style) WebUI for computer science quizzes. Uses **Inter** as the UI font.

## Pages

- `/` — home hero + quiz cards
- `/about-me`
- `/account`, `/profile`, `/profile/:username`
- `/quiz/csharp`, `/quiz/dotnet`, `/quiz/python`, `/quiz/js`, `/quiz/sql-server`
- `/management`, `/management/add-quiz`
- `/management/quiz/:quizName/add-question`, `/management/quiz/:quizName/questions-list`

## Run

```bash
cp .env.example .env
npm install
npm run dev
```

API default: Vite proxies `/api` to `http://127.0.0.1:8090` when `VITE_API_BASE_URL` is empty.

Default login (seeded by the API): username `armin` / password `dopadopa123`.
