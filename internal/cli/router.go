package cli

import (
	"fmt"
	"os"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/commands"
)

// Router gerencia o roteamento de comandos
type Router struct {
	fs      adapters.FileSystem
	version string
}

// NewRouter cria uma nova instância do Router
func NewRouter(fs adapters.FileSystem, version string) *Router {
	return &Router{
		fs:      fs,
		version: version,
	}
}

// Run executa o comando apropriado baseado nos argumentos
func (r *Router) Run(args []string) int {
	if len(args) == 0 {
		r.printUsage()
		return 0
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "version":
		versionCmd := commands.NewVersionCommand(r.fs, r.version)
		return versionCmd.Execute(cmdArgs)
	case "init":
		initCmd := commands.NewInitCommand(r.fs)
		return initCmd.Execute(cmdArgs)
	case "validate":
		validateCmd := commands.NewValidateCommand(r.fs)
		return validateCmd.Execute(cmdArgs)
	case "list":
		listCmd := commands.NewListCommand(r.fs)
		return listCmd.Execute(cmdArgs)
	case "check":
		checkCmd := commands.NewCheckCommand(r.fs)
		return checkCmd.Execute(cmdArgs)
	case "help", "--help", "-h":
		r.printHelp()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "erro: comando desconhecido '%s'\n", cmd)
		fmt.Fprintf(os.Stderr, "Execute 'specs --help' para ver comandos disponíveis\n")
		return 1
	}
}

func (r *Router) printUsage() {
	fmt.Println("specs - CLI para gerenciamento de projetos SDD")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs <comando> [flags]")
	fmt.Println()
	fmt.Println("Comandos:")
	fmt.Println("  init       Inicializa um novo projeto SDD")
	fmt.Println("  list       Lista todas as specs com status")
	fmt.Println("  validate   Valida specs contra checklist formal")
	fmt.Println("  check      Verifica consistência estrutural de specs")
	fmt.Println("  version    Exibe a versão atual")
	fmt.Println("  help       Exibe ajuda")
	fmt.Println()
	fmt.Println("Execute 'specs <comando> --help' para mais informações sobre um comando.")
}

func (r *Router) printHelp() {
	r.printUsage()
}

