---
name: legoloop-compor
description: Use quando o usuário quiser COMPOR um processo como sub-etapa reutilizável (sub-loop) de outro processo existente, via /legoloop-compor. NÃO use para criar processos do zero ou remover processos.
consumes: Artifact:Process
produces: Artifact:Process
artifact: ~/.aidd/canonical/loops/<nome>.processo.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-compor — Composição de Processos em Sub-Loops do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja integrar um processo filho existente (`legoloop-b@versao`) como uma etapa encadeada dentro de um processo pai (`legoloop-a`), reaproveitando o DAG e isolando seus artefatos em `out/<subloop>/`.  
**Quando NÃO usar:** para fabricar processos novos do zero, listar ou remover processos.

Esta skill conduz a composição hierárquica segura. **Toda a lógica reside no comando `legoloop compor`**: validação de ciclo (`ErrCiclo`), isolamento de sub-artefatos e atualização coerente da Tríade de governança.

* **Consumes (Entradas Obrigatórias):**
  - nome do processo pai (`<nome>`)
  - processo filho e versão (`<subprocesso@versão>`)
  - dependência opcional (`--depende <ID_etapa_anterior>`)
* **Produces (Saídas Mecânicas):**
  - processo pai atualizado contendo a nova etapa do tipo `subloop`, com elevação de versão
* **Target Files (escopo de escrita permitido):**
  - nenhum arquivo é alterado manualmente; apenas o comando `legoloop compor` grava

### Procedimento (passos determinísticos)

1. **Inspeção das Partes:**
   - Execute `legoloop listar` para verificar a existência do processo pai e do processo filho com suas respectivas versões.
2. **Prévia com Stop-Gate:**
   - Execute a composição sem `--yes`: `legoloop compor <pai> <filho@versao> [--depende <ID>]`.
   - Analise o diff e as conexões projetadas no DAG.
   - Apresente ao usuário e pare no stop-gate humano, exigindo confirmação explícita.
3. **Execução:**
   - Com aprovação concedida, aplique a composição via `--yes`: `legoloop compor <pai> <filho@versao> [--depende <ID>] --yes`.
4. **Relato:**
   - Confirme a nova versão gerada e a inserção da sub-etapa com seus artefatos devidamente isolados.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** crie nós de sub-loop manualmente nos arquivos `.dag.json` ou `.loop.json`; utilize exclusivamente `legoloop compor`.
> 2. **NUNCA** ignore erros de ciclo (`ErrCiclo`): dependências circulares entre processos são estritamente proibidas.
> 3. **NUNCA** execute a gravação definitiva sem prévia aprovação humana do diff.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop compor legoloop-pipeline legoloop-extrator@1.0.0`
- `legoloop compor legoloop-pipeline legoloop-extrator@1.0.0 --depende e01 --yes`
