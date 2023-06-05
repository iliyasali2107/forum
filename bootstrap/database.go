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
