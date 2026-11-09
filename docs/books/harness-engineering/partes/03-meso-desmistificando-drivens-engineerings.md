# Capítulo 3: O Labirinto dos Acrônimos — Desmistificando os "Drivens" e "Engineerings"

A indústria de tecnologia adora inventar siglas pomposas para vender cursos e posts no Twitter. Quando o desenvolvedor se depara com TDD, BDD, AIDD, Prompt Engineering, Context Engineering e Graph Engineering, a cabeça entra em curto-circuito.

Vamos cortar o ruído de marketing e entender a mecânica prática de cada peça.

---

## 1. TDD Real vs. TDD Falso em Tempos de IA

O Test-Driven Development (Desenvolvimento Orientado a Testes) tradicional diz: *Escreva o teste antes do código de produção*. 

Na era da IA, o TDD ganhou um propósito completamente novo: **ele é o único detector de mentiras mecânico que você tem contra o modelo.**

Quando você pede para um modelo: *"Implemente a função de renderização de vídeo"*, o modelo escreve 200 linhas de código bonito, com nomes de variáveis elegantes. Mas se você não tiver um teste, você é obrigado a:
1. Abrir o código;
2. Ler linha por linha mentalmente;
3. Tentar descobrir onde ele inventou um método que não existe na biblioteca.

Você virou um compilador humano. Isso é exaustivo.

No **TDD com Harness**:
1. Você define o contrato da função (`Input -> Expected Output`).
2. O modelo escreve o teste unitário.
3. O Harness roda `go test` e **exige que o teste falhe (RED)**.
4. O modelo implementa a função de produção.
5. O Harness roda `go test` e **exige que passe (GREEN)**.
6. Se falhar, o erro do compilador volta automaticamente para o modelo corrigir. Você nem precisou olhar o código.

---

## 2. BDD (Behavior-Driven Development) como Contrato Humano-Máquina

O BDD não é sobre ferramentas como Cucumber; é sobre **precisão de fronteira**.

Quando você diz para a IA: *"Quero que o vídeo fique com um visual cartoon legal"*, você usou uma palavra subjetiva (*"legal"*). A IA não sabe o que é legal. Ela vai alucinar.

No BDD, a exigência é binária (Dado que / Quando / Então):
* **Dado que:** Uma lista de 8 cenas com duração de 25 segundos é enviada;
* **Quando:** O renderizador Three.js for invocado com resolução vertical (1080x1920 a 30fps);
* **Então:** Deve produzir um arquivo `.mp4` com codec H.264, áudio AAC normalizado em -14 LUFS e trilha de fundo exatamente a -6 dB em relação à voz.

Não existe margem para debate. O teste afere se o arquivo tem 1080x1920 e se o áudio está em -14 LUFS. Ou atende ou não atende.

---

## 3. AIDD (AI-Driven Development): A Definição Real

AI-Driven Development **NÃO É** abrir o VS Code com Cursor ou Windsurf e pedir: *"Gera um sistema de autenticação pra mim"*. Isso é *AI-Assisted Vibe Coding* e funciona bem para protótipos de fim de semana, mas colapsa em sistemas de produção.

O verdadeiro **AI-Driven Development** é uma disciplina de arquitetura onde:
1. **O Humano é o Arquiteto de Fronteiras:** Você desenha o DAG de execução, os schemas de dados e os Stop-Gates de aprovação.
2. **O Harness é o Fiscal de Linha de Montagem:** Ele executa testes, mutações, checagens estáticas e sandboxing.
3. **A IA é o Motor de Combustão:** Ela preenche a lógica pura de transformação dentro dos limites que o Harness impõe.

---

## 4. Prompt Engineering vs. Context Engineering

Esqueça o "Prompt Engineering" clássico. Ficar ajustando palavras como *"pense passo a passo"* ou *"você é um programador sênior premiado"* tem impacto marginal em tarefas complexas de engenharia.

O que realmente dita o sucesso é a **Engenharia de Contexto (Context Engineering)**:
* **Curadoria Estrita:** Enviar para o modelo apenas os 50 tokens essenciais para a decisão, em vez de 50.000 tokens de arquivos que ele não precisa ler.
* **Graph-First Memory:** Usar o grafo de símbolos do repositório (`codebase-memory-mcp`) para saber exatamente quais funções chamam a função atual, em vez de fazer buscas cegas por texto (`grep` ou `glob`) que trazem 100 arquivos falsos positivos.
* **Ephemeral Rotation:** Destruir a sessão do chat e recomeçar a cada mudança de fase, preservando o estado factual apenas em artefatos no disco (`docs/intents/`, `docs/specs/`, `docs/plans/`).
