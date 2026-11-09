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
#let cor-sucesso  = rgb("#16a34a")                                    // Verde Esmeralda Google
#let cor-teoria   = rgb("#7c3aed")                                    // Violeta DeepMind

#set document(
  title: [Pipeline Canônico de Edição Audiovisual Assistida por IA],
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
  )[MANUAL DIDÁTICO CORPORATIVO · VIDEO AS CODE]

  #v(6pt)

  // Título Principal
      #text(
      font: ("Segoe UI", "Arial"),
      size: 21pt,
      weight: "black",
      fill: cor-tinta,
      tracking: -0.4pt,
    )[Video as Code]
    #v(2pt)
    #text(
      font: ("Segoe UI", "Arial"),
      size: 13pt,
      weight: "bold",
      fill: cor-accent,
    )[Pipeline Canônico de Edição Audiovisual com IA]
  
  #v(4pt)

  // Subtítulo
  #block(width: 95%)[
    #text(
      font: ("Segoe UI", "Arial"),
      size: 8.8pt,
      fill: cor-suave,
      weight: "medium",
    )[Esteira determinística de pós-produção local, separação rígida entre cognição e frames, e entrega híbrida FCPXML]
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
    #let etapas = ("Preparação & Folga +40ms", "Corte Semântico & Notas", "Storyboard Gráfico Remotion", "Câmera YuNet (Olhos a 35%)", "Ducking -18dB & SFX na Frenagem", "Entrega Não-Destrutiva FCPXML", )
    #text(font: ("Segoe UI", "Arial"), size: 6.8pt, weight: "bold", fill: cor-suave, tracking: 1.4pt)[TOPOLOGIA DAG DE 12 NÓS E GOVERNANÇA COM 3 STOP-GATES]
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
    text(font: ("Segoe UI", "Arial"), size: 7.2pt, fill: cor-suave, hyphenate: false)[08/11/2026 · LegoLoop Foundation]
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
  Este livro técnico documenta a arquitetura, a formalização matemática e a operação da esteira industrial de edição audiovisual video-edition do ecossistema LegoLoop.

  Projetada sob o paradigma Video as Code (VaC), a esteira implementa separação estrita entre cognição textual por LLMs e computação de mídia determinística (FFmpeg, WhisperX, OpenCV YuNet e Remotion), eliminando mensalidades de SaaS proprietários (custo R\$ 0,00/mês) e garantindo entrega não-destrutiva em timelines abertas (FCPXML/EDL) governadas por um Grafo Acíclico Dirigido (DAG) de 12 nós e 3 Stop-Gates humanos inegociáveis. Obra fundamentada nos princípios de edição com IA derivados e canonizados por engenharia reversa soberana.

  
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
    Data da compilação: 08/11/2026 · LegoLoop Foundation
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
            text(font: ("Segoe UI", "Arial"), size: 6.2pt, fill: cor-tinta, weight: "bold", tracking: 1.2pt)[VIDEO EDITION],
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
== A quem este livro se destina
<a-quem-este-livro-se-destina>
Esta obra foi estruturada para atender três perfis complementares do ecossistema audiovisual e tecnológico:

+ #strong[Criadores de Conteúdo e Editores de Vídeo:] Que desejam automatizar o trabalho repetitivo e exaustivo de corte de silêncios, enquadramento dinâmico e inserção de cartões gráficos, sem abrir mão do controle artístico final sobre a narrativa e sem depender de ferramentas fechadas com assinaturas abusivas em dólar.
+ #strong[Engenheiros de IA e Desenvolvedores de Agentes:] Que precisam de garantias formais de determinismo, tolerância zero a alucinações de timestamps e contratos de dados estritos entre modelos cognitivos de linguagem (LLMs) e motores pesados de processamento de mídia (FFmpeg, NVENC, Node.js).
+ #strong[Líderes de Operações de Mídia e Produtoras:] Que buscam padronizar esteiras industriais de pós-produção com medição auditável de tokens, custos financeiros por minuto renderizado e entrega não-destrutiva em timelines abertas (FCPXML/EDL) compatíveis com estações DaVinci Resolve e Premiere Pro.

== A Estrutura da Obra (Macro $arrow.r$ Meso $arrow.r$ Micro)
<a-estrutura-da-obra-macro-rightarrow-meso-rightarrow-micro>
A narrativa deste manual segue o rigor didático do LegoLoop:

- #strong[Parte I --- Macro (Fundação & Filosofia do Video as Code):] Contrasta a edição tradicional baseada em monólitos NLE proprietários com a esteira soberana e estabelece a Lei Áurea da Separação Rígida entre cognição textual e precisão matemática de frames.
- #strong[Parte II --- Meso (Os 13 Guias Autorais e a Engenharia de Pós-Produção):] Detalha exaustivamente cada um dos 13 pilares do método de Lucas Leite: alinhamento fonético com pre-roll de +40ms, extirpação automática de notas de gravação, enquadramento facial YuNet pós-motion a 35% do topo, motion design procedural via Remotion e sonorização com ducking inteligente de -18dB.
- #strong[Parte III --- Micro (A Topologia DAG, Stop-Gates e Operação):] Apresenta a matriz executável do Grafo Acíclico Dirigido (12 nós, 11 camadas), os 3 Stop-Gates humanos, os comandos de terminal do CLI `legoloop run` e a matriz de resolução de falhas (troubleshooting).

== O Manifesto da Edição Determinística
<o-manifesto-da-edição-determinística>
A edição de vídeo assistida por IA falhou historicamente por tentar transformar modelos de linguagem em "editores universais". LLMs são calculadoras semânticas probabilísticas; forçá-los a calcular segundos, milissegundos e taxas de quadros por segundo resulta invariavelmente em deriva cumulativa de sincronia labial e cortes mastigados.

O LegoLoop resolve este dilema através de três leis fundamentais:

#painel-teoria("LEI #1: O LLM NUNCA CALCULA FRAMES", [ O modelo de inteligência artificial recebe dados em texto estruturado e toma decisões exclusivamente sobre frases e sentenças completas. O cálculo de timestamps, limites de áudio, janelas de RMS e conversão para timecode de timeline é tarefa exclusiva de scripts matemáticos determinísticos em Python, Go ou Node.js.])

#painel-alerta("LEI #2: TODA FALA EXIGE FOLGA PRÉ-ROLL (+40MS)", [ O ponto de corte de uma sentença falada nunca coincide com o limiar de energia do Voice Activity Detection (VAD). O algoritmo impõe compulsoriamente uma margem de segurança de +40 milissegundos antes da primeira consoante para evitar decapitação fonética.])

#painel-sucesso("LEI #3: A ENTREGA É SEMPRE NÃO-DESTRUTIVA", [ Nenhum pipeline sério de produção entrega apenas um arquivo MP4 achatado. A saída primária da esteira é uma timeline aberta em XML (FCPXML), com handles (sobras de corte) e faixas isoladas, permitindo que o editor humano execute o polimento final no DaVinci Resolve ou Premiere.])

= Parte I --- Macro: Fundação & Filosofia da Edição como Código
<parte-i-macro-fundação-filosofia-da-edição-como-código>
== 1.1 O Colapso da Edição Tradicional de Vídeos
<o-colapso-da-edição-tradicional-de-vídeos>
A produção de conteúdo audiovisual em escala enfrenta um gargalo histórico: a pós-produção manual. Em produções faladas (aulas, podcasts, tutoriais de tecnologia, vídeos corporativos), um editor profissional gasta em média de 3 a 5 horas para cada 15 minutos de vídeo finalizado.

A maior parte desse tempo não é gasta em decisões criativas de alto nível, mas em tarefas puramente mecânicas:

- Ouvir repetições de frases geradas por erros do apresentador;
- Cortar pausas de respiração mortas e silêncios constrangedores;
- Identificar quando o locutor fala #emph["corta isso aqui"] ou #emph["vou recomeçar"] (as chamadas notas de edição);
- Reenquadrar manualmente o vídeo para criar dinamismo visual;
- Sincronizar efeitos sonoros de transição na entrada de cartões gráficos.

A indústria de software respondeu a esse desafio criando plataformas proprietárias em nuvem (Descript, Runway, Adobe Sensei, Riverside). No entanto, essas soluções introduziram três novos problemas graves:

+ #strong[Custo Recorrente Cumulativo:] Mensalidades que variam de R\$ 350 a mais de R\$ 600 por mês, frequentemente cobradas em dólar com limites de minutos de processamento;
+ #strong[Perda de Privacidade e Soberania:] Upload obrigatório de gigabytes de imagens e áudio de gravações confidenciais para servidores de terceiros;
+ #strong[Cortes Rígidos e Destrutivos:] A maioria dos serviços entrega apenas um vídeo MP4 final exportado. Caso o criador discorde de um corte ou queira alterar uma palavra num texto animado, é necessário reenviar o vídeo ou perder o trabalho feito.

#horizontalrule

== 1.2 O Paradigma Video as Code & Soberania Zero Custo
<o-paradigma-video-as-code-soberania-zero-custo>
O ecossistema LegoLoop estabelece o conceito de #strong[Video as Code (VaC)]: a montagem de um vídeo falado deixa de ser um ato artesanal e imprevisível em uma interface gráfica fechada e passa a ser a execução de um pipeline determinístico, composto por blocos de software locais, auditáveis e com custo de licença igual a zero.

#figure(image("img/diagrama_edicao_tradicional_vs_codigo.svg", alt: "Edição Tradicional (NLE Monolítico) vs Metodologia Soberana LegoLoop"),
  caption: [
    Edição Tradicional (NLE Monolítico) vs Metodologia Soberana LegoLoop
  ]
)

A tabela a seguir contrasta a pilha comercial identificada na engenharia reversa do vídeo com a matriz soberana adotada no LegoLoop:

#figure(
  align(center)[#table(
    columns: (22.52%, 33.33%, 27.93%, 16.22%),
    align: (left,left,left,center,),
    table.header([Dimensão Operacional], [Pilha Proprietária Comercial], [Pilha Soberana LegoLoop], [Economia / Ganho],),
    table.hline(),
    [#strong[Ambiente de Edição]], [Adobe Premiere Pro (`US$ 35/mês`)], [#strong[FFmpeg + DaVinci Resolve (Free)]], [100% Gratuito],
    [#strong[Transcrição & Alinhamento]], [Descript / Whisper API (`US$ 24/mês`)], [#strong[Whisper Local (WhisperX)]], [Hardware Local],
    [#strong[Rastreamento Facial]], [Plugins After Effects (`US$ 15/mês`)], [#strong[OpenCV YuNet (Python)]], [Ultra Rápido],
    [#strong[Motion Design]], [Motion Elements / Envato (`US$ 29/mês`)], [#strong[Remotion (React/TypeScript)]], [Código Aberto],
    [#strong[Entrega do Projeto]], [Export MP4 achatado], [#strong[Timeline Universal FCPXML]], [Não-destrutivo],
    [#strong[Custo Consolidado]], [#strong[\~ R\$ 590,00 / mês]], [#strong[R\$ 0,00 / mês]], [#strong[Soberania Total]],
  )]
  , kind: table
  )

#horizontalrule

== 1.3 A Lei Áurea: Separação Rígida entre Cognição e Computação
<a-lei-áurea-separação-rígida-entre-cognição-e-computação>
O pilar mais importante da metodologia de Lucas Leite reside na separação estrita de responsabilidades entre inteligência artificial e scripts determinísticos:

#painel-perigo("O ERRO COMUM: FORÇAR O LLM A EDITAR ÁUDIO", [ Modelos de linguagem não possuem ouvidos nem relógio de sistema. Quando um operador instrui um LLM com #emph["corte o áudio entre 01:23.450 e 01:28.120"], o modelo alucina números baseados em probabilidades estatísticas. O resultado é o desalinhamento progressivo da timeline, cortes no meio de sílabas e perda completa da sincronia labial.])

#painel-sucesso("A SOLUÇÃO LEGOLOOP: DECISÃO NO TEXTO, CÁLCULO NO SCRIPT", [ O LLM recebe exclusivamente a transcrição estruturada por sentenças completas. O modelo decide apenas quais frases permanecem na narrativa e quais devem ser expurgadas. A responsabilidade de mapear essas frases para timestamps absolutos, calcular a folga de pré-roll de +40ms e gerar a lista de cortes para o FFmpeg é de scripts utilitários matemáticos locais.])

#horizontalrule

== 1.4 Isolamento Hermético e Governança de Arquivos
<isolamento-hermético-e-governança-de-arquivos>
Para garantir que a manipulação pesada de mídia nunca polua o repositório Git ou o workspace do usuário, o pipeline impõe que todo arquivo intermediário (áudio WAV de 48kHz extraído, arquivos JSON de alinhamento fonético, frames de OCR e renders temporários de animação) resida exclusivamente dentro da pasta temporária e efêmera `scratch/`.

Apenas três artefatos finais nobres são projetados no diretório `dist/`:

+ `dist/timeline.fcpxml` (Linha do tempo aberta editável);
+ `dist/render_final.mp4` (Render de visualização broadcast em 1080p60);
+ `dist/manifest_final.json` (Certificado criptográfico de execução com telemetria).

= Parte II --- Meso: Os 13 Guias Autorais e a Engenharia de Pós-Produção
<parte-ii-meso-os-13-guias-autorais-e-a-engenharia-de-pós-produção>
== 2. Corte Semântico & Extirpação de Notas de Edição
<corte-semântico-extirpação-de-notas-de-edição>
== 2.1 Decisão por Sentença Estruturada
<decisão-por-sentença-estruturada>
O primeiro passo para um corte automatizado de alta fidelidade é transformar o fluxo contínuo de áudio em uma lista estruturada de sentenças fechadas. Diferente dos algoritmos tradicionais de VAD (Voice Activity Detection), que cortam arbitrariamente por limiar de decibéis em pausas de respiração, o nó `n03_decisao_cortes` processa unidades semânticas completas.

Cada sentença gerada pelo Whisper Local é formatada em um contrato de dados JSON estrito:

```json
{
  "sentence_id": "s_014",
  "text": "Agora nós vamos entender como o pipeline DAG organiza as camadas de nós.",
  "start_ms": 42150,
  "end_ms": 46800,
  "confidence": 0.98,
  "words": [
    {"word": "Agora", "start_ms": 42150, "end_ms": 42400},
    {"word": "nós", "start_ms": 42420, "end_ms": 42600},
    {"word": "vamos", "start_ms": 42610, "end_ms": 42900}
  ]
}
```

O modelo cognitivo (LLM) avalia as sentenças com base em coerência narrativa, eliminando:

- Repetições gagas (#emph["então nós… então nós vamos…"]);
- Falsos inícios de raciocínio descartados pelo orador;
- Pausas excessivas entre ideias desconectadas.

#horizontalrule

== 2.2 O Protocolo de Extirpação de Notas de Edição
<o-protocolo-de-extirpação-de-notas-de-edição>
Durante a gravação de aulas e vídeos falados, é extremamente comum o apresentador errar um conceito, parar de olhar para a câmera e dar orientações explícitas ao futuro editor:

#quote(block: true)[
#emph["Opa, falei errado o nome da função. Corta essa parte aqui e recomeça da explicação do JSON… pronto, voltando."]
]

Em esteiras convencionais de IA, esse trecho de bastidores é transcrito normalmente e inserido na timeline final, expondo erros graves em produção.

O nó `n03_decisao_cortes` implementa uma máquina de estados semântica projetada especificamente para identificar esses delimitadores conversacionais de gravação:

#figure(image("img/diagrama_folga_40ms_extirpacao.svg", alt: "Anatomia da Folga Pré-Fala (+40ms) e Extirpação de Notas"),
  caption: [
    Anatomia da Folga Pré-Fala (+40ms) e Extirpação de Notas
  ]
)

=== Padrões Canônicos de Detecção de Notas
<padrões-canônicos-de-detecção-de-notas>
O analisador semântico busca gatilhos textuais e modulações de prosódia como:

+ #strong[Marcadores de Abertura:] #emph["nota de edição:"], #emph["corta isso"], #emph["editor, atenção:"], #emph["errei, recomeçando de…"], #emph["deixa eu refazer essa frase"]\;
+ #strong[Marcadores de Fechamento:] #emph["fim da nota"], #emph["valeu"], #emph["voltando em 3, 2, 1"], #emph["recomeçando agora"]\;
+ #strong[Pausas de Hesitação:] Silêncios anormais superiores a 1.2 segundos acompanhados de mudança abrupta no tom do locutor.

#painel-alerta("REGRA DE EXPURGO TOTAL", [ Quando uma nota de edição é detectada, o intervalo inteiro --- desde o início da hesitação que originou o erro até a retomada definitiva da frase válida --- é marcado com a flag `status: "purged"`. O script de montagem remove completamente esses blocos de áudio e vídeo, sem deixar respirações residuais na timeline.])

#horizontalrule

== 2.3 Exemplo Real: Entrada Bruta vs Saída Higienizada
<exemplo-real-entrada-bruta-vs-saída-higienizada>
A tabela abaixo ilustra o comportamento do filtro de extirpação operando sobre um trecho real gravado no laboratório:

#figure(
  align(center)[#table(
    columns: (12.61%, 38.74%, 21.62%, 27.03%),
    align: (left,left,center,center,),
    table.header([ID], [Sentença Transcrita Bruta], [Classificação Semântica], [Decisão do Algoritmo],),
    table.hline(),
    [`s_041`], [#emph["O primeiro nó do DAG extrai o arquivo de áudio WAV em 48kHz."]], [Conteúdo Didático Válido], [#strong[KEEP (Manter na timeline)]],
    [`s_042`], [#emph["Nota de edição: adicionar um zoom aqui no comando FFmpeg… fim."]], [#strong[Nota de Gravação Falada]], [🛑 #strong[PURGE (Expurgo Obrigatório)]],
    [`s_043`], [#emph["O segundo nó do DAG executa o alinhamento fonético das palavras."]], [Conteúdo Didático Válido], [#strong[KEEP (Manter na timeline)]],
  )]
  , kind: table
  )

O resultado é uma transição suave e imperceptível entre a sentença `s_041` e a sentença `s_043`, economizando tempo precioso de montagem manual e blindando o vídeo contra gafes de estúdio.

== 3. Alinhamento Fonético & A Regra da Folga Pré-Fala (+40ms)
<alinhamento-fonético-a-regra-da-folga-pré-fala-40ms>
== 3.1 A Física do Ataque Fonético na Fala Humana
<a-física-do-ataque-fonético-na-fala-humana>
Quando um ser humano pronuncia uma palavra iniciada por consoantes plosivas (/p/, /b/, /t/, /d/, /k/, /g/) ou fricativas (/s/, /f/, /v/), a cavidade bucal acumula pressão acústica milissegundos antes que o som seja emitido com energia audível pelo microfone.

Modelos automáticos de transcrição e detectores de voz baseados em limiar de volume (VAD) cometem um erro clássico: eles posicionam a marcação de tempo inicial no ponto exato em que a energia sonora ultrapassa a curva de ruído de fundo.

#painel-perigo("O FENÔMENO DA DECAPITAÇÃO FONÉTICA", [ Se a tesoura digital cortar o vídeo rigorosamente no timestamp reportado pelo alinhador fonético tradicional, o início da consoante é decepado. A palavra #emph["Pipeline"] soa como #emph["-ipeline"], a palavra #emph["Código"] soa como #emph["-ódigo"], e o espectador tem a sensação inconsciente de que o áudio está picotado, truncado ou com emendas grosseiras de fita magnética.])

#horizontalrule

== 3.2 A Margem Dourada Inegociável (+40ms)
<a-margem-dourada-inegociável-40ms>
Para eliminar definitivamente esse problema, o método de Lucas Leite institui a regra mandatória de #strong[\+40ms de pré-roll]. Todo e qualquer corte inicial de fala tem seu ponto de entrada antecipado em 40 milissegundos no arquivo de áudio:

$ T_(upright("corte_início")) = T_(upright("palavra_início")) - 0.040 upright(" s") $

Essa pequena janela temporal é suficiente para capturar:

+ O ruído suave da abertura dos lábios ou dentes do locutor;
+ A respiração preparatória imediata que confere naturalidade orgânica à fala;
+ O ataque transitório completo da forma de onda da primeira sílaba.

```json
{
  "word": "Processo",
  "detected_start_ms": 14250,
  "applied_cut_start_ms": 14210,
  "pre_roll_offset_ms": -40
}
```

#horizontalrule

== 3.3 Algoritmo de Resolução de Conflitos e Sobreposições
<algoritmo-de-resolução-de-conflitos-e-sobreposições>
Um desafio matemático surge quando duas sentenças válidas estão muito próximas uma da outra. Se o locutor faz uma pausa curta de apenas 30 milissegundos entre duas frases, subtrair 40 milissegundos da segunda sentença faria com que o corte invadisse a última sílaba da sentença anterior.

O script de cálculo de cortes em Go implementa uma trava de colisão de limites (#emph[boundary clamping]):

```go
func CalculateSafeCutStart(wordStartMs, prevWordEndMs, preRollMs int) int {
    targetStart := wordStartMs - preRollMs
    // Se a folga invadir o término da palavra anterior, limita ao ponto médio
    if targetStart < prevWordEndMs {
        return prevWordEndMs + (wordStartMs - prevWordEndMs) / 2
    }
    return targetStart
}
```

Dessa forma, o algoritmo garante que nenhuma consoante inicial seja mutilada e nenhuma vogal final seja cortada acidentalmente, produzindo uma cadência rítmica perfeita e fluida.

== 4. Câmera Inteligente YuNet & Linha dos Olhos Estabilizada (35%)
<câmera-inteligente-yunet-linha-dos-olhos-estabilizada-35>
== 4.1 A Ilusão do Zoom Dinâmico Convencional
<a-ilusão-do-zoom-dinâmico-convencional>
Em produções modernas para YouTube e redes sociais, o corte de câmera (zoom in em momentos de ênfase e zoom out em introduções) é amplamente utilizado para quebrar a monotonia visual de um plano fixo. No entanto, a maioria dos plugins comerciais aplica reenquadramento cego no centro geométrico do vídeo (50% de altura e largura).

Quando um apresentador se mexe levemente para o lado ou se inclina para frente, o zoom estático corta o topo da cabeça ou posiciona os olhos muito abaixo da linha de leitura do espectador. Pior ainda: se um cartão gráfico animado é inserido na lateral, o apresentador com zoom centralizado colide visualmente com o texto.

#horizontalrule

== 4.2 A Descoberta da Ordem: Por que o Enquadramento é Pós-Motion?
<a-descoberta-da-ordem-por-que-o-enquadramento-é-pós-motion>
Uma das inovações mais elegantes da metodologia de Lucas Leite reside na #strong[ordem estrita de dependência] dos nós:

#painel-alerta("A ORDEM IMPORTA: MOTION PRIMEIRO, CÂMERA DEPOIS", [ O cálculo do enquadramento de câmera só é executado no nó `n07_camera_yunet` #strong[após] a geração do mapa de motion no nó `n05_motion_mapa`. Se a cena exige um cartão gráfico na lateral direita, o algoritmo de câmera sabe antecipadamente que o apresentador deve ser deslocado suavemente para o terço esquerdo da tela, garantindo que o texto e a pessoa coexistam sem sobreposição.])

#figure(image("img/diagrama_yunet_eyelevel_remotion.svg", alt: "Câmera Inteligente YuNet e Motion Procedural sem Colisão"),
  caption: [
    Câmera Inteligente YuNet e Motion Procedural sem Colisão
  ]
)

#horizontalrule

== 4.3 A Regra Áurea dos 35% (Eye-Level Ancorado)
<a-regra-áurea-dos-35-eye-level-ancorado>
Na cinematografia clássica e na transmissão televisiva profissional, o equilíbrio visual de um plano médio (#emph[medium close-up]) exige que a linha dos olhos do apresentador permaneça ancorada em aproximadamente um terço da altura total do quadro a partir do topo (#strong[35%]):

$ Y_(upright("olhos_alvo")) = 0.35 times upright("Altura_Quadro") $

Para um vídeo em alta definição (1080p, altura de 1080 pixels), os olhos devem residir estavelmente na coordenada vertical $Y = 378 upright(" px")$.

=== Rastreamento com OpenCV YuNet
<rastreamento-com-opencv-yunet>
Para alcançar essa precisão sem queimar recursos da GPU com redes neurais pesadas, o pipeline emprega o detector #strong[OpenCV YuNet], uma rede de detecção facial de altíssima eficiência que roda em CPU local a mais de 120 frames por segundo:

```python
### Trecho representativo do cálculo de enquadramento afim
import cv2

def calcular_enquadramento_camera(frame, face_box, landmarks, largura_alvo=1920, altura_alvo=1080):
    olho_esquerdo = landmarks[0]
    olho_direito = landmarks[1]
    olhos_centro_y = (olho_esquerdo[1] + olho_direito[1]) / 2.0
    
    # Altura desejada: 35% do topo
    y_alvo = altura_alvo * 0.35
    deslocamento_vertical = y_alvo - olhos_centro_y
    
    # Aplica suavização com amortecimento (damping) para evitar tremores
    return aplicar_matriz_afim(frame, deslocamento_vertical, escala=1.15)
```

#painel-sucesso("RESULTADO VISUAL ESTÁVEL E CONFORTÁVEL", [ Mesmo que o apresentador levante a cadeira ou incline o corpo durante a gravação, a câmera reajusta a escala e a posição de forma suave e ponderada. O corte de câmera ganha a sensação de ter sido operado por um diretor de fotografia experiente, transmitindo autoridade e profissionalismo.])

== 5. Motion Design Procedural com Remotion
<motion-design-procedural-com-remotion>
== 5.1 O Fim dos Templates Monolíticos do After Effects
<o-fim-dos-templates-monolíticos-do-after-effects>
A edição comercial tradicional depende fortemente de pacotes de motion graphics pré-renderizados comprados em bancos de mídia (Envato Elements, Motion Array). Esses templates apresentam limitações severas:

- Exigem máquinas potentes para rodar instâncias pesadas do Adobe After Effects;
- Qualquer alteração de uma única linha de texto obriga o editor a abrir um arquivo `.aep` complexo e renderizar novamente em tempo real;
- Dificuldade crônica de versionamento: arquivos binários de projeto não aceitam controle de versão via Git (`git diff`).

#horizontalrule

== 5.2 Animação como Componente React
<animação-como-componente-react>
No LegoLoop, os elementos visuais são codificados como componentes declarativos em #strong[Remotion (React/TypeScript)] ou #strong[HyperFrames]. Cada elemento gráfico (cartões laterais, citações em destaque, números de métricas, infográficos de código) é uma função pura que recebe parâmetros de entrada e renderiza quadros com canal alfa transparente (RGBA):

```tsx
import { interpolate, spring, useCurrentFrame, useVideoConfig } from 'remotion';

export const CartaoDidatico: React.FC<{ titulo: string; itens: string[] }> = ({ titulo, itens }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  // Curva de entrada com frenagem suave (Spring Physics)
  const entrada = spring({
    frame,
    fps,
    config: { damping: 14, mass: 0.8, stiffness: 120 }
  });

  const translateX = interpolate(entrada, [0, 1], [-400, 0]);

  return (
    <div style={{ transform: `translateX(${translateX}px)`, opacity: entrada }} className="cartao-glass-3d">
      <h2 className="titulo-anton">{titulo}</h2>
      <ul>
        {itens.map((it, idx) => (
          <li key={idx}>{it}</li>
        ))}
      </ul>
    </div>
  );
};
```

#horizontalrule

== 5.3 O Contrato de Storyboard Gráfico (`motion_storyboard.json`)
<o-contrato-de-storyboard-gráfico-motion_storyboard.json>
O nó `n05_motion_mapa` analisa a narrativa do vídeo e identifica momentos de alto impacto pedagógico onde um gráfico deve reforçar o aprendizado. Ele gera o manifesto do storyboard gráfico:

```json
{
  "cards": [
    {
      "card_id": "card_003",
      "layout": "side_glass_card",
      "start_ms": 28400,
      "end_ms": 34200,
      "title": "3 PILARES DA SOBERANIA",
      "items": [
        "1. Modelos Locais (Whisper & YuNet)",
        "2. Custo Financeiro R$ 0,00",
        "3. Timeline Universal FCPXML"
      ],
      "icon": "shield-check",
      "target_screen_side": "right"
    }
  ]
}
```

#painel-info("VANTAGENS DO RENDER PROCEDURAL", [ 1. #strong[Automação Total:] O renderizador compila sequências de imagens PNG transparentes ou vídeo ProRes 4444 diretamente pela linha de comando em segundos; 2. #strong[Consistência de Marca:] Tipografia, pesos visuais, paletas de cores e raios de borda são governados por tokens de design reutilizáveis em CSS/TypeScript; 3. #strong[Zero Licenças:] O motor roda em qualquer sistema operacional (Linux, macOS, Windows) sem depender de licenças Adobe.])

== 6. Sonorização, Ducking & Sincronia de Efeitos Sonoros (SFX)
<sonorização-ducking-sincronia-de-efeitos-sonoros-sfx>
== 6.1 A Psicologia Acústica do Efeito Sonoro
<a-psicologia-acústica-do-efeito-sonoro>
O som é responsável por mais de 50% da percepção de qualidade de uma produção audiovisual. Em vídeos educativos e tutoriais técnicos, a inserção de efeitos sonoros (#emph[SFX]) cumpre uma função pedagógica: direcionar a atenção ocular do espectador para a informação mais relevante da tela.

No entanto, o maior vício de esteiras automatizadas de edição é disparar o som de transição (#emph[whoosh]) no primeiro frame em que o cartão começa a se mover na borda da tela.

#painel-alerta("O ERRO DO GATILHO PREMATURO", [ Quando o efeito sonoro de transição toca cedo demais, o cérebro do espectador processa o impacto acústico antes que o elemento visual termine sua trajetória. Isso cria uma sensação incômoda de antecipação desordenada e desvaloriza a animação.])

#horizontalrule

== 6.2 O Alinhamento no Frame de Frenagem
<o-alinhamento-no-frame-de-frenagem>
Na metodologia de Lucas Leite, a curva de envelope do som de transição é invertida ou calculada para que o pico absoluto de pressão sonora coincida exatamente com o #strong[frame de parada/frenagem] (#emph[landing frame]) do elemento visual:

#figure(image("img/diagrama_ducking_entrega_hibrida.svg", alt: "Sonorização Sincronizada (Ducking -18dB) e Entrega Híbrida"),
  caption: [
    Sonorização Sincronizada (Ducking -18dB) e Entrega Híbrida
  ]
)

$ T_(upright("pico_sfx")) = T_(upright("parada_card")) plus.minus 1 upright(" frame (0.016s)") $

Quando o cartão atinge sua posição final com uma frenagem suave, o som de impacto (#emph[whoosh + click]) atinge o volume máximo. O olho e o ouvido confirmam o evento no mesmo milissegundo, gerando uma experiência tátil, sólida e agradável.

#horizontalrule

== 6.3 Ducking Dinâmico de -18dB sob a Fala
<ducking-dinâmico-de--18db-sob-a-fala>
A música de fundo (#emph[background music]) cria atmosfera e mantém a retenção, mas pode comprometer a inteligibilidade das frequências médias da voz do locutor (entre 800Hz e 3.5kHz).

O nó `n10_efeitos_sonoros_sfx` aplica compressão sidechain automatizada via filtro FFmpeg:

```bash
### Trecho de comando FFmpeg para Ducking Dinâmico com atenuação de -18dB
ffmpeg -i voz_limpa.wav -i trilha_fundo.mp3 -filter_complex \
  "[1:a]volume=0.8[music]; \
   [music][0:a]sidechaincompress=threshold=0.08:ratio=6:attack=40:release=250[music_ducked]; \
   [0:a][music_ducked]amix=inputs=2:duration=first" saida_mixada.wav
```

=== Parâmetros Recomendados de Calibração:
<parâmetros-recomendados-de-calibração>
- #strong[Atenuação da Trilha sob Fala:] `-18dB` (garante que a voz permaneça nítida e dominante);
- #strong[Ataque (#emph[Attack]):] `40ms` (a música abaixa rapidamente assim que a primeira sílaba soa);
- #strong[Liberação (#emph[Release]):] `250ms` (a música retorna suavemente após o término da frase, sem oscilações bruscas de volume nos pequenos respiros);
- #strong[Ganho nos Silêncios Narrativos:] `0dB` (preenche pausas deliberadas com musicalidade).

== 7. Entrega Híbrida Não-Destrutiva: FCPXML + MP4
<entrega-híbrida-não-destrutiva-fcpxml-mp4>
== 7.1 A Armadilha da Entrega Fechada
<a-armadilha-da-entrega-fechada>
A imensa maioria das ferramentas de inteligência artificial comercializadas para criadores comete um erro fatal de fluxo de trabalho: elas entregam um único arquivo MP4 compactado. Se o criador discordar de um corte, quiser ajustar o volume de uma palavra ou trocar uma cor em um cartão gráfico, ele enfrenta duas opções ruins:

+ Reenviar o vídeo para a nuvem, esperar a fila de processamento e torcer para o modelo acertar;
+ Importar o MP4 achatado em um editor e tentar recortar sobre uma imagem já compactada, onde a trilha musical e a voz estão misturadas e inseparáveis.

#horizontalrule

== 7.2 A Filosofia da Linha do Tempo Aberta
<a-filosofia-da-linha-do-tempo-aberta>
O ecossistema LegoLoop estabelece que o objetivo de uma esteira autônoma não é substituir o editor humano, mas sim #strong[fazer 90% do trabalho pesado e entregar uma timeline aberta e organizada] para que o humano faça os 10% de refinamento estético e polimento final.

O nó `n11_entrega_final_hibrida` bifurca a entrega em dois artefatos complementares:

```
dist/
├── timeline.fcpxml        # Timeline universal com pistas isoladas e handles
├── render_final.mp4       # Render broadcast imediato em 1080p60
└── manifest_final.json    # Metadados e árvore de cuts para auditoria
```

#horizontalrule

== 7.3 Anatomia do Padrão FCPXML Universal
<anatomia-do-padrão-fcpxml-universal>
O formato #strong[Final Cut Pro XML (FCPXML)] foi adotado como o padrão aberto de interoperabilidade do LegoLoop porque ele é lido nativamente e sem perda pelas principais estações de trabalho profissionais do mercado: #strong[DaVinci Resolve], #strong[Final Cut Pro] e #strong[Adobe Premiere Pro].

=== Estrutura das Pistas na Linha do Tempo
<estrutura-das-pistas-na-linha-do-tempo>
Ao abrir o arquivo `timeline.fcpxml`, o editor encontra uma timeline perfeitamente organizada por camadas:

- #strong[Pista de Vídeo 1 (V1):] O vídeo bruto original cortado nas sentenças selecionadas, preservando #strong[sobras (#emph[handles]) de 0.5 segundo] nas extremidades de cada clipe para permitir que o editor faça transições cruzadas (#emph[crossfades]) ou ajustes manuais;
- #strong[Pista de Vídeo 2 (V2):] Os cartões gráficos animados renderizados pelo Remotion com canal alfa (transparência perfeita), posicionados exatamente nos timestamps corretos;
- #strong[Pista de Áudio 1 (A1):] Áudio limpo da voz do apresentador com os cortes fonéticos de +40ms;
- #strong[Pista de Áudio 2 (A2):] Trilha musical com as curvas de automação de volume de -18dB já desenhadas;
- #strong[Pista de Áudio 3 (A3):] Efeitos sonoros (#emph[SFX]) alinhados milimetricamente aos cartões visuais;
- #strong[Marcadores de Capítulo (Markers):] Marcadores coloridos na linha do tempo indicando exatamente onde a IA removeu notas de gravação ou sugeriu títulos de seções.

#painel-sucesso("LIBERDADE CRIATIVA TOTAL", [ O editor não perde tempo cortando silêncios nem caçando notas de gravação. Ele abre a timeline no DaVinci Resolve ou Premiere e foca imediatamente na narrativa, na gradação de cor (#emph[color grading]) e no acabamento fino.])

= Parte III --- Micro: A Topologia DAG, Stop-Gates & Operação
<parte-iii-micro-a-topologia-dag-stop-gates-operação>
== 8. Arquitetura do Pipeline DAG & Stop-Gates Humanos
<arquitetura-do-pipeline-dag-stop-gates-humanos>
== 8.1 A Estrutura Topológica do Grafo Acíclico
<a-estrutura-topológica-do-grafo-acíclico>
A execução da esteira `video-edition` é governada pelo motor matemático do LegoLoop através de um Grafo Acíclico Dirigido (#strong[DAG]). Isso significa que cada tarefa é isolada, possui dependências explícitas declaradas no campo `depends_on`, e o motor calcula a ordenação topológica via Algoritmo de Kahn, garantindo paralelismo onde for possível e ordem estrita onde houver dependência de dados.

#figure(image("img/diagrama_13_guias_dag.svg", alt: "Topologia do Pipeline DAG de Edição (12 Nós e 3 Stop-Gates Humanos)"),
  caption: [
    Topologia do Pipeline DAG de Edição (12 Nós e 3 Stop-Gates Humanos)
  ]
)

#horizontalrule

== 8.2 Tabela Mestra dos 12 Nós do Pipeline
<tabela-mestra-dos-12-nós-do-pipeline>
A tabela a seguir apresenta o contrato de engenharia completo dos 12 nós distribuídos em 11 camadas:

#figure(
  align(center)[#table(
    columns: (11.97%, 18.8%, 11.97%, 13.68%, 13.68%, 29.91%),
    align: (center,left,center,center,center,left,),
    table.header([Camada], [ID do Nó], [Tipo], [Dependências], [Modelo / Runtime], [Contrato de Saída (Artefato)],),
    table.hline(),
    [#strong[C1]], [`n01_iniciar_preparar`], [`script`], [---], [`code-fast`], [`scratch/audio_48k.wav`],
    [#strong[C2]], [`n02_transcricao`], [`task`], [`n01`], [Whisper Local], [`scratch/transcription_words.json`],
    [#strong[C3]], [`n03_decisao_cortes`], [`skill`], [`n02`], [`code-pro`], [`scratch/cuts_manifest.json`],
    [#strong[C4]], [`n04_sg_revisao_corte`], [`stop_gate`], [`n03`], [Humano], [Aprovação Interativa do Corte],
    [#strong[C5]], [`n05_motion_mapa`], [`skill`], [`n04`], [`code-pro`], [`scratch/motion_storyboard.json`],
    [#strong[C6]], [`n06_sg_revisao_motion`], [`stop_gate`], [`n05`], [Humano], [Aprovação Interativa do Motion],
    [#strong[C7]], [`n07_camera_yunet`], [`script`], [`n06`], [OpenCV (CPU)], [`scratch/yunet_tracking.json`],
    [#strong[C7]], [`n08_assets_download`], [`task`], [`n06`], [`code-fast`], [`output/assets_manifest.json`],
    [#strong[C8]], [`n09_motion_render`], [`task`], [`n07`, `n08`], [Remotion (Node)], [`scratch/motion_alpha_bundle/`],
    [#strong[C9]], [`n10_efeitos_sonoros`], [`script`], [`n09`], [FFmpeg Audio], [`scratch/audio_mix_ducked.wav`],
    [#strong[C10]], [`n11_entrega_final`], [`task`], [`n10`], [FFmpeg + FCPXML], [`dist/timeline.fcpxml` + `MP4`],
    [#strong[C11]], [`n12_sg_entrega_final`], [`stop_gate`], [`n11`], [Humano], [Validação e Liberação da Entrega],
  )]
  , kind: table
  )

#horizontalrule

== 8.3 O Papel Crítico dos 3 Stop-Gates Humanos
<o-papel-crítico-dos-3-stop-gates-humanos>
A automação cega é a causa principal da perda de qualidade em produções com IA. O LegoLoop impõe 3 barreiras onde o processo #strong[pausa compulsoriamente] no terminal ou na interface do Studio, exigindo a decisão deliberada do operador:

#painel-alerta("STOP-GATE 01: REVISÃO DE CORTES (n04)", [ O operador abre o arquivo `scratch/cuts_manifest.json` e visualiza o texto com as marcações de corte e notas de gravação expurgadas. Ele pode remover manualmente qualquer frase que o modelo tenha selecionado incorretamente antes de gerar mídia.])

#painel-alerta("STOP-GATE 02: REVISÃO DO STORYBOARD (n06)", [ O operador revisa o manifesto `scratch/motion_storyboard.json`. Ele valida a tipografia, os textos dos cartões e os momentos de entrada, impedindo que a máquina gaste poder computacional renderizando animações com erros de digitação.])

#painel-alerta("STOP-GATE 03: APROVAÇÃO FINAL DA LINHA DO TEMPO (n12)", [ O operador abre a timeline gerada `dist/timeline.fcpxml` no DaVinci Resolve ou Premiere e dá o play no arquivo `dist/render_final.mp4`. Aprovada a sincronia e a fluidez, o processo é arquivado como concluído.])

== 9. Operação Prática, Telemetria & Troubleshooting
<operação-prática-telemetria-troubleshooting>
== 9.1 Guia de Campo: Execução via Linha de Comando (CLI)
<guia-de-campo-execução-via-linha-de-comando-cli>
O operador pode disparar o pipeline diretamente pelo terminal através do executável compilado do LegoLoop:

=== Disparo Padrão com Variável de Entrada
<disparo-padrão-com-variável-de-entrada>
```bash
legoloop run loops/video-edition/pipeline.dag.json --var input_video="entrada/minha_aula.mp4"
```

=== Disparo Customizado com Ajuste Fino de Parâmetros
<disparo-customizado-com-ajuste-fino-de-parâmetros>
```bash
legoloop run loops/video-edition/pipeline.dag.json \
  --var input_video="gravacao_bruta.mov" \
  --var pre_roll_ms=50 \
  --var ducking_db="-20dB" \
  --var eye_level_ratio=0.35
```

#horizontalrule

== 9.2 Matriz de Resolução de Falhas (Troubleshooting)
<matriz-de-resolução-de-falhas-troubleshooting>
A tabela abaixo mapeia os principais incidentes operacionais encontrados em produção e o protocolo imediato de remediação:

#figure(
  align(center)[#table(
    columns: (30.65%, 33.06%, 36.29%),
    align: (left,left,left,),
    table.header([Incidente Operacional], [Causa Raiz Provável], [Ação Corretiva Imediata],),
    table.hline(),
    [#strong[Corte inicial engolindo sílaba plosiva]], [Locutor tem ataque muito rápido ou o áudio tem reverberação], [Aumente o pre-roll de entrada configurando `--var pre_roll_ms=55`.],
    [#strong[Rosto descentralizado durante o zoom]], [Baixa iluminação ou apresentador virou o rosto de perfil], [Edite `scratch/yunet_tracking.json` e ative o fallback de centro geométrico.],
    [#strong[Trilha sonora abafando a voz do orador]], [Compressão sidechain fraca], [Intensifique o ducking passando a variável `--var ducking_db="-22dB"`.],
    [#strong[Erro de render no nó Remotion (n09)]], [Node.js sem memória ou pacote NPM ausente], [Rode `npm install` dentro de `internal/videomotion` e reinicie o nó.],
    [#strong[Timeline FCPXML sem áudio no DaVinci]], [Taxa de amostragem incompatível], [Garanta que o áudio extraído no nó `n01` esteja rigorosamente em 48.000 Hz WAV.],
    [#strong[Nota de gravação vazou no corte final]], [O locutor não usou marcadores verbais claros], [Edite manualmente o arquivo `scratch/cuts_manifest.json` no Stop-Gate 01.],
  )]
  , kind: table
  )

#horizontalrule

== 9.3 Telemetria & Eficiência Financeira
<telemetria-eficiência-financeira>
Ao término da execução, o motor emite o certificado com o balanço consolidado de recursos:

```
=============================================================
             RELATÓRIO DE TELEMETRIA LEGOLOOP
=============================================================
• Duração do Processamento:       4m 12s
• Tokens Cognitivos Consumidos:   12.400 in / 3.200 out
• Custo de Inteligência (API):    $ 0.041 (R$ 0,24)
• Custo de Infraestrutura Local:  R$ 0,00 (Hardware próprio)
• Economia frente ao SaaS:        94.2% de redução de custo
=============================================================
```

#painel-sucesso("EFICIÊNCIA COMPROVADA", [ A esteira `video-edition` alcança 94.2% de economia frente a soluções em nuvem proprietárias, consolidando o retorno financeiro imediato da infraestrutura local sem qualquer cobrança por minuto processado.])

= Conclusão --- A Soberania Audiovisual em Código
<conclusão-a-soberania-audiovisual-em-código>
== 10.1 A Síntese da Jornada Video as Code
<a-síntese-da-jornada-video-as-code>
A consolidação do pipeline `video-edition` demonstra empiricamente que a pós-produção audiovisual de alta complexidade atingiu o seu ponto de inflexão histórica. O modelo tradicional --- baseado em operadores humanos arrastando agulhas em timelines pesadas ou empresas pagando pedágios recorrentes de centenas de dólares mensais a plataformas em nuvem --- tornou-se obsoleto para o processamento de conteúdo falado em escala.

Ao longo desta obra, demonstramos que a transição para o paradigma #strong[Video as Code (VaC)] repousa sobre três pilares de engenharia:

+ #strong[Separação Rígida de Responsabilidades:] O modelo de linguagem (LLM) atua estritamente no domínio semântico textual (decidindo quais ideias permanecem e quais notas de gravação são expurgadas), enquanto algoritmos determinísticos locais calculam amostras de áudio, janelas de RMS e conversão de timecodes com precisão de milissegundos.
+ #strong[Fidelidade Acústica e Geométrica:] A imposição da folga pré-fala (+40ms) extirpa a decapitação fonética de consoantes plosivas, e o rastreamento YuNet com linha dos olhos ancorada a 35% do topo preserva a autoridade e o conforto visual do orador, mesmo sob zooms agressivos.
+ #strong[Entrega Não-Destrutiva com Soberania Zero Custo:] A bifurcação da esteira em timelines abertas (FCPXML) e renderizações de conferência (MP4) garante liberdade artística total na estação de trabalho do editor (DaVinci Resolve / Premiere), sem lock-in de fornecedor e com custo de infraestrutura igual a zero.

#horizontalrule

== 10.2 Tabela Comparativa de Paradigmas
<tabela-comparativa-de-paradigmas>
#figure(
  align(center)[#table(
    columns: (21.65%, 38.14%, 40.21%),
    align: (left,left,left,),
    table.header([Vetor de Avaliação], [Paradigma Legado (SaaS Fechado)], [Paradigma Soberano (LegoLoop VaC)],),
    table.hline(),
    [#strong[Custo Recorrente]], [R\$ 350 a R\$ 800 / mês por usuário], [R\$ 0,00 / mês (Hardware próprio local)],
    [#strong[Privacidade de Dados]], [Upload obrigatório de mídias para a nuvem], [Processamento 100% hermético no terminal local],
    [#strong[Precisão de Sincronia]], [Deriva temporal e cortes engolindo sílabas], [Alinhamento fonético WhisperX com pre-roll de +40ms],
    [#strong[Edição e Ajustes]], [MP4 achatado ou editor proprietário fechado], [Timeline FCPXML multicamada aberta com handles],
    [#strong[Governança de Erros]], [Automação cega sem pontos de pausa], [3 Stop-Gates humanos de validação inegociável],
    [#strong[Motion Design]], [Templates rígidos de bibliotecas comuns], [Animações procedurais em código React (Remotion)],
  )]
  , kind: table
  )

#horizontalrule

== 10.3 O Humano como Diretor no Loop
<o-humano-como-diretor-no-loop>
A automação agêntica no ecossistema LegoLoop nunca teve como objetivo eliminar a sensibilidade criativa humana. Pelo contrário: seu propósito fundamental é livrar o criador de conteúdo do fardo mecânico, repetitivo e degradante de cortar silêncios, caçar gagueiras e alinhar cartões gráficos manualmente.

#painel-teoria("O PAPEL DO DIRETOR SOBERANO", [ Ao operar os 3 Stop-Gates do pipeline (`n04_sg_revisao_cortes`, `n06_sg_revisao_storyboard` e `n12_sg_entrega_final`), o profissional deixa de ser um "operário de mouse" para reassumir o posto que sempre lhe pertenceu: o de #strong[Diretor de Produção]. A máquina processa os gigabytes em silêncio; o humano dá a palavra final sobre o ritmo, a mensagem e o impacto emocional da narrativa.])

#horizontalrule

== 10.4 Próximos Horizontes de Evolução
<próximos-horizontes-de-evolução>
A esteira `video-edition` foi concebida sob uma topologia modular acíclica dirigida (DAG). Essa modularidade permite que o ecossistema incorpore evoluções contínuas sem quebrar os nós homologados:

- #strong[Síntese de Voz e Dublagem Multilingual:] Inserção de nós locais baseados em XTTS/Kokoro para tradução simultânea com preservação de timbre acústico e lipsync procedural;
- #strong[Corte Vertical Adaptativo (Reels/Shorts/TikTok):] Rastreamento dinâmico 9:16 com reenquadramento inteligente baseado em saliência de olhar e reajuste automático dos componentes Remotion;
- #strong[Indexação Vetorial de Acervo:] Alimentação direta de bancos semânticos vetoriais com os metadados higienizados de cada vídeo para busca instantânea de citações passadas.

#painel-sucesso("CERTIFICADO DE SOBERANIA AUDIOVISUAL", [ O pipeline `video-edition` é um patrimônio de software canônico do LegoLoop. Qualquer organização, produtora ou criador independente que adote seus princípios conquista autonomia tecnológica irrevogável, independência contra flutuações cambiais de SaaS e qualidade cinematográfica determinística governada por código aberto.])
