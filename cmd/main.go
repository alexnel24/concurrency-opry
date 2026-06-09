package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app      "github.com/alexnel24/concurrency-opry/internal/app"
	database "github.com/alexnel24/concurrency-opry/internal/db"
	scrape   "github.com/alexnel24/concurrency-opry/internal/services/scraping"
	session  "github.com/alexnel24/concurrency-opry/internal/session"
	store    "github.com/alexnel24/concurrency-opry/internal/store"
)



func main() {
	fmt.Println("hi Alex")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, _ := database.InitDB()
	stores := store.InitStores(db)
	scraper := scrape.NewScraper(stores)
	sessionManager := session.NewSessionManager()

	app := app.NewApp(scraper, stores, sessionManager)

	app.Run(ctx)
}
