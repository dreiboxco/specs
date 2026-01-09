package update

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/templates"
)

// Service gerencia operações de atualização de templates
type Service struct {
	fs adapters.FileSystem
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{fs: fs}
}

// UpdateOptions contém opções para atualização
type UpdateOptions struct {
	TargetDir  string
	DryRun     bool
	Force      bool
	NoBackup   bool
	Merge      bool
}

// UpdateResult contém resultado da atualização
type UpdateResult struct {
	BackupDir          string
	FilesUpdated       []string
	FilesSkipped       []string
	CursorRulesUpdated bool
	CursorRulesMerged  bool
	HasCustomizations  bool
}

// Update executa atualização de templates
func (s *Service) Update(opts UpdateOptions) (*UpdateResult, error) {
	// Determinar diretório alvo
	targetDir := opts.TargetDir
	if targetDir == "" {
		wd, err := s.fs.Getwd()
		if err != nil {
			return nil, fmt.Errorf("falha ao obter diretório atual: %w", err)
		}
		targetDir = wd
	}

	// Verificar se é projeto SDD
	if !s.isSDDProject(targetDir) {
		return nil, fmt.Errorf("diretório não contém projeto SDD válido")
	}

	// Verificar permissões
	if err := s.checkWritePermissions(targetDir); err != nil {
		return nil, fmt.Errorf("sem permissão de escrita: %w", err)
	}

	result := &UpdateResult{
		FilesUpdated: []string{},
		FilesSkipped: []string{},
	}

	specsDir := filepath.Join(targetDir, "specs")

	// Criar backup se não for dry-run e não estiver desabilitado
	var backupDir string
	if !opts.DryRun && !opts.NoBackup {
		var err error
		backupDir, err = s.createBackup(targetDir, specsDir)
		if err != nil {
			return nil, fmt.Errorf("falha ao criar backup: %w", err)
		}
		result.BackupDir = backupDir
	}

	// Atualizar templates estáticos
	if err := s.updateStaticTemplates(specsDir, result, opts); err != nil {
		return nil, err
	}

	// Atualizar .cursorrules
	if err := s.updateCursorRules(targetDir, result, opts); err != nil {
		return nil, err
	}

	// Limpar backups antigos
	if !opts.DryRun && !opts.NoBackup {
		if err := s.cleanupOldBackups(targetDir); err != nil {
			// Não falhar se limpeza falhar
			_ = err
		}
	}

	return result, nil
}

// isSDDProject verifica se diretório contém projeto SDD válido
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
		if s.fs.Exists(testFile) {
			_ = os.Remove(testFile)
		}
	}()

	if err := s.fs.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("sem permissão de escrita: %w", err)
	}

	return nil
}

// createBackup cria backup dos arquivos que serão atualizados
func (s *Service) createBackup(targetDir, specsDir string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(targetDir, ".specs-backup", timestamp)

	// Criar diretório de backup
	if err := s.fs.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	// Arquivos para fazer backup
	filesToBackup := []string{
		filepath.Join(specsDir, "checklist.md"),
		filepath.Join(specsDir, "template-default.spec.md"),
		filepath.Join(targetDir, ".cursorrules"),
	}

	backupSpecsDir := filepath.Join(backupDir, "specs")
	if err := s.fs.MkdirAll(backupSpecsDir, 0755); err != nil {
		return "", err
	}

	for _, file := range filesToBackup {
		if !s.fs.Exists(file) {
			continue
		}

		data, err := s.fs.ReadFile(file)
		if err != nil {
			continue
		}

		var backupPath string
		if strings.HasPrefix(file, specsDir) {
			relPath, _ := filepath.Rel(specsDir, file)
			backupPath = filepath.Join(backupSpecsDir, relPath)
		} else {
			backupPath = filepath.Join(backupDir, filepath.Base(file))
		}

		if err := s.fs.WriteFile(backupPath, data, 0644); err != nil {
			return "", fmt.Errorf("falha ao criar backup de %s: %w", file, err)
		}
	}

	return backupDir, nil
}

// updateStaticTemplates atualiza templates estáticos
func (s *Service) updateStaticTemplates(specsDir string, result *UpdateResult, opts UpdateOptions) error {
	templatesToUpdate := []string{
		"checklist.md",
		"template-default.spec.md",
	}

	for _, name := range templatesToUpdate {
		template, exists := templates.GetTemplate(name)
		if !exists {
			result.FilesSkipped = append(result.FilesSkipped, name)
			continue
		}

		targetPath := filepath.Join(specsDir, name)

		if opts.DryRun {
			result.FilesUpdated = append(result.FilesUpdated, name)
			continue
		}

		if err := s.fs.WriteFile(targetPath, template, 0644); err != nil {
			return fmt.Errorf("falha ao atualizar %s: %w", name, err)
		}

		result.FilesUpdated = append(result.FilesUpdated, name)
	}

	return nil
}

// updateCursorRules atualiza .cursorrules com detecção de personalizações
func (s *Service) updateCursorRules(targetDir string, result *UpdateResult, opts UpdateOptions) error {
	cursorRulesPath := filepath.Join(targetDir, ".cursorrules")

	// Obter template do boilerplate
	boilerplateContent, err := templates.GetCursorRulesTemplate()
	if err != nil {
		return fmt.Errorf("falha ao obter template .cursorrules: %w", err)
	}

	// Se arquivo não existe, criar diretamente
	if !s.fs.Exists(cursorRulesPath) {
		if opts.DryRun {
			result.FilesUpdated = append(result.FilesUpdated, ".cursorrules")
			return nil
		}

		if err := s.fs.WriteFile(cursorRulesPath, boilerplateContent, 0644); err != nil {
			return fmt.Errorf("falha ao criar .cursorrules: %w", err)
		}

		result.CursorRulesUpdated = true
		return nil
	}

	// Ler conteúdo atual
	currentContent, err := s.fs.ReadFile(cursorRulesPath)
	if err != nil {
		return fmt.Errorf("falha ao ler .cursorrules: %w", err)
	}

	// Detectar personalizações
	hasCustomizations := s.detectCustomizations(currentContent, boilerplateContent)
	result.HasCustomizations = hasCustomizations

	if opts.DryRun {
		if hasCustomizations {
			result.FilesSkipped = append(result.FilesSkipped, ".cursorrules (merge necessário)")
		} else {
			result.FilesUpdated = append(result.FilesUpdated, ".cursorrules")
		}
		return nil
	}

	// Se não há personalizações, atualizar diretamente
	if !hasCustomizations {
		if err := s.fs.WriteFile(cursorRulesPath, boilerplateContent, 0644); err != nil {
			return fmt.Errorf("falha ao atualizar .cursorrules: %w", err)
		}
		result.CursorRulesUpdated = true
		return nil
	}

	// Há personalizações - criar .cursorrules-updated
	updatedPath := filepath.Join(targetDir, ".cursorrules-updated")
	if err := s.fs.WriteFile(updatedPath, boilerplateContent, 0644); err != nil {
		return fmt.Errorf("falha ao criar .cursorrules-updated: %w", err)
	}

	// Tentar merge automático se solicitado
	if opts.Merge {
		mergedContent, err := s.mergeCursorRules(currentContent, boilerplateContent)
		if err == nil {
			mergedPath := filepath.Join(targetDir, ".cursorrules-merged")
			if err := s.fs.WriteFile(mergedPath, mergedContent, 0644); err == nil {
				result.CursorRulesMerged = true
			}
		}
	}

	return nil
}

// detectCustomizations detecta se há personalizações no .cursorrules
func (s *Service) detectCustomizations(current, boilerplate []byte) bool {
	currentStr := string(current)
	boilerplateStr := string(boilerplate)

	// Se são idênticos, não há personalizações
	if currentStr == boilerplateStr {
		return false
	}

	// Dividir em seções
	currentSections := s.extractSections(currentStr)
	boilerplateSections := s.extractSections(boilerplateStr)

	// Verificar seções novas (não existem no boilerplate)
	for section := range currentSections {
		if _, exists := boilerplateSections[section]; !exists {
			return true // Seção personalizada encontrada
		}
	}

	// Verificar seções modificadas (mais de 3 linhas de diferença)
	for section, currentContent := range currentSections {
		if boilerplateContent, exists := boilerplateSections[section]; exists {
			currentLines := strings.Split(strings.TrimSpace(currentContent), "\n")
			boilerplateLines := strings.Split(strings.TrimSpace(boilerplateContent), "\n")

			diff := len(currentLines) - len(boilerplateLines)
			if diff < 0 {
				diff = -diff
			}

			// Se diferença significativa (mais de 3 linhas), considerar personalizado
			if diff > 3 {
				return true
			}

			// Verificar se conteúdo é significativamente diferente
			if strings.TrimSpace(currentContent) != strings.TrimSpace(boilerplateContent) {
				// Contar linhas diferentes
				diffCount := 0
				maxLen := len(currentLines)
				if len(boilerplateLines) > maxLen {
					maxLen = len(boilerplateLines)
				}
				for i := 0; i < maxLen; i++ {
					currentLine := ""
					boilerplateLine := ""
					if i < len(currentLines) {
						currentLine = strings.TrimSpace(currentLines[i])
					}
					if i < len(boilerplateLines) {
						boilerplateLine = strings.TrimSpace(boilerplateLines[i])
					}
					if currentLine != boilerplateLine {
						diffCount++
					}
				}
				if diffCount > 3 {
					return true
				}
			}
		}
	}

	return false
}

// extractSections extrai seções do arquivo
func (s *Service) extractSections(content string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentSection string
	var currentContent []string

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			// Salvar seção anterior
			if currentSection != "" {
				sections[currentSection] = strings.Join(currentContent, "\n")
			}
			// Nova seção
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			currentContent = []string{}
		} else if currentSection != "" {
			currentContent = append(currentContent, line)
		}
	}

	// Salvar última seção
	if currentSection != "" {
		sections[currentSection] = strings.Join(currentContent, "\n")
	}

	return sections
}

// mergeCursorRules tenta fazer merge automático de .cursorrules
func (s *Service) mergeCursorRules(current, boilerplate []byte) ([]byte, error) {
	currentStr := string(current)
	boilerplateStr := string(boilerplate)

	currentSections := s.extractSections(currentStr)
	boilerplateSections := s.extractSections(boilerplateStr)

	// Estratégia de merge:
	// 1. Preservar seções personalizadas do projeto
	// 2. Atualizar seções que existem no boilerplate mas não foram personalizadas
	// 3. Adicionar novas seções do boilerplate

	mergedSections := make(map[string]string)

	// Primeiro, adicionar todas as seções do boilerplate (versão atualizada)
	for section, content := range boilerplateSections {
		mergedSections[section] = content
	}

	// Depois, preservar seções personalizadas do projeto
	for section, content := range currentSections {
		// Se seção não existe no boilerplate, é personalizada - preservar
		if _, exists := boilerplateSections[section]; !exists {
			mergedSections[section] = content
		} else {
			// Seção existe em ambos - verificar se foi personalizada
			currentLines := strings.Split(strings.TrimSpace(content), "\n")
			boilerplateLines := strings.Split(strings.TrimSpace(boilerplateSections[section]), "\n")

			diffCount := 0
			maxLen := len(currentLines)
			if len(boilerplateLines) > maxLen {
				maxLen = len(boilerplateLines)
			}
			for i := 0; i < maxLen; i++ {
				currentLine := ""
				boilerplateLine := ""
				if i < len(currentLines) {
					currentLine = strings.TrimSpace(currentLines[i])
				}
				if i < len(boilerplateLines) {
					boilerplateLine = strings.TrimSpace(boilerplateLines[i])
				}
				if currentLine != boilerplateLine {
					diffCount++
				}
			}

			// Se diferença significativa, preservar versão personalizada
			if diffCount > 3 {
				mergedSections[section] = content
			}
			// Caso contrário, usar versão do boilerplate (já está em mergedSections)
		}
	}

	// Reconstruir arquivo
	var result []string
	result = append(result, "# Cursor Rules - Spec Driven Development")
	result = append(result, "")

	// Ordenar seções (manter ordem do boilerplate, depois adicionar personalizadas)
	sectionOrder := []string{}
	seen := make(map[string]bool)

	// Adicionar seções do boilerplate na ordem original
	boilerplateLines := strings.Split(boilerplateStr, "\n")
	for _, line := range boilerplateLines {
		if strings.HasPrefix(line, "## ") {
			section := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			if !seen[section] {
				sectionOrder = append(sectionOrder, section)
				seen[section] = true
			}
		}
	}

	// Adicionar seções personalizadas que não estão no boilerplate
	for section := range mergedSections {
		if !seen[section] {
			sectionOrder = append(sectionOrder, section)
		}
	}

	// Construir conteúdo
	for _, section := range sectionOrder {
		if content, exists := mergedSections[section]; exists {
			result = append(result, "")
			result = append(result, "## "+section)
			result = append(result, content)
		}
	}

	return []byte(strings.Join(result, "\n")), nil
}

// cleanupOldBackups remove backups antigos, mantendo apenas os últimos 5
func (s *Service) cleanupOldBackups(targetDir string) error {
	backupBaseDir := filepath.Join(targetDir, ".specs-backup")
	if !s.fs.Exists(backupBaseDir) {
		return nil
	}

	// Listar diretórios de backup
	var backups []string
	entries, err := s.fs.ReadDir(backupBaseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	// Se há mais de 5 backups, remover os mais antigos
	if len(backups) > 5 {
		// Ordenar por timestamp (nome do diretório)
		// Remover os mais antigos
		toRemove := len(backups) - 5
		for i := 0; i < toRemove; i++ {
			backupPath := filepath.Join(backupBaseDir, backups[i])
			_ = os.RemoveAll(backupPath)
		}
	}

	return nil
}
