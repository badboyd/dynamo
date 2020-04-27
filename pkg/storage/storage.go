package storage

// Storage interface
type Storage interface {
	Upload() error
	Delete(id string) error
}