CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE NOT NULL
);

INSERT OR IGNORE INTO categories (name) VALUES ("go");
INSERT OR IGNORE INTO categories (name) VALUES ("c++");
INSERT OR IGNORE INTO categories (name) VALUES ("c#");
INSERT OR IGNORE INTO categories (name) VALUES ("java");
