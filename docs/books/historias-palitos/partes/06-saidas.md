# 6. Saídas

Saída é o que o processo entrega ao final de cada etapa ou do vídeo inteiro.

## 6.1 Arquivos de vídeo

| Saída                     | Onde fica                              | Quando sai      |
| :------------------------ | :------------------------------------- | :-------------- |
| Vídeo longo (v1, v2, v3)  | `renders/`                               | etapa 8 e 9     |
| Corte vertical 9:16       | `renders/shorts/<id>.mp4`, 1080 x 1920   | etapa 8b        |
| Legenda `.srt` e `.vtt`       | pasta do projeto                       | etapa 8         |
| Trilha final              | `audio/`                                 | etapa 8         |

Regra: nenhuma versão é apagada.
Correções geram nova versão.
Fonte: `processo/PROCESSO.md`, etapa 9.

## 6.2 Arquivos de publicação

| Saída                   | Onde fica          | Conteúdo                               |
| :---------------------- | :----------------- | :------------------------------------- |
| `publicacao.md`           | pasta do projeto   | títulos, descrição, capítulos e miniaturas |
| Pacote da vitrine       | `publicacao.md`      | cinco títulos com nota e três capas    |
| `acervo/acervo.json`      | `acervo/`            | vídeos aprovados, validade e status    |

Fonte: `processo/PROCESSO.md`, etapas 10 e 11.

## 6.3 Saídas de controle

- `review/`: folhas de contato e relatórios de revisão.
- `acervo/aprendizado.md`: recomendações semanais.
- `cortes/<video>.json`: especificação de cada corte, com título, descrição e status de envio.

## 6.4 Destino

A publicação passa pelo Postiz, instalado numa VPS.
O Postiz envia para YouTube, Instagram, TikTok e LinkedIn.
Fonte: `processo/PROCESSO.md`, etapa 11.
