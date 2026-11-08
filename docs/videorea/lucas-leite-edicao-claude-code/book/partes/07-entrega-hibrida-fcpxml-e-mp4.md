## 7. Entrega Híbrida Não-Destrutiva: FCPXML + MP4

## 7.1 A Armadilha da Entrega Fechada

A imensa maioria das ferramentas de inteligência artificial comercializadas para criadores comete um erro fatal de fluxo de trabalho: elas entregam um único arquivo MP4 compactado. Se o criador discordar de um corte, quiser ajustar o volume de uma palavra ou trocar uma cor em um cartão gráfico, ele enfrenta duas opções ruins:
1. Reenviar o vídeo para a nuvem, esperar a fila de processamento e torcer para o modelo acertar;
2. Importar o MP4 achatado em um editor e tentar recortar sobre uma imagem já compactada, onde a trilha musical e a voz estão misturadas e inseparáveis.

---

## 7.2 A Filosofia da Linha do Tempo Aberta

O ecossistema LegoLoop estabelece que o objetivo de uma esteira autônoma não é substituir o editor humano, mas sim **fazer 90% do trabalho pesado e entregar uma timeline aberta e organizada** para que o humano faça os 10% de refinamento estético e polimento final.

O nó `n11_entrega_final_hibrida` bifurca a entrega em dois artefatos complementares:

```
dist/
├── timeline.fcpxml        # Timeline universal com pistas isoladas e handles
├── render_final.mp4       # Render broadcast imediato em 1080p60
└── manifest_final.json    # Metadados e árvore de cuts para auditoria
```

---

## 7.3 Anatomia do Padrão FCPXML Universal

O formato **Final Cut Pro XML (FCPXML)** foi adotado como o padrão aberto de interoperabilidade do LegoLoop porque ele é lido nativamente e sem perda pelas principais estações de trabalho profissionais do mercado: **DaVinci Resolve**, **Final Cut Pro** e **Adobe Premiere Pro**.

### Estrutura das Pistas na Linha do Tempo
Ao abrir o arquivo `timeline.fcpxml`, o editor encontra uma timeline perfeitamente organizada por camadas:

* **Pista de Vídeo 1 (V1):** O vídeo bruto original cortado nas sentenças selecionadas, preservando **sobras (*handles*) de 0.5 segundo** nas extremidades de cada clipe para permitir que o editor faça transições cruzadas (*crossfades*) ou ajustes manuais;
* **Pista de Vídeo 2 (V2):** Os cartões gráficos animados renderizados pelo Remotion com canal alfa (transparência perfeita), posicionados exatamente nos timestamps corretos;
* **Pista de Áudio 1 (A1):** Áudio limpo da voz do apresentador com os cortes fonéticos de +40ms;
* **Pista de Áudio 2 (A2):** Trilha musical com as curvas de automação de volume de -18dB já desenhadas;
* **Pista de Áudio 3 (A3):** Efeitos sonoros (*SFX*) alinhados milimetricamente aos cartões visuais;
* **Marcadores de Capítulo (Markers):** Marcadores coloridos na linha do tempo indicando exatamente onde a IA removeu notas de gravação ou sugeriu títulos de seções.

#painel-sucesso("LIBERDADE CRIATIVA TOTAL", [
  O editor não perde tempo cortando silêncios nem caçando notas de gravação. Ele abre a timeline no DaVinci Resolve ou Premiere e foca imediatamente na narrativa, na gradação de cor (*color grading*) e no acabamento fino.
])
