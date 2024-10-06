package main

import (
	"api_wappi/interal/handler"
	"api_wappi/interal/serivce"
	"context"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	services := serivce.New(container)
	handlers := handler.New(services)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handlers.Init(),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		<-c

		shutdownCtx, _ := context.WithTimeout(mainCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()

			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatalf("timeout, force exit")
			}
		}()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("server shutdown error: %v", err)
		}

		cancel()
	}()

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-mainCtx.Done()
}
