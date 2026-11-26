# Análise das referências a partir das transcrições (2026-11-26)

Fonte: `projetos/refs/transcricoes.md` (14 transcrições, 43.736 palavras). Lidas na íntegra: Benjamín Cordero, Futurepedia, Matt Penny (motion graphics), Kaye, ProvenTube, Jesse, Eric Michaud, ThreatMind. Lidas por capítulos e termos: MAGI-2, Seedance 4K, Maestros da IA (2), Battisti, Matt Penny (1.000 Reels).
Nada aqui vem de assistir aos vídeos. Números de receita e de views são afirmações dos autores, sem verificação.

## 1. Quatro famílias de técnica

| Família | Quem | Como | Custo | Determinismo |
|---|---|---|---|---|
| A. Código (sem modelo de vídeo) | ProvenTube, Eric Michaud, Futurepedia, Benjamín (claude-animated-shorts, lemo-opuscar, pixelfx) | Claude Code escreve cena em Remotion + Three.js; personagens 2D num mundo 3D; FFmpeg e pocketsphinx para áudio e tempo de cada palavra | 0 crédito de vídeo; custo em tokens do Claude | Alto. Eric: "é o que conserta a slot machine" |
| B. Modelo gerativo | Kaye (Higgsfield MCP + Gemini Omni Flash 1.2), Rourke (MAGI-2), Griffin (Seedance 4K), Jesse (ChatGPT + Google Flow + CapCut), ThreatMind (Claude + Flow MCP) | Imagem de personagem, depois clipe de 8 a 15 s por cena | Créditos (ver 3) | Baixo. Kaye achou "terceira mão" que aparece sozinha e estilo que muda entre cenas |
| C. Híbrido | Eric Michaud | Animação em código como "wireframe", depois restilização vídeo-para-vídeo (Higgsfield Genjutsu, 100 créditos em 1080p) que mantém o movimento | Poucos créditos | Médio, com piso garantido pelo código |
| D. Processo de fábrica | Maestros da IA, Matt Penny | Linha editorial primeiro, etapas com conferência entre elas, cortes dos longos para curtos, modelo de vídeo reaproveitável | Tokens e ferramentas | Alto |

**Nosso motor é da família A.** A conclusão do Eric (determinismo contra slot machine) é o mesmo argumento que usamos.

## 2. O que cada vídeo ensina

- **Benjamín Cordero (Top 10 repos):** Blender MCP (34 mil estrelas, só com GPU; cena de 1 h 10 min no Mac, 47 min numa GPU alugada); lemo-opuscar (cerca de 100 estilos em código, 42 s no Mac); claude-animated-shorts (2D num mundo 3D, tudo em código, WebGL, 25 s no Mac); pixelfx (244 efeitos de pixel art, 11 s); react-three-rapier (física com 10.000 bolinhas, GPU); huashu-motion (aponta o quadro ruim e muda só ele); rea-video (desmonta um vídeo plano a plano: câmera, ritmo e transição, 11 s em CPU); HyperFrames (HTML para MP4, mais de 82 mil estrelas); Earshot (dá "ouvidos" ao Claude e traz mais de 200 efeitos sonoros livres). Regra do autor: HTML, papel e pixel rodam em CPU; luz 3D, física e modelo de vídeo pedem GPU. Aviso dele: se todos usarem o mesmo repositório, todos terminam com o mesmo vídeo; o valor está nas próprias skills.
- **Futurepedia (Three.js):** Claude descreve um mundo (caixa de argila, lâmpada, câmera) e o Three.js calcula luz, sombra e foco. Difere do método 2D quadro a quadro. Atlas grátis com 37 estilos. Sincroniza a animação com a locução por SRT (Gemini TTS). Para marca: pedir ao Claude que transforme o estilo em skill.
- **Matt Penny (motion graphics sem cara de IA):** 4 passos: (1) repos e skills certos (unslop, impeccable, taste, Three.js, react-three-rapier); (2) inspiração real, analisando vídeos de referência; (3) **video manifest**, um documento com cores, fontes, processo e ferramentas do autor; (4) **laço de iteração**: comentários no vídeo viram regras no manifest, e o erro não se repete.
- **Matt Penny (1.000 Reels):** reconstrói um vídeo de referência como modelo, separa partes recombináveis (perguntas, respostas, B-roll, voz) e renderiza em lote.
- **Kaye:** compara 3 caminhos. Blender via Claude gasta 200 a 300 mil tokens por vídeo. Claude Code sozinho dá resultado básico. Claude + Higgsfield MCP + Omni Flash 1.2: 15 s em 1080p custam 48 créditos (720p: 33), contra 240 e 105 no Seedance. Ficha de personagem em 3 painéis primeiro. O Claude **vê o resultado e corrige**. Total declarado: 336 créditos o 2D e 290 o 3D.
- **ProvenTube:** Claude Code (Remotion, Three.js, pocketsphinx, FFmpeg) com pacote pronto. Opus rende bem mais que Sonnet. Vídeo de 1 min leva 30 a 40% de uma sessão de 5 h no Opus. A primeira geração monta uma biblioteca de ativos e aprende para as próximas.
- **Jesse:** ChatGPT grátis com master prompt em estágios (nicho, ideia, duração e formato, personagem, storyboard, animação). Google Flow: imagem grátis (Nano Banana 2.1), vídeo 8 s em 720p, 50 créditos. Montagem no CapCut. Fazer cada cena do storyboard separada deu resultado melhor que o storyboard em lote.
- **Eric Michaud:** Claude Code desenha tudo, sem gerador de vídeo. Ficha de personagem, adereços e eras prontos antes, depois a cena só reaproveita. QC obrigatório (um cavalo ficou chato na versão 1). Restilização no Higgsfield Genjutsu mantém o movimento.
- **ThreatMind:** Claude opera o Google Flow por MCP. Regras de continuidade da master prompt: uma ação por clipe, começa no último frame do anterior, cenário e luz iguais, descrição completa do personagem repetida. Claude revisa frames e corta defeitos. Conta free: 50 créditos por dia, cerca de 15 por clipe de 10 s, ou 30 s por dia. Retoma no dia seguinte por arquivo de estado.
- **Rourke (MAGI-2):** modelo aberto, US$ 0,08 por segundo, até 30 s em 1080p, até 30 referências. Fraquezas dele: cor, tremor entre quadros, pele, espera. Só roda em computador enorme ou no Higgsfield.
- **Griffin (Seedance 4K):** rascunho em 480p custa 15 créditos, 1080p 80, 4K 140. Método: rascunhar barato e pagar o 4K uma vez. Anúncio de 3 cenas: cerca de 600 créditos.
- **Maestros da IA:** 9 tarefas por vídeo; linha editorial é a parte que quase todos pulam; "tudo é conferido" entre etapas; cortes de vídeos longos viram curtos; "a sala onde nenhuma máquina entra" (decisão humana).
- **Battisti (Code Mode):** recurso experimental do Claude Code; gastou menos token e saiu melhor, porém mais lento. Marginal para nós.

## 3. Comparação de custo (cada fonte usa uma moeda diferente; não somar)

- Flow, medido por nós: Veo 3.1 Lite 10 créditos por 8 s; Omni 1.1 Flash 12 créditos por 8 s; imagens 0.
- Flow, segundo ThreatMind: cerca de 15 créditos por 10 s na conta free.
- Higgsfield, segundo Kaye: Omni Flash 1.2 a 48 créditos por 15 s em 1080p.
- Higgsfield, segundo Griffin: Seedance 4K 140 créditos.
- MAGI-2: US$ 0,08 por segundo (Rourke).

## 4. O que agrega ao nosso processo (para decidirmos juntos)

| # | Ideia | Fonte | Esforço | Custo | Risco |
|---|---|---|---|---|---|
| 1 | Camada Three.js no motor (2D num mundo 3D, câmera com perspectiva) | Futurepedia, ProvenTube, Benjamín | Médio | 0 | WebGL em CPU lento; hash dos quadros pode variar entre máquinas. Não verificado na RTX 3050 de 6 GB |
| 2 | Manifesto de vídeo (estilo, cores, fontes, processo) + laço de iteração que grava erros | Matt Penny | Baixo | 0 | Baixo. Já temos `marca.json` e `achados.md` |
| 3 | Revisão visual dos frames pelo agente | Kaye, ThreatMind, Eric | Depende | 0 | **Nosso maior buraco:** neste ambiente eu não leio imagem |
| 4 | Restilizar o render do motor pelo Flow (vídeo-para-vídeo), igual ao Genjutsu | Eric | Baixo | `flow_edit` custa cerca de 20 créditos por 4 s (descrição da ferramenta) | Estilo pode variar entre cenas |
| 5 | Desmontar vídeo de referência plano a plano (rea-video) | Benjamín | Baixo | 0 | Exige baixar o vídeo |
| 6 | Efeitos de pixel (pixelfx) e estilos (lemo-opuscar) como biblioteca | Benjamín | Médio | 0 | Não verificados; licenças a conferir |
| 7 | Física 2D (Matter.js) ou 3D (react-three-rapier) para moedas, balança | Benjamín, Matt Penny | Médio | 0 | GPU no caso 3D |
| 8 | Cortes verticais dos longos | Maestros da IA | Já planejado | 0 | É o TICK-DR-05.3, hoje pausado |
| 9 | Áudio com efeitos sonoros sincronizados (Earshot) | Benjamín | Médio | 0 | Não verificado |

## 5. O que não agrega agora

- **MAGI-2, Seedance e Higgsfield:** pagos ou exigem computador enorme; o nosso Flow já cobre o vídeo gerativo.
- **Remotion:** duplica o nosso motor e tem licença própria (não li os termos).
- **Blender MCP:** 200 a 300 mil tokens por vídeo e GPU pesada.
- **Pacotes pagos e promessas de receita:** sem prova.

## 6. Lacunas do nosso processo reveladas

1. Sem revisão visual do que geramos (só métricas). O rig de peças falhou por isso (A9).
2. Sem manifesto de estilo único que o agente releia a cada cena.
3. Sem biblioteca de ativos reaproveitável entre vídeos (ficha, adereços, eras): Eric e ProvenTube fazem isso desde o primeiro vídeo.
4. Personagens: todos os autores de sucesso usam **2D plano num mundo 3D**, ou stickman. Nenhum monta recortes de IA com juntas.

## 7. Decisões que dependem do usuário

1. Adotar Three.js no motor, com teste 2.5D de 12 s?
2. Criar o manifesto de estilo e a biblioteca de ativos do canal?
3. Como resolver a revisão visual (sua revisão por contato, ou uma ferramenta que devolva imagem ao agente)?
4. Testar a restilização pelo Flow (`flow_edit`) num trecho de 4 s do render do motor?
