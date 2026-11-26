# Espelho Canônico — `legoloop-diorama-roteiro` (Ciclo 1)

> **Nota de redirecionamento (2026-11-26, AC-E11.4):** o processo `legoloop-diorama-roteiro` foi renomeado para `legoloop-video-roteiro`. Este documento e historico e mantido sem alteracao de conteudo.


Sessão da entrevista: `intv_1795520712665396700` (9 turnos; cérebro final: `claude/sonnet`; confiança: 96% / meta 95%; status: `COMPLETED`).
Data: 2026-11-24.

## resultado_esperado

Processo canonizado `legoloop-diorama-roteiro` que executa, em ordem:
1. Etapa 0, configuração guiada: fabrica o `projeto.json` por vídeo, com validação de schema. Campos: `nicho` (texto livre), `perfis_fonte` (`geral` como fallback), `fontes_extras` e `fontes_excluidas` opcionais, `politica_gates` (`humano`, `mecanico` ou `hibrido`, com `gates_humanos`) e `video_referencia` opcional. O `marca.json` é herdado do canal ou criado de template canônico (`legoloop-canal-redc`).
2. Pré-voo determinístico modular, com bloqueio nestes casos: `legoloop version` não retorna exit 0; catálogo declarativo de fontes por nicho ausente ou JSON inválido; pasta do projeto ou `acervo/` sem leitura e escrita reais; nenhum conector compatível com `perfis_fonte` responde na sonda HTTP (timeout <= 8s, status 200); `video_referencia` preenchido e `legoloop-video-rea` ausente do Cofre.
3. Validação estrita: schema inválido, campo obrigatório ausente ou inconsistência entre `projeto.json` e `marca.json` bloqueiam.
4. Descoberta de pautas: catálogo declarativo e extensível de conectores gratuitos, filtrado por `perfis_fonte` e ajustado por `fontes_extras` e `fontes_excluidas`. Temas já presentes em `acervo/temas.json` são bloqueados.
5. Escolha da pauta, conforme `politica_gates`.
6. Geração de `script.json` e `plano.md`: cada linha tem natureza epistêmica (`fato_consolidado`, `evidencia_preliminar`, `alegacao_terceiro` ou `narrativa_transicao`).
7. Gate final da Etapa 2: se reprovar, o Processo tenta corrigir até 2 vezes usando só o relatório estruturado do validador (linha, natureza epistêmica e regra violada). Se a falha persistir, ele para, mesmo no modo `mecanico`, e entrega o relatório completo para decisão humana.

## quem_se_beneficia

O operador do canal, que recebe roteiros auditáveis sem checar fontes manualmente. Os Processos 2 e 3 do Macroprocesso consomem `script.json` e `plano.md` já validados. O próprio Macroprocesso ganha um acervo de temas sem repetição.

## por_que_agora

O Ciclo 1 destrava os Ciclos seguintes, porque sem roteiro validado não há vídeo. Hoje não há bloqueio mecânico contra temas repetidos, nem separação rastreável entre fato, evidência preliminar e alegação de terceiro.

## criterio_de_sucesso

- Um vídeo com `projeto.json` e `marca.json` válidos percorre o Processo e gera `script.json` e `plano.md` aprovados no Gate final.
- Cada linha de `script.json` cumpre a regra da sua natureza:
  - `fato_consolidado`: 1 fonte A ou 2 fontes B.
  - `evidencia_preliminar`: escopo, preprint ou amostra declarados.
  - `alegacao_terceiro`: atribuição a quem afirmou.
  - `narrativa_transicao`: gancho ou pergunta.
- O tema só é gravado em `acervo/temas.json` após a aprovação do Gate final.
- Os três erros bloqueantes e os cinco itens do pré-voo interrompem o Processo antes de qualquer consulta externa. Os itens só de AVISO não interrompem.

## restricoes_inegociaveis

- Nível A: `.gov.br`, `.gov`, `bcb.gov.br`, `cvm.gov.br`, `sec.gov`, `pubmed`, `doi.org`, `who.int`, `ibge.gov.br`, e releases oficiais de `github.com` no nicho software/open source.
- Nível B: imprensa de referência declarada (Reuters, AP, BBC, Bloomberg, CNBC, Folha, Estadão, G1, Valor).
- Preprints: `arxiv`, `medrxiv` e `biorxiv` exigem a natureza `evidencia_preliminar`.
- Nível C: regra fail-closed. Qualquer domínio fora do catálogo e fora de `fontes_extras` é nível C e fica bloqueado como prova. Blogs, Wikipedia e agregadores servem só para descoberta.
- Classificação: declarada no catálogo, por domínio e sufixo institucional, além do nível padrão em conectores primários diretos.
- Fontes: apenas conectores gratuitos.
- Validação: universal para qualquer nicho, sem viés jurídico.
- Gate final: o Processo faz no máximo 2 correções automáticas e para mesmo no modo `mecanico`.
- Tema em `acervo/temas.json`: nunca é gravado antes da aprovação do Gate final, em nenhum modo.

## fora_de_escopo

- Maquete física de diorama.
- Tudo que pertence aos Processos 2 e 3: Chrome do Flow na porta 9333, FFmpeg, Chromium, edge-tts e gitleaks. Esses itens geram só AVISO no pré-voo.
- Uso do `legoloop-concepcao` para preencher cada vídeo. Ele só concebe o Processo.
- Fontes pagas.
- Uso de fonte nível C como prova.
