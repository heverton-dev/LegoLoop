---
name: [nome-da-skill-kebab-case]
description: [Descrição concisa e unívoca de quando usar e quando NÃO usar esta peça]
commands:
  - "/[comando-principal]"
---

# 🧩 [Nome da Skill] — [Propósito Central de Engenharia]

> **Padrão Ouro LegoLoop (Score >= 0.85)**  
> **Autor:** [Nome do Autor / Reputação]  
> **Versão:** 1.0.0  

---

## 1. Contexto & Contratos de I/O

Esta peça opera estritamente sob os seguintes contratos de entrada e saída:

* **Consumes (Entradas Obrigatórias):**
  - `docs/specs/YYYY-MM-DD_SPEC_*.md` (ou outro artefato contratado)
* **Produces (Saídas Mecânicas):**
  - `docs/plans/YYYY-MM-DD_PLAN_*.md` (ou outro artefato persistido)

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
1. **NUNCA** gere funções incompletas com `TODO`, `pass`, `// implementar depois` ou mocks simulados (Zero Stubs).
2. **NUNCA** avance para o próximo passo sem a validação mecânica de terminal aprovada (`exit code == 0`).
3. **NUNCA** toque em arquivos fora do escopo declarado em Target Files.

---

## 3. Comandos de Verificação & Testes de Terminal

Toda alteração produzida por esta peça deve ser imediatamente verificada pelos comandos executáveis abaixo:

```bash
# Validação de compilação / Build
go build -o bin/legoloop.exe ./cmd/legoloop

# Execução de testes focados
go test -v ./internal/[modulo]/...

# Linter e conformidade estática
golangci-lint run ./internal/[modulo]/...
```

---

## 4. Critério de Parada (Stopping Checklist)

A execução desta peça é considerada concluída com sucesso quando:
- [ ] O comando de teste unitário passou com exit code 0.
- [ ] O artefato produzido existe fisicamente no disco e possui tamanho > 0 bytes.
- [ ] Zero resíduos ou arquivos temporários órfãos deixados no sistema.
- [ ] Stop-Gate humano aprovado formalmente caso exigido pelo modo de supervisão.
