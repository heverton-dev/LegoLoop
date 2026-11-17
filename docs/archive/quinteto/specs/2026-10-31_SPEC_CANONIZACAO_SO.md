# 📋 SPEC: Canonização e Desduplicação Global do SO via LegoLoop

> **Referência:** `docs/intents/INTENT_CANONIZACAO_SO.md`  
> **Status:** Especificação com Critérios Binários (Pass/Fail)  
> **Metodologia:** SDD (Spec-Driven Development) & Governança de SO  

---

## 1. Contrato Técnico do Motor de Canonização

O comando `legoloop canonize` opera sob as seguintes regras no Windows:

1. **Varredura Paralela (Worker Pool):**
   * Processamento multithread com pool de 32 goroutines limitadas para varredura de árvores sem esgotamento de handles de arquivos.
2. **Cálculo de Hash Criptográfico (SHA-256):**
   * Cada arquivo é lido e indexado por seu digest SHA-256.
   * Arquivos com o mesmo hash compartilham uma única cópia física real dentro de `~/.aidd/canonical/<sha256[:2]>/<sha256>`.
3. **Mecanismo de Ligação no Windows:**
   * Utilização de hard links / links simbólicos (`mklink`).
   * Arquivos idênticos apontam para a mesma entrada física no disco (economizando espaço real de cluster NTFS).
4. **Proteção de Integridade:**
   * Diretórios `.git` e working trees ativas são ignorados por padrão (`--include-git=false`), garantindo que nenhum link simbólico quebre controle de versão.

---

## 2. Critérios Binários de Aceitação (Pass/Fail)

### Bloco A: Varredura e Dry-Run (Fase 1)
- [ ] **CA-CANON-01:** O comando `legoloop canonize -root C:\Users\trcnologia\.gemini -root C:\Users\trcnologia\.claude -root C:\Users\trcnologia\.agents -vault C:\Users\trcnologia\.aidd\canonical -dry-run -json` executa com sucesso (`exit 0`).
- [ ] **CA-CANON-02:** O JSON de saída reporta com precisão:
  * `total_files`: total de arquivos escaneados.
  * `unique_hashes`: quantidade de arquivos com conteúdo único.
  * `duplicate_files`: quantidade de arquivos redundantes passíveis de canonização.
  * `bytes_reclaimable`: espaço em disco recuperável.
- [ ] **CA-CANON-03:** O relatório estruturado é gravado em `docs/reports/RELATORIO_DRY_RUN_CANONIZACAO.json`.

### Bloco B: Aplicação Segura (Fase 2 — Sob Consentimento)
- [ ] **CA-CANON-04:** O cofre `C:\Users\trcnologia\.aidd\canonical` é inicializado e possui estrutura íntegra de shards SHA-256.
- [ ] **CA-CANON-05:** Após execução com `-apply`, os arquivos nos harnesses mantêm integridade de leitura e permissões de execução.
- [ ] **CA-CANON-06:** O comando `go test ./...` do LegoLoop permanece 100% verde após a canonização.

---

## 3. Protocolo de Transição

Com a aprovação desta especificação, avançamos para o fatiamento no **`docs/plans/PLAN_CANONIZACAO_SO.md`**.
