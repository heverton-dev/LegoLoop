# Espelho Canônico — Distribuição e Forja do LegoLoop

Sessão da entrevista: `intv_1795263519855957000` (8 turnos; cérebro final `claude/sonnet`).
Status: consolidado por Claude a partir das respostas do usuário; aguarda a conferência humana.
Nota: o motor manteve a confiança em 10% e truncou três campos. Este registro substitui o Espelho dele.

## resultado_esperado

Instalação universal do LegoLoop, entregue de uma vez, por todos os canais e em todos os SOs (Windows, Linux e macOS).
O primeiro canal é o script de terminal: baixa o binário, confere o SHA-256 e executa `legoloop instalar`.
Esse comando dispara a Forja (camada 1) e configura todos os harnesses detectados (camada 2), com `AGENTS.md` como SSOT.
MSI, winget, scoop, brew, git clone, agente (linguagem natural e slash command) e MCP chamam o mesmo comando.
A camada 2 vira o Processo canônico Forja, com quatro gatilhos: automático na instalação, `legoloop forjar`, agente (linguagem natural ou `/legoloop-forjar`) e a skill proprietária.
O Processo Forja é fabricado e fechado pelo fluxo oficial: manual, livro LegoPress, skill no Cofre, DAG validado e relatório.
O `GEMINI.md` e o `CLAUDE.md` do usuário recebem um bloco gerenciado que aponta para o `AGENTS.md`, sem perder o conteúdo existente.
A execução usa o Swarm, sob as seis condições da lei 4 do `AGENTS.md`.

## quem_se_beneficia

Usuários e desenvolvedores que instalam o LegoLoop por qualquer canal.
Operadores leigos que pedem a instalação a um agente e veem o resultado no chat.
Quem já tem harnesses com contexto inchado e recebe a correção pela Forja, sem reinstalar.

## por_que_agora

Não existe instalação universal: hoje só há `go install` a partir do clone.
A correção de contexto nos harnesses foi feita à mão e só vale neste Windows.
Um comando único evita divergência entre canais, e a comprovação física elimina o "acho que instalou".

## criterio_de_sucesso

Em máquina limpa, o canal escolhido instala o binário e `legoloop instalar` roda sem erro.
`legoloop skills-custo` mostra cada harness detectado com 1 skill.
Um segundo `--dry-run` não encontra pendência, e todo arquivo alterado tem backup restaurável.
Comprovação ao usuário: resumo no terminal (uma linha por harness, OK ou FALHA), `RELATORIO.md` e `RELATORIO.json` em `~/.legoloop/forja/<data>/` com o MANIFESTO de backups, `legoloop forjar --verificar` e o mesmo resumo no chat quando um agente instala.
Os quatro gatilhos da Forja produzem o mesmo RELATORIO.
O Processo Forja passa em `legoloop validate` e `legoloop fechar`.
A matriz de CI prova cada canal em Windows, Linux e macOS.
Um orçamento de tokens por harness, medido no CI, reprova a Forja se um harness passar do teto.

## restricoes_inegociaveis

Nunca apagar nem sobrescrever arquivo do usuário sem backup restaurável.
Nunca imprimir nem gravar segredos.
Nunca pedir privilégio de administrador sem aviso.
O script baixa o binário só de origem confiável e confere o SHA-256 antes de executar.
Todo canal oferece `--dry-run`.
O núcleo continua em Go, só com a biblioteca padrão.
Commit e push só quando o usuário pedir.

## fora_de_escopo

Nada. O escopo é total, sem onda posterior.

## escopo_expandido (sugestões aceitas pelo usuário)

1. Matriz de CI em Windows, Linux e macOS com contêiner limpo.
2. Orçamento de tokens por harness como gate de CI.
3. Detecção de deriva: `forjar --verificar` agendado ou por hook após instalar plugin.
4. `forjar --desfazer` e `desinstalar` por canal, lendo os MANIFESTOS de backup.
5. Binários assinados (Authenticode no MSI, `cosign`), SBOM e SHA-256 publicado.
6. Perfis de contexto escolhidos na Forja: mínimo, padrão e gateway.
7. Canal de atualização (`legoloop atualizar`, estável e beta) sobre a mesma Forja.
8. Épico de portar para Go o Python do caminho de instalação, deixando em Python só o opcional (vídeo e Whisper).

## ordem_de_implementacao (aprovada pelo usuário)

A ordem segue dependências, para não gerar retrabalho nem débito. Perfis e `--desfazer` pertencem ao motor da Forja.

1. Contrato (INTENT, SPEC) e linha de base medida.
2. CI nos três SOs (Windows, Linux, macOS), para todo ticket nascer testado.
3. Portar para Go o Python do caminho de instalação.
4. Motor da Forja: perfis, SSOT com bloco gerenciado, relatório, `--verificar` e `--desfazer`.
5. `legoloop instalar` chamando a Forja, com `--dry-run`.
6. Gate de orçamento de tokens por harness no CI.
7. Canonizar o Processo Forja (`fechar`, skill no Cofre, `/legoloop-forjar`).
8. Release reproduzível com assinatura, SBOM e SHA-256.
9. Canais: script de terminal primeiro; depois MSI, winget, scoop, brew, git clone, agente e MCP.
10. Detecção de deriva.
11. Canal de atualização (estável e beta).

## execucao_pelo_swarm (autorizada pelo usuário)

O Swarm executa o plano sob as seis condições da lei 4 do `AGENTS.md`:
1. O contexto de cada worker contém só os artefatos do ticket.
2. O ticket lista as ferramentas permitidas.
3. Um worker executa um ticket.
4. A saída fica em disco e passa pelo gate mecânico.
5. O worker não cria sub-workers.
6. O `write_boundary` é `strict`.

## forma_do_plano

Macro plano com iniciativas, épicos, micro tickets e nano tickets, de qualquer tamanho, sem nada de fora.
