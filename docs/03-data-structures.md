# Data Structures in Algorithmic Problem Solving

## Chapter Overview

**Learning Objectives:**

- Master fundamental data structures used in search algorithms
- Understand the relationship between data structures and algorithm behavior
- Implement efficient data management for maze solving
- Apply data structure concepts to real-world problems

**Prerequisites:** Basic Go programming (structs, slices, pointers, methods)

## Chapter 1: The Stack - Foundation of Depth-First Search

### 1.1 Stack Fundamentals

**Definition**: A stack is a **Last-In-First-Out (LIFO)** data structure supporting three primary operations:

| Operation | Description                      | Time Complexity |
| --------- | -------------------------------- | --------------- |
| **Push**  | Add element to top               | O(1)            |
| **Pop**   | Remove top element               | O(1)            |
| **Peek**  | View top element without removal | O(1)            |

#### Real-World Stack Examples

| Context             | Stack Behavior                       | LIFO Demonstration           |
| ------------------- | ------------------------------------ | ---------------------------- |
| **Email Inbox**     | Newest emails appear first           | Latest received, first read  |
| **Browser History** | Back button returns to previous page | Last visited, first returned |
| **Function Calls**  | Most recent call executes first      | Last called, first completed |
| **Undo Operations** | Most recent action reversed first    | Last action, first undone    |

### 1.2 Stack Implementation in Go

```go
type Stack struct {
    items []Point    // Dynamic array for storage
    size  int        // Current number of elements
}

// Push - Add element to top of stack
func (s *Stack) Push(item Point) {
    s.items = append(s.items, item)
    s.size++
}

// Pop - Remove and return top element
func (s *Stack) Pop() (Point, error) {
    if s.IsEmpty() {
        return Point{}, errors.New("stack underflow: cannot pop from empty stack")
    }

    // Retrieve top element
    top := s.items[s.size-1]

    // Remove top element
    s.items = s.items[:s.size-1]
    s.size--

    return top, nil
}

// Peek - View top element without removal
func (s *Stack) Peek() (Point, error) {
    if s.IsEmpty() {
        return Point{}, errors.New("stack empty: no element to peek")
    }
    return s.items[s.size-1], nil
}

// IsEmpty - Check if stack contains no elements
func (s *Stack) IsEmpty() bool {
    return s.size == 0
}

// Size - Return current number of elements
func (s *Stack) Size() int {
    return s.size
}
```

### 1.3 Stack Usage in Depth-First Search

**DFS Exploration Pattern:**

1. **Push** starting position onto stack
2. **Pop** next position to explore
3. **Push** all unvisited neighbors
4. Repeat until goal found or stack empty

#### DFS Stack Trace Example

```
Maze Grid:
[S][.][.]
[#][#][.]
[.][.][G]

Stack Evolution:
Step 1: [S(0,0)]                    // Start position
Step 2: [Right(0,1)]                // Explore right from start
Step 3: [Right(0,1), Right(0,2)]    // Continue right
Step 4: [Right(0,1), Down(1,2)]     // Move down from (0,2)
Step 5: [Right(0,1), Down(2,2)]     // Continue down
Step 6: [Right(0,1)]                // Found goal at (2,2)!
```

## Chapter 2: The Frontier - Managing Search Boundaries

### 2.1 Frontier Concept

**Definition**: The frontier represents the **boundary between explored and unexplored regions** of the search space.

#### Frontier Analogies

| Domain              | Frontier Representation                               |
| ------------------- | ----------------------------------------------------- |
| **Geography**       | Border between mapped and unmapped territory          |
| **Travel Planning** | Destinations you have tickets for but haven't visited |
| **Research**        | Known problems you haven't investigated yet           |
| **Game Strategy**   | Moves you've identified but haven't executed          |

### 2.2 Frontier Implementation

```go
type SearchFrontier struct {
    nodes    []*Node    // Collection of unexplored nodes
    strategy string     // "stack" for DFS, "queue" for BFS
}

// Add - Insert node into frontier
func (f *SearchFrontier) Add(node *Node) {
    f.nodes = append(f.nodes, node)
}

// Remove - Extract next node based on strategy
func (f *SearchFrontier) Remove() (*Node, error) {
    if f.IsEmpty() {
        return nil, errors.New("frontier exhausted: no more nodes to explore")
    }

    var node *Node
    switch f.strategy {
    case "stack":    // LIFO - Depth-First Search
        node = f.nodes[len(f.nodes)-1]
        f.nodes = f.nodes[:len(f.nodes)-1]
    case "queue":    // FIFO - Breadth-First Search
        node = f.nodes[0]
        f.nodes = f.nodes[1:]
    default:
        return nil, errors.New("unknown frontier strategy")
    }

    return node, nil
}

// IsEmpty - Check if frontier contains any nodes
func (f *SearchFrontier) IsEmpty() bool {
    return len(f.nodes) == 0
}

// Contains - Verify if state already exists in frontier
func (f *SearchFrontier) Contains(state Point) bool {
    for _, node := range f.nodes {
        if node.State.Row == state.Row && node.State.Col == state.Col {
            return true
        }
    }
    return false
}
```

### 2.3 Frontier State Management

#### Frontier Lifecycle

| Phase              | Description                  | Frontier Content      |
| ------------------ | ---------------------------- | --------------------- |
| **Initialization** | Add starting state           | [Start]               |
| **Expansion**      | Add neighboring states       | [Start, Neighbors...] |
| **Exploration**    | Remove next state to explore | [Remaining neighbors] |
| **Termination**    | Goal found or frontier empty | [] or [Goal found]    |

## Chapter 3: Nodes - Representing Search States

### 3.1 Node Structure and Purpose

**Node Definition**: A node encapsulates a **search state** along with metadata for path reconstruction and decision making.

```go
type Node struct {
    State  Point    // Current position in search space
    Parent *Node    // Reference to predecessor node
    Action string   // Action that led to this state
    Cost   int      // Path cost from start to this node
    Depth  int      // Distance from start node
}

// Constructor for root node
func NewRootNode(startState Point) *Node {
    return &Node{
        State:  startState,
        Parent: nil,        // Root has no parent
        Action: "",         // No action to reach start
        Cost:   0,          // Zero cost at start
        Depth:  0,          // Zero depth at start
    }
}

// Constructor for child node
func NewChildNode(state Point, parent *Node, action string) *Node {
    return &Node{
        State:  state,
        Parent: parent,
        Action: action,
        Cost:   parent.Cost + 1,    // Increment path cost
        Depth:  parent.Depth + 1,   // Increment depth
    }
}
```

### 3.2 Path Reconstruction Using Parent Pointers

**Path Reconstruction Algorithm:**

1. Start from goal node
2. Follow parent pointers backward
3. Collect actions in reverse order
4. Reverse action sequence for start-to-goal path

```go
// BuildSolutionPath - Reconstruct path from start to goal
func (n *Node) BuildSolutionPath() []string {
    var actions []string
    current := n

    // Traverse backward through parent links
    for current.Parent != nil {
        actions = append(actions, current.Action)
        current = current.Parent
    }

    // Reverse to get start-to-goal sequence
    for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
        actions[i], actions[j] = actions[j], actions[i]
    }

    return actions
}

// GetPathLength - Calculate total path length
func (n *Node) GetPathLength() int {
    return n.Depth
}

// GetPathCost - Return accumulated path cost
func (n *Node) GetPathCost() int {
    return n.Cost
}
```

### 3.3 Node Relationship Visualization

```
Path Reconstruction Example:

Goal Node ← "down" ← Middle Node ← "right" ← Start Node
   (2,2)              (1,2)                  (1,1)

Backward Traversal: ["down", "right"]
Forward Path: ["right", "down"]

Memory Structure:
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Start Node  │◄───┤Middle Node  │◄───┤ Goal Node   │
│ State:(1,1) │    │State:(1,2)  │    │State:(2,2)  │
│ Parent:nil  │    │Parent:Start │    │Parent:Middle│
│ Action:""   │    │Action:"right"│    │Action:"down"│
└─────────────┘    └─────────────┘    └─────────────┘
```

## Chapter 4: Explored Set - Cycle Prevention

### 4.1 The Explored Set Concept

**Purpose**: Track visited states to prevent infinite loops and redundant exploration.

#### Why Cycle Prevention Matters

```
Without Explored Set:
Start → Right → Down → Left → Up → Right → Down → Left...
(Infinite loop)

With Explored Set:
Start → Right → Down ← (Left blocked - already explored)
Continue to new unexplored areas
```

### 4.2 Explored Set Implementation Options

#### Option 1: Map-Based (Efficient)

```go
type ExploredSet struct {
    visited map[Point]bool
}

func NewExploredSet() *ExploredSet {
    return &ExploredSet{
        visited: make(map[Point]bool),
    }
}

func (e *ExploredSet) Add(state Point) {
    e.visited[state] = true
}

func (e *ExploredSet) Contains(state Point) bool {
    return e.visited[state]
}

func (e *ExploredSet) Size() int {
    return len(e.visited)
}
```

#### Option 2: Slice-Based (Simple)

```go
type ExploredSlice struct {
    states []Point
}

func (e *ExploredSlice) Add(state Point) {
    e.states = append(e.states, state)
}

func (e *ExploredSlice) Contains(state Point) bool {
    for _, explored := range e.states {
        if explored.Row == state.Row && explored.Col == state.Col {
            return true
        }
    }
    return false
}
```

#### Performance Comparison

| Implementation  | Add Operation  | Contains Check | Memory Usage    |
| --------------- | -------------- | -------------- | --------------- |
| **Map-based**   | O(1) average   | O(1) average   | Higher overhead |
| **Slice-based** | O(1) amortized | O(n) linear    | Lower overhead  |

### 4.3 Explored Set Integration

```go
// Complete search with cycle prevention
func (dfs *DepthFirstSearch) Solve() ([]string, error) {
    frontier := NewSearchFrontier("stack")
    explored := NewExploredSet()

    // Initialize with start state
    startNode := NewRootNode(dfs.Game.Start)
    frontier.Add(startNode)

    for !frontier.IsEmpty() {
        current, err := frontier.Remove()
        if err != nil {
            return nil, err
        }

        // Check for goal before marking as explored
        if current.State == dfs.Game.Goal {
            return current.BuildSolutionPath(), nil
        }

        // Mark current state as explored
        explored.Add(current.State)

        // Expand neighbors
        for _, neighbor := range dfs.GetNeighbors(current) {
            neighborState := neighbor.State

            // Skip if already explored or in frontier
            if explored.Contains(neighborState) || frontier.Contains(neighborState) {
                continue
            }

            frontier.Add(neighbor)
        }
    }

    return nil, errors.New("no solution found")
}
```

## Chapter 5: Algorithm Data Flow Integration

### 5.1 Complete Search Process

#### Data Structure Orchestration

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Frontier   │    │   Current   │    │  Explored   │
│   (Stack)   │    │    Node     │    │    Set      │
├─────────────┤    ├─────────────┤    ├─────────────┤
│   Remove    │───►│   Process   │───►│     Add     │
│    Next     │    │   State     │    │   State     │
└─────────────┘    └─────────────┘    └─────────────┘
       ▲                  │
       │                  ▼
┌─────────────┐    ┌─────────────┐
│  Generate   │    │    Goal     │
│ Neighbors   │    │    Check    │
└─────────────┘    └─────────────┘
```

### 5.2 Algorithm State Transitions

| Step  | Frontier | Current | Explored      | Action           |
| ----- | -------- | ------- | ------------- | ---------------- |
| **0** | [Start]  | null    | {}            | Initialize       |
| **1** | []       | Start   | {Start}       | Expand neighbors |
| **2** | [N1,N2]  | Start   | {Start}       | Add to frontier  |
| **3** | [N2]     | N1      | {Start,N1}    | Process N1       |
| **4** | [N2,N3]  | N1      | {Start,N1}    | Add N1 neighbors |
| **5** | [N3]     | N2      | {Start,N1,N2} | Process N2       |

### 5.3 Memory and Performance Analysis

#### Space Complexity Analysis

| Component        | Worst Case | Typical Case | Notes                       |
| ---------------- | ---------- | ------------ | --------------------------- |
| **Frontier**     | O(b^d)     | O(bd)        | b=branching factor, d=depth |
| **Explored**     | O(V)       | O(V)         | V = total states visited    |
| **Node Storage** | O(V)       | O(V)         | One node per state          |
| **Total**        | O(V + b^d) | O(V)         | Dominated by explored set   |

#### Time Complexity Analysis

| Operation               | Complexity | Frequency | Total Impact |
| ----------------------- | ---------- | --------- | ------------ |
| **Frontier Add**        | O(1)       | O(V)      | O(V)         |
| **Frontier Remove**     | O(1)       | O(V)      | O(V)         |
| **Explored Check**      | O(1)\*     | O(E)      | O(E)         |
| **Neighbor Generation** | O(b)       | O(V)      | O(bV)        |

\*Using map-based implementation

## Chapter 6: Algorithm Variations and Data Structure Impact

### 6.1 Search Strategy Comparison

#### Stack vs Queue Impact

| Algorithm      | Data Structure | Exploration Pattern | Path Quality         |
| -------------- | -------------- | ------------------- | -------------------- |
| **DFS**        | Stack (LIFO)   | Deep first          | May find suboptimal  |
| **BFS**        | Queue (FIFO)   | Level by level      | Finds shortest path  |
| **Best-First** | Priority Queue | Lowest cost first   | Depends on heuristic |

### 6.2 Data Structure Substitution

```go
// Generic search framework
type SearchAlgorithm struct {
    frontier  FrontierInterface
    explored  ExploredInterface
    generator NeighborGenerator
}

// FrontierInterface allows different strategies
type FrontierInterface interface {
    Add(node *Node)
    Remove() (*Node, error)
    IsEmpty() bool
    Contains(state Point) bool
}

// DFS Implementation
func NewDepthFirstSearch() *SearchAlgorithm {
    return &SearchAlgorithm{
        frontier:  NewStackFrontier(),     // LIFO behavior
        explored:  NewMapExploredSet(),
        generator: NewGridNeighborGenerator(),
    }
}

// BFS Implementation
func NewBreadthFirstSearch() *SearchAlgorithm {
    return &SearchAlgorithm{
        frontier:  NewQueueFrontier(),     // FIFO behavior
        explored:  NewMapExploredSet(),
        generator: NewGridNeighborGenerator(),
    }
}
```

## Chapter 7: Real-World Applications

### 7.1 Web Crawling

```go
type WebCrawler struct {
    frontier []string           // URLs to visit
    visited  map[string]bool    // Already crawled URLs
}

func (w *WebCrawler) Crawl(startURL string) {
    w.frontier = append(w.frontier, startURL)

    for len(w.frontier) > 0 {
        // Stack behavior - depth-first crawling
        url := w.frontier[len(w.frontier)-1]
        w.frontier = w.frontier[:len(w.frontier)-1]

        if w.visited[url] {
            continue
        }

        links := w.extractLinks(url)
        w.frontier = append(w.frontier, links...)
        w.visited[url] = true
    }
}
```

### 7.2 File System Navigation

```go
type FileSystemSearch struct {
    directories []string        // Stack of directories to explore
    processed   map[string]bool // Already processed directories
    results     []string        // Found files matching criteria
}

func (fs *FileSystemSearch) FindFiles(root string, pattern string) []string {
    fs.directories = append(fs.directories, root)

    for len(fs.directories) > 0 {
        // Stack behavior - depth-first directory traversal
        dir := fs.directories[len(fs.directories)-1]
        fs.directories = fs.directories[:len(fs.directories)-1]

        entries := fs.listDirectory(dir)
        for _, entry := range entries {
            if fs.isDirectory(entry) {
                fs.directories = append(fs.directories, entry)
            } else if fs.matchesPattern(entry, pattern) {
                fs.results = append(fs.results, entry)
            }
        }

        fs.processed[dir] = true
    }

    return fs.results
}
```

## Key Takeaways

### Essential Data Structure Principles

| Principle                  | Description                                    | Impact              |
| -------------------------- | ---------------------------------------------- | ------------------- |
| **Separation of Concerns** | Each structure has specific responsibility     | Maintainable code   |
| **Interface Consistency**  | Uniform operations across implementations      | Flexible algorithms |
| **Memory Efficiency**      | Store only necessary information               | Scalable solutions  |
| **Performance Awareness**  | Choose structures based on operation frequency | Efficient execution |

### Algorithm Design Patterns

1. **Frontier Management**: Control exploration order through data structure choice
2. **State Tracking**: Prevent cycles and redundant work through explored sets
3. **Path Reconstruction**: Maintain parent links for solution extraction
4. **Incremental Processing**: Process one state at a time for memory efficiency

### Transferable Concepts

These data structure patterns apply beyond maze solving:

- **Graph Traversal**: Social networks, dependency resolution, route planning
- **Tree Processing**: File systems, organizational hierarchies, decision trees
- **State Space Search**: Game AI, planning problems, optimization
- **Resource Management**: Task scheduling, memory allocation, queue processing

### Preparation for Implementation

Understanding these data structures prepares you for the complete **DFS Algorithm Implementation** in the next module, where you'll see how all components integrate into a working solution.

**Study Recommendation**: Practice implementing stack and queue operations manually before proceeding to see how they orchestrate complex algorithm behavior.

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
