# Go Fundamentals for Maze Solving

## Overview

This module covers the **essential Go programming concepts** required for implementing search algorithms. You'll learn the specific language features, data structures, and patterns used throughout the maze solver project.

### Learning Objectives

By completing this module, you will:

- **Master Go's type system** and struct definitions
- **Understand slice operations** for implementing stacks and queues
- **Apply method receivers** for object-oriented programming
- **Handle errors gracefully** using Go's error patterns
- **Organize code effectively** using packages and interfaces

## Core Data Types and Structures

### Fundamental Types

Go provides several **built-in types** for different kinds of data:

```go
// Numeric types
var height int = 10          // Signed integers
var width int = 15
var distance float64 = 3.14  // Floating-point numbers

// Boolean type
var found bool = false       // true or false values

// String type
var name string = "maze"     // Text data
```

#### Key Concept: Strong Typing

Go is **strongly typed**, meaning:

- Every variable must have a declared type
- Type mismatches cause compilation errors
- This prevents many runtime bugs
- Code is more predictable and maintainable

### Custom Data Structures (Structs)

**Structs** allow you to group related data together:

```go
type Point struct {
    Row int
    Col int
}
```

#### Why Use Structs?

| Before (Separate Variables)  | After (Struct)                   |
| ---------------------------- | -------------------------------- |
| `startRow, startCol := 0, 0` | `start := Point{Row: 0, Col: 0}` |
| Track multiple variables     | Single, cohesive unit            |
| Easy to mix up parameters    | Clear, self-documenting          |
| Hard to pass around          | Easy function parameters         |

#### Working with Structs

```go
// Creating struct instances
start := Point{Row: 0, Col: 0}
goal := Point{Row: 5, Col: 3}

// Accessing struct fields
fmt.Printf("Starting position: (%d, %d)\n", start.Row, start.Col)
fmt.Printf("Goal position: (%d, %d)\n", goal.Row, goal.Col)

// Modifying struct fields
start.Row = 1
start.Col = 2
```

The **Maze struct** demonstrates **composition** - building complex types from simpler ones:

```go
type Maze struct {
    Height int           // Maze dimensions
    Width  int
    Start  Point         // Starting position (uses Point struct)
    Goal   Point         // Target position (uses Point struct)
    Walls  [][]Wall      // 2D grid of wall information
}
```

#### Benefits of Structured Data

- **Organization**: Related data stays together
- **Maintainability**: Easy to add new fields
- **Type Safety**: Prevents mixing up related values
- **Self-Documentation**: Field names explain purpose

## Slice Operations and Data Management

### Understanding Slices

**Slices** are Go's version of dynamic arrays - they can grow and shrink as needed:

```go
// Creating slices
var numbers []int                    // Empty slice
numbers = append(numbers, 1)         // Add item: [1]
numbers = append(numbers, 2, 3)      // Add multiple: [1, 2, 3]

// Accessing slice elements
first := numbers[0]                  // Gets 1 (first element)
length := len(numbers)               // Gets 3 (slice length)
```

#### Slice vs Array Comparison

| Feature         | Arrays                | Slices                   |
| --------------- | --------------------- | ------------------------ |
| **Size**        | Fixed at compile time | Dynamic, can grow/shrink |
| **Declaration** | `var arr [5]int`      | `var slice []int`        |
| **Memory**      | Stack allocated       | Heap allocated           |
| **Use Case**    | Known, fixed size     | Unknown or variable size |

### Stack Implementation with Slices

**Stacks** follow LIFO (Last In, First Out) behavior, essential for DFS:

```go
// Stack operations using slices
stack := []string{}

// Push operation (add to end)
stack = append(stack, "first")
stack = append(stack, "second")
stack = append(stack, "third")
// Result: ["first", "second", "third"]

// Pop operation (remove from end)
if len(stack) > 0 {
    top := stack[len(stack)-1]           // Get "third"
    stack = stack[:len(stack)-1]         // Remove last element
    fmt.Println("Popped:", top)
    // Stack now: ["first", "second"]
}
```

#### Why Stacks Matter for DFS

- **DFS explores deep first** - stack naturally provides this behavior
- **Backtracking** - pop returns to previous decision point
- **Memory efficient** - only stores current exploration path

### 2D Slice Structures (Grids)

**Two-dimensional slices** represent grids, essential for maze representation:

```go
// Creating a 2D grid
height, width := 3, 4
grid := make([][]bool, height)       // Create rows
for i := range grid {
    grid[i] = make([]bool, width)    // Create columns for each row
}

// Working with grid data
grid[1][2] = true                    // Set Row 1, Column 2
value := grid[0][3]                  // Get value at Row 0, Column 3

// Grid visualization:
// grid[0]: [false, false, false, false]
// grid[1]: [false, false, true,  false]
// grid[2]: [false, false, false, false]
```

#### Real-World Grid Applications

- **Game Maps**: Representing terrain, obstacles, player positions
- **Image Processing**: Pixel data manipulation
- **Spreadsheets**: Row and column data organization
- **Geographic Data**: Coordinate-based information systems

## Method Receivers and Object-Oriented Patterns

### Attaching Behavior to Data

**Method receivers** connect functions to specific types:

```go
type Calculator struct {
    result int
}

// Method with receiver - belongs to Calculator
func (c *Calculator) Add(number int) {
    c.result += number
}

func (c *Calculator) GetResult() int {
    return c.result
}

// Usage example
calc := Calculator{result: 0}
calc.Add(5)                         // calc.result becomes 5
calc.Add(3)                         // calc.result becomes 8
fmt.Println(calc.GetResult())       // Prints: 8
```

#### Value vs Pointer Receivers

| Receiver Type | Syntax                 | When to Use                               |
| ------------- | ---------------------- | ----------------------------------------- |
| **Value**     | `func (c Calculator)`  | Read-only operations, small structs       |
| **Pointer**   | `func (c *Calculator)` | Modify data, large structs, avoid copying |

### Maze Solver Method Example

```go
// Method attached to Maze type
func (m *Maze) Load(fileName string) error {
    // This function belongs to Maze
    // Can access: m.Height, m.Width, m.Start, m.Goal, etc.
    f, err := os.Open(fileName)
    if err != nil {
        return fmt.Errorf("cannot open file %s: %w", fileName, err)
    }
    defer f.Close()

    // Process file and populate maze fields...
    return nil
}

// Usage
maze := Maze{}
err := maze.Load("maze.txt")         // Calls the Load method
if err != nil {
    fmt.Println("Error:", err)
}
```

## Pointer Concepts and Memory Management

### Understanding Pointers

**Pointers** store memory addresses rather than values directly:

```go
// Regular variable
number := 42

// Pointer to that variable
pointer := &number               // & means "address of"
value := *pointer               // * means "value at address"

fmt.Println(value)              // Prints: 42
fmt.Println(pointer)            // Prints: 0x... (memory address)
```

### Pointer Applications in Data Structures

**Linked structures** use pointers to connect elements:

```go
type Node struct {
    Value  int
    Parent *Node                // Pointer to another Node
}

// Creating a linked chain
first := Node{Value: 1}
second := Node{Value: 2, Parent: &first}
third := Node{Value: 3, Parent: &second}

// Following the chain backwards
current := &third
for current != nil {
    fmt.Println("Node value:", current.Value)
    current = current.Parent    // Move to parent
}
// Output: 3, 2, 1
```

#### Why Pointers Matter for Algorithms

- **Path Reconstruction**: Trace back from goal to start
- **Memory Efficiency**: Share data instead of copying
- **Dynamic Structures**: Build flexible data relationships
- **Performance**: Avoid expensive data copying

## Error Handling Patterns

### Go's Explicit Error Approach

**Go requires explicit error checking** rather than exceptions:

```go
func openFile(name string) (*os.File, error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, fmt.Errorf("could not open %s: %w", name, err)
    }
    return file, nil
}

// Using functions that return errors
func processFile(filename string) error {
    file, err := openFile(filename)
    if err != nil {
        return err                    // Propagate error
    }
    defer file.Close()               // Always clean up

    // Process file...
    return nil                       // Success
}
```

### Error Handling Best Practices

| Pattern                  | Example                                      | Purpose       |
| ------------------------ | -------------------------------------------- | ------------- |
| **Check immediately**    | `if err != nil { return err }`               | Fail fast     |
| **Add context**          | `fmt.Errorf("processing %s: %w", name, err)` | Debugging     |
| **Resource cleanup**     | `defer file.Close()`                         | Prevent leaks |
| **Graceful degradation** | `if err != nil { useDefault() }`             | Resilience    |

## Constants and Enumerations

### Defining Algorithm Types

**Constants with `iota`** create enumerated values:

```go
const (
    DFS = iota                  // iota starts at 0, so DFS = 0
    BFS                         // Automatically becomes 1
    AStar                       // Automatically becomes 2
    Dijkstra                    // Automatically becomes 3
)
```

#### Benefits of Named Constants

- **Readability**: `searchType == DFS` vs `searchType == 0`
- **Maintainability**: Change values in one place
- **Type Safety**: Prevent invalid values
- **Self-Documentation**: Names explain purpose

### Usage in Algorithm Selection

```go
func (m *Maze) Solve(algorithmType int) error {
    switch algorithmType {
    case DFS:
        return m.solveDFS()
    case BFS:
        return m.solveBFS()
    case AStar:
        return m.solveAStar()
    default:
        return fmt.Errorf("unknown algorithm type: %d", algorithmType)
    }
}
```

## Package Organization and Code Structure

### Package Declaration and Imports

```go
// All files in same directory start with:
package main

// Import statements for external functionality
import (
    "fmt"                       // Standard library - formatting
    "os"                        // Standard library - operating system
    "errors"                    // Standard library - error handling
    "time"                      // Standard library - time operations
)
```

### Visibility Rules

**Go uses capitalization to control visibility**:

```go
// Exported (public) - starts with capital letter
type Point struct {
    Row int                     // Exported field
    Col int                     // Exported field
}

func (p Point) Distance() int { // Exported method
    return p.Row + p.Col
}

// Unexported (private) - starts with lowercase
type internalData struct {
    secret string               // Not visible outside package
}

func (i internalData) process() { // Not visible outside package
    // Internal implementation
}
```

## Key Takeaways

### Essential Go Concepts Mastered

- **Strong Typing**: Prevents runtime errors, improves code reliability
- **Struct Composition**: Build complex types from simple building blocks
- **Slice Operations**: Dynamic arrays for flexible data management
- **Method Receivers**: Attach behavior to custom types
- **Pointer Usage**: Efficient memory management and data sharing
- **Error Handling**: Explicit, predictable error management
- **Package Organization**: Clean code structure and visibility control

### Transferable Patterns

These Go concepts apply to algorithm implementation:

1. **Data Modeling**: Use structs to represent problem states
2. **Memory Management**: Use pointers for efficient data structures
3. **Error Resilience**: Handle edge cases and failures gracefully
4. **Code Organization**: Structure code for maintainability and reuse

### Preparation for Next Module

You now understand the Go foundation needed for implementing search algorithms. The next module, **Data Structures**, will show how these language features combine to create the core components of algorithmic problem-solving.

**Practice Recommendation**: Try implementing a simple stack or queue using the slice operations you've learned before proceeding to the next module.
