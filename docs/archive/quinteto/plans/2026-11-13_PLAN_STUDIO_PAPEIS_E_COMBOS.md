# PLAN: Integração de Papéis Canônicos e Combos Honestos no LegoLoop Studio Web

## 1. Relação de Tickets

### Bloco 1: Backend API e Serialização
- **TICK-ST-01.1:** Atualizar struct `PieceDTO` em `internal/studio/pieces.go` com `PapelPrimario`, `PapeisSecundarios` e `Guilda`, integrando com `taxonomy.CatalogoCanonico`.
- **TICK-ST-01.2:** Adicionar testes em `internal/studio/pieces_test.go` para validar serialização e cálculo de guilda.
- **Stop-Gate Mecânico:** `SG-ST-E01` (Aprovação do backend).

### Bloco 2: Atualização de Templates do Studio
- **TICK-ST-02.1:** Migrar referências de combos em `internal/studio/templates.json` para os combos da I7.
- **TICK-ST-02.2:** Testar integridade de `templates.json` em `internal/studio/templates_test.go`.
- **Stop-Gate Mecânico:** `SG-ST-E02` (Aprovação dos templates).

### Bloco 3: Front-end (Palette, Canvas, Inspector e Estilos)
- **TICK-ST-03.1:** Injetar distintivos e filtros por guilda em `internal/studio/static/app.js`.
- **TICK-ST-03.2:** Injetar paleta de cores para as 7 guildas em `internal/studio/static/style.css`.
- **Stop-Gate Final:** `SG-ST-FINAL` (Testes globais e compilação do binário).
