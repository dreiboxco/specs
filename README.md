# Specs - CLI para Gerenciamento de Projetos SDD

Este repositório contém as especificações (specs) para um CLI multiplataforma desenvolvido em Go que facilita a criação, validação e gerenciamento de aplicações desenvolvidas com **SDD (Spec Driven Development)**.

## 📋 Visão Geral

Este projeto segue a metodologia **Spec Driven Development (SDD)**, onde as especificações são escritas antes da implementação. Este repositório contém:

- **Especificações do projeto CLI** (`specs/`): Definições completas de arquitetura, stack técnica, contexto global e funcionalidades
- **Boilerplate genérico** (`boilerplate/`): Templates reutilizáveis para criar novos projetos seguindo SDD

## 🎯 Objetivo

O CLI `specs` (nome do comando) será uma ferramenta de linha de comando que permite:

- ✅ Inicializar novos projetos seguindo estrutura SDD
- ✅ Validar especificações contra checklist e padrões
- ✅ Listar e gerenciar especificações do projeto
- ✅ Gerar artefatos de software a partir de specs validadas
- ✅ Integrar com pipelines CI/CD para validação automática

## 📁 Estrutura do Projeto

```
specs/
├── specs/                    # Especificações do projeto CLI
│   ├── 00-global-context.spec.md    # Contexto, visão, objetivos
│   ├── 00-architecture.spec.md      # Arquitetura e padrões
│   ├── 00-stack.spec.md             # Stack técnica (Go 1.25.5)
│   └── template-default.spec.md     # Template para novas specs
├── boilerplate/              # Boilerplate genérico
│   ├── specs/                # Templates de specs genéricos
│   └── .cursorrules          # Regras do Cursor para SDD
└── README.md                 # Este arquivo
```

## 🚀 Como Usar

### Para Desenvolvedores do CLI

Este repositório contém as especificações que devem ser seguidas para implementar o CLI. Antes de implementar qualquer funcionalidade:

1. **Consulte as specs** em `specs/` para entender o que precisa ser implementado
2. **Valide a spec** contra o checklist antes de começar a codificar
3. **Implemente conforme a spec** - não adicione funcionalidades não especificadas

### Para Criar Novos Projetos SDD

O diretório `boilerplate/` contém templates genéricos que podem ser usados como base para novos projetos:

1. **Copie o boilerplate** para seu novo projeto
2. **Preencha os arquivos `00-*`** com as informações do seu projeto:
   - `00-global-context.spec.md`: Contexto, visão, objetivos
   - `00-architecture.spec.md`: Arquitetura e padrões
   - `00-stack.spec.md`: Stack técnica escolhida
3. **Use `template-default.spec.md`** como base para criar specs de funcionalidades
4. **Configure o `.cursorrules`** conforme necessário (veja seção abaixo)

## ⚙️ Configuração do Cursor Rules

O arquivo `.cursorrules` contém regras e diretrizes para o Cursor AI seguir durante o desenvolvimento. Este arquivo é uma **base configurável** que pode e deve ser adaptado para cada projeto específico.

### Sobre o `.cursorrules`

O `.cursorrules` fornece:
- **Fluxo SDD**: Regras para consultar specs antes de implementar
- **Guardrails**: Restrições de dependências, convenções de código
- **Padrões**: UX de CLI, segurança, testes, commits
- **Gerenciamento**: Branches, checklists, validação de specs

### Personalização

O `.cursorrules` é uma **base extensível** e pode ser:

- **Estendido**: Adicione novas regras específicas do seu projeto
- **Modificado**: Ajuste regras existentes para se adequar às necessidades
- **Adaptado**: Remova ou altere seções que não se aplicam ao seu contexto

#### Exemplos de Ajustes Comuns

- **Stack específica**: Se usar outra linguagem, ajuste exemplos e referências
- **Convenções de projeto**: Adicione regras específicas de nomenclatura ou estrutura
- **Ferramentas customizadas**: Inclua regras para ferramentas específicas do seu projeto
- **Workflow de equipe**: Adapte regras de branches e commits para o workflow da equipe
- **Integrações**: Adicione regras para integrações específicas (APIs, serviços, etc.)

#### Localização

- **Boilerplate**: `boilerplate/.cursorrules` - Versão genérica para novos projetos
- **Projeto atual**: `.cursorrules` na raiz - Versão específica deste projeto

### Recomendações

1. **Comece com a base**: Use o `.cursorrules` do boilerplate como ponto de partida
2. **Itere conforme necessário**: Ajuste as regras à medida que o projeto evolui
3. **Documente mudanças**: Se fizer ajustes significativos, documente o motivo
4. **Mantenha genérico quando possível**: Evite regras muito específicas que limitem reutilização

O objetivo é ter um conjunto de regras que **facilite o desenvolvimento seguindo SDD**, mas que seja **flexível o suficiente** para se adaptar a diferentes contextos e necessidades de projeto.

## 📚 Especificações

### Especificações Base (00-*)

- **`00-global-context.spec.md`**: Define o contexto global do projeto, visão, objetivos, escopo, requisitos não funcionais, estratégias de distribuição e configuração
- **`00-architecture.spec.md`**: Define o padrão arquitetural, estrutura de diretórios, isolamento de módulos e decisões de design
- **`00-stack.spec.md`**: Define a stack tecnológica, ferramentas e plataformas de build/distribuição

### Stack Técnica Definida

- **Linguagem**: Go 1.25.5
- **Plataformas**: macOS (x64, arm64) e Linux (x64, arm64)
- **Build**: Binário estático único por plataforma
- **Distribuição**: GitHub Releases com checksum SHA256

## 🔧 Desenvolvimento

### Pré-requisitos

- Go 1.25.5 ou superior
- Git

### Estrutura Esperada do Projeto CLI

Quando implementado, o projeto CLI seguirá esta estrutura:

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
└── testdata/             # Arquivos de teste
```

### Comandos Planejados (v1)

- `specs init` - Inicializar projeto SDD
- `specs validate [caminho]` - Validar specs contra checklist
- `specs list` - Listar todas as specs com status
- `specs check [caminho]` - Verificar consistência estrutural
- `specs version` - Exibir versão do CLI

## 📖 Metodologia SDD

**Spec Driven Development (SDD)** é uma metodologia onde:

1. **Especificar antes de codificar**: Todas as funcionalidades são especificadas em arquivos `*.spec.md` antes da implementação
2. **Validar contra checklist**: Cada spec deve passar por um checklist formal antes de ser implementada
3. **Implementar conforme spec**: O código implementa exatamente o que está especificado, sem adicionar funcionalidades não especificadas
4. **Manter specs atualizadas**: As specs são a fonte da verdade e devem ser mantidas atualizadas

## 🔍 Validação de Specs

As specs devem seguir o formato definido e passar pelo checklist. Uma spec é considerada completa quando:

1. Possui todas as seções obrigatórias
2. Todos os itens do checklist estão marcados como concluídos
3. Os critérios de aceite são testáveis e mensuráveis

## 📝 Convenções

- **Especificações**: Arquivos `*.spec.md` em `specs/`
- **Numeração**: `00-*` para specs base, `01-*`, `02-*`, etc. para funcionalidades
- **Formato**: Markdown com seções padronizadas
- **Checklist**: Sempre no final da spec, após "Abertos / Fora de Escopo"

## 🤝 Contribuindo

Este é um projeto de especificações. Para contribuir:

1. Revise as specs existentes
2. Proponha mudanças ou novas specs seguindo o template
3. Certifique-se de que a spec passa pelo checklist antes de solicitar implementação

## 📄 Licença

[Definir licença quando aplicável]

## 🔗 Referências

- Metodologia SDD (Spec Driven Development)
- [Go Documentation](https://go.dev/doc/)
- Arquitetura e padrões definidos em `specs/00-architecture.spec.md`
- Stack técnica detalhada em `specs/00-stack.spec.md`

