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

### ✅ Especificadas (Aguardando Implementação)

- [x] **02-init.spec.md** - Inicialização de projetos SDD
  - **Status:** Especificada, aguardando implementação
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
    - Flags: `--force`, `--with-boilerplate`

- [x] **03-specs-validate.spec.md** - Validação de specs
  - **Status:** Especificada, aguardando implementação
  - **Prioridade:** Alta
  - **Dependências:** Nenhuma
  - **Estimativa:** Alta
  - **Descrição:** Comando `specs validate [caminho]` para validar specs contra checklist
  - **Funcionalidades principais:**
    - Validação contra checklist formal
    - Verificação de seções obrigatórias
    - Validação de formato e estrutura
    - Relatório de erros e warnings
    - Suporte a arquivo único ou diretório completo
    - Identificação de specs completas/incompletas

- [x] **05-specs-check.spec.md** - Verificação de consistência estrutural
  - **Status:** Especificada, aguardando implementação
  - **Prioridade:** Média
  - **Dependências:** Nenhuma
  - **Estimativa:** Média
  - **Descrição:** Comando `specs check [caminho]` para verificar consistência estrutural
  - **Funcionalidades principais:**
    - Validação de numeração sequencial
    - Verificação de links e referências
    - Validação de formato de nomes de arquivos
    - Verificação de referências cruzadas entre specs
    - Detecção de specs órfãs
    - Validação de estrutura de diretórios

### 🔄 Em Andamento

_Nenhuma no momento_

### 📋 Planejadas (v1)

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

- [x] **06-specs-view.spec.md** - Dashboard de visualização
  - **Status:** Especificada, aguardando implementação
  - **Prioridade:** Média
  - **Dependências:** 03-specs-validate, 04-specs-list (reutiliza lógica)
  - **Estimativa:** Média
  - **Descrição:** Comando `specs view` para exibir dashboard interativo com informações agregadas
  - **Funcionalidades principais:**
    - Dashboard com seções organizadas (Summary, Specs em Progresso, Specs Completas, Specifications)
    - Estatísticas agregadas (total de specs, requirements, progresso)
    - Barras de progresso visuais para specs incompletas
    - Contagem de requirements por spec
    - Formatação visual e legível

### 🔄 Em Andamento

_Nenhuma no momento_

### 📋 Planejadas (v1)

#### Prioridade Baixa / Opcional

- [ ] **07-config.spec.md** - Sistema de configuração
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

1. ✅ **01-version-control** - COMPLETO (implementado)
2. 📝 **02-init** - ESPECIFICADA (aguardando implementação)
3. 📝 **03-specs-validate** - ESPECIFICADA (aguardando implementação)
4. 📝 **04-specs-list** - ESPECIFICADA (aguardando implementação)
5. 📝 **05-specs-check** - ESPECIFICADA (aguardando implementação)
6. 📝 **06-specs-view** - ESPECIFICADA (aguardando implementação)
7. 🔜 **07-config** - A ESPECIFICAR (pode ser feito incrementalmente)

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

- **Total de specs v1:** 7
- **Implementadas:** 1 (14.3%)
- **Especificadas (aguardando implementação):** 5 (71.4%)
- **A especificar:** 1 (14.3%)
- **Prioridade alta:** 2 (ambas especificadas)
- **Prioridade média:** 3 (todas especificadas)
- **Prioridade baixa:** 1 (a especificar)

## Atualizações

- **2024-01-07:** Backlog criado
- **2024-01-07:** Spec 01-version-control implementada e completa
- **2024-01-07:** Specs 02-init, 03-specs-validate, 04-specs-list e 05-specs-check especificadas
- **2024-01-07:** Spec 06-specs-view especificada (dashboard de visualização)