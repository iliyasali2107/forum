CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content text NOT NULL,
    created DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
