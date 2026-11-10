# 📋 Relatório de Conclusão: {{TITULO_DA_DEMANDA_OU_PROCESSO}}

> **Status:** {{CONCLUÍDO | PARCIAL | FALHA}}  
> **Data:** {{DATA}}  
> **Autor:** {{AUTOR}}  
> **Iniciativa / Épico:** {{INICIATIVA_OU_EPICO | Opcional}}  
> **Processo / Slug:** {{SLUG_DO_PROCESSO}}  
> **Sessão:** {{UUID_DA_SESSAO}}  
> **Template:** TEMPLATE_REPORT.md v2.2.0  
> **Referências:** [INTENT]({{LINK_INTENT}}) · [SPEC]({{LINK_SPEC}}) · [PLAN]({{LINK_PLAN}}) · [ROUTE]({{LINK_ROUTE}})

<!--
GUIA DE PREENCHIMENTO (apague este trecho ao concluir)
- Todo {{PLACEHOLDER}} deve ser preenchido ou substituído por dados factuais auditáveis.
- A telemetria DEVE conter os números reais orçados e os medidos factualmente.
- Se algum indicador não for aplicável, declare 0 ou N/A com justificativa explícita.
- Nome do arquivo: docs/reports/YYYY-MM-DD_REPORT_<NOME>.md (AGENTS.md, regra 6).
-->

---

## 1. Resumo Executivo
* **Objetivo Alcançado:** {{Descrição direta do objetivo completado e problema resolvido.}}
* **Impacto Arquitetural:** {{Quais subsistemas, pacotes Go, Processos ou skills foram alterados ou criados.}}
* **Veredito Geral:** {{APROVADO_SEM_RESSALVAS | APROVADO_COM_DEBITO_TECNICO | REPROVADO}}

---

## 2. Telemetria Estrita & Balanço de Recursos (Orçado vs Realizado)

### 2.1 Balanço de Tempo & Cronometria
| Dimensão de Tempo | Orçado (Estimado) | Realizado (Medido) | Desvio (%) | Status |
| :--- | :---: | :---: | :---: | :---: |
| **Duração Total da Sessão** | {{TEMPO_ORCADO_TOTAL}} | {{TEMPO_GASTO_TOTAL}} | {{DESVIO_TEMPO_PCT}}% | {{DENTRO_DO_ORCAMENTO | ACIMA}} |
| **Tempo de Inferência Cognitiva (LLM)** | {{TEMPO_ORCADO_LLM}} | {{TEMPO_GASTO_LLM}} | {{DESVIO_LLM_PCT}}% | {{OK | ALERTA}} |
| **Tempo de Compilação & Testes (TDD)** | {{TEMPO_ORCADO_BUILD}} | {{TEMPO_GASTO_BUILD}} | {{DESVIO_BUILD_PCT}}% | {{OK | ALERTA}} |
| **Tempo de I/O, Scripts & Ferramentas** | {{TEMPO_ORCADO_IO}} | {{TEMPO_GASTO_IO}} | {{DESVIO_IO_PCT}}% | {{OK | ALERTA}} |

### 2.2 Balanço de Recursos Computacionais & Tokens
| Recurso / Métrica | Orçado (Budget) | Consumido (Real) | Desvio (%) | Eficiência / Status |
| :--- | :---: | :---: | :---: | :---: |
| **Tokens de Entrada (Prompt)** | {{TOKENS_IN_ORCADOS}} | {{TOKENS_IN_CONSUMIDOS}} | {{DESVIO_TOKENS_IN_PCT}}% | {{OK | ACIMA}} |
| **Tokens de Saída (Completion)** | {{TOKENS_OUT_ORCADOS}} | {{TOKENS_OUT_CONSUMIDOS}} | {{DESVIO_TOKENS_OUT_PCT}}% | {{OK | ACIMA}} |
| **Volume Total de Tokens** | {{TOKENS_TOTAL_ORCADOS}} | {{TOKENS_TOTAL_CONSUMIDOS}} | {{DESVIO_TOKENS_TOTAL_PCT}}% | {{OK | ACIMA}} |
| **Turnos de Conversação / Requisições** | {{TURNOS_ORCADOS}} | {{TURNOS_REALIZADOS}} | {{DESVIO_TURNOS_PCT}}% | {{OK | ACIMA}} |
| **Chamadas de Ferramentas (Tool Calls)** | {{TOOL_CALLS_ORCADAS}} | {{TOOL_CALLS_REALIZADAS}} | {{DESVIO_TOOLS_PCT}}% | {{OK | ACIMA}} |
| **Cache Hits de Contexto (%)** | {{CACHE_HIT_ESTIMADO}}% | {{CACHE_HIT_REAL}}% | {{DESVIO_CACHE_PCT}}% | {{ALTA_EFICIENCIA | MODERADA}} |
| **Modelos Cognitivos Utilizados** | {{MODELOS_PLANEJADOS}} | {{MODELOS_UTILIZADOS}} | N/A | {{CONFORME_ROTA}} |

### 2.3 Balanço Financeiro & Custo (USD / BRL)
| Item Financeiro | Orçado (USD) | Realizado (USD) | Custo em BRL (R$) | Saldo / Status |
| :--- | :---: | :---: | :---: | :---: |
| **Custo de Entrada (Input Tokens)** | ${{CUSTO_IN_ORCADO_USD}} | ${{CUSTO_IN_REAL_USD}} | R$ {{CUSTO_IN_REAL_BRL}} | {{SALDO_IN}} |
| **Custo de Saída (Output Tokens)** | ${{CUSTO_OUT_ORCADO_USD}} | ${{CUSTO_OUT_REAL_USD}} | R$ {{CUSTO_OUT_REAL_BRL}} | {{SALDO_OUT}} |
| **Custo de Ferramentas / MCPs / Infra** | ${{CUSTO_TOOLS_ORCADO_USD}} | ${{CUSTO_TOOLS_REAL_USD}} | R$ {{CUSTO_TOOLS_REAL_BRL}} | {{SALDO_TOOLS}} |
| **CUSTO TOTAL CONSOLIDADO** | **${{CUSTO_TOTAL_ORCADO_USD}}** | **${{CUSTO_TOTAL_REAL_USD}}** | **R$ {{CUSTO_TOTAL_REAL_BRL}}** | **{{ECONOMIA_OU_EXCEDENTE}}** |

> **Taxa Cambial Referencial Aplicada:** 1 USD = R$ {{TAXA_CAMBIO_BRL}}  
> **Custo Médio por 1k Tokens:** ${{CUSTO_MEDIO_1K_TOKENS}} USD  

### 2.4 Análise de Eficiência & Gargalos
* **Taxa de Acerto no Primeiro Tiro (First-Shot Pass Rate):** {{NN}}% (testes que passaram sem retentativa).
* **Gargalos Detectados:** {{Descrição sucinta de timeout, compilação demorada ou ciclos de correção.}}
* **Economia Gerada:** {{Explicação do ganho por prompt caching, VSA estrita ou isolamento por Ticket.}}

---

## 3. Entregas por Épico & Componentes

### Épico 01 — {{NOME_DO_EPICO_01}}
* **Código de Produção:** `{{PACOTE_OU_ARQUIVO_GO}}` ({{DESCRICAO_DA_IMPLEMENTACAO}}).
* **Testes Automatizados:** `{{ARQUIVO_TEST_GO}}` ({{NUM_TESTES}} testes unitários, {{COBERTURA}}% cobertura).
* **Contratos & Schemas:** {{Alterações em structs, JSON schemas ou interfaces}}.

### Épico 02 — {{NOME_DO_EPICO_02}}
* **Código de Produção:** `{{PACOTE_OU_ARQUIVO}}` ({{DESCRICAO}}).
* **CLI & Ferramentas:** `{{SUBCOMANDO_OU_SCRIPT}}`.
* **Skills & Governança:** `{{ARQUIVOS_SKILL}}`.

### Épico 03 — {{NOME_DO_EPICO_03}}
* **Integrações & Grafos:** `{{GRAFOS_OU_PROCESSOS}}`.
* **Universalização nos Harnesses:** {{Vínculos nos harnesses claude, gemini, cursor, codex, etc.}}.

---

## 4. Matriz de Qualidade & Portões Mecânicos de Aceite

| Portão Mecânico | Comando de Validação | Critério Inviolável | Resultado |
| :--- | :--- | :--- | :---: |
| **TDD / Testes Unitários** | `go test -count=1 ./internal/...` | 100% dos testes novos passando | **PASS** |
| **Suíte Global de Regressão** | `go test -count=1 ./...` | Zero falhas em todos os pacotes | **PASS** |
| **Análise Estática (Vet)** | `go vet ./...` | Zero warnings ou erros | **PASS** |
| **Formatação Canônica** | `gofmt -l .` | Saída 100% vazia | **PASS** |
| **Limite de Complexidade** | `go run ./tools/funclen internal` | Nenhuma função > 35 linhas | **PASS** |
| **Compilação Limpa** | `go build -o bin/legoloop.exe ./cmd/legoloop` | Binário gerado sem warnings | **PASS** |
| **Auditoria Geométrica** | `legoloop legopress auditar <pdf>` | 0 linhas órfãs / layout perfeito | **PASS** |
| **Grafo Semântico** | `codebase-memory-mcp cli index_repository` | Grafo atualizado no disco | **PASS** |

---

## 5. Inventário de Artefatos & Linhagem Criptográfica

| Artefato Gerado / Modificado | Tipo | Papel no Ecossistema | Hash SHA-256 |
| :--- | :--- | :--- | :--- |
| `{{CAMINHO_ARTEFATO_1}}` | `{{TIPO_1}}` | {{PAPEL_1}} | `{{SHA256_1}}` |
| `{{CAMINHO_ARTEFATO_2}}` | `{{TIPO_2}}` | {{PAPEL_2}} | `{{SHA256_2}}` |
| `{{CAMINHO_ARTEFATO_3}}` | `{{TIPO_3}}` | {{PAPEL_3}} | `{{SHA256_3}}` |

> **Assinatura de Fechamento Merkle:** `{{MERKLE_ROOT_OU_COMMIT_SHA}}`  

---

## 6. Conclusão, Recomendações & Rollback
* **Veredito de Homologação:** A demanda foi concluída satisfazendo todos os critérios do PLAN e SPEC com aprovação nos portões de qualidade.
* **Próximos Passos Recomendados:** {{Ações pós-fechamento, monitoramento em produção ou evolução de próximos tickets.}}
* **Procedimento de Rollback:**
  ```powershell
  git revert {{COMMIT_SHA}}
  go build -o bin/legoloop.exe ./cmd/legoloop
  go install ./cmd/legoloop
  ```
