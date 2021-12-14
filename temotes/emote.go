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
	Size EmoteSize
	Url  string
}

type Emote struct {
	Provider EmoteProvider
	Code     string
	Urls     []EmoteUrl
}
