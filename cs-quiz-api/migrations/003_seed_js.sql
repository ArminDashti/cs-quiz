INSERT INTO quizzes (name, slug, description)
VALUES
    ('JavaScript', 'js', 'Modern JavaScript language fundamentals')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order
FROM quizzes q
CROSS JOIN (VALUES
    ('Which keyword declares a block-scoped variable?', 'var', 'let', 'function', 'with', 1, 1),
    ('What does === compare?', 'Value only', 'Value and type', 'References only', 'Prototypes only', 1, 2),
    ('Which creates a Promise that resolves after all settle?', 'Promise.race', 'Promise.all', 'Promise.any', 'Promise.resolve only', 1, 3),
    ('What does Array.prototype.map return?', 'The same array mutated', 'A new array of mapped values', 'A boolean', 'An iterator only', 1, 4)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'js'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id);
