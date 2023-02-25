package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"temotes/temotes"
)

type BttvFetcher struct{}

type bttvEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func (t BttvFetcher) FetchGlobalEmotes() []temotes.Emote {
	response, err := temotes.CachedFetcher{}.FetchData("https://api.betterttv.net/3/cached/emotes/global", temotes.GlobalEmotesTtl, "bttv-global-emotes")
	var emotes []temotes.Emote
	if err != nil {
		return emotes
	}

	var bttvEmotes []bttvEmote
	jsonErr := json.Unmarshal(response, &bttvEmotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

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
	response, err := temotes.CachedFetcher{}.FetchData(fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%d", id), temotes.ChannelEmotesTtl, fmt.Sprintf("bttv-channel-emotes-%d", id))
	var emotes []temotes.Emote
	if err != nil {
		return emotes
	}

	var bttvEmotesResponse bttvChannelEmotesResponse
	jsonErr := json.Unmarshal(response, &bttvEmotesResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

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
