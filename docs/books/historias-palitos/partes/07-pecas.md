# 7. Peças

Peça é qualquer item que o processo usa: script, binário, serviço, skill ou dado.
Este capítulo registra cada peça, o que faz e onde está.

## 7.1 Binários

| Peça           | Versão verificada      | Função                      | Verificação                           |
| :------------- | :--------------------- | :-------------------------- | :------------------------------------ |
| Python         | 3.14.7                 | roda todos os scripts       | `python --version`                      |
| FFmpeg         | 7.1 (essentials)       | junta clipes, recorta áudio | `ffmpeg -version`                       |
| ffprobe        | acompanha o FFmpeg     | mede duração e resolução    | usado em `cortes.py`                    |
| Pandoc         | 3.10                   | converte Markdown em Typst  | `legoloop legopress doctor`             |
| Typst          | 0.15.1                 | compila o PDF do livro      | `legoloop legopress doctor`             |
| legoloop       | binário Go em `~/go/bin` | processos e livros          | `which legoloop`                        |
| Node.js        | 18 ou mais             | roda `gerar_cenas_flow.mjs`     | `processo/PROCESSO.md`                  |
| Chromium       | via Playwright         | desenha quadros no render   | `python -m playwright install chromium` |

Versão não verificada: Node.js. A exigência vem do `processo/PROCESSO.md`.

## 7.2 Bibliotecas Python

Arquivo: `ferramentas/requirements.txt`.

- `playwright`: controla o Chromium.
- `numpy` e `scipy`: cálculo de áudio.
- `pillow`: imagens.
- `edge-tts`: voz neural em português.

## 7.3 Serviços externos

| Serviço              | Uso                                    | Custo                                  | Observação                              |
| :------------------- | :------------------------------------- | :------------------------------------- | :-------------------------------------- |
| Google Flow          | imagens, vídeos e fichas de personagem | créditos; o plano grátis dá 50 por dia | consome créditos antes de gerar         |
| edge-tts (Microsoft) | narração                               | grátis                                 | envia o texto ao serviço online         |
| Postiz               | agenda e postagem                      | open source, na VPS                    | conexão de cada rede em `processo/CONECTAR_REDES.md` |
| Earshot              | efeitos sonoros                        | grátis (MIT)                           | pedido em texto, roda em CPU            |

## 7.4 Scripts do projeto

Todos ficam em `ferramentas/`, salvo indicação.

| Script                      | Função                                               |
| :-------------------------- | :--------------------------------------------------- |
| `voz.py`                      | narração por linha e tempos                          |
| `check_timeline.py`            | valida a timeline e mede o ritmo                     |
| `mix.py`                      | trilha, efeitos e volume final                       |
| `engine/render.py`            | render da timeline em 16:9, 9:16 ou 1:1              |
| `marca.py`                    | vinheta, logo e fechamento                           |
| `sheet.py`                    | folhas de contato para revisão                       |
| `key_sprites.py`               | recorta personagem do fundo verde                    |
| `cortes.py`                   | cortes curtos em 9:16                                |
| `miniatura.py`                | miniatura do vídeo                                   |
| `canal.py`                    | kit do canal                                         |
| `srt.py`                      | legendas                                             |
| `align.py`                    | alinhamento de texto e áudio                         |

Fonte: listagem de `ferramentas/` e `ferramentas/LEIAME.md`.

## 7.5 Skills e roteiros de agente

- `legoloop-legopress`: skill que conduz a criação deste livro.
- Os prompts de `processo/PROMPTS.md` são as instruções que o agente lê.

## 7.6 VCP

O termo VCP foi pedido pela pessoa, mas não está definido no código do LegoLoop.
Busca por `VCP` no repositório LegoLoop: nenhum arquivo encontrado.
Definição: não verificado. Falta a pessoa informar o significado.

## 7.7 Dados

| Arquivo                        | Conteúdo                                          |
| :----------------------------- | :------------------------------------------------ |
| `projeto.json`                   | configuração da etapa 0                           |
| `marca.json`                     | kit de marca                                      |
| `roteiro.json`                   | linhas narradas e fontes                          |
| `pronuncia.json`                 | grafia fonética de nomes estrangeiros             |
| `timeline.json`                  | montagem das cenas                                |
| `audio/timings.json`             | tempo de cada linha                               |
| `cortes/<video>.json`            | especificação de cada corte                       |
| `redes.json`                     | limites de legenda por rede                       |
