package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/templates"
)

// Service gerencia operações de inicialização de projetos SDD
type Service struct {
	fs adapters.FileSystem
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{fs: fs}
}

// InitOptions contém opções para inicialização
type InitOptions struct {
	TargetDir        string
	Force            bool
	WithBoilerplate  bool
}

// InitResult contém resultado da inicialização
type InitResult struct {
	SpecsDir      string
	FilesCreated  []string
	DirectoriesCreated []string
}

// Initialize inicializa um novo projeto SDD
func (s *Service) Initialize(opts InitOptions) (*InitResult, error) {
	// Determinar diretório alvo
	targetDir := opts.TargetDir
	if targetDir == "" {
		wd, err := s.fs.Getwd()
		if err != nil {
			return nil, fmt.Errorf("falha ao obter diretório atual: %w", err)
		}
		targetDir = wd
	}

	// Verificar se diretório existe
	if !s.fs.Exists(targetDir) {
		return nil, fmt.Errorf("diretório não existe: %s", targetDir)
	}

	// Verificar se já é projeto SDD
	if s.isSDDProject(targetDir) {
		return &InitResult{
			SpecsDir: filepath.Join(targetDir, "specs"),
		}, nil // Idempotente - retorna sucesso sem criar nada
	}

	// Verificar permissões de escrita
	if err := s.checkWritePermissions(targetDir); err != nil {
		return nil, fmt.Errorf("sem permissão de escrita no diretório: %w", err)
	}

	result := &InitResult{
		FilesCreated:       []string{},
		DirectoriesCreated: []string{},
	}

	// Criar diretório specs/
	specsDir := filepath.Join(targetDir, "specs")
	if err := s.fs.MkdirAll(specsDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório specs: %w", err)
	}
	result.DirectoriesCreated = append(result.DirectoriesCreated, specsDir)
	result.SpecsDir = specsDir

	// Copiar templates de specs
	if err := s.copySpecTemplates(specsDir, opts.Force); err != nil {
		return nil, fmt.Errorf("falha ao copiar templates: %w", err)
	}

	// Criar .cursorrules
	cursorRulesPath := filepath.Join(targetDir, ".cursorrules")
	cursorRulesContent, err := templates.GetCursorRulesTemplate()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter template .cursorrules: %w", err)
	}
	if err := s.createFileIfNotExists(cursorRulesPath, cursorRulesContent, opts.Force); err != nil {
		return nil, fmt.Errorf("falha ao criar .cursorrules: %w", err)
	}
	result.FilesCreated = append(result.FilesCreated, cursorRulesPath)

	// Criar README.md
	readmePath := filepath.Join(targetDir, "README.md")
	if err := s.createFileIfNotExists(readmePath, templates.ReadmeTemplate, opts.Force); err != nil {
		return nil, fmt.Errorf("falha ao criar README.md: %w", err)
	}
	result.FilesCreated = append(result.FilesCreated, readmePath)

	// Criar boilerplate/ se solicitado
	if opts.WithBoilerplate {
		boilerplateDir := filepath.Join(targetDir, "boilerplate", "specs")
		if err := s.fs.MkdirAll(boilerplateDir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretório boilerplate: %w", err)
		}
		result.DirectoriesCreated = append(result.DirectoriesCreated, boilerplateDir)
		// Copiar templates para boilerplate também
		if err := s.copySpecTemplates(boilerplateDir, opts.Force); err != nil {
			return nil, fmt.Errorf("falha ao copiar templates para boilerplate: %w", err)
		}
	}

	return result, nil
}

// isSDDProject verifica se diretório já contém projeto SDD
func (s *Service) isSDDProject(dir string) bool {
	specsDir := filepath.Join(dir, "specs")
	if !s.fs.Exists(specsDir) {
		return false
	}

	// Verificar se existe pelo menos um arquivo 00-*.spec.md
	templateNames := templates.GetAllTemplateNames()
	for _, name := range templateNames {
		if strings.HasPrefix(name, "00-") {
			path := filepath.Join(specsDir, name)
			if s.fs.Exists(path) {
				return true
			}
		}
	}

	return false
}

// checkWritePermissions verifica permissões de escrita
func (s *Service) checkWritePermissions(dir string) error {
	testFile := filepath.Join(dir, ".specs-write-test")
	defer func() {
		// Tentar remover arquivo de teste se existir
		if s.fs.Exists(testFile) {
			_ = os.Remove(testFile)
		}
	}()

	// Tentar criar arquivo de teste
	if err := s.fs.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("sem permissão de escrita: %w", err)
	}

	return nil
}

// copySpecTemplates copia templates de specs para diretório destino
func (s *Service) copySpecTemplates(destDir string, force bool) error {
	templateNames := templates.GetAllTemplateNames()

	for _, name := range templateNames {
		template, exists := templates.GetTemplate(name)
		if !exists {
			continue
		}

		destPath := filepath.Join(destDir, name)
		if err := s.createFileIfNotExists(destPath, template, force); err != nil {
			return fmt.Errorf("falha ao copiar template %s: %w", name, err)
		}
	}

	return nil
}

// createFileIfNotExists cria arquivo se não existir, ou sobrescreve se force=true
func (s *Service) createFileIfNotExists(path string, content []byte, force bool) error {
	if s.fs.Exists(path) && !force {
		// Arquivo existe e não foi solicitado sobrescrever
		return nil
	}

	return s.fs.WriteFile(path, content, 0644)
}
