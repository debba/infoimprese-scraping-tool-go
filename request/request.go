package request

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func PostRequest(httpClient *http.Client, endpointUrl string, formData url.Values) (string, error) {

	resp, err := httpClient.PostForm(endpointUrl, formData)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return string(body), nil

}

func GetRequest(httpClient *http.Client, endpointUrl string) (string, error) {

	resp, err := httpClient.Get(endpointUrl)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return string(body), nil

}
