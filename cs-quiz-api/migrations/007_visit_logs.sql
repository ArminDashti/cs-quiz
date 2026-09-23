CREATE TABLE IF NOT EXISTS visit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ip_address TEXT NOT NULL,
    operating_system TEXT NOT NULL,
    browser TEXT NOT NULL,
    path TEXT NOT NULL,
    visited_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visit_logs_visited_at ON visit_logs (visited_at DESC);
