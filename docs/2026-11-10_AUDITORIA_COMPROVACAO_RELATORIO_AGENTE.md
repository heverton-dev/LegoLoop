# 🔍 Relatório de Auditoria de Comprovação: Alegações do Relatório Externo vs. Código Real

> **Status:** Concluído  
> **Data:** 2026-11-10  
> **Ambiente:** Windows 11 · Go 1.24 / stdlib · Repositório LegoLoop  
> **Tipo:** Auditoria Forense e Comprovação Mecânica  
> **Arquivo Base Auditado:** Relatório de Auditoria sobre `AUDITORIA_FONTE.txt`  

---

## Sumário Executivo

A presente auditoria confrontou, de forma cirúrgica e mecânica, cada uma das alegações apresentadas no relatório externo contra a base de código real do ecossistema **LegoLoop**.

| Seção | Total de Casos | Procedentes | Procedentes c/ Ressalva | Problemas Refutados com Provas |
| :--- | :---: | :---: | :---: | :---: |
| **1. Erros Críticos de Implementação** | 3 | 2 | 1 | 0 |
| **2. Defeitos Estruturais em Planos/Rotas** | 4 | 0 | 0 | 4 |
| **3. Dívida Técnica e Divergências** | 3 | 2 | 0 | 1 |
| **4. Pontos Cegos de Arquitetura/Portabilidade** | 5 | 5 | 0 | 0 |
| **TOTAL** | **15** | **9** | **1** | **5** |

> **Diagnóstico Geral:**  
> O relatório do agente misturou **defeitos reais de implementação** (como o stub de ingestão no `video_rea.go`, a falta de jail em paths e duplicações de esteira) com **alucinações/leituras desatualizadas da governança** (alegação de que faltavam stop-gates nas rotas, que `docs/templates/` não estava versionado ou que faltava a pasta `skills/`, todas refutadas com provas irrefutáveis do repositório).

---

## 1. Erros Críticos de Implementação

### 1.1. Falso `ingest` no Módulo `video-rea`
* **Caso Apresentado:** A função `runIngest` promete fazer a ingestão de um vídeo, mas apenas captura `--url`, instancia `VideoInput{URL: *rawURL}`, salva o JSON e encerra com exit code 0 sem chamar `yt-dlp` ou baixar áudio/metadados.
* **Código Real:** [internal/cli/video_rea.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/cli/video_rea.go#L70-L84) e [loops/legoloop-video-rea/pipeline.dag.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/loops/legoloop-video-rea/pipeline.dag.json#L35-L85).
```go
func runIngest(args []string, d Deps) int {
	fs := newFlags("video-rea ingest", d)
	rawURL := fs.String("url", "", "URL do video")
	out := fs.String("out", "", "arquivo de saida")
	if err := fs.Parse(args); err != nil || *rawURL == "" {
		fmt.Fprintf(d.Stderr, "erro: --url é obrigatório\n")
		return 2
	}
	if _, err := videorea.DetectPlatform(*rawURL); err != nil {
		fmt.Fprintf(d.Stderr, "erro na plataforma: %v\n", err)
		return 1
	}
	vInput := videorea.VideoInput{URL: *rawURL}
	return saveJSONOutput(vInput, *out, d)
}
```
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** O nó `n01_crawler_ingest` em `loops/legoloop-video-rea/pipeline.dag.json` intitula-se *"Ingestão e Download de Metadados e Áudio Bruto"*, mas a função Go invocada (`runIngest`) apenas persiste a URL pura no arquivo de saída sem realizar qualquer download ou chamada a extratores (`yt-dlp`). O motor Go em `internal/videorea/crawler.go` possui apenas o método `DetectPlatform` (checagem de string de host), sem nenhuma implementação de crawler de áudio ou metadados.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-FIX-VREA-01`
  * **Alvo:** `internal/videorea/crawler.go` e `internal/cli/video_rea.go`
  * **Ação:** Implementar em `videorea.DownloadMediaAndMetadata(ctx, rawURL, destDir)` a invocação subprocessada segura de `yt-dlp` para extrair áudio `.aac`/`.m4a` e arquivo `.info.json`, atualizando o struct `VideoInput` com `AudioPath`, `Title`, `Description` e `Duration`.

---

### 1.2. Ausência de Validação Estrita na Forense (`executeForensics`)
* **Caso Apresentado:** Após chamar `videorea.ExtractForensicsWithBrain(...)`, o código apenas verifica se `err != nil` e salva o arquivo com `saveJSONOutput(dossier, out, d)`. Se a inferência falhar silenciosamente ou o input for vazio, aceita dossiês vazios/nulos (`explicit_prompts: null`, etc.).
* **Código Real:** [internal/cli/video_rea.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/cli/video_rea.go#L107-L132) e [internal/videorea/forensics.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/videorea/forensics.go#L48-L75).
```go
	dossier, err := videorea.ExtractForensicsWithBrain(ctx, vInput, text, desc, b)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro extraindo forense: %v\n", err)
		return 1
	}
	return saveJSONOutput(dossier, out, d)
```
* **Veredito:** **PROCEDE COM RESSALVA.**
* **Prova Comprobatória:** Em `internal/videorea/forensics.go:50-52`, a checagem é apenas `if strings.TrimSpace(transcript) == "" && strings.TrimSpace(description) == "" { return nil, ErrEmptyDossier }`. Se a descrição fornecida for apenas `"Vídeo analisado via <url>"` (gerada na linha 125 do CLI), a função avança e produz um `ForensicsDossier` com slices nulas/vazias. Não existe o método `dossier.ValidateStrict()` nem um portão binário de consistência mínima antes de `saveJSONOutput`.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-FIX-VREA-02`
  * **Alvo:** `internal/videorea/types.go` e `internal/cli/video_rea.go`
  * **Ação:** Criar o método `(d *ForensicsDossier) ValidateStrict() error` exigindo ao menos uma das seguintes evidências: `len(ToolsIdentified) > 0`, `len(HiddenTools) > 0` ou `len(OCRFindings) > 0`. No CLI `executeForensics`, abortar com exit code 1 se a validação estrita falhar.

---

### 1.3. Poluição da Raiz e Falta de Path Jail (`bundle.go`)
* **Caso Apresentado:** Ferramentas de I/O aceitam caminhos arbitrários sem sandbox estrito (`filepath.Clean` + checagem de prefixo em `docs/videorea/<slug>/`). O modelo gravou arquivos em `scratch/` diretamente na raiz em vez de confiná-los no bundle.
* **Código Real:** [internal/videorea/bundle.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/videorea/bundle.go#L22-L37) e [loops/legoloop-video-rea/pipeline.dag.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/loops/legoloop-video-rea/pipeline.dag.json#L65-L79).
```go
func NewTargetBundle(baseDir, slug string) *TargetBundle {
	base := filepath.Join(baseDir, slug)
	return &TargetBundle{
		Slug:               slug,
		BaseDir:            base,
		ScratchDir:         filepath.Join(base, "scratch"),
...
```
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** 
  1. Em `internal/videorea/bundle.go`, `NewTargetBundle` concatena diretamente caminhos com `filepath.Join` sem sanitizar com `filepath.Clean` e sem validar se o caminho resultante escapa de `baseDir` via traversals (`../`).
  2. No próprio DAG `loops/legoloop-video-rea/pipeline.dag.json` (linhas 65 e 75), os comandos fixam saída em `"scratch/videorea/raw.json"` (na raiz do repositório) em vez de respeitar a hierarquia do bundle em `docs/videorea/<slug>/scratch/`.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-FIX-VREA-03`
  * **Alvo:** `internal/videorea/bundle.go` e `loops/legoloop-video-rea/pipeline.dag.json`
  * **Ação:** Adicionar função de confinamento seguro `sanitizeBundlePath(baseDir, slug string) (string, error)` aplicando `filepath.Clean` e rejeitando qualquer prefixo fora de `baseDir`. Ajustar os nós do DAG do `video-rea` para utilizarem o caminho relativo ao workspace do bundle.

---

## 2. Defeitos Estruturais em Planos, Rotas e Gates (`PLAN` / `ROUTE`)

### 2.1. Injeção de Comandos Shell dentro da DSL de Gate
* **Caso Apresentado:** Na V1 da rota (`2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json`), os nós declaravam gates como `"exit_code == 0 && go test -v -run ..."`, falhando no compilador de gates que não é um shell.
* **Código Real:** [docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json#L112)
```json
"gate": "exit_code == 0 && file_exists('internal/fabricar/name.go') && file_exists('internal/fabricar/name_test.go') && file_exists('internal/fabricar/doc.go') && file_exists('internal/canon/doc.go')"
```
* **Veredito:** **PROBLEMA REFUTADO COM PROVAS.**
* **Prova Comprobatória:** Uma busca exaustiva no arquivo real `docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json` comprova que **nenhum** nó injetou comandos shell dentro da chave `"gate"`. Em todos os nós (linhas 112, 194, 280, 362, 449, 536, 623, 718, etc.), os comandos de teste estão estritamente contidos na chave `"verify"` como array de argumentos (`["go", "test", "-v", ...]`), e a chave `"gate"` contém unicamente expressões booleanas da DSL (`exit_code == 0 && file_exists(...)`). A alegação do agente é uma alucinação completa.

---

### 2.2. Testes Cosméticos ("Gates que Não Provam Nada") no `FAB-E03.2`
* **Caso Apresentado:** O ticket `FAB-E03.2` (V1) validava a criação de um modal no frontend rodando apenas `TestZeroEmojisInFrontendAndTemplates|TestNoNativeDialogs`, aprovando a esteira sem escrever o modal.
* **Código Real:** [docs/plans/2026-11-03_PLAN_LEGOLOOP_FABRICAR.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-03_PLAN_LEGOLOOP_FABRICAR.md#L351-L368) e [internal/studio/fabricar_modal_test.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/fabricar_modal_test.go).
* **Veredito:** **PROBLEMA REFUTADO COM PROVAS.**
* **Prova Comprobatória:** O relatório do agente copiou a nota histórica de autópsia da V1 e a apresentou como defeito ativo. No plano e esteira vigentes (linhas 351–368 de `2026-11-03_PLAN_LEGOLOOP_FABRICAR.md`), foi criado o ticket cirúrgico `TICK-03.3` com o teste real `TestFabricarModalStatic` no arquivo `internal/studio/fabricar_modal_test.go`, que inspeciona o DOM, os IDs `btn-novo-processo`, `modal-fabricar`, `fabricar-name`, `fabricar-model` e a ausência de diálogos nativos. O problema foi sanado antes da execução da V2.

---

### 2.3. Omissão de *Stop-Gates* Humanos no Grafo Executável
* **Caso Apresentado:** O documento Markdown (`PLAN` V1) prometia paradas humanas obrigatórias ao fim dos Épicos 01 e 02, mas a `ROUTE` V1 em JSON omitiu os nós `stop_gate`, atropelando fases sem autorização humana.
* **Código Real:** [docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json#L721-L745).
```json
    {
      "id": "SG-E01",
      "label": "Stop-Gate Humano do Épico 01",
      "nature": "nano",
      "type": "stop_gate",
      "ref": null,
      "status": "PASSED",
      "execution_mode": "sequential",
      "routing": {
        "mode": "manual",
        "harness": "local",
        "model": "human/operator"
      }
    }
```
* **Veredito:** **PROBLEMA REFUTADO COM PROVAS.**
* **Prova Comprobatória:** A inspeção do JSON real `docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json` revela que todos os 4 stop-gates humanos prometidos constam explicitamente no arquivo:
  1. `SG-E01` (linha 722): *"Stop-Gate Humano do Épico 01"*, `type: "stop_gate"`;
  2. `SG-E02` (linha 1007): *"Stop-Gate Humano do Épico 02"*, `type: "stop_gate"`;
  3. `SG-E03` (linha 1438): *"Stop-Gate Humano do Épico 03 (homologação visual)"*, `type: "stop_gate"`;
  4. `SG-FINAL` (linha 1709): *"Stop-Gate Final: encerramento da demanda"*, `type: "stop_gate"`.  
  Todos configurados com `mode: "manual"` e `model: "human/operator"`. A alegação de omissão é categoricamente falsa.

---

### 2.4. Incompatibilidade de Schemas entre `ROUTE.json` e `dag.Parse`
* **Caso Apresentado:** A SPEC V1 assumia que `ROUTE.json` seria executado por `legoloop run`. Mas `dag.Parse` usa `DisallowUnknownFields` e rejeita `ROUTE.json`, forçando separar `ROUTE.json` de `.dag.json`.
* **Código Real:** [internal/dag/parse.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/parse.go#L47) e [AGENTS.md (Lei #6)](file:///C:/Users/trcnologia/Desktop/LegoLoop/AGENTS.md#L36-L42).
* **Veredito:** **PROBLEMA REFUTADO COM PROVAS (Decisão Arquitetural Canônica Conhecida).**
* **Prova Comprobatória:** Essa separação não é um defeito do código, mas um dos pilares de governança estabelecidos na **Lei #6 do AGENTS.md (Quinteto de Governança ISO)**. O arquivo `ROUTE.json` é um artefato rico de governança humana/agente (contendo rastreabilidade EARS, naturezas, tooling e personas), enquanto o `pipeline.dag.json` é o contrato estrito de execução consumido pelo compilador Go. A ferramenta `legoloop fechar` audita mecanicamente a paridade de nós e SHA-256 entre ambos. O agente tratou a arquitetura do sistema como um erro.

---

## 3. Dívida Técnica e Divergências no Repositório Atual

### 3.1. Violação da Regra de Tamanho de Função (`tools/funclen`)
* **Caso Apresentado:** A governança exige funções < 35 linhas, mas `tools/funclen` já falha hoje em funções legadas em `internal/studio`, `internal/cli` e `internal/canon`. O gate foi restrito a pacotes novos.
* **Código Real:** Execução comprovada de `go run ./tools/funclen internal`.
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** A execução em tempo real de `go run ./tools/funclen internal` retornou **Exit Code 1** acusando exatamente **23 funções** que violam o teto de 35 linhas:
  * `internal/assets/skills.go:170`: `inferOntologyFallback` (43 linhas)
  * `internal/board/client.go:42`: `FetchReadyCards` (48 linhas)
  * `internal/canon/importer.go`: 3 funções (linhas 36, 83, 159 — até 72 linhas)
  * `internal/reputation/scorer.go:70`: `ScorePiece` (102 linhas)
  * `internal/runner/runner.go:161`: `node` (113 linhas)
  * `internal/runner/triggers.go:134`: `Listen` (37 linhas)
  * `internal/studio/`: 11 funções (entre elas `templates.go:20 BuiltinTemplates` com 454 linhas e `handleRunDAG` com 95 linhas)
  * `internal/suggest/suggest.go:188`: `neighborSuggestions` (44 linhas)
  * `internal/wave/`: 2 funções (`ExtractAffectedFiles` e `PlanWaves`)
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-DEBT-FUNCLEN-01`
  * **Alvo:** `internal/canon/importer.go`, `internal/assets/skills.go` e `internal/studio/`
  * **Ação:** Decompor `BuiltinTemplates` (454 linhas) carregando os templates a partir dos arquivos `.json` existentes em `loops/` via `embed.FS`, e decompor as funções de roteador HTTP do Studio em handlers menores.

---

### 3.2. Divergências de Nomenclatura e Estado Pré-Build
* **Caso Apresentado:**
  1. Conflito `canonize` vs `fabricar`: a spec seria `SPEC_LEGOLOOP_CANONIZER_V2.md` com `canonize <nome>` e exigiria patch V2.1.
  2. Templates não versionados: `docs/templates/` não estaria versionado no Git.
  3. Diretório `skills/` inexistente na raiz do repositório.
* **Código Real:** `git ls-files docs/templates`, `git ls-files skills` e [docs/specs/2026-11-04_SPEC_FABRICAR_V2.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-04_SPEC_FABRICAR_V2.md).
* **Veredito:** **PROBLEMA REFUTADO COM PROVAS.**
* **Prova Comprobatória:**
  1. O arquivo `SPEC_LEGOLOOP_CANONIZER_V2.md` sequer existe. O documento canônico é [2026-11-04_SPEC_FABRICAR_V2.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-04_SPEC_FABRICAR_V2.md), aprovado em 04/11/2026, com o comando `legoloop fabricar` totalmente padronizado e unificado.
  2. O comando `git ls-files docs/templates` comprovou que **todos os 26 arquivos de templates** (incluindo `TEMPLATE_ROUTE.json`, `templates.go`, `templates_test.go` e os SVGs) estão rigorosamente versionados e rastreados pelo Git.
  3. O diretório `skills/` **existe**, está versionado no Git com `skills/embed.go` e contém 14 diretórios de skills canônicas universais (`ai-motion-graphics`, `legoloop-compor`, `legoloop-fabricar`, `legoloop-fechar`, etc.).

---

### 3.3. Descompasso no Pipeline `Impeccable`
* **Caso Apresentado:**
  1. Gates fracos no DAG: manual promete 4 gates rigorosos de julgamento, mas `impeccable-design.loop.json` usa apenas `file_exists(...)`.
  2. Roteamento incompleto: últimos 4 nós não declaram modelo.
  3. Duplicação JSON vs Go: esteira declarada em `loops/impeccable-design.loop.json` e hardcoded em `internal/studio/templates.go`.
  4. Typo no manual: `impeccabe-finish-reviewer`.
* **Código Real:** [loops/impeccable-design.loop.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/loops/impeccable-design.loop.json), [internal/studio/templates.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/templates.go#L188-L245) e [docs/manuals/2026-11-01_MANUAL_IMPECCABLE_DESIGN_LOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/manuals/2026-11-01_MANUAL_IMPECCABLE_DESIGN_LOOP.md#L66).
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:**
  1. Todos os nós de `impeccable-design.loop.json` possuem como único gate `file_exists(...)` (ex.: `file_exists('src/styles/theme.css')`).
  2. Os 4 últimos nós omitem a chave `"model"`, registrando apenas `"combo": "code-fast"`.
  3. Em `internal/studio/templates.go:188-245`, a composição `impeccable-design` está inteiramente redigitada em Go, divergindo do JSON (no Go o nó `subloop-avaliar` está com `ModeAuto`, enquanto no JSON está com `ModeSuggested` e `sonnet`).
  4. Na linha 66 de `docs/manuals/2026-11-01_MANUAL_IMPECCABLE_DESIGN_LOOP.md` lê-se textualmente: `D2 --> D3["Revisão pelo impeccabe-finish-reviewer"]`.
* **Micro-Tickets Cirúrgicos:**
  * **Ticket ID:** `TICK-FIX-IMP-01`
  * **Alvo:** `docs/manuals/2026-11-01_MANUAL_IMPECCABLE_DESIGN_LOOP.md:66`
  * **Ação:** Corrigir typo para `impeccable-finish-reviewer`.
  * **Ticket ID:** `TICK-FIX-IMP-02`
  * **Alvo:** `internal/studio/templates.go` e `loops/impeccable-design.loop.json`
  * **Ação:** Unificar a SSOT fazendo `internal/studio` deserializar `impeccable-design.loop.json` via `embed.FS` em vez de manter uma struct duplicada em Go. Evoluir os gates para validar conformidade CSS (`has_pattern('src/styles/theme.css', ':root')`).

---

## 4. Pontos Cegos de Arquitetura, Portabilidade e Segurança

### 4.1. Acoplamento ao Ambiente Windows / PowerShell nas Rotas
* **Caso Apresentado:** As rotas fixam `"runtime": "powershell"` e `"bin/legoloop.exe"`, quebrando a execução em runners Linux ou macOS.
* **Código Real:** [docs/routes/2026-11-04_ROUTE_FABRICAR_V2.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-04_ROUTE_FABRICAR_V2.json).
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** O arquivo `2026-11-04_ROUTE_FABRICAR_V2.json` repete 42 ocorrências do caminho binário fixo `"bin/legoloop.exe"` e define `"runtime": "powershell"` nos nós de comando. Em ambientes Linux/macOS, o binário não possui a extensão `.exe` e o shell nativo é `bash` ou `sh`.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-ARCH-PORT-01`
  * **Alvo:** `internal/plangen/render.go` e geradores de rota
  * **Ação:** Utilizar placeholders dinâmicos `${BIN_EXT}` e `${SHELL_RUNTIME}` resolvidos em tempo de despacho pelo executor Orchestrix conforme o `runtime.GOOS`.

---

### 4.2. Ponto Cego na Verificação de Formatação (`gofmt`)
* **Caso Apresentado:** O comando `gofmt -l .` sempre retorna Exit Code 0, mesmo quando encontra arquivos mal formatados, impedindo seu uso como gate de `exit_code == 0` no DAG.
* **Código Real:** [AGENTS.md (Regra 8)](file:///C:/Users/trcnologia/Desktop/LegoLoop/AGENTS.md#L48-L53).
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** O comportamento padrão do utilitário `gofmt` é imprimir no `stdout` a lista de arquivos divergentes, mas terminar com status `0`. O próprio `AGENTS.md` (Regra 8) ressalta a necessidade de inspeção manual (`gofmt -l . com saída vazia (não falha por exit code: confira a saída)`).
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-TOOL-GOFMT-01`
  * **Alvo:** `tools/checkfmt/main.go`
  * **Ação:** Criar um micro-script em Go (`tools/checkfmt`) que executa `gofmt -l .` e, se qualquer nome de arquivo for emitido na saída padrão, finaliza com `os.Exit(1)`. Adicionar `go run ./tools/checkfmt` aos gates automatizados dos tickets.

---

### 4.3. Detecção de Corrida de Dados (`-race`) Desativada Localmente
* **Caso Apresentado:** O ambiente Windows local opera sem CGO habilitado por padrão; logo, `go test -race` não roda localmente, postergando a detecção de concorrência para o CI Linux.
* **Código Real:** [AGENTS.md (Regra 8)](file:///C:/Users/trcnologia/Desktop/LegoLoop/AGENTS.md#L50) e [docs/reports/2026-11-04_REPORT_FABRICAR_V2.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-04_REPORT_FABRICAR_V2.md#L91).
* **Veredito:** **PROCEDE (Fato Técnico e Decisão Consciente de Governança).**
* **Prova Comprobatória:** Conforme documentado no `AGENTS.md`: `"-race só roda com CGO: obrigatório no CI Linux"`. No Windows sem GCC instalado (MinGW/TDM-GCC), `go test -race` falha na compilação (`CGO_ENABLED=1 requires a C compiler`). Essa limitação do sistema hospedeiro é tratada tornando a checagem obrigatória no pipeline de Integração Contínua (CI) em Linux.
* **Micro-Ticket Cirúrgico:** N/A (Política de governança já formalizada).

---

### 4.4. Fragilidade de Testes de Symlink no Windows
* **Caso Apresentado:** No Windows, `os.Symlink` exige privilégios de Administrador ou Modo Desenvolvedor. A SPEC admite que o teste de proteção contra loops de symlink usa `t.Skip` na máquina local, deixando a proteção sem validação local.
* **Código Real:** [internal/assets/skills_test.go](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/assets/skills_test.go#L101-L103).
```go
	if err := os.Symlink(target, filepath.Join(root, "ligada")); err != nil {
		t.Skipf("sem permissão para criar symlink neste ambiente (Windows sem admin): %v", err)
	}
```
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** A linha 102 de `internal/assets/skills_test.go` salta o teste com `t.Skipf("sem permissão para criar symlink neste ambiente (Windows sem admin): %v", err)` caso o Windows recuse a criação do symlink sem elevação de privilégios.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-TEST-SYMLINK-01`
  * **Alvo:** `internal/assets/skills_test.go`
  * **Ação:** Implementar uma interface de abstração do sistema de arquivos (`FileSystem` ou leitor mock) que permita simular entradas de diretório com `os.ModeSymlink` e grafos circulares em testes unitários puramente em memória, garantindo que o algoritmo anti-loop seja testado em qualquer ambiente Windows sem depender de privilégios de Administrador.

---

### 4.5. Falta de Fallback de Áudio na Validação do `ai-motion-graphics`
* **Caso Apresentado:** O livro `ai-motion-graphics` afirma que a ausência de áudio tem fallback gracioso (sai mudo sem erro), enquanto o livro `harness-engineering` exige áudio AAC estrito a -14 LUFS, havendo contradição de contrato.
* **Código Real:** [docs/books/ai-motion-graphics/partes/06-sync-audio-librosa.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/books/ai-motion-graphics/partes/06-sync-audio-librosa.md#L13-L15) vs. [docs/books/harness-engineering/2026-11-09_ENGENHARIA-DE-HARNESS.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/books/harness-engineering/2026-11-09_ENGENHARIA-DE-HARNESS.md#L253-L255).
* **Veredito:** **PROCEDE TOTALMENTE.**
* **Prova Comprobatória:** 
  * Em `docs/books/harness-engineering/...`: *"O teste afere se o arquivo tem 1080x1920 e se o áudio está em -14 LUFS. Ou atende ou não atende."* (regra binária intransigente).
  * Em `docs/books/ai-motion-graphics/...`: *"Fallback Gracioso: Áudio ausente → sem sync, sem erro."* e *"MP4 sem áudio | Fornecer AudioPath no Request"*.
* **Micro-Ticket Cirúrgico:**
  * **Ticket ID:** `TICK-DOC-AUDIO-01`
  * **Alvo:** `docs/books/ai-motion-graphics/partes/06-sync-audio-librosa.md` e pipeline de render
  * **Ação:** Formalizar o contrato canônico: no modo autônomo de produção broadcast (ex.: `cartoon-shorts`, `video-master`), o áudio a -14 LUFS é obrigatório (gate falha se ausente); no modo de overlay/motion gráfico avulso para timeline de terceiros (ex.: overlays transparentes no DaVinci), o canal de áudio é opcional e gera faixa muda.

---

## Tabela Consolidada de Tickets de Remediação

| Ticket | Prioridade | Pacote / Alvo | Ação de Correção |
| :--- | :---: | :--- | :--- |
| `TICK-FIX-VREA-01` | **P0** | `internal/videorea`, `internal/cli` | Implementar crawler real com `yt-dlp` em `runIngest`. |
| `TICK-FIX-VREA-02` | **P0** | `internal/videorea/types.go` | Adicionar `ValidateStrict()` no `ForensicsDossier` contra achados vazios. |
| `TICK-FIX-VREA-03` | **P0** | `internal/videorea/bundle.go` | Implementar *Path Jail* estrito com sanitização de caminho. |
| `TICK-DEBT-FUNCLEN-01` | **P1** | `internal/studio/templates.go` | Extrair os 454 linhas de templates para arquivos JSON embarcados via `embed.FS`. |
| `TICK-FIX-IMP-01` | **P2** | `docs/manuals/` | Corrigir typo `impeccabe-finish-reviewer` no manual do Impeccable. |
| `TICK-FIX-IMP-02` | **P2** | `internal/studio`, `loops/` | Eliminar código Go duplicado da esteira Impeccable, carregando do JSON canônico. |
| `TICK-TOOL-GOFMT-01` | **P3** | `tools/checkfmt/main.go` | Criar ferramenta wrapper de `gofmt` retornando exit code 1 se houver diffs. |
| `TICK-TEST-SYMLINK-01` | **P3** | `internal/assets/` | Criar teste unitário em memória para proteção contra loop de symlink. |
| `TICK-ARCH-PORT-01` | **P4** | `internal/plangen/` | Parametrizar `.exe` e `powershell` nas rotas para portabilidade POSIX. |
| `TICK-DOC-AUDIO-01` | **P4** | `docs/books/ai-motion-graphics/` | Harmonizar a especificação de áudio (-14 LUFS vs overlay silencioso). |
