-- uploads table schema (SQLite)
CREATE TABLE IF NOT EXISTS uploads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_file TEXT NOT NULL,
    channel TEXT NOT NULL,
    status TEXT NOT NULL,
    telegram_file_id TEXT,
    telegram_message_id INTEGER,
    uploaded_at TEXT,
    UNIQUE(book_file, channel)
);
CREATE INDEX IF NOT EXISTS idx_uploads_book_file ON uploads(book_file);
CREATE INDEX IF NOT EXISTS idx_uploads_channel ON uploads(channel);
