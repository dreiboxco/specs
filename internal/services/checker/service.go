package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
)

// Service gerencia verificação de consistência estrutural
type Service struct {
	fs adapters.FileSystem
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{fs: fs}
}

// CheckOptions contém opções para verificação
type CheckOptions struct {
	Path string
}

// Problem representa um problema encontrado
type Problem struct {
	Category string // "Numeração", "Links", "Formato", etc.
	Severity string // "error", "warning"
	File     string
	Line     int
	Message  string
}

// CheckResult contém resultado da verificação
type CheckResult struct {
	TotalSpecs int
	Problems   []Problem
	Summary    map[string]int // categoria -> quantidade
}

// Check verifica consistência estrutural de specs
func (s *Service) Check(opts CheckOptions) (*CheckResult, error) {
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

	result := &CheckResult{
		TotalSpecs: len(specFiles),
		Problems:   []Problem{},
		Summary:    make(map[string]int),
	}

	// Construir mapeamento de specs (número -> arquivos)
	specMap := s.buildSpecMap(specFiles, path)

	// Validar numeração
	s.checkNumbering(specFiles, path, result, specMap)

	// Validar formato de nomes
	s.checkFileNameFormat(specFiles, path, result)

	// Validar links
	s.checkLinks(specFiles, path, result, specMap)

	// Detectar specs órfãs
	s.checkOrphanedSpecs(specFiles, path, result, specMap)

	// Contar problemas por categoria
	for _, p := range result.Problems {
		result.Summary[p.Category]++
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

// buildSpecMap constrói mapeamento de número para arquivos
func (s *Service) buildSpecMap(files []string, basePath string) map[string][]string {
	specMap := make(map[string][]string)
	for _, file := range files {
		fileName := filepath.Base(file)
		number := s.extractNumber(fileName)
		if number != "" {
			relPath, _ := filepath.Rel(basePath, file)
			if relPath == "" || relPath == "." {
				relPath = fileName
			}
			specMap[number] = append(specMap[number], relPath)
		}
	}
	return specMap
}

// extractNumber extrai número do nome do arquivo
func (s *Service) extractNumber(fileName string) string {
	// Formato esperado: {numero}-{nome}.spec.md
	parts := strings.SplitN(fileName, "-", 2)
	if len(parts) < 2 {
		return ""
	}
	number := parts[0]
	// Validar que é número zero-padded de 2 dígitos
	if matched, _ := regexp.MatchString(`^\d{2}$`, number); matched {
		return number
	}
	return ""
}

// checkNumbering verifica numeração sequencial
func (s *Service) checkNumbering(files []string, basePath string, result *CheckResult, specMap map[string][]string) {
	// Extrair números únicos e ordenar
	numbers := make(map[string]bool)
	for number := range specMap {
		numbers[number] = true
	}

	// Converter para slice e ordenar
	numberList := make([]string, 0, len(numbers))
	for number := range numbers {
		numberList = append(numberList, number)
	}
	sort.Strings(numberList)

	// Verificar duplicatas
	for number, files := range specMap {
		if len(files) > 1 {
			for _, file := range files {
				result.Problems = append(result.Problems, Problem{
					Category: "Numeração",
					Severity: "error",
					File:     file,
					Message:  fmt.Sprintf("Numeração duplicada: %s usado em %d arquivo(s)", number, len(files)),
				})
			}
		}
	}

	// Verificar gaps
	if len(numberList) > 0 {
		// Converter números para inteiros para verificar sequência
		firstNum := s.numberToInt(numberList[0])
		lastNum := s.numberToInt(numberList[len(numberList)-1])

		for i := firstNum; i <= lastNum; i++ {
			numStr := fmt.Sprintf("%02d", i)
			if !numbers[numStr] {
				result.Problems = append(result.Problems, Problem{
					Category: "Numeração",
					Severity: "warning",
					Message:  fmt.Sprintf("Gap detectado - falta %s", numStr),
				})
			}
		}
	}
}

// numberToInt converte string de número para int
func (s *Service) numberToInt(numStr string) int {
	var num int
	fmt.Sscanf(numStr, "%d", &num)
	return num
}

// checkFileNameFormat verifica formato de nomes de arquivos
func (s *Service) checkFileNameFormat(files []string, basePath string, result *CheckResult) {
	fileNameRegex := regexp.MustCompile(`^(\d{2})-(.+?)\.spec\.md$`)
	for _, file := range files {
		fileName := filepath.Base(file)
		relPath, _ := filepath.Rel(basePath, file)
		if relPath == "" || relPath == "." {
			relPath = fileName
		}

		if !fileNameRegex.MatchString(fileName) {
			result.Problems = append(result.Problems, Problem{
				Category: "Formato",
				Severity: "error",
				File:     relPath,
				Message:  fmt.Sprintf("Nome não segue padrão {numero}-{nome}.spec.md"),
			})
		} else {
			matches := fileNameRegex.FindStringSubmatch(fileName)
			if len(matches) >= 3 {
				name := matches[2]
				if name == "" {
					result.Problems = append(result.Problems, Problem{
						Category: "Formato",
						Severity: "error",
						File:     relPath,
						Message:  "Nome descritivo está vazio (apenas número)",
					})
				}
			}
		}
	}
}

// checkLinks verifica links internos
func (s *Service) checkLinks(files []string, basePath string, result *CheckResult, specMap map[string][]string) {
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	
	for _, file := range files {
		data, err := s.fs.ReadFile(file)
		if err != nil {
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		relPath, _ := filepath.Rel(basePath, file)
		if relPath == "" || relPath == "." {
			relPath = filepath.Base(file)
		}

		for lineNum, line := range lines {
			matches := linkRegex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					linkPath := match[2]
					// Verificar se é link interno para spec
					if strings.HasSuffix(linkPath, ".spec.md") {
						// Extrair nome do arquivo
						linkFile := filepath.Base(linkPath)
						// Verificar se arquivo existe no mapeamento
						linkNumber := s.extractNumber(linkFile)
						if linkNumber != "" {
							if _, exists := specMap[linkNumber]; !exists {
								result.Problems = append(result.Problems, Problem{
									Category: "Links",
									Severity: "error",
									File:     relPath,
									Line:     lineNum + 1,
									Message:  fmt.Sprintf("Link para '%s' não encontrado", linkFile),
								})
							}
						} else {
							// Verificar se arquivo existe no diretório
							fullPath := filepath.Join(basePath, linkPath)
							if !s.fs.Exists(fullPath) {
								result.Problems = append(result.Problems, Problem{
									Category: "Links",
									Severity: "error",
									File:     relPath,
									Line:     lineNum + 1,
									Message:  fmt.Sprintf("Link para '%s' não encontrado", linkPath),
								})
							}
						}
					}
				}
			}
		}
	}
}

// checkOrphanedSpecs detecta specs órfãs (referenciadas mas não existem)
func (s *Service) checkOrphanedSpecs(files []string, basePath string, result *CheckResult, specMap map[string][]string) {
	// Construir índice de arquivos existentes
	existingFiles := make(map[string]bool)
	for _, file := range files {
		fileName := filepath.Base(file)
		existingFiles[fileName] = true
	}

	// Construir índice de referências
	referencedSpecs := make(map[string]bool)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	for _, file := range files {
		data, err := s.fs.ReadFile(file)
		if err != nil {
			continue
		}

		content := string(data)
		matches := linkRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				linkPath := match[2]
				if strings.HasSuffix(linkPath, ".spec.md") {
					linkFile := filepath.Base(linkPath)
					referencedSpecs[linkFile] = true
				}
			}
		}
	}

	// Verificar se specs referenciadas existem
	for refFile := range referencedSpecs {
		if !existingFiles[refFile] {
			// Verificar se arquivo existe fisicamente (pode estar em subdiretório)
			fullPath := filepath.Join(basePath, refFile)
			if !s.fs.Exists(fullPath) {
				result.Problems = append(result.Problems, Problem{
					Category: "Órfãs",
					Severity: "error",
					Message:  fmt.Sprintf("Spec referenciada mas não existe: %s", refFile),
				})
			}
		}
	}
}
