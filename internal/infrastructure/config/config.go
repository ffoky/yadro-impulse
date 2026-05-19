package config

import (
	"encoding/json"
	"fmt"
	"os"

	"yadro-impulse/internal/domain"
)

type fileFormat struct {
	Floors    int    `json:"Floors"`
	Monsters  int    `json:"Monsters"`
	OpenAt    string `json:"OpenAt"`
	Duration  int    `json:"Duration"`
}

func Load(path string) (domain.DungeonConfig, error) {
	raw, err := readFile(path)
	if err != nil {
		return domain.DungeonConfig{}, err
	}
	openAt, err := domain.ParseTime(raw.OpenAt)
	if err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("parse OpenAt: %w", err)
	}
	return domain.DungeonConfig{
		Floors:   raw.Floors,
		Monsters: raw.Monsters,
		OpenAt:   openAt,
		Duration: raw.Duration,
	}, nil
}

func readFile(path string) (fileFormat, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileFormat{}, fmt.Errorf("read config: %w", err)
	}
	var raw fileFormat
	if err = json.Unmarshal(data, &raw); err != nil {
		return fileFormat{}, fmt.Errorf("parse config: %w", err)
	}
	return raw, nil
}
