# 🎯 INTENT: Pipeline Master de Produção Audiovisual Soberana

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Autor:** LegoLoop REA Engine  
> **Iniciativa:** Audiovisual Soberano 2026  
> **Alvo:** loops/legoloop-video-master, skills/legoloop-video-master, docs/books/video-master  
> **Sessão:** 7a62d284-717f-4c6c-9e24-47b9c787a9be  
> **Confiança da Entrevista:** 98%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A produção de vídeos técnicos, explainers e animações 3D com IA hoje é fragmentada em dezenas de ferramentas SaaS pagas (Higgsfield, Pika, Runway, CapCut Pro, Premiere, Vyond, Fable, Astra) com custos recorrentes que ultrapassam R$ 1.470,00/mês por criador, gerando dependência externa e falta de reprodutibilidade por código.
* **Quem Sofre:** Desenvolvedores, criadores de conteúdo técnico e equipes que precisam gerar dezenas de explainers, reels e animações por código sem intervenção manual de edição.
* **O Porquê:** Unificar 43 fluxos reais analisados por engenharia reversa em um único Mega-Pipeline soberano e determinístico (`legoloop-video-master`), operando com custo marginal zero (R$ 0,00/mês).
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Criação dispersa, dependente de múltiplos SaaS comerciais, renderizadores lentos com filas pagas e edição manual em timeline.
  * **Amanhã:** Um único pipeline composível em 4 fatias verticais (Motion 3D, Whiteboard SVG, Auto-Editor Reels e REA de Canais) acionado via CLI e orquestrado por código determinístico.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Objetivo 1: Dossiê forense estruturado e matriz de soberania com de-para completo dos 43 vídeos submetidos.
* [x] Objetivo 2: Tríade de Governança ISO v2.0.0 (INTENT, SPEC, PLAN, ROUTE) com 4 fatias verticais e zero stubs.
* [x] Objetivo 3: Grafo acíclico DAG e Loop executáveis validados pelo `legoloop validate`, com Golden Fixture de 3s e Livro LegoPress completo.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não integraremos APIs SaaS pagas com cobrança por minuto (Runway/Sora/HeyGen); apenas motores soberanos locais ou gratuitos.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/legoloop-video-master` | Orquestração soberana da esteira audiovisual unificada |
| `BLOCO-02` | `subloop` | `loops/video-motion` | Render de cartões 3D, tipografia Anton e corte de silêncios |
| `BLOCO-03` | `unit_task` | `internal/videorea` | Dossiê forense, curadoria e ponte de governança |
| `BLOCO-04` | `script` | `go test ./...` | Regressão mecânica e integridade dos testes |
| `BLOCO-05` | `stop_gate` | `SG-CURADORIA` | Validação humana de curadoria e orçamento soberano |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via comando CLI `/legoloop-video-rea` ou `legoloop video-rea run`)
* **Nível de Risco & Blast Radius:** Baixo (workspace local, artefatos versionados em docs/ e loops/)
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente Local / Modelos Gratuitos via 9Router; nenhum segredo ou credencial proprietária exposta
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** DAG validado com exit code 0 por `legoloop validate`, Golden Fixture de 3s presente em `testdata/`, Dossiê forense em JSON e Livro Didático completo no padrão LegoPress auditado com zero órfãs
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Os 43 vídeos dividem-se exatamente em 4 pilares tecnológicos complementares; Remotion e FFmpeg local suprem 100% da demanda de edição
* **Riscos e Perguntas em Aberto:** Tempo de compilação de vídeo pesado em hardware modesto (mitigado pelo uso de Golden Fixture de 3s nos testes)

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Grafo DAG compila e é acíclico sem ciclos órfãos. | Validação nativa do CLI | `legoloop validate loops/legoloop-video-master/pipeline.dag.json` |
| 3 | Livro didático completo estruturado segundo o padrão LegoPress. | Auditoria geométrica com zero quebras e zero órfãs | `legoloop legopress auditar docs/books/video-master/2026-11-07_LIVRO_LEGOLOOP_VIDEO_MASTER.pdf` |
| 4 | Economia de soberania calculada e comprovada em R$ 0,00/mês. | Relatório de curadoria JSON | `docs/videorea/2026-11-07_CURATION_VIDEO_MASTER.json` |

---

## 6. Handoff de Governança
* **Destino:** `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md`
* **Condição de Saída:** Aprovação formal registrada no Stop-Gate `SG-CURADORIA` e conformidade total com `TEMPLATE_INTENT.md`.
