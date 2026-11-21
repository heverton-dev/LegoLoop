# PARTE I — Visão Global e as Três Frentes

## 1. O Problema da Divergência Silenciosa

Sistemas compostos por CLI em Go, servidores MCP, malha de agentes e interface web tendem a degradar nas bordas entre camadas. Testes isolados com dados fictícios escondem quebras reais de caminho, timeouts de primeira varredura e contratos divergentes de navegador.

Na linha de base da auditoria em 14/11/2026, o ecossistema registrou nota média **4,6 / 10**:

| Frente Avaliada                    | Nota Inicial   | Principais Achados na Linha de Base               |
| :--------------------------------- | :------------: | :------------------------------------------------ |
| Arquitetura, Go e Taxonomia        | 4,8            | Stubs em `motion`, violações de `taxlinter` e risco de travessia em endpoints do Studio |
| Canais CLI, Harnesses, MCP e Swarm | 5,2            | Flags divergentes em `swarm`, cliente MCP STDIO com leitura sujeita a corrida e ausência de trava no Board |
| UI/UX Studio Web e Swarm UI        | 3,8            | Erro HTTP 404 de `/favicon.ico`, combos legados em `templates.json` e falhas de contraste e responsividade |

## 2. As Três Frentes e os Sete Princípios de Proteção

O Processo divide a recuperação em três frentes paralelas com fronteiras estritas de escrita:

- **Frente 1 (Arquitetura, Go e Taxonomia):** garante Path Jail com `os.OpenRoot`, validação contra `schemas/route.v2.json`, remoção de stubs e zero violações léxicas.
- **Frente 2 (Canais, CLI, MCP e Swarm):** garante protocolo limpo sobre STDIO, contratos de Board protegidos contra concorrência e paridade de subcomandos CLI.
- **Frente 3 (Studio Web e Swarm UI):** garante renderização limpa em `1440x900`, `1024x768` e `390x844`, temas claro e escuro, preservação de estado no polling e zero erros de console.

## 3. O Quinteto de Governança ISO

Antes de editar código de produção, a Etapa 1 congela cinco artefatos encadeados por assinatura SHA-256 em `docs/`:

1. `INTENT` descreve o propósito e os limites do escopo.
2. `SPEC` fixa critérios binários de aceitação EARS.
3. `PLAN` decompõe os 15 Tickets (`TICK-01.1` a `TICK-07.2`).
4. `ROUTE` amarra dependências, roteamento de modelos certificados (`t-*-128k`) e assinaturas de montante.
5. `DAG` reflete 1:1 o contrato executável para o motor de orquestração.
