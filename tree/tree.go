package tree

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"infoimprese-scraping-tool/settings"
	"math"
	"strconv"
	"strings"
)

func CountFromSearch(tree *html.Node, autoSetting *settings.AutoSetting) {

	totResults := strings.TrimSpace(htmlquery.InnerText(
		htmlquery.FindOne(tree, "//html/body/center/table[2]/tbody/tr[2]/td[1]/table[1]/tbody/tr/td/table[2]/tbody/tr/td[1]/font/text()[2]"))[7:])

	autoSetting.TotResults, _ = strconv.ParseUint(totResults, 10, 32)
	autoSetting.TotPages = uint64(math.Ceil(float64(autoSetting.TotResults) / float64(10)))

}

func GetResultPages(tree *html.Node) []string {

	var i uint64
	var pages []string
	for i = 3; i < 13; i++ {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		onclickFunc := htmlquery.FindOne(tree, "//html/body/center/table[2]/tbody/tr[2]/td[1]/table[1]/tbody/tr/td/table["+strconv.FormatUint(i, 10)+"]/tbody/tr[2]/td/table/tbody/tr/td[2]/a[1]/@onclick")
		if onclickFunc != nil {
			onclickFuncStr := strings.TrimSpace(htmlquery.InnerText(onclickFunc))
			pages = append(pages, onclickFuncStr[14:len(onclickFuncStr)-33])
		}
	}

	return pages
}

func GetValueByField(field string, tree *html.Node) string {
	xpath := "//b[contains(text(), '" + field + "')]/../following-sibling::td"

	if field == "Indirizzo web" || field == "Posta elettronica" || field == "Commercio elettronico" {
		xpath += "/a"
	}

	value := ""

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	fieldValue := htmlquery.FindOne(tree, xpath)
	if fieldValue != nil {
		value = htmlquery.InnerText(fieldValue)
	}
	return value
}

func CreateContact(fields []string, tree *html.Node) map[string]string {
	scrapedContact := make(map[string]string)
	for _, field := range fields {
		scrapedContact[field] = GetValueByField(field, tree)
	}
	return scrapedContact
}

func GetContactByPage(tree *html.Node, fields []string) map[string]string {

	if len(fields) == 0 {
		fields = []string{
			"Denominazione",
			"Sede legale",
			"Attività",
			"Sede operativa",
			"Indirizzo web",
			"Posta elettronica",
			"Commercio elettronico",
			"Chi siamo",
			"Cosa facciamo",
			"Classe di fatturato",
			"Canali di vendita",
			"Marchi",
			"Principali paesi di export",
			"Certificazioni"}
	}

	return CreateContact(fields, tree)

}
