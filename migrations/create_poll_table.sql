CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meeting_id TEXT NOT NULL,
    title TEXT NOT NULL,              
    url TEXT NOT NULL,                 
    created_at INTEGER NOT NULL,
    qr BLOB,
    FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_polls_meeting_id ON polls(meeting_id);