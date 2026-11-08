# Parte II — Meso: Os 13 Guias Autorais e a Engenharia de Pós-Produção

## 2. Corte Semântico & Extirpação de Notas de Edição

## 2.1 Decisão por Sentença Estruturada

O primeiro passo para um corte automatizado de alta fidelidade é transformar o fluxo contínuo de áudio em uma lista estruturada de sentenças fechadas. Diferente dos algoritmos tradicionais de VAD (Voice Activity Detection), que cortam arbitrariamente por limiar de decibéis em pausas de respiração, o nó `n03_decisao_cortes` processa unidades semânticas completas.

Cada sentença gerada pelo Whisper Local é formatada em um contrato de dados JSON estrito:

```json
{
  "sentence_id": "s_014",
  "text": "Agora nós vamos entender como o pipeline DAG organiza as camadas de nós.",
  "start_ms": 42150,
  "end_ms": 46800,
  "confidence": 0.98,
  "words": [
    {"word": "Agora", "start_ms": 42150, "end_ms": 42400},
    {"word": "nós", "start_ms": 42420, "end_ms": 42600},
    {"word": "vamos", "start_ms": 42610, "end_ms": 42900}
  ]
}
```

O modelo cognitivo (LLM) avalia as sentenças com base em coerência narrativa, eliminando:

* Repetições gagas (*"então nós... então nós vamos..."*);
* Falsos inícios de raciocínio descartados pelo orador;
* Pausas excessivas entre ideias desconectadas.

---

## 2.2 O Protocolo de Extirpação de Notas de Edição

Durante a gravação de aulas e vídeos falados, é extremamente comum o apresentador errar um conceito, parar de olhar para a câmera e dar orientações explícitas ao futuro editor:

> *"Opa, falei errado o nome da função. Corta essa parte aqui e recomeça da explicação do JSON... pronto, voltando."*

Em esteiras convencionais de IA, esse trecho de bastidores é transcrito normalmente e inserido na timeline final, expondo erros graves em produção.

O nó `n03_decisao_cortes` implementa uma máquina de estados semântica projetada especificamente para identificar esses delimitadores conversacionais de gravação:

![Anatomia da Folga Pré-Fala (+40ms) e Extirpação de Notas](img/diagrama_folga_40ms_extirpacao.svg)

### Padrões Canônicos de Detecção de Notas
O analisador semântico busca gatilhos textuais e modulações de prosódia como:

1. **Marcadores de Abertura:** *"nota de edição:"*, *"corta isso"*, *"editor, atenção:"*, *"errei, recomeçando de..."*, *"deixa eu refazer essa frase"*;
2. **Marcadores de Fechamento:** *"fim da nota"*, *"valeu"*, *"voltando em 3, 2, 1"*, *"recomeçando agora"*;
3. **Pausas de Hesitação:** Silêncios anormais superiores a 1.2 segundos acompanhados de mudança abrupta no tom do locutor.

#painel-alerta("REGRA DE EXPURGO TOTAL", [
  Quando uma nota de edição é detectada, o intervalo inteiro — desde o início da hesitação que originou o erro até a retomada definitiva da frase válida — é marcado com a flag `status: "purged"`. O script de montagem remove completamente esses blocos de áudio e vídeo, sem deixar respirações residuais na timeline.
])

---

## 2.3 Exemplo Real: Entrada Bruta vs Saída Higienizada

A tabela abaixo ilustra o comportamento do filtro de extirpação operando sobre um trecho real gravado no laboratório:

| ID             | Sentença Transcrita Bruta                   | Classificação Semântica  | Decisão do Algoritmo           |
| :------------- | :------------------------------------------ | :----------------------: | :----------------------------: |
| `s_041`           | *"O primeiro nó do DAG extrai o arquivo de áudio WAV em 48kHz."* | Conteúdo Didático Válido | **KEEP (Manter na timeline)**      |
| `s_042`           | *"Nota de edição: adicionar um zoom aqui no comando FFmpeg... fim."* | **Nota de Gravação Falada**  | 🛑 **PURGE (Expurgo Obrigatório)** |
| `s_043`           | *"O segundo nó do DAG executa o alinhamento fonético das palavras."* | Conteúdo Didático Válido | **KEEP (Manter na timeline)**      |

O resultado é uma transição suave e imperceptível entre a sentença `s_041` e a sentença `s_043`, economizando tempo precioso de montagem manual e blindando o vídeo contra gafes de estúdio.
