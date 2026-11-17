# 📊 REPORT: BUILD do Épico 01 (Núcleo Go da Fábrica) — `SG-E01`

> **Data:** 2026-11-03  
> **Escopo:** `TICK-01.1` a `TICK-01.8` do PLAN `2026-11-03_PLAN_LEGOLOOP_FABRICAR.md`  
> **Estado:** 8 tickets `PASSED`; aguardando a autorização humana do `SG-E01` para iniciar o Épico 02.  

---

## 1. Entregas

| Ticket | Entrega | Arquivos principais |
|:---|:---|:---|
| `TICK-01.1` | `NormalizeProcessName`, `Slug`, `ErrInvalidName`, fuzz; `doc.go` dos pacotes `fabricar` e `canon` | `internal/fabricar/name.go`, `doc.go`, `internal/canon/doc.go` |
| `TICK-01.2` | `DiscoverDrives`, `DiscoverDrivesDetailed`, `ParseMounts`, `MacVolumes`, `SystemRoots`, `DriveProbe` | `internal/fabricar/drives.go` |
| `TICK-01.3` | `ScanStorageLocations`, `DiskScanner`, `Untrusted`, limites com aviso | `scanner.go`, `pieces.go`, `untrusted_windows.go`, `untrusted_other.go` |
| `TICK-01.4` | `Materialize` atômico, `FS` injetável, `Force`, `DryRun`, rollback | `materialize.go` |
| `TICK-01.5` | `plangen.Render` (PLAN e ROUTE sem LLM, data injetada) e `Generator.Now` | `internal/plangen/plangen.go` |
| `TICK-01.6` | `EmitSkill`, `EmitLoop`, `EmitDag` | `emit_skill.go`, `emit_loop.go`, `emit_dag.go` |
| `TICK-01.7` | `Request.Validate`, `Clarifier`, `ArchitectMerge` | `request.go`, `clarifier.go`, `merger.go` |
| `TICK-01.8` | `Fabricar` ponta a ponta, Tríade, Template, link opt-in | `engine.go`, `triad.go`, `template.go`, `link.go` |

Aberturas mínimas nos motores reaproveitados (sem mudar comportamento nem schema): `interview.RenderIntent`, `IntentView` e `SpecRelFor` exportados; `specgen.Generator.Now`; `plangen.Generator.Now` e `plangen.Render`.

## 2. Checklist do `SG-E01`

| Item | Resultado |
|:---|:---|
| 100% dos testes de `internal/fabricar`, incluindo o fuzz de 20 s | **55 testes passando, 0 falhas, 0 pulados**; fuzz 4 rodadas de 20 s (mais de 10 milhões de execuções) sem falha de propriedade |
| `go test ./...` sem regressão | Verde em todos os 21 pacotes, sem cache |
| `go vet ./...` | Limpo |
| `go run ./tools/funclen internal/fabricar` | Limpo (também `internal/plangen`) |
| `gofmt -l .` com saída vazia | **Não é verdade no repositório inteiro**: 22 arquivos preexistentes não estão formatados (`internal/cli/commands.go`, `internal/dag/dag_test.go`...). **Nenhum é meu**: os pacotes que toquei estão limpos. Não os reformatei por estar fora do escopo |
| `go.mod` e `go.sum` inalterados | Confirmado (`git diff --exit-code`) |
| `-race` | Indisponível aqui (exige CGO), como o PLAN previa. Obrigatório no CI Linux |

Também: `go vet` cruzado em `GOOS=linux` e `GOOS=darwin` valida os build tags, e `bin/legoloop.exe` foi recompilado a cada ticket.

## 3. Evidências além dos testes

* **Mutação (os testes detectam regressão):** desliguei o rollback do `Materialize` e 2 testes falharam; fiz o `ArchitectMerge` aceitar `Untrusted` e 2 testes falharam. Código restaurado nos dois casos.
* **Teste de fumaça real:** `Fabricar` com o scanner de disco verdadeiro, nesta máquina. Detectou `C:\`, `D:\` e `F:\`, curou 5 skills reais (a melhor, `aidd-legoloop-mega-pipeline`, tier ouro), gravou 8 arquivos em diretórios temporários e o **binário real** validou o DAG: `✓ ... 5 nós válidos, sem ciclos`, exit 0. O cofre real `~/.aidd` não foi tocado.
* **Laço de diretório:** o teste usa symlink e, no Windows sem privilégio, **junction**. Roda de verdade aqui, em vez de pular.

## 4. Achados e desvios (leia antes de autorizar)

1. **Fuzz: uma falha inicial que não reproduziu.** A primeira execução do fuzz falhou; o caso salvo passa quando reexecutado e 4 rodadas depois não falharam. Não consegui identificar a causa (suspeito do worker do fuzz com a máquina compilando em paralelo). O caso ficou em `testdata/fuzz` como regressão.
2. **Desempenho da varredura real:** a primeira execução a frio levou **60 s** (o `C:\` tem milhares de skills no perfil do usuário); a quente leva **6,6 s**. Nenhuma raiz estourou o prazo. Calcular o SHA-256 só depois do corte de 500 peças reduz leitura de disco, mas **não mudou o tempo**: o custo está na varredura e pontuação. Se 6 a 60 s for demais para a CLI, a saída é limitar por tempo total, não por raiz.
3. **`ScanOptions.MaxDepth` e `MaxFiles`:** como o scanner reaproveita o `internal/assets` (sem segundo walker), esses limites são aplicados sobre o resultado (profundidade relativa ao perfil, itens inspecionados), não durante a caminhada. O `assets` já limita a própria profundidade.
4. **`Untrusted` no Windows:** disco fixo é confiável; removível, rede, CD e RAM disk não. Limite conhecido: alguns HDs externos USB se anunciam como fixos e seriam tratados como confiáveis.
5. **Extensões ao contrato da SPEC:** `Deps.Ask` (a CLI e o Studio entregam as perguntas), `ProcessSpec`/`Step` (entrada dos emissores) e `Merged`/`MergeOptions`. `HarnessDirs` é uma lista de pastas; o harness `claude` casa com a pasta que contém `/.claude/`.
6. **Modo autônomo não entrevista:** a INTENT gerada nasce `DRAFT` com confiança `0%`, em vez de inventar um número.
7. **Peça hostil:** nomes de peças varridas são sanitizados e só viram `ref`/ID (nunca `command`, `slot` ou `hook`); peça de mídia removível fica fora da fusão por padrão.
8. **`internal/canon/doc.go`:** movi o comentário de pacote de `hasher.go` para `doc.go` (sem mudar comportamento), para registrar a fronteira `canon` × `fabricar`.

## 5. Próximo passo

`SG-E01` pede sua autorização explícita para iniciar o **Épico 02** (CLI `legoloop fabricar` e skill do Chat). Nada do Épico 02 foi iniciado.
