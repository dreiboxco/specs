# 04 - Listagem de Specs

Esta especificação define o comando `specs list` para listar todas as especificações do projeto com seu status (completa/incompleta), fornecendo visibilidade sobre o estado das specs e facilitando o gerenciamento do projeto SDD.

## 1. Contexto e Objetivo

- **Contexto:** Desenvolvedores precisam ter visibilidade sobre todas as specs do projeto, saber quais estão completas (prontas para implementação) e quais estão incompletas (aguardando especificação), para gerenciar o progresso do projeto SDD.
- **Objetivo:** 
  - Listar todas as specs do projeto com informações de status
  - Identificar specs completas (checklist totalmente preenchido) e incompletas
  - Exibir informações resumidas de cada spec (nome, numeração, status)
  - Fornecer contadores agregados (total, completas, incompletas)
  - Permitir filtros para visualizar apenas specs completas ou incompletas
  - Formato de saída legível (tabela ou lista)
- **Escopo:** 
  - Listagem de todas as specs em diretório especificado
  - Verificação de status usando lógica de validação (reutiliza `specs validate`)
  - Formatação de saída em tabela ou lista
  - Filtros opcionais (completas, incompletas)
  - Contadores agregados
  - Fora de escopo: ordenação customizada, agrupamento por categoria, exportação para outros formatos (JSON futuro), informações detalhadas de cada spec (usar `specs validate` para detalhes)

## 2. Requisitos Funcionais

- **RF01 - Listagem de Specs:**
  - Listar todos os arquivos `.spec.md` no diretório especificado (recursivo)
  - Ordenar specs por numeração (00-*, 01-*, 02-*, etc.)
  - Exibir nome do arquivo e numeração de cada spec
  - Exibir status de cada spec (completa/incompleta)
  - Processar diretório recursivamente se necessário

- **RF02 - Verificação de Status:**
  - Para cada spec, verificar se checklist está completo (6/6 itens marcados)
  - Usar mesma lógica de validação de `specs validate` (reutilizar serviço)
  - Identificar specs completas (todos os itens do checklist marcados)
  - Identificar specs incompletas (alguns itens do checklist não marcados)
  - Identificar specs com erros (formato inválido, seções faltando)

- **RF03 - Formatação de Saída:**
  - Formato padrão: Tabela legível com colunas (Numeração, Nome, Status)
  - Alternativa: Lista simples (uma linha por spec)
  - Exibir ícones ou símbolos para status (✅ completo, ⚠️ incompleto, ❌ erro)
  - Formatar saída de forma legível e alinhada

- **RF04 - Contadores Agregados:**
  - Contar total de specs encontradas
  - Contar specs completas (6/6 itens do checklist)
  - Contar specs incompletas (< 6 itens do checklist)
  - Contar specs com erros (formato inválido)
  - Exibir resumo no final da listagem

- **RF05 - Filtros:**
  - Flag `--complete` ou `--only-complete`: Listar apenas specs completas
  - Flag `--incomplete` ou `--only-incomplete`: Listar apenas specs incompletas
  - Flag `--errors`: Listar apenas specs com erros
  - Sem flags: Listar todas as specs (completas, incompletas e com erros)

- **RF06 - Informações Adicionais:**
  - Exibir caminho relativo de cada spec
  - Exibir numeração da spec (extraída do nome do arquivo)
  - Exibir nome descritivo da spec (extraído do nome do arquivo)
  - Exibir status detalhado quando aplicável (ex.: "4/6 itens do checklist")

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs list [caminho]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--complete`, `--only-complete`: Lista apenas specs completas
  - `--incomplete`, `--only-incomplete`: Lista apenas specs incompletas
  - `--errors`: Lista apenas specs com erros
  - `--json` (futuro): Output em formato JSON estruturado
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[caminho]` (opcional): Caminho para diretório contendo specs. Se omitido, usa `./specs`
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - listagem exibida corretamente
  - `1`: Erro - falha ao ler diretório ou specs
  - `2`: Erro - input inválido (caminho não existe, não é diretório)
- **Output:**
  - **Sucesso (stdout):** Tabela ou lista de specs com status
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: "erro: diretório não existe: specs/")
- **Exemplos de uso:**
  ```bash
  $ specs list
  Listando specs em ./specs...
  
  Numeração  Nome                    Status
  ──────────  ──────────────────────  ──────────
  00          global-context          ✅ Completa
  00          architecture            ✅ Completa
  00          stack                     ✅ Completa
  01          version-control         ✅ Completa
  02          init                    ⚠️  Incompleta (4/6)
  03          specs-validate          ⚠️  Incompleta (3/6)
  05          specs-check             ⚠️  Incompleta (2/6)
  
  Resumo:
    Total: 7 specs
    Completas: 4
    Incompletas: 3
    Com erros: 0
  
  $ specs list --complete
  Listando specs completas em ./specs...
  
  Numeração  Nome                    Status
  ──────────  ──────────────────────  ──────────
  00          global-context          ✅ Completa
  00          architecture            ✅ Completa
  00          stack                   ✅ Completa
  01          version-control         ✅ Completa
  
  Total: 4 specs completas
  
  $ specs list --incomplete
  Listando specs incompletas em ./specs...
  
  Numeração  Nome                    Status
  ──────────  ──────────────────────  ──────────
  02          init                    ⚠️  Incompleta (4/6)
  03          specs-validate          ⚠️  Incompleta (3/6)
  05          specs-check             ⚠️  Incompleta (2/6)
  
  Total: 3 specs incompletas
  
  $ specs list specs/
  Listando specs em specs/...
  
  [mesmo formato acima]
  
  $ specs list --help
  Lista todas as specs do projeto com status (completa/incompleta).
  
  Uso:
    specs list [caminho] [flags]
  
  Flags:
    --complete, --only-complete     Lista apenas specs completas
    --incomplete, --only-incomplete  Lista apenas specs incompletas
    --errors                         Lista apenas specs com erros
    --help                           Exibe ajuda para este comando
  
  Exemplos:
    specs list                       # Lista todas as specs em specs/
    specs list --complete            # Lista apenas specs completas
    specs list --incomplete          # Lista apenas specs incompletas
    specs list specs/                # Lista specs em diretório específico
  ```

### Arquivos

- **Diretório de specs:**
  - Localização padrão: `specs/` no diretório atual
  - Pode ser especificado via argumento
  - Processamento recursivo (busca em subdiretórios)

- **Arquivos de spec:**
  - Formato: `.spec.md`
  - Numeração: Extraída do nome do arquivo (ex.: `02-init.spec.md` → `02`)
  - Nome: Extraído do nome do arquivo (ex.: `02-init.spec.md` → `init`)

## 4. Fluxos e Estados

### Fluxo Feliz - Listagem de Todas as Specs

1. Usuário executa `specs list` ou `specs list specs/`
2. Sistema identifica diretório `specs/` (padrão ou especificado)
3. Sistema lista todos os arquivos `.spec.md` no diretório (recursivo)
4. Sistema ordena specs por numeração (00, 01, 02, etc.)
5. Para cada spec, sistema verifica status usando lógica de validação:
   - Lê arquivo e verifica checklist
   - Conta itens marcados no checklist
   - Determina status (completa/incompleta/erro)
6. Sistema formata saída em tabela com colunas (Numeração, Nome, Status)
7. Sistema calcula contadores agregados (total, completas, incompletas, erros)
8. Sistema exibe tabela e resumo
9. Comando retorna código 0

### Fluxo - Listagem com Filtro

1. Usuário executa `specs list --complete`
2. Sistema identifica diretório e lista arquivos
3. Sistema verifica status de cada spec
4. Sistema filtra apenas specs completas (6/6 itens marcados)
5. Sistema exibe tabela apenas com specs completas
6. Sistema exibe resumo (total de completas)
7. Comando retorna código 0

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

- **Aviso: Nenhuma spec encontrada:**
  - Mensagem: "Nenhuma spec encontrada em {caminho}"
  - Código de saída: 0 (não é erro, apenas informação)
  - Ação: Verificar se diretório contém arquivos `.spec.md`

- **Aviso: Filtro não retornou resultados:**
  - Mensagem: "Nenhuma spec {tipo} encontrada" (ex.: "Nenhuma spec completa encontrada")
  - Código de saída: 0 (não é erro)
  - Ação: Verificar filtro aplicado

- **Erro: Falha ao validar spec:**
  - Mensagem: "erro: falha ao validar spec {arquivo}: {detalhes}"
  - Spec é marcada como "erro" na listagem
  - Código de saída: 0 (continua processando outras specs)

## 5. Dados

### Cache de Validação (Futuro)

- Na v1, validação é feita em tempo real (sem cache)
- Futuro: Cache de resultados de validação para melhor performance
- Cache seria invalidado quando specs são modificadas

### Resultados de Listagem

- Resultados são temporários (não persistidos)
- Listagem é gerada em tempo real
- Nenhum estado persistente na v1

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Listagem de diretório com 10 specs: < 500ms
  - Listagem de diretório com 100 specs: < 5s
  - Validação de status: Reutiliza lógica de `specs validate` (mesma performance)
  - Processamento pode ser otimizado com cache (futuro)

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Suporta caminhos relativos e absolutos
  - Suporta diretórios com muitos arquivos

- **Segurança:**
  - Valida caminhos para evitar directory traversal
  - Não acessa arquivos fora do diretório especificado
  - Valida tamanho de arquivo antes de ler (limite razoável)

- **Observabilidade:**
  - Mensagens claras sobre o que está sendo listado
  - Formato de saída legível e organizado
  - Resumo agregado fornece visibilidade geral

## 7. Guardrails

- **Restrições:**
  - Lista apenas arquivos com extensão `.spec.md`
  - Reutiliza lógica de validação (não duplica código)
  - Não modifica arquivos (apenas leitura)
  - Não ordena por critérios customizados (apenas numeração)

- **Convenções:**
  - Ordenação: Por numeração (00, 01, 02, etc.)
  - Formato de saída: Tabela com colunas alinhadas
  - Status: ✅ Completa, ⚠️ Incompleta, ❌ Erro

- **Padrão de mensagens:**
  - Cabeçalho: "Listando specs em {caminho}..."
  - Tabela: Colunas (Numeração, Nome, Status)
  - Resumo: "Total: {n} specs, Completas: {n}, Incompletas: {n}, Com erros: {n}"

## 8. Critérios de Aceite

- [ ] Comando `specs list` lista todas as specs em `specs/` quando executado sem argumentos
- [ ] Comando `specs list [diretório]` lista specs em diretório específico
- [ ] Comando exibe specs ordenadas por numeração (00, 01, 02, etc.)
- [ ] Comando identifica status de cada spec (completa/incompleta/erro)
- [ ] Comando exibe tabela formatada com colunas (Numeração, Nome, Status)
- [ ] Comando exibe contadores agregados (total, completas, incompletas, erros)
- [ ] Comando `specs list --complete` lista apenas specs completas
- [ ] Comando `specs list --incomplete` lista apenas specs incompletas
- [ ] Comando `specs list --errors` lista apenas specs com erros
- [ ] Comando retorna código 0 em caso de sucesso
- [ ] Comando retorna código 1 para erros de I/O
- [ ] Comando retorna código 2 para erros de input inválido
- [ ] Comando `specs list --help` exibe ajuda do comando
- [ ] Comando reutiliza lógica de validação (não duplica código)
- [ ] Comando processa diretório recursivamente quando necessário
- [ ] Comando exibe mensagem apropriada quando nenhuma spec é encontrada

## 9. Testes

### Testes de Unidade

- Extração de numeração e nome de arquivo de spec
- Ordenação de specs por numeração
- Formatação de tabela (alinhamento, colunas)
- Cálculo de contadores agregados
- Aplicação de filtros (completas, incompletas, erros)

### Testes de Integração

- Fluxo completo de listagem sem filtros
- Fluxo completo de listagem com filtro `--complete`
- Fluxo completo de listagem com filtro `--incomplete`
- Fluxo completo de listagem com filtro `--errors`
- Integração com serviço de validação (reutilização de lógica)
- Ordenação correta de specs por numeração
- Cálculo correto de contadores agregados

### Testes E2E

- Execução de `specs list` em projeto com specs variadas
- Execução de `specs list --complete` com specs completas e incompletas
- Execução de `specs list --incomplete` com specs completas e incompletas
- Execução de `specs list --errors` com specs com erros
- Execução de `specs list [diretório]` com diretório específico
- Execução de `specs list --help` exibe ajuda corretamente
- Verificação de formato de tabela (alinhamento, legibilidade)
- Verificação de contadores agregados (precisão)
- Validação de performance com muitos arquivos (100+ specs)

### Como Rodar

- Testes unitários: `go test ./internal/services/lister/...`
- Testes de integração: `go test -tags=integration ./internal/services/lister/...`
- Testes E2E: Executar manualmente em projeto de teste com specs variadas

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (comando apenas lê e lista)
- Comando funciona com specs existentes sem modificá-las

### Rollback

- Não há rollback necessário (comando não modifica arquivos)
- Listagem pode ser executada múltiplas vezes sem efeitos colaterais

## 11. Observações Operacionais

### Relação com Outros Comandos

- `specs list` fornece visão geral (status de todas as specs)
- `specs validate` fornece detalhes de validação (erros, warnings)
- `specs check` verifica consistência estrutural (numeração, links)
- Comandos são complementares e podem ser usados juntos

### Uso em CI/CD

- Comando pode ser usado em pipelines CI/CD para verificar status das specs
- Código de saída permite integração com scripts de CI/CD
- Resumo pode ser usado para gerar relatórios de progresso

### Performance

- Validação de status pode ser lenta com muitos arquivos
- Reutilização de lógica de validação evita duplicação
- Cache de validação pode ser adicionado no futuro (v2) para melhor performance

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Ordenação customizada (apenas por numeração)
- Agrupamento por categoria ou tipo
- Informações detalhadas de cada spec (usar `specs validate` para detalhes)
- Exportação para outros formatos (CSV, JSON - JSON fica para v2)
- Histórico de mudanças de status
- Notificações quando specs são completadas

### Decisões em Aberto

- Formato de output JSON (estrutura do relatório)
- Estratégia de cache de validação (se houver no futuro)
- Ordenação alternativa (por nome, por data de modificação)

## Checklist Rápido (preencha antes de gerar código)

- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
