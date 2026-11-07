# PARTE II — Meso: Modos de Superfície e o Loop Mestre

## 4. A ambição depende da tela, não do produto

Um mesmo produto tem telas com objetivos muito diferentes. A página de preços precisa convencer; o painel de configurações precisa ser rápido de usar; a documentação precisa ser lida sem cansaço. O Impeccable mede a ambição visual pelo **modo da superfície**, e não pelo produto como um todo.

| Modo       | O visitante precisa     | Exemplos                    | Prioridade de design             |
| :--------- | :---------------------- | :-------------------------- | :------------------------------- |
| Operate    | Completar uma tarefa    | Dashboards, admin, editores | Escaneabilidade e contenção      |
| Persuade   | Decidir e agir          | Landing pages, preços       | Hierarquia de valor e conversão  |
| Read       | Compreender um conteúdo | Docs, manuais, changelogs   | Ritmo de leitura sem distração   |
| Experience | Imergir na própria obra | Portfólios, galerias        | A interface recua, a obra lidera |

### Como o modo muda as decisões

No modo **Operate**, a ambição aparece como precisão: densidade coerente, alinhamentos exatos, comportamentos que o usuário já espera da plataforma. Um painel administrativo com tipografia expressiva e animações de entrada em cada seção está no modo errado.

No modo **Persuade**, a tipografia pode ser expressiva e a estética pode ser imersiva, porque o objetivo é mover o visitante para uma decisão. É o único modo, ao lado de Experience, em que o comando `overdrive` é permitido.

No modo **Read**, o contraste é descansado e a hierarquia tipográfica faz quase todo o trabalho. O piso de qualidade fixa a medida do corpo de texto entre 65 e 75 caracteres por linha.

No modo **Experience**, a interface sai do caminho. Desde a primeira dobra, quem lidera é a obra exibida: a fotografia, o projeto, a peça interativa.

```{=typst}
#painel("Regra de classificação")[
  Classifique cada tela isoladamente. Um produto SaaS típico tem a landing em Persuade, o app em Operate e a central de ajuda em Read. Aplicar a mesma ambição às três é o erro mais comum de quem começa no método.
]
```

## 5. O loop mestre de cinco estágios

O fluxo completo do Impeccable atravessa cinco estágios. Entre eles existem quatro gates binários de parada obrigatória, chamados P1 a P4. Binário quer dizer que o gate só tem duas respostas: passa ou não passa.

```{=typst}
#esteira((
  ("1. Contexto e Auditoria", "PRODUCT.md, DESIGN.md, critique, audit"),
  ("2. Direção Visual", "classificar modo, shape, preservar ou substituir"),
  ("3. Execução Cirúrgica", "craft-floor.md, layout, typeset, colorize, clarify"),
  ("4. Verificação em Lote", "captura mobile 375px e desktop, polish, revisão"),
  ("5. Documentação e Entrega", "document, extract, proposta de PR"),
), gates: ("P1: achados e severidade", "P2: direção visual aprovada", "P3: antes e depois por item", "P4: veredito VERIFIED"))
```

### Estágio 1: Contexto e Auditoria

O agente carrega os dois documentos de contexto durável, `PRODUCT.md` e `DESIGN.md`, e então roda duas avaliações: `critique`, uma revisão heurística de experiência do usuário, e `audit`, uma checagem técnica de acessibilidade, performance e viewports.

**Gate P1.** O estágio só termina quando os achados estão listados com severidade. Sem severidade, não há como decidir o que corrigir primeiro.

### Estágio 2: Direção Visual

Com o diagnóstico em mãos, o agente classifica o modo da superfície, planeja a experiência antes do código com `shape` e decide o destino do trabalho novo: preservar, expandir ou substituir o que existe.

**Gate P2.** A direção visual precisa de aprovação humana explícita. É o gate mais barato do loop: errar a direção aqui custa uma conversa; errar depois custa uma reimplementação.

### Estágio 3: Execução Cirúrgica

A primeira ação do estágio é a leitura obrigatória de `craft-floor.md`. Depois vêm, nesta ordem, espaçamento e ritmo (`layout`), hierarquia tipográfica (`typeset`), paleta estratégica (`colorize`) e, por fim, textos e casos de borda (`harden` e `clarify`).

A ordem não é arbitrária. Cor aplicada antes de a hierarquia estar resolvida tende a ser usada como muleta para criar destaque que a tipografia deveria criar.

**Gate P3.** Cada item alterado precisa de um antes e depois. Isso torna a revisão do gate objetiva: o revisor compara pares, não opina sobre o todo.

### Estágio 4: Verificação em Lote

Aqui o pilar da verificação delimitada vira prática. Uma captura em lote cobre mobile em 375 pixels e desktop. O comando `polish` faz o acabamento fino e o agente revisor de acabamento julga o resultado.

**Gate P4.** O revisor emite um veredito. Só `VERIFIED` libera o estágio seguinte.

### Estágio 5: Documentação e Entrega

O comando `document` atualiza o `DESIGN.md` a partir do código real, e não das intenções. O comando `extract` isola tokens e componentes para o design system. O estágio termina com uma proposta de entrega ou um pull request.

```{=typst}
#frase-destaque[O DESIGN.md descreve o que foi entregue, não o que se pretendia entregar.]
```

## 6. Por que os gates ficam onde ficam

Cada gate fica no ponto em que o próximo estágio ficaria caro se começasse com uma premissa errada.

| Gate | Fica antes de     | Custo evitado                         |
| :--- | :---------------- | :------------------------------------ |
| P1   | Decidir a direção | Direção escolhida sem diagnóstico     |
| P2   | Escrever código   | Reimplementar uma tela inteira        |
| P3   | Verificar em lote | Revisão subjetiva sem comparação      |
| P4   | Documentar        | Registrar no DESIGN.md algo reprovado |
