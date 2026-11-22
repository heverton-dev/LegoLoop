---
schema: anticorpos.v1
processo: legoloop-forjar
total_anticorpos: 5
globais: 0
dominio: 5
regras_inviolaveis:
  - pacote: "internal/forja"
    antidoto: "Simulação (--dry-run) nunca altera arquivo no disco; hash da árvore do HOME permanece idêntico."
  - pacote: "internal/forja"
    antidoto: "Toda alteração em --apply grava backup prévio com MANIFESTO.tsv antes de sobrescrever."
  - pacote: "internal/forja"
    antidoto: "Todo destino passa por internal/jail.DentroDaRaiz; nomes com '..' ou caminho absoluto são bloqueados."
  - pacote: "internal/forja"
    antidoto: "Valores de chaves com SECRET, TOKEN, KEY ou PASSWORD nunca entram em RELATORIO.md, RELATORIO.json ou terminal."
  - pacote: "internal/forja"
    antidoto: "O caminho de instalação e Forja permanece 100% em Go, sem dependência obrigatória de Python."
---

# 🛡️ Banco de Anticorpos — legoloop-forjar

> Espelho executivo para leitura de agentes. Garante integridade mecânica da Forja em qualquer sistema operacional.

## Regras de Domínio (5)

### `internal/forja` — Simulação estéril
- **Antídoto:** Simulação (`--dry-run`) nunca altera arquivo no disco; hash da árvore do `HOME` permanece idêntico.
- **Erro Evitado:** Gravação acidental de configuração antes de autorização explícita do usuário.

### `internal/forja` — Cópia de segurança com manifesto
- **Antídoto:** Toda alteração em `--apply` grava backup prévio em `~/.aidd/backups/forja_<data>/<harness>/MANIFESTO.tsv`.
- **Erro Evitado:** Perda irreversível de arquivos do usuário e impossibilidade de restaurar com `--desfazer`.

### `internal/forja` — Path Jail estrito
- **Antídoto:** Todo destino passa por `internal/jail.DentroDaRaiz`; nomes com `..` ou caminho absoluto são bloqueados.
- **Erro Evitado:** Escrita fora do `HOME` alvo por travessia de diretórios.

### `internal/forja` — Sanitização de segredos
- **Antídoto:** Valores de variáveis ou chaves com `SECRET`, `TOKEN`, `KEY` ou `PASSWORD` nunca entram em `RELATORIO.md`, `RELATORIO.json` ou stdout.
- **Erro Evitado:** Vazamento de credenciais em relatórios de diagnóstico.

### `internal/forja` — Caminho de instalação sem Python
- **Antídoto:** `internal/forja`, `internal/instalar` e `internal/ste` permanecem livres de chamadas a Python.
- **Erro Evitado:** Quebra do princípio de binário único portátil em máquinas limpas.

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._
