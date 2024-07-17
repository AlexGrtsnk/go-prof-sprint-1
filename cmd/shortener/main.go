package main

import (
	"context"
	"fmt"
	fun "go-prof-sprint-1/internal/functions"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// глобальные переменные флагов
var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	fmt.Printf("version=%s, date=%s, commit=%ss\n", BuildVersion, BuildDate, BuildCommit)
	srv := fun.Run()
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// ошибки старта или остановки Listener
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

}
