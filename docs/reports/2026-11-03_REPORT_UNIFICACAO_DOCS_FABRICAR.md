# 📊 REPORT: Unificação dos Documentos de 2026-11-03 (legoloop fabricar)

> **Data:** 2026-11-03  
> **Escopo:** `docs/intents`, `docs/specs`, `docs/plans`, `docs/routes` (somente os documentos de hoje)  
> **Resultado:** 7 arquivos (com V1 e V2) viraram 4 documentos canônicos, regenerados no formato dos templates v2.0.0.  

---

## 1. De ➔ Para

| Antes | Depois |
|:---|:---|
| `intents/..._INTENT_LEGOLOOP_CANONIZER.md` | `intents/2026-11-03_INTENT_LEGOLOOP_FABRICAR.md` |
| `specs/..._SPEC_LEGOLOOP_CANONIZER.md` (V1) e `..._CANONIZER_V2.md` | `specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md` |
| `plans/..._PLAN_LEGOLOOP_FABRICAR.md` (V1) e `..._FABRICAR_V2.md` | `plans/2026-11-03_PLAN_LEGOLOOP_FABRICAR.md` |
| `routes/..._ROUTE_LEGOLOOP_FABRICAR.json` (V1) e `..._FABRICAR_V2.json` | `routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json` |

Os 7 originais estão em `docs/archive/2026-11-03/` (cópias idênticas, conferidas com `cmp` antes de qualquer remoção). Eles **não estavam no git**, então esse arquivo é a única cópia deles.

## 2. Decisões aplicadas

1. **Nome:** `fabricar` (escolha do autor). Os quatro documentos usam o mesmo nome.
2. **`legoloop run`:** permanece por caminho; a forma `legoloop run legoloop-<nome>` da INTENT original foi descartada (recomendação técnica: exigiria alterar o `run`).
3. **Skill nos harnesses:** só no cofre; link apenas com `--link-harness`.
4. **`.dag.json`:** em `~/.aidd/canonical/loops/`, sem diretório novo no cofre.
5. **Confiança da entrevista:** 95%, informada pelo autor **com incerteza**. A INTENT registra isso em "Riscos e Perguntas em Aberto" e deve ser confirmada no Stop-Gate 01.

## 3. O que mudou em cada documento

* **INTENT:** passou ao `TEMPLATE_INTENT` (6 seções, cabeçalho com confiança e template, tabela de critérios com prova, Stop-Gate 01). Conteúdo e objetivos preservados; `canonize` virou `fabricar` e `run` por nome virou por caminho. Status `UNDER_REVIEW`.
* **SPEC (V1 + V2 ➔ única):** a V2 era a base. Aplicado o template: cabeçalho com `Consome/Produz/Template`, diagrama em **Mermaid** (era ASCII), costuras no formato do template (primária, secundárias, dublês permitidos), `AC-E02.3` virou "Comando Próprio e Legado Intacto", e as 4 decisões viraram a premissa 9.
* **PLAN (V1 + V2 ➔ único):** a V2 era a base e nenhum ticket foi perdido (17 tickets, `Escopo` incluído). Reestruturado nas 4 seções do template (a V2 tinha seções 0 a 5); o registro de revisão V1➔V2 saiu do PLAN e está nos arquivos arquivados; IDs `FAB-*` viraram `TICK-NN.M`; a homologação virou o Épico 04 mais `SG-FINAL`.
* **ROUTE (V1 + V2 ➔ único):** a V2 já tinha as chaves do template. Mudaram só IDs, referências e nomes. 21 nós, 4 Stop-Gates humanos.

## 4. Defeitos encontrados na hora de unificar

| Onde | Defeito | Correção |
|:---|:---|:---|
| SPEC V2 (escrita por mim) | `Request.Validate()` tinha corpo comentado e **nem compilaria**; `Deps`, `Scanner`, `FS`, `ScanOptions` e `ScanResult` eram usados nas costuras sem estarem definidos | Contrato Go completo, com `Validate()` real e todos os tipos definidos |
| PLAN V2 | Seções numeradas 0 a 5 (template: 1 a 4); dependia de "patch V2.1 da SPEC" e de `docs/templates/` "não versionado" | Seções do template; a SPEC unificada já incorpora o patch; os templates agora são embutidos no binário |
| PLAN V2 | Cobria o template apenas "por semelhança" | Verificado por script (seção 5) |
| `TEMPLATE_PLAN` | Não tinha lugar para o `Escopo` de cada ticket, que carrega o detalhe do que o teste prova | Adicionado como campo **opcional** (só no PLAN) e documentado no README dos templates |

## 5. Verificação (script descartável, pacotes reais)

* Seções numeradas de INTENT, SPEC e PLAN na ordem do template; nenhum `{{...}}` sobrando; linha `Template: v2.0.0` presente.
* ROUTE: mesmo conjunto de chaves do `TEMPLATE_ROUTE.json` em todos os níveis (via `templates.RouteKeys()`); todo `output.gate` compila em `gate.Compile`; subagentes sempre `false`; ids únicos, dependências existentes, sem ciclo; `SG-FINAL` por último.
* PLAN espelha ROUTE: 17 tickets nos dois, com natureza, timeout, teto de tokens, modelo, `spec_ref`, `verify`, gate e dependências idênticos.
* Cobertura: os 12 critérios `AC-*` da SPEC são exatamente a união dos `spec_ref` do ROUTE, e as 17 âncoras `arquivo#AC-*` resolvem.
* Todos os links relativos e as referências `intent_ref/spec_ref/plan_ref` apontam para arquivos existentes; nenhuma referência a `_V2`, `CANONIZER_V2` ou `file:///`.

## 6. Pendências (estado final)

Nenhuma decisão em aberto. As três pendências da primeira rodada foram resolvidas em 2026-11-03:

1. **Nomenclatura:** resolvida (seção 7).
2. **Duplicação no ROUTE:** resolvida pela recomendação técnica (seção 8).
3. **Confiança de 95%:** confirmada pelo autor; consta como premissa na INTENT.

Pontos que continuam valendo:
* **Critérios `AC-NF`:** o template prevê esse prefixo para critérios não funcionais. A SPEC mantém `AC-E01.5` (segurança e robustez) para preservar a rastreabilidade dos tickets.
* **Nada foi commitado.** Os documentos novos, o arquivo em `docs/archive/` e o `TEMPLATE_PLAN.md` alterado estão no working tree.
* **Nada foi construído.** Os 4 documentos definem o que será feito; o BUILD só começa depois da aprovação dos Stop-Gates (INTENT, SPEC e PLAN).

## 7. Rodada 2: `canonizer` ➔ `fabricar` (nomenclatura)

**Verificação no repositório (antes de renomear):**
* Não existia nada chamado `canonizer` ou `fabricar` no código: nenhum pacote, skill (nem no repositório nem em `~/.claude/skills`, `~/.aidd/canonical/skills`), rota ou arquivo. Os nomes só existiam nos documentos de hoje.
* O "canonizar" real do repositório é **outra função**: desduplicar arquivos por SHA-256, eleger a versão canônica e trocar as cópias por links para o cofre (`internal/canon`, `legoloop canonize`, `/api/canonize/*`, e o `import` que canoniza skills do GitHub).
* Fabricar é criar um processo novo (entrevista, curadoria multi-drive, fusão, tríade e template). Usar `canonizer` para isso colidiria com `canon` e enganaria quem lê o código.

**Renomeado nos 4 documentos (INTENT, SPEC, PLAN, ROUTE):**

| Antes | Depois |
|:---|:---|
| pacote `internal/canonizer` | `internal/fabricar` |
| função `canonizer.Canonize` | `fabricar.Fabricar` |
| skill `legoloop-canonizer` (`skills/legoloop-canonizer/SKILL.md`, `/legoloop-canonizer`) | `legoloop-fabricar` |
| rota `POST /api/canonizer/run` | `POST /api/fabricar/run` |
| campo `studio.Server.Canonizer` | `studio.Server.Fabricar` |
| arquivos e testes (`canonizer_api.go`, `TestHandleCanonizerRun`, `TestCanonizerSkill`...) | `fabricar_api.go`, `TestHandleFabricarRun`, `TestFabricarSkill`... |
| ids do modal (`modal-canonizer`, `canonizer-name`...) | `modal-fabricar`, `fabricar-name`... |
| "Modal de Canonização", "API de Canonização" | "Modal de Fabricação", "API de Fabricação" |

**Preservado de propósito** (pertence à função legada): `legoloop canonize --root`, `/api/canonize/*`, `internal/canon`, `canon.Linker`, `handleCanonizePiece` (citado só como padrão de código a seguir) e `CA-CANON-01 a 06`.

Como nada foi construído, "criar o pacote, a skill e a rota" significou corrigir esses nomes nos documentos que definem a construção. Os artefatos em si nascem no BUILD (`TICK-01.1`, `TICK-02.4`, `TICK-03.1`).

## 8. Rodada 2: reuso do `plangen` (ex-duplicação no `TICK-01.5`)

Recomendação técnica adotada: a fábrica **não** tem emissor de ROUTE próprio; ela reusa o `plangen`, que já gera ROUTE e PLAN conformes ao template v2.

* **Novo `TICK-01.5`:** exporta `plangen.Render(baseDir, specPath, tickets, level, now)`, que monta PLAN e ROUTE sem LLM e com a data injetada (determinismo). `Generate` passa a ser "interpretar tickets" mais `Render`, sem mudar comportamento nem schema.
* **Renumeração:** o Épico 01 passou de 7 para 8 tickets (`TICK-01.5` novo; os antigos 01.5, 01.6 e 01.7 viraram 01.6, 01.7 e 01.8).
* **Emissores (`TICK-01.6`):** perdem o `emit_route.go`; ficam só SKILL, `.loop.json` e `.dag.json`.
* **Orquestrador (`TICK-01.8`):** passa a depender de 01.5, 01.6 e 01.7; prova que o ROUTE gerado tem as chaves do `TEMPLATE_ROUTE` v2.
* **Alvo:** INTENT e SPEC incluem `internal/plangen/`; a INTENT ganhou o `BLOCO-07` correspondente (o `stop_gate` de build virou `BLOCO-08`).

## 9. Verificação da rodada 2

Mesmo verificador da primeira rodada, com checagens novas. Resultado: 0 falhas.
* 4 documentos nas seções do template; ROUTE com as chaves do template; gates compilam; sem ciclo.
* PLAN espelha ROUTE: 18 tickets e 22 nós (4 Stop-Gates humanos); tabela de épicos do PLAN confere com a contagem real (E01: 8, E02: 4, E03: 4, E04: 2).
* 12 critérios `AC-*` da SPEC cobertos exatamente pelos tickets; âncoras resolvem; links e referências cruzadas existem.
* Nenhum `canonizer` restante; nenhum `emit_route.go` como artefato; as únicas menções a `canonize` são as do comando legado.

Incidente durante a rodada (corrigido): um script de edição meu vazou mensagens para o início da INTENT e da SPEC (`ok: ...` na primeira linha). Foi detectado antes de qualquer uso, reparado e coberto por uma checagem permanente do verificador.