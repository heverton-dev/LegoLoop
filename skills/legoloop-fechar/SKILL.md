---
name: legoloop-fechar
description: Use quando o usuário quiser FECHAR E CANONIZAR UM PROCESSO OU PEÇA APROVADA (gerar manual, compilar livro didático LegoPress, registrar skill universal nos harnesses, validar DAG, recompilar binário LegoLoop, reindexar grafo semântico e emitir relatório de conclusão com telemetria estrita de tempo, recursos e custo), via /legoloop-fechar. NÃO use para criar processos do zero ou para auditorias avulsas.
consumes: Artifact:ProcessSlug
produces: Artifact:CanonizedProcess
artifact: docs/books/<slug>/<slug>.pdf
gate: exit_code == 0
suggested_model: t-chat-128k
---

# 🏁 legoloop-fechar — Processo Universal de Canonização e Fechamento

> **Autor:** LegoLoop  
> **Versão:** 1.1.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Uma implementação de processo, ferramenta ou motor foi aprovada no Stop-Gate Humano e precisa passar pelo fechamento canônico obrigatório em 9 etapas: emissão do Manual Didático em `docs/manuals/`, compilação/auditoria do Livro em `docs/books/<slug>/` via LegoPress, emissão e vinculação universal da Skill em todos os harnesses (`~/.claude/`, `~/.gemini/`, etc.), validação do DAG em `processos/`, recompilação fresca do binário (`bin/legoloop.exe` e `go install`), registro no `MEMORY.md`/`AGENTS.md`, reindexação do grafo semântico (`codebase-memory`), verificação mecânica dos portões de qualidade antes do commit e emissão do Relatório de Conclusão com telemetria estrita de tempo, recursos e custo em `docs/reports/`.  
**Quando NÃO usar:** Para iniciar processos do zero (use `/legoloop-fabricar`), dissecar mídias brutas (use `/legoloop-video-rea`) ou editar livros existentes (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - Identificador ou slug do processo aprovado (ex: `video-motion`, `video-rea` ou nome do processo).
* **Produces (Saídas Mecânicas):**
  - Manual operacional em `docs/manuals/YYYY-MM-DD_MANUAL_FUNCAO_<SLUG>_LEGOLOOP.md`
  - Livro didático compilado e auditado em `docs/books/<slug>/`
  - Skill universal em `skills/legoloop-<slug>/SKILL.md` vinculada aos harnesses do SO
  - DAG acíclico validado em `processos/<slug>/processo.dag.json`
  - Binário atualizado no PATH do sistema (`go install ./cmd/legoloop`)
  - Grafo semântico reindexado
  - Relatório formal com telemetria estrita em `docs/reports/YYYY-MM-DD_REPORT_<SLUG>.md`

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** declarar um processo concluído sem executar as 9 etapas mecânicas de fechamento.
> 2. **NUNCA** comitar no Git sem autorização explícita do operador no Stop-Gate Humano.
> 3. **NUNCA** deixar a skill desvinculada dos harnesses do usuário.
> 4. **NUNCA** aceitar livros com páginas órfãs ou cabeçalhos quebrados na auditoria geométrica.
> 5. **NUNCA** omitir o relatório de conclusão com a telemetria estrita de tempo, tokens e financeiro orçado vs gasto.

---

## 3. Protocolo de Operação Passo a Passo

1. **Execução Direta pela CLI:**
   ```bash
   # Simular e validar sem alterar o disco
   legoloop fechar <slug> --dry-run --json

   # Executar o fechamento completo
   legoloop fechar <slug>
   ```
2. **Pelo Chat do Agente:**
   Acione via slash command passando o slug:
   > `/legoloop-fechar video-motion`

---

## 4. DAG & Nós Canônicos

### Processo Canônico `legoloop-fechar`
| ID do Nó | Nome do Nó | Descrição |
| :--- | :--- | :--- |
| `n01_manual_check` | Auditoria e Geração do Manual Operacional | Verificação e emissão do manual em `docs/manuals/` |
| `n02_diagram_design` | Geração do Diagrama Arquitetural em HTML e SVG | Diagrama visual determinístico com legoloop-diagram-design |
| `n03_book_legopress` | Auditoria Geométrica do Livro LegoPress | Compilação e auditoria geométrica do livro didático |
| `n04_skill_harnesses` | Instalação e Vínculo Universal de Skills nos Harnesses | Registro da skill universal nos harnesses do ecossistema |
| `n05_dag_validation` | Validação Estrutural e Topológica do Grafo DAG | Validação de aciclicidade e schema do DAG |
| `n06_rebuild_binary` | Recompilação e Instalação Fresca do Binário | Build e instalação do binário `legoloop.exe` |
| `n07_governance_sync` | Sincronização da Governança Viva (AGENTS, MEMORY, README) | Atualização dos documentos vivos de governança |
| `n08_report_telemetry` | Emissão de Relatório de Conclusão com Telemetria e Custos | Relatório de conclusão e telemetria estrita em `docs/reports/` |

### Processo Legado / Variável `fechar`
| ID do Nó | Nome do Nó | Descrição |
| :--- | :--- | :--- |
| `etapa-manual` | Manual Operacional da Função | Emissão do manual operacional |
| `etapa-skill` | Skill Canônica nos Harnesses | Instalação da skill nos harnesses |
| `etapa-dag` | Validação Mecânica do Grafo DAG | Verificação mecânica do DAG |
| `etapa-rebuild` | Recompilação e Instalação do Binário | Recompilação do binário |
| `etapa-governanca` | Atualização da Governança Viva | Sincronização da governança |

