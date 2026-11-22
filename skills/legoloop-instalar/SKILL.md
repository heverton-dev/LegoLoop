---
name: legoloop-instalar
description: Use quando o usuário quiser instalar o LegoLoop no ambiente local, simular ou aplicar a Forja via CLI (legoloop instalar) ou pela ferramenta nativa MCP legoloop:instalar.
consumes: Artifact:Home
produces: Artifact:Report
artifact: RELATORIO.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-instalar — Instalação do LegoLoop via CLI e MCP

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja instalar o LegoLoop, configurar harnesses detectados no `HOME`, criar o atalho do grafo (`codebase-memory-mcp`) ou acionar a instalação via servidor MCP `legoloop` (`legoloop:instalar`).  
**Quando NÃO usar:** para fabricar Processo novo (`legoloop fabricar`), listar Processos do Cofre (`legoloop listar`) ou fechar Processo (`legoloop fechar`).

Esta Skill orienta os dois pontos de entrada canônicos da instalação: o comando `legoloop instalar` e a ferramenta nativa `legoloop:instalar` no servidor MCP `legoloop`. Ambos executam o mesmo motor em `internal/instalar` encadeado com `internal/forja`.

* **Consumes (Entradas Obrigatórias):**
  - diretório `HOME` alvo e perfil de contexto (`minimo`, `padrao` ou `gateway`)
* **Produces (Saídas Mecânicas):**
  - `RELATORIO.json`, `RELATORIO.md`, cópia de segurança com `MANIFESTO.tsv` e atalho do grafo quando aplicável
* **Target Files (escopo de escrita permitido):**
  - somente arquivos dentro do `HOME` alvo, validados por Path Jail

### Procedimento via CLI (`legoloop instalar`)

1. **Simulação segura (sem escrita):**
   - Execute `legoloop instalar --dry-run --perfil padrao --json` para inspecionar o plano de instalação sem gravar arquivos no disco.
   - Apresente o resumo ao usuário no Stop-Gate antes de aplicar.
2. **Aplicação controlada:**
   - Após autorização explícita, execute `legoloop instalar --apply --perfil padrao --json`.
   - O comando grava cópia de segurança com `MANIFESTO.tsv`, configura os harnesses detectados, cria o atalho do grafo e grava `RELATORIO.json` e `RELATORIO.md` em `~/.legoloop/forja/<AAAA-MM-DD>/`.
3. **Exclusão de harness (opcional):**
   - Use `--sem <harness>` para ignorar um harness específico sem alterar seus arquivos.

### Procedimento via MCP (`legoloop:instalar`)

1. **Simulação via `call_mcp_tool`:**
   - Chame `call_mcp_tool` com `server: "legoloop"`, `tool: "instalar"` e `arguments: {"dry_run": true, "home": "<HOME>", "perfil": "padrao"}`.
   - A resposta JSON retorna `relatorio_json` (caminho previsto do `RELATORIO.json`) e `veredito: "OK"`, sem gravar nada no disco.
2. **Aplicação via `call_mcp_tool`:**
   - Somente após autorização no Stop-Gate, chame `call_mcp_tool` com `server: "legoloop"`, `tool: "instalar"` e `arguments: {"dry_run": false, "home": "<HOME>", "perfil": "padrao"}`.
   - Confirme o arquivo `RELATORIO.json` gravado no caminho retornado em `relatorio_json` e o `veredito: "OK"`.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** execute `--apply` ou `dry_run: false` sem antes rodar a simulação (`--dry-run` ou `dry_run: true`) e obter autorização no Stop-Gate.
> 2. **NUNCA** sobrescreva arquivo do usuário sem cópia de segurança prévia em `MANIFESTO.tsv`.
> 3. **NUNCA** registre segredos (`SECRET`, `TOKEN`, `KEY`, `PASSWORD`) no relatório ou na saída do terminal.
> 4. **NUNCA** escreva fora do `HOME` alvo.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop instalar --dry-run --perfil padrao --json`
- `legoloop instalar --apply --perfil padrao --json`
- `call_mcp_tool` -> `server: "legoloop"`, `tool: "instalar"`, `arguments: {"dry_run": true, "perfil": "padrao"}`
