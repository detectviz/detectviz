package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SchemaMapping 表示一筆 GVK 與 schema file 的對應。
// zh: 一筆 GVK 與對應 schema 檔案的關聯。
type SchemaMapping struct {
	Group   string `json:"group" yaml:"group"`
	Version string `json:"version" yaml:"version"`
	Kind    string `json:"kind" yaml:"kind"`
	File    string `json:"file" yaml:"file"`
}

// Index 是 index.yaml 解析後的完整資料。
// zh: 載入後的 schema index 結構。
type Index struct {
	Schemas []SchemaMapping `json:"schemas" yaml:"schemas"`
}

// LoadSchemaIndex 讀取指定路徑下的 index.yaml。
// zh: 載入所有已定義的 GVK 與 schema 對應檔案。
func LoadSchemaIndex(path string) (*Index, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema index: %w", err)
	}

	var idx Index
	if err := unmarshalYAML(content, &idx); err != nil {
		return nil, fmt.Errorf("failed to parse schema index: %w", err)
	}
	return &idx, nil
}

// unmarshalYAML 是簡化版解析器。
// 你可以根據使用的 YAML 套件更換此實作。
func unmarshalYAML(data []byte, out interface{}) error {
	return json.Unmarshal(data, out) // 臨時使用 JSON 模擬，實務上應使用 YAML 套件
}

// FindSchemaPath 回傳指定 GVK 的 schema 檔案路徑（含 baseDir）。
func (i *Index) FindSchemaPath(baseDir, group, version, kind string) (string, error) {
	for _, entry := range i.Schemas {
		if entry.Group == group && entry.Version == version && entry.Kind == kind {
			return filepath.Join(baseDir, entry.File), nil
		}
	}
	return "", fmt.Errorf("schema not found for %s/%s/%s", group, version, kind)
}
