package templates

import (
	_ "embed"
	"os"
	"path/filepath"
)

// GetCursorRulesTemplate retorna o template de .cursorrules
func GetCursorRulesTemplate() ([]byte, error) {
	// Tentar encontrar boilerplate relativo ao executável
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		paths := []string{
			filepath.Join(exeDir, "..", "boilerplate", ".cursorrules"),
			filepath.Join(exeDir, "boilerplate", ".cursorrules"),
			filepath.Join(filepath.Dir(exeDir), "boilerplate", ".cursorrules"),
		}
		for _, path := range paths {
			if data, err := os.ReadFile(path); err == nil {
				return data, nil
			}
		}
	}

	// Tentar relativo ao diretório de trabalho atual
	wd, err := os.Getwd()
	if err == nil {
		paths := []string{
			filepath.Join(wd, "boilerplate", ".cursorrules"),
			filepath.Join(wd, "..", "boilerplate", ".cursorrules"),
		}
		for _, path := range paths {
			if data, err := os.ReadFile(path); err == nil {
				return data, nil
			}
		}
	}

	// Fallback: tentar caminho absoluto do projeto (para desenvolvimento)
	projectRoot := findProjectRoot()
	if projectRoot != "" {
		path := filepath.Join(projectRoot, "boilerplate", ".cursorrules")
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}

	// Fallback: template básico embarcado
	return []byte(`# Cursor Rules - Spec Driven Development

## Linguagem e comunicação
- Responder sempre em português (pt-BR).

## Fonte da verdade e fluxo SDD
- Especificar antes de codificar. Só gerar/alterar código a partir de specs em specs/ validadas.
- Não implementar feature sem *.spec.md completa e checada pelo specs/checklist.md.
- Antes de iniciar implementação, passar a spec pelo checklist.

## Fluxo obrigatório de consulta de specs
- SEMPRE, antes de implementar qualquer funcionalidade:
  1. Consultar as specs existentes em specs/
  2. Se estiver especificada: implementar conforme a spec
  3. Se NÃO estiver especificada: perguntar ao usuário antes de implementar

## Estrutura e diretórios
- Respeitar a estrutura definida em specs/00-architecture.spec.md
- specs/ não recebe código executável

## Guardrails de dependências
- Evitar adicionar dependências novas

## UX de CLI
- Mensagens curtas e acionáveis
- Help sempre disponível (--help)
- Códigos de saída padronizados

## Segurança
- Nunca logar ou imprimir segredos/tokens

## Testes e critérios de aceite
- Gerar/rodar testes previstos em cada spec
- Validar critérios de aceite antes de concluir
`), nil
}
