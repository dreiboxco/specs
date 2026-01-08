# Backlog de Specs - CLI Specs

Este arquivo contém o backlog de especificações a serem implementadas para o CLI `specs`, organizadas por prioridade e status.

## Status das Specs

### ✅ Implementadas

- [x] **01-version-control.spec.md** - Controle de versões
  - Comando `specs version`
  - Arquivo VERSION
  - Incremento automático
  - Tags Git automáticas
  - CI/CD com GitHub Actions

### 🔄 Em Andamento

_Nenhuma no momento_

### 📋 Planejadas (v1)

#### Prioridade Alta

- [ ] **02-init.spec.md** - Inicialização de projetos SDD
  - **Prioridade:** Alta
  - **Dependências:** Nenhuma
  - **Estimativa:** Média
  - **Descrição:** Comando `specs init` para criar estrutura base de projeto SDD
  - **Funcionalidades principais:**
    - Criar estrutura de diretórios (`specs/`, `boilerplate/`)
    - Copiar templates de specs base (`00-*.spec.md`)
    - Criar arquivo `.cursorrules` base
    - Criar `README.md` inicial
    - Validar se já existe projeto SDD
    - Opções de template (básico, completo, customizado)

- [ ] **03-specs-validate.spec.md** - Validação de specs
  - **Prioridade:** Alta
  - **Dependências:** Nenhuma (pode usar boilerplate como referência)
  - **Estimativa:** Alta
  - **Descrição:** Comando `specs validate [caminho]` para validar specs contra checklist
  - **Funcionalidades principais:**
    - Validação contra checklist formal
    - Verificação de seções obrigatórias
    - Validação de formato e estrutura
    - Relatório de erros e warnings
    - Suporte a arquivo único ou diretório completo
    - Output formatado (texto/JSON)

#### Prioridade Média

- [ ] **04-specs-list.spec.md** - Listagem de specs
  - **Prioridade:** Média
  - **Dependências:** 03-specs-validate (usa lógica de validação)
  - **Estimativa:** Média
  - **Descrição:** Comando `specs list` para listar todas as specs com status
  - **Funcionalidades principais:**
    - Listar todas as specs do projeto
    - Status de cada spec (completa/incompleta)
    - Verificação automática de checklist
    - Formato de saída legível (tabela ou lista)
    - Filtros opcionais (apenas completas, apenas incompletas)
    - Contadores (total, completas, incompletas)

- [ ] **05-specs-check.spec.md** - Verificação de consistência estrutural
  - **Prioridade:** Média
  - **Dependências:** Nenhuma
  - **Estimativa:** Média
  - **Descrição:** Comando `specs check [caminho]` para verificar consistência
  - **Funcionalidades principais:**
    - Validação de numeração sequencial
    - Verificação de links e referências
    - Validação de formato de nomes de arquivos
    - Verificação de referências cruzadas entre specs
    - Detecção de specs órfãs
    - Validação de estrutura de diretórios

#### Prioridade Baixa / Opcional

- [ ] **06-config.spec.md** - Sistema de configuração
  - **Prioridade:** Baixa (pode ser integrado em outras specs)
  - **Dependências:** Nenhuma
  - **Estimativa:** Baixa
  - **Descrição:** Sistema de configuração XDG-compliant
  - **Funcionalidades principais:**
    - Configuração em `~/.config/specs/config.json`
    - Caminho padrão para specs configurável
    - Validação de configuração
    - Comando para visualizar/editar configuração
    - Valores padrão sensatos

## Roadmap Futuro (v2+)

### v2 - Funcionalidades Avançadas

- [ ] **07-auto-update.spec.md** - Sistema de auto-atualização
  - Checksum SHA256
  - Validação de integridade
  - Rollback automático

- [ ] **08-generate.spec.md** - Geração avançada de artefatos
  - Geração de código a partir de specs
  - Geração de testes
  - Geração de documentação
  - Templates customizáveis

- [ ] **09-templates.spec.md** - Sistema de templates
  - Templates customizáveis
  - Repositório de templates
  - Validação de templates

### v3 - Integrações e Extensibilidade

- [ ] **10-plugins.spec.md** - Sistema de plugins
  - Arquitetura de plugins
  - API de extensão
  - Gerenciamento de plugins

- [ ] **11-ide-integration.spec.md** - Integração com IDEs
  - Language Server Protocol (LSP)
  - Extensões para IDEs populares
  - Autocomplete e validação em tempo real

- [ ] **12-telemetry.spec.md** - Telemetria opcional
  - Coleta de métricas (opt-in)
  - Análise de uso
  - Melhorias baseadas em dados

## Critérios de Priorização

### Fatores de Prioridade

1. **Dependências:** Specs sem dependências têm prioridade
2. **Valor para o usuário:** Funcionalidades core primeiro
3. **Complexidade:** Implementações mais simples primeiro
4. **Bloqueadores:** Specs que bloqueiam outras têm prioridade

### Ordem Recomendada de Implementação

1. ✅ **01-version-control** - COMPLETO
2. 🔜 **02-init** - Permite criar novos projetos
3. 🔜 **03-specs-validate** - Validação essencial
4. 🔜 **04-specs-list** - Visibilidade do status
5. 🔜 **05-specs-check** - Consistência estrutural
6. 🔜 **06-config** - Pode ser feito incrementalmente

## Notas de Implementação

### Padrões a Seguir

- Todas as specs devem seguir o template `template-default.spec.md`
- Checklist deve ser marcado após implementação completa
- Testes devem ser escritos junto com a implementação
- Documentação deve ser atualizada no README.md

### Convenções

- Numeração sequencial: `02-`, `03-`, `04-`, etc.
- Nomenclatura: `{numero}-{nome-descritivo}.spec.md`
- Branches: `feature/{numero}-{nome-sem-extensao}`
- Commits: `feat: implementa spec {numero}-{nome}`

## Métricas de Progresso

- **Total de specs v1:** 6
- **Implementadas:** 1 (16.7%)
- **Pendentes:** 5 (83.3%)
- **Prioridade alta:** 2
- **Prioridade média:** 2
- **Prioridade baixa:** 1

## Atualizações

- **2024-01-07:** Backlog criado
- **2024-01-07:** Spec 01-version-control implementada e completa
