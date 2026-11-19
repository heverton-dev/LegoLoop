---
schema: anticorpos.v1
processo: legoloop-dissecar
total_anticorpos: 3
globais: 0
dominio: 3
regras_inviolaveis:
  - pacote: "legoloop-diagram-design/node-specific-description"
    antidoto: "Cada nó DEVE obrigatoriamente exibir a descrição do que ele faz. O gerador extrai a descrição da Skill (seção de Proc..."
  - pacote: "legoloop-diagram-design/io-clean-names-no-ellipsis"
    antidoto: "Banimento absoluto de reticências (...) em In e Out. Insumos e saídas que são arquivos devem exibir estritamente o no..."
  - pacote: "legoloop-diagram-design/ban-meta-generator-sections"
    antidoto: "Proibição absoluta de meta-seções sobre diretrizes do próprio gerador de diagramas. Todo o conteúdo adicional do docu..."
---

# 🛡️ Banco de Anticorpos — legoloop-dissecar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (3)

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

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

