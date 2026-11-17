# Relatório de Auditoria Universal dos 6 Verbos

- **Data:** 2026-11-17
- **Processo:** auditoria-atomica

| Verbo | Veredito | Código | Lacunas |
| :--- | :--- | :---: | :---: |
| matriz | APROVADO | 0 | 0 |
| listar | APROVADO | 0 | 0 |
| refinar | APROVADO | 0 | 0 |
| remover | APROVADO | 0 | 0 |
| compor | APROVADO | 0 | 0 |
| fechar | APROVADO | 0 | 0 |

# Relatório de Auditoria — matriz (verbo)

## Disseção

- Tipo: verbo
- Item: matriz
- Caminho: internal\matriz
- Arquivos: 18
- Itens: calibragem_test.go, catalogador.go, catalogador_test.go, confirmado.go, confirmado_test.go, frontmatter.go, generate.go, generate_test.go, incremental_test.go, origem_test.go, perfil.go, rank.go, rank_test.go, registro.go, stale.go, stale_test.go, vocab.default.json, vocab.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/matriz/...` — exit 0
  - Saída: ok legoloop/internal/matriz 3.237s


---

# Relatório de Auditoria — listar (verbo)

## Disseção

- Tipo: verbo
- Item: listar
- Caminho: internal\listar
- Arquivos: 5
- Itens: listar.go, listar_test.go, pecas.go, sugestao.go, sugestao_test.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/listar/...` — exit 0
  - Saída: ok legoloop/internal/listar 0.986s


---

# Relatório de Auditoria — refinar (verbo)

## Disseção

- Tipo: verbo
- Item: refinar
- Caminho: internal\refinar
- Arquivos: 9
- Itens: adicionar.go, adicionar_test.go, curado_v2.go, curado_v2_test.go, dag_iso_test.go, layout_canonico_test.go, refinar.go, refinar_test.go, triagem_test.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/refinar/...` — exit 0
  - Saída: ok legoloop/internal/refinar 7.521s


---

# Relatório de Auditoria — remover (verbo)

## Disseção

- Tipo: verbo
- Item: remover
- Caminho: internal\remover
- Arquivos: 3
- Itens: cartao.go, remover.go, remover_test.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/remover/...` — exit 0
  - Saída: ok legoloop/internal/remover 0.977s


---

# Relatório de Auditoria — compor (verbo)

## Disseção

- Tipo: verbo
- Item: compor
- Caminho: internal\compor
- Arquivos: 3
- Itens: compor.go, compor_test.go, layout_canonico_test.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/compor/...` — exit 0
  - Saída: ok legoloop/internal/compor 7.926s


---

# Relatório de Auditoria — fechar (verbo)

## Disseção

- Tipo: verbo
- Item: fechar
- Caminho: internal\fechar
- Arquivos: 12
- Itens: artefatos.go, destilar.go, destilar_test.go, fechar.go, fechar_real_test.go, fechar_stub_test.go, fechar_test.go, natureza.go, quarteto.go, quarteto_test.go, report.go, types.go

## Lacunas

Nenhuma lacuna.

## Tickets

Nenhum ticket.

## Refatoração

Nenhuma refatoração sugerida.

## Veredito

**APROVADO**

### Evidências

- Comando: `go test -count=1 ./internal/fechar/...` — exit 0
  - Saída: ok legoloop/internal/fechar 2.313s


---

## Mutações provadas

Cada mutação forçada fez o teste alvo falhar e foi restaurada em seguida (Green-Red-Green).

- M1: veredito invertido em fecharRelatorio — TestAlvoArbitrario e TestAuditoriaSeisVerbosCriterioBinario falharam — restaurada.
- M2: etapa e05_veredito removida em etapasAuditoria — TestScaffoldAuditoriaAtomica falhou — restaurada.
- M3: seção Refatoração omitida em Renderizar — TestAlvoArbitrario falhou — restaurada.
- M4: exit do critério binário ignorado em Auditar — TestAuditoriaSeisVerbosCriterioBinario falhou — restaurada.
- M5: erro de alvo inexistente ignorado em AuditarEGravar — TestAlvoInexistenteExitDois falhou — restaurada.
- M6: gravação dos relatórios atômicos suprimida em AuditarVerbosEGravar — TestAuditoriaSeisVerbosCriterioBinario falhou — restaurada.

## Não provado

- "-race" não rodou: ambiente sem CGO.
- Auditoria com LLM real não rodou: os testes usam o critério binário simulado.
- Distribuição universal aos harnesses não foi executada nesta sessão.
- A suíte do repositório tem flakes de concorrência pré-existentes (internal/runner, internal/studio): passam isolados e falham de forma intermitente sob a carga da suíte paralela.
