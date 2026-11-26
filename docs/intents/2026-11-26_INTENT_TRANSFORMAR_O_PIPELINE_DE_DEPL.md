---
id: INTENT-intv_1795723484496401600
status: APPROVED
date: 2026-11-26
schema_version: 2
upstream: [INTERVIEW-intv_1795723484496401600]
covers: [Q1]
---

# 🎯 INTENT: Transformar o pipeline de deploy determinístico do Studio Launcher (extração anti-zip-slip de .zip…

> **Status:** APPROVED (Stop-Gate 01)  
> **Data:** 2026-11-26  
> **Autor:** Usuário + Agente (via Entrevista Socrática, motor `legoloop fabricar entrevista`)
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1795723484496401600  
> **Confiança da Entrevista:** 10% (motor estacionou; Espelho Canônico reconciliado manualmente — ver `ESPELHO_CANONICO.md`)
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Transformar o pipeline de deploy determinístico do Studio Launcher (extração anti-zip-slip de .zip, fusão de .env mestre, reescrita SSOT da label Traefik, transação DNS Cloudflare CNAME, docker build e swarm deploy, gates de verificação binária e recibo JSON) em um Processo 100% determinístico do LegoLoop (ex.: legoloop-kraften-deploy) com 5 nós (n01_preflight_zip, n02_sync_metadata, n03_deploy_stream, n04_verify_gates, n05_receipt), zero LLM, para publicar aplicações Go/Node/Python na VPS com subdomínio e TLS automáticos. Documento completo do fluxo: C:/Users/trcnologia/Desktop/APPS/studio-launcher-go/PROCESSO_DEPLOY_DETERMINISTICO_LEGOLOOP.md
* **Quem Sofre:** Desenvolvedores e operadores do ecossistema LegoLoop/Studio Launcher.
* **O Porquê:** Migrar o deploy manual/legado do Studio Launcher para um processo padronizado, auditável e 100% determinístico no LegoLoop.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista (ver O Problema Real).
  * **Amanhã:** Esteira LegoLoop `legoloop-kraften-deploy` funcional com 5 nós determinísticos (`n01_preflight_zip`, `n02_sync_metadata`, `n03_deploy_stream`, `n04_verify_gates`, `n05_receipt`), zero LLM, publicando aplicações de qualquer stack com `Dockerfile` e `docker-stack.yml` (Go, Node, Python, Rust, PHP, estático) com subdomínio e TLS automáticos.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Esteira LegoLoop `legoloop-kraften-deploy` funcional com 5 nós determinísticos (`n01_preflight_zip`, `n02_sync_metadata`, `n03_deploy_stream`, `n04_verify_gates`, `n05_receipt`), zero LLM, publicando aplicações de qualquer stack com `Dockerfile` e `docker-stack.yml` (Go, Node, Python, Rust, PHP, estático) com subdomínio e TLS automáticos.
* [ ] Objetivo 2: critérios de sucesso:
  * Execução completa dos 5 nós em sequência, gerando build, DNS Cloudflare, Docker Swarm deploy, verificação de gates e recibo JSON sem falhas.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Edição do `Dockerfile` da aplicação (nunca mutado pela esteira).
* ⛔ Agentes de IA/LLM intermediários nos nós.
* ⛔ Stacks sem `Dockerfile`/`docker-stack.yml` (fronteira universal do `.zip`).
* ⛔ Redes Docker diferentes de `network_conexao`.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `n01_preflight_zip` | `script` | Validação determinística (Go) | Validar `.zip`/pasta: `Dockerfile`, `docker-stack.yml` na rede `network_conexao`, `HEALTHCHECK`, variáveis padrão. |
| `n02_sync_metadata` | `script` | API do Studio Launcher + Cloudflare | Registrar/atualizar metadados e subdomínio; criar CNAME antes de qualquer consulta DNS. |
| `n03_deploy_stream` | `script` | API do Studio Launcher | Enviar `.zip` e acompanhar `deploy-job` até `status == "done"`. |
| `n04_verify_gates` | `script` | Gates binários (exit code) | Provar `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`. |
| `n05_receipt` | `script` | Recibo auditável | Gravar `DEPLOY_RECEIPT.json` sem acionar agente conversacional. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  * Execução completa dos 5 nós em sequência, gerando build, DNS Cloudflare, Docker Swarm deploy, verificação de gates e recibo JSON sem falhas.
* **Restrições Inegociáveis:**
  * Zero LLM (nenhum agente intermediário nos nós).
  * Implementação da esteira em Go nativo — motor em `internal/kraften/`, comando em `cmd/legoloop` (padrão real do repo).
  * Docker/Cloudflare CLI e REST API v4.
  * Respeito rigoroso à especificação de `PROCESSO_DEPLOY_DETERMINISTICO_LEGOLOOP.md` sem alterações estruturais.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** A esteira atende qualquer stack com `Dockerfile` e `docker-stack.yml`; Go nativo é a implementação da esteira, não restrição das aplicações (divergência do motor reconciliada pelo usuário).
* **Riscos e Perguntas em Aberto:** Confiança da entrevista estacionou em 10% (motor `code-fast`); Espelho Canônico emitido e reconciliado manualmente com o usuário em `docs/intents/ESPELHO_CANONICO.md`.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Execução completa dos 5 nós em sequência, gerando build, DNS Cloudflare, Docker Swarm deploy, verificação de gates e recibo JSON sem falhas. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-26_SPEC_TRANSFORMAR_O_PIPELINE_DE_DEPL.md`](../specs/2026-11-26_SPEC_TRANSFORMAR_O_PIPELINE_DE_DEPL.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 10% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
