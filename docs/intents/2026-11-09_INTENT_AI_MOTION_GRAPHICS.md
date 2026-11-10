# 🎯 INTENT: AI Motion Graphics — Pipeline Soberano com 10 GitHub Repos

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-09  
> **Autor:** MiMoCode (via legoloop-video-rea)  
> **Iniciativa:** I-AIMG — AI Motion Graphics  
> **Alvo:** `loops/ai-motion-graphics/`, `skills/ai-motion-graphics/`, `docs/books/ai-motion-graphics/`  
> **Sessão:** ses_ffe5e3d81734fffeURpuQcduSB  
> **Confiança da Entrevista:** 96%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

* **O Problema Real:** Motion graphics gerados por IA sem ferramentas especializadas resultam em "AI slop" — visual genérico, sem profundidade, sem sincronismo e sem identidade. O vídeo de referência (Matt Penny, "10 GitHub Repos to 10x Your AI Motion Graphics") demonstra que a integração de repos especializados eleva a qualidade de genérica para profissional. Não existe pipeline LegoLoop que orquestre esses repos de forma soberana.
* **Quem Sofre:** Operadores LegoLoop que precisam produzir motion graphics (explicativos, lançamentos de produto, pixel art) e atualmente dependem de ferramentas pagas ou prompts manuais sem estrutura.
* **O Porquê:** Os 10 repos são open-source (R$ 0,00/mês). Integrá-los num pipeline canônico LegoLoop elimina a dependência da comunidade paga "Applied AI Mastermind" e democratiza produção audiovisual profissional.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O operador cola URLs de repos manualmente no Claude, sem orquestração, sem seleção inteligente, sem renderização automatizada e sem governança.
  * **Amanhã:** Um comando (`legoloop run ai-motion-graphics`) recebe um pedido em linguagem natural, sugere os repos ideais, o operador confirma, o pipeline renderiza HTML/WebGL via Playwright, converte para MP4 via FFmpeg e entrega vídeo final com sync de áudio.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Pipeline canônico `ai-motion-graphics` com 5 nós Tier 1 (Three.js, PixelFX, Video ShotCraft, Librosa, React Postprocessing) cobrindo 3 cenários (explicativo, lançamento, pixel art).
* [ ] Objetivo 2: Skill canônica universal distribuída via `legoloop instalar-skills` com seleção híbrida de repos (pipeline sugere, operador confirma).
* [ ] Objetivo 3: Livro Didático LegoPress compilado com os 10 repos, arquitetura do pipeline e guia operacional.
* [ ] Objetivo 4: Renderização Playwright + FFmpeg para MP4 com resolução e duração variáveis conforme o pedido.
* [ ] Objetivo 5: Tier 2 de expansão (5 repos restantes) integrado como nós ativáveis sob demanda.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não implementaremos integração com a comunidade paga "Applied AI Mastermind" ou qualquer ferramenta proprietária.
* ⛔ Não criaremos UI web para o pipeline (CLI e Studio apenas).
* ⛔ Não faremos upload automático para plataformas de vídeo (YouTube, Instagram, etc.).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/ai-motion-graphics` | Skill canônica universal com contratos de I/O e guardrails. |
| `BLOCO-02` | `unit_task` | Motor de seleção de repos | Analisa o pedido e sugere repos Tier 1/Tier 2 (híbrido: sugere + confirma). |
| `BLOCO-03` | `unit_task` | Motor de renderização Playwright | Carrega HTML/WebGL, captura frames via headless browser. |
| `BLOCO-04` | `unit_task` | Motor de conversão FFmpeg | Converte frames para MP4 com resolução/duração variáveis. |
| `BLOCO-05` | `unit_task` | Motor de sync Librosa | Detecta beats de áudio e sincroniza animações. |
| `BLOCO-06` | `script` | `go test ./...` | Regressão mecânica e integridade do código Go do legoloop. |
| `BLOCO-07` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |
| `BLOCO-08` | `subloop` | `legoloop legopress` | Compilação e auditoria do Livro Didático. |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual` — operador invoca via `legoloop run ai-motion-graphics` ou `/legoloop-video-rea`.
* **Nível de Risco & Blast Radius:** `Baixo (workspace local)` — toda produção é local, sem efeitos em sistemas compartilhados.
* **Política de Privacidade & Orçamento Cognitivo:** `Nuvem autorizada` (Claude API para inferência). Dados sensíveis ou segredos envvidos: nenhum.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Um MP4 renderizado de ponta a ponta a partir de pedido em linguagem natural, com os 3 cenários testados (explicativo, lançamento, pixel art), sync de áudio verificado e zero dependências pagas.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * Playwright consegue capturar WebGL/Three.js em headless Chromium sem fallback para software rendering.
  * Librosa 0.11.0 já instalado é suficiente para beat detection sem configuração adicional.
  * Os 10 repos GitHub estão acessíveis e estáveis (sem breaking changes iminentes).
* **Riscos e Perguntas em Aberto:**
  * Blender MCP requer Blender instalado (~300MB) — será Tier 2 opcional, mas verificar viabilidade headless.
  * PixelFX pode ter dependências de browser não suportadas por Playwright headless — testar antes do Tier 1.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Pipeline renderiza MP4 end-to-end para os 3 cenários. | 3 MP4s gerados e verificados | `ffprobe` confirma codec/duração/resolução |
| 2 | Seleção híbrida funciona: pipeline sugere repos, operador confirma. | Fluxo interativo documentado | Log do Studio/CLI com sugestão + confirmação |
| 3 | Sync de áudio Librosa sincroniza animações ao beat. | Comparação visual com/sem sync | Fixture de 3s com e sem beat sync |
| 4 | Skill universal instalada em todos os harnesses. | `legoloop instalar-skills` exit 0 | `skills/ai-motion-graphics/SKILL.md` presente |
| 5 | Livro Didático compilado e auditado. | `legoloop legopress auditar` exit 0 | `docs/books/ai-motion-graphics/` completo |
| 6 | Custo mensal total = R$ 0,00 em ferramentas. | Nenhuma dependência paga | `curation.json` confirma todos open-source |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md`](../specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA (`SG-INTENT`):**  
*A intenção consolidada com 96% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
