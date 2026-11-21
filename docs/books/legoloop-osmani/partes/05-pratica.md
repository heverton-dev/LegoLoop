# PARTE V — Prática: Receitas de Uso

As receitas combinam sub-loops e playbooks para situações frequentes. Rotas, números e mensagens são ilustrativos.

## 21. Receita: a página de checkout está lenta

**Quando usar.** Usuários reclamam da demora no checkout. Ninguém mediu nada ainda.

1. Rode `performance-optimization` só para medir. Registre LCP, INP e CLS da rota `/checkout` em condições reproduzíveis.
2. Compare com os limiares do checklist. Suponha um LCP de 4,6 s: está na faixa ruim.
3. Fixe o orçamento em `CONSTRAINTS.md`, por exemplo LCP de até 2,5 s.
4. Ataque o maior gargalo, uma mudança por vez, com `incremental-implementation`.
5. Meça de novo depois de cada mudança.
6. Rode `browser-testing-with-devtools` para confirmar no navegador real que nada quebrou visualmente.

| Mudança (ilustrativa)                          | LCP antes        | LCP depois       |
| :--------------------------------------------- | :--------------- | :--------------- |
| Imagem principal com `fetchpriority=high`        | 4,6 s            | 3,4 s            |
| Remoção de biblioteca de carrossel             | 3,4 s            | 2,7 s            |
| Fonte com pré-carregamento                     | 2,7 s            | 2,3 s            |

**Erro comum.** Aplicar as três mudanças juntas e medir só no fim. O ganho total aparece, mas não se sabe qual mudança contribuiu, nem se alguma piorou o resultado.

## 22. Receita: um bug relatado por cliente

**Quando usar.** "O cupom de desconto é aplicado duas vezes."

1. Rode `debugging-and-error-recovery`. Reproduza e registre requisição, resposta e status. Sem reprodução, pare e peça dados.
2. Escreva um teste que falha pelo motivo exato: aplicar o cupom duas vezes resulta em desconto dobrado. Este é o gate vermelho.
3. Corrija a causa raiz com a menor mudança possível. O teste fica verde, junto com o restante da suíte da rota.
4. Se a correção toca autenticação, dados ou permissões, rode `security-and-hardening`.
5. Rode `code-review-and-quality` e `doubt-driven-development` sobre o diff.
6. Confirme que o sintoma original não se reproduz mais.

```{=typst}
#painel("Teste que falha pelo motivo certo")[
  Um teste que falha por erro de compilação ou por fixture quebrada não é o vermelho do Prove-It. O vermelho válido mostra o sintoma do bug: o valor esperado do desconto e o valor dobrado lado a lado na mensagem de falha.
]
```

## 23. Receita: auditar uma rota antes de mexer

**Quando usar.** A equipe vai assumir uma rota que não conhece e quer saber o que está errado antes de mudar qualquer coisa.

1. Playbook A, passo 1: resolva a rota em artefatos e rode os testes existentes para registrar a baseline.
2. Defina a barra de qualidade com `constraint-driven-development`, em modo leitura.
3. Varra um eixo de cada vez: correção e arquitetura, contrato, segurança, testes, UI em tela real, performance, visibilidade e complexidade.
4. Rode `doubt-driven-development` sobre os achados críticos e altos.
5. Converta os achados confirmados em tickets verticais com `planning-and-task-breakdown`.

Nenhum passo desta receita altera código. Achados de complexidade são apontados, não aplicados.

## 24. Receita: trocar um campo sem derrubar o sistema

**Quando usar.** A coluna `nome_completo` precisa virar `nome` e `sobrenome`, e há clientes em produção.

1. Rode `deprecation-and-migration` e planeje em expandir e contrair.
2. **Expandir:** crie os campos novos e passe a gravar nos dois formatos.
3. Migre os dados antigos em lotes, com testes de equivalência.
4. Passe a ler dos campos novos, mantendo os antigos.
5. **Contrair:** depois de um período de observação sem erros, remova o campo antigo.
6. Registre a decisão em uma ADR com `documentation-and-adrs`.

Cada passo é implantável e reversível isoladamente. É isso que permite a migração sem tempo de indisponibilidade.

## 25. Como escolher a receita

| Situação                                     | Receita          | Playbook da skill   |
| :------------------------------------------- | :--------------- | :------------------ |
| Lentidão percebida, sem número               | 21               | E, Otimizar         |
| Bug relatado com sintoma claro               | 22               | B, Corrigir         |
| Rota desconhecida, antes de mexer            | 23               | A, Auditar          |
| Mudança de esquema com clientes ativos       | 24               | H, Migrar           |
| Rota inteira, da auditoria à entrega         | —                | L, Ciclo completo   |
