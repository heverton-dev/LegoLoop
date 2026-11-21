// =============================================================================
// TEMPLATE_BOOK.typst — PADRÃO EDITORIAL STRIPE PRESS / LEGOPRESS (170x240mm)
// Direção de arte: Capa branca padronizada com detalhe arquitetural Google,
// ilustração vetorial modular em SVG transparente, miolo arejado sem margem de
// pauta espremida, tipografia Georgia no texto, JetBrains Mono no código e
// cards semânticos com borda colorida na paleta Google Suavizada.
// =============================================================================

#let cor-fundo    = rgb("#ffffff")                                    // Branco puro limpo
#let cor-miolo    = rgb("#fafafa")                                    // Off-white sutil para leitura confortável
#let cor-tinta    = rgb("#0f172a")                                    // Ardósia escuro nobre
#let cor-texto    = rgb("#292524")                                    // Grafite quente stone
#let cor-suave    = rgb("#64748b")                                    // Cinza pedra suave
#let cor-linha    = rgb("#e2e8f0")                                    // Linha divisória sutil
#let cor-accent   = rgb("#16a34a") // Cor da obra (Google Cobalto default)
#let cor-accent2  = cor-accent.darken(25%)                            // Tom profundo da cor da obra

// Paleta Semântica Oficial Google Suavizada
#let cor-info     = rgb("#2563eb")                                    // Azul Cobalto Google
#let cor-alerta   = rgb("#d97706")                                    // Âmbar Google Dourado
#let cor-perigo   = rgb("#dc2626")                                    // Carmesim Google
#let cor-erro     = cor-perigo                                        // Alias canônico para cor de erro
#let cor-sucesso  = rgb("#16a34a")                                    // Verde Esmeralda Google
#let cor-teoria   = rgb("#7c3aed")                                    // Violeta DeepMind

// Helpers de Nós de Diagramas
#let no(titulo, sub: "") = box(stroke: 0.5pt + cor-linha, inset: 5pt, radius: 3pt, fill: white)[#text(weight: "bold", size: 8pt, fill: cor-tinta)[#titulo]#if sub != "" [ \ #text(size: 7pt, fill: cor-suave)[#sub] ]]

// Componente Canônico Esteira de Processo / Pipeline
#let esteira(estagios, gates: ()) = {
  v(1em)
  block(width: 100%, stroke: 0.5pt + cor-linha, inset: 10pt, radius: 4pt, fill: rgb("#F8FAFC"))[
    #text(weight: "bold", size: 8.5pt, fill: cor-tinta, font: ("Segoe UI", "Arial"))[ESTEIRA CANÔNICA DE PROCESSO]
    #v(0.5em)
    #grid(
      columns: (1fr,) * calc.min(estagios.len(), 4),
      gutter: 8pt,
      ..estagios.map(item => {
        let tit = if type(item) == array and item.len() > 0 { item.at(0) } else { str(item) }
        let sub = if type(item) == array and item.len() > 1 { item.at(1) } else { "" }
        box(width: 100%, stroke: 0.5pt + rgb("#CBD5E1"), inset: 6pt, radius: 3pt, fill: white)[
          #text(weight: "bold", size: 8pt, fill: cor-tinta, font: ("Segoe UI", "Arial"))[#tit]
          #if sub != "" [ \ #text(size: 7pt, fill: cor-suave, font: ("Segoe UI", "Arial"))[#sub] ]
        ]
      })
    )
    #if gates.len() > 0 [
      #v(0.4em)
      #text(size: 7.2pt, fill: cor-accent, font: ("Segoe UI", "Arial"), weight: "bold")[Gates Binários: ]
      #text(size: 7.2pt, fill: cor-suave, font: ("Segoe UI", "Arial"))[#gates.join(" · ")]
    ]
  ]
  v(1em)
}

#set document(
  title: [Verbos e Granularidade no LegoLoop],
  author: "Heverton & Ecossistema LegoLoop",
)

#set page(
  width: 170mm,
  height: 240mm,
  fill: cor-fundo,
  margin: (top: 2.0cm, bottom: 2.0cm, left: 1.8cm, right: 1.8cm),
)

// Tipografia: Georgia clássica no miolo com entrelinha arejada
#set text(
  font: ("Georgia", "Segoe UI", "Calibri"),
  size: 10.3pt,
  lang: "pt",
  region: "BR",
  fill: cor-texto,
  hyphenate: true,
)

#set par(justify: true, leading: 1.05em, spacing: 1.35em)
#set smartquote(enabled: false)

// ---------------------------------------------------------------------------
// Painéis Semânticos Oficiais (Borda Lateral Google e Fundo Confortável)
// ---------------------------------------------------------------------------

#let painel(titulo, corpo, cor: cor-accent, bg: rgb("#f8fafc")) = block(
  width: 100%, breakable: false, above: 16pt, below: 16pt,
  fill: bg, radius: 6pt,
  stroke: (left: 3.5pt + cor, rest: 0.6pt + cor-linha),
  inset: (x: 14pt, y: 12pt),
)[
  #text(font: ("Segoe UI", "Arial"), size: 8.5pt, weight: "bold", fill: cor, tracking: 1.2pt)[#upper(titulo)]
  #v(4pt)
  #set text(font: ("Georgia", "Segoe UI"), size: 9.6pt, fill: cor-texto)
  #set par(leading: 0.9em, spacing: 1em)
  #corpo
]

#let painel-info(titulo, corpo)    = painel(titulo, corpo, cor: cor-info, bg: rgb("#f8fafc"))
#let painel-alerta(titulo, corpo)  = painel(titulo, corpo, cor: cor-alerta, bg: rgb("#fffbeb"))
#let painel-perigo(titulo, corpo)  = painel(titulo, corpo, cor: cor-perigo, bg: rgb("#fef2f2"))
#let painel-sucesso(titulo, corpo) = painel(titulo, corpo, cor: cor-sucesso, bg: rgb("#f0fdf4"))
#let painel-teoria(titulo, corpo)  = painel(titulo, corpo, cor: cor-teoria, bg: rgb("#faf5ff"))

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
// Títulos de Seção — APENAS TEXTO LIMPO
// ---------------------------------------------------------------------------

#set heading(numbering: none)
#show heading: set text(font: ("Segoe UI", "Arial"))

// Contador de Partes para Rotação Harmônica dos Blocos LEGO
#let parte-counter = counter("contador-partes")

// Títulos Nível 1 (Partes) — PÁGINA EXCLUSIVA DEDICADA COM TIJOLO LEGO 3D AUTÊNTICO
#show heading.where(level: 1): it => {
  let is-part = repr(it.body).contains("Parte") or repr(it.body).contains("PARTE") or repr(it.body).contains("Part")
  if is-part {
    parte-counter.step()
    context {
      let p-num = parte-counter.get().first()
      let svg-path = if p-num == 1 { "ativos/lego_brick_part1.svg" }
        else if p-num == 2 { "ativos/lego_brick_part2.svg" }
        else if p-num == 3 { "ativos/lego_brick_part3.svg" }
        else { "ativos/lego_brick_part4.svg" }

      let accent-color = if p-num == 1 { cor-sucesso }
        else if p-num == 2 { cor-info }
        else if p-num == 3 { cor-alerta }
        else { cor-teoria }

      pagebreak(weak: true)
      align(center + horizon)[
        #block(width: 7.2cm)[
          #image(svg-path, width: 100%)
        ]
        #v(22pt)
        #text(
          font: ("Segoe UI", "Georgia", "Arial"),
          size: 20pt,
          weight: "bold",
          fill: cor-tinta,
          tracking: 0.2pt,
        )[#it.body]
        #v(14pt)
        #line(length: 3.5cm, stroke: 1.5pt + accent-color)
      ]
      pagebreak(weak: true)
    }
  } else {
    pagebreak(weak: true)
    block(width: 100%, above: 0pt, below: 22pt)[
      #text(size: 19pt, weight: "black", fill: cor-tinta)[#it.body]
      #v(6pt)
      #line(length: 100%, stroke: 1.5pt + cor-linha)
    ]
  }
}

// Nível 2 (Capítulos e Receitas) — Corrido com linha discreta
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
// Código: JetBrains Mono com Formatação Suave e Quebra Suave
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
)[#text(font: ("JetBrains Mono", "Consolas", "Courier New"), size: 8.0pt, fill: cor-accent2, hyphenate: false)[#_quebravel(it.text)]]

#show raw.where(block: true): it => block(
  width: 100%, fill: rgb("#f8fafc"),
  stroke: (left: 3pt + cor-accent, rest: 0.6pt + cor-linha),
  radius: 4pt, inset: (x: 11pt, y: 10pt), above: 13pt, below: 13pt, breakable: true,
)[#text(font: ("JetBrains Mono", "Consolas", "Courier New"), size: 8.0pt, fill: cor-tinta)[#it]]

#show quote.where(block: true): it => block(
  width: 100%, fill: rgb("#f8fafc"), stroke: (left: 3pt + cor-accent),
  inset: (x: 13pt, y: 10pt), radius: 4pt, above: 14pt, below: 14pt, breakable: true,
)[#set text(font: ("Georgia", "Segoe UI"), size: 10pt, fill: cor-suave, style: "italic"); #it.body]

#set list(indent: 8pt, body-indent: 7pt, spacing: 0.88em)
#set enum(indent: 8pt, body-indent: 7pt, spacing: 0.88em)

// ---------------------------------------------------------------------------
// Tabelas: Largura Controlada, Proteção Anti-Overflow e Tipografia Otimizada
// ---------------------------------------------------------------------------

#set table(
  stroke: 0.35pt + cor-linha,
  inset: (x: 5pt, y: 3.2pt),
  fill: (_, y) => if y == 0 { rgb("#f1f5f9") } else if calc.odd(y) { rgb("#f8fafc") } else { white },
)
#show table.cell.where(y: 0): set text(font: ("Segoe UI", "Arial"), fill: cor-tinta, weight: "bold", size: 6.8pt, hyphenate: false)
#show table: set text(size: 6.5pt, hyphenate: false)
#show table: set par(justify: false, leading: 0.50em)
#show table: it => {
  show raw.where(block: false): r => text(font: ("JetBrains Mono", "Consolas"), size: 6.0pt, fill: cor-accent2, hyphenate: false)[#r.text]
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


// ---------------------------------------------------------------------------
// CAPA — DESIGN STRIPE PRESS PADRONIZADO COM ILUSTRAÇÃO SVG TRANSPARENTE
// ---------------------------------------------------------------------------
#page(margin: (x: 1.8cm, top: 1.6cm, bottom: 1.5cm), numbering: none, header: none, footer: none, fill: white)[

  // Topo: 4 traços finos na paleta Google Suavizada
  #grid(
    columns: (1fr, 1fr, 1fr, 1fr),
    gutter: 5pt,
    line(length: 100%, stroke: 2pt + cor-info),
    line(length: 100%, stroke: 2pt + cor-perigo),
    line(length: 100%, stroke: 2pt + cor-alerta),
    line(length: 100%, stroke: 2pt + cor-sucesso),
  )

  #v(8pt)

  // Categoria editorial
  #text(
    font: ("Segoe UI", "Arial"),
    size: 7.5pt,
    fill: cor-suave,
    weight: "bold",
    tracking: 2.0pt,
  )[GUIA DIDÁTICO · BEHAVIOR-DRIVEN DEVELOPMENT]

  #v(6pt)

  // Título Principal
      #text(
      font: ("Segoe UI", "Arial"),
      size: 21pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.4pt,
    )[Verbos & BDD]
    #v(2pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 13pt,
      weight: "bold",
      fill: cor-accent,
    )[Granularidade e Tríade de Triggers no LegoLoop]
  
  #v(4pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 8.8pt,
      fill: cor-suave,
      weight: "medium",
    )[Guia Didático BDD da Tríade de Triggers e do Encadeamento de Verbos]
  ]

  #v(6pt)

  // Ilustração Vetorial da Capa: Peças de Lego 3D modulares autênticas conectadas
  #align(center)[
    #block(width: 80%, inset: (y: 2pt))[
      #image("ativos/lego_cover_blocks.svg", width: 75%)
    ]
  ]

  #v(6pt)

  // Esteira de Estágios na Capa (cover-stages)
    #block(width: 100%, radius: 5pt, clip: true, fill: rgb("#f8fafc"), stroke: 0.5pt + rgb("#cbd5e1"), inset: 8pt)[
    #let etapas = ("Fabricar", "Matriz", "Listar", "Refinar", "Remover", "Compor", )
    #text(font: ("Segoe UI", "Arial"), size: 6.8pt, weight: "bold", fill: cor-suave, tracking: 1.4pt)[OS CINCO VERBOS DO CICLO DE VIDA]
    #v(4pt)
    #grid(
      columns: (1fr,) * 2, gutter: 5pt,
      ..etapas.enumerate().map(((i, e)) => block(
        width: 100%, height: 0.95cm, radius: 3.5pt,
        fill: if calc.even(i) { cor-accent.lighten(92%) } else { white },
        stroke: (left: 2.8pt + cor-accent, rest: 0.4pt + cor-linha),
        inset: (x: 7pt, y: 3pt),
        align(left + horizon)[
          #text(font: ("Segoe UI", "Arial"), size: 6.5pt, weight: "bold", fill: cor-accent)[#if i < 9 [0#str(i + 1)] else [#str(i + 1)]]
          #h(3pt)
          #text(font: ("Segoe UI", "Arial"), size: 7.0pt, weight: "bold", fill: cor-tinta)[#e]
        ]
      ))
    )
  ]
  
  #v(1fr)

  // Rodapé da Capa com Proteção Anti-Quebra
  #line(length: 100%, stroke: 0.6pt + cor-linha)
  #v(4pt)
  #grid(
    columns: (1fr, auto),
    align: (left + horizon, right + horizon),
    gutter: 10pt,
    text(font: ("Segoe UI", "Arial"), size: 7.8pt, weight: "bold", fill: cor-tinta, hyphenate: false)[Heverton & Ecossistema LegoLoop],
    text(font: ("Segoe UI", "Arial"), size: 7.2pt, fill: cor-suave, hyphenate: false)[07/11/2026 · LegoLoop Foundation]
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
  Este guia ensina a operar o LegoLoop por meio de seus verbos canônicos e dos cinco níveis de granularidade que organizam qualquer fluxo: do Mega Pipeline às Nano-Etapas.

  A obra percorre três camadas. No nível Macro, apresenta a filosofia de blocos componíveis governados por um grafo acíclico dirigido e a Tríade de Triggers (Terminal, Chat do Agente e LegoLoop Studio). No nível Meso, documenta as cinco trilhas canônicas de encadeamento de verbos. No nível Micro, especifica sete cenários BDD em que o mesmo comportamento é acionado pelas três vias, com resultado idêntico e auditável.

  
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
    Data da compilação: 07/11/2026 · LegoLoop Foundation
  ]
]

// ---------------------------------------------------------------------------
// CORPO: CABEÇALHO LIMPO E RODAPÉ AREJADO (SEM COLUNA DE ANOTAÇÕES ESTREITA)
// ---------------------------------------------------------------------------
#set page(
  fill: cor-miolo,
  numbering: "1",
  margin: (top: 2.0cm, bottom: 2.0cm, left: 1.8cm, right: 1.8cm),
  header: context {
    if here().page() > 2 {
      let pagina = here().page()
      let partes-nesta = query(heading.where(level: 1)).filter(h => h.location().page() == pagina and (repr(h.body).contains("Parte") or repr(h.body).contains("PARTE")))
      if partes-nesta.len() == 0 {
        let anteriores = query(selector(heading.where(level: 1)).before(here()))
        let titulo-cap = if anteriores.len() > 0 { anteriores.last().body } else { [Sumário] }
        stack(
          spacing: 2.5pt,
          grid(
            columns: (auto, 1fr),
            gutter: 10pt,
            align: (left + horizon, right + horizon),
            text(font: ("Segoe UI", "Arial"), size: 6.2pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[VERBOS BDD],
            text(font: ("Segoe UI", "Arial"), size: 6.2pt, fill: cor-suave, hyphenate: false)[#titulo-cap],
          ),
          line(length: 100%, stroke: 0.5pt + cor-linha),
        )
      }
    }
  },
  footer: context {
    if here().page() > 2 {
      let pagina = here().page()
      let partes-nesta = query(heading.where(level: 1)).filter(h => h.location().page() == pagina and (repr(h.body).contains("Parte") or repr(h.body).contains("PARTE")))
      if partes-nesta.len() == 0 {
        align(center)[
          #text(font: ("Segoe UI", "Arial"), size: 8.2pt, fill: cor-suave)[#counter(page).display("1")]
        ]
      }
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


= Como ler este guia
<como-ler-este-guia>
== Público e pré-requisitos
<público-e-pré-requisitos>
O guia se destina a desenvolvedores, engenheiros de IA e arquitetos de software que operam o ecossistema LegoLoop e AIDD. Pressupõe familiaridade com linha de comando e com a noção de pipeline; nenhum conhecimento prévio do motor interno é exigido.

== Estrutura da obra
<estrutura-da-obra>
#figure(
  align(center)[#table(
    columns: (20%, 20%, 60%),
    align: (left,left,left,),
    table.header([Parte], [Nível], [Conteúdo],),
    table.hline(),
    [#strong[I]], [Macro], [Filosofia Lego, níveis de granularidade e Tríade de Triggers],
    [#strong[II]], [Meso], [Matriz de encadeamento e as cinco trilhas canônicas],
    [#strong[III]], [Micro], [Sete cenários BDD multi-trigger],
    [#strong[IV]], [Referência], [Tabela de consulta rápida e síntese final],
  )]
  , kind: table
  )

== Convenção BDD
<convenção-bdd>
Cada cenário segue a forma #strong[Dado-Quando-Então]. O bloco #emph[Dado] fixa o estado inicial, o bloco #emph[Quando] descreve a ação nas três vias de acionamento, e o bloco #emph[Então] declara o resultado verificável. O comportamento contratado não depende da via escolhida.

= PARTE I --- Macro: Filosofia e Triggers
<parte-i-macro-filosofia-e-triggers>
== 1. A Filosofia Lego
<a-filosofia-lego>
O LegoLoop transforma fluxos de software e de inteligência artificial em #strong[blocos modulares componíveis], governados por um grafo acíclico dirigido (DAG).

Nenhum fluxo é código estático ou monolítico. Qualquer processo pode ser montado, desmontado, inspecionado, refinado ou recombinado sem efeitos colaterais nos demais blocos.

== 2. Os Cinco Níveis de Granularidade
<os-cinco-níveis-de-granularidade>
#figure(image("/img/diagrama_granularidade.svg", width: 88.0%, alt: "Os cinco níveis de granularidade, do Mega Pipeline às Nano-Etapas"),
  caption: [
    Os cinco níveis de granularidade, do Mega Pipeline às Nano-Etapas
  ]
)

#figure(
  align(center)[#table(
    columns: (15.53%, 23.3%, 34.95%, 26.21%),
    align: (left,left,left,left,),
    table.header([Nível], [Escopo], [O que se controla], [Exemplo real],),
    table.hline(),
    [#strong[\1. Mega Pipeline]], [Ciclo de vida do produto], [Orquestração e paralelismo de épicos], [Mega Pipeline ponta a ponta],
    [#strong[\2. Sub-Loops]], [Módulos reutilizáveis], [Injeção hierárquica, saída isolada], [Loop de Auditoria 4F],
    [#strong[\3. Etapas (Nós)]], [Contratos de I/O], [Troca de peças, dependências], [`e01` -\> `e02`],
    [#strong[\4. Micro-Etapas]], [Execução cognitiva], [Modelo de IA, contexto, MCPs], [Redator de testes isolado],
    [#strong[\5. Nano-Etapas]], [Invariantes de código], [Asserções, cobertura, gates `SG-*`], [Função Go com até 35 linhas],
  )]
  , kind: table
  )

== 3. A Tríade de Triggers
<a-tríade-de-triggers>
O LegoLoop não se limita à linha de comando. Todo verbo do ecossistema aceita #strong[três triggers equivalentes de entrada], e todos despacham para o mesmo motor canônico.

#figure(image("/img/diagrama_triggers.svg", width: 95.0%, alt: "A Tríade de Triggers alimentando o motor canônico"),
  caption: [
    A Tríade de Triggers alimentando o motor canônico
  ]
)

+ #strong[Terminal (CLI):] comando direto `legoloop <verbo> [flags]`. Indicado para scripts, CI/CD e automação sem interface.
+ #strong[Chat do Agente (Harness):] linguagem natural no editor (Claude Code, Cursor, Windsurf, MIMO, Antigravity) por meio das skills canônicas ou slash commands, como `/legoloop-fabricar` e `/legoloop-listar`. O agente interpreta o pedido e despacha a ação de forma estéril.
+ #strong[LegoLoop Studio (Interface Web):] painel visual aberto com `legoloop studio`. Indicado para inspeção de grafos DAG em tempo real, cartões, modais assistidos e drag-and-drop.

= PARTE II --- Meso: Encadeamento
<parte-ii-meso-encadeamento>
== 4. Verbos como Elos de uma Corrente
<verbos-como-elos-de-uma-corrente>
Os verbos do LegoLoop não foram projetados para rodar isolados. Cada um é elo de uma corrente lógica: os diagnósticos alimentam ações cirúrgicas, e as ações desembocam em revalidação e entrega.

#figure(image("/img/diagrama_encadeamento.svg", width: 95.0%, alt: "Fluxo entre diagnóstico, ação e validação"),
  caption: [
    Fluxo entre diagnóstico, ação e validação
  ]
)

#figure(
  align(center)[#table(
    columns: (32.53%, 21.69%, 45.78%),
    align: (left,left,left,),
    table.header([Origem], [Destino], [Finalidade da ligação],),
    table.hline(),
    [`listar`], [`refinar`], [Inspecionar os nós antes da troca],
    [`listar`], [`compor`], [Identificar sub-loops reutilizáveis],
    [`listar`], [`remover`], [Auditar resíduos e vínculos],
    [`matriz`], [`refinar`], [Detectar peças defasadas],
    [`matriz`], [`fabricar`], [Validar a aptidão das peças],
    [`fabricar` / `compor` / `remover`], [`listar`], [Revalidar grafo, acoplamento e limpeza],
    [`refinar`], [Auditoria de gates], [Validar o score após a troca],
    [Auditoria de gates], [`legopress`], [Gerar a documentação final],
  )]
  , kind: table
  )

== 5. As Cinco Trilhas Canônicas
<as-cinco-trilhas-canônicas>
As trilhas abaixo esgotam as combinações lógicas da operação diária.

=== Trilha 1: Criação e Publicação
<trilha-1-criação-e-publicação>
#strong[Objetivo:] criar um processo com governança completa sem duplicar trabalho.

#strong[Encadeamento:] `listar` $arrow.r$ `fabricar` $arrow.r$ `matriz`

+ #strong[`listar`:] verifica se já existe processo similar registrado no cofre.
+ #strong[`fabricar`:] gera os quatro artefatos contratuais (`INTENT`, `SPEC`, `PLAN`, `ROUTE`).
+ #strong[`matriz`:] confere se as peças selecionadas pelo motor têm boa reputação (Tier Ouro ou Prata).

=== Trilha 2: Diagnóstico e Refinamento Cirúrgico
<trilha-2-diagnóstico-e-refinamento-cirúrgico>
#strong[Objetivo:] localizar etapas de baixo desempenho e substituir peças defasadas.

#strong[Encadeamento:] `listar` $arrow.r$ `matriz` $arrow.r$ `refinar` $arrow.r$ `matriz`

+ #strong[`listar`:] localiza os identificadores exatos das etapas (`e01`, `e02`).
+ #strong[`matriz`:] filtra peças Bronze ou defasadas e mapeia substitutas Tier Ouro ($gt.eq 0.85$).
+ #strong[`refinar`:] troca atômica da peça na etapa, com backup em `.legoloop-bak` e bump minor de versão.
+ #strong[`matriz`:] revalida o processo e confirma que a saúde geral do cofre subiu.

=== Trilha 3: Composição Hierárquica e Reuso
<trilha-3-composição-hierárquica-e-reuso>
#strong[Objetivo:] reaproveitar um processo existente como sub-etapa de outro maior.

#strong[Encadeamento:] `listar` $arrow.r$ `compor` $arrow.r$ `listar`

+ #strong[`listar`:] seleciona o processo anfitrião (pai) e o sub-loop a ser acoplado (filho).
+ #strong[`compor`:] conecta o sub-loop em posição específica (`legoloop compor <pai> <filho> --pos N`).
+ #strong[`listar`:] valida a árvore hierárquica e os links do DAG na saída JSON.

=== Trilha 4: Ciclo de Vida e Descomissionamento
<trilha-4-ciclo-de-vida-e-descomissionamento>
#strong[Objetivo:] desativar processos depreciados, limpando symlinks de todos os editores sem quebra.

#strong[Encadeamento:] `matriz` $arrow.r$ `listar` $arrow.r$ `remover --dry-run` $arrow.r$ `remover --force`

+ #strong[`matriz`:] identifica processos com peças órfãs ou critérios descontinuados.
+ #strong[`listar`:] analisa o raio de impacto (harnesses e editores vinculados).
+ #strong[`remover --dry-run`:] simula a exclusão para confirmar que nenhum artefato crítico será afetado.
+ #strong[`remover --force`:] executa o descarte com limpeza de symlinks e histórico de rollback preservado.

=== Trilha 5: Governança Editorial e Consolidação Didática
<trilha-5-governança-editorial-e-consolidação-didática>
#strong[Objetivo:] transformar especificações, planos e rotas refinadas em livro técnico ou manual em PDF.

#strong[Encadeamento:] `listar` $arrow.r$ `matriz` $arrow.r$ `legopress`

+ #strong[`listar` e `matriz`:] auditam a estabilidade e a integridade dos nós do processo.
+ #strong[`legopress`:] compila diffs de código, especificações BDD e diagramas em documento Typst/PDF canônico.

== 6. Tabela-Resumo dos Encadeamentos
<tabela-resumo-dos-encadeamentos>
#figure(
  align(center)[#table(
    columns: (21.7%, 32.08%, 20.75%, 25.47%),
    align: (left,left,left,left,),
    table.header([Trilha], [Sequência de verbos], [Gatilho típico], [Resultado entregue],),
    table.hline(),
    [#strong[\1. Criação]], [listar → fabricar → matriz], [Nova demanda], [INTENT, SPEC, PLAN e ROUTE],
    [#strong[\2. Refinamento]], [listar → matriz → refinar → matriz], [Peça lenta ou defasada], [Etapa trocada, bump minor],
    [#strong[\3. Composição]], [listar → compor → listar], [Reuso de fluxo], [Sub-loop isolado no DAG],
    [#strong[\4. Descomissionamento]], [matriz → listar → remover], [Processo legado], [Symlinks limpos, sem órfãos],
    [#strong[\5. Fechamento editorial]], [matriz → legopress], [Fim de ciclo], [Manual e PDF canônico],
  )]
  , kind: table
  )

= PARTE III --- Micro: Cenários BDD
<parte-iii-micro-cenários-bdd>
== 7. Como Ler os Cenários
<como-ler-os-cenários>
Cada cenário demonstra #strong[o mesmo comportamento esperado] acionado pelas três vias. A seção #emph[Quando] lista Terminal, Chat do Agente e Studio; as seções #emph[Então] e #emph[E] valem igualmente para as três.

== 8. Cenário 1: Fabricar um Pipeline do Zero
<cenário-1-fabricar-um-pipeline-do-zero>
#quote(block: true)[
#strong[Objetivo:] o sistema cria automaticamente toda a estrutura de governança para uma nova necessidade de automação.
]

#strong[Dado que] o desenvolvedor precisa de um pipeline para auditar a segurança de contratos inteligentes em Solidity,

#strong[Quando] aciona a fabricação:

#emph[Via Terminal (CLI):]

```bash
legoloop fabricar "auditoria-smart-contracts" \
  --goal "Auditar vulnerabilidades reentrancy e overflow em Solidity" \
  --category build --yes
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-fabricar
Crie um novo processo chamado 'auditoria-smart-contracts' focado em auditar
vulnerabilidades de reentrancy e overflow em Solidity, na categoria build.
```

#emph[Via LegoLoop Studio:]

+ Na aba #strong[Processos], clica em #strong[\+ Novo Processo].
+ No modal, preenche #strong[Nome] `auditoria-smart-contracts`, #strong[Categoria] `build` e #strong[Objetivo].
+ Clica em #strong[Fabricar Processo].

#strong[Então] o motor hermético produz quatro artefatos canônicos, nomeados no padrão `YYYY-MM-DD_<TIPO>_AUDITORIA_SMART_CONTRACTS`:

#figure(
  align(center)[#table(
    columns: (18.18%, 22.08%, 18.18%, 41.56%),
    align: (left,left,left,left,),
    table.header([Artefato], [Pasta], [Extensão], [Conteúdo],),
    table.hline(),
    [INTENT], [`docs/intents/`], [`.md`], [Problema e motivação],
    [SPEC], [`docs/specs/`], [`.md`], [Critérios de aceitação binários],
    [PLAN], [`docs/plans/`], [`.md`], [Tickets em fatias verticais],
    [ROUTE], [`docs/routes/`], [`.json`], [DAG formal executável],
  )]
  , kind: table
  )

#strong[E] o pipeline fica pronto para execução sem código boilerplate manual.

== 9. Cenário 2: Inspecionar o Inventário com `listar`
<cenário-2-inspecionar-o-inventário-com-listar>
#quote(block: true)[
#strong[Objetivo:] saber quais processos existem no cofre, suas versões e a quais harnesses estão vinculados.
]

#strong[Dado que] existem vários loops criados na máquina,

#strong[Quando] o usuário solicita a listagem:

#emph[Via Terminal (CLI):]

```bash
legoloop listar
legoloop listar --json   # saída estruturada para pipelines
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-listar
Liste todos os processos ativos no cofre do LegoLoop e verifique se há
harnesses descalibrados ou vínculos órfãos.
```

#emph[Via LegoLoop Studio:]

+ Acessa a seção #strong[Processos] no menu lateral.
+ Visualiza a grade de cartões com selo de versão (`@1.0.0`), categoria e estado de cada harness vinculado.

#strong[Então] o sistema exibe nome e versão do processo, categoria (`build`, `reverse`, `web`, `mobile`), harnesses vinculados com estado dos links simbólicos e eventuais vínculos órfãos ou advertências de sincronização.

#strong[E] o desenvolvedor obtém visibilidade completa do ecossistema em qualquer uma das vias.

== 10. Cenário 3: Avaliar a Saúde das Peças com `matriz`
<cenário-3-avaliar-a-saúde-das-peças-com-matriz>
#quote(block: true)[
#strong[Objetivo:] garantir que o pipeline usa apenas peças modernas e de alta reputação antes de uma entrega crítica.
]

#strong[Dado que] o desenvolvedor quer auditar a integridade das peças do ecossistema,

#strong[Quando] consulta a matriz:

#emph[Via Terminal (CLI):]

```bash
legoloop matriz --tier ouro
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-matriz
Consulte a matriz de peças do LegoLoop e mostre apenas as ferramentas
de Tier Ouro com score >= 0.85 utilizáveis no pipeline.
```

#emph[Via LegoLoop Studio:]

+ Clica em #strong[Matriz de Peças] na navegação principal.
+ Seleciona o chip #strong[Padrão Ouro] no filtro de topo.
+ Inspeciona a tabela com scores de reputação, cobertura de testes de mutação e conformidade com zero stubs.

#strong[Então] o sistema exibe apenas componentes certificados: score $gt.eq 0.85$, sem stubs e com testes herméticos.

#strong[E] sinaliza as peças Bronze ou defasadas, recomendando substituição cirúrgica.

== 11. Cenário 4: Compor um Sub-Loop em Pipeline Existente
<cenário-4-compor-um-sub-loop-em-pipeline-existente>
#quote(block: true)[
#strong[Objetivo:] anexar um sub-loop de homologação visual já existente a um pipeline de construção.
]

#strong[Dado que] existem o pipeline mestre `pipeline-frontend` e o sub-loop `aidd-homologue@1.0.0`,

#strong[Quando] o usuário conecta o sub-loop após a etapa de build `e02`:

#emph[Via Terminal (CLI):]

```bash
legoloop compor pipeline-frontend aidd-homologue@1.0.0 --depende e02 --yes
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-compor
Componha o sub-loop 'aidd-homologue@1.0.0' dentro do 'pipeline-frontend',
fazendo-o depender da etapa 'e02'.
```

#emph[Via LegoLoop Studio:]

+ Abre `pipeline-frontend` no visualizador de DAG.
+ Clica em #strong[\+ Adicionar Sub-Loop] ou arrasta o bloco `aidd-homologue` da gaveta lateral.
+ Traça a aresta de `e02` até o nó do sub-loop.
+ Clica em #strong[Salvar Composição].

#strong[Então] o LegoLoop:

+ Verifica se a junção criaria ciclos proibidos (`ErrCiclo`).
+ Isola a saída do sub-loop em `out/aidd-homologue/`.
+ Atualiza atomicamente `.processo.json` e `.dag.json` do pipeline pai.

#strong[E] o sub-loop passa a executar como caixa-preta confiável dentro do macro-fluxo.

== 12. Cenário 5: Refinar Cirurgicamente uma Etapa
<cenário-5-refinar-cirurgicamente-uma-etapa>
#quote(block: true)[
#strong[Objetivo:] trocar uma etapa lenta ou com ferramenta antiga sem recriar o pipeline.
]

#strong[Dado que] a etapa `e03_analise_estatica` usa uma ferramenta obsoleta,

#strong[Quando] o usuário solicita a substituição atômica:

#emph[Via Terminal (CLI):]

```bash
legoloop refinar "auditoria-smart-contracts" \
  --trocar e03_analise_estatica --por "slither-analyzer@2.0.0" --yes
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-refinar
No processo 'auditoria-smart-contracts', refine a etapa 'e03_analise_estatica'
substituindo-a pela peça 'slither-analyzer@2.0.0'.
```

#emph[Via LegoLoop Studio:]

+ Clica com o botão direito sobre o nó `e03_analise_estatica`.
+ Seleciona #strong[Refinar / Trocar Peça] no menu de contexto.
+ Escolhe `slither-analyzer@2.0.0` no catálogo de peças aptas sugerido pela matriz.
+ Confirma no modal de validação de contratos.

#strong[Então] o LegoLoop:

+ Valida a compatibilidade dos contratos de entrada e saída com os nós adjacentes.
+ Substitui o nó no DAG preservando todas as conexões.
+ Incrementa a versão menor do processo, de `@1.0.0` para `@1.1.0`.

#strong[E] nenhuma outra etapa do fluxo sofre alteração ou efeito colateral.

== 13. Cenário 6: Remover com Segurança e Rollback
<cenário-6-remover-com-segurança-e-rollback>
#quote(block: true)[
#strong[Objetivo:] excluir um processo sem deixar resíduos nos editores nem quebrar links.
]

#strong[Dado que] o processo `experimento-legado` deve ser desativado,

#strong[Quando] o usuário aciona a remoção segura:

#emph[Via Terminal (CLI):]

```bash
legoloop remover "experimento-legado" --dry-run   # 1. simulação
legoloop remover "experimento-legado" --force     # 2. execução
```

#emph[Via Chat do Agente:]

```markdown
/legoloop-remover
Remova o processo 'experimento-legado' com segurança, executando a simulação
primeiro e limpando os links simbólicos de todos os harnesses.
```

#emph[Via LegoLoop Studio:]

+ No cartão de `experimento-legado`, abre o menu de contexto e seleciona #strong[Remover Processo].
+ O diálogo de #strong[Impacto e Dry-Run] lista arquivos e harnesses afetados.
+ O usuário digita o nome do processo e clica em #strong[Confirmar Remoção].

#strong[Então] o sistema exibe o relatório de desvinculação, limpa os symlinks nos harnesses, arquiva os dados e gera cópia de segurança para restauração imediata.

== 14. Cenário 7: Personalizar Micro e Nano-Etapas no `ROUTE.json`
<cenário-7-personalizar-micro-e-nano-etapas-no-route.json>
#quote(block: true)[
#strong[Objetivo:] controlar custo cognitivo e velocidade ajustando modelo de IA e limites por bloco.
]

#strong[Dado que] o arquivo `docs/routes/YYYY-MM-DD_ROUTE_EXEMPLO.json` define a execução por nós,

#strong[Quando] o usuário ajusta as propriedades declarativas de um nó:

#emph[Via Terminal ou editor de código:]

```json
{
  "id": "ticket-01",
  "titulo": "Gerar Casos de Teste Unitários",
  "modelo_cognitivo": "gemini-flash",
  "contexto_esteril": true,
  "mcps_permitidos": ["codebase-memory-mcp"],
  "stop_gates": ["SG-01-REVISAO-HUMANA"]
}
```

#emph[Via Chat do Agente:]

```markdown
No arquivo de rota mais recente, configure o nó 'ticket-01' para usar o modelo
'gemini-flash', contexto estéril e apenas a ferramenta 'codebase-memory-mcp'.
```

#emph[Via LegoLoop Studio:]

+ Clica no nó `ticket-01` para abrir o painel #strong[Propriedades da Etapa].
+ Seleciona `gemini-flash` no dropdown de modelo.
+ Ativa o switch #strong[Contexto Estéril].
+ Marca apenas o MCP `codebase-memory-mcp` e ativa o portão `SG-01-REVISAO-HUMANA`.
+ Clica em #strong[Aplicar Alterações].

#strong[Então], na execução pela Orchestrix:

+ O nó roda com o modelo econômico `gemini-flash`.
+ Carrega apenas `codebase-memory-mcp`, sem poluir a janela de contexto.
+ Pausa obrigatoriamente no portão humano `SG-01-REVISAO-HUMANA` antes de gravar em disco.

= PARTE IV --- Referência Rápida
<parte-iv-referência-rápida>
== 15. A Tríade de Triggers por Operação
<a-tríade-de-triggers-por-operação>
Na coluna Terminal, todo comando é prefixado por `legoloop`\; na coluna Chat, toda skill é prefixada por `/legoloop-`.

#figure(
  align(center)[#table(
    columns: (24.27%, 13.59%, 19.42%, 13.59%, 29.13%),
    align: (left,left,left,left,left,),
    table.header([Ação desejada], [Nível], [Terminal], [Chat], [LegoLoop Studio],),
    table.hline(),
    [#strong[Criar processo]], [Macro], [`fabricar "<nome>"`], [`fabricar`], [Botão #strong[\+ Novo Processo]],
    [#strong[Consultar inventário]], [Macro], [`listar [--json]`], [`listar`], [Aba #strong[Processos]],
    [#strong[Auditar reputação]], [Meso], [`matriz [--tier ouro]`], [`matriz`], [Aba #strong[Matriz]],
    [#strong[Encaixar sub-loop]], [Meso], [`compor <pai> <filho>`], [`compor`], [Arrastar bloco e traçar aresta],
    [#strong[Substituir etapa]], [Micro], [`refinar "<alvo>"`], [`refinar`], [Botão direito no nó, #strong[Refinar]],
    [#strong[Remover com rollback]], [Macro], [`remover "<nome>"`], [`remover`], [Menu do cartão, #strong[Remover]],
    [#strong[Ajustar modelo e contexto]], [Micro/Nano], [Edição do `ROUTE.json`], [Pedido livre], [Painel de propriedades do nó],
    [#strong[Compilar livro em PDF]], [Macro], [`legopress "<alvo>"`], [`legopress`], [Botão #strong[Exportar Livro Didático]],
  )]
  , kind: table
  )

== 16. Síntese
<síntese>
Com a Tríade de Triggers integrada ao padrão BDD e os encadeamentos canônicos mapeados, o LegoLoop deixa de ser apenas uma CLI:

- Em automação headless ou CI/CD, use o #strong[Terminal].
- Em programação pareada com um modelo de linguagem no editor, use o #strong[Chat do Agente].
- Para visão panorâmica e edição visual, use o #strong[LegoLoop Studio].

Em qualquer dos três caminhos, o comportamento contratado pelo BDD (#strong[Dado-Quando-Então]) é idêntico, determinístico e auditável.
