package file

import (
	"github.com/mosuka/blast/master/registry"
	"github.com/mosuka/blast/master/store"
)

const (
	Name            = "file"
	DefaultBasePath = "/etc/blast"
)

func init() {
	registry.RegisterStore(Name, NewStore)
}

type Store struct {
	BasePath string
}

func NewStore(config map[string]interface{}) (store.Store, error) {
	path, ok := config["base_path"].(string)
	if !ok || path == "" {
		path = DefaultBasePath
	}

	store := Store{
		BasePath: path,
	}

	return &store, nil
}

func (s *Store) Reader() (store.Reader, error) {
	reader, err := NewReader(s)
	if err != nil {
		return nil, err
	}

	return &reader, nil
}

func (s *Store) Writer() (store.Writer, error) {
	writer, err := NewWriter(s)
	if err != nil {
		return nil, err
	}

	return &writer, nil
}

func (s *Store) Close() error {
	return nil
}
