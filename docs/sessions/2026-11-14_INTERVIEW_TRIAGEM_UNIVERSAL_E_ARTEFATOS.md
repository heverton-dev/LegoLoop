# INTERVIEW — Peça de Triagem Universal e Artefatos do Fabricar

- **Data:** 2026-11-14
- **Sessão CLI:** `intv_1794654552891452600`
- **Status:** CONCLUÍDA (validada em prosa com o usuário)

## Decisões Validadas
1. **Função principal do roteador determinístico:** Primeira barreira de entrada nos 3 pontos (CLI, Chat/Skill e Studio) antes de acionar qualquer LLM.
2. **Critério de decisão híbrido:**
   - **Determinístico (zero LLM):** comando/slash/slug explícito, estado de artefato inequívoco no disco ou regra exata com match único.
   - **Probabilístico (LLM):** linguagem natural aberta (ex.: "quero melhorar a interface e experiência da aplicação ABC") ou empate/ambiguidade.
3. **Artefatos do `legoloop fabricar` (13 artefatos):** Remover `template` e gerar `skill`, `loop`, `dag`, `processo`, `intent`, `spec`, `plan`, `route`, `dag_iso`, `triagem`, `manual`, `livro` e `fichas_tecnicas`.
4. **Saneamento:** Adequar todos os processos existentes em `loops/` e remover legados inúteis ou com mau funcionamento.
