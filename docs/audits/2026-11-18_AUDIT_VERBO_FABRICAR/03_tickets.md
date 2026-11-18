# 03 — Tickets de Remediação e Melhoria

> **Alvo:** cmd/legoloop e internal/cli (comando fabricar)
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Metodologia:** BDD / EARS / Zero Stubs

---

## Ticket TICK-FAB-01: Unificação do Despacho de Comandos na Raiz
- **Alvo:** cmd/legoloop/main.go e internal/cli/cli.go.
- **Objetivo:** Mover os desvios diretos de triagem e migrar-processos para o registro canônico da tabela de comandos em cli.Run.
- **Critério de Aceite (EARS):**
  - Dado que qualquer comando ou verbo do LegoLoop seja invocado no terminal;
  - Quando main() for executado;
  - Então todo o despacho deve fluir uniformemente através de cli.Run com contexto e dependências unificadas.
- **Gate Mecânico:** go test -count=1 ./cmd/legoloop/... ./internal/cli/... com exit code 0.

---

## Ticket TICK-FAB-02: Extensibilidade Dinâmica de Categorias no FlagSet
- **Alvo:** internal/cli/fabricar.go:55.
- **Objetivo:** Permitir que o parser da flag --category valide as categorias suportadas contra a lista canônica do pacote fabricar.Category, emitindo mensagens informativas claras se uma categoria personalizada for informada.
- **Critério de Aceite (EARS):**
  - Dado que o usuário informe --category video ou --category audio;
  - Quando legoloop fabricar for executado;
  - Então a categoria deve ser aceita e repassada ao motor sem erro de validação estática no CLI.
- **Gate Mecânico:** go test -v ./internal/cli -run TestRunFabricar passa com cobertura ampliada de categorias.

---

## Ticket TICK-FAB-03: Blindagem de Saída JSON Estrita em Todos os Caminhos
- **Alvo:** internal/cli/fabricar.go:170-177 e fabricar_report.go:15-25.
- **Objetivo:** Assegurar que sob a presença da flag --json, nenhuma string em texto puro seja emitida no os.Stdout, canalizando avisos estritamente para os.Stderr.
- **Critério de Aceite (EARS):**
  - Dado que o comando fabricar seja invocado com a flag --json;
  - Quando ocorrerem erros, avisos ou emissão de resultado;
  - Então o stdout deve conter unicamente um documento JSON válido decodificável por json.Unmarshal.
- **Gate Mecânico:** TestRunFabricar/json_imprime_o_Result_e_deixa_o_stdout_puro e testes adicionais de falha com --json validados com json.Valid.
