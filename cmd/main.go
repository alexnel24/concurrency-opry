package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app 	 "OpryScrape/internal/app"
	database "OpryScrape/internal/db"
	handlers "OpryScrape/internal/handlers"
	scrape   "OpryScrape/internal/services/scraping"
	store 	 "OpryScrape/internal/store"
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
