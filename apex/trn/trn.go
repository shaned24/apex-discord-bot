package trn

import (
	"apex_discord_bot/apex"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	URL        string
	APIKey     string
	HttpClient *http.Client
}

func NewTRNClient(url string, apiKey string) *Client {
	return &Client{
		URL:        url,
		APIKey:     apiKey,
		HttpClient: &http.Client{},
	}
}

func (c *Client) GetPlayer(name string, platform string) (*apex.Player, error) {
	resp, err := c.sendRequest(platform, name)
	if err != nil {
		return nil, err
	}

	playerStats, err := c.newPlayerStatsFromResponse(name, resp)
	if err != nil {
		return nil, err
	}

	return &apex.Player{
		Name:     playerStats.Name,
		Platform: playerStats.Platform,
		Legends:  playerStats.GetLegends(),
	}, nil
}

func (c *Client) newPlayerStatsFromResponse(name string, resp *http.Response) (*PlayerStats, error) {
	defer resp.Body.Close()
	playerStats := &PlayerStats{
		Name:     name,
		Platform: "PC",
	}
	err := json.NewDecoder(resp.Body).Decode(playerStats)
	return playerStats, err
}

func (c *Client) sendRequest(platform string, name string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(c.URL, platform, name), nil)
	req.Header.Set("TRN-Api-Key", c.APIKey)
	return c.HttpClient.Do(req)
}

type PlayerStats struct {
	Data     LegendData
	Name     string
	Platform string
}

func (p *PlayerStats) GetLegends() []*apex.Legend {
	var apexLegends []*apex.Legend
	for _, legend := range p.Data.Children {
		apexLegendStats := legend.GetStats()

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

func (l *Legend) GetStats() []*apex.LegendStatistic {
	var apexLegendStats []*apex.LegendStatistic
	for _, stat := range l.Stats {
		apexLegendStat := &apex.LegendStatistic{
			Name:     stat.MetaData.Name,
			Category: stat.MetaData.CategoryName,
			Value:    stat.Value,
		}
		apexLegendStats = append(apexLegendStats, apexLegendStat)
	}
	return apexLegendStats
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
