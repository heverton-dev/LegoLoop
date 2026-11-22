---
name: legoloop-forjar
description: Use quando o usuário quiser forjar ambiente do LegoLoop nos harnesses locais com simulação, backup restaurável, perfis de contexto (minimo, padrao, gateway), SSOT AGENTS.md e Gate de orçamento de tokens via /legoloop-forjar.
consumes: Artifact:Home
produces: Artifact:Report
artifact: RELATORIO.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-forjar — Processo da Forja de Ambiente do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja forjar ambiente do LegoLoop, configurar harnesses detectados no `HOME`, aplicar perfis de contexto (`minimo`, `padrao`, `gateway`), verificar o orçamento de tokens ou restaurar um backup anterior com `--desfazer`.  
**Quando NÃO usar:** para fabricar Processo novo (`legoloop fabricar`), listar Processos do Cofre (`legoloop listar`) ou fechar Processo (`legoloop fechar`).

Esta Skill conduz as cinco Etapas do DAG do Processo `legoloop-forjar`. Toda a lógica de execução reside no comando `legoloop forjar` (`internal/forja`).

* **Consumes (Entradas Obrigatórias):**
  - diretório `HOME` do usuário e arquivo de orçamento `docs/governanca/orcamento_tokens.json`
* **Produces (Saídas Mecânicas):**
  - `RELATORIO.md`, `RELATORIO.json` e cópia de segurança com `MANIFESTO.tsv`
* **Target Files (escopo de escrita permitido):**
  - somente arquivos de configuração dentro do `HOME` alvo, validados por Path Jail

### Procedimento (Etapas do DAG)

1. **Etapa `e01` (`detectar-harnesses`):**
   - Execute `legoloop forjar --dry-run --json` para detectar os harnesses ativos no `HOME` e validar o contrato.
2. **Etapa `e02` (`simular-forja`):**
   - Execute `legoloop forjar --dry-run --perfil padrao --json` para pré-visualizar todas as alterações pendentes sem gravar nada no disco.
   - Apresente o resumo ao usuário no Stop-Gate antes de aplicar.
3. **Etapa `e03` (`aplicar-forja`):**
   - Após autorização explícita, execute `legoloop forjar --apply --perfil <perfil> --json`.
   - O motor grava backup prévio com `MANIFESTO.tsv`, aplica trechos gerenciados ligados ao `AGENTS.md`, atualiza `settings.json` preservando chaves do usuário e cria lançadores do sistema operacional.
4. **Etapa `e04` (`verificar-orcamento`):**
   - Execute `legoloop forjar --verificar --orcamento docs/governanca/orcamento_tokens.json --json` para confirmar que nenhum harness excedeu o teto de tokens.
5. **Etapa `e05` (`emitir-relatorio`):**
   - Confirme o veredito `OK` no `RELATORIO.json` e entregue os caminhos de relatório e backup ao usuário.
   - Para reverter qualquer aplicação, execute `legoloop forjar --desfazer <AAAA-MM-DD>`.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** execute `--apply` sem antes mostrar a simulação (`--dry-run`) e obter autorização no Stop-Gate.
> 2. **NUNCA** sobrescreva arquivo do usuário sem backup prévio em `MANIFESTO.tsv`.
> 3. **NUNCA** registre segredos (`SECRET`, `TOKEN`, `KEY`, `PASSWORD`) no relatório ou na saída do terminal.
> 4. **NUNCA** escreva fora do `HOME` alvo.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop forjar --dry-run --json`
- `legoloop forjar --apply --perfil padrao --json`
- `legoloop forjar --verificar --orcamento docs/governanca/orcamento_tokens.json --json`
- `legoloop forjar --desfazer <AAAA-MM-DD>`
