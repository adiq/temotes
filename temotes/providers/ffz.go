package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"temotes/temotes"
)

type FfzFetcher struct{}

type ffzEmote struct {
	ID   int               `json:"id"`
	Code string            `json:"name"`
	Urls map[string]string `json:"urls"`
}

type ffzSetsResponse struct {
	Emotes []ffzEmote `json:"emoticons"`
}

type ffzResponse struct {
	Sets map[string]ffzSetsResponse `json:"sets"`
}

func (t FfzFetcher) fetchEmotes(url string) []temotes.Emote {
	response := temotes.FetchData(url)
	var ffzEmotesResponse ffzResponse
	jsonErr := json.Unmarshal(response, &ffzEmotesResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var emotes []temotes.Emote
	for _, emoteSet := range ffzEmotesResponse.Sets {
		for _, ffzEmote := range emoteSet.Emotes {
			emotes = append(emotes, t.parseEmote(ffzEmote))
		}
	}

	return emotes
}

func (t FfzFetcher) FetchGlobalEmotes() []temotes.Emote {
	return t.fetchEmotes("https://api.frankerfacez.com/v1/set/global")
}

func (t FfzFetcher) FetchChannelEmotes(id temotes.TwitchUserId) []temotes.Emote {
	return t.fetchEmotes(fmt.Sprintf("https://api.frankerfacez.com/v1/room/id/%d", id))
}

func (t FfzFetcher) parseEmoteUrls(emote ffzEmote) []temotes.EmoteUrl {
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

	var keys []string
	for k := range emote.Urls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, scale := range keys {
		urls = append(urls, temotes.EmoteUrl{
			Size: getEmoteSize(scale),
			Url:  fmt.Sprintf("https:%s", emote.Urls[scale]),
		})
	}

	return urls
}

func (t FfzFetcher) parseEmote(emote ffzEmote) temotes.Emote {
	return temotes.Emote{
		Provider: temotes.ProviderFfz,
		Code:     emote.Code,
		Urls:     t.parseEmoteUrls(emote),
	}
}
