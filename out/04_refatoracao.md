# 04 — Sugestao de Refatoracao Arquitetural

> **Alvo:** cmd/legoloop e internal/cli (comando fabricar)
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Padrao:** Clean Architecture / Injecao de Dependencias / Funcoes Curtas (<=35 linhas)

---

## 1. Principios Norteadores da Refatoracao

1. **Simetria de Comandos:** Todos os comandos CLI devem possuir exatamente a mesma assinatura e padrao de inicializacao em cli.Run.
2. **Isolamento de Canais de Saida:** Separacao estrita entre canal de dados estruturados (Stdout) e canal de telemetria e diagnostico (Stderr).
3. **Imutabilidade e Limite de Funcoes:** Manutencao estrita da regra de funclen <= 35 linhas sem violar a coesao.

---

## 2. Refatoracao Proposta para cmd/legoloop/main.go

Atualmente o arquivo main.go possui bifurcacao condicional ad-hoc:
- triagem e migrar-processos sao capturados antes de cli.Run.
- A unificacao proposta delega todos os comandos diretamente para cli.Run, simplificando a raiz do projeto e centralizando a logica de roteamento em internal/cli.

---

## 3. Isolamento da FlagSet de Fabricacao

Em internal/cli/fabricar.go, a criacao e parse de flags estao bem construidas, mas podem desacoplar a validacao de regras de negocio em um metodo de conveniencia puro:
- Extrair toRequest(o fabricarOpts, name string) (fabricar.Request, error).
- Garantir que toda mutacao de strings ou flags aconteca antes da resolucao do cerebro LLM.
- Preservar a injecao limpa de FabricarDeps, que hoje permite testes confiaveis sem dependencia de rede externa.
