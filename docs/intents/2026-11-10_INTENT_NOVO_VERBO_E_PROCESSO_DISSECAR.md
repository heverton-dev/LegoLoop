# 🎯 INTENT: Novo verbo e processo dissecar no LegoLoop: inspecao anatomica profunda, engenharia reversa estatica e motion graphics via CLI, Agente IA e Studio

> **Status:** DRAFT  
> **Data:** 2026-11-10  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1794354050170134000  
> **Confiança da Entrevista:** 10% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Novo verbo e processo dissecar no LegoLoop: inspecao anatomica profunda, engenharia reversa estatica e motion graphics via CLI, Agente IA e Studio
* **Quem Sofre:** Engenheiros de software, arquitetos e operadores do LegoLoop que precisam inspecionar, auditar e integrar peças e processos no cofre.
* **O Porquê:** O ecossistema carece de um mecanismo determinístico de inspeção anatômica e visualização gráfica antes de promover peças ao cofre ou rodar pipelines.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Novo verbo e processo dissecar no LegoLoop: inspecao anatomica profunda, engenharia reversa estatica e motion graphics via CLI, Agente IA e Studio
  * **Amanhã:** Implementação industrial do verbo e processo dissecar no LegoLoop (CLI, Agente IA e Studio), com motor estático em internal/dissecar que extrai AST, dependências, variáveis de ambiente, tokens ASD-STE100 e stubs de skills, loops, scripts, hooks e MCPs; gerando ficha anatômica JSON, relatório Markdown em docs/anatomias/, e tríade visual (motion_graph.json para Studio, bundle anatomia.html com Three.js/PixelFX autônomo e anatomia.svg animado).

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Implementação industrial do verbo e processo dissecar no LegoLoop (CLI, Agente IA e Studio), com motor estático em internal/dissecar que extrai AST, dependências, variáveis de ambiente, tokens ASD-STE100 e stubs de skills, loops, scripts, hooks e MCPs; gerando ficha anatômica JSON, relatório Markdown em docs/anatomias/, e tríade visual (motion_graph.json para Studio, bundle anatomia.html com Three.js/PixelFX autônomo e anatomia.svg animado).
* [ ] Objetivo 2: Testes unitários com 100% de cobertura no pacote internal/dissecar (TDD Red-Green), subcomando legoloop dissecar funcional com flags --profundidade, --json e --formato, aba Anatomia no Studio com Canvas 3D e skill canônica legoloop-dissecar registrada.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Execução dinâmica de código inspecionado (sandbox runtime)
* ⛔ renderização de MP4 pelo binário Go
* ⛔ e refatoração de peças legadas fora do escopo da dissecação.

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar. Tipos permitidos: `skill`, `script`, `subloop`, `unit_task`, `ui_task`, `stop_gate`.

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `BLOCO-02` | `unit_task` | Código de produção | Implementação da lógica via TDD Red-Green. |
| `BLOCO-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `BLOCO-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este fluxo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Testes unitários com 100% de cobertura no pacote internal/dissecar (TDD Red-Green), subcomando legoloop dissecar funcional com flags --profundidade, --json e --formato, aba Anatomia no Studio com Canvas 3D e skill canônica legoloop-dissecar registrada.
* **Restrições Inegociáveis:**
  * Obediência estrita às 5 Prisões Mecânicas (Path Jail
  * Schema Validation
  * Exit-Code Gate
  * Sterile Context
  * Test Mutation)
  * Zero Stubs (Lei #1)
  * Graph-First (Lei #14)
  * sem execução de código arbitrário durante análise estática
  * e zero dependências externas novas no go.mod (stdlib pura em Go).
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 10%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Testes unitários com 100% de cobertura no pacote internal/dissecar (TDD Red-Green), subcomando legoloop dissecar funcional com flags --profundidade, --json e --formato, aba Anatomia no Studio com Canvas 3D e skill canônica legoloop-dissecar registrada. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-10_SPEC_NOVO_VERBO_E_PROCESSO_DISSECAR.md`](../specs/2026-11-10_SPEC_NOVO_VERBO_E_PROCESSO_DISSECAR.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 10% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
