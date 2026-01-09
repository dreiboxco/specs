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

- [x] **02-init.spec.md** - Inicialização de projetos SDD
  - Comando `specs init` para criar estrutura base de projeto SDD
  - Criar estrutura de diretórios (`specs/`, `boilerplate/`)
  - Copiar templates de specs base (`00-*.spec.md`)
  - Criar arquivo `.cursorrules` base
  - Criar `README.md` inicial
  - Flags: `--force`, `--with-boilerplate`

- [x] **03-specs-validate.spec.md** - Validação de specs
  - Comando `specs validate [caminho]` para validar specs contra checklist
  - Validação contra checklist formal
  - Verificação de seções obrigatórias
  - Validação de formato e estrutura
  - Relatório de erros e warnings
  - Suporte a arquivo único ou diretório completo

- [x] **04-specs-list.spec.md** - Listagem de specs
  - Comando `specs list` para listar todas as specs com status
  - Status de cada spec (completa/incompleta)
  - Formato de saída legível (tabela)
  - Filtros opcionais (--complete, --incomplete, --errors)
  - Contadores agregados

- [x] **05-specs-check.spec.md** - Verificação de consistência estrutural
  - Comando `specs check [caminho]` para verificar consistência estrutural
  - Validação de numeração sequencial
  - Verificação de links e referências
  - Validação de formato de nomes de arquivos
  - Detecção de specs órfãs e duplicadas

- [x] **06-specs-view.spec.md** - Dashboard de visualização
  - Comando `specs view` para exibir dashboard interativo
  - Dashboard com seções organizadas (Summary, Specs em Progresso, Specs Completas, Specifications)
  - Estatísticas agregadas (total de specs, requirements, progresso)
  - Barras de progresso visuais
  - Exclusão automática de specs de template

### 📋 Especificadas (Aguardando Implementação)

_Nenhuma no momento_

### 📋 Planejadas (v1)

_Nenhuma - todas as specs v1 foram especificadas_

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
2. ✅ **02-init** - COMPLETO (implementado)
3. ✅ **03-specs-validate** - COMPLETO (implementado)
4. ✅ **04-specs-list** - COMPLETO (implementado)
5. ✅ **05-specs-check** - COMPLETO (implementado)
6. ✅ **06-specs-view** - COMPLETO (implementado)
7. ✅ **07-config** - COMPLETO (implementado)
8. ✅ **08-ci-cd-setup** - ESPECIFICADA (workflows criados, aguardando validação)

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

- **Total de specs v1:** 8
- **Implementadas:** 7 (87.5%)
- **Especificadas (aguardando validação):** 1 (12.5%)
- **A especificar:** 0 (0%)
- **Prioridade alta:** 2 (ambas implementadas)
- **Prioridade média:** 3 (todas implementadas)
- **Prioridade baixa:** 2 (ambas especificadas/implementadas)

**Nota:** Todas as specs v1 foram especificadas. Este backlog será removido pois o projeto agora segue apenas as specs em `specs/`.

## Atualizações

- **2024-01-07:** Backlog criado
- **2024-01-07:** Spec 01-version-control implementada e completa
- **2024-01-07:** Specs 02-init, 03-specs-validate, 04-specs-list, 05-specs-check e 06-specs-view especificadas
- **2025-01-08:** Specs 02-init, 03-specs-validate, 04-specs-list, 05-specs-check e 06-specs-view implementadas e completas
- **2025-01-08:** Spec 07-config especificada e implementada
- **2025-01-08:** Spec 08-ci-cd-setup especificada
- **2025-01-08:** Todas as specs v1 especificadas - BACKLOG.md será removido