---
name: legoloop-diorama-ativos
author: LegoLoop
description: Use quando precisar executar o Processo legoloop-diorama-ativos (Ciclo 2 do Macroprocesso legoloop-video-diorama, Etapas 3, 4 e 5): fichas de elenco e recorte de sprites (rembg + despill verde), cenários nos modos codigo (SVG 100% autoral), flow (fila resiliente com backoff, --ref, trava flow_projeto e idempotência por hash) e hibrido, e compilação determinística de cenas.json. NÃO use para roteiro/pautas (Ciclo 1), renderização (Ciclo 3) ou publicação (Ciclo 4).
consumes: Artifact:ValidatedScript
produces: Artifact:DioramaAssets
artifact: cenas.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🎨 legoloop-diorama-ativos — Elenco, Sprites, Cenários e Storyboard

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Motor oficial:** `legoloop diorama-ativos` (implementação em `internal/diorama/ativos` e `internal/cli/diorama_ativos.go`).  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** após o Ciclo 1 (`legoloop-diorama-roteiro`) ter produzido `projeto.json`, `marca.json`, `script.json` e `plano.md` válidos, para fabricar e validar os ativos visuais das Etapas 3, 4 e 5 (`flow/`, `sprites/*.json` e `cenas.json`).  
**Quando NÃO usar:** para descoberta de pauta/roteiro (Ciclo 1), voz/timeline/renderização (Ciclo 3) ou revisão/publicação (Ciclo 4).

* **Consumes (Entradas Obrigatórias):** `projeto.json`, `marca.json`, `script.json` e `plano.md`.
* **Produces (Saídas Mecânicas):** `flow/elenco.json`, `flow/cenarios.json`, `flow/*.svg`, `sprites/*.json` e `cenas.json`.

---

## 2. Procedimento Determinístico e Nós do DAG

0. **Escolha da superfície e via de execução:** ofereça **Interface A2UI** (`sessao_criar` -> `sessao_publicar` -> `sessao_aguardar` no mesmo turno) ou **Chat/CLI direto**. Quando conduzido por agente com MCP ativo, priorize o servidor `flow-mcp` (`flow_generate`/`flow_wait`) para o Google Flow; na execução direta via CLI/DAG/Studio, o motor usa o driver embutido (`flow_driver.py` via CDP `9222`/`9333`).
1. **Pré-voo (`n01_prevoo_ciclo2`):** rode `legoloop diorama-ativos prevoo --projeto <DIR> --json` (valida contratos do Ciclo 1, Path Jail, escrita em `flow/` e `sprites/`, `estilo_visual` e trava `flow_projeto`).
2. **Etapas 3, 4 e 5 (`n02_executar_etapas_3_4_5`):** rode `legoloop diorama-ativos executar --projeto <DIR> [--aprovar g1,g2] --json` (respeita `politica_gates` herdada de `projeto.json` nos Stop-Gates `etapa3_elenco`, `etapa4_quadros_teste` e `etapa5_storyboard`).
3. **Gate final (`n03_validar_storyboard`):** rode `legoloop diorama-ativos validar-storyboard --projeto <DIR> --json` (exige cobertura de 100% das linhas, cenário existente, metáfora visual e `>= 4` elementos móveis por linha).

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

1. **NUNCA** use código sem licença de `RobbiSkill/`: SVGs do modo `codigo` são 100% autorais.
2. **NUNCA** execute modo `flow` ou `hibrido` sem trava `flow_projeto` válida.
3. **NUNCA** regere ativo válido quando o hash SHA-256 de entrada for idêntico.
4. **NUNCA** aceite sprite sem `alfa_valido`, sem `despill_verde`, menor que `64px` ou com `bounding_box` sobreposta.
5. **NUNCA** altere o schema de entrada do Ciclo 1 (`projeto.json`, `marca.json`, `script.json`).

---

## 4. Comandos de Verificação & Testes de Terminal

```bash
legoloop diorama-ativos --help
legoloop validate processos/legoloop-diorama-ativos/*.dag.json
go test -count=1 -cover ./internal/diorama/ativos/...
```

---

## 5. Critério de Parada (Stopping Checklist)

- [ ] Pré-voo do Ciclo 2 sem itens `FALTA` (exit code 0).
- [ ] `flow/elenco.json`, `sprites/*.json`, `flow/cenarios.json` e `cenas.json` gravados e válidos.
- [ ] `validar-storyboard` retorna exit code 0 com 100% das linhas cobertas.
