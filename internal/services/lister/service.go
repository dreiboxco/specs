package lister

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/services/validator"
)

// Service gerencia listagem de specs
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

// ListOptions contém opções para listagem
type ListOptions struct {
	Path       string
	Complete   bool
	Incomplete bool
	Errors     bool
}

// SpecInfo contém informações sobre uma spec
type SpecInfo struct {
	Path       string
	Number     string
	Name       string
	Status     string
	StatusIcon string
	Complete   bool
	HasErrors  bool
	Checklist  validator.ChecklistInfo
}

// ListResult contém resultado da listagem
type ListResult struct {
	Specs      []SpecInfo
	Total      int
	Complete   int
	Incomplete int
	WithErrors int
}

// List lista todas as specs com status
func (s *Service) List(opts ListOptions) (*ListResult, error) {
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

	// Validar cada spec e coletar informações
	result := &ListResult{
		Specs: make([]SpecInfo, 0, len(specFiles)),
	}

	for _, file := range specFiles {
		specInfo := s.getSpecInfo(file)
		result.Specs = append(result.Specs, specInfo)
		result.Total++

		if specInfo.HasErrors {
			result.WithErrors++
		} else if specInfo.Complete {
			result.Complete++
		} else {
			result.Incomplete++
		}
	}

	// Ordenar por numeração
	sort.Slice(result.Specs, func(i, j int) bool {
		return result.Specs[i].Number < result.Specs[j].Number
	})

	// Aplicar filtros
	if opts.Complete || opts.Incomplete || opts.Errors {
		filtered := make([]SpecInfo, 0)
		for _, spec := range result.Specs {
			if opts.Complete && spec.Complete {
				filtered = append(filtered, spec)
			} else if opts.Incomplete && !spec.Complete && !spec.HasErrors {
				filtered = append(filtered, spec)
			} else if opts.Errors && spec.HasErrors {
				filtered = append(filtered, spec)
			}
		}
		result.Specs = filtered
		// Recalcular contadores para specs filtradas
		result.Total = len(filtered)
		result.Complete = 0
		result.Incomplete = 0
		result.WithErrors = 0
		for _, spec := range filtered {
			if spec.HasErrors {
				result.WithErrors++
			} else if spec.Complete {
				result.Complete++
			} else {
				result.Incomplete++
			}
		}
	}

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

// getSpecInfo obtém informações sobre uma spec
func (s *Service) getSpecInfo(filePath string) SpecInfo {
	// Extrair numeração e nome do arquivo
	fileName := filepath.Base(filePath)
	nameWithoutExt := strings.TrimSuffix(fileName, ".spec.md")
	
	// Extrair numeração (ex.: "02-init" -> "02")
	parts := strings.SplitN(nameWithoutExt, "-", 2)
	number := ""
	name := nameWithoutExt
	if len(parts) >= 2 {
		number = parts[0]
		name = strings.Join(parts[1:], "-")
	} else {
		// Se não tem hífen, tentar extrair número do início
		for i, r := range nameWithoutExt {
			if r >= '0' && r <= '9' {
				number += string(r)
			} else {
				if i > 0 {
					name = nameWithoutExt[i:]
				}
				break
			}
		}
	}

	// Validar spec para obter status
	vr, err := s.validator.Validate(validator.ValidateOptions{
		Path: filePath,
	})
	
	var specInfo SpecInfo
	if err == nil && len(vr.Results) > 0 {
		validationResult := vr.Results[0]
		specInfo = SpecInfo{
			Path:      filePath,
			Number:    number,
			Name:      name,
			Complete:  validationResult.Complete,
			HasErrors: len(validationResult.Errors) > 0,
			Checklist: validationResult.Checklist,
		}

		// Determinar status e ícone
		if specInfo.HasErrors {
			specInfo.Status = "Erro"
			specInfo.StatusIcon = "❌"
		} else if specInfo.Complete {
			specInfo.Status = "Completa"
			specInfo.StatusIcon = "✅"
		} else {
			marked := specInfo.Checklist.MarkedCount
			specInfo.Status = fmt.Sprintf("Incompleta (%d/6)", marked)
			specInfo.StatusIcon = "⚠️"
		}
	} else {
		// Fallback se validação falhar
		specInfo = SpecInfo{
			Path:       filePath,
			Number:     number,
			Name:       name,
			Status:     "Erro",
			StatusIcon: "❌",
			HasErrors:  true,
		}
	}

	return specInfo
}
