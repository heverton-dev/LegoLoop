# Prompt — Teste E2E Real do Processo `legoloop-diorama-ativos` (Ciclo 2)

Copie e cole o bloco abaixo em uma nova sessão:

```text
Quero executar um teste ponta a ponta (E2E) real do Processo canonizado `legoloop-diorama-ativos` (Ciclo 2 — Etapas 3, 4 e 5) usando o projeto já criado pelo Ciclo 1 em `C:\Users\trcnologia\Desktop\LegoLoop\projetos\diorama-demo`.

Regras e escopo do teste E2E:
1. Carregue a skill `legoloop-diorama-ativos` (`call_skill`) e use o binário fresco `.\bin\legoloop.exe diorama-ativos`.
2. Conduza toda a interação e os Stop-Gates (`etapa3_elenco`, `etapa4_quadros_teste`, `etapa5_storyboard`) na Interface A2UI com ciclo contínuo MCP (`sessao_criar` -> `sessao_publicar` -> `sessao_aguardar` no mesmo turno), sem encerrar o turno pedindo aviso no chat.
3. Passo 1 — Pré-voo negativo real: rode `prevoo --projeto projetos/diorama-demo --json` com o `projeto.json` atual (`estilo_visual: "diorama-isometrico"`) e comprove que o gate bloqueia `estilo_visual` inválido.
4. Passo 2 — Ajuste controlado e Pré-voo positivo: pergunte na A2UI qual modo visual testar (`codigo`, `flow` com `flow_projeto` ou `hibrido`), atualize `projeto.json` (fazendo backup antes) e comprove `prevoo` com exit code 0.
5. Passo 3 — Execução com Stop-Gates humanos: como `politica_gates.modo` é `"humano"`, execute `executar --projeto projetos/diorama-demo --json` etapa por etapa, apresentando os artefatos reais gerados (`flow/elenco.json`, `sprites/*.json`, `flow/cenarios.json`, `flow/*.svg`, `cenas.json`) na A2UI para aprovação humana de cada Stop-Gate (`--aprovar etapa3_elenco`, depois `etapa4_quadros_teste`, depois `etapa5_storyboard`).
6. Passo 4 — Idempotência e Gate final: reexecute `executar` com todos os gates aprovados para comprovar idempotência por hash SHA-256 e rode `validar-storyboard --projeto projetos/diorama-demo --json` comprovando exit code 0 e 100% das linhas de `script.json` cobertas com metáfora visual e >= 4 elementos móveis.
7. Ao final, entregue o resumo factual das provas obtidas (comandos, exit codes e arquivos gerados) e informe qualquer lacuna encontrada para correção imediata.
```
