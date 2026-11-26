# Análise de impacto: inventário global, nome do macro e estilos de personagem

Data: 2026-11-26. Pedido do usuário: varrer o macro processo e dizer o que muda e o que não muda. Base: leitura de `docs/intents/2026-11-24_MAPA_MACRO_VIDEO_DIORAMA.md`, `docs/progress/2026-11-24_macro_video_diorama/estado.json`, `internal/diorama/**`, `internal/cli/cli.go`, e do projeto `proj--historias-palitos` (`processo/`, `ferramentas/engine/`).
Não verificado: nenhum teste foi rodado nesta análise; as contagens vêm de `grep` e `git ls-files`.

## 1. Onde estamos no macro processo

Macro `legoloop-video-diorama` (`estado.json`), cinco ciclos:

| Ciclo | Processo | Etapas | Estado |
|---|---|---|---|
| 1 | `legoloop-diorama-roteiro` | passo 0, 0, 1, 2 | DONE (canonizado) |
| 2 | `legoloop-diorama-ativos` | 3 (elenco), 4 (cenários e objetos), 5 (storyboard) | DONE (canonizado) |
| 3 | `legoloop-diorama-render` | 6 (voz), 7 (animação), 8 (montagem), 8b (cortes) | PENDING (15 tickets, nenhum aprovado) |
| 4 | `legoloop-diorama-distribuicao` | 9 a 13 | PENDING |
| 5 | `legoloop-video-diorama` | composição | PENDING |

Consequência: a Peça 3 ainda não foi construída, então este é o melhor momento para mudar o desenho. Mudar depois custa 3 a 5 vezes mais.

## 2. O que a análise encontrou (fatos)

### F1. A palavra "diorama" está em quatro lugares com pesos diferentes
- Nome dos processos, skills e pacotes: `internal/diorama/{roteiro,ativos,render}` (33 + 17 + 17 arquivos), `skills/legoloop-diorama-*`, `processos/legoloop-diorama-*`, comandos `diorama-roteiro` e `diorama-ativos` em `internal/cli/cli.go:109-125`.
- Registro: `espinha/registro.json`, `espinha/cobertura.json`, `docs/PROCESSOS_CANONIZADOS.md`, testes de MCP (`internal/mcp/*_test.go` citam `legoloop-diorama-render`).
- Documentos: cerca de 85 arquivos em `docs/` (intents 15, progress 12, books 11, plans 10, specs 8, prompts 6, routes 3, dags 3, reports 2, manuals 2).
- Conteúdo de estilo: o texto do roteiro e dos prompts fala em "diorama isométrico" (`plano.md` do demo, fichas de 4 faixas em `fichas.go`).

### F2. `estilo_visual` mistura duas coisas diferentes
`projeto.json` usa `estilo_visual: codigo|flow|hibrido`, e `internal/diorama/ativos/cenarios.go:39,41,128,135` trata esse campo como modo de produção (de onde vem a imagem). Estilo visual de verdade (palito, mangá, pixel art, massinha) não tem campo. O catálogo de estilos existe só em prosa (`proj--historias-palitos/processo/REPOS.md`).

### F3. O elenco da Peça 2 é de recorte de ficha, e o motor que dá resultado usa especificação
`internal/diorama/ativos/elenco.go`, `fichas.go` e `drivers.go` geram 3 fichas por personagem no Flow, recortam com rembg e fatiam em `sprites/*.json` (um PNG por pose). Foi o caminho que reprovamos em animação (A8, A9). O motor do First Brands (`puppet.js`) não usa imagem: usa um bloco `cast` em JSON (pele, terno, penteado, adornos) e desenha em SVG com 17 poses e 12 humores.

### F4. O contrato `cenas.json` aponta para sprite
Os elementos móveis de personagem em `cenas.json` trazem `sprite_id` e `acao` (por exemplo `narrador_pose_1`, `gesto_enfatiza`). O motor espera `puppet`, `pose`, `mood`, `face`, `talk`. Falta a ponte.

### F5. Não existe inventário
Só fragmentos: "elenco por partes trocáveis" (`processo/PROCESSO.md:33`) e personagem reutilizado no Flow (`CANAIS.md:25`). A SPEC e o PLAN do diorama não citam inventário. Hoje as partes ficam no código: 8 penteados num `switch` e chapéu, bigode, barba, vestido, gravata borboleta como opções (`puppet.js`, `ppSpec` e `ppHair`).

### F6. A Peça 3 foi desenhada como reescrita
O Mapa Macro (linha 95) manda reescrever o motor "de forma limpa e autoral" e o `render.go` virou ffmpeg + zoompan. O motor original já funciona nesta máquina (render de 4 cenas e da cena 1, 0 crédito) e passou no `check_timeline.py`.

## 3. O que NÃO muda

- Peça 1 inteira em lógica: descoberta, fontes A/B, roteiro em batidas, `script.json`, gates epistêmicos, `acervo/`.
- Voz (Etapa 6), visemas, WER, legendas `.srt`, mixagem -14 LUFS, cortes 9:16 (8b), revisor estéril, vitrine, agendador, Postiz.
- Contratos de pasta: `projeto.json`, `marca.json`, `script.json`, `audio/`, `data/timings.json`, `timeline.json`, `renders/`.
- Gates de ritmo (um corte a cada 3 s) e regra de que o render recusa timeline reprovada.
- Regras de Flow: aba do usuário, imagem custa 0, cotação antes de gasto.

## 4. O que muda (raio de impacto)

### M1. Inventário global do macro (decisão do usuário: global para este macro, não para o LegoLoop inteiro)
- Novo diretório versionado `inventario/video/` com `catalogo.json` (índice) e uma pasta por tipo.
- Tipos: pele, traje, penteado, adorno, adereço, pose, humor, personagem (bloco `cast` pronto), cenário (fundo do Flow com prompt e tags), objeto, SFX, trilha.
- Campos de cada item: `id`, `tipo`, `estilo`, `tags` (tema, época, canal), `parametros` (cores, forma, caminho do arquivo), `origem` (projeto e data), `aprovado` (sim/não), `hash`.
- Regras: um vídeo consulta por tema e estilo antes de criar; item novo entra só depois de aprovado; item de forma nova (colar de osso, toca, porrete) exige peça desenhada uma vez no motor e depois fica no inventário.
- Arquivos pesados (fundos de 600 KB cada) ficam em `projetos/inventario/` (fora do git) com hash no catálogo; só o índice JSON é versionado.
- Quem consulta e quem alimenta: Etapa 3 (elenco), Etapa 4 (cenários e objetos), Etapa 5 (storyboard), Etapa 8 (SFX e trilha), Etapa 13 (aprendizado grava o que funcionou).
- Pacote Go novo: `internal/video/inventario` (leitura, busca por tags, gravação atômica, validação de esquema).

### M2. Separar modo de produção de estilo visual
- `projeto.json`: renomear o uso atual para `modo_producao` (`codigo|flow|hibrido`) e criar `estilo` (id do catálogo de estilos).
- Migração: aceitar `estilo_visual` antigo por uma versão, avisando; `cenarios.go` passa a ler `modo_producao`.
- Catálogo de estilos: `inventario/video/estilos/<id>.json` com `renderizador` (qual desenhista de personagem), `prompt_mestre_cenarios`, `paleta`, `partes_permitidas`.

### M3. Estilo como "pele" plugável sobre um núcleo único
Núcleo comum (não muda por estilo): timeline, câmera, luz, efeitos, texto, voz, legenda, cortes. Muda só o desenhista de personagem:

| Estilo | Desenhista | Esforço |
|---|---|---|
| Molde de pílula (atual, `puppet.js`) | existe | zero |
| Palito | variante do molde (corpo de linha) | baixo |
| Quadrinho e mangá em vetor | novo renderizador SVG (cabeça, olhos e cabelo diferentes), mesmas poses e humores | médio |
| Pixel art | quadros de sprite por pose (sheet), troca de quadro no motor | médio |
| Recorte (papercut) | peças separadas com rig; reprovado em A9, só com peças de qualidade | alto |
| 3D | modelos reais no Three.js (glTF) | alto |

Resposta à dúvida do usuário: o macro processo não muda para trocar de estilo, mas o motor ganha um renderizador novo por estilo, e o inventário guarda partes por estilo. Mangá sem código novo não existe: o molde atual não vira mangá só com dados.

### M4. Etapa 3 (elenco) passa a produzir `cast`, e a ficha vira opcional
- Padrão para estilos vetoriais: o agente monta o bloco `cast` a partir do inventário (consulta por tags), grava `elenco.json` com o `cast` e valida contra o esquema do estilo. Sem Flow, 0 crédito.
- Ficha + rembg continua só para estilos de recorte e como referência visual.
- Código afetado: `internal/diorama/ativos/elenco.go`, `fichas.go`, `drivers.go`, `schemas.go` e seus testes; `SKILL.md` e manual do ciclo 2.

### M5. Contrato `cenas.json` com ponte para o motor
Acrescentar a `elementos_moveis` de personagem os campos `personagem_id`, `pose`, `humor`, `olhar` e `fala` (boolean), mantendo `sprite_id` só para estilos de sprite. Afeta `storyboard.go` e `schemas.go` de ativos, e a leitura em `render/timeline.go`.

### M6. Peça 3: portar o motor em vez de reescrevê-lo
- Proposta: `internal/diorama/render` (ou o novo nome) chama o motor Python/Chromium portado para `internal/video/motor/` (cópia versionada de `engine.html`, `puppet.js`, `layers.js`, `sets.js`, `formato.py`, `legendas.py`, `render.py`, `check_timeline.py`), com licença conferida (a RobbiSkill não declara licença; o motor do projeto foi estendido e `puppet.js` é igual ao dela).
- Go continua dono de voz, Whisper, mix, SFX, legendas e orquestração; o motor desenha quadros.
- Isso muda a SPEC da Peça 3 (hoje: reescrita autoral). Entra na Emenda de Concepção já pendente (SPEC v2.0, reabrir Peça 2, regenerar PLAN, ROUTE e DAG).
- Fundos do Flow entram como `bg`. Luz de contorno (`light.rim`) vira parte do estilo (a referência do lab usa `#ffd27a`).

### M7. Nome: "diorama" é estilo, não processo
Diorama isométrico é um estilo (e nem o que estamos fazendo agora). Proposta de nomes neutros:

| Atual | Proposto |
|---|---|
| `legoloop-video-diorama` (macro) | `legoloop-video-animado` |
| `legoloop-diorama-roteiro` | `legoloop-video-roteiro` |
| `legoloop-diorama-ativos` | `legoloop-video-ativos` |
| `legoloop-diorama-render` | `legoloop-video-render` |
| `legoloop-diorama-distribuicao` | `legoloop-video-distribuicao` |
| `internal/diorama/*` | `internal/video/*` |
| comandos `diorama-roteiro`, `diorama-ativos` | `video-roteiro`, `video-ativos` (com alias dos nomes antigos por uma versão) |

Atenção: `legoloop-video-rea` já existe (engenharia reversa de vídeo), então "video" sozinho precisa do segundo termo (`animado`) para não colidir.
O que tocar: 33 + 17 + 17 arquivos em `internal/diorama/`, 6 em `internal/cli`, 2 em `internal/mcp`, `skills/` e `processos/` (4 + 13 arquivos), `espinha/registro.json` e `cobertura.json`, `docs/PROCESSOS_CANONIZADOS.md`, `README.md`, `AGENTS.md`, `MEMORY.md`.
O que NÃO tocar: documentos históricos já arquivados (intents, specs, reports e livros de 2026-11-24): ganham uma nota de redirecionamento no topo, sem reescrever o conteúdo.
Risco: processos já canonizados estão registrados pelo nome antigo. Renomear exige migração do registro e do `espinha`, testes de CLI e MCP atualizados, e `legoloop fechar` outra vez.

### M8. Documentação e registro
- Regra permanente de transição projeto → processo: `TRANSICAO_PROJETO_PROCESSO.md` (já gravada).
- Novo achado A12 em `achados.md` (ver abaixo).
- Atualizar `estado.json` do macro (ciclos, nomes, contratos novos: `inventario/`, `elenco.json` com `cast`).

## 5. Ordem recomendada

1. Decisões do usuário (seção 6).
2. Emenda de Concepção (INTENT → SPEC v2.0 → PLAN/ROUTE/DAG) com M1 a M6; reabrir a Peça 2 para M4 e M5.
3. Ticket de renomeação (M7) separado e mecânico, depois da Emenda, com alias e migração do registro.
4. Peça 3 construída já sobre o motor portado, o inventário e o contrato novo.
5. Piloto: cena 1 e o First Brands (127 cenas) como prova de reprodução.

## 6. Decisões que dependem do usuário

1. Nome do macro: `legoloop-video-animado` (sugestão) ou outro.
2. Inventário: global do macro com tag de canal, índice no git e imagens fora do git (sugestão).
3. Campos novos em `projeto.json`: `modo_producao` e `estilo` (sugestão).
4. Motor: portar o original como núcleo (sugestão) ou manter a reescrita autoral da SPEC.
5. Primeiro estilo novo depois do molde de pílula: palito (baixo esforço, sugestão) ou mangá.
