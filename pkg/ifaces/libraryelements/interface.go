package libraryelements

// pkg/libraryelements/interface.go
type ElementService interface {
	GetElement(id string) (*Element, error)
	CreateElement(kind string, data any) error
}
