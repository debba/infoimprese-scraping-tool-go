package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"infoimprese-scraping-tool/scraper"
	"infoimprese-scraping-tool/settings"
	"io/ioutil"
	"log"
	"os"
)

const configFile string = "conf/config.json"

func main() {

	configJson, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Printf("Error opening config file (%s). Exit.", err.Error())
		os.Exit(0)
	}

	var config settings.Config

	err = json.Unmarshal(configJson, &config)

	if err != nil {
		log.Printf("Error occured while decoding config file (%s). Exit.", err.Error())
		os.Exit(0)
	}

	parser := argparse.NewParser("InfoImprese scraping tool", "Extract data from Infoimprese")
	query := parser.String("q", "query", &argparse.Options{Required: true, Help: "Activites to scrape in Infoimprese"})
	where := parser.String("l", "location", &argparse.Options{Required: true, Help: "Location (it could be a region, city, address etc.)"})
	mode := parser.String("m", "mode", &argparse.Options{Required: false, Help: "Mode", Default: config.Scraper.Mode})
	output := parser.String("o", "output", &argparse.Options{Required: false, Help: "Mode", Default: "export.csv"})

	errn := parser.Parse(os.Args)
	if errn != nil {
		fmt.Print(parser.Usage(errn))
		os.Exit(0)
	}

	config.Scraper.Mode = *mode

	scraper.SetSearch(*query, *where, config, *output)

}
