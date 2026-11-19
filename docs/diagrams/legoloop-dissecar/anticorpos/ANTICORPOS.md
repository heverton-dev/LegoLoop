---
schema: anticorpos.v1
processo: legoloop-dissecar
total_anticorpos: 7
globais: 0
dominio: 7
regras_inviolaveis:
  - pacote: "legoloop-diagram-design/node-specific-description"
    antidoto: "Cada nó DEVE obrigatoriamente exibir a descrição do que ele faz. O gerador extrai a descrição da Skill (seção de Proc..."
  - pacote: "legoloop-diagram-design/io-clean-names-no-ellipsis"
    antidoto: "Banimento absoluto de reticências (...) em In e Out. Insumos e saídas que são arquivos devem exibir estritamente o no..."
  - pacote: "legoloop-diagram-design/ban-meta-generator-sections"
    antidoto: "Proibição absoluta de meta-seções sobre diretrizes do próprio gerador de diagramas. Todo o conteúdo adicional do docu..."
  - pacote: "legoloop-diagram-design/table-type-badge-compact"
    antidoto: "Criar classe .table-type-badge com font-size: 10px, font-family mono, font-weight 700, padding 3px 8px e white-space:..."
  - pacote: "legoloop-diagram-design/lineage-badge-padding"
    antidoto: "Aplicar classe .lineage-badge com padding 4px 10px, border-radius 4px, font-size 11px, display inline-flex, borda sut..."
  - pacote: "legoloop-diagram-design/invariant-card-neutral-border"
    antidoto: "Eliminar qualquer borda lateral colorida em .invariant-card, mantendo contorno perimetral único e neutro 1px solid va..."
  - pacote: "legoloop-diagram-design/header-iso-timestamp"
    antidoto: "Substituir o campo 'Design System' por 'Gerado em: <DATA_HORA_ISO>' com precisão até segundos (YYYY-MM-DDTHH:MM:SS) f..."
---

# 🛡️ Banco de Anticorpos — legoloop-dissecar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (7)

### `legoloop-diagram-design/node-specific-description`
- **Antídoto:** Cada nó DEVE obrigatoriamente exibir a descrição do que ele faz. O gerador extrai a descrição da Skill (seção de Procedimento Passo a Passo), da spec ou sintetiza a ação funcional pelo verbo/objeto do nó. Proibição absoluta de descrição genérica idêntica compartilhada entre nós.
- **Erro Evitado:** Descrição idêntica e genérica replicada em todos os cards de nós do diagrama ('Execução no processo <slug> com isolamento em contexto estéril').
- **Causa Raiz:** O gerador utilizava template estático hardcoded no campo .node-desc em vez de extrair a descrição funcional específica de cada nó a partir de SKILL.md, spec_ref ou do DAG.
- **Assinatura:** `sha256:8c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d`

### `legoloop-diagram-design/io-clean-names-no-ellipsis`
- **Antídoto:** Banimento absoluto de reticências (...) em In e Out. Insumos e saídas que são arquivos devem exibir estritamente o nome limpo do arquivo (ex.: target_resolved.json, ast_metrics.json), e variáveis devem exibir a referência direta limpa (ex.: var.alvo), sem wrappers de dicionário Python/JSON.
- **Erro Evitado:** Uso de reticências (...) e exibição de objetos brutos de dicionário/linguagem em campos In e Out dos cards (ex: In: {'target': "file('as...").
- **Causa Raiz:** Truncamento artificial via val[:max_len-3] + '...' e renderização direta do dicionário bruto de inputs do DAG sem parsing de arquivos.
- **Assinatura:** `sha256:9d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e`

### `legoloop-diagram-design/ban-meta-generator-sections`
- **Antídoto:** Proibição absoluta de meta-seções sobre diretrizes do próprio gerador de diagramas. Todo o conteúdo adicional do documento DEVE descrever exclusivamente dimensões funcionais do processo alvo: Roteamento Cognitivo de Modelos, Linhagem de Dados e Invariantes/Regras Negativas.
- **Erro Evitado:** Inclusão de seção meta-interna do gerador ('04 — Governança de Anticorpos e Paridade com o LegoPress' com regras como Fundo Descansado, Zero Serifa, etc.) no documento do processo alvo.
- **Causa Raiz:** O gerador inseria um bloco HTML estático descrevendo regras internas da ferramenta legoloop-diagram-design em vez de documentar informações técnicas exclusivas do processo alvo.
- **Assinatura:** `sha256:1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b`

### `legoloop-diagram-design/table-type-badge-compact`
- **Antídoto:** Criar classe .table-type-badge com font-size: 10px, font-family mono, font-weight 700, padding 3px 8px e white-space: nowrap, garantindo exibição em linha única compacta.
- **Erro Evitado:** Badge de tipo de nó (UNIT-NÓ) na tabela (Seção 03) renderizado com fonte grande e quebrando em duas linhas.
- **Causa Raiz:** Ausência de classe tipográfica dedicada compacta e white-space: nowrap para a célula de tipo na tabela.
- **Assinatura:** `sha256:2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c`

### `legoloop-diagram-design/lineage-badge-padding`
- **Antídoto:** Aplicar classe .lineage-badge com padding 4px 10px, border-radius 4px, font-size 11px, display inline-flex, borda sutil 1px solid {cor}33 e fundo translúcido {cor}18.
- **Erro Evitado:** Badges de linhagem (ENTRADA, PIPELINE, ENTREGA) na Seção 05 renderizados sem padding, parecendo texto solto e espremido.
- **Causa Raiz:** Spans de badge inline sem padding explícito, border-radius e display inline-flex.
- **Assinatura:** `sha256:3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d`

### `legoloop-diagram-design/invariant-card-neutral-border`
- **Antídoto:** Eliminar qualquer borda lateral colorida em .invariant-card, mantendo contorno perimetral único e neutro 1px solid var(--rule).
- **Erro Evitado:** Cards de invariantes (Seção 06) com bordas laterais coloridas (border-left), poluindo a linguagem visual neutra.
- **Causa Raiz:** Uso residual de border-left: 3px solid var(--cor-accent) na classe .invariant-card.
- **Assinatura:** `sha256:4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e`

### `legoloop-diagram-design/header-iso-timestamp`
- **Antídoto:** Substituir o campo 'Design System' por 'Gerado em: <DATA_HORA_ISO>' com precisão até segundos (YYYY-MM-DDTHH:MM:SS) formatado dinamicamente em tempo de compilação.
- **Erro Evitado:** Cabeçalho exibindo metadado genérico de design system ('LegoPress Editorial') em vez da data e hora real de geração.
- **Causa Raiz:** String fixa inserida no template HTML em vez de timestamp dinâmico de build.
- **Assinatura:** `sha256:5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

