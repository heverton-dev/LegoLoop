# Exemplo real: esteira "clínica veterinária"

Esta pasta guarda uma execução **real** do LegoLoop (31/10/2026), via 9Router da VPS.

## O que aconteceu, em linguagem simples

A "clínica veterinária" **não é um personagem nem um projeto existente**. É só a **ideia de entrada**, uma frase passada na linha de comando:

> App web para uma clínica veterinária agendar consultas, enviar lembretes por WhatsApp e manter o prontuário dos pets

A partir dela, a esteira montou quatro peças Lego em sequência. Cada peça é **uma única chamada a um modelo de IA, isolada das outras**: ela não "lembra" nada da peça anterior e só lê o arquivo que a anterior gravou no disco.

```
ideia (frase)
   │
   ▼
[n1 interview] ── combo code-fast ──► saida/BRIEFING.md   (problema, público, funcionalidades)
   │
   ├──────────────────────────────┐   (estas duas rodam AO MESMO TEMPO)
   ▼                              ▼
[n2a aidd-spec] ── code-pro ──► saida/SPEC.md        [n2b riscos] ── code-free ──► saida/RISCOS.md
   │                (critérios Pass/Fail)                            (tabela de riscos)
   ▼
[n3 check] ── script local (findstr) ──► saida/CHECKLIST.txt  (só as linhas "- [ ]" da SPEC)
```

O que diz a cada peça **o que fazer** é o arquivo de skill dela, em `skills/<nome>/SKILL.md` (um texto curto de instruções).
**Quando a peça está pronta** é decidido pelo "gate" do `pipeline.dag.json`. Exemplo: a SPEC só passa se tiver a seção "Critérios de Aceitação" com itens `- [ ]`. Se um gate falha, as peças seguintes não rodam.

## Arquivos

| Arquivo | O que é |
|---|---|
| `pipeline.dag.json` | o contrato da esteira: peças, ordem, modelo de cada uma e gates |
| `skills/*/SKILL.md` | as instruções de cada peça |
| `saida/BRIEFING.md` | gerado pelo n1 a partir da ideia |
| `saida/SPEC.md` | gerado pelo n2a a partir do BRIEFING |
| `saida/RISCOS.md` | gerado pelo n2b a partir do BRIEFING |
| `saida/CHECKLIST.txt` | extraído da SPEC pelo script n3 |

## Rodar de novo com outra ideia

```bash
bin\legoloop.exe run examples\clinica-veterinaria\pipeline.dag.json --workdir <pasta-de-saida> --skills examples\clinica-veterinaria --allow-scripts --var "idea=<sua ideia>"
```

Para isso, o 9Router precisa estar configurado: `NINEROUTER_URL` e `NINEROUTER_KEY` definidas no ambiente.
