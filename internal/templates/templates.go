package templates

import (
	"os"
	"path/filepath"
)

// GetTemplate retorna o template pelo nome
// Lê do filesystem relativo ao executável ou do diretório de trabalho
func GetTemplate(name string) ([]byte, bool) {
	// Tentar encontrar boilerplate relativo ao executável
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		// Tentar caminhos possíveis
		paths := []string{
			filepath.Join(exeDir, "..", "boilerplate", "specs", name),
			filepath.Join(exeDir, "boilerplate", "specs", name),
			filepath.Join(filepath.Dir(exeDir), "boilerplate", "specs", name),
		}
		for _, path := range paths {
			if data, err := os.ReadFile(path); err == nil {
				return data, true
			}
		}
	}

	// Tentar relativo ao diretório de trabalho atual
	wd, err := os.Getwd()
	if err == nil {
		paths := []string{
			filepath.Join(wd, "boilerplate", "specs", name),
			filepath.Join(wd, "..", "boilerplate", "specs", name),
		}
		for _, path := range paths {
			if data, err := os.ReadFile(path); err == nil {
				return data, true
			}
		}
	}

	// Fallback: tentar caminho absoluto do projeto (para desenvolvimento)
	projectRoot := findProjectRoot()
	if projectRoot != "" {
		path := filepath.Join(projectRoot, "boilerplate", "specs", name)
		if data, err := os.ReadFile(path); err == nil {
			return data, true
		}
	}

	return nil, false
}

// findProjectRoot encontra a raiz do projeto procurando por go.mod
func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	dir := wd
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// GetAllTemplateNames retorna lista de nomes de templates disponíveis
func GetAllTemplateNames() []string {
	return []string{
		"00-global-context.spec.md",
		"00-architecture.spec.md",
		"00-stack.spec.md",
		"checklist.md",
		"template-default.spec.md",
	}
}
