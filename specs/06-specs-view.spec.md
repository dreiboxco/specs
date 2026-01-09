# 06 - Dashboard de Visualização

Esta especificação define o comando `specs view` para exibir um dashboard interativo com informações agregadas sobre o projeto SDD, incluindo resumo de specs, progresso de completude, e visão geral do estado do projeto.

## 1. Contexto e Objetivo

- **Contexto:** Desenvolvedores precisam de uma visão agregada e visual do estado do projeto SDD, incluindo progresso de specs, estatísticas gerais e status de completude, para ter uma compreensão rápida do estado geral do projeto.
- **Objetivo:** 
  - Exibir dashboard interativo com informações agregadas do projeto
  - Mostrar resumo estatístico (total de specs, requirements, progresso)
  - Visualizar progresso de specs incompletas com barras de progresso
  - Listar specs completas e incompletas de forma organizada
  - Exibir contagem de requirements por spec
  - Fornecer visão geral do estado do projeto em um único comando
- **Escopo:** 
  - Dashboard com seções organizadas (Summary, Specs em Progresso, Specs Completas, Lista de Specs)
  - Cálculo de estatísticas agregadas (total de specs, requirements, progresso)
  - Barras de progresso visuais para specs incompletas
  - Contagem de requirements por spec (extraída de seções "Requisitos Funcionais")
  - Formatação visual e legível
  - Exclusão automática de specs de template (00-*.spec.md e template-default.spec.md) do dashboard
  - Fora de escopo: atualização em tempo real, modo interativo com navegação, exportação para outros formatos (JSON futuro), gráficos avançados, histórico de mudanças

## 2. Requisitos Funcionais

- **RF01 - Seção Summary:**
  - Exibir total de specifications encontradas (excluindo specs de template)
  - Exibir total de requirements (contados de todas as specs, excluindo templates)
  - Exibir número de specs em progresso (incompletas, excluindo templates)
  - Exibir número de specs completas (excluindo templates)
  - Calcular e exibir progresso geral (percentual de itens do checklist marcados, excluindo templates)
  - Formatar números de forma destacada (negrito ou cor)
  - Excluir automaticamente specs de template (00-*.spec.md e template-default.spec.md) das estatísticas

- **RF02 - Seção Specs em Progresso:**
  - Listar todas as specs incompletas (< 6 itens do checklist marcados)
  - Exibir nome da spec
  - Calcular percentual de completude (itens marcados / 6)
  - Exibir barra de progresso visual proporcional ao percentual
  - Ordenar por percentual de completude (menor para maior, ou maior para menor)
  - Exibir percentual numérico ao lado da barra

- **RF03 - Seção Specs Completas:**
  - Listar todas as specs completas (6/6 itens do checklist marcados)
  - Exibir nome da spec com ícone de checkmark (✅)
  - Ordenar por numeração (00, 01, 02, etc.)
  - Limitar exibição se houver muitas specs (opção de mostrar todas)

- **RF04 - Seção Specifications:**
  - Listar todas as specs do projeto
  - Exibir nome da spec
  - Contar e exibir número de requirements por spec
  - Ordenar por numeração (00, 01, 02, etc.)
  - Indicar status visual (completa/incompleta) com ícones

- **RF05 - Cálculo de Requirements:**
  - Extrair seção "Requisitos Funcionais" de cada spec
  - Contar itens de requisitos (RF01, RF02, RF03, etc.)
  - Agregar total de requirements de todas as specs (excluindo templates)
  - Tratar casos onde spec não tem seção de requisitos (contar como 0)
  - Excluir specs de template (00-*.spec.md e template-default.spec.md) da contagem

- **RF06 - Cálculo de Progresso:**
  - Para cada spec (excluindo templates), contar itens do checklist marcados (0-6)
  - Calcular progresso individual (itens marcados / 6)
  - Calcular progresso geral (soma de todos os itens marcados / soma de todos os itens possíveis, excluindo templates)
  - Exibir progresso em formato "X/Y (Z% complete)"
  - Excluir specs de template (00-*.spec.md e template-default.spec.md) do cálculo de progresso

- **RF07 - Exclusão de Specs de Template:**
  - Identificar e excluir automaticamente specs de template do dashboard
  - Excluir specs com numeração 00-* (00-global-context.spec.md, 00-architecture.spec.md, 00-stack.spec.md)
  - Excluir template-default.spec.md
  - Não exibir specs de template em nenhuma seção do dashboard (Summary, Specs em Progresso, Specs Completas, Specifications)
  - Não incluir specs de template em cálculos de estatísticas (total, requirements, progresso)
  - Justificativa: Specs de template são base do projeto e não representam funcionalidades a serem implementadas

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs view [caminho]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--json` (futuro): Output em formato JSON estruturado
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[caminho]` (opcional): Caminho para diretório contendo specs. Se omitido, usa `./specs`
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - dashboard exibido corretamente
  - `1`: Erro - falha ao ler diretório ou specs
  - `2`: Erro - input inválido (caminho não existe, não é diretório)
- **Output:**
  - **Sucesso (stdout):** Dashboard formatado com todas as seções
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: "erro: diretório não existe: specs/")
- **Exemplos de uso:**
  ```bash
  $ specs view
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
    cli-list                   7 requirements
    cli-view                   7 requirements
    cli-init                   5 requirements
    cli-update                 5 requirements
    cli-change                 4 requirements
    cli-spec                   4 requirements
    cli-show                   3 requirements
  
  $ specs view specs/
  [mesmo formato acima]
  
  $ specs view --help
  Exibe dashboard interativo com informações agregadas do projeto SDD.
  
  Uso:
    specs view [caminho] [flags]
  
  Flags:
    --help    Exibe ajuda para este comando
  
  Exemplos:
    specs view                    # Dashboard de specs/ no diretório atual
    specs view specs/             # Dashboard de diretório específico
  ```

### Formato de Saída

- **Cabeçalho:** "Specs Dashboard" ou título similar
- **Seções:** Separadas por linhas em branco, com títulos em negrito ou destacados
- **Barras de progresso:** Caracteres visuais (ex.: `[████████  ]`) representando percentual
- **Ícones:** ✅ para completo, ⚠️ para incompleto (quando aplicável)
- **Números destacados:** Valores importantes em negrito ou formato destacado

## 4. Fluxos e Estados

### Fluxo Feliz - Exibição de Dashboard

1. Usuário executa `specs view` ou `specs view specs/`
2. Sistema identifica diretório `specs/` (padrão ou especificado)
3. Sistema lista todos os arquivos `.spec.md` no diretório (recursivo)
4. Para cada spec, sistema:
   - Lê arquivo e extrai seção "Requisitos Funcionais"
   - Conta número de requirements (RF01, RF02, etc.)
   - Verifica checklist e conta itens marcados (0-6)
   - Calcula percentual de completude
5. Sistema agrega estatísticas:
   - Total de specs
   - Total de requirements (soma de todos)
   - Total de specs completas
   - Total de specs em progresso
   - Progresso geral (soma de itens marcados / soma de itens possíveis)
6. Sistema formata e exibe seção Summary
7. Sistema formata e exibe seção Specs em Progresso (com barras de progresso)
8. Sistema formata e exibe seção Specs Completas
9. Sistema formata e exibe seção Specifications (com contagem de requirements)
10. Comando retorna código 0

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
  - Dashboard exibe zeros em todas as estatísticas
  - Código de saída: 0 (não é erro, apenas informação)

- **Aviso: Spec sem seção de requisitos:**
  - Spec é contada como 0 requirements
  - Não interrompe processamento
  - Dashboard exibe normalmente

- **Aviso: Spec com checklist inválido:**
  - Spec é tratada como incompleta (0% de progresso)
  - Não interrompe processamento
  - Dashboard exibe normalmente

## 5. Dados

### Estatísticas Agregadas

- **Total de specs:** Contagem de arquivos `.spec.md` encontrados
- **Total de requirements:** Soma de todos os requirements de todas as specs
- **Specs completas:** Contagem de specs com 6/6 itens do checklist marcados
- **Specs em progresso:** Contagem de specs com < 6 itens do checklist marcados
- **Progresso geral:** (soma de itens marcados) / (total de specs * 6) * 100

### Dados por Spec

- **Nome:** Extraído do nome do arquivo (sem extensão e numeração)
- **Numeração:** Extraída do nome do arquivo (00, 01, 02, etc.)
- **Requirements:** Contados da seção "Requisitos Funcionais"
- **Progresso:** Itens do checklist marcados / 6
- **Status:** Completa (6/6) ou Incompleta (< 6/6)

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Geração de dashboard com 10 specs: < 500ms
  - Geração de dashboard com 100 specs: < 5s
  - Cálculo de requirements: < 50ms por spec
  - Cálculo de progresso: < 50ms por spec

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Suporta caminhos relativos e absolutos
  - Suporta diretórios com muitos arquivos
  - Funciona em terminais com diferentes larguras (ajusta formatação)

- **Segurança:**
  - Valida caminhos para evitar directory traversal
  - Não acessa arquivos fora do diretório especificado
  - Valida tamanho de arquivo antes de ler (limite razoável)

- **Observabilidade:**
  - Formatação visual clara e legível
  - Barras de progresso proporcionais e precisas
  - Números e estatísticas destacados

## 7. Guardrails

- **Restrições:**
  - Processa apenas arquivos com extensão `.spec.md`
  - Reutiliza lógica de validação (não duplica código)
  - Não modifica arquivos (apenas leitura)
  - Barras de progresso têm largura fixa (ex.: 10 caracteres)
  - Exclui automaticamente specs de template (00-*.spec.md e template-default.spec.md) de todas as estatísticas e exibições

- **Convenções:**
  - Ordenação de specs: Por numeração (00, 01, 02, etc.)
  - Ordenação de specs em progresso: Por percentual (configurável)
  - Formato de barra: `[████████  ]` (10 caracteres, preenchido proporcionalmente)
  - Formato de percentual: "X%" (sem casas decimais na v1)

- **Padrão de mensagens:**
  - Cabeçalho: "Specs Dashboard"
  - Seções: Títulos em negrito ou destacados
  - Números: Destacados (negrito ou cor quando suportado)

## 8. Critérios de Aceite

- [x] Comando `specs view` exibe dashboard quando executado sem argumentos
- [x] Comando `specs view [diretório]` exibe dashboard de diretório específico
- [x] Dashboard exibe seção Summary com todas as estatísticas (specs, requirements, progresso)
- [x] Dashboard exibe seção Specs em Progresso com barras de progresso visuais
- [x] Dashboard exibe seção Specs Completas com ícones de checkmark
- [x] Dashboard exibe seção Specifications com contagem de requirements por spec
- [x] Barras de progresso são proporcionais ao percentual de completude
- [x] Percentuais são calculados corretamente (itens marcados / 6)
- [x] Total de requirements é calculado corretamente (soma de todos)
- [x] Progresso geral é calculado corretamente (soma de itens / total possível)
- [x] Comando retorna código 0 em caso de sucesso
- [x] Comando retorna código 1 para erros de I/O
- [x] Comando retorna código 2 para erros de input inválido
- [x] Comando `specs view --help` exibe ajuda do comando
- [x] Dashboard funciona com specs sem seção de requisitos (conta como 0)
- [x] Dashboard funciona com specs com checklist inválido (trata como incompleta)
- [x] Formatação é legível em terminais de diferentes larguras
- [x] Dashboard exclui automaticamente specs de template (00-*.spec.md e template-default.spec.md) de todas as seções
- [x] Estatísticas não incluem specs de template (total, requirements, progresso)

## 9. Testes

### Testes de Unidade

- Contagem de requirements em seção "Requisitos Funcionais"
- Cálculo de percentual de progresso (itens marcados / 6)
- Cálculo de progresso geral (agregação)
- Formatação de barras de progresso (proporção correta)
- Extração de nome e numeração de arquivo de spec

### Testes de Integração

- Fluxo completo de geração de dashboard sem specs
- Fluxo completo de geração de dashboard com specs variadas
- Cálculo correto de estatísticas agregadas
- Formatação correta de todas as seções
- Tratamento de specs sem seção de requisitos
- Tratamento de specs com checklist inválido

### Testes E2E

- Execução de `specs view` em projeto com specs completas e incompletas
- Execução de `specs view [diretório]` com diretório específico
- Execução de `specs view --help` exibe ajuda corretamente
- Verificação de formatação visual (barras, números, seções)
- Verificação de cálculos (requirements, progresso)
- Validação de performance com muitos arquivos (100+ specs)
- Verificação de formatação em terminais de diferentes larguras

### Como Rodar

- Testes unitários: `go test ./internal/services/viewer/...`
- Testes de integração: `go test -tags=integration ./internal/services/viewer/...`
- Testes E2E: Executar manualmente em projeto de teste com specs variadas

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (comando apenas lê e exibe)
- Comando funciona com specs existentes sem modificá-las

### Rollback

- Não há rollback necessário (comando não modifica arquivos)
- Dashboard pode ser executado múltiplas vezes sem efeitos colaterais

## 11. Observações Operacionais

### Relação com Outros Comandos

- `specs view` fornece visão agregada e visual (dashboard)
- `specs list` fornece listagem tabular detalhada
- `specs validate` fornece detalhes de validação (erros, warnings)
- `specs check` verifica consistência estrutural
- Comandos são complementares e podem ser usados juntos

### Uso em CI/CD

- Comando pode ser usado em pipelines CI/CD para gerar relatórios visuais
- Código de saída permite integração com scripts de CI/CD
- Dashboard pode ser capturado e incluído em relatórios

### Performance

- Cálculo de requirements e progresso pode ser otimizado com cache (futuro)
- Processamento de muitos arquivos pode ser paralelizado
- Formatação de barras de progresso é leve (cálculo simples)

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Modo interativo com navegação (apenas exibição estática)
- Atualização em tempo real (apenas snapshot atual)
- Gráficos avançados (apenas barras de progresso simples)
- Exportação para outros formatos (JSON, HTML - fica para v2)
- Histórico de mudanças de progresso
- Comparação entre versões de specs
- Filtros interativos (apenas exibição completa)

### Decisões em Aberto

- Formato de output JSON (estrutura do dashboard)
- Largura de barras de progresso (fixa ou adaptável)
- Ordenação de specs em progresso (por percentual crescente ou decrescente)
- Limite de exibição de specs completas (mostrar todas ou limitar)

## Checklist Rápido (preencha antes de gerar código)

- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
