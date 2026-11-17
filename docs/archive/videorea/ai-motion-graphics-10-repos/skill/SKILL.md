---
name: ai-motion-graphics
description: Use quando o usuário quiser criar motion graphics com IA (explicativos, lançamentos de produto, pixel art) usando os 10 GitHub repos open-source, renderizando MP4 via Playwright+FFmpeg. NÃO use para editar vídeos existentes (use legoloop-video-motion) ou engenharia reversa de vídeos (use legoloop-video-rea).
compatibility: Windows/Linux, Node.js 18+, Python 3.10+, FFmpeg 5+, Playwright (npx), Go 1.22+
allowed-tools: ReadFile WriteFile RunCommand BrowserPlaywright FFmpeg
consumes: Artifact:Spec
produces: Artifact:Code
artifact: SKILL.md
gate: has_binary_criteria('SKILL.md')
suggested_model: code-pro
metadata:
  version: "2.2.0"
  tier: "ouro"
  category: "build"
---

# 🧩 ai-motion-graphics — Pipeline Soberano de Motion Graphics com 10 GitHub Repos

> **Autor:** LegoLoop (derivado de Matt Penny — "10 GitHub Repos to 10x Your AI Motion Graphics")  
> **Versão:** 2.2.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.2.0  
> **Taxonomia:** Peça de domínio `ai-motion-graphics` (não é função nativa `legoloop-*`)

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O operador deseja criar motion graphics em MP4 (explicativos, lançamentos de produto, pixel art) a partir de pedido em linguagem natural, usando repos open-source especializados com renderização automática Playwright+FFmpeg.  
**Quando NÃO usar:** Para editar vídeos existentes (use `legoloop-video-motion`), engenharia reversa de vídeos (use `legoloop-video-rea`), ou conteúdo fora da esteira audiovisual (use `legoloop-legopress`).

Esta peça opera estritamente sob os seguintes contratos de entrada e saída:

* **Consumes (Entradas Obrigatórias):**
  - Request JSON com `prompt`, `scenario`, `resolution`, `duration`
* **Produces (Saídas Mecânicas):**
  - MP4 renderizado com codec h264, resolução e duração conforme Request
  - `loops/ai-motion-graphics/pipeline.dag.json` validado
* **Target Files (Escopo de escrita permitido / Blast Radius):**
  - `internal/motion/`, `scripts/beat_detect.py`, `skills/ai-motion-graphics/`, `docs/books/ai-motion-graphics/`, `loops/ai-motion-graphics/`, `docs/videorea/ai-motion-graphics-10-repos/`

```json
{
  "consumes": ["Request{prompt, scenario, resolution, duration}"],
  "produces": ["output.mp4", "loops/ai-motion-graphics/pipeline.dag.json"],
  "gate": "has_binary_criteria('output.mp4')"
}
```

### Procedimento Determinístico de Execução
1. Validar `Request` (prompt não vazio, duração > 0, cenário válido). Se inválido, abortar com `ErrInvalidRequest`.
2. Executar `SuggestRepos(scenario, prompt)` e apresentar sugestões ao operador para confirmação híbrida.
3. Executar `BuildHTML(repos, prompt, resolution)` → HTML/WebGL com canvas e imports dos repos.
4. Executar `captureFramesWith(playwrightCapturer, html, fps, duration, resolution, outDir)` → N frames PNG.
5. Executar `encodeWith(ffmpegEncoder, frames, audioPath, outPath, resolution)` → MP4 h264.
6. Se `AudioPath` fornecido, executar `detectBeatsWith(librosaDetector, audioPath)` para sync de beats.
7. Verificar MP4 com `ffprobe` (codec h264, duração ±0.1s, resolução). Declarar concluído com evidências.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** gere código incompleto com `TODO`, `pass`, `// implementar depois` ou mocks simulados (Zero Stubs).
> 2. **NUNCA** avance sem confirmação do operador na seleção híbrida de repos (fluxo HITL).
> 3. **NUNCA** escreva fora dos Target Files declarados na seção 1 (violação de write boundary).
> 4. **NUNCA** renderize com Playwright `headless: false` em CI ou smoke tests.
> 5. **NUNCA** declare MP4 entregue sem `ffprobe` confirmar codec h264 e duração esperada.

---

## 3. Comandos de Verificação & Testes de Terminal

```bash
# Validação de compilação / Build
go build -o bin/legoloop.exe ./cmd/legoloop

# Execução de testes focados com cobertura (exit code 0 obrigatório)
go test -v -count=1 ./internal/motion/...

# Verificação estática e formatação (saída vazia obrigatória)
go vet ./internal/motion/...
gofmt -l internal/motion/

# Funclen (nenhuma função > 35 linhas)
go run ./tools/funclen internal

# Validação do DAG
legoloop validate loops/ai-motion-graphics/pipeline.dag.json

# Verificação do MP4 entregue
ffprobe -v quiet -show_entries format=duration -show_entries stream=codec_name,width,height -of json output.mp4
```

---

## 4. Critério de Parada (Stopping Checklist)

A execução desta peça é considerada concluída com sucesso quando:
- [ ] O comando de teste unitário passou com exit code 0 comprovado no terminal (32 testes).
- [ ] O artefato MP4 existe fisicamente no disco e `ffprobe` confirma h264 + duração.
- [ ] A expressão `gate` do frontmatter avalia como verdadeira pelo motor de validação.
- [ ] Zero resíduos ou arquivos temporários órfãos deixados no sistema de arquivos.
- [ ] Stop-Gate humano aprovado formalmente (SG-E01 a SG-FINAL).
