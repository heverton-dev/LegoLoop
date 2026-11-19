# 01 — Disseção Técnica do Comando e Verbo Matriz / Matriciar

> **Alvo:** cmd/legoloop, internal/cli (matriz.go, cli.go) e internal/matriz
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Conformidade:** ASD-STE100 e Quinteto Canônico LegoLoop

---

## 1. Visão Geral da Arquitetura do Alvo

O comando e subsistema de Matriz é o mecanismo central de reputação, catalogação, corte de defasagem e aptidão de peças reutilizáveis no ecossistema LegoLoop. Ele mantém a integridade do inventário canônico (`~/.aidd/canonical/` e repositório local), avaliando histórico de execuções, vocabulário e aderência conceitual.

No entanto, no ciclo de governança denominado **Ciclo dos 7 Verbos de Processo** (`AGENTS.md:124-132`), existe uma assimetria morfológica e taxonômica: seis ações são verbos no infinitivo (`fabricar`, `listar`, `refinar`, `remover`, `compor`, `fechar`), enquanto a ação de catalogação e inspeção na matriz figura sob a forma substantiva `matriz`.

### Ponto de Entrada no Binário e CLI
- `internal/cli/cli.go:110`: A tabela de despacho mapeia `"matriz": cmdMatriz` e `"matrix": cmdMatrix`.
- `internal/cli/matriz.go:76-93`: `cmdMatriz` atua como multiplexador para os subcomandos:
  - `gerar`: varre o inventário e cataloga novas peças ou atualizações.
  - `status`: exibe a saúde, peças catalogadas, pendências e defasagens.
  - `vocab`: gerencia o vocabulário aprovado e termos pendentes.
  - `perfil`: inspeciona ou força o perfil declarativo de uma peça.
  - `ver`: visualiza a matriz com filtros de tier (`ouro`, `prata`, `bronze`) e suporte a `--json`.

---

## 2. Componentes e Estrutura em internal/cli e internal/matriz

A implementação do subsistema é altamente desacoplada e modular:

1. **internal/cli/matriz.go (367 linhas):**
   - Injeção de dependências via `MatrizDeps` (linha 40): desacopla Scanner de assets (`assets.Item`), cérebro LLM (`interview.LLMClient`), nome do modelo e relógio (`time.Time`).
   - Resolução de Cofre via `resolveVaultFlag`: garante isolamento hermético entre testes e o cofre real do usuário (`~/.aidd`).
   - Sistema de Mensageria STE: mensagens canônicas registradas via `ste.Default.Register` para saídas determinísticas.
   - Apresentação de Dados: formatação tabular e serialização JSON estruturada (`matrizVerJSON`).

2. **internal/matriz/ (Core de Domínio):**
   - `matriz.go`: estruturas centrais de dados, pontuações de aptidão, pesos de histórico e reputação.
   - `vocab.go`: armazenamento e validação de vocabulário canônico.
   - `catalogador.go`: pipeline de classificação de peças via LLM ou heurísticas declarativas.

---

## 3. Cobertura e Rastreabilidade de Testes

Os testes cobrem integralmente o comando CLI e o motor de domínio:
- `internal/cli/matriz_test.go`: validação de fluxos de `ver`, `status`, filtros de tier e flags `--json`.
- `internal/cli/dispatch_matriz_test.go`: isolamento hermético do despacho e validação de flags.
- `internal/matriz/matriz_test.go`: testes unitários de pontuação, cálculo de aptidão e corte de defasagem.

Todos os testes unitários do subsistema executam com 100% de sucesso sem vazamento de estado.
