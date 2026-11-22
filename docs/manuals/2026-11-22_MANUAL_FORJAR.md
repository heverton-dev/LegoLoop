# 📘 MANUAL OPERACIONAL: PROCESSO `legoloop-forjar` (FORJA DE AMBIENTE)

> **Data:** 2026-11-22  
> **Autoridade:** Ecossistema LegoLoop (Iniciativa I01 — Distribuição e Forja)  
> **Status:** Canônico  
> **Escopo:** Preparação determinística, segura e reversível do ambiente de harnesses do usuário via CLI, instalação, MCP ou Skill  
> **Fontes de verdade:** `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md` · `internal/forja/` · `processos/legoloop-forjar/` · `skills/legoloop-forjar/SKILL.md`

---

## 1. Visão Geral

O Processo **`legoloop-forjar`** padroniza o ambiente local do usuário em Windows, Linux e macOS. Todos os canais de entrada (`legoloop instalar`, `legoloop forjar`, `call_mcp_tool` e `/legoloop-forjar`) executam o mesmo motor em `internal/forja` e produzem o mesmo `RELATORIO.json`.

### Perfis de Contexto

| Perfil | Efeito principal |
| :--- | :--- |
| `minimo` | Instala apenas o índice lazy de Skills no harness |
| `padrao` | Índice lazy somado às regras de redução de contexto (`permissions.deny`) |
| `gateway` | Modo de ferramenta única via gateway MCP |

---

## 2. Etapas do DAG

1. **`e01` (`detectar-harnesses`):** localiza os harnesses instalados no `HOME` e valida o contrato de entrada e o Path Jail.
2. **`e02` (`simular-forja`):** executa `--dry-run` e lista todas as alterações pendentes sem gravar no disco.
3. **`e03` (`aplicar-forja`):** grava cópia de segurança com `MANIFESTO.tsv`, aplica blocos gerenciados do `AGENTS.md`, faz merge de `settings.json` preservando chaves desconhecidas e cria lançadores do sistema operacional.
4. **`e04` (`verificar-orcamento`):** compara o custo medido por harness com os tetos definidos em `docs/governanca/orcamento_tokens.json`.
5. **`e05` (`emitir-relatorio`):** consolida `RELATORIO.md` e `RELATORIO.json` com sanitização de segredos.

---

## 3. Comandos Operacionais

```bash
# Simulação sem gravação
legoloop forjar --dry-run --perfil padrao --json

# Aplicação com cópia de segurança automática
legoloop forjar --apply --perfil padrao --json

# Excluir um harness específico
legoloop forjar --apply --sem opencode --json

# Verificar orçamento de tokens
legoloop forjar --verificar --orcamento docs/governanca/orcamento_tokens.json --json

# Restaurar estado anterior byte a byte
legoloop forjar --desfazer 2026-11-22
```

---

## 4. Segurança e Reversibilidade

- **Path Jail:** todo destino é validado por `internal/jail.DentroDaRaiz` antes de qualquer escrita.
- **Backup obrigatório:** toda alteração grava antes a versão original e registra o par origem/destino em `MANIFESTO.tsv`.
- **Segredos protegidos:** chaves contendo `SECRET`, `TOKEN`, `KEY` ou `PASSWORD` jamais entram nos relatórios.
- **Idempotência:** executar `--apply` duas vezes consecutivas produz zero alterações na segunda passagem.
