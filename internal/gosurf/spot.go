package gosurf

import (
	"fmt"
	"time"
)

const (
	HighTide = "high"
	AllTide  = "all"
	LowTide  = "low"
)

type Spot struct {
	Name        string  `yaml:"name"`
	Lng         float64 `yaml:"lng"`
	Lat         float64 `yaml:"lat"`
	Tide        Tide    `yaml:"tide"`
	Period      float64 `yaml:"period"`
	WaveHeight  float64 `yaml:"waveHeight"`
	Temperature Temperature
	WindSpeed   float64
}

type Tide struct {
	Position   string    `yaml:"pos"`
	TimeBefore int       `yaml:"before"`
	TimeAfter  int       `yalm:"after"`
	High       time.Time // get via api not config
	Low        time.Time // get via api not config
}

type Temperature struct {
	Air   float64
	Water float64
}

func (s Spot) String() string {
	return fmt.Sprintf("\n\nSpot: %v\nDate: %v\nPeriod: %vs\nHeight: %vm\nTide: %v\nAir Temp: %vc\nWater Temp: %vc",
		s.Name, s.Tide.High.Format("2006-01-02"), s.Period, s.WaveHeight, s.Tide.Position, s.Temperature.Air, s.Temperature.Water)
}
