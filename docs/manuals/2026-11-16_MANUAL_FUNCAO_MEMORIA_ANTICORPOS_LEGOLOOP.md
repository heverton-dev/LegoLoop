# Manual de Função — Memória de Execução e Anticorpos (Lei 11)

- **Data:** 2026-11-16
- **Demanda:** MEMORIA_CONCORRENTE_ANTICORPOS (integrada ao Mega Plano nos Épicos 03/04 e na Lei 11)
- **Pacotes:** `internal/triagem`, `internal/brain`, `internal/fechar`, `internal/runner`, `internal/cli`
- **Norma:** Lei 11 do AGENTS.md (Governança Viva) e Lei 15 (ASD-STE100)

---

## 1. Propósito

Este manual opera a memória de execução e o banco de anticorpos do LegoLoop.
A memória de execução guarda o progresso de cada sessão de um Processo.
O banco de anticorpos guarda erros conhecidos e o antídoto de cada erro.
Juntos, eles dão controle de progresso e evitam reincidência de erro.

## 2. Componentes

| Componente | Pacote | Função |
| :--- | :--- | :--- |
| Memória de execução | `internal/triagem` | Grava `memoria_execucao.json` com lock por arquivo e sessão |
| Anticorpos | `internal/brain` | Persiste `anticorpos.json` no schema `anticorpos.v1` com escopo global e de domínio |
| Quarteto Lei 11 | `internal/fechar` | Zera a memória da sessão ao fechar o Processo |
| Memória em runtime | `internal/runner` | Grava `etapa_atual` e `proximo_passo` a cada nó do DAG |
| CLI de sessões | `internal/cli` | Comando `legoloop sessoes` lista as sessões ativas |

## 3. Operação

### 3.1 Consultar sessões ativas

```
legoloop sessoes <slug> [--json] [--processos-dir DIR]
```

- A saída JSON traz `sessao_id`, `pid`, `etapa_atual` e `proximo_passo`.
- A saída nunca expõe caminho absoluto (Lei 8).
- Processo sem sessão viva devolve `[]` e exit 0.

### 3.2 Ler e gravar a memória de execução

- `triagem.SalvarMemoria(dir, mem)` grava a sessão com lock por arquivo.
- `triagem.LerMemoria(dir, sessaoID)` lê uma sessão; sem ID, lê a legada.
- `triagem.ListarSessoesAtivas(dir)` devolve só as sessões com PID vivo.
- `triagem.ZerarSessao(dir, sessaoID)` zera a sessão no fechamento.
- `triagem.LimparOrfaos(dir)` remove sessões de PID morto.

### 3.3 Persistência de anticorpos

- `brain.NewAntibodiesDBPersistenteComGlobal(dir, dirGlobal)` abre o banco com
  anticorpos globais compartilhadas entre Processos.
- `brain.RecordComEscopo` grava anticorpo global ou de domínio.
- `brain.GateAntiReincidencia(pkg, errSig)` decide o fluxo: `AcaoFluxoNormal`,
  `AcaoAntidoto` ou `AcaoBloqueado`. O bloqueio devolve `ErrGateAntiReincidencia`.
- `brain.InjectPreventiveContext` injeta o antídoto no prompt do nó.

### 3.4 Schemas de persistência

- `memoria_execucao.v1`: envelope com `schema_version`, `processo` e `sessoes`.
- `anticorpos.v1`: banco com `globais` e `dominio`, cada anticorpo com
  `assinatura` SHA-256 no formato `sha256:<hex>`.

## 4. Recuperação de falha

| Sintoma | Causa | Ação |
| :--- | :--- | :--- |
| `ErrLockOcupado` | Outro processo detém o lock | Aguarde ou confira `pid` do lock |
| Lock órfão | PID do lock morto | O runtime toma o lock sozinho |
| `ErrMemoriaCorrompida` | JSON inválido | O arquivo vira `<arquivo>.legoloop-bak` e o sistema segue |
| `ErrAnticorposCorrompidos` | Banco inválido | Backup criado; banco vazio funcional é liberado |
| `ErrGateAntiReincidencia` | Erro conhecido sem antídoto novo | Corrija a causa raiz e registre o antídoto |

## 5. Checklist Lei 11

- [x] Memória de execução com lock por arquivo e sessão.
- [x] Lock órfão tomado quando o PID está morto.
- [x] Backup `.legoloop-bak` em arquivo corrompido.
- [x] Anticorpos persistidos no schema `anticorpos.v1`.
- [x] Escopo global e de domínio de anticorpos.
- [x] Gate anti-reincidência com antídoto ou bloqueio.
- [x] Injeção preventiva de antídoto no prompt.
- [x] Gravação de `etapa_atual` e `proximo_passo` a cada nó.
- [x] CLI `legoloop sessoes` para controle de progresso.
- [x] Zeramento da sessão no Fechamento Canônico.

## 6. Portões de qualidade

- `go test -count=1 ./...`
- `go vet ./...`
- `gofmt -l .` com saída vazia
- `go run ./tools/funclen internal` sem estouro novo
- `go run ./tools/taxlinter internal cmd loops skills docs/templates` sem violação nova
- `git diff --exit-code -- go.mod go.sum`
