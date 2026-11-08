## 6. Sonorização, Ducking & Sincronia de Efeitos Sonoros (SFX)

## 6.1 A Psicologia Acústica do Efeito Sonoro

O som é responsável por mais de 50% da percepção de qualidade de uma produção audiovisual. Em vídeos educativos e tutoriais técnicos, a inserção de efeitos sonoros (*SFX*) cumpre uma função pedagógica: direcionar a atenção ocular do espectador para a informação mais relevante da tela.

No entanto, o maior vício de esteiras automatizadas de edição é disparar o som de transição (*whoosh*) no primeiro frame em que o cartão começa a se mover na borda da tela. 

#painel-alerta("O ERRO DO GATILHO PREMATURO", [
  Quando o efeito sonoro de transição toca cedo demais, o cérebro do espectador processa o impacto acústico antes que o elemento visual termine sua trajetória. Isso cria uma sensação incômoda de antecipação desordenada e desvaloriza a animação.
])

---

## 6.2 O Alinhamento no Frame de Frenagem

Na metodologia de Lucas Leite, a curva de envelope do som de transição é invertida ou calculada para que o pico absoluto de pressão sonora coincida exatamente com o **frame de parada/frenagem** (*landing frame*) do elemento visual:

![Sonorização Sincronizada (Ducking -18dB) e Entrega Híbrida](img/diagrama_ducking_entrega_hibrida.svg)

$$T_{\text{pico\_sfx}} = T_{\text{parada\_card}} \pm 1\text{ frame (0.016s)}$$

Quando o cartão atinge sua posição final com uma frenagem suave, o som de impacto (*whoosh + click*) atinge o volume máximo. O olho e o ouvido confirmam o evento no mesmo milissegundo, gerando uma experiência tátil, sólida e agradável.

---

## 6.3 Ducking Dinâmico de -18dB sob a Fala

A música de fundo (*background music*) cria atmosfera e mantém a retenção, mas pode comprometer a inteligibilidade das frequências médias da voz do locutor (entre 800Hz e 3.5kHz).

O nó `n10_efeitos_sonoros_sfx` aplica compressão sidechain automatizada via filtro FFmpeg:

```bash
### Trecho de comando FFmpeg para Ducking Dinâmico com atenuação de -18dB
ffmpeg -i voz_limpa.wav -i trilha_fundo.mp3 -filter_complex \
  "[1:a]volume=0.8[music]; \
   [music][0:a]sidechaincompress=threshold=0.08:ratio=6:attack=40:release=250[music_ducked]; \
   [0:a][music_ducked]amix=inputs=2:duration=first" saida_mixada.wav
```

### Parâmetros Recomendados de Calibração:
* **Atenuação da Trilha sob Fala:** `-18dB` (garante que a voz permaneça nítida e dominante);
* **Ataque (*Attack*):** `40ms` (a música abaixa rapidamente assim que a primeira sílaba soa);
* **Liberação (*Release*):** `250ms` (a música retorna suavemente após o término da frase, sem oscilações bruscas de volume nos pequenos respiros);
* **Ganho nos Silêncios Narrativos:** `0dB` (preenche pausas deliberadas com musicalidade).
