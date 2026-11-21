# AI Motion Graphics: Pipeline Soberano com 10 GitHub Repos

> **Autor:** LegoLoop  
> **Data:** 2026-11-09  
> **Versão:** 1.0.0  

---

## Sumário

1. Introdução — O Problema do "AI Slop"
2. Os 10 GitHub Repos
3. Arquitetura do Pipeline
4. Seleção Híbrida de Repos
5. Renderização Playwright + FFmpeg
6. Sync de Áudio com Librosa
7. Guia Operacional


# Capítulo 1 — Introdução: O Problema do "AI Slop"

Motion graphics gerados por IA sem ferramentas especializadas resultam em visual genérico — o chamado "AI slop". O vídeo de referência demonstra que integrar repos especializados eleva a qualidade de genérica para profissional.

## Os Três Cenários

| Cenário           | Exemplo                                     | Repo Principal      |
| :---------------- | :------------------------------------------ | :------------------ |
| Explainer         | "Como um painel solar funciona"             | Three.js            |
| Lançamento        | "Vídeo de lançamento do produto X"          | Video ShotCraft     |
| Pixel Art         | "Tempestade em pixel art"                   | PixelFX             |

## O Fluxo

Pedido em linguagem natural → Seleção de repos → Geração HTML/WebGL → Renderização Playwright → Conversão FFmpeg → MP4 final.

Custo total: **R$ 0,00/mês** — todos os 10 repos são open-source.


# Capítulo 2 — Os 10 GitHub Repos

## Tier 1 — Essenciais

| #   | Repo                 | GitHub                             | Função                                 |
| :---- | :------------------- | :--------------------------------- | :------------------------------------- |
| 1   | Three.js             | mrdoob/three.js                    | 3D para web — gráficos, partículas, fly-throughs |
| 2   | PixelFX              | patch-gg/pixelfx                   | Pixel art com biblioteca de efeitos    |
| 3   | Video ShotCraft      | Vincentwei1021/video-shotcraft     | 124 templates de promo/lançamento      |
| 4   | Librosa              | librosa/librosa                    | Detecção de beats para sync de áudio   |
| 5   | React Postprocessing | pmndrs/react-postprocessing        | Glow, DoF, grain, lens, glitch         |

## Tier 2 — Expansão

| #   | Repo                  | GitHub                               | Função                       |
| :---- | :-------------------- | :----------------------------------- | :--------------------------- |
| 6   | Lemo Opuscar          | lemomo-ai/lemo-opuscar               | 66 estilos cinematográficos  |
| 7   | Blender MCP           | ahujasid/mcp-for-blender             | CGI profissional via MCP     |
| 8   | Prompt Packs          | yihui-dev/awesome-opus5-5-videos     | Prompts HTML/WebGL prontos   |
| 9   | Claude Animation Base | JohnHeibel/ClaudeAnimationBase       | Mascote animado + emotions   |
| 10  | React Rapier          | pmndrs/react-three-rapier            | Física real (Jenga, dominós) |


# Capítulo 3 — Arquitetura do Pipeline

O pipeline `ai-motion-graphics` é um DAG acíclico com 5 fases:

```
Seleção → Geração → Renderização → Pós-processamento → Entrega
```

## Contratos Principais

| Tipo               | Campos                                 | Regras                         |
| :----------------- | :------------------------------------- | :----------------------------- |
| Request            | Prompt, Scenario, Resolution, Duration | Duration > 0, prompt não vazio |
| RepoSuggestion     | ID, Tier, Confidence, Reason           | Confidence ≥ 0.7               |
| Result             | MP4Path, ReposUsed, Scenario           | MP4 existe, ffprobe confirma   |

## Costuras de Teste

- **Primária:** `motion.Run(req) (*Result, error)`
- **Injectáveis:** `frameCapturer` (Playwright), `frameEncoder` (FFmpeg), `beatDetector` (Librosa)
- **Dublês:** apenas nas bordas (E/S, processos externos)


# Capítulo 4 — Seleção Híbrida de Repos

## O Modelo Híbrido

1. O pipeline **analisa** o pedido e **sugere** repos rankeados por confiança
2. O operador **confirma** ou ajusta a seleção
3. Os repos selecionados alimentam a geração HTML/WebGL

## Mapeamento Cenário → Repos

| Cenário                | Repos Sugeridos (ordem)                                   |
| :--------------------- | :-------------------------------------------------------- |
| explainer              | threejs → react-postprocessing → librosa → blender-mcp → react-rapier |
| pixel_art               | pixelfx → lemo-opuscar                                    |
| product_launch          | video-shotcraft → librosa → react-postprocessing → threejs → prompt-packs |

## Validação

- Cenário inválido → `ErrInvalidRequest`
- Prompt vazio → `ErrInvalidRequest`
- Nenhum repo selecionado → `ErrNoRepoSelected`


# Capítulo 5 — Renderização Playwright + FFmpeg

## Playwright (Captura de Frames)

- Abre HTML em Chromium headless
- Captura screenshots sequenciais (fps × duração)
- 90 frames para 3s @ 30fps
- `headless: false` proibido em CI

## FFmpeg (Conversão MP4)

```bash
ffmpeg -y -framerate 30 -i frame_%04d.png -c:v libx264 -pix_fmt yuv420p -s 1920x1080 -r 30 output.mp4
```

## Verificação

`ffprobe` confirma:

- Codec: h264
- Duração: ±0.1s do esperado
- Resolução: igual ao Request

## Limite

Renderização de 60s de vídeo ≤ 300s (5 min).


# Capítulo 6 — Sync de Áudio com Librosa

## Por Que Sync Importa

Todo visual change que acontece no beat da música torna o vídeo "premium". É o diferencial entre amador e profissional.

## Como Funciona

1. `scripts/beat_detect.py` carrega áudio via Librosa
2. `librosa.beat.beat_track()` detecta beats
3. Retorna timestamps (segundos) para sincronizar animações

## Fallback Gracioso

O contrato de áudio tem dois modos. O pedido declara o modo.

### Modo Broadcast (Shorts, Reels, Standalone)

- O áudio integrado a -14 LUFS é obrigatório.
- Áudio ausente quebra o gate. O Processo para com erro.
- Arquivo corrompido quebra o gate.

### Modo Motion Overlay (camadas transparentes para DaVinci Resolve)

- A trilha de áudio é opcional.
- Áudio ausente produz um MP4 com canal mudo. O gate não quebra.
- Arquivo corrompido produz um MP4 com canal mudo, sem sync e sem erro.

Os dois modos resolvem a contradição entre o contrato BDD de -14 LUFS (`harness-engineering`) e o fallback silencioso (`ai-motion-graphics`).

## Exemplo (120 BPM)

```json
[0.0, 0.5, 1.0, 1.5, 2.0, 2.5]
```


# Capítulo 7 — Guia Operacional

## Comando Rápido

```bash
legoloop run ai-motion-graphics
```

## Entrada Típica

> "Crie um vídeo explicativo de 30 segundos sobre como um painel solar funciona, 1920x1080"

## Fluxo do Operador

1. Pipeline sugere repos → **confirme** ou ajuste
2. Geração HTML/WebGL → aguarde
3. Renderização Playwright + FFmpeg → aguarde
4. MP4 entregue em `output/`

## Solução de Problemas

| Problema                         | Solução                                         |
| :------------------------------- | :---------------------------------------------- |
| Playwright não inicia            | Verificar `npx playwright install chromium`       |
| FFmpeg não encontrado            | `winget install ffmpeg` ou adicionar ao PATH      |
| Librosa não detecta beats        | Verificar `pip install librosa`                   |
| MP4 sem áudio                    | Fornecer `AudioPath` no Request                   |

## Custo

**R$ 0,00/mês** — todos os 10 repos são open-source.


