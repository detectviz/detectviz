package modules

import (
	"fmt"
)

// DependencyGraph defines module dependencies.
// zh: DependencyGraph 用於描述模組之間的依賴關係。
type DependencyGraph struct {
	nodes map[string][]string
}

// NewDependencyGraph creates a new empty dependency graph.
// zh: 建立新的模組依賴圖。
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string][]string),
	}
}

// AddDependency adds a dependency from module 'from' to module 'to'.
// zh: 新增模組之間的依賴關係，表示 from 依賴 to。
func (g *DependencyGraph) AddDependency(from, to string) {
	g.nodes[from] = append(g.nodes[from], to)
}

// TopologicalSort returns a valid start order of modules based on dependencies.
// It returns an error if there is a cycle.
// zh: 根據依賴關係進行拓撲排序，若有循環依賴則回傳錯誤。
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	result := []string{}

	var visit func(string) error
	visit = func(n string) error {
		if temp[n] {
			return fmt.Errorf("cycle detected at module %q", n)
		}
		if !visited[n] {
			temp[n] = true
			for _, dep := range g.nodes[n] {
				if err := visit(dep); err != nil {
					return err
				}
			}
			temp[n] = false
			visited[n] = true
			result = append(result, n)
		}
		return nil
	}

	for n := range g.nodes {
		if !visited[n] {
			if err := visit(n); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
