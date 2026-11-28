---
name: legoloop-video-roteiro
author: LegoLoop
description: Use quando precisar executar o Processo legoloop-video-roteiro (antigo legoloop-video-roteiro) (Ciclo 1 do processo composto legoloop-video-animado): validação estrita de projeto.json e marca.json, pré-voo modular, descoberta gratuita por nicho com bloqueio em acervo/temas.json e validação epistêmica universal de script.json e plano.md. NÃO use para geração de ativos visuais, renderização ou publicação.
consumes: Artifact:ProjectConfig
produces: Artifact:ValidatedScript
artifact: script.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🎬 legoloop-video-roteiro — Pauta, Fontes e Roteiro Verificado

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Motor oficial:** `legoloop video-roteiro` (implementação em `internal/video/roteiro` e `internal/cli/diorama_roteiro.go`).  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** para validar `projeto.json` e `marca.json`, executar o pré-voo modular, coletar pautas gratuitas filtradas por `perfis_fonte` sem repetir `acervo/temas.json`, e validar `script.json` e `plano.md` sob classificação fail-closed por domínio (`A`, `B`, `preprint`, `C`) e 4 naturezas epistêmicas universais.  
**Quando NÃO usar:** para geração de sprites/Flow (Processo 2), síntese de voz/renderização (Processo 3) ou publicação Postiz (Processo 4).

* **Consumes (Entradas Obrigatórias):**
  - `projeto.json`: `nicho` (texto livre), `perfis_fonte`, `fontes_extras`/`fontes_excluidas`, `formatos`, `duracao_alvo_s`, `politica_gates` (`humano`, `mecanico` ou `hibrido` com `gates_humanos`) e `video_referencia` opcional.
  - `marca.json`: identidade do canal (`nome`, `formatos`, `voz`).
  - `acervo/temas.json` (opcional na 1ª execução): temas já produzidos para bloqueio antirrepetição.
* **Produces (Saídas Mecânicas):**
  - `script.json`: roteiro retrocompatível com `lines[].text`, `natureza`, `sources[]` e lastro verificado.
  - `plano.md`: tabela humana de batidas e dispositivos visuais.
  - `acervo/temas.json`: atualizado de forma atômica somente após aprovação no Gate final da Etapa 2 (temas já produzidos).
  - `acervo/fila_perene.json`: estoque automático apenas de pautas `perene` coletadas para futuras produções (pautas `quente` nunca entram na fila; tema aprovado sai da fila e vai para `acervo/temas.json`).

```json
{
  "schema": "diorama-roteiro.v1",
  "input": { "projeto": "projeto.json", "marca": "marca.json", "acervo": "acervo/temas.json" },
  "output": { "script": "script.json", "plano": "plano.md", "acervo": "acervo/temas.json" }
}
```

---

## 2. Procedimento Determinístico e Nós do DAG

0. **Escolha da superfície interativa (obrigatória ao iniciar por linguagem natural, slash ou CLI):** ofereça sempre ao usuário escolher entre **Interface A2UI** (Agente/Chat: ferramentas MCP `sessao_criar`, `sessao_publicar`, `sessao_status` e `sessao_aguardar`; Studio: painel `Video Roteiro`) ou **Chat textual / CLI direto**, permitindo definir o tempo limite por etapa (`timeout_ms`, padrão 5 min; ajustável e pausável na barra da A2UI). Só avance no modo escolhido.
   - **Ciclo contínuo obrigatório no modo A2UI (MCP):** chame `sessao_criar` (abre o navegador automaticamente) -> `sessao_publicar` -> `sessao_aguardar` **no mesmo turno**, mantendo toda a interação dentro da interface até concluir o Processo. **NUNCA** encerre o turno pedindo ao usuário "clique lá e me avise aqui no chat". Durante execuções e validações, chame `sessao_status` com `mensagem` e `etapas[]` para mostrar o avanço ao vivo no navegador.
1. **Etapa 0 — Configuração guiada (`n01_validar_config`):** se `projeto.json` ou `marca.json` não existirem, colete com o usuário (na superfície escolhida) o nicho, `perfis_fonte`, formato (`16:9`, `9:16`, `1:1`), duração-alvo, `politica_gates` (`humano`, `hibrido` ou `mecanico`) e pasta de trabalho. Grave `projeto.json` e `marca.json` e rode `legoloop video-roteiro validar-config --projeto <DIR> --json` (bloqueia schema inválido, perfil desconhecido ou inconsistência entre `projeto.json` e `marca.json`).
2. **Pré-voo modular (`n02_prevoo_modular`):** rode `legoloop video-roteiro prevoo --projeto <DIR> --json` (bloqueia apenas nos 5 casos `FALTA` do Processo 1; Chrome Flow `9333`, FFmpeg, Chromium, `edge-tts` e `gitleaks` emitem apenas `AVISO`).
3. **Etapa 1 — Coleta de pautas (`n03_coletar_pautas`):** rode `legoloop video-roteiro coletar --projeto <DIR> --json` (filtra conectores gratuitos por `perfis_fonte`, descarta temas em `acervo/temas.json` e apresenta as pautas na superfície escolhida para escolha do usuário quando `politica_gates` exigir).
4. **Etapa 2 — Roteiro e Gate epistêmico (`n04_validar_roteiro`):** gere `script.json` (com `lines[].text`, `natureza`, `sources[]`) e `plano.md` a partir da pauta escolhida, e rode `legoloop video-roteiro validar-roteiro --projeto <DIR> --json` (exige 1 fonte A ou 2 fontes B distintas em `fato_consolidado`; escopo/preprint/amostra em `evidencia_preliminar`; atribuição em `alegacao_terceiro`; bloqueia nível C como prova; no máximo 2 correções automáticas).

| Nó ID | Subcomando | Papel |
| :--- | :--- | :--- |
| `n01_validar_config` | `legoloop video-roteiro validar-config` | Validação estrita de `projeto.json` e `marca.json` |
| `n02_prevoo_modular` | `legoloop video-roteiro prevoo` | Pré-voo modular determinístico (5 bloqueantes `FALTA`, demais `AVISO`) |
| `n03_coletar_pautas` | `legoloop video-roteiro coletar` | Descoberta gratuita por `perfis_fonte` e filtro em `acervo/temas.json` |
| `n04_validar_roteiro` | `legoloop video-roteiro validar-roteiro` | Validação epistêmica universal de `script.json` e `plano.md` |

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** aceite domínio nível C (fora do catálogo e fora de `fontes_extras`) como prova factual.
> 2. **NUNCA** grave o tema em `acervo/temas.json` antes da aprovação do Gate final da Etapa 2, em nenhum modo.
> 3. **NUNCA** ultrapasse 2 correções automáticas no Gate final; se persistir a falha, pare mesmo no modo `mecanico`.
> 4. **NUNCA** consulte conectores pagos ou conectores fora de `perfis_fonte` do projeto.
> 5. **NUNCA** bloqueie o Processo 1 por ausência de dependências exclusivas dos Processos 2 e 3 (`AVISO` apenas).

---

## 4. Comandos de Verificação & Testes de Terminal

```bash
legoloop video-roteiro --help
legoloop validate processos/legoloop-video-roteiro/*.dag.json
go test -count=1 ./internal/video/roteiro/... ./internal/cli/ -run TestDioramaRoteiroCLI
```

---

## 5. Critério de Parada (Stopping Checklist)

- [ ] `projeto.json` e `marca.json` aprovados em `validar-config` (exit code 0).
- [ ] Pré-voo modular com zero `FALTA` (exit code 0).
- [ ] `script.json` aprovado em `validar-roteiro` com zero violações de natureza ou nível C.
- [ ] Tema registrado em `acervo/temas.json` somente após aprovação do Gate final.
- [ ] DAG validado com exit code 0 em `legoloop validate`.
