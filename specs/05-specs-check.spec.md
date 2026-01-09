# 05 - Verificação de Consistência Estrutural

Esta especificação define o comando `specs check` para verificar consistência estrutural de specs, incluindo numeração sequencial, links, referências cruzadas e estrutura de diretórios.

## 1. Contexto e Objetivo

- **Contexto:** Desenvolvedores precisam verificar se as specs estão consistentes estruturalmente, com numeração correta, links válidos, referências cruzadas corretas e sem specs órfãs ou duplicadas.
- **Objetivo:** 
  - Verificar numeração sequencial de specs (00-*, 01-*, 02-*, etc.)
  - Validar links e referências entre specs
  - Detectar specs órfãs (referenciadas mas não existem)
  - Detectar specs duplicadas (mesma numeração)
  - Verificar formato de nomes de arquivos
  - Validar estrutura de diretórios
  - Identificar referências quebradas
- **Escopo:** 
  - Validação de numeração sequencial
  - Validação de links internos (entre specs)
  - Validação de referências cruzadas
  - Detecção de specs órfãs e duplicadas
  - Validação de formato de nomes de arquivos
  - Validação de estrutura de diretórios
  - Fora de escopo: validação de conteúdo semântico, validação de links externos (URLs), validação de checklist (coberto por `specs validate`), validação de formato de spec (coberto por `specs validate`)

## 2. Requisitos Funcionais

- **RF01 - Validação de Numeração Sequencial:**
  - Verificar que specs estão numeradas sequencialmente (00-*, 01-*, 02-*, etc.)
  - Detectar gaps na numeração (ex.: 00, 01, 03 - falta 02)
  - Detectar numeração duplicada (ex.: dois arquivos 01-*.spec.md)
  - Validar formato de numeração (dois dígitos, zero-padded)
  - Reportar numerações inconsistentes

- **RF02 - Validação de Links e Referências:**
  - Extrair todos os links Markdown de cada spec (formato `[texto](caminho)`)
  - Identificar links internos (referências a outras specs, ex.: `00-architecture.spec.md`)
  - Validar que links internos apontam para arquivos existentes
  - Detectar links quebrados (arquivo referenciado não existe)
  - Validar formato de links (caminhos relativos corretos)
  - Reportar links inválidos com localização (arquivo e linha)

- **RF03 - Detecção de Specs Órfãs:**
  - Identificar specs que são referenciadas mas não existem
  - Detectar referências a specs que foram removidas ou renomeadas
  - Reportar specs órfãs com lista de arquivos que as referenciam

- **RF04 - Detecção de Specs Duplicadas:**
  - Identificar múltiplos arquivos com mesma numeração (ex.: `01-test.spec.md` e `01-other.spec.md`)
  - Reportar duplicatas com caminhos completos dos arquivos

- **RF05 - Validação de Formato de Nomes:**
  - Validar que nomes de arquivos seguem padrão `{numero}-{nome-descritivo}.spec.md`
  - Validar que números são zero-padded (00, 01, 02, não 0, 1, 2)
  - Validar que nomes são descritivos (não apenas número)
  - Reportar nomes com formato inválido

- **RF06 - Validação de Estrutura de Diretórios:**
  - Verificar que specs estão no diretório correto (`specs/` por padrão)
  - Validar estrutura de diretórios (não há specs em locais inesperados)
  - Detectar specs fora do diretório padrão (se configurável)

- **RF07 - Validação de Referências Cruzadas:**
  - Mapear todas as referências entre specs (quem referencia quem)
  - Detectar referências circulares (se aplicável)
  - Validar que referências seguem convenções (ex.: specs base 00-* são referenciadas corretamente)
  - Reportar referências inconsistentes

- **RF08 - Geração de Relatório:**
  - Exibir resumo de verificação (total de specs, problemas encontrados)
  - Listar problemas por categoria (numeração, links, órfãs, etc.)
  - Formato de saída legível por padrão (texto)
  - Flag `--json` (futuro) para output estruturado em JSON

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs check [caminho]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--json` (futuro): Output em formato JSON estruturado
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[caminho]` (opcional): Caminho para diretório contendo specs. Se omitido, usa `./specs`
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - todas as verificações passaram (sem problemas encontrados)
  - `1`: Erro - problemas de consistência encontrados (numeração, links, etc.)
  - `2`: Erro - input inválido (caminho não existe, não é diretório)
- **Output:**
  - **Sucesso (stdout):** Mensagem de sucesso e resumo
  - **Problemas (stdout):** Relatório detalhado de problemas encontrados
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: "erro: diretório não existe: specs/")
- **Exemplos de uso:**
  ```bash
  $ specs check
  Verificando consistência estrutural em ./specs...
  
  ✅ Numeração: OK (00, 01, 02, 03, 05)
  ⚠️  Numeração: Gap detectado - falta 04
  ❌ Links: 2 links quebrados encontrados
    - specs/02-init.spec.md:10: link para '03-specs-validate.spec.md' não encontrado
    - specs/03-specs-validate.spec.md:5: link para '00-architecture.spec.md' quebrado
  
  Resumo:
    Total de specs: 5
    Problemas encontrados: 3
    - Numeração: 1 gap
    - Links: 2 quebrados
  
  $ specs check specs/
  Verificando consistência estrutural em specs/...
  
  ✅ Numeração: OK (00, 01, 02, 03)
  ✅ Links: Todos os links válidos
  ✅ Estrutura: OK
  
  Todas as verificações passaram!
  
  $ specs check --help
  Verifica consistência estrutural de specs (numeração, links, referências).
  
  Uso:
    specs check [caminho] [flags]
  
  Flags:
    --help    Exibe ajuda para este comando
  
  Exemplos:
    specs check                    # Verifica specs/ no diretório atual
    specs check specs/             # Verifica diretório specs/
  ```

### Arquivos

- **Arquivos de spec:**
  - Formato de nome: `{numero}-{nome-descritivo}.spec.md`
  - Numeração: Zero-padded, dois dígitos (00, 01, 02, etc.)
  - Localização: Diretório `specs/` por padrão

- **Links internos:**
  - Formato: Markdown `[texto](caminho)`
  - Caminhos: Relativos ao diretório de specs
  - Referências: Podem apontar para outras specs ou seções dentro de specs

## 4. Fluxos e Estados

### Fluxo Feliz - Verificação Sem Problemas

1. Usuário executa `specs check` ou `specs check specs/`
2. Sistema identifica diretório `specs/` (padrão ou especificado)
3. Sistema lista todos os arquivos `.spec.md` no diretório
4. Sistema valida numeração sequencial (sem gaps, sem duplicatas)
5. Sistema extrai e valida todos os links de cada spec
6. Sistema verifica que links apontam para arquivos existentes
7. Sistema valida formato de nomes de arquivos
8. Sistema valida estrutura de diretórios
9. Sistema exibe relatório: "Todas as verificações passaram!"
10. Comando retorna código 0

### Fluxo - Verificação com Problemas

1. Usuário executa `specs check`
2. Sistema identifica diretório e lista arquivos
3. Sistema detecta problemas:
   - Gap na numeração (falta 04)
   - Links quebrados (2 arquivos referenciados não existem)
   - Nome com formato inválido (1 arquivo)
4. Sistema agrega todos os problemas por categoria
5. Sistema exibe relatório detalhado com:
   - Resumo por categoria
   - Lista de problemas com localização (arquivo:linha)
   - Resumo geral
6. Comando retorna código 1 (problemas encontrados)

### Estados Alternativos

- **Erro: Caminho não existe:**
  - Mensagem: "erro: caminho não existe: {caminho}"
  - Código de saída: 2
  - Ação: Verificar caminho fornecido

- **Erro: Caminho não é diretório:**
  - Mensagem: "erro: caminho não é diretório: {caminho}"
  - Código de saída: 2
  - Ação: Verificar tipo do caminho

- **Erro: Falha ao ler diretório:**
  - Mensagem: "erro: falha ao ler diretório {caminho}: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar permissões

- **Problema: Gap na numeração:**
  - Mensagem: "⚠️  Numeração: Gap detectado - falta {numero}"
  - Código de saída: 1 (se houver problemas)
  - Ação: Adicionar spec faltante ou renumerar

- **Problema: Numeração duplicada:**
  - Mensagem: "❌ Numeração: Duplicata detectada - {numero} usado em {arquivo1} e {arquivo2}"
  - Código de saída: 1
  - Ação: Renumerar ou remover duplicata

- **Problema: Link quebrado:**
  - Mensagem: "❌ Links: {arquivo}:{linha}: link para '{caminho}' não encontrado"
  - Código de saída: 1
  - Ação: Corrigir link ou criar arquivo referenciado

- **Problema: Nome com formato inválido:**
  - Mensagem: "❌ Formato: {arquivo} não segue padrão {numero}-{nome}.spec.md"
  - Código de saída: 1
  - Ação: Renomear arquivo para seguir padrão

## 5. Dados

### Mapeamento de Specs

- Sistema constrói mapeamento de todas as specs (número → arquivo)
- Mapeamento é usado para validar links e referências
- Mapeamento é temporário (não persistido)

### Índice de Links

- Sistema constrói índice de todos os links encontrados
- Índice mapeia arquivo referenciado → arquivos que referenciam
- Índice é usado para detectar specs órfãs

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Verificação de diretório com 10 specs: < 200ms
  - Verificação de diretório com 100 specs: < 2s
  - Extração de links: < 50ms por arquivo
  - Validação de links: < 100ms por arquivo

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Suporta caminhos relativos e absolutos
  - Suporta diretórios com muitos arquivos

- **Segurança:**
  - Valida caminhos para evitar directory traversal
  - Não acessa arquivos fora do diretório especificado
  - Valida tamanho de arquivo antes de ler (limite razoável)

- **Observabilidade:**
  - Mensagens claras sobre problemas encontrados
  - Localização precisa de problemas (arquivo:linha)
  - Resumo categorizado de problemas

## 7. Guardrails

- **Restrições:**
  - Verifica apenas arquivos com extensão `.spec.md`
  - Não modifica arquivos (apenas leitura)
  - Não valida conteúdo semântico (apenas estrutura)
  - Não segue links externos (apenas valida existência de arquivos)

- **Convenções:**
  - Numeração: Zero-padded, dois dígitos (00, 01, 02)
  - Formato de nome: `{numero}-{nome-descritivo}.spec.md`
  - Links: Caminhos relativos ao diretório de specs

- **Padrão de mensagens:**
  - Sucesso: "✅ {categoria}: OK"
  - Warning: "⚠️  {categoria}: {problema}"
  - Erro: "❌ {categoria}: {problema}"
  - Resumo: "Total de specs: {n}, Problemas encontrados: {n}"

## 8. Critérios de Aceite

- [x] Comando `specs check` verifica consistência em `specs/` quando executado sem argumentos
- [x] Comando `specs check [diretório]` verifica diretório específico
- [x] Comando detecta gaps na numeração sequencial e reporta
- [x] Comando detecta numeração duplicada e reporta com arquivos envolvidos
- [x] Comando valida todos os links internos e detecta links quebrados
- [x] Comando detecta specs órfãs (referenciadas mas não existem)
- [x] Comando valida formato de nomes de arquivos (padrão correto)
- [x] Comando valida estrutura de diretórios
- [x] Comando exibe relatório categorizado de problemas encontrados
- [x] Comando retorna código 0 quando todas as verificações passam
- [x] Comando retorna código 1 quando há problemas de consistência
- [x] Comando retorna código 2 para erros de input inválido
- [x] Comando `specs check --help` exibe ajuda do comando
- [x] Comando reporta localização precisa de problemas (arquivo:linha quando aplicável)
- [x] Comando processa arquivos eficientemente (performance adequada)

## 9. Testes

### Testes de Unidade

- Validação de numeração sequencial (detecção de gaps)
- Detecção de numeração duplicada
- Extração de links Markdown de arquivos
- Validação de links (verificação de existência)
- Validação de formato de nomes de arquivos
- Detecção de specs órfãs (referenciadas mas não existem)

### Testes de Integração

- Fluxo completo de verificação sem problemas
- Fluxo completo de verificação com gaps na numeração
- Fluxo completo de verificação com links quebrados
- Fluxo completo de verificação com specs órfãs
- Fluxo completo de verificação com numeração duplicada
- Agregação de múltiplos problemas em relatório

### Testes E2E

- Execução de `specs check` em projeto com specs consistentes
- Execução de `specs check` em projeto com gaps na numeração
- Execução de `specs check` em projeto com links quebrados
- Execução de `specs check` em projeto com specs órfãs
- Execução de `specs check` em projeto com numeração duplicada
- Execução de `specs check [diretório]` com diretório específico
- Execução de `specs check --help` exibe ajuda corretamente
- Validação de performance com muitos arquivos (100+ specs)

### Como Rodar

- Testes unitários: `go test ./internal/services/checker/...`
- Testes de integração: `go test -tags=integration ./internal/services/checker/...`
- Testes E2E: Executar manualmente em projeto de teste com specs variadas

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (comando apenas lê e verifica)
- Comando funciona com specs existentes sem modificá-las

### Rollback

- Não há rollback necessário (comando não modifica arquivos)
- Verificação pode ser executada múltiplas vezes sem efeitos colaterais

## 11. Observações Operacionais

### Uso em CI/CD

- Comando pode ser usado em pipelines CI/CD para verificar consistência antes de merge
- Código de saída permite integração com scripts de CI/CD
- Relatório pode ser usado para bloquear merges se houver problemas

### Relação com Outros Comandos

- `specs check` foca em consistência estrutural (numeração, links, referências)
- `specs validate` foca em completude (seções, checklist)
- Comandos são complementares e podem ser usados juntos

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Validação de conteúdo semântico (ex.: se especificação faz sentido)
- Validação de links externos (verificar se URLs estão acessíveis)
- Validação de formato de spec (coberto por `specs validate`)
- Validação de checklist (coberto por `specs validate`)
- Correção automática de problemas (apenas detecta, não corrige)
- Output em JSON (flag `--json` fica para v2)

### Decisões em Aberto

- Estratégia de correção automática (se houver no futuro)
- Formato de output JSON (estrutura do relatório)

## Checklist Rápido (preencha antes de gerar código)

- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
