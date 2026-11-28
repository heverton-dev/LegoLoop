# Relatório de Entrega: Vídeo Autêntico de Tecnologia e IA em 16:9

**Data:** 27/11/2026  
**Projeto:** `projetos/video-animado/ia-tecnologia-16x9`  
**Canal:** TechIA (`@techiaoficial`)  
**Status:** 100% Concluído e Aprovado nos Gates Mecânicos

---

## 1. Diagnóstico e Causa Raiz Eliminada

1. **Formato Travado em 9:16:**
   - *Causa:* Constante `formatoMotor = "9:16"` em `internal/cli/video_render_etapas.go` e verificação estrita em `internal/video/motor/render.go`.
   - *Correção:* Parametrização dinâmica via `FormatoPrimarioProjeto` e liberação de `16:9` (1920x1080), `9:16` (1080x1920) e `1:1` (1080x1080) com testes herméticos.
   - *Prova:* `ffprobe` em `renders/v1.mp4` confirma 1920x1080 a 30.00 fps.

2. **Herança Indevida de Cenários e Signos Automotivos:**
   - *Causa:* Uso de `loja-fachada` mapeando para `sets.js:191` (`shop`), que continha o letreiro textual "AUTOPEÇAS" e pneus da First Brands.
   - *Correção:* Criação e cadastro no inventário de 4 novos cenários de tecnologia em código vetorial nativo: `sala-servidores` (data center), `escritorio-tech` (calendário NOV 2026 sem letreiros), `sala-diretoria` e `cofre-dados`.

3. **Ausência de Motion Graphics e Textos Foco na Tela:**
   - *Causa:* A struct `Corte` na ponte Go omitia o campo `Layers` e `CompilarCenasAutorais` não estava conectada ao pipeline de renderização.
   - *Correção:* Storyboard autoral completo compilado com textos de alto contraste semânticos ("90 HORAS" em `#FF5A5A`, "INTELIGÊNCIA ARTIFICIAL PROMETEU" em `#FFD24A`), camadas gráficas animadas (`stopwatch`, `doc`, `stamp`, `bars`) e cartões pretos digitados para perguntas de virada.

4. **Lacuna do Kit de Canal e Pós-Produção:**
   - *Causa:* Descarte prévio dos utilitários `canal.py`, `marca.py` e `cortes.py` do projeto de referência `proj--historias-palitos`.
   - *Correção:* Portados integralmente para Go com CLI oficial (`video-distribuicao canal textos`, `video-distribuicao canal imagens`, `video-render marca`, `video-render cortes`).

---

## 2. Artefatos Finais Gerados e Provados

1. **Master Oficial 16:9:**
   - Caminho: `projetos/video-animado/ia-tecnologia-16x9/renders/v1.mp4`
   - Resolução: 1920x1080 a 30 fps (Duração: 130.52s, Tamanho: 40.8 MB)
   - Áudio: Mixagem a -14.0 LUFS e pico -2.0 dBTP com voz Thalita Neural.
   - Marca: Vinhetas de abertura e fechamento, logotipo sincronizado no canto superior.

2. **Cortes Verticais 9:16 para Shorts / Reels / TikTok:**
   - `renders/shorts/short-01.mp4` (1080x1920, 51.06s, 18.4 MB)
   - `renders/shorts/short-02.mp4` (1080x1920, 53.20s, 18.5 MB)
   - `renders/shorts/short-03.mp4` (1080x1920, 36.40s, 11.0 MB)
   - Todos com legendas dinâmicas queimadas e CTA final direcionando para o canal.

3. **Kit de Canal e Redes Sociais:**
   - `canal/textos.json` e `canal/colar.md`: textos completos sem extrapolar os limites de YouTube, Instagram, TikTok e LinkedIn.
   - `canal/imagens/`: 7 imagens oficiais (perfil 1000x1000, banners de 2560x1440 com área segura de 1546x423 px e marca d'água 150x150).
   - `imagens/logo.png`: logotipo vetorial sincronizado para as vinhetas e overlay.

4. **Revisão e Vitrine:**
   - Folhas de contato: `review/sheet_01.jpg`, `sheet_02.jpg`, `sheet_03.jpg`.
   - Parecer técnico: `review/v1_parecer.json` (Veredito: ENVIAR).
   - Capas (Thumbnails 1280x720): `renders/capa_v1_01.png`, `capa_v1_02.png`, `capa_v1_03.png`.
   - Publicação: `publicacao.md` gerado e simulado em dry-run no Postiz com sucesso.

---

## 3. Conformidade e Governança

- `legoloop espinha verificar --json`: `[]` (100% íntegro)
- `go test -count=1 ./internal/video/...`: 100% PASS
- `go run ./tools/funclen internal`: 0 violações (todas as funções <= 35 linhas)
- `gofmt -l .`: saída limpa
