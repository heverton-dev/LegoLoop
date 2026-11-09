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
  title: [Pipeline Universal de Engenharia Reversa Audiovisual],
  author: "Ecossistema LegoLoop",
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
  inset: (x: 4pt, y: 2.4pt),
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
  )[MANUAL OPERACIONAL DIDÁTICO · FORENSIC REVERSE ENGINEERING]

  #v(12pt)

  // Título Principal e Destaque
      #text(
      font: ("Segoe UI", "Arial"),
      size: 28pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.5pt,
    )[Video REA]
    #v(4pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 15.5pt,
      weight: "bold",
      fill: cor-accent,
    )[Engenharia Reversa Audiovisual]
  
  #v(8pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 9.6pt,
      fill: cor-suave,
      weight: "medium",
    )[Dissecação Forense de Vídeos, Dedução de Prompts Ocultos e Fábrica Autônoma de Pipelines]
  ]

  #v(16pt)

  // Arte: esteira de sub-loops desenhada em Typst
    #block(width: 100%, radius: 6pt, clip: true, fill: rgb("#f8fafc"), stroke: 0.5pt + rgb("#cbd5e1"), inset: 16pt)[
    #let etapas = ("Crawler & Ingestão Stealth", "Dissecação Forense & OCR", "Curadoria & Soberania R\$ 0", "Stop-Gate Humano Interativo", "Tríade ISO & Golden Fixture", "LegoPress & Skill Canônica", )
    #text(font: ("Segoe UI", "Arial"), size: 7pt, weight: "bold", fill: cor-suave, tracking: 1.6pt)[OS SEIS NÓS DO PIPELINE legoloop-video-rea]
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
    text(font: ("Segoe UI", "Arial"), size: 8.8pt, weight: "bold", fill: cor-tinta)[Ecossistema LegoLoop],
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
  Este livro é o manual técnico oficial do pipeline `legoloop-video-rea`. Ele documenta a esteira industrial que ingere vídeos de YouTube, Instagram e TikTok, extrai metadados, disseca prompts explícitos e deduções ocultas, realiza OCR em telas rápidas de código, cruza ferramentas com a matriz de peças locais Tier Ouro, submete o dossiê ao stop-gate humano e dispara a fábrica de novos pipelines com a Tríade de Governança ISO e seu livro didático correspondente.

  
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
          text(font: ("Segoe UI", "Arial"), size: 5.8pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[VIDEO REA],
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
== Público e pré-requisitos
<público-e-pré-requisitos>
Esta obra destina-se a engenheiros de IA, desenvolvedores de software e operadores do ecossistema LegoLoop interessados em desconstruir tutoriais, hacks e metodologias audiovisuais públicas, convertendo-os em software versionável e determinístico.

== Estrutura da obra
<estrutura-da-obra>
#figure(
  align(center)[#table(
    columns: (20%, 20%, 60%),
    align: (left,left,left,),
    table.header([Parte], [Nível], [Conteúdo],),
    table.hline(),
    [#strong[I]], [Macro], [Visão geral da engenharia reversa e soberania técnica],
    [#strong[II]], [Meso], [Dissecação forense de prompts ocultos e OCR de código],
    [#strong[III]], [Micro], [Curadoria contra a matriz do cofre e Stop-Gate humano],
    [#strong[IV]], [Prática], [Fábrica da Tríade ISO (INTENT, SPEC, PLAN, ROUTE) e Fixtures],
    [#strong[V]], [Governança], [Nomenclatura legoloop-\*, publicação LegoPress e Skills],
  )]
  , kind: table
  )

= Parte I --- Visão Geral e Soberania Técnica
<parte-i-visão-geral-e-soberania-técnica>
== O Desafio da Engenharia Reversa Audiovisual
<o-desafio-da-engenharia-reversa-audiovisual>
Na internet atual, o conhecimento técnico mais avançado sobre inteligência artificial generativa, automação e motion design não está publicado em artigos acadêmicos ou manuais formais, mas em vídeos rápidos postados no YouTube, Instagram Reels e TikTok.

Esses conteúdos, contudo, apresentam três barreiras críticas para a engenharia de software: 1. #strong[Omissão Estratégica:] Autores frequentemente ocultam os prompts mestres ou guardam as ferramentas reais em paywalls e comunidades pagas. 2. #strong[Dependência Comercial:] As receitas ensinadas exigem plataformas SaaS caras, criando dependência de cartões de crédito e serviços fechados. 3. #strong[Falta de Reprodutibilidade:] Vídeos mostram apenas o resultado final editado, sem expor o código-fonte ou as matrizes de dados.

== A Solução Institucional: `legoloop-video-rea`
<a-solução-institucional-legoloop-video-rea>
O pipeline `legoloop-video-rea` foi concebido para fechar essa lacuna. Ele opera como um perito forense automatizado: ingere qualquer vídeo a partir de uma simples URL, quebra as barreiras de login de forma furtiva, lê as entrelinhas e as telas de código congeladas, cruza as necessidades contra o acervo local de código aberto e projeta um novo pipeline canônico pronto para ser fabricado.

#figure(
  align(center)[#table(
    columns: (19.78%, 39.56%, 40.66%),
    align: (left,left,left,),
    table.header([Dimensão], [Análise Manual Convencional], [Esteira Automatizada `video-rea`],),
    table.hline(),
    [#strong[Tempo de Análise]], [2 a 4 horas por vídeo longo], [Menos de 30 segundos por alvo],
    [#strong[Captura de Prompts]], [Apenas o que for digitado na tela], [Reconstrução de prompts implícitos],
    [#strong[Leitura de Código]], [Digitação manual com erros], [OCR local frame a frame em 4K],
    [#strong[Custo de Execução]], [R\$ 180 a R\$ 500/mês em assinaturas], [R\$ 0,00 com soberania local],
    [#strong[Saída Entregue]], [Anotações esparsas em bloco de notas], [Tríade ISO, DAG, CLI e Livro Didático],
  )]
  , kind: table
  )

= Parte II --- Dissecação Forense e Prompts Ocultos
<parte-ii-dissecação-forense-e-prompts-ocultos>
== Decodificando o Não-Dito
<decodificando-o-não-dito>
O motor forense opera sob o postulado de que o áudio falado é apenas a camada superficial da instrução. A engenharia reversa real precisa extrair o que o autor executou no computador.

=== 1. Ingestão Stealth Multi-Plataforma
<ingestão-stealth-multi-plataforma>
- #strong[YouTube:] Captura de streams diretas de áudio em PCM 48kHz e metadados via curl/yt-dlp sem disparo de anti-bot.
- #strong[Instagram e TikTok:] Acoplamento com o navegador furtivo #strong[Obscura] e Chrome CDP do sistema operacional, reaproveitando cookies locais sem risco de CAPTCHA ou bloqueios temporários.

=== 2. Dissecação de Prompts Explícitos e Ocultos
<dissecação-de-prompts-explícitos-e-ocultos>
O dissecador examina a transcrição fonética e a descrição procurando blocos formais de instrução. Quando o autor diz apenas #emph["aqui eu gerei uma textura futurista"], o motor analisa os atributos visuais do frame e reconstrói o prompt mestre: - Parâmetros de iluminação (specular, rim-light, depth). - Estilo de renderização (translucent glass, backdrop-blur 16px). - Restrições geométricas (aspect ratio 9:16, 60fps clean loop).

=== 3. OCR Forense em Telas Rápidas de Código
<ocr-forense-em-telas-rápidas-de-código>
Muitos criadores exibem seu terminal ou editor de código por frações de segundo durante transições de cena: - O módulo de visão multimodal local isola esses frames de alta frequência. - Executa reconhecimento de caracteres com foco em sintaxe de linha de comando (`ffmpeg`, `styled-components`, `css shaders`). - Cataloga essas chamadas na lista `HiddenTools`, revelando dependências omitidas no discurso.

= Parte III --- Curadoria Matriz e o Stop-Gate Humano
<parte-iii-curadoria-matriz-e-o-stop-gate-humano>
== Conexão com o Cofre de Peças Tier Ouro
<conexão-com-o-cofre-de-peças-tier-ouro>
Identificar ferramentas comerciais não basta: a engenharia reversa no LegoLoop visa à independência tecnológica. O módulo curador cruza cada dependência mapeada contra o catálogo de peças do ecossistema (`internal/matriz`):

#figure(
  align(center)[#table(
    columns: (23.76%, 38.61%, 37.62%),
    align: (left,left,left,),
    table.header([Ferramenta do Vídeo], [Peça Local do Cofre], [Ganho de Soberania],),
    table.hline(),
    [#strong[Higgsfield AI]], [Google AI Studio + Texturas Procedurais], [Custo zero, sem filas e determinístico],
    [#strong[After Effects / Premiere]], [Remotion + CSS3 3D Shaders], [Versionamento Git e render em pipeline],
    [#strong[Whisper API Paga]], [Whisper Local / FFmpeg SilenceRemove], [Processamento offline a 48kHz],
    [#strong[Claude Pro Web]], [9Router Local / Modelos Flash], [Execução estéril sem consumo de cota],
  )]
  , kind: table
  )

== O Stop-Gate de Curadoria (`SG-CURADORIA`)
<o-stop-gate-de-curadoria-sg-curadoria>
Nenhum arquivo de código é gerado antes do consentimento soberano do desenvolvedor. O Stop-Gate apresenta o dossiê: 1. Resumo do método desvendado e prompts reconstruídos. 2. Tabela de substituições sugeridas (Comercial vs.~Tier Ouro). 3. Cálculo da Economia de Soberania (ex.: economia de R\$ 180,00/mês). 4. As três sugestões de nomes canônicos no padrão `legoloop-*`.

O operador tem três opções diretas: - #strong[\[1\] Aprovar:] Libera a emissão imediata da governança e fábrica de pipeline. - #strong[\[2\] Ajustar:] Permite trocar manualmente alguma peça antes de construir. - #strong[\[0\] Cancelar:] Descarta o dossiê de forma estéril sem resíduos no disco.

= Parte IV --- Fábrica de Governança e Golden Fixtures
<parte-iv-fábrica-de-governança-e-golden-fixtures>
== Da Curadoria ao Software Executável
<da-curadoria-ao-software-executável>
Assim que o operador aprova a curadoria no Stop-Gate, o pipeline convoca a esteira de fabricação do LegoLoop:

#figure(image("/img/diagrama_373f383a.png", width: 85.0%, alt: "Diagrama de Arquitetura"),
  caption: [
    Diagrama de Arquitetura
  ]
)

=== 1. Emissão da Tríade ISO de Governança
<emissão-da-tríade-iso-de-governança>
A fábrica gera os quatro artefatos canônicos sincronizados: - #strong[`YYYY-MM-DD_INTENT_<SLUG>.md`:] O propósito de negócio e o valor entregue. - #strong[`YYYY-MM-DD_SPEC_<SLUG>.md`:] Contratos formais de dados e cenários BDD binários. - #strong[`YYYY-MM-DD_PLAN_<SLUG>.md`:] Decomposição em tickets verticais com rastreabilidade 1:1. - #strong[`YYYY-MM-DD_ROUTE_<SLUG>.json`:] O grafo de execução estéril com hashes SHA-256 de integridade.

=== 2. A Golden Fixture Sintética (3 Segundos)
<a-golden-fixture-sintética-3-segundos>
Testar pipelines de vídeo não pode depender de renderizar clipes de 10 minutos: - A fábrica sintetiza um clipe MP4 ultraleve de 3.0 segundos a 60fps. - Insere uma faixa de áudio de 48kHz contendo um pico vocal para validar o silêncio de ataque (+40ms). - Permite que a suíte de testes TDD do novo pipeline rode localmente em menos de 2 segundos.

= Parte V --- Nomenclatura, LegoPress e Skills
<parte-v-nomenclatura-legopress-e-skills>
== A Regra de Ouro da Nomenclatura: `legoloop-*`
<a-regra-de-ouro-da-nomenclatura-legoloop->
Todo pipeline concebido a partir de engenharia reversa herda obrigatoriamente a assinatura institucional do ecossistema: - O novo diretório de loops adota `loops/<slug>/pipeline.dag.json` com campo `"name": "legoloop-<slug>"`. - A nova CLI é registrada como subcomando de primeiro nível: `legoloop <slug>`. - A nova skill do agente é distribuída no cofre: `skills/legoloop-<slug>/SKILL.md`.

== Publicação Automática do Livro Didático
<publicação-automática-do-livro-didático>
Nenhum pipeline é entregue como "caixa preta". A fábrica convoca imediatamente o motor #strong[LegoPress]: 1. Cria o manifesto canônico `livro.json` em `docs/books/<slug>/`. 2. Estrutura as partes didáticas (Macro, Meso, Micro, Prática e Governança). 3. Compila o manuscrito consolidado e o documento tipográfico Typst. 4. Produz o arquivo PDF em padrão editorial de 170x240mm e suas imagens de visualização em `preview/`. 5. Roda a auditoria geométrica (`legoloop legopress auditar`) para garantir que nenhum cabeçalho ou página órfã seja entregue.

Dessa forma, o conhecimento extraído de qualquer vídeo público é perenizado, auditável, de custo zero e pronto para ser executado por agentes e humanos.
