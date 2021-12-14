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

type BttvFetcher struct{}

type bttvEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func (t BttvFetcher) fetchEmotes(url string) []byte {
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

	return body
}

func (t BttvFetcher) FetchGlobalEmotes() []temotes.Emote {
	response := t.fetchEmotes("https://api.betterttv.net/3/cached/emotes/global")

	var bttvEmotes []bttvEmote
	jsonErr := json.Unmarshal(response, &bttvEmotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var emotes []temotes.Emote
	for _, bttvEmote := range bttvEmotes {
		emotes = append(emotes, t.parseEmote(bttvEmote))
	}

	return emotes
}

type bttvChannelEmotesResponse struct {
	ChannelEmotes []bttvEmote `json:"channelEmotes"`
	SharedEmotes  []bttvEmote `json:"sharedEmotes"`
}

func (t BttvFetcher) FetchChannelEmotes(id temotes.TwitchUserId) []temotes.Emote {
	response := t.fetchEmotes(fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%d", id))

	var bttvEmotesResponse bttvChannelEmotesResponse
	jsonErr := json.Unmarshal(response, &bttvEmotesResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var emotes []temotes.Emote
	for _, bttvEmote := range bttvEmotesResponse.ChannelEmotes {
		emotes = append(emotes, t.parseEmote(bttvEmote))
	}

	for _, bttvEmote := range bttvEmotesResponse.SharedEmotes {
		emotes = append(emotes, t.parseEmote(bttvEmote))
	}

	return emotes
}

func (t BttvFetcher) parseEmoteUrls(emote bttvEmote) []temotes.EmoteUrl {
	return []temotes.EmoteUrl{
		{
			Size: temotes.Size1x,
			Url:  fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", emote.ID),
		},
		{
			Size: temotes.Size2x,
			Url:  fmt.Sprintf("https://cdn.betterttv.net/emote/%s/2x", emote.ID),
		},
		{
			Size: temotes.Size3x,
			Url:  fmt.Sprintf("https://cdn.betterttv.net/emote/%s/3x", emote.ID),
		},
	}
}

func (t BttvFetcher) parseEmote(emote bttvEmote) temotes.Emote {
	return temotes.Emote{
		Provider: temotes.ProviderBttv,
		Code:     emote.Code,
		Urls:     t.parseEmoteUrls(emote),
	}
}
