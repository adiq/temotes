package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"temotes/temotes"
	"time"
)

type TwitchFetcher struct{}

type twitchEmote struct {
	ID     string   `json:"id"`
	Code   string   `json:"name"`
	Themes []string `json:"theme_mode"`
	Scales []string `json:"scale"`
}

type twitchEmoteResponse struct {
	Data []twitchEmote `json:"data"`
}

func getAuthorizedRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITCH_ACCESS_TOKEN")))

	return req
}

func (t TwitchFetcher) fetchEmotes(url string, ttl time.Duration, cacheKey string) []temotes.Emote {
	request := getAuthorizedRequest(url)
	response := temotes.CachedFetcher{}.FetchDataRequest(request, ttl, cacheKey)
	var twitchEmotes twitchEmoteResponse
	jsonErr := json.Unmarshal(response, &twitchEmotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var emotes []temotes.Emote
	for _, twitchEmote := range twitchEmotes.Data {
		emotes = append(emotes, t.parseEmote(twitchEmote))
	}

	return emotes
}

func (t TwitchFetcher) FetchGlobalEmotes() []temotes.Emote {
	return t.fetchEmotes("https://api.twitch.tv/helix/chat/emotes/global", temotes.GlobalEmotesTtl, "twitch-global-emotes")
}

func (t TwitchFetcher) FetchChannelEmotes(id temotes.TwitchUserId) []temotes.Emote {
	return t.fetchEmotes(fmt.Sprintf("https://api.twitch.tv/helix/chat/emotes?broadcaster_id=%d", id), temotes.ChannelEmotesTtl, fmt.Sprintf("twitch-channel-emotes-%d", id))
}

func (t TwitchFetcher) parseEmoteUrls(emote twitchEmote) []temotes.EmoteUrl {
	var urls []temotes.EmoteUrl

	getEmoteSize := func(scale string) temotes.EmoteSize {
		switch scale {
		case "1.0":
			return temotes.Size1x
		case "2.0":
			return temotes.Size2x
		case "3.0":
			return temotes.Size4x
		default:
			return temotes.Size1x
		}
	}

	getEmoteTheme := func(themes []string) string {
		if len(themes) == 0 {
			panic("Twitch Emote Error: No themes defined")
		}

		if temotes.Contains(emote.Themes, "light") {
			return "light"
		} else {
			return emote.Themes[0]
		}
	}

	theme := getEmoteTheme(emote.Themes)
	for _, scale := range emote.Scales {
		urls = append(urls, temotes.EmoteUrl{
			Size: getEmoteSize(scale),
			Url:  fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/%s/%s", emote.ID, theme, scale),
		})
	}

	return urls
}

func (t TwitchFetcher) parseEmote(emote twitchEmote) temotes.Emote {
	return temotes.Emote{
		Provider: temotes.ProviderTwitch,
		Code:     emote.Code,
		Urls:     t.parseEmoteUrls(emote),
	}
}

type twitchUsersResponse struct {
	Data []twitchUser `json:"data"`
}

type twitchUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"profile_image_url"`
}

type TwitchUser struct {
	ID          temotes.TwitchUserId `json:"id"`
	Login       string               `json:"login"`
	DisplayName string               `json:"display_name"`
	Avatar      string               `json:"avatar"`
}

func (t TwitchFetcher) FetchUserIdentifiers(identifier string) (*TwitchUser, error) {
	id, err := strconv.ParseInt(strings.ToLower(identifier), 10, 64)
	var request *http.Request
	var cacheKey string

	if id == 0 || err != nil {
		request = getAuthorizedRequest(fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", identifier))
		cacheKey = fmt.Sprintf("twitch-user-identifiers-login-%s", identifier)
	} else {
		request = getAuthorizedRequest(fmt.Sprintf("https://api.twitch.tv/helix/users?id=%d", id))
		cacheKey = fmt.Sprintf("twitch-user-identifiers-id-%d", id)
	}

	response := temotes.CachedFetcher{}.FetchDataRequest(request, temotes.TwitchIdTtl, cacheKey)
	var twitchUsers twitchUsersResponse
	jsonErr := json.Unmarshal(response, &twitchUsers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if len(twitchUsers.Data) == 0 {
		return nil, errors.New("user not found")
	}

	userId, err := strconv.ParseInt(twitchUsers.Data[0].ID, 10, 64)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user := &TwitchUser{
		ID:          temotes.TwitchUserId(userId),
		Login:       twitchUsers.Data[0].Login,
		DisplayName: twitchUsers.Data[0].DisplayName,
		Avatar:      twitchUsers.Data[0].Avatar,
	}

	return user, nil
}
