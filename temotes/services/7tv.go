package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"temotes/temotes"
	"time"
)

type SevenTvFetcher struct{}

type sevenTvEmote struct {
	ID   string            `json:"id"`
	Code string            `json:"name"`
	Urls []sevenTvEmoteUrl `json:"urls"`
}

type sevenTvEmoteUrl = [2]string

func (t SevenTvFetcher) fetchEmotes(url string) []temotes.Emote {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
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

	var sevenTvEmotes []sevenTvEmote
	jsonErr := json.Unmarshal(body, &sevenTvEmotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var emotes []temotes.Emote
	for _, sevenTvEmote := range sevenTvEmotes {
		emotes = append(emotes, t.parseEmote(sevenTvEmote))
	}

	return emotes
}

func (t SevenTvFetcher) FetchGlobalEmotes() []temotes.Emote {
	return t.fetchEmotes("https://api.7tv.app/v2/emotes/global")
}

func (t SevenTvFetcher) FetchChannelEmotes(id temotes.TwitchUserId) []temotes.Emote {
	return t.fetchEmotes(fmt.Sprintf("https://api.7tv.app/v2/users/%d/emotes", id))
}

func (t SevenTvFetcher) parseEmoteUrls(emote sevenTvEmote) []temotes.EmoteUrl {
	var urls []temotes.EmoteUrl

	getEmoteSize := func(scale string) temotes.EmoteSize {
		switch scale {
		case "1":
			return temotes.Size1x
		case "2":
			return temotes.Size2x
		case "3":
			return temotes.Size3x
		case "4":
			return temotes.Size4x
		default:
			return temotes.Size1x
		}
	}

	for _, url := range emote.Urls {
		urls = append(urls, temotes.EmoteUrl{
			Size: getEmoteSize(url[0]),
			Url:  url[1],
		})
	}

	return urls
}

func (t SevenTvFetcher) parseEmote(emote sevenTvEmote) temotes.Emote {
	return temotes.Emote{
		Provider: temotes.Provider7tv,
		Code:     emote.Code,
		Urls:     t.parseEmoteUrls(emote),
	}
}
