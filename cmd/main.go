package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app               "github.com/alexnel24/concurrency-opry/internal/app"
	database          "github.com/alexnel24/concurrency-opry/internal/db"
	performancefinder "github.com/alexnel24/concurrency-opry/internal/services/performancefinder"
	scrape            "github.com/alexnel24/concurrency-opry/internal/services/scraping"
	watchlist         "github.com/alexnel24/concurrency-opry/internal/services/watchlist"
	session           "github.com/alexnel24/concurrency-opry/internal/session"
	store             "github.com/alexnel24/concurrency-opry/internal/store"
)



func main() {
	fmt.Println("hi Alex")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, _ := database.InitDB()
	stores := store.InitStores(db)
	scraper := scrape.NewScraper(stores)
	watchlistSvc := watchlist.NewWatchlist(stores)
	performanceFinderSvc := performancefinder.NewPerformanceFinder(stores)
	sessionManager := session.NewSessionManager()

	app := app.NewApp(scraper, watchlistSvc, performanceFinderSvc, stores, sessionManager)

	app.Run(ctx)
}
