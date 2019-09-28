package decrypt

import (
	"github.com/debba/anticaptcha"
	"infoimprese-scraping-tool/settings"
	"log"
	"os"
	"time"
)

func GetCaptcha(autoSetting settings.AutoSetting, url string) (string, error) {
	log.Printf("[CHECKING] captcha (url = %s)", url)
	log.Printf("[CHECKING] captcha (siteKey = %s, apiKey = %s)", autoSetting.SiteKey, autoSetting.ApiKey)
	client := anticaptcha.Client{APIKey: autoSetting.ApiKey}
	balance, berr := client.GetBalance()
	if berr != nil {
		log.Printf("[ERROR] error = %s", berr.Error())
		os.Exit(0)
	}
	if balance < autoSetting.MinimumAmount {
		log.Printf("[ERROR] AntiCaptcha balance is too lower (Current: %2f, Minimum requested: %2f)", balance, autoSetting.MinimumAmount)
		os.Exit(0)
	}
	log.Printf("[ANTICAPTCHA BALANCE] Current: %2f", balance)
	captcha, err := client.SendRecaptcha(
		url,
		autoSetting.SiteKey,
		10*time.Minute)

	if err != nil {
		log.Printf("[ERROR] error = %s", err.Error())
		os.Exit(0)
	}

	return captcha, err
}
