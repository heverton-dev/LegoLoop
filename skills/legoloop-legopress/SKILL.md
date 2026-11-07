---
name: legoloop-legopress
description: Use quando o usuário quiser CRIAR um livro técnico novo, EDITAR ou CORRIGIR o layout de um livro existente, ATUALIZAR capítulos por diffs do código, ou ADEQUAR conteúdos a livros didáticos no padrão determinístico LegoLoop, via /legoloop-legopress. NÃO use para rodar auditorias gerais de código ou criar processos comuns.
consumes: Artifact:Idea | Artifact:Book
produces: Artifact:BookPDF
artifact: docs/books/<slug>/*.pdf
gate: exit_code == 0
suggested_model: sonnet
---

# 📚 legoloop-legopress — Esteira Industrial de Livros Técnicos do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário deseja criar uma nova obra técnica, livro físico ou e-book, corrigir a diagramação visual de um livro gerado, sincronizar partes com o código do Git, ou transformar documentação em material didático de alto padrão editorial (170x240mm).  
**Quando NÃO usar:** Criar processos comuns de código (use `/legoloop-fabricar`), compor processos hierárquicos (use `/legoloop-compor`) ou rodar auditorias gerais da matriz (use `/legoloop-matriz`).

Esta skill conduz a interface no chat e os stop-gates humanos. **Toda a lógica e compilação mora no comando `legoloop legopress`**: particionamento em Macro/Meso/Micro, orquestração Typst + Pandoc, purificação contra IA-slop e auditoria mecânica via PyMuPDF.

* **Consumes (Entradas Obrigatórias):**
  - O pedido do usuário, o manifesto `livro.json`, o arquivo PDF ou os diffs do repositório Git.
* **Produces (Saídas Garantidas):**
  - Pasta própria do livro `docs/books/<slug>/` contendo TUDO dele (ver seção 1.1)
  - Manuscrito estruturado em `docs/books/<slug>/partes/*.md`
  - PDF compilado em padrão editorial 170x240mm (`docs/books/<slug>/<nome_base>.pdf`)
  - Relatório de validação geométrica (`docs/books/<slug>/relatorio_geometria.json`)

### 1.1 Estrutura obrigatória por livro (`docs/books/<slug>/`)

Todo livro nasce dentro da sua própria pasta em `docs/books/`. Nenhum arquivo de livro (fonte, imagem, PDF, preview) vive fora dela e nenhum livro novo é criado em `docs/livro*`.

```
docs/books/<slug>/
  livro.json            manifesto (nome_base, partes[], template: "ativos/livro.typst")
  <nome_base>.md        manuscrito consolidado (gerado)
  <nome_base>.pdf       PDF final (gerado)
  ativos/livro.typst    template Typst do livro (copiar de docs/books/_modelo_livro.typst)
  img/                  capa_arte.jpg, diagramas e capturas (referenciar como `/img/arquivo`)
  partes/NN-*.md        capítulos na ordem do manifesto
  preview/pNN.png       previews das páginas (gerados)
  build/                intermediários (ignorado pelo git)
```

Ferramentas compartilhadas ficam na raiz `docs/books/`: `compilar.sh`, `normalizar_tabelas.py`, `inspecionar.py`, `_modelo_livro.typst`.

---

## 2. REGRAS NEGATIVAS (O que NUNCA fazer)

1. **NUNCA quebrar registros de tabelas em múltiplas linhas:** Todo registro deve caber em uma única linha horizontal no PDF, calibrando a largura e o tamanho da fonte.
2. **NUNCA permitir cabeçalhos quebrados:** O running header no topo da página deve ter rigorosamente uma única linha.
3. **NUNCA permitir páginas em branco ou órfãs:** Quebras de página só acontecem no encerramento de Partes (`level: 1`). Todo o restante do texto corre naturalmente.
4. **NUNCA injetar IA-slop ou saudações prolixas:** Textos didáticos devem ser limpos, objetivos e sem rodeios conversacionais.
5. **NUNCA gerar livro ou ativo fora de `docs/books/<slug>/`:** imagens, partes, manifesto, MD, PDF e preview moram na pasta do próprio livro.
6. **NUNCA compilar sem validar dependências:** Rode sempre `legoloop legopress doctor` antes de disparar grandes compilações.

---

## 3. Protocolo de Operação Passo a Passo

1. **Alinhamento do Modo:** Identifique se a operação é `criar`, `editar`, `atualizar` ou `adequar`.
2. **Execução Estéril via CLI:** Dispare o comando correspondente do binário nativo `legoloop`:
   ```bash
   # Para criar: crie docs/books/<slug>/ (livro.json, partes/, img/, ativos/) e compile:
   legoloop legopress criar "<nome>"
   bash docs/books/compilar.sh <slug>   # gera .md, .pdf e preview DENTRO de docs/books/<slug>/

   # Para auditar a geometria do PDF:
   legoloop legopress auditar docs/books/<slug>/<nome_base>.pdf

   # Para verificar ferramentas:
   legoloop legopress doctor
   ```
3. **Auditoria Visual & Stop-Gate:** Apresente os resultados da auditoria geométrica ao usuário e pause aguardando aprovação explícita antes de qualquer publicação final.
