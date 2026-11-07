# PARTE II — Meso: Encadeamento

## 4. Verbos como Elos de uma Corrente

Os verbos do LegoLoop não foram projetados para rodar isolados. Cada um é elo de uma corrente lógica: os diagnósticos alimentam ações cirúrgicas, e as ações desembocam em revalidação e entrega.

![Fluxo entre diagnóstico, ação e validação](docs/livro-verbos-bdd/img/diagrama_encadeamento.svg){width=95%}

| Origem | Destino | Finalidade da ligação |
| :--- | :--- | :--- |
| `listar` | `refinar` | Inspecionar os nós antes da troca |
| `listar` | `compor` | Identificar sub-loops reutilizáveis |
| `listar` | `remover` | Auditar resíduos e vínculos |
| `matriz` | `refinar` | Detectar peças defasadas |
| `matriz` | `fabricar` | Validar a aptidão das peças |
| `fabricar` / `compor` / `remover` | `listar` | Revalidar grafo, acoplamento e limpeza |
| `refinar` | Auditoria de gates | Validar o score após a troca |
| Auditoria de gates | `legopress` | Gerar a documentação final |

## 5. As Cinco Trilhas Canônicas

As trilhas abaixo esgotam as combinações lógicas da operação diária.

### Trilha 1: Criação e Publicação

**Objetivo:** criar um processo com governança completa sem duplicar trabalho.

**Encadeamento:** `listar` $\rightarrow$ `fabricar` $\rightarrow$ `matriz`

1. **`listar`:** verifica se já existe processo similar registrado no cofre.
2. **`fabricar`:** gera os quatro artefatos contratuais (`INTENT`, `SPEC`, `PLAN`, `ROUTE`).
3. **`matriz`:** confere se as peças selecionadas pelo motor têm boa reputação (Tier Ouro ou Prata).

### Trilha 2: Diagnóstico e Refinamento Cirúrgico

**Objetivo:** localizar etapas de baixo desempenho e substituir peças defasadas.

**Encadeamento:** `listar` $\rightarrow$ `matriz` $\rightarrow$ `refinar` $\rightarrow$ `matriz`

1. **`listar`:** localiza os identificadores exatos das etapas (`e01`, `e02`).
2. **`matriz`:** filtra peças Bronze ou defasadas e mapeia substitutas Tier Ouro ($\ge 0.85$).
3. **`refinar`:** troca atômica da peça na etapa, com backup em `.legoloop-bak` e bump minor de versão.
4. **`matriz`:** revalida o processo e confirma que a saúde geral do cofre subiu.

### Trilha 3: Composição Hierárquica e Reuso

**Objetivo:** reaproveitar um processo existente como sub-etapa de outro maior.

**Encadeamento:** `listar` $\rightarrow$ `compor` $\rightarrow$ `listar`

1. **`listar`:** seleciona o processo anfitrião (pai) e o sub-loop a ser acoplado (filho).
2. **`compor`:** conecta o sub-loop em posição específica (`legoloop compor <pai> <filho> --pos N`).
3. **`listar`:** valida a árvore hierárquica e os links do DAG na saída JSON.

### Trilha 4: Ciclo de Vida e Descomissionamento

**Objetivo:** desativar processos depreciados, limpando symlinks de todos os editores sem quebra.

**Encadeamento:** `matriz` $\rightarrow$ `listar` $\rightarrow$ `remover --dry-run` $\rightarrow$ `remover --force`

1. **`matriz`:** identifica processos com peças órfãs ou critérios descontinuados.
2. **`listar`:** analisa o raio de impacto (harnesses e editores vinculados).
3. **`remover --dry-run`:** simula a exclusão para confirmar que nenhum artefato crítico será afetado.
4. **`remover --force`:** executa o descarte com limpeza de symlinks e histórico de rollback preservado.

### Trilha 5: Governança Editorial e Consolidação Didática

**Objetivo:** transformar especificações, planos e rotas refinadas em livro técnico ou manual em PDF.

**Encadeamento:** `listar` $\rightarrow$ `matriz` $\rightarrow$ `legopress`

1. **`listar` e `matriz`:** auditam a estabilidade e a integridade dos nós do processo.
2. **`legopress`:** compila diffs de código, especificações BDD e diagramas em documento Typst/PDF canônico.

## 6. Tabela-Resumo dos Encadeamentos

| Trilha | Sequência de verbos | Gatilho típico | Resultado entregue |
| :--- | :--- | :--- | :--- |
| **1. Criação** | listar → fabricar → matriz | Nova demanda | INTENT, SPEC, PLAN e ROUTE |
| **2. Refinamento** | listar → matriz → refinar → matriz | Peça lenta ou defasada | Etapa trocada, bump minor |
| **3. Composição** | listar → compor → listar | Reuso de fluxo | Sub-loop isolado no DAG |
| **4. Descomissionamento** | matriz → listar → remover | Processo legado | Symlinks limpos, sem órfãos |
| **5. Fechamento editorial** | matriz → legopress | Fim de ciclo | Manual e PDF canônico |
