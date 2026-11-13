# INTENT: Integração de Papéis Canônicos e Combos Honestos no LegoLoop Studio Web

## 1. Contexto e Motivação
Com a conclusão da Iniciativa I5 (Matriz de Papéis e Casamento Determinístico) e da Iniciativa I7 (Combos Honestos do 9Router e Eliminação do Erro 413), a base de código do LegoLoop conta com 20 papéis canônicos distribuídos em 7 guildas, além de 6 combos honestos certificados com limites mecânicos estritos.
O LegoLoop Studio Web atualmente exibe peças e modelos sem informar o papel primário, os papéis secundários ou a guilda correspondente, e os modelos em `templates.json` ainda referenciam identificadores legados de combos.

## 2. Objetivos
1. Expor `papel_primario`, `papeis_secundarios` e `guilda` na API REST de peças (`/api/pieces`) do Studio.
2. Atualizar a paleta lateral e os nós do Canvas no Studio para exibir distintivos de Guilda e Papel com cores e estilos dedicados.
3. Atualizar o Inspetor lateral para detalhar os limites mecânicos do combo honesto (`max_tokens`, suporte a tools) quando o nó utilizar o 9Router.
4. Atualizar os modelos canônicos em `internal/studio/templates.json` para os combos certificados da Iniciativa I7.
5. Garantir conformidade com os portões mecânicos e cobertura total de testes unitários.

## 3. Não-Objetivos
- Modificar o algoritmo de casamento de peças (`internal/fabricar/merger_fit.go`).
- Adicionar dependências externas ao front-end ou ao `go.mod`.
- Alterar o protocolo de comunicação do `mcp-router`.
