package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func NewSqliteDatabase(env *Env) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open(env.DBDriver, env.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return db

}
