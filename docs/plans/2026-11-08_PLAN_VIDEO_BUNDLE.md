# 📋 PLAN: Implementação do Target Workspace Isolado por Vídeo

> **Status:** APPROVED  
> **Data:** 2026-11-08  
> **Autor:** Antigravity  
> **Iniciativa:** I04 — Arquitetura de Workspaces Isolados e Mídia Soberana  
> **Template:** TEMPLATE_PLAN.md v2.0.0  

---

## 1. Decomposição em Tickets Atômicos

### TICKET-01: Tipos e Testes do Target Bundle
* **Alvo:** `internal/videorea/types.go`, `internal/videorea/bundle_test.go`
* **Critério de Aceitação:** Struct `TargetBundle` definida com métodos de inicialização e cálculo de caminhos relativos/absolutos; teste unitário falhando inicialmente (Red).

### TICKET-02: Materialização da Árvore e Projeção Canônica
* **Alvo:** `internal/videorea/factory_bridge.go`
* **Critério de Aceitação:** Criação das subpastas (`scratch/`, `governance/`, `book/`, `loop/`, `skill/`), gravação dos artefatos e projeção em `loops/<slug>/` e `skills/<slug>/`; teste passando (Green).

### TICKET-03: Integração no CLI e Blindagem do `.gitignore`
* **Alvo:** `internal/cli/video_rea.go`, `.gitignore`
* **Critério de Aceitação:** Subcomando `govern` e `run` aceitam flag `--bundle` (ou usam `docs/videorea/<slug>/` por padrão); `.gitignore` ignora `docs/videorea/**/scratch/`.

### TICKET-04: Portões Mecânicos, Recompilação e Validação
* **Alvo:** `go test`, `go vet`, `funclen`, `gofmt`, `bin/legoloop.exe`
* **Critério de Aceitação:** Testes 100% verdes, saída de `gofmt -l .` vazia, zero warnings no vet, funclen <= 35 linhas, binário compilado com sucesso.
