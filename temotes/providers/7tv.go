package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"temotes/temotes"
)

type SevenTvFetcher struct{}

type sevenTvGlobalEmoteResponse struct {
	Data struct {
		EmoteSet struct {
			Emotes []sevenTvEmoteResponse `json:"emotes"`
		} `json:"namedEmoteSet"`
	} `json:"data"`
}

type sevenTvChannelEmoteResponse struct {
	Data struct {
		UserByConnection struct {
			Id          string `json:"id"`
			Connections []struct {
				Platform   string `json:"platform"`
				EmoteSetID string `json:"emote_set_id"`
			} `json:"connections"`
			EmoteSets []struct {
				Id     string                 `json:"id"`
				Emotes []sevenTvEmoteResponse `json:"emotes"`
			} `json:"emote_sets"`
		} `json:"userByConnection"`
	} `json:"data"`
}

type sevenTvEmoteResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (t SevenTvFetcher) FetchGlobalEmotes() []temotes.Emote {
	query := `
		query NamedEmoteSet {
			namedEmoteSet(name: GLOBAL) {
				emotes {
					id
					name
				}
			}
		}
	`

	response, err := temotes.CachedFetcher{}.FetchGqlData("https://7tv.io/v3/gql", query, temotes.GlobalEmotesTtl, "7tv-global-emotes")
	var emotes []temotes.Emote
	if err != nil {
		return emotes
	}

	var sevenTvEmotes sevenTvGlobalEmoteResponse
	jsonErr := json.Unmarshal(response, &sevenTvEmotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, sevenTvEmote := range sevenTvEmotes.Data.EmoteSet.Emotes {
		emotes = append(emotes, t.parseEmote(sevenTvEmote))
	}

	return emotes
}

func (t SevenTvFetcher) FetchChannelEmotes(id temotes.TwitchUserId) []temotes.Emote {
	query := fmt.Sprintf(`
		query UserByConnection {
		   userByConnection(platform: TWITCH, id: "%d") {
			   id
			   connections(type: TWITCH) {
				   platform
				   emote_set_id
			   }
			   emote_sets {
				   id
				   emotes {
					   id
					   name
				   }
			   }
		   }
		}
	`, id)

	response, err := temotes.CachedFetcher{}.FetchGqlData("https://7tv.io/v3/gql", query, temotes.ChannelEmotesTtl, fmt.Sprintf("7tv-channel-%d", id))
	var emotes []temotes.Emote
	if err != nil {
		return emotes
	}

	var parsedResponse sevenTvChannelEmoteResponse
	jsonErr := json.Unmarshal(response, &parsedResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	activeEmoteSetId := parsedResponse.Data.UserByConnection.Connections[0].EmoteSetID

	for _, emoteSet := range parsedResponse.Data.UserByConnection.EmoteSets {
		if emoteSet.Id != activeEmoteSetId {
			continue
		}

		for _, sevenTvEmote := range emoteSet.Emotes {
			emotes = append(emotes, t.parseEmote(sevenTvEmote))
		}
	}

	return emotes
}

func (t SevenTvFetcher) parseEmoteUrls(emote sevenTvEmoteResponse) []temotes.EmoteUrl {
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

	urlSizes := []string{"1", "2", "3", "4"}
	for _, size := range urlSizes {
		urls = append(urls, temotes.EmoteUrl{
			Size: getEmoteSize(size),
			Url:  fmt.Sprintf("https://cdn.7tv.app/emote/%s/%sx.webp", emote.Id, size),
		})
	}

	return urls
}

func (t SevenTvFetcher) parseEmote(emote sevenTvEmoteResponse) temotes.Emote {
	return temotes.Emote{
		Provider: temotes.Provider7tv,
		Code:     emote.Name,
		Urls:     t.parseEmoteUrls(emote),
	}
}
