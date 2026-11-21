---
id: INTENT-intv_1795263519855957000
status: APPROVED
date: 2026-11-21
schema_version: 2
upstream: [INTERVIEW-intv_1795263519855957000]
covers: [Q1]
---

# 🎯 INTENT: Distribuição e Forja do LegoLoop: instalação universal por vários canais (MSI, script, winget, scoop, brew, git clone, agente, MCP) que dispara a Forja (camada 1) e o processo de configuração lazy de todos os harnesses (camada 2), com AGENTS.md como SSOT

> **Status:** APPROVED  
> **Data:** 2026-11-21  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1795263519855957000  
> **Confiança da Entrevista:** 96% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Distribuição e Forja do LegoLoop: instalação universal por vários canais (MSI, script, winget, scoop, brew, git clone, agente, MCP) que dispara a Forja (camada 1) e o processo de configuração lazy de todos os harnesses (camada 2), com AGENTS.md como SSOT
* **Quem Sofre:** Usuários e desenvolvedores que instalam por qualquer canal, operadores leigos que pedem a instalação a um agente e quem tem harnesses com contexto inchado.
* **O Porquê:** Não existe instalação universal e a correção de contexto nos harnesses ainda não é um Processo reutilizável; um comando único evita divergência entre canais e a comprovação física elimina o "acho que instalou".
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Distribuição e Forja do LegoLoop: instalação universal por vários canais (MSI, script, winget, scoop, brew, git clone, agente, MCP) que dispara a Forja (camada 1) e o processo de configuração lazy de todos os harnesses (camada 2), com AGENTS.md como SSOT
  * **Amanhã:** Instalação universal do LegoLoop por todos os canais (script de terminal primeiro; MSI, winget, scoop, brew, git clone, agente e MCP delegando ao `legoloop instalar`) e SOs, disparando o Processo canônico Forja, que configura os harnesses detectados com `AGENTS.md` como SSOT, é fabricado e fechado pelo fluxo oficial e tem quatro gatilhos com o mesmo RELATORIO.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Instalação universal do LegoLoop por todos os canais (script de terminal primeiro; MSI, winget, scoop, brew, git clone, agente e MCP delegando ao `legoloop instalar`) e SOs, disparando o Processo canônico Forja, que configura os harnesses detectados com `AGENTS.md` como SSOT, é fabricado e fechado pelo fluxo oficial e tem quatro gatilhos com o mesmo RELATORIO.
* [ ] Objetivo 2: Em máquina limpa, `legoloop instalar` roda sem erro, `skills-custo` mostra cada harness com 1 skill, o segundo `--dry-run` não acha pendência, todo arquivo alterado tem backup restaurável, os quatro gatilhos geram o mesmo RELATORIO e o Processo Forja passa em `legoloop validate` e `legoloop fechar`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Nada; o escopo é total
* ⛔ sem onda posterior.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | `unit_task` | Código de produção | Implementação da lógica via TDD Red-Green. |
| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Em máquina limpa, `legoloop instalar` roda sem erro, `skills-custo` mostra cada harness com 1 skill, o segundo `--dry-run` não acha pendência, todo arquivo alterado tem backup restaurável, os quatro gatilhos geram o mesmo RELATORIO e o Processo Forja passa em `legoloop validate` e `legoloop fechar`.
* **Restrições Inegociáveis:**
  * Nunca sobrescrever arquivo do usuário sem backup
  * nunca imprimir segredos
  * nunca pedir administrador sem aviso
  * conferir SHA-256 antes de executar e todo canal com `--dry-run`; valem ainda Zero Stubs
  * TDD estrito
  * recompilação fresca
  * a Forja sem tocar o grafo e execução pelo Swarm em `docs/routes/` sob as seis condições do subagente efêmero (contexto só com artefatos do ticket
  * ferramentas listadas
  * um subagente por ticket
  * saída em disco com gate mecânico
  * sem sub-subagentes
  * `write_boundary` `strict`).
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 96%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Em máquina limpa, `legoloop instalar` roda sem erro, `skills-custo` mostra cada harness com 1 skill, o segundo `--dry-run` não acha pendência, todo arquivo alterado tem backup restaurável, os quatro gatilhos geram o mesmo RELATORIO e o Processo Forja passa em `legoloop validate` e `legoloop fechar`. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-21_SPEC_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md`](../specs/2026-11-21_SPEC_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 96% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*

---

## Anexo A. Ordem de implementação, execução pelo Swarm e escopo expandido

Fonte integral: `docs/intents/ESPELHO_CANONICO.md`. Aprovado pelo usuário.

### ordem_de_implementacao (aprovada pelo usuário)

A ordem segue dependências, para não gerar retrabalho nem débito. Perfis e `--desfazer` pertencem ao motor da Forja.

1. Contrato (INTENT, SPEC) e linha de base medida.
2. CI nos três SOs (Windows, Linux, macOS), para todo ticket nascer testado.
3. Portar para Go o Python do caminho de instalação.
4. Motor da Forja: perfis, SSOT com bloco gerenciado, relatório, `--verificar` e `--desfazer`.
5. `legoloop instalar` chamando a Forja, com `--dry-run`.
6. Gate de orçamento de tokens por harness no CI.
7. Canonizar o Processo Forja (`fechar`, skill no Cofre, `/legoloop-forjar`).
8. Release reproduzível com assinatura, SBOM e SHA-256.
9. Canais: script de terminal primeiro; depois MSI, winget, scoop, brew, git clone, agente e MCP.
10. Detecção de deriva.
11. Canal de atualização (estável e beta).

### execucao_pelo_swarm (autorizada pelo usuário)

O Swarm executa o plano sob as seis condições da lei 4 do `AGENTS.md`:
1. O contexto de cada worker contém só os artefatos do ticket.
2. O ticket lista as ferramentas permitidas.
3. Um worker executa um ticket.
4. A saída fica em disco e passa pelo gate mecânico.
5. O worker não cria sub-workers.
6. O `write_boundary` é `strict`.

### forma_do_plano

Macro plano com iniciativas, épicos, micro tickets e nano tickets, de qualquer tamanho, sem nada de fora.
