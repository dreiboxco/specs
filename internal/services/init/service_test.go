package init

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreibox/specs/internal/adapters"
)

func TestService_Initialize_Success(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	opts := InitOptions{
		TargetDir: tmpDir,
		Force:     false,
	}

	result, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar que specs/ foi criado
	if !fs.Exists(result.SpecsDir) {
		t.Error("diretório specs/ não foi criado")
	}

	// Verificar que templates foram copiados
	templateNames := []string{
		"00-global-context.spec.md",
		"00-architecture.spec.md",
		"00-stack.spec.md",
		"checklist.md",
		"template-default.spec.md",
	}

	for _, name := range templateNames {
		path := filepath.Join(result.SpecsDir, name)
		if !fs.Exists(path) {
			t.Errorf("template %s não foi copiado", name)
		}
	}

	// Verificar que .cursorrules foi criado
	cursorRulesPath := filepath.Join(tmpDir, ".cursorrules")
	if !fs.Exists(cursorRulesPath) {
		t.Error(".cursorrules não foi criado")
	}

	// Verificar que README.md foi criado
	readmePath := filepath.Join(tmpDir, "README.md")
	if !fs.Exists(readmePath) {
		t.Error("README.md não foi criado")
	}

	// Verificar que arquivos foram adicionados ao resultado
	if len(result.FilesCreated) < 2 {
		t.Errorf("esperava pelo menos 2 arquivos criados, obteve %d", len(result.FilesCreated))
	}

	if len(result.DirectoriesCreated) < 1 {
		t.Error("esperava pelo menos 1 diretório criado")
	}
}

func TestService_Initialize_WithBoilerplate(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	opts := InitOptions{
		TargetDir:       tmpDir,
		Force:           false,
		WithBoilerplate: true,
	}

	result, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar que boilerplate/specs/ foi criado
	boilerplateDir := filepath.Join(tmpDir, "boilerplate", "specs")
	if !fs.Exists(boilerplateDir) {
		t.Error("diretório boilerplate/specs/ não foi criado")
	}

	// Verificar que specsDir foi definido
	if result.SpecsDir == "" {
		t.Error("specsDir não foi definido no resultado")
	}

	// Verificar que templates foram copiados para boilerplate
	templateNames := []string{
		"00-global-context.spec.md",
		"00-architecture.spec.md",
		"00-stack.spec.md",
		"checklist.md",
		"template-default.spec.md",
	}

	for _, name := range templateNames {
		path := filepath.Join(boilerplateDir, name)
		if !fs.Exists(path) {
			t.Errorf("template %s não foi copiado para boilerplate", name)
		}
	}
}

func TestService_Initialize_Idempotent(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	opts := InitOptions{
		TargetDir: tmpDir,
		Force:     false,
	}

	// Primeira inicialização
	result1, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro na primeira inicialização: %v", err)
	}

	// Segunda inicialização (deve ser idempotente)
	result2, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro na segunda inicialização: %v", err)
	}

	// Verificar que segunda inicialização não criou novos arquivos
	if len(result2.FilesCreated) > 0 || len(result2.DirectoriesCreated) > 0 {
		t.Error("segunda inicialização não deveria criar novos arquivos (idempotência)")
	}

	// Verificar que specsDir é o mesmo
	if result1.SpecsDir != result2.SpecsDir {
		t.Errorf("specsDir deveria ser o mesmo em ambas as inicializações: %s != %s", result1.SpecsDir, result2.SpecsDir)
	}

	// Usar result2 para evitar unused
	_ = result2
}

func TestService_Initialize_DirectoryNotExists(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	opts := InitOptions{
		TargetDir: "/caminho/inexistente/12345",
		Force:     false,
	}

	_, err := service.Initialize(opts)
	if err == nil {
		t.Error("esperava erro para diretório inexistente")
	}
}

func TestService_Initialize_WithTargetDir(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "meu-projeto")

	opts := InitOptions{
		TargetDir: targetDir,
		Force:     false,
	}

	// Criar diretório alvo
	if err := fs.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("falha ao criar diretório alvo: %v", err)
	}

	result, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar que specs/ foi criado no diretório alvo
	expectedSpecsDir := filepath.Join(targetDir, "specs")
	if result.SpecsDir != expectedSpecsDir {
		t.Errorf("specsDir esperado %s, obteve %s", expectedSpecsDir, result.SpecsDir)
	}

	if !fs.Exists(expectedSpecsDir) {
		t.Error("diretório specs/ não foi criado no diretório alvo")
	}
}

func TestService_Initialize_ForceOverwrite(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	opts := InitOptions{
		TargetDir: tmpDir,
		Force:     false,
	}

	// Primeira inicialização
	_, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro na primeira inicialização: %v", err)
	}

	// Modificar um arquivo existente
	readmePath := filepath.Join(tmpDir, "README.md")
	modifiedContent := []byte("Conteúdo modificado")
	if err := fs.WriteFile(readmePath, modifiedContent, 0644); err != nil {
		t.Fatalf("falha ao modificar arquivo: %v", err)
	}

	// Verificar que arquivo foi modificado
	contentBefore, _ := fs.ReadFile(readmePath)
	if string(contentBefore) != string(modifiedContent) {
		t.Fatal("falha ao modificar arquivo para teste")
	}

	// Remover specs/ para forçar recriação (não idempotente)
	specsDir := filepath.Join(tmpDir, "specs")
	if err := os.RemoveAll(specsDir); err != nil {
		t.Fatalf("falha ao remover specs/: %v", err)
	}

	// Segunda inicialização com --force
	opts.Force = true
	result, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro na segunda inicialização com --force: %v", err)
	}

	// Verificar que specs foi recriado
	if !fs.Exists(result.SpecsDir) {
		t.Error("specs/ não foi recriado")
	}

	// Verificar que arquivo foi sobrescrito
	newContent, err := fs.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("falha ao ler arquivo: %v", err)
	}

	if string(newContent) == string(modifiedContent) {
		t.Error("arquivo não foi sobrescrito com --force")
	}

	// Verificar que conteúdo foi restaurado ao template (não está vazio e não é o modificado)
	if len(newContent) == 0 {
		t.Error("arquivo está vazio após sobrescrita")
	}
}

func TestService_isSDDProject(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	testCases := []struct {
		name     string
		setup    func(string) // função para configurar o diretório
		expected bool
	}{
		{
			name: "projeto SDD existente",
			setup: func(dir string) {
				specsDir := filepath.Join(dir, "specs")
				fs.MkdirAll(specsDir, 0755)
				// Criar arquivo 00-*.spec.md
				templateFile := filepath.Join(specsDir, "00-global-context.spec.md")
				fs.WriteFile(templateFile, []byte("# Test"), 0644)
			},
			expected: true,
		},
		{
			name: "diretório specs existe mas sem arquivos 00-*",
			setup: func(dir string) {
				specsDir := filepath.Join(dir, "specs")
				fs.MkdirAll(specsDir, 0755)
				// Criar arquivo que não é 00-*
				otherFile := filepath.Join(specsDir, "01-test.spec.md")
				fs.WriteFile(otherFile, []byte("# Test"), 0644)
			},
			expected: false,
		},
		{
			name: "diretório specs não existe",
			setup: func(dir string) {
				// Não criar nada
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tc.setup(tmpDir)

			result := service.isSDDProject(tmpDir)
			if result != tc.expected {
				t.Errorf("isSDDProject retornou %v, esperado %v", result, tc.expected)
			}
		})
	}
}

func TestService_checkWritePermissions(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	err := service.checkWritePermissions(tmpDir)
	if err != nil {
		t.Errorf("erro inesperado ao verificar permissões de escrita: %v", err)
	}

	// Verificar que arquivo de teste foi removido
	testFile := filepath.Join(tmpDir, ".specs-write-test")
	if fs.Exists(testFile) {
		t.Error("arquivo de teste não foi removido após verificação")
	}
}

func TestService_copySpecTemplates(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	err := service.copySpecTemplates(tmpDir, false)
	if err != nil {
		t.Fatalf("erro ao copiar templates: %v", err)
	}

	// Verificar que templates foram copiados
	templateNames := []string{
		"00-global-context.spec.md",
		"00-architecture.spec.md",
		"00-stack.spec.md",
		"checklist.md",
		"template-default.spec.md",
	}

	for _, name := range templateNames {
		path := filepath.Join(tmpDir, name)
		if !fs.Exists(path) {
			t.Errorf("template %s não foi copiado", name)
		}

		// Verificar que arquivo não está vazio
		content, err := fs.ReadFile(path)
		if err != nil {
			t.Errorf("falha ao ler template %s: %v", name, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("template %s está vazio", name)
		}
	}
}

func TestService_createFileIfNotExists(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content")

	// Criar arquivo quando não existe
	err := service.createFileIfNotExists(testFile, content, false)
	if err != nil {
		t.Fatalf("erro ao criar arquivo: %v", err)
	}

	if !fs.Exists(testFile) {
		t.Error("arquivo não foi criado")
	}

	// Tentar criar novamente sem force (não deve sobrescrever)
	modifiedContent := []byte("modified content")
	err = service.createFileIfNotExists(testFile, modifiedContent, false)
	if err != nil {
		t.Fatalf("erro ao tentar criar arquivo existente: %v", err)
	}

	// Verificar que conteúdo não foi modificado
	readContent, _ := fs.ReadFile(testFile)
	if string(readContent) != string(content) {
		t.Error("arquivo foi sobrescrito sem --force")
	}

	// Criar com force (deve sobrescrever)
	err = service.createFileIfNotExists(testFile, modifiedContent, true)
	if err != nil {
		t.Fatalf("erro ao sobrescrever arquivo com --force: %v", err)
	}

	// Verificar que conteúdo foi modificado
	readContent, _ = fs.ReadFile(testFile)
	if string(readContent) != string(modifiedContent) {
		t.Error("arquivo não foi sobrescrito com --force")
	}
}

func TestService_Initialize_EmptyTargetDir(t *testing.T) {
	fs := adapters.NewFileSystem()
	service := NewService(fs)

	tmpDir := t.TempDir()

	// Mudar para diretório temporário
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	opts := InitOptions{
		TargetDir: "", // Vazio - deve usar diretório atual
		Force:     false,
	}

	result, err := service.Initialize(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar que specs/ foi criado no diretório atual
	if !fs.Exists(result.SpecsDir) {
		t.Error("diretório specs/ não foi criado")
	}

	// Normalizar caminhos para comparação (resolver symlinks)
	expectedSpecsDir := filepath.Join(tmpDir, "specs")
	expectedAbs, _ := filepath.Abs(expectedSpecsDir)
	resultAbs, _ := filepath.Abs(result.SpecsDir)
	
	// Resolver symlinks
	expectedResolved, _ := filepath.EvalSymlinks(expectedAbs)
	resultResolved, _ := filepath.EvalSymlinks(resultAbs)
	
	if expectedResolved != resultResolved {
		t.Errorf("specsDir esperado %s, obteve %s", expectedResolved, resultResolved)
	}
}
