package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreibox/specs/internal/adapters"
)

// Service gerencia configuração do CLI
type Service struct {
	fs         adapters.FileSystem
	configPath string // Caminho customizado para testes (vazio = usar XDG)
}

// NewService cria uma nova instância do Service
func NewService(fs adapters.FileSystem) *Service {
	return &Service{
		fs: fs,
	}
}

// NewServiceWithPath cria uma nova instância do Service com caminho customizado (para testes)
func NewServiceWithPath(fs adapters.FileSystem, configPath string) *Service {
	return &Service{
		fs:         fs,
		configPath: configPath,
	}
}

// Config representa a estrutura de configuração
type Config struct {
	Specs SpecsConfig `json:"specs"`
}

// SpecsConfig contém configurações relacionadas a specs
type SpecsConfig struct {
	DefaultPath      string `json:"default_path"`
	ExcludeTemplates bool   `json:"exclude_templates"`
}

// DefaultConfig retorna configuração padrão
func DefaultConfig() *Config {
	return &Config{
		Specs: SpecsConfig{
			DefaultPath:      "./specs",
			ExcludeTemplates: true,
		},
	}
}

// GetConfigPath retorna o caminho do arquivo de configuração (XDG-compliant)
func (s *Service) GetConfigPath() (string, error) {
	// Se configPath foi definido (para testes), usar ele
	if s.configPath != "" {
		return s.configPath, nil
	}

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("falha ao obter diretório home: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	configPath := filepath.Join(configDir, "specs", "config.json")
	return configPath, nil
}

// Load carrega configuração do arquivo ou retorna padrão
func (s *Service) Load() (*Config, error) {
	configPath, err := s.GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Verificar se arquivo existe
	if !s.fs.Exists(configPath) {
		return DefaultConfig(), nil
	}

	// Ler arquivo
	data, err := s.fs.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de configuração: %w", err)
	}

	// Parsear JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("arquivo de configuração inválido: %w", err)
	}

	// Validar e aplicar valores padrão para campos ausentes
	if config.Specs.DefaultPath == "" {
		config.Specs.DefaultPath = DefaultConfig().Specs.DefaultPath
	}

	return &config, nil
}

// Save salva configuração no arquivo
func (s *Service) Save(config *Config) error {
	configPath, err := s.GetConfigPath()
	if err != nil {
		return err
	}

	// Validar configuração antes de salvar
	if err := s.Validate(config); err != nil {
		return err
	}

	// Criar diretório se não existir
	configDir := filepath.Dir(configPath)
	if !s.fs.Exists(configDir) {
		if err := s.fs.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("falha ao criar diretório de configuração: %w", err)
		}
	}

	// Serializar para JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("falha ao serializar configuração: %w", err)
	}

	// Escrever arquivo
	if err := s.fs.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("falha ao salvar configuração: %w", err)
	}

	return nil
}

// Validate valida configuração
func (s *Service) Validate(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuração não pode ser nil")
	}

	// Validar default_path
	if config.Specs.DefaultPath == "" {
		return fmt.Errorf("specs.default_path não pode ser vazio")
	}

	// Valores booleanos já são validados pelo JSON unmarshal
	return nil
}

// GetValue obtém valor de uma chave específica
func (s *Service) GetValue(key string) (interface{}, error) {
	config, err := s.Load()
	if err != nil {
		return nil, err
	}

	// Parsear chave (formato: specs.default_path)
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("chave desconhecida: %s", key)
	}

	namespace := parts[0]
	option := parts[1]

	if namespace != "specs" {
		return nil, fmt.Errorf("chave desconhecida: %s", key)
	}

	switch option {
	case "default_path":
		return config.Specs.DefaultPath, nil
	case "exclude_templates":
		return config.Specs.ExcludeTemplates, nil
	default:
		return nil, fmt.Errorf("chave desconhecida: %s", key)
	}
}

// ResolveDefaultPath resolve o caminho padrão para specs baseado na configuração
func (s *Service) ResolveDefaultPath() (string, error) {
	config, err := s.Load()
	if err != nil {
		return "", err
	}

	defaultPath := config.Specs.DefaultPath

	// Se o caminho é relativo, resolver em relação ao diretório atual
	if !filepath.IsAbs(defaultPath) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("falha ao obter diretório atual: %w", err)
		}
		return filepath.Join(wd, defaultPath), nil
	}

	return defaultPath, nil
}

// SetValue define valor de uma chave específica
func (s *Service) SetValue(key string, value interface{}) error {
	// Carregar configuração existente ou usar padrão
	config, err := s.Load()
	if err != nil {
		return err
	}

	// Parsear chave
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return fmt.Errorf("chave desconhecida: %s", key)
	}

	namespace := parts[0]
	option := parts[1]

	if namespace != "specs" {
		return fmt.Errorf("chave desconhecida: %s", key)
	}

	// Validar e definir valor
	switch option {
	case "default_path":
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("valor inválido para %s: deve ser string", key)
		}
		if strValue == "" {
			return fmt.Errorf("valor inválido para %s: não pode ser vazio", key)
		}
		config.Specs.DefaultPath = strValue
	case "exclude_templates":
		boolValue, ok := value.(bool)
		if !ok {
			// Tentar converter string para bool
			if strValue, ok := value.(string); ok {
				switch strings.ToLower(strValue) {
				case "true", "1", "yes":
					boolValue = true
				case "false", "0", "no":
					boolValue = false
				default:
					return fmt.Errorf("valor inválido para %s: deve ser boolean (true/false)", key)
				}
			} else {
				return fmt.Errorf("valor inválido para %s: deve ser boolean", key)
			}
		}
		config.Specs.ExcludeTemplates = boolValue
	default:
		return fmt.Errorf("chave desconhecida: %s", key)
	}

	// Salvar configuração
	return s.Save(config)
}
