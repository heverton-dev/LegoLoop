// Package fabricar é a Fábrica Universal de Processos Canônicos: recebe um pedido (CLI, Chat ou
// Studio) e o conduz pelas 5 camadas (Clarificação, Curadoria multi-drive, Fusão, Tríade e
// Template), materializando um processo `legoloop-<slug>` em três meios: SKILL.md, .loop.json e
// .dag.json.
//
// Fronteira de nomes: este pacote FABRICA processos novos (comando `legoloop fabricar`). Ele não
// tem relação com o pacote internal/canon, que CANONIZA arquivos existentes (desduplica por
// SHA-256 e troca as cópias por links para o cofre; comando `legoloop canonize`).
//
// É um orquestrador fino: reaproveita interview, specgen, plangen, assets, canon, loop, dag, gate
// e reputation, e não reimplementa nada deles. Tudo que depende de LLM, relógio, disco ou unidades
// entra por Deps, o que torna o motor determinístico e testável sem tocar o ~/.aidd real.
package fabricar
