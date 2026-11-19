# 02 — Lista de Lacunas Técnicas do Alvo

> **Alvo:** cmd/legoloop, internal/cli e governança do verbo Matriz / Matriciar
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Classificação:** Risco e Oportunidades de Rigor Arquitetural

---

## 1. Inventário de Lacunas Identificadas

A auditoria atômica identificou quatro lacunas técnicas e de governança na transição do verbo `matriz` para `matriciar`:

### Lacuna 01: Assimetria Gramatical e Taxonômica no Ciclo dos 7 Verbos
- **Evidência:** `AGENTS.md:124-132`.
- **Descrição:** O ecossistema define os 7 Verbos de Processo: Fabricar, Matriz, Listar, Refinar, Remover, Compor e Fechar. O termo `matriz` é um substantivo concreto/abstrato, violando o princípio de que todo comando de ação de ciclo de vida deve ser um verbo no infinitivo. O verbo canônico correspondente à ação é `matriciar`.
- **Impacto:** Confusão semântica entre o artefato/tabela (`a Matriz`) e a ação executável de catalogar/avaliar peças (`matriciar`).
- **Severidade:** Alta (padronização conceitual e governança viva).

### Lacuna 02: Omissão do Verbo no Dicionário Canônico ASD-STE100
- **Evidência:** `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json:71`.
- **Descrição:** O dicionário canônico ASD-STE100 lista `Matriz` estritamente na categoria de substantivo (tabela de aptidão e reputação), proibindo os sinônimos `tabela de aptidão` e `grade`. Não há entrada formal para o verbo de ação correspondente (`matriciar`).
- **Impacto:** O linter taxonômico (`tools/taxlinter`) não possui a regra para fiscalizar o uso verbal incorreto de `matriz` quando este é empregado como verbo em vez do verbo de ação canônico `matriciar`.
- **Severidade:** Média (conformidade com a Lei #12).

### Lacuna 03: Risco de Quebra Retrocompatível de Chamadas de Linha de Comando
- **Evidência:** `internal/cli/cli.go:110`.
- **Descrição:** A tabela de despacho registra `"matriz": cmdMatriz` e `"matrix": cmdMatrix`. Uma substituição abrupta que remova a chave `"matriz"` quebraria scripts de automação, documentações de processos legados e hábitos de digitação estabelecidos no terminal.
- **Impacto:** Quebra em tempo de execução de invocações como `legoloop matriz ver` ou `legoloop matriz status`.
- **Severidade:** Crítica (se implementada sem estratégia de alias retrocompatível).

### Lacuna 04: Discrepância nos Nomes de Skills e Comandos de Chat
- **Evidência:** `skills/legoloop-matriz/SKILL.md` e referências a `/legoloop-matriz`.
- **Descrição:** As skills canônicas embutidas distribuídas para os harnesses utilizam o prefixo `legoloop-matriz`.
- **Impacto:** Necessidade de prover alias transparente para a skill e slash command, garantindo que tanto `/legoloop-matriciar` quanto `/legoloop-matriz` convirjam para o mesmo motor sem quebras.
- **Severidade:** Baixa (ergonomia do usuário).
