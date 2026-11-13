# PARTE VI — Apêndices

## Apêndice A — Glossário

| Termo                        | Significado                                         |
| :--------------------------- | :-------------------------------------------------- |
| `agy`                          | Antigravity CLI, o harness de linha de comando do Google |
| Gateway                      | Endereço HTTPS que recebe as chamadas de modelo no lugar da conta Google |
| Protocolo (`WIRE_PROTOCOL`)     | Formato das chamadas: `genai` ou `openai`               |
| Combo                        | Conjunto nomeado de modelos no 9Router              |
| Rodízio                      | Troca de conta Google por código para somar cotas. Proibido |
| Terminal novo                | Terminal aberto depois da mudança das variáveis     |
| `leak-scan`                    | Hook que bloqueia saída com segredo                 |

## Apêndice B — Relatório de conformidade

| Fato                                   | Prova                                       |
| :------------------------------------- | :------------------------------------------ |
| As seis variáveis `AGY_LLM_GATEWAY_*` são documentadas | Google Cloud, "Connect Antigravity to an Enterprise Gateway" |
| `gw agy 9router` grava as três variáveis | `pkg/driver/driver.go`, `SetRoute`              |
| `gw agy off` apaga as três variáveis     | `pkg/driver/driver.go`, ramo `RouteOff`         |
| Gateway do 9Router respondeu           | `/v1/models` com HTTP 200 em 13/11/2026       |
| O `agy` respondeu pelo 9Router           | confirmação do usuário em terminal novo     |
| Comandos de rodízio estão bloqueados   | execução de `next`, `dev03`, `probe --all` e `sync` com código de saída 2 |
| O Google aprova gateway de terceiros   | **Não verificado**                              |
| `probe --all` causou os bloqueios        | **Não verificado** (hipótese)                   |
| `gw agy list` não usa rede               | leitura de `pkg/driver/list.go`, sem chamadas HTTP |

## Apêndice C — Referência rápida

| Quero                             | Faço                                           |
| :-------------------------------- | :--------------------------------------------- |
| Ligar o 9Router                   | `gw agy 9router`, terminal novo                  |
| Voltar ao `agy` puro                | `gw agy off`, terminal novo                      |
| Testar sem persistir              | variáveis `$env:` na sessão                      |
| Conferir sem mostrar a chave      | `definida: True/False`                           |
| Trocar a chave                    | janela do sistema ou pedido ao agente com variável |

## Apêndice D — Fontes

- Google Cloud, "Connect Antigravity to an Enterprise Gateway": https://docs.cloud.google.com/gemini/enterprise/docs/ai-developer-tools/connect-enterprise-gateway
- Google Cloud, "Configure an enterprise gateway for Google Antigravity": https://docs.cloud.google.com/gemini/enterprise/docs/ai-developer-tools/deploy-enterprise-gateway
- Portkey, "Google Antigravity" (guia de terceiros): https://portkey.ai/docs/integrations/libraries/antigravity
- Manual operacional: `docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md`
