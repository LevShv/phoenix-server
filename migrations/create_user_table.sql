CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    name TEXT,
    surname TEXT,
    patronymic TEXT,
    avatar BLOB,
    password_hash TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN('speaker', 'listener'))
)