---
name: legoloop-fechar
description: Use quando o usuário quiser FECHAR E CANONIZAR UM PROCESSO, PEÇA OU PIPELINE APROVADO (gerar manual, compilar livro didático LegoPress, registrar skill universal nos harnesses, validar grafo DAG, recompilar binário LegoLoop e reindexar grafo semântico), via /legoloop-fechar. NÃO use para criar processos do zero ou para auditorias avulsas.
consumes: Artifact:ProcessSlug
produces: Artifact:CanonizedProcess
artifact: docs/books/<slug>/<slug>.pdf
gate: exit_code == 0
suggested_model: code-fast
---

# 🏁 legoloop-fechar — Pipeline Universal de Canonização e Fechamento

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Uma implementação de processo, ferramenta, motor ou pipeline foi aprovada no Stop-Gate Humano e precisa passar pelo fechamento canônico obrigatório em 8 etapas: emissão do Manual Didático em `docs/manuals/`, compilação/auditoria do Livro em `docs/books/<slug>/` via LegoPress, emissão e vinculação universal da Skill em todos os harnesses (`~/.claude/`, `~/.gemini/`, etc.), validação do DAG em `loops/`, recompilação fresca do binário (`bin/legoloop.exe` e `go install`), registro no `MEMORY.md`/`AGENTS.md`, reindexação do grafo semântico (`codebase-memory`) e verificação mecânica dos portões de qualidade antes do commit.  
**Quando NÃO usar:** Para iniciar processos do zero (use `/legoloop-fabricar`), dissecar mídias brutas (use `/legoloop-video-rea`) ou editar livros existentes (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - Identificador ou slug do processo aprovado (ex: `video-motion`, `video-rea` ou nome do processo).
* **Produces (Saídas Mecânicas):**
  - Manual operacional em `docs/manuals/YYYY-MM-DD_MANUAL_FUNCAO_<SLUG>_LEGOLOOP.md`
  - Livro didático compilado e auditado em `docs/books/<slug>/`
  - Skill universal em `skills/legoloop-<slug>/SKILL.md` vinculada aos harnesses do SO
  - Grafo DAG acíclico validado em `loops/<slug>/pipeline.dag.json`
  - Binário atualizado no PATH do sistema (`go install ./cmd/legoloop`)
  - Grafo semântico reindexado

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** declarar um processo concluído sem executar as 8 etapas mecânicas de fechamento.
> 2. **NUNCA** comitar no Git sem autorização explícita do operador no Stop-Gate Humano.
> 3. **NUNCA** deixar a skill desvinculada dos harnesses do usuário.
> 4. **NUNCA** aceitar livros com páginas órfãs ou cabeçalhos quebrados na auditoria geométrica.

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
