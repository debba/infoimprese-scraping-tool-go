package scraper

import (
	"encoding/csv"
	"fmt"
	"github.com/antchfx/htmlquery"
	"infoimprese-scraping-tool/decrypt"
	"infoimprese-scraping-tool/request"
	"infoimprese-scraping-tool/settings"
	"infoimprese-scraping-tool/tree"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const ApiEndpoint string = "https://www.infoimprese.it/impr"

type Scraper struct {
	settings.AutoQueryParams
	settings.AutoSetting
	settings.Config
	Client *http.Client
	Writer *csv.Writer
	Count  uint64
}

func (s *Scraper) SetSearch(query string, where string, config settings.Config, outputFile string) {

	if config.Scraper.Mode == "" {
		config.Scraper.Mode = "search_by_name"
	}

	log.Printf("[SEARCH] \nQuery: %s, \nWhere: %s, \nMode: %s, \nOutFile: %s", query, where, config.Scraper.Mode, outputFile)

	s.AutoSetting = settings.AutoSetting{Query: query, Mode: config.Scraper.Mode, Where: where, OutputFile: outputFile, ApiKey: config.AntiCaptcha.ApiKey, SiteKey: config.AntiCaptcha.SiteKey}

	jar, _ := cookiejar.New(nil)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	s.Client = &http.Client{Transport: tr, Jar: jar}
	s.Config = config

	s.StartSearch()

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("[CREATE FILE] " + err.Error())
	}
	defer file.Close()

	s.Writer = csv.NewWriter(file)
	defer (s.Writer).Flush()

	s.ScrapePage(1)

	for i := uint64(2); i <= s.AutoSetting.TotPages; i++ {
		s.ScrapePage(uint(i))
	}

}

func (s *Scraper) ScrapePage(page uint) []map[string]string {
	log.Printf("[OPEN PAGE] Page: %s", strconv.Itoa(int(page)))
	var contacts []map[string]string
	endpointUrl := ApiEndpoint + "/ricerca/risultati_globale.jsp"
	queryParams := url.Values{}
	if page != 1 {
		endpointUrl = ApiEndpoint + "/ricerca/pagCaptcha.jsp"
	}

	if page == 1 {
		queryParams = url.Values{
			"cer":          {"1"},
			"statistiche":  {"S"},
			"tipoRicerca":  {s.AutoQueryParams.TipoRicerca},
			"indiceFiglio": {s.AutoQueryParams.IndiceFiglio}}
	} else {
		queryParams = url.Values{
			"pagina":               {"0"},
			"indice":               {strconv.Itoa(int(page))},
			"tipoRicerca":          {s.AutoQueryParams.TipoRicerca},
			"indiceFiglio":         {s.AutoQueryParams.IndiceFiglio},
			"g-recaptcha-response": {""}}
	}

	/* only if page>2 send captcha response
	   what the fuck did they smoke
	   in Camera di Commercio ?
	*/

	if page > 2 {
		captcha, _ := decrypt.GetCaptcha(s.AutoSetting.ApiKey, s.AutoSetting.SiteKey, endpointUrl)
		queryParams.Set("g-recaptcha-response", captcha)
	}

	result, _ := request.PostRequest(s.Client, endpointUrl, queryParams)

	doc, _ := htmlquery.Parse(strings.NewReader(result))

	if page == 1 {
		tree.CountFromSearch(doc, &s.AutoSetting)
		log.Printf("[TOTAL RESULTS] %d", s.AutoSetting.TotResults)
		log.Printf("[TOTAL PAGES] %d", s.AutoSetting.TotPages)
	}

	pages := tree.GetResultPages(doc)

	if len(pages) > 0 {
		for index, crawledPage := range pages {
			log.Printf("[OPEN CONTACT %d] %s", index+1, crawledPage)
			contactUrl := ApiEndpoint + "/ricerca/" + crawledPage
			resp, _ := request.GetRequest(s.Client, contactUrl)
			docResp, _ := htmlquery.Parse(strings.NewReader(resp))
			contacts = append(contacts, tree.GetContactByPage(docResp, s.Config.Scraper.Fields))
		}
	}

	var headers []string
	var row []string

	for index, contact := range contacts {
		row = nil
		for headerName, value := range contact {
			if index == 0 && page == 1 {
				headers = append(headers, headerName)
			}
			row = append(row, value)
		}
		if index == 0 && page == 1 {
			log.Printf("[CREATE HEADER] First time")
			_ = s.Writer.Write(headers)
		}
		if row != nil {
			_ = s.Writer.Write(row)
		}
	}

	return contacts

}

func (s *Scraper) StartSearch() {
	endpointUrl := ApiEndpoint + "/ricerca/lista_globale.jsp"
	captcha, _ := decrypt.GetCaptcha(s.AutoSetting.ApiKey, s.AutoSetting.SiteKey, endpointUrl)
	result, _ := request.PostRequest(s.Client, endpointUrl, url.Values{
		"cer":                  {"1"},
		"pagina":               {"0"},
		"flagDove":             {"true"},
		"dove":                 {s.AutoSetting.Where},
		"ricerca":              {s.AutoSetting.Query},
		"g-recaptcha-response": {captcha}})

	doc, _ := htmlquery.Parse(strings.NewReader(result))
	searchFilterMap := tree.GenerateSearchFilterMap(doc)
	mySearchFilters := searchFilterMap[s.Config.Scraper.Mode]
	s.Count = mySearchFilters.Count
	s.AutoQueryParams.TipoRicerca = mySearchFilters.AutoQueryParams.TipoRicerca
	s.AutoQueryParams.IndiceFiglio = mySearchFilters.AutoQueryParams.IndiceFiglio

	if s.Count == 0 {
		fmt.Println("[ERROR] For the selected mode, no contacts were found.")
		os.Exit(0)
	}

}
