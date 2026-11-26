# Achados da execução — legoloop-diorama-render

## A1 — WER de 5% por linha reprova qualquer erro em linha curta (TICK-DR-02.2)
- **Evidência real (2026-11-25):** `edge-tts` Thalita + `faster-whisper` `small`, linha 0 do `diorama-demo` (17 palavras). Referência "atingiu", hipótese "atingir". WER 5,9%. Linha 1 (23 palavras, com números) deu WER 0,0%.
- **Causa:** 1 palavra em menos de 20 reprova o limite de 5% por linha. O erro é de escuta do Whisper, não da fala; `pronuncia.json` não corrige escuta.
- **Estado:** AC-E02.3 e AC-E02.4 passam com dublês. A SPEC v1.1.0 exige `<= 5%` por linha; o demo real não passa nessa regra.
- **Resolvido (Emenda 1 da SPEC v1.2.0):** limite de 15% por linha e 5% global. Prova real: linha 0 com WER 5,9% passou sem regravar; WER global 2,27% (1 erro em 44 palavras).

## A2 — Regravação não muda o áudio do edge-tts (TICK-DR-02.2)
- **Evidência real:** duas sínteses da mesma linha geraram 245.454 bytes cada. O provedor é determinístico para o mesmo texto e voz.
- **Efeito:** as 3 regravações repetem o mesmo áudio e o mesmo WER; custo medido de 26 s na linha 0.
- **Resolvido:** a regravação para quando o áudio novo é idêntico ao anterior (teste `TestWhisperParaEmAudioIdentico`).

## A3 — `-t` do ffmpeg descartava 28 ms de áudio após o `loudnorm` (TICK-DR-04.2)
- **Evidência real:** `soundtrack.wav` do demo com 52,844 s e a timeline com 52,872 s. Voz, `timings.json` e `timeline.json` estavam exatos em 52,872 s. Sem `-t`, a saída do `loudnorm` mantém todas as amostras; com `-t`, perde 28 ms (os timestamps saem deslocados).
- **Correção:** corte por amostra (`apad` e `atrim=end_sample`). Medido: 0 amostras de diferença, no demo e num teste com `ffmpeg` real de 5,872 s.
- **Estado:** resolvido; o vão de 28 ms passava na tolerância de 1 quadro, mas era margem apertada.

## A4 — Sondagem de loudness em mono dava números errados (relato meu corrigido)
- **Causa:** `-ac 2` converte mono para estéreo com 3 dB de atenuação e distorceu a medição do pico.
- **Efeito:** os valores `-14,7` e `-14,3 LUFS` que citei ao concluir o TICK-DR-04.1 não valem. Com estéreo pleno: alvo de pico `-1,5` entrega `-1,5 dBTP` finais (sem folga); alvo `-2,0` entrega `-14,0 LUFS` e `-2,0 dBTP`.
- **Estado:** resolvido; a mixagem usa `pan` e a SPEC registra a medição válida (premissa 14).

## A5 — TICK-DR-05.2: adereços e efeitos ficam fora do vídeo
- Dado real: `flow/ficha_obj_batida_NN.png` é um quadro cheio e opaco (1376x768, alfa 255 em todo pixel). Não existe recorte do adereço, e a timeline não traz `arquivo` para `aderecos` e `efeito`.
- Efeito: o master compõe cenário, sprite do narrador e legenda karaokê. Colar o quadro cheio cobriria o cenário, por isso adereço e efeito sem arquivo são omitidos.
- Decisão a vetar: gerar recortes de adereço no Ciclo 2 (ativos) ou aceitar o master sem adereços.
- Escopo: `ass.go` (legenda karaokê ASS) e `render_exec.go` foram acrescentados além de `render.go`, porque o ticket exige legenda karaokê e a função de execução passava de 35 linhas.
- Não verificado: aparência visual dos quadros e da legenda (sem inspeção visual nesta sessão).

## A6 — Parecer do usuário sobre o protótipo híbrido (`hibrido/hibrido_v1.mp4`)
- Veredito do usuário: "ficou ruim, ficou horrível". Registrado sem interpretação.
- Bom: a narração e as transições.
- Problema 1: o narrador tem só duas imagens que alternam (uma posição ou outra). Fica feio e tira o foco. O correto é animar o personagem, ou criar a cena já animada. Piscar, ir e vir entre duas imagens não serve.
- Problema 2: usei as fichas (`ficha_obj_batida_NN.png`) como fundo. Isso está errado. Ficha não é cenário.
- Forma correta segundo o usuário: no Google Flow, a cena é criada por **referências**. Referencia-se o narrador e o diorama (e a vista desejada) e envia-se um prompt. O Flow compõe a cena e mantém a consistência pelas referências.
- Pergunta em aberto do usuário: como animar os personagens. Ele sugeriu que o Flow gere a cena animada; pediu que eu diga como fazer.
- Status: nenhuma ação tomada. O TICK-DR-05.3 segue pausado. A análise e a cena-teste de 1 clipe não foram autorizadas ainda.

## A7 — Clipe de teste do Veo saiu falando inglês (2026-11-25)
- **Evidência:** o usuário assistiu `flow/teste8s/veo31lite_8s.mp4` e apontou o inglês. O prompt dizia "while speaking" sem fixar idioma.
- **Causa:** erro meu de prompt; o Flow escolheu o idioma.
- **Correção:** a fala vai entre aspas em português e o prompt pede "Brazilian Portuguese". Regeneração (10 créditos) aguarda "sim" do usuário.
- **Decisão do usuário depois:** o vídeo do produto usa voz Antonio + motor, sem fala dentro do clipe do Flow.

## A8 — Li errado a imagem de poses do narrador (2026-11-25)
- **Evidência:** o usuário enviou as prévias do `rembg`. `flow/narrador_poses.png` mostra 2 bonecos de palito (óculos e gravata-borboleta), não 6 poses e não o narrador terracota.
- **Erro meu:** disse que `u2net` e `isnet` "perderam poses" (2 blocos). Os 2 blocos estavam certos; li como falha de modelo.
- **Achado real:** a ficha do projeto está inconsistente com o narrador atual (origem não verificada).
- **Aprendizado:** recorte de folha de ficha (bandas, textos, vários elementos) falha em qualquer modelo (16 a 18 blocos). Corpo único em fundo liso recorta bem (`corpo_rembg.png`, avaliado TOP pelo usuário).
- **Regra:** gerar um elemento por imagem, em verde chapado. `rembg` só como reserva para imagens prontas.

## A9 — Rig de peças entregue sem checagem visual nem geométrica (2026-11-25)
- **Veredito do usuário:** `elenco_pecas_9x16.mp4` "sem condições". Registrado sem interpretação.
- **Erro meu:** montei e entreguei o rig por regra anatômica, sem conferir o resultado.
- **Diagnóstico geométrico (`rig_diag.py`, máscara das peças por instante):** braços soltos do tronco (vão de 2 a 8 px); analista e cidadão sobrepostos na cena 3 (160 px de distância); formiga atravessando as pernas dos personagens; pernas da formiga quase da altura do corpo.
- **Correção (`elenco_pecas_v2_9x16.mp4`):** encaixe pelo centro real de cada peça (pivo no ombro fica fora do centro), braço 40% sob o tronco e atrás dele, pernas da formiga a 60% do tamanho, balanço menor, personagens separados na cena 3.
- **Prova após a correção (42 instantes, 0 a 4,1 s):** formiga sempre 1 bloco; cena 3 sempre 2 blocos (um por personagem); cena 2 com 2 blocos em 37 de 42 instantes (nos outros 5 o braço do analista toca o cidadão).
- **Não verificado:** o visual (encaixe das peças geradas, estilo, orientação das pernas). A geometria conectada não prova que parece bom.
- **Regra:** rig novo passa por diagnóstico geométrico antes de ir ao usuário.

## A10 — Migração do projeto Histórias Palitos ficou periférica (2026-11-26)
- **Fato:** o processo no LegoLoop cobre roteiro e ativos; render virou zoompan; motor e autoria ficaram no projeto original (`ANALISE_ANIMACAO.md` seção 1; `TRANSICAO_PROJETO_PROCESSO.md`).
- **Efeito:** o render do demo saiu estático e as tentativas de animar personagem no `motor-teste` reinventaram o motor em vez de portar o original.
- **Correção em curso:** portar motor + autoria como núcleo do Ciclo 3 e provar reproduzindo First Brands. Regra e checklist em `TRANSICAO_PROJETO_PROCESSO.md`.
- **Não verificado:** que o motor original roda hoje nesta máquina (próximo passo).

## A11 — Fundo 6 veio errado: o Flow entregou a ficha antiga obj_batida_06 (2026-11-26)
- **Fato:** na folha de contato dos 9 fundos, o quadro do Congresso era o "Character/Asset Master Sheet obj_batida_06", sobra de uma geração antiga do projeto Flow. O usuário apontou.
- **Causa:** o flow-mcp baixou a imagem errada para `scene-06`; esse job foi o único sem a nota "renamed in Flow" no retorno do `flow_wait`.
- **Erro meu:** copiei os 9 arquivos sem conferir nota de renomeação nem hash e mandei a folha sem olhar.
- **Correção:** regenerei o Congresso (`scene-10`, 0 crédito); `flow_download` por título trouxe `clip-01.jpeg`; o errado foi para `fundos_lab/descartados/`.
- **Regra:** job de imagem sem a nota "renamed in Flow" é suspeito; conferir o título, o hash e a folha de contato antes de entregar.
- **A11, atualização:** o 2º arquivo também veio errado (`clip-01.jpeg`, maquete fotográfica). `flow_download` pelo título `...-scene-10` baixou o ativo errado. O usuário achou os dois Congressos certos na grade do Flow. Baixei pelos nomes exatos listados em `flow_assets` ("National congress complex at night" = `bg06_congresso.jpeg`, escolhido pelo usuário; "Congress complex at night" = `bg06_congresso_alt_B.jpeg`). Os dois errados ficam em `fundos_lab/descartados/`.
- **Regra:** após baixar do Flow, conferir a média de cor/estilo contra os outros fundos (fundo vetorial escuro tem média RGB baixa; a foto tem média ~144,130,108) e mostrar a folha antes de usar; nunca confiar em `clip-NN` ou título `scene-NN` sem checar.

## A12 — Macro nomeado por estilo e sem inventário (2026-11-26)
- **Fato:** o macro, os 4 processos, os pacotes Go e os comandos levam "diorama" no nome, mas diorama é um estilo. `projeto.json` usa `estilo_visual` para o modo de produção (`cenarios.go:39,41,128,135`). Não existe inventário de partes de personagem, cenários e objetos entre vídeos.
- **Efeito:** trocar de estilo (mangá, pixel art, pré-história) obriga a mexer em nome, contrato e elenco; cada vídeo recria tudo do zero.
- **Correção proposta:** `ANALISE_IMPACTO_INVENTARIO_ESTILOS.md` (inventário global do macro, `modo_producao` separado de `estilo`, pele plugável, motor portado, renomeação com alias). Aguarda as 5 decisões do usuário (seção 6).

## A13 - DAG do fabricar plano nao roda no swarm (2026-11-26)
- Sintoma: swarm_run no DAG piloto devolveu skipped=[TICK-01.1] e levels=null (nenhum worker subiu).
- Causa: LevelsFromDAG (internal/board/swarm_dag.go:68) so sobe no com role/papel; o fabricar plano nao emite role.
- Causa 2: o command do no e o comando de teste (go test), nao um executor que escreva codigo.
- Contorno no piloto: DAG com role engenheiro-tdd e command de agente headless. Correcao definitiva: ticket para o fabricar plano emitir role e command do executor.

## A13 - CORRIGIDO (2026-11-26)
- Correcao aditiva: Ticket.Command (internal/plangen/plangen.go) e routeNode.Command (internal/plangen/route.go) com omitempty; o routedag ja copiava command e ja atribuia role engenheiro-tdd a nos native_swarm.
- Uso: ticket com "executor":"native_swarm" e "command":[agente...] gera ROUTE/DAG que o swarm sobe.
- Provas: internal/plangen/ticket_command_test.go (TestTicketCommandNativeSwarm, TestTicketCommandAditivo); plangen, routedag, dag, board, studio, mcp e route2dag passam.
- Efeito: so tickets novos com command mudam; ticket sem command gera a mesma ROUTE.

## A14 - CORRIGIDO: espinha em internal/diorama/render (2026-11-26)
- Sintoma: 13 globais mutaveis e 4 erros ignorados (`_ = f()`) reprovavam TestCmdEspinhaVerificarRepositorioReal, TestVerificarRepositorioReal e TestGlobaisDaEspinhaEstaoJustificados.
- Correcao: globais viraram funcoes/switch (resolucaoDe, visemaDeLetra, sinteseSFX, cameraValida e outras); erros de remocao de temporario passaram por removerTemp e errors.Join; RegistroPadrao grava o provedor direto.
- Parte do inventario (tiposComuns, tiposDeEstilo) virou switch em tipoPermitido.
- Prova: `legoloop espinha verificar --json` com 0 violacoes; portao completo verde.

## A15 - CORRIGIDO: Processos diorama sem pipeline.dag.json (2026-11-26)
- Sintoma: TestProcessosCincoPorCinco falhava (internal/processo/conformidade_test.go:46).
- Causa: legoloop-diorama-ativos e legoloop-diorama-roteiro gravaram processo.dag.json; a convencao e pipeline.dag.json.
- Correcao: arquivos renomeados, gates e referencias dos docs atualizados, SKILL.md usa o glob *.dag.json (a palavra pipeline e termo proibido nas skills).
- Prova: `legoloop validate` dos dois DAGs com exit 0; internal/processo e internal/assets passam.

## A16 - CORRIGIDO: TestRegistroValida fixava 36 processos (2026-11-26)
- Causa: numero fixo defasado (registro e pastas tem 38).
- Correcao: o teste compara o registro com as pastas de processos/ (internal/espinha/registro_test.go).

## A17 - Lei da prova tornada mecanica (2026-11-26)
- Pedido do usuario: proibido "Nao verificado", pendencia e proxima etapa; verificar ou corrigir na mesma resposta.
- Correcao: ste.LintEvasao (regra evasao-de-prova), `legoloop stept lint [--hook]`, gancho Stop em .claude/settings.local.json, texto distribuido por internal/ste/prompt_prova.go para 8 harnesses, AGENTS.md regras 8 e 10.
- Prova: mutacao da regra e do gancho quebra os testes; sessao headless com frase evasiva bloqueada (error_max_turns, 5 turnos) e controle com prova concluido (success, 1 turno).

## A18 - CORRIGIDO: swarm run nao reexecuta no ja executado (2026-11-26)
- Sintoma: `swarm run` falhava com "worker TICK-01-1 ja existe" mesmo com o worker killed.
- Correcao: internal/board/swarm.go (semWorkerEncerrado) libera o nome de worker killed; worker vivo continua recusando duplicata.
- Prova: TestAddWorkerReusaNomeDeWorkerEncerrado; a mutacao do filtro quebra o teste; a camada 1 rodou duas vezes em seguida.

## A19 - CORRIGIDO: DAG gerado nao apontava o progresso vivo (2026-11-26)
- Sintoma: `swarm run` terminava com exit 0 e gate verdadeiro, mas progresso.json seguia com 24 pendentes.
- Correcao: `fabricar plano` grava o campo progresso no DAG (internal/cli/fabricar_plano_progresso.go).
- Prova: TestFabricarPlanoApontaProgresso; a mutacao quebra o teste.

## A20 - CORRIGIDO: IDs de ticket fora da convencao (2026-11-26)
- Sintoma: mesmo com o ponteiro, o observador nao casava o no TICK-01.1 com o ticket TICK-VA-01.1.
- Causa: o formato aceito e TICK-<iniciativa letras+digitos>-NN.N (internal/progresso/observador.go:14); VA nao tem digitos.
- Correcao: tickets renomeados para TICK-V01-NN.N; PLAN, ROUTE, DAG e progresso regerados.
- Prova: progresso.json com aprovados=1 e TICK-V01-01.1 em PASSED com a saida do go test.

## Camada 1 do video animado (2026-11-26)
- DAG gerado pelo motor (24 tickets com executor native_swarm e command), `legoloop swarm run` no TICK-01.1: exit 0, gate_ok true.

## A21 - gate manual por aspas no comando de validacao (2026-11-26)
- Achado: validation_cmd com aspas no argumento -run virou prompt(...) em internal/plangen/verify.go; o swarm parou com 'interactive exige prompter'.
- Correcao: fieldsComAspas trata argumento entre aspas duplas como um argv unico; metacaractere de shell e aspa sem fechar seguem para gate manual.
- Efeito colateral achado e corrigido: fabricar plano zerou o progresso (3 PASSED viraram 0); fabricar plano agora preserva estado e prova de ticket nao pendente (internal/cli/fabricar_plano_preserva.go), com teste e mutacao.
- Prova: go test ./internal/plangen/ ok; mutacao falha 1 teste; DAG regerado com gate exit_code e file_exists; camada 3 com gate_ok true.

## A22 - validacao de estilo e item sem prova visual (2026-11-26)
- Achado: o gate humano do TICK-02.1 pedia aprovacao sobre personagens.json (texto). Texto e abstrato para quem aprova.
- Correcao aplicada ao estilo: previa PNG de cada personagem de teste em projetos/video-diorama/molde-de-pilula-preview/ (puppetSVG real, 3 poses e humores por personagem).
- Prova: pixels de pele, terno e contorno presentes nos 3 PNG (analise por PIL).
- Regra proposta: todo estilo e todo item do inventario gera previa visual antes do gate humano (decisao do usuario sobre o comando video-inventario previa).

## A23 - swarm run preso por sessao interrompida (2026-11-26)
- Achado: a sessao caiu durante a camada 6 (TICK-V01-02.2); o worker parou sem editar e o swarm ficou sem c6.done. A nova execucao falhou com 'worker TICK-02-2 ja existe' porque o status era stopped, nao killed.
- Correcao operacional: bin/legoloop.exe swarm kill -name TICK-02-2 e nova execucao. Gate ok.
- Prova: gate_ok true; progresso 6 aprovados; go test ./... 73 pacotes ok.
- Achado de teste: mutacao que aceitava faixa "x" sobreviveu; teste ganhou o caso "x" e a mutacao agora falha 1 teste.

## Exploratorio skins do estilo palito (2026-11-26)
- Nao e ticket do plano: prototipo visual sob pedido do usuario, comparando traco fino (referencia) com o molde-de-pilula.
- Reaproveita PP_POSE, PP_ARM e ppArmGeo do puppet.js original para as mesmas poses (idle, gesture, hip).
- Previa em projetos/video-diorama/palito-preview/skins.png, 3 skins (tinta preta sem acento, tinta azul com acento vermelho, tinta marrom com acento verde) x 3 poses.
- Prova: contagem de pixel por cor confirma as 3 tintas e os 2 acentos presentes na imagem.

## Ajuste do exploratorio: skins preenchidas (2026-11-26)
- Achado: a v1 (so contorno) foi reprovada pelo usuario; pedido e silhueta fina com preenchimento colorido (pele, traje, acento), como nas referencias de mercado.
- Correcao: nova funcao palitoSVG com cabeca, bracos e pernas preenchidos na cor da pele e tronco preenchido na cor do traje, contorno preto fino por cima, mesmas poses do motor.
- Prova: contagem de pixel confirma pele clara, pele morena, 3 tons de traje, 2 acentos e o contorno presentes na imagem.

## Ticket A - vocabulario de 5 faixas etarias (2026-11-26)
- Alteracao: ValidarFaixa aceita bebe, crianca, adolescente, adulto, idoso (internal/video/inventario/genero.go).
- Prova: TestFaixaCincoEtapas cobre as 5 e 4 valores invalidos; mutacao que aceita "x" extra faz 2 testes falharem; go test ./... 73 pacotes ok; vet, funclen, gofmt e espinha limpos.

## Ticket B - exploracao do bebe no motor (2026-11-26)
- Achado: puppetSVG nao tem pose de engatinhar (corpo so fica em pe); rotacionar o corpo todo nao e suportado hoje.
- Alternativa viavel sem mudar motor: escala pequena (0.42) + dress tipo macacao (onesie) + hair bald remove gravata/colarinho automaticamente (if (!dress) so desenha gravata).
- Prova: pele e onesie presentes na imagem; cor da gravata-borboleta ausente (0 pixels) confirmando que o dress supre a gravata.
- Pendente de decisao do usuario: aceitar bebe sentado/em pe (sem engatinhar) como suficiente para o estilo, ou abrir ticket de motor para pose deitada/engatinhando (mudanca maior na geometria do corpo).

## Completar 4 faixas ativas no molde-de-pilula (2026-11-26)
- Acrescentados menina-teste (feminino/crianca), adolescente-teste (masculino/adolescente), adolescente-feminina-teste (feminino/adolescente) e idosa-teste (feminino/idoso) em personagens.json.
- TestGeneroElencoMoldeDePilula ampliado para exigir os 4 (red antes, green depois).
- Prova: go test ./... 73 pacotes ok; video-inventario verificar ok=true; previa visual em projetos/video-diorama/molde-de-pilula-preview/{menina-teste,adolescente-teste,adolescente-feminina-teste,idosa-teste}.png, pixels de pele e traje/vestido confirmados por script.
- Estilo continua rascunho; 10 personagens de teste aguardando aprovacao humana (AC-E07.8).

## Vestuario e cabelo revisados no molde-de-pilula (2026-11-26)
- Regra aplicada: mulheres usam vestido ou camiseta (nunca terno); criancas e adolescentes usam camiseta (nunca terno nem vestido).
- Mudancas em personagens.json: crianca-teste e adolescente-teste (meninos) trocaram suit por dress (camiseta); analista-teste (mulher adulta) trocou suit por dress (camiseta cavada); menina-teste e adolescente-feminina-teste marcadas como camiseta (campo decote) em vez de vestido conceitual.
- 3 tipos de camiseta cobertos: careca (crianca-teste, menina-teste), golaV (adolescente-teste, adolescente-feminina-teste), cavada (analista-teste), via overlay desenhado sobre o torso liso do motor (campo dress), sem mudar puppet.js.
- Cabelo cacheado (ja existente no motor, case 'curly' em ppHair) aplicado a adolescente-teste (menino) e analista-teste (mulher), confirmando 1 de cada genero.
- Prova: go test ./internal/video/inventario/ ok; video-inventario verificar ok=true; folha em projetos/video-diorama/molde-de-pilula-preview/folha-vestuario-cabelo.png; contagem de pixel confirma as 7 cores de roupa e os 5 overlays de decote nas posicoes certas (2 careca, 2 golaV, 1 cavada).

## Correcao: decote vazado (pele) e saia com comprimento (2026-11-26)
- Achado confirmado por pixel: a v1 do decote (overlay) era so uma linha de tinta por cima do tecido, sem pele visivel (amostra de pixel na gola mostrou cor do tecido, nao da pele); vestidinho nao existia (saia nunca tinha comprimento, so o torso).
- Correcao: overlayDecote agora preenche o recorte com a cor da pele do personagem (golaV, cavada, careca); overlaySaia desenha um trapezio na cor do vestido por cima das pernas, com comprimento curto (vestidinho, mulher-teste) ou longo (vestido, idosa-teste).
- Prova: contagem de pixel de pele na regiao do decote: crianca 10446, adolescente 7796, analista 7879 (zero antes); saia longa (idosa, 33953 px) maior que saia curta (mulher, 22396 px), confirmando os dois comprimentos.
- Limite ainda honesto: e um recorte poligonal aproximado sobre a silhueta do motor, nao uma geometria nova no puppet.js original.

## Reprovado pelo usuario: vestido longo e camiseta cavada (2026-11-26)
- Decisao do usuario: remover o vestido longo (idosa-teste) e a camiseta cavada (analista-teste) do conjunto.
- Correcao: idosa-teste voltou ao vestido liso (campo saia removido); analista-teste trocou cavada por golaV (forma ja aprovada em outros 2 personagens).
- Prova: contagem de pixel da idosa (15443) bate com a contagem anterior ao experimento de saia (15414/15416); golaV da analista (7905) na mesma faixa dos golaV ja aceitos (7796-8000).
- go test ./internal/video/inventario/ ok; video-inventario verificar ok=true.
- Folha final em projetos/video-diorama/molde-de-pilula-preview/folha-final.png: 3 tipos de camiseta reduzidos a 2 (careca, golaV), ternos com 2 gravatas, vestido e vestidinho (sem a variante longa).

## Reprovado: vestidinho em curva de sino (2026-11-26)
- Achado: troquei o trapezio reto do vestidinho por uma curva em sino (para melhorar o vestido longo); o usuario reprovou o efeito colateral no vestidinho ("ficou uma bola embaixo").
- Correcao: overlaySaia voltou a forma de trapezio reto original (linhas retas das pontas do quadril ate o barrado, leve curva so no barrado).
- Prova: contagem de pixel do vestidinho (mulher-teste, 22396) bate exatamente com a contagem da primeira versao aprovada.

## Estilo molde-de-pilula validado (2026-11-26)
- Gate humano (AC-E07.8) cumprido: aprovacao.json gravado com os 10 personagens de teste e nota do que foi aprovado (folha-final.png).
- video-inventario validar --estilo molde-de-pilula: status mudou para validado.
- Isso fecha o TICK-V01-02.1 (criacao do estilo) com a validacao humana pendente resolvida; TICK-V01-02.2 (genero/faixa) ja estava aprovado.

## Camada 7 - TICK-V01-03.1 pre-voo le modo e estilo (2026-11-26)
- Entrega: ContratoProducao (modo_producao, estilo), ModoEfetivo/Legado, validarEstiloPreVoo (statusEstilo le inventario/video/estilos/<id>/estilo.json), AvisoEstiloVisualLegado.
- Prova: go test ./internal/diorama/ativos/... ok; mutacao na checagem de status!=validado falha 1 teste; mutacao em ModoVisualValido falha 2 testes; go test ./... 73 pacotes ok; vet, funclen, gofmt, espinha e go.mod limpos.
- Progresso: 7 aprovados.

## Camada 8 - TICK-V01-03.2 storyboard v2 valida elenco e inventario (2026-11-26)
- Entrega: ValidarStoryboardV2, ValidarElenco/validarCampoCast (tipo + partes_suportadas), CarregarInventarioEstilo.
- Achado meu: primeira mutacao na checagem de partes_suportadas quebrou o build (forma nao usada) e meu grep '--- FAIL' contou 0, lendo como sobrevivencia falsa; corrigido mantendo forma usada (if forma && false) e a mutacao real mata 2 testes (TestElencoCastBloqueios ja cobria hair/hat fora de partes_suportadas).
- Prova: go test ./internal/diorama/ativos/... ok; mutacao do personagem_id fora do elenco mata 1 teste; mutacao de partes_suportadas mata 2; go test ./... 73 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 8 aprovados.

## Camada 9 - TICK-V01-04.1 etapa 3 vetorial monta elenco do inventario (2026-11-26)
- Entrega: ExecutarEtapa3Vetorial, lacunasDoCast (item ausente vira lacuna, nao erro), gravarItemNovo (grava com Aprovado=false).
- Prova: go test ./internal/diorama/ativos/... ok; mutacao que ignora item ausente mata 1 teste; mutacao que grava item ja aprovado mata 1 teste; go test ./... 73 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 9 aprovados.

## Camada 10 - TICK-V01-04.2 mantem caminho de recorte sem regressao (2026-11-26)
- Entrega: RenderizadorDeRecorte (recorte ou recorte-<variante>), RecortadorDoEstilo, ExecutarEtapa3Recorte (caminho legado: fichas 4 faixas + rembg), ErrRenderizadorNaoRecorte para estilo vetorial.
- Prova: go test ./internal/diorama/ativos/... ok; mutacao que ignora o renderizador em fichas.go mata 1 teste; mutacao equivalente em drivers.go mata 1 teste; go test ./... 73 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 10 aprovados.

## Camada 11 - TICK-V01-04.3 etapa 4 consulta inventario antes de gerar fundo (2026-11-26)
- Entrega: ObterFundoInventario (busca por estilo+tags antes do Flow), reutilizarFundo (confere hash na origem e na copia), gerarFundoNovo/registrarFundoNovo (prompt_mestre do estilo + prompt da cena, hash, item sem aprovacao).
- Prova: go test ./internal/diorama/ativos/... ok; mutacao que pula a busca no inventario (gera sempre) mata 2 testes.
- Achado sem correcao de codigo: a segunda checagem de hash (da copia recem-escrita no projeto) sobreviveu a mutacao (0 falhas); e defensiva contra falha de IO entre escrever e reler o mesmo arquivo, cenario que um teste de unidade normal nao reproduz sem injetar falha de disco; documentado aqui em vez de simulado.
- go test ./... 73 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos. Progresso: 11 aprovados.

## Camada 12 - TICK-V01-05.1 motor-local com licenca pendente (2026-11-26)
- Entrega: motor-local/LICENCA.json (status pendente, plano B de reescrita de puppet.js/sets.js), motor-local/ ignorado em .gitignore:56, licenca_test.go (TestLicencaDecisaoPorStatus, TestLicencaMotorLocalGovernanca que roda git check-ignore de verdade).
- Prova: go test ./internal/video/motor/... ok; mutacao no case "pendente" mata 2 testes; go test ./... 74 pacotes ok (novo pacote motor); vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 12 aprovados.

## Camada 13 - TICK-V01-05.2 pre-voo do motor (2026-11-26)
- Entrega: PreVoo checa motor-local, python, playwright, chromium (script real via Playwright), ffmpeg, nessa ordem; ErrMotorAusente com instrucao de instalacao por item.
- Achado meu: duas tentativas de mutacao quebraram o build (err nao declarado, import os nao usado) e meu metodo de contagem leu isso como sobrevivencia; corrigido mantendo o build valido, e a mutacao real mata o teste (1 funcao com 5 casos via t.Errorf, confirmado tambem isolando so o caso motor-local).
- Prova: go test ./internal/video/motor/... ok; go test ./... 74 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 13 aprovados.

## Camada 14 - TICK-V01-05.3 ponte cenas+elenco+timings+inventario para timeline.json (2026-11-26)
- Entrega: Montar/Gerar (timeline.json + check_timeline via Checador injetavel), luzDoEstilo (AC-E10.6: light.rim so com estilo.Luz.Rim), validarProporcao (AC-E10.8: cabeca>1, 0<corpo<1), finsEmQuadros (duracao da ultima cena bate com a voz).
- Prova: go test ./internal/video/motor/... ok; mutacao na luz do estilo mata 1 teste; mutacao na proporcao mata 1 (isolada por build valido, as duas primeiras tentativas quebraram o build e foram descartadas); mutacao no CHECK OK mata 1.
- Achado sem correcao (fora do escopo deste ticket): TestAtualizarConcorrenteNaoPerdeEscrita (internal/progresso) falhou uma vez isolado com 'Access is denied' no .lock do Windows; rodado 5x isolado e a suite completa de novo, sempre ok. Teste instavel por concorrencia de arquivo no Windows, nao causado por este ticket.
- go test ./... 74 pacotes ok (2a rodada); vet, funclen, gofmt, go.mod e espinha limpos. Progresso: 14 aprovados.

## Camada 15 - TICK-V01-05.4 render de referencia (2026-11-26)
- Entrega: Renderizar (Executor+Sonda injetaveis), ValidarMidia (1080x1920, fps, audio, duracao da voz +-1 quadro, zero brancos), HashesCache/CompararCaches/ConferirDeterminismo (AC-E10.5), EstadosBoca/ValidarBoca por diferenca de imagem na regiao da boca (AC-E10.7), SondaFFmpeg real via ffprobe/ffmpeg (nao usada nos testes).
- Achado meu: 2 das 4 mutacoes quebraram o build (hash e n nao usados) e foram descartadas; refeitas mantendo o build valido.
- Prova: go test ./internal/video/motor/... ok; mutacao no formato, na duracao, no CompararCaches e no ValidarBoca mata 1 teste cada; go test ./... 74 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos.
- Progresso: 15 aprovados. Epico E10 completo (AC-E10.1 a AC-E10.8).

## Camada 16 - TICK-V01-06.1 render Go ligado ao motor (2026-11-26)
- Entrega: ExecutarRenderMotor/reconfirmarMotor (le timeline.json do motor, reconfirma voz, mix no alvo, SFX e cortes da SPEC v1.4.0), TimelineMotor.conferirCortes (fundo existe, soma das duracoes bate com a voz +-1 quadro), eventosTrocaCena, MotorRender interface injetavel.
- Achado meu: 1 das 5 mutacoes quebrou o build (err nao declarado); refeita mantendo o build valido.
- Prova: go test ./internal/diorama/render/... ok; mutacao na soma de duracoes, no fundo ausente, nos eventos de SFX no mix e na chamada do motor mata 1 teste cada.
- go test ./... 75 pacotes ok; vet, funclen, gofmt, go.mod e espinha limpos. Progresso: 16 aprovados.

## Camada 17 - TICK-V01-04.4 marca entrega suspeita do Flow (2026-11-26)
- Corrigido antes: pulei o TICK-V01-04.4 ao recomendar a camada seguinte (chamei de 06.2, que nao existe no plano); confirmado por progresso.json com PENDING e corrigido executando o ticket real agora.
- Entrega: AvaliarEntregaFlow (3 motivos: sem confirmacao, hash cruzado com item de outro tipo, fora do padrao de cor/dimensao do estilo), ConfirmarSuspeito (exige folha de contato e aprovador humano).
- Prova: go test ./internal/diorama/ativos/... ok; mutacao em cada um dos 3 motivos mata 1 teste; sem erro de build desta vez.
- Achado de governanca, nao corrigido por mim (fora do meu ticket): internal/kraften/{artefatos,cloudflare,preflight}.go e sync_test.go estao no disco, nao commitados, com sync_test.go referenciando NovoCloudflareClient que so existe em cloudflare.go (nao commitado); go build ./internal/kraften/... falha isolado. Esses arquivos ja tinham sido arrastados para o commit 23c7eec por 'git add -A'; desta vez eu os tirei do stage (git restore --staged internal/kraften) antes de commitar, para nao empacotar trabalho alheio inacabado.
- go test -count=1 $(go list ./... | grep -v internal/kraften): 74 pacotes ok, 0 falhas. vet, funclen e gofmt limpos em internal/diorama. Progresso: 17 aprovados.
