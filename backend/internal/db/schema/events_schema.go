package schema

const EventsSchema = `
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    link TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    time DATETIME NOT NULL,
    no_of_performers INTEGER NOT NULL
);
`