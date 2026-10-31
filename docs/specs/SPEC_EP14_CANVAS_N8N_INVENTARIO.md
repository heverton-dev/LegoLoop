# SPEC-EP14: Experiência N8N, Refatoração Ergonômica da Sidebar e Refinamento do Acervo

> **Status:** APROVADA  
> **Data:** 2026-10-31  
> **Escopo:** Varredura recursiva de mono-repos/sub-skills, indexação automática do 9Router na VPS, refatoração da barra lateral esquerda (320px + seções colapsáveis), Canvas estilo N8N (pan 2D infinito, zoom infinito Ctrl+Wheel/Ctrl+/-/Ctrl+0, desconexão de arestas e inserção intermediária), agrupamento semântico de duplicatas no inventário, paginação completa, modal rico com 4 sub-abas e tradução sob demanda em cache local, e sanitização obrigatória de loops em kebab-case.

---

## 1. Critérios de Aceite Binários

### 1.1 Varredura Hiper-Profunda & 9Router VPS
- [ ] **CRIT-SCAN-RECURSIVE-1:** `findSkillFiles` continua a descida recursiva em subdiretórios mesmo após encontrar um `SKILL.md` (suporte a mono-repos como `aidd-osmani/upstream/skills/**`), sem ultrapassar `maxSkillDepth`.
- [ ] **CRIT-SCAN-ROOTS-2:** O scanner inclui todas as raízes do ecossistema: `~/.gemini/antigravity-cli/skills`, `~/.gemini/antigravity/skills`, `~/.gemini/config/skills`, `~/.cline/skills`, `~/.factory/skills`, `~/.config/opencode/skills`.
- [ ] **CRIT-9ROUTER-ENV-1:** O pacote `provider` lê arquivos `.env` conhecidos (`Desktop/9router/.env`, `~/.aidd/.env`, `~/.env`, `.env`) e detecta `NINEROUTER_URL` e `NINEROUTER_KEY` caso não estejam setados no processo.
- [ ] **CRIT-9ROUTER-VPS-2:** Com a URL da VPS `https://9router.vpsconexao.org`, o scan lista com sucesso os modelos disponíveis (`code-pro`, `code-fast`, `code-free`, etc.) no inventário de modelos.

### 1.2 Governança de Loops (Kebab-case Estrito)
- [ ] **CRIT-LOOP-KEBAB-1:** Qualquer salvamento ou criação de loop normaliza compulsoriamente o ID e o Nome para `kebab-case` determinístico (letras minúsculas, números e hifens; espaços e caracteres especiais convertidos em `-`).
- [ ] **CRIT-LOOP-KEBAB-2:** Nomes com acentos ou pontuações (ex: "Criar Código & Testes! V1") são higienizados para "criar-codigo-testes-v1".

### 1.3 Sidebar Ergonômica & Responsiva
- [ ] **CRIT-SIDEBAR-UI-1:** A barra lateral esquerda possui largura respirável (320px) com espaçamento e tipografia legíveis.
- [ ] **CRIT-SIDEBAR-ACCORDION-2:** As seções `Peças Disponíveis` e `Modelos Vivos` possuem botões de colapso/expansão independentes.
- [ ] **CRIT-SIDEBAR-TOGGLE-3:** Há um botão para recolher/ocultar a barra lateral inteira, maximizando a área do canvas.
- [ ] **CRIT-SIDEBAR-PILLS-4:** As pílulas de filtro de categorias de peças são dispostas em grid/flow organizado, sem sobreposição nem corte de texto.

### 1.4 Canvas Estilo N8N (Navegação & Conexão)
- [ ] **CRIT-CANVAS-PAN-1:** O usuário pode navegar infinitamente no canvas 2D arrastando o plano de fundo (ou botão do meio do mouse / espaço + drag).
- [ ] **CRIT-CANVAS-ZOOM-2:** Zoom contínuo com `Ctrl + Wheel` do mouse centrado no ponteiro, e atalhos de teclado `Ctrl +` e `Ctrl -` (com botões flutuantes na tela).
- [ ] **CRIT-CANVAS-FIT-3:** Tecla `Ctrl + 0` e botão de enquadramento ("Ajustar à Tela") recalculam o bounding box de todos os nós e centralizam com padding.
- [ ] **CRIT-CANVAS-DISCONNECT-4:** Ao passar o mouse sobre uma aresta de conexão (ou selecioná-la), um botão visual de desconexão ('x') é exibido e a tecla `Delete`/`Backspace` remove a aresta.
- [ ] **CRIT-CANVAS-INSERT-5:** Ao soltar uma peça sobre uma conexão existente (ou via menu de contexto da aresta), o sistema reconecta intermediariamente: Nó A ➔ Nova Peça ➔ Nó B.

### 1.5 Inventário: Agrupamento Semântico e Paginação
- [ ] **CRIT-INV-GROUP-1:** Itens com mesmo SHA-256 são agrupados por padrão em 1 card consolidado, exibindo tags de todas as origens onde a peça existe e a quantidade de cópias.
- [ ] **CRIT-INV-TOGGLE-2:** Um botão seletor permite alternar entre "Agrupar Cópias" (padrão) e "Exibir Todos os Arquivos".
- [ ] **CRIT-INV-PAGE-3:** A lista de itens possui paginação com seletor (25, 50, 100 itens) e controles de navegação (`« Primeira`, `‹ Anterior`, `X de Y`, `Próxima ›`, `Última »`).

### 1.6 Modal de Detalhes da Peça (4 Sub-abas) e Tradução
- [ ] **CRIT-MODAL-TABS-1:** O modal flutuante exibe 4 sub-abas: `Visão Geral`, `Peças Compatíveis`, `Ocorrências no SO` e `Código / Fonte`.
- [ ] **CRIT-MODAL-COMPAT-2:** A sub-aba `Peças Compatíveis` lista conexões baseadas no contrato ontológico (`inputs`, `outputs`, hooks correlatos e ferramentas de MCP).
- [ ] **CRIT-MODAL-OCCUR-3:** A sub-aba `Ocorrências no SO` lista os caminhos físicos completos de todas as cópias da peça com seus hashes.
- [ ] **CRIT-MODAL-TRANSLATE-4:** Na sub-aba `Código / Fonte`, o botão "Traduzir para PT-BR" invoca o endpoint de tradução sob demanda; o arquivo original permanece **100% inalterado**, e a tradução é salva exclusivamente em cache `~/.aidd/cache/translations/{sha256}.json`.

---

## 2. Critérios de Segurança (SEC)
- [ ] **SEC-EP14-1 (Segredos Preservados):** Chaves de API do 9Router nunca são expostas em payloads de log ou endpoints públicos sem máscara.
- [ ] **SEC-EP14-2 (Anti-Traversal em Tradução e Detalhe):** Endpoints de conteúdo e tradução rejeitam qualquer tentativa de path traversal fora das raízes de ativos autorizadas.
