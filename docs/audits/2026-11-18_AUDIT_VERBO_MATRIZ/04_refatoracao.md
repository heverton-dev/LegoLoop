# 04 — Sugestão de Refatoração Arquitetural

> **Alvo:** cmd/legoloop, internal/cli e governança do verbo Matriz / Matriciar
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Padrão:** Clean Architecture / Injeção de Dependências / Funções Curtas (<=35 linhas)

---

## 1. Princípios Norteadores da Refatoração

1. **Simetria dos 7 Verbos de Ação:** Todas as ações fundamentais de ciclo de vida do LegoLoop devem ser expressas como verbos no infinitivo (`fabricar`, `matriciar`, `listar`, `refinar`, `remover`, `compor`, `fechar`).
2. **Separação Rígida entre Entidade e Ação:** O substantivo canônico `Matriz` refere-se exclusivamente à estrutura de dados, tabela de reputação e pacote de domínio (`internal/matriz`). O verbo `matriciar` refere-se ao ato operacional de inspecionar, avaliar e catalogar peças na dita matriz.
3. **Imutabilidade e Retrocompatibilidade Sem Quebra:** A introdução do novo verbo de comando não deve quebrar nenhuma invocação legada, alias existente (`matrix`) ou automação em produção.

---

## 2. Refatoração Proposta para internal/cli/cli.go

Na tabela de despacho de comandos de `internal/cli/cli.go`:

```go
// Tabela canônica de comandos atualizada
"fabricar": cmdFabricar,
"matriciar": cmdMatriz, // Verbo canônico dos 7 Verbos de Ação
"matriz":    cmdMatriz, // Alias de retrocompatibilidade
"matrix":    cmdMatrix, // Alias em inglês de retrocompatibilidade
"listar":    cmdListar,
```

Essa abordagem garante que:
- O novo comando `legoloop matriciar` seja despachado instantaneamente com zero overhead de conversão.
- Chamadas existentes a `legoloop matriz` permaneçam 100% funcionais.
- A regra de `funclen <= 35 linhas` seja estritamente preservada na função `Run`.

---

## 3. Preservação Arquitetural do Núcleo de Domínio

O pacote de domínio `internal/matriz/` não deve ser renomeado:
- O pacote modela a entidade de negócio `Matriz` (`internal/matriz/matriz.go`).
- Renomear o pacote Go geraria um blast radius desnecessário de importações em dezenas de arquivos (`internal/cli`, `internal/fabricar`, etc.) sem ganho arquitetural.
- Apenas a interface de comando do usuário e a taxonomia normativa de governança adotam a nomenclatura de ação `matriciar`.
