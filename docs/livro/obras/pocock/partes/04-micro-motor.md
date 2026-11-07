# PARTE IV — Micro: Roteamento, Contexto e Integração

## 14. Roteamento cognitivo de modelos

O manual recomenda um modelo primário e um alternativo local para cada sub-loop. A lógica é gastar raciocínio caro onde há ambiguidade e geração rápida onde há execução mecânica.

| Sub-loop          | Complexidade       | Primário    | Alternativo local  |
| :---------------- | :----------------- | :---------- | :----------------- |
| 1. Concepção      | Alta (socrática)   | `sonnet`    | `haiku`            |
| 2. Exploração     | Média              | `code-fast` | `qwen2.5-coder:7b` |
| 3. Especificação  | Alta (arquitetura) | `code-pro`  | `llama3.1:8b`      |
| 4. Construção TDD | Média a alta       | `sonnet`    | `qwen2.5-coder:7b` |
| 5. Diagnóstico    | Alta               | `sonnet`    | `llama3.1:8b`      |
| 6. Review e Retro | Média              | `code-fast` | `haiku`            |

O raciocínio socrático do sub-loop 1 pede um modelo capaz de desafiar premissas. Protótipos descartáveis do sub-loop 2 pedem velocidade, não profundidade. O review do sub-loop 6 é uma verificação binária de regras, em que um modelo rápido basta.

## 15. Higiene de contexto e fronteiras de fase

### As duas regras fixas

1. **Janela unificada de alinhamento.** Grilling, `to-spec` e `to-tickets` ficam na mesma sessão, porque compartilham raciocínio.
2. **Contexto limpo por ticket.** Cada `implement` começa em sessão limpa, consumindo apenas o arquivo do ticket e o repositório.

### A decisão na fronteira de fase

| Situação na fronteira                          | Ação       |
| :--------------------------------------------- | :--------- |
| Fase curta, contexto abaixo de 100 mil tokens  | Continuar  |
| Fase terminou, próximo ticket autocontido      | `/clear`   |
| Alinhamento perto do limite, antes dos tickets | `/compact` |
| Bifurcação lateral ou protótipo                | `/handoff` |

```{=typst}
#painel("Compactar não é limpar")[
  /compact resume a conversa e mantém o resumo; serve quando o raciocínio ainda é necessário, como no fim do alinhamento. /clear descarta tudo; serve quando o próximo passo é autocontido, como um ticket. Usar /compact antes de cada ticket carrega para a implementação resíduos de decisões que já estão no ticket.
]
```

## 16. Integração no LegoLoop

### O arquivo do loop

O arquivo `loops/matt-pocock.loop.json` segue o esquema `dag.v1`, chama-se "Pipeline de Rigor Matt Pocock (6 Sub-Loops)", versão 1.0.0, e tem gatilho manual. Os seis nós apontam para a peça `skills/aidd-pocock` e formam uma cadeia linear. Cada nó declara o artefato que precisa existir para ser concluído.

| Nó                           | Artefato exigido pelo gate    | Modelo sugerido          |
| :--------------------------- | :---------------------------- | :----------------------- |
| `subloop-1-alinhar-dominio`  | `GLOSSARY.md`                 | sonnet, ou haiku         |
| `subloop-2-especificar`      | `docs/specs/SPEC.md`          | não declarado            |
| `subloop-3-planejar-tickets` | `docs/plans/PLAN.md`          | não declarado            |
| `subloop-4-construir-tdd`    | `docs/reports/TDD_REPORT.md`  | sonnet, ou qwen2.5-coder |
| `subloop-5-review-auditar`   | `docs/reports/CODE_REVIEW.md` | não declarado            |
| `subloop-6-retro-ambiente`   | `docs/reports/RETRO.md`       | não declarado            |

O gate de cada nó é do tipo `file_exists(...)`. Ele confirma que o artefato foi produzido, não que está correto. O julgamento continua com os gates humanos G1 a G4.

### O template do Studio

Em `internal/studio/templates.go`, o template "Pipeline Matt Pocock (6 Sub-Loops)" aparece no menu Templates do Studio Canvas com o ícone `git`. Ele carrega o grafo com nós de alinhamento, especificação, tickets, TDD, review e retroalimentação.

### Comandos de verificação

```bash
legoloop validate loops/matt-pocock.loop.json   # contrato, gates e ciclos
legoloop run loops/matt-pocock.loop.json        # execução pelo motor DAG
legoloop studio                                 # Studio em 127.0.0.1:7878
```

Em 07/11/2026, o `validate` respondeu `6 nós válidos, sem ciclos`, em seis camadas de um nó cada, de `subloop-1-alinhar-dominio` a `subloop-6-retro-ambiente`.
