package trn

import (
	"apex_discord_bot/apex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	URL    string
	APIKey string
}

func NewTRNClient(url string, apiKey string) *Client {
	return &Client{
		URL:    url,
		APIKey: apiKey,
	}
}

func (c *Client) Fetch(name string, platform string) (*apex.Player, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf(c.URL, platform, name), nil)
	req.Header.Set("TRN-Api-Key", c.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}

	defer resp.Body.Close()

	playerStats := &PlayerStats{
		Name:     name,
		Platform: "PC",
	}

	err = json.NewDecoder(resp.Body).Decode(playerStats)
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}

	apexLegends := getLegends(playerStats)

	return &apex.Player{
		Name:     playerStats.Name,
		Platform: playerStats.Platform,
		Legends:  apexLegends,
	}, nil
}

func getLegends(playerStats *PlayerStats) []*apex.Legend {
	var apexLegends []*apex.Legend
	for _, legend := range playerStats.Data.Children {
		apexLegendStats := getLegendStats(legend)

		apexLegend := &apex.Legend{
			Name:    legend.MetaData.Name,
			BGImage: legend.MetaData.BGImage,
			Icon:    legend.MetaData.Icon,
			Stats:   apexLegendStats,
		}

		apexLegends = append(apexLegends, apexLegend)
	}
	return apexLegends
}

func getLegendStats(legend Legend) []*apex.LegendStatistic {
	var apexLegendStats []*apex.LegendStatistic
	for _, stat := range legend.Stats {
		apexLegendStat := &apex.LegendStatistic{
			Name:     stat.MetaData.Name,
			Category: stat.MetaData.CategoryName,
			Value:    stat.Value,
		}
		apexLegendStats = append(apexLegendStats, apexLegendStat)
	}
	return apexLegendStats
}

type PlayerStats struct {
	Data     LegendData
	Name     string
	Platform string
}

type LegendData struct {
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
