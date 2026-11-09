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
#let cor-accent   = rgb("#7c3aed") // Azul Cobalto Google Suave
#let cor-accent2  = cor-accent.darken(25%)                            // Tom profundo da cor da obra
#let cor-alerta   = rgb("#d97706")                                    // Âmbar Google Dourado
#let cor-ok       = rgb("#16a34a")                                    // Verde Esmeralda Google
#let cor-erro     = rgb("#dc2626")                                    // Carmesim Google

#set document(
  title: [Impeccable no LegoLoop: Design de Interfaces sem Clichês],
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
  )[ENGENHARIA VISUAL DE INTERFACES]

  #v(12pt)

  // Título Principal e Destaque
      #text(
      font: ("Segoe UI", "Arial"),
      size: 28pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.5pt,
    )[Impeccable]
    #v(4pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 15.5pt,
      weight: "bold",
      fill: cor-accent,
    )[Design de Interfaces sem Clichês]
  
  #v(8pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 9.6pt,
      fill: cor-suave,
      weight: "medium",
    )[O loop de cinco estágios, os seis sub-loops e o piso de qualidade que separa uma tela construída de uma tela montada]
  ]

  #v(16pt)

  // Arte: esteira de sub-loops desenhada em Typst
    #block(width: 100%, radius: 6pt, clip: true, fill: rgb("#f8fafc"), stroke: 0.5pt + rgb("#cbd5e1"), inset: 16pt)[
    #let etapas = ("Avaliar", "Construir", "Melhorar", "Refinar", "Corrigir", "Iterar e Documentar", )
    #text(font: ("Segoe UI", "Arial"), size: 7pt, weight: "bold", fill: cor-suave, tracking: 1.6pt)[OS SEIS NÓS DO LOOP impeccable-design]
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
    text(font: ("Segoe UI", "Arial"), size: 8.2pt, fill: cor-suave)[07/11/2026 · LegoLoop Foundation]
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
  Este livro ensina a operar o método Impeccable, criado por Paul Bakaus (upstream v4.7.0), dentro do ecossistema LegoLoop. O Impeccable é um conjunto de cerca de 25 comandos de design de interface organizados em seis sub-loops: Avaliar, Construir, Melhorar, Refinar, Corrigir e Iterar.

  A obra segue a progressão Macro, Meso e Micro. A Parte I explica por que interfaces geradas por IA convergem para o mesmo visual genérico e quais pilares o método usa para evitar isso. As Partes II e III descrevem o loop mestre de cinco estágios, seus quatro gates de parada e cada comando dos sub-loops. A Parte IV desce ao nível de implementação: o arquivo loops/impeccable-design.loop.json, o template do Studio Canvas e os agentes especializados. A Parte V traz receitas de uso e a Parte VI reúne glossário, referência rápida e um relatório de conformidade entre o manual e o código.

  Fonte primária: docs/manuals/2026-11-01\_MANUAL\_IMPECCABLE\_DESIGN\_LOOP.md, complementada pelo arquivo de referência upstream craft-floor.md.

  
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
          text(font: ("Segoe UI", "Arial"), size: 5.8pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[IMPECCABLE],
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
O livro atende três leitores diferentes, e cada um pode começar por um ponto distinto.

+ #strong[Quem decide sobre produto.] Precisa entender por que uma interface "bonita" gerada por IA ainda parece amadora e como um processo com gates reduz retrabalho. A Parte I responde a isso.
+ #strong[Quem opera o LegoLoop.] Precisa saber qual comando chamar, em que ordem e o que esperar como saída. As Partes II, III e V são o núcleo para esse perfil.
+ #strong[Quem mantém a plataforma.] Precisa saber onde o pipeline está declarado, como o Studio o instancia e onde o manual e o código divergem. A Parte IV e o Apêndice B são escritos para esse leitor.

== A estrutura da obra
<a-estrutura-da-obra>
#figure(
  align(center)[#table(
    columns: (19.23%, 19.23%, 61.54%),
    align: (left,left,left,),
    table.header([Parte], [Nível], [Conteúdo],),
    table.hline(),
    [I], [Macro], [O problema do visual genérico e os quatro pilares do método],
    [II], [Meso], [Os quatro modos de superfície e o loop mestre de cinco estágios],
    [III], [Meso], [Os seis sub-loops e cada um dos seus comandos],
    [IV], [Micro], [Piso de qualidade, agentes especializados e integração no LegoLoop],
    [V], [Prática], [Receitas para situações frequentes],
    [VI], [Apêndices], [Glossário, referência rápida e relatório de conformidade],
  )]
  , kind: table
  )

== As convenções da obra
<as-convenções-da-obra>
Comandos e caminhos de arquivo aparecem em fonte monoespaçada, como `polish` ou `loops/impeccable-design.loop.json`. Exemplos ilustrativos são marcados como tal no texto. Toda afirmação sobre o LegoLoop é rastreável a um arquivo citado; quando manual e código discordam, o livro segue o código e registra a diferença no Apêndice B.

= PARTE I --- Macro: Por que Interfaces de IA Parecem Iguais
<parte-i-macro-por-que-interfaces-de-ia-parecem-iguais>
== 1. O problema do visual médio
<o-problema-do-visual-médio>
Peça a um modelo de linguagem uma "landing page moderna" e o resultado costuma ser previsível: um título enorme sobre um gradiente roxo, um rótulo pequeno em caixa alta acima do título, três cartões do mesmo tamanho com ícone, título e parágrafo, um número grande com legenda pequena e sombras suaves em tudo. Peça de novo, para outro produto, e a página volta quase igual.

Isso não é falta de capacidade do modelo. É o efeito estatístico de como ele foi treinado. Um modelo aprende a produzir o que é mais provável dado o pedido. Para um pedido vago, o mais provável é a média de milhares de páginas parecidas. O resultado é uma interface que não tem erro óbvio, mas também não pertence a nenhum produto em particular.

O manual do Impeccable chama a saída desejada de #emph[out-of-distribution]: uma entrega que sai da média, porque foi decidida e não apenas gerada. O método inteiro existe para trocar o padrão estatístico por decisões explícitas.

=== Três sintomas que o leitor reconhece
<três-sintomas-que-o-leitor-reconhece>
- #strong[A tela montada.] Os componentes estão todos lá, mas nenhum foi pensado para o conteúdo real. Os rótulos dizem "Lorem" ou "Título", os estados de erro não existem e a tela quebra em 375 pixels de largura.
- #strong[O enfeite no lugar da hierarquia.] Gradiente no texto, brilho colorido sem deslocamento, vidro desfocado sem função. Cada um desses recursos tenta compensar a falta de uma hierarquia tipográfica clara.
- #strong[A revisão infinita.] O agente tira uma captura de tela, ajusta um detalhe, tira outra, ajusta outro. Sem um critério de parada, o ciclo consome tokens e tempo sem convergir.

#painel("A diferença entre construir e montar")[
  O arquivo de referência craft-floor.md, que acompanha o upstream do Impeccable, resume a distinção: as superfícies que ninguém desenhou (seleção de texto, cursor, barras de rolagem, anéis de foco, numerais em tabelas) continuam carregando o design. Tematizá-las a partir da paleta é, segundo o próprio documento, o sinal mais barato de que uma página foi construída e não montada, e o que os modelos mais deixam de fazer.
]
== 2. Os quatro pilares do método
<os-quatro-pilares-do-método>
O manual organiza o Impeccable em quatro pilares. Eles não são etapas; são regras que valem em qualquer etapa.

=== Pilar 1: o briefing vence
<pilar-1-o-briefing-vence>
A estética, a paleta, a época e as diretrizes definidas por quem pediu a interface têm prioridade sobre qualquer gosto do agente. O manual é explícito: nunca redirecione uma interface para gosto pessoal arbitrário.

Na prática, isso significa que se o briefing pede uma estética anos 90 com bordas grossas, o agente não "moderniza" a proposta por conta própria. O piso de qualidade, como veremos na Parte IV, trata vários recursos como padrões evitáveis e não como proibições: as palavras do briefing podem autorizar qualquer um deles.

=== Pilar 2: refinamento não é redesenho
<pilar-2-refinamento-não-é-redesenho>
Toda tarefa de interface começa com uma pergunta: estamos refinando ou redesenhando?

#figure(
  align(center)[#table(
    columns: (12.82%, 28.21%, 27.35%, 31.62%),
    align: (left,left,left,left,),
    table.header([Abordagem], [O que preserva], [O que muda], [Quando usar],),
    table.hline(),
    [Refinamento], [Identidade, comportamento e texto], [Alinhamento, ritmo, consistência], [A tela funciona, mas está mal acabada],
    [Redesenho], [Apenas o conteúdo e a função], [Todo o mundo visual], [A tela atual é anti-referência],
  )]
  , kind: table
  )

No redesenho, a interface antiga vira evidência do que não fazer. O novo mundo visual precisa ser documentado, e não apenas aplicado, para que a próxima sessão o reconheça.

=== Pilar 3: o piso de qualidade é inegociável
<pilar-3-o-piso-de-qualidade-é-inegociável>
Antes de qualquer edição de UI, inclusive ajustes pontuais, a leitura do arquivo `craft-floor.md` é obrigatória. O manual cita alguns exemplos do que esse piso proíbe: gradientes sem contraste, bordas excessivas, tipografia desordenada e emojis na interface. A Parte IV dedica um capítulo inteiro ao conteúdo desse arquivo.

=== Pilar 4: verificação delimitada
<pilar-4-verificação-delimitada>
O quarto pilar ataca diretamente o sintoma da revisão infinita. A regra tem três partes:

+ Inspecione em lote: desktop e mobile na mesma rodada.
+ Corrija todos os apontamentos em um único bloco de edição.
+ Finalize com, no máximo, mais uma checagem.

#frase-destaque[Uma rodada de captura, um bloco de correções, uma checagem final. Não existe quarta volta.]
O nome em inglês, #emph[bounded pass], descreve bem a ideia: a verificação tem uma fronteira conhecida antes de começar. Isso torna o custo previsível e força o agente a olhar para todos os problemas de uma vez, em vez de corrigi-los um por um à medida que aparecem.

== 3. O que o leitor leva desta parte
<o-que-o-leitor-leva-desta-parte>
O Impeccable não é uma biblioteca de componentes nem um tema visual. É um processo com quatro regras permanentes: respeitar o briefing, decidir entre refinar e redesenhar, cumprir o piso de qualidade e limitar a verificação. As próximas partes mostram como essas regras viram estágios, gates e comandos concretos.

= PARTE II --- Meso: Modos de Superfície e o Loop Mestre
<parte-ii-meso-modos-de-superfície-e-o-loop-mestre>
== 4. A ambição depende da tela, não do produto
<a-ambição-depende-da-tela-não-do-produto>
Um mesmo produto tem telas com objetivos muito diferentes. A página de preços precisa convencer; o painel de configurações precisa ser rápido de usar; a documentação precisa ser lida sem cansaço. O Impeccable mede a ambição visual pelo #strong[modo da superfície], e não pelo produto como um todo.

#figure(
  align(center)[#table(
    columns: (14.58%, 23.96%, 28.12%, 33.33%),
    align: (left,left,left,left,),
    table.header([Modo], [O visitante precisa], [Exemplos], [Prioridade de design],),
    table.hline(),
    [Operate], [Completar uma tarefa], [Dashboards, admin, editores], [Escaneabilidade e contenção],
    [Persuade], [Decidir e agir], [Landing pages, preços], [Hierarquia de valor e conversão],
    [Read], [Compreender um conteúdo], [Docs, manuais, changelogs], [Ritmo de leitura sem distração],
    [Experience], [Imergir na própria obra], [Portfólios, galerias], [A interface recua, a obra lidera],
  )]
  , kind: table
  )

=== Como o modo muda as decisões
<como-o-modo-muda-as-decisões>
No modo #strong[Operate], a ambição aparece como precisão: densidade coerente, alinhamentos exatos, comportamentos que o usuário já espera da plataforma. Um painel administrativo com tipografia expressiva e animações de entrada em cada seção está no modo errado.

No modo #strong[Persuade], a tipografia pode ser expressiva e a estética pode ser imersiva, porque o objetivo é mover o visitante para uma decisão. É o único modo, ao lado de Experience, em que o comando `overdrive` é permitido.

No modo #strong[Read], o contraste é descansado e a hierarquia tipográfica faz quase todo o trabalho. O piso de qualidade fixa a medida do corpo de texto entre 65 e 75 caracteres por linha.

No modo #strong[Experience], a interface sai do caminho. Desde a primeira dobra, quem lidera é a obra exibida: a fotografia, o projeto, a peça interativa.

#painel("Regra de classificação")[
  Classifique cada tela isoladamente. Um produto SaaS típico tem a landing em Persuade, o app em Operate e a central de ajuda em Read. Aplicar a mesma ambição às três é o erro mais comum de quem começa no método.
]
== 5. O loop mestre de cinco estágios
<o-loop-mestre-de-cinco-estágios>
O fluxo completo do Impeccable atravessa cinco estágios. Entre eles existem quatro gates binários de parada obrigatória, chamados P1 a P4. Binário quer dizer que o gate só tem duas respostas: passa ou não passa.

#esteira((
  ("1. Contexto e Auditoria", "PRODUCT.md, DESIGN.md, critique, audit"),
  ("2. Direção Visual", "classificar modo, shape, preservar ou substituir"),
  ("3. Execução Cirúrgica", "craft-floor.md, layout, typeset, colorize, clarify"),
  ("4. Verificação em Lote", "captura mobile 375px e desktop, polish, revisão"),
  ("5. Documentação e Entrega", "document, extract, proposta de PR"),
), gates: ("P1: achados e severidade", "P2: direção visual aprovada", "P3: antes e depois por item", "P4: veredito VERIFIED"))
=== Estágio 1: Contexto e Auditoria
<estágio-1-contexto-e-auditoria>
O agente carrega os dois documentos de contexto durável, `PRODUCT.md` e `DESIGN.md`, e então roda duas avaliações: `critique`, uma revisão heurística de experiência do usuário, e `audit`, uma checagem técnica de acessibilidade, performance e viewports.

#strong[Gate P1.] O estágio só termina quando os achados estão listados com severidade. Sem severidade, não há como decidir o que corrigir primeiro.

=== Estágio 2: Direção Visual
<estágio-2-direção-visual>
Com o diagnóstico em mãos, o agente classifica o modo da superfície, planeja a experiência antes do código com `shape` e decide o destino do trabalho novo: preservar, expandir ou substituir o que existe.

#strong[Gate P2.] A direção visual precisa de aprovação humana explícita. É o gate mais barato do loop: errar a direção aqui custa uma conversa; errar depois custa uma reimplementação.

=== Estágio 3: Execução Cirúrgica
<estágio-3-execução-cirúrgica>
A primeira ação do estágio é a leitura obrigatória de `craft-floor.md`. Depois vêm, nesta ordem, espaçamento e ritmo (`layout`), hierarquia tipográfica (`typeset`), paleta estratégica (`colorize`) e, por fim, textos e casos de borda (`harden` e `clarify`).

A ordem não é arbitrária. Cor aplicada antes de a hierarquia estar resolvida tende a ser usada como muleta para criar destaque que a tipografia deveria criar.

#strong[Gate P3.] Cada item alterado precisa de um antes e depois. Isso torna a revisão do gate objetiva: o revisor compara pares, não opina sobre o todo.

=== Estágio 4: Verificação em Lote
<estágio-4-verificação-em-lote>
Aqui o pilar da verificação delimitada vira prática. Uma captura em lote cobre mobile em 375 pixels e desktop. O comando `polish` faz o acabamento fino e o agente revisor de acabamento julga o resultado.

#strong[Gate P4.] O revisor emite um veredito. Só `VERIFIED` libera o estágio seguinte.

=== Estágio 5: Documentação e Entrega
<estágio-5-documentação-e-entrega>
O comando `document` atualiza o `DESIGN.md` a partir do código real, e não das intenções. O comando `extract` isola tokens e componentes para o design system. O estágio termina com uma proposta de entrega ou um pull request.

#frase-destaque[O DESIGN.md descreve o que foi entregue, não o que se pretendia entregar.]
== 6. Por que os gates ficam onde ficam
<por-que-os-gates-ficam-onde-ficam>
Cada gate fica no ponto em que o próximo estágio ficaria caro se começasse com uma premissa errada.

#figure(
  align(center)[#table(
    columns: (20.25%, 25.32%, 54.43%),
    align: (left,left,left,),
    table.header([Gate], [Fica antes de], [Custo evitado],),
    table.hline(),
    [P1], [Decidir a direção], [Direção escolhida sem diagnóstico],
    [P2], [Escrever código], [Reimplementar uma tela inteira],
    [P3], [Verificar em lote], [Revisão subjetiva sem comparação],
    [P4], [Documentar], [Registrar no DESIGN.md algo reprovado],
  )]
  , kind: table
  )

= PARTE III --- Meso: Os Seis Sub-Loops
<parte-iii-meso-os-seis-sub-loops>
== 7. Visão geral da taxonomia
<visão-geral-da-taxonomia>
Os cerca de 25 comandos do Impeccable estão agrupados em seis sub-loops reutilizáveis. Cada sub-loop tem um objetivo único, o que permite chamá-lo isoladamente quando a tarefa não exige o fluxo completo.

#figure(
  align(center)[#table(
    columns: (18.99%, 30.38%, 50.63%),
    align: (left,left,left,),
    table.header([Sub-loop], [Nome], [Objetivo em uma frase],),
    table.hline(),
    [1], [Avaliar], [Radiografar o estado atual sem alterar nada],
    [2], [Construir e Direcionar], [Fixar contratos de produto antes de codificar],
    [3], [Melhorar], [Elevar a qualidade pelos fundamentos visuais],
    [4], [Refinar], [Preparar a interface para produção],
    [5], [Corrigir], [Resolver fricções específicas],
    [6], [Iterar], [Explorar alternativas visuais ao vivo],
  )]
  , kind: table
  )

== 8. Sub-loop 1: Avaliar
<sub-loop-1-avaliar>
O objetivo é diagnosticar sem fazer alterações prematuras. Um agente que corrige enquanto avalia perde a visão do conjunto e tende a atacar o primeiro problema que encontra, não o mais grave.

- `critique` faz a revisão heurística: usabilidade, hierarquia de informação e carga cognitiva. Responde à pergunta "o usuário entende o que fazer?".
- `audit` faz a auditoria técnica: acessibilidade segundo a WCAG, performance de renderização e responsividade em vários viewports. Responde à pergunta "a tela funciona para todos?".

Os dois comandos são complementares. Uma tela pode passar no `audit` com contraste correto e foco visível e ainda assim falhar no `critique` por esconder a ação principal.

== 9. Sub-loop 2: Construir e Direcionar
<sub-loop-2-construir-e-direcionar>
Este sub-loop estabelece as fundações conceituais e os contratos duráveis de produto antes de qualquer código.

- `init` captura o contexto durável da aplicação em `PRODUCT.md`: para quem é, o que resolve, em que tom fala.
- `shape` planeja fluxos de tela e interações sem escrever código.
- `document` gera o `DESIGN.md` a partir da verdade do código.
- `extract` garimpa tokens e componentes consolidados para o design system.

#painel("Dois documentos, dois papéis")[
  PRODUCT.md responde "o quê e para quem". DESIGN.md responde "como está construído". O primeiro é escrito antes do código, a partir de decisões. O segundo é escrito depois, a partir do código entregue. Inverter essa ordem produz um DESIGN.md cheio de intenções que o código não cumpre.
]
== 10. Sub-loop 3: Melhorar
<sub-loop-3-melhorar>
O objetivo é elevar a qualidade sensorial e tátil da interface pelos fundamentos visuais. São seis comandos, e a ordem em que aparecem é a ordem recomendada.

#figure(
  align(center)[#table(
    columns: (17.07%, 43.9%, 39.02%),
    align: (left,left,left,),
    table.header([Comando], [Atua sobre], [Pergunta que responde],),
    table.hline(),
    [`layout`], [Respiros, grids, alinhamento, ritmo], [O olho encontra o caminho?],
    [`typeset`], [Proporção, peso, escala tipográfica], [A hierarquia é óbvia sem cor?],
    [`colorize`], [Paleta com significado semântico], [Cada cor quer dizer algo?],
    [`animate`], [Micro-interações funcionais], [O movimento explica uma mudança?],
    [`delight`], [Detalhes memoráveis e estados táteis], [Há um momento que se lembra?],
    [`overdrive`], [Expressividade estética ampliada], [Só em Persuade ou Experience],
  )]
  , kind: table
  )

O comando `colorize` evita dois extremos que o manual nomeia: o monocromatismo insosso e o carnaval de cores. A cor só entra quando carrega significado: estado, categoria, ação primária.

O comando `overdrive` merece atenção especial. Ele é reservado exclusivamente aos modos Persuade e Experience. Usá-lo em uma tela Operate é um erro de classificação, não uma escolha estética.

== 11. Sub-loop 4: Refinar
<sub-loop-4-refinar>
Este sub-loop blinda a interface para ambientes de missão crítica.

- `polish` é o passe de micropolimento e consistência antes da entrega.
- `harden` trata estados de erro, carregamento, internacionalização e limites de dados: o nome com 80 caracteres, a lista vazia, a lista com dez mil itens, a tradução alemã que dobra o tamanho do botão.
- `onboard` cria os fluxos de primeira execução e os estados vazios.
- `distill` remove complexidade visual desnecessária e elementos ruidosos.
- `bolder` e `quieter` calibram o tom visual: o primeiro aumenta contraste, o segundo reduz estimulação excessiva.

#painel("bolder e quieter não são opostos de gosto")[
  Os dois comandos corrigem um desvio em relação ao modo da superfície. Uma tela Persuade sem energia pede bolder. Uma tela Operate que grita pede quieter. A decisão vem do modo, não da preferência de quem revisa.
]
== 12. Sub-loop 5: Corrigir
<sub-loop-5-corrigir>
O objetivo é resolver fricções ergonômicas e falhas operacionais pontuais.

- `clarify` refatora a UX copy: microtexto, rótulos de botões e mensagens de validação. O piso de qualidade dá o critério: controles nomeiam sua ação; erros nomeiam o problema e a forma de recuperação.
- `adapt` adequa a tela a extremos: mobile compacto, ultra-wide e layouts nativos.
- `optimize` diagnostica e corrige gargalos de renderização na árvore de componentes.

Um botão que diz "Enviar" está correto, mas um botão que diz "Publicar artigo" está melhor, porque nomeia a consequência. Uma mensagem "Erro inválido" falha nos dois critérios: não diz qual é o problema nem como sair dele.

== 13. Sub-loop 6: Iterar
<sub-loop-6-iterar>
O último sub-loop serve para explorar alternativas visuais em tempo real.

- `live` permite inspecionar e selecionar elementos diretamente no navegador para gerar variantes ao vivo.
- `generate` cria variações simultâneas de um elemento no navegador para uma decisão imediata.

Iterar é útil quando a direção está aprovada, mas um elemento específico ainda não convence. Em vez de três rodadas de descrição verbal, o agente mostra três versões lado a lado.

== 14. Como escolher o sub-loop certo
<como-escolher-o-sub-loop-certo>
#figure(
  align(center)[#table(
    columns: (48.1%, 22.78%, 29.11%),
    align: (left,left,left,),
    table.header([Situação], [Sub-loop], [Comando de entrada],),
    table.hline(),
    [Não sei o que está errado], [Avaliar], [`critique`],
    [Tela nova, sem código ainda], [Construir], [`shape`],
    [Funciona, mas parece amadora], [Melhorar], [`layout`],
    [Vai para produção amanhã], [Refinar], [`harden`],
    [Usuários erram um formulário], [Corrigir], [`clarify`],
    [Duas opções e nenhuma decisão], [Iterar], [`generate`],
  )]
  , kind: table
  )

= PARTE IV --- Micro: Piso de Qualidade, Agentes e Integração
<parte-iv-micro-piso-de-qualidade-agentes-e-integração>
== 15. O piso de qualidade por dentro
<o-piso-de-qualidade-por-dentro>
O arquivo `craft-floor.md` fica em `upstream/reference/` dentro da skill `aidd-impeccable`. Ele é carregado depois que a direção está fixada e tem duas seções: o que verificar e o que recusar. O próprio arquivo avisa que um briefing fixado ou o mundo visual aprovado tem prioridade sobre qualquer regra dele; o hábito do agente, não.

=== O que verificar
<o-que-verificar>
Cada item é uma checagem sobre o resultado renderizado, e não uma intenção. Todos são conferidos na mesma captura.

#figure(
  align(center)[#table(
    columns: (29.11%, 70.89%),
    align: (left,left,),
    table.header([Eixo], [Critério verificável],),
    table.hline(),
    [Contraste], [Corpo e placeholder ≥ 4,5:1; texto grande ≥ 3:1],
    [Profundidade], [Sombras com deslocamento e desfoque suave],
    [Espaçamento], [Grupos apertados, separação generosa],
    [Tipografia], [Corpo entre 65 e 75ch; display até 6rem],
    [Tracking], [Piso de -0,04em; -0,02 a -0,03em lê melhor],
    [Movimento], [Um momento autoral, ease-out exponencial],
    [Estados], [Hover, disabled, loading, error e empty],
    [Superfícies], [Seleção, cursor, scrollbar e foco tematizados],
    [Texto], [Linguagem do produto; erros com saída],
    [Cobertura], [Todo requisito do briefing achável em segundos],
  )]
  , kind: table
  )

Sobre o contraste em superfícies coloridas, o arquivo acrescenta uma regra fina: o texto secundário deve ser tingido a partir do matiz da superfície ou do primeiro plano, nunca cinza.

=== O que recusar
<o-que-recusar>
O arquivo trata a maior parte desta lista como padrões da categoria, e não como proibições: o briefing pode autorizar qualquer um deles. Recorrer a um deles quando o eixo está livre, diz o texto, significa que o agente não estava decidindo.

#strong[Estruturas de página:]

- Cartões do mesmo tamanho com ícone, título e texto como estrutura da página. Cartões aninhados são sempre errados.
- O modelo de métrica de destaque: número grande, legenda pequena, estatísticas de apoio.
- Números de seção como 01, 02, 03, a menos que a sequência carregue informação.
- Um modal para uma tarefa que não precisa de interrupção nem de foco protegido.

#strong[Hábitos de superfície:]

- Texto com gradiente. A ênfase vem de peso ou tamanho.
- Vidro e desfoque como decoração.
- Borda lateral colorida acima de 1px em cartões, itens de lista ou alertas.
- Sombras duras deslocadas fora de um mundo de fato neobrutalista.
- Fonte monoespaçada como fantasia de "técnico", e não para código ou dados.
- Glifos Unicode ou emoji no lugar de um sistema de ícones.
- Claro ou escuro escolhido pela categoria, e não pela cena de uso.

#painel("A única proibição absoluta", cor: cor-erro)[
  Entre todos os itens, o arquivo marca um como banimento, e não como padrão: o rótulo pequeno (kicker ou eyebrow) acima de um título. Nenhum briefing o recupera. O título carrega o próprio peso; o rótulo deve ser apagado.
]
O arquivo também fixa limites numéricos: raio de cartões entre 12 e 16 pixels, com formato de pílula reservado a controles pequenos, e elevação declarada uma única vez, por borda ou por sombra. Uma borda de 1px sob uma sombra larga e suave é chamada de "cartão fantasma".

Por fim, uma regra de honestidade que vale para qualquer interface: afirmações e configurações vêm de dados fornecidos, e valores ilustrativos são rotulados como tal.

== 16. Os agentes especializados
<os-agentes-especializados>
Para evitar contaminação de contexto, o ecossistema do Impeccable delega quatro papéis a agentes autônomos. Cada um recebe apenas o que precisa para sua função.

#figure(
  align(center)[#table(
    columns: (27.85%, 51.9%, 20.25%),
    align: (left,left,left,),
    table.header([Agente], [Papel], [Momento],),
    table.hline(),
    [Documenter], [Gera DESIGN.md a partir do código entregue], [Estágio 5],
    [Finish reviewer], [Julga o acabamento contra a direção], [Gate P4],
    [Asset producer], [Produz ícones e vetores SVG inline], [Estágio 3],
    [Manual edit applier], [Aplica micro-ajustes da auditoria], [Estágios 3 e 4],
  )]
  , kind: table
  )

O manual descreve o revisor de acabamento como um diretor de arte intransigente: ele decide se o resultado alcançou o nível esperado ou precisa ser refeito. Os SVG do produtor de ativos seguem um padrão semântico, sem emoji e com `stroke="currentColor"`, para herdarem a cor do texto ao redor.

A separação em agentes tem uma razão prática. O agente que implementou uma tela tende a defendê-la. Um revisor que só recebe a direção aprovada e o resultado final julga sem esse viés.

== 17. Integração no LegoLoop
<integração-no-legoloop>
No LegoLoop, o pipeline existe em duas formas canônicas, ambas verificáveis no repositório.

=== A peça de loop no acervo
<a-peça-de-loop-no-acervo>
O arquivo `loops/impeccable-design.loop.json` segue o esquema `dag.v1` e se chama "Pipeline de Design Impeccable (6 Sub-Loops)", versão 1.0.0. Ele declara seis nós em cadeia linear, cada um dependente do anterior, com o artefato exigido e o modelo sugerido em `routing`:

#figure(
  align(center)[#table(
    columns: (33.33%, 48.28%, 18.39%),
    align: (left,left,left,),
    table.header([Nó], [Artefato que o gate exige], [Modelo sugerido],),
    table.hline(),
    [`subloop-avaliar`], [`docs/reports/AUDITORIA_DESIGN.md`], [sonnet, ou haiku],
    [`subloop-construir`], [`docs/specs/DIRECAO_VISUAL.md`], [sonnet, ou haiku],
    [`subloop-melhorar`], [`src/styles/theme.css`], [não declarado],
    [`subloop-refinar`], [`docs/reports/HARDENING_UI.md`], [não declarado],
    [`subloop-corrigir`], [`docs/reports/RESPONSIVE_VERIFICATION.md`], [não declarado],
    [`subloop-iterar-documentar`], [`DESIGN.md`], [não declarado],
  )]
  , kind: table
  )

Cada nó tem um gate mecânico do tipo `file_exists(...)`: o nó só é considerado concluído se o artefato existir no disco. É um gate de existência, não de qualidade; quem julga o conteúdo continua sendo o protocolo humano P1 a P4 conduzido pela skill.

Todos os seis nós apontam para a mesma peça, `skills/impeccable`. O que muda entre eles é o sub-loop que a peça executa.

=== O template do Studio Canvas
<o-template-do-studio-canvas>
Em `internal/studio/templates.go`, o template "Pipeline de Design Impeccable (6 Sub-Loops)" aparece no menu de templates rápidos com o ícone `palette`. Ao instanciá-lo, o canvas recebe a esteira completa com o roteamento cognitivo pré-configurado.

#painel("Do arquivo ao canvas")[
  O JSON do acervo e o template do Studio descrevem a mesma esteira por dois caminhos: o primeiro é executado pelo motor DAG via CLI; o segundo é montado visualmente para quem prefere arrastar nós. Alterar um não altera o outro automaticamente, e por isso vale conferir os dois quando o pipeline evoluir.
]
=== A skill no harness
<a-skill-no-harness>
A skill `aidd-impeccable` funciona como maestro: recebe um pedido em linguagem natural, como "auditar a tela dashboard" ou "polir antes de entregar", e escolhe comando, modo e ordem. Ela carrega a metodologia de uma cópia local do upstream, sem depender de plugin, e por isso roda em qualquer harness que leia arquivos.

== 18. Como conferir a integração
<como-conferir-a-integração>
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

= PARTE V --- Prática: Receitas de Uso
<parte-v-prática-receitas-de-uso>
As receitas a seguir combinam sub-loops para situações frequentes. São roteiros ilustrativos montados a partir do manual; os nomes de tela e os números servem de exemplo e não descrevem um projeto real.

== 18. Receita: auditar uma tela existente
<receita-auditar-uma-tela-existente>
#strong[Quando usar.] Alguém diz que "o dashboard está estranho", mas ninguém sabe dizer o quê.

+ Classifique o modo. Um dashboard é quase sempre Operate.
+ Rode `critique` e `audit` em sequência, sem corrigir nada.
+ Consolide os achados em uma lista com severidade alta, média e baixa.
+ Apresente a lista e pare. Este é o gate P1.

#strong[Saída esperada.] Uma lista como esta, ilustrativa:

#figure(
  align(center)[#table(
    columns: (21.52%, 56.96%, 21.52%),
    align: (left,left,left,),
    table.header([Severidade], [Achado], [Origem],),
    table.hline(),
    [Alta], [Ação principal abaixo da dobra em 375px], [`critique`],
    [Alta], [Contraste de 3,1:1 no texto secundário], [`audit`],
    [Média], [Quatro tamanhos de título sem escala clara], [`critique`],
    [Baixa], [Scrollbar com o estilo padrão do navegador], [`audit`],
  )]
  , kind: table
  )

#strong[Erro comum.] Corrigir o contraste enquanto ainda está avaliando. A correção isolada parece produtiva, mas impede de ver que o problema de hierarquia talvez peça uma reorganização maior.

== 19. Receita: polir antes de uma entrega
<receita-polir-antes-de-uma-entrega>
#strong[Quando usar.] A tela está aprovada e funcionando; falta acabamento.

+ Leia `craft-floor.md`.
+ Rode `harden` para cobrir estados de erro, vazio e carregamento.
+ Rode `clarify` sobre todos os rótulos e mensagens.
+ Faça uma captura em lote, mobile e desktop.
+ Rode `polish` e corrija tudo em um único bloco.
+ Faça no máximo uma checagem final.

#painel("O limite da receita")[
  Se a checagem final encontrar um problema estrutural, como a hierarquia que não funciona no mobile, a receita de polimento não é o lugar para resolvê-lo. Volte ao sub-loop Melhorar com layout e typeset, e só depois retome o polimento.
]
== 20. Receita: redesenhar uma landing page
<receita-redesenhar-uma-landing-page>
#strong[Quando usar.] A página atual é anti-referência: ninguém quer preservar nada dela além do conteúdo.

+ Rode `init` se não existir `PRODUCT.md`.
+ Use a página antiga apenas como evidência: liste o que ela faz mal.
+ Classifique o modo como Persuade.
+ Rode `shape` para planejar o fluxo sem código e peça aprovação da direção (P2).
+ Execute `layout`, `typeset` e `colorize`, nesta ordem.
+ Considere `overdrive`, que só é permitido aqui por causa do modo.
+ Verifique em lote (P4) e rode `document` para registrar o novo mundo visual.

== 21. Receita: decidir entre duas variantes
<receita-decidir-entre-duas-variantes>
#strong[Quando usar.] A direção está aprovada, mas o bloco de preços tem duas propostas e a equipe não decide.

+ Rode `generate` sobre o bloco para gerar as variações lado a lado no navegador.
+ Use `live` para ajustar um detalhe da variante preferida.
+ Registre a escolha e o motivo.
+ Volte ao estágio 4 com a variante escolhida.

== 22. Receita: corrigir um formulário em que usuários erram
<receita-corrigir-um-formulário-em-que-usuários-erram>
#strong[Quando usar.] O formulário de cadastro tem alta taxa de abandono e o suporte recebe perguntas sobre as mesmas mensagens de erro.

+ Rode `critique` focado no formulário para localizar onde o usuário hesita.
+ Rode `clarify` sobre rótulos, textos de ajuda e mensagens de validação.
+ Rode `harden` para cobrir os limites: nome muito longo, campo vazio, colagem de texto com espaços.
+ Faça a captura em lote, mobile e desktop, com os estados de erro visíveis.

O critério do `clarify` vem do piso de qualidade: cada erro precisa nomear o problema e a forma de recuperação. A tabela abaixo mostra o tipo de troca esperada; as mensagens são ilustrativas.

#figure(
  align(center)[#table(
    columns: (30.38%, 69.62%),
    align: (left,left,),
    table.header([Antes], [Depois],),
    table.hline(),
    [Erro inválido], [O CPF precisa ter 11 números. Confira os dígitos.],
    [Campo obrigatório], [Informe um e-mail para receber a confirmação.],
    [Senha fraca], [Use 8 caracteres ou mais, com pelo menos um número.],
    [Enviar], [Criar conta],
  )]
  , kind: table
  )

== 23. Receita: tornar uma tela mais calma
<receita-tornar-uma-tela-mais-calma>
#strong[Quando usar.] Usuários de um painel interno reclamam de cansaço visual.

+ Confirme o modo Operate.
+ Rode `distill` para remover elementos que não ajudam a tarefa.
+ Rode `quieter` para reduzir a estimulação.
+ Revise com `audit` para garantir que o contraste mínimo continua respeitado depois da redução.

O passo 4 existe porque reduzir estímulo visual costuma significar clarear cores, e cores claras demais derrubam o contraste abaixo de 4,5:1.

= PARTE VI --- Apêndices
<parte-vi-apêndices>
== Apêndice A: Glossário
<apêndice-a-glossário>
#strong[Bounded pass (verificação delimitada).] Rodada de revisão com fronteira fixada antes de começar: uma captura em lote, um bloco de correções, no máximo mais uma checagem.

#strong[Briefing.] Conjunto de diretrizes de quem pediu a interface: estética, paleta, época, tom. Tem prioridade sobre o gosto do agente.

#strong[Cartão fantasma.] Cartão com borda de 1px sob uma sombra larga e suave; declara elevação duas vezes.

#strong[craft-floor.md.] Arquivo de referência do upstream que define o piso de qualidade: o que verificar e o que recusar.

#strong[DESIGN.md.] Especificação de design gerada a partir do código entregue.

#strong[Gate binário.] Ponto de parada com apenas duas saídas: passa ou não passa.

#strong[Modo de superfície.] Classificação de uma tela por objetivo: Operate, Persuade, Read ou Experience.

#strong[Out-of-distribution.] Resultado que sai da média estatística do que um modelo produziria por padrão.

#strong[PRODUCT.md.] Contexto durável da aplicação: público, problema e tom.

#strong[Refinamento.] Melhoria que preserva identidade, comportamento e texto.

#strong[Redesenho.] Novo mundo visual em que a interface antiga serve apenas de evidência.

#strong[Sub-loop.] Grupo de comandos com um objetivo único, reutilizável de forma isolada.

== Apêndice B: Relatório de conformidade
<apêndice-b-relatório-de-conformidade>
Este apêndice compara o manual de 01/11/2026 com o estado do repositório em 07/11/2026. O objetivo é que o leitor saiba exatamente onde o texto do manual e o código coincidem.

#figure(
  align(center)[#table(
    columns: (22.45%, 31.63%, 30.61%, 15.31%),
    align: (left,left,left,left,),
    table.header([Ponto verificado], [Manual], [Código], [Situação],),
    table.hline(),
    [Arquivo do loop], [`impeccable-design.loop.json`], [Existe em `loops/`], [Conforme],
    [Número de nós], [6 sub-loops], [6 nós lineares], [Conforme],
    [Sexto nó], [Iterar], [`subloop-iterar-documentar`], [Diverge no nome],
    [Peça executora], [Não especificada], [`skills/impeccable` em todos], [Informativo],
    [Template do Studio], [Ícone `palette`], [Ícone `palette`], [Conforme],
    [Nome do template], [Não especificado], ["…(6 Sub-Loops)"], [Informativo],
    [Estágios e gates P1-P4], [5 estágios, 4 gates], [Gates `file_exists` por nó], [Diverge],
    [Nome do revisor], [`impeccabe-finish-reviewer`], [`impeccable-finish-reviewer`], [Erro de grafia],
  )]
  , kind: table
  )

#strong[Leitura do relatório.] O DAG declara apenas gates de existência de arquivo por nó; os gates de julgamento P1 a P4 são protocolo da skill, conduzido no chat. No sexto nó, o código agrupa iteração e documentação, enquanto o manual trata `document` como comando do sub-loop Construir. Para o revisor, o livro usa a grafia do agente instalado.

== Apêndice C: Referência rápida de comandos
<apêndice-c-referência-rápida-de-comandos>
#figure(
  align(center)[#table(
    columns: (20.25%, 20.25%, 59.49%),
    align: (left,left,left,),
    table.header([Comando], [Sub-loop], [Função resumida],),
    table.hline(),
    [`critique`], [Avaliar], [Revisão heurística de UX],
    [`audit`], [Avaliar], [Acessibilidade, performance, viewports],
    [`init`], [Construir], [Cria PRODUCT.md],
    [`shape`], [Construir], [Planeja fluxos sem código],
    [`document`], [Construir], [Gera DESIGN.md do código],
    [`extract`], [Construir], [Isola tokens e componentes],
    [`layout`], [Melhorar], [Espaçamento, grid, ritmo],
    [`typeset`], [Melhorar], [Escala e pesos tipográficos],
    [`colorize`], [Melhorar], [Paleta com significado],
    [`animate`], [Melhorar], [Micro-interações funcionais],
    [`delight`], [Melhorar], [Detalhes memoráveis],
    [`overdrive`], [Melhorar], [Expressividade; só Persuade e Experience],
    [`polish`], [Refinar], [Micropolimento final],
    [`harden`], [Refinar], [Erros, loading, i18n, limites],
    [`onboard`], [Refinar], [Primeira execução e estados vazios],
    [`distill`], [Refinar], [Remove ruído visual],
    [`bolder`], [Refinar], [Aumenta contraste e presença],
    [`quieter`], [Refinar], [Reduz estimulação],
    [`clarify`], [Corrigir], [UX copy e mensagens],
    [`adapt`], [Corrigir], [Telas extremas],
    [`optimize`], [Corrigir], [Gargalos de renderização],
    [`live`], [Iterar], [Variantes ao vivo no navegador],
    [`generate`], [Iterar], [Variações simultâneas],
  )]
  , kind: table
  )
