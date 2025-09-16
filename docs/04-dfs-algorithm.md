# Depth-First Search Algorithm Implementation

## Understanding DFS Step by Step

This document walks through the complete Depth-First Search implementation in the maze solver. You'll see how every piece works and why each decision was made.

## What DFS Does

**Depth-First Search explores as far as possible along each path before backtracking.**

Real-world analogy: **Exploring a cave system**

- Pick a tunnel and follow it as deep as you can go
- When you hit a dead end, backtrack to the last intersection
- Try a different tunnel
- Continue until you find what you're looking for

## The Complete DFS Structure

```go
type DepthFirstSearch struct {
    Frontier []*Node    // Stack of positions to explore
    Game     *Maze      // Reference to the maze
}
```

**Why this design:**

- **Frontier** acts as our stack (LIFO behavior)
- **Game** gives us access to maze data (walls, start, goal)
- Simple structure focusing on the algorithm logic

## Core DFS Operations

### Adding Positions to Explore

```go
func (dfs *DepthFirstSearch) Add(node *Node) {
    dfs.Frontier = append(dfs.Frontier, node)
}
```

**What this does:** Adds a new position to the end of our stack

**Why append to end:** DFS uses LIFO (Last In, First Out) behavior

### Getting the Next Position

```go
func (dfs *DepthFirstSearch) Remove() (*Node, error) {
    if len(dfs.Frontier) > 0 {
        // Get the last item (top of stack)
        node := dfs.Frontier[len(dfs.Frontier)-1]

        // Remove the last item
        dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1]

        return node, nil
    }
    return nil, errors.New("empty frontier")
}
```

**Key points:**

- Always takes from the **end** of the slice (stack behavior)
- Returns an error if frontier is empty (no more places to explore)
- This is what makes it "depth-first" - newest positions explored first

### Checking if Done

```go
func (dfs *DepthFirstSearch) Empty() bool {
    return len(dfs.Frontier) == 0
}
```

**Simple but important:** When frontier is empty, we've explored all reachable positions

## The Main Algorithm: Solve Method

### Step 1: Initialize

```go
func (dfs *DepthFirstSearch) Solve() {
    fmt.Println("Starting to solve maze with Depth First Search")

    // Create starting node
    start := Node{
        State:  dfs.Game.Start,    // Starting position
        Parent: nil,               // No parent (this is the beginning)
        Action: "",                // No action to get here
    }

    // Add to frontier
    dfs.Add(&start)
}
```

**Why we start here:**

- Every search needs a beginning point
- Starting node has no parent (it's the root)
- Adding to frontier kicks off the exploration

### Step 2: Main Search Loop

```go
for {
    // Check if we have anywhere left to explore
    if dfs.Empty() {
        return    // No solution exists
    }

    // Get next position to explore
    currentNode, err := dfs.Remove()
    if err != nil {
        fmt.Println("Error removing node from frontier:", err)
        return
    }

    // Update game state
    dfs.Game.CurrentNode = *currentNode
    dfs.Game.NumExplored++
```

**The infinite loop pattern:**

- `for {` creates an infinite loop
- We break out when we find the goal or run out of options
- This is a common pattern in search algorithms

### Step 3: Goal Check

```go
// Are we at the goal?
if dfs.Game.Goal == currentNode.State {
    // BUILD SOLUTION PATH
    var actions []string
    var cells []Point

    // Follow parent pointers backwards
    for {
        if currentNode.Parent != nil {
            actions = append(actions, currentNode.Action)
            cells = append(cells, currentNode.State)
            currentNode = currentNode.Parent
        } else {
            break    // Reached the starting node
        }
    }

    // Reverse to get start-to-goal path
    slices.Reverse(actions)
    slices.Reverse(cells)

    // Store solution
    dfs.Game.Solution = Solution{
        Action: actions,
        Cells:  cells,
    }

    return    // Success! We found the goal
}
```

**Path reconstruction explained:**

1. **Follow parents backwards:** Goal → Parent → Parent → Start
2. **Collect actions:** ["down", "right", "right"] (backwards)
3. **Reverse arrays:** ["right", "right", "down"] (correct order)
4. **Result:** Complete path from start to goal

### Step 4: Mark as Explored

```go
// Remember we've been here
dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)
```

**Why this matters:** Prevents infinite loops by avoiding revisiting positions

### Step 5: Explore Neighbors

```go
// Add neighbors to frontier
for _, neighbor := range dfs.Neighbors(currentNode) {
    // Don't add if already in frontier
    if !dfs.ContainsState(neighbor) {
        // Don't add if already explored
        if !inExplored(neighbor.State, dfs.Game.Explored) {
            dfs.Add(&Node{
                State:  neighbor.State,
                Parent: currentNode,        // Remember how we got here
                Action: neighbor.Action,
            })
        }
    }
}
```

**The neighbor exploration process:**

1. **Get all possible moves** from current position
2. **Filter out duplicates** (already in frontier)
3. **Filter out visited positions** (already explored)
4. **Add valid neighbors** with parent links

## Neighbor Generation: Finding Valid Moves

```go
func (dfs *DepthFirstSearch) Neighbors(node *Node) []*Node {
    row := node.State.Row
    col := node.State.Col

    // Define all possible moves
    candidates := []*Node{
        {State: Point{Row: row - 1, Col: col}, Action: "up"},
        {State: Point{Row: row, Col: col - 1}, Action: "left"},
        {State: Point{Row: row, Col: col + 1}, Action: "right"},
        {State: Point{Row: row + 1, Col: col}, Action: "down"},
    }

    var neighbors []*Node

    // Check each candidate
    for _, candidate := range candidates {
        // 1. Boundary check - is it inside the maze?
        if 0 <= candidate.State.Row && candidate.State.Row < dfs.Game.Height {
            if 0 <= candidate.State.Col && candidate.State.Col < dfs.Game.Width {
                // 2. Wall check - can we actually move there?
                if !dfs.Game.Walls[candidate.State.Row][candidate.State.Col].wall {
                    neighbors = append(neighbors, candidate)
                }
            }
        }
    }

    return neighbors
}
```

**The validation process:**

1. **Generate candidates:** Up, Down, Left, Right from current position
2. **Boundary check:** Make sure we don't go outside the maze
3. **Wall check:** Make sure we don't walk through walls
4. **Return valid moves:** Only positions we can actually reach

## Helper Functions

### Checking if State is in Frontier

```go
func (dfs *DepthFirstSearch) ContainsState(node *Node) bool {
    for _, x := range dfs.Frontier {
        if x.State == node.State {
            return true
        }
    }
    return false
}
```

**Why we need this:** Avoids adding the same position to frontier multiple times

### Checking if Position is Explored

```go
func inExplored(needle Point, haystack []Point) bool {
    for _, point := range haystack {
        if point.Row == needle.Row && point.Col == needle.Col {
            return true
        }
    }
    return false
}
```

**Simple linear search:** Check if we've already visited this position

## DFS Behavior Visualization

### Example Maze

```
A . . #
# . B #
# . . #
```

### Step-by-Step Execution

**Initial state:**

```
Frontier: [A(0,0)]
Explored: []
Current: None
```

**Step 1: Process A**

```
Current: A(0,0)
Neighbors: Right(0,1)
Frontier: [Right(0,1)]
Explored: [A(0,0)]
```

**Step 2: Process Right**

```
Current: Right(0,1)
Neighbors: Right(0,2), Down(1,1)
Frontier: [Right(0,2), Down(1,1)]  // Down added last, so explored first
Explored: [A(0,0), Right(0,1)]
```

**Step 3: Process Down (LIFO)**

```
Current: Down(1,1)
Neighbors: Down(2,1)
Frontier: [Right(0,2), Down(2,1)]
Explored: [A(0,0), Right(0,1), Down(1,1)]
```

**Step 4: Process Down**

```
Current: Down(2,1)
Neighbors: Right(2,2)
Frontier: [Right(0,2), Right(2,2)]
Explored: [A(0,0), Right(0,1), Down(1,1), Down(2,1)]
```

**Step 5: Process Right**

```
Current: Right(2,2)
Neighbors: Up(1,2)=B
Frontier: [Right(0,2), B(1,2)]
Explored: [A(0,0), Right(0,1), Down(1,1), Down(2,1), Right(2,2)]
```

**Step 6: Process B (Goal!)**

```
Current: B(1,2)
GOAL FOUND!
Path: A → Right → Down → Down → Right → Up → B
```

## Why DFS Doesn't Find Shortest Path

**DFS found path:** A → Right → Down → Down → Right → Up → B (6 steps)

**Shorter path exists:** A → Right → Right → B (3 steps)

**Why DFS missed it:** The stack processed Down(1,1) before Right(0,2) because Down was added last (LIFO behavior).

**When DFS is good:** When any solution is acceptable and you want to save memory.

## Performance Analysis

### Time Complexity

**O(V + E)** where:

- V = number of positions in maze
- E = number of connections between positions

**In practice:** Visits each reachable position once

### Space Complexity

**O(V)** for:

- Frontier stack: O(V) in worst case
- Explored set: O(V) positions visited
- Node objects: O(V) total nodes created

### Memory Usage Example

**100x100 maze:**

- Maximum 10,000 positions
- Frontier typically holds 10-100 positions
- Explored grows throughout search
- Each node: ~24 bytes (position + pointer + string)

## Common Issues and Solutions

### Problem: Infinite Loops

**Symptom:** Program never terminates

**Cause:** Not checking explored positions

**Solution:**

```go
if !inExplored(neighbor.State, dfs.Game.Explored) {
    dfs.Add(neighbor)
}
```

### Problem: Stack Overflow

**Symptom:** Program crashes with stack overflow

**Cause:** Very deep mazes with recursive implementation

**Solution:** Use iterative approach with explicit stack (which we do)

### Problem: Wrong Path

**Symptom:** Path doesn't lead to goal

**Cause:** Bug in parent linking or path reconstruction

**Solution:** Carefully trace parent pointers and check reconstruction logic

## Debugging Tips

### Add Debug Output

```go
if dfs.Game.Debug {
    fmt.Printf("Exploring: (%d, %d)\n", currentNode.State.Row, currentNode.State.Col)
    fmt.Printf("Frontier size: %d\n", len(dfs.Frontier))
    fmt.Printf("Explored count: %d\n", len(dfs.Game.Explored))
}
```

### Visualize Search Progress

```go
func (dfs *DepthFirstSearch) printState() {
    fmt.Printf("Current: (%d, %d)\n", dfs.Game.CurrentNode.State.Row, dfs.Game.CurrentNode.State.Col)
    fmt.Printf("Frontier: ")
    for _, node := range dfs.Frontier {
        fmt.Printf("(%d, %d) ", node.State.Row, node.State.Col)
    }
    fmt.Println()
}
```

## Key Takeaways

### Algorithm Principles

1. **Systematic exploration:** Never miss any reachable position
2. **State tracking:** Remember what you've explored
3. **Path reconstruction:** Use parent pointers to build solution
4. **Termination conditions:** Stop when goal found or no options left

### Design Patterns

1. **Stack-based iteration:** Avoids recursion limits
2. **State validation:** Check boundaries and walls
3. **Duplicate prevention:** Track explored and frontier separately
4. **Error handling:** Graceful failure when no solution exists

### When to Use DFS

**Good for:**

- Finding any path (not necessarily shortest)
- Memory-constrained environments
- Deep solution paths
- Decision tree problems

**Not good for:**

- Finding shortest paths (use BFS)
- Very deep or infinite graphs
- When you need optimal solutions

## Next Steps

Now that you understand how DFS works completely, you're ready to learn about **Real-World Applications** and see how this same pattern applies to many different programming problems.

The key insight: **DFS is a general problem-solving pattern that you can apply whenever you need to systematically explore possibilities.**
