# Manual Operacional — `legoloop-video-distribuicao`

- **Data:** 2026-11-27
- **Versao:** 1.0.0
- **Motor CLI:** `legoloop video-distribuicao`
- **Pacote Go:** `internal/video/distribuicao` + `internal/cli/video_distribuicao.go`
- **DAG Canonico:** `processos/legoloop-video-distribuicao/pipeline.dag.json`
- **Skill Universal:** `skills/legoloop-video-distribuicao/SKILL.md`

---

## 1. Proposito

O Processo `legoloop-video-distribuicao` e o P4 do processo composto `legoloop-video-animado`. Ele leva o video renderizado ate o canal com revisao, vitrine, agenda, publicacao pelo Postiz, comentarios e aprendizado. A pasta do projeto e `projetos/video-animado/<nome-do-projeto>/`.

## 2. Comandos e Subcomandos

| Subcomando | Entrada | Saida e Gate |
| :--- | :--- | :--- |
| `prevoo --projeto DIR` | ambiente | cria `acervo/`; AVISO para ffmpeg, Pillow e Postiz ausentes |
| `revisao folhas --video REL --versao N [--passo S]` | `renders/vN.mp4`, `script.json` | `review/sheet_NN.jpg` e `review/vN_pacote.json` |
| `revisao parecer --arquivo ARQ` | parecer JSON | `review/vN_parecer.json`; nunca sobrescreve |
| `revisao liberar --versao N` | parecer | `exit 0` so com ENVIAR |
| `vitrine pontuar --palavra P T1 T2...` | titulos | nota de 0 a 10 por titulo |
| `vitrine gravar --arquivo ARQ --titulo K` | `vitrine.json` | `publicacao.md`; valida 5 titulos, 3 capas por formato, fontes, tags e rotulo de IA |
| `capa --arquivo ARQ` | pedido de capa | PNG 1280x720; aviso acima de 2 MB |
| `acervo adicionar|listar|autorizar` | item JSON | `acervo/acervo.json`; regra de originalidade |
| `agenda [--config ARQ] [--agora RFC3339] [--aplicar]` | acervo | plano JSON; `--aplicar` grava e cria reels |
| `publicar canais` | POSTIZ_* | lista de canais |
| `publicar enviar --id ID --arquivo ARQ [--executar]` | pedido | dry-run por padrao; `--executar` envia |
| `comentarios listar|rascunhar|decidir` | `comentarios.json` | `acervo/respostas.json`; so aprovada sai |
| `aprendizado --metricas ARQ` | metricas | `acervo/aprendizado.md` |

## 3. Regras Inegociaveis

1. Visibilidade `private` por padrao; `public` exige `--confirmo-publico` e `autorizado: true`.
2. Rotulo de IA obrigatorio.
3. Segredos so por ambiente ou `.env*`; nunca impressos.
4. Resposta a comentario so apos aprovacao do usuario.
5. Nenhuma versao de render ou parecer e apagada.

## 4. Escopo declarado

Cortes 9:16 (BL-03), coleta automatica de comentarios e metricas, e instalacao do Postiz ficam fora. Limites de tamanho e caracteres por rede seguem nao verificados em pagina oficial.

## 5. Solucao de Problemas

| Sintoma | Causa | Acao |
| :--- | :--- | :--- |
| `segredo ausente: POSTIZ_API_KEY` | chave fora do ambiente | definir a variavel ou criar `.env` na raiz |
| `Postiz respondeu 401` | chave invalida | gerar outra em Configuracoes > Desenvolvedores |
| video fica privado apos publicar | falta auditoria da API da rede | pedir a auditoria (YouTube, TikTok) |
| `tema repetido` | tema quase igual no acervo | escolher outro tema |
| `NAO ENVIAR` no parecer | defeitos apontados | corrigir so os itens citados e gerar a proxima versao |
