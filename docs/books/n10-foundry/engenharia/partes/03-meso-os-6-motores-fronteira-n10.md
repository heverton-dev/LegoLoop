# Parte III: Os 6 Motores de Fronteira Industrial (Nível N10)

## Capítulo 3: Infraestrutura Transversal de Automação

Os 6 motores N10 garantem que o código gerado possua pureza estrutural, sem atalhos ou mocks artificiais.

### 1. Motor I — Anti-Reward Hacking (Red-Team vs Blue-Team)
A separação criptográfica em `internal/runner/redteam_blueteam.go` proíbe que o agente Blue-Team edite arquivos `*_test.go`. O hash SHA-256 do teste é travado no nó anterior, forçando o Blue-Team a escrever lógica real de produção.

### 2. Motor II — Fusão Determinística por AST (`MergeAST`)
O motor em `internal/runner/worktree.go` manipula diretamente a árvore de sintaxe abstrata (`go/ast`), permitindo que múltiplos workers paralelos mesclem declarações e funções sem criar conflitos de merge de git.

### 3. Motor III — Sistema Imunológico Episódico (`antibodies.db`)
O módulo em `internal/brain/antibodies.go` persiste triplas `(Assinatura, Causa Raiz, Diff)` geradas por falhas SBAR. Quando um agente opera no mesmo pacote, o contexto preventivo é automaticamente injetado.

### 4. Motor IV — Cascade Speculative FinOps Routing
O algoritmo em `internal/brain/speculative.go` dispara inferências especulativas em modelos rápidos com timeout restrito. Em caso de falha, aciona o modelo de alta capacidade injetando o erro prévio como contexto negativo.

### 5. Motor V — Caixa-Preta OpenTelemetry GenAI
A struct `FlightRecord` em `internal/runner/telemetry.go` mede o First-Pass Yield (FPY), tokens consumidos e tempo por nó, gravando `flight_record.json` para auditoria contínua.

### 6. Motor VI — Property-Based Testing & Fuzzing Nativo
Todos os crawlers e parsers do LegoLoop são submetidos a `testing.F` com milhares de execuções automatizadas em busca de panics ou vazamentos de memória.
