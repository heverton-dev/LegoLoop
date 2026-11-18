# 01 — Disseção Técnica do Comando Fabricar

> **Alvo:** cmd/legoloop e internal/cli (comando fabricar)
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Conformidade:** ASD-STE100 e Quinteto Canônico LegoLoop

---

## 1. Visão Geral da Arquitetura do Alvo

O comando fabricar é a porta de entrada para geração de novos processos no ecossistema LegoLoop. Ele orquestra o ciclo completo: clarificação socrática (entrevista), leitura da matriz de reputação, fusão e seleção de peças, geração do plano de execução e materialização dos artefatos canônicos.

### Ponto de Entrada no Binário (cmd/legoloop/main.go)
- cmd/legoloop/main.go:13-25: main() captura o contexto com cancelamento por sinal de interrupção (signal.NotifyContext). Despacha comandos específicos (triagem, migrar-processos) e delega os demais para cli.Run(ctx, os.Args[1:], cli.Deps{...}).
- internal/cli/cli.go:110: A tabela de comandos mapeia fabricar para cmdFabricar.

---

## 2. Componentes e Estrutura em internal/cli

A implementação na CLI é dividida em quatro arquivos principais de código e onze arquivos de teste especializados:

1. **internal/cli/fabricar.go (356 linhas):**
   - Estrutura de injeção de dependências: FabricarDeps (linha 23) injeta LLMClient, MatrizFunc, Clock, Harnesses, Link e BrainOpts.
   - Opções de linha de comando: fabricarOpts e newFabricarFlags (linhas 33-65). Suporta flags autônomas (--yes, --sem-entrevista), de modelo (--model, --effort), de escopo (--peca, --category), de caminhos (--vault, --docs) e de saída (--json, --dry-run).
   - Função de entrada: cmdFabricar (linha 83) desvia chamadas de subcomando para cmdEntrevista caso o primeiro argumento seja entrevista.
   - Validação e execução: runFabricar (linha 109) valida o request, escolhe o modelo via brain.Pick, constrói dependências via buildFabricarDeps e executa o motor fabricar.Fabricar.
   - Mensagens canônicas: registrada tabela STE (linhas 294-306) com traduções padronizadas de erros e orientações ao usuário.

2. **internal/cli/fabricar_entrevista.go (214 linhas):**
   - Subcomandos da entrevista: iniciar, responder e intent.
   - Estrutura JSON: passoJSON e espelhoJSON padronizam a troca de mensagens entre o terminal, o Chat e o motor de clarificação.
   - Motor de persistência: armazena sessões em /sessions via interview.Store.

3. **internal/cli/fabricar_interactive.go (135 linhas):**
   - Condução da interface de terminal: prompter (linha 34) utiliza um único bufio.Reader para evitar perda de buffer entre turnos.
   - Limites de tolerância: maxInvalid = 3 aborta execuções após três respostas incorretas consecutivas.
   - Confirmação de confiança e rascunho: checa ConfidenceTarget (95%) e requer confirmação explícita do usuário caso a confiança fique abaixo da meta.

4. **internal/cli/fabricar_report.go (127 linhas):**
   - Apresentação textual do relatório: PrintReportTexto divide a saída nas cinco camadas canônicas:
     - Camada 1: Clarificação (confiança e espelho).
     - Camada 2: Matriz (peças vistas, idade e cortes).
     - Camada 3: Fusão (necessidades, peças aptas e lacunas identificadas).
     - Camada 4: Materialização (artefatos gerados e caminhos).
     - Camada 5: Vínculos (links nos harnesses detectados).

---

## 3. Cobertura e Rastreabilidade de Testes

Os testes da CLI cobrem amplamente todos os caminhos críticos:
- dispatch_fabricar_test.go: isolamento e validação de flags.
- fabricar_binary_test.go: execução ponta a ponta do binário compilado com interatividade mockada.
- fabricar_interactive_test.go: cobertura de fluxos felizes, respostas livres, cancelamentos, timeouts e EOF.
- fabricar_router_test.go: segurança da chave de API e timeout adequado para chamadas remotas ao 9Router.
- fabricar_test.go: testes unitários de integração das opções da CLI com o motor.

Todos os testes unitários executam em 2.8s com 100% de sucesso.
