---
id: INTENT-intv_1794760080406643600
status: UNDER_REVIEW
date: 2026-11-15
schema_version: 2
upstream: [INTERVIEW-intv_1794760080406643600]
covers: [Q1]
---

# 🎯 INTENT: Processo que recebe a URL de um video, extrai metadados e audio (yt-dlp, ffmpeg), transcreve, extrai conceitos, processos, ferramentas, pecas e repositorios citados, classifica cada achado como DEVE, PODE ou NAO PODE integrar ao LegoLoop (com porque e como) e grava o REPORT em docs/reports. Origem: sessao antigravity 3ac86cf1 (analise do video da Lauren Tan), repetida 4 vezes a mao em 14-15/11/2026.

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-15  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1794760080406643600  
> **Confiança da Entrevista:** 100% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Processo que recebe a URL de um video, extrai metadados e audio (yt-dlp, ffmpeg), transcreve, extrai conceitos, processos, ferramentas, pecas e repositorios citados, classifica cada achado como DEVE, PODE ou NAO PODE integrar ao LegoLoop (com porque e como) e grava o REPORT em docs/reports. Origem: sessao antigravity 3ac86cf1 (analise do video da Lauren Tan), repetida 4 vezes a mao em 14-15/11/2026.
* **Quem Sofre:** Mantenedores do LegoLoop e o Grafo Semântico (codebase-memory).
* **O Porquê:** Eliminar a execução manual repetitiva (4x) do padrão da sessão antigravity 3ac86cf1.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Processo que recebe a URL de um video, extrai metadados e audio (yt-dlp, ffmpeg), transcreve, extrai conceitos, processos, ferramentas, pecas e repositorios citados, classifica cada achado como DEVE, PODE ou NAO PODE integrar ao LegoLoop (com porque e como) e grava o REPORT em docs/reports. Origem: sessao antigravity 3ac86cf1 (analise do video da Lauren Tan), repetida 4 vezes a mao em 14-15/11/2026.
  * **Amanhã:** Pipeline automatizado que processa URL de vídeo, extrai áudio/metadados, transcreve via whisper-cli.exe e gera relatório técnico estruturado.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Pipeline automatizado que processa URL de vídeo, extrai áudio/metadados, transcreve via whisper-cli.exe e gera relatório técnico estruturado.
* [ ] Objetivo 2: Geração de arquivo Markdown em `docs/reports` contendo Frontmatter YAML, Resumo Executivo, Tabela de Achados e Classificação (DEVE/PODE/NÃO PODE) validada.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Uso do `TEMPLATE_REPORT.md` (específico para conclusão de processos) e fixação de modelos de LLM específicos na implementação.

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
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Geração de arquivo Markdown em `docs/reports` contendo Frontmatter YAML, Resumo Executivo, Tabela de Achados e Classificação (DEVE/PODE/NÃO PODE) validada.
* **Restrições Inegociáveis:**
  * Uso obrigatório de `whisper-cli.exe` (C:\Users\trcnologia\tools\whisper-cpp\Release)
  * conformidade com as 5 Prisões Mecânicas e lógica de extração baseada em `internal/videorea/curator.go`.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Nenhum registrado.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Geração de arquivo Markdown em `docs/reports` contendo Frontmatter YAML, Resumo Executivo, Tabela de Achados e Classificação (DEVE/PODE/NÃO PODE) validada. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-15_SPEC_PROCESSO_QUE_RECEBE_A_URL_DE_U.md`](../specs/2026-11-15_SPEC_PROCESSO_QUE_RECEBE_A_URL_DE_U.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
