package main

import (
	"fmt"
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

	ctrl := controller.NewController(db)

	route.Setup(env, timeout, db, mux, ctrl)

	fmt.Println("listening")
	http.ListenAndServe(env.ServerAddress, mux)
}
