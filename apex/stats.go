package apex

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
)

type PlayerStatsFetcher interface {
	Fetch(name string, platform string) (*Player, error)
}

type PlayerStatsRenderer interface {
	Render(legend *Player) (string, error)
}

type MarkdownRenderer struct{}

func (m *MarkdownRenderer) Render(player *Player) (string, error) {

	combinedOutput := "```markdown"
	combinedOutput += fmt.Sprintf("\n# %s - %s", player.Name, player.Platform)

	for _, legend := range player.Legends {
		output := fmt.Sprintf("\n## %s \n\n", legend.Name)
		output += fmt.Sprintf("![Icon](%s)\n", legend.Icon)
		output += fmt.Sprintf("![Background Image](%s) \n", legend.BGImage)

		var tableData [][]string
		buffer := new(bytes.Buffer)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader([]string{"Name", "Category", "Value"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, stat := range legend.Stats {
			row := []string{stat.Name, stat.Category, fmt.Sprintf("%.2f", stat.Value)}
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

	return combinedOutput, nil
}
