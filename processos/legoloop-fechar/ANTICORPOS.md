---
schema: anticorpos.v1
processo: legoloop-fechar
total_anticorpos: 1
globais: 0
dominio: 1
regras_inviolaveis:
  - pacote: "fechar/prova-de-afirmacao"
    antidoto: "Todo fato que decide leva comando e resultado. Sem comando, escrever 'Não verificado' e rodar o teste antes de conclu..."
---

# 🛡️ Banco de Anticorpos — legoloop-fechar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (1)

### `fechar/prova-de-afirmacao`
- **Antídoto:** Todo fato que decide leva comando e resultado. Sem comando, escrever 'Não verificado' e rodar o teste antes de concluir. Proibido 'acho', 'parece', 'não há prova mas'.
- **Erro Evitado:** Afirmar causa, ausência de causa ou 'não há prova mas...' sem executar comando na resposta (ex.: agy off persistir no 9Router).
- **Causa Raiz:** Inferência por leitura de código e configuração substituiu execução real. Hipótese enunciada como possibilidade ('vestígios que podem confundir') sem teste.
- **Assinatura:** `sha256:b3459a432352fce2a9649e8ef3c926d64e059386f9a82b2bee00fe16c2079207`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._
