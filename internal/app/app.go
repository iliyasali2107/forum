package app

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"forum/internal/config"
// 	"forum/internal/delivery"
// 	"forum/internal/repository"
// 	"forum/internal/server"
// 	"forum/internal/service"
// 	"forum/pkg"
// )

// func Run(cfgFilePath string) error {
// 	logger := pkg.NewLogger("Forum App")
// 	cfg, err := config.NewConfig(cfgFilePath)
// 	if err != nil {
// 		return err
// 	}

// 	db, err := repository.InitDB(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	if err := repository.CreateTables(db); err != nil {
// 		return err
// 	}

// 	repository := repository.NewRepository(db, cfg)
// 	service := service.NewService(repository)
// 	handler := delivery.NewHandler(service)

// 	server := server.NewServer(cfg, handler)
// 	log.Printf("Starting server at port %v\nhttp://localhost:%v\n", cfg.Http.Addr, cfg.Http.Addr)

// 	go func() {
// 		if err := server.Start(); err != nil {
// 			log.Fatalf("error while running: %s", err.Error())
// 		}
// 	}()
// 	logger.Info("Forum started")

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
// 	<-quit

// 	logger.Info("Forum shutting down")

// 	if err := server.ShutDown(context.Background()); err != nil {
// 		return fmt.Errorf("error while shutting down server: %s", err)
// 	}

// 	if err := db.Close(); err != nil {
// 		return fmt.Errorf("error while closing database connection: %s", err)
// 	}
// 	return nil
// }
