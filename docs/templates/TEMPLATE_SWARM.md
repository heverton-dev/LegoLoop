# 📐 TEMPLATE: legoloop-swarm

> Gabarito contratual das saídas do processo `legoloop-swarm`.  
> Preencha cada `{{PLACEHOLDER}}`. Gate de fechamento: `!has_pattern('<arquivo>', '\{\{')`.

---

## 1. Etapa `e01`: `out/01_etapa.md` (skill `legoloop-piece-upgrader`)

* **Gate:** `file_exists('out/01_etapa.md') && word_count('out/01_etapa.md') >= 30 && !has_pattern('out/01_etapa.md', '\{\{')`

{{CONTEUDO_ETAPA_1}}

## 2. Etapa `e02`: `out/02_etapa.md` (skill `aidd-piece-upgrader`)

* **Gate:** `file_exists('out/02_etapa.md') && word_count('out/02_etapa.md') >= 30 && !has_pattern('out/02_etapa.md', '\{\{')`

{{CONTEUDO_ETAPA_2}}

## 3. Etapa `e03`: `out/03_etapa.md` (skill `legoloop-mega-pipeline`)

* **Gate:** `file_exists('out/03_etapa.md') && word_count('out/03_etapa.md') >= 30 && !has_pattern('out/03_etapa.md', '\{\{')`

{{CONTEUDO_ETAPA_3}}

## 4. Etapa `e04`: `out/04_etapa.md` (skill `aidd-legoloop-mega-pipeline`)

* **Gate:** `file_exists('out/04_etapa.md') && word_count('out/04_etapa.md') >= 30 && !has_pattern('out/04_etapa.md', '\{\{')`

{{CONTEUDO_ETAPA_4}}

## 5. Etapa `e05`: `out/05_etapa.md` (skill `legoloop-video-master`)

* **Gate:** `file_exists('out/05_etapa.md') && word_count('out/05_etapa.md') >= 30 && !has_pattern('out/05_etapa.md', '\{\{')`

{{CONTEUDO_ETAPA_5}}

