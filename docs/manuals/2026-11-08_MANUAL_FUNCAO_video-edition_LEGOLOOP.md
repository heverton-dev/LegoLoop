# 📘 MANUAL OPERACIONAL: `video-edition` (LegoLoop)

> **Função:** Edição Audiovisual Assistida por IA com Governança Soberana  
> **Metodologia Original:** Lucas Leite (*"Como eu edito meus vídeos com Claude Code"*)  
> **Status:** Ativo · Produção · Governança LegoLoop v2.2.0  
> **Data:** 08/11/2026  
> **Custo de Licenciamento:** R$ 0,00 / mês (100% Local e Soberano)  

---

## 1. Visão Geral & Filosofia

O pipeline `video-edition` executa a montagem automatizada e soberana de produções em vídeo faladas (videoaulas, podcasts, tutoriais de tecnologia e apresentações corporativas), combinando inteligência artificial semântica com precisão matemática de áudio e vídeo em 12 nós acíclicos governados por 3 Stop-Gates humanos inegociáveis.

### A Lei Áurea de Separação de Papéis:
1. **O Modelo LLM Nunca Calcula Frames:** O modelo cognitivo opera sobre frases inteiras no texto, selecionando sentenças válidas e expurgando notas de gravação.
2. **Scripts Determinísticos Calculam Mídia:** Scripts utilitários locais (FFmpeg, Python, Node.js) calculam timestamps, convertem frequências, aplicam a folga pré-roll de +40ms e geram a timeline final.
3. **Entrega Sempre Não-Destrutiva:** A esteira entrega simultaneamente um MP4 em 1080p60 e uma timeline aberta em XML (FCPXML) com handles e pistas separadas para DaVinci Resolve ou Premiere.

---

## 2. Instruções de Execução (CLI)

### 2.1 Disparo Básico
```bash
legoloop run loops/video-edition/pipeline.dag.json --var input_video="caminho/do/video.mp4"
```

### 2.2 Disparo com Ajuste Fino de Calibração
```bash
legoloop run loops/video-edition/pipeline.dag.json \
  --var input_video="gravacao_bruta.mov" \
  --var pre_roll_ms=50 \
  --var ducking_db="-20dB" \
  --var eye_level_ratio=0.35
```

---

## 3. Dicionário Estrito de Variáveis de Runtime

| Variável | Tipo | Default | Limites Aceitáveis | Efeito no Processamento |
| :--- | :---: | :---: | :---: | :--- |
| `var.input_video` | `string` | *obrigatório* | Arquivo `.mp4`, `.mov`, `.mkv` | Caminho absoluto ou relativo do vídeo bruto gravado pelo apresentador. |
| `var.pre_roll_ms` | `int` | `40` | `20` a `80` (ms) | Margem de segurança antecipada antes da fala para preservar consoantes plosivas. |
| `var.ducking_db` | `string` | `"-18dB"` | `"-12dB"` a `"-26dB"` | Atenuação dinâmica do volume da trilha musical enquanto houver voz ativa do locutor. |
| `var.eye_level_ratio` | `float` | `0.35` | `0.28` a `0.42` | Posição vertical relativa da linha dos olhos do apresentador a partir do topo do quadro. |
| `var.language` | `string` | `"pt"` | Códigos ISO 639-1 (`pt`, `en`, `es`) | Idioma utilizado pelo Whisper Local no alinhamento fonético das sentenças. |

---

## 4. Protocolo Operacional dos 3 Stop-Gates Humanos

A automação cega é estritamente proibida pelo LegoLoop. A esteira pausa compulsoriamente em 3 momentos:

### 🛑 Stop-Gate 01: `n04_sg_revisao_corte` (Validação Semântica do Corte)
* **Momento:** Ocorre após a transcrição fonética e a seleção das frases pelo LLM.
* **O que o Operador Deve Inspecionar:**
  1. Abrir o arquivo `scratch/cuts_manifest.json`;
  2. Verificar se todas as frases marcadas com `action: "keep"` fazem sentido narrativo;
  3. Confirmar que trechos com instruções de estúdio (*"nota de edição: ... fim da nota"*) foram marcados como `action: "purge"`.
* **Ação no Terminal:** Digite `1` para aprovar e prosseguir, ou edite o JSON diretamente e depois digite `1`.

---

### 🛑 Stop-Gate 02: `n06_sg_revisao_motion` (Aprovação do Storyboard Gráfico)
* **Momento:** Ocorre antes de disparar a renderização pesada de animação do Remotion.
* **O que o Operador Deve Inspecionar:**
  1. Abrir o arquivo `scratch/motion_storyboard.json`;
  2. Ler os textos dos cartões, verificar erros de digitação e acentuação;
  3. Conferir se a posição do cartão (`left` ou `right`) está correta para não cobrir o apresentador.
* **Ação no Terminal:** Digite `1` para autorizar a renderização procedural dos gráficos em lote.

---

### 🛑 Stop-Gate 03: `n12_sg_entrega_final` (Aprovação da Linha do Tempo Aberta)
* **Momento:** Conclusão da montagem da timeline universal FCPXML e renderização do MP4.
* **O que o Operador Deve Inspecionar:**
  1. Abrir o arquivo `dist/timeline.fcpxml` no DaVinci Resolve ou Final Cut Pro;
  2. Tocar a reprodução e conferir a sincronia entre a voz e o pico dos efeitos sonoros de transição;
  3. Validar se as sobras (*handles*) de 0.5s estão presentes para eventuais ajustes manuais de transição.
* **Ação no Terminal:** Digite `1` para declarar a entrega homologada e emitir o certificado de telemetria.

---

## 5. Matriz de Resolução de Falhas (Troubleshooting)

| Incidente Operacional | Causa Raiz | Ação Corretiva Imediata |
| :--- | :--- | :--- |
| **Corte inicial engolindo sílaba plosiva (/p/, /t/, /k/)** | Ataque fonético muito veloz do locutor | Aumente a folga de pré-roll reexecutando com `--var pre_roll_ms=55`. |
| **Rosto descentralizado durante o zoom** | Baixa iluminação ou apresentador virou de perfil | Abra `scratch/yunet_tracking.json` e ative o fallback para centro geométrico `50%`. |
| **Trilha de fundo abafando a voz** | Atenuação sidechain insuficiente | Aumente a compressão passando `--var ducking_db="-22dB"`. |
| **Falha de render no nó Remotion (`n09`)** | Módulo NPM ausente ou memória de Node insuficiente | Execute `npm install` na pasta do renderizador e aumente o timeout do nó. |
| **Áudio mudo na timeline FCPXML no DaVinci** | Frequência de amostragem incompatível | Confirme se o nó `n01` extraiu o áudio em rigorosos 48.000 Hz, 16-bit WAV PCM. |
| **Nota de gravação vazou no corte final** | Locutor não usou marcadores verbais claros | Edite manualmente o arquivo `scratch/cuts_manifest.json` no Stop-Gate 01. |

---

## 6. Estrutura de Diretórios e Saídas Finais

```
dist/
├── timeline.fcpxml        # Linha do tempo universal aberta (DaVinci/Premiere)
├── render_final.mp4       # Render broadcast direto em 1080p60 H.264/AAC
└── manifest_final.json    # Metadados e telemetria estrita de produção
```
