# Guia Oficial Completo: Cotas, Impacto Prático e Catálogo de Slash Commands no Antigravity

> **Data:** 2026-11-16  
> **Status:** Referência Técnica e Operacional  
> **Fontes:** Documentação oficial Google AI Studio, Google Cloud e Antigravity CLI  

---

## 1. O que Significam na Prática as 20 Requisições Diárias do Gemini 4 Hydrogen?

No ecossistema de agentes (como o Antigravity CLI):
* **1 Requisição = 1 Turno de Interação ou Chamada de Ferramenta (Tool Call):**
  Quando você envia uma mensagem e o agente precisa inspecionar o código, rodar um teste e editar um arquivo, cada passo consome **1 requisição da API**.
* **Consumo Real de uma Tarefa Simples:**
  1. O agente lê o seu pedido (Requisição 1).
  2. O agente lista os arquivos do diretório (Requisição 2).
  3. O agente lê o arquivo principal (Requisição 3).
  4. O agente aplica a alteração no código (Requisição 4).
  5. O agente roda o teste para validar (Requisição 5).
  6. O agente responde com o resultado final (Requisição 6).
* **Veredito:** Uma única tarefa pequena consome entre **4 a 8 requisições**. Portanto, **20 requisições diárias se esgotam em apenas 2 ou 3 comandos simples (cerca de 5 a 10 minutos de uso)**. Elas servem estritamente para degustação e testes rápidos de sintaxe, e não para trabalho contínuo.

---

## 2. O que Dá para Fazer com as 1.500 Requisições Diárias do Gemini Flash?

Com **1.500 requisições por dia** no Free Tier da API:
* **Capacidade de Trabalho:** Sustenta facilmente entre **150 a 250 tarefas completas de engenharia de software** por dia.
* **Fluxos Complexos:** Permite rodar pipelines de refatoração, auditorias extensas de repositórios inteiros, geração de documentação e centenas de testes unitários sem atingir o teto.
* **Autonomia:** É o modelo ideal para o papel de "trabalhador contínuo" (*worker*), deixando os modelos de fronteira caros (Hydrogen/Argon) apenas para decisões arquiteturais críticas.

---

## 3. Catálogo Completo dos Slash Commands Nativos do Antigravity

Os *Slash Commands* são atalhos diretos na interface do chat que ativam comportamentos especializados e fluxos autônomos avançados:

| Comando | Função Principal | Quando Usar |
| :--- | :--- | :--- |
| `/swarm` ou `/teamwork` | **Modo Swarm (Enxame de Agentes)** | Quando o trabalho pode ser dividido em frentes paralelas independentes ou exige revisão cruzada entre agentes. |
| `/goal` | **Execução em Modo Noturno / Autônomo Longo** | Para tarefas longas e complexas onde o agente não deve parar até concluir e provar 100% da entrega. |
| `/plan` | **Planejamento Passo a Passo Cirúrgico** | Para decompor arquiteturas complexas antes de alterar qualquer código de produção. |
| `/grill-me` | **Entrevista Socrática Alinhadora** | Quando o agente deve entrevistar o usuário para sanar ambiguidades e validar decisões de design antes de codificar. |
| `/boost` | **Modo Raciocínio Profundo Multiperspectiva** | Para problemas de alta complexidade algorítmica, pesquisa profunda e verificação rigorosa. |
| `/compact` | **Compactação e Limpeza de Contexto** | Quando a conversa fica longa e atinge alto consumo de tokens, resumindo o estado sem perder a memória. |
| `/browser` | **Navegação e Automação Web** | Para tarefas que exigem raspar sites, auditar páginas locais no navegador ou buscar documentações externas. |
| `/schedule` | **Agendamento de Cron e Temporizadores** | Para executar instruções em horários programados ou rodar tarefas periódicas em segundo plano. |
| `/learn` | **Fixação e Aprendizado Persistente** | Quando você ensina uma regra ou preferência ao agente e deseja que ela seja gravada para todas as sessões futuras. |

---

## 4. O Modo Swarm e o Orchestrator em Detalhes

### Como Operar o `/swarm`:
1. **Quadro de Mensagens Compartilhado:** O agente coordenador instancia múltiplos trabalhadores efêmeros em conversas separadas.
2. **Sem Poluição de Contexto:** Os trabalhadores executam tarefas pontuais (ex: um pesquisa na web, outro escreve testes unitários) e postam apenas o resultado final no quadro central.
3. **Economia de Recursos:** Reduz o consumo da janela de contexto principal e acelera a entrega em tarefas paralelizáveis.
