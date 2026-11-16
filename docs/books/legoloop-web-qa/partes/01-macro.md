# Parte I — Visão Macro: Por que Auditar antes de Corrigir

## O problema das correções reativas em interfaces web

Aplicações web modernas falham silenciosamente em quatro dimensões que testes unitários de back-end não enxergam: controles visuais que não produzem efeito ao clique, elementos sobrepostos em telas menores, violações de contraste ou navegação por teclado, e links ou assets que retornam erro HTTP após um build. Quando agentes de código tentam corrigir um defeito visual diretamente a partir de uma queixa genérica, o resultado costuma ser regressão em outra rota ou viewport.

O processo `legoloop-web-qa` resolve esse ciclo separando rigorosamente **diagnóstico determinístico** de **mutação de código**. Primeiro, o sistema prova matematicamente o estado da aplicação através de oito ferramentas reais. Somente depois que o relatório consolidado existe em disco, o operador humano e o agente iniciam a entrevista socrática de afunilamento que origina o pacote de correção governado até `BUILD`.

## Três princípios arquiteturais do processo

1. **Zero Stubs no Preflight:** Se qualquer uma das oito ferramentas obrigatórias não estiver instalada e acessível, o processo aborta no nó `n01_preflight` com `exit_code == 1`. Nenhuma etapa simula resultado positivo ou substitui uma peça ausente por resposta vazia.
2. **Dupla Cobertura de Ambiente:** O runner `tools/webqa/webqa.mjs` aceita tanto uma URL já publicada (`--url`) quanto um comando e porta locais (`--cmd` + `--port`), subindo a aplicação, aguardando a porta responder e encerrando o processo filho ao final da execução.
3. **Relatório como Contrato de Entrada:** A saída em `docs/qa-web/YYYY-MM-DD_<alvo>/` não é um log descartável. Os arquivos `RELATORIO_QA_WEB.json` e `RELATORIO_QA_WEB.md` são insumos obrigatórios para o nó `n11_entrevista_correcao`.
