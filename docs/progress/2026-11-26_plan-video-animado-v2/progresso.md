---
schema: progresso.v1
iniciativa: "V01"
plano: "docs/plans/2026-11-26_PLAN_VIDEO_ANIMADO_V2.md"
atualizado_em: "2026-11-26T19:35:39Z"
---

# Progresso Vivo: V01

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 25 | 7 | 0 | 0 | 0 | 18 |

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
| `TICK-V01-03.2` | E08 | Validar cenas.json v2 e elenco.json contra o inventario e as partes suportadas do renderizador | AC-E08.2, AC-E08.3 |
| `TICK-V01-04.1` | E09 | Etapa 3 vetorial: montar elenco.json com cast a partir do inventario, sem Flow, e gravar itens novos como nao aprovados | AC-E09.1, AC-E09.2 |
| `TICK-V01-04.2` | E09 | Manter o caminho de fichas e rembg para estilos de recorte sem regressao | AC-E09.3 |
| `TICK-V01-04.3` | E09 | Etapa 4: consultar o inventario antes de gerar fundos e registrar fundo novo com prompt, tags e hash | AC-E09.4 |
| `TICK-V01-04.4` | E09 | Marcar entrega suspeita do Flow (sem confirmacao de renomeacao, hash cruzado, fora do padrao de estilo) | AC-E09.5 |
| `TICK-V01-05.1` | E10 | Criar motor-local com LICENCA.json, ignorar no git e teste de governanca da licenca | AC-E10.1 |
| `TICK-V01-05.2` | E10 | Implementar o pre-voo do motor: python, Playwright, Chromium, ffmpeg e motor-local, com ErrMotorAusente | AC-E10.2 |
| `TICK-V01-05.3` | E10 | Implementar a ponte cenas.json v2 + elenco + timings + inventario para timeline.json, com luz de contorno do estilo | AC-E10.3, AC-E10.6, AC-E10.8 |
| `TICK-V01-05.4` | E10 | Renderizar a cena de referencia pelo motor: 9:16, voz, determinismo e boca acompanhando a fala | AC-E10.4, AC-E10.5, AC-E10.7 |
| `TICK-V01-06.1` | E10 | Ligar o render Go ao motor e reconfirmar voz, mix, SFX e cortes da SPEC v1.4.0 sobre a nova timeline | AC-E10.4 |
| `TICK-V01-07.1` | E11 | Mover internal/diorama para internal/video e atualizar imports sem regressao | AC-E11.2 |
| `TICK-V01-07.2` | E11 | Registrar os comandos video-roteiro, video-ativos e video-render com alias diorama-* e aviso de depreciacao | AC-E11.1 |
| `TICK-V01-07.3` | E11 | Migrar skills, processos, espinha e PROCESSOS_CANONIZADOS para os nomes legoloop-video-* com redirecionamento | AC-E11.3 |
| `TICK-V01-07.4` | E11 | Aplicar a nota de redirecionamento nos documentos arquivados sem alterar o conteudo | AC-E11.4 |
| `TICK-V01-08.1` | E12 | Renderizar a cena 1 do diorama-demo pelo processo e colher a aprovacao visual do usuario | AC-E12.1 |
| `TICK-V01-08.2` | E12 | Criar o estilo palito herdando do molde-de-pilula e validar com o usuario | AC-E12.3 |
| `TICK-V01-08.3` | E12 | Reproduzir o First Brands pelo processo: timeline completa com CHECK OK e ritmo no alvo | AC-E12.2 |
| `TICK-V01-09.1` | E12 | Fechamento: cobertura, dependencias, portabilidade, espinha e legoloop fechar 5 de 5 | AC-NF.1, AC-NF.2, AC-NF.3 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| `TICK-V01-01.1` | E07 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/video/inventario	0.016s / | 2026-11-26T15:07:04Z |
| `TICK-V01-01.2` | E07 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/video/inventario	0.113s / | 2026-11-26T15:19:05Z |
| `TICK-V01-01.3` | E07 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/video/inventario	0.140s / | 2026-11-26T16:30:01Z |
| `TICK-V01-01.4` | E07 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/video/inventario	0.087s / ok  	github.com/heverton-dev/LegoLoop/internal/cli	0.097s / | 2026-11-26T16:43:14Z |
| `TICK-V01-02.1` | E07 | 0 | map[estilo:molde-de-pilula ok:true] / | 2026-11-26T16:52:55Z |
| `TICK-V01-02.2` | E07 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/video/inventario	0.077s / map[estilo:molde-de-pilula ok:true] / | 2026-11-26T17:39:52Z |
| `TICK-V01-03.1` | E08 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/diorama/ativos	0.132s / | 2026-11-26T19:35:39Z |
