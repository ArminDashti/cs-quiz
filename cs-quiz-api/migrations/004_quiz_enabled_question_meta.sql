ALTER TABLE quizzes
    ADD COLUMN IF NOT EXISTS enabled BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE questions
    ADD COLUMN IF NOT EXISTS question_type TEXT NOT NULL DEFAULT 'multiple_choice';

ALTER TABLE questions
    ADD COLUMN IF NOT EXISTS difficulty TEXT NOT NULL DEFAULT 'medium';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'questions_difficulty_check'
    ) THEN
        ALTER TABLE questions
            ADD CONSTRAINT questions_difficulty_check
            CHECK (difficulty IN ('easy', 'medium', 'hard'));
    END IF;
END $$;
