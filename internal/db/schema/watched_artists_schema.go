package schema

const WatchedArtistsSchema = `
CREATE TABLE IF NOT EXISTS watched_artists (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    artist_id INTEGER,
    owner_id INTEGER,

    FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE SET NULL,
    UNIQUE(name, owner_id)
);

CREATE INDEX IF NOT EXISTS idx_watched_artists_name ON watched_artists(name);
`
