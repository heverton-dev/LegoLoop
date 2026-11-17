# 🎯 INTENT: Fábrica Universal de Processos Canônicos (legoloop fabricar)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-03  
> **Autor:** Antigravity AI & Everton  
> **Alvo:** `internal/fabricar/`, `internal/plangen/`, `internal/cli/`, `internal/studio/`, `skills/`, `docs/templates/`  
> **Sessão:** 1da38997-8926-43dc-a9c7-4c885fc12303  
> **Confiança da Entrevista:** 95% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A criação de novos processos, esteiras, loops e subetapas no ecossistema assistido por IA ocorria de forma artesanal e fragmentada. Cada autor ou tarefa impunha convenções arbitrárias, o que gera sobrecarga cognitiva, risco de alucinação de modelos, falta de testes padronizados e incompatibilidade entre os meios de execução (CLI, Chat e Canvas Visual).
* **Quem Sofre:** Autores de processos e os agentes de IA que os executam, a cada novo processo criado.
* **O Porquê:** Unificar a criação de qualquer processo sob uma única esteira de 5 camadas determinísticas. A escala não muda o método: vale para um nano-ticket de 5 linhas e para um mega-pipeline enterprise.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Criação manual ou semi-estruturada de artefatos de planejamento; dispersão de onde salvar e como invocar em cada harness; ausência de garantia de prefixo uniforme.
  * **Amanhã:** A Fábrica de Processos recebe uma intenção via CLI, Chat ou Studio, executa as 5 camadas canônicas (Clarificação ➔ Curadoria SO ➔ Fusão ➔ Tríade ➔ Template), crava o prefixo imutável `legoloop-<kebab-case>` e materializa o pacote interoperável pronto para rodar em qualquer um dos 3 meios.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: **Prefixo fixo e imutável.** Todo processo e peça gerados pela fábrica recebem o prefixo `legoloop-` em `kebab-case` (ex: `legoloop-processo-do-joaozinho`).
* [ ] Objetivo 2: **Esteira das 5 camadas canônicas.**
  * **Camada 1 (Clarificação):** entrevista socrática estilo Addy Osmani, adaptada em PT-BR, até ≥ 95% de confiança.
  * **Camada 2 (Curadoria Multi-Drive, agnóstica a SO):** scanner universal (Windows, Linux e macOS) que varre o disco principal e também todos os discos locais, volumes montados e HDs externos ou pen-drives USB (`D:\`, `E:\`, `F:\`, `/Volumes/*`, `/media/*`, `/mnt/*`), catalogando qualquer peça de IA disponível.
  * **Camada 3 (Fusão Arquitetural):** combinação semântica dos pontos fortes das peças curadas.
  * **Camada 4 (Tríade de Governança):** geração de `SPEC.md`, `PLAN.md` e `ROUTE.json` com isolamento estéril, allowlist de ferramentas e gates booleanos.
  * **Camada 5 (Template SSOT):** geração do gabarito contratual imutável de saída.
* [ ] Objetivo 3: **Interoperabilidade tripartida nativa.** Para a CLI, um grafo `.dag.json` executável por `legoloop run <caminho>`; para o Chat/Harness, um manifesto `SKILL.md`; para o LegoLoop Studio, um grafo `.loop.json` renderizável no Canvas e listado no Catálogo.
* [ ] Objetivo 4: **Subcomando CLI Go** `legoloop fabricar <nome>` com modos interativo e autônomo (`--yes`).
* [ ] Objetivo 5: **Interface no Studio:** modal de Fabricação com seletor de modelo cognitivo e plotagem automática no Canvas.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não alterar os contratos das peças existentes no cofre que já estejam no Padrão Ouro.
* ⛔ Não permitir remoção ou substituição do prefixo `legoloop-`.
* ⛔ Não misturar a *Fábrica de Processos* (design-time) com a *Execução da Aplicação* (run-time).
* ⛔ Não alterar o comando legado `legoloop canonize --root` (desduplicação por SHA-256) nem o comportamento de `legoloop run`.

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar. Tipos permitidos: `skill`, `script`, `subloop`, `unit_task`, `ui_task`, `stop_gate`.

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `unit_task` | Go Core (`internal/fabricar/`) | Motor de validação de nomenclatura e sanitização `legoloop-<kebab-case>`. |
| `BLOCO-02` | `unit_task` | Go Core (`internal/fabricar/`) | Orquestrador das 5 camadas (Entrevista ➔ Curadoria ➔ Fusão ➔ Tríade ➔ Template). |
| `BLOCO-03` | `unit_task` | CLI Go (`internal/cli/`) | Subcomando `legoloop fabricar` com `--yes` e modo interativo. |
| `BLOCO-04` | `skill` | `skills/legoloop-fabricar/SKILL.md` | Skill canônica para uso em Chat/NL, com stop-gates, que invoca `legoloop fabricar`. |
| `BLOCO-05` | `ui_task` | Studio Web (`internal/studio/`) | Modal de Entrevista Socrática com seletor de modelos e plotagem no Canvas. |
| `BLOCO-06` | `script` | `go test ./...` | Suíte completa de testes unitários e regressão TDD. |
| `BLOCO-07` | `unit_task` | `internal/plangen/` | Construtor determinístico `plangen.Render`, para a Tríade reusar o motor de PLAN e ROUTE (zero duplicação). |
| `BLOCO-08` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5) e atualização no PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** Tripartido: `cli: legoloop fabricar`, `chat: /legoloop-fabricar`, `studio: botão + Novo Processo`.  
  *(Define como este fluxo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** `Baixo a Médio` (gravação em diretórios canônicos dedicados e em `docs/`; todo conteúdo de peça varrida é tratado como dado não confiável).
* **Política de Privacidade & Orçamento Cognitivo:** Suporte a modelos locais e de nuvem; effort clamping por camada (Clarificação: High; Curadoria: Low/Local; Tríade: Medium/High).
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. O comando `legoloop fabricar "teste do joaozinho" --yes` gera no disco o pacote tripartido `legoloop-teste-do-joaozinho` com `SKILL.md`, `.loop.json` e `.dag.json`, mais os documentos de governança (INTENT, SPEC, PLAN, ROUTE).
  2. O processo aparece na paleta do Studio Canvas imediatamente após a materialização.
  3. A suíte `go test ./internal/fabricar/...` passa com 100% verde e zero stubs.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * O comando se chama `legoloop fabricar` para não colidir com o `legoloop canonize --root` existente (desduplicação). O pacote Go (`internal/fabricar`), a skill (`legoloop-fabricar`) e a rota (`/api/fabricar/run`) usam o mesmo nome. "Canonizar" é outra função do repositório (desduplicar por SHA-256 e ligar ao cofre, em `internal/canon`) e não nomeia a fábrica.
  * O executável da CLI é um `.dag.json` por caminho (`legoloop run <caminho>`). A forma `legoloop run legoloop-<nome>`, que constava na versão anterior desta INTENT, não é adotada: exigiria alterar o `run`.
  * A skill é ligada às pastas dos harnesses somente com `--link-harness` (opt-in), pois o cofre `~/.aidd/canonical/` não é lido pelo Claude Code.
  * A Tríade (SPEC, PLAN e ROUTE) reaproveita os motores existentes: `specgen` e um construtor determinístico exportado do `plangen` (`plangen.Render`), sem reimplementar o ROUTE.
  * A confiança de 95% da entrevista foi confirmada pelo autor em 2026-11-03.
* **Riscos e Perguntas em Aberto:**
  * **Scanner multi-drive:** varrer mídia removível e unidades de rede abre superfície de prompt injection. Mitigação prevista na SPEC: origem `Untrusted`, leitura sem execução e exclusão da fusão por padrão.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks vazios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Prefixo `legoloop-` garantido por teste unitário em todos os fluxos de criação. | Teste em tabela e fuzz do nome | `go test ./internal/fabricar/...` com `FuzzNormalizeProcessName` sem falha |
| 3 | Binário recompilado e atualizado no PATH com sucesso. | Build e instalação limpos | `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop` com exit code 0 |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md`](../specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
