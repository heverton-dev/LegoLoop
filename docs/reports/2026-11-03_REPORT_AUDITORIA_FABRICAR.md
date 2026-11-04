# 🔎 REPORT: Auditoria do `legoloop fabricar` contra o INTENT original

> **Data:** 2026-11-03  
> **Auditado:** commit `19feb9d` (`f783929` + docs)  
> **Referência:** `docs/intents/2026-11-03_INTENT_LEGOLOOP_FABRICAR.md` e `docs/specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md`  
> **Método:** leitura do código e execução real do binário contra o 9Router (`code-fast`). Gravações só em `scratch/entrevista-real/` (ignorado pelo git); o cofre real `~/.aidd` não foi tocado.  
> **Veredito:** a fábrica **grava** um pacote válido, mas **não entrega a intenção**. As três portas produzem um INTENT com 0% de confiança, a fusão não usa o objetivo e a Tríade sai incoerente. Os 12 ACs marcados `[x]` na SPEC verificaram a forma, não o comportamento prometido.

---

## 1. Promessas do INTENT × realidade

| Promessa (linha do INTENT) | Status | Evidência |
| :--- | :--- | :--- |
| Obj. 1: prefixo `legoloop-` fixo (26) | ✅ CUMPRIDO | `NormalizeProcessName`, fuzz; processo real `legoloop-build-engine` |
| Camada 1: entrevista socrática até ≥ 95% (28) | ❌ NÃO CUMPRIDO | ver §2, achados A1–A4 |
| Camada 2: varrer todos os drives e catalogar qualquer peça (29) | 🟡 PARCIAL | varre C:, D:, F:; corta em 500 peças e na profundidade 4 (aviso real: "4848 itens além da profundidade 4", "limite de 500 peças atingido"); quais peças entram é arbitrário |
| Camada 3: "combinação **semântica** dos pontos fortes" (30) | ❌ NÃO CUMPRIDO | `merger.go:29-57`: ordena por `Score` de reputação e pega as 5 primeiras; o objetivo não entra. Requisito perdido na SPEC (linha 48 tornou a montagem "determinística") |
| Camada 4: Tríade com isolamento, allowlist e gates booleanos (31) | 🟡 PARCIAL | a SPEC (LLM) fala do objetivo; o PLAN/ROUTE (`plangen.Render`) é a cadeia das 5 peças; os gates são só `file_exists` |
| Camada 5: gabarito contratual de saída (32) | 🟡 PARCIAL | `TEMPLATE_<NOME>.md` gerado é uma lista de `{{CONTEUDO_PASSO_N}}` sem contrato |
| Obj. 3: pronto para rodar nos 3 meios (33) | 🟡 PARCIAL | `.dag.json` e `.loop.json` válidos; a `SKILL.md` não chega a nenhum harness sem `--link-harness` (premissa da linha 76) |
| Obj. 4: CLI com modos interativo e autônomo (34) | ❌ NÃO CUMPRIDO | modo interativo inalcançável no binário real (A1) |
| Obj. 5: modal no Studio com seletor de modelo (35) | 🟡 PARCIAL | o modelo escolhido só conduz a Tríade; a entrevista usa `code-pro` fixo e o resultado dela é descartado (A3) |
| Effort clamping por camada (68) | ❌ NÃO CUMPRIDO | não existe; `emit_dag.go:96` fixa `EffortMedium` em todo nó |
| Prova 1: `fabricar "teste do joaozinho" --yes` gera o pacote (70) | ✅ CUMPRIDO | reproduzido com LLM real: exit 0, 8 arquivos, ~28 s |
| Prova 2: aparece no Studio (71) | ✅ CUMPRIDO | `Assets.Refresh`, homologado no navegador |
| Prova 3: `go test` verde (72) | ✅ CUMPRIDO | e insuficiente: ver §3 |

## 2. Achados reproduzidos com o sistema real

| # | Achado | Evidência |
| :--- | :--- | :--- |
| A1 | O modo interativo da CLI **nunca funcionou no binário**: `main` não passa o stdin | `cmd/legoloop/main.go:14` (`cli.Deps{Stdout, Stderr}`); saída real: "o modo interativo exige um terminal" |
| A2 | A skill do Chat troca a entrevista por um formulário de 4 campos e usa `--yes`, que chama `autonomous()`: 0 turnos, confiança 0 | `skills/legoloop-fabricar/SKILL.md:35-38`; `clarifier.go:40-56` |
| A3 | O Studio entrevista, mas envia só `mirror.outcome` como `goal`; o servidor monta `Request` sem `Interactive` e cai no `autonomous()`: o resto do espelho é descartado | `fabricar.js:32-37`; `fabricar_api.go:83` |
| A4 | A entrevista é cega ao contexto: para "refinar o BUILD do mega-pipeline" palpitou Docker, GitHub Actions e Node.js; fechou em 3 turnos aceitando palpites errados, sem espelho para confirmar e sem mostrar a confiança | `scratch/entrevista-real/saida.txt` |
| A5 | O INTENT sai `DRAFT`, 0%, com "Não informado na entrevista" em 7 campos, e **mesmo assim** a SPEC e o PLAN são gerados, violando o gate do próprio INTENT | INTENT gerado, linhas 8, 14–57 e 74 |
| A6 | Tríade incoerente: "Épico 01 — Ingestão de Rota" contém "TICK-01.2 — Aplicar a skill aidd-site-dominator" | PLAN gerado |
| A7 | A skill gerada é uma casca: "Aplique a skill X" cinco vezes, encadeamento linear, gates `file_exists` | `SKILL.md` e `.loop.json` gerados |
| A8 | O INTENT gerado traz blocos fixos do próprio LegoLoop (`go build -o bin/legoloop.exe`) qualquer que seja o processo | INTENT gerado, linha 44 |
| A9 | O stop-gate é cego por construção: `Result` não traz espelho, confiança nem conteúdo da Tríade | `engine.go:50-58` |
| A10 | `--docs X` grava em `X/docs/...` (semântica de "raiz do projeto" com nome de "pasta de docs") | árvore de `scratch/entrevista-real/sandbox` |
| A11 | A skill `legoloop-fabricar` não está no cofre nem em `~/.claude/skills`; só em `skills/` do repositório | `ls` |

## 3. Testes verdes que não provam o que dizem

* `TestRunFabricarInteractive` (14 casos): injeta `Deps.Stdin`; nenhum teste roda o `main` real. Não pegou A1.
* `TestFabricarSkill`: confere frontmatter e ausência de `canonize`; não confere que a skill delegue a entrevista ao motor. Não pegou A2.
* Testes do clarificador e do `--yes`: aceitam confiança 0 como comportamento correto do modo autônomo. Não pegaram A5.
* `fabricar_api_test.go` / `fabricar_modal_test.go`: conferem o payload `{name, goal, category, model}`; não conferem que a entrevista chegue ao INTENT. Não pegaram A3.
* Testes da fusão: conferem ordenação estável e exclusões; nenhum confere relevância ao objetivo. Não pegaram A6/A7.
* Nenhum teste avalia o **conteúdo** gerado com LLM real; o smoke usou 9Router falso.

## 4. Causas no método (o que eu fiz errado)

1. **A SPEC perdeu requisitos do INTENT** ("semântica" virou ordenação por `Score`; o espelho da entrevista não tem campo no `Request`) e eu implementei a SPEC, não o INTENT.
2. **ACs de forma, não de comportamento**: cada AC podia ficar verde com a função existindo e o resultado vazio.
3. **Nenhum ponta a ponta pelo binário real em cada porta** antes de declarar pronto.
4. **Limites tratados como aceitáveis**: "modo autônomo gera 0%" e "interativo provado só com LLM falso" foram documentados em vez de bloquear o fechamento.
5. **Tríade montada por duas fontes que não se falam**: SPEC pelo LLM, PLAN/ROUTE pelas peças da fusão.

## 5. O que está sólido e deve ser preservado

Gravação atômica com rollback, sanitização do nome, default-deny de MCP/hook/script/loop varridos, mídia removível `Untrusted`, chave do 9Router só para a URL confiável, `"stream": false`, repetição limitada da SPEC, rota do Studio com `secure`/`TryLock`/`maskResult`, plotagem no Canvas com confirmação.
