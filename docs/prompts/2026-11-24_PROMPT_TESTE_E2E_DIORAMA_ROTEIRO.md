# Prompt de Teste E2E Real — `legoloop-diorama-roteiro`

Cole o bloco abaixo em uma **nova sessão limpa** para testar o Processo `legoloop-diorama-roteiro` em tempo real (vida real), de ponta a ponta, em uma pasta temporária sem sujar o repositório:

```text
Quero fazer um teste end-to-end (E2E) real do Processo `legoloop-diorama-roteiro` usando o binário real `legoloop diorama-roteiro`.

Leia a skill `skills/legoloop-diorama-roteiro/SKILL.md` (via `call_skill` ou leitura direta) e execute o fluxo real de ponta a ponta em uma pasta temporária de teste (ex.: `%TEMP%\e2e_diorama_roteiro`):

1. **Etapa 0 (Configuração real):** Pergunte-me rapidamente o nicho/tema de teste (ou ofereça um padrão como "tecnologia e IA") e o modo de gates (`humano`, `hibrido` ou `mecanico`), crie `projeto.json` e `marca.json` reais na pasta temporária e rode:
   `legoloop diorama-roteiro validar-config --projeto <DIR> --json`
2. **Pré-voo real:** Rode `legoloop diorama-roteiro prevoo --projeto <DIR> --json` e me mostre os itens `OK`, `AVISO` e `FALTA`.
3. **Etapa 1 (Coleta real na internet):** Rode `legoloop diorama-roteiro coletar --projeto <DIR> --json` consultando fontes gratuitas reais para o meu nicho e me apresente as pautas coletadas para eu escolher uma.
4. **Etapa 2 (Roteiro + Gate Epistêmico real):** Com a pauta que eu escolher, gere `script.json` e `plano.md` reais e rode:
   `legoloop diorama-roteiro validar-roteiro --projeto <DIR> --json`
5. **Teste negativo (prova de bloqueio):** Faça um segundo teste injetando uma fonte domínio nível C (ex.: `blog-desconhecido.xyz`) em um `fato_consolidado` e confirme que `validar-roteiro` reprova (exit != 0) e NÃO grava em `acervo/temas.json`.
6. **Relatório final:** Mostre os resultados reais de cada comando (exit code, JSON resumido e eventuais bugs ou atritos encontrados no uso real).
```
