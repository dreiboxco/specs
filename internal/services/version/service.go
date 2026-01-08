package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
)

// Service gerencia operações relacionadas a versão
type Service struct {
	fs adapters.FileSystem
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{fs: fs}
}

// GetVersion lê a versão do arquivo VERSION
func (s *Service) GetVersion() (string, error) {
	// Tentar encontrar o arquivo VERSION na raiz do projeto
	// Primeiro tentar no diretório atual
	versionPath := "VERSION"
	
	// Se não encontrar, tentar subir alguns níveis (para desenvolvimento)
	if !s.fs.Exists(versionPath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("falha ao obter diretório atual: %w", err)
		}
		
		for i := 0; i < 5; i++ {
			testPath := filepath.Join(currentDir, "VERSION")
			if s.fs.Exists(testPath) {
				versionPath = testPath
				break
			}
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				break
			}
			currentDir = parent
		}
	}

	if !s.fs.Exists(versionPath) {
		return "", fmt.Errorf("arquivo VERSION não encontrado")
	}

	data, err := s.fs.ReadFile(versionPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo VERSION: %w", err)
	}

	rawVersion := string(data)
	version := strings.TrimSpace(rawVersion)
	
	// Verificar se havia espaços ou caracteres extras (exceto quebras de linha no final)
	trimmedNewlines := strings.TrimRight(rawVersion, "\n\r")
	trimmedAll := strings.TrimSpace(rawVersion)
	if trimmedNewlines != trimmedAll {
		return "", fmt.Errorf("versão inválida no arquivo VERSION: contém espaços ou caracteres extras")
	}
	
	if version == "" {
		return "", fmt.Errorf("arquivo VERSION está vazio")
	}

	// Validar formato semântico básico
	if !isValidSemanticVersion(version) {
		return "", fmt.Errorf("versão inválida no arquivo VERSION: %s (formato esperado: MAJOR.MINOR.PATCH)", version)
	}

	return version, nil
}

// isValidSemanticVersion valida formato básico de versão semântica
func isValidSemanticVersion(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}
	
	// Verificar se todas as partes são numéricas (validação básica)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return false
		}
		// Verificar se é numérico
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	
	return true
}

// ReadVersionFile lê versão de um arquivo específico (para testes)
func (s *Service) ReadVersionFile(path string) (string, error) {
	if !s.fs.Exists(path) {
		return "", fmt.Errorf("arquivo não encontrado: %s", path)
	}

	data, err := s.fs.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo: %w", err)
	}

	rawVersion := string(data)
	version := strings.TrimSpace(rawVersion)
	
	// Verificar se havia espaços ou caracteres extras (exceto quebras de linha no final)
	trimmedNewlines := strings.TrimRight(rawVersion, "\n\r")
	trimmedAll := strings.TrimSpace(rawVersion)
	if trimmedNewlines != trimmedAll {
		return "", fmt.Errorf("versão inválida: contém espaços ou caracteres extras")
	}
	
	if version == "" {
		return "", fmt.Errorf("arquivo está vazio")
	}

	if !isValidSemanticVersion(version) {
		return "", fmt.Errorf("versão inválida: %s", version)
	}

	return version, nil
}

