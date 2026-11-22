---
id: ESPELHO-intv_1795370552989667600
date: 2026-11-22
sessao: intv_1795370552989667600
confianca: 98
status: RECONCILIADO
---

# Espelho Canônico — Motor Único de Interface Declarativa A2UI no LegoLoop

> **Sessão:** `intv_1795370552989667600`  
> **Confiança da Entrevista:** 98% (meta: 95%, 3 turnos + reconciliação arquitetural com o usuário)  
> **Data:** 2026-11-22  

## 1. resultado_esperado
Construir um **motor único de interface declarativa A2UI** (`@a2ui/lit` embutido no binário Go via `go:embed`) que dê a toda execução de Processo do LegoLoop uma **URL de sessão interativa** compartilhada pelas **3 vias** (Terminal CLI, Agente de IA / Slash Command via MCP `legoloop`, e botão Executar do Studio), com SSE de saída (`GET /api/sessao/{id}/stream`), entrada de ações (`POST /api/sessao/{id}/action`), Design System próprio do LegoLoop, catálogo básico A2UI + componentes de domínio (DAG, progresso/confiança, Stop-Gate, visualizador de artefatos, logs) e suporte a macroprocessos via `surfaceId` (árvore lateral + breadcrumb). A entrega ocorre em duas fases: **Fase 1 (Piloto)** no Processo `legoloop-concepcao` nas 3 vias; **Fase 2 (Unificação)** trazendo o monitoramento do Swarm como superfície A2UI e eliminando a UI web antiga e separada do Swarm, deixando o LegoLoop com apenas duas interfaces: Studio Canvas (construção) e Interface Interativa A2UI (execução, interação e monitoramento).

## 2. quem_se_beneficia
O usuário e operadores do LegoLoop em qualquer sistema operacional (Windows, Linux, macOS) e qualquer harness (Antigravity, Claude Code, OpenCode, MimoCode, CLI pura), que passam a interagir visualmente com Processos, Stop-Gates, artefatos e enxames mantendo fallback textual completo.

## 3. por_que_agora
Hoje, a via Terminal e a via Chat/Slash Command ficam restritas a texto linear, o botão Executar do Studio abre apenas um modal de log SSE textual, e o monitoramento do Swarm exige uma UI web paralela separada. Um motor único A2UI unifica interação e monitoramento para todos os Processos atuais e futuros com um único renderizador seguro.

## 4. criterio_de_sucesso
- Piloto `legoloop-concepcao` rodando de ponta a ponta nas 3 vias (perguntas dinâmicas, barra de confiança, Espelho Canônico editável, Stop-Gate 01 clicável e preview do `INTENT.md`).
- Corrida entre entrada textual (`stdin` / chat) e interface (`POST /api/sessao/{id}/action`) funcionando sem travamento e com fallback textual preservado.
- Zero HTML ou JavaScript executável aceito de agentes (validação estrita sobre catálogo confiável A2UI).
- Macroprocessos isolados por `surfaceId` com navegação por árvore e breadcrumb e foco automático em Stop-Gates.
- Na Fase 2, paridade funcional do monitoramento do Swarm como superfície A2UI comprovada por testes antes da remoção da UI antiga do Swarm.
- Todos os testes e gates do LegoLoop verdes (`go test -count=1 ./...`, `go vet ./...`, `funclen`, `gofmt`, `taxlinter`).

## 5. restricoes_inegociaveis
1. Zero HTML ou script executável enviado por agentes (apenas JSON declarativo A2UI validado por esquema sobre catálogo confiável).
2. Fallback textual garantido nas 3 vias (quem operar só no terminal/chat segue funcionando sem abrir navegador).
3. Isolamento por bloco e por `surfaceId` em macroprocessos (Lei 4, contexto estéril).
4. Design System do próprio LegoLoop no host (`@a2ui/lit` + `go:embed`, sem exigir internet em runtime e sem migrar o Studio para React).
5. O motor/backend do Swarm permanece intacto; apenas a UI antiga duplicada é removida na Fase 2 após paridade provada.
6. TDD estrito (Red-Green-Refactor) e respeito à governança viva (Leis 5, 8, 12, 16 e 18).

## 6. fora_de_escopo
- Criar um formato proprietário de UI fora do protocolo aberto A2UI.
- Reescrever o Studio Canvas inteiro em React, Angular ou Flutter.
- Quebrar o formato atual do DAG ou os endpoints SSE já existentes antes da paridade provada.
- Remover o motor/backend do Swarm (`legoloop swarm`, quadro de mensagens, workers).
- Integrar dashboards externos fora do LegoLoop.
