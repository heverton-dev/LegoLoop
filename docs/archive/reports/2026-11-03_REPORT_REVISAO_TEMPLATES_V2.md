# 📊 REPORT: Revisão Profunda dos Templates Canônicos (V1 ➔ V2)

> **Data:** 2026-11-03  
> **Escopo:** `docs/templates/` (11 arquivos), V1 ➔ V2 (a V2 substituiu a V1 no mesmo caminho)  
> **Método:** cada template foi testado contra o componente do repositório que deveria lê-lo, não apenas lido.  

---

## 1. Veredito

Os templates V1 são bem-intencionados e a estrutura documental (INTENT ➔ SPEC ➔ PLAN ➔ ROUTE) é boa. O defeito central é que **as peças JSON e o ROUTE descrevem formatos que nenhum componente do repositório lê**, e o README promete "Padrão Ouro" sem prova. Nenhum template V1 chega a ouro pelo scorer do próprio projeto.

## 2. Evidências medidas (V1)

| Teste | Resultado |
|:---|:---|
| Score `internal/reputation` | Melhor: SPEC `0.77` (prata). `HOOK` e `MCP` `0.18`, `GATE` `0.30`: **tóxicas**. `ROUTE`/`LOOP` bronze. Nenhum ouro. |
| `dag.Parse` (estrito) nos 5 JSON | Todos rejeitados: `unknown field "tier"`, `"protocol_version"`, `"date"`. |
| `gate.Compile` em 7 expressões dos templates | 5 falham: `go test -v ...` dentro da expressão, `output.artifact`, `has_no_stubs(...)` (função inexistente) e `human_approved == true` (identificador inexistente). |
| Loader de hooks (`assets/hooks.go`) | Lê `hooks.<Evento>[{matcher,hooks:[{type,command}]}]`; o template usa `execution`, `security`, `conditions`: **0 hooks reconhecidos**. |
| Loader de MCP (`assets/mcps.go`) | Lê `mcpServers.<nome>`; o template usa `transport`, `tools`, `capabilities`: **0 servidores reconhecidos**. |
| Script `TEMPLATE_PIECE_SCRIPT.ps1` | Sem BOM (PS 5.1 lê como ANSI). `exit 1` e `exit 2` inalcançáveis: `Write-Error` com `$ErrorActionPreference='Stop'` lança exceção; falha de runtime devolve **1**, não 2. Execução "real" devolve **0 sem fazer nada** (falso sucesso). Convenção oposta à da CLI `legoloop` (1 falha, 2 uso). |
| PLAN ↔ ROUTE ("espelho 1:1") | Não eram espelho: `TICK-01.2` com modelo, workspace, MCPs e subagentes diferentes; timeout padrão 180s no PLAN e 300s no ROUTE; ROUTE sem `spec_ref`; PLAN sem timeout, teto de tokens e rollback. |
| ROUTE x AGENTS.md | Exemplo com `subagents_allowed: true`, contra a regra 4. Subloop apontava para `docs/routes/2026-11-01_ROUTE_REA_DISSECTOR.json`, que **não existe**. |
| Gerador `plangen` x ROUTE | O gerador emite um subconjunto do template e `output.gate` no formato antigo, que não compila. O "SSOT" não é imposto por nada. |
| Ferramenta de lint | PLAN e SPEC exigem `golangci-lint`, que **não está instalado** nem configurado no repositório. |

## 3. Achados por template e mudança na V2

| Template | Principais defeitos na V1 | V2 |
|:---|:---|:---|
| `README.md` | Atribui autoridades não verificáveis ("OpenAI / Interceptadores Seguros", "OpenWorkflows"); "100% universais"; promete ouro. | Coluna "quem lê no código"; teto de score medido ($0.82$ anônimo sem histórico); seções 5–8 (placeholders, dicionário, verificação, limitações). |
| `INTENT` | Tabela de blocos acoplada a Go; vocabulário de tipos diferente do DAG/ROUTE; sem premissas, riscos, confiança; critérios sem forma de prova. | Confiança ≥ 95% no cabeçalho; premissas e riscos; critérios com "como se prova"; mapeamento de vocabulário. |
| `SPEC` | Link absoluto `file:///C:/Users/...`; lint inexistente; só caminho feliz; sem AC de falha ou não funcional; stub no exemplo. | Links relativos; `go vet` + `gofmt`; AC de falha obrigatório e `AC-NF`; contrato Go com validação real; dublês só nas bordas. |
| `PLAN` | Não espelhava o ROUTE; sem limites por ticket; gate misturando comando e expressão. | Espelho 1:1 verificado; `Verificação (argv)` separada de `Gate (DSL)`; exemplo com subloop e os stop-gates `SG-E01` e `SG-FINAL`. |
| `ROUTE` | Gates não compiláveis; subagentes `true`; subloop inexistente; timeouts incoerentes; sem `spec_ref`. | `output.verify` (argv) + `output.gate` (DSL compilável); `prompt(...)` para gate humano; subagentes sempre `false`; `spec_ref`; `template_version`. |
| `PIECE_SKILL` | `commands:` não lido por ninguém; sem procedimento (só guardrails); campos de ontologia ausentes; comandos presos a Go. | Frontmatter que o Studio e o grafo leem; seção de procedimento; Target Files; exemplo de I/O; regra contra instrução embutida em dados. |
| `PIECE_SCRIPT` | Ver seção 2. | Exit codes alcançáveis e alinhados à CLI; confinamento à raiz; escrita atômica; `-DryRun` real; transformação real e idempotente; BOM; exceção própria de uso. 7 cenários testados. |
| `PIECE_HOOK` | Formato inventado; evento `pre_tool_call` e `timeout_ms` inexistentes. | Trecho nativo de `settings.json` (`PreToolUse`, `matcher`, `timeout` em segundos), reconhecido pelo inventário, sem segredo no comando. |
| `PIECE_MCP` | Mistura configuração e manifesto do servidor; fixa `protocol_version` antigo. | Trecho nativo de `.mcp.json`; segredo por `${VAR}`; contrato de ferramenta movido para o README. |
| `PIECE_GATE` | `has_no_stubs` inexistente; `output.artifact` inválido; `retry_or_block` ambíguo. | Cinco receitas, todas compiláveis, cobrindo o DSL (mecânico, artefato, contrato, schema, humano). |
| `PIECE_LOOP` | Modelo próprio (`trigger`/`piece`/`edges`) rejeitado por `dag.Parse`. | DAG v2 válido: passa em `legoloop validate`, `dag.Parse` e `loop.Store`. |

## 4. Defeitos que eu mesmo introduzi e corrigi durante a verificação

Registrados porque a verificação existe para isso:
1. Exemplo Go de `Validate()` na SPEC com corpo-stub (violava a Lei #1) ➔ corpo real.
2. `Move-Item -Force` sobrescrevia arquivo somente leitura e devolvia 0 ➔ checagem explícita + `File.Replace`.
3. `File.Replace(..., $null)` falhava porque o PowerShell converte `$null` em `""` ➔ `[NullString]::Value`.
4. `ArgumentException` do .NET era classificada como erro de uso (exit 2) ➔ exceção própria `UsageError`.

## 5. O que a V2 NÃO resolve (decisões suas)

1. ~~`plangen` continua emitindo o ROUTE antigo.~~ **Resolvido na seção 7.**
2. **Scorer textual.** Hook e MCP nativos pontuam como tóxicos por não terem texto; o score não mede a qualidade real do JSON.
3. **Hooks de outros harnesses** não têm template.
4. **Adoção:** concluída. A V2 substituiu a V1 em `docs/templates/`; a V1 está no histórico do git.
5. ~~Teste de conformidade permanente.~~ **Resolvido:** `go test ./docs/templates/` (seção 7).

## 6. Como reproduzir

Comandos na seção 7 de `docs/templates/README.md`.

## 7. Motores atualizados para seguir os templates

Pedido posterior: todos os motores passam a seguir os templates. Os templates viraram um pacote Go (legoloop/docs/templates, go:embed), a fonte única que os motores leem; cada motor tem teste de conformidade (escrito antes do código, visto falhar).

| Motor | Antes | Depois |
|:---|:---|:---|
| interview | Seções próprias (Contexto & Dor Original, Espelho Socrático, épicos E01/E02 fixos); gravava Status: Aprovado pelo Usuário sem passar pelo gate humano | Seções do TEMPLATE_INTENT; UNDER_REVIEW (ou DRAFT abaixo de 95%); nada inventado: o que a entrevista não respondeu fica "a definir" |
| specgen | Estrutura copiada em string no prompt (desatualizável); cabeçalho escrito pelo modelo | Prompt embute o template; cabeçalho carimbado pelo motor; placeholder sem preencher é rejeitado |
| plangen | ROUTE parcial; `output.gate` não compilável; `SupervisionEpic` nunca gerava gate; IDs TICKET-E01.1 que não casavam com as dependências; PLAN com a resposta crua do modelo | ROUTE com o conjunto de chaves do template; gates compiláveis; SG-FINAL; gates por nível; IDs e dependências canônicos; PLAN espelho; avisos de cobertura |
| upgrade | "Aplique o template correspondente" | Aponta o template do tipo da peça |

Mudanças de comportamento deliberadas (os testes antigos foram atualizados): IDs TICK-01.1; todo ROUTE termina em SG-FINAL (modo autônomo passa de 1 para 2 nós); o nível epic passa a gerar o gate de épico que prometia.

Fora de escopo, por decisão: canon.importer copia skills de terceiros verbatim (reescrevê-las seria errado) e studio/compose já gera DAG pelos tipos do dag.