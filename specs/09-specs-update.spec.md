# 09 - Atualização de Templates e Arquivos Base

Esta especificação define o comando `specs update` para atualizar arquivos de template, cursor rules e checklist com as versões mais recentes do boilerplate interno do CLI.

## 1. Contexto e Objetivo

- **Contexto:** Projetos SDD inicializados com versões anteriores do CLI podem ter templates e arquivos base desatualizados. Desenvolvedores precisam de uma forma automatizada de atualizar esses arquivos sem perder personalizações feitas no `.cursorrules`.
- **Objetivo:**
  - Atualizar templates estáticos (checklist.md, template-default.spec.md) com versões mais recentes
  - Atualizar `.cursorrules` preservando regras personalizadas do projeto
  - Detectar conflitos entre versão atual e versão do boilerplate
  - Fornecer mecanismo seguro de merge para `.cursorrules` quando houver personalizações
- **Escopo:**
  - Atualização de arquivos base do projeto (templates e checklists)
  - Detecção e merge de `.cursorrules` com regras personalizadas
  - Backup automático de arquivos antes de atualização
  - Validação de pré-condições (projeto SDD válido)
  - Fora de escopo: atualização de specs de funcionalidades (01-*.spec.md, 02-*.spec.md, etc.), atualização de código do projeto, atualização de dependências

## 2. Requisitos Funcionais

- **RF01 - Validação de Pré-condições:**
  - Comando deve verificar se diretório atual contém projeto SDD válido (presença de `specs/` com pelo menos um arquivo `00-*.spec.md`)
  - Se não for projeto SDD, exibir erro e retornar código de saída 2
  - Verificar permissões de escrita no diretório `specs/` e na raiz do projeto
  - Comando deve funcionar apenas em projetos SDD existentes

- **RF02 - Backup de Arquivos Antes de Atualização:**
  - Criar backup de todos os arquivos que serão atualizados antes de modificá-los
  - Backup deve ser criado em diretório `.specs-backup/` na raiz do projeto
  - Estrutura de backup: `.specs-backup/{timestamp}/` contendo cópias dos arquivos
  - Timestamp no formato: `YYYYMMDD-HHMMSS` (ex.: `20240115-143022`)
  - Manter apenas últimos 5 backups (remover backups mais antigos)
  - Exibir mensagem informando localização do backup

- **RF03 - Atualização de Templates Estáticos:**
  - Atualizar `specs/checklist.md` com versão do boilerplate interno
  - Atualizar `specs/template-default.spec.md` com versão do boilerplate interno
  - Arquivos devem ser sobrescritos completamente (não há personalização esperada)
  - Exibir mensagem informando quais templates foram atualizados

- **RF04 - Detecção de Regras Personalizadas em .cursorrules:**
  - Comparar `.cursorrules` do projeto com versão do boilerplate interno
  - Identificar seções que existem apenas no projeto (regras personalizadas)
  - Identificar seções que foram modificadas no projeto em relação ao boilerplate
  - Considerar como personalização:
    - Seções completamente novas (não existem no boilerplate)
    - Seções existentes com conteúdo significativamente diferente (mais de 3 linhas de diferença)
    - Comentários ou regras adicionais dentro de seções existentes
  - Se não houver personalizações detectadas, atualizar `.cursorrules` diretamente

- **RF05 - Merge de .cursorrules com Personalizações:**
  - Se personalizações forem detectadas, criar arquivo `.cursorrules-updated` com versão do boilerplate
  - Preservar `.cursorrules` original intacto
  - Exibir mensagem informando que merge manual é necessário
  - Opcionalmente (se implementado), tentar merge automático:
    - Preservar seções personalizadas do projeto
    - Atualizar seções que existem no boilerplate mas não foram personalizadas
    - Adicionar novas seções do boilerplate que não existem no projeto
    - Criar arquivo `.cursorrules-merged` com resultado do merge automático
    - Exibir mensagem informando que merge automático foi tentado e pedir revisão

- **RF06 - Validação Pós-Atualização:**
  - Verificar que todos os arquivos foram atualizados corretamente
  - Validar que arquivos atualizados são legíveis e têm formato válido
  - Exibir resumo do que foi atualizado
  - Exibir instruções sobre próximos passos (se merge manual for necessário)

- **RF07 - Modo Dry-Run:**
  - Flag `--dry-run` deve exibir o que seria atualizado sem fazer alterações
  - Listar arquivos que seriam atualizados
  - Indicar se `.cursorrules` tem personalizações e se merge seria necessário
  - Não criar backups nem modificar arquivos

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs update [diretório]`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--dry-run`: Exibe o que seria atualizado sem fazer alterações
  - `--force`: Força atualização mesmo se não houver diferenças detectadas
  - `--no-backup`: Não cria backup antes de atualizar (não recomendado)
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - `[diretório]` (opcional): Diretório do projeto SDD a atualizar. Se omitido, usa diretório atual
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - arquivos atualizados ou já estão atualizados
  - `1`: Erro - falha ao atualizar arquivos ou criar backups
  - `2`: Erro - input inválido (não é projeto SDD, diretório não existe, sem permissão)
- **Output:**
  - **Sucesso (stdout):** Mensagem informativa sobre o que foi atualizado
  - **Aviso (stderr):** Avisos sobre merge necessário ou arquivos que não foram atualizados
  - **Erro (stderr):** Mensagem de erro descritiva
- **Exemplos de uso:**
  ```bash
  $ specs update
  Backup criado em .specs-backup/20240115-143022
  Atualizando templates...
  ✓ checklist.md atualizado
  ✓ template-default.spec.md atualizado
  Atualizando .cursorrules...
  ⚠ Regras personalizadas detectadas em .cursorrules
  Arquivo .cursorrules-updated criado com versão do boilerplate
  Execute merge manual ou use: specs update --merge
  
  $ specs update --dry-run
  Arquivos que seriam atualizados:
  - specs/checklist.md
  - specs/template-default.spec.md
  - .cursorrules (merge necessário - regras personalizadas detectadas)
  
  $ specs update --merge
  Backup criado em .specs-backup/20240115-143022
  Tentando merge automático de .cursorrules...
  ✓ Merge automático concluído
  Arquivo .cursorrules-merged criado
  Revise o arquivo e substitua .cursorrules se estiver correto
  
  $ specs update --help
  Atualiza templates e arquivos base do projeto SDD.
  
  Uso:
    specs update [diretório] [flags]
  
  Flags:
    --dry-run      Exibe o que seria atualizado sem fazer alterações
    --force        Força atualização mesmo se não houver diferenças
    --no-backup    Não cria backup antes de atualizar
    --merge        Tenta merge automático de .cursorrules (experimental)
    --help         Exibe esta ajuda
  ```

## 4. Fluxos e Estados

### Fluxo Feliz (Sem Personalizações)

1. Usuário executa `specs update` em projeto SDD válido
2. CLI valida pré-condições (projeto SDD, permissões)
3. CLI cria backup em `.specs-backup/{timestamp}/`
4. CLI atualiza `checklist.md` e `template-default.spec.md`
5. CLI compara `.cursorrules` com boilerplate (sem personalizações detectadas)
6. CLI atualiza `.cursorrules` diretamente
7. CLI valida arquivos atualizados
8. CLI exibe resumo de atualizações
9. Retorna código 0

### Fluxo com Personalizações em .cursorrules

1. Usuário executa `specs update` em projeto SDD válido
2. CLI valida pré-condições
3. CLI cria backup
4. CLI atualiza templates (checklist.md, template-default.spec.md)
5. CLI compara `.cursorrules` e detecta personalizações
6. CLI cria `.cursorrules-updated` com versão do boilerplate
7. CLI preserva `.cursorrules` original
8. CLI exibe aviso sobre merge necessário
9. Se flag `--merge` for usada:
   - CLI tenta merge automático
   - CLI cria `.cursorrules-merged` com resultado
   - CLI exibe instruções para revisão
10. Retorna código 0 (com avisos)

### Estados Alternativos

- **Erro: Não é projeto SDD**
  - Mensagem: "erro: diretório não contém projeto SDD válido"
  - Código: 2
  - Ação: Verificar se está no diretório correto ou executar `specs init` primeiro

- **Erro: Sem permissão de escrita**
  - Mensagem: "erro: sem permissão de escrita em {caminho}"
  - Código: 2
  - Ação: Verificar permissões do diretório

- **Erro: Falha ao criar backup**
  - Mensagem: "erro: falha ao criar backup: {detalhes}"
  - Código: 1
  - Ação: Verificar espaço em disco e permissões

- **Erro: Falha ao atualizar arquivo**
  - Mensagem: "erro: falha ao atualizar {arquivo}: {detalhes}"
  - Código: 1
  - Ação: Verificar permissões e espaço em disco

- **Aviso: Arquivo não existe no boilerplate**
  - Mensagem: "aviso: {arquivo} não encontrado no boilerplate, pulando"
  - Código: 0 (continua atualização)
  - Ação: Nenhuma (arquivo pode ter sido removido do boilerplate)

- **Aviso: Merge automático falhou**
  - Mensagem: "aviso: merge automático de .cursorrules falhou, merge manual necessário"
  - Código: 0 (arquivos base foram atualizados)
  - Ação: Fazer merge manual comparando `.cursorrules` e `.cursorrules-updated`

## 5. Dados

- **Backup:**
  - Localização: `.specs-backup/{timestamp}/`
  - Estrutura:
    ```
    .specs-backup/
    └── 20240115-143022/
        ├── specs/
        │   ├── checklist.md
        │   └── template-default.spec.md
        └── .cursorrules
    ```
  - Retenção: Manter últimos 5 backups, remover mais antigos
  - Permissões: 644 para arquivos, 755 para diretórios

- **Arquivos de Merge:**
  - `.cursorrules-updated`: Versão do boilerplate (criado apenas se houver personalizações)
  - `.cursorrules-merged`: Resultado do merge automático (criado apenas se `--merge` for usado)
  - Permissões: 644
  - Localização: Raiz do projeto

- **Estado:**
  - Não persiste estado entre execuções
  - Não modifica configuração do usuário
  - Não altera specs de funcionalidades (01-*.spec.md, etc.)

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Atualização completa: < 1s
  - Criação de backup: < 200ms
  - Detecção de personalizações: < 100ms
  - Merge automático (se implementado): < 500ms

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Respeita permissões do sistema de arquivos
  - Funciona com caminhos relativos e absolutos
  - Compatível com projetos SDD criados com versões anteriores do CLI

- **Segurança:**
  - Sempre cria backup antes de modificar arquivos
  - Não sobrescreve arquivos sem backup (exceto com `--no-backup`)
  - Valida permissões antes de modificar arquivos
  - Não modifica arquivos fora do diretório do projeto

- **Observabilidade:**
  - Mensagens claras sobre o que está sendo atualizado
  - Logs de erros descritivos em stderr
  - Resumo de atualizações em stdout
  - Flag `--dry-run` para preview sem alterações

## 7. Guardrails

- **Restrições:**
  - Não atualiza specs base (00-*.spec.md) - são arquivos customizados
  - Não atualiza specs de funcionalidades (01-*.spec.md, 02-*.spec.md, etc.)
  - Não modifica código do projeto
  - Não atualiza dependências
  - Não remove arquivos (apenas atualiza existentes)
  - Não cria arquivos que não existem no boilerplate

- **Convenções:**
  - Backup sempre criado antes de modificações (exceto com `--no-backup`)
  - Arquivos atualizados mantêm permissões originais (644)
  - Estrutura de diretórios não é modificada
  - Nomes de arquivos seguem convenções SDD

- **Padrão de mensagens:**
  - Sucesso: "✓ {arquivo} atualizado"
  - Aviso: "⚠ {mensagem}" em stderr
  - Erro: "erro: {descrição}" em stderr
  - Info: "{mensagem}" em stdout

## 8. Critérios de Aceite

- [ ] Comando `specs update` atualiza templates estáticos (checklist.md, template-default.spec.md)
- [ ] Comando `specs update` cria backup antes de atualizar arquivos
- [ ] Comando `specs update` detecta regras personalizadas em .cursorrules
- [ ] Comando `specs update` cria .cursorrules-updated quando há personalizações
- [ ] Comando `specs update` preserva .cursorrules original quando há personalizações
- [ ] Comando `specs update --merge` tenta merge automático de .cursorrules
- [ ] Comando `specs update --dry-run` exibe preview sem fazer alterações
- [ ] Comando `specs update --no-backup` atualiza sem criar backup
- [ ] Comando `specs update` valida que diretório é projeto SDD válido
- [ ] Comando `specs update` retorna código 0 em caso de sucesso
- [ ] Comando `specs update` retorna código 1 para erros de I/O
- [ ] Comando `specs update` retorna código 2 para erros de input inválido
- [ ] Comando `specs update --help` exibe ajuda do comando
- [ ] Backup mantém apenas últimos 5 backups (remove mais antigos)
- [ ] Arquivos atualizados mantêm permissões corretas (644)
- [ ] Templates atualizados são idênticos aos do boilerplate interno

## 9. Testes

### Testes de Unidade

- Validação de pré-condições (detecção de projeto SDD)
- Detecção de personalizações em .cursorrules
- Comparação de arquivos (diferenças detectadas corretamente)
- Criação e limpeza de backups
- Merge automático de .cursorrules (se implementado)

### Testes de Integração

- Fluxo completo de atualização em projeto SDD temporário
- Atualização com personalizações em .cursorrules
- Atualização com flag --dry-run
- Atualização com flag --merge
- Criação e limpeza de backups (manter últimos 5)
- Validação de arquivos atualizados

### Testes E2E

- Execução de `specs update` em projeto SDD real
- Execução de `specs update --dry-run` exibe preview correto
- Execução de `specs update --merge` tenta merge automático
- Verificação de backup criado corretamente
- Verificação de arquivos atualizados são idênticos ao boilerplate
- Verificação de .cursorrules-updated criado quando há personalizações
- Verificação de permissões dos arquivos atualizados

### Como Rodar

- Testes unitários: `go test ./internal/services/update/...`
- Testes de integração: `go test -tags=integration ./internal/services/update/...`
- Testes E2E: Executar manualmente em projeto SDD de teste

## 10. Migração / Rollback

### Migração

- Não há migração necessária (comando atualiza arquivos existentes)
- Usuário pode executar `specs update` a qualquer momento
- Primeira execução cria backup antes de atualizar

### Rollback

- Usuário pode restaurar arquivos do backup em `.specs-backup/{timestamp}/`
- CLI não fornece comando de rollback automático (restauração manual)
- Backup contém cópias completas dos arquivos antes da atualização
- Usuário pode comparar arquivos atualizados com backup para verificar mudanças

## 11. Observações Operacionais

### Distribuição

- Templates do boilerplate devem ser embarcados no binário ou distribuídos junto com o CLI
- Versão do boilerplate deve corresponder à versão do CLI
- Atualizações do CLI podem incluir novos templates ou melhorias em templates existentes

### Versionamento

- Comando `specs update` atualiza para versão do boilerplate do CLI atual
- Usuário pode executar `specs update` após atualizar o CLI para obter templates mais recentes
- Não há verificação de versão (sempre atualiza para versão do boilerplate do CLI)

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Atualização de specs base (00-*.spec.md) - são arquivos customizados e não devem ser atualizados
- Atualização de specs de funcionalidades (01-*.spec.md, 02-*.spec.md, etc.)
- Atualização de código do projeto
- Atualização de dependências
- Verificação de versão do boilerplate vs versão do projeto
- Rollback automático
- Merge interativo de .cursorrules
- Atualização seletiva (escolher quais arquivos atualizar)
- Suporte a múltiplos projetos simultaneamente

### Decisões em Aberto

- Estratégia de merge automático de .cursorrules (heurística vs diff3)
- Formato de backup (tar.gz vs diretório)
- Política de retenção de backups (5 vs configurável)
- Detecção de personalizações (threshold de diferenças)

## Checklist Rápido (preencha antes de gerar código)

- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?
