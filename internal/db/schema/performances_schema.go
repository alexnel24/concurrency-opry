package schema

const PerformancesSchema = `
CREATE TABLE IF NOT EXISTS performances (
    id INTEGER PRIMARY KEY,
    event_link TEXT NOT NULL,
    artist_name TEXT NOT NULL,
    combo_string TEXT NOT NULL UNIQUE,

    FOREIGN KEY(event_link) REFERENCES events(link) ON DELETE CASCADE,
    FOREIGN KEY(artist_name) REFERENCES artists(name) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_performers_artist_name ON performances(artist_name);
`