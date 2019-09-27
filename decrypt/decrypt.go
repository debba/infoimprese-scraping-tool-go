package decrypt

import (
	"github.com/debba/anticaptcha"
	"log"
	"os"
	"time"
)

func GetCaptcha(apiKey string, siteKey string, url string) (string, error) {
	log.Printf("[CHECKING] captcha (url = %s)", url)
	log.Printf("[CHECKING] captcha (siteKey = %s, apiKey = %s)", siteKey, apiKey)
	client := anticaptcha.Client{APIKey: apiKey}
	balance, berr := client.GetBalance()
	if berr != nil {
		log.Printf("[ERROR] error = %s", berr.Error())
		os.Exit(0)
	}
	log.Printf("[ANTICAPTCHA BALANCE] Current: %2f", balance)
	captcha, err := client.SendRecaptcha(
		url,
		siteKey,
		10*time.Minute)

	if err != nil {
		log.Printf("[ERROR] error = %s", err.Error())
		os.Exit(0)
	}

	return captcha, err
}
