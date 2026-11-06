# 📘 MANUAL OPERACIONAL: FUNÇÃO REFINAR NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Substituir peças, casos de uso ou títulos de etapas em processos existentes  
> **Fontes de verdade:** `internal/refinar/` · `internal/cli/refinar.go` · `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`

---

## 1. Visão Geral

O comando **`legoloop refinar`** permite alterar cirurgicamente uma etapa (`e01`, `e02`...) de um processo já existente no cofre. Ele recalcula a projeção da Tríade de governança, preserva o histórico de versões anteriores, sobe a versão minor semântica (ex.: `1.0.0` para `1.1.0`) e mantém backup atômico em `.legoloop-bak` com rollback automático.

---

## 2. Comandos e Uso

### 2.1 Refinamento Interativo com Prévia
```bash
# Apresenta o diff projetado e pede confirmação humana [s/N]
legoloop refinar legoloop-extrator --etapa e01 --peca nova-skill-coleta
```

### 2.2 Refinamento com Aplicação Direta
```bash
# Aplica a alteração com consentimento explícito prévio
legoloop refinar legoloop-extrator --etapa e01 --peca nova-skill-coleta --yes
```

### 2.3 Alteração de Casos de Uso ou Título
```bash
legoloop refinar legoloop-software --etapa e02 --caso tdd --caso qualidade --titulo "Construção com TDD" --yes
```

### 2.4 Saída Estruturada JSON
```bash
legoloop refinar legoloop-extrator --etapa e01 --peca nova-skill-coleta --yes --json
```

---

## 3. Invariantes
- A etapa alvo deve existir no processo; caso contrário, a execução falha com exit code 1.
- O histórico de alterações é preservado no arquivo de processo (`loops/<nome>.processo.json`).
- Qualquer falha no sistema de arquivos reverte automaticamente o processo para a versão anterior.
