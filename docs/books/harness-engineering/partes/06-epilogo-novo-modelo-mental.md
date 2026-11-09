# Capítulo 6: Epílogo — O Novo Modelo Mental do Engenheiro com IA

A frustração que você sentiu antes de pedir este livro é o marco divisor na jornada de qualquer programador que decide trabalhar com inteligência artificial.

A maioria das pessoas desiste e volta para o desenvolvimento manual dizendo: *"IA é uma porcaria, só serve para fazer autocomplete"*. Outros continuam sofrendo em silêncio, passando 12 horas por dia corrigindo alucinações e refazendo arquivos.

Os poucos que constroem coisas monumentais e escaláveis fazem uma única mudança psicológica:

## A Mudança de Pergunta

Quando um modelo comete um erro ridículo, o desenvolvedor frustrado pergunta para o modelo:
> *"Por que você fez isso de novo? Eu já não te avisei cinco vezes no prompt?"*

Essa pergunta não leva a lugar nenhum. Ela pressupõe que o modelo tem memória, vergonha ou bom senso.

O engenheiro de Harness faz uma pergunta completamente diferente para si mesmo:
> **"Qual brecha mecânica no meu software permitiu que essa ação fosse executada sem quebrar o build?"**

* O modelo gerou `null`? $\rightarrow$ *Faltou um `ValidateStrict()` no struct.*
* O modelo gravou na pasta errada? $\rightarrow$ *Faltou um Path Jail no wrapper de I/O.*
* O modelo inventou um método que não existe? $\rightarrow$ *Faltou um teste de compilação antes de entregar a resposta.*
* O modelo esqueceu um detalhe do repositório? $\rightarrow$ *O contexto enviado para ele estava poluído com arquivos inúteis.*

Quando você para de tratar a IA como um estagiário preguiçoso e passa a tratá-la como um **motor a combustão**, tudo muda. Um motor a combustão sem câmbio e sem freio explode o carro contra o muro. Mas quando você coloca pistões de aço, freios ABS e um chassi rígido (o Harness), o mesmo motor leva o carro a 300 km/h com segurança absoluta.

## O Próximo Nível

O LegoLoop já tem a base arquitetural mais avançada que você poderia desejar:
* Você tem o **LegoPress** para gerar livros e manuais sem alucinação visual;
* Você tem o **DAG acíclico** para garantir que cada nó só rode depois do anterior;
* Você tem o **Quinteto de Governança ISO** para rastrear cada decisão;
* Você tem o **Video-REA** para dissecar ferramentas e adotar código aberto.

O que faltava não era mais uma sigla ou mais um framework. Faltava apenas endurecer o aço do arreio: fechar as portas de saída falsa, travar os caminhos no disco e exigir que o código comprove o que diz antes de dar o ticket como pronto.

Respire fundo. A confusão passou. Agora você tem o mapa, tem a chave e sabe exatamente onde cada parafuso precisa ser apertado.
