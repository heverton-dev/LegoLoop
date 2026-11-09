# 🎯 INTENT: Target Workspace Isolado por Vídeo no Pipeline Video REA

> **Status:** APPROVED  
> **Data:** 2026-11-08  
> **Autor:** Antigravity  
> **Iniciativa:** I04 — Arquitetura de Workspaces Isolados e Mídia Soberana  
> **Alvo:** `internal/videorea/`, `internal/cli/video_rea.go`, `.gitignore`  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Atualmente, a dissecação de vídeos pelo pipeline `legoloop-video-rea` gera arquivos voláteis e persistentes misturados em diretórios genéricos (`scratch/`, `docs/intents/`, etc.). Analisar dezenas de vídeos resulta em poluição da raiz do projeto, acúmulo de arquivos pesados não versionáveis e risco de órfãos ao descartar um vídeo analisado.
* **Quem Sofre:** Desenvolvedores, operadores de conteúdo audiovisual e agentes que consomem o ecossistema e precisam de isolamento estéril por alvo.
* **O Porquê:** Ao encapsular cada análise em seu próprio Target Bundle (`docs/videorea/<slug>/`), garantimos isolamento cirúrgico de ciclo de vida (camada volátil em `scratch/`, forense em `dossier.json`, curadoria em `curation.json`, governança em `governance/`, livro em `book/` e grafo em `loop/`), sem quebrar a interoperabilidade do LegoLoop.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Arquivos de análise espalhados por diretórios raiz globais.
  * **Amanhã:** Cada vídeo possui sua pasta autocontida e, após aprovação do Stop-Gate, realiza a projeção canônica leve em `loops/<slug>/` e `skills/<slug>/`.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo:
* [x] Definir e implementar a estrutura `TargetBundle` em `internal/videorea/`.
* [x] Atualizar `factory_bridge.go` para materializar a árvore de diretórios do bundle (`scratch/`, `governance/`, `book/`, `loop/`, `skill/`).
* [x] Implementar a projeção canônica automática (registro leve em `loops/` e `skills/` para consumo transparente pelo LegoLoop e pelos harnesses).
* [x] Atualizar `internal/cli/video_rea.go` para adotar a nova estrutura por padrão ou via flag explícita.
* [x] Blindar o repositório adicionando `docs/videorea/**/scratch/` ao `.gitignore`.
* [x] Cobrir 100% com testes unitários em `internal/videorea/`.

### O que NÃO ESTÁ no Escopo:
* ⛔ Não alterar os motores globais do LegoLoop (`fechar`, `instalar-skills`, `matriz`, `listar`).
* ⛔ Não alterar o formato dos templates canônicos v2 de governança.
* ⛔ Não fazer download de arquivos de mídia reais em ambiente de teste (permanecer com Golden Fixtures sintéticas herméticas).

---

## 3. Critérios de Sucesso
* Criação completa da árvore do bundle em uma única chamada.
* Existência da Golden Fixture (3s) e dos manifestos dentro do bundle.
* Projeção dos arquivos canônicos nos diretórios raiz sem erros.
* Testes unitários verdes e conformidade estrita com `funclen <= 35` e `gofmt`.
