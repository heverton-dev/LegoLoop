# Arquitetura Meso

A arquitetura divide-se em seis módulos coesos dentro de `internal/video/roteiro`: `schemas.go` (validação e consistência cruzada), `prevoo.go` (5 bloqueantes `FALTA` e itens de P2/P3 como `AVISO`), `catalogo.go` e `catalogo_padrao.go` (classificação fail-closed por domínio), `coletor.go` (HTTP/RSS/JSON limitado a 10 MB e timeout <= 8s), `validador.go` (4 naturezas epistêmicas universais) e `processo.go` com `acervo.go` (limite de 2 correções automáticas e gravação atômica de `acervo/temas.json` somente após aprovação do Gate final).
