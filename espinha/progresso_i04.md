---
schema: progresso.v1
iniciativa: "I04"
plano: "docs/plans/2026-11-21_PLAN_INICIATIVA_4_I04_DO_MACROPROJE.md"
atualizado_em: "2026-11-22T14:38:10Z"
---

# Progresso Vivo: I04

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 56 | 9 | 0 | 1 | 0 | 46 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| `TICK-I4-01.1` | E01 | FAILED | 1 | 0 / |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-I4-03.1` | E03 | Mover as constantes de paths para caminhos e renomear o jail para internal/jail | AC-E03.1 |
| `TICK-I4-03.2` | E03 | Migrar os chamadores e apagar as funções delegadas de paths | AC-E03.2 |
| `TICK-I4-03.3` | E03 | Ligar o jail aos verbos que recebem nome de Processo | AC-E03.3 |
| `TICK-I4-03.4` | E03 | Remover o pacote vault sem importador em produção | AC-E03.4 |
| `TICK-I4-03.5` | E03 | Criar o verbo legoloop renomear com dry-run, backup, colisão e registro | AC-E03.5 |
| `TICK-I4-03.6` | E03 | Renomear os 9 Processos sem colisão com o verbo renomear | AC-E03.6 |
| `TICK-I4-03.7` | E03 | Aposentar o Processo fechar e passar os anticorpos para legoloop-fechar | AC-E03.7 |
| `TICK-I4-03.8` | E03 | Provar que nenhum arquivo cita os nomes antigos dos Processos | AC-E03.8 |
| `TICK-I4-04.1` | E04 | Criar a descoberta de unidades fixas por SO com provedor injetável | AC-E04.1 |
| `TICK-I4-04.2` | E04 | Fazer o TargetDrives restringir a busca e remover o aviso drives_ignorados | AC-E04.2 |
| `TICK-I4-04.3` | E04 | Fazer o fabricar atualizar a Matriz defasada em modo sem cérebro | AC-E04.3 |
| `TICK-I4-04.4` | E04 | Provar que Matriz em dia não dispara atualização | AC-E04.4 |
| `TICK-I4-04.5` | E04 | Provar que a descoberta é determinística e somente leitura, com relatório | AC-E04.5 |
| `TICK-I4-04.6` | E04 | Provar que matriz status fica sem diferença depois do fabricar | AC-E04.6 |
| `TICK-I4-05.1` | E05 | Criar o Registro da Espinha e o schema espinha.v1 | AC-E05.1 |
| `TICK-I4-05.2` | E05 | Exportar cli.Verbos e provar que todo verbo está no Registro | AC-E05.2 |
| `TICK-I4-05.3` | E05 | Ligar cada lei do AGENTS.md a um verificador ou justificativa | AC-E05.3 |
| `TICK-I4-05.4` | E05 | Criar legoloop espinha verificar com a conferência de motores e Processos | AC-E05.4 |
| `TICK-I4-05.5` | E05 | Detectar função duplicada entre pacotes | AC-E05.5 |
| `TICK-I4-05.6` | E05 | Provar as doze regras com fixture com defeito e exit 2 | AC-E05.6 |
| `TICK-I4-05.7` | E05 | Rodar o gate no fechar e provar que a espinha não usa LLM nem relógio | AC-E05.7, AC-NF.1 |
| `TICK-I4-06.1` | E06 | Completar o nascimento 5 de 5 dos três Processos sem ANTICORPOS.md | AC-E06.1 |
| `TICK-I4-06.2` | E06 | Remover o verbo matrix descontinuado | AC-E06.2 |
| `TICK-I4-06.3` | E06 | Criar tools/cobertura e subir board, triagem e cli a 80 por cento | AC-E06.3 |
| `TICK-I4-06.4` | E06 | Criar escrita.GravarAtomico e usá-la nas gravações de estado crítico | AC-E06.4 |
| `TICK-I4-06.5` | E06 | Relatar a falha de telemetria e fixar o limite de erros ignorados | AC-E06.5 |
| `TICK-I4-06.6` | E06 | Provar que os globais e init da espinha estão justificados | AC-E06.6 |
| `TICK-I4-06.7` | E06 | Corrigir o relato do Andon quando o nó é recusado antes de executar | AC-E06.7 |
| `TICK-I4-06.8` | E06 | Remover o arquivo temporário rastreado em scratch | AC-E06.8 |
| `TICK-I4-07.1` | E07 | Fazer o jail recusar caminho hostil em qualquer SO | AC-E07.6 |
| `TICK-I4-07.2` | E07 | Tornar o lock da memória da triagem exclusivo em Linux | AC-E07.7 |
| `TICK-I4-07.3` | E07 | Normalizar os separadores na guarda de deleção | AC-E07.8 |
| `TICK-I4-07.4` | E07 | Tornar herméticos os seis testes que dependem de ferramenta instalada | AC-E07.9 |
| `TICK-I4-07.5` | E07 | Trocar powershell e python direto por iniciador portável nos Processos | AC-E07.10 |
| `TICK-I4-07.6` | E07 | Resolver fonte, python e shell por SO | AC-E07.4 |
| `TICK-I4-07.7` | E07 | Criar o workflow de CI nos três SOs com os gates e go test -race | AC-E07.1, AC-E07.2, AC-E07.3 |
| `TICK-I4-07.8` | E07 | Conferir o CI verde nos três SOs | AC-E07.5 |
| `TICK-I4-08.1` | E08 | Criar o tools/mutar só com biblioteca padrão | AC-E08.1 |
| `TICK-I4-08.2` | E08 | Implementar os sete operadores de mutação | AC-E08.2 |
| `TICK-I4-08.3` | E08 | Reprovar com exit 2 quando um mutante sobrevive | AC-E08.3 |
| `TICK-I4-08.4` | E08 | Matar ou justificar todo mutante dos sete pacotes da espinha | AC-E08.4 |
| `TICK-I4-09.1` | E09 | Revisar a SPEC e o PLAN da I01, da I02 e da I03 pela Emenda 2 | AC-E09.1 |
| `TICK-I4-09.2` | E09 | Tirar o Épico de CI do plano revisado da I01 e listar G01, G02 e G05 como pré-requisito | AC-E09.2 |
| `TICK-I4-09.3` | E09 | Manter o congelamento até o Stop-Gate final e removê-lo só com a aprovação do usuário | AC-E09.3 |
| `TICK-I4-09.4` | E09 | Gravar o relatório de fechamento da I04 com o desfecho das lacunas | AC-E09.4 |
| `TICK-I4-09.5` | E09 | Atualizar os documentos vivos e rodar o portão completo | AC-E09.5, AC-NF.2 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| `TICK-I4-01.2` | E01 | 0 | ok  	legoloop/internal/cli	1.500s / | 2026-11-22T14:37:50Z |
| `TICK-I4-01.3` | E01 | 0 | ok  	legoloop/internal/cli	1.511s / | 2026-11-22T14:37:53Z |
| `TICK-I4-02.1` | E02 | 0 | ok  	legoloop/internal/dag	1.030s / | 2026-11-22T14:37:55Z |
| `TICK-I4-02.2` | E02 | 0 | - | 2026-11-22T14:37:57Z |
| `TICK-I4-02.3` | E02 | 0 | ok  	legoloop/internal/dag	0.710s / | 2026-11-22T14:38:05Z |
| `TICK-I4-02.4` | E02 | 0 | ok  	legoloop/internal/dag	0.713s / | 2026-11-22T14:38:06Z |
| `TICK-I4-02.5` | E02 | 0 | ok  	legoloop/internal/dag	1.183s / | 2026-11-22T14:38:08Z |
| `TICK-I4-02.6` | E02 | 0 | ok  	legoloop/internal/cli	1.177s / | 2026-11-22T14:38:10Z |
| `TICK-I4-08.5` | E08 | 0 | ok  	legoloop/internal/espinha	1.670s / | 2026-11-22T14:36:10Z |
