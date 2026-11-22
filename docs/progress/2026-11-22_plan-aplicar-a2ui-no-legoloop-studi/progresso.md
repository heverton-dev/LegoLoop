---
schema: progresso.v1
iniciativa: "A2UI"
plano: "docs/plans/2026-11-22_PLAN_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md"
atualizado_em: "2026-11-22T00:00:00Z"
---

# Progresso Vivo: A2UI

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 13 | 0 | 0 | 0 | 0 | 13 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-A2UI-01.1` | E01 | Criar contratos, validador de catálogo confiável e bloqueio de payload hostil em internal/a2ui | AC-E01.2, AC-E03.3 |
| `TICK-A2UI-01.2` | E01 | Implementar Hub de sessões interativas com SSE, snapshot de reconexão e entrega de Action sem vazamento | AC-E01.1, AC-E01.3, AC-NF.3 |
| `TICK-A2UI-01.3` | E01 | Expor rotas HTTP /api/sessao, /stream, /action e /artefato protegidas por s.secure e internal/jail no Studio | AC-E01.1, AC-E02.3, AC-NF.1 |
| `TICK-A2UI-02.1` | E02 | Empacotar bundle offline de @a2ui/lit e @a2ui/web_core em internal/studio/static/components via go:embed | AC-E02.1 |
| `TICK-A2UI-02.2` | E02 | Implementar catálogo de domínio LegoLoop sobre @a2ui/lit com Design System do Studio | AC-E02.2 |
| `TICK-A2UI-03.1` | E03 | Implementar hierarquia de macroprocessos isolados por surfaceId com árvore, breadcrumb e foco automático em Stop-Gate | AC-E03.1, AC-E03.2 |
| `TICK-A2UI-04.1` | E04 | Implementar AwaitAction (corrida stdin vs navegador) e integração da Via 1 (Terminal CLI) | AC-E04.1 |
| `TICK-A2UI-04.2` | E04 | Adicionar ferramentas MCP sessao_criar, sessao_publicar e sessao_aguardar no legoloop mcp-router (Via 2) | AC-E04.2 |
| `TICK-A2UI-04.3` | E04 | Conectar botão Executar do Studio Canvas à rota de sessão interativa A2UI preservando POST /api/run (Via 3) | AC-E04.3 |
| `TICK-A2UI-05.1` | E05 | Implementar adaptador A2UI do piloto legoloop-concepcao (confiança, perguntas, Espelho, Stop-Gate 01 e INTENT.md) | AC-E05.1, AC-E05.2 |
| `TICK-A2UI-06.1` | E06 | Fase 2: Implementar superfície LegoSwarmBoard alimentada por swarm_api.go com paridade funcional e segurança | AC-E06.1 |
| `TICK-A2UI-06.2` | E06 | Fase 2: Remover UI web antiga do Swarm após paridade provada, mantendo motor/backend e redirecionamento | AC-E06.2 |
| `TICK-A2UI-06.3` | E06 | Portão final global: race, vet, funclen, gofmt, taxlinter, espinha verificar e recompilação fresca | AC-NF.2, AC-NF.3 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| nenhum | - | - | - | - |
