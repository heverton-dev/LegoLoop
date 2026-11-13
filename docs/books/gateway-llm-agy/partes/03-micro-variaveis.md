# PARTE III — Variáveis, protocolos e precedência

## 9. As seis variáveis

| Variável                      | Obrigatória     | Função                               |
| :---------------------------- | :-------------: | :----------------------------------- |
| `AGY_LLM_GATEWAY_URL`              | Sim             | URL HTTPS do gateway                 |
| `AGY_LLM_GATEWAY_API_KEY`           | Sim             | Chave que o gateway exige            |
| `AGY_LLM_GATEWAY_WIRE_PROTOCOL`     | Não             | `genai` (padrão) ou `openai`             |
| `AGY_LLM_GATEWAY_CA_CERT`           | Não             | Caminho do pacote de certificados raiz |
| `AGY_LLM_GATEWAY_MODELS`           | Não             | Modelos permitidos, separados por vírgula |
| `AGY_LLM_GATEWAY_HEADERS`          | Não             | Cabeçalhos extras no formato `CHAVE: VALOR` |

A chave segue como token no cabeçalho `Authorization` de cada chamada.

## 10. Os dois protocolos

| Protocolo      | Quando usar                           | Observação                   |
| :------------- | :------------------------------------ | :--------------------------- |
| `genai`          | Gateway que fala a API do Gemini      | Padrão do `agy`                |
| `openai`         | Gateway ou provedor com API Chat Completions | Não funciona no servidor ACP |

O servidor ACP é o que o JetBrains e o Zed usam. Nesses editores, só `genai` funciona.

## 11. Qual chave serve

A chave é a do gateway escolhido. Não há "chave Google" obrigatória.

| Destino                        | Protocolo      | URL                                    | Chave             |
| :----------------------------- | :------------: | :------------------------------------- | :---------------- |
| 9Router da empresa             | `openai`         | `https://roteador.signa-docs.com/v1`     | chave do 9Router  |
| Provedor compatível com OpenAI | `openai`         | URL base do provedor                   | chave do provedor |
| Gateway que fala Gemini        | `genai`          | URL do gateway                         | chave do gateway  |
| Roteador próprio               | `openai`         | URL do roteador                        | chave do roteador |

Um provedor que não fala `genai` nem `openai` não funciona direto. Coloque um roteador na frente dele.

## 12. Precedência

1. O `admin_settings.json` do administrador vence a variável de ambiente. Em máquina gerenciada, a variável pode não valer.
2. No JetBrains e no Zed, defina `"auth": { "type": "gateway" }` em `~/.gemini/antigravity-acp/settings.json` e ponha as variáveis no campo `env` do agente.
3. Programas abertos pelo menu Iniciar ou pelo Dock não leem o perfil do shell.
4. O `agy` na web não aceita gateway.

## 13. Terminal novo, sempre

Variável de usuário só chega aos processos que nascem depois dela.
Todo terminal aberto antes da mudança mantém o ambiente antigo.
Esse é o motivo mais comum de "mudei e não funcionou".
