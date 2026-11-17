# Capítulo 4: Prática & Operação na Tríade Completa

O pipeline **LegoLoop Cartoon Shorts** oferece três portas de entrada com perfeita paridade operacional:

## 1. Operação via Linha de Comando (CLI)
A interface de terminal permite integração com scripts de automação, cron jobs e pipelines de lote:
```bash
# Execução direta informando o tema
legoloop cartoon-shorts --theme "O Enigma do ET de Varginha"

# Ajuste de duração e voz neural
legoloop cartoon-shorts --theme "A Invenção do Avião" --duration 25 --voice pt-BR-FranciscaNeural

# Validação do contrato em modo dry-run
legoloop cartoon-shorts --theme "Ouro Preto" --dry-run
```

## 2. Operação via Grafo DAG no Looping Engine
Execução sequencial por nós isolados em contexto estéril:
```bash
legoloop run loops/legoloop-cartoon-shorts/pipeline.dag.json
```

## 3. Invocação via Chat dos Agentes
Para criar um novo short em meio a qualquer sessão de desenvolvimento:
```text
/legoloop-cartoon-shorts Como o jogo GTA 6 vazou por uma TV de hotel?
```
O agente acionará o motor, gerará o roteiro de curiosidades com as cenas calibradas, sintetizará o áudio Edge-TTS e devolverá o arquivo `.mp4` compilado pronto para publicação.
