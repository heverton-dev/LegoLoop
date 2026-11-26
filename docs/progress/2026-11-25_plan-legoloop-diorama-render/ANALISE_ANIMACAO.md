# Análise: por que o render do diorama ficou estático e como corrigir

Data: 2026-11-25. Status: análise para decisão do usuário. Nenhum código foi alterado depois do TICK-DR-05.2.

## 1. Causa raiz (provada)

1. O `render.go` do TICK-DR-05.2 monta cenário + sprite estático + legenda com `zoompan` do ffmpeg (`internal/diorama/render/render.go`, função `expressaoCamera`). Resultado: a câmera balança e nada se move dentro da cena.
2. O motor original do "Bilhões em Jogo" é outro tipo de coisa: HTML/SVG/JS desenhado quadro a quadro no Chromium via Playwright (`proj--historias-palitos/ferramentas/engine/{engine.html,puppet.js,layers.js,sets.js,render.py}`).
3. O motor original tem bonecos com pose, humor e fala (`puppet.js:1-4`; 17 poses e 12 humores), camadas animáveis com entrada, saída, loop e keyframes (`layers.js:1-8`), câmera por keyframes com easing, luz e efeitos (`lab-flow/timeline.json`: `cam`, `fx`, `light`) e cache por cena.
4. O vídeo completo original tem 127 cenas (`timeline.json` da raiz) e 459 s (`renders/draft_thalita.mp4`, ffprobe). O laboratório tem 10 cortes de 3 s com 10 fundos diferentes (`lab-flow/plates/p01..p10`, `lab_flow_30s.mp4`, 30 s).
5. O Mapa Macro previa "Motor Diorama Autoral" na Peça 3. Eu reduzi isso a ffmpeg + zoompan. Foi erro meu de leitura, não limite técnico.
6. O demo declara `estilo_visual: "hibrido"` (`diorama-demo/projeto.json`, linha 15). O render ignorou isso.
7. "Sempre o mesmo fundo": o demo tem 3 cenários (`flow/cen_01..03.png`) para 17 batidas. O lab usa 10 fundos em 30 s.
8. "Personagem estático": o sprite é um PNG único por pose, sem boca nem gesto no render. Os visemas existem (`data/timings.json`, TICK-DR-02.3), mas nada os consome.

## 2. Três caminhos de animação (já definidos no processo original)

| Modo | Como | Custo | Determinismo |
| :--- | :--- | :--- | :--- |
| A. Motor em código | Chromium renderiza fundos, bonecos, objetos, textos, câmera e efeitos | 0 crédito, local | Alto: mesma entrada, mesmo vídeo |
| B. Vídeo do Flow por cena | Imagem para vídeo no Flow | Créditos; teto padrão 25 por cena (`PROCESSO.md:168`) | Baixo: cada geração varia |
| C. Híbrido | Flow só nos momentos-chave (gancho, virada, final); motor no resto | Poucos créditos | Médio, com piso garantido pelo motor |

## 3. Gestão do limite do Flow (conta grátis)

- Plano grátis: 50 créditos por dia (`PROCESSO.md:225`, e dito pelo usuário). Não verifiquei o valor atual no site.
- Imagens em lote pelo modo Agente custam 0 crédito segundo o processo (`PROCESSO.md:150`). Não verifiquei hoje.
- Conta: com teto de 25 por clipe, 50 créditos dão 2 clipes por dia.
- Demo de 17 batidas: modo B = 17 clipes = 425 créditos = 9 dias. Modo C com 3 clipes-chave = 75 créditos = 2 dias.
- Conclusão: B puro é inviável na conta grátis. A é o piso universal; C é o aprimoramento opcional.
- Usuário sem conta: o processo cai em A sem perder nenhuma etapa e entrega o mesmo contrato de arquivos.

Gestão proposta (`flow_creditos.json`, dentro do projeto):
- Saldo do dia, teto por clipe, gasto acumulado e data do último reset, preenchidos pelo usuário ou lidos da tela do Flow.
- Plano de gasto calculado antes de gerar: lista de clipes por prioridade (gancho, virada, final), custo previsto e dia previsto.
- Portão de créditos: nada é gerado sem o custo aparecer e o usuário aprovar (já existe no processo, `PROCESSO.md:225`).
- Retomada: clipe já gerado nunca é pedido de novo; o plano continua no dia seguinte.
- Falha ou falta de crédito: a cena usa o motor (A). O vídeo sai completo em qualquer cenário.
- Risco: automatizar a interface web do Flow pode violar os termos do Google (aviso do próprio `flow.py`). A geração fica opcional e manual-assistida.

## 4. Requisito agnóstico (Claude Code, Antigravity, OpenCode, Mimocode; Windows, Linux, macOS)

- O motor A roda com Python, Playwright e Chromium. O Playwright já está instalado aqui e o `prevoo.go` já exige Python.
- O `go.mod` não tem dependências. Um cliente CDP em Go puro exigiria WebSocket próprio. Recomendo manter o Playwright em Python atrás da interface `Executor`, como já fazemos com edge-tts e faster-whisper.
- Para igualdade entre sistemas: fontes embutidas (o original tem `fonts_embedded`), versão do Chromium registrada no relatório, cache por hash (o original já usa `ENGINE_HASH`).
- Prova de igualdade: render de uma cena de referência com hash dos quadros comparado entre sistemas. Não executado ainda.

## 5. O que precisa ser refeito

Peça 2 (ativos, ciclo já fechado):
- Cenários: no mínimo 1 fundo distinto a cada 2 a 3 batidas (hoje 3 para 17).
- Storyboard: cada batida passa a emitir um "shot" com fundo, camadas (texto, documento, dinheiro, carimbo...), boneco (pose, humor, fala) e câmera por keyframes.
- Sprites do narrador: conjunto de poses e bocas por visema (no mínimo fechada e aberta; ideal 4), gerados no Flow em imagem (0 crédito).

Peça 3 (render, ciclo em andamento):
- Mantém: voz, Whisper, visemas, legendas ASS/SRT, mixagem, gate de loudness, encoder com fallback, versionamento e gate de integridade (tickets 01.1 a 05.2 continuam válidos).
- Troca: a composição ffmpeg/zoompan do `render.go` passa a chamar o motor em código por cena, com cache. O ffmpeg fica só para montar clipes, áudio e legenda.
- `timeline.json` ganha o shot por batida. Os critérios AC-E03 e AC-E05 precisam de emenda na SPEC e o PLAN é regenerado.
- Cortes 9:16 (05.3): re-render pelo motor com a janela do narrador, mais simples que recortar vídeo pronto.

## 6. Limites desta análise

- Não executei o motor original nesta sessão.
- Não assisti ao "Bilhões em Jogo"; só medi arquivo e li código.
- Não verifiquei créditos, custo por clipe, nem regras atuais do Flow.
- Não verifiquei se o Flow gera boca sincronizada com a voz.

## 7. Decisões que dependem do usuário

1. Adotar o modo C (motor em código como base, Flow opcional nos momentos-chave) como padrão do produto.
2. Boneco do narrador: sprite do Flow com poses e bocas por visema, ou boneco desenhado em código (RobbiSkill, já gesticula e fala).
3. Abrir emenda de Concepção (SPEC v2.0 da Peça 3) e reabrir a Peça 2 para cenários e shots.

## 8. Adendo: o MCP `flow-mcp` e o modo só-Flow (análise de 2026-11-25)

Fontes: README do repositório `GTAI-1/flow-mcp` (baixado e lido) e a lista de ferramentas do servidor `flow-mcp`, que já está registrado no cofre do LegoLoop (`mcp_search`). A skill `legoloop-diorama-ativos` já o cita (SKILL.md, linha 33).

### 8.1 O que o flow-mcp faz

- Dirige o Chrome do Flow por CDP (porta 9333), com perfil próprio. Não digita senha.
- Gera imagem (0 crédito) e vídeo; cria personagem reutilizável (`flow_character`, `asset:<nome>`); encadeia clipes pelo último quadro (`flow_continue`); baixa o que já existe sem gastar (`flow_download`); junta clipes com ffmpeg (`flow_assemble`).
- Preços por vídeo, segundo a descrição da ferramenta `flow_generate`: Omni 1.1 Flash 7 a 12, Veo 3.1 Lite 10, Fast 20, Quality 100. Imagens: 0.
- Antes de gastar, lê a cota do próprio Flow e pula a cena se passar de `max_credits` (padrão 25).
- Roda uma geração por vez, com pausa aleatória de 25 a 70 s.
- Tem painel local (`npm run studio`, porta 8787) com abas Cast, Frames, Shots, Voice, Restyle, Cut, Captions e Library.
- Também faz voz (Kokoro ou Gemini) e legendas (whisper.cpp). O `flow_voices` só lista vozes em inglês instaladas no PC; não verifiquei pt-BR. O LegoLoop já tem voz pt-BR (edge-tts) e legendas ASS, então não usamos essas partes.

### 8.2 Limites que afetam o requisito agnóstico

- README: "macOS or Windows 10/11"; o Windows "hasn't been run on a real PC yet". Não há Linux.
- Exige Node 22+, Chrome e uma tela ligada (não roda sem display).
- O README avisa que automatizar o Flow pode violar os termos do Google e quebrar quando a interface mudar.
- Conclusão: o flow-mcp serve como adaptador **opcional** de clipes. A base do produto (motor em código) não pode depender dele.

### 8.3 Quantos clipes o modo só-Flow precisa

Não precisa de 1 clipe por batida. Um clipe de 8 a 10 s cobre várias batidas da mesma cena, com a câmera e o movimento dentro do próprio clipe. A duração máxima por clipe e o custo por duração não foram verificados; o README só mostra o exemplo "10-second shots, 720p, no more than 15 credits each".

Estimativa (preços da ferramenta, clipes de 10 s; NÃO medida):

| Vídeo | Clipes | Créditos com Lite (10) | Com Fast (20) |
| :--- | :--- | :--- | :--- |
| Demo de 53 s, só-Flow | 6 | 60 | 120 |
| 10 min, só-Flow | 60 | 600 | 1200 |
| 53 s, híbrido (3 a 6 clipes-chave) | 3 a 6 | 30 a 60 | 60 a 120 |
| 10 min, híbrido (10 a 15 clipes-chave) | 10 a 15 | 100 a 150 | 200 a 300 |

Leitura:
- Com 1000 créditos por mês, um vídeo de 10 min só-Flow consome de 60% do mês (Lite) a mais que o mês inteiro (Fast). Não escala.
- O híbrido permite de 3 a 10 vídeos longos por mês na conta Pro.
- O modo só-Flow fica indicado para vídeos curtos (até cerca de 1 a 2 min) e para teste de qualidade.
- Sobre as 14 contas: usar contas diferentes para contornar o limite de créditos de um mesmo projeto tende a violar os termos do Google. Recomendo não automatizar rodízio de contas; usar cada conta para seus próprios projetos.

### 8.4 Gestão de créditos proposta

- Orçamento por vídeo em `flow_creditos.json`: teto total, teto por clipe (`max_credits`), modelo, saldo lido do `flow_status` e gasto acumulado.
- Storyboard e quadros-chave em imagem (0 crédito) antes de qualquer vídeo. Só depois de aprovados viram clipe.
- Fila priorizada: gancho, virada, final; depois o resto, até o teto.
- Portão: mostrar custo previsto e pedir "sim" antes de cada lote. Nunca regerar o que já existe (`flow_download` recupera sem custo).
- Cena sem crédito, sem conta ou sem flow-mcp usa o motor em código. O vídeo sai completo sempre.

### 8.5 Resultado do protótipo híbrido (motor em código, sem Flow ainda)

- Arquivo: `projetos/video-diorama/diorama-demo/hibrido/hibrido_v1.mp4` (1920x1080, 30 fps, 1586 quadros, áudio 52,872 s).
- Feito com o motor original (`render.py`, Playwright) lendo uma timeline convertida do demo. Renderizou em cerca de 2 min 33 s.
- Gate: 0 quadros pretos e 0 brancos. Na 1ª versão a transição `fade` do motor gerou 4 quadros pretos; troquei por `focus`.
- 13 fundos distintos (antes: 3), câmera por keyframes, transições, textos animados, dinheiro, documento e carimbo, narrador que alterna pose e movimenta com a fala.
- Limite: a boca do narrador alterna entre `narrador_pose_1.png` e `narrador_expressao_2.png`, que não são bocas de verdade. Falta o conjunto de bocas por visema na Peça 2.
- Limite: os textos e adereços da tela foram escritos à mão para este protótipo. No processo real vêm do storyboard.
- Não verifiquei o visual: não consigo assistir ao vídeo aqui. Extraí 5 quadros em `hibrido/quadros/`.
