---
name: legoloop-refinar
description: Use quando o usuário quiser REFINAR uma etapa de um processo existente no cofre (trocar peça, atualizar título ou casos de uso, elevando a versão), via /legoloop-refinar. NÃO use para criar processos do zero ou compor novos processos.
consumes: Artifact:Process
produces: Artifact:Process
artifact: ~/.aidd/canonical/processos/<nome>.processo.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-refinar — Refinamento de Etapas de Processos do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja alterar uma etapa específica de um processo existente (`e01`, `e02`...), trocando a peça executora, atualizando casos de uso ou título, gerando uma nova versão minor (ex.: 1.0.0 -> 1.1.0) com backup automático.  
**Quando NÃO usar:** para fabricar processos novos do zero, remover ou compor processos como sub-etapa.

Esta skill conduz o diálogo e o stop-gate de confirmação. **Toda a lógica reside no comando `legoloop refinar`**: checagem da matriz, regeneração atômica da Tríade de governança, preservação de histórico e backup em `.legoloop-bak`.

* **Consumes (Entradas Obrigatórias):**
  - nome do processo (`<nome>`)
  - identificador da etapa (`--etapa <ID>`)
  - novos parâmetros opcionais (`--peca <nome>`, `--caso <caso>`, `--titulo <texto>`)
* **Produces (Saídas Mecânicas):**
  - nova versão atômica do processo (`processos/<nome>.processo.json`, `.studio.json`, `.dag.json`, Tríade de docs e skill)
* **Target Files (escopo de escrita permitido):**
  - nenhum arquivo é alterado manualmente; apenas o comando `legoloop refinar` grava

### Procedimento (passos determinísticos)

1. **Inspeção Prévia:**
   - Liste as etapas do processo com `legoloop listar` para confirmar o ID exato da etapa (`e01`, `e02`, etc.) e a versão atual.
2. **Prévia com Stop-Gate:**
   - Execute o comando em modo interativo (sem `--yes`): `legoloop refinar <nome> --etapa <ID> [--peca <P>] [--caso <C>]...`.
   - Analise o diff projetado e mostre as alterações ao usuário.
   - Solicite aprovação explícita ("sim").
3. **Execução:**
   - Com confirmação concedida, aplique as mudanças via flag `--yes`: `legoloop refinar <nome> --etapa <ID> [--peca <P>] [--caso <C>]... --yes`.
4. **Relato:**
   - Informe a nova versão do processo (ex.: 1.1.0), confirme a preservação da versão anterior no histórico e exiba os artefatos atualizados.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** edite arquivos `.processo.json`, `.dag.json` ou docs de governança manualmente; use sempre `legoloop refinar`.
> 2. **NUNCA** use `--yes` sem prévia aprovação humana do diff das alterações.
> 3. **NUNCA** altere a etapa se o ID informado não existir no processo.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop refinar legoloop-exemplo --etapa e01 --peca nova-skill`
- `legoloop refinar legoloop-exemplo --etapa e02 --caso tdd --yes`
