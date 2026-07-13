-- 001_initial_schema.sql
CREATE TABLE IF NOT EXISTS schema_migrations(
    version INTEGER PRIMARY KEY,
    applied_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS jobs(
    id                 TEXT PRIMARY KEY,
    created_at         INTEGER NOT NULL,
    updated_at         INTEGER NOT NULL,
    state              TEXT NOT NULL,
    source_image_path  TEXT NOT NULL DEFAULT '',
    product_hint       TEXT,
    tone               TEXT NOT NULL DEFAULT 'clear and practical',
    variant_count      INTEGER NOT NULL DEFAULT 4,
    seed               INTEGER,
    warning            TEXT,
    error              TEXT
);

CREATE TABLE IF NOT EXISTS artifacts(
    id          TEXT PRIMARY KEY,
    job_id      TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    kind        TEXT NOT NULL,
    seq         INTEGER NOT NULL DEFAULT 0,
    path        TEXT,
    payload     TEXT,
    seed        INTEGER,
    prompt      TEXT,
    style_ver   TEXT,
    created_at  INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_artifacts_job ON artifacts(job_id, kind, seq);
CREATE INDEX IF NOT EXISTS idx_jobs_state ON jobs(state, created_at);
