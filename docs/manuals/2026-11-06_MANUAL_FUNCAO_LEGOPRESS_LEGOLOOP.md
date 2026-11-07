# 📘 MANUAL OPERACIONAL: FUNÇÃO LEGOPRESS NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Esteira industrial de criação, edição, auditoria geométrica, atualização e adequação didática de livros técnicos  
> **Fontes de verdade:** `internal/legopress/` · `internal/cli/legopress.go` · `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md`

---

## 1. Visão Geral

O comando **`legoloop legopress`** (e seu slash command canônico `/legoloop-legopress`) é a esteira oficial de livros técnicos do LegoLoop. Ele orquestra os 4 modos operacionais:
1. **Criar (Greenfield):** Concepção socrática, fatiamento hierárquico na Tríade Macro/Meso/Micro, redação estéril e compilação em PDF 170x240mm.
2. **Editar (Hotfix Cirúrgico):** Diagnóstico e reparo automático de tabelas espremidas, palavras hifenizadas, colisões em cabeçalhos e páginas órfãs.
3. **Atualizar (Sync Vivo):** Leitura de diffs do Git para regenerar exclusivamente os capítulos e tabelas defasados.
4. **Adequar (Transmutação Didática):** Purificação de textos brutos, códigos ou transcrições, removendo IA-slop e convertendo para o estilo didático LegoLoop.

---

## 2. Comandos da CLI

```bash
# Diagnóstico de dependências mecânicas (pandoc e typst)
legoloop legopress doctor

# Auditoria geométrica de layout no PDF gerado
legoloop legopress auditar caminho/do/livro.pdf

# Inicialização de novo livro a partir de briefing
legoloop legopress criar "Nome do Livro"

# Edição e calibração de livro existente
legoloop legopress editar docs/books/<slug>/

# Sincronização viva com diffs do Git
legoloop legopress atualizar docs/books/<slug>/

# Adequação didática sem IA-slop
legoloop legopress adequar conteudo_bruto.md

# Geração de preview em alta resolução (PNG por página)
legoloop legopress preview docs/books/<slug>/
```

---

## 3. Padrões Editoriais Invioláveis

* **Formato Físico Padrão:** 170mm x 240mm (largura útil 134mm, entrelinha 1.45 a 1.5).
* **Tabelas de Linha Única:** Cada registro da tabela ocupa rigorosamente uma única linha, com fontes calibradas (5.9pt a 6.2pt) e `hyphenate: false`.
* **Cabeçalhos sem Quebra:** Running headers ocupam sempre uma única linha horizontal com fonte discreta (5.8pt).
* **Fluxo Corrido:** Quebras de página ocorrem estritamente ao final de Partes (`level: 1`), mantendo o miolo corrido e sem páginas órfãs.
