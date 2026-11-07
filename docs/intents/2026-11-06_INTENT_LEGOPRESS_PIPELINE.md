# 🎯 INTENT: LegoPress™ — Pipeline Industrial e Template Canônico de Livros Técnicos

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-06  
> **Autor:** Heverton & Ecossistema LegoLoop  
> **Alvo:** `internal/legopress/`, `docs/templates/`, `cmd/legoloop/`  
> **Sessão:** 088c6d3f-4908-4fd5-b40e-896ee8fcf556  
> **Confiança da Entrevista:** 96%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A geração de livros técnicos com inteligência artificial tradicional falha por acúmulo de contexto, alucinações, layouts quebrados, textos espremidos em tabelas, cabeçalhos com sobreposição e quebras de página aleatórias. Durante a criação real da obra *LegoLoop: A Revolução da IA Modular em Peças Simples*, foi necessário um processo manual e iterativo de calibração geométrica (formato editorial 170x240mm, tipografia Typst, tabelas de linha única, SVGs com alto contraste e eliminação de páginas em branco). Falta transformar esse método em uma esteira de engenharia reutilizável, automatizada e auditável.
* **Quem Sofre:** Engenheiros de plataforma, autores técnicos, product managers e líderes de tecnologia que precisam produzir documentações de alto padrão ou publicar livros físicos/digitais sem retrabalho tipográfico nem explosão de tokens.
* **O Porquê:** Sistematizar a inteligência adquirida em um pipeline oficial do LegoLoop garante que qualquer livro ou apostila didática nasça com padrão editorial Stripe Press / O'Reilly em minutos, com determinismo matemático e zero esforço de diagramação manual.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** A geração de livros exige invocar scripts dispersos (`livro.py`), ajustar manualmente quebras no Typst e rodar inspeções avulsas no PyMuPDF para caçar bugs visuais.
  * **Amanhã:** O ecossistema LegoLoop possui o comando canônico `legoloop legopress` (e seu template visual no Studio) que executa os 4 modos operacionais (*Criar*, *Editar*, *Atualizar*, *Adequar*) de forma estéril por bloco, validando cada etapa por testes mecânicos.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] **Modo 1 — Criação (Greenfield):** Esteira ponta a ponta desde o briefing estruturado até o particionamento em Macro/Meso/Micro, geração de SVGs arquiteturais e compilação de PDF 170x240mm.
* [ ] **Modo 2 — Edição & Correção Cirúrgica (Hotfix):** Motor automatizado de inspeção visual (PyMuPDF) que caça e repara tabelas espremidas, palavras cortadas, colisões em cabeçalhos e páginas em branco sem reescrever o texto íntegro.
* [ ] **Modo 3 — Atualização Viva (Sync):** Analisador de diffs do Git que detecta alterações no código-fonte e regenera apenas as seções e tabelas técnicas defasadas.
* [ ] **Modo 4 — Adequação Didática (Transmutação):** Ingestor de documentações caóticas, READMEs ou transcrições de aulas, transmutando-os para a estrutura didática padrão LegoLoop (Linguagem Acessível, Sem Jargões, Rastreabilidade e Zero IA-slop).
* [ ] **Template Canônico no LegoLoop:** Registro do Grafo Acíclico Dirigido em `docs/templates/TEMPLATE_ROUTE_LEGOPRESS.json` e integração com a CLI e o Looping Studio.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos suporte a formatos proprietários de editoração fechados (ex: Adobe InDesign `.indd`, QuarkXPress). O padrão ouro aberto é Markdown + Typst + Pandoc.
* ⛔ Não alteraremos o núcleo de execução do DAG (`internal/dag/`) ou do Looping Studio; utilizaremos as interfaces já homologadas.
* ⛔ Não criaremos dependências externas em linguagens pesadas além da stdlib de Go, do Python com PyMuPDF para testes de auditoria geométrica, e do binário leve do Typst.
* ⛔ Não permitiremos modelos de linguagem gerando código ou diagramas sem validação estrita por esquemas formais.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `E01-CONCEPCAO` | `subloop` | `interview-me` / `aidd-spec` | Entrevista estruturada para definir público, macro-tema e índice. |
| `E02-PARTICIONAMENTO` | `unit_task` | `internal/legopress/partitioner.go` | Fatiamento da obra na Tríade Macro (Conceito), Meso (Processos) e Micro (Código). |
| `E03-REDAÇÃO-ESTERIL` | `unit_task` | `internal/legopress/writer.go` | Redação estéril nó por nó, sem inchaço de memória nem IA-slop. |
| `E04-DIAGRAMACAO-SVG` | `unit_task` | `internal/legopress/diagram.go` | Geração de diagramas SVG com contraste WCAG AAA e respiro horizontal. |
| `E05-COMPILAÇÃO-TYPST` | `unit_task` | `internal/legopress/compiler.go` | Compilação via Pandoc + Typst no molde 170x240mm e entrelinha 1.5. |
| `E06-AUDITORIA-GEOMETRICA` | `script` | `internal/legopress/audit_geometry.py` | Checagem de tabelas de linha única, ausência de páginas em branco e cabeçalhos. |
| `E07-STOP-GATE-HUMANO` | `stop_gate` | `internal/gate/gate.go` | Parada com renderização de páginas no Studio para aprovação humana final. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via CLI `legoloop legopress [criar|editar|atualizar|adequar]` ou botão no Looping Studio).
* **Nível de Risco & Blast Radius:** Baixo (opera em pastas isoladas de projeto `docs/livro/` ou diretórios dedicados sem tocar em código de infraestrutura).
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem ou Local autorizado. Modelos com clamping de tokens: `--effort high` para arquitetura/índice, `--effort medium` para redação técnica, e `--effort low` para lint e inspeção geométrica.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. Compilação bem-sucedida de um PDF com dimensões físicas 170x240mm.
  2. Script de auditoria mecânica (`audit_geometry.py`) retornando código de saída 0:
     - Zero páginas em branco ou órfãs (< 15 linhas) no miolo.
     - 100% dos cabeçalhos em linha única horizontal (sem quebras de linha).
     - 100% das tabelas com registros em linha única sem hifenização forçada.
     - SVGs com distância mínima garantida de 25pt entre texto e badges de tokens.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  - O motor Typst (CLI ou pacote pip) e Pandoc continuam sendo a fundação de compilação sem dependência de LaTeX.
  - A esteira pode rodar tanto dentro do repositório LegoLoop quanto exportada para qualquer projeto que deseje publicar sua documentação técnica.
* **Riscos e Perguntas em Aberto:**
  - Desempenho na geração de livros muito longos (> 200 páginas): o particionamento em arquivos modulares independentes (`partes/*.md`) anula o inchaço de memória e permite compilação linear sem degradação.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte de testes em Go 100% verde | `go test -count=1 ./internal/legopress/...` |
| 2 | Verificação mecânica de conformidade geométrica do PDF. | Script de inspeção de bounding boxes via PyMuPDF | `python audit_geometry.py livro.pdf` (código 0) |
| 3 | Validação estrutural do grafo da rota no ecossistema. | Comando nativo de validação de DAG | `legoloop validate docs/templates/TEMPLATE_ROUTE_LEGOPRESS.json` |
| 4 | Recompilação do binário da CLI com novo comando exposto. | Binário atualizado com teste de versão e help | `go build -o bin/legoloop.exe ./cmd/legoloop` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md`
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat, antes de qualquer linha de especificação técnica (SPEC) ou plano de tickets (PLAN) ser gerada.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 96% de certeza factual para a esteira LegoPress atende à sua expectativa?*  
*Somente após sua aprovação explícita, avançaremos para o Loop 2 & 3: a redação da SPEC técnica e da Tríade de Planejamento (`SPEC.md` + `PLAN.md` + `ROUTE.json`).*
