package main

import (
	"flag"
	"fmt"
	"log"
	"pacman/config"
	"pacman/game"
	"time"
)

var configFile = flag.String("configfile", "config.json", "Config file path")
var mazePath = flag.String("mazefile", "maze.txt", "maze file path")

func main() {
	// Parse command-line flags to get configuration file and maze file paths
	flag.Parse()

	// Create a new GameConfiguration instance to load game configuration from JSON
	gc := config.NewGameConfiguration(*configFile)
	gameConfig, err := gc.LoadConfiguration()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Create a new MazeConfiguration instance to load maze layout from a file
	mc := config.NewMazeConfiguration(*mazePath)
	maze, err := mc.LoadMaze()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Create a new game instance with the loaded maze and game configuration
	g := game.NewGame(maze, gameConfig)
	g.Init()
	defer g.Clear()

	// Game loop
	for {
		// Move all sprites concurrently
		for _, s := range g.Sprites {
			go s.Move()
		}

		// Check for game over conditions
		if g.Player().NumDots() == 0 || g.PlayerLives() <= 0 {
			// Display game over message if player runs out of lives
			if g.PlayerLives() <= 0 {
				g.MoveCursor(g.Player().Pos())
				fmt.Print(gameConfig.Death)
				g.MoveCursor(len(maze)+2, 0)
			}
			break
		}

		// Display the updated game screen and wait before rendering the next frame
		g.PrintScreen()
		time.Sleep(1000 / gameConfig.FrameRate * time.Millisecond)
	}
}
