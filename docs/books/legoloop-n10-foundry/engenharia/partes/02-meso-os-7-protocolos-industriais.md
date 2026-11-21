# Parte II: Os 7 Protocolos Industriais de Missão Crítica

## Capítulo 2: Engenharia de Harness de Alta Confiabilidade

Para erradicar a fragilidade de pipelines gerados por LLM, o LegoLoop implementa transversalmente sete protocolos da aviação, medicina e indústria militar.

| Protocolo | Origem | Mecanismo Go | Propósito no LegoLoop |
|---|---|---|---|
| **STANAG 2014** | OTAN | `internal/cli` | Estrutura militar de 5 parágrafos (Situação, Missão, Execução, Logística, Sinal). |
| **ASD-STE100** | Aeroespacial | `internal/gate/ste100.go` | Frases $\le 20$ palavras, voz ativa e banimento de *hedging*. |
| **EARS** | Requisitos | `internal/specgen` | Templates formais de máquina de estados (`WHEN... SHALL`). |
| **FMEA** | Automotiva | `internal/gate` | Matriz de falhas pré-código; severidade $\ge 3$ exige Poka-Yoke. |
| **ICAO Readback** | Aviação Civil | `internal/cli/icao.go` | Cartão de confirmação antes de qualquer ação destrutiva. |
| **Toyota Jidoka** | Manufatura | `internal/runner/retry.go` | Limite de 2 retries; disjuntor Andon Cord puxado na 2ª falha. |
| **Protocolo SBAR** | Medicina Crítica | `internal/runner/sbar.go` | Relatório estruturado (`[S]`, `[B]`, `[A]`, `[R]`) em qualquer falha de nó. |

### O Disjuntor Andon Cord em Ação
O algoritmo `ExecuteWithAndon()` intercepta retries infinitos. Se um nó falha na primeira tentativa, é permitida uma segunda execução. Caso persista a falha, o disjuntor interrompe o pipeline e grava o diagnóstico SBAR, bloqueando alucinações repetidas.
