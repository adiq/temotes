package temotes

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Fetcher struct{}

func (f Fetcher) GetRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	return req
}

func (f Fetcher) GetGqlRequest(url string, query string) *http.Request {
	cleanQuery := strings.ReplaceAll(query, "\n", " ")
	cleanQuery = strings.ReplaceAll(cleanQuery, "\r", " ")
	cleanQuery = strings.ReplaceAll(cleanQuery, "\t", " ")

	spaceRegexp, err := regexp.Compile("\\s+")
	if err != nil {
		panic(err)
	}

	cleanQuery = strings.TrimSpace(spaceRegexp.ReplaceAllString(cleanQuery, " "))

	payload := map[string]interface{}{
		"query": cleanQuery,
	}

	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(parsedPayload)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (f Fetcher) FetchData(url string) ([]byte, error) {
	return f.FetchDataRequest(f.GetRequest(url))
}

func (f Fetcher) FetchDataRequest(req *http.Request) ([]byte, error) {
	timeout, timeoutErr := strconv.ParseInt(os.Getenv("FETCHER_TIMEOUT"), 10, 64)
	if timeoutErr != nil {
		log.Print("FETCHER_TIMEOUT not specified or defined incorrectly. Defaulting to 3 seconds.")
		timeout = 3 // sane default
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
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
