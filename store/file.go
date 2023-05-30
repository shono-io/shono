package store

import "fmt"

func NewFileStore(conceptCode, code, directory string) *FileStore {
	return &FileStore{
		store: &store{
			conceptCode: conceptCode,
			code:        code,
			name:        fmt.Sprintf("%s File Store", code),
			description: fmt.Sprintf("A %s File Store storing data relative to the %q directory", code, directory),
		},
		directory: directory,
	}
}

type FileStore struct {
	*store
	directory string
}

func (s *FileStore) AsBenthosComponent() (map[string]interface{}, error) {
	return map[string]interface{}{
		"file": map[string]interface{}{
			"directory": s.directory,
		},
	}, nil
}
