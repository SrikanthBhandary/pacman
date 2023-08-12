// Package config provides functionality to load and parse configuration files for a game.
package config

import (
	"errors"
	"strings"
	"testing"
)

// TestConfig is the test suite for the config package.
func TestGameConfig(t *testing.T) {
	// Define test cases with file paths and expected error values
	var testData = []struct {
		fileName string
		result   error
	}{
		{
			"../test_files/valid.json",
			nil,
		},
		{
			"../test_files/invalid.json",
			errors.New("json: cannot unmarshal string into Go value of type config.GameConfig"),
		},
		{
			"../test_files/invalids.json",
			errors.New("open ../test_files/invalids.json: no such file or directory"),
		},
	}

	// Loop through test cases
	for _, data := range testData {
		// Create a GameConfig instance for the test case
		parser := NewGameConfiguration(data.fileName)
		// Load the configuration and get the error
		_, err := parser.LoadConfiguration()

		// Compare the error with the expected error, if applicable
		if err != nil && !strings.Contains(err.Error(), data.result.Error()) {
			t.Errorf("Test case failed for file %s. Expected %q in error message but found %q", data.fileName, data.result, err)
		}
	}
}

func TestMazeConfig(t *testing.T) {
	// Define test cases with file paths and expected error values
	var testData = []struct {
		fileName string
		result   error
		length   int
	}{
		{
			"../test_files/maze.txt",
			nil,
			24,
		},
	}

	// Loop through test cases
	for _, data := range testData {
		// Create a Maze instance for the test case
		parser := NewMazeConfiguration(data.fileName)
		// Load the configuration and get the error
		maze, err := parser.LoadMaze()

		// Compare the error with the expected error, if applicable
		if err != nil && !strings.Contains(err.Error(), data.result.Error()) {
			t.Errorf("Test case failed for file %s. Expected %q in error message but found %q", data.fileName, data.result, err)
		}
		if len(maze) != data.length {
			t.Errorf("Test case failed for file %s. Expected length of maze is %d  but found %d", data.fileName, data.length, len(maze))
		}
	}
}
