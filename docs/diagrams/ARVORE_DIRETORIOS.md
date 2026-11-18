# Arvore de Diretorios Canonica — legoloop-diagram-design

> Cada diagrama gerado pela skill recebe sua propria arvore de diretorios isolada.
> O slug deriva do titulo em kebab-case.

## Estrutura

```
docs/diagrams/
└── <slug>/                        # ex: arquitetura-fluxo-verbos
    ├── src/                       # fonte principal — HTML/SVG gerado
    │   └── YYYY-MM-DD_<NOME>.html
    ├── assets/                    # fontes locais, imagens, tokens SVG reutilizaveis
    ├── export/                    # saidas derivadas
    │   ├── <slug>.svg             # SVG puro exportado via export_svg.py
    │   └── <slug>.png             # PNG via export_png.py (2x)
    ├── animation/                 # micro-animacoes Remotion
    │   ├── composition.tsx        # composicao Remotion principal
    │   ├── package.json           # deps Remotion isoladas
    │   └── out/                   # MP4/WebM renderizados
    └── anticorpos/                # memoria de erros do diagrama
        └── anticorpos.json        # schema anticorpos.v1
```

## Regras

- O slug e sempre kebab-case do titulo do diagrama.
- O arquivo `src/` recebe prefixo ISO: `YYYY-MM-DD_<NOME>.html`.
- `export/` e gerado apos inspecao impeccable-design aprovada.
- `animation/` e obrigatorio quando Remotion esta disponivel (gate: `npx remotion --version`).
- `anticorpos/anticorpos.json` e criado no inicio, antes de gerar o primeiro arquivo.

## Exemplo de slugs

| Titulo do diagrama                        | Slug                          |
|:------------------------------------------|:------------------------------|
| Arquitetura e Fluxo dos Verbos LegoLoop   | arquitetura-fluxo-verbos      |
| Mapa de Processos do Mega-Pipeline        | mapa-processos-mega-pipeline  |
| DAG do Processo legoloop-fabricar         | dag-legoloop-fabricar         |
| Swimlane de Auditoria Atomica             | swimlane-auditoria-atomica    |
