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
	srv, enableHTTPS := fun.Run()
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if !enableHTTPS {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		fmt.Println(srv.Addr)
		/*
			if err := srv.ListenAndServeTLS("certificate", "key"); err != http.ErrServerClosed {
				// ошибки старта или остановки Listener
				log.Fatalf("HTTP server ListenAndServe: %v", err)

			}*/
	}
	<-idleConnsClosed

}
