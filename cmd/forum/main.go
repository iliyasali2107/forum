package main

import (
	"net/http"
	"time"

	"forum/api/controller"
	"forum/api/route"
	"forum/bootstrap"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.DB

	timeout := time.Duration(env.ContextTimeout) * time.Second

	mux := http.NewServeMux()

	ctrl := controller.NewController()

	route.Setup(env, timeout, db, mux, ctrl)

	http.ListenAndServe(env.ServerAddress, mux)
}
