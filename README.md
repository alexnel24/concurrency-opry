ConcurrencyOpry

Purpose(s) of this application:
1. Keep learning about Effective Go (specifically concurrency) and Software Engineering design principles
2. Keep track of who is coming to the Opry so that I can watch for people of interest

What is the Opry: 
A radio show that happens in Nashville ~250 times per year. Every show has 5-8 artists/groups that perform 3-5 songs. Shows are announced well before they announce the performers. It is common for the "big" name acts to not be linked to a show until a week or two before the show happens.


Three scripts:
build - ./scripts/build
server - ./scripts/server
test - ./scripts/test

Six Endpoints:
/health - basic healthcheck
/scrape - scrape the opry websites for new annouced performances, artists, and events. Updates DB at end
/update-db - force an update to db
/mark-past-events - mark events whose time has passed as no longer upcoming
/artist-performances - return performances for one or more artists (?artist=, ?filter=all|upcoming|past)
/sessions - POST creates a session; DELETE (X-Session-ID header) destroys it

Database: SQLite
Only allows one connection at a time. There is a background worker listening for new events, artists, and performances. The items are inserted to the db in batches. Both batch size and partial batches time limits are controlled via env vars 

Needed setup:
An empty data/opry.db file is required

Future Goals:
See TODO.md

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
