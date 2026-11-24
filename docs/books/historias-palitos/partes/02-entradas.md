# 2. Entradas

Entrada é tudo que o processo recebe antes de produzir qualquer coisa.
Há entradas de decisão e entradas de dados.

## 2.1 Entradas de decisão

Estas vêm da pessoa, na etapa 0.

| Entrada                | Exemplo                                | Onde fica           |
| :--------------------- | :------------------------------------- | :------------------ |
| Nicho                  | finanças, odontologia                  | `projeto.json`        |
| Estrutura narrativa    | explicador, investigação, lista        | `projeto.json`        |
| Formato                | 16:9, 9:16, 1:1                        | `projeto.json`        |
| Duração                | de 30 s a 20 min                       | `projeto.json`        |
| Voz                    | Thalita, Francisca, Antonio            | `projeto.json`        |
| Marca                  | fontes, cores, logo                    | `marca.json`          |
| Estilo visual          | dioramas, palitos, whiteboard          | `projeto.json`        |
| Animação               | motor em código, Flow ou híbrido       | `projeto.json`        |
| Distribuição           | canais e volume por dia                | `projeto.json`        |

Fonte: `processo/PROCESSO.md`, seção "Etapa 0".

## 2.2 Entradas de dados

Estas nascem durante o processo e alimentam a montagem.

| Entrada               | Conteúdo                               | Quem gera            |
| :-------------------- | :------------------------------------- | :------------------- |
| `roteiro.json`          | linhas narradas e fontes com link, data e nível | agente, etapa 2      |
| `plano.md`              | batidas, elenco, cenários, linha a linha | agente, etapas 2 a 4 |
| `prompts/*.txt`          | prompts prontos para o Flow            | agente               |
| `flow/`                 | imagens e vídeos baixados do Flow      | Google Flow          |
| `pronuncia.json`        | grafia fonética de nomes estrangeiros  | agente, etapa 6      |
| `acervo/temas.json`     | temas já feitos, para não repetir      | agente, etapa 1      |

## 2.3 Regra de entrada

Nada entra sem fonte.
Cada número do roteiro tem link, data e nível A ou B em `roteiro.json`.
Nível C serve só para achar pistas.
Fonte: `processo/PROCESSO.md`, etapa 2, e `processo/FONTES.md`.

## 2.4 Exemplo real

O projeto `curiosidades-malucas` tem estes arquivos de entrada na própria pasta:
`projeto.json`, `marca.json`, `roteiro.json`, `script.json`, `cenarios.json`, `elenco.json`, `storyboard.json` e `pronuncia.json`.
Verificado por listagem da pasta `curiosidades-malucas/`.
