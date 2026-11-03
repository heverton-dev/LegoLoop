// Package canon implementa a Canonização Transparente (Symlink Federation): varre ativos
// .md/.json, agrupa duplicatas exatas por SHA-256, elege a versão canônica, move-a para o
// cofre (~/.aidd/canonical) e substitui cada cópia por link, com backup e rollback.
//
// Fronteira de nomes: este pacote CANONIZA arquivos que já existem (comando `legoloop canonize`
// e rotas /api/canonize/*). Ele não cria processos novos: isso é do pacote internal/fabricar
// (comando `legoloop fabricar`).
package canon
