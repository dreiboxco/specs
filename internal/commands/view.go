package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	configSvc "github.com/dreibox/specs/internal/services/config"
	viewerSvc "github.com/dreibox/specs/internal/services/viewer"
)

// ViewCommand implementa o comando view
type ViewCommand struct {
	fs          adapters.FileSystem
	viewerSvc   *viewerSvc.Service
	configSvc   *configSvc.Service
}

// NewViewCommand cria uma nova instância do ViewCommand
func NewViewCommand(fs adapters.FileSystem) *ViewCommand {
	return &ViewCommand{
		fs:        fs,
		viewerSvc: viewerSvc.NewService(fs),
		configSvc: configSvc.NewService(fs),
	}
}

// Execute executa o comando view
func (c *ViewCommand) Execute(args []string) int {
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

	// Resolver caminho padrão se não fornecido
	path := opts.Path
	if path == "" {
		resolvedPath, err := c.configSvc.ResolveDefaultPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			return 1
		}
		path = resolvedPath
	}

	// Executar visualização
	result, err := c.viewerSvc.View(viewerSvc.ViewOptions{
		Path: path,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir dashboard
	c.printDashboard(result, opts)

	return 0
}

// viewOptions contém opções do comando view
type viewOptions struct {
	Path string
	Help bool
}

// parseArgs parseia argumentos e flags
func (c *ViewCommand) parseArgs(args []string) (*viewOptions, error) {
	opts := &viewOptions{}

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

// printDashboard exibe dashboard formatado
func (c *ViewCommand) printDashboard(result *viewerSvc.DashboardResult, opts *viewOptions) {
	fmt.Println("Specs Dashboard")
	fmt.Println()

	// Seção Summary
	fmt.Println("Summary:")
	fmt.Printf("  Specifications: %d specs, %d requirements\n", result.TotalSpecs, result.TotalRequirements)
	fmt.Printf("  Specs em Progresso: %d\n", result.SpecsInProgress)
	fmt.Printf("  Specs Completas: %d\n", result.SpecsComplete)
	fmt.Printf("  Progresso Geral: %s\n", result.OverallProgressStr)
	fmt.Println()

	// Separar specs completas e em progresso
	var inProgress []viewerSvc.SpecStats
	var complete []viewerSvc.SpecStats

	for _, spec := range result.Specs {
		if spec.Complete {
			complete = append(complete, spec)
		} else {
			inProgress = append(inProgress, spec)
		}
	}

	// Ordenar specs em progresso por percentual (menor para maior)
	sort.Slice(inProgress, func(i, j int) bool {
		return inProgress[i].Progress < inProgress[j].Progress
	})

	// Seção Specs em Progresso
	if len(inProgress) > 0 {
		fmt.Println("Specs em Progresso:")
		for _, spec := range inProgress {
			percent := int(spec.Progress * 100)
			bar := c.generateProgressBar(spec.Progress, 10)
			fmt.Printf("  %-30s %s %d%%\n", spec.Name, bar, percent)
		}
		fmt.Println()
	}

	// Seção Specs Completas
	if len(complete) > 0 {
		fmt.Println("Specs Completas:")
		for _, spec := range complete {
			fmt.Printf("  ✅ %s\n", spec.Name)
		}
		fmt.Println()
	}

	// Seção Specifications
	fmt.Println("Specifications:")
	for _, spec := range result.Specs {
		fmt.Printf("  %-30s %d requirements\n", spec.Name, spec.Requirements)
	}
}

// generateProgressBar gera barra de progresso visual
func (c *ViewCommand) generateProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat(" ", empty))
}

func (c *ViewCommand) printHelp() {
	fmt.Println("Exibe dashboard interativo com informações agregadas do projeto SDD.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs view [caminho] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help    Exibe ajuda para este comando")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs view                    # Dashboard de specs/ no diretório atual")
	fmt.Println("  specs view specs/             # Dashboard de diretório específico")
}
