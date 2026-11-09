# 📐 SPEC: Target Workspace Isolado por Vídeo (`docs/videorea/<slug>/`)

> **Status:** APPROVED  
> **Data:** 2026-11-08  
> **Autor:** Antigravity  
> **Iniciativa:** I04 — Arquitetura de Workspaces Isolados e Mídia Soberana  
> **Alvo:** `internal/videorea/`, `internal/cli/video_rea.go`, `.gitignore`  
> **Template:** TEMPLATE_SPEC.md v2.0.0  

---

## 1. Sumário Executivo & Decisões Arquiteturais

Esta especificação define o padrão arquitetural de **Target Workspace (Bundle Autocontido)** para o pipeline de engenharia reversa de vídeos (`legoloop-video-rea`).
Cada vídeo analisado recebe um diretório dedicado em `docs/videorea/<slug>/` contendo:
1. `scratch/`: Camada volátil para downloads brutos, áudio PCM e frames de OCR.
2. `dossier.json`: Dossiê forense estruturado com transcrição e prompts.
3. `curation.json`: Relatório de curadoria e soberania com peças do cofre.
4. `governance/`: Tríade ISO (`INTENT.md`, `SPEC.md`, `PLAN.md`, `ROUTE.json`).
5. `book/`: Livro didático LegoPress (`livro.json`, `partes/`).
6. `loop/`: Grafo executável (`pipeline.dag.json`) e Golden Fixture (`fixture_3s.mp4`).
7. `skill/`: Skill canônica instalável (`SKILL.md`).

Além do isolamento, a fábrica executa a **Projeção Canônica**: espelha ou registra o grafo em `loops/<slug>/` e a skill em `skills/<slug>/`, garantindo consumo natural e sem atrito pelo LegoLoop e seus agentes.

---

## 2. Contratos de Dados & I/O

### 2.1 Estrutura de Diretórios do Bundle
```
docs/videorea/<slug>/
├── scratch/
├── dossier.json
├── curation.json
├── governance/
│   ├── INTENT.md
│   ├── SPEC.md
│   ├── PLAN.md
│   └── ROUTE.json
├── book/
│   └── livro.json
├── loop/
│   ├── pipeline.dag.json
│   └── fixture_3s.mp4
└── skill/
    └── SKILL.md
```

### 2.2 Struct TargetBundle em Go
```go
type TargetBundle struct {
    Slug          string `json:"slug"`
    BaseDir       string `json:"base_dir"`
    ScratchDir    string `json:"scratch_dir"`
    DossierPath   string `json:"dossier_path"`
    CurationPath  string `json:"curation_path"`
    GovernanceDir string `json:"governance_dir"`
    BookDir       string `json:"book_dir"`
    LoopDir       string `json:"loop_dir"`
    SkillDir      string `json:"skill_dir"`
}
```

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> 1. **NUNCA** permitir que arquivos dentro de `docs/videorea/**/scratch/` sejam rastreados pelo Git (bloqueio via `.gitignore`).
> 2. **NUNCA** quebrar a retrocompatibilidade do motor global `fechar` e dos harnesses (projeção canônica obrigatória).
> 3. **NUNCA** ultrapassar o limite estrito de 35 linhas por função (`funclen <= 35`).

---

## 4. Costuras de Teste (Test Seams — Pocock)

* **Costura Primária:** Validação da materialização completa da árvore de diretórios via `InitTargetBundle` e `MaterializeBundle`.
* **Costura Secundária:** Verificação da projeção canônica em `loops/` e `skills/` em diretório temporário (`t.TempDir()`).
* **Dublês:** Mídia sintética de 3s gerada deterministicamente via `SynthesizeGoldenFixture`.
