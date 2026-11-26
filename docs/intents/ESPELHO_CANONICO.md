# Espelho Canônico — legoloop-kraften-deploy

> **Sessão:** `intv_1795723484496401600` (INTERVIEW gravada no cofre de sessões)
> **Data:** 2026-11-26
> **Reconciliação humana:** divergências do motor corrigidas e confirmadas pelo usuário.

## Os 6 Campos Consolidados

### 1. resultado_esperado
Esteira LegoLoop `legoloop-kraften-deploy` funcional com 5 nós determinísticos (`n01_preflight_zip`, `n02_sync_metadata`, `n03_deploy_stream`, `n04_verify_gates`, `n05_receipt`), zero LLM, publicando aplicações de qualquer stack (Go, Node, Python, Rust, PHP, estático) com subdomínio e TLS automáticos na VPS.

### 2. quem_se_beneficia
Desenvolvedores e operadores do ecossistema LegoLoop/Studio Launcher.

### 3. por_que_agora
Migrar o deploy manual/legado do Studio Launcher para um processo padronizado, auditável e 100% determinístico no LegoLoop.

### 4. criterio_de_sucesso
Execução completa dos 5 nós em sequência, gerando build, DNS Cloudflare, Docker Swarm deploy, verificação de gates binários e recibo JSON sem falhas.

### 5. restricoes_inegociaveis
- Zero LLM (nenhum agente intermediário nos nós).
- Implementação da esteira em Go nativo (motor em `internal/kraften/`, comando em `cmd/legoloop` — padrão real do repo, como `internal/diorama/roteiro` + `legoloop diorama-roteiro`).
- Chamadas Docker/Cloudflare CLI e REST API v4.
- Respeito rigoroso à especificação de `PROCESSO_DEPLOY_DETERMINISTICO_LEGOLOOP.md` sem alterações estruturais.

### 6. fora_de_escopo
- Edição do `Dockerfile` da aplicação (nunca mutado pelo Launcher).
- Agentes de IA/LLM intermediários nos nós.
- Stacks sem `Dockerfile`/`docker-stack.yml` (fronteira universal do `.zip`).
- Suporte a redes diferentes de `network_conexao`.

## Notas de Reconciliação
1. **Correção sobre stacks suportadas:** o motor sugeriu excluir Node e Python do escopo; incorreto — a esteira atende QUALQUER stack com `Dockerfile` e `docker-stack.yml` (seção 4 do documento de referência). Go nativo é a implementação da esteira, não restrição das aplicações.
2. **Correção sobre estrutura de diretórios:** o motor sugeriu `pkg/fabrica/`; este caminho não existe no repo. Padrão canônico real: motor em `internal/kraften/`, comando CLI em `cmd/legoloop`.
