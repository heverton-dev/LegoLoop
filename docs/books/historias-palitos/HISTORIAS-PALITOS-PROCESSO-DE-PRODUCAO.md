---
title: "Histórias Palitos"
subtitle: "O processo completo de produção de vídeo animado, do nicho ao YouTube"
author:
  - Projeto Histórias Palitos
date: "23/11/2026"
lang: pt-BR
institute: "Bilhões em Jogo"
eyebrow: "PROCESSO DE PRODUÇÃO"
cover-title: "Histórias Palitos"
cover-accent: "Do nicho ao YouTube, etapa por etapa"
cover-caption: "ENTRADAS · PROCESSAMENTO · RENDER · SAÍDAS"
header-mark: "HISTÓRIAS PALITOS"
accent: "2563eb"
tagline: |
  Um processo que qualquer agente de IA executa com os mesmos arquivos.
toc: true
toc-depth: 2
abstract: |
  Este livro descreve como o projeto Histórias Palitos cria um vídeo.
  Mostra o que entra, o que cada etapa faz, o que sai e quais peças fazem o trabalho.
  Também registra o que foi testado e o que ainda falta.
---


# 1. Visão geral

O projeto transforma um tema em um vídeo animado de canal.
O vídeo usa bonecos de palito, cenários, voz neural em português e trilha.
Todo o processo está descrito em `processo/PROCESSO.md` no repositório.

## 1.1 Princípio

O processo não depende de um modelo de IA nem de um ambiente.
Qualquer agente que leia o processo e os prompts consegue executá-lo.
Cada decisão fica em arquivo de texto (JSON ou Markdown).
Assim, outro agente retoma o trabalho só com a pasta do projeto.

Fonte: `processo/PROCESSO.md`, seção "Regra inegociável".

## 1.2 Quem faz o quê

Três atores executam o processo.

- **Pessoa (dona do canal).** Escolhe nicho, formato e duração. Aprova cada portão.
- **Agente de IA.** Pesquisa, escreve o roteiro, monta os arquivos e roda os scripts.
- **Serviços externos.** Google Flow gera imagens e vídeos. edge-tts gera a voz.

## 1.3 Dois modos de execução

No modo A, o agente tem terminal e navegador.
Ele roda os scripts e opera o Google Flow.

No modo B, o agente só conversa.
Ele escreve os arquivos e os prompts.
A pessoa roda os scripts e cola os prompts no Flow.

Os dois modos geram os mesmos arquivos e o mesmo vídeo.
Fonte: `processo/PROCESSO.md`, seção "Dois modos de execução".

## 1.4 Mapa do processo

O fluxo tem catorze etapas, de 0 a 13, com portões de aprovação.

```
Etapa 0  configuração  →  1 pauta  →  2 roteiro  →  3 elenco
   →  4 cenários e estilo  →  5 storyboard  →  6 voz
   →  7 animação  →  8 montagem  →  8b cortes curtos
   →  9 revisão  →  10 vitrine  →  11 publicação  →  12 comentários  →  13 aprendizado
```

O detalhamento de cada etapa está no capítulo 3.


# 2. Entradas

Entrada é tudo que o processo recebe antes de produzir qualquer coisa.
Há entradas de decisão e entradas de dados.

## 2.1 Entradas de decisão

Estas vêm da pessoa, na etapa 0.

| Entrada                | Exemplo                                | Onde fica           |
| :--------------------- | :------------------------------------- | :------------------ |
| Nicho                  | finanças, odontologia                  | `projeto.json`        |
| Estrutura narrativa    | explicador, investigação, lista        | `projeto.json`        |
| Formato                | 16:9, 9:16, 1:1                        | `projeto.json`        |
| Duração                | de 30 s a 20 min                       | `projeto.json`        |
| Voz                    | Thalita, Francisca, Antonio            | `projeto.json`        |
| Marca                  | fontes, cores, logo                    | `marca.json`          |
| Estilo visual          | dioramas, palitos, whiteboard          | `projeto.json`        |
| Animação               | motor em código, Flow ou híbrido       | `projeto.json`        |
| Distribuição           | canais e volume por dia                | `projeto.json`        |

Fonte: `processo/PROCESSO.md`, seção "Etapa 0".

## 2.2 Entradas de dados

Estas nascem durante o processo e alimentam a montagem.

| Entrada               | Conteúdo                               | Quem gera            |
| :-------------------- | :------------------------------------- | :------------------- |
| `roteiro.json`          | linhas narradas e fontes com link, data e nível | agente, etapa 2      |
| `plano.md`              | batidas, elenco, cenários, linha a linha | agente, etapas 2 a 4 |
| `prompts/*.txt`          | prompts prontos para o Flow            | agente               |
| `flow/`                 | imagens e vídeos baixados do Flow      | Google Flow          |
| `pronuncia.json`        | grafia fonética de nomes estrangeiros  | agente, etapa 6      |
| `acervo/temas.json`     | temas já feitos, para não repetir      | agente, etapa 1      |

## 2.3 Regra de entrada

Nada entra sem fonte.
Cada número do roteiro tem link, data e nível A ou B em `roteiro.json`.
Nível C serve só para achar pistas.
Fonte: `processo/PROCESSO.md`, etapa 2, e `processo/FONTES.md`.

## 2.4 Exemplo real

O projeto `curiosidades-malucas` tem estes arquivos de entrada na própria pasta:
`projeto.json`, `marca.json`, `roteiro.json`, `script.json`, `cenarios.json`, `elenco.json`, `storyboard.json` e `pronuncia.json`.
Verificado por listagem da pasta `curiosidades-malucas/`.


# 3. Etapas

Cada etapa lê entradas, faz um trabalho e gera saídas.
Ao fim, há um portão.
Sem aprovação, o agente não avança.
Fonte: `processo/PROCESSO.md`, seção "Portões de aprovação".

| Etapa          | Nome              | Entrada             | Trabalho                               | Saída                      | Portão         |
| :------------- | :---------------- | :------------------ | :------------------------------------- | :------------------------- | :------------- |
| 0              | Configuração      | respostas da pessoa | pergunta uma a uma                     | `projeto.json`               | sim            |
| 1              | Pauta             | nicho e fontes      | busca tendências e demanda             | 10 ideias, 3 marcadas      | escolha        |
| 2              | Roteiro           | pauta aprovada      | pesquisa em fontes A ou B e escreve em batidas | `roteiro.json`, batidas      | sim            |
| 3              | Elenco            | roteiro             | fichas do personagem no Flow           | personagem salvo           | sim            |
| 4              | Estilo e cenários | elenco              | prompt mestre, cenários, objetos       | prompts e quadros de teste | sim            |
| 5              | Storyboard        | cenários e elenco   | um quadro por cena                     | `storyboard.json`            | sim            |
| 6              | Voz               | roteiro             | narração por linha com tempos          | `audio/timings.json`         | checagem       |
| 7              | Animação          | storyboard          | motor em código ou vídeo do Flow       | cenas                      | cena teste     |
| 8              | Montagem          | cenas, voz, marca   | timeline com marca e legendas          | `timeline.json`              | não            |
| 8b             | Cortes curtos     | timeline            | recorta trechos em 9:16                | `renders/shorts/`            | sim            |
| 9              | Revisão           | render              | folhas de contato e revisor            | "enviar" ou "não enviar"   | sim            |
| 10             | Vitrine           | vídeo aprovado      | títulos, capas, descrição              | pacote                     | sim            |
| 11             | Publicação        | acervo              | agenda no Postiz                       | postagem                   | autorizada     |
| 12             | Comentários       | comentários         | rascunhos de resposta                  | respostas                  | sim            |
| 13             | Aprendizado       | números             | recomendações semanais                 | `acervo/aprendizado.md`      | não            |

## 3.1 Etapa 0: a régua da duração

A duração define quantos cortes e cenas o vídeo precisa.
A regra dos 3 segundos vale para o vídeo longo: um corte ou movimento a cada 3 s.
Em 9:16 e 1:1, a regra vira 2 a 2,5 s.
Fonte: `processo/PROCESSO.md`, tabela de ritmo.

## 3.2 Etapa 2: regras do roteiro

Abra com o momento mais dramático.
Ligue as batidas com "mas" ou "portanto", nunca com "e então".
Use linhas de 8 a 25 palavras.
Separe fato de alegação com "segundo os promotores".
Fonte: `processo/PROCESSO.md`, etapa 2.

## 3.3 Etapa 3: o elenco vem primeiro

Cada personagem tem três fichas no Flow: referência, expressões e poses.
A descrição é igual em todos os prompts do mesmo personagem.
O personagem aprovado vai para o recurso "Personagens" do Flow.
Fonte: `processo/PROCESSO.md`, etapa 3.

## 3.4 Etapa 7: três formas de animar

A: motor em código, grátis e local.
B: vídeo do Flow por cena, que consome créditos.
C: híbrido, com vídeo do Flow só nos momentos-chave.
Antes de tudo, faça uma cena teste de 10 s.
Fonte: `processo/PROCESSO.md`, etapa 7.


# 4. Processamento

Processar é transformar as entradas em material pronto para montar.
Esta parte mostra os scripts e a ordem em que rodam.
Os caminhos são de `ferramentas/` no repositório `proj--historias-palitos`.

## 4.1 Voz

`voz.py` lê as linhas do roteiro.
Gera uma narração por linha com edge-tts.
Grava os tempos de cada linha em `audio/timings.json`.
Comando: `python voz.py <projeto> --voz pt-BR-FranciscaNeural`.
Aviso: edge-tts envia o texto ao serviço online da Microsoft.
Fonte: `ferramentas/LEIAME.md`, linha 7; `processo/PROCESSO.md`, etapa 6.

## 4.2 Validação da timeline

`check_timeline.py` valida `timeline.json`.
Também mede o ritmo em cortes por minuto.
Fonte: `ferramentas/LEIAME.md`, linha 8.

## 4.3 Legendas

`srt.py` e `gerar_legendas.py` geram legendas a partir de `data/lines.json`.
No vídeo longo, a legenda sai em `.srt`.
No 9:16, a legenda é queimada no vídeo, em blocos de 3 palavras.
Fonte: `ferramentas/cortes.py` (docstring) e `ferramentas/LEIAME.md`, linha 17.

## 4.4 Trilha e volume

`mix.py` monta a trilha e os efeitos.
O volume final vai para -14 LUFS, com pico em -1,5 dB.
Comando: `python mix.py <projeto> --music dark`.
Fonte: `ferramentas/LEIAME.md`, linha 9; `processo/PROCESSO.md`, etapa 8.

## 4.5 Montagem da timeline

A timeline é um arquivo JSON com a lista de cenas (`shots`).
Cada cena tem duração, linhas de texto, cenário e personagem.
Scripts como `gerar_timeline_completa.py` montam esse arquivo.
Fonte: `curiosidades-malucas/gerar_timeline_completa.py` e `build_timeline.py` na raiz.

## 4.6 Recorte de cortes curtos

`cortes.py` lê uma especificação em `cortes/<video>.json`.
Cada corte é um intervalo de linhas do roteiro.
O script recorta cenas, voz e trilha, refaz a legenda e acrescenta um cartão final.
Não grava nada de novo.
Fonte: `ferramentas/cortes.py`.


# 5. Render

Render é a etapa que transforma a timeline em arquivo de vídeo.
A timeline descreve o que aparece em cada instante.
O render desenha cada quadro e junta tudo num arquivo MP4.

## 5.1 O que entra

- `timeline.json`: lista de cenas, duração e texto.
- Imagens de cenário e personagens, em `flow/` e `sprites/`.
- Áudio de voz e trilha, em `audio/`.
- `marca.json`: fontes, cores e logo.
- Formato: 16:9, 9:16 ou 1:1.

## 5.2 O que o render faz

1. Abre o navegador Chromium, controlado pelo Playwright.
2. Cada cena vira um clipe curto. O log mostra as cenas em paralelo.
3. Os clipes ficam em `.render_cache/`, com `list.txt` e `silent.mp4`.
4. O FFmpeg junta os clipes e mistura o áudio.
5. O arquivo final vai para `renders/`.

Comando: `python engine/render.py <projeto>/timeline.json --formato 9:16 --out x.mp4`.
Fonte: `ferramentas/LEIAME.md`, linha 10.

## 5.3 Evidência do teste

Render do corte `02-invasao-muscular` (9:16).
Saída: 1080 x 1920 pixels, 31,7 s, em 1 min 20 s.
O log mostrou 5 cenas em paralelo e a linha "Done".
Arquivo: `renders/shorts/02-invasao-muscular.mp4`.
Verificado com `ffprobe` na própria execução do `cortes.py --render`.

## 5.4 Aviso sobre fontes

Em duas cenas do teste, o log mostrou "marca: falhou" para as fontes.
O render terminou mesmo assim.
A causa não foi verificada.
Hipótese a testar: falha ao carregar fontes do Google Fonts.

## 5.5 Pós-produção

`marca.py` adiciona vinheta de abertura, logo e vinheta de fechamento.
Comando: `python marca.py <projeto> x.mp4 --out final.mp4`.
Fonte: `ferramentas/LEIAME.md`, linha 11.


# 6. Saídas

Saída é o que o processo entrega ao final de cada etapa ou do vídeo inteiro.

## 6.1 Arquivos de vídeo

| Saída                     | Onde fica                              | Quando sai      |
| :------------------------ | :------------------------------------- | :-------------- |
| Vídeo longo (v1, v2, v3)  | `renders/`                               | etapa 8 e 9     |
| Corte vertical 9:16       | `renders/shorts/<id>.mp4`, 1080 x 1920   | etapa 8b        |
| Legenda `.srt` e `.vtt`       | pasta do projeto                       | etapa 8         |
| Trilha final              | `audio/`                                 | etapa 8         |

Regra: nenhuma versão é apagada.
Correções geram nova versão.
Fonte: `processo/PROCESSO.md`, etapa 9.

## 6.2 Arquivos de publicação

| Saída                   | Onde fica          | Conteúdo                               |
| :---------------------- | :----------------- | :------------------------------------- |
| `publicacao.md`           | pasta do projeto   | títulos, descrição, capítulos e miniaturas |
| Pacote da vitrine       | `publicacao.md`      | cinco títulos com nota e três capas    |
| `acervo/acervo.json`      | `acervo/`            | vídeos aprovados, validade e status    |

Fonte: `processo/PROCESSO.md`, etapas 10 e 11.

## 6.3 Saídas de controle

- `review/`: folhas de contato e relatórios de revisão.
- `acervo/aprendizado.md`: recomendações semanais.
- `cortes/<video>.json`: especificação de cada corte, com título, descrição e status de envio.

## 6.4 Destino

A publicação passa pelo Postiz, instalado numa VPS.
O Postiz envia para YouTube, Instagram, TikTok e LinkedIn.
Fonte: `processo/PROCESSO.md`, etapa 11.


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


# 8. Verificação e limites

Este capítulo registra o que foi testado e o que ainda falta.
Cada teste mostra o comando e o resultado.

## 8.1 Testes executados

| Teste                              | Comando                                    | Resultado                           |
| :--------------------------------- | :----------------------------------------- | :---------------------------------- |
| Dependências                       | `legoloop legopress doctor`                  | pandoc e typst disponíveis (exit 0) |
| Preparação dos cortes first-brands | `python ferramentas/cortes.py cortes/first-brands.json` | 5 cortes preparados (exit 0)        |
| Preparação dos cortes curiosidades | `python ferramentas/cortes.py cortes/curiosidades-malucas.json` | 3 cortes preparados (exit 0)        |
| Render de um corte                 | `cortes.py --render --so 02-invasao-muscular` | 1080 x 1920, 31,7 s, 1 min 20 s     |

## 8.2 Lacunas encontradas

- `cortes.py` diz que confere o gancho. O código não confere o gancho (`ferramentas/cortes.py`).
- O JSON de curiosidades-malucas não tem link do vídeo completo nem status de envio (`cortes/curiosidades-malucas.json`).
- O corte 03-mordida-25cm ainda não foi renderizado.
- A correção de caminhos em `ferramentas/cortes.py` não está commitada.
- O comando `legoloop legopress criar` só imprime uma linha. Não cria a pasta do livro.

## 8.3 Limites conhecidos

- Os cenários foram desenhados em 16:9. No vertical, alguns elementos cortam ou ficam perto da legenda.
- A trilha do vídeo longo serve de base para os cortes. Há um fade curto na emenda.
- Cada corte leva alguns minutos para renderizar. Não cobre publicação real nem geração no Flow.


# 9. Glossário

**Batida.** Frase curta do roteiro que pede uma imagem.

**Cartão final.** Tela preta com texto de chamada, no fim de cada corte.

**Corte (de cena).** Troca de plano na animação. Regra: um a cada 3 s.

**Corte curto (Short).** Trecho vertical de 15 a 59 s, recortado do vídeo longo.

**Flow.** Google Flow, o gerador de imagens e vídeos usado no projeto.

**Folha de contato.** Imagem com quadros do vídeo, para revisão rápida.

**Portão.** Ponto em que a pessoa aprova antes do próximo passo.

**Render.** Transformar a timeline em arquivo de vídeo.

**Timeline.** Arquivo JSON com a lista de cenas, a voz, a trilha e a marca.

**VCP.** Termo pedido pela pessoa. Definição não verificada.


