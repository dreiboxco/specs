package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
	configSvc "github.com/dreibox/specs/internal/services/config"
)

// ConfigCommand implementa o comando config
type ConfigCommand struct {
	fs          adapters.FileSystem
	configSvc   *configSvc.Service
}

// NewConfigCommand cria uma nova instância do ConfigCommand
func NewConfigCommand(fs adapters.FileSystem) *ConfigCommand {
	return &ConfigCommand{
		fs:        fs,
		configSvc: configSvc.NewService(fs),
	}
}

// Execute executa o comando config
func (c *ConfigCommand) Execute(args []string) int {
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

	// Executar subcomando
	switch opts.Subcommand {
	case "show", "":
		return c.executeShow()
	case "get":
		if opts.Key == "" {
			fmt.Fprintf(os.Stderr, "erro: chave não especificada\n")
			fmt.Fprintf(os.Stderr, "Uso: specs config get <chave>\n")
			return 2
		}
		return c.executeGet(opts.Key)
	case "set":
		if opts.Key == "" || opts.Value == "" {
			fmt.Fprintf(os.Stderr, "erro: chave e valor devem ser especificados\n")
			fmt.Fprintf(os.Stderr, "Uso: specs config set <chave> <valor>\n")
			return 2
		}
		return c.executeSet(opts.Key, opts.Value)
	default:
		fmt.Fprintf(os.Stderr, "erro: subcomando desconhecido '%s'\n", opts.Subcommand)
		c.printHelp()
		return 2
	}
}

// configOptions contém opções do comando config
type configOptions struct {
	Subcommand string
	Key        string
	Value      string
	Help       bool
}

// parseArgs parseia argumentos e flags
func (c *ConfigCommand) parseArgs(args []string) (*configOptions, error) {
	opts := &configOptions{}

	for i, arg := range args {
		switch arg {
		case "--help", "-h":
			opts.Help = true
			return opts, nil
		case "show", "get", "set":
			if opts.Subcommand == "" {
				opts.Subcommand = arg
				if arg == "get" && i+1 < len(args) {
					opts.Key = args[i+1]
				} else if arg == "set" && i+1 < len(args) && i+2 < len(args) {
					opts.Key = args[i+1]
					opts.Value = args[i+2]
				}
			}
		default:
			if strings.HasPrefix(arg, "-") {
				return nil, fmt.Errorf("flag desconhecida: %s", arg)
			}
		}
	}

	return opts, nil
}

// executeShow exibe configuração atual
func (c *ConfigCommand) executeShow() int {
	config, err := c.configSvc.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 1
	}

	configPath, err := c.configSvc.GetConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 1
	}

	// Verificar se arquivo existe
	exists := c.fs.Exists(configPath)

	// Exibir caminho do arquivo
	fmt.Printf("Configuração em: %s\n", configPath)
	if !exists {
		fmt.Println("(arquivo não existe, usando valores padrão)")
	}
	fmt.Println()

	// Exibir configuração formatada
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: falha ao formatar configuração: %v\n", err)
		return 1
	}

	fmt.Println(string(data))
	return 0
}

// executeGet obtém valor de uma chave específica
func (c *ConfigCommand) executeGet(key string) int {
	value, err := c.configSvc.GetValue(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir valor
	switch v := value.(type) {
	case string:
		fmt.Println(v)
	case bool:
		fmt.Printf("%t\n", v)
	default:
		fmt.Printf("%v\n", v)
	}

	return 0
}

// executeSet define valor de uma chave
func (c *ConfigCommand) executeSet(key string, valueStr string) int {
	// Determinar tipo do valor baseado na chave
	var value interface{}
	if strings.HasSuffix(key, "exclude_templates") {
		// Boolean
		value = strings.ToLower(valueStr) == "true" || valueStr == "1" || strings.ToLower(valueStr) == "yes"
	} else {
		// String
		value = valueStr
	}

	// Definir valor
	if err := c.configSvc.SetValue(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return 2
	}

	// Exibir confirmação
	fmt.Printf("Configuração atualizada: %s = %v\n", key, value)
	return 0
}

// printHelp exibe ajuda do comando
func (c *ConfigCommand) printHelp() {
	fmt.Println("Gerencia configuração do CLI specs.")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  specs config [subcomando] [flags]")
	fmt.Println()
	fmt.Println("Subcomandos:")
	fmt.Println("  show              Exibe configuração atual (padrão)")
	fmt.Println("  get <chave>       Obtém valor de uma chave específica")
	fmt.Println("  set <chave> <valor>  Define valor de uma chave")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help            Exibe ajuda para este comando")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  specs config                    # Exibe configuração completa")
	fmt.Println("  specs config get specs.default_path  # Obtém caminho padrão")
	fmt.Println("  specs config set specs.default_path ./specs  # Define caminho padrão")
	fmt.Println()
	fmt.Println("Chaves disponíveis:")
	fmt.Println("  specs.default_path       Caminho padrão para diretório de specs (string)")
	fmt.Println("  specs.exclude_templates  Excluir specs de template do dashboard (boolean)")
}
