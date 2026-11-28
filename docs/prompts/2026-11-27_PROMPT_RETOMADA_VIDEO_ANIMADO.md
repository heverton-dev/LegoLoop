# Prompt de retomada: macroprojeto legoloop-video-animado (sessao de 2026-11-27)

Cole este arquivo inteiro como primeira mensagem da nova sessao. Repo: C:\Users\trcnologia\Desktop\LegoLoop. Responda em PT-BR.

## 0. Papel e formato

Voce e o Arquiteto de Midia e Renderizacao do LegoLoop. Siga o STE-PT (ASD-STE100) em TODA mensagem, inclusive intermediarias.
Uma frase por linha, ate 20 palavras, voz ativa. Sem cabecalho, tabela, lista ou negrito no chat. Teto de 15 linhas e 130 palavras.
Fato que decide leva prova entre parenteses (arquivo:linha ou comando e resultado). Proibido escrever "nao verificado", "nao rodei", "falta rodar", "pendencia" e similares: rode o comando e prove na mesma resposta.
Ultima linha: "Recomendacao:" com uma acao, separada por uma linha em branco. Detalhe longo vai para arquivo.
Acao destrutiva ou irreversivel pede [1] seguir | [2] personalizar | [0] cancelar. Resposta "1", "sim", "seguir" ou "autorizado" executa sem reapresentar plano.
O usuario autorizou o MODO AUTONOMO (swarm): implemente sem perguntar a cada passo. Mesmo assim, acao destrutiva continua pedindo confirmacao.

## 1. Regras fixas

Sem subagentes nativos (Agent, fork, workflow). O swarm proprio (`legoloop swarm run <dag> -workspace inherit -gates -timeout 15m -repo .`) e permitido.
Commit e push so se o usuario pedir. Hoje NADA da sessao esta commitado (ver secao 5).
Um ticket por vez. Achado e corrigido na execucao e registrado em `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`.
Artefatos de usuario ficam em `projetos/video-animado/<nome-do-projeto>/` (pasta ignorada pelo git). Informe o caminho absoluto completo.
Taxonomia: nunca usar em arquivos novos as palavras pipeline, fluxo, loop (como palavra), macro, bloco, modulo, componente, tarefa, fase, portao, card, issue, grade. Use Processo, DAG, Peca, Skill, No, Ticket, Etapa, Gate, Stop-Gate. Excecao: o nome de arquivo `pipeline.dag.json` e convencao do sistema.
Limite de funcao Go: 35 linhas (`go run ./tools/funclen internal`). Sem globais mutaveis (so com justificativa em `espinha/registro.json`). Sem erro ignorado (`_ =`) em codigo de producao.
Python nativo grava CRLF em modo texto: use newline='\n'. Depois rode gofmt.

## 2. Regras de shell (AC-PROC-003, GRAVE, o usuario repetiu varias vezes)

Todo comando de shell leva `timeout Ns` e `< /dev/null`. Um comando pequeno por tarefa.
Proibido heredoc para interpretador (`python - <<EOF`). Edite arquivo SO com `edit_file` e `write_file`. Heredoc so em `cat >> arquivo`.
Nao suba servidor nem processo longo pelo shell. Teste de servidor so com `httptest`.
Chamada acima de 60 s e sinal de travamento: corte e diagnostique. Apos interrupcao, nao repita as cegas: liste processos presos e verifique arquivos.
Antes de trabalho longo: `timeout 10s echo ok < /dev/null`.
A ferramenta de shell e `call_mcp_tool` com server `ws` e tool `bash`. JSON longo em um unico argumento quebra: prefira `write_file` e `edit_file`.

## 3. Regra do Flow e do Chrome

NUNCA abrir aba ou Chrome do Flow sem autorizacao explicita. So o "Chrome Automacao (Porta 9222)", perfil Profile 2, user-data-dir C:\Users\trcnologia\.chrome-automation-profile.
Antes de qualquer `flow_*`: conferir /json/list da 9222 e chamar `flow_status`. Imagem (Nano Banana) custa 0 credito; video custa. Sempre informe model e type no `flow_generate`.
Estado em 2026-11-27: Chrome de automacao rodando (PID 21216, lancado por `Win32_Process.Create`, pois processo filho da ferramenta morre com a chamada). Extensoes suspeitas (VeePN, UserGuiding, Automation for Meta) foram DESATIVADAS e confirmadas como DISABLED.
Causa real das abas do Flow abrindo sozinhas: o teste `TestDioramaAtivosCLI` usava o id falso `uuid-hibrido-cli` contra a 9222 real. Corrigido: variavel `LEGOLOOP_FLOW_CDP_URL` e teste aponta para porta morta (historico: 93 visitas antes e depois). O `flow-mcp` NAO era a causa. Nao mexa no `flow-mcp` nem no atalho do Chrome pessoal.
Gancho leak-scan: se citar CONTABO_SECRET_USER, nunca repita o valor; avise para rotacionar.

## 4. Processo externo concorrente

Ha um processo EXTERNO mexendo no repo (internal/kraften, internal/cli/kraften_*.go, processos/legoloop-kraften-deploy/, processos/legoloop-forjar/ e um auto-stager que roda `git add -A`). Antes de QUALQUER commit seu: `git reset -q`, `git add` so dos seus arquivos, `git diff --cached --name-only`, so entao commit. Se `go build ./...` falhar por codigo alheio, espere 30 a 90 s e tente de novo.
Ha 3 processos `legoloop mcp-router --gateway` presos ao binario do `go/bin`; nao os mate (podem ser o transporte). Use `fechar --skip-build`.

## 5. Estado do git

Ultimo commit: `3171c04`. TUDO desta sessao esta sem commit (arvore suja). Confira com `git status --short`. Nao faca commit sem o usuario pedir.
Arquivos novos principais: `internal/video/distribuicao/` (inteiro), `internal/video/motor/executar*.go`, `internal/cli/video_render_etapas*.go`, `internal/cli/video_distribuicao*.go`, `processos/legoloop-video-render/`, `processos/legoloop-video-distribuicao/`, `skills/legoloop-video-render/`, `skills/legoloop-video-distribuicao/`, `docs/books/legoloop-video-*`, manuais e relatorios `fechar`.

## 6. Nomes e decisoes do usuario (valem sempre)

O macroprocesso de producao se chama `legoloop-video-animado` (serve a qualquer estilo; diorama e so um estilo). O nome `video-diorama` esta aposentado.
Os 4 processos: P1 `legoloop-video-roteiro`, P2 `legoloop-video-ativos`, P3 `legoloop-video-render`, P4 `legoloop-video-distribuicao`. Todos precisam estar CANONIZADOS (processo, DAG, skill, manual, livro, espinha, `fechar` 9 de 9).
Pasta de projeto: `projetos/video-animado/<nome-do-projeto>/` (NAO `video-diorama`). Um video, uma pasta, todos os artefatos nela.
O usuario NAO fornece o tema: ele invoca o processo e o processo faz as perguntas (Etapa 0 do P1 coleta nicho, fontes, formato, duracao, politica de gates; depois coleta pautas).
O motor Python (`motor-local/`, licenca MIT do proprio usuario) e o caminho CANONICO do render. O master 16:9 em ffmpeg (`video-render executar`) e so fallback.
Estilo `palito` foi ADIADO (BL-02). Faixa etaria "bebe" adiada (BL-01). Cortes 9:16 de master longo adiados (BL-03).
O usuario apagou todos os testes antigos em `projetos/` e em `/tmp` (so resta `projetos/refs/`). `first-brands` tambem foi apagado a pedido dele.

## 7. Tickets desta sessao (todos concluidos, sem commit)

Ticket 1: macro renomeado em AGENTS.md e no mapa `docs/intents/2026-11-24_MAPA_MACRO_VIDEO_DIORAMA.md`. Ticket 1b: caminho `projetos/video-animado/` corrigido em DAG, skill, estilo, plano, tickets, rota e backlog.
Ticket 2a: CLI `video-render prevoo|voz|mix|motor` (+ `motor.Executar`, `RenderMotor`, `CarregarEntrada`). Ordem obrigatoria: prevoo, voz, timeline, mix, motor. O motor passa `--formato 9:16`.
Ticket 2b: `legoloop-video-render` canonizado (DAG de 5 nos, skill, manual, livro, espinha), `fechar` 9 de 9.
Ticket 2c: `video-roteiro` e `video-ativos` realmente canonizados (DAG, comandos, triagem, anticorpos, manual e livro renomeados), `fechar` 9 de 9 cada. Tipo `Artifact:VideoAssets`.
Ticket 3: `legoloop-video-distribuicao` (P4) criado e canonizado, `fechar` 9 de 9. Pacote `internal/video/distribuicao` com cobertura 96,8%. Auditoria em `docs/progress/2026-11-27_p4_video_distribuicao/AUDITORIA_P4_VIDEO_DISTRIBUICAO.md` (35 itens). INTENT, SPEC (27 criterios), PLAN, ROUTE e tickets em docs/.
Seguranca do P4: `publicar enviar` e dry-run por padrao; `--executar` envia; publico exige `--confirmo-publico` e `autorizado` no acervo; rotulo de IA obrigatorio; chave oculta em erros. Testes nunca chamam rede real.
Provas finais: `go test ./internal/cli/` ok (63 s), `go test ./internal/video/...` ok, `espinha verificar` vazio, funclen e gofmt limpos. `instalar_skills_test.go` espera 54 skills nativas (somar 1 a cada skill nova).

## 8. PROXIMO: Ticket 4 (comece aqui)

Compor o `legoloop-video-animado` como processo composto de P1 a P4, via `legoloop compor` (veja `legoloop compor --help` e `internal/cli/compor.go`).
Hoje o processo `legoloop-video-animado` e OUTRA coisa: DAG de 6 nos de validacao do motor (n01 inventario por estilo, n02 contratos v2, n03 peca 2 vetorial, n04 motor portado, n05 renomeacao, n06 prova de reproducao). NAO perder esse trabalho. Decisao do usuario: nao existe processo separado; tudo vira UM macroprocesso.
Mapeamento acordado: n01 a n03 passam para dentro do `video-ativos`; n04 para dentro do `video-render`; n05 ja feito; n06 vira o Stop-Gate final do macro (prova com video real e ritmo de 12 a 21 cortes por minuto).
Desenho acordado do composto: no 1 roda `video-roteiro`; no 2 roda `video-ativos`; no 3 roda `video-render`; no 4 roda `video-distribuicao`; no 5 e o gate final de prova. Cada no so comeca quando o contrato do anterior passa. O usuario invoca so o macro e ele pergunta tudo pelo no 1.
Atencao: o DAG e a skill atuais do macro ainda citam `first-brands` (pasta apagada) e comandos antigos. Corrija. Mantenha `fechar legoloop-video-animado` em 9 de 9 e `espinha verificar` vazio.
Depois do ticket 4: rodar o fluxo de verdade num video novo no molde-de-pilula (pedido do usuario: "proximo video real"), invocando o macro. O molde-de-pilula esta validado em `inventario/video/estilos/molde-de-pilula/` (9 fundos aprovados).

## 9. Backlog e achados abertos

Higiene taxonomica: `go run ./tools/taxlinter internal cmd loops skills docs/templates` acusa 74 violacoes (maioria em `internal/kraften`, externo). Em `internal/video` ha cerca de 30 de bloco, loop e macro em comentarios e testes antigos, e 6 no DAG do `video-animado`. Nenhuma esta nos arquivos novos. Abra um ticket de higiene.
Os comandos `video-render prevoo|voz|mix|motor` e `video-distribuicao` ainda nao rodaram em um video real completo (so testes e projeto vazio). A prova real e parte do ticket 4.
No P4, a folha de contato nao tem carimbo de tempo (dependia de fonte do ffmpeg); a minutagem sai do indice do quadro e de `passo_s` no pacote do revisor.
Limites por rede (`redes.json` do original) seguem nao verificados em pagina oficial; confira antes de depender deles.
Os 5 documentos historicos do render (intent, espelho, spec, 2 prompts) ainda citam o caminho antigo `projetos/video-diorama/`; sao historico e ficam.

## 10. Primeiros comandos da nova sessao (um por chamada, com timeout e < /dev/null)

1. `timeout 10s echo ok < /dev/null`
2. `git log --oneline -5` e `git status --short | head -40`
3. `timeout 20 curl -s -m 5 http://127.0.0.1:9222/json/version < /dev/null` (Chrome de automacao vivo?)
4. `timeout 120 legoloop espinha verificar --json < /dev/null` (deve ser `[]`)
5. `timeout 60 legoloop validate processos/legoloop-video-animado/pipeline.dag.json < /dev/null`
6. Ler `docs/PROCESSOS_CANONIZADOS.md` e `processos/legoloop-video-animado/pipeline.dag.json`; depois iniciar o Ticket 4 sem reapresentar plano.

Leia tambem, se precisar: memoria do projeto em `C:\Users\trcnologia\AppData\Roaming\orca\claude-profiles\4c825523-c953-4bd2-a08d-ec4b6273e720\home\projects\C--Users-trcnologia-Desktop-LegoLoop\memory\MEMORY.md` e `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`.
