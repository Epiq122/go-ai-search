# Go Fundamentals for Maze Solving

## Essential Go Concepts You Need

This document covers the Go programming concepts used in the maze solver. If you're new to Go, focus on understanding these building blocks.

## Data Types and Structures

### Basic Types

```go
// Simple data types
var height int = 10        // Whole numbers
var width int = 15
var found bool = false     // true or false
var name string = "maze"   // Text
```

**Key concept:** Go is **strongly typed** - you must declare what type of data each variable holds.

### Custom Types (Structs)

```go
type Point struct {
    Row int
    Col int
}
```

**What this does:** Creates a new data type that groups related information together.

**Real-world analogy:** Like creating a "Contact" that has both a name and phone number, a Point has both a row and column.

**Why we use it:** Instead of tracking row and column separately, we can pass around one Point object.

```go
// Using the Point struct
start := Point{Row: 0, Col: 0}
goal := Point{Row: 5, Col: 3}

// Accessing the data
fmt.Println("Start row:", start.Row)
fmt.Println("Start column:", start.Col)
```

### Complex Structures

```go
type Maze struct {
    Height int           // How tall the maze is
    Width  int           // How wide the maze is
    Start  Point         // Starting position
    Goal   Point         // Target position
    Walls  [][]Wall      // 2D grid of wall information
}
```

**What this does:** Creates a container that holds all maze-related information in one place.

**Real-world analogy:** Like a "User Profile" that contains name, email, address, and preferences all together.

## Slices (Dynamic Arrays)

### Basic Slice Operations

```go
// Creating slices
var numbers []int                    // Empty slice
numbers = append(numbers, 1)         // Add item: [1]
numbers = append(numbers, 2, 3)      // Add multiple: [1, 2, 3]

// Accessing items
first := numbers[0]                  // Gets 1
length := len(numbers)               // Gets 3
```

### Stack Behavior with Slices

```go
// Using slice as a stack (Last In, First Out)
stack := []string{}

// Push (add to end)
stack = append(stack, "first")
stack = append(stack, "second")
stack = append(stack, "third")
// Now stack is: ["first", "second", "third"]

// Pop (remove from end)
if len(stack) > 0 {
    last := stack[len(stack)-1]          // Get "third"
    stack = stack[:len(stack)-1]         // Remove it
    // Now stack is: ["first", "second"]
}
```

**Why this matters:** The DFS algorithm uses this exact pattern to explore the maze.

### 2D Slices (Grid)

```go
// Creating a 2D grid
height, width := 3, 4
grid := make([][]bool, height)       // Create rows
for i := range grid {
    grid[i] = make([]bool, width)    // Create columns for each row
}

// Accessing grid items
grid[1][2] = true                    // Row 1, Column 2
value := grid[0][3]                  // Get value at Row 0, Column 3
```

**Real-world analogy:** Like a spreadsheet with rows and columns, or a chessboard with squares.

## Methods and Receivers

### Attaching Functions to Types

```go
type Calculator struct {
    result int
}

// Method with receiver - belongs to Calculator
func (c *Calculator) Add(number int) {
    c.result += number
}

// Using the method
calc := Calculator{result: 0}
calc.Add(5)                         // calc.result is now 5
calc.Add(3)                         // calc.result is now 8
```

**What this does:** Connects functions to specific data types.

**Why it's useful:** Groups related functionality with the data it operates on.

### In the Maze Solver

```go
// Method attached to Maze type
func (m *Maze) Load(fileName string) error {
    // This function belongs to Maze
    // It can access and modify m.Height, m.Width, etc.
}

// Using it
maze := Maze{}
err := maze.Load("maze.txt")         // Calls the Load method
```

## Pointers

### What Pointers Are

```go
// Regular variable
number := 42

// Pointer to that variable
pointer := &number               // & means "address of"
value := *pointer               // * means "value at address"

fmt.Println(value)              // Prints 42
```

**Simple explanation:** A pointer is like a sticky note that tells you where to find something, rather than being the thing itself.

### Why We Use Pointers

```go
type Node struct {
    Value  int
    Parent *Node                // Pointer to another Node
}

// Creating a chain
first := Node{Value: 1}
second := Node{Value: 2, Parent: &first}
third := Node{Value: 3, Parent: &second}

// Following the chain backwards
current := &third
for current != nil {
    fmt.Println(current.Value)
    current = current.Parent    // Move to parent
}
// Prints: 3, 2, 1
```

**Why this matters:** This is exactly how the maze solver remembers the path it took to reach the goal.

## Error Handling

### Go's Error Pattern

```go
func openFile(name string) (*File, error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, fmt.Errorf("could not open %s: %w", name, err)
    }
    return file, nil
}

// Using functions that return errors
file, err := openFile("maze.txt")
if err != nil {
    fmt.Println("Error:", err)
    return
}
defer file.Close()              // Always clean up
```

**Key points:**

- Functions return both a result and an error
- Always check if `err != nil`
- `defer` ensures cleanup happens even if something goes wrong

## Constants and Enums

### Defining Constants

```go
const (
    DFS = iota                  // iota starts at 0
    BFS                         // Automatically becomes 1
    AStar                       // Automatically becomes 2
)
```

**What `iota` does:** Automatically assigns increasing numbers to related constants.

**Why use this:** Instead of magic numbers in your code, you have readable names.

```go
// Instead of this:
if searchType == 0 {            // What does 0 mean?

// You write this:
if searchType == DFS {          // Clear and readable
```

## Package Organization

### File Structure

```go
// All files in the same directory start with:
package main

// Import statements
import (
    "fmt"                       // Standard library
    "os"                        // Standard library
    "errors"                    // Standard library
)
```

**Key rule:** All `.go` files in the same folder must have the same package name.

### Visibility Rules

```go
// Exported (visible outside package) - starts with capital letter
type Point struct {
    Row int                     // Exported
    Col int                     // Exported
}

// Unexported (private to package) - starts with lowercase
type internalData struct {
    secret string               // Not visible outside package
}
```

## Common Patterns Used in Maze Solver

### Initialization Pattern

```go
func main() {
    var maze Maze               // Create empty maze

    err := maze.Load("file.txt") // Initialize with data
    if err != nil {
        fmt.Println(err)
        os.Exit(1)              // Exit with error code
    }

    // Use the initialized maze
}
```

### Loop Until Done Pattern

```go
for {                           // Infinite loop
    if shouldStop() {
        break                   // Exit the loop
    }

    // Do work
    processNextItem()
}
```

### Safe Slice Access Pattern

```go
if len(slice) > 0 {
    item := slice[len(slice)-1] // Get last item
    slice = slice[:len(slice)-1] // Remove last item
} else {
    // Handle empty slice
}
```

## Putting It All Together

Here's how these concepts work together in a simple example:

```go
// Define data structure
type Stack struct {
    items []string
}

// Add method to push items
func (s *Stack) Push(item string) {
    s.items = append(s.items, item)
}

// Add method to pop items
func (s *Stack) Pop() (string, error) {
    if len(s.items) == 0 {
        return "", errors.New("stack is empty")
    }

    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, nil
}

// Use the stack
func main() {
    stack := Stack{}

    stack.Push("first")
    stack.Push("second")

    item, err := stack.Pop()
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Popped:", item)    // Prints "second"
    }
}
```

## Next Steps

Now that you understand these Go fundamentals, you're ready to learn about **Data Structures** and how they're used in the maze solver.

The key takeaway: **Go's simplicity makes it perfect for learning algorithms** because you spend time understanding the logic, not fighting with complex syntax.
