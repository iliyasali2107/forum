CREATE TABLE votes_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);