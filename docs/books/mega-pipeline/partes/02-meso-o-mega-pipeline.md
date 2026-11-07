# PARTE II — NÍVEL MESO: O MEGA PIPELINE E A ORCHESTRIX

## 1. Visão Panorâmica dos Cinco Sub-Loops

O **Mega Pipeline** organiza o ciclo de vida completo de engenharia de software em cinco sub-loops lineares e complementares. Nenhum sub-loop avança para o próximo sem que seu artefato terminal tenha sido aprovado por um portão de qualidade:

```{=typst}
#esteira(
  no-claro("LOOP 1", sub: "Concepção & Intenção"),
  no-claro("LOOP 2", sub: "Design & Engenharia Reversa"),
  no-claro("LOOP 3", sub: "Tríade de Planejamento"),
  no("LOOP 4", sub: "Construção TDD por Épico"),
)
#v(-4pt)
#esteira(
  no("LOOP 5", sub: "Homologação E2E & Release"),
)
```

1. **Loop 1 — Concepção & Alinhamento:** Transforma o pedido do usuário em uma declaração clara de intenção (`INTENT.md`) através de uma entrevista socrática.
2. **Loop 2 — Engenharia Reversa & Design Pass:** Inspeciona referências visuais ou aplicações legadas, extraindo tokens limpos e wireframes técnicos (`DESIGN.md`).
3. **Loop 3 — Especificação & Tríade de Planejamento:** Formaliza a engenharia em critérios BDD executáveis (`SPEC.md`), quebra o trabalho em tickets (`PLAN.md`) e compila a rota acíclica (`ROUTE.json`).
4. **Loop 4 — Construção TDD por Épico:** Implementa o código em ciclo estrito de testes unitários (Red ➔ Green ➔ Refactor), módulo a módulo, sem tolerância a stubs.
5. **Loop 5 — Homologação E2E & Fechamento:** Audita o sistema com navegadores reais headless, recompila os binários de produção, consolida métricas de tokens e gera o relatório final.

```{=typst}
#frase-destaque[Nenhum bloco herda o contexto do anterior. A Orchestrix cria e destrói diretórios estéreis em milissegundos, garantindo Zero Context Bloat.]
```

![A Esteira Mestre dos Cinco Sub-Loops do Mega Pipeline](/img/diagrama_mega_pipeline_cinco_subloops.svg)

---

## 2. A Orchestrix Cirúrgica e o Handoff File-Based

O maior ofensor de custos e alucinações em sistemas agênticos convencionais é o compartilhamento indiscriminado de histórico de conversação entre tarefas. Quando um agente recebe centenas de mensagens acumuladas de etapas passadas, o ruído semântico destrói a atenção do modelo.

O motor **Orchestrix** do LegoLoop resolve esse problema impondo o isolamento cirúrgico por bloco:

```{=typst}
#ficha(
  ("Mecanismo", "Diretório temporário e isolamento hermético"),
  ("Contrato de Entrada", "inputs: declaração estrita de arquivos consumidos"),
  ("Função de Transferência", "transferInputs(): copia apenas o necessário"),
  ("Contrato de Saída", "output.artifact: gravação no diretório pai"),
  ("Limpeza de Memória", "defer os.RemoveAll(): zero lixo residual no SO"),
)
```

### O Ciclo de Execução Estéril de um Sub-Loop

Quando a Orchestrix é instruída a executar um nó do tipo `"type": "subloop"`, o código em `internal/orchestrix/subloop.go` executa o seguinte protocolo de seis passos:

1. **Geração do Workspace Estéril:** O runtime chama `os.MkdirTemp("", "legoloop-subloop-*")`, instanciando uma pasta limpa no disco.
2. **Transferência Cirúrgica de Entradas:** A função `transferInputs` inspeciona a lista de arquivos declarados em `inputs` e copia apenas esses arquivos para dentro da pasta temporária. Nenhum arquivo estranho é incluído.
3. **Execução Isolada do Sub-Grafo:** O sub-grafo interno é parsed e executado pelo `runner.Run`. As etapas internas consomem e produzem dados exclusivamente dentro daquela bolha.
4. **Avaliação dos Gates Mecânicos:** O gate de saída de cada etapa é testado localmente. Se houver falha, o processo é imediatamente interrompido.
5. **Extração do Artefato Final:** A função `finalizeSubloop` copia o arquivo terminal contratado (`sub_artifact`) de volta para o diretório do projeto principal, preservando seu hash e integridade.
6. **Destruição do Ambiente:** O diretório temporário é totalmente apagado via `os.RemoveAll()`, garantindo que resíduos de compilação ou arquivos intermediários não contaminem as etapas futuras.

---

## 3. Roteamento Cognitivo e Effort Clamping

Nem toda tarefa exige o mesmo poder de fogo computacional. Enviar a correção de um erro de lint para um modelo de raciocínio profundo de última geração é um desperdício financeiro inaceitável; da mesma forma, confiar a concepção da arquitetura de um sistema complexo a um modelo leve leva a colapsos de lógica.

O Mega Pipeline introduz a matriz de **Roteamento Cognitivo com Clamping de Esforço**:

| Nível de Tarefa   | Escopo e Complexidade                       | Modelo Cognitivo Sugerido               | Flag de Esforço |
| :---------------- | :------------------------------------------ | :-------------------------------------- | :-------------- |
| **Épico (Macro)** | Arquitetura, contratos BDD, DAG e modelagem | `Claude 3.7 Sonnet`, `Gemini 2.0 Pro`   | `--effort high` |
| **Micro-ticket**  | Módulo funcional, suíte de testes unitários | `Claude 3.5 Sonnet`, `Gemini 2.0 Flash` | `--effort med`  |
| **Nano-ticket**   | Correção de lint, tipo TypeScript, assert   | `Gemini Flash-Lite`, `Ollama Local 8B`  | `--effort low`  |

```{=typst}
#painel("Regra de Clamping de Tokens")[
  A flag `--effort` define o teto de tokens de raciocínio interno (thinking budget) alocados pela API. Em nano-tickets, o modelo é instruído a responder instantaneamente sem divagações, preservando a cota semanal para o que realmente importa: o design arquitetural e os contratos de teste.
]
```

---

## 4. As Sete Leis Invioláveis da Esteira

O Mega Pipeline opera debaixo de sete mandamentos mandatórios que nenhum agente ou executor tem autorização para violar:

1. **Zero Context Bloat:** Nenhum bloco herda o histórico conversacional de blocos anteriores.
2. **Stop-Gates de Auditoria Humana:** Paradas obrigatórias entre concepção, planejamento e ao término de cada Épico de código.
3. **A Tríade Inegociável:** O código só pode começar se existirem no disco a `SPEC.md`, o `PLAN.md` e a rota executável `ROUTE.json`.
4. **Provider-Awareness:** Roteamento de modelos adaptado ao ambiente local do desenvolvedor (`9Router`, `FreeLLMAPI` ou chave direta).
5. **TDD Estrito & Zero Stubs:** O teste unitário deve falhar comprovadamente antes do código de produção. Proibição absoluta de funções vazias com `TODO` ou `pass`.
6. **Padronização ISO de Nomenclatura:** Todos os artefatos de documentação seguem o prefixo `YYYY-MM-DD_<NOME>.<ext>`.
7. **Executor Silencioso:** Respostas concisas no chat e pensamento interno telegráfico no padrão Caveman.

## Rastreabilidade

* Implementação do isolamento de sub-loops: `internal/orchestrix/subloop.go`.
* Runner determinístico do DAG: `internal/runner/runner.go`.
* Regras de conduta de agentes: `AGENTS.md` e `docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`.
