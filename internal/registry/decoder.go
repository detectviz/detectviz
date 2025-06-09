package registry

import (
	"fmt"
	"os"

	"github.com/detectviz/detectviz/pkg/ifaces/registry"
	"github.com/detectviz/detectviz/pkg/registry/kinds"
	"gopkg.in/yaml.v3"
)

// DecodedResult 表示解碼後的 GVK 與原始資料。
// zh: 解碼後取得的 GVK 與 YAML 原文。
type DecodedResult struct {
	GVK     registry.GVK
	RawYAML []byte
}

// DecodeAndValidate 讀取指定檔案並解析為 GVK，驗證其符合 schema。
// zh: 讀取 YAML 檔案，解析出 GVK 並呼叫 schema validator。
func DecodeAndValidate(schemaPath, filePath string) (*DecodedResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var meta struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
	}
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse YAML header: %w", err)
	}

	group, version, err := parseAPIVersion(meta.APIVersion)
	if err != nil {
		return nil, err
	}

	gvk := registry.GVK{
		Group:   group,
		Version: version,
		Kind:    meta.Kind,
	}

	if err := kinds.Validate(schemaPath, data); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	return &DecodedResult{GVK: gvk, RawYAML: data}, nil
}

// parseAPIVersion 將 "group/version" 拆解為 GVK 所需欄位。
func parseAPIVersion(apiVersion string) (string, string, error) {
	var group, version string
	n, _ := fmt.Sscanf(apiVersion, "%[^/]/%s", &group, &version)
	if n != 2 {
		return "", "", fmt.Errorf("invalid apiVersion: %s", apiVersion)
	}
	return group, version, nil
}
