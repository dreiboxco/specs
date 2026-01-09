package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	updateSvc "github.com/dreibox/specs/internal/services/update"
)

// UpdateCommand implementa o comando update
type UpdateCommand struct {
	fs        adapters.FileSystem
	updateSvc *updateSvc.Service
}

// NewUpdateCommand cria uma nova instância do UpdateCommand
func NewUpdateCommand(fs adapters.FileSystem) *UpdateCommand {
	return &UpdateCommand{
		fs:        fs,
		updateSvc: updateSvc.NewService(fs),
	}
}

// Execute executa o comando update
func (c *UpdateCommand) Execute(args []string) int {
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

	// Executar atualização
	result, err := c.updateSvc.Update(opts.UpdateOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 1
	}

	// Exibir resultados
	c.printResults(result, opts.UpdateOptions.DryRun)

	return 0
}

// updateOptions contém opções parseadas do comando
type updateOptions struct {
	UpdateOptions updateSvc.UpdateOptions
	Help          bool
}

// parseArgs parseia argumentos e flags
func (c *UpdateCommand) parseArgs(args []string) (*updateOptions, error) {
	opts := &updateOptions{
		UpdateOptions: updateSvc.UpdateOptions{},
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			opts.Help = true
			return opts, nil
		}

		if arg == "--dry-run" {
			opts.UpdateOptions.DryRun = true
			continue
		}

		if arg == "--force" {
			opts.UpdateOptions.Force = true
			continue
		}

		if arg == "--no-backup" {
			opts.UpdateOptions.NoBackup = true
			continue
		}

		if arg == "--merge" {
			opts.UpdateOptions.Merge = true
			continue
		}

		if strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("flag desconhecida: %s", arg)
		}

		// Argumento posicional: diretório
		if opts.UpdateOptions.TargetDir == "" {
			opts.UpdateOptions.TargetDir = arg
			continue
		}

		return nil, fmt.Errorf("argumento extra: %s", arg)
	}

	return opts, nil
}

// printResults exibe resultados da atualização
func (c *UpdateCommand) printResults(result *updateSvc.UpdateResult, dryRun bool) {
	if dryRun {
		fmt.Println("Arquivos que seriam atualizados:")
		for _, file := range result.FilesUpdated {
			fmt.Printf("  - specs/%s\n", file)
		}
		if len(result.FilesSkipped) > 0 {
			fmt.Println("\nArquivos que requerem atenção:")
			for _, file := range result.FilesSkipped {
				fmt.Printf("  - %s\n", file)
			}
		}
		return
	}

	// Exibir backup
	if result.BackupDir != "" {
		relPath, _ := filepath.Rel(".", result.BackupDir)
		if relPath == "" || relPath == "." {
			relPath = result.BackupDir
		}
		fmt.Printf("Backup criado em %s\n", relPath)
	}

	// Exibir arquivos atualizados
	if len(result.FilesUpdated) > 0 {
		fmt.Println("Atualizando templates...")
		for _, file := range result.FilesUpdated {
			fmt.Printf("  ✓ %s atualizado\n", file)
		}
	}

	// Exibir .cursorrules
	if result.CursorRulesUpdated {
		fmt.Println("Atualizando .cursorrules...")
		fmt.Println("  ✓ .cursorrules atualizado")
	} else if result.HasCustomizations {
		fmt.Println("Atualizando .cursorrules...")
		fmt.Fprintf(os.Stderr, "  ⚠ Regras personalizadas detectadas em .cursorrules\n")
		fmt.Println("  Arquivo .cursorrules-updated criado com versão do boilerplate")
		if result.CursorRulesMerged {
			fmt.Println("  Arquivo .cursorrules-merged criado com merge automático")
			fmt.Println("  Revise o arquivo e substitua .cursorrules se estiver correto")
		} else {
			fmt.Println("  Execute merge manual ou use: specs update --merge")
		}
	}

	// Exibir arquivos pulados
	if len(result.FilesSkipped) > 0 {
		fmt.Fprintf(os.Stderr, "\naviso: alguns arquivos não foram atualizados:\n")
		for _, file := range result.FilesSkipped {
			fmt.Fprintf(os.Stderr, "  - %s\n", file)
		}
	}
}

// printHelp exibe ajuda do comando
func (c *UpdateCommand) printHelp() {
	fmt.Println("Atualiza templates e arquivos base do projeto SDD.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs update [diretório] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --dry-run      Exibe o que seria atualizado sem fazer alterações")
	fmt.Println("  --force        Força atualização mesmo se não houver diferenças")
	fmt.Println("  --no-backup    Não cria backup antes de atualizar")
	fmt.Println("  --merge        Tenta merge automático de .cursorrules (experimental)")
	fmt.Println("  --help, -h     Exibe esta ajuda")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs update                    # Atualiza templates no diretório atual")
	fmt.Println("  specs update --dry-run           # Preview sem alterações")
	fmt.Println("  specs update --merge             # Tenta merge automático de .cursorrules")
}
