# Relatório de Análise Técnica — Rethinking Claude Code Best Practices (Lydia Hallie, Anthropic)

## 1. Metadados do Vídeo

| Campo | Valor |
| :--- | :--- |
| **Título** | Rethinking Claude Code Best Practices — Lydia Hallie, Anthropic |
| **URL** | https://www.youtube.com/watch?v=vlfCjWuv4es |
| **Canal** | AI Engineer |
| **Palestrante** | Lydia Hallie (Member of Technical Staff, Developer Relations para Claude Code, Anthropic) |
| **Duração** | 17:31 |
| **Modelos Citados** | Opus 5.5 (execução/threads) e Fable 5.2 (planejamento, verificação e julgamento) |
| **Data da Análise** | 2026-11-14 |

---

## 2. Resumo Executivo

A palestra apresenta a evolução das boas práticas no uso do Claude Code dentro da Anthropic, que saltou de ~2.500 mudanças/dia em maio para mais de 30.000/dia em setembro. A tese central é a transição do desenvolvedor *"in the loop"* (aprovando cada comando manualmente) para *"manager of dozens of sessions"* (gerente de dezenas de sessões autônomas em paralelo). Para sustentar essa escala sem perder segurança nem qualidade, Lydia detalha cinco pilares: enxugamento radical de instruções (`CLAUDE.md` + `doctor prompt audit`), permissões em camadas (`auto mode` + `deny rules` determinísticas), isolamento real via `Cloud Sessions`, verificação mecânica com **hooks retornando exit code `2`**, e revisão cruzada por sessão limpa (`Cross-Session Review`) antes do merge.

---

## 3. Achados Completos (Conceitos, Processos e Peças/Ferramentas)

| Categoria | Achado | Descrição Técnica |
| :--- | :--- | :--- |
| **Conceito** | Manager of dozens of sessions | O humano deixa de aprovar tool calls unitárias e passa a orquestrar múltiplas sessões paralelas, revisando apenas saídas e PRs. |
| **Conceito** | Enxugamento do `CLAUDE.md` | Instruções antigas (`IMPORTANT!`, `NEVER DO THIS!`, descrição de stack/pastas) degradam modelos novos (geram mais tool calls e piores resultados). O arquivo deve conter **apenas** o que não é visível no código. |
| **Ferramenta** | `doctor prompt audit` | Comando que audita `CLAUDE.md` e `skills` para detectar e remover instruções redundantes ou prejudiciais. |
| **Processo** | `Auto Mode` (classificador de permissões) | Padrão desde agosto; aprova comandos seguros (testes, deps) e bloqueia incertos; configurável em linguagem natural no `settings`. Em estudo de agosto, humanos pegaram 14% dos comandos perigosos contra 89% do `auto mode`. |
| **Peça / Gate** | `Deny Rules` no Harness | Regras determinísticas checadas **antes** do classificador do `auto mode` (ex.: nunca ler `.env`, nunca descartar uncommitted changes). Critério: *"se errar uma vez, dá para reverter?"*. |
| **Conceito / Infra** | Sandbox do SO vs `Cloud Sessions` | Sandbox local não cobre arquivos do Claude, MCPs ou hooks, e child processes (ex.: `npm test`) podem ler `.env`/SSH. `Cloud Sessions` rodam em máquina descartável isolada, sem segredos locais, com rede restrita e segredos injetados fora do alcance do modelo. |
| **Peça / Gate** | Pontos de Verificação & Hooks com **Exit Code `2`** | Hooks determinísticos rodam linter ao editar, testes ao concluir tarefa e suíte ao encerrar sessão. **Crítico:** exit code `1` faz o Claude marcar a tarefa como concluída; **exit code `2`** força o agente a voltar e corrigir. |
| **Processo** | `CI Autofix` & `Merge When Ready` | CI executa as mesmas checagens em todo PR; `CI autofix` corrige falhas automaticamente; `merge when ready` aplica auto-merge por nível de risco (docs: 0 humanos; produto comum: 1; arriscado: 2). |
| **Processo** | `Cross-Session Review` | Todo PR é revisado primeiro por uma **sessão separada** de Claude (contexto limpo) com subgrupo de agentes verificando cada bug antes da revisão humana. |
| **Ferramenta / Processo** | `Claude Code Projects` (Orchestrator + Threads) | Sessão principal (ex.: Fable 5.2) atua apenas como orquestradora (sem escrever código), divide o trabalho em threads (Opus 5.5) e valida o retorno sem compactação manual pelo humano. |

---

## 4. Matriz de Integração ao LegoLoop

### 4.1 DEVEM Integrar (Alta Prioridade — Alinhamento Direto com as 5 Prisões Mecânicas)

| Achado | Por quê | Como Integrar no LegoLoop |
| :--- | :--- | :--- |
| **Hooks com Exit Code `2` como Gate** | Reforça a Prisão Mecânica #3 (*Exit-Code Gate*). Evita falso-positivo de conclusão quando um script falha com `1`. | Padronizar nos hooks e scripts de Gate do LegoLoop o retorno `exit 2` para falhas bloqueantes de validação (linter, `taxlinter`, `funclen`, testes). |
| **`Deny Rules` determinísticas antes de classificadores** | Reforça a Prisão Mecânica #1 (*Path Jail*) e proteção de segredos (`AGENTS.md` Lei #8). | Manter regras `deny` explícitas nos harnesses (`claude`, `agy`, `cursor`) bloqueando leitura de `.env`, chaves SSH e comandos destrutivos antes de qualquer aprovação automática. |
| **Enxugamento de instruções (`doctor prompt audit`)** | Alinha com a Lei #3 (*Economia Severa de Tokens*) e evita degradação de modelos novos por excesso de regras óbvias. | Criar auditoria periódica de `SKILL.md` e `AGENTS.md` para remover instruções que o próprio código já expressa, mantendo apenas invariantes não óbvias. |
| **Cross-Session Review (contexto estéril)** | Espelha a Prisão Mecânica #4 (*Sterile Context*) e Lei #4 (*Orchestrix Cirúrgica*). | Formalizar nó de revisão em sessão limpa (sem histórico de implementação) nos DAGs de código antes do Stop-Gate humano. |

### 4.2 PODEM Integrar (Opcional / Evolução Futura)

| Achado | Por quê | Como Integrar no LegoLoop |
| :--- | :--- | :--- |
| **Roteamento Orquestrador (Fable 5.2) vs Threads (Opus 5.5)** | Complementa a Matriz de Peças e roteamento cognitivo já existente no LegoLoop. | Mapear papéis de orquestração/julgamento para modelos de raciocínio e execução de tickets para modelos rápidos na Matriz de Peças. |
| **`Cloud Sessions` isoladas** | Garante isolamento total de rede e segredos além do `git worktree` local. | Avaliar integração futura de runners efêmeros em container/cloud para execução de DAGs não confiáveis ou em background. |
| **`CI Autofix`** | Acelera correção de falhas mecânicas em PRs. | Adicionar etapa opcional pós-CI que abre ticket de correção automática antes de notificar o Stop-Gate humano. |

### 4.3 NÃO PODEM Integrar (Conflito Direto com Leis Invioláveis)

| Achado | Por quê | Regra Violada no LegoLoop |
| :--- | :--- | :--- |
| **`Merge When Ready` com 0 humanos (auto-merge cego)** | Permite merge sem aprovação humana explícita. | Viola `AGENTS.md` Lei #8 (*Stop-gates `SG-*` são humanos; commit e push só quando pedidos*). Todo merge e entrega no LegoLoop exige Stop-Gate humano. |
| **Substituir Gate determinístico apenas por `Auto Mode`** | O próprio vídeo comprova que o classificador deixa passar ~1 em 10 comandos perigosos (89% de bloqueio). | Viola as 5 Prisões Mecânicas (*Path Jail* e *Exit-Code Gate* exigem garantia determinística de 100%, nunca probabilística). |
