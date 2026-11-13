# SPEC: Integração de Papéis Canônicos e Combos Honestos no LegoLoop Studio Web

## 1. Escopo e Requisitos

### Requisito 1: Serialização de Papéis e Guilda na API (`internal/studio/pieces.go`)
- O tipo `PieceDTO` retornado por `/api/pieces` deve conter:
  - `papel_primario` (string)
  - `papeis_secundarios` ([]string)
  - `guilda` (string, derivada de `taxonomy.CatalogoCanonico[papel_primario].Guilda`)
- Caso o papel primário não esteja definido, a guilda deve ser uma string vazia ou omitida de forma segura sem quebra de contrato.

### Requisito 2: Apresentação no Front-end (`internal/studio/static/`)
- `static/app.js`:
  - Cartões de peças na paleta lateral devem exibir a tag da guilda e o papel primário.
  - Filtro interativo na paleta lateral para selecionar por guilda (`governanca`, `reversa`, `construcao`, `infra`, `audiovisual`, `qualidade`, `publicacao`).
  - Nós no Canvas devem exibir um distintivo com o papel canônico.
- `static/style.css`:
  - Estilos de cores acessíveis para cada uma das 7 guildas.

### Requisito 3: Atualização de Modelos do Studio (`internal/studio/templates.json`)
- Modelos existentes devem ter suas propriedades de roteamento migradas para os combos certificados da Iniciativa I7 (`t-codigo-128k`, `t-chat-128k`, `t-agente-128k`).
- `templates_test.go` deve validar a ausência de identificadores legados não mapeados.

## 2. Critérios de Aceitação Binários (EARS)
1. **[UBÍQUO]** Todo item retornado pela rota `/api/pieces` que possui papel canônico válido deve apresentar o campo `guilda` correspondente à sua guilda oficial.
2. **[ORIENTADO A EVENTO]** Quando o usuário filtrar a paleta por uma guilda específica, apenas peças pertencentes àquela guilda devem ser listadas.
3. **[ORIENTADO A ESTADO]** Enquanto um nó do Canvas estiver selecionado e seu modelo usar 9Router, o Inspetor deve exibir os metadados do combo honesto.
4. **[OPCIONAL]** Se uma peça não tiver papel definido no frontmatter, o Studio deve classificar como `geral` sem falhar.
