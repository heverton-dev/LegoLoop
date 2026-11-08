---
name: legoloop-video-rea
description: Use quando o usuário quiser FAZER ENGENHARIA REVERSA de um vídeo (YouTube, Instagram Reels, TikTok ou arquivo local), extraindo o método, ferramentas, prompts explícitos/ocultos e gerando a governança canônica e o pipeline do LegoLoop, via /legoloop-video-rea. NÃO use para rodar pipelines já existentes ou compilar livros avulsos.
consumes: Artifact:VideoURL
produces: Artifact:Pipeline
artifact: loops/<slug>/pipeline.dag.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🎥 legoloop-video-rea — Pipeline Universal de Engenharia Reversa Audiovisual

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário fornece uma URL de vídeo (YouTube, Instagram Reels, TikTok) ou vídeo local e deseja extrair todo o conhecimento prático, decodificar prompts explícitos e ocultos, fazer OCR de trechos rápidos de código, curar as ferramentas contra peças locais Tier Ouro do cofre e fabricar um novo pipeline canônico (`legoloop-*`) com seu Livro Didático.  
**Quando NÃO usar:** Executar pipelines que já foram construídos (use `legoloop run` ou `legoloop video-motion`) ou para criar livros de textos arbitrários (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - URL do vídeo ou caminho local de arquivo multimídia.
* **Produces (Saídas Mecânicas):**
  - Target Workspace isolado (`docs/videorea/<slug>/`) com camada volátil (`scratch/`), forense (`dossier.json`), matriz original (`original_matrix.json`), curadoria (`curation.json`), governança local (`governance/`), livro didático (`book/`), pipeline (`loop/`) e skill (`skill/`).
  - Projeção canônica automática na raiz (`loops/<slug>/pipeline.dag.json` e `skills/<slug>/SKILL.md`).
  - Golden Fixture determinística de 3s (`testdata/fixture_3s.mp4`).
  - Livro Didático compilado pelo LegoPress (`docs/books/<slug>/`).

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
1. **NUNCA** construir o pipeline sem aprovação explícita no Stop-Gate de Curadoria (`SG-CURADORIA`).
2. **NUNCA** gerar um novo pipeline sem o prefixo canônico institucional **`legoloop-*`**.
3. **NUNCA** criar um pipeline de vídeo sem a sua Golden Fixture sintética (3s) para testes rápidos locais.
4. **NUNCA** finalizar a entrega sem compilar e auditar a geometria do livro didático via LegoPress.

---

## 3. Protocolo de Operação Passo a Passo

1. **Ingestão e Extração:**
   ```bash
   legoloop video-rea run "<URL>"
   ```
2. **Auditoria no Stop-Gate:** Apresente o relatório com as ferramentas do vídeo, substituições recomendadas e a Economia de Soberania (R$ 0,00).
3. **Fabricação e Validação:** Após a aprovação do operador, o comando gera a Tríade de Governança, a Golden Fixture e o livro didático LegoPress em `docs/books/<slug>/`.
