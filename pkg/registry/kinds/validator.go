package kinds

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
)

// SchemaValidationError represents a standardized error for schema validation failures.
// zh: SchemaValidationError 用來包裝驗證不通過的錯誤。
type SchemaValidationError struct {
	Msg string
}

func (e *SchemaValidationError) Error() string {
	return fmt.Sprintf("schema validation failed: %s", e.Msg)
}

// IsSchemaValidationError checks if the given error is a SchemaValidationError.
// zh: 判斷錯誤是否為 SchemaValidationError。
func IsSchemaValidationError(err error) bool {
	var e *SchemaValidationError
	return errors.As(err, &e)
}

// Validate checks if given YAML data conforms to the provided schema file.
// zh: 驗證指定 YAML 是否符合對應的 schema 格式。
func Validate(schemaPath string, data []byte) error {
	// 載入 schema 為 cue.Instance
	schemaInstance := load.Instances([]string{schemaPath}, nil)
	if len(schemaInstance) == 0 || schemaInstance[0].Err != nil {
		return fmt.Errorf("invalid cue schema instance: %v", schemaInstance[0].Err)
	}

	var ctx cue.Context
	schema := ctx.BuildInstance(schemaInstance[0])

	// 將 YAML 轉為 cue.Value
	value := ctx.CompileString(string(data))

	// 合併 schema 與資料並驗證
	if err := schema.Unify(value).Validate(); err != nil {
		return &SchemaValidationError{Msg: err.Error()}
	}

	return nil
}
