# 07 - Sistema de Configuração

Esta especificação define o sistema de configuração XDG-compliant para o CLI `specs`, permitindo que usuários personalizem comportamentos padrão do CLI através de arquivo de configuração.

## 1. Contexto e Objetivo

- **Contexto:** Usuários podem precisar personalizar comportamentos padrão do CLI, como caminho padrão para specs, preferências de formatação, ou outras opções que afetam o comportamento dos comandos.
- **Objetivo:** 
  - Fornecer sistema de configuração XDG-compliant para personalização do CLI
  - Permitir configuração de caminho padrão para specs
  - Permitir configuração de preferências de formatação e comportamento
  - Fornecer comando para visualizar e editar configuração
  - Manter valores padrão sensatos quando configuração não existe
  - Validar configuração e fornecer mensagens de erro claras
- **Escopo:** 
  - Arquivo de configuração em formato JSON
  - Localização XDG-compliant (`~/.config/specs/config.json`)
  - Comando para visualizar configuração atual
  - Comando para editar configuração (opcional na v1)
  - Validação de formato e valores de configuração
  - Valores padrão quando configuração não existe
  - Fora de escopo: configuração por projeto (`.specs/config.json`), migração automática de versões antigas, interface interativa de edição, validação de schema avançada (v2)

## 2. Requisitos Funcionais

- **RF01 - Arquivo de Configuração:**
  - Criar arquivo de configuração em `~/.config/specs/config.json` (XDG-compliant)
  - Formato JSON com estrutura definida
  - Valores padrão sensatos quando arquivo não existe
  - Criar diretório de configuração automaticamente se não existir
  - Validar formato JSON ao carregar configuração

- **RF02 - Opções de Configuração:**
  - `specs.default_path`: Caminho padrão para diretório de specs (padrão: `./specs`)
  - `specs.exclude_templates`: Excluir specs de template do dashboard (padrão: `true`)
  - Estrutura extensível para futuras opções (v2+)
  - Valores padrão aplicados quando opção não está presente

- **RF03 - Comando de Visualização:**
  - Comando `specs config` ou `specs config show` para exibir configuração atual
  - Exibir configuração em formato legível (JSON formatado ou tabela)
  - Indicar valores padrão quando configuração não existe
  - Exibir caminho do arquivo de configuração

- **RF04 - Comando de Edição (Opcional v1):**
  - Comando `specs config set <chave> <valor>` para definir valor
  - Comando `specs config get <chave>` para obter valor específico
  - Validação de valores antes de salvar
  - Criação automática de arquivo se não existir

- **RF05 - Validação de Configuração:**
  - Validar formato JSON ao carregar
  - Validar tipos de valores (string, boolean, etc.)
  - Validar valores permitidos para cada opção
  - Reportar erros de validação com mensagens claras
  - Fallback para valores padrão em caso de erro de validação

- **RF06 - Integração com Comandos:**
  - Comandos existentes devem ler e usar configuração quando disponível
  - Caminho padrão de specs deve usar `specs.default_path` se configurado
  - Dashboard deve respeitar `specs.exclude_templates` se configurado
  - Valores padrão aplicados quando configuração não existe ou opção não está presente

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs config [subcomando] [flags]`
- **Aliases:** Nenhum na v1
- **Subcomandos:**
  - `show` (padrão): Exibe configuração atual
  - `get <chave>`: Obtém valor de uma chave específica
  - `set <chave> <valor>`: Define valor de uma chave
- **Flags:**
  - `--help`: Exibe ajuda do comando
- **Argumentos:**
  - Para `get`: `<chave>` - nome da chave de configuração (ex.: `specs.default_path`)
  - Para `set`: `<chave> <valor>` - nome da chave e valor a definir
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso
  - `1`: Erro - falha ao ler/escrever configuração
  - `2`: Erro - input inválido (chave não existe, valor inválido)
- **Output:**
  - **Sucesso (stdout):** Configuração formatada ou valor solicitado
  - **Erro (stderr):** Mensagem de erro descritiva
- **Exemplos de uso:**
  ```bash
  $ specs config
  Configuração em: ~/.config/specs/config.json
  
  {
    "specs": {
      "default_path": "./specs",
      "exclude_templates": true
    }
  }
  
  $ specs config get specs.default_path
  ./specs
  
  $ specs config set specs.default_path ./minhas-specs
  Configuração atualizada: specs.default_path = ./minhas-specs
  
  $ specs config --help
  Gerencia configuração do CLI specs.
  
  Uso:
    specs config [subcomando] [flags]
  
  Subcomandos:
    show              Exibe configuração atual (padrão)
    get <chave>       Obtém valor de uma chave específica
    set <chave> <valor>  Define valor de uma chave
  
  Flags:
    --help            Exibe ajuda para este comando
  
  Exemplos:
    specs config                    # Exibe configuração completa
    specs config get specs.default_path  # Obtém caminho padrão
    specs config set specs.default_path ./specs  # Define caminho padrão
  ```

### Arquivos

- **Arquivo de configuração:**
  - Localização: `~/.config/specs/config.json` (XDG_CONFIG_HOME/specs/config.json)
  - Formato: JSON
  - Encoding: UTF-8
  - Permissões: 0600 (rw-------)
  - Estrutura:
    ```json
    {
      "specs": {
        "default_path": "./specs",
        "exclude_templates": true
      }
    }
    ```

- **Valores padrão:**
  - `specs.default_path`: `"./specs"`
  - `specs.exclude_templates`: `true`

## 4. Fluxos e Estados

### Fluxo Feliz - Visualizar Configuração

1. Usuário executa `specs config` ou `specs config show`
2. Sistema verifica se arquivo de configuração existe
3. Se existe, sistema lê e valida JSON
4. Sistema exibe configuração formatada com caminho do arquivo
5. Comando retorna código 0

### Fluxo - Configuração Não Existe

1. Usuário executa `specs config`
2. Sistema verifica se arquivo existe
3. Se não existe, sistema exibe valores padrão
4. Sistema informa que arquivo não existe e onde será criado
5. Comando retorna código 0

### Fluxo - Definir Configuração

1. Usuário executa `specs config set specs.default_path ./minhas-specs`
2. Sistema valida chave (deve existir na estrutura)
3. Sistema valida valor (tipo e formato correto)
4. Sistema lê configuração existente ou cria nova
5. Sistema atualiza valor da chave
6. Sistema salva arquivo de configuração
7. Sistema exibe confirmação
8. Comando retorna código 0

### Estados Alternativos

- **Erro: Arquivo de configuração corrompido:**
  - Mensagem: "erro: arquivo de configuração inválido: {caminho} ({detalhes})"
  - Código de saída: 1
  - Ação: Corrigir JSON ou remover arquivo para usar padrões

- **Erro: Chave não existe:**
  - Mensagem: "erro: chave desconhecida: {chave}"
  - Código de saída: 2
  - Ação: Verificar chaves disponíveis com `specs config`

- **Erro: Valor inválido:**
  - Mensagem: "erro: valor inválido para {chave}: {valor} ({detalhes})"
  - Código de saída: 2
  - Ação: Verificar formato esperado

- **Erro: Falha ao criar diretório:**
  - Mensagem: "erro: falha ao criar diretório de configuração: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar permissões

- **Erro: Falha ao escrever arquivo:**
  - Mensagem: "erro: falha ao salvar configuração: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar permissões

## 5. Dados

### Estrutura de Configuração

- **Formato:** JSON
- **Localização:** `~/.config/specs/config.json` (ou `$XDG_CONFIG_HOME/specs/config.json`)
- **Estrutura:**
  ```json
  {
    "specs": {
      "default_path": string,        // Caminho padrão para specs (padrão: "./specs")
      "exclude_templates": boolean   // Excluir templates do dashboard (padrão: true)
    }
  }
  ```

### Valores Padrão

- Quando arquivo não existe, valores padrão são aplicados
- Valores padrão não são persistidos automaticamente (apenas quando usuário edita)
- Valores padrão são documentados no help do comando

### Políticas de Retenção

- Configuração persiste indefinidamente até ser modificada ou removida
- Não há expurgo automático
- Usuário pode remover arquivo manualmente para voltar aos padrões

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Carregamento de configuração: < 10ms
  - Salvamento de configuração: < 50ms
  - Validação de configuração: < 5ms

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Respeita XDG Base Directory Specification
  - Suporta `$XDG_CONFIG_HOME` quando definido
  - Fallback para `~/.config` quando `XDG_CONFIG_HOME` não está definido

- **Segurança:**
  - Arquivo de configuração com permissões 0600 (apenas leitura/escrita pelo dono)
  - Não armazena informações sensíveis (tokens, senhas)
  - Validação de entrada para prevenir injeção

- **Observabilidade:**
  - Mensagens claras sobre localização do arquivo
  - Indicação quando valores padrão estão sendo usados
  - Logs de erro descritivos em stderr

## 7. Guardrails

- **Restrições:**
  - Formato JSON obrigatório (não suporta outros formatos na v1)
  - Apenas chaves definidas na estrutura são permitidas
  - Valores devem ser do tipo correto (string, boolean)
  - Não modifica configuração sem confirmação explícita do usuário

- **Convenções:**
  - Nomes de chaves: `specs.{opcao}` (namespace `specs`)
  - Caminhos: Relativos ao diretório atual ou absolutos
  - Booleanos: `true` ou `false` (não `1`/`0` ou `yes`/`no`)

- **Padrão de mensagens:**
  - Sucesso: "Configuração atualizada: {chave} = {valor}"
  - Erro: "erro: {descrição do problema}"
  - Info: "Configuração em: {caminho}"

## 8. Critérios de Aceite

- [ ] Comando `specs config` exibe configuração atual quando arquivo existe
- [ ] Comando `specs config` exibe valores padrão quando arquivo não existe
- [ ] Comando `specs config get <chave>` retorna valor da chave específica
- [ ] Comando `specs config set <chave> <valor>` atualiza configuração
- [ ] Sistema cria diretório de configuração automaticamente se não existir
- [ ] Sistema valida formato JSON ao carregar configuração
- [ ] Sistema valida tipos de valores (string, boolean)
- [ ] Sistema reporta erros claros para configuração inválida
- [ ] Comandos existentes usam `specs.default_path` quando configurado
- [ ] Dashboard respeita `specs.exclude_templates` quando configurado
- [ ] Arquivo de configuração tem permissões 0600
- [ ] Comando retorna código 0 em caso de sucesso
- [ ] Comando retorna código 1 para erros de I/O
- [ ] Comando retorna código 2 para erros de input inválido
- [ ] Comando `specs config --help` exibe ajuda do comando
- [ ] Sistema funciona quando `XDG_CONFIG_HOME` está definido
- [ ] Sistema funciona quando `XDG_CONFIG_HOME` não está definido (fallback para ~/.config)

## 9. Testes

### Testes de Unidade

- Carregamento de configuração válida
- Carregamento de configuração inválida (JSON malformado)
- Aplicação de valores padrão quando arquivo não existe
- Validação de tipos de valores
- Validação de chaves desconhecidas
- Criação de diretório de configuração
- Salvamento de configuração
- Permissões de arquivo (0600)

### Testes de Integração

- Fluxo completo de visualização de configuração existente
- Fluxo completo de visualização quando arquivo não existe
- Fluxo completo de definição de valor
- Fluxo completo de obtenção de valor
- Integração com comandos existentes (uso de default_path)
- Integração com dashboard (uso de exclude_templates)
- Respeito a XDG_CONFIG_HOME quando definido

### Testes E2E

- Execução de `specs config` em ambiente sem configuração
- Execução de `specs config` em ambiente com configuração
- Execução de `specs config set` e verificação de persistência
- Execução de `specs config get` para obter valor
- Verificação de que comandos usam configuração corretamente
- Verificação de permissões de arquivo criado
- Teste com `XDG_CONFIG_HOME` definido

### Como Rodar

- Testes unitários: `go test ./internal/services/config/...`
- Testes de integração: `go test -tags=integration ./internal/services/config/...`
- Testes E2E: Executar manualmente em ambiente de teste

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (sistema novo)
- Primeira execução cria arquivo de configuração se usuário editar
- Valores padrão aplicados até configuração ser criada

### Rollback

- Usuário pode remover arquivo `~/.config/specs/config.json` para voltar aos padrões
- Sistema funciona normalmente sem arquivo de configuração
- Não há migração de versões antigas na v1

## 11. Observações Operacionais

### Localização XDG

- Respeita XDG Base Directory Specification
- Usa `$XDG_CONFIG_HOME/specs/config.json` se `XDG_CONFIG_HOME` estiver definido
- Fallback para `~/.config/specs/config.json` se `XDG_CONFIG_HOME` não estiver definido
- Cria diretório `specs/` automaticamente se não existir

### Integração com Comandos Existentes

- Comandos devem ler configuração uma vez por execução
- Cache de configuração pode ser adicionado no futuro (v2)
- Comandos devem funcionar normalmente mesmo se configuração estiver corrompida (usar padrões)

### Segurança

- Arquivo com permissões restritas (0600)
- Não armazena informações sensíveis
- Validação de entrada para prevenir problemas

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Configuração por projeto (`.specs/config.json` no diretório do projeto)
- Migração automática de versões antigas de configuração
- Interface interativa de edição (apenas comandos CLI)
- Validação de schema avançada (JSON Schema)
- Configuração de múltiplos perfis
- Importação/exportação de configuração
- Histórico de mudanças de configuração

### Decisões em Aberto

- Estrutura de configuração para futuras opções (como organizar)
- Estratégia de cache de configuração (se houver no futuro)
- Suporte a configuração por projeto (v2?)

## Checklist Rápido (preencha antes de gerar código)

- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
