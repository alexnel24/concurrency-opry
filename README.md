OpryScrape

Trying to see Brad Paisley or Reba (for Sam)

What is the Opry: 
A radio show that happens in Nashville ~250 times per year. Every show has 5-8 artists/groups that perform 3-5 songs. Shows are announced well before they announce the performers. It is common for the "big" name acts to not be linked to a show until a week or two before the show happens.

Purpose(s) of this scraper:
1. Mostly fun, building a pet project to keep learning about Effective Go and Software Engineering
2. Keep track of who is coming to the Opry so that I can watch for people of interest

Three scripts:
build - ./scripts/build
server - ./scripts/server
test - ./scripts/test

Three Endpoints:
/health - basic healthcheck
/scrape - scrape the opry websites for new annouced performances, artists, and events. Updates DB at end
/update-db - force an update to db

Database: SQLite
Only allows one connection at a time. There is a background worker listening for new events, artists, and performances. The items are inserted to the db in batches. Both batch size and partial batches time limits are controlled via env vars 

Needed setup:
An empty data/opry.db file is required

Future Goals:
1. Add the people of interest model and incorporate sql select statements to find the wanted artists
2. Capture the date of the event (playwright vs colly) and mark past vs upcoming events
3. Send a text to myself when artists of interest are found
4. Testing - build out unit tests for scraping using mockery and fake values from Opry web page
5. Setup to run automatically (cron as I've done in other personal project: BasketballBubbleScraper)
6. Update any print statements to logs (i.e. implement and improve logging)
7. Expand Error Handling throughout application

Things still in question:
1. Figure out if deploying is goal (Pipeline work, logging, cost)
2. Would anybody else use this? Do I need users that have different wanted artists?
3. Should I utilize event.NoOfPerformers? Has potential to stop scrape early

Things I would do differently (hindsight bias):
1. For event scraping, I would attempt using playwright instead of colly
2. Switch from sqlite and split clunky stores.StartBackgroundWorker into 3

Things I'm proud of:
1. Concurrency
2. Project Architecture
3. SQLite db, tried to minimize duplication, utilized indexing

Notes for Rockbot:
Commenting: If you would like a sample of traditional commenting, I have provided it on internal/servcies/scraping/event_scrape.go