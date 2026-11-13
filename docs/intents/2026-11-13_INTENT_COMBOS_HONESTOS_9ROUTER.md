---
id: INTENT-I7
status: DRAFT
date: 2026-11-13
schema_version: 2
upstream: []
covers: [Q1, Q2, Q3]
---

# 🎯 INTENT: Combos Honestos do 9Router (Iniciativa I7)

> **Status:** DRAFT  
> **Data:** 2026-11-13  
> **Autor:** i7-architect (Arquiteto I7)  
> **Iniciativa:** I7 — Combos Honestos do 9Router  
> **Alvo:** `internal/provider`, rotas de modelo do 9Router e catálogo canônico de combos  
> **Sessão:** `2269091e-b05b-4acd-a6ba-199f7248b3aa`  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

* **O Problema Real:** Harnesses e agentes emitindo requisições com histórico longo (~30k a 115k tokens) falham sistematicamente com erro HTTP 413: `Request too large for model openai/gpt-oss-120b ... TPM: Limit 8000, Requested 114246`. A evidência empírica coletada no diagnóstico em `docs/reports/2026-11-13_REPORT_DIAGNOSTICO_413_9ROUTER.md` comprova que o plano gratuito do Groq impõe teto de 8.000 tokens por minuto. Esse modelo foi posicionado na linha de frente dos combos `code-free`, `code-fast`, `code-pro`, `micro-antigo`, `free-raciocinio-M` e `free-ferramentas-M`. Além disso, modelos Cloudflare como Llama-3.1-70b e Llama-3.3-70b sofrem de janela restrita com falha 413 em 100k tokens.
* **Quem Sofre:** Todos os agentes de harness que operam com contexto amplo (sistema, ferramentas e histórico de sessão), bem como os Processos automatizados do LegoLoop executados via 9Router.
* **O Porquê:** O erro inviabiliza execuções complexas em harnesses de contexto grande. Como o 9Router devolve o erro do primeiro membro recusado, a falha 413 do Groq mascara eventuais membros vivos subsequentes e derruba a execução do agente. É indispensável estruturar combos honestos, testados e comprovados em 100k tokens.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Modelos Groq com teto de 8k TPM e modelos Cloudflare com janela curta ficam na frente de combos gerais, gerando falhas 413 frequentes em pedidos de 30k-115k tokens.
  * **Amanhã:** Seis combos honestos operam com membros certificados em testes reais de 100k tokens (`t-agente-128k`, `t-chat-128k`, `t-raciocinio-128k`, `t-visao-128k`, `t-codigo-128k`), enquanto modelos Groq ficam restritos exclusivamente ao combo `t-micro-8k` para requisições de até 8.000 tokens.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Formalizar a arquitetura e os contratos dos 6 combos honestos certificados por testes reais em 100k tokens (`t-agente-128k`, `t-chat-128k`, `t-raciocinio-128k`, `t-visao-128k`, `t-codigo-128k` e `t-micro-8k`).
* [ ] Objetivo 2: Estabelecer barreira mecânica e validação em Go para impedir o roteamento de prompts maiores que 8.000 tokens para modelos Groq ou combos micro.
* [ ] Objetivo 3: Criar manifesto de sincronização e mapeamento de compatibilidade para os nomes de combos legados (`code-pro`, `code-fast`, `code-free`), eliminando Groq e Cloudflare da frente de rotas de contexto amplo.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não alterar a infraestrutura física de rede ou o deploy da VPS sem homologação prévia.
* ⛔ Não reintroduzir modelos reprovados em 100k tokens (Groq geral, Llama-3.1/3.3-70b Cloudflare, Gemini 3.8 Flash em 429 persistente).
* ⛔ Não alterar contratos de templates canônicos existentes em `docs/templates/`.
* ⛔ Não violar a regra de Zero Stubs ou quebrar a suíte de testes de `internal/provider` e `internal/brain`.

---

## 3. Mapeamento Preliminar de Peças

Identificação das Peças que o Processo acionará.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_task` | Catálogo de Combos em Go | Definição dos tipos, membros e janelas de contexto |
| `PECA-02` | `unit_task` | Validador TDD de TPM/Janela | Teste e filtro mecânico contra erro HTTP 413 |
| `PECA-03` | `script` | Gerador de Manifesto de Combos | Produção do JSON de sincronização para o 9Router |
| `PECA-04` | `script` | `go test ./...` | Validação de integridade e regressão na suíte Go |
| `PECA-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória do binário |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** `Médio` (afeta a resolução de modelos e combos de LLM no LegoLoop e gateway).
* **Política de Privacidade & Orçamento Cognitivo:** `Nuvem autorizada` via 9Router com chave confiável (`provider.TrustedRouterKey`), sem logar chaves ou segredos.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Catálogo formal dos 6 combos validado em testes unitários Go, manifesto JSON gerado e compatibilidade de combos legados verificada sem chamadas Groq em prompts > 8k tokens.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  1. Os testes empíricos registrados em `docs/reports/data/2026-11-13_homologacao_membros.json` e `..._reteste.json` constituem a fonte da verdade para a seleção dos membros.
  2. Requisições de agentes em harnesses padrão Claude Code exigem no mínimo 100k tokens de volume sem estouro de TPM.
  3. Modelos Groq são úteis para execuções atômicas (< 8k tokens) devido à baixa latência (1,6s - 2,0s), mas não devem receber tráfego de agentes com contexto amplo.
* **Riscos e Perguntas em Aberto:** Nenhum risco bloqueante; dados empíricos já foram levantados e consolidados no relatório de diagnóstico.

---

## 5. Critérios de Sucesso de Alto Nível

Cada critério possui forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todos os 6 combos honestos possuem contrato formal com membros e tetos declarados. | Presença da definição formal e validação em Go | `go test -v ./internal/...` |
| 2 | O combo `t-micro-8k` rejeita mecanicamente requisições acima de 8.000 tokens. | Teste unitário de barreira de tokens falhando antes da correção e verde após | `go test -v -run TestTetoTokens ./...` |
| 3 | Mapeamento de retrocompatibilidade de `code-pro`, `code-fast` e `code-free` não aponta Groq na primeira posição. | Verificação estática de ordem de membros | Asserção binária em teste Go |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md`](../specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após validação e geração da SPEC correspondente.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada para eliminação do erro 413 e criação dos 6 combos honestos atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
