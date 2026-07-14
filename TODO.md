# TODO

## Features
- [ ] Build unit tests for scraping (mockery + fake Opry HTML)
- [ ] Set up automatic scheduling (cron, like BasketballBubbleScraper)
- [ ] Replace print statements with structured logging
- [ ] Expand error handling throughout
- [ ] Enable SMS-based interaction — send texts when artists of interest are found, allow triggering scrapes via text
- [ ] Update ArtistId when watched artist is found (backfill `watched_artists.artist_id` once a previously-unmatched entry resolves to a scraped artist)

## In Progress

## Bugs

## Open Questions
- Investigate adding Foreign Key artist_id to performances model as primary identifier rather than name. Evaluate ordering risk: can performances be flushed to the DB before their artist exists? If so, enforce artist-before-performance ordering in the flush function.
- Evaluate scrape order of operations: Sync functions after FlushToDb
- Deploy or keep local? (pipeline work, logging, cost)
- Use `Event.NoOfPerformers` to short-circuit scraping early?
