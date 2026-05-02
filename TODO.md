# TODO

## Features
- [ ] Add tests for ArtistStore and PerformanceStore
- [ ] Add people-of-interest model + SQL queries to find wanted artists
- [ ] Capture actual event date (playwright vs colly) and mark past vs upcoming
- [ ] Build unit tests for scraping (mockery + fake Opry HTML)
- [ ] Set up automatic scheduling (cron, like BasketballBubbleScraper)
- [ ] Replace print statements with structured logging
- [ ] Expand error handling throughout
- [ ] Enable SMS-based interaction — send texts when artists of interest are found, allow triggering scrapes via text
- [ ] Switch from SQLite to PostgreSQL (or another full SQL server)

## Bugs
(currently no known bugs)

## Open Questions
- Deploy or keep local? (pipeline work, logging, cost)
- Multi-user support? (different wanted-artist lists per user)
- Use `Event.NoOfPerformers` to short-circuit scraping early?
