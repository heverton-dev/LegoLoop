# 1. Visão geral

O projeto transforma um tema em um vídeo animado de canal.
O vídeo usa bonecos de palito, cenários, voz neural em português e trilha.
Todo o processo está descrito em `processo/PROCESSO.md` no repositório.

## 1.1 Princípio

O processo não depende de um modelo de IA nem de um ambiente.
Qualquer agente que leia o processo e os prompts consegue executá-lo.
Cada decisão fica em arquivo de texto (JSON ou Markdown).
Assim, outro agente retoma o trabalho só com a pasta do projeto.

Fonte: `processo/PROCESSO.md`, seção "Regra inegociável".

## 1.2 Quem faz o quê

Três atores executam o processo.

- **Pessoa (dona do canal).** Escolhe nicho, formato e duração. Aprova cada portão.
- **Agente de IA.** Pesquisa, escreve o roteiro, monta os arquivos e roda os scripts.
- **Serviços externos.** Google Flow gera imagens e vídeos. edge-tts gera a voz.

## 1.3 Dois modos de execução

No modo A, o agente tem terminal e navegador.
Ele roda os scripts e opera o Google Flow.

No modo B, o agente só conversa.
Ele escreve os arquivos e os prompts.
A pessoa roda os scripts e cola os prompts no Flow.

Os dois modos geram os mesmos arquivos e o mesmo vídeo.
Fonte: `processo/PROCESSO.md`, seção "Dois modos de execução".

## 1.4 Mapa do processo

O fluxo tem catorze etapas, de 0 a 13, com portões de aprovação.

```
Etapa 0  configuração  →  1 pauta  →  2 roteiro  →  3 elenco
   →  4 cenários e estilo  →  5 storyboard  →  6 voz
   →  7 animação  →  8 montagem  →  8b cortes curtos
   →  9 revisão  →  10 vitrine  →  11 publicação  →  12 comentários  →  13 aprendizado
```

O detalhamento de cada etapa está no capítulo 3.
