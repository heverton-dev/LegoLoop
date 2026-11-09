# PARTE IV — TRANSVERSAL: TEMPLATES E COMPOSIÇÃO FRACTAL

Os templates do LegoLoop não são sugestões estéticas ou formulários vazios; eles são os **contratos fundacionais** que garantem que o sistema opere como uma máquina de estados determinística. Se uma ferramenta ou modelo não consegue interpretar o formato de um template, aquele arquivo deixa de ser engenharia e vira mera decoração.

---

## 1. A Cascata Canônica em Quatro Estágios

Toda demanda que entra no ecossistema passa compulsoriamente por quatro estágios de refinamento sucessivo antes que a primeira linha de código de produção seja autorizada:

```{=typst}
#esteira(
  no-claro("1. INTENT", sub: "O Porquê"),
  no-claro("2. SPEC", sub: "O Quê"),
  no-claro("3. PLAN", sub: "O Como"),
  no("4. ROUTE", sub: "A Máquina"),
)
```

| Estágio de Governança | Arquivo de Template    | Função Primária                          | Validação no Motor    |
| :-------------------- | :--------------------- | :--------------------------------------- | :-------------------- |
| **1. INTENT**             | `TEMPLATE_INTENT.md`      | Registra a dor do usuário, hipóteses e premissas | `Human-Gate 01`         |
| **2. SPEC**               | `TEMPLATE_SPEC.md`        | Define premissas, boundaries e critérios BDD `AC-XX` | `has_binary_criteria`     |
| **3. PLAN**               | `TEMPLATE_PLAN.md`        | Fatiamento vertical hierárquico (Épicos ➔ Nano) | `Human-Gate 02`         |
| **4. ROUTE**              | `TEMPLATE_ROUTE.json`     | Grafo acíclico executável em formato machine-readable | `dag.Parse` e Schema    |

```{=typst}
#frase-destaque[O PLAN para os olhos humanos é o espelho exato da ROUTE para os nós de máquina. Divergência entre os dois é defeito impeditivo.]
```

![A Tríade Canônica da Fabricação Autônoma](/img/diagrama_triade_canonico.svg)

```{=typst}
#painel("A Regra de Espelho Inegociável")[
  O PLAN (para humanos) e a ROUTE (para as máquinas) são duas faces da mesma moeda. Eles compartilham obrigatoriamente o mesmo conjunto de identificadores de nós, as mesmas dependências causais, os mesmos limites de timeout e o mesmo roteamento cognitivo. Qualquer divergência entre o PLAN e a ROUTE constitui defeito impeditivo de compilação.
]
```

---

## 2. A Composição Fractal via CLI: O Comando `compor`

A grande inovação arquitetural do LegoLoop é a capacidade de transformar qualquer processo isolado em uma etapa acoplável de um processo maior. 

Imagine que você fabricou um processo chamado `subloop-tdd` contendo o ciclo completo de Red-Green-Refactor. Quando você estiver desenhando um novo sistema corporativo, você não precisa redesenhar essas etapas:

```{=typst}
#ficha(
  ("Comando CLI", "legoloop compor <alvo> <subloop> [--pos N]"),
  ("Ação no DAG", "Insere nó do tipo 'subloop' com ref ao loop filho"),
  ("Prevenção de Ciclos", "cycle.DetectCycle(): recusa loops recursivos"),
  ("Isolamento de I/O", "Separação de caminhos via sub_artifact"),
  ("Versionamento", "Elevação semântica automática de versão (ex: 1.0.0 -> 1.1.0)"),
)
```

### Como a Composição Funciona nos Bastidores

1. **Validação Estrutural:** O comando verifica se ambos os loops existem no cofre (`~/.aidd/loops/` ou repositório) e confere as assinaturas de entrada e saída.
2. **Prevenção Mecânica de Ciclos:** O resolvedor de grafos simula a injeção. Se o loop `A` compuser o loop `B`, e o loop `B` já contiver referências ao loop `A`, o sistema aborta com erro fatal `ErrCiclo`.
3. **Mapeamento de Artefatos com `sub_artifact`:** Para impedir que o artefato final do sub-loop colida com o artefato do pipeline pai (um bug clássico em geradores de código), o motor grava o campo `sub_artifact`. O sub-loop opera em sua pasta própria (`out/<subloop>/`) e a Orchestrix entrega apenas o resultado final no destino do pai.

---

## 3. O Ciclo de Vida dos Cinco Verbos de Processo

O LegoLoop padronizou a gestão de processos no ecossistema através de cinco verbos de comando universais, disponíveis tanto na CLI quanto na interface visual do Studio:

| Verbo Operacional | Comando CLI                          | Papel no Ciclo de Vida                   |
| :---------------- | :----------------------------------- | :--------------------------------------- |
| **Fabricar**          | `legoloop fabricar "<nome>"`           | Cria um processo novo do zero via cascata socrática |
| **Compor**            | `legoloop compor "<alvo>" "<subloop>"` | Acopla um processo como sub-etapa de outro |
| **Listar**            | `legoloop listar`                      | Audita o inventário de processos, versões e vínculos |
| **Refinar**           | `legoloop refinar "<nome>"`            | Substitui atomicamente um nó defasado por um melhor |
| **Remover**           | `legoloop remover "<nome>"`            | Desvincula um processo com garantia de rollback seguro |

```{=typst}
#painel("Canonização e Hashing SHA-256")[
  Toda peça ou template importado ou fabricado passa imediatamente pelo crivo da canonização (legoloop canonize). O motor calcula a impressão digital SHA-256 de cada ativo, desduplicando cópias idênticas e vinculando a reputação daquela peça no cofre central do desenvolvedor.
]
```

## Rastreabilidade

* Implementação do comando compor: `internal/compor/compor.go` e testes em `internal/cli/compor_binary_test.go`.
* Validação de sub-artefatos sem colisão: `internal/dag/subartifact_test.go`.
* Contratos e templates SSOT: `docs/templates/templates.go`.
