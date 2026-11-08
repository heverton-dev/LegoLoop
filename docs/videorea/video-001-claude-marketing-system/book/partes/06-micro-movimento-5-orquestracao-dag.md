# Parte IV — Micro: Orquestração em DAG e Governança

## 4.1 O Cérebro Central do Sistema

O quinto movimento não cria conteúdo nem conversa com clientes: ele é o maestro que gerencia o fluxo de dados entre os 4 movimentos anteriores.

```mermaid
graph TD
  M1["1. Research"] -->|"dossier.json"| M2["2. Attract"]
  M2 -->|"leads capturados"| M3["3. Convert"]
  M3 -->|"clientes ativos"| M4["4. Retain"]
  M4 -->|"reviews & referrals"| M5["5. Orchestrator"]
  M5 -->|"feedback de mercado"| M1
```

### Contrato do Grafo Acíclico Dirigido:
* **Entrada:** `market_brief.json`
* **Processamento:** 5 nós estéreis isolados
* **Saída:** `daily_marketing_report.json` com telemetria de leads, conversão, custo por aquisição e gargalos do funil.
