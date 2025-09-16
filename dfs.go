package main

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
)

type DepthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

// Return the current frontier
func (dfs *DepthFirstSearch) GetFrontier() []*Node {
	return dfs.Frontier
}

// Add a node to the end of the slice (LIFO)
func (dfs *DepthFirstSearch) Add(i *Node) {
	dfs.Frontier = append(dfs.Frontier, i)
}

// Remove and return the last node (LIFO)
func (dfs *DepthFirstSearch) ContainsState(i *Node) bool {
	for _, x := range dfs.Frontier {
		if x.State == i.State {
			return true
		}
	}
	return false
}

// Check if the frontier is empty
func (dfs *DepthFirstSearch) Empty() bool {
	return len(dfs.Frontier) == 0
}

// Remove and return the last node (LIFO)
func (dfs *DepthFirstSearch) Remove() (*Node, error) {
	if len(dfs.Frontier) > 0 {
		if dfs.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range dfs.Frontier {
				fmt.Println("Node:", x.State)
			}
		}
		node := dfs.Frontier[len(dfs.Frontier)-1]         // get last item
		dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1] // delete last item
		return node, nil
	}
	return nil, errors.New("empty frontier")

}

// this is what solves the maze
func (dfs *DepthFirstSearch) Solve() {
	fmt.Println("Starting to solve  maze with Depth First Search")
	dfs.Game.NumExplored = 0
	start := Node{
		State:  dfs.Game.Start,
		Parent: nil,
		Action: "",
	}
	dfs.Add(&start)

	// where am i
	dfs.Game.CurrentNode = start

	for {
		if dfs.Empty() {
			return
		}
		currentNode, err := dfs.Remove()
		if err != nil {
			fmt.Println("Error removing node from frontier:", err)
			return
		}
		if dfs.Game.Debug {
			fmt.Println("Removed:", currentNode.State)
			fmt.Println("---------")
			fmt.Println("")
		}
		dfs.Game.CurrentNode = *currentNode
		dfs.Game.NumExplored++

		// have we reached the goal?
		if dfs.Game.Goal == currentNode.State {
			var actions []string
			var cells []Point

			for {
				if currentNode.Parent != nil {
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}
			// reverse the actions and cells
			slices.Reverse(actions)
			slices.Reverse(cells)

			dfs.Game.Solution = Solution{
				Action: actions,
				Cells:  cells,
			}

			// add the starting point to the solution path
			dfs.Game.Explored = append(dfs.Game.Explored, dfs.Game.CurrentNode.State)
			return
		}
		dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)

		// add neighbors to the frontier
		for _, x := range dfs.Neighbors(currentNode) {
			if !dfs.ContainsState(x) {
				if !inExplored(x.State, dfs.Game.Explored) {
					dfs.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}
	}
}

func (dfs *DepthFirstSearch) Neighbors(node *Node) []*Node {
	row := node.State.Row
	col := node.State.Col
	candidates := []*Node{
		{
			State:  Point{Row: row - 1, Col: col},
			Parent: node,
			Action: "up",
		},
		{
			State:  Point{Row: row, Col: col - 1},
			Parent: node,
			Action: "left",
		},
		{
			State:  Point{Row: row, Col: col + 1},
			Parent: node,
			Action: "right",
		},
		{
			State:  Point{Row: row + 1, Col: col},
			Parent: node,
			Action: "down",
		},
	}
	var neighbors []*Node
	for _, x := range candidates {
		if 0 <= x.State.Row && x.State.Row < dfs.Game.Height {
			if 0 <= x.State.Col && x.State.Col < dfs.Game.Width {
				if !dfs.Game.Walls[x.State.Row][x.State.Col].wall {
					neighbors = append(neighbors, x)
				}
			}
		}
	}
	for i := range neighbors {
		j := rand.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}
	return neighbors
}
