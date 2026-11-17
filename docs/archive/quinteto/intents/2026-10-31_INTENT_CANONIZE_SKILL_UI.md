# 🎯 INTENT: Canonização Sob Demanda sem Fricção (Slash Command & Studio UI)

> **Data:** 31/10/2026  
> **Status:** Proposta de Intenção (Aguardando Aprovação para Spec)  
> **Interfaces:** Chat (`/canonize`) e Looping Studio (`:7878`)  

---

## 1. Motivação e Problema Observado

Executar a canonização via linha de comando exige conhecimento de flags complexas (`-root`, `-vault`, `-workers`, `-apply`). Além disso, projetos em discos externos (`D:\`, backups) ficam isolados da desduplicação por falta de uma interface simples e sem barreiras.

---

## 2. Intenção Primária

Provisionar duas interfaces amigáveis para acionar a canonização do LegoLoop sob demanda:
1. **Via Chat (Slash Command `/canonize`):** Permite disparar a varredura e a aplicação em qualquer pasta ou drive apenas digitando `/canonize [caminho]`.
2. **Via Studio Visual (`:7878`):** Adicionar um painel/modal de **Cofre & Volumes** onde o desenvolvedor visualiza os discos conectados e clica em **"Canonizar"** com confirmação prévia em 1 clique.

---

## 3. Escopo e Entregáveis

* **Entregável A (Skill / Slash Command):** Atualizar o `SKILL.md` mestre do LegoLoop registrando o comando `/canonize` e o protocolo conversacional de Dry-Run ➔ Confirmação ➔ Apply.
* **Entregável B (Studio UI):** Adicionar no topo do Studio ou na sidebar um botão/modal de auditoria de volumes do SO.
* **Entregável C (Script Auxiliar / CLI):** Garantir que `legoloop canonize` aceite aliases simples como `legoloop canonize D:\`.

---

## 4. Critérios de Sucesso

- [ ] Digitar `/canonize D:\` no chat aciona automaticamente o dry-run e solicita confirmação simples para o apply.
- [ ] O Studio exibe o estado do cofre canônico e a lista de volumes auditáveis.
- [ ] Zero quebra de contratos de testes (`go test ./...` 100% verde).
