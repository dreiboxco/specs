package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	initSvc "github.com/dreibox/specs/internal/services/init"
)

// InitCommand implementa o comando init
type InitCommand struct {
	fs        adapters.FileSystem
	initSvc   *initSvc.Service
}

// NewInitCommand cria uma nova instância do InitCommand
func NewInitCommand(fs adapters.FileSystem) *InitCommand {
	return &InitCommand{
		fs:      fs,
		initSvc: initSvc.NewService(fs),
	}
}

// Execute executa o comando init
func (c *InitCommand) Execute(args []string) int {
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

	// Executar inicialização
	result, err := c.initSvc.Initialize(opts.InitOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 1
	}

	// Verificar se projeto já existia (idempotente)
	if len(result.FilesCreated) == 0 && len(result.DirectoriesCreated) == 0 {
		relPath, _ := filepath.Rel(".", result.SpecsDir)
		if relPath == "" || relPath == "." {
			relPath = "./specs"
		}
		fmt.Printf("Projeto SDD já existe em %s. Nada a fazer.\n", relPath)
		return 0
	}

	// Exibir mensagem de sucesso
	relPath, _ := filepath.Rel(".", result.SpecsDir)
	if relPath == "" || relPath == "." {
		relPath = "./specs"
	}
	fmt.Printf("Projeto SDD inicializado com sucesso em %s\n", relPath)

	if opts.InitOptions.WithBoilerplate {
		boilerplatePath := filepath.Join(filepath.Dir(result.SpecsDir), "boilerplate")
		relBoilerplate, _ := filepath.Rel(".", boilerplatePath)
		if relBoilerplate == "" || relBoilerplate == "." {
			relBoilerplate = "./boilerplate"
		}
		fmt.Printf("Boilerplate criado em %s\n", relBoilerplate)
	}

	return 0
}

// parseArgs parseia argumentos e flags
func (c *InitCommand) parseArgs(args []string) (*initOptions, error) {
	opts := &initOptions{
		InitOptions: initSvc.InitOptions{},
	}

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			opts.Help = true
			return opts, nil
		case "--force":
			opts.InitOptions.Force = true
		case "--with-boilerplate":
			opts.InitOptions.WithBoilerplate = true
		default:
			// Se não é flag, é o diretório alvo
			if !strings.HasPrefix(arg, "-") && opts.InitOptions.TargetDir == "" {
				opts.InitOptions.TargetDir = arg
			} else if strings.HasPrefix(arg, "-") {
				return nil, fmt.Errorf("flag desconhecida: %s", arg)
			}
		}
	}

	return opts, nil
}

type initOptions struct {
	InitOptions initSvc.InitOptions
	Help        bool
}

func (c *InitCommand) printHelp() {
	fmt.Println("Inicializa um novo projeto SDD no diretório especificado.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs init [diretório] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --force              Sobrescreve arquivos existentes sem confirmação")
	fmt.Println("  --with-boilerplate   Cria também diretório boilerplate com templates genéricos")
	fmt.Println("  --help               Exibe ajuda para este comando")
}
