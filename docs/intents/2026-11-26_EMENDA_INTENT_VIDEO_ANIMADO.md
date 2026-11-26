---
id: EMENDA-INTENT-VIDEO-ANIMADO-20261126
status: APPROVED
date: 2026-11-26
schema_version: 2
emenda_de: [INTENT-a2ui_sess_3_20261125_render]
upstream: [ANALISE_IMPACTO_INVENTARIO_ESTILOS]
---

# EMENDA da INTENT: macro `legoloop-video-animado` (estilo plugável, inventário global, motor portado)

> **Status:** APPROVED (Stop-Gate 01 aprovado em 2026-11-26, resposta "[1]").
> **Emenda de:** `docs/intents/2026-11-25_INTENT_LEGOLOOP_DIORAMA_RENDER.md` (APPROVED em 2026-11-25) e do Mapa Macro `docs/intents/2026-11-24_MAPA_MACRO_VIDEO_DIORAMA.md`.
> **Base de fatos:** `docs/progress/2026-11-25_plan-legoloop-diorama-render/ANALISE_IMPACTO_INVENTARIO_ESTILOS.md` e `achados.md` (A7 a A12).

## 1. Decisões do usuário (2026-11-26, resposta "sugestões")

| # | Decisão | Valor aprovado |
|---|---|---|
| D1 | Nome do macro | `legoloop-video-animado`; processos `legoloop-video-{roteiro,ativos,render,distribuicao}`; pacotes `internal/video/*`; comandos `video-*` com alias dos nomes `diorama-*` por uma versão |
| D2 | Inventário | global do macro, **hierárquico por estilo** (ajuste do usuário em 2026-11-26): `inventario/video/estilos/<estilo>/` guarda `estilo.json` (referências, renderizador, prompt mestre, paleta, status) e os inventários do estilo (personagens, skins, trajes, penteados, adornos, adereços, poses, humores, cenários, objetos); `inventario/video/comum/` guarda o que não depende de estilo (SFX, trilhas, vozes). Cada item pertence a um só estilo. Tag de canal, índice JSON no git e arquivos pesados fora do git (com hash) |
| D3 | Campos de `projeto.json` | `modo_producao` (`codigo`, `flow`, `hibrido`) e `estilo` (id do catálogo de estilos) |
| D4 | Motor | portar o motor original como núcleo da Peça 3, em vez da reescrita autoral |
| D5 | Primeiro estilo novo | palito (esforço baixo); mangá depois |

Contexto dado pelo usuário: o inventário é consultado e alimentado em cada vídeo (pele, traje, penteado, adorno, adereço, pose, humor, cenário, objeto), vale para o macro inteiro, e cada vídeo tem o seu elenco, sempre no estilo escolhido.

## 2. O que muda na INTENT aprovada

### Entra no escopo (antes estava fora ou não existia)
* Alterar os contratos `projeto.json` (D3), `cenas.json` (campos `personagem_id`, `pose`, `humor`, `olhar`, `fala`) e o formato do elenco (`elenco.json` com bloco `cast`). A INTENT de 25/11 proibia mexer nesses contratos; esta emenda revoga a proibição só para esses três pontos.
* Reabrir a Peça 2 (`internal/diorama/ativos`: `elenco.go`, `fichas.go`, `drivers.go`, `schemas.go`, `storyboard.go`) para gerar `cast` a partir do inventário.
* Criar o pacote de inventário e o inventário de estilos (`internal/video/inventario`, `inventario/video/estilos/`): cada estilo tem referências (imagens e vídeos com origem), `status` (`rascunho` ou `validado`) e só aceita itens `aprovado` depois de validado. O estilo atual (molde de pílula) é o primeiro a validar.
* Portar o motor (`engine.html`, `puppet.js`, `layers.js`, `sets.js`, `formato.py`, `legendas.py`, `render.py`, `check_timeline.py`) para `internal/video/motor/`, chamado pelo Go.
* Renomear macro, processos, pacotes e comandos (D1), com alias e migração do registro.

### Continua fora do escopo
* Etapas 9 a 13 (Ciclo 4) e a composição ponta a ponta (Ciclo 5).
* Segundo estilo depois do palito (mangá, pixel art, recorte, 3D): ficam no catálogo como candidatos, sem implementação.
* Reescrever documentos históricos arquivados: recebem só uma nota de redirecionamento.
* Qualquer gasto de crédito do Flow sem cotação e "sim" do usuário.

## 3. Risco que bloqueia a D4: licença do motor

O Mapa Macro (seção 4) e a INTENT de 25/11 proíbem código sem licença da `RobbiSkill`. O motor do projeto Histórias Palitos descende dela: `puppet.js` é igual ao original, e os outros arquivos foram estendidos (`proj--historias-palitos/processo/AUDITORIA.md`, frente 3: "O README não declara licença. Não publique este repositório sem checar com o autor").

Portar o motor como está contradiz essa proibição. Opções:

| Opção | O que é | Custo |
|---|---|---|
| A (sugestão) | Portar o motor para uma pasta local fora do git (`projetos/`), usar e testar agora, e abrir em paralelo o contato com o autor para licença; se negada, executar a opção B | baixo agora |
| B | Reescrever `puppet.js` e `sets.js` do zero com desenho próprio (clean-room), mantendo as extensões nossas (`layers.js`, `formato.py`, `legendas.py`, `render.py`) | médio |
| C | Reescrever tudo (a SPEC original) | alto |

Esta emenda adota a **opção A**, com B como plano reserva obrigatório antes de qualquer publicação do repositório. Nova decisão **D6** para o usuário confirmar.
Não verificado: o que exatamente em `layers.js`, `sets.js` e `engine.html` é código da RobbiSkill e o que é nosso; precisa de `diff` contra `RobbiSkill/RobbiSkill/engine/`.

## 4. Critérios de sucesso da emenda

1. Cena de referência renderizada pelo processo a partir de `cenas.json` e inventário, com `check_timeline` aprovado e revisão visual do usuário (cena 1 do `diorama-demo`).
2. Reprodução do First Brands (127 cenas, `CHECK OK`, mediana ~3,3 s) a partir do roteiro e do inventário.
3. Elenco do segundo vídeo (tema diferente) criado consultando o inventário e gravando itens novos aprovados.
4. `go test ./...` e `legoloop espinha verificar` verdes; cobertura >= 95% nos pacotes novos.
5. Nomes antigos continuam funcionando por alias por uma versão, com aviso.

## 5. Próximos artefatos (após aprovação)

SPEC v2.0 do macro (com FMEA) → PLAN, ROUTE e DAG regenerados → tickets TDD, um por vez com autorização.
Ordem sugerida dos tickets: (1) esquema e pacote de inventário, (2) catálogo de estilos + estilo molde-de-pílula, (3) reabrir Peça 2 com `cast`, (4) motor portado + ponte `cenas.json` → `timeline.json`, (5) render e voz, (6) renomeação com alias, (7) estilo palito.

## 6. Stop-Gate 01

Aprovado: **[1] seguir** (2026-11-26). Registro em `docs/intents/2026-11-26_EMENDA_INTENT_APROVADA_VIDEO_ANIMADO.md`.
