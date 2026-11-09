## 3. Alinhamento Fonético & A Regra da Folga Pré-Fala (+40ms)

## 3.1 A Física do Ataque Fonético na Fala Humana

Quando um ser humano pronuncia uma palavra iniciada por consoantes plosivas (/p/, /b/, /t/, /d/, /k/, /g/) ou fricativas (/s/, /f/, /v/), a cavidade bucal acumula pressão acústica milissegundos antes que o som seja emitido com energia audível pelo microfone.

Modelos automáticos de transcrição e detectores de voz baseados em limiar de volume (VAD) cometem um erro clássico: eles posicionam a marcação de tempo inicial no ponto exato em que a energia sonora ultrapassa a curva de ruído de fundo.

#painel-perigo("O FENÔMENO DA DECAPITAÇÃO FONÉTICA", [
  Se a tesoura digital cortar o vídeo rigorosamente no timestamp reportado pelo alinhador fonético tradicional, o início da consoante é decepado. A palavra *"Pipeline"* soa como *"-ipeline"*, a palavra *"Código"* soa como *"-ódigo"*, e o espectador tem a sensação inconsciente de que o áudio está picotado, truncado ou com emendas grosseiras de fita magnética.
])

---

## 3.2 A Margem Dourada Inegociável (+40ms)

Para eliminar definitivamente esse problema, o método de Lucas Leite institui a regra mandatória de **+40ms de pré-roll**. Todo e qualquer corte inicial de fala tem seu ponto de entrada antecipado em 40 milissegundos no arquivo de áudio:

$$T_{\text{corte\_início}} = T_{\text{palavra\_início}} - 0.040\text{ s}$$

Essa pequena janela temporal é suficiente para capturar:

1. O ruído suave da abertura dos lábios ou dentes do locutor;
2. A respiração preparatória imediata que confere naturalidade orgânica à fala;
3. O ataque transitório completo da forma de onda da primeira sílaba.

```json
{
  "word": "Processo",
  "detected_start_ms": 14250,
  "applied_cut_start_ms": 14210,
  "pre_roll_offset_ms": -40
}
```

---

## 3.3 Algoritmo de Resolução de Conflitos e Sobreposições

Um desafio matemático surge quando duas sentenças válidas estão muito próximas uma da outra. Se o locutor faz uma pausa curta de apenas 30 milissegundos entre duas frases, subtrair 40 milissegundos da segunda sentença faria com que o corte invadisse a última sílaba da sentença anterior.

O script de cálculo de cortes em Go implementa uma trava de colisão de limites (*boundary clamping*):

```go
func CalculateSafeCutStart(wordStartMs, prevWordEndMs, preRollMs int) int {
    targetStart := wordStartMs - preRollMs
    // Se a folga invadir o término da palavra anterior, limita ao ponto médio
    if targetStart < prevWordEndMs {
        return prevWordEndMs + (wordStartMs - prevWordEndMs) / 2
    }
    return targetStart
}
```

Dessa forma, o algoritmo garante que nenhuma consoante inicial seja mutilada e nenhuma vogal final seja cortada acidentalmente, produzindo uma cadência rítmica perfeita e fluida.
