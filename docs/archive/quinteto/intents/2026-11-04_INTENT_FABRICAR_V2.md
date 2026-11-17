# 🎯 INTENT: Fábrica de Processos v2 (correção e expansão do `fabricar`)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-04  
> **Autor:** Everton & Claude  
> **Alvo:** `internal/fabricar/`, `internal/interview/`, `internal/reputation/`, `internal/canon/`, `internal/provider/`, `internal/plangen/`, `internal/specgen/`, `internal/cli/`, `internal/studio/`, `skills/`, `cmd/legoloop/`  
> **Sessão:** 872376ce-a2c0-45d2-aa66-989fbaa25a97  
> **Confiança da Entrevista:** 95% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A v1 do `fabricar` (commit `f783929`) grava um pacote válido, mas não entrega a intenção do INTENT de 2026-11-03. A auditoria (`docs/reports/2026-11-03_REPORT_AUDITORIA_FABRICAR.md`) provou com o binário real e o 9Router: (a) o modo interativo da CLI nunca funcionou, porque o `main` não liga o stdin; (b) a skill do Chat troca a entrevista por um formulário e usa `--yes`, e o Studio descarta o espelho da entrevista, então o INTENT gerado sai com 0% de confiança e "Não informado" em 7 campos; (c) a Fusão ordena só por `Score` de reputação, sem usar o objetivo, e o `build-engine` ficou com `aidd-site-dominator` e sem `aidd-tdd`; (d) SPEC, PLAN, ROUTE e DAG contam histórias diferentes; (e) os gates são só `file_exists`; (f) a skill não chega aos harnesses; (g) o cérebro só aceita o 9Router; (h) a matriz de reputação do cofre vem de um script avulso em `scratch/` e nenhum comando do produto a gera.
* **Quem Sofre:** O autor de processos, a cada processo fabricado; e o usuário que roda o resultado, que recebe um processo vazio.
* **O Porquê:** O processo fabricado é uma peça durável e reutilizável. Se a base é anêmica, toda peça derivada herda a falha. Quanto mais processos existirem, maior o custo de corrigir depois.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** `legoloop fabricar` cria processos com INTENT em 0%, peças escolhidas por reputação, Tríade incoerente, gates triviais, só no 9Router, sem editar, listar, remover ou compor, e com saída que explica pouco.
  * **Amanhã:** A fábrica entrevista de verdade, escolhe peças pelo caso de uso (com o cérebro explicando), gera uma Tríade coerente a partir de um contrato por etapa, entrega o processo a todos os harnesses, funciona com qualquer cérebro disponível, explica cada ação antes de executá-la e tem os cinco verbos (`fabricar`, `listar`, `refinar`, `remover`, `compor`) em CLI, Chat e Studio.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: **Processo universal.** A Tríade descreve o próprio processo, e PLAN e ROUTE são o grafo que roda. A natureza do processo (software, extração de dados, documentação, outras) vem da entrevista; não há lista fechada de categorias nem blocos de Go fixos. O processo é durável e reutilizável; cada execução é efêmera e estéril.
* [ ] Objetivo 2: **Entrevista socrática contextual** em CLI, Chat e Studio: lê o contexto antes de perguntar, mostra o espelho e a confiança a cada turno, e só fecha em ≥ 95% depois que o usuário confirma o espelho final. O resultado alimenta INTENT, SPEC, PLAN e ROUTE completos, sem "Não informado".
* [ ] Objetivo 3: **Motor de matrizes** próprio, com a matriz geral e a matriz por tipo de peça, e um **catalogador** que preenche o perfil de uso de cada peça. O vocabulário de casos de uso é fechado, versionado e cresce com aprovação do usuário. O motor roda sempre no fim do `canonize --apply` e sob demanda nas três portas, de forma incremental por hash. O `fabricar` só consulta e alerta quando a matriz está defasada ou ausente.
* [ ] Objetivo 4: **Fusão ~95% determinística.** A aptidão ao caso de uso filtra e ordena; o `Score` só desempata e rebaixa peças tóxicas; o cérebro explica "entrou porque" e "ficou de fora porque" e pode contestar. O usuário revê e troca peças antes de gravar. Peça sem perfil entra como "aptidão desconhecida" e nunca é escolhida em silêncio; lacunas sem peça são avisadas.
* [ ] Objetivo 5: **Tríade coerente por contrato único de etapa** (entrada, saída, peça executora, gate). SPEC, PLAN, ROUTE, DAG, loop e skill são projeções do mesmo contrato. Os gates vêm dos critérios de sucesso da entrevista e de uma biblioteca de gates por tipo de artefato; gate só de existência é recusado quando a etapa promete conteúdo.
* [ ] Objetivo 6: **Cérebro plugável.** Qualquer executor do inventário pode conduzir a fabricação (9Router, Ollama, CLI de harness, o agente do próprio Chat). O 9Router deixa de ser requisito. A escolha automática mostra qual foi usada; sem cérebro, o modo determinístico continua e a mensagem explica o que falta.
* [ ] Objetivo 7: **Alcance universal.** O cofre é a fonte única e o processo é ligado a todos os harnesses detectados, com a lista na prévia, `--no-link` e `--link-harness X`. As skills do próprio LegoLoop seguem o mesmo caminho.
* [ ] Objetivo 8: **Cinco funções em fatia vertical** (`fabricar`, `listar`, `refinar`, `remover`, `compor`), cada uma com pacote, comando, skill, rota, ação no Studio, manual e ACs próprios, sobre um núcleo compartilhado mínimo e sem regra de negócio. `refinar` mostra o diff, faz backup, sobe a versão e guarda o histórico. `remover` exige confirmação, tem rollback e desvincula os harnesses. `compor` liga um processo a outro por nome e versão, sem ciclo, e bloqueia a remoção de um processo em uso.
* [ ] Objetivo 9: **Comunicação clara e personalizável.** Toda ação com efeito é explicada antes, em português simples no perfil STE-PT (inspirado no ASD-STE100: frases curtas, um termo por coisa, aviso de risco antes). Plano único no início, com pausa extra só para ação irreversível ou fora do projeto. Tudo é personalizável pelo usuário. A mesma explicação aparece nas três portas. Um linter confere as mensagens do produto contra o glossário (`AC-STEPT100`).
* [ ] Objetivo 10: **Relatório por camada e prévia com conteúdo.** Cada camada entrega o seu resumo (espelho e confiança; peças vistas, cortes e idade da matriz; ranking com motivos; ACs, tickets e gates; arquivos e harnesses). É um único objeto, impresso pela CLI, repassado pelo Chat e desenhado pelo Studio.
* [ ] Objetivo 11: **Prova de comportamento.** Os ACs descrevem comportamento observável e são provados pelo binário instalado, nas três portas, em pelo menos três cenários de natureza diferente, com mais de um cérebro, com mutação nos pontos críticos e homologação humana por épico.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto (só stdlib em Go).
* ⛔ Não aplicaremos o STE-PT à saída do chat dos agentes: é demanda futura, só o `AC-STEPT100` ficará definido.
* ⛔ Não trataremos o `go test -race` no CI Linux; segue como pendência separada.
* ⛔ Não faremos os verbos como processos fabricados: são funções do binário.
* ⛔ Não regeneraremos a matriz dentro do `fabricar`: ele só consulta e alerta.
* ⛔ Não removeremos o prefixo `legoloop-` nem o comportamento de `legoloop run`; `canonize` segue desduplicando por SHA-256 e apenas passa a chamar o motor de matrizes ao final do `--apply`.
* ⛔ Não removeremos as garantias da v1: gravação atômica com rollback, default-deny de MCP, hook, script e loop varridos, mídia removível `Untrusted`, chave do roteador só para a URL confiável.

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar. Tipos permitidos: `skill`, `script`, `subloop`, `unit_task`, `ui_task`, `stop_gate`.

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `unit_task` | Núcleo compartilhado (`internal/processo/`) | Modelo do processo, gravação atômica, vínculo a harnesses e consulta à matriz, sem regra de negócio. |
| `BLOCO-02` | `unit_task` | Motor de matrizes (`internal/reputation/`, `internal/canon/`) | Matriz geral e por tipo, catalogador, vocabulário fechado, gatilho no `canonize --apply`, defasagem por hash. |
| `BLOCO-03` | `unit_task` | Cérebro plugável (`internal/provider/`, `internal/fabricar/`) | Abstração de executor, escolha automática e mensagem sem cérebro. |
| `BLOCO-04` | `unit_task` | Entrevista contextual (`internal/interview/`) | Leitura de contexto, espelho e confiança por turno, confirmação final. |
| `BLOCO-05` | `unit_task` | Fusão e contrato de etapa (`internal/fabricar/`, `internal/plangen/`, `internal/specgen/`) | Seleção por aptidão, explicação, contrato único e projeções coerentes, biblioteca de gates. |
| `BLOCO-06` | `unit_task` | Verbos `listar`, `refinar`, `remover`, `compor` | Uma fatia vertical por verbo. |
| `BLOCO-07` | `unit_task` | Comunicação STE-PT (`internal/` pacote de mensagens) | Cartão de ação, glossário e linter das mensagens. |
| `BLOCO-08` | `ui_task` | Studio (`internal/studio/`) | Modal e ações para cada verbo, com a mesma prévia e o mesmo relatório. |
| `BLOCO-09` | `skill` | `skills/legoloop-*` | Uma skill de Chat por função, casca fina sobre o binário, ligada a todos os harnesses. |
| `BLOCO-10` | `script` | `go test ./...` | Regressão mecânica; os ACs de comportamento rodam o binário instalado. |
| `BLOCO-11` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop` | Recompilação fresca e binário do PATH idêntico (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`: `cli: legoloop fabricar|listar|refinar|remover|compor|matriz`, `chat: /legoloop-<função>`, `studio: ações do painel de processos`. O motor de matrizes também dispara no fim de `canonize --apply`.  
  *(Define como este fluxo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** `Médio`. Grava no cofre `~/.aidd/canonical/`, no projeto e nas pastas dos harnesses (links). Mitigação: gravação atômica com rollback, plano explicado antes, pausa extra em ação irreversível, testes hermeticos que nunca tocam o `~/.aidd` real e provas em pasta descartável.
* **Política de Privacidade & Orçamento Cognitivo:** O cérebro é o que o usuário escolher: nuvem autorizada ou exclusivamente local (Ollama). Conteúdo de peças varridas é dado não confiável e nunca é instrução. Dados sensíveis: a chave do roteador, que fica na configuração do usuário e só é enviada à URL confiável.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. No `build-engine` (cenário que originou esta INTENT), a entrevista real conduzida pelo binário instalado termina em ≥ 95% com espelho confirmado, e o INTENT gerado não contém "Não informado".
  2. O processo gerado escolhe peças relevantes ao objetivo (por exemplo `aidd-tdd`) e exclui as irrelevantes (por exemplo `aidd-site-dominator`), com o motivo de cada decisão.
  3. SPEC, PLAN, ROUTE, DAG, loop e skill têm o mesmo conjunto de etapas, verificado por teste.
  4. O processo aparece nos harnesses detectados e no Catálogo do Studio, e `legoloop run` executa o `.dag.json`.
  5. Os três cenários de referência passam nas três portas, com 9Router e com um segundo cérebro.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * `compor` referencia por nome e versão, sem ciclo, e bloqueia a remoção de um processo em uso.
  * A categoria (`build|reverse|web|mobile`) vira rótulo opcional e a natureza vem da entrevista.
  * A sequência de entregas é: (a) núcleo, motor de matrizes e cérebro; (b) `fabricar` refeito; (c) `listar`, `refinar`, `remover`, `compor`; (d) Studio e homologação final, cada épico com um stop-gate humano.
  * Hoje a matriz do cofre (1.588 peças, de 02/11) é o único acervo catalogado e quase nenhuma peça tem perfil de uso; o catalogador parte do frontmatter e do texto e carimba o perfil como "inferido" até a revisão do usuário.
* **Riscos e Perguntas em Aberto:**
  * **Catalogar ~1.600 peças com LLM** é lento e caro; mitigação: incremental por hash e perfil determinístico primeiro, cérebro só nas lacunas.
  * **O vocabulário fechado** pode não cobrir peças existentes; mitigação: crescimento por proposta aprovada.
  * **Cérebros fracos** podem gerar entrevistas ruins; a prova de aceite exige um segundo executor, e a escolha automática avisa qual foi usada.
  * **Qual é o segundo cérebro de referência** para a prova (Ollama local ou a CLI de um harness): a definir na SPEC com o usuário.
  * **Migração dos processos já fabricados** (ex.: `legoloop-*` gravados pela v1): fora do escopo até decisão.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0; `go vet ./...`; `gofmt -l .` vazio; `tools/funclen` limpo |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green), com mutação nos pontos críticos. | Teste novo falhando antes do código; a mutação do código é detectada | Relatório de cada ticket com a mutação aplicada e o teste que a pegou |
| 3 | A entrevista é real em todas as portas e nunca sai com 0% sem aviso. | Binário instalado, 9Router real, saída capturada | Transcrição da entrevista e INTENT gerado, sem "Não informado" |
| 4 | A Fusão usa o caso de uso. | Cenários de referência com peças esperadas e proibidas | Relatório por camada com aptidão, `Score` e motivos |
| 5 | As seis visões do processo são coerentes. | Teste que compara etapas de SPEC, PLAN, ROUTE, DAG, loop e skill | `go test` do pacote de contrato |
| 6 | Funciona sem o 9Router. | Execução com um segundo cérebro e mensagem clara sem nenhum | Saída capturada nos dois casos |
| 7 | O processo chega a todos os harnesses detectados. | Verificação dos links criados e do rollback | Listagem dos harnesses e teste de rollback |
| 8 | Os cinco verbos funcionam nas três portas. | CLI, harness real e navegador real por verbo | Relatório de homologação por épico com imagens |
| 9 | As mensagens do produto seguem o perfil STE-PT. | Linter com glossário e limite de frase | `AC-STEPT100` verde |
| 10 | A matriz é gerada pelo produto e o `fabricar` avisa quando ela está defasada. | `canonize --apply` atualiza a matriz; alteração de peça gera o alerta | Teste hermético e execução real em pasta descartável |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`](../specs/2026-11-04_SPEC_FABRICAR_V2.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
