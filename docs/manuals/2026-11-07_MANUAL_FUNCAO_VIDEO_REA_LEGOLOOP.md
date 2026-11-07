# 📘 MANUAL OPERACIONAL: FUNÇÃO VIDEO REA NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Pipeline universal de engenharia reversa audiovisual, dissecação forense e fábrica de governança  
> **Fontes de verdade:** `internal/videorea/` · `skills/legoloop-video-rea/SKILL.md` · `loops/video-rea/pipeline.dag.json`

---

## 1. Visão Geral

A função **`legoloop-video-rea`** ingere URLs de vídeos (YouTube, Instagram Reels, TikTok) ou arquivos locais e executa a dissecação técnica forense: extrai prompts declarados e implícitos, faz OCR em trechos ultrarrápidos de código, mapeia as ferramentas proprietárias contra a Matriz LegoLoop (R$ 0,00/mês de custo de licença), impõe o Stop-Gate humano de aprovação e fabrica a Tríade de Governança ISO (`INTENT`, `SPEC`, `PLAN`, `ROUTE`), a Golden Fixture de 3s e o livro didático via LegoPress.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop-video-rea <URL>`
* `/video-rea <URL>`

### 2.2 Pela CLI Mecânica
```bash
# Execução completa de ponta a ponta
legoloop video-rea run "https://www.youtube.com/watch?v=..."

# Etapas modulares
legoloop video-rea ingest --url "<URL>" --out scratch/raw.json
legoloop video-rea forensics --input scratch/raw.json --out scratch/dossier.json
legoloop video-rea curate --dossier scratch/dossier.json --out scratch/curation.json
legoloop video-rea gate --report scratch/curation.json
legoloop video-rea govern --approval scratch/approval.json --out docs/generated/
```
