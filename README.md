# Specs

**Spec-driven development (SDD) CLI** para gerenciar projetos que seguem metodologia de especificação antes de implementação.

[![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey.svg)](https://github.com/dreibox/specs)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## O que é Specs?

**Specs** é uma ferramenta de linha de comando que facilita a criação, validação e gerenciamento de projetos desenvolvidos com **SDD (Spec Driven Development)**. Com Specs, você pode:

- ✅ Inicializar novos projetos seguindo estrutura SDD padronizada
- ✅ Validar especificações contra checklist formal
- ✅ Listar e visualizar status de todas as specs do projeto
- ✅ Verificar consistência estrutural (numeração, links, referências)
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

### 2. Validar Especificações

Verifique se suas specs estão completas e prontas para implementação:

```bash
# Validar todas as specs no diretório padrão (specs/)
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

### 3. Listar Specs com Status

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

### 4. Verificar Consistência Estrutural

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

### 5. Visualizar Dashboard

Exiba um dashboard interativo com informações agregadas do projeto:

```bash
# Dashboard de specs/ no diretório atual
specs view

# Dashboard de diretório específico
specs view ./minhas-specs
```

**O que é exibido:**

* 📊 **Summary**: Total de specs, requirements, progresso geral
* 📈 **Specs em Progresso**: Lista com barras de progresso visuais
* ✅ **Specs Completas**: Lista de specs finalizadas
* 📋 **Specifications**: Lista completa com contagem de requirements por spec

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

### `specs validate [caminho]`

Valida specs contra checklist formal e verifica estrutura.

**Exemplos:**
```bash
specs validate                    # Valida specs/ no diretório atual
specs validate specs/             # Valida diretório específico
specs validate specs/01-test.spec.md  # Valida arquivo único
```

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
specs list --incomplete      # Apenas specs incompletas
specs list specs/             # Lista specs em diretório específico
```

### `specs check [caminho]`

Verifica consistência estrutural (numeração, links, referências).

**Exemplos:**
```bash
specs check                   # Verifica specs/ no diretório atual
specs check specs/            # Verifica diretório específico
```

### `specs view [caminho]`

Exibe dashboard interativo com informações agregadas do projeto SDD.

**Exemplos:**
```bash
specs view                    # Dashboard de specs/ no diretório atual
specs view specs/             # Dashboard de diretório específico
```

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

### `specs version`

Exibe a versão atual do CLI.

**Exemplos:**
```bash
specs version
# 0.0.3
```

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

- Go 1.25.5 ou superior
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
```

### Estrutura do Código

```
specs/
├── cmd/
│   └── specs/            # Entry point
├── internal/
│   ├── cli/              # Parser, roteamento
│   ├── commands/         # Comandos
│   ├── services/         # Lógica de negócio
│   ├── adapters/         # I/O abstrato
│   └── config/           # Configuração
├── pkg/                  # Código exportável
├── testdata/             # Arquivos de teste
└── specs/                # Especificações do projeto
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

## Licença

MIT

## Referências

- [Go Documentation](https://go.dev/doc/)
- Arquitetura e padrões: `specs/00-architecture.spec.md`
- Stack técnica: `specs/00-stack.spec.md`
- Contexto global: `specs/00-global-context.spec.md`
