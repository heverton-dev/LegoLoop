# PARTE IV — Micro: Números, Eixos, Playbooks e Integração

## 15. Os limiares de Core Web Vitals

O arquivo `references/performance-checklist.md`, no upstream da skill `aidd-osmani`, fixa as faixas usadas para classificar as três métricas de experiência do usuário.

| Métrica                     | Bom      | Precisa melhorar | Ruim     |
| :-------------------------- | :------- | :--------------- | :------- |
| LCP, maior elemento visível | ≤ 2,5 s  | ≤ 4,0 s          | > 4,0 s  |
| INP, resposta à interação   | ≤ 200 ms | ≤ 500 ms         | > 500 ms |
| CLS, deslocamento acumulado | ≤ 0,1    | ≤ 0,25           | > 0,25   |

O mesmo checklist liga cada métrica a ações concretas. Imagens e elementos `<source>` com largura e altura explícitas evitam CLS. A imagem principal da página usa `fetchpriority="high"` e não é carregada sob demanda, para não atrasar o LCP. Tarefas longas, acima de 50 ms, são quebradas para liberar a thread principal, que o checklist aponta como a principal alavanca do INP.

## 16. Os cinco eixos de revisão

O upstream da habilidade `code-review-and-quality` define cinco eixos e afirma que toda mudança é revisada antes do merge, sem exceções.

| Eixo                        | Pergunta central                                      |
| :-------------------------- | :---------------------------------------------------- |
| Correção                    | O código faz o que a spec pede, inclusive nas bordas? |
| Legibilidade e simplicidade | Outra pessoa entende sem explicação?                  |
| Arquitetura                 | A mudança respeita os limites dos módulos?            |
| Segurança                   | Introduz alguma vulnerabilidade?                      |
| Performance                 | Introduz algum problema de desempenho?                |

Os eixos de segurança e performance remetem às habilidades especializadas, `security-and-hardening` e `performance-optimization`, para a análise detalhada.

## 17. Roteamento cognitivo e esforço

| Sub-loop          | Papel cognitivo           | Modelo             | Esforço       |
| :---------------- | :------------------------ | :----------------- | :------------ |
| Alinhamento       | Raciocínio socrático      | Sonnet ou Code-Pro | Médio         |
| Spec e restrições | Síntese de contratos      | Sonnet             | Alto          |
| Tickets           | Planejamento de DAG       | Code-Pro           | Médio         |
| Construção        | Engenharia TDD            | Sonnet             | Alto          |
| Diagnóstico       | Isolamento de falhas      | Sonnet             | Alto          |
| Performance       | Telemetria e profiling    | Code-Pro           | Médio         |
| Review            | Auditoria adversarial     | Sonnet             | Alto          |
| Entrega           | Governança e documentação | Code-Pro           | Baixo a médio |

## 18. Higiene de contexto e retomada

1. **Janela de alinhamento.** Intenção, restrições, spec e plano acontecem em uma janela de contexto progressiva.
2. **Contexto estéril por ticket.** Cada ticket é implementado apenas com a spec e as fronteiras do ticket, sem rascunhos anteriores.
3. **Arquivo de estado.** Tarefas longas mantêm `docs/osmani/<alvo>/ESTADO.md` com o estágio atual, os gates satisfeitos com evidência (`arquivo:linha` ou log) e os tickets pendentes. Se o arquivo já existe, a skill o lê, informa o estágio e continua dali.

## 19. Os playbooks da skill

A skill `aidd-osmani` organiza as habilidades em doze playbooks. Em cada um, o agente para em cada gate e resume em uma linha a passagem entre fases.

| Playbook | Nome           | Uso típico                                          |
| :------- | :------------- | :-------------------------------------------------- |
| A        | Auditar        | Varredura somente leitura por eixo, com contestação |
| B        | Corrigir       | Bug com reprodução, Prove-It e correção mínima      |
| C        | Construir      | Funcionalidade da entrevista à revisão              |
| D        | Simplificar    | Menos código com comportamento idêntico             |
| E        | Otimizar       | Medir, atacar o maior gargalo, medir de novo        |
| F        | Proteger       | Segurança com teste negativo por correção           |
| G        | Entregar       | Review final, observabilidade, ADR e ship           |
| H        | Migrar         | Expandir e contrair com testes e rollback           |
| I        | Definir        | Entrevista até spec ou decisão de não fazer         |
| J        | Documentar     | ADRs e contratos                                    |
| K        | Planejar       | Tickets verticais com critério binário              |
| L        | Ciclo completo | Encadeia os demais em cinco estágios, P1 a P4       |

O Playbook L é comparado ao loop mestre no capítulo 6.

## 20. Integração no LegoLoop

O arquivo `loops/addy-osmani.loop.json` segue o esquema `dag.v1`, tem gatilho manual e declara sete nós em cadeia linear, todos apontando para `skills/aidd-osmani`.

| Nó                                 | Artefato exigido pelo gate    | Modelo sugerido  |
| :--------------------------------- | :---------------------------- | :--------------- |
| `subloop-1-alinhar-contexto`       | `docs/intents/INTENT.md`      | sonnet, ou haiku |
| `subloop-2-especificar-restricoes` | `docs/specs/SPEC.md`          | não declarado    |
| `subloop-3-planejar-tickets`       | `docs/plans/PLAN.md`          | não declarado    |
| `subloop-4-construir-tdd`          | `docs/reports/TDD_REPORT.md`  | sonnet, ou qwen  |
| `subloop-5-runtime-webperf`        | `docs/reports/PERF_REPORT.md` | não declarado    |
| `subloop-6-review-hardening`       | `docs/reports/CODE_REVIEW.md` | sonnet, ou haiku |
| `subloop-7-observabilidade-ship`   | `docs/reports/SHIP_REPORT.md` | não declarado    |

Cada gate é do tipo `file_exists(...)`: confirma que o artefato existe, não que está correto. No Studio, o template se chama "Addy Osmani Web Performance & Rigor Loop", em `internal/studio/templates.go`.

Os comandos reais do CLI para operar a esteira são:

```bash
legoloop validate loops/addy-osmani.loop.json
legoloop run loops/addy-osmani.loop.json --var target_route=/checkout
legoloop studio        # padrão: 127.0.0.1:7878
```

Em 07/11/2026, o `validate` respondeu `7 nós válidos, sem ciclos`, em sete camadas de um nó cada.
