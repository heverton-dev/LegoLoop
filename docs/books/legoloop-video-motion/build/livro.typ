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
#let cor-accent   = rgb("#2563eb") // Cor da obra (Google Cobalto default)
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
  title: [Pipeline Canônico de Motion Video],
  author: "Ecossistema LegoLoop",
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
  )[MANUAL OPERACIONAL DIDÁTICO · VIDEO AS CODE]

  #v(6pt)

  // Título Principal
      #text(
      font: ("Segoe UI", "Arial"),
      size: 21pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.4pt,
    )[Motion Video]
    #v(2pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 13pt,
      weight: "bold",
      fill: cor-accent,
    )[Pipeline Canônico por Código]
  
  #v(4pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 8.8pt,
      fill: cor-suave,
      weight: "medium",
    )[Engenharia Reversa de Motion Design, Sincronia Fonética e Renderização por Código]
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
    #let etapas = ("Silêncio & Ataque +40ms", "Matriz Words & Punch RMS", "Vidro 3D & Tipografia Anton", "Cache SHA-256 & Bridge", "Multi-Aspecto (9:16, 16:9, 1:1)", "Broadcast Bundle & LegoPress", )
    #text(font: ("Segoe UI", "Arial"), size: 6.8pt, weight: "bold", fill: cor-suave, tracking: 1.4pt)[OS SEIS NÓS DO PIPELINE video-motion]
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
    text(font: ("Segoe UI", "Arial"), size: 7.8pt, weight: "bold", fill: cor-tinta, hyphenate: false)[Ecossistema LegoLoop],
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
  Este livro é o manual técnico e didático oficial do pipeline `video-motion` do LegoLoop. Ele documenta a engenharia reversa do método de motion video editing, a substituição integral de ferramentas proprietárias por código aberto e a arquitetura determinística dos quatro pilares: relógio fonético com corte de silêncio (+40ms), motion de vidro 3D em Remotion/CSS, cache criptográfico SHA-256 e entrega broadcast multi-aspecto (9:16, 16:9, 1:1).

  
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
            text(font: ("Segoe UI", "Arial"), size: 6.2pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[MOTION VIDEO],
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


= Como ler este livro
<como-ler-este-livro>
== Público e pré-requisitos
<público-e-pré-requisitos>
Esta obra destina-se a engenheiros de software, criadores de conteúdo técnico e operadores do ecossistema LegoLoop. Pressupõe noções básicas de linha de comando e processamento de mídia.

== Estrutura da obra
<estrutura-da-obra>
#figure(
  align(center)[#table(
    columns: (20%, 20%, 60%),
    align: (left,left,left,),
    table.header([Parte], [Nível], [Conteúdo],),
    table.hline(),
    [#strong[I]], [Macro], [Fundação filosófica, determinismo estrito e zero dependência],
    [#strong[II]], [Meso], [Arquitetura dos quatro pilares técnicos],
    [#strong[III]], [Micro], [Guia prático de execução e stop-gates de QA instantâneo],
    [#strong[IV]], [Prática], [Receitas multi-aspecto (9:16, 16:9, 1:1) e broadcast bundle],
    [#strong[V]], [Governança], [Rotina de atualização contínua por diff (LegoPress)],
  )]
  , kind: table
  )

= Parte I --- Fundação e Princípio Determinístico
<parte-i-fundação-e-princípio-determinístico>
== A Lei Fundamental: Vídeo é Código
<a-lei-fundamental-vídeo-é-código>
O pipeline `video-motion` foi concebido sob a lei de conduta inegociável do LegoLoop: #strong[tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de uma IA generativa.]

Na produção audiovisual contemporânea, a maioria dos criadores depende de pacotes fechados (After Effects, Premiere, CapCut) ou de plataformas SaaS baseadas em créditos e assinaturas opacas. Essa dependência cria pontos únicos de falha, custos imprevisíveis e impossibilidade de automação em lote.

No LegoLoop, a premissa é invertida:

- O layout visual é composto em #strong[código declarativo] (CSS3 3D, Web-Canvas e Remotion).
- Os tempos de transição são guiados pela matemática da fala (análise de áudio RMS e timestamps milimétricos do Whisper).
- Os assets visuais secundários são gerados via APIs públicas gratuitas (Google AI Studio) ou sintetizados deterministicamente por shaders procedurais em Canvas.

== Comparativo de Eficiência Arquitetural
<comparativo-de-eficiência-arquitetural>
#figure(
  align(center)[#table(
    columns: (21.28%, 40.43%, 38.3%),
    align: (left,left,left,),
    table.header([Dimensão], [Fluxo Tradicional (SaaS / NLE)], [Pipeline LegoLoop (`video-motion`)],),
    table.hline(),
    [#strong[Corte de Respiração]], [Manual em timeline ou app proprietário], [Automático (-30dB RMS + 40ms ataque)],
    [#strong[Sincronia de Legenda]], [Digitação ou exportação com drift], [`words.json` milimétrico sem desvio],
    [#strong[Estilo Visual 3D]], [Plugins pagos de After Effects], [Remotion / CSS3 com blur procedural],
    [#strong[Custo por Minuto]], [Assinaturas mensais recorrentes], [Custo zero de licença e 100% local],
    [#strong[Versionamento]], [Binários pesados (.aep / \.prproj)], [Código puro versionado em Git],
  )]
  , kind: table
  )

= Parte II --- A Arquitetura dos Quatro Pilares
<parte-ii-a-arquitetura-dos-quatro-pilares>
O pipeline é estruturado em quatro motores independentes, integrados por um grafo acíclico dirigido (DAG):

#figure(image("/img/diagrama_436e5a57.png", width: 85.0%, alt: "Diagrama de Arquitetura"),
  caption: [
    Diagrama de Arquitetura
  ]
)

== 1. Relógio Mestre Fonético (Speech-Sync Engine)
<relógio-mestre-fonético-speech-sync-engine>
A linha do tempo do vídeo não é definida por segundos arbitrários, mas pela cadência real da locução:

+ #strong[Corte Cirúrgico de Pausas:] Todo silêncio abaixo do threshold de -30dB com duração superior a 400ms é removido, preservando estritamente #strong[\+40ms de ataque pré-fala] para que nenhuma consoante oclusiva seja deglutida.
+ #strong[Alinhamento Fonético Milimétrico:] Geração da matriz `words.json`, contendo cada vocábulo com seus marcadores de `start_ms` e `end_ms`.
+ #strong[Detecção de Ênfase Vocal:] O motor calcula a energia RMS por palavra. Vocábulos com amplitude 40% superior à média do segmento recebem a flag `is_punch_word = true`.
+ #strong[Normalização Broadcast EBU R128:] Ajuste de ganho para teto de transmissão em #strong[\-14 LUFS] com true peak a -1.0 dBFS.

== 2. Compositor Visual Reativo (Remotion & CSS 3D)
<compositor-visual-reativo-remotion-css-3d>
Em vez de renderizar frames estáticos pesados, o visual é gerado via DOM reativo acelerado por hardware:

- #strong[Cartões 3D de Vidro:] Superfícies translúcidas com `backdrop-filter: blur(16px)`, bordas sutis com iluminação especular e transição suave de push-in.
- #strong[Tipografia Cinética Anton:] Exibição com destaque cromático instantâneo sobre as palavras de ênfase identificadas na etapa fonética.
- #strong[Auto-Framing Adaptativo:] Recorte inteligente mantendo o sujeito no terço superior de segurança.

== 3. Pipeline de Assets Determinístico com Cache Criptográfico
<pipeline-de-assets-determinístico-com-cache-criptográfico>
- #strong[Chave de Cache SHA-256:] Cada solicitação de imagem secundária gera uma assinatura criptográfica a partir do prompt e dimensões. Se já existe no disco, o custo de rede é zero.
- #strong[Bridge com Fallback Local:] Conexão nativa com a API do Google AI Studio. Em caso de indisponibilidade de rede ou ausência de chave, o gerador procedural sintetiza texturas abstratas locais sem interromper o fluxo.

== 4. Renderizador Multi-Aspecto e Bundle Broadcast
<renderizador-multi-aspecto-e-bundle-broadcast>
- #strong[Exportação Paralela:] Suporte nativo aos formatos `9:16` (Vertical), `16:9` (Widescreen) e `1:1` (Quadrado).
- #strong[Entrega Broadcast Completa:] Geração simultânea do pacote em `dist/<job_id>/` contendo `.mp4`, legendas `.srt` e `.vtt`, arquivo de timecode `fcpxml` (v1.9) e copy pronta para redes em `POST_COPY.md`.

= Parte III --- Guia Prático de Execução
<parte-iii-guia-prático-de-execução>
== Tríade de Acionamento
<tríade-de-acionamento>
O pipeline pode ser executado pelas três vias canônicas do ecossistema LegoLoop, produzindo exatamente os mesmos artefatos e garantindo auditoria completa.

=== 1. Via Linha de Comando (CLI)
<via-linha-de-comando-cli>
O binário `legoloop` expõe o comando unificado `video-motion` e seus subcomandos:

```bash
# Execução completa da esteira de ponta a ponta
legoloop video-motion run --input video.mp4 --aspect 9:16 --output dist/

# Execução cirúrgica isolada por etapas:
legoloop video-motion silence --input audio.wav --threshold -30 --pad 40
legoloop video-motion punch   --words words.json --audio audio.wav
legoloop video-motion assets  --prompt "cyberpunk glass ui" --cache .cache/
legoloop video-motion preview --input video.mp4 --output qa_sheet.jpg
legoloop video-motion bundle  --input video.mp4 --output dist/job_01/
```

=== 2. Parâmetros e Modos de Operação
<parâmetros-e-modos-de-operação>
#figure(
  align(center)[#table(
    columns: (20.51%, 25.64%, 53.85%),
    align: (left,left,left,),
    table.header([Subcomando], [Flag Principal], [Descrição e Comportamento],),
    table.hline(),
    [`silence`], [`--pad 40`], [Margem de ataque em milissegundos pré-fala],
    [`punch`], [`--words words.json`], [Análise de amplitude RMS e destaque Anton],
    [`assets`], [`--prompt "..."`], [Consulta ao cache SHA-256 e bridge AI Studio],
    [`preview`], [`--output qa.jpg`], [Renderização ultra-rápida dos 12 frames-chave],
    [`bundle`], [`--output dist/`], [Consolidação do master, legendas e post copy],
  )]
  , kind: table
  )

=== 3. Via Grafo DAG e LegoLoop Studio
<via-grafo-dag-e-legoloop-studio>
O pipeline é formalmente orquestrado pelo arquivo `processos/video-motion/pipeline.dag.json`:

```bash
legoloop run processos/video-motion/pipeline.dag.json --var input_media="video.mp4"
```

No painel visual do #strong[LegoLoop Studio], o operador acompanha o progresso nó a nó:

- #strong[Nó 1 (Audio Clock):] Inspeção dos marcadores de silêncio detectados.
- #strong[Nó 2 (Phonetic RMS):] Gráfico de energia sonora com palavras de ênfase destacadas.
- #strong[Nó 3 (Compositor 3D):] Pré-visualização dos shaders e fontes compiladas.
- #strong[Nó 4 (Export):] Geração do bundle final e telemetria de renderização.

== O Stop-Gate de Qualidade Visual (`qa_sheet.jpg`)
<o-stop-gate-de-qualidade-visual-qa_sheet.jpg>
Antes do render final em alta resolução, o motor gera em menos de 2 segundos a folha de contato `qa_sheet.jpg`:

- Grade matricial 4x3 contendo os 12 frames mais significativos do vídeo.
- Inspeção visual imediata das legendas, do enquadramento facial e dos reflexos dos cartões de vidro.
- O pipeline pausa aguardando a confirmação do operador para liberar o empacotamento broadcast.

=== Protocolo de Validação e Ajustes no Gate
<protocolo-de-validação-e-ajustes-no-gate>
Caso algum artefato visual apresente inconsistência na folha de contato:

+ #strong[Recorte Facial Desalinhado:] Ajustar a sensibilidade do auto-framing via parâmetro `--face-margin`.
+ #strong[Ênfase Vocal Excessiva:] Recalibrar o limiar de RMS com a flag `--rms-boost 1.2`.
+ #strong[Cores ou Reflexos:] Modificar o shader CSS diretamente no arquivo de estilos antes de autorizar o render final.

= Parte IV --- Receitas Multi-Aspect Ratio e Broadcast Bundle
<parte-iv-receitas-multi-aspect-ratio-e-broadcast-bundle>
== Matriz de Formatos e Áreas de Segurança
<matriz-de-formatos-e-áreas-de-segurança>
Cada rede social possui restrições visuais específicas (sobreposição de botões de like, comentários e legenda de rodapé). O motor ajusta o layout com base na matriz abaixo:

#figure(
  align(center)[#table(
    columns: (19.23%, 19.23%, 32.05%, 29.49%),
    align: (left,left,left,left,),
    table.header([Formato], [Resolução], [Destino Principal], [Margem de Segurança],),
    table.hline(),
    [#strong[9:16]], [1080x1920], [TikTok, Reels, Shorts], [120px topo, 280px base],
    [#strong[16:9]], [1920x1080], [YouTube Desktop, Cursos], [60px uniforme],
    [#strong[1:1]], [1080x1080], [Instagram Feed, LinkedIn], [80px topo e base],
  )]
  , kind: table
  )

== Anatomia do Pacote de Entrega (`dist/<job_id>/`)
<anatomia-do-pacote-de-entrega-distjob_id>
Ao concluir a renderização, a pasta de entrega é montada com todos os ativos requeridos para distribuição profissional:

```
dist/job_20261107_01/
├── video_final_9x16.mp4      # Master de transmissão codificado em H.264 / AAC
├── legendas.srt              # Legenda universal para players externos
├── legendas.vtt              # Legenda web com estilização CSS
├── timeline.fcpxml           # XML v1.9 importável no DaVinci Resolve Studio
├── qa_sheet.jpg              # Prova de contato dos 12 frames-chave
└── POST_COPY.md              # Texto de postagem com hashtags e timing de publicação
```

= Parte V --- Rotina de Atualização por Diff (LegoPress)
<parte-v-rotina-de-atualização-por-diff-legopress>
== Governança Viva de Documentação
<governança-viva-de-documentação>
No ecossistema LegoLoop, nenhuma documentação pode se tornar obsoleta em relação ao código real. Sempre que um arquivo do pipeline (`internal/videomotion/` ou `processos/video-motion/`) for alterado via commit ou refatoração, o motor #strong[LegoPress] sincroniza a obra técnica automaticamente.

== Comandos de Manutenção do Livro
<comandos-de-manutenção-do-livro>
```bash
# Auditar a integridade geométrica do PDF gerado
legoloop legopress auditar docs/books/video-motion/2026-11-07_PIPELINE-MOTION-VIDEO.pdf

# Atualizar capítulos a partir dos diffs recentes do Git
legoloop legopress atualizar docs/books/video-motion/

# Recompilar o livro completo em caso de alterações nas partes
bash docs/books/compilar.sh video-motion
```

A esteira de publicação gera automaticamente o manuscrito consolidado em Markdown, o documento tipográfico final em Typst, o PDF de alta resolução (170x240mm) e as imagens de visualização para cada página na pasta `preview/`.
