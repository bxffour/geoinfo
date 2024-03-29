package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
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

func (c *CountryModel) GetAll(ctx context.Context, filters Filters) ([]*Country, Metadata, error) {
	ctx, span := tracer.Start(ctx, "database-get-all")
	defer span.End()

	query := `
			SELECT COUNT(*) OVER(), country 
			FROM countries
			LIMIT $1 OFFSET $2`

	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(dbCtx, query, filters.limit(), filters.offset())
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

func (c *CountryModel) GetByName(otelCtx context.Context, name string, filers Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-name")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		WHERE EXISTS (
			SELECT 1 FROM jsonb_each(country -> 'name') as j(key, value)
			WHERE j.value::text ILIKE '%' || $1 || '%'
			UNION
			SELECT 1 FROM jsonb_each(country -> 'name' -> 'nativeName') as n(key, value)
			WHERE n.value::text ILIKE '%' || $1 || '%'
		)
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, name, filers)
}

func (c CountryModel) GetByCode(otelCtx context.Context, code string) (*Country, error) {
	_, span := tracer.Start(otelCtx, "database-get-by-code")
	defer span.End()

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

func (c CountryModel) GetByCapital(otelCtx context.Context, capital string) (*Country, error) {
	_, span := tracer.Start(otelCtx, "database-get-by-capital")
	defer span.End()

	query := `
		SELECT country FROM countries c
		WHERE (to_tsvector('simple', c.country -> 'capital') @@ plainto_tsquery('simple', $1))
	`

	var country Country

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, capital).Scan(&country)
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

func (c CountryModel) GetByTranslation(otelCtx context.Context, translation string) (*Country, error) {
	_, span := tracer.Start(otelCtx, "database-get-by-translation")
	defer span.End()

	query := `
		SELECT country FROM countries c
		WHERE EXISTS (
			SELECT 1
			FROM jsonb_each(c.country -> 'translations') as n(key, value)
			where value::text ilike '%' || $1 || '%'
		)
	`

	var country Country

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, translation).Scan(&country)
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

func (c CountryModel) GetByCodes(otelCtx context.Context, codes []string) ([]*Country, error) {
	ctx, span := tracer.Start(otelCtx, "database-get-by-codes")
	defer span.End()

	var countries []*Country
	var queried = map[string]uint8{}

	for _, code := range codes {
		code = strings.ToUpper(code)
		country, err := c.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}

		if _, ok := queried[country.Name.Common]; ok {
			continue
		}

		countries = append(countries, country)
		queried[country.Name.Common] = 1
	}

	return countries, nil
}

func (c CountryModel) GetByCurrency(otelCtx context.Context, currency string, filters Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-currency")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		CROSS JOIN LATERAL jsonb_each(c.country->'currencies') AS j(key, value)
		WHERE j.key = $1 
		OR j.value::text ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, currency, filters)
}

func (c CountryModel) GetByLanguage(otelCtx context.Context, language string, filters Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-currency")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		CROSS JOIN LATERAL jsonb_each(c.country->'languages') AS j(key, value)
		WHERE j.key = $1
		OR j.value::text ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, language, filters)
}

func (c CountryModel) GetByRegion(otelCtx context.Context, region string, filters Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-currency")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		WHERE c.country ->> 'region' ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
 	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, region, filters)
}

func (c CountryModel) GetBySubregion(otelCtx context.Context, subregion string, filters Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-currency")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		WHERE c.country ->> 'subregion' ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, subregion, filters)
}

func (c CountryModel) GetByDemonyms(otelCtx context.Context, demonyms string, filters Filters) ([]*Country, Metadata, error) {
	cctx, span := tracer.Start(otelCtx, "database-get-by-currency")
	defer span.End()

	query := `
		SELECT COUNT(*) OVER(), country FROM countries c
		CROSS JOIN LATERAL jsonb_each(c.country->'demonyms') AS j(key, value)
		WHERE j.value::text ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(cctx, 3*time.Second)
	defer cancel()

	return c.multiRows(ctx, query, demonyms, filters)
}

func (c CountryModel) multiRows(ctx context.Context, query, placeholder string, filters Filters) ([]*Country, Metadata, error) {
	_, span := tracer.Start(ctx, "muiltirow-func")
	defer span.End()

	rows, err := c.DB.QueryContext(ctx, query, placeholder, filters.limit(), filters.offset())
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

	if len(countries) == 0 {
		return nil, Metadata{}, ErrRecordNotFound
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return countries, metadata, nil

}
