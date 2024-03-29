package settings

type AntiCaptchaConfig struct {
	ApiKey        string `json:"api_key"`
	SiteKey       string `json:"site_key"`
	MinimumAmount string `json:"minimum_amount"`
}

type ScraperConfig struct {
	Mode   string   `json:"mode"`
	Fields []string `json:"fields"`
}

type Config struct {
	AntiCaptcha AntiCaptchaConfig `json:"anticaptcha"`
	Scraper     ScraperConfig     `json:"scraper"`
}

type AutoSetting struct {
	Query         string
	Where         string
	Mode          string
	OutputFile    string
	ApiKey        string
	SiteKey       string
	MinimumAmount float64
	TotResults    uint64
	TotPages      uint64
}

type AutoQueryParams struct {
	TipoRicerca  string
	IndiceFiglio string
}

type SearchFilterMap struct {
	AutoQueryParams
	Count uint64
}
