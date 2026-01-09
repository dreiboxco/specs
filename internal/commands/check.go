package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	checkerSvc "github.com/dreibox/specs/internal/services/checker"
)

// CheckCommand implementa o comando check
type CheckCommand struct {
	fs          adapters.FileSystem
	checkerSvc  *checkerSvc.Service
}

// NewCheckCommand cria uma nova instância do CheckCommand
func NewCheckCommand(fs adapters.FileSystem) *CheckCommand {
	return &CheckCommand{
		fs:         fs,
		checkerSvc: checkerSvc.NewService(fs),
	}
}

// Execute executa o comando check
func (c *CheckCommand) Execute(args []string) int {
	// Parsear flags e argumentos
	opts, err := c.parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Verificar flag --help
	if opts.Help {
		c.printHelp()
		return 0
	}

	// Executar verificação
	result, err := c.checkerSvc.Check(checkerSvc.CheckOptions{
		Path: opts.Path,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir resultados
	c.printResults(result, opts)

	// Determinar código de saída
	if len(result.Problems) > 0 {
		return 1
	}
	return 0
}

// checkOptions contém opções do comando check
type checkOptions struct {
	Path string
	Help bool
}

// parseArgs parseia argumentos e flags
func (c *CheckCommand) parseArgs(args []string) (*checkOptions, error) {
	opts := &checkOptions{}

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			opts.Help = true
			return opts, nil
		case "--json":
			// Flag para futuro (v2)
			// Por enquanto ignorar
		default:
			if strings.HasPrefix(arg, "-") {
				if arg != "--json" {
					return nil, fmt.Errorf("flag desconhecida: %s", arg)
				}
			} else if opts.Path == "" {
				opts.Path = arg
			}
		}
	}

	return opts, nil
}

// printResults exibe resultados da verificação
func (c *CheckCommand) printResults(result *checkerSvc.CheckResult, opts *checkOptions) {
	path := opts.Path
	if path == "" {
		path = "./specs"
	}
	relPath, _ := filepath.Rel(".", path)
	if relPath == "" || relPath == "." {
		relPath = path
	}

	fmt.Printf("Verificando consistência estrutural em %s...\n\n", relPath)

	if len(result.Problems) == 0 {
		fmt.Println("✅ Numeração: OK")
		fmt.Println("✅ Links: Todos os links válidos")
		fmt.Println("✅ Estrutura: OK")
		fmt.Println()
		fmt.Println("Todas as verificações passaram!")
		return
	}

	// Agrupar problemas por categoria
	problemsByCategory := make(map[string][]checkerSvc.Problem)
	for _, p := range result.Problems {
		problemsByCategory[p.Category] = append(problemsByCategory[p.Category], p)
	}

	// Exibir problemas por categoria
	categories := []string{"Numeração", "Links", "Formato", "Órfãs"}
	for _, category := range categories {
		problems := problemsByCategory[category]
		if len(problems) == 0 {
			fmt.Printf("✅ %s: OK\n", category)
			continue
		}

		// Determinar severidade geral
		hasError := false
		for _, p := range problems {
			if p.Severity == "error" {
				hasError = true
				break
			}
		}

		icon := "⚠️"
		if hasError {
			icon = "❌"
		}

		// Contar problemas
		errorCount := 0
		warningCount := 0
		for _, p := range problems {
			if p.Severity == "error" {
				errorCount++
			} else {
				warningCount++
			}
		}

		var msg string
		if errorCount > 0 && warningCount > 0 {
			msg = fmt.Sprintf("%d problema(s) encontrado(s) (%d erro(s), %d aviso(s))", len(problems), errorCount, warningCount)
		} else if errorCount > 0 {
			msg = fmt.Sprintf("%d problema(s) encontrado(s)", errorCount)
		} else {
			msg = fmt.Sprintf("%d aviso(s) encontrado(s)", warningCount)
		}

		fmt.Printf("%s %s: %s\n", icon, category, msg)

		// Exibir detalhes dos problemas
		for _, p := range problems {
			if p.File != "" {
				if p.Line > 0 {
					fmt.Printf("  - %s:%d: %s\n", p.File, p.Line, p.Message)
				} else {
					fmt.Printf("  - %s: %s\n", p.File, p.Message)
				}
			} else {
				fmt.Printf("  - %s\n", p.Message)
			}
		}
		fmt.Println()
	}

	// Exibir resumo
	fmt.Println("Resumo:")
	fmt.Printf("  Total de specs: %d\n", result.TotalSpecs)
	fmt.Printf("  Problemas encontrados: %d\n", len(result.Problems))
	for category, count := range result.Summary {
		fmt.Printf("  - %s: %d\n", category, count)
	}
}

func (c *CheckCommand) printHelp() {
	fmt.Println("Verifica consistência estrutural de specs (numeração, links, referências).")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs check [caminho] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help    Exibe ajuda para este comando")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs check                    # Verifica specs/ no diretório atual")
	fmt.Println("  specs check specs/             # Verifica diretório specs/")
}
