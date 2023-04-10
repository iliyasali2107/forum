CREATE TABLE votes_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES comments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
