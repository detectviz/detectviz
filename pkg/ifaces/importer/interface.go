package importer

import "context"

// pkg/importer/interface.go
type Importer interface {
	CanImport(src string) bool
	Import(ctx context.Context, raw []byte) (*DomainObject, error)
}
