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
#let cor-accent   = rgb("#0f766e") // Azul Cobalto Google Suave
#let cor-accent2  = cor-accent.darken(25%)                            // Tom profundo da cor da obra
#let cor-alerta   = rgb("#d97706")                                    // Âmbar Google Dourado
#let cor-ok       = rgb("#16a34a")                                    // Verde Esmeralda Google
#let cor-erro     = rgb("#dc2626")                                    // Carmesim Google

#set document(
  title: [Matt Pocock no LegoLoop: Rigor de Engenharia com Agentes],
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
  )[ENGENHARIA DE SOFTWARE COM AGENTES]

  #v(12pt)

  // Título Principal e Destaque
      #text(
      font: ("Segoe UI", "Arial"),
      size: 28pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.5pt,
    )[Matt Pocock]
    #v(4pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 15.5pt,
      weight: "bold",
      fill: cor-accent,
    )[Rigor de Engenharia com Agentes]
  
  #v(8pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 9.6pt,
      fill: cor-suave,
      weight: "medium",
    )[Alinhar antes de construir, testar antes de implementar e corrigir o ambiente em vez do passado]
  ]

  #v(16pt)

  // Arte: esteira de sub-loops desenhada em Typst
    #block(width: 100%, radius: 6pt, clip: true, fill: rgb("#f8fafc"), stroke: 0.5pt + rgb("#cbd5e1"), inset: 16pt)[
    #let etapas = ("Alinhar Domínio", "Especificar", "Planejar Tickets", "Construir com TDD", "Revisar em Dois Eixos", "Retroalimentar o Ambiente", )
    #text(font: ("Segoe UI", "Arial"), size: 7pt, weight: "bold", fill: cor-suave, tracking: 1.6pt)[OS SEIS NÓS DO LOOP matt-pocock]
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
  Este livro ensina a operar o método de Matt Pocock (upstream sob licença MIT) dentro do ecossistema LegoLoop. O método substitui a geração cega de código por um ciclo disciplinado: entrevistar a ideia até não restarem premissas ocultas, transformar o consenso em especificação, fatiar a especificação em tickets com dependências explícitas e implementar cada ticket com testes primeiro, em contexto limpo.

  A Parte I apresenta os cinco pilares do método. A Parte II descreve o loop mestre e seus quatro gates. A Parte III detalha os seis sub-loops e as habilidades de cada um. A Parte IV desce ao nível de implementação: roteamento de modelos, higiene de contexto e o arquivo loops/matt-pocock.loop.json. A Parte V traz receitas e a Parte VI reúne glossário, relatório de conformidade e referência das habilidades.

  Fonte primária: docs/manuals/2026-11-01\_MANUAL\_MATT\_POCOCK\_SKILLS\_LOOP.md, complementada pelas habilidades upstream tdd, code-review e diagnosing-bugs.

  
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
          text(font: ("Segoe UI", "Arial"), size: 5.8pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[MATT POCOCK],
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
+ #strong[Quem lidera times que usam agentes.] Precisa entender por que "deixar a IA programar" produz código que ninguém consegue manter, e onde o humano precisa continuar decidindo. A Parte I responde a isso.
+ #strong[Quem desenvolve com agentes no dia a dia.] Precisa saber qual habilidade chamar, quando limpar o contexto e o que cada gate exige. As Partes II, III e V são o núcleo para esse perfil.
+ #strong[Quem mantém a plataforma LegoLoop.] Precisa saber onde o pipeline está declarado, quais artefatos cada nó exige e onde o manual e o código divergem. A Parte IV e o Apêndice B são escritos para esse leitor.

== A estrutura da obra
<a-estrutura-da-obra>
#figure(
  align(center)[#table(
    columns: (20.25%, 20.25%, 59.49%),
    align: (left,left,left,),
    table.header([Parte], [Nível], [Conteúdo],),
    table.hline(),
    [I], [Macro], [O problema da geração cega e os cinco pilares do método],
    [II], [Meso], [O loop mestre, da ideia à produção, com quatro gates],
    [III], [Meso], [Os seis sub-loops e suas habilidades],
    [IV], [Micro], [Roteamento de modelos, higiene de contexto e integração],
    [V], [Prática], [Receitas para situações frequentes],
    [VI], [Apêndices], [Glossário, relatório de conformidade e referência rápida],
  )]
  , kind: table
  )

== As convenções da obra
<as-convenções-da-obra>
Nomes de habilidades e caminhos de arquivo aparecem em fonte monoespaçada, como `grill-with-docs` ou `loops/matt-pocock.loop.json`. Exemplos ilustrativos são marcados como tal no texto. Toda afirmação sobre o LegoLoop é rastreável a um arquivo citado; quando manual e código discordam, o livro segue o código e registra a diferença no Apêndice B.

= PARTE I --- Macro: O Fim da Geração Cega
<parte-i-macro-o-fim-da-geração-cega>
== 1. O problema: código que chega antes da decisão
<o-problema-código-que-chega-antes-da-decisão>
Um agente de IA consegue escrever centenas de linhas em segundos. Esse é o problema. Quando o pedido é "adicione um sistema de assinaturas", o agente não para para perguntar se uma assinatura pertence a um usuário ou a uma conta, se o cancelamento vale no fim do ciclo ou imediatamente, ou o que acontece com um pagamento recusado. Ele escolhe uma resposta plausível para cada pergunta, sem avisar, e escreve o código.

O resultado parece pronto. Os testes passam, porque foram escritos pelo mesmo agente, a partir das mesmas suposições. Semanas depois, alguém descobre que "conta" significa coisas diferentes em três módulos e que a regra de cancelamento contradiz o contrato comercial.

O método de Matt Pocock parte de um diagnóstico simples: o gargalo do desenvolvimento com agentes não é escrever código, é decidir o que escrever. Por isso o método gasta esforço deliberado antes da primeira linha de produção e mantém o humano no controle das decisões que são difíceis de desfazer.

=== Três sintomas da geração cega
<três-sintomas-da-geração-cega>
- #strong[Premissas silenciosas.] O agente preencheu lacunas do pedido com suposições que ninguém aprovou.
- #strong[Testes que concordam com o código por construção.] O teste recalcula o resultado do mesmo jeito que a implementação e, por isso, nunca discorda dela.
- #strong[A sessão que apodrece.] Depois de horas de conversa, o agente começa a esquecer decisões do início, a misturar tentativas descartadas com a versão atual e a errar coisas que acertava antes.

== 2. Os cinco pilares do método
<os-cinco-pilares-do-método>
=== Pilar 1: alinhar antes de construir
<pilar-1-alinhar-antes-de-construir>
Em inglês, #emph[grill before build]. Nenhuma linha de código de produção é gerada antes de a ideia passar por uma entrevista implacável que exponha premissas ocultas, casos de borda e ambiguidades de domínio. A entrevista não é uma formalidade: ela termina com decisões registradas em arquivo.

=== Pilar 2: o desenvolvedor no controle
<pilar-2-o-desenvolvedor-no-controle>
O agente opera em execuções curtas, de três a cinco passos, com gates de parada obrigatória. Ações críticas, como push, abertura de pull request, deploy e deleção, exigem autorização humana expressa.

#painel("Por que execuções curtas")[
  Um agente que executa vinte passos sem parar acumula erros que se apoiam uns nos outros. Em três a cinco passos, o desvio é pequeno o bastante para ser notado e barato o bastante para ser desfeito.
]
=== Pilar 3: módulos profundos
<pilar-3-módulos-profundos>
O pilar vem de John Ousterhout e do livro #emph[A Philosophy of Software Design]. Um módulo profundo tem uma interface estreita e simples que esconde uma implementação rica. O contrário, o módulo raso, expõe quase tanta complexidade na interface quanto contém por dentro.

O método usa a palavra #emph[seam], costura, para a fronteira pública em que o comportamento é observado sem olhar para dentro. É nas costuras que os testes moram.

=== Pilar 4: higiene severa de contexto
<pilar-4-higiene-severa-de-contexto>
Modelos de linguagem degradam o raciocínio quando a janela de contexto cresce além de uma zona de clareza que o manual estima em cerca de 150 mil tokens. O método chama a faixa saudável de #emph[Smart Zone] e organiza o trabalho para ficar dentro dela.

A regra prática é dupla: alinhamento, especificação e tickets acontecem em uma janela única e contínua, porque compartilham raciocínio; a implementação de cada ticket começa em contexto limpo, porque só precisa do ticket e do repositório.

=== Pilar 5: retroalimentar o ambiente
<pilar-5-retroalimentar-o-ambiente>
Retrospectivas não mudam o código passado. Elas mudam as regras, os linters, as checagens automáticas e o `AGENTS.md`, para que a próxima sessão não repita o mesmo erro.

#frase-destaque[Um erro corrigido no código volta. Um erro corrigido no ambiente não volta.]
== 3. O que o leitor leva desta parte
<o-que-o-leitor-leva-desta-parte>
O método não é um conjunto de prompts melhores. É uma divisão de trabalho: o humano decide o que é difícil de reverter, o agente executa fatias pequenas e verificáveis, e cada falha vira uma regra automática. A tabela abaixo antecipa onde cada pilar se materializa no loop.

#figure(
  align(center)[#table(
    columns: (41.77%, 58.23%),
    align: (left,left,),
    table.header([Pilar], [Onde aparece no loop],),
    table.hline(),
    [Alinhar antes de construir], [Sub-loop 1 e gate G1],
    [Desenvolvedor no controle], [Gates G1 a G4 e ações críticas autorizadas],
    [Módulos profundos], [`codebase-design` e testes em costuras],
    [Higiene de contexto], [Janela única até os tickets, `/clear` depois],
    [Retroalimentar o ambiente], [`retro` e o gate G4],
  )]
  , kind: table
  )

= PARTE II --- Meso: O Loop Mestre de Rigor
<parte-ii-meso-o-loop-mestre-de-rigor>
== 4. Da ideia à produção em cinco estágios
<da-ideia-à-produção-em-cinco-estágios>
O loop mestre atravessa cinco estágios. Quatro gates binários separam os estágios principais: G1 a G4. O segundo estágio, a prova de conceito, é opcional e funciona como um desvio que volta ao primeiro.

#esteira((
  ("1. Alinhamento e Domínio", "grill-with-docs, domain-modeling"),
  ("3. Especificação e Grafo", "to-spec, to-tickets"),
  ("4. Construção TDD", "implement, tdd, codebase-design"),
  ("5. Fechamento e Entrega", "code-review, pr, retro"),
), gates: ("G1: decisões aprovadas", "G2: tickets e bloqueios aprovados", "G3: testes 100% verdes"))
Depois do estágio 5, o gate G4, review aprovado, marca a entrega como concluída.

=== O desvio empírico: estágio 2
<o-desvio-empírico-estágio-2>
Às vezes uma pergunta do alinhamento não se resolve conversando. "Esse modelo de estado aguenta edição concorrente?" só se responde com código rodando. Nesse caso o fluxo sai do estágio 1, abre uma sessão isolada com `handoff`, constrói um protótipo descartável com `prototype` e volta com a resposta por outro `handoff`. O protótipo não entra na árvore principal.

#painel("Protótipo não é rascunho de produção")[
  O protótipo responde uma pergunta e é jogado fora. Se o código do protótipo for aproveitado na implementação, ele entra sem spec, sem tickets e sem testes primeiro, que é exatamente o que o método existe para evitar.
]
== 5. Os gates e o que cada um exige
<os-gates-e-o-que-cada-um-exige>
#figure(
  align(center)[#table(
    columns: (17.5%, 33.75%, 48.75%),
    align: (left,left,left,),
    table.header([Gate], [Entre], [Condição para passar],),
    table.hline(),
    [G1], [Alinhamento e Especificação], [Ambiguidades de domínio sanadas e aprovadas],
    [G2], [Especificação e Construção], [DAG de tickets aprovado, bloqueios resolvidos],
    [G3], [Construção e Fechamento], [Testes 100% verdes, sem testes enfraquecidos],
    [G4], [Fechamento e Entrega], [Review sem notas bloqueadoras, ambiente atualizado],
  )]
  , kind: table
  )

=== G1: decisões aprovadas
<g1-decisões-aprovadas>
O estágio de alinhamento termina quando o desenvolvedor aprova as decisões. Não basta o agente achar que entendeu; as decisões precisam estar escritas no `GLOSSARY.md`, em um documento de intenção ou em ADRs.

=== G2: tickets e bloqueios aprovados
<g2-tickets-e-bloqueios-aprovados>
A especificação vira um grafo acíclico de tickets. O gate exige que cada aresta de bloqueio esteja explícita: o ticket 4 depende do 2, o 5 não depende de ninguém. Um ticket com dependência implícita quebra a execução em contexto limpo, porque a sessão que o implementa não sabe o que veio antes.

=== G3: testes 100% verdes
<g3-testes-100-verdes>
O manual detalha o que conta como verde: nenhum teste enfraquecido e nenhuma supressão de compilador. Um teste que foi afrouxado para passar, ou um `@ts-ignore` colocado para calar um erro de tipo, reprovam o gate mesmo com a suíte verde.

=== G4: review aprovado
<g4-review-aprovado>
O review em dois eixos não pode ter notas bloqueadoras, e a retrospectiva precisa ter atualizado o ambiente. Uma entrega sem retro passa no código, mas reprova no método.

== 6. Por que a ordem importa
<por-que-a-ordem-importa>
Cada estágio produz o insumo exato do seguinte. A especificação depende de um vocabulário fixado; os tickets dependem da especificação; a implementação depende de tickets autocontidos; o review depende da especificação para o segundo eixo.

#figure(
  align(center)[#table(
    columns: (26.58%, 46.84%, 26.58%),
    align: (left,left,left,),
    table.header([Estágio], [Produz], [Consumido por],),
    table.hline(),
    [Alinhamento], [`GLOSSARY.md`, intenção, ADRs], [Especificação],
    [Especificação], [Spec e DAG de tickets], [Construção],
    [Construção], [Código, testes, relatório TDD], [Fechamento],
    [Fechamento], [Review, corpo do PR, `AGENTS.md`], [A próxima sessão],
  )]
  , kind: table
  )

= PARTE III --- Meso: Os Seis Sub-Loops
<parte-iii-meso-os-seis-sub-loops>
== 7. Visão geral
<visão-geral>
O manual decompõe o ecossistema de Matt Pocock em seis sub-loops, disponíveis como blocos conectáveis no Studio Canvas e em pipelines DAG.

#figure(
  align(center)[#table(
    columns: (16.67%, 35.71%, 47.62%),
    align: (left,left,left,),
    table.header([Sub-loop], [Nome], [Gate de saída],),
    table.hline(),
    [1], [Concepção, Grilling e Domínio], [Ambiguidades sanadas e aprovadas],
    [2], [Exploração Empírica e Pesquisa], [Resposta empírica devolvida via `handoff`],
    [3], [Especificação e Grafo], [DAG aprovado, bloqueios resolvidos],
    [4], [Construção Cirúrgica e TDD], [Testes verdes, sem testes enfraquecidos],
    [5], [Diagnóstico e Arquitetura], [Falha reproduzida, corrigida e protegida],
    [6], [Auditoria, Entrega e Ambiente], [Review aprovado e ambiente atualizado],
  )]
  , kind: table
  )

== 8. Sub-loop 1: Concepção, Grilling e Domínio
<sub-loop-1-concepção-grilling-e-domínio>
O objetivo é fixar um vocabulário unívoco, desafiar premissas e registrar decisões arquiteturais antes de qualquer código.

- `grill-with-docs` é a entrevista socrática focada no repositório. Deixa rastro durável em `GLOSSARY.md` e em ADRs.
- `grill-me` é a variante sem estado, para alinhamento rápido sem pasta de projeto.
- `grilling` é o primitivo da entrevista: rodadas numeradas e separação clara entre o que o agente sabe e o que o usuário decide.
- `domain-modeling` audita termos sobrecarregados, como "usuário" e "conta", e documenta decisões difíceis de reverter.

#strong[Saídas:] `GLOSSARY.md`, `docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md` e ADRs em `docs/adr/`.

#painel("Fatos do agente e decisões do usuário")[
  A separação que o grilling impõe é a parte mais valiosa da entrevista. O agente pode afirmar fatos que leu no repositório: o campo tal existe, a tabela tal tem tal índice. Ele não pode afirmar decisões: se o cancelamento é imediato é uma escolha de produto, e só o usuário a faz.
]
== 9. Sub-loop 2: Exploração Empírica e Pesquisa
<sub-loop-2-exploração-empírica-e-pesquisa>
O objetivo é responder perguntas que só código executável ou fontes primárias resolvem, sem poluir a árvore principal.

- `prototype` cria um programa descartável na branch `prototype/<nome>`, para testar um modelo de estado ou a sensação de uma interface.
- `handoff` passa o contexto de forma estéril para uma nova pasta, harness ou sessão.
- `research` roda uma pesquisa profunda em segundo plano sobre documentação oficial e fontes primárias.

#strong[Saídas:] relatório técnico com citações e protótipo referenciado na branch isolada.

== 10. Sub-loop 3: Especificação e Grafo
<sub-loop-3-especificação-e-grafo>
O objetivo é sintetizar o consenso em uma especificação formal e fatiá-la em um grafo acíclico de tarefas independentes.

- `to-spec` consolida a discussão em uma spec executável e determinística, sem entrevistar de novo.
- `to-tickets` quebra a spec em tickets verticais, os #emph[tracer bullets], cada um com suas arestas de bloqueio (`blocked-by`).
- `wayfinder` traça um mapa de tickets de decisão para iniciativas nebulosas ou projetos grandes começando do zero.

#strong[Saídas:] `docs/specs/YYYY-MM-DD_SPEC_<TEMA>.md` e tickets em `.scratch/<feature>/issues/NN-slug.md` ou `docs/plans/`.

Um ticket vertical atravessa todas as camadas necessárias para entregar um comportamento observável, da interface ao banco. O contrário, o ticket horizontal ("criar todas as tabelas"), não entrega nada testável sozinho.

== 11. Sub-loop 4: Construção Cirúrgica e TDD
<sub-loop-4-construção-cirúrgica-e-tdd>
O objetivo é implementar orientado a testes, com contexto isolado por ticket.

- `implement` consome um ticket por vez, a partir de uma janela limpa, atacando nós desbloqueados.
- `implement-spec` distribui subagentes em paralelo sobre os tickets prontos e unifica em uma branch de integração.
- `tdd` conduz o ciclo vermelho-verde por fatia de comportamento, sem mocks de produção e sem stubs vazios.
- `codebase-design` fornece o vocabulário de módulos profundos, interfaces estreitas e costuras.

#strong[Saídas:] código de produção, testes e `docs/reports/TDD_REPORT.md`.

=== O que a habilidade tdd exige
<o-que-a-habilidade-tdd-exige>
O texto upstream da habilidade `tdd` é mais preciso do que o resumo do manual em três pontos.

+ #strong[Testes só em costuras combinadas.] Antes de escrever qualquer teste, o agente lista as costuras que serão testadas e confirma com o usuário. Nenhum teste é escrito em costura não confirmada.
+ #strong[Três anti-padrões nomeados.] O teste acoplado à implementação, que quebra em refatorações sem mudança de comportamento; o teste tautológico, que recalcula o esperado do mesmo jeito que o código; e o fatiamento horizontal, que escreve todos os testes antes de toda a implementação.
+ #strong[Refatorar não faz parte do loop.] O ciclo é vermelho e verde. A refatoração pertence ao estágio de review, conduzido pela habilidade `code-review`.

#painel("O teste tautológico", cor: cor-erro)[
  expect(soma(a, b)).toBe(a + b) passa por construção: o esperado é calculado do mesmo jeito que o código calcula. O valor esperado precisa vir de uma fonte independente: um literal sabidamente correto, um exemplo resolvido à mão ou a especificação.
]
== 12. Sub-loop 5: Diagnóstico, Arquitetura e Triagem
<sub-loop-5-diagnóstico-arquitetura-e-triagem>
O objetivo é resolver falhas complexas e manter a saúde arquitetural.

- `diagnosing-bugs` proíbe teorizar antes de obter um loop de feedback reproduzível: um comando que fique vermelho por causa deste bug. A correção vem acompanhada de teste de regressão.
- `improve-codebase-architecture` procura módulos rasos e oportunidades de aprofundamento.
- `triage` qualifica e enriquece issues externas antes de despachá-las para agentes.

O texto upstream de `diagnosing-bugs` chama a construção do loop de feedback de "a habilidade em si" e lista formas de obtê-lo, em ordem: teste falhando, script HTTP contra o servidor local, chamada de CLI comparando a saída com um snapshot, script de navegador headless e reprodução de um trace capturado. A mesma habilidade exige ocultar todo segredo como `<REDACTED>` antes de mostrar comandos e saídas.

== 13. Sub-loop 6: Auditoria, Entrega e Ambiente
<sub-loop-6-auditoria-entrega-e-ambiente>
O objetivo é revisar a conformidade, redigir o PR e blindar o ambiente.

- `code-review` audita em dois eixos ortogonais: #emph[Standards], as regras documentadas do repositório, e #emph[Spec], a fidelidade à tarefa de origem.
- `pr` estrutura o pull request com a menor evidência visual necessária, antes e depois, e classifica o risco como porta de mão única ou de mão dupla.
- `retro` converte falhas mecânicas em checagens automáticas e divergências em regras de linter ou do `AGENTS.md`.
- `wizard` gera scripts interativos para etapas que exigem presença humana, como tokens e senhas.
- `wait-what` re-explica um conceito sem jargão quando o usuário se perde.

#strong[Saídas:] `docs/reports/CODE_REVIEW.md`, corpo do PR e atualizações no `AGENTS.md`.

O upstream de `code-review` acrescenta um detalhe de mecânica: os dois eixos rodam em subagentes paralelos, para que um não contamine o contexto do outro, sobre o diff de três pontos entre um ponto fixo informado pelo usuário e o `HEAD`.

= PARTE IV --- Micro: Roteamento, Contexto e Integração
<parte-iv-micro-roteamento-contexto-e-integração>
== 14. Roteamento cognitivo de modelos
<roteamento-cognitivo-de-modelos>
O manual recomenda um modelo primário e um alternativo local para cada sub-loop. A lógica é gastar raciocínio caro onde há ambiguidade e geração rápida onde há execução mecânica.

#figure(
  align(center)[#table(
    columns: (24.36%, 25.64%, 20.51%, 29.49%),
    align: (left,left,left,left,),
    table.header([Sub-loop], [Complexidade], [Primário], [Alternativo local],),
    table.hline(),
    [1. Concepção], [Alta (socrática)], [`sonnet`], [`haiku`],
    [2. Exploração], [Média], [`code-fast`], [`qwen2.5-coder:7b`],
    [3. Especificação], [Alta (arquitetura)], [`code-pro`], [`llama3.1:8b`],
    [4. Construção TDD], [Média a alta], [`sonnet`], [`qwen2.5-coder:7b`],
    [5. Diagnóstico], [Alta], [`sonnet`], [`llama3.1:8b`],
    [6. Review e Retro], [Média], [`code-fast`], [`haiku`],
  )]
  , kind: table
  )

O raciocínio socrático do sub-loop 1 pede um modelo capaz de desafiar premissas. Protótipos descartáveis do sub-loop 2 pedem velocidade, não profundidade. O review do sub-loop 6 é uma verificação binária de regras, em que um modelo rápido basta.

== 15. Higiene de contexto e fronteiras de fase
<higiene-de-contexto-e-fronteiras-de-fase>
=== As duas regras fixas
<as-duas-regras-fixas>
+ #strong[Janela unificada de alinhamento.] Grilling, `to-spec` e `to-tickets` ficam na mesma sessão, porque compartilham raciocínio.
+ #strong[Contexto limpo por ticket.] Cada `implement` começa em sessão limpa, consumindo apenas o arquivo do ticket e o repositório.

=== A decisão na fronteira de fase
<a-decisão-na-fronteira-de-fase>
#figure(
  align(center)[#table(
    columns: (73.42%, 26.58%),
    align: (left,left,),
    table.header([Situação na fronteira], [Ação],),
    table.hline(),
    [Fase curta, contexto abaixo de 100 mil tokens], [Continuar],
    [Fase terminou, próximo ticket autocontido], [`/clear`],
    [Alinhamento perto do limite, antes dos tickets], [`/compact`],
    [Bifurcação lateral ou protótipo], [`/handoff`],
  )]
  , kind: table
  )

#painel("Compactar não é limpar")[
  /compact resume a conversa e mantém o resumo; serve quando o raciocínio ainda é necessário, como no fim do alinhamento. /clear descarta tudo; serve quando o próximo passo é autocontido, como um ticket. Usar /compact antes de cada ticket carrega para a implementação resíduos de decisões que já estão no ticket.
]
== 16. Integração no LegoLoop
<integração-no-legoloop>
=== O arquivo do loop
<o-arquivo-do-loop>
O arquivo `loops/matt-pocock.loop.json` segue o esquema `dag.v1`, chama-se "Pipeline de Rigor Matt Pocock (6 Sub-Loops)", versão 1.0.0, e tem gatilho manual. Os seis nós apontam para a peça `skills/aidd-pocock` e formam uma cadeia linear. Cada nó declara o artefato que precisa existir para ser concluído.

#figure(
  align(center)[#table(
    columns: (35.71%, 35.71%, 28.57%),
    align: (left,left,left,),
    table.header([Nó], [Artefato exigido pelo gate], [Modelo sugerido],),
    table.hline(),
    [`subloop-1-alinhar-dominio`], [`GLOSSARY.md`], [sonnet, ou haiku],
    [`subloop-2-especificar`], [`docs/specs/SPEC.md`], [não declarado],
    [`subloop-3-planejar-tickets`], [`docs/plans/PLAN.md`], [não declarado],
    [`subloop-4-construir-tdd`], [`docs/reports/TDD_REPORT.md`], [sonnet, ou qwen2.5-coder],
    [`subloop-5-review-auditar`], [`docs/reports/CODE_REVIEW.md`], [não declarado],
    [`subloop-6-retro-ambiente`], [`docs/reports/RETRO.md`], [não declarado],
  )]
  , kind: table
  )

O gate de cada nó é do tipo `file_exists(...)`. Ele confirma que o artefato foi produzido, não que está correto. O julgamento continua com os gates humanos G1 a G4.

=== O template do Studio
<o-template-do-studio>
Em `internal/studio/templates.go`, o template "Pipeline Matt Pocock (6 Sub-Loops)" aparece no menu Templates do Studio Canvas com o ícone `git`. Ele carrega o grafo com nós de alinhamento, especificação, tickets, TDD, review e retroalimentação.

=== Comandos de verificação
<comandos-de-verificação>
```bash
legoloop validate loops/matt-pocock.loop.json   # contrato, gates e ciclos
legoloop run loops/matt-pocock.loop.json        # execução pelo motor DAG
legoloop studio                                 # Studio em 127.0.0.1:7878
```

Em 07/11/2026, o `validate` respondeu `6 nós válidos, sem ciclos`, em seis camadas de um nó cada, de `subloop-1-alinhar-dominio` a `subloop-6-retro-ambiente`.

= PARTE V --- Prática: Receitas de Uso
<parte-v-prática-receitas-de-uso>
As receitas combinam sub-loops para situações frequentes. Os nomes de funcionalidades e os números são ilustrativos.

== 17. Receita: uma funcionalidade nova, do zero
<receita-uma-funcionalidade-nova-do-zero>
#strong[Quando usar.] Há uma ideia, por exemplo "permitir que clientes pausem a assinatura", e nenhuma linha de código.

+ Rode `grill-with-docs`. Responda as rodadas até que "pausar" tenha um significado único: quanto tempo, quem paga, o que acontece com o ciclo.
+ Confirme que as decisões estão em `GLOSSARY.md` e, se alguma for difícil de reverter, em uma ADR. Este é o G1.
+ Na mesma sessão, rode `to-spec` e depois `to-tickets`.
+ Revise o DAG: cada ticket precisa dizer de quem depende. Este é o G2.
+ Limpe o contexto. Para cada ticket desbloqueado, abra uma sessão limpa e rode `implement` com `tdd`.
+ Com a suíte verde e sem testes enfraquecidos (G3), rode `code-review`, `pr` e `retro`.

#strong[Erro comum.] Pular o passo 5 e implementar os tickets na mesma sessão do alinhamento. A sessão já carrega toda a conversa, e o agente passa a implementar a partir dela, e não do ticket.

== 18. Receita: um bug difícil
<receita-um-bug-difícil>
#strong[Quando usar.] "O total do carrinho às vezes sai errado." Ninguém sabe reproduzir.

+ Rode `diagnosing-bugs` e não aceite hipóteses antes do passo 2.
+ Construa um loop de feedback: um comando que fique vermelho por causa deste bug. Tente, nesta ordem, teste, script HTTP, CLI com snapshot, navegador headless e reprodução de trace.
+ Com o vermelho em mãos, investigue.
+ Corrija com o menor patch possível e veja o mesmo comando ficar verde.
+ Transforme o comando em teste de regressão permanente.
+ Rode `retro`: se o bug veio de uma classe de erro detectável, converta-a em checagem automática.

#painel("Sem loop, sem diagnóstico")[
  Se nenhuma das formas de construir o loop funcionar, o passo seguinte é dizer isso ao usuário e pedir mais informação, não teorizar. Teorizar sem um sinal vermelho produz correções que parecem plausíveis e não se pode provar que funcionaram.
]
== 19. Receita: uma dúvida que só código responde
<receita-uma-dúvida-que-só-código-responde>
#strong[Quando usar.] No meio do alinhamento surge a pergunta "uma fila em memória aguenta o pico de pedidos?".

+ Interrompa o grilling e rode `handoff` para uma sessão isolada.
+ Na sessão nova, rode `prototype` na branch `prototype/fila-pico`.
+ Meça, anote o resultado e descarte o código.
+ Volte com `handoff`, levando só a resposta e a evidência.
+ Retome o grilling com a decisão informada.

== 20. Receita: revisar o trabalho de outro agente
<receita-revisar-o-trabalho-de-outro-agente>
#strong[Quando usar.] Um agente entregou uma branch dizendo que está pronta.

+ Fixe o ponto de comparação, por exemplo `main`, e confirme que ele existe e que o diff não está vazio.
+ Localize a spec de origem: referência de issue nos commits, caminho passado como argumento ou arquivo em `docs/specs/`.
+ Rode `code-review`. Os dois eixos, Standards e Spec, rodam em subagentes separados.
+ Trate cada nota bloqueadora antes do PR. Notas não bloqueadoras viram tickets.
+ Rode `retro` sobre o que o review encontrou que uma checagem automática poderia ter pego.

O passo 1 existe porque um ponto de comparação inválido deve falhar logo, e não dentro de dois subagentes paralelos.

= PARTE VI --- Apêndices
<parte-vi-apêndices>
== Apêndice A: Glossário
<apêndice-a-glossário>
#strong[ADR.] Registro de decisão arquitetural; documenta uma escolha difícil de reverter e o motivo.

#strong[Costura (seam).] Fronteira pública em que o comportamento é observado sem olhar para dentro do módulo. Os testes moram nas costuras.

#strong[DAG.] Grafo acíclico direcionado; aqui, o grafo de tickets com suas arestas de bloqueio.

#strong[Grilling.] Entrevista socrática em rodadas que separa fatos do agente de decisões do usuário.

#strong[Handoff.] Passagem estéril de contexto para outra sessão, pasta ou harness.

#strong[Módulo profundo.] Módulo com interface estreita que esconde implementação rica.

#strong[Porta de mão única.] Mudança difícil ou impossível de reverter; exige mais evidência no PR.

#strong[Smart Zone.] Faixa da janela de contexto em que o modelo raciocina com clareza; o manual a estima em até cerca de 150 mil tokens.

#strong[Teste tautológico.] Teste cujo valor esperado é calculado do mesmo jeito que o código, e que por isso nunca falha.

#strong[Tracer bullet.] Ticket vertical que entrega um comportamento observável atravessando todas as camadas.

== Apêndice B: Relatório de conformidade
<apêndice-b-relatório-de-conformidade>
Comparação entre o manual de 01/11/2026 e o repositório em 07/11/2026.

#figure(
  align(center)[#table(
    columns: (25.27%, 27.47%, 31.87%, 15.38%),
    align: (left,left,left,left,),
    table.header([Ponto verificado], [Manual], [Código ou upstream], [Situação],),
    table.hline(),
    [Arquivo do loop], [`matt-pocock.loop.json`], [Existe em `loops/`], [Conforme],
    [Validação do DAG], [Não citada], [6 nós, sem ciclos], [Conforme],
    [Sub-loop 2, Exploração], [Sub-loop próprio], [Sem nó no JSON], [Diverge],
    [Sub-loop 5, Diagnóstico], [Sub-loop próprio], [Sem nó no JSON], [Diverge],
    [Spec e tickets], [Um sub-loop], [Dois nós separados], [Diverge],
    [Modelos dos nós 1 e 4], [sonnet/haiku, sonnet/qwen], [Iguais em `routing`], [Conforme],
    [Modelos dos nós 2,3,5,6], [Declarados na matriz], [Ausentes em `routing`], [Diverge],
    [Template do Studio], [Ícone `git`], [Ícone `git`], [Conforme],
    [Ciclo do TDD], [Red-Green-Refactor], [Upstream: só vermelho e verde], [Diverge],
  )]
  , kind: table
  )

#strong[Leitura do relatório.] O JSON reorganiza os seis sub-loops do manual: Exploração e Diagnóstico ficam fora da esteira declarada, por serem desvios sob demanda, e Especificação vira dois nós, spec e tickets. O upstream da habilidade `tdd` tira a refatoração do ciclo e a coloca no review; este livro segue o upstream.

== Apêndice C: Referência rápida das habilidades
<apêndice-c-referência-rápida-das-habilidades>
As habilidades estão preservadas na cópia local da skill `aidd-pocock`, em `upstream/skills/`: 20 em `engineering/` e 7 em `productivity/`.

#figure(
  align(center)[#table(
    columns: (32.35%, 15.69%, 13.73%, 38.24%),
    align: (left,left,left,left,),
    table.header([Habilidade], [Pasta], [Etapa], [Função resumida],),
    table.hline(),
    [`grill-with-docs`], [engineering], [1], [Entrevista com rastro em arquivo],
    [`grill-me`], [productivity], [1], [Entrevista rápida sem estado],
    [`grilling`], [productivity], [1], [Primitivo da entrevista],
    [`domain-modeling`], [engineering], [1], [Termos sobrecarregados e ADRs],
    [`prototype`], [engineering], [2], [Código descartável em branch isolada],
    [`handoff`], [productivity], [2], [Passagem estéril de contexto],
    [`research`], [engineering], [2], [Pesquisa em fontes primárias],
    [`to-spec`], [engineering], [3], [Spec a partir do consenso],
    [`to-tickets`], [engineering], [3], [Tickets verticais com bloqueios],
    [`wayfinder`], [engineering], [3], [Mapa de decisões para iniciativas vagas],
    [`implement`], [engineering], [4], [Um ticket por sessão limpa],
    [`implement-spec`], [engineering], [4], [Tickets em paralelo com integração],
    [`tdd`], [engineering], [4], [Vermelho e verde em costuras combinadas],
    [`codebase-design`], [engineering], [4], [Vocabulário de módulos profundos],
    [`diagnosing-bugs`], [engineering], [5], [Loop de feedback antes de teorizar],
    [`improve-codebase-architecture`], [engineering], [5], [Caça a módulos rasos],
    [`triage`], [engineering], [5], [Qualifica issues externas],
    [`code-review`], [engineering], [6], [Review em Standards e Spec],
    [`pr`], [engineering], [6], [PR com evidência mínima e risco],
    [`retro`], [engineering], [6], [Falhas viram regras do ambiente],
    [`wizard`], [engineering], [6], [Scripts para etapas humanas],
    [`wait-what`], [productivity], [6], [Re-explicação sem jargão],
  )]
  , kind: table
  )

Fora dos sub-loops ficam `ask-matt` e `setup-matt-pocock-skills`, em `engineering/`, e `teach`, `to-questionnaire` e `writing-for-agents`, em `productivity/`.
