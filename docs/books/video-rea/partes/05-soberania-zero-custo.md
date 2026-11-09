# Parte V — Nomenclatura, LegoPress e Skills

## A Regra de Ouro da Nomenclatura: `legoloop-*`

Todo pipeline concebido a partir de engenharia reversa herda obrigatoriamente a assinatura institucional do ecossistema:
- O novo diretório de loops adota `loops/<slug>/pipeline.dag.json` com campo `"name": "legoloop-<slug>"`.
- A nova CLI é registrada como subcomando de primeiro nível: `legoloop <slug>`.
- A nova skill do agente é distribuída no cofre: `skills/legoloop-<slug>/SKILL.md`.

## Publicação Automática do Livro Didático

Nenhum pipeline é entregue como "caixa preta". A fábrica convoca imediatamente o motor **LegoPress**:
1. Cria o manifesto canônico `livro.json` em `docs/books/<slug>/`.
2. Estrutura as partes didáticas (Macro, Meso, Micro, Prática e Governança).
3. Compila o manuscrito consolidado e o documento tipográfico Typst.
4. Produz o arquivo PDF em padrão editorial de 170x240mm e suas imagens de visualização em `preview/`.
5. Roda a auditoria geométrica (`legoloop legopress auditar`) para garantir que nenhum cabeçalho ou página órfã seja entregue.

Dessa forma, o conhecimento extraído de qualquer vídeo público é perenizado, auditável, de custo zero e pronto para ser executado por agentes e humanos.
