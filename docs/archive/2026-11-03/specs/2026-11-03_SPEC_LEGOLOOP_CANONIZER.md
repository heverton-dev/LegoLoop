# 📋 SPEC: Fábrica Universal de Processos Canônicos (legoloop-canonizer)

> **Status:** Proposta para Revisão (UNDER_REVIEW)  
> **Data:** 2026-11-03  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Everton (Fusão Addy Osmani + Matt Pocock + AIDD-Spec)  
> **Referência da Intenção:** [docs/intents/2026-11-03_INTENT_LEGOLOOP_CANONIZER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-03_INTENT_LEGOLOOP_CANONIZER.md)  
> **Alvo Primário:** `internal/canonizer/`, `internal/cli/`, `internal/studio/`, `cmd/legoloop/`  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)

> [!IMPORTANT]
> A IA lista abaixo todas as premissas técnicas adotadas para evitar decisões silenciosas:
1. **Normalização Inviolável de Prefixo:** Toda peça, esteira ou processo fabricado recebe obrigatoriamente o prefixo `legoloop-` e seu slug é estritamente formatado em `kebab-case` minúsculo e sem caracteres especiais. Nenhuma opção de configuração pode desativar este prefixo.
2. **Scanner Multi-Drive Agnóstico a SO:** A curadoria de peças não se restringe à unidade principal do sistema. Em ambiente Windows, o motor detecta e varre todas as unidades lógicas ativas (`C:\`, `D:\`, `E:\`, `F:\`, etc.), incluindo discos rígidos secundários, partições e unidades USB externas. Em Linux e macOS, varre a raiz `/` e todos os pontos de montagem externos (`/Volumes/*`, `/media/*`, `/mnt/*`).
3. **Interoperabilidade Tripartida Isomórfica:** Todo processo materializado produz simultaneamente três artefatos canônicos acoplados:
   * **Chat / Harness:** Manifesto canônico `SKILL.md` depositado na pasta de skills do usuário para invocação via linguagem natural ou `/legoloop-<slug>`;
   * **LegoLoop Studio Canvas:** Arquivo declarativo `<slug>.loop.json` renderizável visualmente no Canvas e indexado no catálogo;
   * **Linha de Comando (CLI):** Grafo formal executável `ROUTE_<slug>.json` compatível com o executor determinístico `legoloop run`.
4. **Isolamento de Contexto (Orchestrix):** Cada etapa do processo gerado declara sua allowlist estrita de ferramentas (`runtime` e `mcps_allowed`). Nenhuma ferramenta fora da allowlist é injetada no nó.
5. **Zero Dependências Pesadas:** Toda a implementação no Go utiliza tipos nativos e a biblioteca padrão (`os`, `path/filepath`, `encoding/json`, `strings`, `regexp`), mantendo a compilação instantânea e sem bloat.

---

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)

* **Sempre Fazer:**
  * Desenvolver via TDD estrito (teste unitário falhando comprovadamente antes de qualquer linha de produção).
  * Manter funções pequenas, enxutas e de responsabilidade única (< 35 linhas).
  * Validar e sanitizar nomes de processos com expressões regulares determinísticas (`^legoloop-[a-z0-9-]+$`).
  * Tratar caminhos e separadores de diretório de forma compatível com Windows e POSIX (`filepath.ToSlash` / `filepath.Join`).
* **Perguntar Antes:**
  * Adicionar qualquer nova dependência externa no `go.mod`.
  * Modificar a estrutura de diretórios do cofre canônico `~/.aidd/canonical/`.
  * Alterar o schema JSON de `pipeline.dag.json` ou `ROUTE.json`.
* **Nunca Fazer:**
  * Permitir que um processo seja salvo sem o prefixo `legoloop-`.
  * Criar funções vazias com `TODO`, `pass` ou mocks falsos (Lei #1: Zero Stubs).
  * Usar diálogos nativos como `alert()`, `confirm()` ou emojis no frontend do Studio.
  * Ignorar unidades montadas ou travar em discos lentos/desconectados durante o scanner.

---

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)

A Fábrica de Processos (`legoloop-canonizer`) é um motor unificado em Go que recebe uma solicitação (via CLI, Chat ou Studio) e a conduz pelas 5 camadas estruturais:

```text
┌────────────────────────────────────────────────────────────────────────┐
│             ARQUITETURA DA FÁBRICA CANÔNICA (legoloop-canonizer)        │
├────────────────────────────────────────────────────────────────────────┤
│ ENTRADAS TRIPARTIDAS:                                                  │
│   [CLI: legoloop canonize] │ [Chat: /legoloop-*] │ [Studio: Modal Web] │
└────────────────────────────────────┬───────────────────────────────────┘
                                     │
                                     ▼
┌────────────────────────────────────────────────────────────────────────┐
│ NÚCLEO GO: package canonizer (internal/canonizer/)                     │
│                                                                        │
│ 1. NameSanitizer: Normaliza para legoloop-<kebab-case>                 │
│ 2. Clarifier: Modela a intenção socrática e gera INTENT.md             │
│ 3. MultiDriveScanner: Localiza peças em C:\, D:\, E:\, USBs e mounts   │
│ 4. ArchitectMerger: Seleciona e funde as melhores peças locais         │
│ 5. TriadGenerator: Produz SPEC.md, PLAN.md e ROUTE.json                │
│ 6. TemplateGenerator: Gera o gabarito contratual do artefato de saída  │
└────────────────────────────────────┬───────────────────────────────────┘
                                     │
                                     ▼
┌────────────────────────────────────────────────────────────────────────┐
│ MATERIALIZAÇÃO TRIPARTIDA INTEROPERÁVEL:                               │
│   ├── [CLI]:    docs/routes/YYYY-MM-DD_ROUTE_legoloop-<slug>.json      │
│   ├── [Chat]:   ~/.aidd/canonical/skills/legoloop-<slug>/SKILL.md      │
│   └── [Studio]: ~/.aidd/canonical/loops/legoloop-<slug>.loop.json      │
└────────────────────────────────────────────────────────────────────────┘
```

### Costuras de Teste (Test Seams — Matt Pocock)

* **Costura Primária (Highest Seam do Motor):**  
  `canonizer.Canonize(ctx context.Context, req canonizer.Request) (*canonizer.Result, error)`  
  *Justificativa:* Testa a esteira inteira de ponta a ponta através da sua interface pública mais alta. O teste fornece os parâmetros brutos e valida a geração dos 3 artefatos sem se acoplar às funções internas de cada camada.
* **Costura Secundária (Sanitização e Prefixo):**  
  `canonizer.NormalizeProcessName(rawName string) (string, error)`  
  *Justificativa:* Garante que nenhuma string suja gere processos fora do padrão `legoloop-<kebab-case>`.
* **Costura Terciária (Scanner Multi-Drive):**  
  `canonizer.DiscoverDrives() []string` e `canonizer.ScanStorageLocations(roots []string) ([]canonizer.DiscoveredPiece, error)`  
  *Justificativa:* Permite injetar raízes simuladas em testes de unidade e verificar a descoberta de peças em múltiplos volumes sem depender de pendrives físicos espetados no ambiente de CI.

---

## 4. Contratos de Dados & Schemas Formais

```go
package canonizer

import "time"

// Request representa o payload de entrada para a fábrica de processos.
type Request struct {
	RawName     string   `json:"raw_name"`     // Ex: "Processo do Joãozinho"
	Goal        string   `json:"goal"`         // Objetivo em linguagem natural
	Category    string   `json:"category"`     // Ex: "build", "reverse", "web", "mobile"
	TargetDrives []string `json:"target_drives"` // Vazio = todos os drives do SO
	ModelChoice string   `json:"model_choice"` // Modelo cognitivo para a entrevista
	Interactive bool     `json:"interactive"`  // true = solicita prompts; false = autônomo (--yes)
}

// Result representa a entrega canônica da fábrica de processos.
type Result struct {
	CanonicalName string            `json:"canonical_name"` // Ex: "legoloop-processo-do-joaozinho"
	Slug          string            `json:"slug"`           // "processo-do-joaozinho"
	CreatedAt     time.Time         `json:"created_at"`
	Artifacts     CanonicalArtifacts `json:"artifacts"`
	CuratedPieces []DiscoveredPiece  `json:"curated_pieces"`
	Success       bool              `json:"success"`
}

// CanonicalArtifacts consolida os caminhos gerados para os 3 meios de execução.
type CanonicalArtifacts struct {
	IntentPath   string `json:"intent_path"`   // docs/intents/YYYY-MM-DD_INTENT_...
	SpecPath     string `json:"spec_path"`     // docs/specs/YYYY-MM-DD_SPEC_...
	PlanPath     string `json:"plan_path"`     // docs/plans/YYYY-MM-DD_PLAN_...
	RoutePath    string `json:"route_path"`    // docs/routes/YYYY-MM-DD_ROUTE_...
	SkillPath    string `json:"skill_path"`    // ~/.aidd/canonical/skills/legoloop-.../SKILL.md
	LoopPath     string `json:"loop_path"`     // ~/.aidd/canonical/loops/legoloop-....loop.json
	TemplatePath string `json:"template_path"` // docs/templates/TEMPLATE_...
}

// DiscoveredPiece representa uma peça localizada em qualquer drive montado.
type DiscoveredPiece struct {
	Name     string  `json:"name"`
	Kind     string  `json:"kind"`     // skill, script, hook, gate, loop, mcp
	Drive    string  `json:"drive"`    // Ex: "C:", "D:", "/Volumes/Pendrive"
	Path     string  `json:"path"`
	Score    float64 `json:"score"`    // 0.00 a 1.00
	Tier     string  `json:"tier"`     // ouro, prata, bronze, toxica
}
```

---

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — Núcleo Go da Fábrica Canônica (`internal/canonizer`)

* [ ] `AC-E01.1`: **Sanitização Determinística e Prefixo Imutável**
  * **Dado que:** O usuário fornece nomes em formatos diversos (ex: `"Processo do Joãozinho"`, `"legoloop_teste"`, `"MEU PROCESSO"`, `"legoloop-ja-tem"`).
  * **Quando:** A função `NormalizeProcessName` for executada.
  * **Então:** Todas as saídas devem conter o prefixo obrigatório `legoloop-`, formato `kebab-case`, sem acentos ou caracteres especiais (ex: `legoloop-processo-do-joaozinho`, `legoloop-teste`, `legoloop-meu-processo`, `legoloop-ja-tem`).

* [ ] `AC-E01.2`: **Descoberta Multi-Drive e Agnóstica a SO**
  * **Dado que:** O sistema está rodando em Windows, Linux ou macOS com unidades locais e volumes externos montados.
  * **Quando:** `DiscoverDrives()` for invocado.
  * **Então:** Em Windows, deve retornar todas as unidades de letras montadas (`C:\`, `D:\`, etc.); em Linux/macOS deve retornar `/` e diretórios sob `/Volumes`, `/media` e `/mnt`. Se uma unidade estiver inacessível, deve registrar log silencioso e continuar sem travar.

* [ ] `AC-E01.3`: **Geração do Pacote Tripartido Interoperável**
  * **Dado que:** Uma solicitação válida `canonizer.Request` é submetida em modo autônomo (`--yes`).
  * **Quando:** `canonizer.Canonize()` for executado.
  * **Então:** Deve gerar no disco os arquivos:
    1. `SKILL.md` com YAML frontmatter válido e comando `commands: ["/legoloop-<slug>"]`;
    2. `.loop.json` com schema DAG acíclico compatível com o Studio Canvas;
    3. `ROUTE_<slug>.json` com isolamento estéril, inputs, outputs e gates booleanos.

---

### Épico 02 — Subcomando CLI (`internal/cli/canonize.go`)

* [ ] `AC-E02.1`: **Invocação por Linha de Comando e Modos de Operação**
  * **Dado que:** O usuário digita no terminal `legoloop canonize "extrator de audio" --yes`.
  * **Quando:** O subcomando `canonize` for processado.
  * **Então:** Deve exibir no stdout o progresso das 5 camadas, criar a pasta `legoloop-extrator-de-audio` e retornar `exit code 0` com os caminhos dos 3 artefatos gerados.

* [ ] `AC-E02.2`: **Modo Interativo com Perguntas no Terminal**
  * **Dado que:** O usuário executa `legoloop canonize "novo fluxo"` sem a flag `--yes`.
  * **Quando:** A camada de Clarificação iniciar.
  * **Então:** Deve apresentar opções numeradas via terminal (readline), afunilar a intenção e confirmar os drives antes da gravação final.

---

### Épico 03 — Interface Web & Canvas no Studio (`internal/studio`)

* [ ] `AC-E03.1`: **Modal de Canonização com Seletor de Cérebro**
  * **Dado que:** O usuário clica no botão `+ Novo Processo` na barra superior do Studio (`http://127.0.0.1:7878`).
  * **Quando:** O modal abrir.
  * **Então:** Deve exibir campo para o nome do processo, dropdown com os modelos disponíveis no inventário, e botão "Iniciar Entrevista".

* [ ] `AC-E03.2`: **Plotagem Automática e Atualização do Catálogo**
  * **Dado que:** A materialização do processo foi finalizada via API `/api/canonize`.
  * **Quando:** O backend responder com o `canonizer.Result`.
  * **Então:** O Canvas do Studio deve plotar automaticamente os nós gerados, o painel Inspetor deve abrir com os dados da Tríade e a nova peça deve aparecer imediatamente na barra lateral sob a categoria de `Loops`.

---

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)

### Comandos de Verificação Mecânica
* **Testes TDD Estritos:** `go test -v ./internal/canonizer/... ./internal/cli/... ./internal/studio/...`
* **Regressão Global:** `go test ./...`
* **Compilação do Binário:** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Instalação Global no PATH:** `go install ./cmd/legoloop`

### Critério de Done Global
A especificação técnica é considerada **100% CONCLUÍDA** quando:
- [ ] 100% dos critérios binários `AC-E01.1` até `AC-E03.2` estiverem marcados como atendidos (`[x]`).
- [ ] O comando de **Test** passar com exit code 0 em todos os pacotes.
- [ ] O subcomando `legoloop canonize --help` e `legoloop canonize "teste-smoke" --yes` funcionarem no terminal com sucesso.
- [ ] O binário `bin/legoloop.exe` estiver compilado e atualizado no PATH.
- [ ] Zero arquivos temporários ou órfãos residuais.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**  
*As premissas técnicas, fronteiras (boundaries), arquitetura de costura de teste, contratos Go e critérios binários BDD atendem à sua expectativa?*  
*Somente após a aprovação explícita do usuário, iniciaremos o PLAN (`docs/plans/2026-11-03_PLAN_LEGOLOOP_CANONIZER.md`) e a ROUTE (`docs/routes/2026-11-03_ROUTE_LEGOLOOP_CANONIZER.json`).*
