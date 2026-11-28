# Manual Operacional — `legoloop-video-roteiro`

- **Data:** 2026-11-24
- **Versão:** 1.0.0
- **Motor CLI:** `legoloop video-roteiro`
- **Pacote Go:** `internal/video/roteiro` + `internal/cli/diorama_roteiro.go`
- **DAG Canônico:** `processos/legoloop-video-roteiro/pipeline.dag.json`
- **Skill Universal:** `skills/legoloop-video-roteiro/SKILL.md`

---

## 1. Propósito

O Processo `legoloop-video-roteiro` corresponde ao Ciclo 1 do processo composto `legoloop-video-animado`. Ele garante que nenhum roteiro avance sem configuração consistente, ambiente válido, pauta inédita e lastro factual verificável.

## 2. Comandos e Subcomandos

| Subcomando | Entrada | Saída e Gate |
| :--- | :--- | :--- |
| `legoloop video-roteiro validar-config --projeto DIR [--catalogo ARQ] [--json]` | `projeto.json`, `marca.json` | Valida campos obrigatórios, `perfis_fonte`, `politica_gates` e consistência cruzada (`exit 0`) |
| `legoloop video-roteiro prevoo --projeto DIR [--acervo DIR] [--catalogo ARQ] [--json]` | `projeto.json`, ambiente local | Bloqueia nos 5 itens `FALTA`; emite `AVISO` para itens de P2 e P3 (`exit 0` se zero `FALTA`) |
| `legoloop video-roteiro coletar --projeto DIR [--acervo ARQ] [--catalogo ARQ] [--json]` | `projeto.json`, `acervo/temas.json` | Consulta apenas conectores gratuitos aderentes a `perfis_fonte` e remove temas repetidos |
| `legoloop video-roteiro validar-roteiro --projeto DIR [--catalogo ARQ] [--json]` | `projeto.json`, `script.json` | Valida as 4 naturezas epistêmicas e bloqueia domínio nível C como prova (`exit 0` se aprovado) |

## 3. Regras Inegociáveis

1. **Classificação fail-closed por domínio:** domínios oficiais e científicos declarados no catálogo são nível `A`; imprensa de referência declarada é nível `B`; preprints exigem `evidencia_preliminar`; qualquer domínio fora do catálogo e fora de `fontes_extras` é nível `C` e fica bloqueado como prova.
2. **Naturezas epistêmicas universais (sem viés jurídico):** `fato_consolidado` (>= 1 fonte A ou >= 2 fontes B distintas), `evidencia_preliminar` (escopo, preprint ou amostra declarados), `alegacao_terceiro` (atribuição explícita a quem afirmou) e `narrativa_transicao` (gancho ou pergunta).
3. **Política híbrida de Gates:** `projeto.json` define `politica_gates.modo` (`humano`, `mecanico` ou `hibrido` com `gates_humanos`).
4. **Limite de correções e proteção do acervo:** o Gate final executa no máximo 2 correções automáticas; se persistir a falha, o Processo para mesmo no modo `mecanico` e nunca grava o tema em `acervo/temas.json`.

## 4. Verificação Mecânica

```bash
go test -count=1 -cover ./internal/video/roteiro/...
legoloop validate processos/legoloop-video-roteiro/pipeline.dag.json
legoloop video-roteiro --help
```
