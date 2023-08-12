// Package point provides utility functions for manipulating positions in a grid.
package game

import "errors"

var errInvalidPosition = errors.New("invalid position")

// Point represents a position with row and column coordinates.
type Point struct {
	row int
	col int
}

// Up moves the point up one row within the given maze, wrapping around if needed.
// It returns the new position after moving and an error if the position is invalid.
func (p Point) Up(maze []string) (Point, error) {
	p.row--
	if p.row < 0 {
		p.row = len(maze) - 1
	}
	if !IsLegal(p, maze) {
		p.row++
		return p, errInvalidPosition
	}
	return p, nil
}

// Down moves the point down one row within the given maze, wrapping around if needed.
// It returns the new position after moving and an error if the position is invalid.
func (p Point) Down(maze []string) (Point, error) {
	p.row++
	if p.row == len(maze) {
		p.row = 0
	}
	if !IsLegal(p, maze) {
		p.row--
		return p, errInvalidPosition
	}
	return p, nil
}

// Left moves the point left one column within the given maze, wrapping around if needed.
// It returns the new position after moving and an error if the position is invalid.
func (p Point) Left(maze []string) (Point, error) {
	p.col--
	if p.col < 0 {
		p.col = len(maze[0]) - 1
	}
	if !IsLegal(p, maze) {
		p.col++
		return p, errInvalidPosition
	}
	return p, nil
}

// Right moves the point right one column within the given maze, wrapping around if needed.
// It returns the new position after moving and an error if the position is invalid.
func (p Point) Right(maze []string) (Point, error) {
	p.col++
	if p.col == len(maze[0]) {
		p.col = 0
	}
	if !IsLegal(p, maze) {
		p.col--
		return p, errInvalidPosition
	}
	return p, nil
}

// IsLegal checks if the given position is within the bounds of the maze and not blocked
// by an obstacle.
func IsLegal(pos Point, maze []string) bool {
	return pos.row >= 0 && pos.row < len(maze) &&
		pos.col >= 0 && pos.col < len(maze[0]) &&
		maze[pos.row][pos.col] != '#'
}
