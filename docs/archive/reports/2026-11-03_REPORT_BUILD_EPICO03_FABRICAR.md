# 📊 REPORT: BUILD do Épico 03 (Studio: API, modal e plotagem no Canvas) — `SG-E03`

> **Data:** 2026-11-03  
> **Escopo:** `TICK-03.1` a `TICK-03.4` do PLAN `2026-11-03_PLAN_LEGOLOOP_FABRICAR.md`, mais a repetição limitada da SPEC aprovada junto com a autorização do `SG-E02`  
> **Estado:** 4 tickets `PASSED`; aguardando a **sua** homologação no navegador e a autorização do `SG-E03` para o Épico 04.  

---

## 1. Entregas

| Ticket | Entrega | Arquivos principais |
|:---|:---|:---|
| Pré-requisito (decisão do `SG-E02`) | **Repetição limitada** da SPEC: até 3 tentativas, devolvendo ao modelo a resposta rejeitada com o motivo; só repete quando o erro é de contrato (`specgen.ErrContract`), nunca em falha de transporte | `internal/fabricar/repair.go`, `internal/specgen/specgen.go` |
| `TICK-03.1` | `POST /api/fabricar/run` atrás do middleware `secure`: 405 para GET, 403 origem cruzada, 421 host inválido, 415 content-type, 400 nome/modelo inválido ou campo desconhecido, 409 duplicata, resposta sem caminhos absolutos | `internal/studio/fabricar_api.go`, `server.go` |
| `TICK-03.2` | Exclusão mútua (409), `context.WithoutCancel` com teto de 10 min, prazo de escrita estendido (o `WriteTimeout` de 30 s cortaria a resposta), 422/502/500 sem expor o erro bruto, `Assets.Refresh` após o sucesso | `fabricar_api.go` |
| `TICK-03.3` | Botão `+ Novo Processo` e `<dialog id="modal-fabricar">` com nome, categoria, modelo (só 9Router habilitado), entrevista via `/api/interview/*` e erros no próprio modal | `static/index.html`, `style.css`, `app.js`, **`fabricar.js`** |
| `TICK-03.4` | `fabricarPlotModel`: o Canvas planta os nós e arestas do loop gerado, o Inspetor abre com a Tríade e os três meios, a peça aparece em `Loops` sem recarregar | `static/fabricar.js`, `inspector.js`, `app.js` |

Peças compartilhadas criadas para não duplicar regra de segurança: `provider.TrustedRouterKey` (a chave só vai à URL confiável) e `fabricar.BrainFor` (o cérebro e o prazo), usadas pela CLI e pelo Studio.

## 2. Checklist do `SG-E03`

| Item | Resultado |
|:---|:---|
| `go test ./internal/studio/...` verde, incluindo emojis, diálogos nativos e CSP | **129 passando, 0 falhas, 0 pulados** (os testes de lógica JS rodam em `node` v26, não são pulados) |
| `go test -count=1 ./...` | Verde em 21 pacotes |
| `go vet ./...`, cross-vet `linux`/`darwin` | Limpos |
| `funclen` | Limpo em `fabricar`, `plangen`, `specgen`, `provider` e nos meus arquivos de `studio` e `cli` |
| `gofmt` | Meus arquivos limpos; `gofmt -l .` global lista os mesmos 22 preexistentes |
| `go.mod` e `go.sum` | Inalterados |
| Homologação visual no navegador com captura | **Feita por mim em ambiente isolado** (seção 3); `docs/reports/img/2026-11-03_E03_modal_plotagem.jpg`. **A homologação humana do PLAN continua pendente** |
| A peça nova aparece sob `Loops` sem recarregar | Confirmado no navegador (categoria `Loops` passou de 0 para 1, e o nome está na paleta) |
| Autorização explícita do usuário | **Pendente** |

## 3. Evidências além dos testes

* **Navegador real (Chrome), ambiente isolado:** Studio em porta alternativa, `USERPROFILE` e projeto temporários, 9Router **falso** local (LLM determinístico), varredura **real** dos discos. Fluxo executado: abrir o modal, digitar o nome, iniciar a entrevista (98%), clicar em **Fabricar processo**, esperar. Resultado: 8 arquivos gravados só nos diretórios temporários; o modal fechou; o Canvas plantou 5 nós e 4 arestas; o Inspetor mostrou INTENT, SPEC, PLAN, ROUTE, TEMPLATE, Skill, Loop e DAG com caminhos mascarados (`<projeto>`, `~`) e as 5 peças curadas; `Loops` passou a listar o processo. Console sem erros e sem violação de CSP. Erros aparecem **dentro do modal**: nome `!!!` ("Use ao menos uma letra ou número") e duplicata (409, "já existe um processo com este nome").
* **Defeito que só o navegador mostrou:** o painel da entrevista tem `hidden`, mas o meu `display: grid` o sobrepunha, então "Sua resposta", **Responder** e **Fabricar processo** apareciam antes da hora. Corrigi com regra `[hidden]` e teste (a mutação sem a regra falha).
* **Mutações (os testes detectam regressão):** resposta sem máscara de caminhos; decodificação aceitando campo desconhecido; `WithoutCancel` removido; prazo de escrita removido; chave enviada a qualquer URL; repetição acima de 3 tentativas; repetição também em falha de transporte; `alert(` e `innerHTML` em `fabricar.js`; aresta para dependência inexistente; nós sobrepostos. Todas falham e o código foi restaurado. Num caso (a máscara de caminhos), o meu **próprio teste estava errado**: a variante de barra escapada era um no-op e não detectaria vazamento; só a mutação revelou, e eu o corrigi.

## 4. Achados e desvios (leia antes de autorizar)

1. **A repetição limitada não foi validada com um LLM real.** Foi provada com LLM falso (corrige na 2ª tentativa; desiste na 3ª sem gravar nada; não repete falha de transporte). Não refiz a chamada ao seu 9Router remoto para não gastar mais; fica como verificação do `TICK-04.1` ou sua, se quiser.
2. **Defeito preexistente corrigido:** `/api/interview`, `/api/spec` e `/api/plan` do Studio montavam o cliente do 9Router **sem a chave** (o mesmo `401` da CLI). Agora usam `provider.TrustedRouterKey`. Sem isso, o botão "Iniciar Entrevista" não funcionaria contra um router remoto.
3. **A entrevista no Studio ainda usa `code-pro` fixo**, e não o modelo escolhido no dropdown (`/api/interview/start` só recebe `idea`). O modelo escolhido conduz a **Tríade**. Alinhar a entrevista exigiria um motor por modelo; deixei fora.
4. **Arquivo novo fora da lista do PLAN:** a lógica pura do modal e da plotagem está em `static/fabricar.js`, e não dentro de `app.js`. É o padrão do repositório (`palette.js`, `search.js`), e é o que permite testar em `node`. Também incluí `fabricar.js` na lista do `TestNoInnerHTMLInScripts`.
5. **A plotagem substitui o conteúdo do Canvas** (com `pushState`, então Ctrl+Z desfaz). Se houver uma esteira não salva, ela sai da tela. Pode ser melhor perguntar antes; fica à sua escolha.
6. **O modal não oferece `force`.** O servidor aceita o campo, mas a interface não o envia: em duplicata o usuário vê o erro e precisa trocar o nome (ou usar a CLI com `--force`).
7. **Respostas de erro da entrevista são texto puro** (`http.Error` do código preexistente), então o modal mostra a mensagem genérica por status. Não alterei esses handlers.
8. **Fui descuidado ao encerrar o ambiente de teste:** usei `taskkill /IM legoloop.exe`, que encerra **todo** processo com esse nome, inclusive um Studio seu que estivesse rodando. Se você tinha um aberto (por exemplo na porta 7878), ele foi encerrado; basta reabrir. Não verifiquei se havia um em execução, então não sei se algo foi afetado.
9. **Ainda abertos do relatório do Épico 02:** o ponto-a-ponto com LLM real; `-race` indisponível (CGO); 22 arquivos preexistentes sem `gofmt`.

## 5. Próximo passo

`SG-E03` pede (a) a **sua homologação** do modal no navegador, que é o gate humano do PLAN (abrir, nomear, escolher o modelo, entrevistar, plotar), e (b) a autorização explícita para o **Épico 04** (homologação mecânica final, smoke hermético em `scratch/smoke-fabricar`, relatório final e `SG-FINAL`). **Nada do Épico 04 foi iniciado.**
