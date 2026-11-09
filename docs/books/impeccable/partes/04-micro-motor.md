# PARTE IV — Micro: Piso de Qualidade, Agentes e Integração

## 15. O piso de qualidade por dentro

O arquivo `craft-floor.md` fica em `upstream/reference/` dentro da skill `aidd-impeccable`. Ele é carregado depois que a direção está fixada e tem duas seções: o que verificar e o que recusar. O próprio arquivo avisa que um briefing fixado ou o mundo visual aprovado tem prioridade sobre qualquer regra dele; o hábito do agente, não.

### O que verificar

Cada item é uma checagem sobre o resultado renderizado, e não uma intenção. Todos são conferidos na mesma captura.

| Eixo                    | Critério verificável                                     |
| :---------------------- | :------------------------------------------------------- |
| Contraste               | Corpo e placeholder ≥ 4,5:1; texto grande ≥ 3:1          |
| Profundidade            | Sombras com deslocamento e desfoque suave                |
| Espaçamento             | Grupos apertados, separação generosa                     |
| Tipografia              | Corpo entre 65 e 75ch; display até 6rem                  |
| Tracking                | Piso de -0,04em; -0,02 a -0,03em lê melhor               |
| Movimento               | Um momento autoral, ease-out exponencial                 |
| Estados                 | Hover, disabled, loading, error e empty                  |
| Superfícies             | Seleção, cursor, scrollbar e foco tematizados            |
| Texto                   | Linguagem do produto; erros com saída                    |
| Cobertura               | Todo requisito do briefing achável em segundos           |

Sobre o contraste em superfícies coloridas, o arquivo acrescenta uma regra fina: o texto secundário deve ser tingido a partir do matiz da superfície ou do primeiro plano, nunca cinza.

### O que recusar

O arquivo trata a maior parte desta lista como padrões da categoria, e não como proibições: o briefing pode autorizar qualquer um deles. Recorrer a um deles quando o eixo está livre, diz o texto, significa que o agente não estava decidindo.

**Estruturas de página:**

* Cartões do mesmo tamanho com ícone, título e texto como estrutura da página. Cartões aninhados são sempre errados.
* O modelo de métrica de destaque: número grande, legenda pequena, estatísticas de apoio.
* Números de seção como 01, 02, 03, a menos que a sequência carregue informação.
* Um modal para uma tarefa que não precisa de interrupção nem de foco protegido.

**Hábitos de superfície:**

* Texto com gradiente. A ênfase vem de peso ou tamanho.
* Vidro e desfoque como decoração.
* Borda lateral colorida acima de 1px em cartões, itens de lista ou alertas.
* Sombras duras deslocadas fora de um mundo de fato neobrutalista.
* Fonte monoespaçada como fantasia de "técnico", e não para código ou dados.
* Glifos Unicode ou emoji no lugar de um sistema de ícones.
* Claro ou escuro escolhido pela categoria, e não pela cena de uso.

```{=typst}
#painel("A única proibição absoluta", cor: cor-erro)[
  Entre todos os itens, o arquivo marca um como banimento, e não como padrão: o rótulo pequeno (kicker ou eyebrow) acima de um título. Nenhum briefing o recupera. O título carrega o próprio peso; o rótulo deve ser apagado.
]
```

O arquivo também fixa limites numéricos: raio de cartões entre 12 e 16 pixels, com formato de pílula reservado a controles pequenos, e elevação declarada uma única vez, por borda ou por sombra. Uma borda de 1px sob uma sombra larga e suave é chamada de "cartão fantasma".

Por fim, uma regra de honestidade que vale para qualquer interface: afirmações e configurações vêm de dados fornecidos, e valores ilustrativos são rotulados como tal.

## 16. Os agentes especializados

Para evitar contaminação de contexto, o ecossistema do Impeccable delega quatro papéis a agentes autônomos. Cada um recebe apenas o que precisa para sua função.

| Agente                 | Papel                                     | Momento          |
| :--------------------- | :---------------------------------------- | :--------------- |
| Documenter             | Gera DESIGN.md a partir do código entregue | Estágio 5        |
| Finish reviewer        | Julga o acabamento contra a direção       | Gate P4          |
| Asset producer         | Produz ícones e vetores SVG inline        | Estágio 3        |
| Manual edit applier    | Aplica micro-ajustes da auditoria         | Estágios 3 e 4   |

O manual descreve o revisor de acabamento como um diretor de arte intransigente: ele decide se o resultado alcançou o nível esperado ou precisa ser refeito. Os SVG do produtor de ativos seguem um padrão semântico, sem emoji e com `stroke="currentColor"`, para herdarem a cor do texto ao redor.

A separação em agentes tem uma razão prática. O agente que implementou uma tela tende a defendê-la. Um revisor que só recebe a direção aprovada e o resultado final julga sem esse viés.

## 17. Integração no LegoLoop

No LegoLoop, o pipeline existe em duas formas canônicas, ambas verificáveis no repositório.

### A peça de loop no acervo

O arquivo `loops/impeccable-design.loop.json` segue o esquema `dag.v1` e se chama "Pipeline de Design Impeccable (6 Sub-Loops)", versão 1.0.0. Ele declara seis nós em cadeia linear, cada um dependente do anterior, com o artefato exigido e o modelo sugerido em `routing`:

| Nó                            | Artefato que o gate exige                  | Modelo sugerido  |
| :---------------------------- | :----------------------------------------- | :--------------- |
| `subloop-avaliar`               | `docs/reports/AUDITORIA_DESIGN.md`            | sonnet, ou haiku |
| `subloop-construir`             | `docs/specs/DIRECAO_VISUAL.md`                | sonnet, ou haiku |
| `subloop-melhorar`              | `src/styles/theme.css`                       | não declarado    |
| `subloop-refinar`               | `docs/reports/HARDENING_UI.md`                | não declarado    |
| `subloop-corrigir`              | `docs/reports/RESPONSIVE_VERIFICATION.md`     | não declarado    |
| `subloop-iterar-documentar`     | `DESIGN.md`                                  | não declarado    |

Cada nó tem um gate mecânico do tipo `file_exists(...)`: o nó só é considerado concluído se o artefato existir no disco. É um gate de existência, não de qualidade; quem julga o conteúdo continua sendo o protocolo humano P1 a P4 conduzido pela skill.

Todos os seis nós apontam para a mesma peça, `skills/impeccable`. O que muda entre eles é o sub-loop que a peça executa.

### O template do Studio Canvas

Em `internal/studio/templates.go`, o template "Pipeline de Design Impeccable (6 Sub-Loops)" aparece no menu de templates rápidos com o ícone `palette`. Ao instanciá-lo, o canvas recebe a esteira completa com o roteamento cognitivo pré-configurado.

```{=typst}
#painel("Do arquivo ao canvas")[
  O JSON do acervo e o template do Studio descrevem a mesma esteira por dois caminhos: o primeiro é executado pelo motor DAG via CLI; o segundo é montado visualmente para quem prefere arrastar nós. Alterar um não altera o outro automaticamente, e por isso vale conferir os dois quando o pipeline evoluir.
]
```

### A skill no harness

A skill `aidd-impeccable` funciona como maestro: recebe um pedido em linguagem natural, como "auditar a tela dashboard" ou "polir antes de entregar", e escolhe comando, modo e ordem. Ela carrega a metodologia de uma cópia local do upstream, sem depender de plugin, e por isso roda em qualquer harness que leia arquivos.

## 18. Como conferir a integração

Três comandos do CLI `legoloop` bastam para confirmar que a esteira está íntegra antes de usá-la.

```bash
# 1. Valida contrato, gates e ciclos do DAG
legoloop validate loops/impeccable-design.loop.json

# 2. Executa a esteira pelo motor DAG
legoloop run loops/impeccable-design.loop.json

# 3. Abre o Looping Studio (padrão: 127.0.0.1:7878)
legoloop studio
```

Em 07/11/2026, o primeiro comando respondeu `6 nós válidos, sem ciclos`, em seis camadas de um nó cada. Uma camada reúne nós que podem rodar ao mesmo tempo; como cada sub-loop depende do anterior, uma camada com dois nós aqui seria erro de declaração.
