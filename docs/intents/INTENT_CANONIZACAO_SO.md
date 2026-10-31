# 🎯 INTENT: Canonização e Desduplicação do Sistema Operacional

> **Data:** 31/10/2026  
> **Status:** Proposta de Intenção (Aguardando Aprovação para Spec)  
> **Comando Principal:** `legoloop canonize`  

---

## 1. Motivação e Problema Observado

Os múltiplos harnesses instalados no sistema operacional (`.gemini`, `.claude`, `.agents`, `.cursor`) acumularam milhares de cópias físicas redundantes de skills, scripts e definições.
* **Impacto negativo:** Divergência de versões, desperdício de espaço em disco e dificuldade de manter as ferramentas atualizadas de forma centralizada.

---

## 2. Intenção Primária

Executar a **canonização global** dos diretórios de IA do SO através do motor `legoloop canonize`, consolidando todos os arquivos físicos idênticos em um **Cofre Canônico Único (SSOT)** indexado por hash SHA-256 (`~/.aidd/canonical`), substituindo as cópias redundantes por **links simbólicos (`mklink`)** transparentes para os harnesses.

---

## 3. Escopo e Diretórios-Alvo

1. **Raízes de Varredura (`-root`):**
   * `C:\Users\trcnologia\.gemini`
   * `C:\Users\trcnologia\.claude`
   * `C:\Users\trcnologia\.agents`
2. **Cofre Canônico (`-vault`):**
   * `C:\Users\trcnologia\.aidd\canonical`
3. **Fases de Execução:**
   * **Fase 1 (Auditoria / Dry-Run):** Execução passiva com `-dry-run --json` para quantificar duplicatas e validar a árvore sem alterar arquivos.
   * **Fase 2 (Aplicação com Consentimento):** Execução com `-apply` somente após aprovação expressa do relatório de auditoria.

---

## 4. Fora de Escopo & Proteções Invioláveis

* **Working trees do Git:** Protegidas contra links acidentais (sem `-include-git` a menos que aprovado).
* **Segredos e Tokens:** Arquivos de credenciais (`.env`, `chaves/`, `credentials.json`) são preservados sem centralização pública.
* **Zero Risco de Quebra:** Caso um arquivo tenha permissão restrita ou seja incompatível com link, o motor pula o item e reporta no log.

---

## 5. Critérios de Sucesso (Definition of Done)

- [ ] Relatório de Dry-Run gerado em `docs/reports/RELATORIO_DRY_RUN_CANONIZACAO.json`.
- [ ] Confirmação de quantidade de arquivos únicos vs. duplicatas eliminadas.
- [ ] Spec formalizada em `docs/specs/SPEC_CANONIZACAO_SO.md`.
