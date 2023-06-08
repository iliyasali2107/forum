package main

import (
	"forum/api/controller"
	"forum/api/route"
	"forum/bootstrap"
	"net/http"
	"time"
)

func main() {
	app := bootstrap.App()
	defer app.CloseDBConnection() ////////

	env := app.Env

	db := app.DB

	timeout := time.Duration(env.ContextTimeout) * time.Second

	mux := http.NewServeMux()

	ctrl := controller.NewController(db)

	route.Setup(env, timeout, db, mux, ctrl)
	
	http.ListenAndServe(env.ServerAddress, mux)
}
