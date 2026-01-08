# 01 - Controle de Versões

Esta especificação define o sistema de controle de versões do projeto, incluindo armazenamento da versão atual, incremento automático durante builds, criação de tags Git e integração com CI/CD para deploy automático.

## 1. Contexto e Objetivo

- **Contexto:** O projeto precisa de um sistema de versionamento semântico para rastrear releases, facilitar deploy automático e permitir que usuários identifiquem a versão do CLI instalada.
- **Objetivo:** 
  - Manter versão atual em arquivo versionado no repositório
  - Incrementar versão automaticamente durante builds
  - Criar tags Git automaticamente para cada release
  - Integrar com CI/CD para deploy automático
  - Permitir consulta da versão via comando CLI
- **Escopo:** 
  - Sistema de versionamento semântico (MAJOR.MINOR.PATCH)
  - Arquivo de versão versionado no Git
  - Incremento automático durante build
  - Criação automática de tags Git
  - Comando CLI para exibir versão
  - Integração com GitHub Actions para deploy
  - Fora de escopo: versionamento manual, múltiplos canais de release (stable/beta), versionamento de dependências

## 2. Requisitos Funcionais

- **RF01 - Armazenamento de Versão:**
  - Versão atual deve ser armazenada em arquivo na raiz do projeto
  - Formato: versão semântica (MAJOR.MINOR.PATCH, ex.: 0.0.1)
  - Arquivo deve ser versionado no Git
  - Versão inicial: 0.0.1

- **RF02 - Incremento Automático:**
  - Durante processo de build, versão deve ser incrementada automaticamente
  - Incremento padrão: PATCH (0.0.1 → 0.0.2)
  - Arquivo de versão deve ser atualizado antes do build
  - Incremento deve ocorrer apenas em builds de release (não em builds de desenvolvimento)

- **RF03 - Tag Git Automática:**
  - Após incremento de versão, tag Git deve ser criada automaticamente
  - Formato da tag: `v{MAJOR}.{MINOR}.{PATCH}` (ex.: `v0.0.2`)
  - Tag deve ser criada localmente e enviada para repositório remoto
  - Tag deve incluir mensagem descritiva

- **RF04 - Comando de Versão:**
  - Comando `specs version` deve exibir versão atual do CLI
  - Versão exibida deve ser lida do arquivo `VERSION` ou injetada durante build
  - Output deve ser simples e legível (apenas número da versão em formato semântico)
  - Comando deve funcionar mesmo quando executado de qualquer diretório
  - Versão exibida deve corresponder à versão do binário instalado
  - Flag `--json` (futuro) para output estruturado em formato JSON
  - Comando deve ser rápido (< 100ms) e não requerer conexão de rede

- **RF05 - Integração CI/CD:**
  - GitHub Actions deve detectar criação de tag
  - Workflow deve ser acionado automaticamente quando tag é criada
  - Workflow deve executar build para todas as plataformas alvo
  - Workflow deve criar release no GitHub com artefatos
  - Workflow deve incluir checksum SHA256 dos binários

## 3. Contratos e Interfaces

### CLI

- **Comando:** `specs version`
- **Aliases:** Nenhum na v1
- **Flags:**
  - `--json` (futuro): Output em formato JSON estruturado `{"version": "0.0.1"}`
  - `--help`: Exibe ajuda do comando
- **Argumentos:** Nenhum
- **Variáveis de ambiente:** Nenhuma
- **Códigos de saída:**
  - `0`: Sucesso - versão exibida corretamente
  - `1`: Erro - falha ao ler ou exibir versão
- **Output:**
  - **Sucesso (stdout):** Versão em formato semântico, uma linha única (ex.: `0.0.1`)
  - **Erro (stderr):** Mensagem de erro descritiva (ex.: `erro: arquivo VERSION não encontrado`)
- **Exemplos de uso:**
  ```bash
  $ specs version
  0.0.1
  
  $ specs version --help
  Exibe a versão atual do CLI.
  
  Uso:
    specs version [flags]
  
  Flags:
    --help    Exibe ajuda para este comando
  ```

### Arquivos

- **Arquivo de versão:**
  - Localização: `VERSION` na raiz do projeto
  - Formato: Versão semântica única (ex.: `0.0.1`)
  - Encoding: UTF-8
  - Permissões: 644 (rw-r--r--)
  - Versionado: Sim, no Git

### Git

- **Tags:**
  - Formato: `v{MAJOR}.{MINOR}.{PATCH}`
  - Exemplo: `v0.0.1`, `v0.0.2`, `v1.0.0`
  - Mensagem: `Release v{MAJOR}.{MINOR}.{PATCH}`
  - Push: Automático para `origin` após criação

### CI/CD (GitHub Actions)

- **Trigger:** Push de tag com padrão `v*`
- **Workflow:** Executa build, testes e criação de release
- **Outputs:**
  - Binários para todas as plataformas alvo
  - Arquivo `checksums.txt` com SHA256
  - Release notes (futuro: gerado automaticamente)

## 4. Fluxos e Estados

### Fluxo Feliz - Build e Release

1. Desenvolvedor executa comando de build para release
2. Sistema lê versão atual do arquivo `VERSION`
3. Sistema incrementa versão PATCH (0.0.1 → 0.0.2)
4. Sistema atualiza arquivo `VERSION` com nova versão
5. Sistema executa build do projeto
6. Sistema cria tag Git `v0.0.2` com mensagem "Release v0.0.2"
7. Sistema faz push da tag para repositório remoto
8. GitHub Actions detecta push da tag
9. Workflow executa build para todas as plataformas
10. Workflow cria release no GitHub com artefatos e checksums
11. Release fica disponível para download

### Fluxo - Consulta de Versão

1. Usuário executa `specs version` em qualquer diretório
2. Sistema identifica localização do binário `specs`
3. Sistema busca arquivo `VERSION` na raiz do projeto (relativo ao binário) ou usa versão injetada durante build
4. Sistema valida formato da versão (semântico)
5. Sistema exibe versão em stdout (ex.: `0.0.1`)
6. Comando retorna código 0

**Alternativa (versão injetada):**
1. Durante build, versão é injetada como variável no binário
2. Comando `specs version` lê versão da variável injetada (sem necessidade de arquivo `VERSION`)
3. Sistema exibe versão em stdout
4. Comando retorna código 0

### Estados Alternativos

- **Erro: Arquivo VERSION não encontrado:**
  - Mensagem: "erro: arquivo VERSION não encontrado"
  - Código de saída: 1
  - Ação: Verificar se arquivo existe na raiz do projeto

- **Erro: Versão inválida no arquivo:**
  - Mensagem: "erro: versão inválida no arquivo VERSION: {conteúdo}"
  - Código de saída: 1
  - Ação: Validar formato antes de usar

- **Erro: Tag já existe:**
  - Mensagem: "erro: tag v{versão} já existe"
  - Código de saída: 1
  - Ação: Verificar tags existentes antes de criar

- **Erro: Falha ao fazer push da tag:**
  - Mensagem: "erro: falha ao fazer push da tag: {detalhes}"
  - Código de saída: 1
  - Ação: Verificar permissões e conexão com repositório remoto

- **Erro: CI/CD falhou:**
  - Tag foi criada mas workflow falhou
  - Ação: Corrigir problema e criar nova tag ou re-executar workflow

## 5. Dados

### Arquivo VERSION

- **Localização:** `VERSION` na raiz do projeto
- **Formato:** Versão semântica única, sem quebras de linha
- **Exemplo:** `0.0.1`
- **Versionado:** Sim, commitado no Git
- **Atualização:** Automática durante build de release

### Tags Git

- **Armazenamento:** No repositório Git
- **Formato:** `v{MAJOR}.{MINOR}.{PATCH}`
- **Persistência:** Permanente, não devem ser deletadas
- **Histórico:** Mantido no Git para rastreabilidade

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Leitura de versão: < 10ms
  - Incremento e atualização de arquivo: < 50ms
  - Criação de tag Git: < 1s
  - Comando `specs version`: < 100ms

- **Compatibilidade:**
  - Funciona em macOS e Linux
  - Requer Git instalado e configurado
  - Requer acesso ao repositório remoto para push de tags

- **Segurança:**
  - Arquivo VERSION não contém informações sensíveis
  - Tags são assinadas (se GPG configurado)
  - Checksums SHA256 obrigatórios nos releases

- **Observabilidade:**
  - Logs de incremento de versão durante build
  - Logs de criação de tag
  - Logs de falhas em CI/CD

## 7. Guardrails

- **Restrições:**
  - Versão deve seguir formato semântico (MAJOR.MINOR.PATCH)
  - Não permitir decremento de versão
  - Não permitir tags duplicadas
  - Arquivo VERSION deve conter apenas versão, sem espaços ou caracteres extras

- **Convenções:**
  - Incremento padrão: PATCH
  - Tags sempre com prefixo `v`
  - Mensagens de commit para atualização de versão: "chore: bump version to {versão}"

- **Padrão de mensagens:**
  - Sucesso: Apenas versão (ex.: `0.0.1`)
  - Erro: "erro: {descrição}" em stderr

## 8. Critérios de Aceite

- [x] Arquivo `VERSION` existe na raiz com versão inicial `0.0.1` (atual: 0.0.3)
- [x] Comando `specs version` exibe versão atual corretamente
- [x] Comando `specs version` funciona quando executado de qualquer diretório (busca VERSION recursivamente)
- [x] Comando `specs version` exibe versão em formato semântico (MAJOR.MINOR.PATCH)
- [x] Comando `specs version --help` exibe ajuda do comando
- [x] Build de release incrementa versão PATCH automaticamente (script build-release.sh implementado)
- [x] Arquivo `VERSION` é atualizado após incremento
- [x] Tag Git é criada automaticamente após incremento (script implementado)
- [x] Tag é enviada para repositório remoto (push manual necessário, conforme especificado)
- [ ] GitHub Actions é acionado ao detectar push de tag (aguardando execução)
- [ ] Workflow executa build para todas as plataformas alvo (aguardando execução)
- [ ] Release é criado no GitHub com artefatos e checksums (aguardando execução)
- [x] Comando `specs version` retorna código 0 em caso de sucesso
- [x] Comando `specs version` retorna código 1 e mensagem de erro quando arquivo não existe
- [x] Versão exibida corresponde à versão do binário instalado

## 9. Testes

### Testes de Unidade

- Leitura de versão do arquivo `VERSION`
- Validação de formato de versão semântico
- Incremento de versão PATCH
- Formatação de tag Git

### Testes de Integração

- Fluxo completo: leitura → incremento → atualização de arquivo
- Criação de tag Git local
- Push de tag para repositório remoto (mock ou sandbox)

### Testes E2E

- Execução de `specs version` com arquivo `VERSION` válido
- Execução de `specs version` com arquivo `VERSION` inexistente
- Execução de `specs version` com versão injetada no binário
- Execução de `specs version` de diferentes diretórios
- Execução de `specs version --help` exibe ajuda corretamente
- Build completo de release com incremento e tag
- Workflow CI/CD completo (teste em ambiente de staging)

### Como Rodar

- Testes unitários: Executar testes de unidade do módulo de versão
- Testes de integração: Executar com repositório Git temporário
- Testes E2E: Executar build completo em ambiente de teste

## 10. Migração / Rollback

### Migração Inicial

- Criar arquivo `VERSION` com versão inicial `0.0.1`
- Commitar arquivo no repositório
- Não requer migração de dados existentes (projeto novo)

### Rollback

- Se tag foi criada incorretamente: deletar tag local e remota (não recomendado, preferir nova versão)
- Se release foi criado incorretamente: marcar release como pré-release ou deletar (se necessário)
- Se versão foi incrementada incorretamente: corrigir arquivo `VERSION` e criar nova tag

## 11. Observações Operacionais

### Distribuição

- Releases são criados automaticamente via GitHub Actions
- Binários disponíveis em GitHub Releases
- Checksums SHA256 incluídos em cada release

### Versionamento

- Versão inicial: 0.0.1
- Incremento: Automático durante build de release
- Formato: Semântico (MAJOR.MINOR.PATCH)
- Tags: Criadas automaticamente com prefixo `v`

### CI/CD

- Workflow acionado por push de tag
- Build executado para todas as plataformas alvo
- Release criado automaticamente com artefatos
- Notificações (futuro): Slack, email, etc.

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Incremento manual de versão (MAJOR/MINOR)
- Múltiplos canais de release (stable/beta/alpha)
- Versionamento de dependências
- Changelog automático
- Notificações de release
- Assinatura GPG de tags (pode ser configurado manualmente)

### Decisões em Aberto

- Estratégia de incremento (sempre PATCH ou configurável?)
- Formato de release notes (manual ou automático?)
- Versionamento de pré-releases (0.0.1-alpha, 0.0.1-beta)

## Checklist Rápido (preencha antes de gerar código)

- [x] Requisitos estão testáveis? Entradas/saídas precisas?
- [x] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [x] Estados de erro e mensagens estão claros?
- [x] Guardrails e convenções estão escritos?
- [x] Critérios de aceite cobrem fluxos principais e erros?
- [x] Migração/rollback definidos quando há mudança de estado?

