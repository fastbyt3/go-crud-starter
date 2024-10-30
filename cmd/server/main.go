//go:debug x509negativeserial=1
package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastbyt3/diy-mssql-test/internal/server"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()

	utils.Logger.Info("Received termination signal... Attempting to gracefully shutdown")

	serverTerminationCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(serverTerminationCtx)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error gracefully shutting down server: %s", err.Error()))
	}

	done <- true
}

func main() {
	utils.InitLogger()
	defer utils.FlushLogger()

	done := make(chan bool, 1)

	server := server.New()
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil {
		utils.Logger.Fatal(fmt.Sprintf("Error in server: %s", err.Error()))
	}

	<-done
	utils.Logger.Info("Application terminated")
}
