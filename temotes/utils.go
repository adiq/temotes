package temotes

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"temotes/temotes/services"
	"time"
)

func FetchData(url string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return FetchDataRequest(req)
}

func FetchDataRequest(req *http.Request) []byte {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetTwitchUserId(input string) TwitchUserId {
	id, err := strconv.ParseInt(input, 10, 64)

	if id == 0 || err != nil {
		return services.TwitchFetcher{}.FetchUserId(input)
	}

	return TwitchUserId(id)
}
