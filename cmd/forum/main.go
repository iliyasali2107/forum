package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/sqlite"
)

func main() {
	mux := http.NewServeMux()
	db, err := sqlite.Connect("./db/forum.db")
	if err != nil {
		log.Fatalln(err)
	}
	rp := repository.NewRepository(db)
	svc := service.NewService(rp)
	handler := delivery.NewHandler(svc)
	handler.InitRoutes(mux)

	fmt.Println("http://localhost:8080/")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}
