// Package config provides functionality to load and parse configuration files for a game.
package config

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

// GameConfig represents the configuration settings for a game.
type GameConfig struct {
	Player    string        `json:"player"`
	Ghost     string        `json:"ghost"`
	Wall      string        `json:"wall"`
	Dot       string        `json:"dot"`
	Pill      string        `json:"pill"`
	Death     string        `json:"death"`
	Chaser    string        `json:"chaser"`
	Space     string        `json:"space"`
	UseEmoji  bool          `json:"use_emoji"`
	FrameRate time.Duration `json:"frame_rate"`
}

// GameConfiguration represents a configuration loader for a game.
type GameConfiguration struct {
	FilePath string
}

// NewGameConfiguration creates a new GameConfiguration instance with the provided file path.
func NewGameConfiguration(filePath string) GameConfiguration {
	return GameConfiguration{FilePath: filePath}
}

// LoadConfiguration loads the game configuration from a JSON file.
// It returns a GameConfig struct and an error if any.
func (j GameConfiguration) LoadConfiguration() (GameConfig, error) {
	var conf GameConfig
	file, err := os.Open(j.FilePath)
	if err != nil {
		return conf, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	return conf, err
}

// MazeConfiguration represents a configuration loader for a maze.
type MazeConfiguration struct {
	Filepath string
}

// NewMazeConfiguration creates a new MazeConfiguration instance with the provided file path.
func NewMazeConfiguration(path string) MazeConfiguration {
	return MazeConfiguration{Filepath: path}
}

// LoadMaze loads the maze configuration from a text file.
// It returns a slice of strings representing each line of the maze and an error if any.
func (m MazeConfiguration) LoadMaze() ([]string, error) {
	var maze []string
	file, err := os.Open(m.Filepath)
	if err != nil {
		return maze, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}
	return maze, nil
}
