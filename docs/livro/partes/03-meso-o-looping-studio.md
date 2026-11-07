# PARTE III — Meso: O Looping Studio Visual

## 13. O Ambiente Visual: Desmistificando Grafos de IA para Negócios

Durante décadas, a história da computação demonstrou repetidamente que tecnologias transformadoras só atingem seu verdadeiro potencial quando superam a barreira dos terminais de comando pretos e códigos inacessíveis para se tornarem interfaces visuais intuitivas. A planilha eletrônica fez isso com o cálculo financeiro, a interface gráfica de janelas fez isso com os sistemas operacionais e as ferramentas *node-based* (como n8n, Blender e Unreal Engine) fizeram isso com a automação e a computação gráfica.

O **Looping Studio** representa esse mesmo salto evolutivo para o universo dos sistemas agênticos de inteligência artificial.

Em vez de exigir que gestores, designers ou analistas de processos passem semanas aprendendo sintaxes obscuras de arquivos JSON ou comandos de terminal em Go, o Looping Studio transporta todo o rigor matemático dos grafos acíclicos para um canvas interativo, elegante e responsivo no navegador:

![Modal de Plotagem e Configuração de Nós no Looping Studio](docs/livro/img/2026-11-03_E03_modal_plotagem.jpg)

No Studio, uma esteira de inteligência artificial não é um bloco de código impenetrável: é um mapa visual claro onde cada nó executa apenas uma tarefa bem definida e cada aresta conecta o resultado de uma estação de trabalho à próxima.

---

## 14. Anatomia do Canvas: Nós, Arestas e Modais de Inspeção

Ao abrir o Looping Studio através do comando `legoloop studio`, o usuário é recebido por uma interface espacial moderna projetada para foco e produtividade:

### Os Quatro Elementos Estruturais da Interface

1. **Os Nós de Execução (Cards Modulares):** Cada retângulo colorido no canvas representa uma estação de trabalho atômica. O card exibe em destaque o nome da função, o selo de categoria (como engenharia de dados, testes unitários ou geração de relatórios), o tempo médio de processamento e seu estado em tempo real:
   * **Pendente (Cinza):** O nó aguarda a finalização das tarefas anteriores.
   * **Em Execução (Azul com pulso):** O nó está processando ativamente seus dados em contexto estéril.
   * **Aprovado (Verde):** A etapa foi concluída com sucesso e todos os testes passaram com código de saída zero.
   * **Parado em Gate (Âmbar):** A esteira pausou voluntariamente aguardando auditoria humana.
   * **Falha (Vermelho):** Um contrato ou teste unitário foi violado; a esteira congelou de forma segura.
2. **As Arestas Direcionadas (Cabos de Dados):** As linhas que interligam os nós estabelecem a precedência temporal e o fluxo de artefatos. Uma aresta só conduz sinal quando o nó de origem emitir seu artefato final devidamente validado.
3. **A Gaveta Lateral de Inspeção (Inspector Drawer):** Ao clicar sobre qualquer nó do fluxo, um painel lateral se abre revelando a radiografia completa daquela peça: o arquivo de especificação (`SPEC.md`), os critérios binários de aceitação, a pontuação de qualidade na Matriz de Reputação (Tier Ouro, Prata ou Bronze) e o consumo exato de tokens medido nas últimas execuções.
4. **O Catálogo de Peças (Toolbar Superior):** Permite arrastar novas peças do cofre diretamente para o canvas ou disparar a fábrica autônoma através do botão `+ Novo Processo`.

---

## 15. Construindo a Primeira Esteira do Zero Passo a Passo

Para vivenciar a simplicidade operacional do Looping Studio, vamos acompanhar o passo a passo da construção de uma esteira corporativa para processamento automático de contratos e faturas de fornecedores:

### Passo 1: Definição da Intenção Inicial
O operador clica no botão `+ Novo Processo` no canto superior direito. Um assistente embutido abre um modal amigável perguntando: *"Qual o objetivo de negócio desta esteira?"*. O operador digita: *"Ler faturas em PDF recebidas por e-mail, extrair CNPJ, valor total e itens discriminados, e validar se o CNPJ consta na base de fornecedores homologados"*.

### Passo 2: O Desdobramento da Tríade
Em segundo plano, o motor do LegoLoop gera os arquivos `INTENT.md` e `SPEC.md`, convertendo as regras em critérios puramente binários. Três nós surgem automaticamente dispostos no canvas:
* **Nó 1 (Extrator OCR):** Responsável por ler o PDF e emitir um JSON estruturado.
* **Nó 2 (Validador Fiscal):** Responsável por checar se a soma dos itens bate com o valor total da nota.
* **Nó 3 (Auditor de Cadastro):** Responsável por consultar a base local de fornecedores autorizados.

### Passo 3: Inserindo um Portão Humano (Stop-Gate)
O operador clica na aresta que conecta o *Validador Fiscal* ao *Auditor de Cadastro* e ativa a opção **Stop-Gate Humano**. Isso garante que faturas com divergência de valores ou notas acima de R$ 50.000,00 nunca serão aprovadas silenciosamente pela inteligência artificial: a esteira entrará em pausa obrigatória para que um analista financeiro revise os números com um clique.

### Passo 4: Otimização Contínua com um Clique
Caso uma versão aprimorada do nó extrator de PDF seja cadastrada no cofre central, o Studio exibe um sutil selo de recomendação amarelo. Ao clicar em `Refinar`, o modal de substituição atômica exibe o diff exato da melhoria e aplica a atualização sem quebrar as outras conexões:

![Modal de Confirmação Cirúrgica de Substituição de Peças no Studio](docs/livro/img/2026-11-03_E03_confirmar_substituicao.jpg)

---

## 16. Stop-Gates Humanos: O Equilíbrio entre IA e Decisão

O maior erro dos sistemas de automação que tentam implementar "piloto automático total" é desconsiderar que no mundo real existem decisões de negócios que carregam implicações éticas, legais e financeiras irreversíveis. Permitir que uma IA emita ordens de pagamento ou altere códigos críticos de infraestrutura sem supervisão humana é uma negligência inaceitável.

No LegoLoop, os **Stop-Gates Humanos (`SG-*`)** são cidadãos de primeira classe:

* **Ponto de Checagem Factual:** Quando um nó marcado como stop-gate atinge sua conclusão, a esteira congela seu estado e emite uma notificação visual no Studio.
* **Dossiê com Evidências Computáveis:** A interface não exibe apenas um botão genérico de "aprovar". Ela apresenta um dossiê com evidências frias: o arquivo gerado, os testes que foram executados, os limites avaliados e o impacto financeiro ou operacional daquela mudança.
* **Três Portas de Decisão:** O operador humano pode escolher entre:
  * **[1] Seguir:** Aprova a etapa e autoriza a continuidade automática da esteira.
  * **[2] Personalizar:** Insere um ajuste cirúrgico ou observação antes de liberar o próximo nó.
  * **[3] Cancelar / Reverter:** Aciona o rollback automático, restaurando o sistema ao estado inicial.

---

## 17. Concorrência e a Barreira de Junção (Join Barrier)

Um dos maiores desafios em esteiras de software e processos corporativos é o paralelismo: como permitir que três tarefas independentes executem ao mesmo tempo sem que uma atropele a outra?

Em chats convencionais de inteligência artificial, o paralelismo é impossível: a conversa é forçada a ser linear e sequencial, fazendo com que uma tarefa de 10 minutos atrase todo o restante do fluxo.

O LegoLoop implementa o mecanismo de **Barreira de Junção (*Join Barrier*)**:

* Se uma esteira precisa gerar simultaneamente o código de uma funcionalidade, os testes unitários e a documentação técnica, esses três nós são disparados em paralelo no canvas.
* Cada nó roda em seu próprio processo estéril e independente.
* A Barreira de Junção posiciona-se no ponto de convergência: ela monitora as três tarefas e aguarda deterministicamente até que todas tenham concluído com sucesso e código de saída zero.
* Se dois nós terminarem em 3 segundos, mas o terceiro nó falhar em um teste, a barreira não abre: a esteira trava com segurança, impedindo que um código sem testes correspondentes seja integrado ao sistema principal.

---

## 18. O Protocolo STE-PT: Transparência Radical na IA

Assistentes de inteligência artificial tradicionais oscilam entre dois extremos problemáticos: ou são verborrágicos demais — despejando longos parágrafos de conversa fiada que cansam o usuário —, ou agem como caixas-pretas silenciosas que alteram arquivos sem explicar o que fizeram.

Para estabelecer um padrão de cooperação transparente de alta fidelidade entre humanos e máquinas, o LegoLoop impõe o **Protocolo STE-PT (Situação, Teoria, Execução e Evidência)**.

### A Tríade Canônica de Comunicação

#### 1. Pré-Ação (O Que Será Feito Antes de Tocar no Sistema)
Antes de modificar um único arquivo ou disparar uma esteira, a IA deve apresentar formalmente suas intenções:
* **`**[S] Situação:**`** Diagnóstico factual do estado atual, sem introduções ou saudações.
* **`**[T] Teoria:**`** A explicação lógica e de engenharia pela qual aquela mudança específica resolverá o problema.
* **`**[E] Execução:**`** O plano de ação com a lista de arquivos afetados, a porta de rollback e as opções de decisão humana explícitas (`[1] seguir`, `[2] personalizar`, `[0] cancelar`).

#### 2. Pós-Ação (O Que Foi Feito e Quais as Provas)
Após concluir a operação, a IA apresenta o balanço de fechamento com evidências reais:
* **`**[S] Situação:**`** O novo estado do sistema após as intervenções.
* **`**[T] Tese:**`** A confirmação de que a hipótese inicial foi comprovada na prática.
* **`**[E] Evidência:**`** Provas computáveis e inquestionáveis: logs de testes que passaram, código de saída zero e verificação de conformidade de esquemas.

### O Protocolo em Ação: Diálogo Real no Terminal

Veja a diferença prática de comunicação quando um desenvolvedor solicita a criação de uma nova tabela de faturamento com o protocolo STE-PT ativo:

```text
**[S] Situação:**
A demanda exige adicionar a tabela 'faturas' com chave estrangeira para 'clientes', porém o arquivo de esquema atual ('internal/db/schema.sql') possui apenas a tabela 'usuários'.

**[T] Teoria:**
Criar uma nova migração SQL idempotente ('migrations/004_faturas.up.sql') garante que a alteração seja versionada e possa ser desfeita via rollback ('migrations/004_faturas.down.sql') sem perda de dados existentes.

**[E] Execução:**
· Ação: Criar migrações 004_faturas e rodar 'go test ./internal/db/...'
· Rollback: 'git checkout migrations/' e 'migrate down 1'
· Decisão: [1] seguir | [2] personalizar | [0] cancelar
```

Observe que a IA não fez rodeios, não disse "Olá, tudo bem? Claro, vou te ajudar com isso!", não executou nada silenciosamente e ofereceu uma porta clara de decisão humana. Após a aprovação do usuário, ela retorna com o fechamento rigoroso:

```text
**[S] Situação:**
Migração '004_faturas' aplicada com sucesso no banco de dados local de testes.

**[T] Tese:**
A tabela 'faturas' foi criada com índices de busca por CNPJ e restrição de integridade referencial ativa.

**[E] Evidência:**
· Testes: 'go test -count=1 ./internal/db/...' passou com sucesso (exit code: 0).
· Build: Binário compilado sem warnings.
· Verificação: Consulta 'SELECT count(*) FROM faturas' executada com retorno vazio (0 registros).
```

```{=typst}
#painel("Controle Centralizado com um Único Comando")[
  O comando 'legoloop stept on' ou 'legoloop stept off' permite ligar e desligar a governança do protocolo de transparência em todos os assistentes instalados no computador (Claude, Cursor, Windsurf, OpenCode) simultaneamente.
]
```

### Rastreabilidade Técnica da Parte III

* `internal/studio/server.go` — Servidor HTTP nativo do Looping Studio com barramento de eventos em tempo real.
* `internal/studio/assets/` — Canvas interativo em HTML5/JavaScript sem frameworks externos pesados.
* `internal/stept/stept.go` — Motor de auditoria, chaveamento e verificação de conformidade do protocolo STE-PT.
* `internal/dag/runner.go` — Executor topológico com suporte nativo a paralelismo e barreiras de junção.