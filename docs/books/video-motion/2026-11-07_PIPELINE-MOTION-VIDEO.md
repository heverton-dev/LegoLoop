---
title: "Pipeline Canônico de Motion Video"
subtitle: "Engenharia Reversa de Motion Design, Sincronia Fonética e Renderização por Código"
author:
  - Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "MANUAL OPERACIONAL DIDÁTICO · VIDEO AS CODE"
cover-title: "Motion Video"
cover-accent: "Pipeline Canônico por Código"
cover-caption: "OS SEIS NÓS DO PIPELINE video-motion"
cover-stages:
  - "Silêncio & Ataque +40ms"
  - "Matriz Words & Punch RMS"
  - "Vidro 3D & Tipografia Anton"
  - "Cache SHA-256 & Bridge"
  - "Multi-Aspecto (9:16, 16:9, 1:1)"
  - "Broadcast Bundle & LegoPress"
header-mark: "MOTION VIDEO"
accent: "2563eb"
tagline: |
  Do tutorial ao código executável: pipeline determinístico, multi-aspecto e com zero dependência proprietária
toc: true
toc-depth: 2
abstract: |
  Este livro é o manual técnico e didático oficial do pipeline `video-motion` do LegoLoop. Ele documenta a engenharia reversa do método de motion video editing, a substituição integral de ferramentas proprietárias por código aberto e a arquitetura determinística dos quatro pilares: relógio fonético com corte de silêncio (+40ms), motion de vidro 3D em Remotion/CSS, cache criptográfico SHA-256 e entrega broadcast multi-aspecto (9:16, 16:9, 1:1).
---

# Como ler este livro

## Público e pré-requisitos

Esta obra destina-se a engenheiros de software, criadores de conteúdo técnico e operadores do ecossistema LegoLoop. Pressupõe noções básicas de linha de comando e processamento de mídia.

## Estrutura da obra

| Parte   | Nível      | Conteúdo                                                     |
| :------ | :--------- | :----------------------------------------------------------- |
| **I**   | Macro      | Fundação filosófica, determinismo estrito e zero dependência |
| **II**  | Meso       | Arquitetura dos quatro pilares técnicos                      |
| **III** | Micro      | Guia prático de execução e stop-gates de QA instantâneo      |
| **IV**  | Prática    | Receitas multi-aspecto (9:16, 16:9, 1:1) e broadcast bundle  |
| **V**   | Governança | Rotina de atualização contínua por diff (LegoPress)          |


# Parte I — Fundação e Princípio Determinístico

## A Lei Fundamental: Vídeo é Código

O pipeline `video-motion` foi concebido sob a lei de conduta inegociável do LegoLoop: **tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de uma IA generativa.**

Na produção audiovisual contemporânea, a maioria dos criadores depende de pacotes fechados (After Effects, Premiere, CapCut) ou de plataformas SaaS baseadas em créditos e assinaturas opacas. Essa dependência cria pontos únicos de falha, custos imprevisíveis e impossibilidade de automação em lote.

No LegoLoop, a premissa é invertida:
- O layout visual é composto em **código declarativo** (CSS3 3D, Web-Canvas e Remotion).
- Os tempos de transição são guiados pela matemática da fala (análise de áudio RMS e timestamps milimétricos do Whisper).
- Os assets visuais secundários são gerados via APIs públicas gratuitas (Google AI Studio) ou sintetizados deterministicamente por shaders procedurais em Canvas.

## Comparativo de Eficiência Arquitetural

| Dimensão                 | Fluxo Tradicional (SaaS / NLE)         | Pipeline LegoLoop (`video-motion`)   |
| :----------------------- | :------------------------------------- | :----------------------------------- |
| **Corte de Respiração**  | Manual em timeline ou app proprietário | Automático (-30dB RMS + 40ms ataque) |
| **Sincronia de Legenda** | Digitação ou exportação com drift      | `words.json` milimétrico sem desvio  |
| **Estilo Visual 3D**     | Plugins pagos de After Effects         | Remotion / CSS3 com blur procedural  |
| **Custo por Minuto**     | Assinaturas mensais recorrentes        | Custo zero de licença e 100% local   |
| **Versionamento**        | Binários pesados (.aep / .prproj)      | Código puro versionado em Git        |


# Parte II — A Arquitetura dos Quatro Pilares

O pipeline é estruturado em quatro motores independentes, integrados por um grafo acíclico dirigido (DAG):


![Diagrama de Arquitetura](/img/diagrama_436e5a57.svg){width=95%}


## 1. Relógio Mestre Fonético (Speech-Sync Engine)

A linha do tempo do vídeo não é definida por segundos arbitrários, mas pela cadência real da locução:
1. **Corte Cirúrgico de Pausas:** Todo silêncio abaixo do threshold de -30dB com duração superior a 400ms é removido, preservando estritamente **+40ms de ataque pré-fala** para que nenhuma consoante oclusiva seja deglutida.
2. **Alinhamento Fonético Milimétrico:** Geração da matriz `words.json`, contendo cada vocábulo com seus marcadores de `start_ms` e `end_ms`.
3. **Detecção de Ênfase Vocal:** O motor calcula a energia RMS por palavra. Vocábulos com amplitude 40% superior à média do segmento recebem a flag `is_punch_word = true`.
4. **Normalização Broadcast EBU R128:** Ajuste de ganho para teto de transmissão em **-14 LUFS** com true peak a -1.0 dBFS.

## 2. Compositor Visual Reativo (Remotion & CSS 3D)

Em vez de renderizar frames estáticos pesados, o visual é gerado via DOM reativo acelerado por hardware:
- **Cartões 3D de Vidro:** Superfícies translúcidas com `backdrop-filter: blur(16px)`, bordas sutis com iluminação especular e transição suave de push-in.
- **Tipografia Cinética Anton:** Exibição com destaque cromático instantâneo sobre as palavras de ênfase identificadas na etapa fonética.
- **Auto-Framing Adaptativo:** Recorte inteligente mantendo o sujeito no terço superior de segurança.

## 3. Pipeline de Assets Determinístico com Cache Criptográfico

- **Chave de Cache SHA-256:** Cada solicitação de imagem secundária gera uma assinatura criptográfica a partir do prompt e dimensões. Se já existe no disco, o custo de rede é zero.
- **Bridge com Fallback Local:** Conexão nativa com a API do Google AI Studio. Em caso de indisponibilidade de rede ou ausência de chave, o gerador procedural sintetiza texturas abstratas locais sem interromper o fluxo.

## 4. Renderizador Multi-Aspecto e Bundle Broadcast

- **Exportação Paralela:** Suporte nativo aos formatos `9:16` (Vertical), `16:9` (Widescreen) e `1:1` (Quadrado).
- **Entrega Broadcast Completa:** Geração simultânea do pacote em `dist/<job_id>/` contendo `.mp4`, legendas `.srt` e `.vtt`, arquivo de timecode `fcpxml` (v1.9) e copy pronta para redes em `POST_COPY.md`.


# Parte III — Guia Prático de Execução

## Tríade de Acionamento

O pipeline pode ser executado pelas três vias canônicas do ecossistema LegoLoop, produzindo exatamente os mesmos artefatos e garantindo auditoria completa.

### 1. Via Linha de Comando (CLI)

O binário `legoloop` expõe o comando unificado `video-motion` e seus subcomandos:

```bash
# Execução completa da esteira de ponta a ponta
legoloop video-motion run --input video.mp4 --aspect 9:16 --output dist/

# Execução cirúrgica isolada por etapas:
legoloop video-motion silence --input audio.wav --threshold -30 --pad 40
legoloop video-motion punch   --words words.json --audio audio.wav
legoloop video-motion assets  --prompt "cyberpunk glass ui" --cache .cache/
legoloop video-motion preview --input video.mp4 --output qa_sheet.jpg
legoloop video-motion bundle  --input video.mp4 --output dist/job_01/
```

### 2. Parâmetros e Modos de Operação

| Subcomando | Flag Principal       | Descrição e Comportamento                     |
| :--------- | :------------------- | :-------------------------------------------- |
| `silence`  | `--pad 40`           | Margem de ataque em milissegundos pré-fala    |
| `punch`    | `--words words.json` | Análise de amplitude RMS e destaque Anton     |
| `assets`   | `--prompt "..."`     | Consulta ao cache SHA-256 e bridge AI Studio  |
| `preview`  | `--output qa.jpg`    | Renderização ultra-rápida dos 12 frames-chave |
| `bundle`   | `--output dist/`     | Consolidação do master, legendas e post copy  |

### 3. Via Grafo DAG e LegoLoop Studio

O pipeline é formalmente orquestrado pelo arquivo `loops/video-motion/pipeline.dag.json`:

```bash
legoloop run loops/video-motion/pipeline.dag.json --var input_media="video.mp4"
```

No painel visual do **LegoLoop Studio**, o operador acompanha o progresso nó a nó:
- **Nó 1 (Audio Clock):** Inspeção dos marcadores de silêncio detectados.
- **Nó 2 (Phonetic RMS):** Gráfico de energia sonora com palavras de ênfase destacadas.
- **Nó 3 (Compositor 3D):** Pré-visualização dos shaders e fontes compiladas.
- **Nó 4 (Export):** Geração do bundle final e telemetria de renderização.

## O Stop-Gate de Qualidade Visual (`qa_sheet.jpg`)

Antes do render final em alta resolução, o motor gera em menos de 2 segundos a folha de contato `qa_sheet.jpg`:
- Grade matricial 4x3 contendo os 12 frames mais significativos do vídeo.
- Inspeção visual imediata das legendas, do enquadramento facial e dos reflexos dos cartões de vidro.
- O pipeline pausa aguardando a confirmação do operador para liberar o empacotamento broadcast.

### Protocolo de Validação e Ajustes no Gate

Caso algum artefato visual apresente inconsistência na folha de contato:
1. **Recorte Facial Desalinhado:** Ajustar a sensibilidade do auto-framing via parâmetro `--face-margin`.
2. **Ênfase Vocal Excessiva:** Recalibrar o limiar de RMS com a flag `--rms-boost 1.2`.
3. **Cores ou Reflexos:** Modificar o shader CSS diretamente no arquivo de estilos antes de autorizar o render final.


# Parte IV — Receitas Multi-Aspect Ratio e Broadcast Bundle

## Matriz de Formatos e Áreas de Segurança

Cada rede social possui restrições visuais específicas (sobreposição de botões de like, comentários e legenda de rodapé). O motor ajusta o layout com base na matriz abaixo:

| Formato  | Resolução | Destino Principal        | Margem de Segurança    |
| :------- | :-------- | :----------------------- | :--------------------- |
| **9:16** | 1080x1920 | TikTok, Reels, Shorts    | 120px topo, 280px base |
| **16:9** | 1920x1080 | YouTube Desktop, Cursos  | 60px uniforme          |
| **1:1**  | 1080x1080 | Instagram Feed, LinkedIn | 80px topo e base       |

## Anatomia do Pacote de Entrega (`dist/<job_id>/`)

Ao concluir a renderização, a pasta de entrega é montada com todos os ativos requeridos para distribuição profissional:

```
dist/job_20261107_01/
├── video_final_9x16.mp4      # Master de transmissão codificado em H.264 / AAC
├── legendas.srt              # Legenda universal para players externos
├── legendas.vtt              # Legenda web com estilização CSS
├── timeline.fcpxml           # XML v1.9 importável no DaVinci Resolve Studio
├── qa_sheet.jpg              # Prova de contato dos 12 frames-chave
└── POST_COPY.md              # Texto de postagem com hashtags e timing de publicação
```


# Parte V — Rotina de Atualização por Diff (LegoPress)

## Governança Viva de Documentação

No ecossistema LegoLoop, nenhuma documentação pode se tornar obsoleta em relação ao código real. Sempre que um arquivo do pipeline (`internal/videomotion/` ou `loops/video-motion/`) for alterado via commit ou refatoração, o motor **LegoPress** sincroniza a obra técnica automaticamente.

## Comandos de Manutenção do Livro

```bash
# Auditar a integridade geométrica do PDF gerado
legoloop legopress auditar docs/books/video-motion/2026-11-07_PIPELINE-MOTION-VIDEO.pdf

# Atualizar capítulos a partir dos diffs recentes do Git
legoloop legopress atualizar docs/books/video-motion/

# Recompilar o livro completo em caso de alterações nas partes
bash docs/books/compilar.sh video-motion
```

A esteira de publicação gera automaticamente o manuscrito consolidado em Markdown, o documento tipográfico final em Typst, o PDF de alta resolução (170x240mm) e as imagens de visualização para cada página na pasta `preview/`.


