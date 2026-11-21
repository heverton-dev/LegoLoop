#set page(
  paper: "a4",
  flipped: false,
  margin: (top: 2.5cm, bottom: 2.5cm, left: 2.5cm, right: 2.5cm),
  header: none,
  footer: none,
  fill: rgb("#fafafa"),
)

#set text(
  font: ("Georgia", "Segoe UI", "Arial"),
  size: 11pt,
  lang: "pt",
)

#align(center + horizon)[
  #block(width: 8.5cm)[
    #image("../ativos/lego_brick_part1.svg", width: 100%)
  ]

  #v(28pt)

  #text(
    size: 26pt,
    weight: "bold",
    fill: rgb("#0f172a"),
    tracking: 0.5pt,
  )[Macro: Fundação e Filosofia\ da Edição como Código]

  #v(14pt)

  #text(
    size: 12pt,
    style: "italic",
    fill: rgb("#475569"),
  )[A transição do paradigma NLE manual para a esteira determinística de alta fidelidade]

  #v(32pt)

  #line(length: 4cm, stroke: 1.5pt + rgb("#10b981"))
]
