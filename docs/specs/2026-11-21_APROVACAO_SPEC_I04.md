# Stop-Gate 02: SPEC da I04 aprovada

Decisão: **APROVADA** pelo usuário em 2026-11-21 ("[1]", com duas ordens junto: inserir o prefixo `legoloop-` nos 10 Processos e verificar o que estava "Não verificado").
Artefato: `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md` (60 critérios de aceite, 9 Épicos, 9 decisões).
Decisões: (1) `caminhos` absorve as constantes de `paths`, as funções delegadas saem e o jail vira `internal/jail`; (2) o jail passa a proteger os verbos que recebem nome de Processo; (3) a descoberta liga `matriz.Staleness`, `assets.Scanner` e `matriz gerar -sem-cerebro` ao `fabricar`; (4) o verbo `matrix` é removido; (5) os 9 Processos sem colisão ganham o prefixo pelo verbo novo `legoloop renomear`, e `fechar` é aposentado porque colide com `legoloop-fechar`; (6) o `-race` roda no CI, e o WSL serve para rodar antes; (7) o mutador roda por pacote, com sobreviventes justificados em arquivo; (8) as falhas reais em Linux são corrigidas e os seis testes dependentes de ferramenta ficam herméticos; (9) os Processos que chamam `powershell` ou `python` direto passam a usar iniciador portável.
Resposta à pendência da INTENT ("como tratar os Processos fora do prefixo"): inserir o prefixo, sem lista de exceções.
Verificado antes da aprovação (ver `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md`, seção 4A): `go test ./...` em Linux falha em 9 testes de 5 pacotes; `-race` limpo em 12 pacotes; execução de ponta a ponta determinística; Gate de aceite aprova a SPEC.
Continua não verificado: macOS, unidades de rede e removíveis, desempenho sob DAG grande.
Próxima etapa liberada: PLAN e ROUTE (gerados: 56 Tickets). A execução de I01 a I03 continua congelada.
