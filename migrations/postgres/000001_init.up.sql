CREATE TABLE IF NOT EXISTS url_mappings (
    id VARCHAR(255) PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    created_at BIGINT NOT NULL
);

CREATE INDEX idx_long_url ON url_mappings(long_url);
CREATE INDEX idx_short_url ON url_mappings(short_url);
CREATE INDEX idx_created_at ON url_mappings(created_at);