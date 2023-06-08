package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var queries = []string{
	`CREATE TABLE IF NOT EXISTS users  (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		token TEXT,
		expiry DATETIME
	);`,
	`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content text NOT NULL,
		created DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`,
	`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE NOT NULL
	);
	INSERT OR IGNORE INTO categories (name) VALUES ("basketball");
	INSERT OR IGNORE INTO categories (name) VALUES ("tennis");
	INSERT OR IGNORE INTO categories (name) VALUES ("football");`,
	`CREATE TABLE IF NOT EXISTS categories_posts (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`,
	`CREATE TABLE IF NOT EXISTS reactions_posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		type INTEGER NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE (post_id, user_id)
	);`,
	`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		parent_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (parent_id) REFERENCES comments(id)
	);
	
	`,
	`CREATE TABLE IF NOT EXISTS reactions_comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		type INTEGER NOT NULL,
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`,
}

func NewSqliteDatabase(env *Env) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filePath := env.DBPath

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatal("couldn't create a db file")
		}
	}

	db, err := sql.Open(env.DBDriver, env.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := CreateTable(db, "./db/migrations"); err != nil {
		log.Fatalln(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("successfully connected to db")

	return db
}

func CreateTable(db *sql.DB, path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range dir {
		info, err := file.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, info.Name()))
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(data)); err != nil {
			log.Println(err)
		}
	}
	return nil
}
