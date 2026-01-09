# 02 - Inicialização de Projetos SDD

Esta especificação define o comando `specs init` para inicializar novos projetos seguindo a metodologia SDD (Spec Driven Development), criando a estrutura de diretórios e arquivos base necessários.

## 1. Contexto e Objetivo

- **Contexto:** Desenvolvedores precisam de uma forma rápida e padronizada de inicializar novos projetos seguindo a metodologia SDD, com estrutura de diretórios e templates de specs base.
- **Objetivo:** 
  - Criar estrutura de diretórios padrão para projetos SDD
  - Copiar templates de specs base (00-*.spec.md) e checklist
  - Criar arquivo `.cursorrules` base para desenvolvimento
  - Criar `README.md` inicial com estrutura básica
  - Validar se já existe projeto SDD no diretório atual
  - Permitir inicialização idempotente (executar múltiplas vezes produz mesmo resultado)
- **Escopo:** 
  - Criação de estrutura de diretórios (`specs/`, `boilerplate/` se necessário)
  - Cópia de templates de specs base do boilerplate interno
  - Criação de arquivos de configuração base
  - Validação de pré-condições (diretório não vazio, já existe projeto SDD)
  - Fora de escopo: templates customizados, inicialização interativa, integração com Git (criação de repositório), configuração de CI/CD

## 2. Requisitos Funcionais

- **RF01 - Validação de Pré-condições:**
  - Comando deve verificar se diretório atual já contém projeto SDD (presença de `specs/` com arquivos `00-*.spec.md`)
  - Se projeto SDD já existe, exibir mensagem informativa e retornar código de saída 0 (idempotente)
  - Se diretório não está vazio e não é projeto SDD, perguntar confirmação (flag `--force` para sobrescrever)
  - Comando deve funcionar em diretório vazio ou não vazio

- **RF02 - Criação de Estrutura de Diretórios:**
  - Criar diretório `specs/` na raiz do projeto
  - Criar diretório `boilerplate/specs/` (opcional, apenas se flag `--with-boilerplate` for usada)
  - Criar diretórios com permissões apropriadas (755 para diretórios)
  - Não falhar se diretórios já existem (idempotência)

- **RF03 - Cópia de Templates de Specs Base:**
  - Copiar `00-global-context.spec.md` do boilerplate interno para `specs/`
  - Copiar `00-architecture.spec.md` do boilerplate interno para `specs/`
  - Copiar `00-stack.spec.md` do boilerplate interno para `specs/`
  - Copiar `checklist.md` do boilerplate interno para `specs/`
  - Copiar `template-default.spec.md` do boilerplate interno para `specs/`
  - Arquivos copiados devem ter permissões 644 (rw-r--r--)
  - Não sobrescrever arquivos existentes a menos que flag `--force` seja usada

- **RF04 - Criação de Arquivos de Configuração:**
  - Criar arquivo `.cursorrules` na raiz com regras base para SDD
  - Criar arquivo `README.md` na raiz com estrutura básica e instruções
  - Arquivos devem ter permissões 644 (rw-r--r--)
  - Não sobrescrever arquivos existentes a menos que flag `--force` seja usada

- **RF05 - Validação Pós-Inicialização:**
  - Após inicialização, validar que estrutura foi criada corretamente
  - Verificar que arquivos base existem e são legíveis
  - Exibir mensagem de sucesso com resumo do que foi criado

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs init [diretório]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--force`: Sobrescreve arquivos existentes sem confirmação
  - `--with-boilerplate`: Cria também diretório `boilerplate/` com templates genéricos
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[diretório]` (opcional): Diretório onde inicializar projeto. Se omitido, usa diretório atual
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - projeto inicializado ou já existe
  - `1`: Erro - falha ao criar estrutura ou arquivos
  - `2`: Erro - input inválido (diretório não existe, sem permissão de escrita)
- **Output:**
  - **Sucesso (stdout):** Mensagem informativa sobre o que foi criado (ex.: "Projeto SDD inicializado com sucesso em ./specs")
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: "erro: diretório não existe: /caminho/invalido")
- **Exemplos de uso:**
  ```bash
  $ specs init
  Projeto SDD inicializado com sucesso em ./specs
  
  $ specs init ./meu-projeto
  Projeto SDD inicializado com sucesso em ./meu-projeto/specs
  
  $ specs init --force
  Arquivos existentes serão sobrescritos.
  Projeto SDD inicializado com sucesso em ./specs
  
  $ specs init --with-boilerplate
  Projeto SDD inicializado com sucesso em ./specs
  Boilerplate criado em ./boilerplate
  
  $ specs init --help
  Inicializa um novo projeto SDD no diretório especificado.
  
  Uso:
    specs init [diretório] [flags]
  
  Flags:
    --force              Sobrescreve arquivos existentes sem confirmação
    --with-boilerplate   Cria também diretório boilerplate com templates genéricos
    --help               Exibe ajuda para este comando
  ```

### Arquivos

- **Diretórios criados:**
  - `specs/`: Diretório principal de especificações
  - `boilerplate/specs/`: Diretório de templates genéricos (se `--with-boilerplate`)
  - Permissões: 755 (rwxr-xr-x)

- **Arquivos criados:**
  - `specs/00-global-context.spec.md`: Template de contexto global
  - `specs/00-architecture.spec.md`: Template de arquitetura
  - `specs/00-stack.spec.md`: Template de stack técnica
  - `specs/checklist.md`: Checklist de validação
  - `specs/template-default.spec.md`: Template para novas specs
  - `.cursorrules`: Regras do Cursor para SDD
  - `README.md`: Documentação inicial do projeto
  - Permissões: 644 (rw-r--r--)

- **Fonte dos templates:**
  - Templates são copiados do boilerplate interno do CLI (embutido no binário ou em diretório conhecido)
  - Templates devem ser versionados junto com o CLI

## 4. Fluxos e Estados

### Fluxo Feliz - Inicialização em Diretório Vazio

1. Usuário executa `specs init` em diretório vazio ou novo
2. Sistema verifica se diretório atual já contém projeto SDD
3. Sistema verifica permissões de escrita no diretório
4. Sistema cria diretório `specs/`
5. Sistema copia templates de specs base para `specs/`
6. Sistema cria arquivo `.cursorrules` na raiz
7. Sistema cria arquivo `README.md` na raiz
8. Sistema valida que estrutura foi criada corretamente
9. Sistema exibe mensagem de sucesso
10. Comando retorna código 0

### Fluxo - Inicialização com Diretório Específico

1. Usuário executa `specs init ./meu-projeto`
2. Sistema verifica se diretório `./meu-projeto` existe
3. Sistema verifica se diretório já contém projeto SDD
4. Sistema verifica permissões de escrita
5. Sistema cria estrutura dentro de `./meu-projeto/`
6. Sistema exibe mensagem de sucesso com caminho completo
7. Comando retorna código 0

### Fluxo - Projeto SDD Já Existe (Idempotência)

1. Usuário executa `specs init` em diretório que já contém projeto SDD
2. Sistema detecta presença de `specs/00-global-context.spec.md` (ou similar)
3. Sistema exibe mensagem: "Projeto SDD já existe em ./specs. Nada a fazer."
4. Comando retorna código 0 (sucesso, idempotente)

### Estados Alternativos

- **Erro: Diretório não existe:**
  - Mensagem: "erro: diretório não existe: {caminho}"
  - Código de saída: 2
  - Ação: Verificar caminho fornecido

- **Erro: Sem permissão de escrita:**
  - Mensagem: "erro: sem permissão de escrita no diretório: {caminho}"
  - Código de saída: 2
  - Ação: Verificar permissões do diretório

- **Erro: Falha ao criar diretório:**
  - Mensagem: "erro: falha ao criar diretório {caminho}: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar espaço em disco e permissões

- **Erro: Falha ao copiar template:**
  - Mensagem: "erro: falha ao copiar template {nome}: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar integridade do binário e espaço em disco

- **Aviso: Arquivo já existe (sem --force):**
  - Mensagem: "aviso: arquivo {nome} já existe. Use --force para sobrescrever."
  - Comando continua criando outros arquivos
  - Código de saída: 0 (sucesso parcial)

- **Confirmação: Diretório não vazio (sem --force):**
  - Mensagem: "Diretório não está vazio. Continuar? (s/N)"
  - Se usuário confirmar, prossegue
  - Se não confirmar, cancela com código 0

## 5. Dados

### Templates Embarcados

- Templates são parte do binário do CLI ou carregados de diretório conhecido
- Templates não devem ser modificados pelo usuário após inicialização (são apenas base)
- Templates devem ser versionados junto com o CLI

### Estrutura Criada

- **Persistência:** Estrutura criada é permanente até que usuário a remova manualmente
- **Versionamento:** Arquivos criados devem ser versionados no Git pelo usuário (CLI não inicializa Git)
- **Permissões:** Diretórios 755, arquivos 644

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Inicialização completa: < 500ms
  - Criação de diretórios: < 50ms
  - Cópia de templates: < 200ms
  - Validação pós-inicialização: < 100ms

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Respeita permissões do sistema de arquivos
  - Funciona com caminhos relativos e absolutos

- **Segurança:**
  - Não sobrescreve arquivos existentes sem confirmação explícita (`--force`)
  - Valida permissões antes de criar arquivos
  - Não cria arquivos em diretórios sem permissão de escrita

- **Observabilidade:**
  - Mensagens claras sobre o que está sendo criado
  - Logs de erros descritivos em stderr
  - Flag `--debug` (futuro) para logs verbosos

## 7. Guardrails

- **Restrições:**
  - Não inicializa repositório Git (fora de escopo)
  - Não modifica arquivos existentes sem `--force`
  - Não cria arquivos fora do diretório especificado
  - Não requer privilégios de administrador

- **Convenções:**
  - Estrutura de diretórios segue padrão definido em `00-architecture.spec.md`
  - Nomes de arquivos seguem convenções SDD (00-*.spec.md, etc.)
  - Templates são read-only após criação (usuário pode editar, mas CLI não modifica)

- **Padrão de mensagens:**
  - Sucesso: "Projeto SDD inicializado com sucesso em {caminho}"
  - Erro: "erro: {descrição}" em stderr
  - Aviso: "aviso: {descrição}" em stderr

## 8. Critérios de Aceite

- [ ] Comando `specs init` cria estrutura `specs/` com todos os templates base
- [ ] Comando `specs init` cria arquivo `.cursorrules` na raiz
- [ ] Comando `specs init` cria arquivo `README.md` na raiz
- [ ] Comando `specs init` é idempotente (executar duas vezes não falha)
- [ ] Comando `specs init` detecta projeto SDD existente e informa sem falhar
- [ ] Comando `specs init --force` sobrescreve arquivos existentes
- [ ] Comando `specs init --with-boilerplate` cria também diretório `boilerplate/`
- [ ] Comando `specs init [diretório]` inicializa em diretório específico
- [ ] Comando `specs init` retorna código 0 em caso de sucesso
- [ ] Comando `specs init` retorna código 1 para erros de I/O
- [ ] Comando `specs init` retorna código 2 para erros de input inválido
- [ ] Comando `specs init --help` exibe ajuda do comando
- [ ] Todos os arquivos criados têm permissões corretas (644 para arquivos, 755 para diretórios)
- [ ] Templates copiados são idênticos aos templates do boilerplate interno

## 9. Testes

### Testes de Unidade

- Validação de pré-condições (detecção de projeto SDD existente)
- Validação de permissões de escrita
- Formatação de mensagens de sucesso/erro
- Parsing de flags e argumentos

### Testes de Integração

- Fluxo completo de inicialização em diretório temporário
- Inicialização idempotente (executar duas vezes)
- Inicialização com `--force` sobrescrevendo arquivos
- Inicialização com `--with-boilerplate`
- Inicialização em diretório específico
- Validação de estrutura criada (arquivos e diretórios existem)

### Testes E2E

- Execução de `specs init` em ambiente limpo
- Execução de `specs init` em diretório com projeto SDD existente
- Execução de `specs init --force` sobrescrevendo arquivos
- Execução de `specs init [diretório]` com caminho relativo
- Execução de `specs init [diretório]` com caminho absoluto
- Execução de `specs init --help` exibe ajuda corretamente
- Verificação de permissões dos arquivos criados
- Verificação de conteúdo dos templates copiados

### Como Rodar

- Testes unitários: `go test ./internal/services/init/...`
- Testes de integração: `go test -tags=integration ./internal/services/init/...`
- Testes E2E: Executar manualmente ou via script de teste em diretório temporário

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (comando cria estrutura do zero)
- Usuário pode executar `specs init` em qualquer momento para criar estrutura

### Rollback

- Usuário pode remover manualmente diretório `specs/` e arquivos criados
- CLI não fornece comando de "desinicialização" (remoção manual)
- Arquivos criados são independentes e podem ser removidos sem afetar o CLI

## 11. Observações Operacionais

### Distribuição

- Templates devem ser embarcados no binário ou distribuídos junto com o CLI
- Templates devem ser versionados e atualizados junto com releases do CLI

### Versionamento

- Templates seguem versão do CLI
- Usuário pode atualizar templates manualmente editando arquivos
- CLI não atualiza templates automaticamente (fora de escopo v1)

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Templates customizados (usuário deve editar após criação)
- Inicialização interativa (perguntas ao usuário)
- Integração com Git (criação de repositório, commit inicial)
- Configuração de CI/CD
- Atualização automática de templates
- Múltiplos templates de projeto (apenas um template base)

### Decisões em Aberto

- Estratégia de embarcamento de templates (binário vs diretório externo)
- Formato de templates (markdown puro vs templates com variáveis)

## Checklist Rápido (preencha antes de gerar código)

- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
