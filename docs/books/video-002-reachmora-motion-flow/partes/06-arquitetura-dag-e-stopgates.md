## 4.2 Topologia DAG e Stop-Gates

```mermaid
graph TD
  A["1. Script Input"] --> B["2. Master Prompt Parser"]
  B --> C["3. Beat Slicer (10s)"]
  C --> D["4. Asset Generator (ComfyUI/Remotion)"]
  D --> E["5. Audio Synthesizer (Local TTS)"]
  E --> F["6. Final Assembly (FFmpeg)"]
  
  C -.->|"SG-BEATS: Aprovação de Tomadas"| C
  D -.->|"SG-RENDER: Qualidade de Imagem"| D
```

### Stop-Gates Invioláveis:
* **SG-BEATS:** O operador aprova a lista de tomadas e palavras-chave antes de iniciar a renderização gráfica.
* **SG-RENDER:** Verificação de nitidez e ausência de artefatos visuais nos primeiros 3 segundos de cada clipe.
