# 📋 SPEC: Canonização Sob Demanda sem Fricção (Slash Command & Studio UI)

> **Referência:** `docs/intents/INTENT_CANONIZE_SKILL_UI.md`  
> **Status:** Especificação com Critérios Binários (Pass/Fail)  
> **Metodologia:** SDD (Spec-Driven Development) & UX Impeccable  

---

## 1. Arquitetura dos Dois Pontos de Entrada

### A. Slash Command no Chat (`/canonize [alvo]`)
1. **Registro:** Adicionado à lista de `commands` na skill canônica `aidd-legoloop-mega-pipeline`.
2. **Protocolo Conversacional Estrito (3 Etapas):**
   * *Etapa 1 (Disparo):* Usuário digita `/canonize D:\` (ou pasta específica, ou sem argumentos para varredura global).
   * *Etapa 2 (Auditoria & Resumo):* O agente executa `legoloop canonize` com `-dry-run -json` e exibe um card resumo:
     * Alvo auditado, total de duplicatas encontradas e estimativa de espaço recuperável em MB.
     * Pergunta de fechamento: *"Deseja aplicar a desduplicação canônica em [Alvo]? (Sim / Não)"*.
   * *Etapa 3 (Aplicação Segura):* Ao receber o "sim", o agente dispara `-apply -json`, valida `failed == null` e reporta os links criados.

### B. Painel Visual no LegoLoop Studio (`:7878`)
1. **Endpoint REST Backend (`GET /api/canon/status` e `POST /api/canon/apply`):**
   * `GET /api/canon/status`: Retorna o caminho do cofre ativo (`~/.aidd/canonical`), quantidade de arquivos e lista de drives do sistema com status de desduplicação.
   * `POST /api/canon/apply`: Recebe `{ "root": "D:\\caminho" }` e executa a canonização em streaming/JSON.
2. **Componente de UI no Studio:**
   * Botão sutil na barra superior ou no rodapé do Studio: `[ 🗄️ Cofre Canônico ]`.
   * Modal com:
     * Card do Cofre Canônico (status do SSOT).
     * Lista de Drives detectados com botão individual `[ Auditar ]` e `[ Canonizar ]`.

---

## 2. Critérios Binários de Aceitação (Pass/Fail)

### Bloco A: Slash Command & Skill
- [ ] **CA-SKILL-01:** O arquivo `~/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md` declara o comando `/canonize` no frontmatter YAML.
- [ ] **CA-SKILL-02:** A seção operacional de `/canonize` define o protocolo de Dry-Run obrigatório antes do Apply.

### Bloco B: Backend Go do Studio (`internal/studio`)
- [ ] **CA-API-01:** O servidor HTTP expõe a rota `GET /api/canon/status` retornando JSON com cofre e métricas.
- [ ] **CA-API-02:** O servidor HTTP expõe a rota `POST /api/canon/run` aceitando `{ "root": "...", "apply": bool }` e devolvendo o relatório padronizado.
- [ ] **CA-API-03:** Teste unitário em `internal/studio/canon_api_test.go` valida as rotas com resposta HTTP 200 e schema correto.

### Bloco C: Interface do Studio & Segurança
- [ ] **CA-UI-01:** O Studio exibe o botão do Cofre Canônico sem quebrar a proporção do canvas ou da sidebar.
- [ ] **CA-UI-02:** O clique no botão abre o modal com estatísticas reais do cofre.
- [ ] **CA-UI-03:** Zero `innerHTML` desprotegido (auditado per linter de segurança).
- [ ] **CA-INT-01:** `go test ./...` 100% verde e binário `bin/legoloop.exe` recompilado.

---

## 3. Protocolo de Transição

Com a aprovação desta especificação, avançamos para o fatiamento no **`docs/plans/PLAN_CANONIZE_SKILL_UI.md`**.
