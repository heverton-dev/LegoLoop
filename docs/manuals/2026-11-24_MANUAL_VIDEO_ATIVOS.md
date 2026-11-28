# Manual Operacional — `legoloop-video-ativos` (Ciclo 2)

Data: 2026-11-24. Versão: 1.0.0.

## 1. Finalidade
Executar as Etapas 3, 4 e 5 do processo composto `legoloop-video-animado`:
- **Pré-voo Ciclo 2:** valida `projeto.json`, `marca.json`, `script.json` e `plano.md` do Ciclo 1, permissões em `flow/` e `sprites/`, `estilo_visual` (`codigo`, `flow`, `hibrido`) e trava `flow_projeto`.
- **Etapa 3 (Elenco e Sprites):** valida 3 fichas obrigatórias por personagem (`turnaround`, `expressoes_3x2`, `poses_6x`), executa recorte transparente (`rembg` + despill verde) e grava `flow/elenco.json` e `sprites/*.json` sem sobreposição de `bounding_box`.
- **Etapa 4 (Cenários e Objetos):** modo `codigo` gera SVGs 100% autorais; modo `flow` usa fila resiliente com retry/backoff, `--ref`, trava `flow_projeto` e idempotência por hash SHA-256; modo `hibrido` combina ambos; Gate exige `>= 1` cenário distinto a cada 2 linhas.
- **Etapa 5 (Storyboard):** compila `cenas.json` cobrindo 100% das linhas de `script.json`, com metáfora visual e `>= 4` elementos móveis por linha, respeitando `politica_gates`.

## 2. Comandos CLI
```bash
legoloop video-ativos prevoo --projeto <DIR> --json
legoloop video-ativos executar --projeto <DIR> [--aprovar etapa3_elenco,etapa4_quadros_teste,etapa5_storyboard] --json
legoloop video-ativos validar-storyboard --projeto <DIR> --json
```
