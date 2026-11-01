# 📋 SPEC_EP22: Telemetria em Tempo Real (Por Card e Header Global no Canvas)

> **Status:** Proposta para Execução  
> **Data:** 2026-11-01  
> **Épico Alvo:** ÉPICO 22 (Telemetria em Tempo Real no Looping Studio)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

O **ÉPICO 22** estabelece a visibilidade total da execução em tempo real dentro do Studio Canvas, sincronizada via Server-Sent Events (SSE) do motor Go (`/api/run`):

1. **Header Global de Telemetria no Canvas (`#canvas-telemetry-bar`):**
   * Barra flutuante ou fixa no topo da área do canvas (`.stage`), com visual sóbrio corporativo e 100% SVG (Zero Emojis).
   * Métricas consolidadas em tempo real:
     * **Tempo Total de Execução:** Duração decorrida ou finalizada (`Total: X.Xs`).
     * **Tokens de Entrada Totais ($\sum \text{Input}$):** Soma dos tokens de entrada de todos os nós.
     * **Tokens de Saída Totais ($\sum \text{Output}$):** Soma dos tokens gerados pelos modelos.
     * **Custo Consolidado:** Estimativa de custo agregado com base nos tokens consumidos pelos modelos.
     * **Status Global da Esteira:** Indicador visual de estado (`Pronto`, `Executando`, `Sucesso`, `Falha`).

2. **Telemetria Individual por Card de Nó (`.node__telemetry`):**
   * No rodapé do card de cada nó no canvas:
     * Duração individual em milissegundos/segundos (`120ms` ou `1.8s`).
     * Consumo de tokens pontual (`In: 1.4k / Out: 320`).
     * Status do nó e do gate com badges visuais nítidos (`passed`, `failed`, `running`).
     * Modelo efetivamente despachado ou executado.

3. **Eventos SSE Enriquecidos no Servidor Go (`internal/studio/run_api.go`):**
   * `dag_finish` transmite métricas consolidadas: tempo total em ms, soma de tokens e custo estimado em centavos/dólares.
   * `node_finish` transmite detalhes de tokens (input/output separados) e duração.

4. **Preservação Estrita da Regra Zero Emojis:**
   * Nenhum emoji em headers, cartões, logs ou tooltips da telemetria.

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

* [ ] **CA-22.1 (Header Global de Telemetria no Canvas):** O HTML e CSS do Studio definem a barra `#canvas-telemetry-bar` no topo do canvas com contadores de tempo total, tokens in/out e status da esteira.
* [ ] **CA-22.2 (Transmissão de Métricas Consolidadas via SSE):** O evento `dag_finish` do backend Go (`run_api.go`) envia `duration_ms`, `tokens.input`, `tokens.output`, `tokens.total` e `cost_est`.
* [ ] **CA-22.3 (Atualização em Tempo Real no Canvas):** O frontend (`app.js`) atualiza o header global e os cards individuais de nós em tempo real à medida que os eventos `node_start`, `node_finish` e `dag_finish` chegam.
* [ ] **CA-22.4 (Telemetria no Card do Nó):** Cada card no canvas renderiza sua seção `.node__telemetry` quando há dados de execução disponíveis.
* [ ] **CA-22.5 (Barreira Zero Emojis 100% Preservada):** Testes automatizados confirmam ausência total de emojis na camada de telemetria.
* [ ] **CA-22.6 (Zero Regressões & 100% Testes Verdes):** Toda a suíte de testes passa com sucesso (`go test ./...`) e o binário é recompilado.
