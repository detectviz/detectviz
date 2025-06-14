package composition

import (
	"fmt"
	"sort"

	"detectviz/pkg/platform/contracts"
)

// DependencyResolver implements plugin dependency resolution with topological sort.
// zh: DependencyResolver 實作插件依賴解析與拓撲排序。
type DependencyResolver struct {
	graph map[string][]string // dependency graph: node -> [dependencies]
	nodes map[string]contracts.ComponentInfo
}

// NewDependencyResolver creates a new dependency resolver.
// zh: NewDependencyResolver 建立新的依賴解析器。
func NewDependencyResolver() contracts.DependencyResolver {
	return &DependencyResolver{
		graph: make(map[string][]string),
		nodes: make(map[string]contracts.ComponentInfo),
	}
}

// ResolveDependencies resolves component dependencies and returns them in start order.
// zh: ResolveDependencies 解析組件依賴關係並回傳啟動順序。
func (dr *DependencyResolver) ResolveDependencies(components []contracts.ComponentInfo) ([]contracts.ComponentInfo, error) {
	// Reset the resolver state
	dr.graph = make(map[string][]string)
	dr.nodes = make(map[string]contracts.ComponentInfo)

	// Build the dependency graph
	for _, comp := range components {
		dr.nodes[comp.Name] = comp
		dr.graph[comp.Name] = comp.Dependencies
	}

	// Validate dependencies
	if err := dr.ValidateDependencies(components); err != nil {
		return nil, err
	}

	// Perform topological sort
	order, err := dr.topologicalSort()
	if err != nil {
		return nil, err
	}

	// Convert ordered names back to ComponentInfo
	orderedComponents := make([]contracts.ComponentInfo, 0, len(order))
	for _, name := range order {
		orderedComponents = append(orderedComponents, dr.nodes[name])
	}

	return orderedComponents, nil
}

// ValidateDependencies validates that all dependencies exist and there are no cycles.
// zh: ValidateDependencies 驗證所有依賴關係存在且無循環依賴。
func (dr *DependencyResolver) ValidateDependencies(components []contracts.ComponentInfo) error {
	// Check that all dependencies exist
	for _, comp := range components {
		for _, dep := range comp.Dependencies {
			if _, exists := dr.nodes[dep]; !exists {
				return fmt.Errorf("component %s depends on non-existent component %s", comp.Name, dep)
			}
		}
	}

	// Check for cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range dr.nodes {
		if !visited[name] {
			if dr.hasCycleDFS(name, visited, recStack) {
				return fmt.Errorf("circular dependency detected involving component %s", name)
			}
		}
	}

	return nil
}

// GetDependencyGraph returns the internal dependency graph.
// zh: GetDependencyGraph 回傳內部依賴關係圖。
func (dr *DependencyResolver) GetDependencyGraph() contracts.DependencyGraph {
	return &SimpleDependencyGraph{
		nodes: dr.nodes,
		edges: dr.graph,
	}
}

// topologicalSort performs Kahn's algorithm for topological sorting.
// zh: topologicalSort 執行 Kahn 演算法進行拓撲排序。
func (dr *DependencyResolver) topologicalSort() ([]string, error) {
	// Calculate in-degrees
	inDegree := make(map[string]int)
	for name := range dr.nodes {
		inDegree[name] = 0
	}

	for _, deps := range dr.graph {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	// Initialize queue with nodes having no dependencies
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// Sort the initial queue for deterministic ordering
	sort.Strings(queue)

	result := make([]string, 0, len(dr.nodes))

	// Process queue
	for len(queue) > 0 {
		// Remove first element
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree of dependent nodes
		dependents := dr.getDependents(current)
		newZeroDegreeNodes := make([]string, 0)

		for _, dependent := range dependents {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				newZeroDegreeNodes = append(newZeroDegreeNodes, dependent)
			}
		}

		// Sort and add new zero-degree nodes to maintain deterministic order
		sort.Strings(newZeroDegreeNodes)
		queue = append(queue, newZeroDegreeNodes...)
	}

	// Check if all nodes were processed (no cycles)
	if len(result) != len(dr.nodes) {
		return nil, fmt.Errorf("circular dependency detected in component graph")
	}

	return result, nil
}

// getDependents returns all components that depend on the given component.
// zh: getDependents 回傳所有依賴於指定組件的組件。
func (dr *DependencyResolver) getDependents(name string) []string {
	dependents := make([]string, 0)
	for node, deps := range dr.graph {
		for _, dep := range deps {
			if dep == name {
				dependents = append(dependents, node)
				break
			}
		}
	}
	return dependents
}

// hasCycleDFS uses depth-first search to detect cycles.
// zh: hasCycleDFS 使用深度優先搜尋偵測循環依賴。
func (dr *DependencyResolver) hasCycleDFS(node string, visited, recStack map[string]bool) bool {
	visited[node] = true
	recStack[node] = true

	// Visit all dependencies
	for _, dep := range dr.graph[node] {
		if !visited[dep] {
			if dr.hasCycleDFS(dep, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[node] = false
	return false
}

// SimpleDependencyGraph implements a simple dependency graph.
// zh: SimpleDependencyGraph 實作簡單的依賴關係圖。
type SimpleDependencyGraph struct {
	nodes map[string]contracts.ComponentInfo
	edges map[string][]string
}

// AddNode adds a node to the dependency graph.
// zh: AddNode 新增節點到依賴關係圖。
func (sdg *SimpleDependencyGraph) AddNode(name string, info contracts.ComponentInfo) error {
	if sdg.nodes == nil {
		sdg.nodes = make(map[string]contracts.ComponentInfo)
	}
	if sdg.edges == nil {
		sdg.edges = make(map[string][]string)
	}

	sdg.nodes[name] = info
	sdg.edges[name] = info.Dependencies
	return nil
}

// AddEdge adds a dependency edge between two nodes.
// zh: AddEdge 在兩個節點間新增依賴邊。
func (sdg *SimpleDependencyGraph) AddEdge(from, to string) error {
	if sdg.edges == nil {
		sdg.edges = make(map[string][]string)
	}

	// Add dependency: 'from' depends on 'to'
	sdg.edges[from] = append(sdg.edges[from], to)
	return nil
}

// GetTopologicalOrder returns nodes in topological order.
// zh: GetTopologicalOrder 回傳拓撲排序的節點順序。
func (sdg *SimpleDependencyGraph) GetTopologicalOrder() ([]string, error) {
	resolver := &DependencyResolver{
		graph: sdg.edges,
		nodes: sdg.nodes,
	}
	return resolver.topologicalSort()
}

// HasCycle checks if the dependency graph has cycles.
// zh: HasCycle 檢查依賴關係圖是否有循環。
func (sdg *SimpleDependencyGraph) HasCycle() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range sdg.nodes {
		if !visited[name] {
			resolver := &DependencyResolver{
				graph: sdg.edges,
				nodes: sdg.nodes,
			}
			if resolver.hasCycleDFS(name, visited, recStack) {
				return true
			}
		}
	}
	return false
}

// GetDependents returns all nodes that depend on the given node.
// zh: GetDependents 回傳所有依賴於指定節點的節點。
func (sdg *SimpleDependencyGraph) GetDependents(name string) []string {
	dependents := make([]string, 0)
	for node, deps := range sdg.edges {
		for _, dep := range deps {
			if dep == name {
				dependents = append(dependents, node)
				break
			}
		}
	}
	return dependents
}

// GetDependencies returns all dependencies of the given node.
// zh: GetDependencies 回傳指定節點的所有依賴。
func (sdg *SimpleDependencyGraph) GetDependencies(name string) []string {
	if deps, exists := sdg.edges[name]; exists {
		return append([]string(nil), deps...) // return a copy
	}
	return []string{}
}
