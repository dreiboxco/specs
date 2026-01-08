package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_GetVersion(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	// Criar diretório temporário para teste
	tmpDir := t.TempDir()
	versionFile := filepath.Join(tmpDir, "VERSION")

	// Teste com versão válida
	testVersion := "1.2.3"
	err := os.WriteFile(versionFile, []byte(testVersion), 0644)
	if err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	// Mudar para o diretório temporário
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	version, err := service.ReadVersionFile(versionFile)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if version != testVersion {
		t.Errorf("versão esperada %s, obtida %s", testVersion, version)
	}
}

func TestService_InvalidVersion(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	testCases := []struct {
		name    string
		content string
	}{
		{"vazio", ""},
		{"formato inválido", "1.2"},
		{"formato inválido 2", "1.2.3.4"},
		{"caracteres não numéricos", "1.2.a"},
		{"espaços", "1.2.3 "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			versionFile := filepath.Join(tmpDir, "VERSION")

			err := os.WriteFile(versionFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("falha ao criar arquivo de teste: %v", err)
			}

			_, err = service.ReadVersionFile(versionFile)
			if err == nil {
				t.Errorf("esperava erro para conteúdo '%s', mas não houve erro", tc.content)
			}
		})
	}
}

func TestService_FileNotFound(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	_, err := service.ReadVersionFile("/caminho/inexistente/VERSION")
	if err == nil {
		t.Error("esperava erro para arquivo não encontrado, mas não houve erro")
	}
}
