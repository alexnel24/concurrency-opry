package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
    "os"

	"github.com/alexnel24/concurrency-opry/internal/db/schema"
)

const defalutDbPath = "data/opry.db"

func InitDB() (*sql.DB, error) {
    db, err := openDB()
    if err != nil {
        return nil, err
    }
    
    applySchemas(db)
    return db, nil
}

func openDB() (*sql.DB, error) {
    dbPath := os.Getenv("DATABASE_PATH")
    if dbPath == "" {dbPath = defalutDbPath}

    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)

    db.Exec(`PRAGMA journal_mode=WAL;`)
    db.Exec(`PRAGMA synchronous=NORMAL;`)
    db.Exec(`PRAGMA busy_timeout=3000;`)

    if err = db.Ping(); err != nil {return nil, err}
    
    return db, nil
}

func applySchemas(db *sql.DB) {
    schemas := map[string]string{
        "Events":      schema.EventsSchema,
        "Artists":     schema.ArtistsSchema,
        "Performances": schema.PerformancesSchema,
    }

    for name, sqlText := range schemas {
        if _, err := db.Exec(sqlText); err != nil {
            fmt.Println("Error applying schema: ", name)
            fmt.Println("Error: ", err.Error())
            continue
        }
    }

    fmt.Println("Done applying all schemas")
}

