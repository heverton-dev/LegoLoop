#let cor-tinta  = rgb("#0f172a") // Ardósia nobre
#let cor-fundo  = rgb("#ffffff")

#set page(
  width: 170mm,
  height: 240mm,
  fill: cor-fundo,
  margin: (top: 2.0cm, bottom: 2.0cm, left: 1.8cm, right: 1.8cm),
)

// Folha Divisória de Parte: 100% centralizada na horizontal e vertical, sem cabeçalho e sem rodapé
#page(header: none, footer: none, numbering: none, margin: 2.2cm)[
  #align(center + horizon)[
    // Peça de Lego 3D Real com iluminação, studs cilíndricos e relevo autêntico
    #block(width: 100%, inset: (bottom: 10pt))[
      #image("ativos/lego_brick_part1.svg", width: 52%)
    ]
    #v(16pt)
    #block(width: 90%)[
      #text(font: ("Segoe UI", "Arial"), size: 21pt, weight: "bold", fill: cor-tinta, tracking: -0.4pt)[
        Macro: Fundação & Filosofia da Edição como Código
      ]
    ]
    #v(10pt)
    #block(width: 80%)[
      #text(font: ("Georgia", "Segoe UI"), size: 10.5pt, fill: rgb("#64748b"), style: "italic")[
        O colapso da montagem manual e a esteira soberana a custo zero
      ]
    ]
  ]
]

// Página seguinte com o conteúdo normal
#page(numbering: "1")[
  #text(font: ("Segoe UI", "Arial"), size: 13pt, weight: "bold", fill: cor-tinta)[1.1 O Colapso da Edição Tradicional]
  #v(8pt)
  #text(font: ("Georgia", "Segoe UI"), size: 10.3pt, fill: rgb("#292524"))[
    A produção de conteúdo audiovisual em escala enfrenta um gargalo histórico: a pós-produção manual. Em produções faladas, um editor profissional gasta horas em tarefas repetitivas.
  ]
]
