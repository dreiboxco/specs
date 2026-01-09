# 03 - Validação de Specs

Esta especificação define o comando `specs validate` para validar especificações contra o checklist formal e verificar se atendem aos critérios de completude definidos na metodologia SDD.

## 1. Contexto e Objetivo

- **Contexto:** Desenvolvedores precisam validar se suas specs estão completas e prontas para implementação, verificando se todas as seções obrigatórias estão presentes e se o checklist foi preenchido corretamente.
- **Objetivo:** 
  - Validar specs contra checklist formal definido em `checklist.md`
  - Verificar presença de seções obrigatórias em cada spec
  - Verificar formato e estrutura de arquivos de spec
  - Gerar relatório de validação com erros e warnings
  - Suportar validação de arquivo único ou diretório completo
  - Identificar specs completas e incompletas
- **Escopo:** 
  - Validação contra checklist formal (seções obrigatórias, formato)
  - Validação de estrutura de arquivos Markdown
  - Validação de formato de checklist (marcação de itens)
  - Geração de relatório de validação
  - Suporte a arquivo único ou diretório recursivo
  - Fora de escopo: validação semântica de conteúdo (ex.: se especificação faz sentido), validação de links externos, validação de código de exemplo, validação de integridade de referências cruzadas (coberto por `specs check`)

## 2. Requisitos Funcionais

- **RF01 - Validação de Seções Obrigatórias:**
  - Verificar presença de todas as seções obrigatórias do template (1-12)
  - Seções obrigatórias: Contexto e Objetivo, Requisitos Funcionais, Contratos e Interfaces, Fluxos e Estados, Dados, NFRs, Guardrails, Critérios de Aceite, Testes, Migração/Rollback, Observações Operacionais, Abertos/Fora de Escopo
  - Detectar seções faltantes e reportar em relatório
  - Validar que seções estão no formato correto (títulos com `##`)

- **RF02 - Validação de Checklist:**
  - Verificar presença de checklist no final da spec (após seção "Abertos / Fora de Escopo")
  - Validar formato do checklist (itens com `- [ ]` ou `- [x]`)
  - Contar número de itens do checklist (deve ter exatamente 6 itens)
  - Identificar se todos os itens estão marcados (spec completa) ou não (spec incompleta)
  - Validar que checklist está no formato correto (markdown de lista)

- **RF03 - Validação de Formato de Arquivo:**
  - Verificar que arquivo tem extensão `.spec.md`
  - Verificar que arquivo é Markdown válido (estrutura básica)
  - Verificar que arquivo não está vazio
  - Verificar encoding UTF-8
  - Detectar caracteres inválidos ou problemas de encoding

- **RF04 - Validação de Estrutura:**
  - Verificar que arquivo começa com título principal (`#`)
  - Verificar numeração no título (ex.: `# 02 - Nome da Spec`)
  - Validar hierarquia de títulos (não pular níveis, ex.: `##` após `#`)
  - Verificar que seções têm conteúdo (não apenas título)

- **RF05 - Validação de Múltiplos Arquivos:**
  - Suportar validação de arquivo único (caminho para arquivo `.spec.md`)
  - Suportar validação de diretório (valida todos os `.spec.md` recursivamente)
  - Se caminho não especificado, validar diretório `specs/` no diretório atual
  - Processar arquivos em paralelo quando possível (para performance)
  - Agregar resultados de múltiplos arquivos em relatório único

- **RF06 - Geração de Relatório:**
  - Exibir resumo de validação (total de specs, completas, incompletas, com erros)
  - Listar erros encontrados por spec (seções faltantes, formato inválido)
  - Listar warnings (checklist incompleto, mas spec válida)
  - Formato de saída legível por padrão (texto)
  - Flag `--json` (futuro) para output estruturado em JSON

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs validate [caminho]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--json` (futuro): Output em formato JSON estruturado
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[caminho]` (opcional): Caminho para arquivo `.spec.md` ou diretório contendo specs. Se omitido, usa `./specs`
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - todas as specs validadas são válidas (podem estar completas ou incompletas, mas sem erros)
  - `1`: Erro - falha ao ler arquivos ou specs com erros de formato
  - `2`: Erro - input inválido (caminho não existe, não é arquivo nem diretório)
- **Output:**
  - **Sucesso (stdout):** Relatório de validação com resumo e detalhes
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: "erro: arquivo não encontrado: specs/01-test.spec.md")
- **Exemplos de uso:**
  ```bash
  $ specs validate
  Validando specs em ./specs...
  
  ✅ specs/00-global-context.spec.md: Completa (6/6 itens do checklist)
  ✅ specs/00-architecture.spec.md: Completa (6/6 itens do checklist)
  ⚠️  specs/01-version-control.spec.md: Incompleta (4/6 itens do checklist)
  ❌ specs/02-init.spec.md: Erro - seção "Testes" faltando
  
  Resumo:
    Total: 4 specs
    Completas: 2
    Incompletas: 1
    Com erros: 1
  
  $ specs validate specs/01-version-control.spec.md
  ✅ specs/01-version-control.spec.md: Completa (6/6 itens do checklist)
  
  $ specs validate ./minhas-specs
  Validando specs em ./minhas-specs...
  
  ✅ minhas-specs/01-test.spec.md: Completa (6/6 itens do checklist)
  
  $ specs validate --help
  Valida specs contra checklist formal e verifica estrutura.
  
  Uso:
    specs validate [caminho] [flags]
  
  Flags:
    --help    Exibe ajuda para este comando
  
  Exemplos:
    specs validate                    # Valida specs/ no diretório atual
    specs validate specs/             # Valida diretório specs/
    specs validate specs/01-test.spec.md  # Valida arquivo específico
  ```

### Arquivos

- **Arquivos de spec:**
  - Formato: Markdown (`.spec.md`)
  - Encoding: UTF-8
  - Estrutura: Deve seguir template definido em `template-default.spec.md`
  - Checklist: Deve estar no final, após seção "Abertos / Fora de Escopo"

- **Checklist:**
  - Formato: Lista Markdown com itens `- [ ]` (pendente) ou `- [x]` (concluído)
  - Número de itens: Exatamente 6 itens
  - Localização: Final da spec, após seção "Abertos / Fora de Escopo"

## 4. Fluxos e Estados

### Fluxo Feliz - Validação de Diretório

1. Usuário executa `specs validate` ou `specs validate specs/`
2. Sistema identifica diretório `specs/` (padrão ou especificado)
3. Sistema lista todos os arquivos `.spec.md` no diretório (recursivo)
4. Para cada arquivo, sistema valida:
   - Presença de seções obrigatórias
   - Formato e estrutura
   - Checklist (presença, formato, completude)
5. Sistema agrega resultados de todas as specs
6. Sistema exibe relatório com resumo e detalhes
7. Comando retorna código 0 (todas válidas, mesmo que incompletas)

### Fluxo - Validação de Arquivo Único

1. Usuário executa `specs validate specs/01-test.spec.md`
2. Sistema verifica que caminho é arquivo `.spec.md`
3. Sistema valida arquivo único:
   - Seções obrigatórias
   - Formato e estrutura
   - Checklist
4. Sistema exibe resultado da validação
5. Comando retorna código 0 (válida) ou 1 (com erros)

### Estados Alternativos

- **Erro: Caminho não existe:**
  - Mensagem: "erro: caminho não existe: {caminho}"
  - Código de saída: 2
  - Ação: Verificar caminho fornecido

- **Erro: Caminho não é arquivo nem diretório:**
  - Mensagem: "erro: caminho inválido: {caminho} (não é arquivo nem diretório)"
  - Código de saída: 2
  - Ação: Verificar tipo do caminho

- **Erro: Arquivo não é .spec.md:**
  - Mensagem: "erro: arquivo deve ter extensão .spec.md: {caminho}"
  - Código de saída: 2
  - Ação: Verificar extensão do arquivo

- **Erro: Falha ao ler arquivo:**
  - Mensagem: "erro: falha ao ler arquivo {caminho}: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar permissões e integridade do arquivo

- **Erro: Seção obrigatória faltando:**
  - Mensagem: "❌ {arquivo}: Erro - seção '{seção}' faltando"
  - Código de saída: 1 (se houver erros)
  - Ação: Adicionar seção faltante na spec

- **Warning: Checklist incompleto:**
  - Mensagem: "⚠️  {arquivo}: Incompleta ({marcados}/6 itens do checklist)"
  - Código de saída: 0 (não é erro, apenas incompleta)
  - Ação: Completar checklist na spec

- **Erro: Checklist com formato inválido:**
  - Mensagem: "❌ {arquivo}: Erro - checklist com formato inválido (esperado 6 itens, encontrado {n})"
  - Código de saída: 1
  - Ação: Corrigir formato do checklist

- **Erro: Encoding inválido:**
  - Mensagem: "erro: arquivo {caminho} não está em UTF-8"
  - Código de saída: 1
  - Ação: Converter arquivo para UTF-8

## 5. Dados

### Checklist de Referência

- Checklist formal está definido em `checklist.md` no diretório `specs/`
- Checklist define seções obrigatórias e critérios de validação
- Checklist é lido uma vez por execução e usado para validar todas as specs

### Resultados de Validação

- Resultados são temporários (não persistidos)
- Relatório é gerado em tempo real durante validação
- Nenhum cache de validação na v1 (sempre valida do zero)

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Validação de arquivo único: < 100ms
  - Validação de diretório com 10 specs: < 500ms
  - Validação de diretório com 100 specs: < 5s
  - Processamento paralelo quando possível (múltiplos arquivos)

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Suporta caminhos relativos e absolutos
  - Suporta diretórios com muitos arquivos

- **Segurança:**
  - Valida caminhos para evitar directory traversal
  - Não lê arquivos fora do diretório especificado
  - Valida tamanho de arquivo antes de ler (limite razoável)

- **Observabilidade:**
  - Mensagens claras sobre o que está sendo validado
  - Progresso visível durante validação de múltiplos arquivos
  - Logs de erros descritivos em stderr

## 7. Guardrails

- **Restrições:**
  - Valida apenas arquivos com extensão `.spec.md`
  - Não valida conteúdo semântico (apenas estrutura e formato)
  - Não segue links externos ou valida referências
  - Não modifica arquivos (apenas leitura)

- **Convenções:**
  - Seções obrigatórias definidas no template
  - Checklist deve ter exatamente 6 itens
  - Formato de checklist: `- [ ]` ou `- [x]`

- **Padrão de mensagens:**
  - Sucesso: "✅ {arquivo}: Completa ({n}/6 itens do checklist)"
  - Incompleta: "⚠️  {arquivo}: Incompleta ({n}/6 itens do checklist)"
  - Erro: "❌ {arquivo}: Erro - {descrição}"
  - Resumo: "Total: {n} specs, Completas: {n}, Incompletas: {n}, Com erros: {n}"

## 8. Critérios de Aceite

- [ ] Comando `specs validate` valida todas as specs em `specs/` quando executado sem argumentos
- [ ] Comando `specs validate [arquivo]` valida arquivo único especificado
- [ ] Comando `specs validate [diretório]` valida todas as specs no diretório recursivamente
- [ ] Comando detecta seções obrigatórias faltando e reporta como erro
- [ ] Comando valida formato do checklist (6 itens, formato correto)
- [ ] Comando identifica specs completas (6/6 itens marcados) e incompletas (< 6 itens)
- [ ] Comando exibe relatório com resumo (total, completas, incompletas, com erros)
- [ ] Comando retorna código 0 quando todas as specs são válidas (mesmo que incompletas)
- [ ] Comando retorna código 1 quando há specs com erros de formato
- [ ] Comando retorna código 2 para erros de input inválido
- [ ] Comando `specs validate --help` exibe ajuda do comando
- [ ] Comando valida apenas arquivos com extensão `.spec.md`
- [ ] Comando processa arquivos em paralelo quando possível (performance)
- [ ] Comando valida encoding UTF-8 e reporta erros de encoding

## 9. Testes

### Testes de Unidade

- Validação de seções obrigatórias (detecção de seções faltantes)
- Validação de formato de checklist (contagem de itens, formato)
- Validação de estrutura de Markdown (títulos, hierarquia)
- Parsing de checklist (detecção de itens marcados/não marcados)
- Validação de encoding UTF-8

### Testes de Integração

- Fluxo completo de validação de arquivo único válido
- Fluxo completo de validação de arquivo único com erros
- Fluxo completo de validação de diretório com múltiplos arquivos
- Validação de specs completas e incompletas
- Agregação de resultados de múltiplos arquivos
- Tratamento de erros (arquivo não encontrado, permissões)

### Testes E2E

- Execução de `specs validate` em projeto com specs válidas
- Execução de `specs validate` em projeto com specs incompletas
- Execução de `specs validate` em projeto com specs com erros
- Execução de `specs validate [arquivo]` com arquivo específico
- Execução de `specs validate [diretório]` com diretório específico
- Execução de `specs validate --help` exibe ajuda corretamente
- Validação de performance com muitos arquivos (100+ specs)

### Como Rodar

- Testes unitários: `go test ./internal/services/validator/...`
- Testes de integração: `go test -tags=integration ./internal/services/validator/...`
- Testes E2E: Executar manualmente em projeto de teste com specs variadas

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (comando apenas lê e valida)
- Comando funciona com specs existentes sem modificá-las

### Rollback

- Não há rollback necessário (comando não modifica arquivos)
- Validação pode ser executada múltiplas vezes sem efeitos colaterais

## 11. Observações Operacionais

### Uso em CI/CD

- Comando pode ser usado em pipelines CI/CD para validar specs antes de merge
- Código de saída permite integração com scripts de CI/CD
- Relatório pode ser usado para bloquear merges se houver erros

### Performance

- Validação de muitos arquivos pode ser lenta; processamento paralelo ajuda
- Cache de validação pode ser adicionado no futuro (v2)

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Validação semântica de conteúdo (ex.: se especificação faz sentido)
- Validação de links externos (verificar se URLs estão acessíveis)
- Validação de código de exemplo (syntax highlighting, execução)
- Validação de integridade de referências cruzadas (coberto por `specs check`)
- Cache de resultados de validação
- Output em JSON (flag `--json` fica para v2)

### Decisões em Aberto

- Estratégia de cache de validação (se houver no futuro)
- Formato de output JSON (estrutura do relatório)

## Checklist Rápido (preencha antes de gerar código)

- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
