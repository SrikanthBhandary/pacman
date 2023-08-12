package game

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"pacman/config"

	"github.com/danicat/simpleansi"
)

// Sprite represents any game character that can move and has a position
type Sprite interface {
	Move()
	Pos() (int, int)
	Img() string
}

// Game manages the game state and operations
type Game struct {
	maze    []string          // The maze layout
	cfg     config.GameConfig // Game configuration
	player  *Player           // The player character
	Sprites []Sprite          // List of all sprites in the game
}

// NewGame creates a new game instance
func NewGame(maze []string, cfg config.GameConfig) *Game {
	return &Game{maze: maze, cfg: cfg}
}

// PlayerLives returns the remaining lives of the player
func (g *Game) PlayerLives() int {
	return g.player.lives
}

// Player returns the player character
func (g *Game) Player() *Player {
	return g.player
}

// Init initializes the game by setting up the player and sprites
func (g *Game) Init() {
	g.initialise()

	totalDots := 0
	for row, line := range g.maze {
		for col, char := range line {
			switch char {
			case 'P':
				g.player = NewPlayer(row, col, 1, g.cfg.Player, g.maze)
				g.Sprites = append(g.Sprites, g.player)
			// Add other sprite initialization logic here
			case '.':
				totalDots++
			}
		}
	}
	g.player.numDots = &totalDots
}

// PrintScreen prints the current game screen
func (g *Game) PrintScreen() {
	simpleansi.ClearScreen()

	// Print maze with characters based on their representation
	for _, line := range g.maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Print("\x1b[1;35;45m" + g.cfg.Wall + "\x1b[0m")
			case '.':
				fmt.Print(g.cfg.Dot)
			case 'X':
				fmt.Print(g.cfg.Pill)
			default:
				fmt.Print(g.cfg.Space)
			}
		}
		fmt.Println()
	}

	// Print sprites on the maze
	for _, s := range g.Sprites {
		g.MoveCursor(s.Pos())
		fmt.Print(s.Img())
	}

	// Print game-related information
	g.MoveCursor(len(g.maze)+1, 0)
	fmt.Println("Score:", g.player.score, "\tLives:", g.player.lives, "\t Total Dots:", g.player.numDots)
}

// MoveCursor moves the console cursor to the specified row and column
func (g *Game) MoveCursor(row, col int) {
	if g.cfg.UseEmoji {
		simpleansi.MoveCursor(row, col*2)
	} else {
		simpleansi.MoveCursor(row, col)
	}
}

// initialise sets up the console for non-blocking input
func (g *Game) initialise() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin
	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cbreak mode:", err)
	}
}

// Clear resets the console back to its default state
func (g *Game) Clear() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cooked mode:", err)
	}
}
