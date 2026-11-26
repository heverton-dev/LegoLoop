# Transição projeto → processo (nota permanente)

Data: 2026-11-26. Origem: pedido do usuário. Vale para qualquer projeto futuro que vire processo do LegoLoop. O caminho inverso (processo → projeto) fica para outra conversa.

## Definições
- **Projeto:** pasta viva com ativos, roteiro e renders de UM vídeo ou canal (ex.: `proj--historias-palitos`).
- **Processo:** capacidade repetível do LegoLoop: skill em `skills/<nome>/SKILL.md`, código em `internal/`, tickets, verificadores e relatório. Registrado em `docs/PROCESSOS_CANONIZADOS.md`.

## O que aconteceu com Histórias Palitos (fatos lidos em 2026-11-26)
- O projeto original prova o resultado: First Brands, 127 cenas, 459 s, zero imagens do Flow (`proj--historias-palitos/processo/POSTMORTEM.md`).
- O que gera o resultado é a autoria de cenas (`processo/AUTORIA.md`, `build_timeline.py` com 282 linhas à mão) mais o motor (`ferramentas/engine/`: `puppet.js` 157 linhas, `layers.js` 261, `sets.js` 229).
- O LegoLoop migrou as etapas periféricas: skills `legoloop-diorama-roteiro` e `legoloop-diorama-ativos`, pacotes `internal/diorama/{roteiro,ativos,render}`.
- O render migrou como ffmpeg + zoompan (`ANALISE_ANIMACAO.md`, seção 1). O motor e a autoria ficaram fora do processo. Por isso o render saiu estático.

## Lição (vale como regra)
Migrar um projeto para processo exige portar a capacidade que PRODUZ a qualidade (aqui: autoria + motor), não só as etapas em volta. Etapa periférica aprovada não prova que o processo reproduz o projeto.

## Roteiro de transição (checklist reutilizável)
1. Inventário do projeto: etapas, ferramentas, motor, autoria, dados, licenças (anotar licença ausente, ex.: RobbiSkill sem licença; não publicar).
2. Marcar o núcleo de qualidade: o que, se removido, piora o vídeo (autoria, motor, verificador de ritmo, revisão por folha de contato).
3. Portar o núcleo primeiro, como código e skill versionados no LegoLoop. Dependência externa fica declarada, nunca implícita.
4. Critério de aceite: o processo reproduz o resultado de referência a partir do roteiro (First Brands: 127 cenas, CHECK OK, mediana ~3,3 s) e a revisão visual aprova.
5. Só então canonizar (`legoloop fechar`). Registrar o que ficou de fora e por quê.
6. Estilo = pele plugável sobre núcleo único (timeline, câmera, voz, legenda, cortes). Não criar um motor por estilo.
