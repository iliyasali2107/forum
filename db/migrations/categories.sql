CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE NOT NULL
);

INSERT INTO categories (name) VALUES ("basketball");
INSERT INTO categories (name) VALUES ("ps");
INSERT INTO categories (name) VALUES ("tennis");
INSERT INTO categories (name) VALUES ("football");
