# 1. Macro — Fundação da Forja

A Forja unifica a preparação do ambiente do LegoLoop sob um único motor em Go (`internal/forja`). Todos os quatro gatilhos — instalação (`legoloop instalar`), CLI (`legoloop forjar`), MCP (`call_mcp_tool`) e Skill (`/legoloop-forjar`) — convergem para `forja.Executar` e emitem o mesmo par `RELATORIO.md` e `RELATORIO.json`.

Princípios invioláveis:
- Simulação (`--dry-run`) nunca grava arquivos no disco.
- Toda alteração em `--apply` grava antes cópia de segurança com `MANIFESTO.tsv`.
- Todo caminho destino passa por `internal/jail.DentroDaRaiz`.
- Nenhum segredo entra nos relatórios.
