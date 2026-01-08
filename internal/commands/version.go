package commands

import (
	"fmt"
	"os"

	"github.com/dreibox/specs/internal/adapters"
	"github.com/dreibox/specs/internal/services/version"
)

// VersionCommand implementa o comando version
type VersionCommand struct {
	fs           adapters.FileSystem
	injectedVer  string
	versionSvc   *version.Service
}

// NewVersionCommand cria uma nova instância do VersionCommand
func NewVersionCommand(fs adapters.FileSystem, injectedVer string) *VersionCommand {
	return &VersionCommand{
		fs:          fs,
		injectedVer: injectedVer,
		versionSvc:  version.NewService(fs),
	}
}

// Execute executa o comando version
func (c *VersionCommand) Execute(args []string) int {
	// Verificar flag --help
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		c.printHelp()
		return 0
	}

	// Se versão foi injetada durante build, usar ela
	if c.injectedVer != "" && c.injectedVer != "dev" {
		fmt.Println(c.injectedVer)
		return 0
	}

	// Tentar ler do arquivo VERSION
	ver, err := c.versionSvc.GetVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 1
	}

	fmt.Println(ver)
	return 0
}

func (c *VersionCommand) printHelp() {
	fmt.Println("Exibe a versão atual do CLI.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs version [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help    Exibe ajuda para este comando")
}

