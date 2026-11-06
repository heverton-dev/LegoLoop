# 📘 MANUAL OPERACIONAL: FUNÇÃO MATRIZ NO LEGOLOOP

> **Data:** 2026-11-06  
> **Autoridade:** Ecossistema LegoLoop & AIDD  
> **Status:** Canônico  
> **Escopo:** Gerar, consultar e auditar a matriz de peças (qualidade, aptidão, vocabulário e defasagem)  
> **Fontes de verdade:** `internal/matriz/` · `internal/cli/matriz.go` · `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md`

---

## 1. Visão Geral

O comando **`legoloop matriz`** gerencia o catálogo central de peças do ecossistema LegoLoop (`~/.aidd/canonical/registry_matrix.json` e `~/.aidd/canonical/matrix/*.json`). Ele quantifica reputação, audita critérios objetivos de qualidade, mapeia perfis de uso (etapas, naturezas e casos de uso) e rastreia defasagem incremental por hash SHA-256.

---

## 2. Comandos e Subcomandos

### 2.1 Visualização e Filtros (`matriz ver`)
Substitui o antigo comando `matrix` com suporte avançado a agrupamento e filtragem:
```bash
# Ver tabela completa consolidada
legoloop matriz ver

# Filtrar por tier (bronze, prata, ouro, platina)
legoloop matriz ver --tier ouro

# Filtrar por tipo de peça
legoloop matriz ver --tipo skill

# Filtrar por perfil de uso
legoloop matriz ver --natureza software --etapa construir --caso tdd

# Exibir todas as cópias físicas instaladas no SO
legoloop matriz ver --copias

# Saída estruturada JSON
legoloop matriz ver --json
```

### 2.2 Status e Defasagem (`matriz status`)
Inspeciona alterações nas peças sem reprocessar conteúdos:
```bash
legoloop matriz status
```
Relata quantas peças estão sincronizadas, novas, alteradas ou removidas.

### 2.3 Geração da Matriz (`matriz gerar`)
```bash
# Geração rápida sem gastar modelo de IA (peças sem perfil aguardam cérebro)
legoloop matriz gerar --sem-cerebro

# Geração forçada reprocessando tudo
legoloop matriz gerar --forcar
```

### 2.4 Vocabulário Fechado (`matriz vocab`)
```bash
# Listar vocabulário canônico ativo
legoloop matriz vocab listar

# Listar termos propostos pendentes de revisão
legoloop matriz vocab pendentes

# Aprovar um termo proposto
legoloop matriz vocab aprovar <termo>
```

### 2.5 Definição e Confirmação de Perfil (`matriz perfil`)
```bash
# Atribuir perfil manual a uma peça
legoloop matriz perfil definir <peça> --etapa construir --natureza software --caso tdd

# Confirmar perfil existente (declarado/confirmado tem precedência na Fusão)
legoloop matriz perfil confirmar <peça>
```

---

## 3. Invariantes
- A matriz opera em modo incremental com gravação atômica via arquivo temporário e rename.
- Perfis declarados/confirmados pelo usuário sobrevivem a regerações (`matriz gerar --forcar`).
- Falhas em gravações acionam reversão automática mantendo os arquivos intactos.
