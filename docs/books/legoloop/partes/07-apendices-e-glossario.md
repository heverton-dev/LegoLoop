# PARTE VII — Apêndices & Referências

## Apêndice A: Glossário Completo de IA para Leigos

Para que qualquer profissional, independentemente de sua formação técnica, compreenda com clareza as discussões da literatura moderna sobre sistemas inteligentes, reunimos as definições essenciais do ecossistema:

* **Token:** A menor fração de texto processada por um modelo de linguagem. Em português, uma palavra média equivale a aproximadamente 1,5 a 2 tokens. Os provedores de inteligência artificial cobram seus serviços com base na quantidade exata de tokens de entrada e saída.
* **Janela de Contexto (*Context Window*):** O limite máximo de texto que um modelo de IA consegue processar simultaneamente em uma única chamada. Embora janelas modernas alcancem até 1 milhão de tokens, enchê-las degrada a atenção e gera erros de raciocínio.
* **Inchaço de Contexto (*Context Bloat*):** A perda progressiva de velocidade, precisão e foco que ocorre quando uma conversa de chat acumula mensagens antigas, logs e tentativas descartadas.
* **Contexto Estéril:** Princípio fundamental do LegoLoop no qual cada bloco de trabalho é executado em isolamento completo, recebendo apenas o arquivo de entrada estritamente necessário e encerrando sua memória imediatamente após gerar a saída validada.
* **Harness (Ambiente de Execução):** O software de trabalho do desenvolvedor ou operador onde o agente atua. Exemplos comuns incluem Claude Code, Cursor IDE, Windsurf, OpenCode e MIMO Code.
* **MCP (*Model Context Protocol*):** Protocolo padronizado da indústria que permite a modelos de IA conectarem-se a bancos de dados, terminais de comando e ferramentas externas através de interfaces padronizadas.
* **Lazy Discovery (Descoberta Sob Demanda):** Técnica em que ferramentas e habilidades não são injetadas antecipadamente no prompt inicial do modelo, sendo consultadas dinamicamente apenas quando a esteira identifica uma necessidade concreta de negócio.
* **Grafo Acíclico Dirigido (DAG):** Estrutura matemática que organiza tarefas em nós interligados por flechas unidirecionais sem loops infinitos, determinando a ordem correta de execução e viabilizando o paralelismo seguro.
* **TDD (*Test-Driven Development*):** Metodologia de desenvolvimento de software onde os testes automatizados que comprovam a falha de um requisito são escritos antes do código de produção.
* **Stop-Gate Humano:** Ponto de parada programado em uma esteira automatizada que exige inspeção, auditoria e autorização explícita de uma pessoa antes que o processo prossiga.
* **Zero Stubs:** Regra de ouro da engenharia que proíbe terminantemente a entrega de funções simuladas, esqueletos vazios ou códigos fictícios contendo comentários do tipo `// implementar depois`.
* **Alucinação:** Fenômeno estatístico no qual um modelo de linguagem gera respostas incorretas, cita fontes inexistentes ou ignora restrições com alta convicção gramatical, provocado por inchaço de contexto ou falta de ancoragem factual.
* **Temperatura e Top-P:** Parâmetros matemáticos que controlam a aleatoriedade das previsões do modelo. Em esteiras de alta precisão do LegoLoop, a temperatura é configurada próxima de zero para garantir determinismo estrito.
* **Injeção de Prompt (*Prompt Injection*):** Vetor de ataque onde um texto malicioso fornecido por um usuário externo tenta burlar as regras de sistema do modelo. O LegoLoop mitiga esse risco isolando dados e regras em nós com contextos estéreis.
* **Barreira de Junção (*Join Barrier*):** Mecanismo de sincronização determinística que aguarda a conclusão com sucesso de todos os nós paralelos concorrentes antes de liberar a etapa seguinte.
* **Algoritmo de Kahn:** Método matemático de ordenação topológica de grafos que detecta dependências circulares e determina a sequência de execução ótima de uma esteira sem loops infinitos.
* **Rollback Atômico:** Capacidade de reverter o sistema operacional ao estado anterior exato à execução da IA caso um teste falhe ou a inspeção humana rejeite a entrega.
* **SHA-256 (Hash Criptográfico):** Algoritmo que calcula uma impressão digital matemática única de 256 bits para o conteúdo de um arquivo, permitindo detectar cópias idênticas e garantir a integridade dos artefatos.
* **JSON Schema:** Contrato formal baseado em padrões da web que define a estrutura de dados, campos obrigatórios e tipos aceitos em entradas e saídas de cada ferramenta.
* **AST (*Abstract Syntax Tree*):** Árvore sintática abstrata que representa a estrutura do código-fonte na memória, permitindo analisar e manipular software sem executá-lo.
* **Hermeticidade:** Princípio onde um teste ou esteira executa de forma 100% isolada em diretórios temporários, sem depender de internet, sem poluir o sistema e sem tocar em dados de outros processos.
* **SSOT (*Single Source of Truth*):** A prática de manter uma única fonte canônica da verdade para todas as ferramentas e regras de governança, evitando discrepâncias entre múltiplos assistentes.
* **Poka-Yoke (À Prova de Erros):** Filosofia industrial japonesa incorporada no LegoLoop através de portões mecânicos que impedem fisicamente o avanço de uma esteira caso uma pré-condição não seja atendida.

---

## Apêndice B: Manual de Referência Completo da CLI

O executável central do LegoLoop (`legoloop`) oferece comandos determinísticos para gerenciamento completo do ciclo de vida de esteiras e governança de ambientes:

Todos os comandos são executados a partir do binário `legoloop` no terminal. A tabela abaixo sintetiza os principais subcomandos, opções de execução e efeitos no sistema:

| Subcomando         | Opções e Argumentos        | Finalidade Operacional no Ecossistema   |
| :----------------- | :------------------------- | :-------------------------------------- |
| `studio`             | `[--porta 8080] [--abrir]`   | Abre o canvas gráfico interativo no navegador |
| `fabricar`           | `<nome> [--yes] [--model M]` | Dispara entrevista socrática para gerar a Tríade |
| `matriz`             | `[--tier ouro/prata/bronze]` | Audita reputação e testes das ferramentas do cofre |
| `listar`             | `[--json] [--verbose]`       | Exibe o inventário de processos e harnesses ativos |
| `refinar`            | `<nome> [--trocar P]`        | Substitui cirurgicamente uma peça na esteira |
| `remover`            | `<nome> [--force]`           | Desvincula e arquiva o processo com rollback limpo |
| `compor`             | `<alvo> <sub-loop>`          | Encaixa um sub-loop como passo dentro da esteira |
| `sync-harnesses`     | `[--dry-run] [--apply]`      | Conecta assistentes ao Meta-MCP Router sob demanda |
| `stept status`       | `[--harness H] [--json]`     | Audita conformidade do protocolo nos assistentes |
| `stept on`           | `[--harness todos]`          | Ativa a Tríade STE-PT compulsoriamente nas IAs |
| `stept off`          | `[--harness todos]`          | Desativa a exigência formal para exploração livre |

---

## Apêndice C: As Leis Invioláveis do LegoLoop

Para garantir que a plataforma permaneça para sempre um porto seguro de engenharia de software e não degenere no caos dos assistentes comuns, o ecossistema opera sob leis inegociáveis:

1. **Lei #1 (Zero Stubs):** Todo código entregue deve ser real, funcional e testado. É proibido simular comportamentos ou entregar funções vazias.
2. **Lei #2 (Contexto Estéril):** Cada nó da esteira deve nascer limpo e viver apenas o tempo necessário para executar sua tarefa contratada.
3. **Lei #3 (Gramática Fechada):** Nenhuma operação de ciclo de vida pode ser inventada fora dos Cinco Verbos Canônicos (`fabricar`, `matriz`, `listar`, `refinar`, `remover` e `compor`).
4. **Lei #4 (Governança Lazy de MCPs):** Nenhuma ferramenta deve ser injetada no prompt inicial de um assistente sem que haja uma demanda explícita da esteira.
5. **Lei #5 (Transparência STE-PT):** Toda intervenção em código deve apresentar seu diagnóstico prévio e comprovar seu sucesso com evidências computáveis.
6. **Lei #6 (Soberania do Julgamento Humano):** Nenhum stop-gate humano pode ser aprovado automaticamente por uma IA. Decisões estratégicas pertencem exclusivamente a pessoas.

---

## Apêndice D: Rastreabilidade Técnica & Estado Honesto da Plataforma

Seguindo o compromisso inegociável de verdade técnica do LegoLoop, declaramos com absoluta clareza o estado atual de implementação de cada componente no repositório:

* **Motor de Topologia e Grafos (DAG):** Totalmente implementado em `internal/dag/`. Validação formal de dependências cíclicas, ordenação topológica via algoritmo de Kahn e testes unitários cobrindo mais de 90% do código.
* **Fábrica de Processos (`fabricar`):** Operando em produção com geração automática da Tríade Canônica v2 (`INTENT.md`, `SPEC.md`, `PLAN.md` e `ROUTE.json`). Os esquemas formais impedem desvios de sintaxe em saídas de modelos reais.
* **Matriz de Reputação:** Algoritmo ponderado ativo em `internal/matrix/`, categorizando peças em Tiers Ouro, Prata e Bronze a partir de evidências frias do código.
* **Looping Studio Visual:** Servidor nativo em Go (`internal/studio/server.go`) sem dependências externas de pacotes Node pesados em runtime de produção, servindo canvas interativo com suporte a nós, arestas e modais de substituição.
* **Meta-MCP Router e Sincronização:** Validado e operando com sucesso em Claude Code, Cursor IDE, Windsurf e OpenCode, com redução comprovada de até 95% no consumo inicial de tokens.
* **Protocolo STE-PT:** Alternador universal funcional em `internal/stept/`, chaveando diretrizes de transparência em tempo real.

O LegoLoop continuará evoluindo como um projeto aberto, dedicado a todos os profissionais que acreditam que a inteligência artificial deve ser uma esteira de engenharia rigorosa, e não uma loteria de palpites.