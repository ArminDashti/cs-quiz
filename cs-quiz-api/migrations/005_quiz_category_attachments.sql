ALTER TABLE quizzes
    ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT '';

ALTER TABLE quizzes
    ADD COLUMN IF NOT EXISTS attachments JSONB NOT NULL DEFAULT '[]'::jsonb;

UPDATE quizzes SET category = 'Language' WHERE slug IN ('csharp', 'python', 'js') AND category = '';
UPDATE quizzes SET category = 'Platform' WHERE slug IN ('dotnet') AND category = '';
UPDATE quizzes SET category = 'Database' WHERE slug IN ('sql-server') AND category = '';
