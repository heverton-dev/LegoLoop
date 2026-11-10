# Capítulo 4 — Seleção Híbrida de Repos

## O Modelo Híbrido

1. O pipeline **analisa** o pedido e **sugere** repos rankeados por confiança
2. O operador **confirma** ou ajusta a seleção
3. Os repos selecionados alimentam a geração HTML/WebGL

## Mapeamento Cenário → Repos

| Cenário | Repos Sugeridos (ordem) |
|:---|:---|
| explainer | threejs → react-postprocessing → librosa → blender-mcp → react-rapier |
| pixel_art | pixelfx → lemo-opuscar |
| product_launch | video-shotcraft → librosa → react-postprocessing → threejs → prompt-packs |

## Validação

- Cenário inválido → `ErrInvalidRequest`
- Prompt vazio → `ErrInvalidRequest`
- Nenhum repo selecionado → `ErrNoRepoSelected`
