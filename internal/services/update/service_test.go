package update

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_Update_ValidatesSDDProject(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	opts := UpdateOptions{
		TargetDir: tmpDir,
	}

	_, err := service.Update(opts)
	if err == nil {
		t.Error("esperado erro para diretório que não é projeto SDD")
	}

	if !strings.Contains(err.Error(), "projeto SDD válido") {
		t.Errorf("erro esperado sobre projeto SDD, recebido: %v", err)
	}
}

func TestService_Update_UpdatesStaticTemplates(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	// Criar estrutura de projeto SDD
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar specs dir: %v", err)
	}

	// Criar arquivo 00-*.spec.md para validar como SDD
	archFile := filepath.Join(specsDir, "00-architecture.spec.md")
	if err := fs.WriteFile(archFile, []byte("# Arquitetura"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	// Criar checklist.md existente
	checklistFile := filepath.Join(specsDir, "checklist.md")
	if err := fs.WriteFile(checklistFile, []byte("# Checklist antigo"), 0644); err != nil {
		t.Fatalf("falha ao criar checklist: %v", err)
	}

	opts := UpdateOptions{
		TargetDir: tmpDir,
		DryRun:    true,
	}

	result, err := service.Update(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar que checklist.md está na lista de atualizados
	found := false
	for _, file := range result.FilesUpdated {
		if file == "checklist.md" {
			found = true
			break
		}
	}

	if !found {
		t.Error("checklist.md deveria estar na lista de arquivos atualizados")
	}
}

func TestService_Update_CreatesBackup(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	// Criar estrutura de projeto SDD
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar specs dir: %v", err)
	}

	archFile := filepath.Join(specsDir, "00-architecture.spec.md")
	if err := fs.WriteFile(archFile, []byte("# Arquitetura"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	opts := UpdateOptions{
		TargetDir: tmpDir,
		NoBackup:  false,
	}

	result, err := service.Update(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.BackupDir == "" {
		t.Error("backup deveria ter sido criado")
	}

	if !fs.Exists(result.BackupDir) {
		t.Error("diretório de backup deveria existir")
	}
}

func TestService_DetectCustomizations_NoCustomizations(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	current := []byte("# Cursor Rules\n## Seção 1\nConteúdo")
	boilerplate := []byte("# Cursor Rules\n## Seção 1\nConteúdo")

	hasCustomizations := service.detectCustomizations(current, boilerplate)
	if hasCustomizations {
		t.Error("não deveria detectar personalizações em arquivos idênticos")
	}
}

func TestService_DetectCustomizations_WithCustomizations(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	current := []byte("# Cursor Rules\n## Seção 1\nConteúdo personalizado\nLinha extra\nMais uma\nE outra")
	boilerplate := []byte("# Cursor Rules\n## Seção 1\nConteúdo")

	hasCustomizations := service.detectCustomizations(current, boilerplate)
	if !hasCustomizations {
		t.Error("deveria detectar personalizações quando há diferenças significativas")
	}
}

func TestService_DetectCustomizations_NewSection(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	current := []byte("# Cursor Rules\n## Seção 1\nConteúdo\n## Seção Personalizada\nConteúdo personalizado")
	boilerplate := []byte("# Cursor Rules\n## Seção 1\nConteúdo")

	hasCustomizations := service.detectCustomizations(current, boilerplate)
	if !hasCustomizations {
		t.Error("deveria detectar personalizações quando há seção nova")
	}
}

