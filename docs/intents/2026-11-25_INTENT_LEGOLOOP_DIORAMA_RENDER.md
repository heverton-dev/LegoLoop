---
id: INTENT-a2ui_sess_3_20261125_render
status: APPROVED
date: 2026-11-25
schema_version: 2
upstream: [INTERVIEW-a2ui-sess_3-20261125]
covers: [Q1]
---

> **Nota de redirecionamento (2026-11-26, AC-E11.4):** o processo `legoloop-diorama-render` foi renomeado para `internal/video/render (ligado ao render Go existente, AC-E10.4)`. Este documento e historico e mantido sem alteracao de conteudo.


# 🎯 INTENT: Criar o Processo canonizado legoloop-diorama-render (Ciclo 3 — Etapas 6, 7, 8 e 8b)

> **Status:** APPROVED  
> **Data:** 2026-11-25  
> **Autor:** Claude Code + Usuário (via Entrevista Socrática A2UI)  
> **Alvo:** `internal/diorama/render`, `internal/cli`, `processos/legoloop-diorama-render/`, `skills/legoloop-diorama-render/`, `internal/videomotion/loudness.go` (detalhado na SPEC)  
> **Sessão A2UI:** sess_3  
> **Confiança da Entrevista:** 95% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Conduzir a Concepção A2UI do Processo `legoloop-diorama-render` (Peça 3 do macroprocesso `legoloop-video-diorama`), cobrindo voz neural com conferência Whisper (Etapa 6), mixagem de áudio com SFX (Etapa 7), motor de composição de timeline em camadas (Etapa 8) e renderização multiformato com cortes verticais (Etapa 8b), consumindo `cenas.json` e `script.json`.
* **Quem Sofre:** O operador do canal (hoje monta voz, mix e render à mão), o Processo 4 (`legoloop-diorama-distribuicao`, que consome renders MP4, SRT e shorts) e o Macroprocesso `legoloop-video-diorama`.
* **O Porquê:** O Ciclo 2 já entrega `cenas.json` e sprites. Sem o render não existe vídeo, e o Ciclo 5 fica bloqueado.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** `projetos/video-diorama/diorama-demo` tem `script.json` (6 linhas) e `cenas.json` (6 cenas, 4 elementos móveis por linha). Não há motor de voz, timeline, mix nem render.
  * **Amanhã:** Processo canonizado `legoloop-diorama-render` (pacote `internal/diorama/render` + CLI `legoloop diorama-render` + DAG 5 de 5), executando por dependência:
    1. **Etapa 6 (Voz e Whisper):** provedor de voz plugável (padrão `edge-tts` `pt-BR-ThalitaMultilingualNeural`; Gemini TTS, ElevenLabs, Fish Audio e outros via adaptador, com gate de créditos antes de chamada paga). `faster-whisper` sobre o áudio final gera `timings.json` por palavra e WER `<= 5%` global e `<= 15%` por linha, com até 3 regravações só das linhas reprovadas, parando se o áudio regravado for idêntico. Gera fonemas/visemas para sincronia labial simples.
    2. **Etapa 8 (Timeline 2.5D):** leitura estrita de `cenas.json` e `script.json`; camadas de cenário, adereços em profundidade, sprites com troca de expressão/pose, câmera (zoom, shake, pan) e legendas karaokê. Grava `timeline.json` validado por `check_timeline` (12–21 cortes/min, mediana 2,5–4 s).
    3. **Etapa 7 (Mix e SFX):** BGM com ducking dinâmico, SFX por evento da timeline, normalização `-14 LUFS ± 0,5` e `-1,5 dBTP`.
    4. **Etapa 8b (Render e cortes):** master 16:9 em 720p, 1080p (padrão) e 4K (opt-in em `render.json`) com NVENC/QSV/AMF detectados automaticamente e fallback `libx264`; cortes 9:16 re-renderizados pela timeline, centrados no narrador e no elemento de ação, legenda fora da caixa do narrador, 15–59 s.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Implementar `internal/diorama/render` e o comando `legoloop diorama-render` cobrindo as Etapas 6, 8, 7 e 8b, na ordem de dependência, mantendo a numeração do pedido nos títulos.
* [ ] Objetivo 2: Gates mecânicos de bloqueio:
  * Gate Etapa 6: WER `<= 5%` global e `<= 15%` por linha; `timings.json` completo por palavra; troca de provedor não altera o contrato.
  * Gate Etapa 8: cobertura de 100% das linhas; `check_timeline` aprovado.
  * Gate Etapa 7: LUFS em `-14 ± 0,5`, TP `<= -1,5 dBTP`, duração áudio == vídeo.
  * Gate Etapa 8b: zero frames em branco; cortes 9:16 entre 15 s e 59 s com área segura.
* [ ] Objetivo 3: Cobertura `>= 95%` nos pacotes novos, testes herméticos (ffmpeg, Whisper e provedores injetáveis), zero stubs, Espinha Dorsal verde e `legoloop fechar` 5 de 5.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Alterar contratos de entrada dos Ciclos 1 e 2 (`projeto.json`, `marca.json`, `script.json`, `cenas.json`, `sprites/*.json`).
* ⛔ Geração de ativos e storyboard (Ciclo 2) e roteiro (Ciclo 1).
* ⛔ Etapas 9 a 13: revisão estéril, vitrine, Postiz, comentários e aprendizado (Ciclo 4).
* ⛔ Composição ponta a ponta `P1->P4` (Ciclo 5).
* ⛔ Visão computacional ou crop por saliency no 9:16.
* ⛔ Painel web paralelo fora do Looping Studio.
* ⛔ Código sem licença de `RobbiSkill/`.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `legoloop-especificacao` | Elaborar a SPEC formal com critérios BDD binários e FMEA. |
| `PECA-02` | `skill` | `legoloop-planejamento` | Elaborar `PLAN`, `ROUTE` e `DAG` de execução. |
| `PECA-03` | `unit_task` | `internal/diorama/render` + `internal/cli` | Implementação via TDD Red-Green com cobertura `>= 95%`. |
| `PECA-04` | `script` | `go test ./...` + `cobertura` + `espinha verificar` | Regressão mecânica, piso de cobertura e integridade. |
| `PECA-05` | `stop_gate` | `legoloop fechar "legoloop-diorama-render"` | Fechamento 5 de 5. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (ou encadeado pelo Ciclo 5 `legoloop-video-diorama`).
* **Nível de Risco & Blast Radius:** Médio (executa ffmpeg e Whisper locais e pode chamar provedores de voz pagos); isolado na pasta `<projeto>/` e no pacote novo `internal/diorama/render`.
* **Política de Privacidade & Orçamento Cognitivo:** Testes 100% herméticos e offline; provedores pagos só atrás do gate de créditos; chaves de API nunca gravadas em artefato.
* **Contrato de Evidência Factual (Definition of Done):**
  * `diorama-demo` percorre as Etapas 6, 8, 7 e 8b e passa em todos os gates mecânicos.
  * Cobertura `>= 95%` (`go test -cover ./internal/diorama/render/...`).
  * `legoloop espinha verificar`, `funclen`, `taxlinter`, `gofmt` e `legoloop fechar` verdes.
* **Restrições Inegociáveis:**
  * Provedor de voz plugável por interface; padrão `edge-tts` Thalita.
  * Alinhamento por palavra sempre pelo `faster-whisper` sobre o áudio final; timestamps nativos do provedor são ignorados.
  * Fallback `libx264` obrigatório; encoder usado registrado no relatório.
  * True Peak `-1,5 dBTP` (Mapa Macro prevalece sobre o `-1 dBFS` do pedido, por decisão do usuário no turno 1).
  * Motor autoral, sem código do `RobbiSkill/`.
  * Teto de True Peak com fonte única em `videomotion.ValidateLoudnessConformity`, corrigido com TDD (autorizado no Stop-Gate 02).
  * Configuração do Ciclo 3 em `render.json` opcional; `projeto.json` permanece intacto.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * A numeração do pedido (7 = mixagem, 8 = timeline) difere do Mapa Macro (7 = timeline, 8 = mix); vale a ordem de dependência confirmada no turno 6.
  * 720p é perfil leve de primeira classe, ao lado de 1080p (padrão) e 4K (opt-in).
* **Riscos e Perguntas em Aberto:** disponibilidade de GPU e de `faster-whisper` na máquina do operador; o Passo 0 do Ciclo 1 já verifica Whisper e ffmpeg.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Invariantes passam com zero stubs. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura `>= 95%` nos pacotes novos. | Medição de cobertura | `go test -cover ./internal/diorama/render/...` |
| 3 | Gates das Etapas 6, 8, 7 e 8b bloqueiam WER global `> 5%` ou por linha `> 15%`, timeline fora de ritmo, LUFS fora de `-14 ± 0,5`, TP acima de `-1,5 dBTP`, frames em branco e cortes fora de 15–59 s. | Casos negativos e positivos em teste | A detalhar na SPEC como `AC-XX` binário |
| 4 | Fechamento 5 de 5 sem lacunas. | Espinha Dorsal e `legoloop fechar` | `legoloop espinha verificar` e `legoloop fechar` com exit code 0 |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md`
* **Gate de Transição:** aprovação humana explícita (Stop-Gate 01) e confiança `>= 95%` antes de qualquer linha da SPEC.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
