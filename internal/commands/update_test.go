package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestUpdateCommand_Execute_Help(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewUpdateCommand(fs)

	exitCode := cmd.Execute([]string{"--help"})
	if exitCode != 0 {
		t.Errorf("esperado código 0 para --help, recebido %d", exitCode)
	}
}

func TestUpdateCommand_Execute_InvalidFlag(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewUpdateCommand(fs)

	exitCode := cmd.Execute([]string{"--invalid-flag"})
	if exitCode != 2 {
		t.Errorf("esperado código 2 para flag inválida, recebido %d", exitCode)
	}
}

func TestUpdateCommand_Execute_DryRun(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewUpdateCommand(fs)

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

	// Mudar para o diretório temporário
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	exitCode := cmd.Execute([]string{"--dry-run"})
	if exitCode != 0 {
		t.Errorf("esperado código 0 para dry-run, recebido %d", exitCode)
	}
}

func TestUpdateCommand_ParseArgs(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewUpdateCommand(fs)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
		},
		{
			name:    "dry-run flag",
			args:    []string{"--dry-run"},
			wantErr: false,
		},
		{
			name:    "multiple flags",
			args:    []string{"--dry-run", "--force"},
			wantErr: false,
		},
		{
			name:    "invalid flag",
			args:    []string{"--invalid"},
			wantErr: true,
		},
		{
			name:    "target directory",
			args:    []string{"/tmp/test"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := cmd.parseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && opts == nil {
				t.Error("parseArgs() retornou nil sem erro")
			}
		})
	}
}
