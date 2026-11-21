---
id: INTENT-intv_1795277863681062100
status: APPROVED
date: 2026-11-21
schema_version: 2
upstream: [INTERVIEW-intv_1795277863681062100]
covers: [Q1]
---

# 🎯 INTENT: Iniciativa 4 (I04) — Espinha dorsal previsível, universal e determinística

> **Status:** APPROVED  
> **Data:** 2026-11-21  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** `internal/espinha` (novo), `internal/descoberta` (novo), `internal/dag`, `internal/runner`, `internal/fabricar`, `internal/matriz`, `internal/caminhos`, `internal/paths`, `tools/mutar` (novo), `.github/workflows`  
> **Sessão:** intv_1795277863681062100  
> **Confiança da Entrevista:** 96% reportada pelo motor com o modelo `sonnet`, depois de o motor com `code-fast` parar em 85%. As decisões foram tomadas pelo usuário uma a uma ("ACEITO", "[1]"); o Espelho abaixo foi conferido contra elas.  
> **Insumo:** `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md` (14 lacunas).  
> **Macroprojeto:** Distribuição e Forja do LegoLoop. As INTENTs da I01, da I02 e da I03 não são alteradas; as mudanças entram por Emenda (`docs/intents/2026-11-21_EMENDA_2_MACROPLANO.md`).  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** a espinha dorsal do LegoLoop (verbos, Processos, motores e requisitos) nasceu sem arquitetura prévia e acumulou dívida. A auditoria de 2026-11-21 achou 14 lacunas. As duas críticas: a Matriz de peças está defasada e nada a atualiza (505 novas, 2 alteradas, 6.146 removidas), e o `fabricar` não varre drives, embora a premissa do produto seja usar as peças do SO e dos discos do usuário. As altas: `isolation.executor` é decorativo, o vocabulário de tipo de nó está duplicado (`unit_no` e `unit_task`) e existem duas fachadas de caminhos (`caminhos` e `paths`).
* **Quem Sofre:** quem mantém e evolui o LegoLoop, o usuário nos três SOs e as iniciativas I01, I02 e I03, que se apoiam na espinha.
* **O Porquê (por que agora):** a I01, a I02 e a I03 constroem sobre vocabulário, fachadas e descoberta que a auditoria mostrou falhos. Corrigir depois obrigaria a refazer o que a outra sessão acabou de construir. Por isso a I04 executa primeiro e as três ficam congeladas.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** a leitura e a geração são determinísticas e medidas (6 verbos de leitura com hashes idênticos em duas execuções; 33 de 33 Processos validam), mas a descoberta de peças é manual e defasada, o vocabulário não é amarrado, e a universalidade só foi compilada para Linux e macOS, nunca executada.
  * **Amanhã:** um Registro da Espinha conferido por gate mecânico, descoberta automática por hash de conteúdo, vocabulário amarrado por uma tabela única, CI verde nos três SOs e mutação cobrindo as regras e o código da espinha.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: **Registro da Espinha.** Arquivo único e versionado com quatro listas: verbos (camada espinha ou domínio, motor e contrato), motores (pacote, interface e se é determinístico), Processos (nascimento 5 de 5) e requisitos (cada lei do `AGENTS.md` ligada ao verificador mecânico). Comando `legoloop espinha verificar`, que roda no CI e no `legoloop fechar`.
* [ ] Objetivo 2: **Descoberta de peças.** Motor determinístico e somente leitura que procura no Cofre, no HOME e nas unidades fixas do SO (Windows pelo tipo de unidade, Linux em `/mnt`, macOS em `/Volumes`). Unidade removível só com `TrustRemovable`. O `TargetDrives` passa a valer. A Matriz se atualiza por hash de conteúdo, não por relógio, antes de o `fabricar` montar o Processo.
* [ ] Objetivo 3: **Vocabulário amarrado.** O executor é derivado do tipo do nó e validado (`stop_gate` é `local_process`; `script`, `skill`, `unit_no`, `mcp` e `subprocesso` são `orchestrix_hermetic`; `native_swarm` só explícito). Uma tabela única liga `unit_task` da ROUTE a `unit_no` do DAG. Saem o alias `TypeSubloop` e o literal `dag.NodeType("unit_no")` do `runner`. Os DAGs com rótulo errado ou ausente são corrigidos.
* [ ] Objetivo 4: **Fachada única de caminhos.** `paths` é absorvido por `caminhos` (35 importadores contra 15), e o pacote duplicado sai.
* [ ] Objetivo 5: **Universalidade provada.** CI em Windows, Linux e macOS, trazido da I01 (Épico 02), rodando `go test ./...`, `go vet`, `funclen` e `legoloop espinha verificar`. Correções que o CI exigir: caminho de fontes por SO, `python3` onde o SO só o tem, escolha entre `powershell` e `cmd` por SO.
* [ ] Objetivo 6: **Mutação em duas camadas.** (a) Cada regra do Registro tem um mutante que cria o defeito numa cópia temporária, e o gate exige 100% de mutantes mortos, cada um com exit 2. (b) Mutador de código `tools/mutar`, só com `go/ast`, que troca operadores, inverte condições e remove `return err` nos pacotes da espinha (`dag`, `runner`, `gate`, `plangen`, `triagem`, `upstream`, `interview`); mutante sobrevivente reprova o gate, salvo justificativa registrada como anticorpo.
* [ ] Objetivo 7: **As 14 lacunas e a medição que faltou.** Cada lacuna tem teste Red e anticorpo. A rodada de medição inclui `go test -race` e a execução de ponta a ponta (`run` e `exec` com DAG só de scripts) duas vezes, com relatório idêntico.
* [ ] Objetivo 8: **Revisão das iniciativas congeladas.** A Emenda 2 revisa a I01 (sai o Épico 02, que o CI da I04 entrega; entram G01, G02 e G05 como pré-requisito), a I02 (vocabulário) e a I03 (a validação passa a ser coerência com o tipo).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Refatorar os verbos de domínio (`video-motion`, `video-rea`, `canal`, `cartoon-shorts`, `web-qa`, `legopress`) além de classificá-los no Registro.
* ⛔ Mudar o formato do DAG.
* ⛔ Alterar as INTENTs da I01, da I02 e da I03.
* ⛔ Tornar determinístico o texto que o modelo de linguagem devolve. O determinismo vale para as decisões da espinha (validar, rotear, planejar, despachar), não para a saída do modelo.
* ⛔ Dependências externas novas além das já autorizadas.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_task` | `internal/espinha` e `legoloop espinha verificar` | Registro e gate da espinha. |
| `PECA-02` | `unit_task` | `internal/descoberta`, `internal/fabricar`, `internal/matriz` | Descoberta e Matriz por hash. |
| `PECA-03` | `unit_task` | `internal/dag`, `internal/runner`, `tools/route2dag` | Vocabulário e executor derivado do tipo. |
| `PECA-04` | `unit_task` | `internal/caminhos`, `internal/paths` | Fachada única de caminhos. |
| `PECA-05` | `unit_task` | `tools/mutar` | Mutador de código com `go/ast`. |
| `PECA-06` | `script` | `.github/workflows` | CI nos três SOs. |
| `PECA-07` | `script` | `go test ./...` | Regressão mecânica e medição. |
| `PECA-08` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** alto. Mexe no `dag`, no `runner`, em dois pacotes de caminhos usados por 50 arquivos e na descoberta que lê discos. A I04 executa antes da I01, da I02 e da I03.
* **Política de Privacidade & Orçamento Cognitivo:** a descoberta é somente leitura, não executa o que encontra e relata o que leu e o que recusou. Nenhum LLM na decisão da espinha. Nenhum segredo é lido nem gravado.
* **Contrato de Evidência Factual (Definition of Done):** as cinco provas da seção 5 e o portão completo (`go test ./...`, `go vet ./...`, `gofmt -l` vazio, `funclen`, `taxlinter`, `go.mod` intacto).
* **Restrições Inegociáveis:**
  * TDD com Red primeiro, e fixture com defeito sai com exit 2 (a mesma convenção da trava `upstream_ok`).
  * A descoberta é determinística e somente leitura.
  * A Matriz se atualiza por hash de conteúdo, nunca por relógio.
  * O executor é derivado do tipo do nó, e `native_swarm` só entra de forma explícita.
  * Existe uma única tabela de ligação `unit_task` e `unit_no`, e uma única fachada de caminhos.
  * A verificação da espinha roda no CI e no `fechar`.
  * INTENTs aprovadas da I01, da I02 e da I03 intactas.
  * Zero débito: todo achado se corrige na execução e vira anticorpo.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** o GitHub Actions está disponível para o repositório (a I01 já o assumia); o mutador sobre os sete pacotes termina em tempo aceitável no CI quando roda só nos pacotes alterados.
* **Riscos e Perguntas em Aberto:** como tratar os Processos fora do prefixo `legoloop-` (10 de 33) no Registro; custo de tempo do mutador completo; como a descoberta se comporta em unidades de rede (a SPEC define).

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | A espinha se confere sozinha: sai com 0 no repositório e com exit 2 em fixtures com defeito (verbo sem registro, lei sem verificador, dois pacotes para a mesma função, Processo incompleto, executor incoerente, tipo de nó fora da tabela). | Testes de mutação das regras, 100% dos mutantes mortos | `legoloop espinha verificar` e `go test ./internal/espinha/...` |
| 2 | A descoberta acha uma peça nova sem `matriz gerar`; `matriz status` mostra zero diferença depois do `fabricar`; `TargetDrives` restringe a busca. | Teste com raiz de teste e peça colocada nela | `go test ./internal/descoberta/... ./internal/fabricar/...` |
| 3 | O executor é derivado do tipo; os 33 Processos e os DAGs de `docs/dags` validam; nenhum `stop_gate` com rótulo errado. | Testes do `dag` e `legoloop validate` em todos os DAGs | `go test ./internal/dag/...` |
| 4 | O CI fica verde em Windows, Linux e macOS. | Execução do workflow | `gh run list --workflow ci.yml --limit 1 --json conclusion` |
| 5 | As 14 lacunas têm teste Red e anticorpo; `go test -race` passa; a execução de ponta a ponta duas vezes gera relatório idêntico; nenhum mutante de código sobrevive sem justificativa. | Relatório de fechamento da I04 | `go test -race ./...`, `tools/mutar`, `legoloop anticorpos status` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md`](../specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat, antes de qualquer linha de especificação ser redigida. A trava `upstream_ok` exige `status: APPROVED` no cabeçalho.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
