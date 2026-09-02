INSERT INTO quizzes (name, slug, description)
VALUES
    ('C#', 'csharp', 'Core C# language fundamentals'),
    ('.NET', 'dotnet', 'ASP.NET Core and .NET platform'),
    ('Python', 'python', 'Python language and standard library'),
    ('SQL Server', 'sql-server', 'T-SQL and SQL Server concepts')
ON CONFLICT (slug) DO NOTHING;

-- C#
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order
FROM quizzes q
CROSS JOIN (VALUES
    ('What keyword declares a constant in C#?', 'const', 'readonly', 'static', 'fixed', 0, 1),
    ('Which type is a reference type?', 'int', 'bool', 'string', 'decimal', 2, 2),
    ('What does async methods typically return?', 'void only', 'Task or Task<T>', 'Thread', 'IEnumerable', 1, 3),
    ('Which access modifier is most restrictive?', 'public', 'internal', 'protected', 'private', 3, 4)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'csharp'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id);

-- .NET
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order
FROM quizzes q
CROSS JOIN (VALUES
    ('Which package builds ASP.NET Core web APIs?', 'System.Web', 'Microsoft.AspNetCore.App', 'WinForms', 'WPF', 1, 1),
    ('What registers services for DI in Program.cs?', 'builder.Services', 'app.Use', 'MapGet', 'AddControllers only', 0, 2),
    ('EF Core is primarily used for?', 'UI rendering', 'ORM / data access', 'CSS bundling', 'DNS lookup', 1, 3),
    ('Minimal APIs typically map routes with?', 'MapGet/MapPost', 'Page_Load', 'Global.asax', 'WebForms', 0, 4)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'dotnet'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id);

-- Python
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order
FROM quizzes q
CROSS JOIN (VALUES
    ('Which creates a virtual environment (3.3+)?', 'pip freeze', 'python -m venv', 'npm init', 'go mod', 1, 1),
    ('What is the correct list comprehension for squares?', '[x^2 for x in range(5)]', '[x**2 for x in range(5)]', '{x**2 for x}', '(x**2)', 1, 2),
    ('Which keyword handles exceptions?', 'catch', 'except', 'rescue', 'trap', 1, 3),
    ('What does PEP 8 primarily cover?', 'Networking', 'Style guide', 'Packaging only', 'Async IO only', 1, 4)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'python'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id);

-- SQL Server
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order
FROM quizzes q
CROSS JOIN (VALUES
    ('Which clause filters rows before grouping?', 'HAVING', 'WHERE', 'ORDER BY', 'GROUP BY', 1, 1),
    ('What does NOLOCK roughly mean?', 'Serializable', 'Read uncommitted / dirty reads allowed', 'Exclusive lock', 'Snapshot only', 1, 2),
    ('Which statement creates a stored procedure?', 'CREATE PROC / PROCEDURE', 'MAKE SP', 'NEW PROCEDURE', 'DEFINE PROC', 0, 3),
    ('What is a clustered index typically tied to?', 'Every nonclustered key', 'The physical row order of the table', 'Only views', 'Only temp tables', 1, 4)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'sql-server'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id);
