# PARTE I — Macro: Filosofia e Triggers

## 1. A Filosofia Lego

O LegoLoop transforma fluxos de software e de inteligência artificial em **blocos modulares componíveis**, governados por um grafo acíclico dirigido (DAG).

Nenhum fluxo é código estático ou monolítico. Qualquer processo pode ser montado, desmontado, inspecionado, refinado ou recombinado sem efeitos colaterais nos demais blocos.

## 2. Os Cinco Níveis de Granularidade

![Os cinco níveis de granularidade, do Mega Pipeline às Nano-Etapas](/img/diagrama_granularidade.svg){width=88%}

| Nível            | Escopo                   | O que se controla                    | Exemplo real                |
| :--------------- | :----------------------- | :----------------------------------- | :-------------------------- |
| **1. Mega Pipeline** | Ciclo de vida do produto | Orquestração e paralelismo de épicos | Mega Pipeline ponta a ponta |
| **2. Sub-Loops**     | Módulos reutilizáveis    | Injeção hierárquica, saída isolada   | Loop de Auditoria 4F        |
| **3. Etapas (Nós)**  | Contratos de I/O         | Troca de peças, dependências         | `e01` -> `e02`                  |
| **4. Micro-Etapas**  | Execução cognitiva       | Modelo de IA, contexto, MCPs         | Redator de testes isolado   |
| **5. Nano-Etapas**   | Invariantes de código    | Asserções, cobertura, gates `SG-*`      | Função Go com até 35 linhas |

## 3. A Tríade de Triggers

O LegoLoop não se limita à linha de comando. Todo verbo do ecossistema aceita **três triggers equivalentes de entrada**, e todos despacham para o mesmo motor canônico.

![A Tríade de Triggers alimentando o motor canônico](/img/diagrama_triggers.svg){width=95%}

1. **Terminal (CLI):** comando direto `legoloop <verbo> [flags]`. Indicado para scripts, CI/CD e automação sem interface.
2. **Chat do Agente (Harness):** linguagem natural no editor (Claude Code, Cursor, Windsurf, MIMO, Antigravity) por meio das skills canônicas ou slash commands, como `/legoloop-fabricar` e `/legoloop-listar`. O agente interpreta o pedido e despacha a ação de forma estéril.
3. **LegoLoop Studio (Interface Web):** painel visual aberto com `legoloop studio`. Indicado para inspeção de grafos DAG em tempo real, cartões, modais assistidos e drag-and-drop.
