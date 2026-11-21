---
title: "Processo Canônico de Gestão da Matriz de Peças"
subtitle: "Governança matemática, catalogação atômica, tiers de reputação e vocabulário controlado no LegoLoop"
author:
  - Heverton & Ecossistema LegoLoop
date: "20/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "CATÁLOGO E REPUTAÇÃO DE PEÇAS"
cover-title: "Matriz de Peças"
cover-accent: "Governança e Tiers de Reputação"
cover-caption: "ARQUITETURA CANÔNICA DO SUBSISTEMA MATRIZ"
cover-stages:
  - Ingestão Incremental
  - Cálculo de Hashes
  - Auditoria de Qualidade
  - Tiers de Reputação
  - Vocabulário ASD-STE100
  - Consulta Multidimensional
header-mark: "LEGOLOOP MATRIZ"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro documenta a arquitetura, os contratos e a operação do processo canônico `legoloop-matriz` no ecossistema LegoLoop. O subsistema governa o inventário central de peças (`registry_matrix.json`), assegurando que cada componente possua proveniência, integridade mecânica por SHA-256 e reputação auditável.

  A Parte I apresenta a fundação conceitual da Matriz, o modelo matemático de pontuação de peças e a hierarquia dos tiers de reputação. A Parte II disseca a arquitetura técnica, o vocabulário controlado ASD-STE100 e as dimensões de perfil (etapa, natureza e caso de uso). A Parte III detalha a operação da interface de linha de comando (`legoloop matriz`), as garantias transacionais de gravação atômica e os procedimentos de recuperação.

  Fonte primária: `docs/manuals/2026-11-06_MANUAL_FUNCAO_MATRIZ_LEGOLOOP.md`, `skills/legoloop-matriz/SKILL.md` e `internal/matriz/`.
---

# Como ler este livro

## A quem este livro se destina

1. **Operadores de Agentes e Desenvolvedores.** Profissionais que precisam consultar peças aptas, inspecionar tiers de reputação e garantir reutilização sem duplicidade. A Parte I e a Parte III são leituras prioritárias.
2. **Arquitetos de Governança.** Responsáveis pela integridade do cofre canônico, aprovação de vocabulário e certificação de novos componentes. A Parte II detalha os contratos invariantes.
3. **Engenheiros de Plataforma LegoLoop.** Mantenedores do núcleo em Go que estendem os algoritmos de score, filtros de busca e integração com os harnesses. Todas as partes constituem especificação formal.

## Estrutura da obra

| Parte          | Título                    | Conteúdo Principal                       |
| :------------- | :------------------------ | :--------------------------------------- |
| I              | Fundação e Modelo Mental  | Problema da defasagem, pontuação e tiers de reputação |
| II             | Arquitetura e Vocabulário | Subsistemas, metadados, ASD-STE100 e perfis de uso |
| III            | Operação e Invariantes    | Comandos da CLI, gravação atômica e recuperação |


# Parte I — Fundação e Modelo Mental da Matriz de Peças

## O desafio da governança em ecossistemas de agentes

Em ecossistemas autônomos orientados a agentes de inteligência artificial, o acúmulo desordenado de componentes produz degradação contínua da base de código. Agentes operando de forma isolada tendem a duplicar peças já existentes, reescrever ferramentas conhecidas com pequenas variações ou empregar componentes obsoletos que não atendem aos critérios de segurança vigentes.

O processo `legoloop-matriz` elimina essa entropia ao estabelecer uma fonte única de verdade matemática para todas as peças armazenadas no cofre canônico (`~/.aidd/canonical/`). Cada ferramenta, script, skill ou hook é submetido a uma inspeção estrita que afere sua proveniência, analisa seus contratos de entrada e saída e calcula sua integridade criptográfica por meio de hashes SHA-256 incrementais.

Com a Matriz, o ecossistema substitui palpites subjetivos sobre a qualidade de uma peça por uma métrica determinística reproduzível em qualquer ambiente de execução.

## Os quatro pilares do cálculo de reputação

A pontuação de cada peça na Matriz decorre da avaliação ponderada de quatro dimensões de engenharia:

1. **Rigor Estrutural:** Conformidade absoluta com os padrões de arquitetura, presença de documentação formal e ausência de stubs ou implementações simuladas.
2. **Contratos Mecânicos:** Definição binária dos artefatos consumidos e produzidos, com validação de schemas JSON e garantias de isolamento de caminho.
3. **Proveniência Rastreada:** Origem auditada da peça, histórico de autoria e verificação de compatibilidade com os protocolos industriais do LegoLoop.
4. **Histórico de Execução:** Comportamento observado em execuções anteriores, ausência de falhas não tratadas e estabilidade em ambientes estéreis.

## A hierarquia dos tiers de reputação

Com base na pontuação consolidada, cada peça é categorizada em um tier de reputação que define seu grau de maturidade e autorização de uso em produção:

| Tier           | Faixa de Pontuação | Critério de Qualidade                  | Autorização de Uso                     |
| :------------- | :----------------- | :------------------------------------- | :------------------------------------- |
| **Bronze**         | 0.00 a 0.49        | Peça inicial ou legada com contratos parciais | Ambientes experimentais isolados       |
| **Prata**          | 0.50 a 0.74        | Conformidade básica e testes mecânicos válidos | Uso interno e desenvolvimento          |
| **Ouro**           | 0.75 a 0.89        | Padrão industrial com documentação e TDD estrito | Produção padrão homologada             |
| **Platina**        | 0.90 a 1.00        | Peça mission-critical certificada e blindada | Sistemas críticos e pipelines centrais |

Peças abaixo do tier Ouro são impedidas de entrar em pipelines centrais até que sejam elevadas pelo processo de melhoria contínua de peças.


# Parte II — Arquitetura Sistêmica e Vocabulário Controlado

## Os cinco subsistemas da Matriz

A arquitetura do subsistema Matriz é composta por cinco módulos desacoplados implementados no pacote `internal/matriz/`:

1. **Detector de Defasagem (`matriz status`):** Examina o estado do disco frente ao manifesto central sem consumir tokens de LLM. Compara hashes SHA-256 de cada peça com a última foto registrada, identificando peças novas, alteradas, inalteradas ou removidas.
2. **Consolidador do Registro Canônico (`matriz gerar`):** Processa as peças elegíveis e sintetiza o catálogo unificado `registry_matrix.json` e os arquivos detalhados em `matrix/*.json`. Opera em dois modos: determinístico sem IA (`--sem-cerebro`) ou enriquecido com análise semântica.
3. **Guardião do Vocabulário ASD-STE100 (`matriz vocab`):** Mantém a tabela fechada de termos técnicos padronizados, garantindo que sinônimos proibidos não entrem nas descrições de peças e contratos de grafo.
4. **Gerenciador de Perfis de Uso (`matriz perfil`):** Vincula cada componente a uma tríade unívoca: Natureza (tipo do artefato), Etapa (fase do ciclo de engenharia) e Caso de Uso (contexto específico de aplicação).
5. **Motor de Filtros e Projeções (`matriz ver`):** Executa buscas parametrizadas em tempo real na memória com suporte a paginação, agrupamento e serialização estruturada em formato JSON ou tabelas de terminal.

## Vocabulário controlado e Lei Taxonômica ASD-STE100

A integridade do ecossistema depende da precisão semântica. A Matriz atua como autoridade central na fiscalização do dicionário canônico, rejeitando termos coloquiais ou ambíguos:

| Conceito Canônico         | Termo Homologado  | Sinônimos Proibidos                   |
| :------------------------ | :---------------- | :------------------------------------ |
| Fluxo de alto nível       | **Processo**          | esteira, fluxo, pipeline, loop        |
| Grafo de nós              | **DAG**               | árvore de execução, cadeia, grafo     |
| Componente atômico        | **Peça**              | bloco, módulo, brick, componente      |
| Instrução para agentes    | **Skill**             | habilidade, capacidade, prompt        |
| Ponto de execução         | **Nó**                | etapa de grafo, passo, task           |
| Barreira automática       | **Gate**              | trava, validação automática, check    |
| Ponto de decisão humana   | **Stop-Gate**         | portão humano, checkpoint manual      |
| Repositório central       | **Cofre**             | storage, base central, warehouse      |

O subsistema armazena propostas de termos em quarentena e exige aprovação explícita do operador antes de integrá-los ao catálogo ativo.

## Dimensões do Perfil de Peça

Cada peça homologada na Matriz recebe metadados taxonômicos que orientam a montagem dinâmica de grafos de execução:

- **Natureza:** Identifica a modalidade construtiva do artefato (`software`, `documento`, `video`, `infraestrutura`).
- **Etapa:** Indica a fase do ciclo de vida onde a peça atua (`concepcao`, `especificacao`, `construir`, `testar`, `fechar`).
- **Caso de Uso:** Descreve a especialidade técnica da peça (`tdd`, `qa-web`, `auditoria-ast`, `sintese-voz`).


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


