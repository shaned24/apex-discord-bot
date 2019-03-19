package apex

type PlayerReader interface {
	GetPlayer(name string, platform string) (*Player, error)
}

type Player struct {
	Name     string
	Platform string
	Legends  []*Legend
}

type Legend struct {
	Name    string
	Icon    string
	BGImage string
	Stats   []*LegendStatistic
}

type LegendStatistic struct {
	Name     string
	Value    float64
	Category string
}
