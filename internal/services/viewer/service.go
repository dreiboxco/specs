package viewer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/services/validator"
)

// Service gerencia visualização de dashboard
type Service struct {
	fs        adapters.FileSystem
	validator *validator.Service
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{
		fs:        fs,
		validator: validator.NewService(fs),
	}
}

// ViewOptions contém opções para visualização
type ViewOptions struct {
	Path string
}

// SpecStats contém estatísticas de uma spec
type SpecStats struct {
	Path         string
	Number       string
	Name         string
	Requirements int
	Progress     float64 // 0.0 a 1.0
	Complete     bool
	MarkedItems  int
	TotalItems   int
}

// DashboardResult contém resultado do dashboard
type DashboardResult struct {
	TotalSpecs        int
	TotalRequirements  int
	SpecsComplete      int
	SpecsInProgress    int
	OverallProgress    float64 // 0.0 a 1.0
	OverallProgressStr string  // "X/Y (Z%)"
	Specs              []SpecStats
}

// View gera dashboard de visualização
func (s *Service) View(opts ViewOptions) (*DashboardResult, error) {
	// Determinar caminho
	path := opts.Path
	if path == "" {
		wd, err := s.fs.Getwd()
		if err != nil {
			return nil, fmt.Errorf("falha ao obter diretório atual: %w", err)
		}
		path = filepath.Join(wd, "specs")
	}

	// Verificar se caminho existe
	if !s.fs.Exists(path) {
		return nil, fmt.Errorf("caminho não existe: %s", path)
	}

	// Verificar se é diretório
	stat, err := s.fs.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter informações do caminho: %w", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("caminho não é diretório: %s", path)
	}

	// Listar todos os arquivos .spec.md
	specFiles, err := s.findSpecFiles(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar arquivos: %w", err)
	}

	result := &DashboardResult{
		TotalSpecs: len(specFiles),
		Specs:      make([]SpecStats, 0, len(specFiles)),
	}

	// Processar cada spec
	totalMarkedItems := 0
	totalPossibleItems := 0

	for _, file := range specFiles {
		stats := s.getSpecStats(file, path)
		result.Specs = append(result.Specs, stats)

		result.TotalRequirements += stats.Requirements
		totalMarkedItems += stats.MarkedItems
		totalPossibleItems += stats.TotalItems

		if stats.Complete {
			result.SpecsComplete++
		} else {
			result.SpecsInProgress++
		}
	}

	// Calcular progresso geral
	if totalPossibleItems > 0 {
		result.OverallProgress = float64(totalMarkedItems) / float64(totalPossibleItems)
		percent := int(result.OverallProgress * 100)
		result.OverallProgressStr = fmt.Sprintf("%d/%d (%d%% complete)", totalMarkedItems, totalPossibleItems, percent)
	} else {
		result.OverallProgressStr = "0/0 (0% complete)"
	}

	// Ordenar specs por numeração
	sort.Slice(result.Specs, func(i, j int) bool {
		return result.Specs[i].Number < result.Specs[j].Number
	})

	return result, nil
}

// findSpecFiles encontra todos os arquivos .spec.md recursivamente
func (s *Service) findSpecFiles(root string) ([]string, error) {
	var files []string
	err := s.fs.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".spec.md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// getSpecStats obtém estatísticas de uma spec
func (s *Service) getSpecStats(filePath string, basePath string) SpecStats {
	// Extrair numeração e nome
	fileName := filepath.Base(filePath)
	nameWithoutExt := strings.TrimSuffix(fileName, ".spec.md")
	
	parts := strings.SplitN(nameWithoutExt, "-", 2)
	number := ""
	name := nameWithoutExt
	if len(parts) >= 2 {
		number = parts[0]
		name = strings.Join(parts[1:], "-")
	}

	stats := SpecStats{
		Path:   filePath,
		Number: number,
		Name:   name,
		TotalItems: 6, // Checklist sempre tem 6 itens
	}

	// Ler arquivo
	data, err := s.fs.ReadFile(filePath)
	if err != nil {
		return stats
	}

	content := string(data)

	// Contar requirements
	stats.Requirements = s.countRequirements(content)

	// Validar spec para obter progresso
	vr, err := s.validator.Validate(validator.ValidateOptions{
		Path: filePath,
	})
	if err == nil && len(vr.Results) > 0 {
		validationResult := vr.Results[0]
		stats.MarkedItems = validationResult.Checklist.MarkedCount
		stats.Complete = validationResult.Complete
	}

	// Calcular progresso
	if stats.TotalItems > 0 {
		stats.Progress = float64(stats.MarkedItems) / float64(stats.TotalItems)
	}

	return stats
}

// countRequirements conta requirements na seção "Requisitos Funcionais"
func (s *Service) countRequirements(content string) int {
	// Procurar seção "Requisitos Funcionais"
	lines := strings.Split(content, "\n")
	inSection := false
	count := 0

	// Regex para detectar RF01, RF02, etc.
	rfRegex := regexp.MustCompile(`^-\s+\*\*RF\d+`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Verificar se entrou na seção
		if strings.HasPrefix(trimmed, "##") {
			if strings.Contains(trimmed, "Requisitos Funcionais") {
				inSection = true
				continue
			}
			// Se encontrou outra seção, parar
			if inSection {
				break
			}
		}

		// Se está na seção, contar requirements
		if inSection {
			if rfRegex.MatchString(trimmed) {
				count++
			}
		}
	}

	return count
}
