package bootstrap

import (
	"database/sql"
	"fmt"
)

type Application struct {
	Env *Env
	DB  *sql.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	fmt.Println(app.Env)
	app.DB = NewSqliteDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	app.DB.Close()
}
