# Prompt — Teste E2E Real Encadeado: Ciclo 1 (`legoloop-diorama-roteiro`) + Ciclo 2 (`legoloop-diorama-ativos`)

Copie e cole o bloco abaixo em uma nova sessão:

```text
Quero executar um teste ponta a ponta (E2E) real encadeando o Ciclo 1 (`legoloop-diorama-roteiro`) e o Ciclo 2 (`legoloop-diorama-ativos`) em uma pasta nova de projeto: `projetos/diorama-e2e-p1-p2`.

Regras obrigatórias de condução:
1. Carregue as skills `legoloop-diorama-roteiro` e `legoloop-diorama-ativos` via `call_skill`.
2. Conduza toda a interação e todos os Stop-Gates na Interface A2UI com ciclo contínuo MCP (`sessao_criar` -> `sessao_publicar` -> `sessao_aguardar` no mesmo turno), usando formatação estruturada (subtítulos, listas, negrito e código inline). Nunca encerre o turno pedindo aviso no chat.
3. Use o binário fresco `.\bin\legoloop.exe`.

Roteiro de execução E2E:
- **Fase A — Ciclo 1 (`legoloop-diorama-roteiro`):**
  1. Crie `projetos/diorama-e2e-p1-p2/projeto.json` e `marca.json` com `estilo_visual: "hibrido"`, `politica_gates.modo: "humano"` e nicho escolhido na A2UI.
  2. Rode `legoloop diorama-roteiro validar-config` e `legoloop diorama-roteiro prevoo` (exit code 0).
  3. Rode `legoloop diorama-roteiro coletar`, apresente as pautas na A2UI para escolha humana (Stop-Gate Etapa 1).
  4. Gere `script.json` (incluindo `ativos[]` descritivos com `personagem`, `objeto` e `motion`) e `plano.md` (com coluna `Dispositivo Visual (Diorama)`), valide com `legoloop diorama-roteiro validar-roteiro` e aprove na A2UI (Stop-Gate Etapa 2).

- **Fase B — Ciclo 2 (`legoloop-diorama-ativos`):**
  5. Rode `legoloop diorama-ativos criar-projeto-flow --projeto projetos/diorama-e2e-p1-p2 --json` (conectando automaticamente no Chrome persistente `9222`/`9333` para gravar o UUID `flow_projeto` em `projeto.json`; se a porta estiver fechada, use `abrir-flow` e aguarde o login na A2UI).
  6. Rode `legoloop diorama-ativos prevoo --projeto projetos/diorama-e2e-p1-p2 --json` (10/10 itens `OK`, exit code 0).
  7. Execute `legoloop diorama-ativos executar --projeto projetos/diorama-e2e-p1-p2 --json` etapa por etapa, aprovando na A2UI os 3 Stop-Gates humanos (`etapa3_elenco`, `etapa4_quadros_teste`, `etapa5_storyboard`), comprovando:
     - Folhas mestres de 4 faixas (`Turnaround`, `Expressoes`, `Poses`, `BuildNotes`) em `flow/fichas/`.
     - Sprites recortados com `rembg` + despill em `sprites/`.
     - Cenários reais no Google Flow (`flow/cen_*.png` + variantes `_var2.png`, sem SVG no modo `hibrido`).
     - Idempotência por hash na reexecução e `cenas.json` aprovado com `legoloop diorama-ativos validar-storyboard` (exit code 0).
  8. Ao final, entregue o resumo factual das provas obtidas (comandos, exit codes e tamanhos reais dos arquivos gerados).
```
