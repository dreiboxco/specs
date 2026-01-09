package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/dreibox/specs/internal/adapters"
)

// Service gerencia validação de specs
type Service struct {
	fs adapters.FileSystem
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{fs: fs}
}

// RequiredSections são as seções obrigatórias de uma spec
var RequiredSections = []string{
	"Contexto e Objetivo",
	"Requisitos Funcionais",
	"Contratos e Interfaces",
	"Fluxos e Estados",
	"Dados",
	"NFRs",
	"Guardrails",
	"Critérios de Aceite",
	"Testes",
	"Migração",
	"Observações Operacionais",
	"Abertos",
}

// ValidateOptions contém opções para validação
type ValidateOptions struct {
	Path string // Caminho para arquivo ou diretório
}

// ValidationResult contém resultado da validação de uma spec
type ValidationResult struct {
	Path      string
	Valid     bool
	Complete  bool // Todos os itens do checklist marcados
	Errors    []string
	Warnings  []string
	Checklist ChecklistInfo
}

// ChecklistInfo contém informações sobre o checklist
type ChecklistInfo struct {
	Found      bool
	ItemCount  int
	MarkedCount int
	ValidFormat bool
}

// ValidateResult contém resultado agregado de validação
type ValidateResult struct {
	Results      []ValidationResult
	Total        int
	Complete     int
	Incomplete   int
	WithErrors   int
}

// Validate valida um arquivo ou diretório de specs
func (s *Service) Validate(opts ValidateOptions) (*ValidateResult, error) {
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

	// Verificar se é arquivo ou diretório
	stat, err := s.fs.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter informações do caminho: %w", err)
	}

	var specFiles []string
	if stat.IsDir() {
		// Listar todos os arquivos .spec.md recursivamente
		specFiles, err = s.findSpecFiles(path)
		if err != nil {
			return nil, fmt.Errorf("falha ao listar arquivos: %w", err)
		}
	} else {
		// Verificar se é arquivo .spec.md
		if !strings.HasSuffix(path, ".spec.md") {
			return nil, fmt.Errorf("arquivo deve ter extensão .spec.md: %s", path)
		}
		specFiles = []string{path}
	}

	// Validar cada arquivo
	result := &ValidateResult{
		Results: make([]ValidationResult, 0, len(specFiles)),
	}

	for _, file := range specFiles {
		vr := s.validateFile(file)
		result.Results = append(result.Results, vr)
		result.Total++

		if len(vr.Errors) > 0 {
			result.WithErrors++
		} else if vr.Complete {
			result.Complete++
		} else {
			result.Incomplete++
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

// validateFile valida um arquivo de spec
func (s *Service) validateFile(path string) ValidationResult {
	result := ValidationResult{
		Path:   path,
		Valid:  true,
		Errors: []string{},
		Warnings: []string{},
	}

	// Ler arquivo
	data, err := s.fs.ReadFile(path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("falha ao ler arquivo: %v", err))
		return result
	}

	// Verificar encoding UTF-8
	if !utf8.Valid(data) {
		result.Valid = false
		result.Errors = append(result.Errors, "arquivo não está em UTF-8")
		return result
	}

	content := string(data)

	// Verificar se arquivo não está vazio
	if strings.TrimSpace(content) == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "arquivo está vazio")
		return result
	}

	// Validar estrutura básica
	if err := s.validateStructure(content); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, err.Error())
	}

	// Validar seções obrigatórias
	missingSections := s.validateRequiredSections(content)
	if len(missingSections) > 0 {
		result.Valid = false
		for _, section := range missingSections {
			result.Errors = append(result.Errors, fmt.Sprintf("seção '%s' faltando", section))
		}
	}

	// Validar checklist
	checklistInfo := s.validateChecklist(content)
	result.Checklist = checklistInfo

	if !checklistInfo.Found {
		result.Valid = false
		result.Errors = append(result.Errors, "checklist não encontrado")
	} else if !checklistInfo.ValidFormat {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("checklist com formato inválido (esperado 6 itens, encontrado %d)", checklistInfo.ItemCount))
	} else if checklistInfo.MarkedCount < 6 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("checklist incompleto (%d/6 itens)", checklistInfo.MarkedCount))
	}

	// Determinar se está completa (sem erros e checklist completo)
	result.Complete = result.Valid && checklistInfo.MarkedCount == 6

	return result
}

// validateStructure valida estrutura básica do arquivo
func (s *Service) validateStructure(content string) error {
	lines := strings.Split(content, "\n")

	// Verificar se começa com título principal (#)
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "# ") {
		return fmt.Errorf("arquivo deve começar com título principal (#)")
	}

	// Verificar hierarquia de títulos (não pular níveis)
	prevLevel := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for _, r := range trimmed {
				if r == '#' {
					level++
				} else {
					break
				}
			}
			if level > 0 && level <= 6 {
				if prevLevel > 0 && level > prevLevel+1 {
					return fmt.Errorf("hierarquia de títulos inválida: pulou do nível %d para %d", prevLevel, level)
				}
				prevLevel = level
			}
		}
	}

	return nil
}

// validateRequiredSections verifica se todas as seções obrigatórias estão presentes
func (s *Service) validateRequiredSections(content string) []string {
	missing := []string{}
	found := make(map[string]bool)

	// Procurar por seções (## N. Nome da Seção ou ## Nome da Seção)
	sectionRegex := regexp.MustCompile(`^##\s+(?:\d+\.\s*)?(.+)$`)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		matches := sectionRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			sectionName := strings.TrimSpace(matches[1])
			// Normalizar nome da seção (remover " / " e variações)
			normalized := s.normalizeSectionName(sectionName)
			found[normalized] = true
		}
	}

	// Verificar seções obrigatórias
	for _, required := range RequiredSections {
		if !found[required] {
			missing = append(missing, required)
		}
	}

	return missing
}

// normalizeSectionName normaliza nome de seção para comparação
func (s *Service) normalizeSectionName(name string) string {
	// Remover variações comuns
	name = strings.TrimSpace(name)
	
	// Mapear variações para nomes padrão
	variations := map[string]string{
		"Abertos / Fora de Escopo": "Abertos",
		"Abertos/Fora de Escopo":   "Abertos",
		"Migração / Rollback":      "Migração",
		"Migração/Rollback":         "Migração",
		"NFRs (Não Funcionais)":     "NFRs",
		"NFRs":                      "NFRs",
	}

	if normalized, ok := variations[name]; ok {
		return normalized
	}

	// Verificar se contém palavras-chave das seções obrigatórias
	if strings.Contains(name, "Abertos") || strings.Contains(name, "Fora de Escopo") {
		return "Abertos"
	}
	if strings.Contains(name, "Migração") || strings.Contains(name, "Rollback") {
		return "Migração"
	}
	if strings.HasPrefix(name, "NFRs") || strings.Contains(name, "Não Funcionais") {
		return "NFRs"
	}
	if strings.Contains(name, "Contexto") && strings.Contains(name, "Objetivo") {
		return "Contexto e Objetivo"
	}
	if strings.Contains(name, "Requisitos Funcionais") {
		return "Requisitos Funcionais"
	}
	if strings.Contains(name, "Contratos") && strings.Contains(name, "Interfaces") {
		return "Contratos e Interfaces"
	}
	if strings.Contains(name, "Fluxos") && strings.Contains(name, "Estados") {
		return "Fluxos e Estados"
	}
	if name == "Dados" {
		return "Dados"
	}
	if strings.Contains(name, "Guardrails") {
		return "Guardrails"
	}
	if strings.Contains(name, "Critérios") && strings.Contains(name, "Aceite") {
		return "Critérios de Aceite"
	}
	if name == "Testes" {
		return "Testes"
	}
	if strings.Contains(name, "Observações Operacionais") {
		return "Observações Operacionais"
	}

	return name
}

// validateChecklist valida o checklist da spec
func (s *Service) validateChecklist(content string) ChecklistInfo {
	info := ChecklistInfo{
		Found:       false,
		ItemCount:   0,
		MarkedCount: 0,
		ValidFormat: false,
	}

	// Procurar por seção "Checklist" (Checklist Rápido) após seção "Abertos"
	lines := strings.Split(content, "\n")
	foundAbertos := false
	checklistSectionFound := false
	checklistStart := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Verificar se encontrou seção "Abertos"
		if strings.HasPrefix(trimmed, "##") && strings.Contains(strings.ToLower(trimmed), "abertos") {
			foundAbertos = true
		}
		// Após "Abertos", procurar por seção "Checklist"
		if foundAbertos && strings.HasPrefix(trimmed, "##") && strings.Contains(strings.ToLower(trimmed), "checklist") {
			checklistSectionFound = true
			continue
		}
		// Procurar início do checklist (linha com "- [") após seção "Checklist"
		if checklistSectionFound && strings.HasPrefix(trimmed, "- [") {
			checklistStart = i
			info.Found = true
			break
		}
	}

	if !info.Found {
		return info
	}

	// Contar itens do checklist
	itemRegex := regexp.MustCompile(`^-\s+\[([ x])\]\s+(.+)$`)
	for i := checklistStart; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		
		// Parar se encontrar próxima seção (##) ou linha vazia seguida de seção
		if strings.HasPrefix(trimmed, "##") {
			break
		}

		// Verificar se é item de checklist
		matches := itemRegex.FindStringSubmatch(trimmed)
		if len(matches) > 0 {
			info.ItemCount++
			if strings.TrimSpace(matches[1]) == "x" {
				info.MarkedCount++
			}
		} else if trimmed != "" && !strings.HasPrefix(trimmed, "-") {
			// Se encontrou linha não vazia que não é item de checklist, pode ser fim do checklist
			// Mas continuar para pegar todos os itens consecutivos
		}
	}

	// Validar formato (deve ter exatamente 6 itens)
	info.ValidFormat = info.ItemCount == 6

	return info
}
