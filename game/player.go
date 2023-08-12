package game

import (
	"fmt"
	"log"
	"os"
)

// Player is the player character \o/
type Player struct {
	position Point
	origin   Point
	img      string
	lives    int
	score    int
	numDots  *int
	maze     []string
}

// NewPlayer creates a new player
func NewPlayer(row, col, lives int, img string, maze []string) *Player {
	var p Player
	p.position = Point{row, col}
	p.origin = Point{row, col}
	p.lives = lives
	p.img = img
	p.maze = maze
	return &p
}

// Pos returns the current position of the player
func (p *Player) Pos() (row, col int) {
	row = p.position.row
	col = p.position.col
	return
}

// NumDots returns the number of remaining dots in the maze
func (p *Player) NumDots() int {
	return *p.numDots
}

// Img returns the image representation of the player
func (p *Player) Img() string {
	return p.img
}

// Kill decreases player lives and resets position if lives are remaining
func (p *Player) Kill() {
	p.lives--
	if p.lives > 0 {
		p.position = p.origin
	}
}

// Move processes player input
func (p *Player) Move() {
	input, err := readInput()
	if err != nil {
		log.Print("Error reading input:", err)
		p.lives = 0
		return
	}

	if input == "ESC" {
		p.lives = 0
		fmt.Println("Exiting the game..")
		os.Exit(1)
	}
	p.movePlayer(input)
}

// movePlayer handles player movement based on input
func (p *Player) movePlayer(dir string) {
	p.position = makeMove(p.position, dir, p.maze)
	row := p.position.row
	col := p.position.col

	removeDot := func(row, col int, maze []string) {
		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
	}

	switch p.maze[row][col] {
	case '.':
		*p.numDots--
		p.score++
		removeDot(row, col, p.maze)
	case 'X':
		p.score += 10
		removeDot(row, col, p.maze)
	}
}

// readInput reads player input from the console
func readInput() (string, error) {
	buffer := make([]byte, 100)

	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}

	return "", nil
}

// makeMove calculates the new position based on the current position and direction
func makeMove(oldPos Point, dir string, maze []string) Point {
	var pos Point
	switch dir {
	case "UP":
		pos, _ = oldPos.Up(maze)
		return pos
	case "DOWN":
		pos, _ = oldPos.Down(maze)
		return pos
	case "LEFT":
		pos, _ = oldPos.Left(maze)
		return pos
	case "RIGHT":
		pos, _ = oldPos.Right(maze)
		return pos
	default:
		return oldPos
	}
}
