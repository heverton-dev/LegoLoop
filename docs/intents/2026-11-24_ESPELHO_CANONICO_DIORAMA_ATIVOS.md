# Espelho Canônico — `legoloop-diorama-ativos` (Ciclo 2)

> **Data:** 2026-11-24
> **Sessão da Entrevista:** `intv_1795554379812376300`
> **Sessão A2UI:** `sess_3`
> **Confiança:** 98% (meta: 95%)
> **Status:** APROVADO PELO USUÁRIO NA INTERFACE A2UI

## 1. Resultado Esperado
Processo canonizado `legoloop-diorama-ativos` (pacote `internal/diorama/ativos` + CLI `legoloop diorama-ativos` + DAG 5 de 5) que consome `projeto.json`, `marca.json`, `script.json` e `plano.md` do Ciclo 1 e executa, em ordem:
1. **Pré-voo do Ciclo 2:** valida os contratos do Ciclo 1, permissões de leitura e escrita em `flow/` e `sprites/`, e os pré-requisitos do modo visual escolhido (`codigo`, `flow` ou `hibrido`).
2. **Etapa 3 (Elenco e Sprites):** gera e valida 3 fichas por personagem (`turnaround`, `expressoes_3x2`, `poses_6x`) com `descricao_imutavel`, executa recorte transparente via `rembg` + despill verde e fatiamento automático em `sprites/*.json` (`id`, `personagem_id`, `tipo`, `arquivo_png`, `largura`, `altura`, `alfa_valido`, `bounding_box` sem sobreposição).
3. **Etapa 4 (Cenários e Objetos):** no modo `codigo` gera SVGs 100% autorais sem dependência externa; no modo `flow` opera fila resiliente com retry/backoff, suporte a `--ref`, trava de UUID `flow_projeto` e idempotência por hash (não regera ativo já válido); no modo `hibrido` combina ambos; exige mínimo de 1 cenário distinto a cada 2 linhas.
4. **Etapa 5 (Storyboard):** compila `cenas.json` cobrindo 100% das linhas de `script.json`, com metáfora visual e `>= 4` elementos móveis por linha.
5. **Governança de Gates:** respeita `politica_gates` (`humano`, `mecanico` ou `hibrido`) herdada de `projeto.json` nos Stop-Gates de elenco, quadros-teste e storyboard.

## 2. Quem se Beneficia
O operador do canal (recebe elenco consistente, sprites recortados e storyboard auditável sem intervenção manual repetitiva), o Ciclo 3 (`legoloop-diorama-render`, que consome `flow/`, `sprites/` e `cenas.json` validados como contrato imutável) e o Macroprocesso `legoloop-video-diorama`.

## 3. Por que Agora
O Ciclo 1 (`legoloop-diorama-roteiro`) já está entregue e canonizado. Sem os ativos visuais validados e o contrato `cenas.json` do Ciclo 2, o Ciclo 3 (renderização, animação e timeline) não pode consumir dados reais nem ser integrado ao Macroprocesso.

## 4. Critérios de Sucesso
1. Projeto válido vindo do Ciclo 1 percorre Etapas 3, 4 e 5 gerando `flow/`, `sprites/*.json` e `cenas.json` aprovados em todos os gates mecânicos.
2. Gate da Etapa 3 bloqueia ficha ausente, canal alfa inválido, dimensão abaixo do mínimo ou sobreposição de `bounding_box`.
3. Gate da Etapa 4 bloqueia densidade `< 1` cenário distinto a cada 2 linhas ou violação da trava `flow_projeto`; reexecução com mesmo hash é idempotente.
4. Gate da Etapa 5 bloqueia cobertura `< 100%` das linhas de `script.json`, ausência de metáfora visual ou `< 4` elementos móveis por linha.
5. Cobertura de testes `>= 95%` no pacote `internal/diorama/ativos`, zero stubs, gates da Espinha Dorsal verdes e fechamento completo via `legoloop fechar`.

## 5. Restrições Inegociáveis
1. Zero código sem licença de `RobbiSkill`: SVGs do modo `codigo` são 100% autorais.
2. Contrato de entrada imutável: consome `projeto.json`, `marca.json`, `script.json` e `plano.md` sem alterar o schema do Ciclo 1.
3. Trava `flow_projeto` obrigatória no modo `flow`/`hibrido` e idempotência por hash.
4. Política de gates herdada de `projeto.json` (`humano`/`mecanico`/`hibrido`).
5. Testes 100% herméticos com cobertura `>= 95%` (adaptadores externos como `rembg` e Flow injetáveis via interface).
6. Fechamento completo 5 de 5 (manual, livro, skill, DAG, build/install, governança viva, grafo e relatório).

## 6. Fora de Escopo
1. Etapas 6, 7, 8 e 8b (voz neural, Whisper, motor diorama HTML/SVG de animação, `timeline.json`, mixagem LUFS, render MP4/SRT e Shorts — pertencem ao Ciclo 3).
2. Etapas 9 a 13 (revisor estéril, vitrine, acervo/Postiz, comentários e aprendizado — pertencem ao Ciclo 4).
3. Composição ponta a ponta `P1->P4` (Ciclo 5).
4. Painel web paralelo fora do Looping Studio.
5. Alterações no schema do Ciclo 1 (`projeto.json`/`script.json`).
