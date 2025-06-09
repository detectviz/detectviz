package registry_test

import (
	"path/filepath"
	"testing"

	"github.com/detectviz/detectviz/internal/registry"
)

// TestDecodeAndValidate_ValidHost tests schema validation with valid host YAML.
// zh: 測試 valid_host.yaml 是否能正確通過 schema 驗證並取得正確 GVK。
func TestDecodeAndValidate_ValidHost(t *testing.T) {
	schema := filepath.Join("pkg", "registry", "schemas", "host.schema.yaml")
	file := filepath.Join("pkg", "registry", "kinds", "testdata", "valid_host.yaml")

	result, err := registry.DecodeAndValidate(schema, file)
	if err != nil {
		t.Fatalf("expected valid YAML, got error: %v", err)
	}

	if result.GVK.Group != "core" || result.GVK.Version != "v1" || result.GVK.Kind != "Host" {
		t.Errorf("unexpected GVK: %+v", result.GVK)
	}
	if len(result.RawYAML) == 0 {
		t.Error("expected non-empty RawYAML")
	}
}

// TestDecodeAndValidate_InvalidHost tests schema validation failure.
// zh: 測試 invalid_host.yaml 是否會正確拋出驗證錯誤。
func TestDecodeAndValidate_InvalidHost(t *testing.T) {
	schema := filepath.Join("pkg", "registry", "schemas", "host.schema.yaml")
	file := filepath.Join("pkg", "registry", "kinds", "testdata", "invalid_host.yaml")

	_, err := registry.DecodeAndValidate(schema, file)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
