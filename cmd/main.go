package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dilap54/voronov_idor/internal/app"
	"github.com/dilap54/voronov_idor/internal/config"
	"github.com/dilap54/voronov_idor/internal/repository"
	"github.com/dilap54/voronov_idor/internal/server"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := repository.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := repository.CreateInitialData(db, config.GetCode()); err != nil {
		log.Fatal(err)
	}

	appSvc := app.New(db)
	addr := config.GetHost()
	srv := server.New(addr, appSvc)

	log.Printf("Запуск сервера на %s", addr)
	go srv.Run()

	<-done
	log.Printf("Завершение работы сервера")
	srv.Shutdown(context.Background())
}
