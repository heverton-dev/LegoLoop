# 📘 MANUAL OPERACIONAL: FUNÇÃO LISTAR NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Inspecionar e auditar os processos ativos do cofre, versões, naturezas, harnesses e lacunas  
> **Fontes de verdade:** `internal/listar/` · `internal/cli/listar.go` · `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`

---

## 1. Visão Geral

O comando **`legoloop listar`** inspeciona o cofre (`~/.aidd/canonical/loops/`) e exibe todos os processos gerenciados (`legoloop-*`). Ele detalha versão atual, natureza do processo, quais harnesses possuem a skill vinculada, dependências de uso cruzado (`Usado por`) e etapas que apresentam lacuna.

---

## 2. Comandos e Uso

### 2.1 Listagem Padrão
```bash
legoloop listar
```
Exibe tabela com:
- **PROCESSO**: nome canônico do processo (ex: `legoloop-extrator`).
- **VERSÃO**: versão semântica ativa (ex: `1.0.0`).
- **NATUREZA**: natureza declarada do processo (ex: `software`, `dados`, `documentação`).
- **HARNESSES**: quais harnesses detectados têm o link ativo (ex: `claude, gemini, agents`).
- **USADO POR**: se o processo é referenciado como sub-etapa de outro processo.

### 2.2 Listagem com Cofre Customizado
```bash
legoloop listar --vault /caminho/do/cofre
```

### 2.3 Saída Estruturada JSON
```bash
legoloop listar --json
```

### 2.4 Resolução de Lacunas
Quando um processo possui etapas sem peça adequada no momento da fabricação, o comando marca `LACUNA` e sugere os comandos prontos para preenchimento via matriz:
```bash
# Exemplo de sugestão emitida:
legoloop matriz perfil definir <peça> --caso <caso>
```

---

## 3. Invariantes
- A listagem é estritamente somente-leitura; nenhuma alteração de estado ocorre.
- Processos legados da v1 (sem `.processo.json`) aparecem com `natureza: desconhecida`.
- Esteiras no acervo sem o prefixo `legoloop-` não são consideradas processos.
