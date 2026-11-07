# PARTE I — NÍVEL MACRO: ONTOLOGIA E PROCESSO

## 1. O Conceito Formal de Processo na Engenharia de IA

Na ciência da computação clássica e na engenharia de sistemas de missão crítica, a definição de um **Processo** é rigorosa:

```{=typst}
#painel("Definição Axiomática de Processo")[
  Um processo é uma sequência finita, ordenada e governada de transições de estado, que transforma entradas definidas em saídas verificáveis, operando sob regras explícitas, com recursos previamente alocados e sob critérios determinísticos de validação e parada.
]
```

Quando a inteligência artificial generativa foi introduzida no desenvolvimento de software, a indústria cometeu o erro de tratá-la como um chat interativo contínuo. Em conversas longas e desestruturadas, os modelos sofrem de alucinação cumulativa, esquecimento de diretrizes anteriores e poluição descontrolada de contexto (*context bloat*). 

O LegoLoop nasce para refundar a engenharia de software com IA sobre bases fabris determinísticas. Em vez de prompts soltos, todo trabalho é encapsulado em uma **máquina de estados finita** representada por um Grafo Acíclico Dirigido (DAG).

```{=typst}
#frase-destaque[Um processo é uma máquina de estados finita: cada nó transforma entradas declaradas em saídas verificáveis, sem depender de memória mágica.]
```

![O Monólito Tradicional versus a Esteira Modular LegoLoop](docs/livro-mega-pipeline/img/diagrama_monolito_vs_legoloop.svg)

```{=typst}
#ficha(
  ("Natureza", "Máquina de Estados Finita e Auditável"),
  ("Grafo Base", "Directed Acyclic Graph (DAG v2) em JSON Schema"),
  ("Garantia de Parada", "Proibição estrita de ciclos infinitos (ErrCiclo)"),
  ("Handoff", "Exclusivamente mediado por sistema de arquivos (SSOT)"),
  ("Ambiente", "Workspaces temporários herméticos criados pelo Orchestrix"),
)
```

### Os Cinco Pilares de Governança do Processo

1. **Sequência Finita e Ordenada:** Cada processo possui um início formal demarcado por um gatilho (`trigger`) e um término inequívoco. A ordem de execução não é estocástica: o resolvedor topológico do motor garante que uma etapa só execute após todos os seus nós predecessores declarados em `depends_on` terem sido concluídos com sucesso.
2. **Transformação de Estados Auditável:** Um nó nunca "acha" o que deve fazer. Ele recebe uma lista fechada de arquivos de entrada (`inputs`), aplica uma transformação controlada e materializa um artefato no disco (`output.artifact`).
3. **Regras Explícitas e Ausência de Premissas Ocultas:** A IA é contratualmente proibida de inferir comportamento arquitetural em silêncio. Toda decisão deve estar ancorada nos documentos da tríade de planejamento.
4. **Recursos Rigidamente Alocados:** Cada bloco do processo declara com antecedência sua cota de recursos: qual modelo cognitivo será acionado, qual o teto de tokens de saída (`max_output_tokens`), qual o limite de tempo (`timeout_sec`) e quais ferramentas locais estão liberadas.
5. **Critérios Determinísticos de Validação:** A aprovação de uma etapa não depende da avaliação subjetiva do próprio modelo. Ela é julgada por expressões booleanas matemáticas puras compiladas em DSL (Domain Specific Language), tais como conferência de código de retorno (`exit_code == 0`) e existência física do arquivo (`file_exists`).

---

## 2. A Arquitetura Fractal do LegoLoop

Uma das propriedades mais poderosas da engenharia do LegoLoop é a sua **natureza fractal e homotética**. A mesma estrutura conceitual e contratual que rege a esteira inteira se repete em cada uma de suas subdivisões:

```{=typst}
#esteira(
  no("MEGA PIPELINE", sub: "Macro-Orquestrador"),
  no("SUB-LOOPS", sub: "Unidades Modulares"),
  no("LOOPS", sub: "DAGs Acíclicos"),
  no("NÓS / ETAPAS", sub: "Átomos Executáveis"),
)
```

* **Mega Pipeline (A Fábrica Completa):** Encadeia a visão global de ponta a ponta. Ele transforma uma ideia bruta ou solicitação humana em um produto de software compilado, testado em múltiplos dispositivos e documentado com rastreabilidade de sessão.
* **Sub-Loops (As Linhas de Montagem):** São especializações funcionais do processo (ex.: Concepção, Engenharia Reversa, Planejamento BDD, Construção TDD e Homologação). Cada sub-loop tem vida própria: pode ser invocado de forma autônoma ou encadeado como uma macro-peça.
* **Loops e DAGs (O Contrato Isomórfico):** No nível de dados e código, **não há diferença estrutural entre um loop simples e um mega pipeline**. Ambos são manifestos JSON aderentes ao mesmo esquema formal (`https://legoloop.dev/schemas/dag.v1.json`). O Mega Pipeline é apenas um loop cujos nós referenciam outros loops.
* **Nós e Etapas (`dag.Node`):** É a menor partícula operacional. Um nó nasce em um diretório temporário, recebe os artefatos de entrada, executa as ferramentas atribuídas, submete seu resultado ao crivo mecânico de um gate e encerra sua existência sem deixar resíduos de memória.

---

## 3. A Sextilha Canônica de Peças

Para que uma esteira seja flexível e componível, as etapas do processo utilizam um catálogo ontológico padronizado em seis naturezas de peças:

| Natureza da Peça | Formato Nativo                     | Função Primária no Fluxo                                | Motor de Execução |
| :--------------- | :--------------------------------- | :------------------------------------------------------ | :---------------- |
| **1. Skill**     | Markdown (`SKILL.md`)              | Prompt especialista com heurísticas e raciocínio de IA  | `Orchestrix`      |
| **2. Script**    | PowerShell / Shell / Go            | Transformação determinística de dados e arquivos locais | `Runtime Local`   |
| **3. Hook**      | Configuração JSON                  | Interceptação de ciclo de vida pré e pós-execução       | `Harness Hooks`   |
| **4. MCP**       | Configuração JSON (`.mcp.json`)    | Ferramental vivo externo (navegador, banco de dados)    | `Protocolo MCP`   |
| **5. Gate**      | Expressão booleana em DSL          | Validação mecânica binária e aprovação de transição     | `internal/gate`   |
| **6. Loop**      | Manifesto JSON (`.loop.json`)      | Encadeamento de nós em DAG para reúso modular           | `dag.Parse`       |

```{=typst}
#painel("Equilíbrio Cognitivo vs. Determinístico")[
  O segredo de uma esteira robusta está no casamento entre peças cognitivas (Skills) e peças determinísticas (Scripts e Gates). A inteligência artificial só deve ser chamada para raciocínio, escrita e síntese; tudo o que puder ser medido por um compilador, formatador ou script local deve ser executado fora da LLM para economizar tokens e garantir tolerância zero a falhas.
]
```

## Rastreabilidade

* Código do grafo e contrato de nós: `internal/dag/types.go` e `internal/dag/parse.go`.
* Validação de ciclos infinitos: `internal/dag/cycle.go`.
* Definições de templates canônicos: `docs/templates/README.md` e `docs/templates/TEMPLATE_PIECE_LOOP.json`.
