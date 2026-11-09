# 🎯 INTENT: Pipeline Universal de Engenharia Reversa Audiovisual (legoloop-video-rea)

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Autor:** Heverton & Ecossistema LegoLoop  
> **Iniciativa:** LegoLoop Autonomous Evolution  
> **Alvo:** internal/videorea, internal/cli, loops/video-rea, docs/books/video-rea, skills/legoloop-video-rea  
> **Sessão:** 217415b3-cbf6-432f-bbe6-b267fa3fddea  
> **Confiança da Entrevista:** 98% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criadores de conteúdo publicam tutoriais avançados de vídeo, motion design e automação em redes sociais (YouTube, Instagram, TikTok), porém o conhecimento prático fica fragmentado em vídeos longos, assets distribuídos em links externos e, frequentemente, com segredos técnicos omitidos (prompts implícitos, ferramentas não reveladas, truques de timeline).
* **Quem Sofre:** Desenvolvedores e engenheiros de IA do LegoLoop que precisam manualmente assistir, transcrever, caçar links, descompactar arquivos, descobrir parâmetros ocultos e converter o método em código.
* **O Porquê:** Criar uma esteira autônoma e determinística de engenharia reversa que recebe a URL do vídeo, extrai todo o método explícito e oculto, curadoria contra as peças locais do cofre, apresenta um Stop-Gate humano e dispara a fabricação completa do novo pipeline com seu livro didático.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** A engenharia reversa de vídeo foi feita de modo ad-hoc no terminal e chat durante a concepção do pipeline `video-motion`.
  * **Amanhã:** O comando nativo `legoloop video-rea "<URL>"` automatiza todo o ciclo com relatório forense, curadoria contra a `legoloop matriz`, Stop-Gate humano (`SG-CURADORIA`) e convocação da fábrica de governança LegoLoop (`INTENT`, `SPEC`, `PLAN`, `ROUTE`, `BOOK`, TDD, DAG e Skill).

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Objetivo 1: Ingestão universal de vídeo via Playwright, Obscura e transcrição literal para YouTube, Instagram e TikTok.
* [x] Objetivo 2: Dissecação forense profunda de assets (`.skill`, `.zip`, presets), prompts explícitos e ocultos, ferramentas implícitas e OCR em telas rápidas.
* [x] Objetivo 3: Curadoria semântica contra a `legoloop matriz` (peças Tier Ouro/Prata do cofre local) e cálculo da Economia de Soberania (R$ 0/mês).
* [x] Objetivo 4: Stop-Gate Humano obrigatório (`SG-CURADORIA`) com dossiê técnico antes da fabricação do pipeline.
* [x] Objetivo 5: Convocação da fábrica de governança gerando a Tríade ISO (`INTENT`, `SPEC`, `PLAN`, `ROUTE`), Golden Fixture (3s) e Livro Didático LegoPress em `docs/books/<slug>/`.
* [x] Objetivo 6: Padronização mandatória de nomes com o prefixo `legoloop-*`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não faremos publicação ou agendamento direto em redes sociais.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/legoloop-video-rea` | Ingestão, forense e curadoria contra a matriz do cofre. |
| `BLOCO-02` | `unit_task` | `internal/videorea` | Implementação do motor forense e curador via TDD. |
| `BLOCO-03` | `script` | `go test ./internal/videorea/...` | Regressão mecânica e integridade de parsing e curadoria. |
| `BLOCO-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória e teste de CLI. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via CLI `legoloop video-rea <URL>` ou chat `/legoloop-video-rea`).
* **Nível de Risco & Blast Radius:** Baixo (geração estéril em disco, hermético).
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem autorizada para APIs públicas de transcrição/visão; scraping furtivo local via Obscura; fallback para modelos locais.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Execução do comando `legoloop video-rea` gerando dossiê forense, aprovação no Stop-Gate, emissão dos artefatos de governança ISO e compilação do livro didático com auditoria de geometria válida.
* **Premissas Assumidas:** O usuário possui acesso à internet ou fixture local; o cofre `~/.aidd/skills` e o inventário de matriz estão disponíveis.
* **Riscos e Perguntas em Aberto:** Vídeos em plataformas fechadas com anti-bot exigem sessão persistente via Obscura Chrome CDP.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | `go test -v ./internal/videorea/...` |
| 3 | Pipeline DAG acíclico validado e compilável no LegoLoop. | DAG validado sem ciclos e com camadas | `legoloop validate loops/video-rea/pipeline.dag.json` |
| 4 | Livro didático compilado em PDF e auditado pelo LegoPress. | Auditoria de geometria com zero cabeçalhos quebrados | `legoloop legopress auditar docs/books/video-rea/*.pdf` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_LEGOLOOP_VIDEO_REA.md`](../specs/2026-11-07_SPEC_LEGOLOOP_VIDEO_REA.md)
* **Gate de Transição:** Este Intent foi aprovado pelo usuário no chat antes da redação da SPEC.
