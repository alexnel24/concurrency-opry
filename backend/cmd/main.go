package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app 	 "github.com/alexnel24/concurrency-opry/internal/app"
	database "github.com/alexnel24/concurrency-opry/internal/db"
	handlers "github.com/alexnel24/concurrency-opry/internal/handlers"
	scrape   "github.com/alexnel24/concurrency-opry/internal/services/scraping"
	store 	 "github.com/alexnel24/concurrency-opry/internal/store"
)



func main() {
	fmt.Println("hi Alex")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	db, _ := database.InitDB()
	stores := store.InitStores(db)
	scraper := scrape.NewScraper(stores)

	handler := handlers.New(scraper, stores)
	
	app := app.NewApp(handler, stores)

	app.Run(ctx)
}
