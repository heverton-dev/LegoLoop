# Capítulo 3 — Arquitetura do Pipeline

O pipeline `ai-motion-graphics` é um DAG acíclico com 5 fases:

```
Seleção → Geração → Renderização → Pós-processamento → Entrega
```

## Contratos Principais

| Tipo | Campos | Regras |
|:---|:---|:---|
| Request | Prompt, Scenario, Resolution, Duration | Duration > 0, prompt não vazio |
| RepoSuggestion | ID, Tier, Confidence, Reason | Confidence ≥ 0.7 |
| Result | MP4Path, ReposUsed, Scenario | MP4 existe, ffprobe confirma |

## Costuras de Teste

- **Primária:** `motion.Run(req) (*Result, error)`
- **Injectáveis:** `frameCapturer` (Playwright), `frameEncoder` (FFmpeg), `beatDetector` (Librosa)
- **Dublês:** apenas nas bordas (E/S, processos externos)
