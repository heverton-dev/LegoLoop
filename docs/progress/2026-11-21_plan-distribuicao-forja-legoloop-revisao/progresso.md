---
schema: progresso.v1
iniciativa: "I01"
plano: "docs/plans/2026-11-21_PLAN_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md"
atualizado_em: "2026-11-22T18:19:03Z"
---

# Progresso Vivo: I01

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 38 | 0 | 0 | 0 | 0 | 38 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-01.1` | E01 | Conferir a SPEC revisada aprovada contra os Gates de lacuna e de aceite | AC-E01.1 |
| `TICK-01.2` | E01 | Renomear o módulo para github.com/heverton-dev/LegoLoop (pré-requisitos I04: G01 Matriz por hash, G02 descoberta de peças, G05 fachada única de caminhos) | AC-E01.3 |
| `TICK-01.3` | E01 | Gravar a linha de base de contexto por harness | AC-E01.2 |
| `TICK-03.1` | E03 | Criar ste.AplicarBloco idempotente para o bloco gerenciado | AC-E03.1 |
| `TICK-03.2` | E03 | Provar a paridade com o script Python e remover o script | AC-E03.2 |
| `TICK-03.3` | E03 | Provar que o caminho de instalação não usa Python e listar Python como opcional no doctor | AC-E03.3 |
| `TICK-04.1` | E04 | Criar os contratos da Forja e a validação das opções | AC-E04.6 |
| `TICK-04.2` | E04 | Implementar a simulação --dry-run que não grava nada | AC-E04.1 |
| `TICK-04.3` | E04 | Implementar o aplicar com backup, MANIFESTO e idempotência | AC-E04.2 |
| `TICK-04.4` | E04 | Implementar o --desfazer que restaura byte a byte | AC-E04.3 |
| `TICK-04.5` | E04 | Gravar o bloco gerenciado que aponta para o AGENTS.md no GEMINI.md e no CLAUDE.md | AC-E04.4 |
| `TICK-04.6` | E04 | Mesclar o settings.json do Claude preservando as chaves desconhecidas | AC-E04.5 |
| `TICK-04.7` | E04 | Implementar --sem e recusar perfil inválido | AC-E04.6 |
| `TICK-04.8` | E04 | Gerar os lançadores por SO sem PATH fixo | AC-E04.7 |
| `TICK-04.9` | E04 | Garantir que segredos não entram no relatório nem na saída | AC-NF.1 |
| `TICK-04.10` | E04 | Aplicar path jail na resolução da pasta do harness | AC-NF.2 |
| `TICK-04.11` | E04 | Medir o desempenho do dry-run com 24 harnesses | AC-NF.3 |
| `TICK-05.1` | E05 | Criar legoloop instalar encadeando a Forja e entregando o relatório | AC-E05.1 |
| `TICK-05.2` | E05 | Criar o atalho do grafo semântico no PATH, com dry-run | AC-E05.2 |
| `TICK-05.3` | E05 | Marcar o bootstrap como substituído pelo instalar | AC-E05.3 |
| `TICK-06.1` | E06 | Implementar forjar --verificar --orcamento com reprovação por teto estourado | AC-E06.1 |
| `TICK-06.2` | E06 | Aprovar quando o custo respeita o teto | AC-E06.2 |
| `TICK-07.1` | E07 | Criar o Processo legoloop-forjar: DAG, skill e slash command | AC-E07.2 |
| `TICK-07.2` | E07 | Provar que os quatro gatilhos geram o mesmo RELATORIO.json | AC-E07.3 |
| `TICK-07.3` | E07 | Fechar o Processo com legoloop fechar | AC-E07.1 |
| `TICK-08.1` | E08 | Criar tools/release com build reproduzível | AC-E08.1 |
| `TICK-08.2` | E08 | Gerar pacotes por SO e arquitetura, SHA256SUMS e SBOM | AC-E08.2 |
| `TICK-08.3` | E08 | Declarar a assinatura no relatório do release e assinar com cosign no CI | AC-E08.3 |
| `TICK-09.1` | E09 | Criar install.sh e install.ps1 que conferem o SHA-256 e abortam se divergir | AC-E09.1 |
| `TICK-09.2` | E09 | Fazer os scripts instalarem o binário e dispararem a Forja | AC-E09.2 |
| `TICK-09.3` | E09 | Gerar os manifestos de winget, scoop e brew | AC-E09.3 |
| `TICK-09.4` | E09 | Criar a definição WiX do MSI por usuário, sem administrador | AC-E09.6 |
| `TICK-09.5` | E09 | Provar git clone e go install remoto | AC-E09.4 |
| `TICK-09.6` | E09 | Expor o instalar como skill, slash command e ferramenta do MCP | AC-E09.5 |
| `TICK-10.1` | E10 | Implementar forjar --verificar com detecção de deriva | AC-E10.1 |
| `TICK-10.2` | E10 | Implementar --corrigir que devolve a skill ao Cofre com backup | AC-E10.2 |
| `TICK-11.1` | E11 | Criar legoloop atualizar que aborta com SHA-256 divergente | AC-E11.1 |
| `TICK-11.2` | E11 | Guardar o binário anterior e revalidar a Forja depois de atualizar | AC-E11.2 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| nenhum | - | - | - | - |
