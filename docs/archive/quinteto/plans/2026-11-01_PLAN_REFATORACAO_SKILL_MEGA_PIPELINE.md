# 🗺️ PLAN: Refatoração da Skill `aidd-legoloop-mega-pipeline`

> **Status:** Proposta de Plano Hierárquico Completo para Aprovação Humana  
> **Data:** 2026-11-01  
> **Alvo:** `C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`  
> **Especificação Base:** [docs/specs/2026-11-01_SPEC_REFATORACAO_SKILL_MEGA_PIPELINE.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-01_SPEC_REFATORACAO_SKILL_MEGA_PIPELINE.md)

---

## 🏛️ ÉPICO 01: Reestruturação do Cabeçalho, Leis Fundamentais e Blindagem Anti-Afobação
> **Objetivo:** Estabelecer a base doutrinária da skill, blindando o comportamento do agente contra afobação cognitiva e fixando as 7 Leis Invioláveis.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-01.1: Atualização do Frontmatter e Redação das 7 Leis Invioláveis [CONCLUÍDO]
* **Tipo:** Micro-ticket (Engenharia de Prompt / Governança)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Frontmatter YAML completo (`name`, `description`, `commands`).
  2. Nova Seção 1 com as 7 Leis Fundamentais:
     * Lei 1: Zero Context Bloat (Orchestrix Estéril).
     * Lei 2: Stop-Gates de Auditoria Humana Obrigatórios (Anti-Afobação).
     * Lei 3: Plano Único e Hierárquico (Épicos ➔ Micro/Nano-Tickets).
     * Lei 4: Provider-Awareness & Roteamento Cognitivo com Clamping.
     * Lei 5: TDD Estrito & Zero Stubs.
     * Lei 6: Nomenclatura ISO (`YYYY-MM-DD_*`).
     * Lei 7: Idioma PT-BR & Estilo Silencioso (Caveman).

### 🔸 Nano-ticket TICK-01.2: Cláusula Rígida de Sanção Anti-Afobação (Rule 10 & Stop-Gates) [CONCLUÍDO]
* **Tipo:** Nano-ticket (Polimento / Guardrail de Sessão)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Bloco em destaque (`> [!CAUTION]`) explicitando que qualquer tentativa de gerar código sem aprovação explícita do plano ou sem auditoria do épico anterior constitui falha primária com abort imediato.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 01
* [x] Frontmatter e 7 Leis conferidos e implementados.
* [x] Cláusula anti-afobação aprovada pelo usuário para desbloqueio do Épico 02.

---

## 🏛️ ÉPICO 02: Matriz Hierárquica de Decomposição e Roteamento com Effort Clamping
> **Objetivo:** Definir com precisão cirúrgica a separação entre Épicos, Micro-tickets e Nano-tickets, parametrizando cotas de tokens e flags de reasoning.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-02.1: Matriz de Níveis com Effort Clamping e Cotas de Tokens [CONCLUÍDO]
* **Tipo:** Micro-ticket (Arquitetura de Roteamento)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Tabela com as 3 camadas hierárquicas:
     * **Épico (Macro):** > 20k tokens, Pro/Sonnet, `--effort high`.
     * **Micro-ticket:** 4k a 20k tokens, Sonnet/Flash, `--effort medium`.
     * **Nano-ticket:** < 4k tokens, Flash-Lite/Haiku, `--effort low`.
  2. Guia de aplicabilidade prática para o agente classificar tarefas automaticamente.

### 🔸 Nano-ticket TICK-02.2: Formalização das Modalidades de Despacho (Sugerido, Manual, 9Router) [CONCLUÍDO]
* **Tipo:** Nano-ticket (Configuração de Governança)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Especificação do comportamento nos modos: Sugerido (Heurística), Manual (Lock) e Autônomo (9Router SLA/Budget com clamping de max_tokens e effort).

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 02
* [x] Tabela de roteamento e flags de effort implementadas no SKILL.md.
* [x] Auditoria e aprovação expressa do usuário para desbloqueio do Épico 03.

---

## 🏛️ ÉPICO 03: Arquitetura dos 5 Sub-Loops, Stop-Gates e Templates Canônicos
> **Objetivo:** Mapear o ciclo de vida completo dos 5 sub-loops, incorporando o diagrama Mermaid atualizado, os gates binários/humanos e os templates padronizados de artefatos.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-03.1: Reengenharia dos 5 Sub-Loops e Diagrama com Stop-Gates [CONCLUÍDO]
* **Tipo:** Micro-ticket (Engenharia de Processos)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Diagrama Mermaid atualizado exibindo os Stop-Gates Humanos explicitamente entre Loops 1, 2/3 e 4.
  2. Detalhamento de cada Sub-Loop: Peças Envolvidas, Artefatos Gerados e Gates de Saída (Mecânicos e Humanos).

### 🔹 Micro-ticket TICK-03.2: Definição dos Templates Canônicos Unificados [CONCLUÍDO]
* **Tipo:** Micro-ticket (Padronização Documental)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Template oficial do `INTENT.md` (com sessão socrática).
  2. Template oficial da `SPEC.md` (com critérios G1..Gn).
  3. Template oficial do `PLAN.md` (estrutura única hierárquica com Épicos e Stop-Gates).
  4. Template oficial do `REPORT.md` (fechamento de Loop 5).

### 🔸 Nano-ticket TICK-03.3: Protocolo de Handoff e Tratamento de Falhas com Rollback [CONCLUÍDO]
* **Tipo:** Nano-ticket (Resiliência Operacional)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. JSON padronizado de transição de bloco com metadados de telemetria e effort.
  2. Política estrita de no máximo 3 retries em falha de gate mecânico antes do transbordo humano.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 03
* [x] Fluxo dos 5 sub-loops, stop-gates e templates canônicos implementados.
* [x] Auditoria e aprovação expressa do usuário para desbloqueio do Épico 04.

---

## 🏛️ ÉPICO 04: Auditoria de Conformidade, Homologação e Fechamento
> **Objetivo:** Verificar a conformidade do arquivo `SKILL.md` contra todos os critérios de aceite (G1 a G6), validar lint e emitir relatório de homologação.  
> **Modelo Recomendado:** `Claude 3.5 Sonnet` ou `Gemini 2.0 Flash` (`--effort medium`)

### 🔹 Micro-ticket TICK-04.1: Auditoria Binária dos Critérios G1 a G6 [CONCLUÍDO]
* **Tipo:** Micro-ticket (Validação de Conformidade)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Check individual de cada critério da SPEC (G1 a G6 validados como PASS).
  2. Zero inconsistências no documento.

### 🔸 Nano-ticket TICK-04.2: Emissão do Relatório de Conclusão e Fechamento do Loop 5 [CONCLUÍDO]
* **Tipo:** Nano-ticket (Documentação e Telemetria)
* **Status:** ✅ Concluído
* **Arquivo Alvo:** [`docs/reports/2026-11-01_REPORT_REFATORACAO_SKILL_MEGA_PIPELINE.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-01_REPORT_REFATORACAO_SKILL_MEGA_PIPELINE.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Relatório factual ISO resumindo as alterações e atestando a blindagem da skill.

### 🛑 STOP-GATE DE AUDITORIA HUMANA FINAL — HOMOLOGAÇÃO COMPLETA
* [x] Auditoria final de todos os 4 épicos realizada com sucesso.
* [x] Relatório de homologação emitido.
* [ ] Auditoria final do usuário aprovada para início da Evolução do Motor Go.
