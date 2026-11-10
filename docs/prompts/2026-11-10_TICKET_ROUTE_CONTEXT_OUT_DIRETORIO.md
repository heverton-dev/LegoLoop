# TICKET: write_boundary rejeita context_out de diretório (TICK-03.3)

Data: 2026-11-10. Origem: sessão orquestradora da trilha Lei #15. Destino: outra sessão (engenharia do motor/gerador).

## Fato
- O nó TICK-03.3 declara `context_out: ["internal/studiodoc/"]` (diretório com barra final).
- O pacote tem 9 arquivos. O limite de 12 arquivos não causa a falha.
- O motor abortou com `write_boundary: "<arquivo>" fora das saídas permitidas` para 6 arquivos:
  hierarchy.go, ports.go, sanitize.go, sanitize_test.go, store.go, types.go.
- Custo da tentativa: US$ 0,38; 11155 tokens de saída; nenhum arquivo alterado.
- O gate do nó falha por 149 violações do termo "loop" em `internal/studiodoc/*`.

## Hipóteses (não verificadas)
1. O motor compara `context_out` por igualdade de caminho e ignora prefixo de diretório.
2. O gerador `tools/route2dag` deveria expandir o diretório em arquivos.

## Tarefa
1. Reproduza com o DAG de um nó só: `legoloop run <dag> --allow-unit-tasks --parallel 1`.
2. Decida o dono da correção: motor (aceitar prefixo de diretório) ou gerador (listar arquivos).
3. Aplique a correção com teste que cobre `context_out` de diretório.
4. Verifique outros nós com `context_out` de diretório no DAG da Lei #15.
5. Devolva a trilha: rodar o TICK-03.3 e seguir até o SG-E03.

## Restrições
- O claude do nó não tem shell. Não use unit_task para renomear pacote.
- Commit só dos arquivos da correção. Sem push.
