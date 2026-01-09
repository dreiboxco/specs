package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

// mockFileSystem com caminho de configuração customizável
type mockFileSystem struct {
	adapters.FileSystem
	configPath string
}

func TestService_GetConfigPath(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	path, err := service.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() retornou erro: %v", err)
	}

	// Verificar que o caminho contém .config/specs/config.json
	if !filepath.IsAbs(path) {
		t.Errorf("GetConfigPath() retornou caminho relativo: %s", path)
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("GetConfigPath() não retorna config.json: %s", path)
	}
}

func TestService_Load_FileNotExists(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	config, err := service.Load()
	if err != nil {
		t.Fatalf("Load() retornou erro: %v", err)
	}

	// Verificar valores padrão
	if config.Specs.DefaultPath != "./specs" {
		t.Errorf("DefaultPath esperado './specs', obtido '%s'", config.Specs.DefaultPath)
	}

	if !config.Specs.ExcludeTemplates {
		t.Errorf("ExcludeTemplates esperado true, obtido %v", config.Specs.ExcludeTemplates)
	}
}

func TestService_Load_ValidFile(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "specs-config-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo de configuração válido
	configPath := filepath.Join(tmpDir, "config.json")
	configData := Config{
		Specs: SpecsConfig{
			DefaultPath:      "./minhas-specs",
			ExcludeTemplates: false,
		},
	}

	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatalf("Falha ao serializar configuração: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Falha ao escrever arquivo: %v", err)
	}

	// Criar serviço com caminho temporário
	fs := adapters.NewFileSystem()
	service := NewServiceWithPath(fs, configPath)

	config, err := service.Load()
	if err != nil {
		t.Fatalf("Load() retornou erro: %v", err)
	}

	if config.Specs.DefaultPath != "./minhas-specs" {
		t.Errorf("DefaultPath esperado './minhas-specs', obtido '%s'", config.Specs.DefaultPath)
	}

	if config.Specs.ExcludeTemplates {
		t.Errorf("ExcludeTemplates esperado false, obtido %v", config.Specs.ExcludeTemplates)
	}
}

func TestService_Load_InvalidJSON(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "specs-config-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo de configuração inválido
	configPath := filepath.Join(tmpDir, "config.json")
	invalidJSON := `{ "specs": { "default_path": } }`

	if err := os.WriteFile(configPath, []byte(invalidJSON), 0600); err != nil {
		t.Fatalf("Falha ao escrever arquivo: %v", err)
	}

	fs := adapters.NewFileSystem()
	service := NewServiceWithPath(fs, configPath)

	_, err = service.Load()
	if err == nil {
		t.Error("Load() deveria retornar erro para JSON inválido")
	}
}

func TestService_Save(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "specs-config-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")

	fs := adapters.NewFileSystem()
	service := NewServiceWithPath(fs, configPath)

	config := &Config{
		Specs: SpecsConfig{
			DefaultPath:      "./test-specs",
			ExcludeTemplates: true,
		},
	}

	if err := service.Save(config); err != nil {
		t.Fatalf("Save() retornou erro: %v", err)
	}

	// Verificar que arquivo foi criado
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Arquivo de configuração não foi criado")
	}

	// Verificar permissões (aproximadamente)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Falha ao obter informações do arquivo: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Permissões esperadas 0600, obtidas %o", info.Mode().Perm())
	}

	// Verificar conteúdo
	loadedConfig, err := service.Load()
	if err != nil {
		t.Fatalf("Falha ao carregar configuração salva: %v", err)
	}

	if loadedConfig.Specs.DefaultPath != "./test-specs" {
		t.Errorf("DefaultPath esperado './test-specs', obtido '%s'", loadedConfig.Specs.DefaultPath)
	}
}

func TestService_Validate(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "config válido",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "config com default_path vazio",
			config: &Config{
				Specs: SpecsConfig{
					DefaultPath:      "",
					ExcludeTemplates: true,
				},
			},
			wantErr: true,
		},
		{
			name:    "config nil",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() erro = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetValue(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "chave válida default_path",
			key:     "specs.default_path",
			wantErr: false,
		},
		{
			name:    "chave válida exclude_templates",
			key:     "specs.exclude_templates",
			wantErr: false,
		},
		{
			name:    "chave desconhecida",
			key:     "specs.unknown",
			wantErr: true,
		},
		{
			name:    "namespace desconhecido",
			key:     "unknown.option",
			wantErr: true,
		},
		{
			name:    "formato inválido",
			key:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetValue(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() erro = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SetValue(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "specs-config-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")

	fs := adapters.NewFileSystem()
	service := NewServiceWithPath(fs, configPath)

	tests := []struct {
		name    string
		key     string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "definir default_path válido",
			key:     "specs.default_path",
			value:   "./new-specs",
			wantErr: false,
		},
		{
			name:    "definir exclude_templates true",
			key:     "specs.exclude_templates",
			value:   true,
			wantErr: false,
		},
		{
			name:    "definir exclude_templates false",
			key:     "specs.exclude_templates",
			value:   false,
			wantErr: false,
		},
		{
			name:    "chave desconhecida",
			key:     "specs.unknown",
			value:   "value",
			wantErr: true,
		},
		{
			name:    "default_path vazio",
			key:     "specs.default_path",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SetValue(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetValue() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verificar que valor foi salvo
				value, err := service.GetValue(tt.key)
				if err != nil {
					t.Errorf("GetValue() após SetValue() retornou erro: %v", err)
				}
				if value != tt.value {
					t.Errorf("Valor salvo não corresponde: esperado %v, obtido %v", tt.value, value)
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Specs.DefaultPath != "./specs" {
		t.Errorf("DefaultPath esperado './specs', obtido '%s'", config.Specs.DefaultPath)
	}

	if !config.Specs.ExcludeTemplates {
		t.Errorf("ExcludeTemplates esperado true, obtido %v", config.Specs.ExcludeTemplates)
	}
}
