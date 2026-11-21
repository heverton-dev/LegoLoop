---
schema: anticorpos.v1
processo: legoloop-agent-desktop
total_anticorpos: 2
globais: 0
dominio: 2
regras_inviolaveis:
  - pacote: "legoloop-agent-desktop/accessibility-tree"
    antidoto: "Utilizar estritamente os comandos de snapshot e referências semânticas (@ref) via UIA no Windows e AXUI no macOS/Linu..."
  - pacote: "legoloop-agent-desktop/powershell-quoting"
    antidoto: "Sempre envolver referências em aspas simples (ex.: '@s8f3k2p9:e1') em comandos Windows PowerShell."
---

# 🛡️ Banco de Anticorpos — legoloop-agent-desktop

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (2)

### `legoloop-agent-desktop/accessibility-tree`
- **Antídoto:** Utilizar estritamente os comandos de snapshot e referências semânticas (@ref) via UIA no Windows e AXUI no macOS/Linux, garantindo determinismo sem depender de OCR ou coordenadas de pixels.
- **Erro Evitado:** Tentativa de automação desktop através de reconhecimento visual por pixels ou screenshot em vez da árvore de acessibilidade nativa.
- **Causa Raiz:** Suposição incorreta de que o agente desktop precisa de visão computacional (CV) para interagir com a interface.
- **Assinatura:** `sha256:d1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2`

### `legoloop-agent-desktop/powershell-quoting`
- **Antídoto:** Sempre envolver referências em aspas simples (ex.: '@s8f3k2p9:e1') em comandos Windows PowerShell.
- **Erro Evitado:** Execução de comandos com referências cruas (@ref) no PowerShell sem aspas, provocando erro de splatting do shell.
- **Causa Raiz:** O caractere @ no PowerShell é o operador de splatting de variáveis quando colocado no início de um argumento não citado.
- **Assinatura:** `sha256:d2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

