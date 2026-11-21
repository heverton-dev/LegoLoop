# ESPELHO CANÔNICO — Iniciativa 4 (I04)

Sessão: `intv_1795277863681062100`. Data: 2026-11-21. Decisões tomadas pelo usuário uma a uma ("ACEITO", "[1]").
O motor concluiu com o modelo `sonnet` (96%) depois de parar em 85% com `code-fast`. Este registro é a versão consolidada e conferida contra as decisões.

1. **resultado_esperado:** espinha dorsal previsível, universal e determinística: Registro da Espinha conferido por `legoloop espinha verificar`; descoberta determinística de peças com Matriz atualizada por hash; executor derivado do tipo do nó e tabela única `unit_task` e `unit_no`; fachada única de caminhos; CI nos três SOs; mutação em duas camadas.
2. **quem_se_beneficia:** quem mantém o LegoLoop, o usuário nos três SOs e as iniciativas I01, I02 e I03.
3. **por_que_agora:** a auditoria achou 14 lacunas e as três iniciativas planejadas se apoiam em vocabulário, fachadas e descoberta que estão falhos.
4. **criterio_de_sucesso:** espinha verificar com 0 no repositório e exit 2 em fixtures com defeito; descoberta acha peça nova sem `matriz gerar`; DAGs validam com executor coerente; CI verde nos três SOs; 14 lacunas com Red e anticorpo; `-race` e ponta a ponta duas vezes idênticos; mutantes das regras 100% mortos e mutantes de código sem sobrevivente injustificado.
5. **restricoes_inegociaveis:** TDD com Red e exit 2 em defeito; descoberta determinística e somente leitura; Matriz por hash; executor pelo tipo; uma tabela e uma fachada; verificação no CI e no `fechar`; INTENTs da I01 a I03 intactas; zero débito.
6. **fora_de_escopo:** refatorar verbos de domínio além de classificá-los; mudar o formato do DAG; alterar as INTENTs da I01 a I03; determinismo do texto do modelo.

Decisões do usuário: (1) escopo das 14 lacunas mais a medição, execução antes da I01 a I03; (2) Registro da Espinha com gate; (3) descoberta por hash de conteúdo; (4) executor derivado do tipo e tabela única; (5) CI de três SOs trazido da I01; (6) provas de pronto e limites; (7) mutação em duas camadas e exit 2.
