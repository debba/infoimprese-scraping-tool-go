package scraper

import (
	"encoding/csv"
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

func SetSearch(query string, where string, config settings.Config, outputFile string) {

	if config.Scraper.Mode == "" {
		config.Scraper.Mode = "search_by_name"
	}

	log.Printf("[SEARCH] Query: %s, Where: %s, OutFile: %s", query, where, outputFile)

	AutoSetting := settings.AutoSetting{Query: query, Mode: config.Scraper.Mode, Where: where, OutputFile: outputFile, ApiKey: config.AntiCaptcha.ApiKey, SiteKey: config.AntiCaptcha.SiteKey}

	var TipoRicerca int
	var IndiceFiglio string

	switch config.Scraper.Mode {
	case "search_by_desc":
		TipoRicerca = 1
		IndiceFiglio = ""
	case "with_dash":
		TipoRicerca = 1
		IndiceFiglio = "0"
	case "with_cert":
		TipoRicerca = 1
		IndiceFiglio = "1"
	case "with_ecom":
		TipoRicerca = 1
		IndiceFiglio = "2"
	case "with_email":
		TipoRicerca = 1
		IndiceFiglio = "3"
	case "with_website":
		TipoRicerca = 1
		IndiceFiglio = "4"
	case "with_export":
		TipoRicerca = 1
		IndiceFiglio = "5"
	default:
		TipoRicerca = 0
		IndiceFiglio = ""
	}

	AutoQueryParams := settings.AutoQueryParams{TipoRicerca: TipoRicerca, IndiceFiglio: IndiceFiglio}

	jar, _ := cookiejar.New(nil)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr, Jar: jar}

	StartSearch(client, AutoSetting)

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("[CREATE FILE] " + err.Error())
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	ScrapePage(client, AutoQueryParams, &AutoSetting, config, w, 1)

	for i := uint64(2); i <= AutoSetting.TotPages; i++ {
		ScrapePage(client, AutoQueryParams, &AutoSetting, config, w, uint(i))
	}

}

func ScrapePage(httpClient *http.Client, autoQueryParams settings.AutoQueryParams, autoSetting *settings.AutoSetting, config settings.Config, w *csv.Writer, page uint) []map[string]string {
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
			"tipoRicerca":  {strconv.Itoa(autoQueryParams.TipoRicerca)},
			"indiceFiglio": {autoQueryParams.IndiceFiglio}}
	} else {
		queryParams = url.Values{
			"pagina":               {"0"},
			"indice":               {strconv.Itoa(int(page))},
			"tipoRicerca":          {strconv.Itoa(autoQueryParams.TipoRicerca)},
			"indiceFiglio":         {autoQueryParams.IndiceFiglio},
			"g-recaptcha-response": {""}}
	}

	/* only if page>2 send captcha response
	   what the fuck did they smoke
	   in Camera di Commercio ?
	*/

	if page > 2 {
		captcha, _ := decrypt.GetCaptcha(autoSetting.ApiKey, autoSetting.SiteKey, endpointUrl)
		queryParams.Set("g-recaptcha-response", captcha)
	}

	result, _ := request.PostRequest(httpClient, endpointUrl, queryParams)

	doc, _ := htmlquery.Parse(strings.NewReader(result))

	if page == 1 {
		tree.CountFromSearch(doc, autoSetting)
		log.Printf("[TOTAL RESULTS] %d", autoSetting.TotResults)
		log.Printf("[TOTAL PAGES] %d", autoSetting.TotPages)
	}

	pages := tree.GetResultPages(doc)

	if len(pages) > 0 {
		for index, crawledPage := range pages {
			log.Printf("[OPEN CONTACT %d] %s", index+1, crawledPage)
			contactUrl := ApiEndpoint + "/ricerca/" + crawledPage
			resp, _ := request.GetRequest(httpClient, contactUrl)
			docResp, _ := htmlquery.Parse(strings.NewReader(resp))
			contacts = append(contacts, tree.GetContactByPage(docResp, config.Scraper.Fields))
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
			_ = w.Write(headers)
		}
		if row != nil {
			_ = w.Write(row)
		}
	}

	return contacts

}

func StartSearch(httpClient *http.Client, autoSetting settings.AutoSetting) {
	endpointUrl := ApiEndpoint + "/ricerca/lista_globale.jsp"
	captcha, _ := decrypt.GetCaptcha(autoSetting.ApiKey, autoSetting.SiteKey, endpointUrl)
	_, _ = request.PostRequest(httpClient, endpointUrl, url.Values{
		"cer":                  {"1"},
		"pagina":               {"0"},
		"flagDove":             {"true"},
		"dove":                 {autoSetting.Where},
		"ricerca":              {autoSetting.Query},
		"g-recaptcha-response": {captcha}})
}
