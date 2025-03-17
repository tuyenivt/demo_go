CREATE TABLE IF NOT EXISTS patients (
    id VARCHAR(255) PRIMARY KEY,
    data JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS api_keys (
    key VARCHAR(255) PRIMARY KEY,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

-- For testing purposes
INSERT INTO api_keys (key, active) VALUES ('test-api-key-1', TRUE) ON CONFLICT (key) DO NOTHING;
