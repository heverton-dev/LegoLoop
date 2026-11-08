# Parte IV — Meso: Render 4K e Canal Alpha

## 4.1 O Segredo da Transparência Perfeita

Para que as animações gráficas possam ser sobrepostas em qualquer vídeo de fundo sem mascaramentos imperfeitos ou telas pretas, a exportação deve conter o canal Alpha transparente.

### Configuração de Renderização Headless:
```bash
npx remotion render src/index.ts PromptBar dist/overlay.mov   --pixel-format yuva444p10le   --codec prores   --prores-profile 4444   --image-format png
```

#painel-alerta("POR QUE PRORES 4444?", [
  O formato ProRes 4444 preserva 100% da transparência com qualidade sem perdas (lossless) em 4K. O arquivo pode ser solto diretamente na timeline do CapCut, Premiere ou DaVinci sem qualquer conversão.
])
