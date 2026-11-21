# PARTE II — Execução em Swarm e Remediação Técnica

## 1. Blindagem de Concorrência, MCP STDIO e Path Jail

Durante a Etapa 4 e o fechamento dos tickets estruturais, falhas intermitentes sob carga concorrente foram eliminadas na raiz:

### Inventário Assíncrono sem Timeout EOF
A varredura inicial de harnesses e binários (`internal/assets/watch.go` e `internal/assets/system.go`) podia expirar conexões HTTP no Studio.
A correção introduziu três garantias:

- `Watcher.EnsureInitial` aplica double-checked locking para impedir varreduras frias duplicadas.
- `scanHarnesses` e `scanBinaries` executam em paralelo via `wg.Go`.
- Os handlers HTTP de `/api/assets` estendem o prazo de escrita (`extendWriteDeadline`) para três minutos e `newStudioWatcher` isola variáveis de ambiente por `d.Home`.

### Sincronização no Executor MCP STDIO
Em `internal/orchestrix/backends.go` e `internal/orchestrix/mcp_client.go`, `limitedBuffer` passou a usar `sync.Mutex` com cópia defensiva de slice em `Bytes()`.
A função `finishMCPError` aguarda o encerramento do processo filho (`cmd.Wait()`, teto de `500ms`) antes de formatar `stderr`, eliminando leituras incompletas em falhas súbitas.

### Substituição Real em Mídia e Path Jail no Studio

- Em `internal/motion/sync.go`, `librosaDetector.Detect` deixou de retornar `nil, nil` e passou a acionar `scripts/beat_detect.py` com leitura JSON de batidas (`ErrBeatDetectFailed` em caso de falha).
- Em `internal/studio/translate.go`, `POST /api/assets/translate` foi confinado às raízes de `revealRoots()` via `os.OpenRoot` e `checkNoLinksRoot`, mascarando caminhos com `<projeto>`, `~` e `<cofre>`.
- Em `internal/studio/server.go`, `writeOrchestratorSkill` aplica `isSafeSkillSlug` e `os.OpenRoot`.

## 2. Paridade de Contratos, CLI e Taxonomia

Para garantir rastreabilidade nos pacotes `internal/orchestrix` e `internal/studiodoc/`, os campos `$schema`, `description`, `schema_version`, `date`, `upstream_signatures` e `execution_strategy` foram preservados na serialização.

Nos comandos de terminal (`internal/cli/`):

- `legoloop swarm` ganhou suporte a `--help` e `-h` retornando código `0`.
- `legoloop sync-harnesses` passou a emitir saída estruturada com `--json`.
- `internal/cli/instalar_skills_test.go` passou a validar todas as 23 skills nativas, incluindo `legoloop-homologacao`.

Em `docs/routes/` e `docs/dags/`, as ocorrências restantes de termos fora do dicionário foram normalizadas com paridade de hash SHA-256 validada por `tools/taxlinter`.
