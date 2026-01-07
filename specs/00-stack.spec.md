# 00 - Especificação de Stack Técnica

Esta especificação define a stack tecnológica, ferramentas e plataformas de build/distribuição do projeto. Use-a como referência para todas as decisões técnicas de implementação.

## 1. Contexto e Objetivo
- **Contexto:** Projeto precisa de uma stack tecnológica definida e justificada para garantir consistência, portabilidade e manutenibilidade.
- **Objetivo:** Estabelecer linguagem, runtime, ferramentas de build, empacotamento e plataformas alvo de forma clara e testável.
- **Escopo:** Stack técnica completa (linguagem, ferramentas, build, distribuição). Decisões arquiteturais de alto nível estão em `00-architecture.spec.md`. Contexto global do projeto está em `00-global-context.spec.md`.

## 2. Requisitos Funcionais
- **Artefato portável:** Binário único estático, sem dependências externas
- **Build reproduzível:** Builds devem ser determinísticos e reproduzíveis
- **Empacotamento com checksum:** SHA256 obrigatório para todos os artefatos
- **Cross-compilation:** Suporte a build para múltiplas plataformas (macOS x64/arm64, Linux x64/arm64)
- **Ferramentas padronizadas:** Uso de ferramentas padrão do ecossistema Go

## 3. Stack e Plataformas

### 3.1 Linguagem e Runtime
- **Linguagem:** Go 1.25.5
- **Justificativa:** 
  - **Performance:** Compilação para binário nativo, execução rápida
  - **Portabilidade:** Cross-compilation nativa para múltiplas plataformas
  - **Simplicidade:** Linguagem simples e produtiva, stdlib robusta
  - **Binário único:** Gera binário estático sem dependências externas
  - **Ecossistema:** Ferramentas maduras (gofmt, go test, goreleaser)
  - **Concorrência:** Suporte nativo a goroutines (útil para validação paralela)
- **Runtime:** Binário estático compilado (sem runtime externo)
- **Versão mínima:** Go 1.25.5 (versão atual instalada)

### 3.2 Ferramentas de Build e Desenvolvimento
- **Gerenciamento de dependências:** `go.mod`/`go.sum` (módulos Go)
- **Build tool:** `go build` (stdlib)
- **Formatação:** `gofmt` (stdlib, obrigatório)
- **Linting:** `golangci-lint` (ferramenta externa recomendada, mas opcional na v1)
- **Testes:** `go test` (stdlib)
- **Coverage:** `go test -cover` e `go test -coverprofile=coverage.out` (stdlib)
- **Versionamento/Release:** `goreleaser` (ferramenta externa para releases automatizados)
- **Git Hooks:** Scripts em `.git/hooks/` (pre-commit, pre-push)

### 3.3 Empacotamento e Distribuição
- **Formato:** Binário estático único por plataforma/arquitetura
- **Empacotamento:** 
  - macOS (x64, arm64): binário `specs` em arquivo `.tar.gz`
  - Linux (x64, arm64): binário `specs` em arquivo `.tar.gz`
  - Todos os artefatos incluem arquivo `checksums.txt` com SHA256
- **Checksum:** SHA256 obrigatório para todos os binários (arquivo `checksums.txt`)
- **Ferramenta:** `goreleaser` para releases automatizados (futuro) ou scripts manuais na v1
- **Build flags/opções:** 
  - `CGO_ENABLED=0`: Desabilita CGO para binário estático
  - `-ldflags="-s -w"`: Remove símbolos de debug e reduz tamanho
  - `-trimpath`: Remove caminhos absolutos do binário (build reproduzível)
  - `-tags`: Tags de build se necessário (ex.: `release`)

### 3.4 Plataformas Alvo
- **Plataformas suportadas:** 
  - macOS (x64/amd64, arm64)
  - Linux (x64/amd64, arm64)
- **Arquiteturas:** x64 (amd64), arm64
- **Versões mínimas:** 
  - macOS: 10.15+ (Catalina)
  - Linux: glibc 2.17+ (compatível com distribuições modernas)
- **Fora de escopo:** 
  - Windows (não suportado na v1)
  - Outras arquiteturas (ppc64le, s390x, etc.)
  - Navegadores (não aplicável para CLI)

### 3.5 Compatibilidade e Dependências
- **Runtime:** Binário estático, sem runtime externo necessário
- **Dependências externas:** Sem dependências externas na v1 (apenas stdlib Go)
- **Bibliotecas do sistema:** 
  - Binário estático com `CGO_ENABLED=0` não requer bibliotecas C
  - Compatível com sistemas que não têm glibc (via build estático)
- **Compatibilidade:** 
  - Compatível com todas as distribuições Linux modernas
  - Compatível com macOS 10.15+ em Intel e Apple Silicon
  - Binário único por plataforma/arquitetura

## 4. Estrutura de Build

### 4.1 Diretórios de Build
```
specs/
├── cmd/specs/            # Entry point
├── internal/             # Código interno
├── pkg/                  # Código exportável (se necessário)
├── testdata/             # Arquivos de teste
├── go.mod                # Dependências
├── go.sum                # Checksums de dependências
└── bin/                  # Binários gerados (gitignored)
```

### 4.2 Comandos de Build
- **Desenvolvimento local:** 
  - `go build -o bin/specs ./cmd/specs` (build local)
  - `go run ./cmd/specs` (executar sem build)
- **Build para produção:** 
  - `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/specs ./cmd/specs`
  - Build com flags de otimização e tamanho reduzido
- **Build cross-platform:** 
  - `GOOS=darwin GOARCH=amd64 go build ...` (macOS x64)
  - `GOOS=darwin GOARCH=arm64 go build ...` (macOS arm64)
  - `GOOS=linux GOARCH=amd64 go build ...` (Linux x64)
  - `GOOS=linux GOARCH=arm64 go build ...` (Linux arm64)
  - Ou usar `goreleaser` para build automatizado de todas as plataformas
- **Testes:** 
  - `go test ./...` (todos os testes)
  - `go test -v ./...` (verbose)
- **Coverage:** 
  - `go test -cover ./...` (cobertura resumida)
  - `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out` (relatório HTML)

## 5. Dependências e Bibliotecas

### 5.1 Princípios
- **Preferir stdlib:** Usar bibliotecas padrão sempre que possível
- **Justificar dependências externas:** Cada dependência externa deve ser justificada na spec da feature que a utiliza
- **Minimizar dependências:** Evitar dependências pesadas ou com muitas sub-dependências
- **Versionamento:** Todas as dependências versionadas com versões específicas (sem `latest`)

### 5.2 Bibliotecas Sugeridas (quando necessário)
- **CLI/Parsing:** 
  - Preferir `flag` (stdlib) para flags simples
  - Se necessário parsing complexo: `cobra` ou `urfave/cli` (avaliar necessidade)
- **Configuração:** 
  - `encoding/json` (stdlib) para JSON
  - Sem bibliotecas externas na v1
- **Validação:** 
  - Lógica própria na v1
  - Se necessário: `go-playground/validator` (avaliar necessidade)
- **Justificativa:** 
  - Preferir stdlib sempre que possível
  - Adicionar dependências externas apenas se justificado na spec da feature
  - Minimizar dependências para reduzir tamanho do binário e complexidade

### 5.3 Git Hooks e Ferramentas de Qualidade
- **Git Hooks:** Scripts em `.git/hooks/`
  - `pre-commit`: Executa `gofmt -l` e verifica formatação, executa `golangci-lint` (se configurado)
  - `pre-push`: Executa `go test ./...` e bloqueia push se testes falharem
- **Ferramentas de qualidade:** 
  - `gofmt`: Formatação obrigatória (stdlib)
  - `golangci-lint`: Linting (opcional na v1, recomendado)
  - `go vet`: Análise estática (stdlib, sempre executar)
- **Configuração:** 
  - `.golangci.yml` para configuração do golangci-lint (se usado)
  - Scripts de hooks em `.git/hooks/` (não versionados, mas documentados)

## 6. NFRs (Não Funcionais)

### 6.1 Desempenho
- **Build local:** < 5s para build de desenvolvimento
- **Build completo:** < 30s para build de produção com todas as plataformas (via goreleaser)
- **Tamanho do artefato:** < 10MB por binário (objetivo: < 5MB com otimizações)

### 6.2 Compatibilidade
- **Versão mínima:** Go 1.25.5
- **Compatibilidade de runtime:** Binário estático, sem runtime externo
- **Dependências do sistema:** Sem dependências C (CGO_ENABLED=0), binário totalmente estático

### 6.3 Segurança
- **Auditoria de dependências:** `govulncheck` (ferramenta oficial do Go) antes de releases
- **Checksums/Assinaturas:** SHA256 obrigatório para todos os binários, arquivo `checksums.txt` em cada release
- **Builds reproduzíveis:** 
  - Flags `-trimpath` para remover caminhos absolutos
  - Versionamento fixo de dependências em `go.sum`
  - Builds em containers para garantir reproduzibilidade

### 6.4 Observabilidade
- **Version info:** Injetado via `-ldflags` com variáveis `-X main.version`, `-X main.commit`, `-X main.date`
- **Build info:** 
  - Versão (semantic versioning)
  - Commit hash (git)
  - Data de build
  - Disponível via comando `specs version`

## 7. Guardrails

### 7.1 Restrições de Stack
- **Dependências:** 
  - Nunca adicionar dependência sem justificativa na spec da feature
  - Preferir stdlib Go sempre que possível
  - Evitar dependências que requerem CGO
  - Evitar dependências pesadas ou com muitas sub-dependências
- **Versionamento:** 
  - Sempre usar versões fixas em `go.mod` (sem `latest` ou ranges)
  - Atualizar dependências com `go get -u` e commitar `go.sum`

### 7.2 Convenções de Código
- **Formatação:** `gofmt` obrigatório, código deve passar `gofmt -l` sem mudanças
- **Linting:** `golangci-lint` recomendado (opcional na v1), `go vet` sempre executar
- **Nomenclatura:** 
  - Seguir convenções Go: `PascalCase` para exportados, `camelCase` para privados
  - Interfaces: `PascalCase`, geralmente terminando com `er` se apropriado
  - Arquivos: `snake_case.go` ou `camelCase.go`
- **Estrutura:** Seguir estrutura de diretórios definida em `00-architecture.spec.md`

### 7.3 Build e Release
- **Builds reproduzíveis:** 
  - Usar flags `-trimpath`, `-ldflags="-s -w"`, `CGO_ENABLED=0`
  - Versionamento fixo de dependências
- **Releases:** 
  - v1: Releases manuais com build cross-platform via scripts
  - Futuro: `goreleaser` para releases automatizados
- **Checksums:** Obrigatórios para todos os artefatos (SHA256 em `checksums.txt`)

## 8. Critérios de Aceite

- [ ] Linguagem e versão definidas e documentadas
- [ ] Ferramentas de build e desenvolvimento especificadas
- [ ] Estrutura de diretórios definida e alinhada com arquitetura
- [ ] Plataformas alvo especificadas
- [ ] Formato de empacotamento definido (com checksum/assinatura)
- [ ] Git hooks configurados (se aplicável)
- [ ] Guardrails de dependências e build estabelecidos
- [ ] Critérios de build reproduzível definidos

## 9. Testes

### 9.1 Testes de Build
- Build local funciona (`go build ./cmd/specs`)
- Build de produção funciona com flags otimizadas
- Build cross-platform funciona para todas as plataformas alvo
- Artefato gerado é executável e válido
- Tamanho do binário está dentro dos limites (< 10MB)

### 9.2 Testes de Compatibilidade
- Binário funciona em macOS x64 (testado em CI ou manualmente)
- Binário funciona em macOS arm64 (testado em CI ou manualmente)
- Binário funciona em Linux x64 (testado em CI)
- Binário funciona em Linux arm64 (testado em CI)
- Comando `specs version` retorna informações corretas

### 9.3 Testes de Qualidade
- Formatação ok (`gofmt -l` não retorna mudanças)
- Lint passa (`golangci-lint run` ou `go vet ./...`)
- Testes passam com cobertura adequada (`go test -cover ./...` com cobertura >= 80%)
- Sem vulnerabilidades conhecidas (`govulncheck ./...`)

## 10. Migração / Rollback

### 10.1 Atualização de Versão da Linguagem/Runtime
- **TODO:** Processo de atualização (ex.: atualizar arquivo de dependências, testar em todas as plataformas, atualizar documentação, etc.)

### 10.2 Mudança de Ferramentas
- Documentar mudança em ADR (Architecture Decision Record)
- Atualizar scripts e documentação
- Garantir compatibilidade com builds existentes

## 11. Observações Operacionais

### 11.1 CI/CD
- Build deve rodar em CI para todas as plataformas alvo
- Testes devem rodar em CI antes de merge
- Releases devem ser automatizados no CI

### 11.2 Desenvolvimento Local
- **Requisitos:**
  - Go 1.25.5 instalado e no PATH
  - Git configurado
  - Git hooks configurados (scripts em `.git/hooks/`)
  - (Opcional) `golangci-lint` instalado para linting
- **Setup inicial:**
  - `go mod init` (se necessário)
  - Configurar git hooks (copiar scripts para `.git/hooks/`)
  - Executar `go mod tidy` para sincronizar dependências

## 12. Abertos / Fora de Escopo

- **Plataformas não suportadas:** Windows, outras arquiteturas (ppc64le, s390x, etc.)
- **Ferramentas não usadas:** 
  - Docker (não necessário para CLI estático)
  - Kubernetes (não aplicável)
  - Ferramentas de CI/CD específicas (usar GitHub Actions ou similar)
- **Dependências externas:** Evitar dependências que requerem CGO ou bibliotecas nativas

## Checklist Rápido (preencha antes de gerar código)
- [ ] Stack tecnológica definida e justificada
- [ ] Ferramentas de build e desenvolvimento especificadas
- [ ] Plataformas alvo definidas e testáveis
- [ ] Formato de empacotamento e distribuição definido
- [ ] Guardrails de dependências e build estabelecidos
- [ ] Critérios de aceite cobrem build, compatibilidade e qualidade

