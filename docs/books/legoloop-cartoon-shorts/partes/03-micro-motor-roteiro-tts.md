# Capítulo 3: O Motor de Roteiro & Síntese Soberana Edge-TTS

A alma do short está no gancho dos primeiros 3 segundos. Se a frase de abertura for previsível ou pedir para o usuário curtir o vídeo, a retenção colapsa para menos de 20%.

## A Fórmula do Gancho Imediato (No-Setup Hook)

O motor `BuildCartoonScript` do LegoLoop implementa a diretiva de **Impacto Imediato**:
* A primeira frase é um paradoxo afirmado com naturalidade histórica:
  > *"Uma cidade de Minas ficou famosa por um visitante que ninguém conseguiu provar que existiu."*
  > *"Cientistas deram a um remédio o nome do computador mais perigoso do cinema."*
* Não há saudações (*"Olá pessoal, bem-vindos ao meu canal"*). O vídeo começa com o acontecimento no primeiro frame (Frame 0).

## Síntese de Voz com Edge-TTS

Substituímos o ElevenLabs comercial pelo motor local **Edge-TTS 7.3.1**, utilizando a voz neural `pt-BR-AntonioNeural`:
* **Zero Custo:** Nenhuma cobrança de créditos ou limite de caracteres.
* **Expressividade:** As indicações de interpretação entre colchetes (`[reflexivo]`, `[animado]`) são filtradas para o arquivo de áudio enquanto preservam o timing de respiração e pausas naturais do narrador.
* **Hermetismo:** Em caso de interrupção ou falha de rede, o método remove qualquer arquivo corrompido, evitando áudios quebrados na montagem.
