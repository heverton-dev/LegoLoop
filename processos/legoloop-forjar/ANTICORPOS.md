---
schema: anticorpos.v1
processo: legoloop-forjar
total_anticorpos: 6
globais: 1
dominio: 5
regras_inviolaveis:
  - pacote: ""
    antidoto: "Verificação: len(harnesses) >= 1"
  - pacote: ""
    antidoto: "Verificação: hash(home) antes == hash(home) depois"
  - pacote: ""
    antidoto: "Verificação: alterados do segundo run == 0"
  - pacote: ""
    antidoto: "Verificação: exit_code == 1 quando o teto estoura"
  - pacote: ""
    antidoto: "Verificação: RELATORIO.json sem segredos e veredito == OK"
  - pacote: "auto-provider-onboarding/cmd/gw (binarios Go locais)"
    antidoto: "Exclusoes permanentes de path e processo no Defender via scripts/defender_exclusoes.ps1 do repo auto-provider-onboard..."
---

# 🛡️ Banco de Anticorpos — legoloop-forjar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (5)

### ``
- **Antídoto:** Verificação: len(harnesses) >= 1
- **Erro Evitado:** Detectar os harnesses instalados no HOME antes de planejar a Forja
- **Assinatura:** `sha256:882c4c23b63ed81de6c2214a9a51c13baf9354c1a534c580ec32fd8000cfe663`

### ``
- **Antídoto:** Verificação: hash(home) antes == hash(home) depois
- **Erro Evitado:** Simular a Forja sem gravar nada e listar as Etapas como pendente
- **Assinatura:** `sha256:cd4cdf1e2a19e838b6c15a130bf91df3273805c8cb30a6ff55291ad951c7ab34`

### ``
- **Antídoto:** Verificação: alterados do segundo run == 0
- **Erro Evitado:** Aplicar a Forja com backup MANIFESTO.tsv e idempotencia no segundo run
- **Assinatura:** `sha256:059f51a1aa70fc924717f3312c3c85aaaa86ed72118079f9baad34cba3a8a8ee`

### ``
- **Antídoto:** Verificação: exit_code == 1 quando o teto estoura
- **Erro Evitado:** Reprovar o harness cujo custo medido excede o teto do orcamento de tokens
- **Assinatura:** `sha256:a1372bf8fcf521529a9941b80ca669892e043ee8efb8699d49149dfbec72d27a`

### ``
- **Antídoto:** Verificação: RELATORIO.json sem segredos e veredito == OK
- **Erro Evitado:** Emitir RELATORIO.json sem valores de segredos e com veredito OK
- **Assinatura:** `sha256:1a61e89c54644daf32c16254ba4760f037de61b557d4bb7649afefc671085f37`

## Regras Globais (1)

### `auto-provider-onboarding/cmd/gw (binarios Go locais)`
- **Antídoto:** Exclusoes permanentes de path e processo no Defender via scripts/defender_exclusoes.ps1 do repo auto-provider-onboarding (requer admin; idempotente). Rebuild se perder: go build -o $env:USERPROFILE/.local/bin/gw.exe ./cmd/gw no repo auto-provider-onboarding
- **Erro Evitado:** gw.exe/legoloop*.exe compilados localmente somem de ~/.local/bin e bin/ (claude-gateway exit 127: gw.exe: No such file or directory)
- **Causa Raiz:** Windows Defender com heuristica ML classifica binarios Go locais como Trojan:Win32/Bearfoos.A!ml e os poe em quarentena (confirmado em Get-MpThreatDetection e evento 5007)
- **Assinatura:** `sha256:1f29af1d81349af5c43998d6b0620844d0d179807c5bf83002b01f9c3ddc4ed7`

