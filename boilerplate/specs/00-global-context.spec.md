# 00 - Especificação de Contexto Global

Esta especificação define o contexto global do projeto: visão, objetivos, escopo, requisitos não funcionais, estratégias de distribuição, configuração, integrações e testes. Use-a como referência para entender o projeto como um todo antes de implementar features específicas.

## 1. Visão e Objetivos

### 1.1 Propósito do Projeto
- **TODO:** Descrever o propósito do projeto (ex.: CLI para desenvolvedores, API REST para integração, aplicação web para usuários finais, etc.)
- **Inspiração/Referência:** TODO (ex.: inspirado em Heroku CLI, AWS CLI, Stripe API, etc.)

### 1.2 Usuário-alvo
- **TODO:** Definir quem são os usuários principais (ex.: desenvolvedores, usuários finais, sistemas externos, pipelines CI/CD, etc.)
- **Casos de uso primários:** TODO (ex.: automação, integração, gerenciamento, etc.)

### 1.3 Resultados Esperados
- **TODO:** Listar resultados mensuráveis esperados (ex.: onboarding em < 2 minutos, latência < 100ms, etc.)
- **Métricas de sucesso:** TODO (ex.: tempo de resposta, taxa de erro, satisfação do usuário, etc.)

## 2. Escopo

### 2.1 Escopo Inicial (v1)
- **TODO:** Listar funcionalidades/core que estarão na primeira versão
- **Entregas mínimas:** TODO (ex.: autenticação, CRUD básico, comandos core, etc.)

### 2.2 Fora de Escopo (v1)
- **TODO:** Listar explicitamente o que NÃO estará na v1 (ex.: plugins dinâmicos, OAuth interativo, features avançadas, etc.)
- **Justificativa:** TODO (ex.: complexidade, priorização, dependências externas, etc.)

### 2.3 Roadmap Futuro
- **TODO:** Itens planejados para versões futuras (se aplicável)
- **Priorização:** TODO (ex.: v2, v3, etc.)

## 3. Requisitos Não Funcionais Globais

### 3.1 Desempenho
- **Latência:** TODO (ex.: resposta < 150ms, tempo de build < 10s, etc.)
- **Throughput:** TODO (ex.: requisições por segundo, operações simultâneas, etc.)
- **Timeouts:** TODO (ex.: rede 30s, operações 60s, etc.)

### 3.2 Robustez
- **Idempotência:** TODO (ex.: operações devem ser idempotentes quando fizer sentido)
- **Retries:** TODO (ex.: backoff exponencial 1s, 2s, 4s para rede)
- **Tolerância a falhas:** TODO (ex.: degradação graciosa, fallbacks, etc.)

### 3.3 Observabilidade
- **Logs:** TODO (ex.: níveis error/warn/info/debug, formato estruturado, request-id)
- **Métricas:** TODO (ex.: latência, taxa de erro, throughput, etc.)
- **Rastreamento:** TODO (ex.: correlation IDs, distributed tracing, etc.)

### 3.4 Portabilidade
- **Plataformas:** TODO (ex.: macOS, Linux, Windows, navegadores, etc.)
- **Arquiteturas:** TODO (ex.: x64, arm64, etc.)
- **Dependências:** TODO (ex.: binário único, runtime mínimo, etc.)

### 3.5 Segurança
- **Autenticação:** TODO (ex.: tokens, OAuth, API keys, etc.)
- **Autorização:** TODO (ex.: RBAC, permissões, etc.)
- **Armazenamento seguro:** TODO (ex.: keychain, secret service, criptografia, etc.)
- **Transporte:** TODO (ex.: HTTPS obrigatório, TLS mínimo, etc.)

### 3.6 Recuperação
- **Backup:** TODO (ex.: estratégia de backup, frequência, etc.)
- **Rollback:** TODO (ex.: estratégia de rollback, validação, etc.)
- **Integridade:** TODO (ex.: checksums, assinaturas, validação, etc.)

## 4. Distribuição e Instalação

### 4.1 Estratégia de Distribuição
- **Formato:** TODO (ex.: binário, pacote npm, imagem Docker, etc.)
- **Canais:** TODO (ex.: GitHub Releases, npm registry, Docker Hub, etc.)
- **Instalador:** TODO (ex.: `curl | sh`, `npm install`, `docker pull`, etc.)

### 4.2 Verificações Pós-Instalação
- **TODO:** Listar verificações necessárias após instalação (ex.: versão ok, permissões, PATH, etc.)

### 4.3 Desinstalação
- **TODO:** Estratégia de desinstalação (ex.: comando `uninstall`, script, remoção manual, etc.)

## 5. Auto-Update (se aplicável)

### 5.1 Fonte de Versões
- **TODO:** De onde vêm as versões (ex.: GitHub Releases, registry, API, etc.)
- **Canais:** TODO (ex.: stable, beta, alpha, etc.)

### 5.2 Estratégia
- **TODO:** Como funciona o update (ex.: in-place, download + substituição, etc.)
- **Validação:** TODO (ex.: checksum, assinatura, validação pós-update, etc.)
- **Rollback:** TODO (ex.: backup automático, restauração em caso de falha, etc.)

### 5.3 Segurança
- **TODO:** Medidas de segurança (ex.: checksum SHA256, HTTPS, assinatura, etc.)

## 6. Configuração, Estado e Cache

### 6.1 Arquivos e Localização
- **Config:** TODO (ex.: `~/.projeto/config.json`, XDG `~/.config/projeto/`, etc.)
- **Estado:** TODO (ex.: banco de dados, arquivos de estado, etc.)
- **Cache:** TODO (ex.: `~/.projeto/cache/`, XDG `~/.cache/projeto/`, etc.)
- **Lock files:** TODO (ex.: para evitar concorrência, etc.)

### 6.2 Formato
- **TODO:** Formato de configuração (ex.: JSON, YAML, TOML, etc.)
- **Schema:** TODO (ex.: onde está definido o schema, validação, etc.)

### 6.3 Campos Sensíveis
- **TODO:** Como são protegidos (ex.: nunca em config, armazenamento seguro, permissões 600, etc.)
- **Rotação:** TODO (ex.: como rotacionar credenciais, etc.)

### 6.4 Cache
- **TODO:** Estratégia de cache (ex.: TTL padrão, invalidação, limpeza, etc.)

## 7. Integrações Externas

### 7.1 APIs Externas
- **Base URL:** TODO (ex.: `https://api.exemplo.com`, configurável via env/config)
- **Autenticação:** TODO (ex.: header `Authorization: Bearer <token>`, API key, etc.)
- **Timeouts:** TODO (ex.: 30s padrão, retries 3x, etc.)
- **Versão de API:** TODO (ex.: header `X-API-Version: v1`, etc.)

### 7.2 Proxies/Corporate
- **TODO:** Suporte a proxies (ex.: `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, etc.)

### 7.3 Outras Integrações
- **TODO:** Outras integrações necessárias (ex.: banco de dados, serviços de terceiros, etc.)

## 8. Estratégia de Testes

### 8.1 Tipos de Testes
- **Unidade:** TODO (ex.: com mocks, interfaces testáveis, cobertura alvo 80%, etc.)
- **Integração:** TODO (ex.: serviços reais/sandbox, fluxos completos, etc.)
- **E2E:** TODO (ex.: ambiente completo, containers, etc.)
- **Contratos:** TODO (ex.: schemas, validação request/response, etc.)

### 8.2 Ferramentas
- **TODO:** Ferramentas de teste (ex.: framework de testes, mocks, fixtures, etc.)

### 8.3 Git Hooks (se aplicável)
- **TODO:** Hooks configurados (ex.: pre-push executa testes, pre-commit executa lint, etc.)

## 9. Convenções e Guardrails Globais

### 9.1 Mensagens e Comunicação
- **TODO:** Padrão de mensagens (ex.: curtas e acionáveis, help sempre disponível, etc.)
- **Output:** TODO (ex.: texto amigável, `--json` para output estruturado, etc.)
- **Códigos de saída:** TODO (ex.: 0 sucesso, 1 erro genérico, 2 input inválido, etc.)

### 9.2 Logs
- **TODO:** Padrão de logging (ex.: stdout para output, stderr para erros, níveis, etc.)
- **Nunca logar:** TODO (ex.: segredos, tokens, senhas, etc.)

### 9.3 Dependências
- **TODO:** Política de dependências (ex.: evitar libs pesadas, justificar cada nova lib, preferir stdlib, etc.)

### 9.4 Segurança
- **TODO:** Guardrails de segurança (ex.: nunca logar segredos, armazenamento seguro, etc.)

## 10. Riscos e Decisões em Aberto

### 10.1 Riscos Identificados
- **TODO:** Listar riscos conhecidos (ex.: tamanho do binário, compatibilidade, dependências nativas, etc.)

### 10.2 Decisões em Aberto
- **TODO:** Decisões que ainda precisam ser tomadas (ex.: telemetria, feature flags, etc.)

## 11. Referências a Outras Specs

### 11.1 Arquitetura
- **Referência:** Detalhes de padrão arquitetural, estrutura de diretórios e isolamento estão em `00-architecture.spec.md`.

### 11.2 Stack Técnica
- **Referência:** Detalhes de linguagem, ferramentas, build e empacotamento estão em `00-stack.spec.md`.

## Critérios de Aceite (Contexto Global)

- [ ] Visão e objetivos definidos e claros
- [ ] Escopo inicial e fora de escopo explicitamente listados
- [ ] Requisitos não funcionais globais mensuráveis e testáveis
- [ ] Estratégia de distribuição e instalação definida
- [ ] Estratégia de configuração, estado e cache definida
- [ ] Integrações externas documentadas
- [ ] Estratégia de testes por camada definida
- [ ] Convenções e guardrails globais estabelecidos
- [ ] Riscos e decisões em aberto identificados

## Checklist Rápido (preencha antes de gerar código)

- [ ] Visão e objetivos estão claros e mensuráveis?
- [ ] Escopo inicial e fora de escopo estão explicitamente definidos?
- [ ] Requisitos não funcionais são testáveis e mensuráveis?
- [ ] Estratégias de distribuição, configuração e integração estão definidas?
- [ ] Convenções e guardrails globais estão escritos?
- [ ] Riscos e decisões em aberto estão identificados?

