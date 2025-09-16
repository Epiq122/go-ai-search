# Data Structures in Maze Solving

## Understanding the Core Data Structures

This document explains the key data structures that make maze solving possible. Think of these as the **tools** your algorithm uses to organize and process information.

## The Stack: Your Algorithm's Memory

### What is a Stack?

A stack follows **LIFO** (Last In, First Out) - like a stack of plates:

- **Push:** Add a new plate to the top
- **Pop:** Remove the top plate
- **Peek:** Look at the top plate without removing it

### Stack in Daily Life

```
Email inbox (newest on top)
Browser back button (last page first)
Undo function (last action first)
Function calls (last called, first returned)
```

### Stack Implementation in Go

```go
type Stack struct {
    items []string
}

// Push - add to top
func (s *Stack) Push(item string) {
    s.items = append(s.items, item)
}

// Pop - remove from top
func (s *Stack) Pop() (string, error) {
    if len(s.items) == 0 {
        return "", errors.New("stack is empty")
    }

    // Get the last item
    item := s.items[len(s.items)-1]

    // Remove the last item
    s.items = s.items[:len(s.items)-1]

    return item, nil
}

// Peek - look at top without removing
func (s *Stack) Peek() (string, error) {
    if len(s.items) == 0 {
        return "", errors.New("stack is empty")
    }
    return s.items[len(s.items)-1], nil
}
```

### Why DFS Uses a Stack

**DFS explores as deep as possible before backtracking.** A stack naturally provides this behavior:

```
Path exploration:
Start → Right → Right → Down → DEAD END
                        ↑
                    Pop back here
                → Try Down → Down → GOAL!
```

**Stack contents during exploration:**

```
1. [Start]
2. [Start, Right1]
3. [Start, Right1, Right2]
4. [Start, Right1, Right2, Down1]
5. [Start, Right1, Right2] (popped dead end)
6. [Start, Right1, Right2, Down2]
7. [Start, Right1, Right2, Down2, Down3] (found goal!)
```

## The Frontier: Tracking What to Explore

### What is the Frontier?

The frontier represents **places you know about but haven't visited yet.**

Think of it like your travel plans:

- **Explored:** Cities you've already visited
- **Frontier:** Cities you have tickets to visit
- **Unknown:** Cities you don't know about yet

### Frontier in the Maze Solver

```go
type DepthFirstSearch struct {
    Frontier []*Node    // Stack of positions to explore
    Game     *Maze      // The maze we're solving
}
```

**Visual representation:**

```
Maze:
A . . #
# . B #
# . . #

Frontier after starting at A:
- Position(0,1) going right
```

### Frontier Operations

```go
// Add a position to explore
func (dfs *DepthFirstSearch) Add(node *Node) {
    dfs.Frontier = append(dfs.Frontier, node)
}

// Get the next position to explore (LIFO for DFS)
func (dfs *DepthFirstSearch) Remove() (*Node, error) {
    if len(dfs.Frontier) == 0 {
        return nil, errors.New("no more positions to explore")
    }

    // Take from end (stack behavior)
    node := dfs.Frontier[len(dfs.Frontier)-1]
    dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1]

    return node, nil
}

// Check if we still have places to explore
func (dfs *DepthFirstSearch) Empty() bool {
    return len(dfs.Frontier) == 0
}
```

## Node: Representing Search States

### What is a Node?

A Node represents **one step in your search** with memory of how you got there.

```go
type Node struct {
    State  Point    // Where you are now
    Parent *Node    // Where you came from
    Action string   // How you got here ("up", "down", "left", "right")
}
```

### Why Nodes Track Parents

**To remember the path back to the start:**

```
Node chain:
Goal ← Down ← Right ← Right ← Start
 B  ←  .   ←   .   ←   .   ←  A

When you reach the goal, follow Parent pointers backwards
to reconstruct the complete path.
```

### Building the Solution Path

```go
func buildSolution(goalNode *Node) []string {
    var path []string
    current := goalNode

    // Follow parent links backwards
    for current.Parent != nil {
        path = append(path, current.Action)
        current = current.Parent
    }

    // Reverse to get start-to-goal path
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }

    return path
}
```

## The Explored Set: Avoiding Cycles

### Why We Need an Explored Set

**Without tracking visited positions, you'd go in circles:**

```
A → Right → Down → Left → Up → Right → Down → Left...
```

### Implementation

```go
// Track positions we've already visited
explored := make(map[Point]bool)

// Before exploring a position
if explored[position] {
    continue    // Skip - already been here
}

// After exploring
explored[position] = true
```

### Checking for Visited Positions

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

## Putting It All Together: The Search Loop

### The Complete Data Flow

```go
func (dfs *DepthFirstSearch) Solve() {
    // 1. Initialize frontier with starting position
    start := &Node{State: dfs.Game.Start}
    dfs.Add(start)

    // 2. Keep exploring until frontier is empty
    for !dfs.Empty() {
        // 3. Get next position to explore
        current, err := dfs.Remove()
        if err != nil {
            break
        }

        // 4. Check if we reached the goal
        if current.State == dfs.Game.Goal {
            return buildSolution(current)
        }

        // 5. Mark as explored
        dfs.Game.Explored = append(dfs.Game.Explored, current.State)

        // 6. Add neighbors to frontier
        for _, neighbor := range dfs.Neighbors(current) {
            if !alreadyExplored(neighbor) && !inFrontier(neighbor) {
                dfs.Add(neighbor)
            }
        }
    }
}
```

### Data Structure Interaction

```
1. Frontier (Stack)     → Next position to explore
2. Current Node         → Where we are now
3. Explored Set         → Where we've been
4. Neighbor Generation  → Where we can go next
5. Goal Check          → Are we done?
6. Path Reconstruction → How did we get here?
```

## Visualizing the Data Structures

### Step-by-Step Example

**Initial State:**

```
Maze:
A . B
# # .

Frontier: [A(0,0)]
Explored: []
```

**Step 1: Explore A**

```
Current: A(0,0)
Neighbors: Right(0,1)

Frontier: [Right(0,1)]
Explored: [A(0,0)]
```

**Step 2: Explore Right**

```
Current: Right(0,1)
Neighbors: Right(0,2)=B, Down(1,1)=wall

Frontier: [B(0,2)]
Explored: [A(0,0), Right(0,1)]
```

**Step 3: Explore B (Goal!)**

```
Current: B(0,2)
Goal reached!

Path reconstruction:
B(0,2) → Parent: Right(0,1) → Parent: A(0,0) → Parent: nil
Actions: [right, right]
```

## Memory and Performance

### Space Complexity

```go
Frontier:  O(V) in worst case (V = number of positions)
Explored:  O(V) positions visited
Nodes:     O(V) node objects created
Total:     O(V) space complexity
```

### Why This Matters

**In a 100x100 maze:**

- Maximum 10,000 positions
- Frontier might hold 100-1000 positions at once
- Explored set grows as you search
- Each node is small (just coordinates and a pointer)

## Data Structure Variations

### BFS vs DFS: Just Change the Data Structure

**DFS uses Stack (LIFO):**

```go
// Remove from end
node := frontier[len(frontier)-1]
frontier = frontier[:len(frontier)-1]
```

**BFS uses Queue (FIFO):**

```go
// Remove from beginning
node := frontier[0]
frontier = frontier[1:]
```

**Same algorithm, different exploration order!**

### Priority Queue for A\*

```go
type PriorityQueue []*Node

// Nodes sorted by cost + heuristic
func (pq *PriorityQueue) Pop() *Node {
    // Return node with lowest cost
}
```

## Common Pitfalls and Solutions

### Problem: Infinite Loops

**Cause:** Not checking if positions are already explored

**Solution:**

```go
if !inExplored(neighbor.State, dfs.Game.Explored) {
    dfs.Add(neighbor)
}
```

### Problem: Memory Growth

**Cause:** Keeping all nodes in memory forever

**Solution:** Only keep essential data, clean up unused nodes

### Problem: Deep Recursion

**Cause:** Using recursive function calls instead of explicit stack

**Solution:** Use iterative approach with explicit data structures

## Real-World Applications

### Web Crawling

```go
frontier := []string{"https://example.com"}
visited := make(map[string]bool)

for len(frontier) > 0 {
    url := frontier[len(frontier)-1]    // Stack behavior
    frontier = frontier[:len(frontier)-1]

    if visited[url] {
        continue
    }

    links := crawlPage(url)
    frontier = append(frontier, links...)
    visited[url] = true
}
```

### File System Search

```go
frontier := []string{"/Users/rob/Documents"}
results := []string{}

for len(frontier) > 0 {
    dir := frontier[len(frontier)-1]    // Stack behavior
    frontier = frontier[:len(frontier)-1]

    files := listDirectory(dir)
    for _, file := range files {
        if isDirectory(file) {
            frontier = append(frontier, file)
        } else if matches(file, "*.txt") {
            results = append(results, file)
        }
    }
}
```

## Key Takeaways

### Essential Concepts

1. **Stack = LIFO** - Perfect for DFS exploration
2. **Frontier = Known but unexplored** - Your to-do list
3. **Nodes = States with history** - Remember how you got somewhere
4. **Explored = Visited positions** - Avoid going in circles

### Design Principles

1. **Separation of concerns** - Each data structure has one job
2. **Simple operations** - Push, pop, check membership
3. **Memory efficiency** - Only store what you need
4. **Flexibility** - Same pattern works for different algorithms

### Next Steps

Now that you understand how data structures work in search problems, you're ready to see the complete **DFS Algorithm** implementation and understand how all these pieces work together.
