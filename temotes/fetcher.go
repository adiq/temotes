package temotes

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Fetcher struct{}

func (f Fetcher) GetRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func (f Fetcher) FetchData(url string) ([]byte, error) {
	return f.FetchDataRequest(f.GetRequest(url))
}

func (f Fetcher) FetchDataRequest(req *http.Request) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.StatusCode != 200 {
		return nil, errors.New("request returned non-successful response response")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}
