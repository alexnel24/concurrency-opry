package app

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/alexnel24/concurrency-opry/internal/handlers"
	"github.com/alexnel24/concurrency-opry/internal/server"
	"github.com/alexnel24/concurrency-opry/internal/services/scraping"
	"github.com/alexnel24/concurrency-opry/internal/session"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

const defaultDbBatchSize = 100
const defaultFlushSeconds = 120

type App struct {
	handler        *handlers.Handler
	stores         *store.Stores
	sessionManager *session.SessionManager
}

func NewApp(scraper *scraping.Scraper, stores *store.Stores, sessionManager *session.SessionManager) *App {
	return &App{
		handler:        handlers.New(scraper, stores, sessionManager),
		stores:         stores,
		sessionManager: sessionManager,
	}
}

func (a *App) Run(ctx context.Context) {
	batchSize := envInt("DB_BATCH_SIZE", defaultDbBatchSize)
	flushDbEverySeconds := envInt("FLUSH_SECONDS", defaultFlushSeconds)

	a.sessionManager.StartBackgroundSessionCleanup(ctx)
	a.stores.StartBackgroundDBWorker(ctx, batchSize, flushDbEverySeconds)

	mux := http.NewServeMux()
	httpHandler := server.Routes(mux, a.handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := server.New(":"+port, httpHandler)
	server.Run(ctx)
	a.stores.WaitForBackgroundDBWorkerToFlush()
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
