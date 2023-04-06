package main

import (
	"fmt"
	"forum/internal/repository"
	"forum/pkg/sqlite"
	"log"
	"net/http"

	"forum/internal/delivery"
	"forum/internal/service"
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

	fmt.Println("OOKKK")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}
