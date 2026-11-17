# 📋 SPEC_EP24: Importador Global do GitHub com Canonização HDSM & Fork/Versionamento

> **Status:** Proposta para Execução  
> **Data:** 2026-11-01  
> **Épico Alvo:** ÉPICO 24 (Importador GitHub, Canonização Imediata e Fork Copy-on-Write)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

O **ÉPICO 24** consolida a soberania do Sistema Operacional sobre ativos externos e esteiras compostas, garantindo que qualquer ativo importado seja desduplicado, versionado e canonizado:

1. **Importador Global do GitHub (`internal/canon/importer.go`):**
   * Aceita URLs do GitHub (repositórios completos ou subpastas de skills).
   * Validação estrita de segurança: validação de domínio (`github.com`), sanitização de nomes contra *path traversal* (`../`) e controle de timeout de download.
   * **Instalação Global Obrigatória:** Gravação prioritária no diretório global do SO: `~/.gemini/config/skills/<skill-name>` (com fallback para `~/.aidd/canonical/skills/<skill-name>`).
   * **Canonização Automática:** Varredura imediata, extração de frontmatter ontológico, cálculo do hash SHA-256 e integração ao catálogo de peças.

2. **Endpoint no Studio API (`POST /api/assets/import`):**
   * Recebe JSON: `{"url": "https://github.com/...", "name": "nome-opcional"}`.
   * Executa a importação, canoniza e atualiza o inventário com `s.Assets.Refresh(ctx)`.
   * Retorna os metadados do ativo importado (`name`, `sha256`, `target_dir`, `files_count`).

3. **Comando CLI `legoloop import` (`internal/cli`):**
   * Permite importar diretamente via terminal: `legoloop import <url> [--name <custom-name>]`.

4. **Fork e Versionamento de Loops (Copy-on-Write):**
   * Endpoint `POST /api/loops/{id}/fork`: clona o sub-loop incrementando a versão ou adicionando sufixo de fork sem tocar no original.
   * Proteção de integridade dos templates embutidos (`mega-pipeline`, `tdd-subloop`, etc.).

5. **Interface do Usuário (Zero Emojis):**
   * Modal de importação no Studio acionado por botão SVG limpo no Catálogo.
   * Feedback imediato via toast e atualização reativa do inventário.

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

* [ ] **CA-24.1 (Validação de URL e Segurança do Importador):** `ParseGitHubURL` rejeita URLs fora do padrão `github.com` e previne injeção de caracteres de travessia de diretório.
* [ ] **CA-24.2 (Instalação e Canonização SHA-256):** `ImportSkill` baixa o conteúdo para `~/.gemini/config/skills/<nome>`, calcula o hash SHA-256 dos arquivos e valida o `SKILL.md`.
* [ ] **CA-24.3 (Endpoint Studio `/api/assets/import`):** A rota HTTP POST processa a requisição, invoca o importador e dispara o refresh do inventário.
* [ ] **CA-24.4 (Comando CLI `legoloop import`):** A CLI expõe o comando `import` com saída clara em stdout e código de saída 0 em caso de sucesso.
* [ ] **CA-24.5 (Endpoint de Fork `/api/loops/{id}/fork`):** Duplicação imutável Copy-on-Write cria nova versão com ID derivado sem alterar o original.
* [ ] **CA-24.6 (Zero Regressões & 100% Testes Verdes):** Toda a suíte de testes passa com sucesso (`go test ./...`) e o binário é recompilado.
