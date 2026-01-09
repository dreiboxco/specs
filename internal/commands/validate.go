package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	validatorSvc "github.com/dreibox/specs/internal/services/validator"
)

// ValidateCommand implementa o comando validate
type ValidateCommand struct {
	fs          adapters.FileSystem
	validatorSvc *validatorSvc.Service
}

// NewValidateCommand cria uma nova instância do ValidateCommand
func NewValidateCommand(fs adapters.FileSystem) *ValidateCommand {
	return &ValidateCommand{
		fs:          fs,
		validatorSvc: validatorSvc.NewService(fs),
	}
}

// Execute executa o comando validate
func (c *ValidateCommand) Execute(args []string) int {
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

	// Executar validação
	result, err := c.validatorSvc.Validate(validatorSvc.ValidateOptions{
		Path: opts.Path,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir resultados
	c.printResults(result)

	// Determinar código de saída
	if result.WithErrors > 0 {
		return 1
	}
	return 0
}

// validateOptions contém opções do comando validate
type validateOptions struct {
	Path string
	Help bool
}

// parseArgs parseia argumentos e flags
func (c *ValidateCommand) parseArgs(args []string) (*validateOptions, error) {
	opts := &validateOptions{}

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

// printResults exibe resultados da validação
func (c *ValidateCommand) printResults(result *validatorSvc.ValidateResult) {
	// Determinar se é diretório ou arquivo único
	isDir := false
	if len(result.Results) > 1 {
		isDir = true
	}

	if isDir {
		// Mostrar caminho sendo validado
		if len(result.Results) > 0 {
			dir := filepath.Dir(result.Results[0].Path)
			relPath, _ := filepath.Rel(".", dir)
			if relPath == "" || relPath == "." {
				relPath = "./specs"
			}
			fmt.Printf("Validando specs em %s...\n\n", relPath)
		}
	}

	// Exibir resultado de cada spec
	for _, vr := range result.Results {
		relPath, _ := filepath.Rel(".", vr.Path)
		if relPath == "" || relPath == "." {
			relPath = vr.Path
		}

		if len(vr.Errors) > 0 {
			// Spec com erros
			fmt.Printf("❌ %s: Erro", relPath)
			if len(vr.Errors) == 1 {
				fmt.Printf(" - %s\n", vr.Errors[0])
			} else {
				fmt.Println()
				for _, err := range vr.Errors {
					fmt.Printf("   - %s\n", err)
				}
			}
		} else if vr.Complete {
			// Spec completa
			fmt.Printf("✅ %s: Completa (%d/6 itens do checklist)\n", relPath, vr.Checklist.MarkedCount)
		} else {
			// Spec incompleta
			fmt.Printf("⚠️  %s: Incompleta (%d/6 itens do checklist)\n", relPath, vr.Checklist.MarkedCount)
		}
	}

	// Exibir resumo se houver múltiplos arquivos
	if isDir && result.Total > 0 {
		fmt.Println()
		fmt.Println("Resumo:")
		fmt.Printf("  Total: %d specs\n", result.Total)
		fmt.Printf("  Completas: %d\n", result.Complete)
		fmt.Printf("  Incompletas: %d\n", result.Incomplete)
		fmt.Printf("  Com erros: %d\n", result.WithErrors)
	}
}

func (c *ValidateCommand) printHelp() {
	fmt.Println("Valida specs contra checklist formal e verifica estrutura.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs validate [caminho] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help    Exibe ajuda para este comando")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs validate                    # Valida specs/ no diretório atual")
	fmt.Println("  specs validate specs/             # Valida diretório specs/")
	fmt.Println("  specs validate specs/01-test.spec.md  # Valida arquivo específico")
}
