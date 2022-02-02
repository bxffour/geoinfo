package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type CountryModel struct {
	DB *sql.DB
}

type Item struct {
	ID      int
	Version int32
	Country Country
}

type Country struct {
	Name         name                   `json:"name"`
	Tld          []string               `json:"tld,omitempty"`
	CCA2         string                 `json:"cca2,omitempty"`
	CCN3         string                 `json:"ccn3,omitempty"`
	CCA3         string                 `json:"cca3,omitempty"`
	CIOC         string                 `json:"cioc,omitempty"`
	Independent  bool                   `json:"independent,omitempty"`
	Status       string                 `json:"status,omitempty"`
	UNMember     bool                   `json:"unMember,omitempty"`
	Currencies   map[string]currency    `json:"currencies,omitempty"`
	IDD          *idd                   `json:"idd,omitempty"`
	Capital      []string               `json:"capital,omitempty"`
	AltSpellings []string               `json:"altSpellings,omitempty"`
	Region       string                 `json:"region,omitempty"`
	SubRegion    string                 `json:"subregion,omitempty"`
	Languages    map[string]string      `json:"languages,omitempty"`
	Translations map[string]translation `json:"translations,omitempty"`
	LatLng       []float32              `json:"latlng,omitempty"`
	LandLocked   bool                   `json:"landlocked,omitempty"`
	Borders      []string               `json:"borders,omitempty"`
	Area         float32                `json:"area,omitempty"`
	Demonyms     map[string]demonyms    `json:"demonyms,omitempty"`
	Flag         string                 `json:"flag,omitempty"`
	Maps         map[string]string      `json:"maps,omitempty"`
	Population   int                    `json:"population,omitempty"`
	Gini         map[string]float32     `json:"gini,omitempty"`
	Fifa         string                 `json:"fifa,omitempty"`
	Car          car                    `json:"car,omitempty"`
	Timezones    []string               `json:"timezones,omitempty"`
	Continents   []string               `json:"continents,omitempty"`
	Flags        graphicsFormat         `json:"flags,omitempty"`
	CoatOfArms   *graphicsFormat        `json:"coatOfArms,omitempty"`
	StartOfWeek  string                 `json:"startOfWeek,omitempty"`
	CapitalInfo  *capitalInfo           `json:"capitalInfo,omitempty"`
	PostalCode   *postalCode            `json:"postalCode,omitempty"`
}

type name struct {
	Official   string            `json:"official,omitempty"`
	Common     string            `json:"common,omitempty"`
	NativeName map[string]native `json:"nativeName,omitempty"`
}

type native struct {
	Official string `json:"official,omitempty"`
	Common   string `json:"common,omitempty"`
}

type currency struct {
	Name   string `json:"name,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}

type idd struct {
	Root     string   `json:"root,omitempty"`
	Suffixes []string `json:"suffixes,omitempty"`
}

type translation struct {
	Official string `json:"official,omitempty"`
	Common   string `json:"common,omitempty"`
}

type demonyms struct {
	Female string `json:"f,omitempty"`
	Male   string `json:"m,omitempty"`
}

type car struct {
	Signs []string `json:"signs,omitempty"`
	Side  string   `json:"side,omitempty"`
}

type graphicsFormat struct {
	PNG string `json:"png,omitempty"`
	SVG string `json:"svg,omitempty"`
}

type capitalInfo struct {
	LatLng []float32 `json:"latlng,omitempty"`
}

type postalCode struct {
	Format string `json:"format,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

func (c Country) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Country) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c *CountryModel) Insert(item *Item) error {
	query := `
			INSERT INTO countries(country) 
			VALUES($1)
			RETURNING id, version
			`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, item.Country).Scan(&item.ID, &item.Version)
}

func (c *CountryModel) GetAll(filters Filters) ([]*Country, Metadata, error) {
	query := `
			SELECT COUNT(*) OVER(), country 
			FROM countries
			LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	countries := []*Country{}

	for rows.Next() {
		var country Country

		err := rows.Scan(
			&totalRecords,
			&country,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		countries = append(countries, &country)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return countries, metadata, nil
}

func (c *CountryModel) GetByName(name string) ([]*Country, error) {
	query := `
		SELECT country FROM countries
		WHERE (to_tsvector('simple', country->'name') @@ plainto_tsquery('simple', $1))
		ORDER BY country->'name'->'official'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	countries := []*Country{}

	for rows.Next() {
		var country Country

		err := rows.Scan(&country)
		if err != nil {
			return nil, err
		}

		countries = append(countries, &country)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(countries) == 0 {
		return nil, ErrRecordNotFound
	}

	return countries, nil

}

func (c CountryModel) GetByCode(code string) (*Country, error) {
	query := `
		SELECT country FROM countries
		WHERE country->'cca2' ? $1
		or country->'ccn3' ? $1
		or country->'cca3' ? $1
		or country->'cioc' ? $1
	`

	var country Country

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, code).Scan(&country)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &country, nil
}
