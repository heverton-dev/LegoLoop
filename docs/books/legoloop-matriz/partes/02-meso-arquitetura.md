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
