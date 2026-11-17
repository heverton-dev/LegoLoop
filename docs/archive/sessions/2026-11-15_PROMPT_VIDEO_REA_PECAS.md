Você trabalha no repositório C:\Users\trcnologia\Desktop\LegoLoop (Go, binário legoloop).
Leia AGENTS.md antes de qualquer ação. Responda em PT-BR, uma frase por linha. Sem subagentes. Sem commit e sem push sem autorização do operador.

Estado herdado (provado em 15/11/2026)
- FABRICAR agora cria peça única: legoloop fabricar "<nome>" --peca script|skill|gate|hook|mcp|processo (internal/fabricar/peca_unica.go). Gera INTENT, SPEC, PLAN, ROUTE, DAG, DAG ISO, memória e swarm.plan.json. Não gera TRIAGEM.
- REFINAR adiciona etapa: legoloop refinar <nome> --adicionar --etapa e0N --peca P --caso C [--depende ID] --yes.
- COMPOR liga processo: legoloop compor <alvo> <sub@versão> [--depende ID].
- Regra: triagem só com mais de um nó. validate, fechar, refinar e compor exigem triagem.json cobrindo todos os nós (triagem.Exigida em internal/triagem/exigencia.go).
- DAG e Studio ficam na pasta processos/<nome>/ (caminhos.DestinosDoProcesso).
- Fechamento da capacidade: docs/reports/2026-11-15_REPORT_FABRICAR_PECA_UNICA.md. Sem commit feito.
- Use sempre --model code-pro. O modelo code-fast devolve INTENT vazio.
- Entrevista socrática: cada pergunta traz fato provado (arquivo:linha ou comando) e a sua recomendação. Não responda pela pessoa.

Objetivo (autorizado pelo operador)
Fabricar duas peças para o processo legoloop-video-rea com a capacidade nova:
1. transcrever (script): converte o áudio aac com ffmpeg e transcreve com C:\Users\trcnologia\tools\whisper-cpp\Release\whisper-cli.exe e o modelo C:\Users\trcnologia\.cache\whisper-models\ggml-small.bin. Plano B: read_file multimodal. Verifique se os dois caminhos existem antes de usar.
2. relatório (skill): grava docs/reports/AAAA-MM-DD_REPORT_<slug>.md com frontmatter YAML e as classes DEVE, PODE e NÃO PODE (cada item com o porquê e o como).
Depois refine os nós:
- n02 forensics: acrescentar conceitos, processos e repositórios (internal/videorea/types.go:46-56).
- n03 curate: classificar todos os achados em DEVE, PODE e NÃO PODE.
Atualize junto internal/golden/testdata/processos/legoloop-video-rea.dag.json.
Já existe a entrevista intv_1794760080406643600 (100%) e o INTENT docs/intents/2026-11-15_INTENT_PROCESSO_QUE_RECEBE_A_URL_DE_U.md para a análise de vídeo.

Como conduzir
- Um ticket por contexto. TDD Red-Green em Go. Recompile com go build -o bin/legoloop.exe ./cmd/legoloop após alterar Go.
- Antes de mexer no video-rea, leia o estado atual: git status, o DAG real e o golden. Não reverifique o que está acima sem motivo.
- Portão por ticket: go test -count=1 ./..., go vet ./..., go run ./tools/funclen internal, gofmt -l . (saída vazia), git diff --exit-code -- go.mod go.sum, go run ./tools/taxlinter internal cmd loops skills docs/templates.
- Teste novo só vale se uma mutação do código o faz falhar. Restaure após cada mutação.
- Feche com legoloop fechar <slug>. Nunca declare concluído sem as 9 etapas. Declare o que não foi provado.
- Não use taskkill /IM legoloop.exe. Stop-gates SG-* só com aprovação explícita.
- Ao fim, entregue a próxima ação sugerida e o prompt de ativação.
