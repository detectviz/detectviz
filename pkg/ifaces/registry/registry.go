package registry

import (
	"context"
)

// GVK represents a Group-Version-Kind identifier for a resource.
// zh: GVK 是資源的分組、版本與類型，用於唯一標識一類型資源。
type GVK struct {
	Group   string
	Version string
	Kind    string
}

// Resource defines the interface that each registered resource must implement.
// zh: 所有可被註冊的資源都必須實作此介面。
type Resource interface {
	GetName() string
}

// Registry defines the behavior of a resource registry, providing GVK registration and CRUD operations.
// zh: Registry 定義資源註冊表的操作行為，提供 GVK 註冊與 CRUD 功能。
type Registry interface {
	Register(gvk GVK, handler ResourceHandler) error

	Get(ctx context.Context, gvk GVK, name string) (Resource, error)

	List(ctx context.Context, gvk GVK) ([]Resource, error)

	Create(ctx context.Context, gvk GVK, res Resource) error

	Update(ctx context.Context, gvk GVK, res Resource) error

	Delete(ctx context.Context, gvk GVK, name string) error
}

// ResourceHandler defines CRUD operations for a specific resource type.
// zh: ResourceHandler 負責某一資源類型的實際 CRUD 操作實作。
type ResourceHandler interface {
	Get(ctx context.Context, name string) (Resource, error)

	List(ctx context.Context) ([]Resource, error)

	Create(ctx context.Context, res Resource) error

	Update(ctx context.Context, res Resource) error

	Delete(ctx context.Context, name string) error
}
