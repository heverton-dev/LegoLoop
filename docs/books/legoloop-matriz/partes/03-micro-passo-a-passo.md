# Parte III — Operação Passo a Passo e Invariantes Mecânicas

## Comandos operacionais da CLI

O subsistema Matriz é acessado por meio do comando `legoloop matriz`, oferecendo cinco fluxos principais de trabalho:

### 1. Inspecionar e filtrar peças catalogadas

Para visualizar as peças disponíveis e aplicar filtros por tier ou perfil de uso:
```bash
# Exibir todas as peças consolidadas
legoloop matriz ver

# Filtrar por nível de qualidade e reputação
legoloop matriz ver --tier ouro

# Filtrar por tipo e caso de uso específico
legoloop matriz ver --tipo skill --caso tdd

# Exportar visão estruturada em JSON para consumo de scripts
legoloop matriz ver --json
```

### 2. Avaliar status e defasagem incremental

Para auditar o cofre sem gastar cotas de modelo de inteligência artificial:
```bash
legoloop matriz status
```
O comando calcula a assinatura digital de todos os arquivos físicos e relata a quantidade exata de peças sincronizadas, novas, modificadas ou removidas.

### 3. Gerar e regenerar o inventário central

Para consolidar o registro canônico de peças em disco:
```bash
# Geração mecânica rápida sem consulta a LLM
legoloop matriz gerar --sem-cerebro

# Geração completa reprocessando todos os itens
legoloop matriz gerar --forcar
```

### 4. Gerenciar o vocabulário controlado

Para aprovar novos termos técnicos submetidos por agentes:
```bash
# Listar vocabulário em vigor
legoloop matriz vocab listar

# Inspecionar termos pendentes de aprovação
legoloop matriz vocab pendentes

# Homologar um termo proposto
legoloop matriz vocab aprovar <termo>
```

### 5. Atribuir e confirmar perfis de peça

Para fixar a classificação de uma peça no cofre:
```bash
# Atribuir perfil manual
legoloop matriz perfil definir <peça> --etapa construir --natureza software --caso tdd

# Confirmar perfil existente garantindo permanência
legoloop matriz perfil confirmar <peça>
```

## Garantias transacionais e invariantes

1. **Gravação Atômica com Arquivo Temporário:** Todas as operações de escrita no cofre utilizam escrita em arquivo efêmero seguida de renomeação atômica pelo sistema operacional. Em caso de interrupção abrupta, o manifesto anterior permanece intacto.
2. **Preservação de Confirmações Humanas:** Perfis e termos explicitamente confirmados pelo operador humano recebem prioridade absoluta na árvore de precedência e jamais são sobrescritos por regerações automáticas.
3. **Isolamento de Erro e Rollback:** Caso ocorra erro de validação durante o processamento de uma peça individual, o pipeline aborta com código de saída diferente de zero e reverte os arquivos para o estado imediatamente anterior.
