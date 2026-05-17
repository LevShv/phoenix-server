CREATE TABLE IF NOT EXISTS meetings (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    presentation_url TEXT,
    speaker_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (speaker_id) REFERENCES users(id) ON DELETE CASCADE
)