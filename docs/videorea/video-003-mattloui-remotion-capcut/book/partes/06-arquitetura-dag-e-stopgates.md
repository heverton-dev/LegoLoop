# Parte VI — Micro: Arquitetura DAG e Governança

## 6.1 Topologia do Grafo Acíclico Dirigido

```mermaid
graph TD
  In["1. Ingest de Telas (Video/Print)"] --> Prompt["2. Claude Code Prompt Engine"]
  Prompt --> TSX["3. Remotion TSX Generator"]
  TSX --> Render["4. 4K ProRes Alpha Render"]
  Render --> FX["5. FFmpeg FX Engine (Chromatic & Blur)"]
  FX --> SFX["6. Procedural Sound Design"]
  SFX --> Out["7. Master Output MP4"]

  TSX -.->|"SG-CODE: Validação de Sintaxe"| TSX
  Render -.->|"SG-ALPHA: Verificação de Canal Transparente"| Render
```

### Stop-Gates Human-in-the-Loop:
* **SG-CODE:** Validação estática de TypeScript para garantir zero erros de tipagem antes de acionar a GPU.
* **SG-ALPHA:** Smoke test determinístico de 1 segundo para certificar que o fundo do vídeo é 100% transparente.
