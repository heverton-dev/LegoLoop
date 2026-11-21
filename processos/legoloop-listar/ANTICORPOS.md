---
schema: anticorpos.v1
processo: legoloop-listar
total_anticorpos: 6
globais: 0
dominio: 6
regras_inviolaveis:
  - pacote: ""
    antidoto: "Verificação: os.path.isdir(vault) && os.access(vault, os.R_OK)"
  - pacote: ""
    antidoto: "Verificação: json.loads(output).get('Processos') is not None"
  - pacote: ""
    antidoto: "Verificação: len(re.findall(r'class=\"node\"', html)) >= 6"
  - pacote: ""
    antidoto: "Verificação: '@keyframes flow-dash' in html"
  - pacote: ""
    antidoto: "Verificação: verify_geometry(html) == True"
  - pacote: ""
    antidoto: "Verificação: os.path.exists('docs/reports/2026-11-19_REPORT_LISTAR.md')"
---

# 🛡️ Banco de Anticorpos — legoloop-listar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (6)

### ``
- **Antídoto:** Verificação: os.path.isdir(vault) && os.access(vault, os.R_OK)
- **Erro Evitado:** Validar existência e acessibilidade do cofre canônico antes de listar processos
- **Assinatura:** `sha256:8f1ac3f6cc0c8538bfaddb2c3802197f8f301369acb268d1ff1d7cc4fe3d791c`

### ``
- **Antídoto:** Verificação: json.loads(output).get('Processos') is not None
- **Erro Evitado:** Garantir que a leitura de processos retorna estrutura JSON válida com campos obrigatórios
- **Assinatura:** `sha256:2817ca42f8be55ff3f1b359f18b59f6bc18e771bc6d1e30c3810df074b8a2a5f`

### ``
- **Antídoto:** Verificação: len(re.findall(r'class="node"', html)) >= 6
- **Erro Evitado:** Verificar que o diagrama HTML gerado contém todos os nós do processo listar
- **Assinatura:** `sha256:78431e6f6f47f5961e5578c610390f38e368cb1e733a83becafd92f838f60de4`

### ``
- **Antídoto:** Verificação: '@keyframes flow-dash' in html
- **Erro Evitado:** Confirmar que a animação CSS @keyframes flow-dash está presente e funcional
- **Assinatura:** `sha256:270d6d67933e458c857eed2e9d615031949452a1eeae62ca0c6d291baefe2f0f`

### ``
- **Antídoto:** Verificação: verify_geometry(html) == True
- **Erro Evitado:** Validar geometria: zero overflow, grid ortogonal, máscaras opacas
- **Assinatura:** `sha256:dba2ce7a4f09f2f6a1729c8ac541f511fc5c080865becf581374b5ffb0485bbe`

### ``
- **Antídoto:** Verificação: os.path.exists('docs/reports/2026-11-19_REPORT_LISTAR.md')
- **Erro Evitado:** Consolidar achados em anticorpos.json e emitir relatório em docs/reports/
- **Assinatura:** `sha256:87288f9d5d4935b284326fc62fa22f1d965ec81265ade5e7460dc4f1e3d4bb26`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

