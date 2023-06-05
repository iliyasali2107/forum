CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE NOT NULL
);

INSERT OR IGNORE INTO categories (name) VALUES ("basketball");
INSERT OR IGNORE INTO categories (name) VALUES ("tennis");
INSERT OR IGNORE INTO categories (name) VALUES ("football");
