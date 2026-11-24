---
id: INTENT-intv_1795520712665396700
status: APPROVED
date: 2026-11-24
schema_version: 2
upstream: [INTERVIEW-intv_1795520712665396700]
covers: [Q1]
---

# 🎯 INTENT: Criar o Processo canonizado legoloop-diorama-roteiro (Ciclo 1 do Macroprocesso…

> **Status:** APPROVED  
> **Data:** 2026-11-24  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1795520712665396700  
> **Confiança da Entrevista:** 96% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criar o Processo canonizado legoloop-diorama-roteiro (Ciclo 1 do Macroprocesso legoloop-video-diorama) para executar pre-voo deterministico de ambiente, configuracao validada em projeto.json e marca.json, descoberta de pautas por fontes gratuitas com bloqueio contra temas repetidos em acervo/temas.json e geracao validada de script.json e plano.md com fontes nivel A ou B e separacao entre fato e alegacao.
* **Quem Sofre:** o operador do canal, que recebe roteiros auditáveis sem checar fontes manualmente. Os Processos 2 e 3 do Macroprocesso consomem `script.json` e `plano.md` já validados. O próprio Macroprocesso ganha um acervo de temas sem repetição.
* **O Porquê:** o Ciclo 1 destrava os Ciclos seguintes, porque sem roteiro validado não há vídeo. Hoje não há bloqueio mecânico contra temas repetidos, nem separação rastreável entre fato, evidência preliminar e alegação de terceiro.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista (ver O Problema Real).
  * **Amanhã:** Processo canonizado `legoloop-diorama-roteiro` que executa, em ordem: 1. Etapa 0, configuração guiada: fabrica o `projeto.json` por vídeo, com validação de schema. Campos: `nicho` (texto livre), `perfis_fonte` (`geral` como fallback), `fontes_extras` e `fontes_excluidas` opcionais, `politica_gates` (`humano`, `mecanico` ou `hibrido`, com `gates_humanos`) e `video_referencia` opcional. O `marca.json` é herdado do canal ou criado de template canônico (`legoloop-canal-redc`). 2. Pré-voo determinístico modular, com bloqueio nestes casos: `legoloop version` não retorna exit 0. Catálogo declarativo de fontes por nicho ausente ou JSON inválido. Pasta do projeto ou `acervo/` sem leitura e escrita reais. Nenhum conector compatível com `perfis_fonte` responde na sonda HTTP (timeout ≤ 8s, status 200). `video_referencia` preenchido e `legoloop-video-rea` ausente do Cofre. 3. Validação estrita: schema inválido, campo obrigatório ausente ou inconsistência entre `projeto.json` e `marca.json` bloqueiam. 4. Descoberta de pautas: catálogo declarativo e extensível de conectores gratuitos, filtrado por `perfis_fonte` e ajustado por `fontes_extras` e `fontes_excluidas`. Temas já presentes em `acervo/temas.json` são bloqueados. 5. Escolha da pauta, conforme `politica_gates`. 6. Geração de `script.json` e `plano.md`: cada linha tem natureza epistêmica (`fato_consolidado`, `evidencia_preliminar`, `alegacao_terceiro` ou `narrativa_transicao`). 7. Gate final da Etapa 2: se reprovar, o Processo tenta corrigir até 2 vezes usando só o relatório estruturado do validador (linha, natureza epistêmica e regra violada). Se a falha persistir, ele para, mesmo no modo `mecanico`, e entrega o relatório completo para decisão humana.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Processo canonizado `legoloop-diorama-roteiro` que executa, em ordem: 1. Etapa 0, configuração guiada: fabrica o `projeto.json` por vídeo, com validação de schema. Campos: `nicho` (texto livre), `perfis_fonte` (`geral` como fallback), `fontes_extras` e `fontes_excluidas` opcionais, `politica_gates` (`humano`, `mecanico` ou `hibrido`, com `gates_humanos`) e `video_referencia` opcional. O `marca.json` é herdado do canal ou criado de template canônico (`legoloop-canal-redc`). 2. Pré-voo determinístico modular, com bloqueio nestes casos: `legoloop version` não retorna exit 0. Catálogo declarativo de fontes por nicho ausente ou JSON inválido. Pasta do projeto ou `acervo/` sem leitura e escrita reais. Nenhum conector compatível com `perfis_fonte` responde na sonda HTTP (timeout ≤ 8s, status 200). `video_referencia` preenchido e `legoloop-video-rea` ausente do Cofre. 3. Validação estrita: schema inválido, campo obrigatório ausente ou inconsistência entre `projeto.json` e `marca.json` bloqueiam. 4. Descoberta de pautas: catálogo declarativo e extensível de conectores gratuitos, filtrado por `perfis_fonte` e ajustado por `fontes_extras` e `fontes_excluidas`. Temas já presentes em `acervo/temas.json` são bloqueados. 5. Escolha da pauta, conforme `politica_gates`. 6. Geração de `script.json` e `plano.md`: cada linha tem natureza epistêmica (`fato_consolidado`, `evidencia_preliminar`, `alegacao_terceiro` ou `narrativa_transicao`). 7. Gate final da Etapa 2: se reprovar, o Processo tenta corrigir até 2 vezes usando só o relatório estruturado do validador (linha, natureza epistêmica e regra violada). Se a falha persistir, ele para, mesmo no modo `mecanico`, e entrega o relatório completo para decisão humana.
* [ ] Objetivo 2: critérios de sucesso:
  * Um vídeo com `projeto.json` e `marca.json` válidos percorre o Processo e gera `script.json` e `plano.md` aprovados no Gate final.
  * Cada linha de `script.json` cumpre a regra da sua natureza:
  * `fato_consolidado`: 1 fonte A ou 2 fontes B.
  * `evidencia_preliminar`: escopo, preprint ou amostra declarados.
  * `alegacao_terceiro`: atribuição a quem afirmou.
  * `narrativa_transicao`: gancho ou pergunta.
  * O tema só é gravado em `acervo/temas.json` após a aprovação do Gate final.
  * Os três erros bloqueantes e os cinco itens do pré-voo interrompem o Processo antes de qualquer consulta externa. Os itens só de AVISO não interrompem.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Maquete física de diorama.
* ⛔ Tudo que pertence aos Processos 2 e 3: Chrome do Flow na porta 9333, FFmpeg, Chromium, edge-tts e gitleaks. Esses itens geram só AVISO no pré-voo.
* ⛔ Uso do `legoloop-concepcao` para preencher cada vídeo. Ele só concebe o Processo.
* ⛔ Fontes pagas.
* ⛔ Uso de fonte nível C como prova.

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
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  * Um vídeo com `projeto.json` e `marca.json` válidos percorre o Processo e gera `script.json` e `plano.md` aprovados no Gate final.
  * Cada linha de `script.json` cumpre a regra da sua natureza:
  * `fato_consolidado`: 1 fonte A ou 2 fontes B.
  * `evidencia_preliminar`: escopo, preprint ou amostra declarados.
  * `alegacao_terceiro`: atribuição a quem afirmou.
  * `narrativa_transicao`: gancho ou pergunta.
  * O tema só é gravado em `acervo/temas.json` após a aprovação do Gate final.
  * Os três erros bloqueantes e os cinco itens do pré-voo interrompem o Processo antes de qualquer consulta externa. Os itens só de AVISO não interrompem.
* **Restrições Inegociáveis:**
  * Nível A: `.gov.br`, `.gov`, `bcb.gov.br`, `cvm.gov.br`, `sec.gov`, `pubmed`, `doi.org`, `who.int`, `ibge.gov.br`, e releases oficiais de `github.com` no nicho software/open source.
  * Nível B: imprensa de referência declarada (Reuters, AP, BBC, Bloomberg, CNBC, Folha, Estadão, G1, Valor).
  * Preprints: `arxiv`, `medrxiv` e `biorxiv` exigem a natureza `evidencia_preliminar`.
  * Nível C: regra fail-closed. Qualquer domínio fora do catálogo e fora de `fontes_extras` é nível C e fica bloqueado como prova. Blogs, Wikipedia e agregadores servem só para descoberta.
  * Classificação: declarada no catálogo, por domínio e sufixo institucional, além do nível padrão em conectores primários diretos.
  * Fontes: apenas conectores gratuitos.
  * Validação: universal para qualquer nicho, sem viés jurídico.
  * Gate final: o Processo faz no máximo 2 correções automáticas e para mesmo no modo `mecanico`.
  * Tema em `acervo/temas.json`: nunca é gravado antes da aprovação do Gate final, em nenhum modo.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 96%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | (1) Um vídeo com `projeto.json` e `marca.json` válidos percorre o Processo e gera `script.json` e `plano.md` aprovados no Gate final. (2) Cada linha de `script.json` cumpre a regra da sua natureza: (3) `fato_consolidado`: 1 fonte A ou 2 fontes B. (4) `evidencia_preliminar`: escopo, preprint ou amostra declarados. (5) `alegacao_terceiro`: atribuição a quem afirmou. (6) `narrativa_transicao`: gancho ou pergunta. (7) O tema só é gravado em `acervo/temas.json` após a aprovação do Gate final. (8) Os três erros bloqueantes e os cinco itens do pré-voo interrompem o Processo antes de qualquer consulta externa. Os itens só de AVISO não interrompem. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md`](../specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 96% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
