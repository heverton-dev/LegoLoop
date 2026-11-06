# 📋 SPEC: Governança Lazy Universal de Ferramentas e Conhecimentos no LegoLoop

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-06  
> **Versão:** 1.0.0  
> **Autor:** Antigravity & Usuário (Fusão Addy Osmani + Matt Pocock + AIDD-Spec)  
> **Referência da Intenção:** [docs/intents/2026-11-06_INTENT_GOVERNANCA_LAZY_UNIVERSAL.md](../intents/2026-11-06_INTENT_GOVERNANCA_LAZY_UNIVERSAL.md)  
> **Alvo Primário:** `internal/mcp`, `internal/assets`, `cmd/legoloop`  
> **Consome:** INTENT aprovada · **Produz:** PLAN e ROUTE  
> **Template:** TEMPLATE_SPEC.md v2.0.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)

> [!IMPORTANT]
> A IA lista abaixo todas as premissas técnicas adotadas para evitar decisões silenciosas. Alerte no Stop-Gate 02 caso alguma premissa esteja incorreta:
1. **Runtime Go Estrito:** 100% Go puro conforme `go.mod` (Go 1.24), utilizando exclusivamente a biblioteca padrão (`stdlib`), sem dependências externas adicionadas.
2. **Transporte STDIO Nativo:** O servidor `legoloop mcp-router` se comunica exclusivamente via entrada e saída padrão (`os.Stdin` e `os.Stdout`) utilizando o protocolo padrão JSON-RPC 2.0 (especificação oficial do MCP). Sem portas de rede HTTP, sem daemon residente, sem risco de conflito de portas ou regras de firewall no Windows.
3. **Cofre Canônico de MCPs:** Os servidores MCP upstream (como REA, Playwright, Mobbin) residem de forma declarativa e segura em `~/.config/mcp-vault/profiles/full.json` (ou `.mcp.json` de projeto), permanecendo desligados por padrão até serem invocados por demanda.
4. **Reuso de Infraestrutura Existente:** O roteador utiliza o scanner existente `internal/assets` para localizar skills e MCPs, e integra com o `codebase-memory-mcp` já presente para análise de AST e skeleton, sem reescrever do zero analisadores sintáticos.

---

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)

O pipeline e os agentes devem respeitar estritamente o sistema de fronteiras em 3 níveis:
* **Sempre Fazer:**
  * Desenvolver via TDD estrito (teste unitário falhando comprovadamente antes de qualquer código de produção).
  * Manter funções pequenas, enxutas e de responsabilidade única (< 35 linhas, verificado por `tools/funclen`).
  * Validar rigorosamente todos os inputs e decodificações de JSON-RPC, tratando chamadas desconhecidas com mensagens de erro JSON-RPC 2.0 bem formadas.
  * Tratar caminhos de forma compatível com Windows e POSIX (`filepath.Join` / `filepath.ToSlash`).
  * Recompilar o binário `bin/legoloop.exe` ao final de qualquer alteração em Go.
* **Perguntar Antes:**
  * Adicionar qualquer nova dependência externa no `go.mod`.
  * Alterar schemas persistidos de `~/.config/mcp-vault/` sem migração retrocompatível.
  * Sobrescrever arquivos de configuração de harnesses sem criar backup de segurança `.bak_legoloop`.
* **Nunca Fazer:**
  * Criar funções vazias com `TODO`, `pass`, `// implementar depois` ou mocks simulados em código de produção (Lei #1: Zero Stubs).
  * Abrir portas de rede locais TCP/HTTP na Fase 1 (restrito a stdio).
  * Commitar arquivos temporários, segredos ou chaves no git.
  * Ignorar testes falhando ou quebrar a suíte global `go test -count=1 ./...`.

---

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)

A solução cria o subsistema `internal/mcp` que atua como **Gateway / Proxy Inteligente JSON-RPC 2.0**.
Quando um harness de IA (Claude Code, MIMO Code, Cursor) inicializa, ele executa o comando `legoloop mcp-router` via stdio.

### Fluxo de Comunicação:
```
[ Harness de IA ]
       │
       ▼ (stdio JSON-RPC 2.0)
[ legoloop mcp-router ]
  ├── 1. Handshake: Exibe apenas 2 ferramentas (mcp_search, call_mcp_tool)
  ├── 2. mcp_search: Consulta ~/.config/mcp-vault/ e retorna schemas pontuais
  └── 3. call_mcp_tool: Spawna processo real (ex: rea.cmd mcp) e repassa resultado
```

### Costuras de Teste (Test Seams — Pocock)
* **Costura Primária (Highest Seam):**
  `func ServeSTDIO(r io.Reader, w io.Writer, vault VaultRegistry, runner ProcessRunner) error`
  Permite testar o ciclo de vida completo JSON-RPC (initialize, tools/list, tools/call) em memória injetando `bytes.Buffer` sem tocar em hardware real.
* **Costura de Busca de Skills:**
  `func (s *SkillEngine) Search(query string, topK int) []SkillMatch`
  Testa a pontuação BM25 e a injeção dos metadados de Tier, Natureza e Categoria.
* **Dublês Permitidos:**
  Apenas `io.Reader`/`io.Writer` simulados e um `ProcessRunner` falso para interceptar chamadas a processos externos nos testes unitários.

---

## 4. Contratos de Dados & Schemas Formais

```go
package mcp

import (
	"errors"
	"strings"
)

var (
	ErrRequisicaoInvalida = errors.New("mcp: requisicao json-rpc invalida")
	ErrFerramentaNaoAchada = errors.New("mcp: ferramenta nao encontrada no cofre")
)

// JSONRPCRequest modelo de entrada do protocolo MCP.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  map[string]any  `json:"params,omitempty"`
}

// JSONRPCResponse modelo de saida do protocolo MCP.
type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// SkillMatch modelo de retorno enriquecido da busca de skills.
type SkillMatch struct {
	ID        string  `json:"skill_id"`
	Nome      string  `json:"name"`
	Score     float64 `json:"score"`
	Natureza  string  `json:"natureza"`
	Tier      string  `json:"tier"`      // Ouro, Prata, Bronze
	Categoria string  `json:"categoria"` // build, reverse, web, mobile
	Resumo    string  `json:"resumo"`
}

func (s SkillMatch) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Nome) == "" {
		return ErrRequisicaoInvalida
	}
	return nil
}
```

---

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — Meta-MCP Router STDIO (Núcleo JSON-RPC)

* [ ] `AC-E01.1`: **Handshake e Listagem Mínima de Ferramentas**
  * **Dado que:** O harness envia mensagem JSON-RPC `{"method":"initialize"}` seguido de `{"method":"tools/list"}`.
  * **Quando:** O motor `ServeSTDIO` processa a entrada.
  * **Então:** Retorna exit code 0, protocolo `2.0` e o array `tools` contém **estritamente 2 ferramentas**: `mcp_search` e `call_mcp_tool`.
* [ ] `AC-E01.2`: **Busca Dinâmica de Ferramentas sob Demanda (`mcp_search`)**
  * **Dado que:** O cofre contém o servidor `rea` com 207 ferramentas desativadas.
  * **Quando:** O harness chama `mcp_search` com argumento `{"query":"decompile"}`.
  * **Então:** Retorna apenas as ferramentas do REA correlacionadas (ex: `decompile_jeb_item`), sem expor as outras 200 ferramentas não relacionadas.
* [ ] `AC-E01.3`: **Tratamento de Método Desconhecido**
  * **Dado que:** Uma mensagem com método inexistente `{"method":"invalido"}` é enviada.
  * **Quando:** O motor processa a requisição.
  * **Então:** Retorna erro JSON-RPC com código `-32601` (Method not found) sem quebrar o loop de leitura.

### Épico 02 — Lazy Skill Engine com Metadados Enriquecidos

* [ ] `AC-E02.1`: **Busca de Skill com Retorno de Taxonomia Completa**
  * **Dado que:** O catálogo de skills canônicas está indexado com metadados de Tier e Natureza.
  * **Quando:** A função `SkillEngine.Search("testes tdd", 1)` é invocada.
  * **Então:** O retorno contém a skill `aidd-tdd` com campos preenchidos: `tier: "Ouro"`, `natureza: "Processo de Engenharia"`, `categoria: "build"` e `resumo` de 1 linha.
* [ ] `AC-E02.2`: **Poda de Pastas Upstream Ocultas**
  * **Dado que:** Existem pastas com prefixo de ponto como `.upstream`.
  * **Quando:** O scanner de skills indexa o diretório.
  * **Então:** Nenhuma skill dentro de pastas com ponto entra no índice retornado.

### Épico 03 — Sincronizador de Harnesses (`legoloop sync-harnesses`)

* [ ] `AC-E03.1`: **Distribuição das Configurações para Harnesses**
  * **Dado que:** Os arquivos `~/.claude.json`, `.mimocode/mimocode.jsonc` e `.cursor/mcp.json` existem no SO.
  * **Quando:** O comando `legoloop sync-harnesses --dry-run=false` é disparado.
  * **Então:** Gera backup `.bak_legoloop` de cada arquivo e atualiza a chave de MCPs para apontar exclusivamente para `legoloop mcp-router`.

---

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)

### Comandos de Verificação Mecânica
* **Build:** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Test:** `go test -v -count=1 ./...`
* **Estática e Linters:** `go vet ./...` e `go run ./tools/funclen internal` (0 funções > 35 linhas).
* **Formatação:** `gofmt -l .` (saída vazia).
* **Dependências intactas:** `git diff --exit-code go.mod go.sum`

### Critério de Done Global
A especificação técnica é considerada **100% CONCLUÍDA** quando:
- [ ] 100% dos critérios binários `AC-E{N}.{M}` estiverem atendidos (`[x]`).
- [ ] O comando de **Test** passar com exit code 0 em todos os pacotes.
- [ ] `go run ./tools/funclen internal` validar 100% das funções com <= 35 linhas.
- [ ] O comando de **Build** recompilar com sucesso o binário `bin/legoloop.exe`.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**  
*As premissas assumidas, as fronteiras (boundaries), os contratos e os critérios binários acima atendem à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar o PLAN (`TEMPLATE_PLAN.md`) e a ROUTE (`TEMPLATE_ROUTE.json`).*
