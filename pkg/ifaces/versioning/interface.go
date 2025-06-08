package versioning

import "context"

// pkg/versioning/interface.go
type VersionStore interface {
	SaveVersion(ctx context.Context, item VersionedResource) error
	ListVersions(ctx context.Context, id string) ([]VersionInfo, error)
}
