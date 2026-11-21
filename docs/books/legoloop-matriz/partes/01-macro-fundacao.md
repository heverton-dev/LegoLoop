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
