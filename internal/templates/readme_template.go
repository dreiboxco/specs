package templates

// ReadmeTemplate é o template básico para README.md
var ReadmeTemplate = []byte(`# Nome do Projeto

Descrição do projeto aqui.

## Estrutura

Este projeto segue a metodologia **SDD (Spec Driven Development)**.

- specs/: Especificações do projeto
- .cursorrules: Regras do Cursor para desenvolvimento

## Como Usar

1. Preencha os arquivos 00-*.spec.md em specs/ com informações do seu projeto
2. Use template-default.spec.md como base para criar novas specs
3. Execute specs validate para verificar se suas specs estão completas

## Especificações

Consulte as specs em specs/ para entender a arquitetura e funcionalidades do projeto.

`)
