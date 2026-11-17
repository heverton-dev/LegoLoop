# 🎯 INTENT-2026-11-01: Refatoração Canônica e Factual do Modal de Peça (#modal-item)

> **Status:** AGUARDANDO_GATE_USUARIO  
> **Data:** 2026-11-01  
> **Origem:** Loop 1 — Concepção & Entrevista Socrática (`/aidd-legoloop-mega-pipeline`)  
> **Alvo:** Modal de Detalhes e Inspeção de Peças no Studio e Inventário (`#modal-item`)

---

## 1. Contexto & Diagnóstico Factual

A engenharia reversa ponta a ponta sobre o elemento `#modal-item` e suas fontes de dados (`internal/graph/ontology.go`, `internal/assets/`, `internal/gate/`, `internal/loop/` e `palette.js`) identificou distorções e omissões críticas:

1. **Achatamento e Mutilação de Dados:**
   * A função `entryFromBlock` em `palette.js` descartava campos essenciais do contrato do bloco (`consumes`, `produces`, `artifact`, `gate`, `cost`, `suggested_model`).
   * Como consequência, o modal caía em fallbacks incorretos (ex: o gatilho `github-board-trigger` exibia entrada `"Briefing / Requisitos"` e saída `"Report"`, quando de fato sua entrada é vazia e sua saída é `Spec` gerando `WAVE_PLAN.json`).
2. **Textos Genéricos Hardcoded:**
   * Frases fixas como *"Executa tarefas técnicas determinísticas com isolamento estrito de ambiente."* eram injetadas indistintamente para qualquer peça, mascarando o comportamento real.
3. **Ausência de Especialização por Tipo de Peça:**
   * Peças com naturezas ontológicas distintas (Gatilhos, Controles de Fluxo, Gates, Sub-Loops, MCPs, Hooks) eram tratadas como se fossem todas a mesma coisa.
4. **Gates sem o Tipo de Artefato Auditado:**
   * Não ficava explícito qual artefato o Gate inspeciona (ex: Markdown com critérios binários, Código-Fonte/Build, PR, Relatório de Testes).
5. **Compatibilidade Plana:**
   * As conexões do grafo não distinguiam a tríade de posicionamento da peça (**Antecessora ⬅️, Substituta/Swap 🔄, Sucessora ➡️**).

---

## 2. Princípios Inegociáveis da Refatoração

1. **Zero Textos Inventados:** Todo dado exibido deve provir do manifesto da peça, do catálogo formal ou do arquivo físico no disco.
2. **Especialização Tipada:** O modal adapta layout, diagrama e métricas conforme a família da peça:
   * **Skill / Bloco:** Fluxo IPO, Modelo Sugerido, Custo em Tokens, Artefato Físico e Gate acoplado.
   * **Gatilho (Trigger):** Diagrama de Disparo (Sem entrada ➔ Evento/Filtro ➔ Artefato Inicial).
   * **Controle (Flow):** Grafo de Decisão (Entrada ➔ IF/Loop/Merge ➔ Ramos TRUE/FALSE ou Sincronização).
   * **Gate:** Julgamento Binário (`PASS`/`FAIL`), com **destaque explícito para o tipo de artefato auditado**, assinatura formal e ação em caso de reprovação.
   * **Sub-Loop:** Macro-peça com portas derivadas (`DerivePorts`), contagem de nós internos e artefato terminal.
   * **MCP:** Conector de ferramentas com transporte (`stdio`/`sse`/`http`), host e comando base sanitizado.
   * **Hook / Script:** Interceptador de eventos no SO, runtime e evento disparador.
3. **Tríade de Compatibilidade:**
   * A aba de compatibilidade classifica as conexões em:
     * **⬅️ Antecessoras:** Peças que alimentam a entrada da peça atual.
     * **🔄 Substitutas (Equivalentes/Swap):** Peças com o mesmo contrato de I/O para troca imediata no Canvas.
     * **➡️ Sucessoras:** Peças que consomem a entrega da peça atual.
4. **Fonte Única da Verdade (SSOT):** O modal deve ser um componente desacoplado e reutilizável compartilhado entre o Canvas Studio (`index.html`) e o Inventário (`inventario.html`).

---

## 3. Matriz de Comportamento das 4 Sub-Abas

### Sub-Aba 1: Visão Geral
* **Diagrama Topo Dinâmico:**
  * Bloco/Skill: `Input ➔ Processamento ➔ Output`
  * Gatilho: `Disparo ➔ Evento ➔ Artefato Inicial`
  * Controle: `Entrada ➔ Condição/Barreira ➔ Saídas TRUE/FALSE/Merge`
  * Gate: `Artefato Auditado ➔ Regra Binária ➔ PASS / FAIL`
  * Sub-Loop: `Raiz ➔ Pipeline Interno ➔ Artefato Terminal`
* **Blocos de Propriedades Factuais:**
  * Artefato Produzido / Auditado (com link para convenção de nomes ISO).
  * Gate / Critério de Aceite formal.
  * O que é, o que faz e o que entrega (extraídos da documentação real da peça).
  * Custo em Tokens e Modelo Cognitivo Recomendado.

### Sub-Aba 2: Peças Compatíveis (Tríade Ontológica)
* Filtros por Relação:
  * **[Todas]** · **[⬅️ Antecessoras]** · **[🔄 Substitutas / Swap]** · **[➡️ Sucessoras]**
* Filtros por Categoria:
  * Skills, MCPs, Scripts, Hooks, Gates, Sub-Loops.
* Ações nos cards de compatibilidade:
  * *"Conectar Antes"*, *"Trocar (Swap)"*, *"Conectar Depois"*, *"Adicionar ao Canvas"*.

### Sub-Aba 3: Ocorrências no SO
* Localização física real de cada cópia/instância:
  * Se arquivo no disco (Skill, Script, Loop, Hook): caminho absoluto com botão *"Abrir no Explorer"*.
  * Se MCP: arquivo de configuração de origem (`.claude.json`, `settings.json`, etc.).
  * Se builtin do core (Controle, Gate builtin, Trigger nativo): identificação clara de *"Nativo do Core LegoLoop"*.
* Status de canonização e desduplicação no cofre central.

### Sub-Aba 4: Código / Fonte
* Renderização contextual do fonte:
  * Skills: `SKILL.md` com YAML frontmatter + botão de tradução didática PT-BR.
  * Gatilhos e Controles: Contrato declarativo JSON do bloco.
  * Gates: Implementação/expressão da regra de checagem.
  * Sub-Loops: Manifesto `.loop.json` com DAG de nós internos.
  * MCPs: Bloco de configuração JSON sanitizado.
  * Scripts: Código-fonte real no interpretador nativo.
* Ações: Copiar código, copiar JSON da peça e traduzir.

---

## 4. Critérios de Aceitação Binários (Gates do Loop 1)

- [ ] **G1 (Preservação de Dados da Paleta):** A função de mapeamento de blocos e itens (`palette.js`) preserva 100% dos campos de contrato (`consumes`, `produces`, `artifact`, `gate`, `cost`, `suggested_model`, `meta`) sem descartar dados.
- [ ] **G2 (Fim dos Fallbacks Mentirosos):** O modal exibe apenas dados reais; para gatilhos sem entrada, exibe "Disparo Autônomo" em vez de "Briefing / Requisitos".
- [ ] **G3 (Tipo de Artefato do Gate):** Para qualquer Gate, o modal explicita com clareza a tipagem do artefato auditado (ex: Markdown, Código, PR, Decisão).
- [ ] **G4 (Tríade de Compatibilidade Ativa):** A sub-aba de compatibilidade permite alternar e filtrar entre Antecessoras, Substitutas (Swap) e Sucessoras.
- [ ] **G5 (Zero Textos Hardcoded):** Eliminação completa das descrições estáticas genéricas.
- [ ] **G6 (Unificação sem Duplicação):** O Canvas Studio e o Inventário consomem o mesmo componente sem duplicar HTML ou lógica.

---

## 5. Gate de Aprovação do Usuário
Para avançarmos para o **Loop 2 / Loop 3** (Geração da Rota em `docs/routes/` e da Especificação com Micro-Tickets), você confirma e aprova os termos e critérios deste `INTENT.md`?
