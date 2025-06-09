package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/detectviz/detectviz/pkg/registry/kinds"
)

// TestValidateHostSchema_Valid validates a known-valid host YAML against its schema.
// zh: 驗證合法的 host YAML 是否能正確通過對應 schema 驗證。
func TestValidateHostSchema_Valid(t *testing.T) {
	path := filepath.Join("pkg", "registry", "schemas", "host.schema.yaml")
	input := filepath.Join("pkg", "registry", "kinds", "testdata", "valid_host.yaml")

	data, err := os.ReadFile(input)
	if err != nil {
		t.Fatalf("failed to read input: %v", err)
	}

	if err := kinds.Validate(path, data); err != nil {
		t.Errorf("expected valid input, got error: %v", err)
	}
}

// TestValidateHostSchema_Invalid checks if an invalid host YAML fails schema validation as expected.
// zh: 驗證不合法的 host YAML 是否會正確被驗證器擋下並回傳錯誤。
func TestValidateHostSchema_Invalid(t *testing.T) {
	path := filepath.Join("pkg", "registry", "schemas", "host.schema.yaml")
	input := filepath.Join("pkg", "registry", "kinds", "testdata", "invalid_host.yaml")

	data, err := os.ReadFile(input)
	if err != nil {
		t.Fatalf("failed to read input: %v", err)
	}

	err = kinds.Validate(path, data)
	if err == nil {
		t.Errorf("expected error for invalid input, got nil")
	} else if !kinds.IsSchemaValidationError(err) {
		t.Errorf("unexpected error type: %v", err)
	}
}
