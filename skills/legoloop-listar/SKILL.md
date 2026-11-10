---
name: legoloop-listar
description: Use quando o usuário quiser LISTAR e INSPECIONAR os processos do cofre (versão, natureza, harnesses, vínculos, lacunas e sugestões), via /legoloop-listar. NÃO use para fabricar, remover ou compor processos.
consumes: Artifact:Vault
produces: Artifact:Report
artifact: stdout
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-listar — Inspeção de Processos do LegoLoop

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** o usuário deseja inspecionar todos os processos fabricados no cofre, verificar versões, natureza declarada, harnesses onde estão vinculados, processos compostos (subprocessos), dependências de uso e peças faltantes (lacunas).  
**Quando NÃO usar:** para fabricar processos novos, alterar etapas, remover ou compor processos.

Esta skill apenas apresenta dados consultados. **Toda a lógica reside no comando `legoloop listar`**: inspeção do cofre, mapeamento de vínculos e sugestão de resolução de lacunas.

* **Consumes (Entradas Obrigatórias):**
  - diretório do cofre (opcional, padrão `~/.aidd/canonical`)
* **Produces (Saídas Mecânicas):**
  - tabela formatada ou JSON com o inventário dos processos
* **Target Files (escopo de escrita permitido):**
  - nenhum arquivo é alterado pelo listar

### Procedimento (passos determinísticos)

1. **Listagem em Texto:**
   - Execute `legoloop listar`.
   - Apresente ao usuário a lista com nome, versão, natureza, harnesses onde o processo está ativo e processos que o utilizam (`Usado por`).
2. **Listagem em JSON:**
   - Para integrações ou ferramentas estruturadas, execute `legoloop listar --json`.
3. **Tratamento de Lacunas:**
   - Se o comando exibir processos com `LACUNA`, apresente as sugestões de `perfil definir` informadas pelo comando (peças candidatas do cofre).
   - Se o usuário desejar preencher a lacuna, ofereça rodar o comando sugerido: `legoloop matriz perfil definir <peça> --caso <caso>`.

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** filtre ou omita processos que estão em estado de lacuna ou versão desatualizada sem que o usuário peça.
> 2. **NUNCA** tente modificar arquivos do cofre durante a listagem.
> 3. **NUNCA** execute comandos de alteração sem autorização explícita do usuário.

---

## 3. Comandos de Verificação & Testes de Terminal

- `legoloop listar`
- `legoloop listar --json`
