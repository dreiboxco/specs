package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestInitCommand_Execute_Success(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	tmpDir := t.TempDir()

	// Mudar para diretório temporário
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Executar comando sem argumentos (usa diretório atual)
	code := cmd.Execute([]string{})
	if code != 0 {
		t.Errorf("esperava código 0, obteve %d", code)
	}

	// Verificar que estrutura foi criada
	if !fs.Exists("specs") {
		t.Error("diretório specs/ não foi criado")
	}

	if !fs.Exists(".cursorrules") {
		t.Error(".cursorrules não foi criado")
	}

	if !fs.Exists("README.md") {
		t.Error("README.md não foi criado")
	}
}

func TestInitCommand_Execute_WithTargetDir(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "meu-projeto")

	// Criar diretório alvo
	if err := fs.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório alvo: %v", err)
	}

	// Executar comando com diretório específico
	code := cmd.Execute([]string{targetDir})
	if code != 0 {
		t.Errorf("esperava código 0, obteve %d", code)
	}

	// Verificar que estrutura foi criada no diretório alvo
	specsDir := filepath.Join(targetDir, "specs")
	if !fs.Exists(specsDir) {
		t.Error("diretório specs/ não foi criado no diretório alvo")
	}
}

func TestInitCommand_Execute_WithBoilerplate(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	tmpDir := t.TempDir()

	// Mudar para diretório temporário
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Executar comando com --with-boilerplate
	code := cmd.Execute([]string{"--with-boilerplate"})
	if code != 0 {
		t.Errorf("esperava código 0, obteve %d", code)
	}

	// Verificar que boilerplate foi criado
	boilerplateDir := filepath.Join(tmpDir, "boilerplate", "specs")
	if !fs.Exists(boilerplateDir) {
		t.Error("diretório boilerplate/specs/ não foi criado")
	}
}

func TestInitCommand_Execute_Idempotent(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	tmpDir := t.TempDir()

	// Mudar para diretório temporário
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Primeira execução
	code1 := cmd.Execute([]string{})
	if code1 != 0 {
		t.Errorf("primeira execução: esperava código 0, obteve %d", code1)
	}

	// Segunda execução (deve ser idempotente)
	code2 := cmd.Execute([]string{})
	if code2 != 0 {
		t.Errorf("segunda execução: esperava código 0, obteve %d", code2)
	}
}

func TestInitCommand_Execute_Help(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	// Executar com --help
	code := cmd.Execute([]string{"--help"})
	if code != 0 {
		t.Errorf("esperava código 0 para --help, obteve %d", code)
	}
}

func TestInitCommand_Execute_InvalidDirectory(t *testing.T) {
	cmd := NewInitCommand(adapters.NewFileSystem())

	// Executar com diretório inexistente
	code := cmd.Execute([]string{"/caminho/inexistente/12345"})
	if code != 1 {
		t.Errorf("esperava código 1 para diretório inexistente, obteve %d", code)
	}
}

func TestInitCommand_Execute_InvalidFlag(t *testing.T) {
	cmd := NewInitCommand(adapters.NewFileSystem())

	// Executar com flag inválida
	code := cmd.Execute([]string{"--invalid-flag"})
	if code != 2 {
		t.Errorf("esperava código 2 para flag inválida, obteve %d", code)
	}
}

func TestInitCommand_parseArgs(t *testing.T) {
	fs := adapters.NewFileSystem()
	cmd := NewInitCommand(fs)

	testCases := []struct {
		name     string
		args     []string
		wantHelp bool
		wantErr  bool
	}{
		{
			name:     "sem argumentos",
			args:     []string{},
			wantHelp: false,
			wantErr:  false,
		},
		{
			name:     "com --help",
			args:     []string{"--help"},
			wantHelp: true,
			wantErr:  false,
		},
		{
			name:     "com --force",
			args:     []string{"--force"},
			wantHelp: false,
			wantErr:  false,
		},
		{
			name:     "com --with-boilerplate",
			args:     []string{"--with-boilerplate"},
			wantHelp: false,
			wantErr:  false,
		},
		{
			name:     "com diretório",
			args:     []string{"./meu-projeto"},
			wantHelp: false,
			wantErr:  false,
		},
		{
			name:     "com flag inválida",
			args:     []string{"--invalid"},
			wantHelp: false,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := cmd.parseArgs(tc.args)
			if tc.wantErr && err == nil {
				t.Error("esperava erro, mas não houve")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
			if !tc.wantErr && opts != nil {
				if opts.Help != tc.wantHelp {
					t.Errorf("Help esperado %v, obteve %v", tc.wantHelp, opts.Help)
				}
			}
		})
	}
}
