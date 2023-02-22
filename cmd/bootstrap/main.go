package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"
	"time"

	"github.com/bxffour/crest-countries/internal/data"

	_ "github.com/lib/pq"
)

func main() {
	var dsn string
	var filepath string

	flag.StringVar(&dsn, "dsn", "", "postgres database connection string")
	flag.StringVar(&filepath, "path", "", "path to the countries json file")
	flag.Parse()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cleanedPath := path.Clean(filepath)

	raw, err := os.ReadFile(cleanedPath)
	if err != nil {
		log.Fatal(err)
	}

	var count []data.Country

	err = json.Unmarshal(raw, &count)
	if err != nil {
		log.Fatal(err)
	}

	for _, country := range count {
		item := &item{
			Country: country,
		}

		err := insert(item, db)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("operation was successful")

}

type item struct {
	ID      int
	Version int32
	Country data.Country
}

func insert(item *item, db *sql.DB) error {
	query := `
		INSERT INTO countries(country)
		VALUES($1)
		RETURNING id, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, item.Country).Scan(&item.ID, &item.Version)
}