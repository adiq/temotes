package temotes

type EmoteSize string

const (
	Size1x EmoteSize = "1x"
	Size2x           = "2x"
	Size3x           = "3x"
	Size4x           = "4x"
)

type EmoteProvider int

const (
	ProviderTwitch EmoteProvider = 0
	Provider7tv                  = 1
	ProviderBttv                 = 2
	ProviderFfz                  = 3
)

type EmoteUrl struct {
	Size EmoteSize `json:"size"`
	Url  string    `json:"url"`
}

type Emote struct {
	ProviderEmoteID string        `json:"provider_emote_id"`
	Provider        EmoteProvider `json:"provider"`
	Code            string        `json:"code"`
	Urls            []EmoteUrl    `json:"urls"`
}
