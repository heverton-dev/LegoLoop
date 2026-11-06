---
name: legoloop-matriz
description: Use quando o usuário quiser GERAR, CONSULTAR ou GERENCIAR a matriz de peças (qualidade, perfil de uso, vocabulário e defasagem), via /legoloop-matriz. NÃO use para fabricar processos novos, remover ou compor processos.
consumes: Artifact:Vault
produces: Artifact:Matrix
artifact: ~/.aidd/canonical/registry_matrix.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-matriz — Gestão da Matriz de Peças do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja gerar a matriz de peças do cofre, auditar a qualidade e tiers de reputação, listar peças por natureza/etapa/caso, aprovar novos termos de vocabulário ou confirmar perfis de uso.  
**Quando NÃO usar:** fabricar novos processos do zero, remover ou compor processos existentes.

Esta skill conduz a interação no chat e o gate humano. **Toda a lógica reside no comando `legoloop matriz`**: geração atômica, cálculo incremental por hash, pontuação e perfis. Não reproduza a catalogação em prompt.

* **Consumes (Entradas Obrigatórias):**
  - comando solicitado pelo usuário (`gerar`, `ver`, `status`, `vocab`, `perfil`)
* **Produces (Saídas Mecânicas):**
  - `~/.aidd/canonical/registry_matrix.json` e `~/.aidd/canonical/matrix/*.json`
* **Target Files (escopo de escrita permitido):**
  - nenhum arquivo é manipulado diretamente pela skill; apenas o comando `legoloop matriz` grava

### Procedimento (passos determinísticos)

1. **Consulta e Visão:**
   - Para inspecionar peças catalogadas, execute `legoloop matriz ver [--tipo <tipo>] [--tier <tier>] [--natureza <nat>] [--etapa <etapa>] [--caso <caso>]`.
2. **Status e Defasagem:**
   - Para checar peças novas ou alteradas sem reprocessar, execute `legoloop matriz status`.
3. **Geração da Matriz:**
   - Para gerar sem gastar modelo de IA: `legoloop matriz gerar --sem-cerebro`.
   - Para catalogação completa com cérebro e stop-gate: apresente prévia com `legoloop matriz status`, peça confirmação explícita ao usuário e execute `legoloop matriz gerar`.
4. **Vocabulário:**
   - Listar termos: `legoloop matriz vocab listar`.
   - Listar termos pendentes: `legoloop matriz vocab pendentes`.
   - Aprovar termo: execute `legoloop matriz vocab aprovar <termo>`.
5. **Perfil de Peça:**
   - Definir perfil: `legoloop matriz perfil definir <peça> --etapa <e> --natureza <n> --caso <c>`.
   - Confirmar perfil: `legoloop matriz perfil confirmar <peça>`.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** edite `registry_matrix.json` ou `matrix/*.json` manualmente; utilize exclusivamente `legoloop matriz`.
> 2. **NUNCA** rode `matriz gerar --forcar` sem consentimento explícito do usuário.
> 3. **NUNCA** invente termos fora do vocabulário fechado; termos desconhecidos devem passar por proposta pendente e aprovação.
> 4. **NUNCA** colete credenciais ou execute scripts não autorizados.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop matriz status`
- `legoloop matriz ver --tier ouro`
- `legoloop matriz vocab listar`
