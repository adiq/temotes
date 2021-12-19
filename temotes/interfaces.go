package temotes

type TwitchUserId int

type EmoteFetcher interface {
	FetchGlobalEmotes() []Emote
	FetchChannelEmotes(id TwitchUserId) []Emote
}

//type PersonalDataFetcher interface {
//	FetchPersonalAvatar(id TwitchUserId)
//	FetchPersonalEmotes(id TwitchUserId) []Emote
//}
