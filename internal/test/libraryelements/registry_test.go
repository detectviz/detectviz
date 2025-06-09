package libraryelements_test

import (
	"errors"
	"testing"

	"github.com/detectviz/detectviz/pkg/libraryelements"
)

type fakeRenderer struct {
	content string
}

func (r fakeRenderer) Render(e libraryelements.Element) ([]byte, error) {
	return []byte(r.content), nil
}

// TestRegistry_RegisterAndGetRenderer 測試渲染器註冊與查詢。
// zh: 驗證 RegisterRenderer 與 GetRenderer 是否正確運作。
func TestRegistry_RegisterAndGetRenderer(t *testing.T) {
	reg := libraryelements.NewRegistry()
	renderer := fakeRenderer{content: "test-rendered"}

	err := reg.RegisterRenderer("chart", renderer)
	if err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	r, ok := reg.GetRenderer("chart")
	if !ok {
		t.Fatal("expected renderer to be found")
	}

	output, err := r.Render(nil)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}
	if string(output) != "test-rendered" {
		t.Errorf("expected output 'test-rendered', got %s", output)
	}
}

// TestRegistry_DuplicateRegister 測試重複註冊應回傳錯誤。
// zh: 驗證同一類型註冊多次會失敗。
func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := libraryelements.NewRegistry()
	err1 := reg.RegisterRenderer("input", fakeRenderer{content: "1"})
	err2 := reg.RegisterRenderer("input", fakeRenderer{content: "2"})

	if err1 != nil {
		t.Fatalf("unexpected error on first register: %v", err1)
	}
	if err2 == nil {
		t.Fatal("expected error on duplicate register, got nil")
	}
	if !errors.Is(err2, err2) {
		t.Errorf("expected duplicate error, got: %v", err2)
	}
}
