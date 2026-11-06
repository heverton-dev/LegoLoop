# 📘 MANUAL OPERACIONAL: FUNÇÃO COMPOR NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Compor processos existentes como sub-etapas modulares (sub-loops) em outros processos  
> **Fontes de verdade:** `internal/compor/` · `internal/cli/compor.go` · `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`

---

## 1. Visão Geral

O comando **`legoloop compor`** integra um processo filho existente (`legoloop-filho@1.0.0`) como uma nova sub-etapa do tipo `subloop` dentro de um processo pai (`legoloop-pai`). O DAG resultante é validado contra ciclos e os artefatos intermediários do filho são isolados em `out/<subloop>/`, garantindo confinamento hermético.

---

## 2. Comandos e Uso

### 2.1 Composição com Prévia Interativa
```bash
# Apresenta a topologia e as conexões projetadas, solicitando confirmação [s/N]
legoloop compor legoloop-pai legoloop-filho@1.0.0
```

### 2.2 Composição com Dependência Específica
```bash
# Liga a nova sub-etapa dependendo de uma etapa anterior (ex: e02)
legoloop compor legoloop-pai legoloop-filho@1.0.0 --depende e02 --yes
```

### 2.3 Saída Estruturada JSON
```bash
legoloop compor legoloop-pai legoloop-filho@1.0.0 --depende e02 --yes --json
```

---

## 3. Invariantes
- **Detecção Estrita de Ciclo (`ErrCiclo`):** Compor A em B quando B já depende de A é estritamente recusado.
- **Isolamento de Artefatos:** Artefatos gerados pelo sub-processo residem no sub-diretório de saída (`out/<filho>/`), evitando colisões de nomes de arquivos comuns (`01_construir.md`, etc.).
- **Projeção Automática:** A inserção do sub-loop atualiza de forma atômica o `.processo.json`, `.loop.json`, `.dag.json` e os documentos de governança correspondentes.
