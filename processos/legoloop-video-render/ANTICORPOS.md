---
schema: anticorpos.v1
processo: legoloop-video-render
total_anticorpos: 4
globais: 0
dominio: 4
regras_inviolaveis:
  - pacote: "internal/cli"
    antidoto: "Passar --formato 9:16 ao run_render.py; o comando video-render motor ja faz isso"
  - pacote: "internal/video/motor"
    antidoto: "Rodar na ordem prevoo, voz, timeline, mix, motor; o motor grava por ultimo"
  - pacote: "motor-local"
    antidoto: "motor-local/sets.js reduz a fonte para texto com mais de 2 caracteres"
  - pacote: "internal/video/ativos"
    antidoto: "Passar model Nano Banana e type image explicitamente em todo flow_generate de fundo"
---

# Banco de Anticorpos — legoloop-video-render

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Dominio (4)

### `internal/cli`

- **Erro:** Render do motor saiu em 16:9 e o gate 1080x1920 reprovou
- **Antidoto:** Passar --formato 9:16 ao run_render.py; o comando video-render motor ja faz isso
- **Causa raiz:** a ponte nao grava formato na timeline e o motor assume 16:9

### `internal/video/motor`

- **Erro:** timeline.json sobrescrita por dois produtores
- **Antidoto:** Rodar na ordem prevoo, voz, timeline, mix, motor; o motor grava por ultimo
- **Causa raiz:** video-render timeline e motor.Gerar gravam o mesmo arquivo com formatos diferentes

### `motor-local`

- **Erro:** Data do calendario do set office estourava o cartao
- **Antidoto:** motor-local/sets.js reduz a fonte para texto com mais de 2 caracteres
- **Causa raiz:** o valor padrao 1925 era maior que o cartao

### `internal/video/ativos`

- **Erro:** flow_generate saiu como video e gastou 12 creditos
- **Antidoto:** Passar model Nano Banana e type image explicitamente em todo flow_generate de fundo
- **Causa raiz:** o composer do projeto Flow estava em modo video
