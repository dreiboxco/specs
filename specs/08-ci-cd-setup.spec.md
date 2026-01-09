# 08 - Setup e CI/CD

Esta especificação define o processo de setup inicial do projeto, configuração de CI/CD, criação de releases e validação de badges no README.

## 1. Contexto e Objetivo

- **Contexto:** O projeto precisa de documentação clara sobre como configurar o ambiente de desenvolvimento, CI/CD, criar releases e validar que todos os componentes estão funcionando corretamente.
- **Objetivo:** 
  - Documentar processo de setup inicial do projeto
  - Configurar workflows de CI/CD no GitHub Actions
  - Garantir que releases sejam criadas automaticamente
  - Validar que badges no README funcionam corretamente
  - Fornecer checklist para validação de setup completo
- **Escopo:** 
  - Workflows de CI/CD (testes, lint, build)
  - Workflow de release automático
  - Documentação de setup inicial
  - Validação de badges e links
  - Checklist de validação
  - Fora de escopo: configuração de ambientes de staging/produção, deploy automático, monitoramento

## 2. Requisitos Funcionais

- **RF01 - Workflow de CI:**
  - Workflow deve executar testes em push/PR para branches `main` e `development`
  - Workflow deve executar linting de código
  - Workflow deve fazer build para todas as plataformas suportadas
  - Workflow deve gerar relatório de cobertura de testes
  - Workflow deve falhar se testes ou lint falharem

- **RF02 - Workflow de Release:**
  - Workflow deve ser acionado quando tag `v*` é criada
  - Workflow deve ler versão do arquivo `VERSION`
  - Workflow deve compilar binários para todas as plataformas (Linux/macOS, amd64/arm64)
  - Workflow deve gerar arquivos `.tar.gz` com binários
  - Workflow deve gerar checksums SHA256 para todos os arquivos
  - Workflow deve criar release no GitHub automaticamente
  - Workflow deve incluir instruções de instalação na release

- **RF03 - Documentação de Setup:**
  - Deve existir documentação clara sobre setup inicial
  - Deve incluir checklist passo a passo
  - Deve incluir comandos prontos para executar
  - Deve incluir troubleshooting de problemas comuns
  - Deve documentar como criar primeira release

- **RF04 - Validação de Badges:**
  - Badge de CI deve apontar para workflow correto
  - Badge de release deve funcionar após primeira release
  - Todos os badges devem ter links funcionais
  - Badges devem aparecer corretamente no README

## 3. Contratos e Interfaces

### GitHub Actions

- **Workflow de CI:**
  - Arquivo: `.github/workflows/ci.yml`
  - Triggers: `push` e `pull_request` em branches `main` e `development`
  - Jobs: `test`, `lint`, `build`
  - Output: Status de CI visível no GitHub

- **Workflow de Release:**
  - Arquivo: `.github/workflows/release.yml`
  - Trigger: `push` de tags `v*`
  - Jobs: `build` (matriz de plataformas), `release`
  - Output: Release criada no GitHub com arquivos `.tar.gz`

### Documentação

- **Arquivo de Setup:**
  - Localização: `specs/08-ci-cd-setup.spec.md` (esta spec)
  - Formato: Markdown seguindo template padrão
  - Conteúdo: Checklist, comandos, troubleshooting

### Badges

- **Badge de CI:**
  - URL: `https://img.shields.io/github/actions/workflow/status/dreibox/specs/ci.yml`
  - Link: `https://github.com/dreibox/specs/actions/workflows/ci.yml`

- **Badge de Release:**
  - URL: `https://img.shields.io/github/v/release/dreibox/specs`
  - Link: `https://github.com/dreibox/specs/releases`

## 4. Fluxos e Estados

### Fluxo Feliz - Setup Inicial

1. Repositório existe no GitHub (`dreibox/specs`)
2. Remote `origin` está configurado corretamente
3. Código é feito push para GitHub
4. Workflow de CI executa automaticamente
5. CI passa com sucesso
6. Tag `v0.0.3` é criada e enviada
7. Workflow de release executa automaticamente
8. Release é criada no GitHub
9. Badges aparecem corretamente no README

### Fluxo - Criar Release

1. Versão é atualizada no arquivo `VERSION`
2. Mudanças são commitadas
3. Tag é criada: `git tag v{VERSION}`
4. Tag é enviada: `git push origin v{VERSION}`
5. Workflow de release detecta tag
6. Binários são compilados para todas as plataformas
7. Arquivos `.tar.gz` são criados
8. Checksums são gerados
9. Release é criada no GitHub
10. Arquivos ficam disponíveis para download

### Estados Alternativos

- **Erro: Repositório não existe:**
  - Mensagem: Verificar se repositório `dreibox/specs` existe no GitHub
  - Ação: Criar repositório no GitHub

- **Erro: Remote não configurado:**
  - Mensagem: `git remote -v` não mostra origin
  - Ação: Configurar remote: `git remote add origin https://github.com/dreibox/specs.git`

- **Erro: Workflow de CI falha:**
  - Mensagem: Verificar logs no GitHub Actions
  - Ação: Corrigir testes ou lint que estão falhando

- **Erro: Release não é criada:**
  - Mensagem: Verificar se tag foi criada corretamente
  - Ação: Verificar logs do workflow de release

- **Erro: Badge mostra "repo not found":**
  - Mensagem: Repositório não existe ou nome está incorreto
  - Ação: Verificar nome do repositório no GitHub

## 5. Dados

### Arquivos de Workflow

- **`.github/workflows/ci.yml`:**
  - Formato: YAML
  - Localização: `.github/workflows/ci.yml`
  - Conteúdo: Definição de jobs de CI

- **`.github/workflows/release.yml`:**
  - Formato: YAML
  - Localização: `.github/workflows/release.yml`
  - Conteúdo: Definição de jobs de release

### Arquivo de Versão

- **`VERSION`:**
  - Formato: Texto simples (ex: `0.0.3`)
  - Localização: Raiz do projeto
  - Uso: Lido pelos workflows para determinar versão

## 6. NFRs (Não Funcionais)

- **Desempenho:**
  - Workflow de CI deve completar em < 10 minutos
  - Workflow de release deve completar em < 15 minutos
  - Badges devem carregar em < 2 segundos

- **Compatibilidade:**
  - Workflows devem funcionar no GitHub Actions
  - Suporta Go 1.25+
  - Suporta Linux e macOS runners

- **Segurança:**
  - Secrets não devem ser expostos nos workflows
  - Checksums SHA256 para validação de integridade
  - Permissões mínimas necessárias nos workflows

- **Observabilidade:**
  - Logs claros em cada step dos workflows
  - Status visível no GitHub Actions
  - Notificações de falha (se configurado)

## 7. Guardrails

- **Restrições:**
  - Workflows devem usar actions oficiais ou amplamente utilizadas
  - Não adicionar dependências externas desnecessárias
  - Manter workflows simples e legíveis

- **Convenções:**
  - Nomes de workflows: `CI`, `Release`
  - Nomes de jobs: descritivos (ex: `test`, `lint`, `build`)
  - Tags Git: formato `v{MAJOR}.{MINOR}.{PATCH}`

- **Padrão de mensagens:**
  - Commits: Conventional Commits
  - Releases: Incluir changelog básico

## 8. Critérios de Aceite

- [ ] Workflow de CI existe em `.github/workflows/ci.yml`
- [ ] Workflow de CI executa em push para `main` e `development`
- [ ] Workflow de CI executa em pull requests
- [ ] Workflow de CI executa testes com sucesso
- [ ] Workflow de CI executa lint com sucesso
- [ ] Workflow de CI faz build para todas as plataformas
- [ ] Workflow de release existe em `.github/workflows/release.yml`
- [ ] Workflow de release é acionado por tags `v*`
- [ ] Workflow de release compila binários para todas as plataformas
- [ ] Workflow de release gera arquivos `.tar.gz`
- [ ] Workflow de release gera checksums SHA256
- [ ] Workflow de release cria release no GitHub
- [ ] Release contém todos os arquivos necessários
- [ ] Badge de CI aparece e funciona no README
- [ ] Badge de release aparece após primeira release
- [ ] Todos os links de badges funcionam
- [ ] Documentação de setup está completa e clara
- [ ] Checklist de validação está completo

## 9. Testes

### Testes de Validação

- Verificar se workflows existem e estão no formato correto
- Validar sintaxe YAML dos workflows
- Testar criação de tag e execução de release
- Verificar se badges aparecem corretamente

### Como Validar

```bash
# Validar sintaxe YAML (se tiver yamllint)
yamllint .github/workflows/*.yml

# Verificar se workflows existem
ls -la .github/workflows/

# Testar criação de tag (localmente, sem push)
git tag v0.0.3-test
git tag -d v0.0.3-test

# Verificar badges (manualmente no README renderizado)
```

## 10. Migração / Rollback

### Migração Inicial

- Não há migração necessária (setup inicial)
- Workflows são adicionados ao projeto
- Primeira release cria estrutura inicial

### Rollback

- Remover workflows: deletar arquivos `.github/workflows/*.yml`
- Remover tags: `git tag -d v{VERSION}` e `git push origin :refs/tags/v{VERSION}`
- Remover releases: manualmente no GitHub

## 11. Observações Operacionais

### Setup Inicial

1. **Criar/verificar repositório no GitHub:**
   - Nome: `dreibox/specs`
   - Visibilidade: Pública (para badges funcionarem)

2. **Configurar remote:**
   ```bash
   git remote add origin https://github.com/dreibox/specs.git
   # ou verificar se já existe
   git remote -v
   ```

3. **Fazer push inicial:**
   ```bash
   git push -u origin development
   git push -u origin main
   ```

4. **Criar primeira release:**
   ```bash
   git tag v0.0.3
   git push origin v0.0.3
   ```

### Manutenção

- **Incrementar versão:**
  - Atualizar arquivo `VERSION`
  - Commit e push
  - Criar nova tag

- **Monitorar CI/CD:**
  - Verificar status em GitHub Actions
  - Corrigir falhas rapidamente
  - Manter workflows atualizados

## 12. Abertos / Fora de Escopo

### Fora de Escopo (v1)

- Configuração de ambientes de staging/produção
- Deploy automático para produção
- Notificações de falha (email, Slack, etc.)
- Integração com outros serviços (Codecov, etc.)
- Versionamento de dependências
- Changelog automático
- Release notes automáticos baseados em commits

### Decisões em Aberto

- Estratégia de versionamento (automático vs manual)
- Integração com Codecov ou similar
- Notificações de falha de CI/CD

## Checklist Rápido (preencha antes de gerar código)

- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?
