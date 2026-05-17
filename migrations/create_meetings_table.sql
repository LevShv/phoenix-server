CREATE TABLE IF NOT EXISTS meetings (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    start_date TEXT NOT NULL, 
    speaker_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    status TEXT NOT NULL,
    avatar BLOB,
    presentation_url TEXT,
    FOREIGN KEY (speaker_id) REFERENCES users(id) ON DELETE CASCADE
);