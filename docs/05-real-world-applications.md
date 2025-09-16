# Real-World Applications of Search Algorithms

## Taking Your Knowledge Beyond Mazes

Now that you understand how DFS works, let's see how the same concepts power real-world applications. The patterns you've learned apply to many programming challenges you'll encounter in your career.

## Web Development Applications

### 1. Website Crawling and SEO Analysis

**Problem:** A marketing company needs to analyze all pages on a website to check for broken links and SEO issues.

**Solution using DFS:**

```go
type WebCrawler struct {
    frontier   []string          // URLs to visit (stack)
    visited    map[string]bool   // URLs already crawled
    results    []PageData        // Data collected from each page
    baseURL    string           // Only crawl within this domain
}

type PageData struct {
    URL         string
    Title       string
    Links       []string
    StatusCode  int
    LoadTime    time.Duration
}

func (crawler *WebCrawler) CrawlWebsite(startURL string) {
    // Initialize with starting URL
    crawler.frontier = append(crawler.frontier, startURL)

    // DFS exploration
    for len(crawler.frontier) > 0 {
        // Pop from stack (DFS behavior)
        currentURL := crawler.frontier[len(crawler.frontier)-1]
        crawler.frontier = crawler.frontier[:len(crawler.frontier)-1]

        // Skip if already visited
        if crawler.visited[currentURL] {
            continue
        }

        fmt.Printf("Crawling: %s\n", currentURL)
        crawler.visited[currentURL] = true

        // Fetch and analyze page
        pageData, err := crawler.fetchPage(currentURL)
        if err != nil {
            fmt.Printf("Error crawling %s: %v\n", currentURL, err)
            continue
        }

        crawler.results = append(crawler.results, pageData)

        // Add internal links to frontier (go deeper)
        for _, link := range pageData.Links {
            if crawler.isInternalLink(link) && !crawler.visited[link] {
                crawler.frontier = append(crawler.frontier, link)
            }
        }
    }
}

func (crawler *WebCrawler) isInternalLink(url string) bool {
    return strings.HasPrefix(url, crawler.baseURL)
}
```

**Why DFS works here:**

- **Memory efficient** - doesn't store all URLs at once
- **Goes deep** - fully explores one section before moving to another
- **Natural behavior** - similar to how humans browse websites

### 2. API Dependency Resolution

**Problem:** A microservices system where services depend on each other. You need to start services in the correct order.

**Solution using DFS:**

```go
type ServiceManager struct {
    services     map[string]*Service
    startOrder   []string
    visited      map[string]bool
    starting     map[string]bool  // Detect circular dependencies
}

type Service struct {
    Name         string
    Dependencies []string
    Started      bool
}

func (sm *ServiceManager) StartAllServices() error {
    for serviceName := range sm.services {
        if !sm.services[serviceName].Started {
            err := sm.startService(serviceName)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (sm *ServiceManager) startService(serviceName string) error {
    // Check for circular dependency
    if sm.starting[serviceName] {
        return fmt.Errorf("circular dependency detected: %s", serviceName)
    }

    // Already started
    if sm.visited[serviceName] {
        return nil
    }

    service := sm.services[serviceName]
    sm.starting[serviceName] = true

    // Start all dependencies first (DFS goes deep)
    for _, dep := range service.Dependencies {
        err := sm.startService(dep)
        if err != nil {
            return err
        }
    }

    // Now start this service
    fmt.Printf("Starting service: %s\n", serviceName)
    service.Started = true
    sm.visited[serviceName] = true
    sm.starting[serviceName] = false
    sm.startOrder = append(sm.startOrder, serviceName)

    return nil
}
```

**Real-world example:**

```
Frontend API depends on Auth Service
Auth Service depends on Database Service
Database Service depends on Config Service

DFS ensures: Config → Database → Auth → Frontend (correct order)
```

### 3. File System Organization

**Problem:** A content management system needs to find all images in a directory structure to generate thumbnails.

**Solution using DFS:**

```go
type FileFinder struct {
    results    []string          // Found files
    extensions map[string]bool   // File types to find
}

func (ff *FileFinder) FindFiles(rootDir string, extensions []string) ([]string, error) {
    // Setup
    ff.extensions = make(map[string]bool)
    for _, ext := range extensions {
        ff.extensions[ext] = true
    }

    // DFS through directory structure
    stack := []string{rootDir}

    for len(stack) > 0 {
        // Pop current directory
        currentDir := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        fmt.Printf("Scanning: %s\n", currentDir)

        entries, err := os.ReadDir(currentDir)
        if err != nil {
            continue  // Skip inaccessible directories
        }

        for _, entry := range entries {
            fullPath := filepath.Join(currentDir, entry.Name())

            if entry.IsDir() {
                // Add subdirectory to stack (go deeper first)
                stack = append(stack, fullPath)
            } else {
                // Check if file matches our criteria
                ext := filepath.Ext(entry.Name())
                if ff.extensions[ext] {
                    ff.results = append(ff.results, fullPath)
                }
            }
        }
    }

    return ff.results, nil
}

// Usage
finder := FileFinder{}
imageFiles, err := finder.FindFiles("/Users/rob/Photos", []string{".jpg", ".png", ".gif"})
```

## Backend System Applications

### 4. Database Query Optimization

**Problem:** A complex database with many related tables. Find the optimal path to join tables for a query.

**Solution using DFS:**

```go
type QueryOptimizer struct {
    tables      map[string]*Table
    joinCosts   map[string]map[string]int  // Cost to join table A to table B
}

type Table struct {
    Name        string
    ForeignKeys map[string]string  // column -> referenced table
}

type JoinPath struct {
    Tables []string
    Cost   int
}

func (qo *QueryOptimizer) FindOptimalJoinPath(startTable, endTable string) *JoinPath {
    // DFS to find path with minimum cost
    stack := []*JoinPath{{Tables: []string{startTable}, Cost: 0}}
    visited := make(map[string]bool)
    bestPath := (*JoinPath)(nil)

    for len(stack) > 0 {
        // Pop current path
        currentPath := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        lastTable := currentPath.Tables[len(currentPath.Tables)-1]

        if visited[lastTable] {
            continue
        }
        visited[lastTable] = true

        // Found target table
        if lastTable == endTable {
            if bestPath == nil || currentPath.Cost < bestPath.Cost {
                bestPath = currentPath
            }
            continue
        }

        // Explore connected tables
        table := qo.tables[lastTable]
        for _, foreignTable := range table.ForeignKeys {
            if !visited[foreignTable] {
                newPath := &JoinPath{
                    Tables: append(currentPath.Tables, foreignTable),
                    Cost:   currentPath.Cost + qo.joinCosts[lastTable][foreignTable],
                }
                stack = append(stack, newPath)
            }
        }
    }

    return bestPath
}
```

### 5. Social Network Analysis

**Problem:** A social media platform wants to find mutual connections between users (friend suggestions).

**Solution using DFS:**

```go
type SocialNetwork struct {
    users       map[string]*User
    connections map[string][]string  // user -> list of friends
}

type User struct {
    ID       string
    Name     string
    Mutual   []string  // Mutual friends found
}

func (sn *SocialNetwork) FindMutualFriends(userA, userB string, maxDepth int) []string {
    pathsFromA := sn.findPathsFromUser(userA, maxDepth)
    pathsFromB := sn.findPathsFromUser(userB, maxDepth)

    // Find users reachable from both
    mutual := make(map[string]bool)

    for user := range pathsFromA {
        if pathsFromB[user] && user != userA && user != userB {
            mutual[user] = true
        }
    }

    var result []string
    for user := range mutual {
        result = append(result, user)
    }

    return result
}

func (sn *SocialNetwork) findPathsFromUser(startUser string, maxDepth int) map[string]bool {
    reachable := make(map[string]bool)
    stack := []struct {
        user  string
        depth int
    }{{startUser, 0}}

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if current.depth > maxDepth || reachable[current.user] {
            continue
        }

        reachable[current.user] = true

        // Add friends to stack
        for _, friend := range sn.connections[current.user] {
            if !reachable[friend] {
                stack = append(stack, struct {
                    user  string
                    depth int
                }{friend, current.depth + 1})
            }
        }
    }

    return reachable
}
```

## Game Development Applications

### 6. Game AI Pathfinding

**Problem:** An NPC in a game needs to navigate around obstacles to reach the player.

**Solution using DFS:**

```go
type GameAI struct {
    gameMap    [][]bool  // true = walkable, false = obstacle
    width      int
    height     int
}

type Position struct {
    X, Y int
}

func (ai *GameAI) FindPathToPlayer(npcPos, playerPos Position) []Position {
    stack := []struct {
        pos  Position
        path []Position
    }{{npcPos, []Position{npcPos}}}

    visited := make(map[Position]bool)

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if visited[current.pos] {
            continue
        }
        visited[current.pos] = true

        // Found player
        if current.pos == playerPos {
            return current.path
        }

        // Try all 4 directions
        directions := []Position{
            {current.pos.X + 1, current.pos.Y},     // right
            {current.pos.X - 1, current.pos.Y},     // left
            {current.pos.X, current.pos.Y + 1},     // down
            {current.pos.X, current.pos.Y - 1},     // up
        }

        for _, next := range directions {
            if ai.isValidMove(next) && !visited[next] {
                newPath := make([]Position, len(current.path))
                copy(newPath, current.path)
                newPath = append(newPath, next)

                stack = append(stack, struct {
                    pos  Position
                    path []Position
                }{next, newPath})
            }
        }
    }

    return nil  // No path found
}

func (ai *GameAI) isValidMove(pos Position) bool {
    return pos.X >= 0 && pos.X < ai.width &&
           pos.Y >= 0 && pos.Y < ai.height &&
           ai.gameMap[pos.Y][pos.X]  // Check if walkable
}
```

## Data Processing Applications

### 7. Log Analysis and Pattern Detection

**Problem:** Analyze web server logs to trace user sessions and detect suspicious behavior.

**Solution using DFS:**

```go
type LogAnalyzer struct {
    logEntries  []LogEntry
    sessions    map[string]*UserSession
}

type LogEntry struct {
    Timestamp time.Time
    UserID    string
    Action    string
    IPAddress string
    UserAgent string
}

type UserSession struct {
    UserID    string
    Actions   []LogEntry
    Suspicious bool
}

func (la *LogAnalyzer) TraceUserSessions() {
    // Group by user
    userLogs := make(map[string][]LogEntry)
    for _, entry := range la.logEntries {
        userLogs[entry.UserID] = append(userLogs[entry.UserID], entry)
    }

    // DFS through each user's action sequence
    for userID, logs := range userLogs {
        session := la.analyzeUserSession(userID, logs)
        la.sessions[userID] = session
    }
}

func (la *LogAnalyzer) analyzeUserSession(userID string, logs []LogEntry) *UserSession {
    session := &UserSession{UserID: userID}

    // DFS through action sequences
    stack := []struct {
        index   int
        pattern []string
    }{{0, []string{}}}

    suspiciousPatterns := [][]string{
        {"login", "admin_access", "delete_user"},
        {"rapid_requests", "rapid_requests", "rapid_requests"},
    }

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if current.index >= len(logs) {
            continue
        }

        entry := logs[current.index]
        newPattern := append(current.pattern, entry.Action)
        session.Actions = append(session.Actions, entry)

        // Check for suspicious patterns
        for _, suspicious := range suspiciousPatterns {
            if la.matchesPattern(newPattern, suspicious) {
                session.Suspicious = true
                fmt.Printf("Suspicious activity detected for user %s\n", userID)
            }
        }

        // Continue exploring
        stack = append(stack, struct {
            index   int
            pattern []string
        }{current.index + 1, newPattern})
    }

    return session
}

func (la *LogAnalyzer) matchesPattern(actions, pattern []string) bool {
    if len(actions) < len(pattern) {
        return false
    }

    for i := 0; i <= len(actions)-len(pattern); i++ {
        match := true
        for j := 0; j < len(pattern); j++ {
            if actions[i+j] != pattern[j] {
                match = false
                break
            }
        }
        if match {
            return true
        }
    }
    return false
}
```

## DevOps and Infrastructure

### 8. Configuration Management

**Problem:** Deploy configuration changes across a cluster of servers with dependencies.

**Solution using DFS:**

```go
type ConfigManager struct {
    servers      map[string]*Server
    dependencies map[string][]string  // server -> depends on these servers
}

type Server struct {
    Name       string
    Config     map[string]string
    Updated    bool
    Deploying  bool
}

func (cm *ConfigManager) DeployConfigurations() error {
    for serverName := range cm.servers {
        if !cm.servers[serverName].Updated {
            err := cm.deployToServer(serverName)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (cm *ConfigManager) deployToServer(serverName string) error {
    server := cm.servers[serverName]

    // Prevent circular dependencies
    if server.Deploying {
        return fmt.Errorf("circular dependency detected: %s", serverName)
    }

    if server.Updated {
        return nil
    }

    server.Deploying = true

    // Deploy to dependencies first (DFS)
    for _, dep := range cm.dependencies[serverName] {
        err := cm.deployToServer(dep)
        if err != nil {
            return err
        }
    }

    // Now deploy to this server
    fmt.Printf("Deploying configuration to %s\n", serverName)
    err := cm.performDeployment(server)
    if err != nil {
        return err
    }

    server.Updated = true
    server.Deploying = false

    return nil
}
```

## Key Patterns You Can Apply

### 1. The Search Template

This pattern works for any exploration problem:

```go
type SearchProblem struct {
    start    State
    goal     State
    frontier []State
    visited  map[State]bool
}

func (sp *SearchProblem) solve() Solution {
    sp.frontier = append(sp.frontier, sp.start)

    for len(sp.frontier) > 0 {
        current := sp.removeFromFrontier()  // DFS: from end, BFS: from beginning

        if current == sp.goal {
            return sp.buildSolution(current)
        }

        sp.visited[current] = true

        for _, neighbor := range sp.getNeighbors(current) {
            if !sp.visited[neighbor] {
                sp.addToFrontier(neighbor)
            }
        }
    }
    return nil
}
```

### 2. Dependency Resolution Pattern

```go
func resolveDependencies(item string, dependencies map[string][]string, resolved []string, resolving map[string]bool) error {
    if resolving[item] {
        return errors.New("circular dependency")
    }

    if contains(resolved, item) {
        return nil
    }

    resolving[item] = true

    for _, dep := range dependencies[item] {
        err := resolveDependencies(dep, dependencies, resolved, resolving)
        if err != nil {
            return err
        }
    }

    resolved = append(resolved, item)
    resolving[item] = false

    return nil
}
```

### 3. Tree/Graph Traversal Pattern

```go
func traverse(root *Node, visited map[*Node]bool, action func(*Node)) {
    stack := []*Node{root}

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if visited[current] {
            continue
        }

        visited[current] = true
        action(current)

        for _, child := range current.Children {
            if !visited[child] {
                stack = append(stack, child)
            }
        }
    }
}
```

## Career Applications

### Skills You've Gained

1. **Problem decomposition** - Breaking complex problems into searchable states
2. **Algorithm thinking** - Systematic exploration strategies
3. **Data structure usage** - Choosing stacks vs queues vs priority queues
4. **Performance analysis** - Understanding time/space tradeoffs

### Interview Questions You Can Handle

- "How would you crawl a website?"
- "Design a system to resolve service dependencies"
- "Find all files matching a pattern in a directory tree"
- "Detect cycles in a dependency graph"
- "Implement a recommendation system"

### Projects You Can Build

- **Web scraper** for job listings or product prices
- **File organizer** that categorizes files by type and size
- **Social network analyzer** to find connection paths
- **Build system** that compiles projects in dependency order
- **Game AI** with pathfinding and decision making

## Next Steps

You now understand how DFS applies to real-world programming challenges. The final document, **Next Steps**, will show you how to:

- Extend the maze solver with new algorithms
- Build your own projects using these patterns
- Continue learning advanced AI concepts
- Apply these skills in your career

**Remember:** The maze solver wasn't just about mazes - it was about learning a fundamental approach to systematic problem-solving that you'll use throughout your programming career.
