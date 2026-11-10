---
name: legoloop-remover
description: Use quando o usuário quiser REMOVER um processo do cofre e limpar seus vínculos nos harnesses com segurança, via /legoloop-remover. NÃO use para listar, refinar ou compor processos.
consumes: Artifact:Process
produces: Artifact:Vault
artifact: ~/.aidd/canonical/processos/<nome>.processo.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-remover — Remoção Segura de Processos do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja excluir permanentemente um processo do cofre (`legoloop-*`), removendo seus documentos do Studio, DAGs, documentos de governança e desfazendo todos os links nos harnesses.  
**Quando NÃO usar:** para remover arquivos isolados de código, canonizar arquivos (`canonize`) ou desinstalar harnesses.

Esta skill protege contra exclusões acidentais ou quebras de dependência. **Toda a lógica reside no comando `legoloop remover`**: verificação de integridade referencial (`ErrEmUso`), confirmação interativa e remoção atômica com rollback.

* **Consumes (Entradas Obrigatórias):**
  - nome do processo (`<nome>`)
* **Produces (Saídas Mecânicas):**
  - remoção atômica de arquivos do cofre e desfazimento de links nos harnesses
* **Target Files (escopo de escrita permitido):**
  - nenhum arquivo é deletado manualmente pela skill; apenas o comando `legoloop remover` remove

### Procedimento (passos determinísticos)

1. **Checagem de Dependências:**
   - Execute `legoloop listar` para verificar se o processo está em uso como sub-etapa por algum outro processo (`Usado por`).
   - Se estiver em uso, alerte o usuário: o comando falhará com `ErrEmUso` até que o processo pai seja removido ou refinado.
2. **Prévia e Confirmação com Stop-Gate:**
   - Execute `legoloop remover <nome>` sem `--yes`. O comando listará os arquivos e links que serão excluídos e solicitará `[s/N]`.
   - Apresente os itens ao usuário e pare no stop-gate humano. Exija confirmação explícita ("sim").
3. **Execução:**
   - Com confirmação concedida, aplique a remoção definitiva com `legoloop remover <nome> --yes`.
4. **Relato:**
   - Confirme a remoção dos artefatos e a limpeza dos links nos harnesses.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** apague arquivos do cofre (`rm`, `Remove-Item`) manualmente; utilize sempre `legoloop remover`.
> 2. **NUNCA** use `--yes` sem confirmação humana explícita prévia.
> 3. **NUNCA** force a remoção de um processo que outro processo declara como dependência.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop remover legoloop-antigo`
- `legoloop remover legoloop-antigo --yes`
