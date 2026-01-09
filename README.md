<p align="center">
  <a href="https://github.com/dreibox/specs">
    <h1>Specs</h1>
  </a>
</p>

<p align="center">Spec-driven development CLI for managing SDD projects.</p>

<p align="center">
  <a href="https://github.com/dreibox/specs/actions"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/dreibox/specs/ci.yml?style=flat-square" /></a>
  <a href="https://golang.org/"><img alt="Go version" src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go" /></a>
  <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" /></a>
  <a href="https://conventionalcommits.org"><img alt="Conventional Commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=flat-square" /></a>
  <a href="https://github.com/dreibox/specs/releases"><img alt="GitHub release" src="https://img.shields.io/github/v/release/dreibox/specs?style=flat-square" /></a>
  <a href="https://github.com/dreibox/specs"><img alt="Platform" src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey.svg?style=flat-square" /></a>
</p>

---

## O que é Specs?

**Specs** é uma ferramenta de linha de comando que facilita a criação, validação e gerenciamento de projetos desenvolvidos com **SDD (Spec Driven Development)**. Com Specs, você pode:

- ✅ Inicializar novos projetos seguindo estrutura SDD padronizada
- ✅ Validar especificações contra checklist formal
- ✅ Listar e visualizar status de todas as specs do projeto
- ✅ Verificar consistência estrutural (numeração, links, referências)
- ✅ Gerenciar configuração personalizada do CLI
- ✅ Visualizar dashboard interativo com progresso do projeto
- ✅ Gerenciar o ciclo de vida completo de especificações

## Instalação

### Via GitHub Releases (Recomendado)

```bash
# Download do binário para sua plataforma
# macOS (arm64)
curl -L https://github.com/dreibox/specs/releases/latest/download/specs-darwin-arm64.tar.gz | tar -xz
sudo mv specs /usr/local/bin/

# macOS (x64)
curl -L https://github.com/dreibox/specs/releases/latest/download/specs-darwin-amd64.tar.gz | tar -xz
sudo mv specs /usr/local/bin/

# Linux (arm64)
curl -L https://github.com/dreibox/specs/releases/latest/download/specs-linux-arm64.tar.gz | tar -xz
sudo mv specs /usr/local/bin/

# Linux (x64)
curl -L https://github.com/dreibox/specs/releases/latest/download/specs-linux-amd64.tar.gz | tar -xz
sudo mv specs /usr/local/bin/

# Verificar instalação
specs version
```

### Build Local

```bash
git clone https://github.com/dreibox/specs.git
cd specs
go build -o bin/specs ./cmd/specs
sudo mv bin/specs /usr/local/bin/
```

## Quick Start

### 1. Inicializar um Projeto SDD

Crie a estrutura base de um novo projeto seguindo SDD:

```bash
specs init
```

**O que acontece durante a inicialização:**

* Cria diretório `specs/` com templates de specs base
* Copia templates de especificações (`00-global-context.spec.md`, `00-architecture.spec.md`, `00-stack.spec.md`)
* Cria arquivo `checklist.md` para validação
* Cria arquivo `.cursorrules` base para desenvolvimento
* Cria `README.md` inicial com estrutura do projeto

**Após a inicialização:**

* Preencha os arquivos `00-*.spec.md` com informações do seu projeto
* Use `template-default.spec.md` como base para criar novas specs
* Execute `specs validate` para verificar se suas specs estão completas

### 2. Configurar o CLI (Opcional)

Personalize o comportamento do CLI através de arquivo de configuração:

```bash
# Visualizar configuração atual
specs config

# Obter valor de uma opção
specs config get specs.default_path

# Definir caminho padrão para specs
specs config set specs.default_path ./minhas-specs

# Configurar exclusão de templates no dashboard
specs config set specs.exclude_templates false
```

**Localização da configuração:**
- `~/.config/specs/config.json` (ou `$XDG_CONFIG_HOME/specs/config.json`)

**Opções disponíveis:**
- `specs.default_path`: Caminho padrão para diretório de specs (padrão: `./specs`)
- `specs.exclude_templates`: Excluir specs de template do dashboard (padrão: `true`)

**Exemplo de configuração:**

```json
{
  "specs": {
    "default_path": "./documentation/specs",
    "exclude_templates": true
  }
}
```

### 3. Validar Especificações

Verifique se suas specs estão completas e prontas para implementação:

```bash
# Validar todas as specs no diretório padrão (ou configurado)
specs validate

# Validar diretório específico
specs validate ./minhas-specs

# Validar arquivo único
specs validate specs/01-feature.spec.md
```

**O que é validado:**

* ✅ Presença de todas as seções obrigatórias (1-12)
* ✅ Formato do checklist (6 itens)
* ✅ Completude do checklist (specs completas/incompletas)
* ✅ Estrutura e formato de arquivos Markdown

### 4. Listar Specs com Status

Visualize todas as specs do projeto e seu status:

```bash
# Listar todas as specs
specs list

# Listar apenas specs completas
specs list --complete

# Listar apenas specs incompletas
specs list --incomplete

# Listar specs com erros
specs list --errors
```

**Output exemplo:**

```
Numeração  Nome                    Status
──────────  ──────────────────────  ──────────
00          global-context          ✅ Completa
00          architecture            ✅ Completa
00          stack                   ✅ Completa
01          version-control         ✅ Completa
02          init                    ⚠️  Incompleta (4/6)
03          specs-validate          ⚠️  Incompleta (3/6)

Resumo:
  Total: 6 specs
  Completas: 4
  Incompletas: 2
  Com erros: 0
```

### 5. Verificar Consistência Estrutural

Valide numeração, links e referências entre specs:

```bash
# Verificar consistência em specs/
specs check

# Verificar diretório específico
specs check ./minhas-specs
```

**O que é verificado:**

* ✅ Numeração sequencial (detecta gaps e duplicatas)
* ✅ Links internos válidos (detecta links quebrados)
* ✅ Specs órfãs (referenciadas mas não existem)
* ✅ Formato de nomes de arquivos
* ✅ Estrutura de diretórios

### 6. Visualizar Dashboard

Exiba um dashboard interativo com informações agregadas do projeto:

```bash
# Dashboard de specs/ no diretório atual (ou configurado)
specs view

# Dashboard de diretório específico
specs view ./minhas-specs
```

**O que é exibido:**

* 📊 **Summary**: Total de specs, requirements, progresso geral
* 📈 **Specs em Progresso**: Lista com barras de progresso visuais
* ✅ **Specs Completas**: Lista de specs finalizadas
* 📋 **Specifications**: Lista completa com contagem de requirements por spec

**Output exemplo:**

```
Specs Dashboard

Summary:
  Specifications: 10 specs, 64 requirements
  Specs em Progresso: 3
  Specs Completas: 4
  Progresso Geral: 30/41 (73% complete)

Specs em Progresso:
  make-validation-scope-aware        [          ] 0%
  remove-diff-command                 [█████████ ] 90%
  improve-deterministic-tests        [█████████ ] 92%

Specs Completas:
  ✅ add-slash-command-support
  ✅ sort-active-changes-by-progress
  ✅ update-agent-file-name
  ✅ update-agent-instructions

Specifications:
  cli-archive              10 requirements
  openspec-conventions     10 requirements
  cli-validate              9 requirements
  ...
```

## Comandos Disponíveis

### `specs init [diretório]`

Inicializa um novo projeto SDD criando estrutura base e templates.

**Flags:**
- `--force`: Sobrescreve arquivos existentes sem confirmação
- `--with-boilerplate`: Cria também diretório `boilerplate/` com templates genéricos

**Exemplos:**
```bash
specs init                    # Inicializa no diretório atual
specs init ./meu-projeto      # Inicializa em diretório específico
specs init --force            # Sobrescreve arquivos existentes
```

### `specs config [subcomando]`

Gerencia configuração do CLI. Permite personalizar comportamento padrão.

**Subcomandos:**
- `show` (padrão): Exibe configuração atual
- `get <chave>`: Obtém valor de uma chave específica
- `set <chave> <valor>`: Define valor de uma chave

**Flags:**
- `--help`: Exibe ajuda do comando

**Chaves disponíveis:**
- `specs.default_path`: Caminho padrão para diretório de specs (string, padrão: `./specs`)
- `specs.exclude_templates`: Excluir specs de template do dashboard (boolean, padrão: `true`)

**Exemplos:**
```bash
# Exibir configuração completa
specs config
specs config show

# Obter valor de uma chave
specs config get specs.default_path
specs config get specs.exclude_templates

# Definir valores
specs config set specs.default_path ./documentation/specs
specs config set specs.exclude_templates false

# Ajuda
specs config --help
```

**Localização da configuração:**
- `~/.config/specs/config.json` (ou `$XDG_CONFIG_HOME/specs/config.json` se `XDG_CONFIG_HOME` estiver definido)
- Arquivo criado automaticamente quando você define valores
- Permissões: 0600 (apenas leitura/escrita pelo dono)

**Formato do arquivo de configuração:**

```json
{
  "specs": {
    "default_path": "./specs",
    "exclude_templates": true
  }
}
```

**Notas:**
- Valores padrão são aplicados quando arquivo não existe
- Configuração é validada ao carregar (formato JSON e tipos)
- Erros de configuração resultam em fallback para valores padrão
- Todos os comandos que aceitam caminho usam `specs.default_path` quando não especificado

### `specs validate [caminho]`

Valida specs contra checklist formal e verifica estrutura.

**Exemplos:**
```bash
specs validate                    # Valida specs/ no diretório atual (ou configurado)
specs validate specs/             # Valida diretório específico
specs validate specs/01-test.spec.md  # Valida arquivo único
```

**O que é validado:**
- Presença de todas as seções obrigatórias (1-12)
- Formato do checklist (6 itens)
- Completude do checklist
- Estrutura e formato de arquivos Markdown

**Códigos de saída:**
- `0`: Sucesso (sem erros)
- `1`: Erros encontrados
- `2`: Erro de input inválido

### `specs list [caminho]`

Lista todas as specs com status (completa/incompleta/erro).

**Flags:**
- `--complete`, `--only-complete`: Lista apenas specs completas
- `--incomplete`, `--only-incomplete`: Lista apenas specs incompletas
- `--errors`: Lista apenas specs com erros

**Exemplos:**
```bash
specs list                    # Lista todas as specs
specs list --complete         # Apenas specs completas
specs list --incomplete       # Apenas specs incompletas
specs list --errors           # Apenas specs com erros
specs list specs/             # Lista specs em diretório específico
```

### `specs check [caminho]`

Verifica consistência estrutural (numeração, links, referências).

**Exemplos:**
```bash
specs check                   # Verifica specs/ no diretório atual (ou configurado)
specs check specs/            # Verifica diretório específico
```

**O que é verificado:**
- Numeração sequencial (detecta gaps e duplicatas)
- Links internos válidos (detecta links quebrados)
- Specs órfãs (referenciadas mas não existem)
- Formato de nomes de arquivos
- Estrutura de diretórios

**Códigos de saída:**
- `0`: Sem problemas encontrados
- `1`: Problemas encontrados
- `2`: Erro de input inválido

### `specs view [caminho]`

Exibe dashboard interativo com informações agregadas do projeto SDD.

**Exemplos:**
```bash
specs view                    # Dashboard de specs/ no diretório atual (ou configurado)
specs view specs/             # Dashboard de diretório específico
```

**O que é exibido:**
- **Summary**: Total de specs, requirements, progresso geral
- **Specs em Progresso**: Lista com barras de progresso visuais
- **Specs Completas**: Lista de specs finalizadas
- **Specifications**: Lista completa com contagem de requirements por spec

**Notas:**
- Respeita configuração `specs.exclude_templates` (exclui `00-*.spec.md` e `template-default.spec.md` por padrão)
- Calcula progresso baseado em itens do checklist marcados

### `specs version`

Exibe a versão atual do CLI.

**Exemplos:**
```bash
specs version
# 0.0.3
```

## Configuração

O CLI Specs suporta configuração personalizada através de arquivo JSON em localização XDG-compliant.

### Localização

- **Linux/macOS**: `~/.config/specs/config.json`
- **Com XDG_CONFIG_HOME**: `$XDG_CONFIG_HOME/specs/config.json`

### Opções de Configuração

#### `specs.default_path`

Caminho padrão para diretório de specs usado quando nenhum caminho é especificado nos comandos.

- **Tipo**: string
- **Padrão**: `"./specs"`
- **Exemplo**: `"./documentation/specs"`

**Uso:**
```bash
specs config set specs.default_path ./minhas-specs
```

#### `specs.exclude_templates`

Controla se specs de template devem ser excluídas do dashboard e cálculos de progresso.

- **Tipo**: boolean
- **Padrão**: `true`
- **Valores**: `true` ou `false`

**Uso:**
```bash
specs config set specs.exclude_templates false
```

**Specs excluídas quando `true`:**
- Arquivos com prefixo `00-*` (ex: `00-architecture.spec.md`)
- Arquivo `template-default.spec.md`

### Exemplo Completo de Configuração

```json
{
  "specs": {
    "default_path": "./documentation/specs",
    "exclude_templates": true
  }
}
```

### Gerenciamento de Configuração

```bash
# Visualizar configuração atual
specs config

# Obter valor específico
specs config get specs.default_path

# Definir valores
specs config set specs.default_path ./custom-path
specs config set specs.exclude_templates false

# Remover arquivo para voltar aos padrões
rm ~/.config/specs/config.json
```

### Valores Padrão

Quando o arquivo de configuração não existe, os seguintes valores padrão são aplicados:

- `specs.default_path`: `"./specs"`
- `specs.exclude_templates`: `true`

## Estrutura de Projeto SDD

Após executar `specs init`, sua estrutura será:

```
projeto/
├── specs/                          # Diretório de especificações
│   ├── 00-global-context.spec.md  # Contexto, visão, objetivos
│   ├── 00-architecture.spec.md    # Arquitetura e padrões
│   ├── 00-stack.spec.md           # Stack técnica
│   ├── 01-feature.spec.md         # Specs de funcionalidades
│   ├── 02-outra-feature.spec.md
│   ├── checklist.md               # Checklist de validação
│   └── template-default.spec.md   # Template para novas specs
├── .cursorrules                   # Regras do Cursor para SDD
└── README.md                       # Documentação do projeto
```

## Metodologia SDD

**Spec Driven Development (SDD)** é uma metodologia onde:

1. **Especificar antes de codificar**: Todas as funcionalidades são especificadas em arquivos `*.spec.md` antes da implementação
2. **Validar contra checklist**: Cada spec deve passar por um checklist formal antes de ser implementada
3. **Implementar conforme spec**: O código implementa exatamente o que está especificado, sem adicionar funcionalidades não especificadas
4. **Manter specs atualizadas**: As specs são a fonte da verdade e devem ser mantidas atualizadas

### Fluxo de Trabalho SDD

```
┌─────────────┐
│ Especificar │  ← Criar spec seguindo template
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Validar    │  ← specs validate (verificar checklist)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Implementar │  ← Codificar conforme spec
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Testar    │  ← Validar critérios de aceite
└─────────────┘
```

## Convenções de Specs

### Numeração

- `00-*`: Specs base (contexto global, arquitetura, stack)
- `01-*`, `02-*`, `03-*`, etc.: Specs de funcionalidades (sequencial)

### Formato

- Arquivos: `{numero}-{nome-descritivo}.spec.md`
- Encoding: UTF-8
- Formato: Markdown com seções padronizadas
- Checklist: Sempre no final, após "Abertos / Fora de Escopo"

### Seções Obrigatórias

Toda spec deve conter:

1. Contexto e Objetivo
2. Requisitos Funcionais
3. Contratos e Interfaces
4. Fluxos e Estados
5. Dados
6. NFRs (Não Funcionais)
7. Guardrails
8. Critérios de Aceite
9. Testes
10. Migração / Rollback
11. Observações Operacionais
12. Abertos / Fora de Escopo
13. Checklist Rápido

## Desenvolvimento

### Pré-requisitos

- Go 1.25+ ou superior
- Git

### Build Local

```bash
# Build de desenvolvimento
go build -o bin/specs ./cmd/specs

# Build de produção (otimizado)
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/specs ./cmd/specs

# Executar sem build
go run ./cmd/specs
```

### Testes

```bash
# Todos os testes
go test ./...

# Com cobertura
go test -cover ./...

# Testes verbosos
go test -v ./...

# Testes de um pacote específico
go test ./internal/services/config/...
```

### Estrutura do Código

```
specs/
├── cmd/
│   └── specs/            # Entry point
├── internal/
│   ├── cli/              # Parser, roteamento
│   ├── commands/         # Comandos
│   ├── services/        # Lógica de negócio
│   │   ├── config/      # Serviço de configuração
│   │   ├── validator/   # Validação de specs
│   │   ├── lister/      # Listagem de specs
│   │   ├── checker/     # Verificação estrutural
│   │   ├── viewer/      # Dashboard
│   │   └── init/        # Inicialização de projetos
│   ├── adapters/        # I/O abstrato
│   └── templates/       # Templates de arquivos
├── specs/               # Especificações do projeto
└── boilerplate/         # Templates para novos projetos
```

## Compatibilidade

- **Plataformas**: macOS (10.15+), Linux (glibc 2.17+)
- **Arquiteturas**: x64 (amd64), arm64
- **Runtime**: Binário estático único, sem dependências externas

## Contribuindo

Este projeto segue metodologia SDD. Para contribuir:

1. **Revise as specs** em `specs/` para entender o que precisa ser implementado
2. **Valide a spec** contra o checklist antes de começar a codificar
3. **Implemente conforme a spec** - não adicione funcionalidades não especificadas
4. **Mantenha specs atualizadas** quando fizer mudanças

### Processo de Contribuição

1. Consulte `BACKLOG.md` para ver specs pendentes
2. Escolha uma spec para implementar
3. Crie branch `feature/{numero}-{nome-sem-extensao}`
4. Implemente conforme a spec
5. Execute testes e validações
6. Marque checklist da spec como completo
7. Abra Pull Request

### Convenções de Commit

Este projeto usa [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: adiciona nova funcionalidade
fix: corrige bug
docs: atualiza documentação
style: formatação de código
refactor: refatoração
test: adiciona testes
chore: tarefas de manutenção
```

## Licença

MIT License - veja [LICENSE](./LICENSE) para detalhes.

## Referências

- [Go Documentation](https://go.dev/doc/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- Arquitetura e padrões: `specs/00-architecture.spec.md`
- Stack técnica: `specs/00-stack.spec.md`
- Contexto global: `specs/00-global-context.spec.md`
