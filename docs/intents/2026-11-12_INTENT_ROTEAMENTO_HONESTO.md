---
id: INTENT-I3
status: UNDER_REVIEW
date: 2026-11-12
schema_version: 2
upstream: [INTERVIEW-I3]
covers: [Q1, Q2, Q3, Q4, Q5, Q6, Q7, Q8, Q9, Q10, Q11, Q12]
---

# 🎯 INTENT: Roteamento honesto (Iniciativa I3)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-12  
> **Autor:** Claude  
> **Iniciativa:** I3 Roteamento honesto (programa Lei #15: I1 Taxonomia, I2 Sexteto e templates, I3 Roteamento honesto, I4 DISSECAR, I5 MATRIZ_PAPEIS, I6 Migração dos legados)  
> **Alvo:** `internal/roteamento`, `internal/orchestrix`, `internal/dag`, `internal/mcp`, `internal/provider`, `internal/cli`, `tools/route2dag`, `docs/templates`  
> **Sessão:** 240fe656-1081-4519-b96d-4dd4985fff63  
> **Confiança da Entrevista:** 95% ([INTERVIEW](../interviews/2026-11-12_INTERVIEW_ROTEAMENTO_HONESTO.md))  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** a ROUTE declara harness, modelo e esforço, mas o motor não cumpre o que ela declara. Seis fatos provam o problema.
  * O nó de código traduz `code-fast` e `code-pro` para `sonnet` e `flash` para `haiku` (F1, `internal/orchestrix/unit_run.go:321`).
  * O nó de código roda só no claude CLI (F5, `internal/orchestrix/unit_claude.go:43`).
  * O pacote `internal/roteamento` classifica o nó, lê o inventário e resolve a escolha, mas nenhum arquivo o importa (F7).
  * Nenhuma ferramenta escolhe o roteamento da ROUTE. O agente o escreve à mão (F13 e F14).
  * Não existe `skill_search` nem `call_skill` em Go (F10).
  * Não existe bateria que prove se um harness roda o nó de código (F15).
* **Quem Sofre:** o usuário que orquestra uma rodada, em toda rodada. A ROUTE mostra uma escolha e o motor executa outra, sem aviso.
* **O Porquê:** o LegoLoop quer ser agnóstico de harness e de modelo. Sem roteamento honesto, a ROUTE perde o valor como fonte da decisão (Q2 da I2).
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** o agente preenche o roteamento à mão. O motor colapsa os combos em `sonnet`. Só o claude CLI roda nó de código.
  * **Amanhã:** o comando `legoloop rota` grava a ROUTE com harness, modelo, esforço e `via` resolvidos contra o inventário do sistema operacional. O motor executa exatamente o que a ROUTE declara. Todo harness aprovado na bateria de conformidade pode rodar nó de código.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1 (Q8, Q9, Q12): o comando `legoloop rota` lê o PLAN, o perfil do usuário (global e por projeto, com padrão embutido) e o inventário, e grava a ROUTE. O contrato v2 do DAG ganha `routing.via`, obrigatório.
* [ ] Objetivo 2 (Q2, Q4): o motor executa o harness, o modelo, o esforço e o `via` declarados. Ele comuta o gateway por nó, só no processo filho. Em falha, os nós `nano` e `micro` degradam com registro, e os nós de classe grande e `extreme` falham.
* [ ] Objetivo 3 (Q3, Q10, Q11): a bateria de conformidade homologa harnesses por par harness e `via`. Ela roda na instalação e a cada harness novo, com aviso do teto de gasto e confirmação. O `internal/roteamento` e o `provider.Inventory` viram um inventário só.
* [ ] Objetivo 4 (Q5): o `skill-router` entrega `skill_search` e `call_skill` sobre o mesmo cofre do roteador de MCP. O `call_skill` devolve o texto da skill.
* [ ] Objetivo 5 (Q6, Q7): o motor limita os workers pela memória livre e pelo peso de cada harness. O alarme de banimento marca o harness como reprovado. O E05 da I1 é absorvido: ficam 05.1, 05.5 e 05.6 adaptados. Os nós 05.2 e 05.4 são substituídos. O nó 05.3 é removido.
* [ ] Objetivo 6 (regra de prova): um validador em Go recusa pergunta de entrevista sem fato citado.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia. O campo `routing.via` passa pelo Stop-Gate `SG-I3-CONTRATO`.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não reativaremos o rodízio de contas do agy (banido por ToS em 2026-11-10).
* ⛔ Não tocaremos o Cofre real `~/.aidd`. Os testes usam diretórios temporários.
* ⛔ Não migraremos os arquivos legados (iniciativa I6), nem o MANUAL_TEMPLATES_LEGOLOOP.
* ⛔ Não cobriremos o verbo DISSECAR (I4) nem a MATRIZ_PAPEIS (I5).

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | skills/aidd-spec | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | nó de código | Go (stdlib) com TDD Red-Green | Implementar o inventário unificado, o `legoloop rota`, a bateria, o `skill-router` e o limite de workers. |
| `PECA-03` | `script` | `go test -count=1 ./...` | Regressão mecânica e integridade. |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop_i2.exe ./cmd/legoloop` | Recompilação fresca. |
| `PECA-05` | `stop_gate` | `SG-I3-CONTRATO` | Aprovação humana do campo `routing.via` no contrato do DAG. |
| `PECA-06` | `stop_gate` | `SG-I3-FINAL` | Aprovação humana do fechamento da iniciativa. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`.
* **Nível de Risco & Blast Radius:** Alto. A iniciativa muda o contrato público do DAG e o executor do motor. Exige rollback `stash_changes` e commit por nó, só dos arquivos do `context_out`.
* **Política de Privacidade & Orçamento Cognitivo:** nuvem autorizada para o código da iniciativa. Teto de US$ 1,50 por nó. A bateria de conformidade faz chamadas reais só depois da confirmação do usuário (Q10). Segredos: a chave do 9Router vem da variável `NINEROUTER_KEY`. Nenhum artefato grava o valor.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** a ROUTE gerada pelo `legoloop rota` traz harness, modelo, esforço e `via` iguais aos do DAG projetado e aos executados. O REPORT factual lista o harness e o modelo reais de cada nó. Um nó em harness não claude passa no gate com o arquivo do `context_out` presente.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * Pelo menos um harness além do claude CLI passa na bateria nesta máquina (F17 está PARCIAL).
  * O peso de memória de cada harness é mensurável pela bateria (F16 está NÃO REPRODUZIDO).
  * O mimo, o omp e o opencode usam o 9Router como provider configurado, sem comutação (F8).
* **Riscos e Perguntas em Aberto:**
  * Qual harness passa na bateria é desconhecido até a primeira execução.
  * O formato exato do arquivo de perfil nasce na SPEC.
  * A bateria gasta dinheiro. O usuário confirma o teto antes.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura de testes expandida (TDD Red-Green). | Teste novo falhando antes do código | Cada nó de código declara o teste-alvo e a mutação que o faz falhar |
| 3 | O motor executa o roteamento declarado. | Teste de paridade entre a ROUTE, o DAG e a execução | `go test ./internal/orchestrix/...` e `go test ./tools/route2dag/...` |
| 4 | O `legoloop rota` não deixa nó sem roteamento resolvido. | Teste com inventário simulado | `go test ./internal/cli/...` |
| 5 | A bateria homologa e reprova harness por prova. | Teste com harness falso que passa e outro que falha | `go test ./internal/roteamento/...` |
| 6 | O `skill-router` carrega skill sob demanda. | Teste do `skill_search` e do `call_skill` sobre um cofre temporário | `go test ./internal/mcp/...` |
| 7 | O taxlinter não regride. | Contagem de violações abaixo da linha de base | `go run ./tools/ticketgate -tax docs/templates -tax-max 35` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-12_SPEC_ROTEAMENTO_HONESTO.md`](../specs/2026-11-12_SPEC_ROTEAMENTO_HONESTO.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (Stop-Gate) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
