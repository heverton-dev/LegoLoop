# PARTE II — Meso: Os Cinco Verbos Canônicos

## 6. A Gramática Fechada: Por que Criar uma Linguagem Própria?

Na história da engenharia moderna, grandes saltos de produtividade e confiabilidade sempre coincidiram com a criação de uma linguagem formal com vocabulário restrito e regras inegociáveis. Quando a navegação marítima adotou a rosa dos ventos, a aviação instituiu a fraseologia padronizada de controle de tráfego aéreo e a manufatura desenvolveu as normas técnicas de calibração, o objetivo era sempre o mesmo: **eliminar o ruído interpretativo entre quem dá o comando e quem executa a tarefa**.

No entanto, quando observamos o mercado contemporâneo de agentes de inteligência artificial, deparamo-nos com o oposto do rigor: um festival de scripts improvisados, ferramentas que inventam seus próprios comandos sem qualquer padrão e desenvolvedores que passam horas tentando descobrir como fazer o agente A conversar com o agente B.

O LegoLoop extingue essa babel operacional estabelecendo a **Lei do Vocabulário Fechado**: toda e qualquer operação dentro do ecossistema deve ser compulsoriamente expressa através de apenas **Cinco Verbos Canônicos de Processo**, complementados por um operador hierárquico fundamental:

![O Ciclo de Vida Contínuo dos Cinco Verbos Canônicos no LegoLoop](docs/livro/img/diagrama_ciclo_cinco_verbos.svg)

Essa gramática minimalista cria um modelo mental transparente e compartilhado: tanto o desenvolvedor sênior, quanto o analista de negócios e os próprios modelos de linguagem entendem exatamente o que cada verbo faz, quais são seus efeitos colaterais e como reverter a operação caso algo saia do esperado.

---

## 7. O Verbo Fabricar e a Tríade Canônica

Como transformar um desejo humano abstrato em um sistema de software determinístico que nunca falha e não queima tokens?

A prática convencional de desenvolvimento com IA consiste em abrir o chat do modelo e colar uma descrição longa e genérica, como: *"crie um sistema que leia extratos bancários em PDF, classifique as despesas e gere uma planilha em Excel"*. Essa abordagem ingênua falha invariavelmente: o modelo esquece regras de negócio não declaradas, inventa tratamentos de erro inexistentes e gera um código monolítico impossível de testar.

No LegoLoop, o comando `legoloop fabricar "<nome>"` dispara uma linha de montagem determinística conhecida como a **Tríade Canônica**, composta por quatro artefatos sequenciais padronizados:

![A Tríade Canônica: Da Intenção Pura à Esteira Matemática](docs/livro/img/diagrama_triade_canonico.svg)

### Os Quatro Artefatos Fundamentais

1. **`INTENT.md` (A Intenção Cristalizada):** O sistema não aceita premissas tácitas. Ele inicia uma entrevista socrática de alinhamento com o usuário, formulando perguntas focadas: *Quais formatos de arquivo serão aceitos? Qual o comportamento em caso de CNPJ inválido? Qual o limite aceitável de memória?* O documento resultante consolida o escopo inegociável da demanda.
2. **`SPEC.md` (A Especificação Contratual):** Traduz a intenção em critérios puramente binários. Uma cláusula de especificação no LegoLoop não pode ser interpretativa; ela deve permitir apenas dois estados: **passou** ou **falhou**. É aqui que leis invariáveis são estabelecidas (por exemplo: *"proibido gravar credenciais em disco sem criptografia"*).
3. **`PLAN.md` (O Plano de Micro-Tarefas em TDD):** Decompõe a especificação em tickets verticais atômicos. Sob a disciplina estrita do Desenvolvimento Orientado a Testes (TDD), cada ticket exige que o teste unitário que comprova a falha seja implementado antes de qualquer código de produção. Stop-gates humanos são cravados nos pontos de decisão mais sensíveis.
4. **`ROUTE.json` (O Grafo Acíclico Dirigido — DAG):** É a partitura matemática executada pelas máquinas. O arquivo JSON lista todos os nós de trabalho, suas entradas e saídas esperadas e suas dependências topológicas exatas.

### Exemplo Real de Contrato da Tríade Canônica

Para compreender a precisão cirúrgica exigida pela fábrica do LegoLoop, veja um trecho real do documento de especificação gerado automaticamente durante a fabricação de um validador de cartões de crédito corporativos:

```markdown
# ESPECIFICAÇÃO TÉCNICA: validador-cartao-corporativo

## 1. Critérios Binários de Aceitação (Pass/Fail)
- [ ] O algoritmo de Luhn DEVE validar o número do cartão com código de saída 0.
- [ ] Cartões expirados DEVEM ser rejeitados com erro `ERR_CARTAO_EXPIRADO`.
- [ ] NENHUM número de cartão pode ser gravado em texto puro nos arquivos de log.
- [ ] O tempo total de validação NÃO PODE ultrapassar 50 milissegundos.

## 2. Contrato de Entrada (JSON Schema)
{
  "tipo": "objeto",
  "obrigatorios": ["numero", "validade_mes", "validade_ano"],
  "propriedades": {
    "numero": { "tipo": "string", "tamanho_minimo": 13, "tamanho_maximo": 19 }
  }
}
```

A partir dessa especificação cristalina, o gerador de planos (`plangen`) deriva o plano de tickets com fases Vermelha e Verde estritas:

```markdown
# PLANO DE EXECUÇÃO: tickets-validador-cartao

### Ticket 01: Teste do Algoritmo de Luhn (Fase Vermelha)
- Ação: Criar `internal/cartao/luhn_test.go` com 10 números de cartão válidos e 5 inválidos.
- Condição de Sucesso: O comando `go test ./...` DEVE FALHAR (código diferente de 0).

### Ticket 02: Implementação Mínima (Fase Verde)
- Ação: Criar `internal/cartao/luhn.go` contendo apenas a lógica necessária para o teste passar.
- Condição de Sucesso: O comando `go test ./...` DEVE PASSAR (código igual a 0).
- Stop-Gate Humano: Revisão de segurança criptográfica antes de liberar o merge.
```

```{=typst}
#painel("A Regra Zero Stubs")[
  É terminantemente proibido no ecossistema LegoLoop a criação de funções simuladas contendo comentários como 'TODO', 'pass' ou '// implementar depois'. Todo artefato gerado pela fábrica deve ser funcional, hermético, acompanhado de testes de verdade e pronto para o ambiente de produção.
]
```

---

## 8. O Verbo Matriz: Reputação, Tiers e Qualidade

Em um ecossistema com centenas de ferramentas disponíveis, como garantir que uma esteira financeira ou de auditoria fiscal utilize apenas peças seguras, testadas e que não consumam tokens desnecessariamente?

O LegoLoop não confia em reputações subjetivas ou número de estrelas em plataformas sociais. O motor de **Matriz de Reputação** (`legoloop matriz`) avalia continuamente a qualidade das peças a partir de critérios frios, objetivos e matematicamente auditáveis:

| Critério Automatizado da Peça | O Que é Auditado no Código | Peso no Score |
| :----------------------------------- | :------------------------------------------------ | :-----------: |
| **Cobertura e Integridade de Testes** | Testes unitários reais que falham sob mutação | +0.35 |
| **Zero Stubs (Sem Mocks Simulados)** | Sem funções vazias, TODOs ou simulações fictícias | +0.25 |
| **Conformidade de Esquema Formal** | Validação estrita de contratos via JSON Schema | +0.20 |
| **Economia de Contexto e Tokens** | Prompts concisos e isolamento cirúrgico estéril | +0.20 |

### A Hierarquia dos Tiers de Qualidade

A nota final acumulada varia de 0.00 a 1.00 e define em qual nível a peça opera:

* **Tier Ouro (Score $\ge$ 0.85):** O padrão industrial máximo. Possui bateria completa de testes de regressão, contratos blindados, documentação completa e consome o mínimo indispensável de tokens. Apenas peças Tier Ouro podem ser utilizadas em esteiras que tocam dados sensíveis ou dinheiro.
* **Tier Prata (Score entre 0.60 e 0.84):** Peças funcionais e estáveis, porém com pequenas defasagens — como documentação aprimorável ou prompts que podem ser enxugados via `/refinar`.
* **Tier Bronze (Score $<$ 0.60):** Peças experimentais ou rascunhos em desenvolvimento. O sistema bloqueia seu uso automático e emite alertas caso um usuário tente adicioná-las a uma esteira sem autorização explícita.

### A Telemétrica da Matriz em Pipelines de Integração (CI/CD)

Além da visualização em tabela no terminal, o comando aceita a flag `--json`, permitindo que esteiras automatizadas de integração contínua (CI/CD no GitHub Actions ou GitLab) barrem automaticamente qualquer alteração que reduza a nota de uma peça:

```json
{
  "total_pecas_auditadas": 52,
  "distribuicao_tiers": { "ouro": 41, "prata": 9, "bronze": 2 },
  "defasagens_criticas": 0,
  "conformidade_global": 0.92,
  "bloqueio_ci": false
}
```

Se um desenvolvedor enviar uma alteração de código que introduza uma função vazia com `TODO` (violando a Lei #1 de Zero Stubs), a pontuação daquela peça cai instantaneamente para menos de 0.60, a chave `bloqueio_ci` torna-se `true` e o pipeline de entrega é travado antes de atingir o servidor de produção.

---

## 9. O Verbo Listar: Inventário Vivo e Visibilidade Operacional

Em arquiteturas convencionais de automação, é comum que ferramentas antigas fiquem esquecidas no sistema, acumulando vulnerabilidades e gerando custos residuais.

O comando `legoloop listar` atua como o painel de bordo central do operador, fornecendo em fração de segundo o inventário completo de todos os processos cadastrados no cofre, suas ferramentas vinculadas e o status de conexão em cada um dos ambientes de desenvolvimento (*harnesses*) ativos no computador:

```bash
$ legoloop listar
[PROCESSO] analisador-extratos-v2    | Tier: OURO   | Harnesses: Claude, Cursor, Windsurf
[PROCESSO] auditoria-seguranca-go   | Tier: OURO   | Harnesses: Claude, OpenCode
[PROCESSO] extrator-metadados-pdf   | Tier: PRATA  | Harnesses: Windsurf (Defasagem: v2 disponível)
[PROCESSO] gerador-copy-social      | Tier: BRONZE | Harnesses: Nenhum (Em desenvolvimento)
```

Essa visibilidade em tempo real impede que processos abandonados continuem rodando ou que desenvolvedores utilizem ferramentas desatualizadas sem perceber.

---

## 10. O Verbo Refinar: A Cirurgia Atômica de Peças

O que acontece quando uma biblioteca externa é atualizada, um modelo de IA ganha uma versão mais econômica ou uma nova vulnerabilidade de segurança é descoberta? Em sistemas tradicionais, a equipe é obrigada a reescrever ou reconstruir a esteira inteira, correndo o risco iminente de quebrar integrações que já estavam funcionando perfeitamente.

O LegoLoop resolve esse dilema através do verbo `refinar` (`legoloop refinar "<nome>"`).

Em vez de reconstruir o todo, o comando `refinar` executa uma **cirurgia atômica**:

1. Ele isola o nó que possui uma recomendação de melhoria ou defasagem detectada pela Matriz.
2. Apresenta ao usuário a comparação precisa (*diff*) entre a peça antiga e a peça aprimorada.
3. Roda os testes de regressão do contrato de entrada e saída.
4. Se todos os testes passarem, substitui a peça antiga pela nova sem alterar um único milímetro das conexões e arestas do restante da esteira.

---

## 11. O Verbo Remover: Desvinculação Limpa e Rollback

Remover ferramentas e processos de sistemas de IA tradicionais é uma das tarefas mais frustrantes da rotina de tecnologia: arquivos de configuração permanecem esquecidos em pastas ocultas do usuário, atalhos quebrados continuam poluindo as mensagens do assistente e variáveis de ambiente desnecessárias permanecem ativas.

O comando `legoloop remover "<nome>"` opera sob a garantia de **desvinculação limpa**:

* Remove os vínculos simbólicos de todos os ambientes de desenvolvimento conectados (Claude, Cursor, Windsurf, OpenCode) simultaneamente.
* Arquiva o processo no cofre com um registro temporal imutável.
* Gera um ponto de restauração (*snapshot* de rollback). Caso o usuário descubra no dia seguinte que precisava daquele processo antigo, basta executar `legoloop remover --rollback "<nome>"` para que o estado do sistema seja restaurado em milissegundos.

---

## 12. O Operador Compor: A Arquitetura de Sub-Loops Hierárquicos

Uma das leis fundamentais da engenharia de sistemas complexos afirma que **qualquer sistema grande que funciona de forma confiável evoluiu invariavelmente a partir de um sistema pequeno que funcionava de forma impecável**.

Tentar construir uma esteira gigantesca de 40 passos de uma só vez é garantia de falha: a complexidade combinatória e a quantidade de variáveis tornam a depuração humanamente impossível.

O LegoLoop resolve o desafio da escala através do operador de composição hierárquica `compor` (`legoloop compor "<alvo>" "<sub-loop>"`).

Da mesma forma que na indústria mecânica um motor completo é montado a partir de centenas de componentes pré-testados e entregue para a linha de montagem final como um único bloco fechado, o operador `compor` permite transformar um processo inteiro em um único nó dentro de uma esteira superior:

* Uma esteira de **Auditoria de Código** (composta por 5 nós internos) pode ser encaixada como um único nó dentro de uma esteira de **Lançamento de Software**.
* A esteira principal enxerga o sub-loop apenas através do seu contrato de entrada e saída, sem precisar gerenciar seus nós internos.
* Isso viabiliza a criação de sistemas corporativos de complexidade infinita mantendo a simplicidade visual e a clareza de depuração em cada nível.

### Rastreabilidade Técnica da Parte II

* `internal/process/fabricar.go` — Motor de criação e entrevista socrática da Tríade Canônica.
* `internal/matrix/matrix.go` — Algoritmo de auditoria de conformidade e pontuação dos Tiers de qualidade.
* `internal/process/refinar.go` — Mecanismo de substituição cirúrgica com preservação de topologia.
* `internal/process/remover.go` — Desvinculação atômica multi-harness com cofre de rollback.
* `internal/process/compor.go` — Resolução matemática e aninhamento de sub-loops no grafo acíclico.