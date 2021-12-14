package main

//id i06ys0o81lgxcdmqxvpdxat5p4oinw
//key dtxlvblbnl0y3kzn7q9aoghw8qtcee
//app access token y7out057fc0q9rcw75g0yy96iwv702

import (
	"fmt"
	"temotes/temotes/services"
)

func main() {
	twitch := services.TwitchFetcher{}
	fmt.Println(twitch.FetchGlobalEmotes())
	fmt.Println(twitch.FetchChannelEmotes(27187817))

	bttv := services.BttvFetcher{}
	fmt.Println(bttv.FetchGlobalEmotes())
	fmt.Println(bttv.FetchChannelEmotes(27187817))

	sevenTv := services.SevenTvFetcher{}
	fmt.Println(sevenTv.FetchGlobalEmotes())
	fmt.Println(sevenTv.FetchChannelEmotes(27187817))

	ffz := services.FfzFetcher{}
	fmt.Println(ffz.FetchGlobalEmotes())
	fmt.Println(ffz.FetchChannelEmotes(27187817))
}
