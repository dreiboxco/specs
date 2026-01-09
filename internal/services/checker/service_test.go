package checker

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_Check_NumberingGap(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar specs com gap (00, 01, 03 - falta 02)
	spec00 := `# 00 Test Spec`
	spec01 := `# 01 Test Spec`
	spec03 := `# 03 Test Spec`

	if err := fs.WriteFile(filepath.Join(specsDir, "00-test.spec.md"), []byte(spec00), 0644); err != nil {
		t.Fatalf("falha ao criar spec00: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "01-test.spec.md"), []byte(spec01), 0644); err != nil {
		t.Fatalf("falha ao criar spec01: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "03-test.spec.md"), []byte(spec03), 0644); err != nil {
		t.Fatalf("falha ao criar spec03: %v", err)
	}

	result, err := service.Check(CheckOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Deve detectar gap
	foundGap := false
	for _, p := range result.Problems {
		if p.Category == "Numeração" && strings.Contains(p.Message, "Gap") {
			foundGap = true
			break
		}
	}
	if !foundGap {
		t.Error("deveria detectar gap na numeração")
	}
}

func TestService_Check_DuplicateNumbering(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar specs com numeração duplicada
	spec1 := `# 01 Test Spec`
	spec2 := `# 01 Other Spec`

	if err := fs.WriteFile(filepath.Join(specsDir, "01-test.spec.md"), []byte(spec1), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "01-other.spec.md"), []byte(spec2), 0644); err != nil {
		t.Fatalf("falha ao criar spec2: %v", err)
	}

	result, err := service.Check(CheckOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Deve detectar duplicata
	foundDuplicate := false
	for _, p := range result.Problems {
		if p.Category == "Numeração" && strings.Contains(p.Message, "duplicada") {
			foundDuplicate = true
			break
		}
	}
	if !foundDuplicate {
		t.Error("deveria detectar numeração duplicada")
	}
}

func TestService_Check_InvalidFileNameFormat(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar spec com formato inválido
	invalidSpec := `# Test Spec`

	if err := fs.WriteFile(filepath.Join(specsDir, "invalid-format.spec.md"), []byte(invalidSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.Check(CheckOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Deve detectar formato inválido
	foundFormatError := false
	for _, p := range result.Problems {
		if p.Category == "Formato" {
			foundFormatError = true
			break
		}
	}
	if !foundFormatError {
		t.Error("deveria detectar formato de nome inválido")
	}
}

func TestService_Check_BrokenLink(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar spec com link quebrado
	specWithLink := `# 01 Test Spec

Veja também [02-other.spec.md](02-other.spec.md)
`

	if err := fs.WriteFile(filepath.Join(specsDir, "01-test.spec.md"), []byte(specWithLink), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.Check(CheckOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Deve detectar link quebrado
	foundBrokenLink := false
	for _, p := range result.Problems {
		if p.Category == "Links" && strings.Contains(p.Message, "não encontrado") {
			foundBrokenLink = true
			break
		}
	}
	if !foundBrokenLink {
		t.Error("deveria detectar link quebrado")
	}
}

func TestService_Check_ValidLinks(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar specs com links válidos
	spec1 := `# 01 Test Spec

Veja também [02-other.spec.md](02-other.spec.md)
`
	spec2 := `# 02 Other Spec`

	if err := fs.WriteFile(filepath.Join(specsDir, "01-test.spec.md"), []byte(spec1), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "02-other.spec.md"), []byte(spec2), 0644); err != nil {
		t.Fatalf("falha ao criar spec2: %v", err)
	}

	result, err := service.Check(CheckOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Não deve ter problemas de links
	for _, p := range result.Problems {
		if p.Category == "Links" {
			t.Errorf("não deveria ter problemas de links, mas encontrou: %s", p.Message)
		}
	}
}

func TestService_Check_InvalidPath(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	_, err := service.Check(CheckOptions{
		Path: "/caminho/inexistente/12345",
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho inexistente")
	}
}

func TestService_Check_NotDirectory(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")

	if err := fs.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	_, err := service.Check(CheckOptions{
		Path: file,
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho que não é diretório")
	}
}
