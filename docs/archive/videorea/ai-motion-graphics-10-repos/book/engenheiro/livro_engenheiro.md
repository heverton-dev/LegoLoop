# AI Motion Graphics — Livro do Engenheiro

> **Público:** Desenvolvedores, arquitetos de software, operadores LegoLoop  
> **Nível:** Avançado · **Data:** 2026-11-09 · **Versão:** 1.0.0  

---

# PARTE I — FUNDAMENTOS

## Capítulo 1 — O Problema e a Solução

### 1.1 O Custo do "AI Slop"

Motion graphics gerados por IA sem ferramentas especializados produzem o chamado *AI slop* — visual genérico, sem profundidade, sem sincronismo e sem identidade visual. A evidência empírica (Matt Penny, "10 GitHub Repos to 10x Your AI Motion Graphics", 2026-11-09) demonstra que a integração de **repos especializados** eleva a qualidade de genérica para profissional.

A tese central é simples: **a qualidade vem de ativos curados por humanos (templates, estilos, shaders) combinados com seleção e execução por IA**.

### 1.2 Os 3 Cenários de Uso

| Cenário | Prompt Exemplo | Repos Primários | Output Esperado |
|:---|:---|:---|:---|
| **Explainer** | "Explique como um painel solar funciona" | Three.js + Postprocessing | Vídeo 3D infográfico |
| **Lançamento** | "Vídeo de lançamento do app X" | Video ShotCraft + Librosa | Promo sincronizada |
| **Pixel Art** | "Tempestade em pixel art" | PixelFX + Lemo Opuscar | Animação retrô |

### 1.3 Princípio de Soberania

Todos os 10 repos são **open-source** (MIT/ISC). Custo mensal de ferramentas: **R$ 0,00**. Não há dependência de comunidades pagas, SaaS proprietários ou licenças.

---

## Capítulo 2 — Os 10 Repos: Análise Técnica

### 2.1 Tier 1 — Essenciais (5 repos)

#### Three.js (mrdoob/three.js)
- **Licença:** MIT · **Estrelas:** 100k+ · **Linguagem:** JavaScript
- **Função:** Renderização 3D via WebGL para web
- **Capacidades:** Bar charts 3D, modelos procedurais, efeitos de partícula (galáxias), fly-throughs de câmera
- **Por que é Tier 1:** Base de todo motion graphic 3D. Sem ele, tudo é 2D plano.
- **Evidência no vídeo:** 2:34–3:59 — gráfico de barras 3D com pilhas de dinheiro, modelo 3D criado pelo Claude, galáxia espiral, fly-through de cena

#### PixelFX (patch-gg/pixelfx)
- **Licença:** Open-source · **Linguagem:** HTML/Canvas
- **Função:** Pixel art motion graphics com biblioteca de efeitos embutida
- **Capacidades:** Explosões, clima, fogo, fogos de artifício, aurora boreal, relâmpagos
- **Por que é Tier 1:** Cobre inteiramente o cenário pixel art. Efeitos prontos que IA pura não reproduz.
- **Evidência no vídeo:** 4:00–4:52 — cena de tempestade pixel art, explosão, explicativo de sistema climático

#### Video ShotCraft (Vincentwei1021/video-shotcraft)
- **Licença:** Open-source · **Linguagem:** Templates/HTML
- **Função:** 124 receitas de vídeo promo/lançamento de produto
- **Capacidades:** Templates testados em lançamentos reais, seleção automática do melhor template
- **Por que é Tier 1:** Reprodutibilidade de vídeos de lançamento. O humano fez o trabalho criativo; a IA seleciona e executa.
- **Evidência no vídeo:** 11:38–13:04 — demonstração de templates de produto

#### Librosa (librosa/librosa)
- **Licença:** ISC · **Linguagem:** Python · **Instalado:** v0.11.0
- **Função:** Análise de áudio — detecção de beats, estimação de tempo
- **Capacidades:** Beat detection ±50ms, feature extraction, sincronização visual-musical
- **Por que é Tier 1:** O diferencial "premium" que separa amador de profissional. Tudo acontece no beat.
- **Evidência no vídeo:** 10:49–11:37 — demo de sync de animações com beats de música

#### React Postprocessing (pmndrs/react-postprocessing)
- **Licença:** MIT · **Linguagem:** React + GLSL
- **Função:** Pós-processamento de luz e efeitos visuais
- **Capacidades:** Glow/bloom, depth of field (DoF), film grain, lens coloring, glitch
- **Por que é Tier 1:** Polimento visual universal. Eleva qualquer renderização de "flat" para "premium".
- **Evidência no vídeo:** 9:56–10:48 — glow em átomos de sódio, DoF em cenas 3D

### 2.2 Tier 2 — Expansão (5 repos)

| Repo | Licença | Gatilho de Ativação | Capacidade-Chave |
|:---|:---|:---|:---|
| Lemo Opuscar | Open | Pedido de estilo cinematográfico | 66 estilos (paper cuts, ink wash, retro terminal) |
| Blender MCP | Open | Pedido de CGI profissional | Cenas 3D completas via MCP |
| Prompt Packs | Open | Pedido de efeitos HTML/WebGL prontos | Coleções de prompts testadas |
| Claude Animation Base | Open | Pedido de mascote/explicativo | Mascote "Clawd" com emotions |
| React Rapier | MIT | Pedido de simulação física | Jenga, dominós, Galton board |

---

# PARTE II — ARQUITETURA

## Capítulo 3 — Pipeline DAG

### 3.1 Topologia

```
TICK-01.1 → TICK-01.2 → TICK-01.3 → SG-E01
                                      ↓
TICK-02.1 → TICK-02.2 → SG-E02
                          ↓
TICK-03.1 → TICK-03.2 → TICK-03.3 → SG-E03
                                      ↓
TICK-04.1 → TICK-04.2 → SG-E04
                          ↓
TICK-05.1 → TICK-05.2 → SG-FINAL
```

### 3.2 Contratos de Dados

```go
// Erros sentinela
var ErrInvalidRequest  = errors.New("motion: requisição inválida")
var ErrNoRepoSelected  = errors.New("motion: nenhum repo selecionado")
var ErrRenderFailed    = errors.New("motion: falha na renderização")
var ErrNoFrames        = errors.New("motion: nenhum frame para codificar")

// Cenários válidos
const (
    ScenarioExplainer     Scenario = "explainer"
    ScenarioProductLaunch Scenario = "product_launch"
    ScenarioPixelArt      Scenario = "pixel_art"
)

// Request — contrato de entrada
type Request struct {
    Prompt     string   `json:"prompt"`
    Scenario   Scenario `json:"scenario"`
    Resolution string   `json:"resolution"`
    Duration   float64  `json:"duration"`
    AudioPath  string   `json:"audio_path,omitempty"`
}

// RepoSuggestion — saída do motor de seleção
type RepoSuggestion struct {
    ID         string  `json:"id"`
    Name       string  `json:"name"`
    Tier       int     `json:"tier"`
    Reason     string  `json:"reason"`
    Confidence float64 `json:"confidence"`
}
```

### 3.3 Costuras de Teste (Test Seams)

| Costura | Tipo | Justificativa |
|:---|:---|:---|
| `motion.Run(req)` | Primária | Contrato externo do pacote |
| `frameCapturer` | Interface injectável | Separa Playwright da lógica |
| `frameEncoder` | Interface injectável | Separa FFmpeg da lógica |
| `beatDetector` | Interface injectável | Separa Librosa da lógica |

**Dublês permitidos:** apenas nas bordas (E/S, processos externos, filesystem). Nunca dublê de módulo interno.

---

## Capítulo 4 — Motor de Seleção Híbrida

### 4.1 Algoritmo

```
SuggestRepos(scenario, prompt) → []RepoSuggestion
1. Validar scenario (enum check)
2. Mapear scenario → lista ordenada de repo IDs
3. Para cada repo, atribuir confidence do catálogo
4. Gerar Reason contextual
5. Retornar lista rankeada
```

### 4.2 Tabela de Mapeamento

| Cenário | Repos (ordem de prioridade) |
|:---|:---|
| explainer | threejs(0.90) → react-postprocessing(0.80) → librosa(0.85) → blender-mcp(0.75) → react-rapier(0.70) |
| pixel_art | pixelfx(0.95) → lemo-opuscar(0.75) |
| product_launch | video-shotcraft(0.90) → librosa(0.85) → react-postprocessing(0.80) → threejs(0.90) → prompt-packs(0.70) |

### 4.3 Fluxo HITL (Human-in-the-Loop)

1. Pipeline sugere repos → apresenta ao operador
2. Operador confirma ou ajusta
3. Repos confirmados alimentam `BuildHTML()`
4. **Pular confirmação = violação de guardrail**

---

## Capítulo 5 — Renderização

### 5.1 Playwright (Captura)

```
captureFramesWith(capturer, html, fps, duration, resolution, outDir)
→ N = fps × duration frames PNG
→ 90 frames para 3s @ 30fps
→ Cada frame com dimensões do Request
```

### 5.2 FFmpeg (Conversão)

```bash
# Args gerados por ffmpegArgs():
ffmpeg -y \
  -framerate 30 \
  -i frame_0000.png -i frame_0001.png ... \
  [-i audio.wav] \
  -c:v libx264 \
  -pix_fmt yuv420p \
  -s 1920x1080 \
  -r 30 \
  [-c:a aac -shortest] \
  output.mp4
```

### 5.3 Verificação Mecânica

```bash
ffprobe -v quiet \
  -show_entries format=duration \
  -show_entries stream=codec_name,width,height \
  -of json output.mp4

# Esperado:
# codec_name: "h264"
# width: 1920, height: 1080
# duration: ±0.1s do Request
```

### 5.4 Limites de Performance

| Métrica | Limite | Verificação |
|:---|:---|:---|
| Renderização 60s vídeo | ≤ 300s | Timeout no teste |
| Frames por segundo | 30 fps | Contagem de frames |
| Tamanho máximo HTML | 50KB | `len(html) < 51200` |
| Memória | `GOMEMLIMIT` | Monitoramento Go |

---

## Capítulo 6 — Sync de Áudio (Librosa)

### 6.1 Pipeline de Beat Detection

```
detectBeatsWith(detector, audioPath) → []float64
1. Validar audioPath (não vazio)
2. Verificar existência do arquivo
3. Invocar detector (Librosa via Python)
4. Retornar timestamps de beats
5. Fallback gracioso: qualquer erro → nil, nil
```

### 6.2 Script Python

```python
# scripts/beat_detect.py
import librosa, json, sys

def detect_beats(audio_path):
    y, sr = librosa.load(audio_path, sr=22050)
    tempo, beats = librosa.beat.beat_track(y=y, sr=sr)
    return [float(b * 512 / sr) for b in beats]

# Saída: JSON array de timestamps [0.0, 0.5, 1.0, ...]
```

### 6.3 Precisão

- **120 BPM:** beats a cada 0.5s, tolerância ±50ms
- **Detecção mínima:** ≥ 4 beats para 3s de áudio
- **Fallback:** áudio ausente/corrompido → sem sync, sem erro

---

# PARTE III — OPERAÇÃO

## Capítulo 7 — Comandos Operacionais

```bash
# === DESENVOLVIMENTO ===
go test -count=1 ./internal/motion/...     # Testes unitários
go vet ./internal/motion/...               # Estática
gofmt -l internal/motion/                  # Formatação (saída vazia)
go run ./tools/funclen internal            # Funções ≤ 35 linhas

# === BUILD ===
go build -o bin/legoloop.exe ./cmd/legoloop

# === VALIDAÇÃO ===
legoloop validate loops/ai-motion-graphics/pipeline.dag.json

# === SKILL ===
legoloop instalar-skills --dry-run --json

# === VERIFICAÇÃO DE OUTPUT ===
ffprobe -v quiet -show_entries format=duration -show_entries stream=codec_name -of json output.mp4
```

## Capítulo 8 — Solução de Problemas

| Erro | Causa | Solução |
|:---|:---|:---|
| `ErrInvalidRequest` | Prompt vazio ou duration ≤ 0 | Validar entrada |
| `ErrNoRepoSelected` | Lista de repos vazia | Confirmar seleção |
| `ErrRenderFailed` | Playwright/FFmpeg falhou | Verificar PATH |
| `ErrNoFrames` | Nenhum frame gerado | Verificar HTML |
| Playwright não inicia | Chromium não instalado | `npx playwright install chromium` |
| FFmpeg não encontrado | PATH incompleto | `winget install ffmpeg` |
| Librosa não detecta beats | Python/librosa ausente | `pip install librosa` |
| MP4 sem áudio | AudioPath não fornecido | Incluir AudioPath no Request |

## Capítulo 9 — Métricas de Qualidade

| Portão | Comando | Critério |
|:---|:---|:---|
| Testes | `go test ./internal/motion/...` | 32/32 pass, exit 0 |
| Estática | `go vet ./internal/motion/...` | 0 warnings |
| Formatação | `gofmt -l internal/motion/` | saída vazia |
| Funclen | `go run ./tools/funclen internal` | 0 violações em motion |
| Build | `go build -o bin/legoloop.exe` | exit 0 |
| DAG | `legoloop validate` | exit 0 |
| Deps | `git diff --exit-code go.mod go.sum` | sem mudanças |
