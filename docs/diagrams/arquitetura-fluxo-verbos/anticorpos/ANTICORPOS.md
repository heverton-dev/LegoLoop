---
schema: anticorpos.v1
processo: legoloop-diagram-design
total_anticorpos: 11
globais: 0
dominio: 11
regras_inviolaveis:
  - pacote: "legoloop-diagram-design/svg-layout"
    antidoto: "Usar grid CSS externo ao SVG para distribuir nós da Camada 2 em duas linhas explícitas (linha A: listar, matriz, fabr..."
  - pacote: "legoloop-diagram-design/text-overflow"
    antidoto: "Quebrar o sublabel em duas linhas SVG usando dois elementos <text> com dy offset, ou aumentar o width do nó para >= 3..."
  - pacote: "legoloop-diagram-design/animation"
    antidoto: "Adicionar nó obrigatorio 'animacao' ao DAG da skill com gate: arquivo animation/composition.tsx deve existir. Verific..."
  - pacote: "legoloop-diagram-design/output-path"
    antidoto: "Sempre criar /docs/diagrams/<slug>/{src,assets,export,animation,anticorpos}/ antes de gerar qualquer arquivo. O slug ..."
  - pacote: "legoloop-diagram-design/line-card-collision"
    antidoto: "Roteamento Perimetral Obrigatório: arestas de longo curso e retornos devem correr exclusivamente pelas calhas externa..."
  - pacote: "legoloop-diagram-design/text-box-overflow"
    antidoto: "Dimensionamento Responsivo de Cards: W_min >= X_pad + (N_chars * Char_W) + R_pad. Card fechar alargado para 520px com..."
  - pacote: "legoloop-diagram-design/native-browser-animation"
    antidoto: "Fusão Didática Impeccable: o documento preserva a grade explicativa dos 13 verbos e a jornada passo a passo de 5 etap..."
  - pacote: "legoloop-diagram-design/anti-emoji-and-grid-3x2"
    antidoto: "Proibição estrita de emojis: substituir por micro-ícones SVG vetoriais inline. Adoção de grid ortogonal 3x2: coluna e..."
  - pacote: "legoloop-diagram-design/badge-alignment-and-sans-rule"
    antidoto: "Regra Anti-Serif Estrita: proibido Instrument Serif e Georgia, usar apenas Geist Sans e Geist Mono. Badges Horizontai..."
  - pacote: "legoloop-diagram-design/geometric-shapes-legend"
    antidoto: "Legenda Dupla Estruturada: bloco explícito com as 5 formas geométricas dos nós (hexágono, losango, círculo, triângulo..."
  - pacote: "legoloop-diagram-design/subtle-legend-micro-svg"
    antidoto: "Isolamento Estrito de Canvas e Micro-Vetores: o SVG principal do diagrama recebe obrigatoriamente a classe .diagram-c..."
---

# 🛡️ Banco de Anticorpos — legoloop-diagram-design

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (11)

### `legoloop-diagram-design/svg-layout`
- **Antídoto:** Usar grid CSS externo ao SVG para distribuir nós da Camada 2 em duas linhas explícitas (linha A: listar, matriz, fabricar, refinar, compor / linha B: fechar, remover). Dentro do SVG, garantir que o viewBox seja calculado dinamicamente com base no número de nós + padding, não em valores hardcoded.
- **Erro Evitado:** Layout SVG da Camada 2 empilha nos em duas linhas sem alinhamento ortogonal. Os verbos 'fechar' e 'remover' ficam deslocados para uma segunda linha quebrada, cortando labels como 'ROUTE — DAG' no nó 'fabricar'.
- **Causa Raiz:** O viewBox e os valores de translate() usam coordenadas absolutas fixas sem calcular a altura real de cada nó. O nó 'fabricar' (height=68) é mais alto que os outros (height=52), causando desalinhamento vertical da grade.
- **Assinatura:** `sha256:a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2`

### `legoloop-diagram-design/text-overflow`
- **Antídoto:** Quebrar o sublabel em duas linhas SVG usando dois elementos <text> com dy offset, ou aumentar o width do nó para >= 340px. Sempre executar scripts/verify_geometry.py antes de salvar o artefato final.
- **Erro Evitado:** Texto de sublabel truncado no nó 'fabricar' (Camada 2): 'INTERVIEW — INTENT — SPEC — PLAN — ROUTE — DAG' ultrapassa a borda direita do nó e fica cortado visivelmente.
- **Causa Raiz:** O width do nó 'fabricar' (260px) é insuficiente para a string de 50 chars em font-size 8px. A regra de orçamento de largura (verify_geometry.py) não foi executada antes da entrega.
- **Assinatura:** `sha256:b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3`

### `legoloop-diagram-design/animation`
- **Antídoto:** Adicionar nó obrigatorio 'animacao' ao DAG da skill com gate: arquivo animation/composition.tsx deve existir. Verificar disponibilidade de Remotion via 'npx remotion --version' antes de gerar.
- **Erro Evitado:** Ausencia de micro-animacoes. Diagramas entregues sem transicoes CSS/Remotion mesmo quando o processo exige inspecao impeccable-design completa.
- **Causa Raiz:** O pipeline de geracao nao incluia a etapa de animacao como nó obrigatorio. A skill nao declarava dependencia de legoloop-impeccable-design nem de Remotion.
- **Assinatura:** `sha256:c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4`

### `legoloop-diagram-design/output-path`
- **Antídoto:** Sempre criar /docs/diagrams/<slug>/{src,assets,export,animation,anticorpos}/ antes de gerar qualquer arquivo. O slug deriva do titulo do diagrama em kebab-case.
- **Erro Evitado:** Artefato gerado salvo em docs/manuals/ e docs/diagrams/ raiz, sem estrutura de slug dedicada /docs/diagrams/<slug>/.
- **Causa Raiz:** A regra de diretorio canonico (guardrail 8) foi adicionada depois da geracao. Nao havia slug tree definida no processo.
- **Assinatura:** `sha256:d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5`

### `legoloop-diagram-design/line-card-collision`
- **Antídoto:** Roteamento Perimetral Obrigatório: arestas de longo curso e retornos devem correr exclusivamente pelas calhas externas (margins/gutters x<=50 ou x>=1000). Proibido cruzar bounding boxes de nós intermediários. Adicionar mascara de papel opaco <rect fill='var(--paper)'/> em nós translúcidos para prevenir bleeding-through.
- **Erro Evitado:** Linhas e arestas atravessando cards intermediarios (ex: aresta vertical listar->remover cortando refinar; aresta fechar->canonize cortando compor e upgrade).
- **Causa Raiz:** Roteamento de arestas de longo curso usando coordenadas absolutas que intersectam bounding boxes de nós intermediários em vez de usar as calhas perimetrais desobstruídas.
- **Assinatura:** `sha256:e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6`

### `legoloop-diagram-design/text-box-overflow`
- **Antídoto:** Dimensionamento Responsivo de Cards: W_min >= X_pad + (N_chars * Char_W) + R_pad. Card fechar alargado para 520px com centralização. Todo rótulo de aresta deve ter elemento <rect class='mask'> opaco com padding de 6-10px afastado do traço.
- **Erro Evitado:** Texto transbordando borda do card fechar (9 etapas: Manual - Livro - Skill...) e rótulos sobrepondo linhas.
- **Causa Raiz:** Card com largura de 320px insuficiente para texto de 85 caracteres mono (~430px) sem quebra de linha. Rótulos de arestas sem placa de máscara retangular opaca de proteção.
- **Assinatura:** `sha256:f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1`

### `legoloop-diagram-design/native-browser-animation`
- **Antídoto:** Fusão Didática Impeccable: o documento preserva a grade explicativa dos 13 verbos e a jornada passo a passo de 5 etapas. Animações operam nativamente no navegador via CSS com keyframes (stroke-dashoffset contínuo para fluxo de elétrons e pulso suave no nó fabricar).
- **Erro Evitado:** Animações invisíveis no navegador e perda da didática tutorial ao isolar o SVG puro.
- **Causa Raiz:** Dependência de renderizador Remotion externo em vez de CSS nativo. Supressão do contexto em prosa e da tabela de camadas que existiam na versão anterior.
- **Assinatura:** `sha256:0a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b`

### `legoloop-diagram-design/anti-emoji-and-grid-3x2`
- **Antídoto:** Proibição estrita de emojis: substituir por micro-ícones SVG vetoriais inline. Adoção de grid ortogonal 3x2: coluna esq (refinar/remover), centro (validate/exec), dir (compor/upgrade), deixando calhas centrais de 170px para rotas verticais sem sobreposição.
- **Erro Evitado:** Presença de emojis decorativos no cabeçalho e sobreposição de rota entre exec->refinar e nó remover.
- **Causa Raiz:** Violação da regra de emojis na tipografia e posicionamento assimétrico de nós na camada de execução sem colunas ortogonais fixas.
- **Assinatura:** `sha256:1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c`

### `legoloop-diagram-design/badge-alignment-and-sans-rule`
- **Antídoto:** Regra Anti-Serif Estrita: proibido Instrument Serif e Georgia, usar apenas Geist Sans e Geist Mono. Badges Horizontais Centrados: posicionar badges nos segmentos retos livres sem rotação de 90 graus, com alinhamento matemático (text-anchor=middle). Alargamento do nó fechar para 540px garantindo folga interna de 40px à direita.
- **Erro Evitado:** Badges com texto vertical desalinhado fora do retângulo (A e B); texto de 'fechar' transbordando o card à direita (C); uso de fontes serifadas proibidas.
- **Causa Raiz:** Badges rotacionados em 90 graus sem dimensionamento proporcional do rect e texto fora do eixo. Largura do nó 'fechar' (440px) insuficiente para 82 caracteres mono. Uso de Instrument Serif herdado do template.
- **Assinatura:** `sha256:2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d`

### `legoloop-diagram-design/geometric-shapes-legend`
- **Antídoto:** Legenda Dupla Estruturada: bloco explícito com as 5 formas geométricas dos nós (hexágono, losango, círculo, triângulo, retângulo) e bloco com os 4 tipos de aresta. Nó fechar ampliado para 560px com subtítulo compacto usando separadores de ponto central (·), garantindo mais de 120px de respiro interno.
- **Erro Evitado:** Ausência da legenda semântica de formas geométricas dos nós e persistência de risco de transbordamento no card fechar.
- **Causa Raiz:** Legenda limitada apenas a arestas de conexão. String de 9 etapas com travessões longos gerando largura superior a 500px em renderers de fonte no Windows.
- **Assinatura:** `sha256:3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e`

### `legoloop-diagram-design/subtle-legend-micro-svg`
- **Antídoto:** Isolamento Estrito de Canvas e Micro-Vetores: o SVG principal do diagrama recebe obrigatoriamente a classe .diagram-canvas vinculada a .diagram-card > svg.diagram-canvas { width: 100% }. Todo SVG secundário de legenda recebe estilos inline rígidos (style='width:13px;height:13px;max-width:13px;max-height:13px;flex-shrink:0;') e classe dedicada (.shape-icon ou .arrow-icon), impedindo qualquer herança ou expansão indesejada.
- **Erro Evitado:** Ícones da legenda renderizados gigantescos (esticamento para 100% ou sem restrição no flexbox).
- **Causa Raiz:** O seletor CSS .diagram-card svg aplicava width:100% globalmente a qualquer SVG filho, inclusive aos micro-ícones de 16px da legenda, e os elementos SVG da legenda não possuíam estilos inline estritos de contenção.
- **Assinatura:** `sha256:4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

