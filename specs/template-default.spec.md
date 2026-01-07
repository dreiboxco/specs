# Template de Especificação (SDD)

Use este template para qualquer entrega (arquitetura, comandos, integrações). Substitua os blocos `TODO` por conteúdo concreto e mensurável. Evite termos vagos; prefira formatos, contratos e critérios testáveis.

**IMPORTANTE:** Este template é para specs de **funcionalidades** (comandos, features). **NÃO inclua detalhes de implementação técnica** como ferramentas de build, comandos de build específicos, ou bibliotecas de implementação. Detalhes técnicos devem estar apenas em `00-stack.spec.md`. Foque em **comportamento funcional**, **contratos** e **requisitos não funcionais** sem mencionar como será implementado.

## 1. Contexto e Objetivo
- **Contexto:** TODO (ex.: CLI specs precisa expor comandos para gestão de contas/projetos)
- **Objetivo:** TODO (o que muda para o usuário/negócio)
- **Escopo:** TODO (o que está dentro/fora nesta entrega)

## 2. Requisitos Funcionais
- TODO listar comportamentos observáveis. Ex.: comando, flags, entradas, saídas, mensagens, efeitos colaterais.
- Para cada requisito, torne-o testável (entrada → saída/efeito).

## 3. Contratos e Interfaces
- **CLI:** comando, subcomando, aliases, flags, args posicionais, variáveis de ambiente, códigos de saída, formato de output (ex.: tabela/JSON), mensagens de erro/sucesso.
- **APIs chamadas (se houver):** endpoint, método, payload, headers, auth, timeouts, códigos de status, esquemas de request/response.
- **Arquivos/SO:** caminhos tocados, permissões, formato esperado de config/state/cache, compatibilidade (Linux/macOS/Windows, x64/arm).

## 4. Fluxos e Estados
- Fluxo feliz (happy path) passo a passo.
- Estados alternativos: erros previsíveis (rede, auth, input inválido), reentrância/idempotência, retriability/backoff, comportamento offline (se aplicável).
- Mensagens exibidas por estado (curtas, orientadas a ação).

## 5. Dados
- Estruturas persistidas (config, cache, credenciais) com formato e localização.
- Políticas de retenção/expurgo. Campos sensíveis e como são protegidos.

## 6. NFRs (Não Funcionais)
- Desempenho (latência alvo por comando), limites (tamanho de payload, retries).
- Compatibilidade de plataforma (SO/arquitetura), versões mínimas (ex.: Node/go/python).
- Segurança: armazenamento seguro, transport layer (TLS), assinatura/verificação (checksums), princípio do menor privilégio.
- Observabilidade: logs mínimos, níveis, correlação (request-id), métricas/eventos.

## 7. Guardrails
- Restrições de stack/dependências (o que pode/não pode adicionar).
- Convenções de diretório/nome.
- Padrão de mensagens de help/erro.
- Política de feature flags/experimentos (se existir).

## 8. Critérios de Aceite
- Lista objetiva de verificações. Ex.: "`specs --help` exibe comandos core em <150ms", "`specs login --token X` retorna 0 e persiste credencial em <path>".
- Para cada critério, especifique como validar (teste automatizado/manual) e oráculo esperado.

## 9. Testes
- Tipos: unidade, integração (com mocks/sandbox), e2e local, contratos.
- Casos mínimos por requisito. Como isolar efeitos (fixtures, temp dirs).
- Como rodar: comandos, variáveis de ambiente necessárias, dependências externas (mocks/fakes).

## 10. Migração / Rollback
- Se altera estado local, definir passo de migração, verificação de sucesso e rollback seguro.
- Como detectar e recuperar de instalações corrompidas/versões quebradas.

## 11. Observações Operacionais
- Como distribuir: ex.: script `curl | sh`, pacote homebrew, binário estático.
- Auto-update: estratégia, verificação de integridade, fallback em caso de falha.

## 12. Abertos / Fora de Escopo
- Itens não cobertos nesta entrega (para evitar ambiguidade).

## Checklist Rápido (preencha antes de gerar código)
- [ ] Requisitos estão testáveis? Entradas/saídas precisas?
- [ ] Contratos de CLI/APIs têm formatos e códigos de saída definidos?
- [ ] Estados de erro e mensagens estão claros?
- [ ] Guardrails e convenções estão escritos?
- [ ] Critérios de aceite cobrem fluxos principais e erros?
- [ ] Migração/rollback definidos quando há mudança de estado?

