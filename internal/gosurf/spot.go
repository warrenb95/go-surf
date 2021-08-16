package gosurf

import "fmt"

const (
	HighTide = "HIGH"
	AllTide  = "ALL"
	LowTide  = "LOW"
)

type Spot struct {
	Name          string  `yaml:"name"`
	Lng           float64 `yaml:"lng"`
	Lat           float64 `yaml:"lat"`
	Tide          Tide    `yaml:"tide"`
	Period        int     `yaml:"period"`
	MinWaveHeight int     `yaml:"minWaveHeight"`
}

type Tide struct {
	Position   string `yaml:"pos"`
	TimeBefore int    `yaml:"before"`
	TimeAfter  int    `yalm:"after"`
}

func (s Spot) String() string {
	return fmt.Sprintf("Spot: %v\nStart: %v\nEnd: %v\nPeriod: %v\nHeight: %v",
		s.Name, s.Tide.TimeBefore, s.Tide.TimeAfter, s.Period, s.MinWaveHeight)
}
