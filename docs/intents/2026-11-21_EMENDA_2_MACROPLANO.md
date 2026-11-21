# EMENDA 2 — Macroplano Distribuição e Forja do LegoLoop

Data: 2026-11-21. Status: UNDER_REVIEW (aprovação junto com a INTENT da I04).
Motivo: a auditoria da espinha dorsal (`docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md`) achou 14 lacunas que as iniciativas I01, I02 e I03 herdariam. As INTENTs aprovadas dessas três não mudam; esta Emenda muda a ordem e registra as revisões que a SPEC e o PLAN de cada uma devem receber.

| Ordem | Iniciativa | INTENT | Situação |
| :--- | :--- | :--- | :--- |
| 1 | I04 — Espinha dorsal | `2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md` | a executar primeiro |
| 2 | I01 — Distribuição e Forja | `2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md` | congelada; revisão abaixo |
| 3 | I02 — Especificação e planejamento | `2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md` | congelada; revisão abaixo |
| 4 | I03 — Despachante de executores | `2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md` | congelada; revisão abaixo |

Revisões que as SPEC e os PLAN congelados devem receber depois da I04:
- **I01:** o Épico 02 (CI nos três SOs) sai, porque a I04 o entrega. A Forja passa a depender de G01 (Matriz por hash), G02 (descoberta de peças) e G05 (fachada única de caminhos) resolvidos antes do TICK-01.2 (renomear o módulo).
- **I02:** o vocabulário de tipo de nó usa a tabela única `unit_task` e `unit_no` da I04. Os Processos novos nascem 5 de 5 e registrados no Registro da Espinha.
- **I03:** a validação de executor passa de "valor conhecido" para "coerente com o tipo do nó", e o `dag.Parse` deriva o valor ausente pelo tipo. O valor do Swarm continua `native_swarm`. A cobertura do `board` (65,3%) sobe antes da I03.

Regras:
- A execução da I01, da I02 e da I03 continua congelada até o Stop-Gate final da I04.
- Cada revisão acima vira uma SPEC e um PLAN novos (nova versão do arquivo), passando pela trava `upstream_ok`. Nenhuma INTENT aprovada é editada.
