package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

const TimeMillis = "15:04:05.000"

func ParseConfig(pathToFile string) Config {
	data, err := os.ReadFile(pathToFile)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config
}
