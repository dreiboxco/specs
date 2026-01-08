package adapters

import (
	"os"
	"path/filepath"
)

// FileSystem interface para abstração de operações de sistema de arquivos
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Exists(path string) bool
}

// fileSystem implementa FileSystem usando os padrão
type fileSystem struct{}

// NewFileSystem cria uma nova instância de FileSystem
func NewFileSystem() FileSystem {
	return &fileSystem{}
}

func (fs *fileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *fileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (fs *fileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fs *fileSystem) GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

