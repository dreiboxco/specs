package validator

import (
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_Validate_ValidSpec(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	// Criar spec válida em diretório temporário
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "test.spec.md")

	validSpec := `# Test Spec

## 1. Contexto e Objetivo
Teste

## 2. Requisitos Funcionais
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

	if err := fs.WriteFile(specPath, []byte(validSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.Validate(ValidateOptions{
		Path: specPath,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("esperado 1 resultado, obtido %d", len(result.Results))
	}

	vr := result.Results[0]
	if !vr.Valid {
		t.Errorf("spec deveria ser válida, mas tem erros: %v", vr.Errors)
	}
	if !vr.Complete {
		t.Error("spec deveria estar completa")
	}
	if vr.Checklist.ItemCount != 6 {
		t.Errorf("esperado 6 itens no checklist, obtido %d", vr.Checklist.ItemCount)
	}
	if vr.Checklist.MarkedCount != 6 {
		t.Errorf("esperado 6 itens marcados, obtido %d", vr.Checklist.MarkedCount)
	}
}

func TestService_Validate_IncompleteSpec(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "test.spec.md")

	incompleteSpec := `# Test Spec

## 1. Contexto e Objetivo
Teste

## 2. Requisitos Funcionais
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
- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
`

	if err := fs.WriteFile(specPath, []byte(incompleteSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.Validate(ValidateOptions{
		Path: specPath,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("esperado 1 resultado, obtido %d", len(result.Results))
	}

	vr := result.Results[0]
	if !vr.Valid {
		t.Errorf("spec deveria ser válida (apenas incompleta), mas tem erros: %v", vr.Errors)
	}
	if vr.Complete {
		t.Error("spec não deveria estar completa")
	}
	if len(vr.Warnings) == 0 {
		t.Error("deveria ter warning sobre checklist incompleto")
	}
}

func TestService_Validate_MissingSections(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "test.spec.md")

	invalidSpec := `# Test Spec

## 1. Contexto e Objetivo
Teste

## 2. Requisitos Funcionais
Teste

## Checklist Rápido (preencha antes de gerar código)
- [ ] Item 1
- [ ] Item 2
- [ ] Item 3
- [ ] Item 4
- [ ] Item 5
- [ ] Item 6
`

	if err := fs.WriteFile(specPath, []byte(invalidSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec: %v", err)
	}

	result, err := service.Validate(ValidateOptions{
		Path: specPath,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("esperado 1 resultado, obtido %d", len(result.Results))
	}

	vr := result.Results[0]
	if vr.Valid {
		t.Error("spec deveria ser inválida (faltam seções)")
	}
	if len(vr.Errors) == 0 {
		t.Error("deveria ter erros sobre seções faltantes")
	}
}

func TestService_Validate_Directory(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar duas specs
	spec1 := filepath.Join(specsDir, "01-test.spec.md")
	spec2 := filepath.Join(specsDir, "02-test.spec.md")

	validSpec := `# Test Spec

## 1. Contexto e Objetivo
Teste

## 2. Requisitos Funcionais
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

	if err := fs.WriteFile(spec1, []byte(validSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}
	if err := fs.WriteFile(spec2, []byte(validSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec2: %v", err)
	}

	result, err := service.Validate(ValidateOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("esperado 2 specs, obtido %d", result.Total)
	}
	if result.Complete != 2 {
		t.Errorf("esperado 2 specs completas, obtido %d", result.Complete)
	}
	if result.WithErrors != 0 {
		t.Errorf("esperado 0 specs com erros, obtido %d", result.WithErrors)
	}
}

func TestService_Validate_InvalidPath(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	_, err := service.Validate(ValidateOptions{
		Path: "/caminho/inexistente/12345",
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho inexistente")
	}
}

func TestService_Validate_InvalidExtension(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "test.txt")

	if err := fs.WriteFile(invalidFile, []byte("test"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	_, err := service.Validate(ValidateOptions{
		Path: invalidFile,
	})
	if err == nil {
		t.Error("deveria retornar erro para arquivo sem extensão .spec.md")
	}
}
