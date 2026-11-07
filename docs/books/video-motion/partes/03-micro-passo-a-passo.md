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
