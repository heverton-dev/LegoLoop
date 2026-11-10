---
name: legoloop-video-master
description: Use quando precisar orquestrar ou executar a produção audiovisual completa e soberana com IA em 4 pilares tecnológicos (Motion Graphics 3D com código, Whiteboard SVG no piloto automático, Auto-Editor de Reels e Engenharia Reversa de Canais). NÃO use para edição manual arbitrária de timeline ou para criar livros didáticos avulsos.
compatibility: Windows/Linux, PowerShell 5.1+, Go 1.23+, FFmpeg, Node.js
allowed-tools: read_file write_file edit_file run_command
consumes: Artifact:Spec
produces: Artifact:Code
artifact: processos/legoloop-video-master/processo.dag.json
gate: file_exists('processos/legoloop-video-master/processo.dag.json')
suggested_model: code-pro
metadata:
  version: "2.2.0"
  tier: "ouro"
  category: "build"
---

# 🧩 legoloop-video-master — Mega Processo Soberano de Produção Audiovisual

> **Autor:** LegoLoop REA Engine  
> **Versão:** 2.2.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.2.0  
> **Taxonomia:** Função nativa universal do ecossistema LegoLoop para automação audiovisual fim a fim.

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando o operador desejar produzir vídeos técnicos de alta fidelidade, animações 3D com código, explainers estilo whiteboard ou cortes curtos verticais para redes sociais com custo operacional R$ 0,00/mês.  
**Quando NÃO usar:** Quando a demanda for apenas a dissecação forense de um único vídeo externo (use `/legoloop-video-rea`) ou a compilação pura de um livro didático isolado (use `/legoloop-legopress`).

Esta peça opera estritamente sob os seguintes contratos de entrada e saída:

* **Consumes (Entradas Obrigatórias):**
  - `docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json`
  - `docs/videorea/2026-11-07_CURATION_VIDEO_MASTER.json`
* **Produces (Saídas Mecânicas):**
  - `processos/legoloop-video-master/processo.dag.json`
  - `docs/books/video-master/livro.json`
  - `testdata/fixtures/video_master_fixture_3s.mp4`
* **Target Files (Escopo de escrita permitido / Blast Radius):**
  - `processos/legoloop-video-master/*`
  - `docs/books/video-master/*`
  - `skills/legoloop-video-master/*`
  - `testdata/fixtures/*`

```json
{
  "consumes": ["docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json"],
  "produces": ["processos/legoloop-video-master/processo.dag.json"],
  "gate": "file_exists('processos/legoloop-video-master/processo.dag.json')"
}
```

### Procedimento Determinístico de Execução
1. Inspecionar o dossiê forense e validar os 4 pilares tecnológicos contratados.
2. Garantir que as substituições soberanas de custo zero estejam ativas (Remotion, FFmpeg, Whisper local, Three.js).
3. Executar o DAG validando a ausência de ciclos e portas de saída únicas por Nó.
4. Auditar a integridade geométrica do livro técnico gerado via LegoPress.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
1. **NUNCA** gere código incompleto com `TODO`, `pass`, `// implementar depois` ou mocks simulados (Zero Stubs).
2. **NUNCA** avance para o próximo passo sem a validação mecânica de terminal aprovada (`exit code == 0`).
3. **NUNCA** escreva fora dos Target Files declarados na seção 1 (violação de write boundary).
4. **NUNCA** utilize APIs comerciais com custos recorrentes ou créditos de nuvem; preserve a soberania absoluta de R$ 0,00/mês.
5. **NUNCA** utilize tags HTML ou quebras não escapadas em definições de grafos Mermaid.

---

## 3. Comandos de Verificação & Testes de Terminal

Toda alteração produzida por esta peça deve ser imediatamente verificada pelos comandos executáveis abaixo:

```bash
# Validação do DAG acíclico
legoloop validate processos/legoloop-video-master/processo.dag.json

# Verificação mecânica de integridade dos testes
go test -count=1 ./...

# Auditoria estrutural do livro LegoPress
legoloop legopress auditar docs/books/video-master/livro.json
```

---

## 4. Critério de Parada (Stopping Checklist)

- [x] DAG acíclico validado sem ciclos e com 4 camadas resolvidas.
- [x] Golden Fixture de 3s presente em `testdata/fixtures/video_master_fixture_3s.mp4`.
- [x] Livro didático estruturado e auditado pelo LegoPress.
- [x] Suíte de testes do ecossistema 100% verde.
