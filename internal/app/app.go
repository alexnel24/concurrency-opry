package app

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"OpryScrape/internal/handlers"
	"OpryScrape/internal/server"
	"OpryScrape/internal/store"
)

const defaultDbBatchSize = 100
const defaultFlushSeconds = 120

type App struct {
	handler *handlers.Handler
	stores  *store.Stores
}

func NewApp(handler *handlers.Handler, stores *store.Stores) *App {
	return &App{
		handler: handler,
		stores:  stores,
	}
}

func (a *App) Run(ctx context.Context) {
	batchSize := envInt("DB_BATCH_SIZE", defaultDbBatchSize)
	flushDbEverySeconds := envInt("FLUSH_SECONDS", defaultFlushSeconds)

	a.stores.StartBackgroundDBWorker(ctx, batchSize, flushDbEverySeconds)

	mux := http.NewServeMux()
	httpHandler := server.Routes(mux, a.handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := server.New(":"+port, httpHandler)
	server.Run(ctx)
	//ToDo: Error handling
}

func envInt(key string, defaultInt int) int {
	if envString := os.Getenv(key); envString != "" {
		num, err := strconv.Atoi(envString)
		if err == nil && num > 0 {
			return num
		}
	}

	return defaultInt
}
