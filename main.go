package main

import (
	"apex_discord_bot/apex"
	"apex_discord_bot/apex/trn"
	"flag"
	"fmt"
)

const (
	// platform/player_name
	TRNStatsUrl = "https://public-api.tracker.gg/apex/v1/standard/profile/%s/%s"
)

var (
	APIKey     string
	Platform   string
	PlayerName string
)

func init() {
	flag.StringVar(&Platform, "platform", "5", "1 = XBOX 2 = PSN 5 = Origin / PC")
	flag.StringVar(&PlayerName, "player", "snakerd", "The Player Name")
	flag.StringVar(&APIKey, "token", "", "The API Token")
	flag.Parse()
}

func GetPlayer(reader apex.PlayerReader) (*apex.Player, error) {
	return reader.GetPlayer(PlayerName, Platform)
}

func main() {
	renderer := &apex.MarkdownRenderer{}
	trnClient := trn.NewTRNClient(TRNStatsUrl, APIKey)

	player, _ := GetPlayer(trnClient)
	markdownOutput, _ := renderer.Render(player)

	fmt.Println(markdownOutput)
}
