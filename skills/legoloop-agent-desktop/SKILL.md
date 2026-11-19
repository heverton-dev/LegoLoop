---
name: legoloop-agent-desktop
description: "Use quando precisar operar ou automatizar aplicações desktop com a peça agent-desktop via árvore de acessibilidade nativa do sistema operacional (Windows, macOS, Linux) sem capturas de tela nem pixels, executando snapshots de acessibilidade, cliques semânticos e digitação em elementos com referências estáveis."
---

# LegoLoop Agent Desktop — Automação de Aplicações Desktop por Acessibilidade

Esta skill orienta a operação da peça `agent-desktop` no ecossistema LegoLoop.
A peça interage diretamente com aplicações desktop sem depender de visão computacional ou reconhecimento de pixels.
Ela inspeciona a árvore de acessibilidade nativa do sistema operacional e devolve referências estáveis em JSON.

## 1. Regras Operacionais

1. Inspecione o estado do adaptador com `agent-desktop status` antes de acionar comandos desktop.
2. No Windows PowerShell, envolva sempre as referências entre aspas simples (exemplo: `'@s8f3k2p9:e1'`) para evitar interpretação de splatting.
3. Utilize a sequência canônica de três passos: inspecionar (`snapshot`), acionar a referência (`click`, `type`, `press`) e confirmar (`snapshot` ou `wait`).
4. Quando a interface mudar, referências anteriores expiram com o erro `STALE_REF`. Gere um novo `snapshot` imediatamente.
5. Prefira a entrega semântica padrão. Adicione `--headed` apenas quando a aplicação exigir foco físico ou movimentação de cursor real.
6. Nunca exponha dados confidenciais ou conteúdos pessoais presentes em telas inspecionadas.

## 2. Comandos Principais da Peça

| Operação | Comando | Finalidade |
| :--- | :--- | :--- |
| Diagnóstico | `agent-desktop status` | Exibe saúde da peça, sistema operacional e permissões |
| Listar Aplicações | `agent-desktop list-apps` | Lista programas com interface gráfica em execução |
| Listar Janelas | `agent-desktop list-windows` | Lista janelas ativas com identificadores e títulos |
| Iniciar Aplicação | `agent-desktop launch <app> --activate` | Inicia o programa e aguarda a janela aparecer |
| Focar Janela | `agent-desktop focus-window --app <Nome>` | Traz a janela para primeiro plano |
| Inspecionar Árvore | `agent-desktop snapshot --app <Nome> -i --compact` | Gera árvore de elementos interativos com referências |
| Acionar Elemento | `agent-desktop click '<ref>'` | Executa ação de clique semântico no elemento |
| Digitar Texto | `agent-desktop type '<ref>' "<texto>"` | Insere texto no campo selecionado |
| Pressionar Teclas | `agent-desktop press "<tecla>"` | Envia atalhos de teclado como enter, esc ou ctrl+c |
| Aguardar Condição | `agent-desktop wait --element '<ref>' --timeout 5000` | Aguarda prontidão do elemento sem pausas cegas |
| Área de Transferência | `agent-desktop clipboard-get` | Lê texto na área de transferência |
| Encerrar Aplicação | `agent-desktop close-app --app <Nome>` | Fecha a aplicação de forma ordenada |

## 3. Modo Servidor MCP

A peça `agent-desktop` opera como servidor de ferramentas MCP via entrada e saída padrão:
```bash
agent-desktop mcp
```
Esse comando publica todas as operações de automação desktop para consumo direto por agentes conectados.
