package decrypt

import (
	"github.com/nuveo/anticaptcha"
	"log"
	"time"
)

func GetCaptcha(apiKey string, siteKey string, url string) (string, error) {
	log.Printf("[CHECKING] captcha (siteKey = %s, apiKey = %s)", siteKey, apiKey)
	client := anticaptcha.Client{APIKey: apiKey}
	captcha, err := client.SendRecaptcha(
		url,
		siteKey,
		10*time.Minute)

	if err != nil {
		log.Printf("[ERROR] error = %s", err.Error())
	}

	return captcha, err
}
