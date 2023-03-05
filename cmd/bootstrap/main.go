package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/bxffour/crest-countries/internal/data"
	flag "github.com/spf13/pflag"

	_ "github.com/lib/pq"
)

//go:embed countries.json
var configFS embed.FS

func main() {
	var dsn string

	flag.StringVarP(&dsn, "database", "d", "", "database connection string")
	flag.Parse()

	// if !validateDsn(dsn) {
	// 	log.Fatal(errors.New("err invalid dsn: valid format -> postgres://username:password@hostname:port/database_name?optional_params"))
	// }

	db, err := sql.Open("postgres", "passfile=/home/sxntana/.pgpass")
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

	file, err := configFS.Open("countries.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	raw, err := io.ReadAll(file)
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

// func validateDsn(dsn string) bool {
// 	r := regexp.MustCompile(`^postgres:\/\/(\w+):([^@]+)@([\w.-]+):(\d+)\/(\w+)(?:\?(.+))?$`)

// 	return r.MatchString(dsn)
// }