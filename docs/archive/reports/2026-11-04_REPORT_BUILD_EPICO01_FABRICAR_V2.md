# 📊 REPORT: BUILD do Épico 01 (Núcleo, motor de matrizes e cérebro plugável) — `SG-E01`

> **Data:** 2026-11-04  
> **Escopo:** `TICK-01.1` a `TICK-01.11` do PLAN `2026-11-04_PLAN_FABRICAR_V2.md`  
> **Estado:** 11 tickets `PASSED`, `AC-E01.1` a `AC-E01.10` marcados `[x]` na SPEC (os dois últimos fecham as suas respostas de 04/11, seção 7). Aguardando a **sua** homologação e a autorização do `SG-E01` para o Épico 02.  
> **Executor:** sessão principal, sem subagentes.  

---

## 1. Entregas

| Ticket | Entrega | Onde |
|:---|:---|:---|
| `01.1` | Perfil **STE-PT**: `Lint` (≤ 20 palavras, glossário, palavras ambíguas), `Cartao` com o aviso de risco **antes** dos passos, `Decisao`, catálogo de mensagens | `internal/ste` |
| `01.2` | **Núcleo do processo**: contrato `Processo`/`Etapa`, `Validate` (nome, semver, ciclo, etapa sem conexão, gate que compila, artefato dentro do projeto) e `Store.Write` atômico com rollback. O motor de gravação foi **movido** do `fabricar` para cá; `fabricar.Materialize` virou camada fina sobre ele | `internal/processo` |
| `01.3` | **Cérebro plugável**: `brain.Pick` escolhe 9Router, Ollama ou CLI de harness pelo inventário e **explica** a escolha; `Model()` e `Timeout()` | `internal/brain` |
| `01.4` | **Vocabulário fechado** e versionado (9 etapas, 7 naturezas, 22 casos de uso), fila de propostas, `Aprovar`, **catalogador** (declarado, depois palavras-chave, depois cérebro só na lacuna) | `internal/matriz` |
| `01.5` | `Generate`: matriz geral + `matrix/<tipo>.json`, compatível com a leitura da v1 (`ratings`) | `internal/matriz` |
| `01.6` | **Incremental por hash**; peça sem pista fica esperando o cérebro | `internal/matriz` |
| `01.7` | `Staleness`, `Previa` e as mensagens de defasagem em STE-PT | `internal/matriz` |
| `01.8` | Comando **`legoloop matriz`** (`gerar`, `status`, `vocab listar|pendentes|aprovar`) e o gancho no `canonize` | `internal/cli` |
| `01.9` | `--model` validado pelo `brain.Pick`: modelo inexistente sai com **2** e lista os disponíveis; o 9Router **deixou de ser requisito** (CLI e Studio) | `internal/cli`, `internal/studio` |
| `01.10` | **`matriz perfil confirmar` e `definir`**: a palavra do usuário fica em `perfil.confirmado.json` (por hash), atualiza a matriz geral e a do tipo numa gravação atômica, sobrevive a `matriz gerar --forcar` **sem chamar o cérebro** para a peça, e perde a validade (com aviso) se o conteúdo mudar | `internal/matriz`, `internal/cli` |
| `01.11` | **`matriz ver`** (fusão do `matrix`): tabela com tier, nota, tipo, naturezas, etapas, casos, origem e cópias agrupadas; filtros `--tipo --tier --min --natureza --etapa --caso --origem --copias --json`. `legoloop matrix` virou **alias descontinuado** (avisa no stderr; `--kind` e `--path` são traduzidos). A `RunMatrix` de 96 linhas foi desfeita | `internal/cli` |

## 2. Checklist do `SG-E01`

| Item | Resultado |
|:---|:---|
| `go test -count=1 ./...` | **25 de 25 pacotes**, **658** testes de topo + **104** subtestes, **0 falhas**, 1 pulado (o preexistente de symlink em `assets`). Contagem medida com `go test -json` ao fechar o `TICK-01.11` |
| `go vet ./...`, cross-vet `linux` e `darwin` | Limpos (reexecutados ao fechar o `TICK-01.11`) |
| `gofmt -l .` | Vazio |
| `funclen` | Limpo em `processo`, `brain`, `matriz`, `ste`, `orchestrix`, `fabricar`, `provider`. **Correção de um relato meu anterior:** `internal/cli` (`cmdRun`, `cmdImport`, `RunUpgrade`; a `RunMatrix` **saiu da lista** no `TICK-01.11`) e `internal/studio` (`loopItems`, `handleAssetImport`, `buildDAG`) **já violavam no `HEAD`**; meus relatórios de "funclen limpo" valiam só para os pacotes que eu tinha checado. Não mexi nessas funções |
| `go.mod` | Intacto (o projeto não tem `go.sum`: não há dependências) |
| `go install` | Binário do PATH idêntico ao `bin/` |
| Cofre real `~/.aidd` | Até a seção 6, a v1 ficou intacta (`67bf193acbd825e0`). Depois, com a sua autorização (seção 7), a v2 foi gerada; nesta rodada o `matriz ver` só **leu** o arquivo real e o hash dele ficou **idêntico antes e depois** |
| Mutações | ≈ **170** mutações deliberadas em 18 suítes, mais **58** nesta rodada (20 no motor do perfil confirmado, 8 na CLI do perfil, 30 no `matriz ver` e no alias); **todas detectadas** ao final. Na 1ª passada, 10 dessas 58 sobreviveram ou não compilaram; cada uma virou teste ou mutante corrigido (seção 8) |
| Autorização do usuário | **Pendente** |

## 3. Prova pelos critérios de aceite

| AC | Prova de comportamento |
|:---|:---|
| `AC-E01.1` | Cofre de 12 skills (4 declaradas), 2 scripts e 1 hook gera `registry_matrix.json` 2.0.0 com 15 entradas e `matrix/skill.json` com 12; todo termo pertence ao vocabulário |
| `AC-E01.2` | Sem mudança: **0** reprocessadas, **0** chamadas ao cérebro, **0** recálculos de nota, arquivos byte a byte iguais; com 1 skill alterada: exatamente 1 |
| `AC-E01.3` | Pelo **binário real** com `USERPROFILE` isolado: `--dry-run` informa "3 peças seriam recatalogadas" e não cria matriz; `--apply` gera a matriz e relata; o JSON do `canonize` segue puro |
| `AC-E01.4` | `raspagem-web` vira proposta pendente e só entra com `matriz vocab aprovar` (a versão sobe) |
| `AC-E01.5` | 1 nova, 1 alterada e 1 removida são detectadas; 1.600 peças em ≪ 2 s |
| `AC-E01.6` | Só Ollama, 9Router+Ollama, só `claude`, vazio: cada caso escolhe o executor certo e explica; vazio dá `ErrSemCerebro` com o que falta |
| `AC-E01.7` | **Sistema real**: `fabricar x --model inexistente` com o inventário real sai com 2, mostra os 15 primeiros dos 243 modelos e "e mais 228. Veja todos com legoloop doctor", e grava **0 arquivos** |
| `AC-E01.8` | Falha na 3ª escrita de 8 arquivos deixa a árvore com hash idêntico, sem `.legoloop-tmp` |
| `AC-E01.9` | **Binário instalado**, `HOME` isolado: `definir z --etapa construir --natureza software --caso tdd` e `confirmar` gravam `perfil.confirmado.json`; `matriz gerar --sem-cerebro --forcar` mantém `z` como `confirmado`; `definir --etapa inventada` sai **2**; `confirmar` de peça sem perfil e peça inexistente saem **1**; ao alterar o conteúdo de `z`, o resumo avisa e a peça é catalogada de novo |
| `AC-E01.10` | **Binário instalado**, `HOME` isolado, 5 caminhos (3 cópias de `x`): `matriz ver` mostra **3 linhas** com `x` e a contagem 3; `--copias` mostra **5**; `--natureza`, `--etapa`, `--origem`, `--tier`, `--min`, `--caso`, `--tipo` e `--json` filtram; `--tier platina`, `--natureza inventada`, `--min 2` e argumento sobrando saem **2**; sem matriz sai **1** e manda rodar `matriz gerar`; `legoloop matrix --kind skill` dá o mesmo resultado e avisa "descontinuado". **Cofre real, só leitura:** 6.207 caminhos viram **898 linhas** agrupadas |

## 4. O que o sistema real mostrou (e os testes sintéticos não)

Rodei o binário e o 9Router reais. Cada achado virou teste vermelho antes da correção.

| # | Achado | Correção |
|:---|:---|:---|
| 1 | **`Ref` colidia** em MCP e hook, que dividem o mesmo arquivo de configuração (`~/.claude.json`): a defasagem acusava "22 novas, 4 alteradas" sem mudança nenhuma | `Ref` própria `arquivo#nome` para MCP e hook |
| 2 | **Cópias idênticas**: 6.174 skills por caminho, só **869 conteúdos únicos**. Esperando o cérebro: 3.366 caminhos = **507 conteúdos**. O cérebro seria chamado ~7× mais que o necessário | Perfil por conteúdo (tipo + hash): cópias compartilham o perfil, sem nova chamada. `Summary.Copias`. Cópia **não** herda falha |
| 3 | **Bug meu de escape** (`\p{L}` virou `p{L}` numa edição por `node`): a "fronteira de palavra" passou a casar no meio das palavras. **Meu "ganho de 50×" no desempenho era mentira**: a regex quebrada era barata. O cache de regex vale, mas o ganho real é outro | Restaurada pelo `Edit`; teste de meio de palavra; guarda de 10 s para a 1ª geração de 1.600 peças (sem cache: 25 s) |
| 4 | **Cérebro liberal demais**: o `aidd-site-dominator` saiu com 6 casos de uso, 3 etapas e 3 naturezas | Prompt pede só o propósito principal; **limite em Go**: 3 casos, 2 etapas, 2 naturezas, 3 propostas por peça. O declarado pelo autor não é cortado |
| 5 | **230 modelos despejados** na mensagem de erro | Lista resumida: 15 e "e mais N. Veja todos com legoloop doctor" |
| 6 | **Teste antigo varria o HOME real** (38 s, não hermético) depois do gancho do `canonize` | O helper `deps(t)` injeta varredura vazia |
| 7 | `orchBrain` usava o `Timeout` do cliente HTTP como prazo: `&http.Client{}` cancelaria na hora | O prazo vem da configuração do cérebro |
| 8 | Prévia do `canonize` contava caminhos; o `--apply` conta conteúdos | `Previa` conta conteúdos distintos |

## 5. Mudanças de contrato em relação à v1 e ao PLAN

* **O 9Router deixou de ser requisito.** Dois testes da v1 codificavam o contrário e foram **substituídos** (um recusava modelo que não fosse do 9Router; outro explicava só o 9Router). Dois subtestes foram **atualizados** (modelo inexistente agora é exit 2; "sem LLM" agora é "sem nenhum modelo"). `fabricar.BrainFor`, `ErrBrain` e `BrainTimeout` foram removidos.
* **Studio:** modelo desconhecido vira 422 e "sem cérebro" vira 503, com mensagens próprias (campos de teste novos: `BrainOpts`, `FabricarScanner`).
* **`Staleness` devolve `(Stale, error)` e recebe as peças e `Deps`**, e não `(ctx, reg, snap) Stale` como a SPEC rascunhava: precisava propagar o cancelamento e reusar a identidade da geração.
* **`canonize --apply` gera a matriz sem cérebro** (decisão minha): o `canonize` não gasta modelo sem você saber. Peças sem pista ficam esperando o cérebro, que `matriz gerar` completa. Com `--json`, o aviso da matriz vai ao `stderr`, e o contrato JSON do `canonize` não muda.
* **Entrevista passo a passo, Fusão e Tríade:** nada disso mudou (é o Épico 02). Por isso o `fabricar` real ainda escolhe as mesmas "5 peças curadas" da v1.

## 6. O que **não** foi provado

* **Os executores por CLI** (`claude`, `agy`, `omp`...) foram provados só na escolha e no roteamento. Nenhum processo real foi disparado por eles (entra no `TICK-04.2`).
* **O cérebro real catalogou só 6 skills.** O catálogo completo (507 conteúdos esperando) **não foi rodado**. Estimativa grosseira de 15 a 30 minutos, **não medida**.
* **A matriz real do cofre (`~/.aidd/canonical/registry_matrix.json`) continua a v1, intacta.** Toda a verificação usou cofres temporários. Gerar a real exige a sua autorização.
* **Qualidade do perfil do cérebro:** mesmo limitado, o `aidd-site-dominator` ainda leva `software` (o modelo lê "reconstruir uma aplicação" como software). A defesa é a aptidão por casos de uso (`TICK-02.5`). Todo perfil do cérebro nasce `inferido`; **não existe ainda comando para confirmá-lo**.
* **`-race`** segue só no CI Linux.
* **Alertas de texto suspeito** têm falsos positivos esperados: skills de segurança que citam `rm -rf` ou "ignore previous instructions" como exemplo (5 no cofre real). Elas são citadas como dado, sem bloquear nada.
* **Singular/plural:** as mensagens dizem "1 peças". Polimento pendente.

## 7. Respostas do usuário (2026-11-04, depois deste relatório)

| # | Pergunta | Decisão |
|:---|:---|:---|
| 1 | Gerar a matriz real? | **Sim, seguir a recomendação:** `matriz gerar --sem-cerebro` agora; a catalogação com o cérebro fica para uma sessão separada. **Feito em 04/11:** 6.207 peças, 2.423 cópias, 3.393 esperando o cérebro. Cópia de segurança da v1 em `~/.aidd/canonical/registry_matrix.v1.2026-11-02.bak.json` (hash `67bf193acbd825e0`) |
| 2 | Matriz por caminho ou por conteúdo? | **Por caminho** (recomendação técnica): cada caminho é um arquivo real que um harness lê, e a defasagem e a "chegada a todos os harnesses" precisam do estado por caminho. A deduplicação acontece no **consumo** (Fusão por nome e hash; `matriz ver` agrupa as cópias) |
| 3 | Comando para confirmar o perfil? | **Criar junto, no Épico 01:** `TICK-01.10`, `matriz perfil confirmar` e `matriz perfil definir`. A confirmação fica num arquivo próprio do cofre (por hash), para sobreviver às regerações |
| 4 | `matrix` e `matriz`? | **Fusão** (recomendação técnica): um só comando `matriz`, com `matriz ver` herdando o melhor do `matrix` (tabela de reputação, `--tier`, `--min`, `--json`) e ganhando tipo real, perfil de uso e cópias agrupadas. `legoloop matrix` vira alias descontinuado de `matriz ver`. `TICK-01.11` |
| 5 | Natureza vazia = "qualquer"? | **Sim** (recomendação técnica): vazia vale como coringa na Fusão, e o casamento **explícito** de natureza desempata a favor. Peça **sem perfil** é "desconhecida" e nunca é escolhida. Entra como premissa 9 da SPEC e no escopo do `TICK-02.5` |

**Feito (seção 8):** `AC-E01.9` e `AC-E01.10` na SPEC, `TICK-01.10` e `TICK-01.11` no PLAN e na ROUTE, e os dois implementados. O `SG-E01` volta a ser apresentado a você.

## 8. Fechamento do Épico 01 (`TICK-01.10` e `TICK-01.11`)

**O que a prova real mostrou** (cada item virou teste ou correção):

| # | Achado | Correção |
|:---|:---|:---|
| 1 | O aviso de confirmação perdida dizia "o perfil voltou a inferido", mas sem cérebro a peça fica **sem** perfil | Mensagem agora diz "a peça foi catalogada de novo"; a SPEC (`AC-E01.9`) diz o mesmo |
| 2 | O `help` ainda descrevia o `matrix` como comando vivo e a coluna do `matriz` ficou torta | `help` marca o `matrix` como DESCONTINUADO e realinha |
| 3 | O exemplo da SPEC usava `--caso "criar api"`, que **não é termo do vocabulário** (kebab-case, fechado) | Exemplo trocado por `--caso tdd` |
| 4 | Mutações que sobreviveram à 1ª passada: nome ambíguo, `definir` sem termo, `consome/produz` apagados, `retentar` não limpo, peça removida avisando à toa, `--tier`, `--min` inclusivo, filtro de termo deixando passar peça sem perfil, hash vazio da v1 agrupando, mesmo hash em tipos diferentes, maiúsculas, cabeçalho `CAMINHO` | Um teste para cada, incluindo um registro montado à mão com notas e tiers variados |

**Decisões minhas, para você revisar:**

* **Nome ambíguo:** se o mesmo nome tem conteúdos diferentes (hashes distintos), `perfil confirmar|definir` recusa e lista os caminhos. Cópias idênticas contam como uma peça só e **herdam** a confirmação.
* **Validade da confirmação:** vale por hash do conteúdo. Se o conteúdo muda, a confirmação **não** é apagada do arquivo (se o conteúdo voltar, ela volta a valer), mas a geração avisa a cada rodada enquanto durar a diferença.
* **`definir` exige ao menos um termo** e cada termo precisa estar na categoria certa do vocabulário. `confirmar` exige que a peça já tenha perfil.
* **`matriz ver` valida os termos de `--natureza`, `--etapa` e `--caso` contra o vocabulário** (exit 2), para um erro de digitação não parecer "nenhuma peça".
* **Tabela do `matriz ver`:** deixei de mostrar AUTOR, RIGOR, CONTRATOS, PROVENIÊNCIA e HISTÓRICO para caber o perfil de uso. **Continuam no `--json`** (campo `ratings`, no formato antigo, mais o campo novo `entradas`).
* **`matrix` com matriz da v1:** continua legível; o tipo vem do caminho, como antes.

**O que não foi provado:**

* **Perfil `inferido` por cérebro real** no `matriz ver`: a prova real usou `--sem-cerebro`. O cérebro real (507 conteúdos à espera) segue para a sessão separada que você escolheu.
* **Peças sem perfil são comuns sem cérebro:** no cofre de prova, uma skill com descrição clara de TDD ficou **sem perfil** (as palavras-chave não a pegaram). É comportamento do `TICK-01.4`, não deste ticket, mas mostra por que a catalogação com cérebro importa antes do Épico 02.
* **Windows apenas** nas provas pelo binário; o `vet` cruzado cobre linux e darwin, mas nenhum teste rodou lá.
* **Estado do git:** **nada foi commitado** (31 entradas em `git status`, contando pastas novas inteiras). Isso é decisão sua.

## 9. Decisões e perguntas originais (arquivo)

1. **Gerar a matriz real agora?** Grava em `~/.aidd/canonical/` (os arquivos antigos viram `.legoloop-bak`). Sugiro `matriz gerar` **sem cérebro** primeiro (≈ 30 s), depois a catalogação com o cérebro em uma sessão separada.
2. **Armazenar a matriz por caminho (6.207 entradas, 5,4 MB) ou por conteúdo (≈ 869)?** Hoje é por caminho, para o Studio e o `matrix` continuarem funcionando. Por conteúdo é menor e mais limpo, mas muda a leitura.
3. **Comando para confirmar um perfil** (`inferido` para `confirmado`): crio no Épico 03, junto de `listar`, ou como ticket extra aqui?
4. **`legoloop matrix` (v1, reputação) e `legoloop matriz` (v2)** ficam parecidos demais. Quer um alias ou fusão?
5. **No `fabricar`, naturezas vazias significam "qualquer"?** (Skills como `aidd-spec` e `aidd-grill` saíram sem natureza.) É o que a Fusão do `TICK-02.5` vai assumir, salvo objeção.

## 10. Autorização

- [ ] Evidências aprovadas pelo usuário.
- [ ] Autorização para iniciar o Épico 02 (`SG-E01`).
