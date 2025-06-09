package versioning

import (
	"fmt"
	"strings"
	"sync"
)

// SimpleVersioned 是預設版本結構實作。
// zh: 基礎版本資源結構，可作為內嵌使用。
type SimpleVersioned struct {
	Version string
}

// GetVersion 回傳物件的版本。
// zh: 取得該資源的版本標記。
func (v SimpleVersioned) GetVersion() string {
	return v.Version
}

// VersionMap 提供簡單的版本註冊與比較邏輯。
// zh: 使用 map 實作的版本控制器。
type VersionMap struct {
	mu       sync.RWMutex
	registry map[string]Versioned
}

// NewVersionMap 建立空白版本註冊器。
// zh: 回傳新的 VersionMap。
func NewVersionMap() *VersionMap {
	return &VersionMap{
		registry: make(map[string]Versioned),
	}
}

// Register 加入一筆命名版本。
// zh: 註冊一筆資源與其版本。
func (vm *VersionMap) Register(name string, v Versioned) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if _, exists := vm.registry[name]; exists {
		return fmt.Errorf("version already registered for: %s", name)
	}
	vm.registry[name] = v
	return nil
}

// Get 查詢已註冊的版本。
// zh: 根據名稱取得對應的版本物件。
func (vm *VersionMap) Get(name string) (Versioned, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	v, ok := vm.registry[name]
	return v, ok
}

// Compare 比較兩個版本字串。
// zh: 回傳 a 是否大於 b：1 表 a > b，-1 表 a < b，0 表相等。
func (vm *VersionMap) Compare(a, b string) int {
	as := normalizeVersion(a)
	bs := normalizeVersion(b)

	for i := 0; i < len(as) && i < len(bs); i++ {
		if as[i] > bs[i] {
			return 1
		} else if as[i] < bs[i] {
			return -1
		}
	}

	if len(as) > len(bs) {
		return 1
	} else if len(as) < len(bs) {
		return -1
	}
	return 0
}

func normalizeVersion(v string) []int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	out := make([]int, len(parts))
	for i, p := range parts {
		fmt.Sscanf(p, "%d", &out[i])
	}
	return out
}
