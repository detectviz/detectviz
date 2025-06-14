package navtree

import (
	"fmt"
	"sort"
	"sync"

	"detectviz/pkg/platform/contracts"
)

// Builder implements NavTreeBuilder interface for building navigation trees.
// zh: Builder 實作 NavTreeBuilder 介面用於建構導覽樹。
type Builder struct {
	nodes map[string]*contracts.NavNode
	tree  []contracts.NavNode
	mutex sync.RWMutex
}

// NewBuilder creates a new navigation tree builder.
// zh: NewBuilder 建立新的導覽樹建構器。
func NewBuilder() *Builder {
	return &Builder{
		nodes: make(map[string]*contracts.NavNode),
		tree:  make([]contracts.NavNode, 0),
	}
}

// AddNode adds a navigation node to the tree.
// zh: AddNode 向樹中添加導覽節點。
func (b *Builder) AddNode(id string, node contracts.NavNode) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, exists := b.nodes[id]; exists {
		return fmt.Errorf("navigation node %s already exists", id)
	}

	// Set the ID if not already set
	if node.ID == "" {
		node.ID = id
	}

	// Store the node
	b.nodes[id] = &node

	// Add to root level if no parent
	b.tree = append(b.tree, node)

	// Sort by order
	b.sortNodes(b.tree)

	return nil
}

// AddChildNode adds a child node to an existing parent node.
// zh: AddChildNode 向現有父節點添加子節點。
func (b *Builder) AddChildNode(parentID string, id string, node contracts.NavNode) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if child node already exists
	if _, exists := b.nodes[id]; exists {
		return fmt.Errorf("child navigation node %s already exists", id)
	}

	// Find parent node
	parentNode, exists := b.nodes[parentID]
	if !exists {
		return fmt.Errorf("parent navigation node %s not found", parentID)
	}

	// Set the ID if not already set
	if node.ID == "" {
		node.ID = id
	}

	// Store the child node
	b.nodes[id] = &node

	// Add to parent's children
	if parentNode.Children == nil {
		parentNode.Children = make([]contracts.NavNode, 0)
	}
	parentNode.Children = append(parentNode.Children, node)

	// Sort children by order
	b.sortNodes(parentNode.Children)

	// Update the parent in the tree
	b.updateNodeInTree(parentID, *parentNode)

	return nil
}

// RemoveNode removes a navigation node from the tree.
// zh: RemoveNode 從樹中移除導覽節點。
func (b *Builder) RemoveNode(id string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, exists := b.nodes[id]; !exists {
		return fmt.Errorf("navigation node %s not found", id)
	}

	// Remove from nodes map
	delete(b.nodes, id)

	// Remove from tree (recursive)
	b.tree = b.removeNodeFromSlice(b.tree, id)

	// Remove from all parent nodes' children
	for _, node := range b.nodes {
		if node.Children != nil {
			node.Children = b.removeNodeFromSlice(node.Children, id)
		}
	}

	return nil
}

// GetNode retrieves a navigation node by ID.
// zh: GetNode 根據 ID 取得導覽節點。
func (b *Builder) GetNode(id string) (*contracts.NavNode, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	node, exists := b.nodes[id]
	if !exists {
		return nil, fmt.Errorf("navigation node %s not found", id)
	}

	// Return a copy to prevent external modification
	nodeCopy := *node
	return &nodeCopy, nil
}

// SetNodePermission sets the permission requirement for a navigation node.
// zh: SetNodePermission 設定導覽節點的權限要求。
func (b *Builder) SetNodePermission(id string, permission string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	node, exists := b.nodes[id]
	if !exists {
		return fmt.Errorf("navigation node %s not found", id)
	}

	node.Permission = permission

	// Update the node in the tree
	b.updateNodeInTree(id, *node)

	return nil
}

// Additional methods for navigation tree management
// zh: 導覽樹管理的額外方法

// BuildTree builds the complete navigation tree.
// zh: BuildTree 建構完整的導覽樹。
func (b *Builder) BuildTree() []contracts.NavNode {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Return a copy of the tree
	treeCopy := make([]contracts.NavNode, len(b.tree))
	copy(treeCopy, b.tree)

	return treeCopy
}

// BuildTreeForUser builds a navigation tree filtered by user permissions.
// zh: BuildTreeForUser 建構依使用者權限篩選的導覽樹。
func (b *Builder) BuildTreeForUser(user *contracts.UserInfo) []contracts.NavNode {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if user == nil {
		return make([]contracts.NavNode, 0)
	}

	var filteredTree []contracts.NavNode
	for _, node := range b.tree {
		if filteredNode := b.filterNodeByPermission(node, user); filteredNode != nil {
			filteredTree = append(filteredTree, *filteredNode)
		}
	}

	return filteredTree
}

// GetStats returns statistics about the navigation tree.
// zh: GetStats 回傳導覽樹的統計資訊。
func (b *Builder) GetStats() map[string]any {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	stats := map[string]any{
		"total_nodes":            len(b.nodes),
		"root_nodes":             len(b.tree),
		"nodes_with_children":    b.countNodesWithChildren(),
		"nodes_with_permissions": b.countNodesWithPermissions(),
		"enabled_nodes":          b.countEnabledNodes(),
		"visible_nodes":          b.countVisibleNodes(),
	}

	return stats
}

// Helper methods
// zh: 輔助方法

// sortNodes sorts navigation nodes by their order.
// zh: sortNodes 按順序排序導覽節點。
func (b *Builder) sortNodes(nodes []contracts.NavNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Order < nodes[j].Order
	})
}

// updateNodeInTree updates a node in the tree structure.
// zh: updateNodeInTree 更新樹結構中的節點。
func (b *Builder) updateNodeInTree(nodeID string, updatedNode contracts.NavNode) {
	// Update in root level
	for i, node := range b.tree {
		if node.ID == nodeID {
			b.tree[i] = updatedNode
			return
		}
		// Update in children recursively
		if b.updateNodeInChildren(node.Children, nodeID, updatedNode) {
			return
		}
	}
}

// updateNodeInChildren updates a node in children slice recursively.
// zh: updateNodeInChildren 遞歸更新子節點片段中的節點。
func (b *Builder) updateNodeInChildren(children []contracts.NavNode, nodeID string, updatedNode contracts.NavNode) bool {
	for i, child := range children {
		if child.ID == nodeID {
			children[i] = updatedNode
			return true
		}
		if b.updateNodeInChildren(child.Children, nodeID, updatedNode) {
			return true
		}
	}
	return false
}

// removeNodeFromSlice removes a node from a slice by ID.
// zh: removeNodeFromSlice 根據 ID 從切片中移除節點。
func (b *Builder) removeNodeFromSlice(nodes []contracts.NavNode, nodeID string) []contracts.NavNode {
	var result []contracts.NavNode
	for _, node := range nodes {
		if node.ID != nodeID {
			// Keep the node but also check its children
			if node.Children != nil {
				node.Children = b.removeNodeFromSlice(node.Children, nodeID)
			}
			result = append(result, node)
		}
	}
	return result
}

// filterNodeByPermission filters a node based on user permissions.
// zh: filterNodeByPermission 根據使用者權限篩選節點。
func (b *Builder) filterNodeByPermission(node contracts.NavNode, user *contracts.UserInfo) *contracts.NavNode {
	// Check if node requires permission
	if node.Permission != "" {
		hasPermission := false
		for _, permission := range user.Permissions {
			if b.matchesPermission(node.Permission, permission) {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return nil
		}
	}

	// Check if node is enabled and visible
	if !node.Enabled || !node.Visible {
		return nil
	}

	// Filter children recursively
	var filteredChildren []contracts.NavNode
	for _, child := range node.Children {
		if filteredChild := b.filterNodeByPermission(child, user); filteredChild != nil {
			filteredChildren = append(filteredChildren, *filteredChild)
		}
	}

	// Create a copy with filtered children
	filteredNode := node
	filteredNode.Children = filteredChildren

	return &filteredNode
}

// matchesPermission checks if a required permission matches a user permission.
// zh: matchesPermission 檢查所需權限是否符合使用者權限。
func (b *Builder) matchesPermission(required string, userPermission contracts.Permission) bool {
	// Simple permission matching - can be enhanced for more complex rules
	// Format: action.resource (e.g., "view.system", "admin.plugins")

	if required == userPermission.Resource {
		return true
	}

	// Check if user has wildcard permission
	if userPermission.Resource == "*" {
		return true
	}

	// Check scope
	for _, scope := range userPermission.Scope {
		if scope == required || scope == "*" {
			return true
		}
	}

	return false
}

// Statistics helper methods
// zh: 統計輔助方法

// countNodesWithChildren counts nodes that have children.
// zh: countNodesWithChildren 計算有子節點的節點數量。
func (b *Builder) countNodesWithChildren() int {
	count := 0
	for _, node := range b.nodes {
		if len(node.Children) > 0 {
			count++
		}
	}
	return count
}

// countNodesWithPermissions counts nodes that have permission requirements.
// zh: countNodesWithPermissions 計算有權限要求的節點數量。
func (b *Builder) countNodesWithPermissions() int {
	count := 0
	for _, node := range b.nodes {
		if node.Permission != "" {
			count++
		}
	}
	return count
}

// countEnabledNodes counts enabled nodes.
// zh: countEnabledNodes 計算啟用的節點數量。
func (b *Builder) countEnabledNodes() int {
	count := 0
	for _, node := range b.nodes {
		if node.Enabled {
			count++
		}
	}
	return count
}

// countVisibleNodes counts visible nodes.
// zh: countVisibleNodes 計算可見的節點數量。
func (b *Builder) countVisibleNodes() int {
	count := 0
	for _, node := range b.nodes {
		if node.Visible {
			count++
		}
	}
	return count
}
