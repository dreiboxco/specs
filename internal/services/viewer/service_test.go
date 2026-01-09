package viewer

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_View_Success(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar spec completa com requirements
	completeSpec := `# 01 Test Spec

## 2. Requisitos Funcionais

- **RF01 - Teste 1:**
  Teste

- **RF02 - Teste 2:**
  Teste

## 1. Contexto e Objetivo
Teste

## 3. Contratos e Interfaces
Teste

## 4. Fluxos e Estados
Teste

## 5. Dados
Teste

## 6. NFRs (Não Funcionais)
Teste

## 7. Guardrails
Teste

## 8. Critérios de Aceite
Teste

## 9. Testes
Teste

## 10. Migração / Rollback
Teste

## 11. Observações Operacionais
Teste

## 12. Abertos / Fora de Escopo
Teste

## Checklist Rápido (preencha antes de gerar código)
- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
`

	if err := fs.WriteFile(filepath.Join(specsDir, "01-test.spec.md"), []byte(completeSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.View(ViewOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.TotalSpecs != 1 {
		t.Errorf("esperado 1 spec, obtido %d", result.TotalSpecs)
	}
	if result.TotalRequirements != 2 {
		t.Errorf("esperado 2 requirements, obtido %d", result.TotalRequirements)
	}
	if result.SpecsComplete != 1 {
		t.Errorf("esperado 1 spec completa, obtido %d", result.SpecsComplete)
	}
	if len(result.Specs) != 1 {
		t.Errorf("esperado 1 spec na lista, obtido %d", len(result.Specs))
	}
	if result.Specs[0].Requirements != 2 {
		t.Errorf("esperado 2 requirements na spec, obtido %d", result.Specs[0].Requirements)
	}
}

func TestService_View_CountRequirements(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	// Testar contagem de requirements
	content := `## 2. Requisitos Funcionais

- **RF01 - Teste 1:**
  Teste

- **RF02 - Teste 2:**
  Teste

- **RF03 - Teste 3:**
  Teste
`

	count := service.countRequirements(content)
	if count != 3 {
		t.Errorf("esperado 3 requirements, obtido %d", count)
	}
}

func TestService_View_NoRequirements(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	// Testar spec sem seção de requisitos
	content := `# Test Spec

## 1. Contexto e Objetivo
Teste
`

	count := service.countRequirements(content)
	if count != 0 {
		t.Errorf("esperado 0 requirements, obtido %d", count)
	}
}

func TestService_View_InvalidPath(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	_, err := service.View(ViewOptions{
		Path: "/caminho/inexistente/12345",
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho inexistente")
	}
}

func TestService_View_NotDirectory(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")

	if err := fs.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	_, err := service.View(ViewOptions{
		Path: file,
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho que não é diretório")
	}
}

func TestService_View_ExcludeTemplates(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar spec de template (00-*)
	templateSpec := `# 00 Template Spec

## 2. Requisitos Funcionais
- **RF01 - Teste:** Teste

## 1. Contexto e Objetivo
Teste

## 3. Contratos e Interfaces
Teste

## 4. Fluxos e Estados
Teste

## 5. Dados
Teste

## 6. NFRs (Não Funcionais)
Teste

## 7. Guardrails
Teste

## 8. Critérios de Aceite
Teste

## 9. Testes
Teste

## 10. Migração / Rollback
Teste

## 11. Observações Operacionais
Teste

## 12. Abertos / Fora de Escopo
Teste

## Checklist Rápido (preencha antes de gerar código)
- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
`

	// Criar spec normal
	normalSpec := `# 01 Normal Spec

## 2. Requisitos Funcionais
- **RF01 - Teste:** Teste

## 1. Contexto e Objetivo
Teste

## 3. Contratos e Interfaces
Teste

## 4. Fluxos e Estados
Teste

## 5. Dados
Teste

## 6. NFRs (Não Funcionais)
Teste

## 7. Guardrails
Teste

## 8. Critérios de Aceite
Teste

## 9. Testes
Teste

## 10. Migração / Rollback
Teste

## 11. Observações Operacionais
Teste

## 12. Abertos / Fora de Escopo
Teste

## Checklist Rápido (preencha antes de gerar código)
- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
`

	// Criar template-default.spec.md
	templateDefault := `# Template Default`

	if err := fs.WriteFile(filepath.Join(specsDir, "00-template.spec.md"), []byte(templateSpec), 0644); err != nil {
		t.Fatalf("falha ao criar template spec: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "template-default.spec.md"), []byte(templateDefault), 0644); err != nil {
		t.Fatalf("falha ao criar template-default: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(specsDir, "01-normal.spec.md"), []byte(normalSpec), 0644); err != nil {
		t.Fatalf("falha ao criar normal spec: %v", err)
	}

	result, err := service.View(ViewOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Deve ter apenas 1 spec (excluindo templates)
	if result.TotalSpecs != 1 {
		t.Errorf("esperado 1 spec (excluindo templates), obtido %d", result.TotalSpecs)
	}
	if len(result.Specs) != 1 {
		t.Errorf("esperado 1 spec na lista, obtido %d", len(result.Specs))
	}
	if result.Specs[0].Name != "normal" {
		t.Errorf("esperado spec 'normal', obtido '%s'", result.Specs[0].Name)
	}
	// Verificar que templates não estão na lista
	for _, spec := range result.Specs {
		if strings.HasPrefix(spec.Number, "00") || spec.Name == "template-default" {
			t.Errorf("template não deveria estar na lista: %s", spec.Name)
		}
	}
}
