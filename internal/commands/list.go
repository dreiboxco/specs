package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	listerSvc "github.com/dreibox/specs/internal/services/lister"
)

// ListCommand implementa o comando list
type ListCommand struct {
	fs          adapters.FileSystem
	listerSvc   *listerSvc.Service
}

// NewListCommand cria uma nova instância do ListCommand
func NewListCommand(fs adapters.FileSystem) *ListCommand {
	return &ListCommand{
		fs:        fs,
		listerSvc: listerSvc.NewService(fs),
	}
}

// Execute executa o comando list
func (c *ListCommand) Execute(args []string) int {
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

	// Executar listagem
	result, err := c.listerSvc.List(listerSvc.ListOptions{
		Path:       opts.Path,
		Complete:   opts.Complete,
		Incomplete: opts.Incomplete,
		Errors:     opts.Errors,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir resultados
	c.printResults(result, opts)

	return 0
}

// listOptions contém opções do comando list
type listOptions struct {
	Path       string
	Complete   bool
	Incomplete bool
	Errors     bool
	Help       bool
}

// parseArgs parseia argumentos e flags
func (c *ListCommand) parseArgs(args []string) (*listOptions, error) {
	opts := &listOptions{}

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			opts.Help = true
			return opts, nil
		case "--complete", "--only-complete":
			opts.Complete = true
		case "--incomplete", "--only-incomplete":
			opts.Incomplete = true
		case "--errors":
			opts.Errors = true
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

// printResults exibe resultados da listagem
func (c *ListCommand) printResults(result *listerSvc.ListResult, opts *listOptions) {
	// Verificar se há specs
	if len(result.Specs) == 0 {
		path := opts.Path
		if path == "" {
			path = "./specs"
		}
		
		var msg string
		if opts.Complete {
			msg = "Nenhuma spec completa encontrada"
		} else if opts.Incomplete {
			msg = "Nenhuma spec incompleta encontrada"
		} else if opts.Errors {
			msg = "Nenhuma spec com erros encontrada"
		} else {
			msg = fmt.Sprintf("Nenhuma spec encontrada em %s", path)
		}
		fmt.Println(msg)
		return
	}

	// Determinar caminho para exibição
	path := opts.Path
	if path == "" {
		path = "./specs"
	}
	relPath, _ := filepath.Rel(".", path)
	if relPath == "" || relPath == "." {
		relPath = path
	}

	// Determinar tipo de listagem
	var listType string
	if opts.Complete {
		listType = "completas"
	} else if opts.Incomplete {
		listType = "incompletas"
	} else if opts.Errors {
		listType = "com erros"
	} else {
		listType = ""
	}

	if listType != "" {
		fmt.Printf("Listando specs %s em %s...\n\n", listType, relPath)
	} else {
		fmt.Printf("Listando specs em %s...\n\n", relPath)
	}

	// Exibir tabela
	c.printTable(result.Specs)

	// Exibir resumo
	fmt.Println()
	if opts.Complete {
		fmt.Printf("Total: %d specs completas\n", result.Complete)
	} else if opts.Incomplete {
		fmt.Printf("Total: %d specs incompletas\n", result.Incomplete)
	} else if opts.Errors {
		fmt.Printf("Total: %d specs com erros\n", result.WithErrors)
	} else {
		fmt.Println("Resumo:")
		fmt.Printf("  Total: %d specs\n", result.Total)
		fmt.Printf("  Completas: %d\n", result.Complete)
		fmt.Printf("  Incompletas: %d\n", result.Incomplete)
		fmt.Printf("  Com erros: %d\n", result.WithErrors)
	}
}

// printTable exibe tabela formatada
func (c *ListCommand) printTable(specs []listerSvc.SpecInfo) {
	// Calcular larguras das colunas
	maxNumLen := len("Numeração")
	maxNameLen := len("Nome")
	maxStatusLen := len("Status")

	for _, spec := range specs {
		if len(spec.Number) > maxNumLen {
			maxNumLen = len(spec.Number)
		}
		if len(spec.Name) > maxNameLen {
			maxNameLen = len(spec.Name)
		}
		statusLen := len(spec.StatusIcon + " " + spec.Status)
		if statusLen > maxStatusLen {
			maxStatusLen = statusLen
		}
	}

	// Garantir larguras mínimas
	if maxNumLen < 10 {
		maxNumLen = 10
	}
	if maxNameLen < 20 {
		maxNameLen = 20
	}
	if maxStatusLen < 10 {
		maxStatusLen = 10
	}

	// Cabeçalho
	fmt.Printf("%-*s  %-*s  %s\n", maxNumLen, "Numeração", maxNameLen, "Nome", "Status")
	
	// Separador
	separator := strings.Repeat("─", maxNumLen) + "  " + strings.Repeat("─", maxNameLen) + "  " + strings.Repeat("─", maxStatusLen)
	fmt.Println(separator)

	// Linhas da tabela
	for _, spec := range specs {
		status := spec.StatusIcon + " " + spec.Status
		fmt.Printf("%-*s  %-*s  %s\n", maxNumLen, spec.Number, maxNameLen, spec.Name, status)
	}
}

func (c *ListCommand) printHelp() {
	fmt.Println("Lista todas as specs do projeto com status (completa/incompleta).")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs list [caminho] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --complete, --only-complete     Lista apenas specs completas")
	fmt.Println("  --incomplete, --only-incomplete  Lista apenas specs incompletas")
	fmt.Println("  --errors                         Lista apenas specs com erros")
	fmt.Println("  --help                           Exibe ajuda para este comando")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs list                       # Lista todas as specs em specs/")
	fmt.Println("  specs list --complete            # Lista apenas specs completas")
	fmt.Println("  specs list --incomplete          # Lista apenas specs incompletas")
	fmt.Println("  specs list specs/                # Lista specs em diretório específico")
}
