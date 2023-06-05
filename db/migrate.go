package db

import (
	"database/sql"
	"fmt"
	"forum/pkg/sqlite"
	"log"
	"os"
)

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

func SetupDB() {
	if len(os.Args) <= 1 || len(os.Args) >= 3 {
		log.Fatal("Usage: go run change_db.go <argument>")
	}
	flag := os.Args[1]
	db, err := sqlite.Connect("./db/forum.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	switch flag {
	case "up":
		if err := CreateTable(db, "./db/migrations"); err != nil {
			log.Fatalln(err)
		}
		log.Println("Successful")
	case "down":
		if err := sqlite.DropAllDB(db); err != nil {
			log.Fatalln(err)
		}
		log.Println("Successful")
	default:
		log.Fatalf("%s: unknown flag. Use: 'up' or 'down'", flag)
	}
}
