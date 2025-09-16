package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	DFS      = iota // Depth-First Search
	BFS             // Breadth-First Search
	GBFS            // Greedy Best-First Search
	AStar           // A* Search
	DIJKSTRA        // Dijkstra's Algorithm
)

// Directions for moving in the maze
type Point struct {
	Row int
	Col int
}

// Used to keep track of potential nodes that are walls and cannot be explored
type Wall struct {
	State Point
	wall  bool
}

type Node struct {
	index  int
	State  Point
	Parent *Node
	Action string
}

type Solution struct {
	Action []string
	Cells  []Point
}

// Maze structure to hold the maze data
type Maze struct {
	Height      int      // how tall is the maze
	Width       int      // how wide is the maze
	Start       Point    // starting point
	Goal        Point    // goal point
	Walls       [][]Wall // slice of slices of wall type;
	CurrentNode Node
	Solution    Solution
	Explored    []Point
	Steps       int
	NumExplored int
	Debug       bool
	SearchType  int
}

func main() {
	var m Maze
	var maze, searchType string

	flag.StringVar(&maze, "file", "maze.txt", "maze file")
	flag.StringVar(&searchType, "search", "dfs", "search type")
	flag.Parse()

	err := m.Load(maze)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startTime := time.Now()
	switch searchType {
	case "dfs":
		m.SearchType = DFS
		solveDFS(&m)
	default:
		fmt.Println("Unknown search type:", searchType)
		os.Exit(1)
	}

	if len(m.Solution.Action) > 0 {
		fmt.Println("Solution: ")
		m.printMaze()
		fmt.Println("Solution is", len(m.Solution.Cells), "steps.")
		fmt.Println("Time to solve:", time.Since(startTime))
	} else {
		fmt.Println("No solution found")
	}
	fmt.Println("Explored", len(m.Explored), "nodes")

}
func (g *Maze) printMaze() {
	for i, row := range g.Walls {
		for j, col := range row {
			if col.wall {
				fmt.Print("▉")
			} else if g.Start.Row == col.State.Row && g.Start.Col == col.State.Col {
				fmt.Print("A")
			} else if g.Goal.Row == col.State.Row && g.Goal.Col == col.State.Col {
				fmt.Print("B")
			} else if g.IsSolution(Point{Row: i, Col: j}) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Maze) IsSolution(x Point) bool {
	for _, step := range g.Solution.Cells {
		if step.Row == x.Row && step.Col == x.Col {
			return true
		}
	}
	return false
}

func solveDFS(m *Maze) {
	var s DepthFirstSearch
	s.Game = m
	fmt.Println("Goal is", s.Game.Goal)
	s.Solve()
}

// load the maze
func (g *Maze) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", fileName, err)

	}
	defer f.Close()

	var fileContents []string

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New(fmt.Errorf("cannot open file %s: %v", fileName, err).Error())
		}
		fileContents = append(fileContents, line)

	}

	foundStart, foundEnd := false, false
	for _, line := range fileContents {
		if strings.Contains(line, "A") {
			foundStart = true
		}
		if strings.Contains(line, "B") {
			foundEnd = true
		}
	}
	if !foundStart {
		return errors.New("no start point 'A' found in the maze")
	}
	if !foundEnd {
		return errors.New("no end point 'B' found in the maze")
	}
	g.Height = len(fileContents)
	g.Width = len(fileContents[0])

	var rows [][]Wall

	for i, row := range fileContents {
		var cols []Wall
		for j, col := range row {
			curLetter := fmt.Sprintf("%c", col)
			var wall Wall
			switch curLetter {
			case "A":
				g.Start = Point{Row: i, Col: j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "B":
				g.Goal = Point{Row: i, Col: j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case " ":
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "#":
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = true
			default:
				continue
			}
			cols = append(cols, wall)
		}
		rows = append(rows, cols)
	}
	g.Walls = rows
	return nil
}
