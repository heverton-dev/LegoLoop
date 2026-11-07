# PARTE III — Micro: Cenários BDD

## 7. Como Ler os Cenários

Cada cenário demonstra **o mesmo comportamento esperado** acionado pelas três vias. A seção *Quando* lista Terminal, Chat do Agente e Studio; as seções *Então* e *E* valem igualmente para as três.

## 8. Cenário 1: Fabricar um Pipeline do Zero

> **Objetivo:** o sistema cria automaticamente toda a estrutura de governança para uma nova necessidade de automação.

**Dado que** o desenvolvedor precisa de um pipeline para auditar a segurança de contratos inteligentes em Solidity,

**Quando** aciona a fabricação:

*Via Terminal (CLI):*

```bash
legoloop fabricar "auditoria-smart-contracts" \
  --goal "Auditar vulnerabilidades reentrancy e overflow em Solidity" \
  --category build --yes
```

*Via Chat do Agente:*

```markdown
/legoloop-fabricar
Crie um novo processo chamado 'auditoria-smart-contracts' focado em auditar
vulnerabilidades de reentrancy e overflow em Solidity, na categoria build.
```

*Via LegoLoop Studio:*

1. Na aba **Processos**, clica em **+ Novo Processo**.
2. No modal, preenche **Nome** `auditoria-smart-contracts`, **Categoria** `build` e **Objetivo**.
3. Clica em **Fabricar Processo**.

**Então** o motor hermético produz quatro artefatos canônicos, nomeados no padrão `YYYY-MM-DD_<TIPO>_AUDITORIA_SMART_CONTRACTS`:

| Artefato | Pasta | Extensão | Conteúdo |
| :--- | :--- | :--- | :--- |
| INTENT | `docs/intents/` | `.md` | Problema e motivação |
| SPEC | `docs/specs/` | `.md` | Critérios de aceitação binários |
| PLAN | `docs/plans/` | `.md` | Tickets em fatias verticais |
| ROUTE | `docs/routes/` | `.json` | DAG formal executável |

**E** o pipeline fica pronto para execução sem código boilerplate manual.

## 9. Cenário 2: Inspecionar o Inventário com `listar`

> **Objetivo:** saber quais processos existem no cofre, suas versões e a quais harnesses estão vinculados.

**Dado que** existem vários loops criados na máquina,

**Quando** o usuário solicita a listagem:

*Via Terminal (CLI):*

```bash
legoloop listar
legoloop listar --json   # saída estruturada para pipelines
```

*Via Chat do Agente:*

```markdown
/legoloop-listar
Liste todos os processos ativos no cofre do LegoLoop e verifique se há
harnesses descalibrados ou vínculos órfãos.
```

*Via LegoLoop Studio:*

1. Acessa a seção **Processos** no menu lateral.
2. Visualiza a grade de cartões com selo de versão (`@1.0.0`), categoria e estado de cada harness vinculado.

**Então** o sistema exibe nome e versão do processo, categoria (`build`, `reverse`, `web`, `mobile`), harnesses vinculados com estado dos links simbólicos e eventuais vínculos órfãos ou advertências de sincronização.

**E** o desenvolvedor obtém visibilidade completa do ecossistema em qualquer uma das vias.

## 10. Cenário 3: Avaliar a Saúde das Peças com `matriz`

> **Objetivo:** garantir que o pipeline usa apenas peças modernas e de alta reputação antes de uma entrega crítica.

**Dado que** o desenvolvedor quer auditar a integridade das peças do ecossistema,

**Quando** consulta a matriz:

*Via Terminal (CLI):*

```bash
legoloop matriz --tier ouro
```

*Via Chat do Agente:*

```markdown
/legoloop-matriz
Consulte a matriz de peças do LegoLoop e mostre apenas as ferramentas
de Tier Ouro com score >= 0.85 utilizáveis no pipeline.
```

*Via LegoLoop Studio:*

1. Clica em **Matriz de Peças** na navegação principal.
2. Seleciona o chip **Padrão Ouro** no filtro de topo.
3. Inspeciona a tabela com scores de reputação, cobertura de testes de mutação e conformidade com zero stubs.

**Então** o sistema exibe apenas componentes certificados: score $\ge 0.85$, sem stubs e com testes herméticos.

**E** sinaliza as peças Bronze ou defasadas, recomendando substituição cirúrgica.

## 11. Cenário 4: Compor um Sub-Loop em Pipeline Existente

> **Objetivo:** anexar um sub-loop de homologação visual já existente a um pipeline de construção.

**Dado que** existem o pipeline mestre `pipeline-frontend` e o sub-loop `aidd-homologue@1.0.0`,

**Quando** o usuário conecta o sub-loop após a etapa de build `e02`:

*Via Terminal (CLI):*

```bash
legoloop compor pipeline-frontend aidd-homologue@1.0.0 --depende e02 --yes
```

*Via Chat do Agente:*

```markdown
/legoloop-compor
Componha o sub-loop 'aidd-homologue@1.0.0' dentro do 'pipeline-frontend',
fazendo-o depender da etapa 'e02'.
```

*Via LegoLoop Studio:*

1. Abre `pipeline-frontend` no visualizador de DAG.
2. Clica em **+ Adicionar Sub-Loop** ou arrasta o bloco `aidd-homologue` da gaveta lateral.
3. Traça a aresta de `e02` até o nó do sub-loop.
4. Clica em **Salvar Composição**.

**Então** o LegoLoop:

1. Verifica se a junção criaria ciclos proibidos (`ErrCiclo`).
2. Isola a saída do sub-loop em `out/aidd-homologue/`.
3. Atualiza atomicamente `.processo.json` e `.dag.json` do pipeline pai.

**E** o sub-loop passa a executar como caixa-preta confiável dentro do macro-fluxo.

## 12. Cenário 5: Refinar Cirurgicamente uma Etapa

> **Objetivo:** trocar uma etapa lenta ou com ferramenta antiga sem recriar o pipeline.

**Dado que** a etapa `e03_analise_estatica` usa uma ferramenta obsoleta,

**Quando** o usuário solicita a substituição atômica:

*Via Terminal (CLI):*

```bash
legoloop refinar "auditoria-smart-contracts" \
  --trocar e03_analise_estatica --por "slither-analyzer@2.0.0" --yes
```

*Via Chat do Agente:*

```markdown
/legoloop-refinar
No processo 'auditoria-smart-contracts', refine a etapa 'e03_analise_estatica'
substituindo-a pela peça 'slither-analyzer@2.0.0'.
```

*Via LegoLoop Studio:*

1. Clica com o botão direito sobre o nó `e03_analise_estatica`.
2. Seleciona **Refinar / Trocar Peça** no menu de contexto.
3. Escolhe `slither-analyzer@2.0.0` no catálogo de peças aptas sugerido pela matriz.
4. Confirma no modal de validação de contratos.

**Então** o LegoLoop:

1. Valida a compatibilidade dos contratos de entrada e saída com os nós adjacentes.
2. Substitui o nó no DAG preservando todas as conexões.
3. Incrementa a versão menor do processo, de `@1.0.0` para `@1.1.0`.

**E** nenhuma outra etapa do fluxo sofre alteração ou efeito colateral.

## 13. Cenário 6: Remover com Segurança e Rollback

> **Objetivo:** excluir um processo sem deixar resíduos nos editores nem quebrar links.

**Dado que** o processo `experimento-legado` deve ser desativado,

**Quando** o usuário aciona a remoção segura:

*Via Terminal (CLI):*

```bash
legoloop remover "experimento-legado" --dry-run   # 1. simulação
legoloop remover "experimento-legado" --force     # 2. execução
```

*Via Chat do Agente:*

```markdown
/legoloop-remover
Remova o processo 'experimento-legado' com segurança, executando a simulação
primeiro e limpando os links simbólicos de todos os harnesses.
```

*Via LegoLoop Studio:*

1. No cartão de `experimento-legado`, abre o menu de contexto e seleciona **Remover Processo**.
2. O diálogo de **Impacto e Dry-Run** lista arquivos e harnesses afetados.
3. O usuário digita o nome do processo e clica em **Confirmar Remoção**.

**Então** o sistema exibe o relatório de desvinculação, limpa os symlinks nos harnesses, arquiva os dados e gera cópia de segurança para restauração imediata.

## 14. Cenário 7: Personalizar Micro e Nano-Etapas no `ROUTE.json`

> **Objetivo:** controlar custo cognitivo e velocidade ajustando modelo de IA e limites por bloco.

**Dado que** o arquivo `docs/routes/YYYY-MM-DD_ROUTE_EXEMPLO.json` define a execução por nós,

**Quando** o usuário ajusta as propriedades declarativas de um nó:

*Via Terminal ou editor de código:*

```json
{
  "id": "ticket-01",
  "titulo": "Gerar Casos de Teste Unitários",
  "modelo_cognitivo": "gemini-flash",
  "contexto_esteril": true,
  "mcps_permitidos": ["codebase-memory-mcp"],
  "stop_gates": ["SG-01-REVISAO-HUMANA"]
}
```

*Via Chat do Agente:*

```markdown
No arquivo de rota mais recente, configure o nó 'ticket-01' para usar o modelo
'gemini-flash', contexto estéril e apenas a ferramenta 'codebase-memory-mcp'.
```

*Via LegoLoop Studio:*

1. Clica no nó `ticket-01` para abrir o painel **Propriedades da Etapa**.
2. Seleciona `gemini-flash` no dropdown de modelo.
3. Ativa o switch **Contexto Estéril**.
4. Marca apenas o MCP `codebase-memory-mcp` e ativa o portão `SG-01-REVISAO-HUMANA`.
5. Clica em **Aplicar Alterações**.

**Então**, na execução pela Orchestrix:

1. O nó roda com o modelo econômico `gemini-flash`.
2. Carrega apenas `codebase-memory-mcp`, sem poluir a janela de contexto.
3. Pausa obrigatoriamente no portão humano `SG-01-REVISAO-HUMANA` antes de gravar em disco.
