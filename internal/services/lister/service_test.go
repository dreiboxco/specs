package lister

import (
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_List_Success(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	// Criar specs em diretório temporário
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	// Criar spec completa
	completeSpec := `# 01 Test Spec

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

	// Criar spec incompleta
	incompleteSpec := `# 02 Test Spec Incomplete

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

	spec1 := filepath.Join(specsDir, "01-test.spec.md")
	spec2 := filepath.Join(specsDir, "02-test-incomplete.spec.md")

	if err := fs.WriteFile(spec1, []byte(completeSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}
	if err := fs.WriteFile(spec2, []byte(incompleteSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec2: %v", err)
	}

	result, err := service.List(ListOptions{
		Path: specsDir,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("esperado 2 specs, obtido %d", result.Total)
	}
	if result.Complete != 1 {
		t.Errorf("esperado 1 spec completa, obtido %d", result.Complete)
	}
	if result.Incomplete != 1 {
		t.Errorf("esperado 1 spec incompleta, obtido %d", result.Incomplete)
	}

	// Verificar ordenação
	if len(result.Specs) != 2 {
		t.Fatalf("esperado 2 specs na lista, obtido %d", len(result.Specs))
	}
	if result.Specs[0].Number != "01" {
		t.Errorf("primeira spec deveria ter numeração 01, obtido %s", result.Specs[0].Number)
	}
	if result.Specs[1].Number != "02" {
		t.Errorf("segunda spec deveria ter numeração 02, obtido %s", result.Specs[1].Number)
	}
}

func TestService_List_FilterComplete(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	completeSpec := `# 01 Test Spec

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

	incompleteSpec := `# 02 Test Spec Incomplete

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

	spec1 := filepath.Join(specsDir, "01-test.spec.md")
	spec2 := filepath.Join(specsDir, "02-test-incomplete.spec.md")

	if err := fs.WriteFile(spec1, []byte(completeSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}
	if err := fs.WriteFile(spec2, []byte(incompleteSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec2: %v", err)
	}

	result, err := service.List(ListOptions{
		Path:     specsDir,
		Complete: true,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("esperado 1 spec completa, obtido %d", result.Total)
	}
	if len(result.Specs) != 1 {
		t.Errorf("esperado 1 spec na lista, obtido %d", len(result.Specs))
	}
	if !result.Specs[0].Complete {
		t.Error("spec deveria estar completa")
	}
}

func TestService_List_FilterIncomplete(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := fs.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório: %v", err)
	}

	incompleteSpec := `# 02 Test Spec Incomplete

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

	spec1 := filepath.Join(specsDir, "02-test-incomplete.spec.md")

	if err := fs.WriteFile(spec1, []byte(incompleteSpec), 0644); err != nil {
		t.Fatalf("falha ao criar spec1: %v", err)
	}

	result, err := service.List(ListOptions{
		Path:       specsDir,
		Incomplete: true,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("esperado 1 spec incompleta, obtido %d", result.Total)
	}
	if len(result.Specs) != 1 {
		t.Errorf("esperado 1 spec na lista, obtido %d", len(result.Specs))
	}
	if result.Specs[0].Complete {
		t.Error("spec não deveria estar completa")
	}
}

func TestService_List_InvalidPath(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	_, err := service.List(ListOptions{
		Path: "/caminho/inexistente/12345",
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho inexistente")
	}
}

func TestService_List_NotDirectory(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")

	if err := fs.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo: %v", err)
	}

	_, err := service.List(ListOptions{
		Path: file,
	})
	if err == nil {
		t.Error("deveria retornar erro para caminho que não é diretório")
	}
}
