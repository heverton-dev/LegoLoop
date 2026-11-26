---
id: INTENT-intv_1795554379812376300
status: APPROVED
date: 2026-11-24
schema_version: 2
upstream: [INTERVIEW-intv_1795554379812376300]
covers: [Q1]
---

# 🎯 INTENT: Criar o Processo canonizado legoloop-diorama-ativos (Ciclo 2 — Etapas 3, 4 e 5)

> **Status:** APPROVED  
> **Data:** 2026-11-24  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática A2UI)  
> **Alvo:** `internal/diorama/ativos`, `internal/cli`, `processos/legoloop-diorama-ativos/`, `skills/legoloop-diorama-ativos/` (detalhado na SPEC)  
> **Sessão:** intv_1795554379812376300  
> **Sessão A2UI:** sess_3  
> **Confiança da Entrevista:** 98% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criar o Processo canonizado `legoloop-diorama-ativos` (Ciclo 2 do Macroprocesso `legoloop-video-diorama`), cobrindo as Etapas 3, 4 e 5: fichas de elenco e recorte de sprites (`rembg` e despill verde) com descritores JSON estritos, cenários e objetos nos modos `codigo` (SVG 100% autoral), `flow` (fila resiliente com backoff, referências `--ref`, trava `flow_projeto` e idempotência por hash) e `hibrido`, além da composição de storyboard com compilação determinística de `cenas.json` cobrindo 100% das linhas de `script.json`.
* **Quem Sofre:** O operador do canal (recebe elenco consistente, sprites recortados e storyboard auditável sem intervenção manual repetitiva), o Ciclo 3 (`legoloop-diorama-render`, que consome `flow/`, `sprites/` e `cenas.json` validados como contrato imutável) e o Macroprocesso `legoloop-video-diorama`.
* **O Porquê:** O Ciclo 1 (`legoloop-diorama-roteiro`) já está entregue e canonizado. Sem os ativos visuais validados e o contrato `cenas.json` do Ciclo 2, o Ciclo 3 (renderização, animação e timeline) não pode consumir dados reais nem ser integrado ao Macroprocesso.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Apenas o Ciclo 1 (`legoloop-diorama-roteiro`) está canonizado no LegoLoop, produzindo `projeto.json`, `marca.json`, `script.json` e `plano.md`. Não existe motor determinístico para elenco, sprites, cenários/objetos ou compilação de `cenas.json`.
  * **Amanhã:** Processo canonizado `legoloop-diorama-ativos` (pacote `internal/diorama/ativos` + CLI `legoloop diorama-ativos` + DAG 5 de 5) que consome `projeto.json`, `marca.json`, `script.json` e `plano.md` do Ciclo 1 e executa, em ordem:
    1. **Pré-voo do Ciclo 2:** valida contratos do Ciclo 1, permissões de leitura e escrita em `flow/` e `sprites/`, e pré-requisitos do modo visual escolhido (`codigo`, `flow` ou `hibrido`).
    2. **Etapa 3 (Elenco e Sprites):** gera e valida 3 fichas por personagem (`turnaround`, `expressoes_3x2`, `poses_6x`) com `descricao_imutavel`, executa recorte transparente via `rembg` + despill verde e fatiamento automático em `sprites/*.json` (`id`, `personagem_id`, `tipo`, `arquivo_png`, `largura`, `altura`, `alfa_valido`, `bounding_box` sem sobreposição). Gate bloqueia ficha ausente, alfa inválido, dimensão abaixo do mínimo ou sobreposição.
    3. **Etapa 4 (Cenários e Objetos):** modo `codigo` gera cenários e objetos SVG 100% autorais sem dependência externa; modo `flow` exige trava de UUID `flow_projeto`, suporte a `--ref`, fila resiliente com retry/backoff e idempotência por hash (não regera ativo já válido); modo `hibrido` combina ambos. Gate exige mínimo de 1 cenário distinto a cada 2 linhas.
    4. **Etapa 5 (Storyboard):** compila `cenas.json` cobrindo 100% das linhas de `script.json`, com metáfora visual e `>= 4` elementos móveis por linha.
    5. **Governança de Gates:** respeita `politica_gates` (`humano`, `mecanico` ou `hibrido`) herdada de `projeto.json`, parando nos Stop-Gates de elenco, quadros-teste e storyboard quando configurados como humanos.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Implementar o pacote `internal/diorama/ativos` e o comando CLI `legoloop diorama-ativos` cobrindo Pré-voo do Ciclo 2, Etapa 3 (Elenco e Sprites), Etapa 4 (Cenários e Objetos nos modos `codigo`, `flow` e `hibrido`) e Etapa 5 (Storyboard e `cenas.json`).
* [ ] Objetivo 2: Garantir contratos determinísticos e gates mecânicos de bloqueio:
  * Gate Etapa 3: esquema estrito de personagem (`id`, `nome`, `descricao_imutavel`, 3 fichas obrigatórias) e `sprites/*.json` (`id`, `personagem_id`, `tipo`, `arquivo_png`, `largura`, `altura`, `alfa_valido`, `bounding_box`), bloqueando ficha ausente, canal alfa inválido, dimensão abaixo do mínimo ou sobreposição de `bounding_box`.
  * Gate Etapa 4: densidade `>= 1` cenário distinto a cada 2 linhas, trava `flow_projeto` obrigatória no modo `flow`/`hibrido`, suporte a `--ref`, retry/backoff na fila e idempotência por hash (reexecução não regera ativo válido).
  * Gate Etapa 5: `cenas.json` cobrindo 100% das linhas de `script.json`, com metáfora visual e `>= 4` elementos móveis por linha.
  * Política de gates herdada de `projeto.json` (`humano`, `mecanico` ou `hibrido`) nos Stop-Gates de elenco, quadros-teste e storyboard.
* [ ] Objetivo 3: Cobertura de testes `>= 95%` em `internal/diorama/ativos`, testes 100% herméticos (adaptadores externos como `rembg` e Flow injetáveis via interface), zero stubs, gates da Espinha Dorsal verdes e fechamento completo 5 de 5 via `legoloop fechar`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON do Ciclo 1 (`projeto.json`, `marca.json`, `script.json`, `plano.md`).
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Etapas 6, 7, 8 e 8b (voz neural, Whisper, motor diorama HTML/SVG de animação, `timeline.json`, mixagem LUFS, render MP4/SRT e Shorts — pertencem ao Ciclo 3).
* ⛔ Etapas 9 a 13 (revisor estéril, vitrine, acervo/Postiz, comentários e aprendizado — pertencem ao Ciclo 4).
* ⛔ Composição ponta a ponta `P1->P4` (pertence ao Ciclo 5).
* ⛔ Painel web paralelo fora do Looping Studio.
* ⛔ Código sem licença de `RobbiSkill/`.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills e Peças que o futuro Processo precisará acionar.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `legoloop-especificacao` | Elaborar a especificação formal (`SPEC`) com critérios BDD binários. |
| `PECA-02` | `skill` | `legoloop-planejamento` | Elaborar `PLAN`, `ROUTE` e `DAG` de execução. |
| `PECA-03` | `unit_task` | `internal/diorama/ativos` + `internal/cli` | Implementação via TDD Red-Green com cobertura `>= 95%`. |
| `PECA-04` | `script` | `go test ./...` + `cobertura` + `espinha verificar` | Regressão mecânica, piso de cobertura e integridade. |
| `PECA-05` | `stop_gate` | `legoloop fechar "legoloop-diorama-ativos"` | Fechamento 5 de 5 (manual, livro, skill, DAG, build/install, governança viva e relatório). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (ou encadeado pelo Ciclo 5 `legoloop-video-diorama`).
* **Nível de Risco & Blast Radius:** Baixo (opera isoladamente na pasta `<projeto>/` do vídeo e no pacote novo `internal/diorama/ativos`).
* **Política de Privacidade & Orçamento Cognitivo:** Testes 100% herméticos e offline; chamadas Flow só via interface injetável com trava `flow_projeto` e idempotência por hash.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  * Projeto válido vindo do Ciclo 1 percorre Etapas 3, 4 e 5 gerando `flow/`, `sprites/*.json` e `cenas.json` aprovados em todos os gates mecânicos.
  * Gates das Etapas 3, 4 e 5 bloqueiam todas as violações de contrato declaradas.
  * Cobertura de testes `>= 95%` em `internal/diorama/ativos` (`go test -cover ./...`).
  * `legoloop espinha verificar`, `funclen`, `taxlinter`, `gofmt` e `legoloop fechar` verdes.
* **Restrições Inegociáveis:**
  * Zero código sem licença de `RobbiSkill`: SVGs do modo `codigo` são 100% autorais.
  * Contrato de entrada imutável: consome `projeto.json`, `marca.json`, `script.json` e `plano.md` sem alterar o schema do Ciclo 1.
  * Trava `flow_projeto` obrigatória no modo `flow`/`hibrido` e idempotência por hash.
  * Política de gates herdada de `projeto.json` (`humano`, `mecanico` ou `hibrido`).
  * Testes 100% herméticos com cobertura `>= 95%` (adaptadores externos injetáveis via interface).
  * Fechamento completo 5 de 5 via `legoloop fechar`.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * O campo `estilo_visual` de `projeto.json` define o modo visual (`codigo`, `flow` ou `hibrido`), com `codigo` como fallback determinístico quando aplicável ou bloqueio se inválido.
* **Riscos e Perguntas em Aberto:** Confiança de 98%: validar a INTENT no Stop-Gate 01 antes de gerar a SPEC.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura de testes `>= 95%` em `internal/diorama/ativos` (TDD Red-Green). | Medição de cobertura do pacote | `go test -cover ./internal/diorama/ativos/...` (`>= 95%`) |
| 3 | Gates mecânicos das Etapas 3, 4 e 5 bloqueiam ficha ausente, alfa inválido, sobreposição de `bounding_box`, densidade `< 1` cenário a cada 2 linhas, violação de `flow_projeto`, cobertura `< 100%` das linhas, falta de metáfora visual ou `< 4` elementos móveis por linha. | Casos negativos e positivos em teste unitário e de integração | A detalhar na SPEC como `AC-XX` binário |
| 4 | Fechamento completo 5 de 5 sem lacunas. | Verificação da Espinha Dorsal e `legoloop fechar` | `legoloop espinha verificar` e `legoloop fechar` com exit code 0 |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md`](../specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) na Interface A2UI ou chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
