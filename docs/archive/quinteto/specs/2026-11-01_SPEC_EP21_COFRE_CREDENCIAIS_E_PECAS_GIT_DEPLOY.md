# 📋 SPEC_EP21: Cofre de Credenciais (`env-vault`) & Peças de Operações Git e Deploy

> **Status:** Proposta para Execução  
> **Data:** 2026-11-01  
> **Épico Alvo:** ÉPICO 21 (Cofre de Credenciais, Automação Git & Esteira de Deploy)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

O **ÉPICO 21** introduz as camadas fundamentais de segurança de segredos e fechamento do ciclo de entrega contínua (*Continuous Delivery*) no LegoLoop:

1. **Cofre Criptografado de Credenciais (`internal/vault`):**
   * Motor criptográfico nativo com cifra simétrica AES-256-GCM.
   * Derivação de chave mestra resiliente (via variável de ambiente `LEGOLOOP_VAULT_KEY` ou semente determinística da máquina do desenvolvedor).
   * Persistência em disco em arquivo criptografado (`.legoloop/vault.enc` ou configurável).
   * Sanitização e mascaramento estrito de valores confidenciais (`***MASKED***`) para evitar vazamento em logs, telemetria e saída de terminais.
   * Resolução dinâmica de variáveis de ambiente com injeção segura nos executores de nós sem exposição no grafo DAG público.

2. **Peça Ontológica `env-vault` (`internal/graph/ontology.go`):**
   * Bloco registrado com `Kind: "security_vault"`.
   * Resolve e valida a presença das chaves de API e variáveis exigidas pela esteira, gerando o manifesto de conformidade `VAULT_ENV.json`.

3. **Peças Canônicas de Versionamento Git (`internal/graph/ontology.go` & `internal/runner`):**
   * `git-commit` (`Kind: "vcs_git_commit"`): Staging de artefatos e criação de commit convencional atômico, gerando `GIT_COMMIT.json`.
   * `git-merge` (`Kind: "vcs_git_merge"`): Integração de branches locais/remotas com verificação de conflitos, gerando `GIT_MERGE.json`.
   * `git-tag-release` (`Kind: "vcs_git_tag"`): Criação de tags semânticas de versão (`vX.Y.Z`) e changelog resumido em `GIT_RELEASE.json`.

4. **Peças Canônicas de Deploy e Publicação (`internal/graph/ontology.go` & `internal/runner`):**
   * `deploy-docker` (`Kind: "deploy_docker"`): Validação de `Dockerfile`, build e verificação de imagem conteinerizada, gerando `DOCKER_BUILD.json`.
   * `deploy-vps` (`Kind: "deploy_vps"`): Despacho de artefatos para VPS Linux remota via SSH/SCP determinístico com chave do Vault, gerando `DEPLOY_VPS.json`.
   * `deploy-vercel` (`Kind: "deploy_vercel"`): Disparo de deploy estático/serverless em ambientes de borda (Edge), gerando `DEPLOY_VERCEL.json`.

5. **Exposição no Studio Canvas (`internal/studio/pieces.go`):**
   * Cores semânticas ontológicas:
     * Vault/Segurança: Dourado Âmbar (`#F59E0B`).
     * Git / VCS: Azul Royal (`#2563EB`).
     * Deploy / Entrega: Verde Esmeralda (`#10B981`).
   * Biblioteca de ícones SVG canônica integrada (Zero Emojis).

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

* [ ] **CA-21.1 (Vault Criptográfico AES-256-GCM):** O pacote `internal/vault` criptografa e decriptografa pares chave-valor com autenticação GCM sem corrupção e persiste no disco de forma cifrada.
* [ ] **CA-21.2 (Mascaramento de Segredos):** O método `MaskSecrets` substitui qualquer ocorrência de valores sensíveis conhecidos por `***MASKED***` em strings de log e saídas.
* [ ] **CA-21.3 (Peça Ontológica `env-vault`):** O bloco `env-vault` está registrado em `SecurityBlocks()` / `ControlBlocks()`, consumindo qualquer artefato e produzindo `Code` / `Audit` com artefato `VAULT_ENV.json`.
* [ ] **CA-21.4 (Peças Git Canônicas):** `git-commit`, `git-merge` e `git-tag-release` estão registradas no catálogo ontológico com portas e validações de gates declaradas.
* [ ] **CA-21.5 (Peças Deploy Canônicas):** `deploy-docker`, `deploy-vps` e `deploy-vercel` estão registradas com portas `Consumes: []ArtifactType{Code, Audit}` e `Produces: Audit`.
* [ ] **CA-21.6 (Studio Palette & Categorização):** As novas peças aparecem categorizadas e com metadados corretos nas rotas da API do Studio (`/api/pieces` e `/api/blocks`).
* [ ] **CA-21.7 (Zero Regressões & 100% Testes Verdes):** Todos os testes unitários (`go test ./...`) passam limpos e o binário é recompilado em `bin/legoloop.exe`.
