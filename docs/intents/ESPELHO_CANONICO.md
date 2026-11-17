# 🪞 Espelho Canônico — Mega Plano LegoLoop (Produção)

> **Sessão:** `intv_1794873379697350300` · **Turnos:** 3 · **Confiança:** 94% (premonição com 3 perguntas)
> **Data:** 2026-11-16 · **Processo:** `legoloop-concepcao`

## Campos Consolidados (6)

1. **resultado_esperado:** LegoLoop em estado de produção com arquitetura VSA, saneamento de 1,1 GB do docs/, memória concorrente segura e auditoria SWARM funcional.
2. **quem_se_beneficia:** Desenvolvedores e operadores do ecossistema LegoLoop (CLI, Skill, Studio).
3. **por_que_agora:** Eliminação de dívida técnica crítica, resolução de conflito legal (Lei 15) e necessidade de escala industrial com concorrência.
4. **criterio_de_sucesso:** Saneamento concluído (docs/ reduzido), testes `-race` e `taxlinter` aprovados, auditoria universal validando os 6 verbos sem colisões de PID.
5. **restricoes_inegociaveis:**
   - Sem middleware centralizado para locks; gestão de memória via biblioteca compartilhada com lock por arquivo (Session+PID).
   - Conformidade estrita com a Lei 15 (taxonomia ASD-STE100).
   - Anticorpos: globais para infraestrutura/LLM; locais por processo para domínio.
6. **fora_de_escopo:** Refatoração da interface visual do Studio; criação de novos verbos além dos estabelecidos.

## Decisões dos Turnos

| Turno | Pergunta | Decisão |
| :--- | :--- | :--- |
| 1 | Foco do saneamento de docs/ | Deleção de redundâncias + indexação no grafo semântico; docs/ = referência viva + arquivo morto (`docs/archive/`); SSOT = codebase-memory + MEMORY.md + AGENTS.md + README |
| 2 | Gestão concorrente de memória | **Ressalva do arquiteto adotada:** biblioteca compartilhada no ponto de escrita (`internal/triagem/memoria.go`), lock por arquivo (Session+PID); sem middleware no mcp-router (servidor opcional) |
| 3 | Escopo dos anticorpos | Compartilhamento global para erros de infraestrutura/LLM; local (por processo) para regras de domínio |
