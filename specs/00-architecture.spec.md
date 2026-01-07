# 00 - Especificação de Arquitetura

Esta especificação define o padrão arquitetural, estrutura de diretórios, isolamento de módulos e decisões de design do sistema. Use-a como blueprint para implementar features e garantir consistência arquitetural.

## 1. Contexto e Objetivo

### 1.1 Contexto
- **Referência:** Contexto global do projeto, visão, objetivos e escopo estão em `00-global-context.spec.md`.
- **Stack técnica:** Detalhes de linguagem, ferramentas e build estão em `00-stack.spec.md`.

### 1.2 Objetivo
- Estabelecer padrão arquitetural claro e testável para CLI em Go
- Definir estrutura de diretórios e organização do código seguindo convenções Go
- Garantir isolamento e testabilidade dos módulos (comandos, serviços, adapters)
- Estabelecer convenções arquiteturais específicas para CLI SDD

## 2. Padrão Arquitetural

### 2.1 Forma Arquitetural
- **Padrão:** Arquitetura em camadas com separação clara de responsabilidades (similar a Clean Architecture adaptada para CLI)
- **Estrutura:**
  - **Camada de entrada:** CLI parser e roteamento de comandos
  - **Camada de aplicação:** Comandos que orquestram operações
  - **Camada de domínio:** Lógica de negócio e validação de specs
  - **Camada de infraestrutura:** I/O (sistema de arquivos, leitura/escrita)
- **Justificativa:** 
  - Simplicidade para CLI pequeno/médio
  - Testabilidade através de interfaces e mocks
  - Manutenibilidade com separação clara de responsabilidades
  - Escalabilidade para adicionar novos comandos sem afetar existentes

### 2.2 Módulos e Componentes
- **`cmd/specs/`:** Entry point da aplicação, inicialização e setup
- **`internal/cli/`:** Parser de argumentos, roteamento de comandos, help
- **`internal/commands/`:** Implementação de cada comando (init, specs validate, specs list, etc.)
- **`internal/services/`:** Lógica de negócio (validação de specs, geração de artefatos, etc.)
- **`internal/adapters/`:** Abstrações de I/O (FileSystem, Config, etc.)
- **`internal/config/`:** Gerenciamento de configuração (carregamento, validação, XDG)
- **`pkg/`:** Código exportável (se necessário para bibliotecas futuras)
- **Responsabilidades:**
  - **cli:** Parsing, roteamento, formatação de output
  - **commands:** Orquestração de operações, validação de inputs
  - **services:** Lógica de negócio pura (validação, transformação, geração)
  - **adapters:** Abstração de I/O para testabilidade

### 2.3 Isolamento e Dependências
- **Isolamento:** 
  - Comandos não acessam sistema de arquivos diretamente, sempre via interfaces
  - Lógica de negócio não conhece detalhes de I/O
  - Adaptadores implementam interfaces definidas em services
- **Injeção de dependência:** 
  - Dependências injetadas via construtores de structs
  - Interfaces definidas no mesmo pacote que as consomem (ou em pacote compartilhado)
  - Facilita criação de mocks para testes
- **Testabilidade:** 
  - Interfaces para FileSystem, Config, Logger
  - Mocks implementados manualmente ou via interfaces vazias
  - Testes de unidade isolados sem I/O real
  - Testes de integração com sistema de arquivos temporário

## 3. Estrutura de Diretórios

### 3.1 Estrutura Base
```
specs/
├── cmd/
│   └── specs/            # Entry point da aplicação
│       └── main.go
├── internal/
│   ├── cli/              # Parser de argumentos, roteamento, help
│   │   ├── parser.go
│   │   ├── router.go
│   │   └── help.go
│   ├── commands/         # Implementação de comandos
│   │   ├── init.go       # Comando init
│   │   ├── specs/        # Comandos de specs
│   │   │   ├── validate.go
│   │   │   ├── list.go
│   │   │   └── check.go
│   │   └── version.go     # Comando version
│   ├── services/         # Lógica de negócio
│   │   ├── validator.go  # Validação de specs
│   │   ├── generator.go  # Geração de artefatos
│   │   └── checker.go    # Verificação de consistência
│   ├── adapters/         # Abstrações de I/O
│   │   ├── filesystem.go # Interface e implementação de filesystem
│   │   └── config.go     # Interface e implementação de config
│   └── config/           # Gerenciamento de configuração
│       ├── loader.go     # Carregamento de config
│       └── xdg.go        # Helpers XDG
├── pkg/                  # Código exportável (se necessário)
├── testdata/             # Arquivos de teste e fixtures
│   └── specs/            # Exemplos de specs para testes
├── go.mod
├── go.sum
└── README.md
```

### 3.2 Convenções de Organização
- **Um arquivo por comando:** Cada comando em arquivo separado (ex.: `init.go`, `validate.go`)
- **Agrupamento por feature:** Comandos relacionados agrupados em subdiretórios (ex.: `commands/specs/`)
- **Nomenclatura:**
  - Arquivos: `snake_case.go` ou `camelCase.go` (seguir convenção Go)
  - Diretórios: `lowercase`
  - Structs: `PascalCase`
  - Interfaces: `PascalCase` (geralmente terminando com `er` se apropriado, ex.: `FileSystem`, `Validator`)
  - Funções: `PascalCase` para exportadas, `camelCase` para privadas

## 4. Padrões de Design

### 4.1 Padrões Aplicados
- **Adapter Pattern:** 
  - Abstração de I/O (FileSystem, Config) para permitir mocks em testes
  - Justificativa: Testabilidade sem dependência de sistema de arquivos real
- **Command Pattern:** 
  - Cada comando é uma struct com método `Execute()`
  - Justificativa: Isolamento de comandos, fácil adicionar novos comandos
- **Dependency Injection:** 
  - Dependências injetadas via construtores
  - Justificativa: Facilita testes e reduz acoplamento
- **Strategy Pattern (futuro):** 
  - Diferentes estratégias de validação ou geração
  - Justificativa: Extensibilidade para diferentes tipos de specs

### 4.2 Abstrações e Interfaces
- **`FileSystem` interface:**
  - Responsabilidade: Abstração de operações de sistema de arquivos
  - Métodos: `ReadFile(path)`, `WriteFile(path, data)`, `Exists(path)`, `MkdirAll(path)`, `Walk(root, walkFn)`
  - Mockabilidade: Implementação mock para testes sem I/O real
- **`Config` interface:**
  - Responsabilidade: Abstração de carregamento e salvamento de configuração
  - Métodos: `Load()`, `Save()`, `Get(key)`, `Set(key, value)`
  - Mockabilidade: Implementação mock para testes
- **`Validator` interface (futuro):**
  - Responsabilidade: Validação de specs contra regras
  - Métodos: `Validate(spec)`, `ValidateAll(specs)`
  - Mockabilidade: Implementação mock para testes

## 5. Fluxo de Dados

### 5.1 Fluxo Principal
- **Entrada:** Argumentos de linha de comando parseados pelo CLI
- **Roteamento:** CLI identifica comando e delega para implementação correspondente
- **Validação de Input:** Comando valida argumentos e flags
- **Processamento:** Comando chama serviços que executam lógica de negócio
- **I/O:** Serviços usam adapters para ler/escrever arquivos
- **Saída:** Resultado formatado e enviado para stdout/stderr
- **Exemplo (comando `specs validate`):**
  1. CLI parseia `specs validate ./specs`
  2. Router identifica comando `specs validate`
  3. Comando `ValidateCommand` valida argumentos
  4. Serviço `Validator` lê arquivos via `FileSystem` adapter
  5. Serviço valida cada spec contra regras
  6. Resultado formatado e impresso em stdout

### 5.2 Tratamento de Erros
- **Tipos de erro:**
  - `ErrInvalidInput`: Argumentos ou flags inválidos (código 2)
  - `ErrFileNotFound`: Arquivo não encontrado (código 1)
  - `ErrValidationFailed`: Validação de spec falhou (código 1)
  - `ErrConfig`: Erro de configuração (código 1)
- **Propagação:** 
  - Erros retornados como `error` do Go
  - Não usar panic exceto para erros de programação
  - Erros logados em stderr com mensagem clara
- **Logging:** 
  - Erros sempre logados em stderr
  - Mensagens de erro acionáveis com sugestões
  - Flag `--debug` para stack traces detalhados

## 6. Convenções Arquiteturais

### 6.1 Separação de Responsabilidades
- **Comandos (`commands/`):** 
  - Apenas orquestração e validação de inputs
  - Não contêm lógica de negócio
  - Delegam para serviços
- **Serviços (`services/`):** 
  - Contêm toda lógica de negócio
  - Não conhecem detalhes de I/O (usam interfaces)
  - Testáveis sem I/O real
- **Adaptadores (`adapters/`):** 
  - Implementam interfaces de I/O
  - Única camada que acessa sistema de arquivos diretamente
  - Podem ser substituídos por mocks em testes

### 6.2 Acesso a Recursos
- **Sistema de arquivos:** Sempre via interface `FileSystem`, nunca `os.Open`, `ioutil.ReadFile` diretamente em services
- **Configuração:** Sempre via interface `Config`, nunca leitura direta de arquivos JSON
- **Regra geral:** Services e commands nunca importam pacotes de I/O (`os`, `io/ioutil`, `encoding/json` diretamente), apenas interfaces

### 6.3 Configuração
- **Carregamento:** Centralizado no módulo `config/`
- **Injeção:** Config carregada no `main.go` e injetada nos comandos via construtor
- **Padrão:** Config opcional (valores padrão se não existir)
- **Localização:** XDG-compliant (`~/.config/specs/config.json`)

## 7. Escalabilidade e Manutenibilidade

### 7.1 Extensibilidade
- **Novos comandos:** 
  - Adicionar arquivo em `internal/commands/` ou subdiretório apropriado
  - Registrar comando no router em `internal/cli/router.go`
  - Seguir interface de comando existente
- **Novos serviços:** 
  - Adicionar em `internal/services/`
  - Definir interfaces se necessário para testabilidade
- **Plugins (futuro):** 
  - Sistema de plugins via interface comum
  - Plugins carregados dinamicamente (futuro)

### 7.2 Manutenibilidade
- **Código modular:** Cada módulo com responsabilidade única e clara
- **Documentação:** 
  - Comentários em funções exportadas
  - README.md atualizado com cada feature
  - Exemplos de uso em help dos comandos
- **Testes:** 
  - Cobertura mínima de 80%
  - Testes de unidade para lógica de negócio
  - Testes de integração para comandos completos
- **Refatoração:** 
  - Arquivos não devem exceder 300 linhas
  - Funções não devem exceder 50 linhas
  - Extrair lógica complexa para funções auxiliares

## 8. Referências

### 8.1 Contexto Global
- **Referência:** Visão, objetivos, escopo e requisitos não funcionais estão em `00-global-context.spec.md`.

### 8.2 Stack Técnica
- **Referência:** Linguagem, ferramentas e estrutura de build estão em `00-stack.spec.md`.

## Critérios de Aceite (Arquitetura)

- [ ] Padrão arquitetural definido e justificado
- [ ] Estrutura de diretórios acordada e documentada
- [ ] Módulos e componentes identificados com responsabilidades claras
- [ ] Isolamento e testabilidade garantidos (interfaces, adapters, mocks)
- [ ] Padrões de design aplicados e documentados
- [ ] Fluxo de dados descrito
- [ ] Convenções arquiteturais estabelecidas
- [ ] Estratégia de escalabilidade e manutenibilidade definida

## Checklist Rápido (preencha antes de gerar código)

- [ ] Padrão arquitetural está claro e justificado?
- [ ] Estrutura de diretórios está definida e alinhada com o padrão?
- [ ] Isolamento e testabilidade estão garantidos?
- [ ] Padrões de design estão documentados?
- [ ] Convenções arquiteturais estão escritas?
- [ ] Fluxo de dados está descrito?
