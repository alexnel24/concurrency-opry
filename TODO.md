# TODO

## Features
- [ ] Add people-of-interest model + SQL queries to find wanted artists
- [ ] Build unit tests for scraping (mockery + fake Opry HTML)
- [ ] Set up automatic scheduling (cron, like BasketballBubbleScraper)
- [ ] Replace print statements with structured logging
- [ ] Expand error handling throughout
- [ ] Enable SMS-based interaction — send texts when artists of interest are found, allow triggering scrapes via text
- [ ] Switch from SQLite to PostgreSQL (or another full SQL server)

## Bugs

## Open Questions
- Investigate adding Foreign Key artist_id to performances model as primary identifier rather than name. Evaluate ordering risk: can performances be flushed to the DB before their artist exists? If so, enforce artist-before-performance ordering in the flush function.
- Evaluate scrape order of operations: Sync functions after FlushToDb
- Deploy or keep local? (pipeline work, logging, cost)
- Multi-user support? (different wanted-artist lists per user)
- Use `Event.NoOfPerformers` to short-circuit scraping early?
