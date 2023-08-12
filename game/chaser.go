package game

// Chaser represents a ghost character that chases a player in a maze-based game.
type Chaser struct {
	position Point    // Current position of the ghost.
	img      string   // Image filename for the ghost.
	path     []string // List of directions for the ghost to follow.
	sprites  []Sprite // Sprites associated with the ghost.
	maze     []string // Maze layout.
	player   *Player  // Reference to the player being chased.
}

// NewChaser creates a new ghost character at the specified position in the maze.
func NewChaser(row, col int, img string, maze []string) *Chaser {
	return &Chaser{position: Point{row, col}, img: img, maze: maze}
}

// Pos returns the current row and column indices of the ghost's position.
func (c *Chaser) Pos() (int, int) {
	return c.position.row, c.position.col
}

// SetPlayer sets the player that the ghost will chase.
func (c *Chaser) SetPlayer(player *Player) {
	c.player = player
}

// SetSprites sets the sprites associated with the ghost.
func (c *Chaser) SetSprites(sprites []Sprite) {
	c.sprites = sprites
}

// Img returns the image filename associated with the ghost.
func (c *Chaser) Img() string {
	return c.img
}

// Move updates the ghost's position based on its chasing behavior.
func (c *Chaser) Move() {
	dir := c.drawDirection()                       // Calculate the direction to move.
	c.position = makeMove(c.position, dir, c.maze) // Update position.

	// Check for collisions with other sprites, specifically the player.
	for _, s := range c.sprites {
		switch p := s.(type) {
		case *Player:
			if p.position == c.position {
				p.Kill()     // Kill the player on collision.
				c.path = nil // Reset the ghost's path to recalculate.
			}
		}
	}
}

// drawDirection calculates the direction in which the ghost should move to chase the player.
func (c *Chaser) drawDirection() string {
	if len(c.path) == 0 {
		target := c.player.position
		c.path = c.find(c.position, target) // Calculate the path to the player.
	}
	dir := c.path[0]
	c.path = c.path[1:]
	return dir
}

// find implements the A* pathfinding algorithm to find the shortest path between two points.
func (c *Chaser) find(origin Point, target Point) []string {
	var pf PathFinder
	pf.maze = c.maze
	path := pf.walk(origin, target) // Calculate the path using A* algorithm.

	var directions []string
	current := origin
	for len(path) > 0 {
		next := path[len(path)-1]
		path = path[:len(path)-1]
		directions = append(directions, giveDirection(current, next))
		current = next
	}

	return directions
}

// giveDirection determines the direction from current to destination point.
func giveDirection(curr, dest Point) string {
	switch {
	case curr.row-dest.row == 0 && dest.col-curr.col == 1:
		return "RIGHT"
	case curr.row-dest.row == 0 && dest.col-curr.col == -1:
		return "LEFT"
	case curr.col-dest.col == 0 && dest.row-curr.row == 1:
		return "DOWN"
	case curr.col-dest.col == 0 && dest.row-curr.row == -1:
		return "UP"
	default:
		return "NOP" // No operation.
	}
}

// PathFinder manages the A* pathfinding process.
type PathFinder struct {
	open   PointSet
	closed PointSet
	table  map[Point]PointInfo
	maze   []string
}

// PointSet represents a set of points.
type PointSet map[Point]bool

// PointInfo holds information about a point's cost and parent for pathfinding.
type PointInfo struct {
	g      int
	h      int
	parent *Point
}

// Cost calculates the total cost for the point.
func (p PointInfo) Cost() int {
	return p.g + p.h
}

// nextPoint finds the next point with the lowest cost from the open set.
func (pf *PathFinder) nextPoint() *Point {
	var min *Point
	for k := range pf.open {
		if min == nil {
			min = &k
			continue
		}
		if pf.table[k].Cost() < pf.table[*min].Cost() {
			min = &k
		}
	}
	return min
}

// abs calculates the absolute value of an integer.
func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

// distance calculates the Manhattan distance between two points.
func distance(p1, p2 Point) int {
	return abs(p1.row-p2.row) + abs(p1.col-p2.col)
}

// isClosed checks if a point is in the closed set.
func (pf *PathFinder) isClosed(p Point) bool {
	_, ok := pf.closed[p]
	return ok
}

// isOpen checks if a point is in the open set.
func (pf *PathFinder) isOpen(p Point) bool {
	_, ok := pf.open[p]
	return ok
}

// walk implements the A* pathfinding algorithm to find a path from start to target.
func (pf *PathFinder) walk(start Point, target Point) []Point {
	var path []Point

	pf.open = make(PointSet)
	pf.closed = make(PointSet)
	pf.table = make(map[Point]PointInfo)

	pf.open[start] = true
	pf.table[start] = PointInfo{
		h: distance(start, target),
	}
	for current := pf.nextPoint(); current != nil; current = pf.nextPoint() {
		// Determine valid neighboring points.
		var neighbors []Point
		if up, err := current.Up(pf.maze); err == nil {
			neighbors = append(neighbors, up)
		}
		if down, err := current.Down(pf.maze); err == nil {
			neighbors = append(neighbors, down)
		}
		if left, err := current.Left(pf.maze); err == nil {
			neighbors = append(neighbors, left)
		}
		if right, err := current.Right(pf.maze); err == nil {
			neighbors = append(neighbors, right)
		}

		// Process neighboring points.
		for _, n := range neighbors {
			switch {
			case n == target:
				// Found path to target.
				path = append(path, n)
				for pf.table[*current].parent != nil {
					path = append(path, *current)
					current = pf.table[*current].parent
				}
				return path
			case pf.isClosed(n):
				// Point is in the closed set, skip.
			case pf.isOpen(n):
				info := PointInfo{
					g:      pf.table[*current].g + 1,
					h:      distance(n, target),
					parent: current,
				}
				if pf.table[n].Cost() > info.Cost() {
					// Found a better cost to n, update.
					pf.table[n] = info
				}
			default:
				pf.open[n] = true
				pf.table[n] = PointInfo{
					g:      pf.table[*current].g + 1,
					h:      distance(n, target),
					parent: current,
				}
			}
		}

		// Done with current point.
		// Add to closed set, remove from open set.
		pf.closed[*current] = true
		delete(pf.open, *current)
	}
	return path
}
