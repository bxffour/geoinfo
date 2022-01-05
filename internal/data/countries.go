package data

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
