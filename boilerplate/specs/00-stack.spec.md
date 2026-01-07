# 00 - Especificação de Stack Técnica

Esta especificação define a stack tecnológica, ferramentas e plataformas de build/distribuição do projeto. Use-a como referência para todas as decisões técnicas de implementação.

## 1. Contexto e Objetivo
- **Contexto:** Projeto precisa de uma stack tecnológica definida e justificada para garantir consistência, portabilidade e manutenibilidade.
- **Objetivo:** Estabelecer linguagem, runtime, ferramentas de build, empacotamento e plataformas alvo de forma clara e testável.
- **Escopo:** Stack técnica completa (linguagem, ferramentas, build, distribuição). Decisões arquiteturais de alto nível estão em `00-architecture.spec.md`. Contexto global do projeto está em `00-global-context.spec.md`.

## 2. Requisitos Funcionais
- **TODO:** Listar requisitos funcionais da stack (ex.: artefato portável, build reproduzível, empacotamento com checksum, cross-compilation, ferramentas padronizadas, etc.)

## 3. Stack e Plataformas

### 3.1 Linguagem e Runtime
- **Linguagem:** TODO (ex.: Go 1.25.5, Node.js 20.x, Python 3.11, Rust 1.75, etc.)
- **Justificativa:** 
  - TODO: Listar justificativas para escolha da linguagem (ex.: performance, portabilidade, ecossistema, etc.)
- **Runtime:** TODO (ex.: binário estático, Node.js, Python interpreter, JVM, etc.)
- **Versão mínima:** TODO (ex.: Go 1.25.5, Node.js 20.x, Python 3.11, etc.)

### 3.2 Ferramentas de Build e Desenvolvimento
- **Gerenciamento de dependências:** TODO (ex.: `go.mod`/`go.sum`, `package.json`, `requirements.txt`, `Cargo.toml`, etc.)
- **Build tool:** TODO (ex.: `go build`, `npm run build`, `webpack`, `cargo build`, etc.)
- **Formatação:** TODO (ex.: `gofmt`, `prettier`, `black`, `rustfmt`, etc.)
- **Linting:** TODO (ex.: `golangci-lint`, `eslint`, `pylint`, `clippy`, etc.)
- **Testes:** TODO (ex.: `go test`, `jest`, `pytest`, `cargo test`, etc.)
- **Coverage:** TODO (ex.: `go test -cover`, `jest --coverage`, `pytest-cov`, `cargo tarpaulin`, etc.)
- **Versionamento/Release:** TODO (ex.: `goreleaser`, `semantic-release`, `bumpversion`, etc.)
- **Git Hooks:** TODO (ex.: `pre-commit`, `pre-push`, ferramentas específicas, etc.)

### 3.3 Empacotamento e Distribuição
- **Formato:** TODO (ex.: binário estático, pacote npm, imagem Docker, wheel Python, etc.)
- **Empacotamento:** TODO (ex.: tar.gz por plataforma, .deb/.rpm, .dmg/.pkg, etc.)
  - TODO: Listar formatos por plataforma/arquitetura
- **Checksum:** TODO (ex.: SHA256 obrigatório, assinatura GPG, etc.)
- **Ferramenta:** TODO (ex.: `goreleaser`, `npm pack`, `docker build`, etc.)
- **Build flags/opções:** 
  - TODO: Listar flags/opções de build específicas (ex.: `-ldflags`, `CGO_ENABLED=0`, `-trimpath`, etc.)

### 3.4 Plataformas Alvo
- **TODO:** Listar plataformas suportadas (ex.: macOS x64/arm64, Linux x64/arm64, Windows, navegadores, etc.)
- **Arquiteturas:** TODO (ex.: x64, arm64, etc.)
- **Versões mínimas:** TODO (ex.: macOS 10.15+, Linux glibc 2.17+, Windows 10+, etc.)
- **Fora de escopo:** TODO (ex.: Windows na v1, navegadores antigos, etc.)

### 3.5 Compatibilidade e Dependências
- **Runtime:** TODO (ex.: binário estático, requer Node.js, requer Python, etc.)
- **Dependências externas:** TODO (ex.: sem dependências, requer runtime X, etc.)
- **Bibliotecas do sistema:** TODO (ex.: glibc, musl, libc, etc.)
- **Compatibilidade:** TODO (ex.: compatível com distribuições X, versões Y, etc.)

## 4. Estrutura de Build

### 4.1 Diretórios de Build
```
TODO: Definir estrutura de diretórios específica da stack

Exemplos:
- Go: cmd/, internal/, pkg/, go.mod, go.sum
- Node.js: src/, dist/, package.json, node_modules/
- Python: src/, tests/, requirements.txt, setup.py
- Rust: src/, Cargo.toml, target/
```

### 4.2 Comandos de Build
- **Desenvolvimento local:** TODO (ex.: `go build`, `npm run dev`, `python -m build`, etc.)
- **Build para produção:** TODO (ex.: `go build -o bin/app`, `npm run build`, etc.)
- **Build cross-platform:** TODO (ex.: `goreleaser`, `electron-builder`, etc.)
- **Testes:** TODO (ex.: `go test ./...`, `npm test`, `pytest`, etc.)
- **Coverage:** TODO (ex.: `go test -cover`, `jest --coverage`, `pytest-cov`, etc.)

## 5. Dependências e Bibliotecas

### 5.1 Princípios
- **Preferir stdlib:** Usar bibliotecas padrão sempre que possível
- **Justificar dependências externas:** Cada dependência externa deve ser justificada na spec da feature que a utiliza
- **Minimizar dependências:** Evitar dependências pesadas ou com muitas sub-dependências
- **Versionamento:** Todas as dependências versionadas com versões específicas (sem `latest`)

### 5.2 Bibliotecas Sugeridas (quando necessário)
- **TODO:** Listar bibliotecas sugeridas por categoria (ex.: CLI/Parsing, Configuração, HTTP Client, Validação, etc.)
- **Justificativa:** TODO (ex.: por que cada biblioteca foi sugerida, quando usar, etc.)

### 5.3 Git Hooks e Ferramentas de Qualidade
- **Git Hooks:** TODO (ex.: scripts em `.git/hooks/`, ferramenta `pre-commit`, Husky, etc.)
  - `pre-commit`: TODO (ex.: executa formatação e lint antes de commit)
  - `pre-push`: TODO (ex.: executa testes antes de push, bloqueia se falhar)
- **Ferramentas de qualidade:** TODO (ex.: linter, formatter, type checker, etc.)
- **Configuração:** TODO (ex.: arquivos de configuração, onde estão, etc.)

## 6. NFRs (Não Funcionais)

### 6.1 Desempenho
- **Build local:** TODO (ex.: < 10s, < 30s, etc.)
- **Build completo:** TODO (ex.: < 2min, < 5min, etc.)
- **Tamanho do artefato:** TODO (ex.: < 20MB, < 100MB, etc.)

### 6.2 Compatibilidade
- **Versão mínima:** TODO (ex.: Go 1.25.5, Node.js 20.x, Python 3.11, etc.)
- **Compatibilidade de runtime:** TODO (ex.: binário estático, requer runtime X, etc.)
- **Dependências do sistema:** TODO (ex.: sem dependências C, requer bibliotecas Y, etc.)

### 6.3 Segurança
- **Auditoria de dependências:** TODO (ex.: `govulncheck`, `npm audit`, `safety`, `cargo audit`, etc.)
- **Checksums/Assinaturas:** TODO (ex.: SHA256 obrigatório, GPG, etc.)
- **Builds reproduzíveis:** TODO (ex.: via flags específicas, versionamento fixo, etc.)

### 6.4 Observabilidade
- **Version info:** TODO (ex.: injetado via build flags, disponível via comando, etc.)
- **Build info:** TODO (ex.: timestamp, commit hash, versão, etc.)

## 7. Guardrails

### 7.1 Restrições de Stack
- **TODO:** Listar restrições de dependências (ex.: nunca adicionar sem justificativa, preferir stdlib, evitar dependências nativas, etc.)
- **Versionamento:** TODO (ex.: sempre usar versões fixas, nunca `latest`, etc.)

### 7.2 Convenções de Código
- **Formatação:** TODO (ex.: `gofmt` obrigatório, `prettier`, `black`, etc.)
- **Linting:** TODO (ex.: `golangci-lint`, `eslint`, `pylint`, etc.)
- **Nomenclatura:** TODO (ex.: seguir convenções da linguagem, padrão do projeto, etc.)
- **Estrutura:** Seguir estrutura de diretórios definida em `00-architecture.spec.md`

### 7.3 Build e Release
- **Builds reproduzíveis:** TODO (ex.: usar flags específicas, versionamento fixo, etc.)
- **Releases:** TODO (ex.: sempre via ferramenta X para garantir consistência)
- **Checksums:** TODO (ex.: obrigatórios para todos os artefatos)

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
- **TODO:** Listar testes de build (ex.: build local funciona, build cross-platform funciona, artefato gerado é válido, etc.)

### 9.2 Testes de Compatibilidade
- **TODO:** Listar testes de compatibilidade por plataforma (ex.: funciona em macOS x64, Linux x64, Windows, navegadores, etc.)

### 9.3 Testes de Qualidade
- **TODO:** Listar testes de qualidade (ex.: formatação ok, lint passa, testes passam com cobertura adequada, etc.)

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
- **TODO:** Requisitos para desenvolvimento local (ex.: linguagem X instalada, ferramentas Y configuradas, Git hooks configurados, etc.)

## 12. Abertos / Fora de Escopo

- **TODO:** Listar itens fora de escopo da stack técnica (ex.: plataformas não suportadas, ferramentas não usadas, etc.)

## Checklist Rápido (preencha antes de gerar código)
- [ ] Stack tecnológica definida e justificada
- [ ] Ferramentas de build e desenvolvimento especificadas
- [ ] Plataformas alvo definidas e testáveis
- [ ] Formato de empacotamento e distribuição definido
- [ ] Guardrails de dependências e build estabelecidos
- [ ] Critérios de aceite cobrem build, compatibilidade e qualidade

