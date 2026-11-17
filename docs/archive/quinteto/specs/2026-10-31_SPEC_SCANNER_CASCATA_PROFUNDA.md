# 🛰️ Espec — Motor de Varredura Cascata Hiper-Profunda do SO

> **Documento:** `docs/specs/SPEC_SCANNER_CASCATA_PROFUNDA.md`  
> **Status:** Proposta de Implementação · 31/10/2026  
> **Conformidade:** `AGENTS.md` (Zero Stubs, TDD estrito, PT-BR, Funclen ≤ 35 linhas, Segurança)

---

## 1. Objetivo & Visão Geral

Substituir os coletores pontuais de caminhos fixos do pacote `internal/assets` por um **Motor de Varredura Cascata Hiper-Profunda** que descobre, indexa e monitora de ponta a ponta todos os ativos de desenvolvimento presentes no sistema operacional do usuário, sem exigir configuração manual.

A varredura opera estritamente nas **5 Camadas Canônicas** do ecossistema:
1. **Sistema (`system` / `builtin`):** Executáveis no `PATH`, modelos locais (Ollama), gates nativos.
2. **Global (`global`):** Configurações globais em `~` (`.claude`, `.gemini`, `.aidd`, `.cursor`, `.config`, `.omp`, `.mimo`, `.opencode`).
3. **Scaffolds (`scaffold`):** Templates de esteiras, kits de início e boilerplates em `~/.aidd/scaffolds/**`, `~/.aidd/canonical/**`.
4. **Projetos (`project`):** Auto-descoberta concorrente em árvores de desenvolvimento do usuário (`Desktop`, `Projects`, `Workspace`, `repos`, `src`, `dev`, `code`, `Documents`).
5. **Local (`local`):** Configurações de override do desenvolvedor que não vão para o controle de versão (`settings.local.json`).

---

## 2. Requisitos & Critérios Binários de Aceite

### 2.1 Descoberta e Profundidade do SO
- **CRIT-SCAN-1 (Auto-descoberta de Projetos):** O scanner deve descobrir automaticamente repositórios e pastas de projeto a partir das raízes padrão do usuário (`~/Desktop`, `~/Projects`, `~/Workspace`, `~/src`, `~/repos`, `~/code`, `~/dev`, `~/Documents`) até profundidade 4, identificando marcadores canônicos (`.git`, `.claude/`, `.gemini/`, `package.json`, `go.mod`, `pyproject.toml`, `Cargo.toml`).
- **CRIT-SCAN-2 (Poda Segura):** Pastas de dependências e caches pesados (`node_modules`, `.git/objects`, `vendor`, `.venv`, `__pycache__`, `dist`, `build`, `.cache`) **nunca** devem ser percorridas recursivamente.
- **CRIT-SCAN-3 (Tratamento de Links e Cíclicos):** Symlinks e junctions cíclicos nunca devem travar a varredura (resolução atômica por conjunto de inodes/caminhos canônicos reais já visitados).

### 2.2 Cascata Completa de Peças (Todas as Categorias)
- **CRIT-PECA-1 (Skills em Cascata):** Indexa `SKILL.md` nos 3 escopos: global (`~/.claude/skills`, `~/.gemini/antigravity-cli/skills`, `~/.aidd/canonical/skills`), scaffold (`~/.aidd/scaffolds/**/skills/**`) e projeto (`<projeto>/**/skills/**` ou `<projeto>/.claude/skills/**`).
- **CRIT-PECA-2 (MCPs em Cascata):** Indexa servidores MCP de arquivos globais (`~/.claude.json`, `~/.gemini.json`, `~/.cursor/mcp.json`, `~/.config/opencode/`) e de arquivos locais em cada projeto (`.mcp.json`, `.claude/settings.json`, `.cursor/mcp.json`), atribuindo os escopos `global` ou `project`.
- **CRIT-PECA-3 (Hooks em Cascata):** Coleta hooks de:
  - Global: `~/.claude/settings.json`, `~/.gemini/settings.json`.
  - Scaffold: `~/.aidd/scaffolds/**/settings.json`.
  - Projeto: `<projeto>/.claude/settings.json`.
  - Local: `<projeto>/.claude/settings.local.json`.
- **CRIT-PECA-4 (Scripts em Cascata):** Indexa executáveis (`.ps1`, `.sh`, `.py`, `.js`, `.bat`) de:
  - Global: `~/.aidd/scripts/`, `~/.claude/scripts/`.
  - Scaffold: `~/.aidd/scaffolds/**/scripts/`.
  - Projeto: pastas `scripts/`, `tools/`, `bin/` de cada projeto descoberto.
- **CRIT-PECA-5 (Loops no Acervo & Projetos):** Indexa arquivos `*.loop.json` tanto do cofre canônico (`~/.aidd/canonical/loops/`, escopo `global`) quanto de pastas de scaffolds (`escopo scaffold`) e de projetos locais (`escopo project`).
- **CRIT-PECA-6 (Harnesses e Binários Expandidos):** Além do `PATH`, varre pastas conhecidas de ferramentas (`~/.cargo/bin`, `~/.local/bin`, `AppData/Local/Programs`, `~/go/bin`).

### 2.3 Desempenho e Governança de Recursos
- **CRIT-PERF-1 (Worker Pool Limitado):** A varredura de árvores deve rodar com pool de goroutines limitado (`runtime.NumCPU() * 2`), sem abrir descritores de arquivo desenfreados.
- **CRIT-PERF-2 (Cache Incremental):** Re-varreduras só re-analisam arquivos com alteração de `(mtime, size)`. Tempo de re-varredura incremental em repouso ≤ 150 ms para 30.000 arquivos monitorados.
- **CRIT-PERF-3 (Piso de Memória):** Working set do Studio em repouso deve permanecer ≤ 45 MB de baseline Go + 15 MB de índices.

---

## 3. Segurança & Sanitização
- **SEC-SCAN-1:** Segredos de ambiente (`env`, `headers`, `Bearer`, senhas) em hooks e MCPs continuam bloqueados e redigidos (`redact`) antes de entrar na memória pública do inventário.
- **SEC-SCAN-2:** Comandos de hooks só residem no servidor (`json:"-"`), expondo apenas identificador opaco e descrição auditada.
