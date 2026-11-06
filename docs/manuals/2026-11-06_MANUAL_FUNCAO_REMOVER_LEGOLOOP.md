# 📘 MANUAL OPERACIONAL: FUNÇÃO REMOVER NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Exclusão segura de processos, integridade referencial e limpeza de vínculos nos harnesses  
> **Fontes de verdade:** `internal/remover/` · `internal/cli/remover.go` · `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`

---

## 1. Visão Geral

O comando **`legoloop remover`** remove de forma atômica e segura um processo gerenciado (`legoloop-*`) do cofre do usuário. Ele apaga os loops, pipelines, documentos de governança associados e desfaz os links nos harnesses detectados. Se outro processo depender do processo alvo como sub-etapa, o comando bloqueia a exclusão imediatamente com erro claro (`ErrEmUso`).

---

## 2. Comandos e Uso

### 2.1 Remoção com Confirmação Interativa
```bash
# Apresenta todos os arquivos e links que serão deletados e solicita confirmação [s/N]
legoloop remover legoloop-processo-antigo
```

### 2.2 Remoção com Confirmação Explícita
```bash
legoloop remover legoloop-processo-antigo --yes
```

### 2.3 Remoção em Cofre Específico
```bash
legoloop remover legoloop-processo-antigo --vault /caminho/do/cofre --yes
```

---

## 3. Invariantes
- **Bloqueio de Dependência Ativa (`ErrEmUso`):** Se o processo for usado por outro processo como subloop, a remoção é recusada com código 1 sem alterar nenhum arquivo.
- **Rollback em Falha:** Se ocorrer erro durante a exclusão, todos os arquivos e links são mantidos íntegros.
- **Limpeza de Vínculos:** Ao remover um processo, seus symlinks/cópias em `~/.claude/skills`, `~/.gemini/skills`, etc., são removidos sem afetar outras skills.
