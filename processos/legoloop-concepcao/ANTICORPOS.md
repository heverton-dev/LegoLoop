---
schema: anticorpos.v1
processo: legoloop-concepcao
total_anticorpos: 18
globais: 0
dominio: 18
regras_inviolaveis:
  - pacote: "internal/interview"
    antidoto: "internal/interview/espelho_campos.go lê o realce antes ou depois dos dois-pontos e as linhas seguintes; testes espelh..."
  - pacote: "internal/mcp"
    antidoto: "internal/mcp/busca_tokens.go: pontua palavras, exige metade, ordena; testes busca_tokens_test.go com mutação provada."
  - pacote: "internal/gate"
    antidoto: "Regex ampliada em internal/gate/builtins.go; teste aceite_titulo_test.go lê o template e a SPEC da Forja."
  - pacote: "internal/processo"
    antidoto: "layoutsLegados em internal/processo/vinculo.go; TestInventarioPorHarness protege."
  - pacote: "internal/cli"
    antidoto: "internal/cli/cli.go corrigido; usage_gateway_test.go protege."
  - pacote: "internal/cli"
    antidoto: "parseIntercalado em internal/cli/anticorpos_registrar.go; o subcomando registrar passou a existir."
  - pacote: "operação"
    antidoto: "Reiniciar o roteador e os harnesses depois de go install (AGENTS.md, regra 5); conferir a busca pelo comando e não pe..."
  - pacote: "docs/specs"
    antidoto: "Escrever o padrão como classe de caracteres: grep -c com colchetes em volta de cada chave."
  - pacote: "operação"
    antidoto: "Atalhos codebase-memory-mcp e .cmd em ~/.local/bin; a Forja da E05 repete a criação do atalho."
  - pacote: "operação"
    antidoto: "Gravar com newline LF ou rodar gofmt -w no arquivo; conferir gofmt -l antes de fechar o ticket."
  - pacote: "operação"
    antidoto: "Gravar arquivos Go com a ferramenta de escrita, ou montar a barra com chr(92) em Python."
  - pacote: "internal/plangen"
    antidoto: "materialize grava o PLAN antes de montar a ROUTE; testes assinatura_plano_test.go provam a assinatura real e a detecç..."
  - pacote: "internal/cli"
    antidoto: "legoloop plano --spec S --tickets T.json [--estrito]; testes plano_test.go; a lista de Tickets fica versionada em doc..."
  - pacote: "internal/cli"
    antidoto: "Mesma correção do materialize; plano_test.go confere que duas execuções geram os mesmos bytes"
  - pacote: "internal/refinar"
    antidoto: "Expectativa ROUTE {3, 3} com comentário no teste; a mudança é consequência correta e não regressão"
  - pacote: "internal/upstream"
    antidoto: "Pacote internal/upstream (VerificarSpec, IntencaoAprovada, Recusa) ligado a plano, specgen, plangen.Generate, route2d..."
  - pacote: "internal/studio"
    antidoto: "upstream.Recusa(err) mapeia a recusa da trava para 422; testes Devolve422 em spec_test.go e plan_test.go"
  - pacote: "internal/plangen"
    antidoto: "Helper gravarIntencaoAprovada e constante specMinima; toda fixture que gera PLAN, SPEC ou DAG traz a INTENT aprovada"
---

# 🛡️ Banco de Anticorpos — legoloop-concepcao

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (18)

### `internal/interview`
- **Antídoto:** internal/interview/espelho_campos.go lê o realce antes ou depois dos dois-pontos e as linhas seguintes; testes espelho_formato_test.go com Red provado.
- **Erro Evitado:** Espelho Canônico truncado em ** e confiança presa em 10% com o claude
- **Causa Raiz:** Os regex aceitavam **Campo**: mas não **Campo:** texto, e liam só a primeira linha de cada campo.
- **Assinatura:** `sha256:6ba5c27c050b3b8cb32a5501efdc0c174183ff2a6e745688717c00be084a0f48`

### `internal/mcp`
- **Antídoto:** internal/mcp/busca_tokens.go: pontua palavras, exige metade, ordena; testes busca_tokens_test.go com mutação provada.
- **Erro Evitado:** mcp_search não acha ferramenta com várias palavras e devolve ordem aleatória
- **Causa Raiz:** matchQuery exigia a frase inteira colada no texto, e searchExternal percorria um map.
- **Assinatura:** `sha256:75d8205c203343f86b8c62c3eb045dba74e8302c4c4d4dca542dbd7e75351321`

### `internal/gate`
- **Antídoto:** Regex ampliada em internal/gate/builtins.go; teste aceite_titulo_test.go lê o template e a SPEC da Forja.
- **Erro Evitado:** has_acceptance_criteria reprova o título canônico Critérios Binários de Aceite
- **Causa Raiz:** acceptanceRe exigia 'critérios de aceita' e não aceitava a palavra Binários nem 'Aceite'.
- **Assinatura:** `sha256:2ce32da32c24a211b6fabf6b53659c1a7532f95cc9de0a55f90995018b6cc764`

### `internal/processo`
- **Antídoto:** layoutsLegados em internal/processo/vinculo.go; TestInventarioPorHarness protege.
- **Erro Evitado:** opencode deixa de ser detectado quando só existe ~/.opencode
- **Causa Raiz:** O layout novo (.config/opencode) substituiu o antigo na detecção.
- **Assinatura:** `sha256:3686f2efd420a8d8e99f916eaec99eb9db7dc52487904c00a2a3d16ca9c43b45`

### `internal/cli`
- **Antídoto:** internal/cli/cli.go corrigido; usage_gateway_test.go protege.
- **Erro Evitado:** Texto de uso do mcp-router --gateway dizia 4 ferramentas
- **Causa Raiz:** O texto não acompanhou o gateway, que expõe só call_mcp_tool.
- **Assinatura:** `sha256:afa010fdf38ffa196dae4fff4ef7e418df0f174f7bb648c5562f7a9e25634d96`

### `internal/cli`
- **Antídoto:** parseIntercalado em internal/cli/anticorpos_registrar.go; o subcomando registrar passou a existir.
- **Erro Evitado:** anticorpos sync <slug> --raiz DIR ignorava --raiz
- **Causa Raiz:** O pacote flag para de ler opções na primeira palavra solta.
- **Assinatura:** `sha256:b632096d0466ce69ff06a3db42ec45ebf1272dda2aedeed11faa628f2d601489`

### `operação`
- **Antídoto:** Reiniciar o roteador e os harnesses depois de go install (AGENTS.md, regra 5); conferir a busca pelo comando e não pelo processo aberto.
- **Erro Evitado:** Roteador MCP ou harness aberto continua com o binário antigo depois do go install
- **Causa Raiz:** O processo em execução não recarrega o binário; o skill_search mostrou o ranking velho.
- **Assinatura:** `sha256:3a31a0623781a421b0ae41c098e09ef3e6ffdf2a1a0c01f57fefbb52d3507d56`

### `docs/specs`
- **Antídoto:** Escrever o padrão como classe de caracteres: grep -c com colchetes em volta de cada chave.
- **Erro Evitado:** SPEC que cita o padrão de lacuna no próprio texto reprova o Gate de lacunas
- **Causa Raiz:** O Gate reprova qualquer ocorrência literal do padrão, inclusive dentro de um comando de exemplo.
- **Assinatura:** `sha256:d10e5cf09fc0d68d61cba3b7da4adc6f2404ffef076c9416ee7362deeb4a9867`

### `operação`
- **Antídoto:** Atalhos codebase-memory-mcp e .cmd em ~/.local/bin; a Forja da E05 repete a criação do atalho.
- **Erro Evitado:** codebase-memory-mcp não resolvia pelo nome no bash
- **Causa Raiz:** O binário fica em ~/tools/codebase-memory-mcp, fora do PATH.
- **Assinatura:** `sha256:8e11d6fec8ae2351948bc94cfea1a68d8d3bcc6a7de126b1c392d0a8adfc5b63`

### `operação`
- **Antídoto:** Gravar com newline LF ou rodar gofmt -w no arquivo; conferir gofmt -l antes de fechar o ticket.
- **Erro Evitado:** Python no Windows grava CRLF e o gofmt -l lista o arquivo
- **Causa Raiz:** open() em modo texto converte quebra de linha no Windows.
- **Assinatura:** `sha256:3ec523f0b6eda237aa0d7a0b4a79b1b4600f6f786f0db8a78f9df9704eb0b33c`

### `operação`
- **Antídoto:** Gravar arquivos Go com a ferramenta de escrita, ou montar a barra com chr(92) em Python.
- **Erro Evitado:** Barras invertidas e quebras de linha de literais Go corrompidas no transporte de heredoc
- **Causa Raiz:** O transporte consome um nível de barra invertida em strings de comandos.
- **Assinatura:** `sha256:b389171a3ac019cf765d216605533b7fce34d7113df69a2a938f1989e8a258b9`

### `internal/plangen`
- **Antídoto:** materialize grava o PLAN antes de montar a ROUTE; testes assinatura_plano_test.go provam a assinatura real e a detecção de adulteração
- **Erro Evitado:** ROUTE nasce com a assinatura do PLAN zerada na primeira geração
- **Causa Raiz:** materialize montava a ROUTE antes de gravar o PLAN, e a assinatura zerada é ignorada na verificação, o que deixava um PLAN adulterado passar
- **Assinatura:** `sha256:940e2ca880627e51effecd16a5fabbaff5ce2065142f2c9e1fde129dd4ddafda`

### `internal/cli`
- **Antídoto:** legoloop plano --spec S --tickets T.json [--estrito]; testes plano_test.go; a lista de Tickets fica versionada em docs/plans
- **Erro Evitado:** Não havia comando para renderizar PLAN e ROUTE a partir de uma lista de Tickets
- **Causa Raiz:** Render existia só como API Go e o caminho da CLI exigia LLM, então a geração dependia de um programa descartável
- **Assinatura:** `sha256:9ff0a90cbd8a422e7b2cc0f48019fa6e0c8f3e7ccccc999f5a27389ca656081f`

### `internal/cli`
- **Antídoto:** Mesma correção do materialize; plano_test.go confere que duas execuções geram os mesmos bytes
- **Erro Evitado:** Teste do legoloop plano falhava por geração não determinística
- **Causa Raiz:** A ROUTE era gerada com o PLAN ausente do disco, então a segunda execução produzia bytes diferentes
- **Assinatura:** `sha256:37386c76fe5134b86a5e9b58813efd1aa5cc3cf35a6e6fcf232578dbd9dfb5d5`

### `internal/refinar`
- **Antídoto:** Expectativa ROUTE {3, 3} com comentário no teste; a mudança é consequência correta e não regressão
- **Erro Evitado:** TestRefinar esperava ROUTE +2 -2 e passou a dar +3 -3 depois da correção da assinatura
- **Causa Raiz:** O teste codificava o hash do PLAN zerado, que nunca mudava; com a assinatura real, reescrever o PLAN muda o hash na ROUTE
- **Assinatura:** `sha256:3ae8f038b07d91ab8a49ba6e1f22f61e4f489eeacda70107956f60a0146fb305`

### `internal/upstream`
- **Antídoto:** Pacote internal/upstream (VerificarSpec, IntencaoAprovada, Recusa) ligado a plano, specgen, plangen.Generate, route2dag e Studio; exit 2 na CLI e HTTP 422 no Studio; testes com Red provado
- **Erro Evitado:** SPEC, PLAN e DAG eram gerados sem INTENT aprovada, então a cadeia da concepção era só convenção
- **Causa Raiz:** plangen.Render, specgen e route2dag não conferiam o upstream; legoloop plano gerou ROUTE sem INTENT e saiu com 0
- **Assinatura:** `sha256:2ec34291e87fa1b8a1e2eff3a004e37b8803dc894c7a5d23bd1b199ff2e1db62`

### `internal/studio`
- **Antídoto:** upstream.Recusa(err) mapeia a recusa da trava para 422; testes Devolve422 em spec_test.go e plan_test.go
- **Erro Evitado:** Studio devolvia 500 para recusa de regra do processo
- **Causa Raiz:** Os handlers tratavam todo erro do gerador como falha do servidor
- **Assinatura:** `sha256:9be662b9afc67a2dc4bb61d298502d42a90863872ba4e88b6406e39b6b98e077`

### `internal/plangen`
- **Antídoto:** Helper gravarIntencaoAprovada e constante specMinima; toda fixture que gera PLAN, SPEC ou DAG traz a INTENT aprovada
- **Erro Evitado:** Fixtures de teste geravam PLAN a partir de SPEC sem INTENT
- **Causa Raiz:** Os testes antigos escreviam SPEC mínima, sem Referência da Intenção nem INTENT aprovada
- **Assinatura:** `sha256:743a5df489742305a0dd74a292f3a1a37284243456c917e8fd41f3241dc5a4ce`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

