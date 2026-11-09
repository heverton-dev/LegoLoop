// =============================================================================
// TEMPLATE TYPST CORPORATIVO — ESTILO EDITORIAL MINIMALISTA COM COLUNA DE ANOTAÇÕES
// Direção de arte: Capa branca padronizada, tipografia clássica,
// coluna lateral direita com linhas pautadas de caderno para anotações do leitor,
// cabeçalho discreto em linha única e paleta Google Suavizada.
// =============================================================================

#let cor-fundo    = rgb("#ffffff")                                    // Branco puro limpo
#let cor-miolo    = rgb("#fafafa")                                    // Off-white sutil para leitura confortável
#let cor-tinta    = rgb("#0f172a")                                    // Ardósia escuro nobre
#let cor-texto    = rgb("#292524")                                    // Grafite quente stone
#let cor-suave    = rgb("#64748b")                                    // Cinza pedra suave
#let cor-linha    = rgb("#e2e8f0")                                    // Linha divisória sutil
#let cor-accent   = rgb("#1d4ed8") // Azul Cobalto Google Suave
#let cor-accent2  = cor-accent.darken(25%)                            // Tom profundo da cor da obra
#let cor-alerta   = rgb("#d97706")                                    // Âmbar Google Dourado
#let cor-ok       = rgb("#16a34a")                                    // Verde Esmeralda Google
#let cor-erro     = rgb("#dc2626")                                    // Carmesim Google

#set document(
  title: [O Mega Pipeline do LegoLoop],
  author: "Heverton & Ecossistema LegoLoop",
)

#set page(
  width: 170mm,
  height: 240mm,
  fill: cor-fundo,
  margin: (top: 2.0cm, bottom: 2.0cm, left: 1.8cm, right: 1.8cm),
)

// Tipografia: Georgia clássica no miolo com entrelinha arejada 1.5
#set text(
  font: ("Georgia", "Segoe UI", "Calibri"),
  size: 10.3pt,
  lang: "pt",
  region: "BR",
  fill: cor-texto,
  hyphenate: true,
)

// Entrelinha arejada 1.45 a 1.5
#set par(justify: true, leading: 1.05em, spacing: 1.35em)
#set smartquote(enabled: false)

// ---------------------------------------------------------------------------
// Blocos Utilitários de Destaque (Design Clean e Discreto)
// ---------------------------------------------------------------------------

// Painel com cantos arredondados, fundo sutil e barra lateral limpa
#let painel(titulo, corpo, cor: cor-accent) = block(
  width: 100%, breakable: false, above: 16pt, below: 16pt,
  fill: rgb("#f8fafc"), radius: 6pt,
  stroke: (left: 3.5pt + cor, rest: 0.6pt + cor-linha),
  inset: (x: 14pt, y: 12pt),
)[
  #text(font: ("Segoe UI", "Arial"), size: 8.5pt, weight: "bold", fill: cor, tracking: 1.2pt)[#upper(titulo)]
  #v(4pt)
  #set text(font: ("Georgia", "Segoe UI"), size: 9.6pt, fill: cor-texto)
  #set par(leading: 0.9em, spacing: 1em)
  #corpo
]

// Frase de Impacto com barra lateral limpa
#let frase-destaque(corpo, cor: cor-accent) = block(
  width: 100%, above: 18pt, below: 18pt,
  stroke: (left: 4pt + cor), inset: (left: 14pt, y: 4pt),
)[
  #text(
    font: ("Segoe UI", "Arial"),
    size: 13pt, weight: "bold", fill: cor-tinta,
  )[#corpo]
]

#let chip(txt, cor: cor-accent) = box(
  fill: cor.lighten(90%), stroke: 0.5pt + cor.lighten(50%),
  inset: (x: 5pt, y: 2.5pt), radius: 3pt, outset: (y: 1.5pt),
)[#text(font: ("Segoe UI", "Arial"), size: 7.8pt, weight: "bold", fill: cor.darken(15%))[#txt]]

#let no(titulo, sub: none, cor: cor-tinta) = block(
  width: 100%, fill: cor, radius: 4pt, inset: (x: 4pt, y: 7pt),
)[
  #align(center)[
    #text(font: ("Segoe UI", "Arial"), size: 7.8pt, weight: "bold", fill: white, hyphenate: false)[#titulo]
    #if sub != none [ #linebreak() #text(font: ("Segoe UI", "Arial"), size: 6.6pt, fill: rgb("#cbd5e1"), hyphenate: false)[#sub] ]
  ]
]

#let no-claro(titulo, sub: none, cor: cor-accent) = block(
  width: 100%, fill: cor.lighten(92%), radius: 4pt,
  stroke: 0.6pt + cor.lighten(40%), inset: (x: 4pt, y: 7pt),
)[
  #align(center)[
    #text(font: ("Segoe UI", "Arial"), size: 7.8pt, weight: "bold", fill: cor.darken(25%), hyphenate: false)[#titulo]
    #if sub != none [ #linebreak() #text(font: ("Segoe UI", "Arial"), size: 6.6pt, fill: cor-suave, hyphenate: false)[#sub] ]
  ]
]

#let seta = align(horizon + center)[#text(size: 11pt, fill: cor-accent)[#sym.arrow.r]]

#let ficha(..pares) = block(width: 100%, above: 12pt, below: 14pt)[
  #table(
    columns: (auto, 1fr),
    stroke: 0.5pt + cor-linha,
    fill: (_, y) => if calc.odd(y) { rgb("#f8fafc") } else { white },
    inset: (x: 6pt, y: 4.5pt),
    ..pares.pos().map(p => (
      text(font: ("Segoe UI", "Arial"), size: 7.5pt, weight: "bold", fill: cor-suave)[#p.at(0)],
      text(font: ("Segoe UI", "Arial"), size: 7.8pt, fill: cor-tinta)[#p.at(1)],
    )).flatten()
  )
]

// ---------------------------------------------------------------------------
// Títulos de Seção — APENAS TEXTO LIMPO (Sem estilo botão, sem blocos escuros)
// ---------------------------------------------------------------------------

#set heading(numbering: none)
#show heading: set text(font: ("Segoe UI", "Arial"))

// Títulos de Seção — Quebra de página APENAS em Nível 1 (Partes)
#show heading.where(level: 1): it => {
  pagebreak(weak: true)
  block(width: 100%, above: 0pt, below: 22pt)[
    #text(size: 19pt, weight: "black", fill: cor-tinta)[#it.body]
    #v(6pt)
    #line(length: 100%, stroke: 1.5pt + cor-linha)
  ]
}

// Nível 2 (Capítulos e Receitas) é TUDO CORRIDO, sem quebra de página
#show heading.where(level: 2): it => block(above: 20pt, below: 8pt, width: 100%)[
  #text(size: 12.5pt, weight: "bold", fill: cor-tinta)[#it.body]
  #v(2.5pt)
  #line(length: 100%, stroke: 0.5pt + cor-linha)
]

#show heading.where(level: 3): it => block(above: 14pt, below: 5pt)[
  #text(size: 10.5pt, weight: "bold", fill: cor-accent)[#it.body]
]

#show heading.where(level: 4): it => block(above: 10pt, below: 3pt)[
  #text(size: 9.2pt, weight: "bold", fill: cor-suave)[#it.body]
]

// ---------------------------------------------------------------------------
// Código, Tabelas, Citações e Imagens
// ---------------------------------------------------------------------------

#let _quebravel(t) = {
  let saida = t
  for sep in ("/", "_", ".", "-") {
    saida = saida.replace(sep, sep + "\u{200B}")
  }
  saida
}

#show raw.where(block: false): it => highlight(
  fill: rgb("#f1f5f9"), extent: 1.5pt, radius: 2.5pt,
  top-edge: "ascender", bottom-edge: "descender",
)[#text(font: ("Consolas", "Courier New"), size: 8.0pt, fill: cor-accent2, hyphenate: false)[#_quebravel(it.text)]]

#show raw.where(block: true): it => block(
  width: 100%, fill: rgb("#f8fafc"),
  stroke: (left: 3pt + cor-accent, rest: 0.6pt + cor-linha),
  radius: 4pt, inset: (x: 11pt, y: 10pt), above: 13pt, below: 13pt, breakable: true,
)[#text(font: ("Consolas", "Courier New"), size: 8.0pt, fill: cor-tinta)[#it]]

#show quote.where(block: true): it => block(
  width: 100%, fill: rgb("#f8fafc"), stroke: (left: 3pt + cor-accent),
  inset: (x: 13pt, y: 10pt), radius: 4pt, above: 14pt, below: 14pt, breakable: true,
)[#set text(font: ("Georgia", "Segoe UI"), size: 10pt, fill: cor-suave, style: "italic"); #it.body]

#set list(indent: 8pt, body-indent: 7pt, spacing: 0.88em)
#set enum(indent: 8pt, body-indent: 7pt, spacing: 0.88em)

#set table(
  stroke: 0.35pt + cor-linha,
  inset: (x: 4pt, y: 3pt),
  fill: (_, y) => if y == 0 { rgb("#f1f5f9") } else if calc.odd(y) { rgb("#f8fafc") } else { white },
)
#show table.cell.where(y: 0): set text(font: ("Segoe UI", "Arial"), fill: cor-tinta, weight: "bold", size: 6.2pt, hyphenate: false)
#show table: set text(size: 5.9pt, hyphenate: false)
#show table: set par(justify: false, leading: 0.45em)
#show table: it => {
  show raw.where(block: false): r => text(font: ("Consolas", "Courier New"), size: 5.6pt, fill: cor-accent2, hyphenate: false)[#r.text]
  it
}
#show figure.where(kind: table): set block(breakable: true)
#show figure.where(kind: table): set align(left)
#show figure.where(kind: table): set figure(numbering: none)
#show figure.where(kind: table): set figure.caption(position: top)
#show figure.caption: it => text(font: ("Segoe UI", "Arial"), size: 8.4pt, style: "italic", fill: cor-suave)[#it.body]

#show link: it => text(fill: cor-accent, weight: "medium")[#it]
#let horizontalrule = align(center)[#block(above: 16pt, below: 16pt)[#line(length: 25%, stroke: 1pt + cor-linha)]]

#show strong: set text(hyphenate: false)

#show figure.where(kind: image): it => block(
  width: 100%, above: 16pt, below: 16pt, align(center)[
    #it.body
    #v(4pt)
    #text(font: ("Segoe UI", "Arial"), size: 8.4pt, fill: cor-suave)[#it.caption]
  ]
)

#let esteira(..itens, gates: ()) = {
  let pos = itens.pos()
  if pos.len() == 1 and type(pos.first()) == array {
    let etapas = pos.first()
    block(width: 100%, above: 14pt, below: 14pt, breakable: false)[
      #for (i, e) in etapas.enumerate() {
        block(width: 100%, radius: 4pt, fill: rgb("#f8fafc"), stroke: (left: 3pt + cor-accent, rest: 0.5pt + cor-linha), inset: (x: 10pt, y: 6pt), below: 0pt)[
          #text(font: ("Segoe UI", "Arial"), size: 8.6pt, weight: "bold", fill: cor-tinta)[#e.at(0)]
          #h(6pt)
          #text(font: ("Segoe UI", "Arial"), size: 7.6pt, fill: cor-suave)[#e.at(1)]
        ]
        if i < etapas.len() - 1 {
          align(center)[#block(above: 3pt, below: 3pt)[#text(font: ("Segoe UI", "Arial"), size: 7.2pt, fill: cor-accent2, weight: "bold")[↓ #if i < gates.len() { gates.at(i) }]]]
        }
      }
    ]
  } else {
    let n = pos.len()
    block(width: 100%, above: 14pt, below: 14pt)[
      #grid(
        columns: range(n * 2 - 1).map(i => if calc.even(i) { 1fr } else { auto }),
        align: horizon, column-gutter: 4pt,
        ..pos.enumerate().map(((i, it)) => if i == n - 1 { (it,) } else { (it, seta) }).flatten()
      )
    ]
  }
}


// ---------------------------------------------------------------------------
// CAPA — DESIGN STRIPE PRESS / EDITORIAL MINIMALISTA CONTEMPORÂNEO (17x24cm)
// ---------------------------------------------------------------------------
#page(margin: (x: 1.8cm, top: 2.0cm, bottom: 1.8cm), numbering: none, header: none, footer: none, fill: white)[

  // Detalhe arquitetural no topo: 4 traços finos elegantes na paleta Google Suavizada
  #grid(
    columns: (1fr, 1fr, 1fr, 1fr),
    gutter: 5pt,
    line(length: 100%, stroke: 2pt + cor-accent),
    line(length: 100%, stroke: 2pt + cor-erro),
    line(length: 100%, stroke: 2pt + cor-alerta),
    line(length: 100%, stroke: 2pt + cor-ok),
  )

  #v(16pt)

  // Categoria editorial
  #text(
    font: ("Segoe UI", "Arial"),
    size: 8pt,
    fill: cor-suave,
    weight: "bold",
    tracking: 2.2pt,
  )[ARQUITETURA DE SISTEMAS AGÊNTICOS]

  #v(12pt)

  // Título Principal e Destaque
      #text(
      font: ("Segoe UI", "Arial"),
      size: 28pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.5pt,
    )[Mega Pipeline]
    #v(4pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 15.5pt,
      weight: "bold",
      fill: cor-accent,
    )[A Engenharia Determinística dos 5 Sub-Loops]
  
  #v(8pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 9.6pt,
      fill: cor-suave,
      weight: "medium",
    )[A Engenharia Determinística dos 5 Sub-Loops, Contexto Estéril e Governança de Processos com IA]
  ]

  #v(16pt)

  // Arte: esteira de sub-loops desenhada em Typst
    #block(width: 100%, radius: 6pt, clip: true, fill: rgb("#f8fafc"), stroke: 0.5pt + rgb("#cbd5e1"), inset: 16pt)[
    #let etapas = ("Sub-Loop 1: Planejar", "Sub-Loop 2: Roteamento", "Sub-Loop 3: Executar", "Sub-Loop 4: Auditar", "Sub-Loop 5: Concluir", "Handoff Imutável", )
    #text(font: ("Segoe UI", "Arial"), size: 7pt, weight: "bold", fill: cor-suave, tracking: 1.6pt)[OS CINCO SUB-LOOPS DO MEGA PIPELINE]
    #v(8pt)
    #grid(
      columns: (1fr,) * 2, gutter: 8pt,
      ..etapas.enumerate().map(((i, e)) => block(
        width: 100%, height: 1.7cm, radius: 4pt,
        fill: if calc.even(i) { cor-accent.lighten(88%) } else { white },
        stroke: (left: 3pt + cor-accent, rest: 0.5pt + cor-linha),
        inset: (x: 8pt, y: 6pt),
      )[
        #text(font: ("Consolas", "Courier New"), size: 6.5pt, fill: cor-accent2)[#str(i + 1).]
        #h(3pt)
        #text(font: ("Segoe UI", "Arial"), size: 8pt, weight: "bold", fill: cor-tinta)[#e]
      ])
    )
  ]
  
  #v(1fr)

  // Rodapé da Capa
  #line(length: 100%, stroke: 0.6pt + cor-linha)
  #v(6pt)
  #grid(
    columns: (1fr, auto),
    align: (left + horizon, right + horizon),
    text(font: ("Segoe UI", "Arial"), size: 8.8pt, weight: "bold", fill: cor-tinta)[Heverton & Ecossistema LegoLoop],
    text(font: ("Segoe UI", "Arial"), size: 8.2pt, fill: cor-suave)[06/11/2026 · LegoLoop Foundation]
  )
]

// ---------------------------------------------------------------------------
// FOLHA DE ROSTO / COLOFÃO
// ---------------------------------------------------------------------------
#page(numbering: none, header: none, footer: none, fill: cor-miolo)[
  #v(2cm)
  #text(font: ("Segoe UI", "Arial"), size: 8.8pt, weight: "bold", fill: cor-accent, tracking: 1.8pt)[SOBRE ESTA OBRA]
  #v(6pt)
  #line(length: 100%, stroke: 0.8pt + cor-linha)
  #v(16pt)
    #set text(size: 9.8pt, fill: cor-texto)
  #set par(leading: 0.95em, spacing: 1.2em)
  Este livro documenta a arquitetura, ontologia e implementação do Mega Pipeline do LegoLoop --- o orquestrador mestre determinístico para desenvolvimento de software auxiliado por inteligência artificial. A obra formaliza o conceito de processo na engenharia agêntica como uma máquina de estados finita, governada e orientada a contratos, eliminando o caos probabilístico de execuções soltas.

  Percorremos a esteira completa dividida em três níveis canônicos: o nível macro (a ontologia de processos e o paradigma fractal), o nível meso (o fluxo dos cinco sub-loops com isolamento hermético via Orchestrix) e o nível micro (o detalhamento das etapas internas, seus contratos de entrada, processamento, peças acopladas, gates mecânicos e handoffs imutáveis).

  A obra aborda também a sextilha canônica de peças e a cascata de governança em quatro estágios, oferecendo aos arquitetos, engenheiros e desenvolvedores o guia definitivo para construir esteiras industriais previsíveis, de altíssimo rigor e livre de alucinações.

  
  #v(1fr)
  #line(length: 100%, stroke: 0.8pt + cor-linha)
  #v(8pt)
  #text(font: ("Segoe UI", "Arial"), size: 8.2pt, fill: cor-suave)[
    Documento gerado deterministicamente a partir das fontes
    citadas no próprio texto, via pandoc + typst. Toda afirmação técnica desta obra é
    rastreável a um arquivo nomeado no capítulo correspondente.
  ]
  #v(2pt)
  #text(font: ("Segoe UI", "Arial"), size: 7.8pt, fill: cor-suave)[
    Data da compilação: 06/11/2026 · LegoLoop Foundation
  ]
]

// ---------------------------------------------------------------------------
// CORPO: CABEÇALHO, COLUNA DE ANOTAÇÕES (MARGINALIA PONTILHADA) E RODAPÉ
// ---------------------------------------------------------------------------
#set page(
  fill: cor-miolo,
  numbering: "1",
  margin: (top: 2.0cm, bottom: 2.0cm, left: 1.6cm, right: 3.8cm),
  background: context {
    if here().page() > 2 {
      // Linha vertical divisória da coluna de anotações
      place(top + left, dx: 13.2cm, dy: 2.0cm, line(angle: 90deg, length: 20cm, stroke: 0.4pt + rgb("#e2e8f0")))
      // Rótulo discreto de anotações
      place(top + left, dx: 13.5cm, dy: 2.1cm, text(font: ("Segoe UI", "Arial"), size: 6.2pt, fill: rgb("#94a3b8"), weight: "bold", tracking: 1.6pt)[ANOTAÇÕES])
      // Pauta de caderno pontilhada sutil para anotações do leitor
      for i in range(21) {
        place(top + left, dx: 13.5cm, dy: 2.6cm + i * 8.5mm, line(length: 2.7cm, stroke: (paint: rgb("#e2e8f0"), dash: "dotted", thickness: 0.45pt)))
      }
    }
  },
  header: context {
    if here().page() > 2 {
      let pagina = here().page()
      let nesta = query(heading.where(level: 1)).filter(h => h.location().page() == pagina)
      let anteriores = query(selector(heading.where(level: 1)).before(here()))
      let titulo-cap = if nesta.len() > 0 { nesta.first().body }
        else if anteriores.len() > 0 { anteriores.last().body }
        else { [Sumário] }
      stack(
        spacing: 2.5pt,
        grid(
          columns: (auto, 1fr),
          gutter: 10pt,
          align: (left + horizon, right + horizon),
          text(font: ("Segoe UI", "Arial"), size: 5.8pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[MEGA PIPELINE],
          text(font: ("Segoe UI", "Arial"), size: 5.8pt, fill: cor-suave, hyphenate: false)[#titulo-cap],
        ),
        line(length: 100%, stroke: 0.5pt + cor-linha),
      )
    }
  },
  footer: context {
    if here().page() > 2 {
      align(center)[
        #text(font: ("Segoe UI", "Arial"), size: 8.2pt, fill: cor-suave)[#counter(page).display("1")]
      ]
    }
  },
)
#counter(page).update(1)

#block(above: 0pt, below: 18pt)[
  #text(font: ("Segoe UI", "Arial"), size: 18pt, weight: "bold", fill: cor-tinta)[Sumário]
  #v(5pt)
  #line(length: 100%, stroke: 1.5pt + cor-linha)
]
#show outline.entry.where(level: 1): it => {
  v(6pt, weak: true)
  text(font: ("Segoe UI", "Arial"), size: 9.6pt, weight: "bold", fill: cor-tinta)[#it]
}
#show outline.entry.where(level: 2): it => text(font: ("Segoe UI", "Arial"), size: 8.8pt, fill: cor-suave)[#it]
#show outline: set par(leading: 0.45em, spacing: 0.5em)
#show outline.entry: set block(spacing: 0.42em)
#outline(title: none, depth: 2, indent: 1.2em)


= Como ler este livro
<como-ler-este-livro>
== A quem este livro se destina
<a-quem-este-livro-se-destina>
Esta obra foi projetada para três perfis centrais da engenharia contemporânea:

+ #strong[Arquitetos de Sistemas e Agentes:] que necessitam desenhar fluxos de trabalho autônomos com garantias matemáticas de não colisão, orçamentação estrita de tokens e eliminação de vazamento de contexto (#emph[Zero Context Bloat]). Comece pelo Capítulo 1 e siga para o Capítulo 2.
+ #strong[Engenheiros de Plataforma e DevOps:] responsáveis por integrar ferramentas, MCPs, scripts determinísticos e esteiras de CI/CD à governança agêntica. Dê foco especial ao Capítulo 3 e às convenções de contratos no Capítulo 4.
+ #strong[Desenvolvedores Full-Stack e Especialistas de TDD:] que operam o LegoLoop Studio e a interface de linha de comando no dia a dia. Explore detalhadamente as fichas técnicas do Capítulo 3 e o caderno prático no Capítulo 4.

== A estrutura da obra
<a-estrutura-da-obra>
O livro organiza o conhecimento técnico em três camadas de profundidade analítica:

#esteira(
  no("NÍVEL MACRO", sub: "Ontologia & Processo"),
  no("NÍVEL MESO", sub: "Mega Pipeline & 5 Loops"),
  no("NÍVEL MICRO", sub: "Etapas, I/O & Gates"),
  no("TRANSVERSAL", sub: "Templates & Governança"),
)
Cada nível responde de forma consistente a três perguntas fundamentais: \* #strong[Qual é o contrato?] Quais arquivos entram, quais saem e qual regra aprova a transição. \* #strong[Qual é o motor?] Quais cérebros cognitivos (LLMs) ou scripts locais realizam o trabalho. \* #strong[Qual é a segurança?] Como o isolamento estéril e os Stop-Gates humanos protegem a integridade do código.

== As convenções visuais
<as-convenções-visuais>
O livro adota elementos gráficos padronizados para acelerar a leitura e evidenciar decisões:

#painel("Painel de Governança & Regras Críticas")[
  Caixas com esta barra lateral azul representam regras inquebráveis, invariantes de arquitetura ou alertas de segurança. Quando um ponto do sistema exigir parada para inspeção humana, ele será destacado dessa forma.
]
- #strong[Caixas Escuras (`#no`):] Representam etapas puramente mecânicas e determinísticas (scripts locais, verificações de arquivos, linters e compilação de binários) que rodam sem qualquer uso de inteligência artificial.
- #strong[Caixas Claras (`#no-claro`):] Representam etapas cognitivas onde modelos de linguagem são acionados sob limites estritos de raciocínio (#emph[effort clamping]).
- #strong[Fichas Técnicas (`#ficha`):] Tabelas sintéticas que abrem cada seção apresentando a carteira de identidade daquele componente.

= PARTE I --- NÍVEL MACRO: ONTOLOGIA E PROCESSO
<parte-i-nível-macro-ontologia-e-processo>
== 1. O Conceito Formal de Processo na Engenharia de IA
<o-conceito-formal-de-processo-na-engenharia-de-ia>
Na ciência da computação clássica e na engenharia de sistemas de missão crítica, a definição de um #strong[Processo] é rigorosa:

#painel("Definição Axiomática de Processo")[
  Um processo é uma sequência finita, ordenada e governada de transições de estado, que transforma entradas definidas em saídas verificáveis, operando sob regras explícitas, com recursos previamente alocados e sob critérios determinísticos de validação e parada.
]
Quando a inteligência artificial generativa foi introduzida no desenvolvimento de software, a indústria cometeu o erro de tratá-la como um chat interativo contínuo. Em conversas longas e desestruturadas, os modelos sofrem de alucinação cumulativa, esquecimento de diretrizes anteriores e poluição descontrolada de contexto (#emph[context bloat]).

O LegoLoop nasce para refundar a engenharia de software com IA sobre bases fabris determinísticas. Em vez de prompts soltos, todo trabalho é encapsulado em uma #strong[máquina de estados finita] representada por um Grafo Acíclico Dirigido (DAG).

#frase-destaque[Um processo é uma máquina de estados finita: cada nó transforma entradas declaradas em saídas verificáveis, sem depender de memória mágica.]
#figure(image("/img/diagrama_monolito_vs_legoloop.svg", alt: "O Monólito Tradicional versus a Esteira Modular LegoLoop"),
  caption: [
    O Monólito Tradicional versus a Esteira Modular LegoLoop
  ]
)

#ficha(
  ("Natureza", "Máquina de Estados Finita e Auditável"),
  ("Grafo Base", "Directed Acyclic Graph (DAG v2) em JSON Schema"),
  ("Garantia de Parada", "Proibição estrita de ciclos infinitos (ErrCiclo)"),
  ("Handoff", "Exclusivamente mediado por sistema de arquivos (SSOT)"),
  ("Ambiente", "Workspaces temporários herméticos criados pelo Orchestrix"),
)
=== Os Cinco Pilares de Governança do Processo
<os-cinco-pilares-de-governança-do-processo>
+ #strong[Sequência Finita e Ordenada:] Cada processo possui um início formal demarcado por um gatilho (`trigger`) e um término inequívoco. A ordem de execução não é estocástica: o resolvedor topológico do motor garante que uma etapa só execute após todos os seus nós predecessores declarados em `depends_on` terem sido concluídos com sucesso.
+ #strong[Transformação de Estados Auditável:] Um nó nunca "acha" o que deve fazer. Ele recebe uma lista fechada de arquivos de entrada (`inputs`), aplica uma transformação controlada e materializa um artefato no disco (`output.artifact`).
+ #strong[Regras Explícitas e Ausência de Premissas Ocultas:] A IA é contratualmente proibida de inferir comportamento arquitetural em silêncio. Toda decisão deve estar ancorada nos documentos da tríade de planejamento.
+ #strong[Recursos Rigidamente Alocados:] Cada bloco do processo declara com antecedência sua cota de recursos: qual modelo cognitivo será acionado, qual o teto de tokens de saída (`max_output_tokens`), qual o limite de tempo (`timeout_sec`) e quais ferramentas locais estão liberadas.
+ #strong[Critérios Determinísticos de Validação:] A aprovação de uma etapa não depende da avaliação subjetiva do próprio modelo. Ela é julgada por expressões booleanas matemáticas puras compiladas em DSL (Domain Specific Language), tais como conferência de código de retorno (`exit_code == 0`) e existência física do arquivo (`file_exists`).

#horizontalrule

== 2. A Arquitetura Fractal do LegoLoop
<a-arquitetura-fractal-do-legoloop>
Uma das propriedades mais poderosas da engenharia do LegoLoop é a sua #strong[natureza fractal e homotética]. A mesma estrutura conceitual e contratual que rege a esteira inteira se repete em cada uma de suas subdivisões:

#esteira(
  no("MEGA PIPELINE", sub: "Macro-Orquestrador"),
  no("SUB-LOOPS", sub: "Unidades Modulares"),
  no("LOOPS", sub: "DAGs Acíclicos"),
  no("NÓS / ETAPAS", sub: "Átomos Executáveis"),
)
- #strong[Mega Pipeline (A Fábrica Completa):] Encadeia a visão global de ponta a ponta. Ele transforma uma ideia bruta ou solicitação humana em um produto de software compilado, testado em múltiplos dispositivos e documentado com rastreabilidade de sessão.
- #strong[Sub-Loops (As Linhas de Montagem):] São especializações funcionais do processo (ex.: Concepção, Engenharia Reversa, Planejamento BDD, Construção TDD e Homologação). Cada sub-loop tem vida própria: pode ser invocado de forma autônoma ou encadeado como uma macro-peça.
- #strong[Loops e DAGs (O Contrato Isomórfico):] No nível de dados e código, #strong[não há diferença estrutural entre um loop simples e um mega pipeline]. Ambos são manifestos JSON aderentes ao mesmo esquema formal (`https://legoloop.dev/schemas/dag.v1.json`). O Mega Pipeline é apenas um loop cujos nós referenciam outros loops.
- #strong[Nós e Etapas (`dag.Node`):] É a menor partícula operacional. Um nó nasce em um diretório temporário, recebe os artefatos de entrada, executa as ferramentas atribuídas, submete seu resultado ao crivo mecânico de um gate e encerra sua existência sem deixar resíduos de memória.

#horizontalrule

== 3. A Sextilha Canônica de Peças
<a-sextilha-canônica-de-peças>
Para que uma esteira seja flexível e componível, as etapas do processo utilizam um catálogo ontológico padronizado em seis naturezas de peças:

#figure(
  align(center)[#table(
    columns: (15.69%, 28.43%, 39.22%, 16.67%),
    align: (left,left,left,left,),
    table.header([Natureza da Peça], [Formato Nativo], [Função Primária no Fluxo], [Motor de Execução],),
    table.hline(),
    [#strong[\1. Skill]], [Markdown (`SKILL.md`)], [Prompt especialista com heurísticas e raciocínio de IA], [`Orchestrix`],
    [#strong[\2. Script]], [PowerShell / Shell / Go], [Transformação determinística de dados e arquivos locais], [`Runtime Local`],
    [#strong[\3. Hook]], [Configuração JSON], [Interceptação de ciclo de vida pré e pós-execução], [`Harness Hooks`],
    [#strong[\4. MCP]], [Configuração JSON (`.mcp.json`)], [Ferramental vivo externo (navegador, banco de dados)], [`Protocolo MCP`],
    [#strong[\5. Gate]], [Expressão booleana em DSL], [Validação mecânica binária e aprovação de transição], [`internal/gate`],
    [#strong[\6. Loop]], [Manifesto JSON (`.loop.json`)], [Encadeamento de nós em DAG para reúso modular], [`dag.Parse`],
  )]
  , kind: table
  )

#painel("Equilíbrio Cognitivo vs. Determinístico")[
  O segredo de uma esteira robusta está no casamento entre peças cognitivas (Skills) e peças determinísticas (Scripts e Gates). A inteligência artificial só deve ser chamada para raciocínio, escrita e síntese; tudo o que puder ser medido por um compilador, formatador ou script local deve ser executado fora da LLM para economizar tokens e garantir tolerância zero a falhas.
]
== Rastreabilidade
<rastreabilidade>
- Código do grafo e contrato de nós: `internal/dag/types.go` e `internal/dag/parse.go`.
- Validação de ciclos infinitos: `internal/dag/cycle.go`.
- Definições de templates canônicos: `docs/templates/README.md` e `docs/templates/TEMPLATE_PIECE_LOOP.json`.

= PARTE II --- NÍVEL MESO: O MEGA PIPELINE E A ORCHESTRIX
<parte-ii-nível-meso-o-mega-pipeline-e-a-orchestrix>
== 1. Visão Panorâmica dos Cinco Sub-Loops
<visão-panorâmica-dos-cinco-sub-loops>
O #strong[Mega Pipeline] organiza o ciclo de vida completo de engenharia de software em cinco sub-loops lineares e complementares. Nenhum sub-loop avança para o próximo sem que seu artefato terminal tenha sido aprovado por um portão de qualidade:

#esteira(
  no-claro("LOOP 1", sub: "Concepção & Intenção"),
  no-claro("LOOP 2", sub: "Design & Engenharia Reversa"),
  no-claro("LOOP 3", sub: "Tríade de Planejamento"),
  no("LOOP 4", sub: "Construção TDD por Épico"),
)
#v(-4pt)
#esteira(
  no("LOOP 5", sub: "Homologação E2E & Release"),
)
+ #strong[Loop 1 --- Concepção & Alinhamento:] Transforma o pedido do usuário em uma declaração clara de intenção (`INTENT.md`) através de uma entrevista socrática.
+ #strong[Loop 2 --- Engenharia Reversa & Design Pass:] Inspeciona referências visuais ou aplicações legadas, extraindo tokens limpos e wireframes técnicos (`DESIGN.md`).
+ #strong[Loop 3 --- Especificação & Tríade de Planejamento:] Formaliza a engenharia em critérios BDD executáveis (`SPEC.md`), quebra o trabalho em tickets (`PLAN.md`) e compila a rota acíclica (`ROUTE.json`).
+ #strong[Loop 4 --- Construção TDD por Épico:] Implementa o código em ciclo estrito de testes unitários (Red ➔ Green ➔ Refactor), módulo a módulo, sem tolerância a stubs.
+ #strong[Loop 5 --- Homologação E2E & Fechamento:] Audita o sistema com navegadores reais headless, recompila os binários de produção, consolida métricas de tokens e gera o relatório final.

#frase-destaque[Nenhum bloco herda o contexto do anterior. A Orchestrix cria e destrói diretórios estéreis em milissegundos, garantindo Zero Context Bloat.]
#figure(image("/img/diagrama_mega_pipeline_cinco_subloops.svg", alt: "A Esteira Mestre dos Cinco Sub-Loops do Mega Pipeline"),
  caption: [
    A Esteira Mestre dos Cinco Sub-Loops do Mega Pipeline
  ]
)

#horizontalrule

== 2. A Orchestrix Cirúrgica e o Handoff File-Based
<a-orchestrix-cirúrgica-e-o-handoff-file-based>
O maior ofensor de custos e alucinações em sistemas agênticos convencionais é o compartilhamento indiscriminado de histórico de conversação entre tarefas. Quando um agente recebe centenas de mensagens acumuladas de etapas passadas, o ruído semântico destrói a atenção do modelo.

O motor #strong[Orchestrix] do LegoLoop resolve esse problema impondo o isolamento cirúrgico por bloco:

#ficha(
  ("Mecanismo", "Diretório temporário e isolamento hermético"),
  ("Contrato de Entrada", "inputs: declaração estrita de arquivos consumidos"),
  ("Função de Transferência", "transferInputs(): copia apenas o necessário"),
  ("Contrato de Saída", "output.artifact: gravação no diretório pai"),
  ("Limpeza de Memória", "defer os.RemoveAll(): zero lixo residual no SO"),
)
=== O Ciclo de Execução Estéril de um Sub-Loop
<o-ciclo-de-execução-estéril-de-um-sub-loop>
Quando a Orchestrix é instruída a executar um nó do tipo `"type": "subloop"`, o código em `internal/orchestrix/subloop.go` executa o seguinte protocolo de seis passos:

+ #strong[Geração do Workspace Estéril:] O runtime chama `os.MkdirTemp("", "legoloop-subloop-*")`, instanciando uma pasta limpa no disco.
+ #strong[Transferência Cirúrgica de Entradas:] A função `transferInputs` inspeciona a lista de arquivos declarados em `inputs` e copia apenas esses arquivos para dentro da pasta temporária. Nenhum arquivo estranho é incluído.
+ #strong[Execução Isolada do Sub-Grafo:] O sub-grafo interno é parsed e executado pelo `runner.Run`. As etapas internas consomem e produzem dados exclusivamente dentro daquela bolha.
+ #strong[Avaliação dos Gates Mecânicos:] O gate de saída de cada etapa é testado localmente. Se houver falha, o processo é imediatamente interrompido.
+ #strong[Extração do Artefato Final:] A função `finalizeSubloop` copia o arquivo terminal contratado (`sub_artifact`) de volta para o diretório do projeto principal, preservando seu hash e integridade.
+ #strong[Destruição do Ambiente:] O diretório temporário é totalmente apagado via `os.RemoveAll()`, garantindo que resíduos de compilação ou arquivos intermediários não contaminem as etapas futuras.

#horizontalrule

== 3. Roteamento Cognitivo e Effort Clamping
<roteamento-cognitivo-e-effort-clamping>
Nem toda tarefa exige o mesmo poder de fogo computacional. Enviar a correção de um erro de lint para um modelo de raciocínio profundo de última geração é um desperdício financeiro inaceitável; da mesma forma, confiar a concepção da arquitetura de um sistema complexo a um modelo leve leva a colapsos de lógica.

O Mega Pipeline introduz a matriz de #strong[Roteamento Cognitivo com Clamping de Esforço]:

#figure(
  align(center)[#table(
    columns: (15.53%, 35.92%, 33.98%, 14.56%),
    align: (left,left,left,left,),
    table.header([Nível de Tarefa], [Escopo e Complexidade], [Modelo Cognitivo Sugerido], [Flag de Esforço],),
    table.hline(),
    [#strong[Épico (Macro)]], [Arquitetura, contratos BDD, DAG e modelagem], [`Claude 3.7 Sonnet`, `Gemini 2.0 Pro`], [`--effort high`],
    [#strong[Micro-ticket]], [Módulo funcional, suíte de testes unitários], [`Claude 3.5 Sonnet`, `Gemini 2.0 Flash`], [`--effort med`],
    [#strong[Nano-ticket]], [Correção de lint, tipo TypeScript, assert], [`Gemini Flash-Lite`, `Ollama Local 8B`], [`--effort low`],
  )]
  , kind: table
  )

#painel("Regra de Clamping de Tokens")[
  A flag `--effort` define o teto de tokens de raciocínio interno (thinking budget) alocados pela API. Em nano-tickets, o modelo é instruído a responder instantaneamente sem divagações, preservando a cota semanal para o que realmente importa: o design arquitetural e os contratos de teste.
]

#horizontalrule

== 4. As Sete Leis Invioláveis da Esteira
<as-sete-leis-invioláveis-da-esteira>
O Mega Pipeline opera debaixo de sete mandamentos mandatórios que nenhum agente ou executor tem autorização para violar:

+ #strong[Zero Context Bloat:] Nenhum bloco herda o histórico conversacional de blocos anteriores.
+ #strong[Stop-Gates de Auditoria Humana:] Paradas obrigatórias entre concepção, planejamento e ao término de cada Épico de código.
+ #strong[A Tríade Inegociável:] O código só pode começar se existirem no disco a `SPEC.md`, o `PLAN.md` e a rota executável `ROUTE.json`.
+ #strong[Provider-Awareness:] Roteamento de modelos adaptado ao ambiente local do desenvolvedor (`9Router`, `FreeLLMAPI` ou chave direta).
+ #strong[TDD Estrito & Zero Stubs:] O teste unitário deve falhar comprovadamente antes do código de produção. Proibição absoluta de funções vazias com `TODO` ou `pass`.
+ #strong[Padronização ISO de Nomenclatura:] Todos os artefatos de documentação seguem o prefixo `YYYY-MM-DD_<NOME>.<ext>`.
+ #strong[Executor Silencioso:] Respostas concisas no chat e pensamento interno telegráfico no padrão Caveman.

== Rastreabilidade
<rastreabilidade-1>
- Implementação do isolamento de sub-loops: `internal/orchestrix/subloop.go`.
- Runner determinístico do DAG: `internal/runner/runner.go`.
- Regras de conduta de agentes: `AGENTS.md` e `docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`.

= PARTE III --- NÍVEL MICRO: ANATOMIA DOS CINCO SUB-LOOPS
<parte-iii-nível-micro-anatomia-dos-cinco-sub-loops>
Esta seção detalha o nível atômico do Mega Pipeline. Cada um dos cinco sub-loops é dissecado em suas etapas operacionais, apresentando os contratos rigorosos de entrada, processamento, ferramental acoplado, validação mecânica e handoff imutável.

#horizontalrule

== 1. Sub-Loop 1: Concepção & Alinhamento de Intenção
<sub-loop-1-concepção-alinhamento-de-intenção>
O primeiro sub-loop tem por objetivo extrair a verdadeira necessidade do usuário, eliminando a ambiguidade natural da linguagem humana antes de qualquer especificação técnica.

#ficha(
  ("Papel", "Alinhamento socrático de escopo e premissas"),
  ("Disparo", "Linguagem natural do usuário via chat ou trigger manual"),
  ("Entrada", "Briefing inicial ou descrição da demanda"),
  ("Saída Principal", "docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md"),
  ("Portão de Saída", "STOP-GATE HUMANO 01 (Aprovação da Intenção)"),
)
#esteira(
  no-claro("1.1 ENTREVISTA", sub: "aidd-grill socrático"),
  no("1.2 ESPELHO", sub: "template ISO"),
  no("1.3 HUMAN-GATE", sub: "parada obrigatória"),
)
=== Etapa 1.1 --- Entrevista Socrática (Briefing Ativo)
<etapa-1.1-entrevista-socrática-briefing-ativo>
- #strong[Consome:] Linguagem natural do usuário via chat ou comando CLI.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-grill/SKILL.md")[`aidd-grill`] + Hook de telemetria de tokens.
- #strong[Processamento:] O agente assume uma postura socrática empática e incisiva, formulando uma pergunta por vez. Cada intervenção traz quatro marcadores formais: `HIPÓTESE`, `CONFIANÇA [X%]`, `PERGUNTA` afuniladora e `PALPITE` sugerido. O interrogatório se encerra quando a confiança atinge $gt.eq 95 %$ ou o modelo prevê com certeza as próximas respostas.
- #strong[Gate:] Confiança calculada $gt.eq 95 %$ ou encerramento socrático acordado.
- #strong[Handoff:] Transfere a estrutura de hipóteses validadas em memória para a Etapa 1.2.

=== Etapa 1.2 --- Geração do Espelho Canônico e Persistência
<etapa-1.2-geração-do-espelho-canônico-e-persistência>
- #strong[Consome:] Dados consolidados da Etapa 1.1.
- #strong[Peças:] Script gerador de template ISO + Gate mecânico.
- #strong[Processamento:] Materializa o documento formal estruturado em seis seções fixas: Resultado Esperado, Quem se Beneficia, Por que Agora, Critérios de Sucesso, Restrições Inegociáveis e Fora de Escopo.
- #strong[Gate:] Expressão DSL: `file_exists('docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md')`.
- #strong[Handoff:] Grava o arquivo físico persistido no disco.

=== Etapa 1.3 --- Human-Gate 01 (Validação da Intenção)
<etapa-1.3-human-gate-01-validação-da-intenção>
- #strong[Consome:] Documento de intenção gerado na Etapa 1.2.
- #strong[Peças:] Controle interativo `human-gate`.
- #strong[Processamento:] O runtime entra em pausa física obrigatória. Exibe o espelho canônico e aguarda o comando humano (`[1] seguir`, `[2] personalizar`, `[0] cancelar`).
- #strong[Gate:] Confirmação explícita do usuário.
- #strong[Handoff:] Destrava a execução para o Sub-Loop 2 ou Sub-Loop 3.

#horizontalrule

== 2. Sub-Loop 2: Engenharia Reversa & Design Pass #emph[\(Condicional)]
<sub-loop-2-engenharia-reversa-design-pass-condicional>
Acionado quando a demanda envolve recriar telas existentes, reconstruir interfaces legadas ou extrair design systems de produtos de terceiros.

#ficha(
  ("Papel", "Captura de runtime, dissecção DOM e tokens de design"),
  ("Disparo", "Condicional: presença de URL, mock, binário ou tela legada"),
  ("Entrada", "URL ou arquivos de interface existentes"),
  ("Saída Principal", "DESIGN.md e TOKENS.css"),
  ("Portão de Saída", "STOP-GATE HUMANO 02 (Validação Visual)"),
)
#esteira(
  no("2.1 INSPEÇÃO", sub: "Playwright / Obscura"),
  no-claro("2.2 NORMALIZAÇÃO", sub: "impeccable + OKLCH"),
  no("2.3 GATE DESIGN", sub: "validação de contraste"),
)
=== Etapa 2.1 --- Inspeção de Runtime & Captura DOM
<etapa-2.1-inspeção-de-runtime-captura-dom>
- #strong[Consome:] URL de referência externa, binário compilado ou mock de tela.
- #strong[Peças:] MCP `playwright`, MCP `chrome-devtools` ou motor #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/obscura/SKILL.md")[`obscura`] em Rust.
- #strong[Processamento:] Inicia o navegador em ambiente isolado, captura a árvore DOM semântica, calcula os estilos reais de CSS e realiza capturas de tela multi-viewport.
- #strong[Gate:] `exit_code == 0 && file_exists('scratch/raw_dom.json')`.
- #strong[Handoff:] Grava o payload JSON da estrutura da tela em `scratch/raw_dom.json`.

=== Etapa 2.2 --- Normalização e Extração de Tokens
<etapa-2.2-normalização-e-extração-de-tokens>
- #strong[Consome:] `scratch/raw_dom.json` da Etapa 2.1.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/impeccable/SKILL.md")[`impeccable`] + Script de paleta OKLCH.
- #strong[Processamento:] Erradica 100% de emojis amadores, extrai tokens semânticos de espaçamento, tipografia e cores, e valida contrastes em conformidade com as normas WCAG.
- #strong[Gate:] Zero emojis detectados e contraste de cores aprovado pelo script de validação.
- #strong[Handoff:] Grava os artefatos `DESIGN.md` e `TOKENS.css`.

=== Etapa 2.3 --- Gate de Design & Human-Gate 02
<etapa-2.3-gate-de-design-human-gate-02>
- #strong[Consome:] `DESIGN.md` e `TOKENS.css`.
- #strong[Peças:] Controle `human-gate` + Script linter de design.
- #strong[Processamento:] Apresentação visual da paleta e layout propostos para aprovação do desenvolvedor.
- #strong[Gate:] Aval humano formal registrado no Studio ou terminal.
- #strong[Handoff:] Libera os tokens aprovados para a especificação do Sub-Loop 3.

#horizontalrule

== 3. Sub-Loop 3: Engenharia & Tríade de Planejamento
<sub-loop-3-engenharia-tríade-de-planejamento>
É o coração arquitetural do LegoLoop. Ele converte a intenção humana e o design visual em um contrato de engenharia matematicamente executável.

#ficha(
  ("Papel", "Especificação BDD, decomposição hierárquica e rota DAG"),
  ("Disparo", "Conclusão e aprovação dos Sub-Loops 1 e 2"),
  ("Entrada", "INTENT.md e opcionalmente DESIGN.md"),
  ("Saída Principal", "A Tríade: SPEC.md, PLAN.md e ROUTE.json"),
  ("Portão de Saída", "STOP-GATE HUMANO 03 (Auditoria da Tríade)"),
)
#esteira(
  no-claro("3.1 SPEC BDD", sub: "aidd-spec + boundaries"),
  no-claro("3.2 TICKETS", sub: "fatiamento vertical"),
  no("3.3 ROTA JSON", sub: "compilador DAG"),
  no("3.4 STOP-GATE", sub: "bloqueio pré-código"),
)
=== Etapa 3.1 --- Especificação BDD (SPEC)
<etapa-3.1-especificação-bdd-spec>
- #strong[Consome:] `docs/intents/INTENT.md` (+ `DESIGN.md` se presente).
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-spec/SKILL.md")[`aidd-spec`].
- #strong[Processamento:] Redige as premissas técnicas explícitas, define as fronteiras em 3 níveis (#emph[Sempre Fazer], #emph[Perguntar Antes], #emph[Nunca Fazer]), demarca as costuras de teste (#emph[Test Seams]) e formaliza os critérios de aceite binários `AC-E{N}.{M}` no formato rigoroso Dado/Quando/Então.
- #strong[Gate:] `file_exists('docs/specs/YYYY-MM-DD_SPEC_<TEMA>.md') && has_binary_criteria(...)`.
- #strong[Handoff:] Grava o documento de especificação técnica no disco.

=== Etapa 3.2 --- Decomposição Hierárquica em Tickets (PLAN)
<etapa-3.2-decomposição-hierárquica-em-tickets-plan>
- #strong[Consome:] `docs/specs/SPEC.md` da Etapa 3.1 #emph[\(Dependência estrita)].
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-tickets/SKILL.md")[`aidd-tickets`].
- #strong[Processamento:] Fatiamento vertical da especificação em uma hierarquia de tickets: Épicos (\>20k tokens de contexto), Micro-tickets (4k a 20k tokens para TDD pontual) e Nano-tickets (\<4k tokens para correções de lint).
- #strong[Gate:] `file_exists('docs/plans/YYYY-MM-DD_PLAN_<TEMA>.md')`.
- #strong[Handoff:] Grava o plano unificado com tabela estruturada.

=== Etapa 3.3 --- Compilação da Rota Executável (ROUTE)
<etapa-3.3-compilação-da-rota-executável-route>
- #strong[Consome:] `docs/plans/PLAN.md` e `docs/specs/SPEC.md`.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-orchestrate-pipeline-creator/SKILL.md")[`aidd-orchestrate-pipeline-creator`] + Script validador JSON Schema.
- #strong[Processamento:] Converte a tabela legível do plano em um mapa formal JSON com declaração de IDs únicos, dependências, modelo roteado (`code-pro`, `code-fast`, `code-nano`), esforço (`--effort`) e comandos de verificação mecânica.
- #strong[Gate:] Validação contra o schema oficial (`exit_code == 0`).
- #strong[Handoff:] Grava `docs/routes/YYYY-MM-DD_ROUTE_<TEMA>.json`.

=== Etapa 3.4 --- Human-Gate 03 (Auditoria da Tríade)
<etapa-3.4-human-gate-03-auditoria-da-tríade>
- #strong[Consome:] Tríade completa (`SPEC.md` + `PLAN.md` + `ROUTE.json`).
- #strong[Peças:] Controle `human-gate`.
- #strong[Processamento:] Parada obrigatória antes de tocar em qualquer arquivo da aplicação. O desenvolvedor valida se o plano reflete as decisões arquiteturais corretas.
- #strong[Gate:] Aprovação explícita do usuário.
- #strong[Handoff:] Destrava a inicialização da construção no Sub-Loop 4.

#horizontalrule

== 4. Sub-Loop 4: Construção TDD por Épico #emph[\(Orchestrix)]
<sub-loop-4-construção-tdd-por-épico-orchestrix>
É o sub-loop onde o código real de software é escrito. Ele executa cada micro-ticket isoladamente dentro de um ciclo rigoroso de Test-Driven Development.

#ficha(
  ("Papel", "Construção de código de produção orientado a testes unitários"),
  ("Disparo", "Aprovação formal da Tríade de Planejamento"),
  ("Entrada", "ROUTE.json e âncoras da SPEC.md"),
  ("Saída Principal", "Código de produção testado e TDD_REPORT.md"),
  ("Portão de Saída", "STOP-GATE POR ÉPICO (Revisão Humana de Código)"),
)
#esteira(
  no("4.1 RED", sub: "teste falha primeiro"),
  no-claro("4.2 GREEN", sub: "implementação real"),
  no("4.3 REFACTOR", sub: "linters e funclen"),
  no("4.4 STOP-GATE", sub: "auditoria do épico"),
)
=== Etapa 4.1 --- Ciclo RED (Teste Falhando Primeiro)
<etapa-4.1-ciclo-red-teste-falhando-primeiro>
- #strong[Consome:] Seção âncora do critério de aceite na SPEC (`spec_anchor`) e arquivo de teste alvo.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-tdd/SKILL.md")[`aidd-tdd`] + Script executor de testes locais (`go test`).
- #strong[Processamento:] O agente implementa o teste unitário estrito reproduzindo fielmente o cenário do critério BDD.
- #strong[Gate:] O teste #strong[deve falhar comprovadamente] com erro de asserção funcional (`exit_code != 0`). Uma falha de compilação invalida o teste.
- #strong[Handoff:] Grava o código do teste e o log de falha em `scratch/test_red.log`.

=== Etapa 4.2 --- Ciclo GREEN (Implementação Mínima de Produção)
<etapa-4.2-ciclo-green-implementação-mínima-de-produção>
- #strong[Consome:] Teste falhando gerado na Etapa 4.1 #emph[\(Dependência estrita)].
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/implementacao/SKILL.md")[`implementacao`] + Script de compilação.
- #strong[Processamento:] Escreve o código de produção mínimo e limpo para fazer o teste unitário passar. Proibição absoluta de stubs, mocks falsos, `TODO` ou `pass`.
- #strong[Gate:] A suíte de teste daquele micro-ticket deve retornar código de sucesso (`exit_code == 0`).
- #strong[Handoff:] Salva o código de produção funcional no workspace do ticket.

=== Etapa 4.3 --- Ciclo REFACTOR & Quality Gates
<etapa-4.3-ciclo-refactor-quality-gates>
- #strong[Consome:] Código funcional aprovado na Etapa 4.2.
- #strong[Peças:] Scripts locais de qualidade: `go vet ./...`, `funclen` (nenhuma função \> 35 linhas), `gofmt -l .`.
- #strong[Processamento:] Limpeza de redundâncias, padronização idiomática de código e redução da complexidade ciclomática.
- #strong[Gate:] Zero violações de linter e saída vazia no verificador de formatação.
- #strong[Handoff:] Consolida os commits locais do micro-ticket no branch do Épico.

=== Etapa 4.4 --- Human-Gate por Épico
<etapa-4.4-human-gate-por-épico>
- #strong[Consome:] Conjunto de micro-tickets finalizados do Épico ativo.
- #strong[Peças:] Controle `human-gate`.
- #strong[Processamento:] Pausa obrigatória ao fechar um Épico. Apresenta o sumário de testes verdes e commits para auditoria humana antes de avançar para o próximo Épico.
- #strong[Gate:] Aprovação explícita do desenvolvedor.
- #strong[Handoff:] Encaminha para o próximo Épico ou, ao término de todos, para o Sub-Loop 5.

#horizontalrule

== 5. Sub-Loop 5: Homologação E2E, Recompilação & Fechamento
<sub-loop-5-homologação-e2e-recompilação-fechamento>
O último sub-loop garante que o sistema construído opera em perfeita harmonia no mundo real, atualiza os binários executáveis e encerra a rastreabilidade da sessão.

#ficha(
  ("Papel", "Homologação ponta a ponta, compilação de binários e fechamento"),
  ("Disparo", "Conclusão e aprovação de 100% dos Épicos de construção"),
  ("Entrada", "Código de produção completo no repositório"),
  ("Saída Principal", "Binário executável no PATH e REPORT.md"),
  ("Portão de Saída", "STOP-GATE FINAL (Homologação e Entrega)"),
)
#esteira(
  no("5.1 E2E SMOKE", sub: "Playwright headless"),
  no("5.2 RECOMPILAÇÃO", sub: "go build legoloop.exe"),
  no-claro("5.3 RELATÓRIO", sub: "sessão e métricas"),
  no("5.4 STOP-GATE", sub: "entrega homologada"),
)
=== Etapa 5.1 --- Homologação E2E Multi-Viewport & Smoke Test
<etapa-5.1-homologação-e2e-multi-viewport-smoke-test>
- #strong[Consome:] Aplicação viva rodando na porta local do ambiente de desenvolvimento.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-homologue/SKILL.md")[`aidd-homologue`] + MCP `playwright`.
- #strong[Processamento:] Navega autonomamente por todos os fluxos críticos de usuário em resoluções mobile, tablet e desktop, auditando erros de console HTTP e regressões visuais.
- #strong[Gate:] Zero erros de console HTTP e suíte E2E 100% verde.
- #strong[Handoff:] Grava o relatório de execução e telas de evidência em `scratch/e2e/`.

=== Etapa 5.2 --- Recompilação Obrigatória de Binários
<etapa-5.2-recompilação-obrigatória-de-binários>
- #strong[Consome:] Código-fonte final aprovado no repositório.
- #strong[Peças:] Script local de compilação: `go build -o bin/legoloop.exe ./cmd/legoloop`.
- #strong[Processamento:] Compilação nativa de release do executável para garantir que as alterações no motor Go estejam ativas no PATH do sistema operacional.
- #strong[Gate:] `exit_code == 0 && file_exists('bin/legoloop.exe')`.
- #strong[Handoff:] Atualiza o binário oficial no diretório `bin/`.

=== Etapa 5.3 --- Consolidação de Relatório & Sincronização de Sessão
<etapa-5.3-consolidação-de-relatório-sincronização-de-sessão>
- #strong[Consome:] Logs de execução, telemetria de tokens de todos os blocos e status dos gates.
- #strong[Peças:] Skill #link("file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-sessao/SKILL.md")[`aidd-sessao`] + Script CLI `faz-session`.
- #strong[Processamento:] Consolida métricas de tokens de entrada e saída, custo financeiro aproximado, modelos utilizados e balanço de conformidade contra a SPEC inicial.
- #strong[Gate:] `file_exists('docs/reports/YYYY-MM-DD_REPORT_<TEMA>.md')`.
- #strong[Handoff:] Grava o relatório final consolidado em `docs/reports/` e salva o histórico da sessão.

=== Etapa 5.4 --- Human-Gate Final (Entrega do Produto)
<etapa-5.4-human-gate-final-entrega-do-produto>
- #strong[Consome:] Relatório final consolidado e binário compilado.
- #strong[Peças:] Controle `human-gate`.
- #strong[Processamento:] Apresentação formal de encerramento da entrega para o cliente ou líder técnico.
- #strong[Gate:] Aceite final assinado pelo usuário humano.
- #strong[Handoff:] Conclusão bem-sucedida do Mega Pipeline.

== Rastreabilidade
<rastreabilidade-2>
- Esquema canônico de nós: `internal/dag/types.go`.
- Implementação do runner e compilação de gates: `internal/gate/gate.go`.
- Manual mestre dos sub-loops: `docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`.

= PARTE IV --- TRANSVERSAL: TEMPLATES E COMPOSIÇÃO FRACTAL
<parte-iv-transversal-templates-e-composição-fractal>
Os templates do LegoLoop não são sugestões estéticas ou formulários vazios; eles são os #strong[contratos fundacionais] que garantem que o sistema opere como uma máquina de estados determinística. Se uma ferramenta ou modelo não consegue interpretar o formato de um template, aquele arquivo deixa de ser engenharia e vira mera decoração.

#horizontalrule

== 1. A Cascata Canônica em Quatro Estágios
<a-cascata-canônica-em-quatro-estágios>
Toda demanda que entra no ecossistema passa compulsoriamente por quatro estágios de refinamento sucessivo antes que a primeira linha de código de produção seja autorizada:

#esteira(
  no-claro("1. INTENT", sub: "O Porquê"),
  no-claro("2. SPEC", sub: "O Quê"),
  no-claro("3. PLAN", sub: "O Como"),
  no("4. ROUTE", sub: "A Máquina"),
)
#figure(
  align(center)[#table(
    columns: (20.19%, 21.15%, 38.46%, 20.19%),
    align: (left,left,left,left,),
    table.header([Estágio de Governança], [Arquivo de Template], [Função Primária], [Validação no Motor],),
    table.hline(),
    [#strong[\1. INTENT]], [`TEMPLATE_INTENT.md`], [Registra a dor do usuário, hipóteses e premissas], [`Human-Gate 01`],
    [#strong[\2. SPEC]], [`TEMPLATE_SPEC.md`], [Define premissas, boundaries e critérios BDD `AC-XX`], [`has_binary_criteria`],
    [#strong[\3. PLAN]], [`TEMPLATE_PLAN.md`], [Fatiamento vertical hierárquico (Épicos ➔ Nano)], [`Human-Gate 02`],
    [#strong[\4. ROUTE]], [`TEMPLATE_ROUTE.json`], [Grafo acíclico executável em formato machine-readable], [`dag.Parse` e Schema],
  )]
  , kind: table
  )

#frase-destaque[O PLAN para os olhos humanos é o espelho exato da ROUTE para os nós de máquina. Divergência entre os dois é defeito impeditivo.]
#figure(image("/img/diagrama_triade_canonico.svg", alt: "A Tríade Canônica da Fabricação Autônoma"),
  caption: [
    A Tríade Canônica da Fabricação Autônoma
  ]
)

#painel("A Regra de Espelho Inegociável")[
  O PLAN (para humanos) e a ROUTE (para as máquinas) são duas faces da mesma moeda. Eles compartilham obrigatoriamente o mesmo conjunto de identificadores de nós, as mesmas dependências causais, os mesmos limites de timeout e o mesmo roteamento cognitivo. Qualquer divergência entre o PLAN e a ROUTE constitui defeito impeditivo de compilação.
]

#horizontalrule

== 2. A Composição Fractal via CLI: O Comando `compor`
<a-composição-fractal-via-cli-o-comando-compor>
A grande inovação arquitetural do LegoLoop é a capacidade de transformar qualquer processo isolado em uma etapa acoplável de um processo maior.

Imagine que você fabricou um processo chamado `subloop-tdd` contendo o ciclo completo de Red-Green-Refactor. Quando você estiver desenhando um novo sistema corporativo, você não precisa redesenhar essas etapas:

#ficha(
  ("Comando CLI", "legoloop compor <alvo> <subloop> [--pos N]"),
  ("Ação no DAG", "Insere nó do tipo 'subloop' com ref ao loop filho"),
  ("Prevenção de Ciclos", "cycle.DetectCycle(): recusa loops recursivos"),
  ("Isolamento de I/O", "Separação de caminhos via sub_artifact"),
  ("Versionamento", "Elevação semântica automática de versão (ex: 1.0.0 -> 1.1.0)"),
)
=== Como a Composição Funciona nos Bastidores
<como-a-composição-funciona-nos-bastidores>
+ #strong[Validação Estrutural:] O comando verifica se ambos os loops existem no cofre (`~/.aidd/loops/` ou repositório) e confere as assinaturas de entrada e saída.
+ #strong[Prevenção Mecânica de Ciclos:] O resolvedor de grafos simula a injeção. Se o loop `A` compuser o loop `B`, e o loop `B` já contiver referências ao loop `A`, o sistema aborta com erro fatal `ErrCiclo`.
+ #strong[Mapeamento de Artefatos com `sub_artifact`:] Para impedir que o artefato final do sub-loop colida com o artefato do pipeline pai (um bug clássico em geradores de código), o motor grava o campo `sub_artifact`. O sub-loop opera em sua pasta própria (`out/<subloop>/`) e a Orchestrix entrega apenas o resultado final no destino do pai.

#horizontalrule

== 3. O Ciclo de Vida dos Cinco Verbos de Processo
<o-ciclo-de-vida-dos-cinco-verbos-de-processo>
O LegoLoop padronizou a gestão de processos no ecossistema através de cinco verbos de comando universais, disponíveis tanto na CLI quanto na interface visual do Studio:

#figure(
  align(center)[#table(
    columns: (18.28%, 38.71%, 43.01%),
    align: (left,left,left,),
    table.header([Verbo Operacional], [Comando CLI], [Papel no Ciclo de Vida],),
    table.hline(),
    [#strong[Fabricar]], [`legoloop fabricar "<nome>"`], [Cria um processo novo do zero via cascata socrática],
    [#strong[Compor]], [`legoloop compor "<alvo>" "<subloop>"`], [Acopla um processo como sub-etapa de outro],
    [#strong[Listar]], [`legoloop listar`], [Audita o inventário de processos, versões e vínculos],
    [#strong[Refinar]], [`legoloop refinar "<nome>"`], [Substitui atomicamente um nó defasado por um melhor],
    [#strong[Remover]], [`legoloop remover "<nome>"`], [Desvincula um processo com garantia de rollback seguro],
  )]
  , kind: table
  )

#painel("Canonização e Hashing SHA-256")[
  Toda peça ou template importado ou fabricado passa imediatamente pelo crivo da canonização (legoloop canonize). O motor calcula a impressão digital SHA-256 de cada ativo, desduplicando cópias idênticas e vinculando a reputação daquela peça no cofre central do desenvolvedor.
]
== Rastreabilidade
<rastreabilidade-3>
- Implementação do comando compor: `internal/compor/compor.go` e testes em `internal/cli/compor_binary_test.go`.
- Validação de sub-artefatos sem colisão: `internal/dag/subartifact_test.go`.
- Contratos e templates SSOT: `docs/templates/templates.go`.

= APÊNDICES: RASTREABILIDADE, GLOSSÁRIO E ESTADO HONESTO
<apêndices-rastreabilidade-glossário-e-estado-honesto>
== Apêndice A --- Matriz de Rastreabilidade Factual
<apêndice-a-matriz-de-rastreabilidade-factual>
Para garantir a auditabilidade total desta obra, a tabela abaixo mapeia cada conceito arquitetural apresentado diretamente ao código-fonte ou documento normativo correspondente no repositório:

#figure(
  align(center)[#table(
    columns: (20.83%, 46.67%, 32.5%),
    align: (left,left,left,),
    table.header([Conceito ou Mecanismo], [Arquivo de Implementação no LegoLoop], [Responsabilidade Factual],),
    table.hline(),
    [#strong[Contrato do Nó DAG]], [`internal/dag/types.go`], [Estrutura do nó, modos de roteamento e handoffs],
    [#strong[Parser e Schema]], [`internal/dag/parse.go`], [Validação sintática contra `schemas/dag.v1.json`],
    [#strong[Prevenção de Ciclos]], [`internal/dag/cycle.go`], [Algoritmo de DFS que impede grafos circulares],
    [#strong[Isolamento de Sub-Loop]], [`internal/orchestrix/subloop.go`], [Workspaces temporários herméticos e `transferInputs`],
    [#strong[Compilação de Gates]], [`internal/gate/gate.go`], [Compilador de expressões booleanas em DSL],
    [#strong[Composição de Loops]], [`internal/compor/compor.go`], [Injeção de nós sub-loop e elevação de versão],
    [#strong[Sub-Artefatos sem Colisão]], [`internal/dag/subartifact_test.go`], [Isolamento de saídas entre loops pai e filho],
    [#strong[Templates Canônicos]], [`docs/templates/README.md`], [Modelos oficiais da sextilha de peças],
    [#strong[Mega Pipeline Master]], [`loops/mega-pipeline.loop.json`], [Manifesto canônico da esteira dos 5 sub-loops],
    [#strong[Manual dos Sub-Loops]], [`docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`], [Especificação operacional e stop-gates],
  )]
  , kind: table
  )

#horizontalrule

== Apêndice B --- Glossário Terminológico
<apêndice-b-glossário-terminológico>
- #strong[Cascata Canônica em 4 Estágios:] Sequência mandatória de planejamento (`INTENT` ➔ `SPEC` ➔ `PLAN` ➔ `ROUTE`) que precede a escrita de qualquer código de produção.
- #strong[Context Bloat:] Degradação severa da qualidade de raciocínio de um modelo provocada pelo acúmulo desnecessário de históricos de mensagens e arquivos não relacionados no prompt.
- #strong[Directed Acyclic Graph (DAG):] Estrutura matemática de dados unidirecional e sem laços recursivos, utilizada para orquestrar etapas de processos com dependências lineares ou paralelas.
- #strong[Effort Clamping:] Imposição de um teto rígido de tokens de raciocínio interno (`--effort low|med|high`) a uma chamada de IA, prevenindo gasto excessivo de cota em tarefas simples.
- #strong[Gate Mecânico:] Expressão booleana avaliada deterministicamente pelo motor local sem intervenção de IA (ex.: `exit_code == 0 && file_exists(...)`).
- #strong[Handoff File-Based:] Transferência de estado entre etapas realizada exclusivamente através da gravação e leitura de arquivos físicos no disco, garantindo reprodutibilidade e isolamento.
- #strong[Human-Gate (Stop-Gate):] Ponto de controle onde o motor paralisa fisicamente a execução e exige uma decisão humana explícita antes de liberar o próximo estágio.
- #strong[Mega Pipeline:] Macro-orquestrador do ecossistema que conecta múltiplos sub-loops para levar uma demanda do briefing à entrega final de produção.
- #strong[Orchestrix:] O motor hermético de execução do LegoLoop, responsável por criar workspaces temporários, filtrar entradas e destruir ambientes após o término de cada nó.
- #strong[Sextilha Canônica:] O conjunto das seis naturezas padronizadas de peças reutilizáveis no LegoLoop: Skills, Scripts, Hooks, MCPs, Gates e Loops.
- #strong[Sub-Loop:] Sub-grafo acíclico autocontido que encapsula uma linha de montagem especializada (ex.: ciclo TDD ou engenharia reversa).

#horizontalrule

== Apêndice C --- Estado Honesto do Ecossistema
<apêndice-c-estado-honesto-do-ecossistema>
Em consonância com a Lei da Honestidade de Rotulagem, declaramos formalmente o estado real e as limitações conhecidas da implementação na data de publicação desta obra:

#figure(
  align(center)[#table(
    columns: (20.37%, 13.89%, 26.85%, 38.89%),
    align: (left,left,left,left,),
    table.header([Recurso / Componente], [Status Atual], [Evidência de Funcionamento], [Limitação Conhecida],),
    table.hline(),
    [#strong[Validação de DAG]], [Concluído], [100% verde em `go test ./...`], [Suporta apenas grafos estáticos em tempo de carga],
    [#strong[Isolamento Orchestrix]], [Concluído], [Provado em testes temporários], [No Windows, remoção de pasta pode exigir retentativas se aberta],
    [#strong[Composição de Loops]], [Concluído], [Provado no ticket TICK-03.4], [O sub-loop reescreve apenas o gate do artefato final],
    [#strong[Templates Canônicos V2]], [Concluído], [Testado em `templates_test.go`], [Scorer de reputação ainda avalia estrutura por regex],
    [#strong[Execução com LLM Real]], [Operacional], [Testado com Claude e 9Router], [Exige chave configurada e conexão ativa de rede],
  )]
  , kind: table
  )
