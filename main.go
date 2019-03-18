package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"net/http"
)

const (
	// platform/player_name
	StatsUrl = "https://public-api.tracker.gg/apex/v1/standard/profile/%s/%s"
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

type PlayerStats struct {
	Data Data
	Name string
	Platform string
}

type Data struct {
	ID       string
	Type     string
	Children []Legend
}

type Legend struct {
	ID       string
	Type     string
	MetaData LegendMeta
	Stats    []LegendStats
}

type LegendMeta struct {
	Name    string `json:"legend_name"`
	Icon    string `json:"icon"`
	BGImage string `json:"bgimage"`
}

type LegendStats struct {
	Value        float64
	Percentile   float64
	DisplayValue string
	DisplayRank  string
	MetaData     LegendStatsMeta
}

type LegendStatsMeta struct {
	Key          string
	Name         string
	CategoryKey  string
	CategoryName string
	IsReverse    bool
}

func renderStatsTable(playerStats *PlayerStats) string {
	combinedOutput := "```markdown"
	combinedOutput += fmt.Sprintf("\n# %s - %s", playerStats.Name, playerStats.Platform)

	for _, legend := range playerStats.Data.Children {
		output := fmt.Sprintf("\n## %s \n\n", legend.MetaData.Name)
		output += fmt.Sprintf("![Icon](%s)\n", legend.MetaData.Icon)
		output += fmt.Sprintf("![Background Image](%s) \n", legend.MetaData.BGImage)

		var tableData [][]string
		buffer := new(bytes.Buffer)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader([]string{"Name", "Category", "Value"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, stat := range legend.Stats {
			row := []string{stat.MetaData.Name, stat.MetaData.CategoryName, fmt.Sprintf("%.2f", stat.Value)}
			tableData = append(tableData, row)
		}

		// Create a buffer so we can capture table output to a string later
		table.AppendBulk(tableData) // Add Bulk Data
		table.Render()

		output += "## Stats \n"
		output += buffer.String()

		combinedOutput += output
	}

	combinedOutput += "\n```"

	return combinedOutput
}

func getStats(name string, platform string) (*PlayerStats, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf(StatsUrl, platform, name), nil)
	req.Header.Set("TRN-Api-Key", APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}

	defer resp.Body.Close()

	playerStats := &PlayerStats{
		Name: name,
		Platform: "PC",
	}

	err = json.NewDecoder(resp.Body).Decode(playerStats)
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}


	//for _, legend := range playerStats.Data.Children {
	//	fmt.Println("========= Legend =========")
	//	fmt.Printf("\nID    : %s", legend.ID)
	//	fmt.Printf("\nType  : %s", legend.Type)
	//	fmt.Printf("\nBGImage  : %s", legend.MetaData.BGImage)
	//	fmt.Printf("\nIcon  : %s", legend.MetaData.Icon)
	//	fmt.Printf("\nName  : %s", legend.MetaData.Name)
	//	fmt.Println("========= Stats =========")
	//	for index, stat := range legend.Stats {
	//		fmt.Printf("\n|=======> Stats[%d]", index)
	//
	//		fmt.Printf("\nDisplayRank  : %s", stat.DisplayRank)
	//		fmt.Printf("\nDisplayValue : %s", stat.DisplayValue)
	//		fmt.Printf("\nValue        : %f", stat.Value)
	//		fmt.Printf("\nMetaData.Key     : %s", stat.MetaData.Key)
	//		fmt.Printf("\nMetaData.Name     : %s", stat.MetaData.Name)
	//		fmt.Printf("\nMetaData.CategoryKey     : %s", stat.MetaData.CategoryKey)
	//		fmt.Printf("\nMetaData.CategoryName     : %s", stat.MetaData.CategoryName)
	//		fmt.Printf("\nMetaData.IsReverse     : %t", stat.MetaData.IsReverse)
	//		fmt.Printf("\nPercentile   : %f", stat.Percentile)
	//
	//	}
	//}
	//fmt.Printf("\nResponse: %v", playerStats.Data.Children[0].ID)


	return playerStats, nil
}

func main() {
	playerStats, _ := getStats(PlayerName, Platform)
	fmt.Println(renderStatsTable(playerStats))
}
