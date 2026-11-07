# 📘 MANUAL OPERACIONAL: FUNÇÃO PIECE UPGRADER NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Elevar qualquer peça da caixa de ferramentas (Skill, Script, Hook, MCP, Gate ou Loop) ao Padrão Ouro (Score >= 0.85)  
> **Fontes de verdade:** `internal/cli/upgrade.go` · `skills/legoloop-piece-upgrader/SKILL.md` · `docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md`

---

## 1. Visão Geral

A função **`legoloop-piece-upgrader`** é o motor de alquimia do LegoLoop para elevar peças anêmicas ou em rascunho para o Padrão Ouro da indústria (Score $\ge 0.85$). Ela combina diagnóstico estático mecânico via CLI com uma esteira determinística em 6 estágios orientada por TDD e Stop-Gates humanos.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
Ative a função usando qualquer um dos comandos canônicos ou aliases:
* `/legoloop-upgrade <caminho_da_peça>`
* `/upgrade-piece <caminho_da_peça>`
* `/elevar-peca <caminho_da_peça>`

### 2.2 Pela CLI Mecânica
```bash
# Diagnóstico de gaps sem alterar o disco
legoloop upgrade "skills/minha-skill/SKILL.md" --diagnose --json

# Validação e selo de Padrão Ouro
legoloop upgrade "skills/minha-skill/SKILL.md" --apply
```

### 2.3 Pelo Looping Studio (Web)
Na aba **Palette** ou na lista de peças, clique no botão de ação **Upgrade** (ícone de seta ascendente) ao lado da peça desejada.

---

## 3. Fluxo de Elevação em 6 Estágios

1. **Trigger:** Disparo via comando ou UI.
2. **Diagnóstico Mecânico:** Execução de `legoloop upgrade <alvo> --diagnose` com avaliação de Score e gaps.
3. **Entrevista Socrática Addy Osmani:** Alinhamento de hipótese e confiança ($\ge 95\%$) ➔ Stop-Gate 01.
4. **Especificação Técnica Formal (SPEC):** Premissas, boundaries e critérios BDD ➔ Stop-Gate 02.
5. **Planejamento & Rota (PLAN + ROUTE):** Fatiamento estéril (`orchestrix_hermetic`) ➔ Stop-Gate 03.
6. **BUILD & Fechamento:** Reescrita da peça conforme templates e validação final com `legoloop upgrade <alvo> --apply`.
