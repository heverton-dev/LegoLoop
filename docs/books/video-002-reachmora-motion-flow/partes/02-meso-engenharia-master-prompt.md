# Parte II — Meso: Engenharia do Master Prompt

## 2.1 Anatomia de um Master Prompt Resiliente

O Master Prompt é a peça central que governa a coerência de todas as tomadas geradas:

```
[ESTILO VISUAL BASE]
Minimalist 3D isometric flat design, clean corporate palette (#0f172a, #8b5cf6, #ffffff).

[REGRAS DE CÂMERA]
Slow dynamic push-in, 45 degree top-down angle, isometric projection, zero chromatic aberration.

[PARÂMETROS DE MOVIMENTO]
Smooth spring animation, ease-out-cubic velocity curve, floating geometric cards.

[SLOT DE CENA]
Beat narrativo: {cena_descricao}
Palavra de ênfase: {palavra_chave}
```

#painel-sucesso("ZERO ALUCINAÇÃO DE ESTILO", [
  Ao fixar a paleta e os ângulos de câmera nas regras negativas do prompt, o modelo cognitivo é proibido de inventar novos estilos visuais no meio da narrativa.
])
