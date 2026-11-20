# 📘 MANUAL OPERACIONAL: FUNÇÃO CANONIZAR NO LEGOLOOP

> **Data:** 2026-11-20  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Processo de industrialização e canonização universal de peças  
> **Fontes de verdade:** `processos/legoloop-canonizar/` · `skills/legoloop-canonizar/SKILL.md` · `skills/legoloop-piece-upgrader/SKILL.md` · `skills/legoloop-fechar/SKILL.md`

---

## 1. Visão Geral

O processo **`legoloop-canonizar`** é a esteira canônica de industrialização de peças do LegoLoop. Ele une em uma cadeia determinística:
1. **Elevação ao Padrão Ouro (`legoloop-piece-upgrader`):** diagnóstico mecânico de gaps, refatoração orientada a TDD e alcance compulsório de Score >= 0.85.
2. **Fechamento e Canonização Universal (`legoloop-fechar`):** execução determinística da esteira de fechamento com manual operacional, livro didático LegoPress, diagrama visual SVG, skill universal nos harnesses e validação de DAG.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop-canonizar <alvo>`
* `/canonizar <alvo>`

### 2.2 Pela CLI Mecânica
```bash
# Simulação sem efeitos colaterais
legoloop run processos/legoloop-canonizar/pipeline.dag.json --dry-run

# Execução real do processo de canonização
legoloop run processos/legoloop-canonizar/pipeline.dag.json
```

---

## 3. As Duas Camadas do Grafo

1. **Camada 1 (`n01_upgrade`):** consome o contrato da peça em `TEMPLATE_PIECE_UNIT.json`, avalia o score com `legoloop upgrade <alvo> --diagnose` e emite `out/upgrade_result.json`.
2. **Camada 2 (`n02_fechar`):** recebe o resultado aprovado e aciona o sub-processo `legoloop-fechar`, materializando toda a governança canônica no cofre.
