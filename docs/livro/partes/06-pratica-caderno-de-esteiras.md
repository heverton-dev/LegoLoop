# PARTE VI — Prática: O Caderno de 10 Esteiras

## 26. Do Conceito à Prática: Como Montar Sistemas de Alta Fidelidade

Um dos maiores problemas da literatura sobre inteligência artificial é o excesso de teoria abstrata desacompanhada de exemplos do mundo real. Leitores terminam livros sabendo a definição de tokens e redes neurais, mas continuam sem saber como aplicar a tecnologia para resolver problemas concretos na manhã de segunda-feira em suas empresas.

Esta parte do livro foi concebida como um **Caderno de Receitas de Engenharia**. Cada capítulo apresenta uma esteira completa do início ao fim, detalhando o problema de negócios, a composição de nós modulares, a localização estratégica dos stop-gates humanos e a partitura executável `ROUTE.json`.

---

## 27. Receita 1: Esteira de Especificação de Software em BDD

### O Problema de Negócios
Equipes de desenvolvimento perdem semanas implementando sistemas com base em requisitos vagos recebidos por mensagens de chat ou conversas informais. O desenvolvedor acha que entendeu uma coisa, mas o cliente desejava outra completamente diferente.

### A Composição Modular
* **Nó 1 (Entrevista Socrática):** Formula perguntas de fechamento sobre cenários de exceção.
* **Nó 2 (Extrator de Regras):** Identifica entidades, relacionamentos e restrições.
* **Nó 3 (Gerador de Cenários BDD):** Converte as regras na sintaxe clássica *Dado que / Quando / Então*.
* **Stop-Gate Humano:** O Product Owner (PO) lê os cenários e aprova com um clique.

```json
{
  "nome": "especificacao-bdd",
  "versao": "2.0",
  "nos": [
    { "id": "entrevista", "tarefa": "resolver-premissas", "dependencias": [] },
    { "id": "extracao", "tarefa": "isolar-entidades", "dependencias": ["entrevista"] },
    { "id": "gerar-bdd", "tarefa": "escrever-gherkin", "dependencias": ["extracao"] },
    { "id": "gate-po", "tarefa": "aprovacao-humana", "dependencias": ["gerar-bdd"] }
  ]
}
```

---

## 28. Receita 2: Esteira de Desenvolvimento Blindado em TDD

### O Problema de Negócios
Assistentes de IA tradicionais geram códigos aparentemente funcionais, mas que falham miseravelmente em casos de borda ou quebram funcionalidades existentes do sistema.

### A Composição Modular
* **Nó 1 (Criador de Testes):** Lê o arquivo `SPEC.md` e escreve testes unitários rigorosos. **A esteira exige que o teste falhe (Fase Vermelha do TDD).**
* **Nó 2 (Implementador Estéril):** Recebe exclusivamente o teste falho e a assinatura da função. Implementa o código mínimo necessário até que o teste passe com sucesso (Fase Verde).
* **Nó 3 (Auditor de Limites):** Verifica se alguma função possui mais de 35 linhas ou se houve criação de dependências não autorizadas.

```json
{
  "nome": "desenvolvimento-tdd",
  "versao": "2.0",
  "nos": [
    { "id": "escrever-teste", "tarefa": "fase-vermelha", "dependencias": [] },
    { "id": "implementar", "tarefa": "fase-verde", "dependencias": ["escrever-teste"] },
    { "id": "auditoria-limites", "tarefa": "validar-regras", "dependencias": ["implementar"] }
  ]
}
```

---

## 29. Receita 3: Esteira de Auditoria e Limpeza de Código Legado

### O Problema de Negócios
Empresas mantêm sistemas legados com milhares de linhas de código antigo contendo funções duplicadas, comentários obsoletos e potenciais falhas de segurança. Quando um desenvolvedor tenta limpar o código manualmente, frequentemente introduz regressões que quebram sistemas em produção.

### A Composição Modular
* **Nó 1 (Mapeador de Árvore Sintática):** Analisa o código através de análise estática (AST) sem executar o binário, localizando variáveis não utilizadas e código morto (*dead code*).
* **Nó 2 (Auditor de Vulnerabilidades):** Verifica se chaves de API, senhas ou chamadas de sistema inseguras estão presentes no repositório.
* **Nó 3 (Formatador Determinístico):** Aplica as regras canônicas de formatação da linguagem (`gofmt` ou similar) com tolerância zero para desvios de estilo.
* **Stop-Gate Humano:** Um arquiteto de software revisa o relatório detalhado de alterações e autoriza o commit de modernização.

```json
{
  "nome": "limpeza-legado",
  "versao": "2.0",
  "nos": [
    { "id": "mapear-ast", "tarefa": "analisar-funcoes-mortas", "dependencias": [] },
    { "id": "auditoria-sec", "tarefa": "procurar-vulnerabilidades", "dependencias": ["mapear-ast"] },
    { "id": "formatacao", "tarefa": "aplicar-gofmt", "dependencias": ["auditoria-sec"] },
    { "id": "gate-arquiteto", "tarefa": "aprovacao-humana", "dependencias": ["formatacao"] }
  ]
}
```

---

## 30. Receita 4: Esteira de Engenharia Reversa de APIs Web

### O Problema de Negócios
Uma equipe precisa integrar seu sistema a um serviço externo antigo cuja documentação oficial está desatualizada ou simplesmente desapareceu dos servidores da empresa parceira.

### A Composição Modular
* **Nó 1 (Captura de Tráfego de Rede):** Registra requisições e respostas HTTP legítimas em ambiente de homologação.
* **Nó 2 (Extrator de Esquema JSON):** Identifica tipos de dados, parâmetros obrigatórios e cabeçalhos de autenticação de cada endpoint.
* **Nó 3 (Gerador de Contrato OpenAPI):** Cria automaticamente a especificação Swagger/OpenAPI 3.0 documentando toda a interface.
* **Nó 4 (Gerador de SDK Cliente):** Produz uma biblioteca cliente em Go ou TypeScript com tipagem estrita para consumo imediato.

```json
{
  "nome": "engenharia-reversa-api",
  "versao": "2.0",
  "nos": [
    { "id": "captura-rede", "tarefa": "interceptar-trafego", "dependencias": [] },
    { "id": "extrair-json", "tarefa": "inferir-esquema", "dependencias": ["captura-rede"] },
    { "id": "gerar-openapi", "tarefa": "produzir-swagger", "dependencias": ["extrair-json"] },
    { "id": "gerar-sdk", "tarefa": "compilar-cliente", "dependencias": ["gerar-openapi"] }
  ]
}
```

---

## 31. Receita 5: Esteira de Processamento de Notas Fiscais

### O Problema de Negócios
O setor de contas a pagar de uma empresa gasta dezenas de horas por semana digitando manualmente dados de faturas em PDF recebidas por e-mail, gerando frequentes erros de cálculo e atrasos de pagamento.

### A Composição Modular
* **Nó 1 (Extrator OCR Estruturado):** Converte o PDF da fatura em dados tabulares com reconhecimento de campos chave.
* **Nó 2 (Validador Matemático):** Confere se a soma de produtos, frete e tributos é matematicamente idêntica ao valor total faturado.
* **Nó 3 (Auditor de CNPJ):** Checa a situação cadastral do fornecedor na base pública da Receita Federal.
* **Stop-Gate Humano:** Faturas com valores superiores a R$ 20.000,00 ou fornecedores recém-cadastrados exigem confirmação explícita de um gestor financeiro.
* **Nó 4 (Emissor de Lançamento):** Dispara a ordem de pagamento no sistema ERP corporativo via webhook seguro.

```json
{
  "nome": "processamento-faturas",
  "versao": "2.0",
  "nos": [
    { "id": "extracao-ocr", "tarefa": "ler-pdf", "dependencias": [] },
    { "id": "validacao-matematica", "tarefa": "conferir-somas", "dependencias": ["extracao-ocr"] },
    { "id": "auditoria-cnpj", "tarefa": "consultar-receita", "dependencias": ["validacao-matematica"] },
    { "id": "gate-financeiro", "tarefa": "aprovacao-gestor", "dependencias": ["auditoria-cnpj"] },
    { "id": "emissao-erp", "tarefa": "lancar-pagamento", "dependencias": ["gate-financeiro"] }
  ]
}
```

---

## 32. Receita 6: Esteira de Livros e Manuais Corporativos (Typst)

### O Problema de Negócios
Produzir manuais técnicos ou livros corporativos no Microsoft Word gera problemas crônicos de diagramação: imagens que saltam de página, tabelas cortadas e fontes desconfiguradas a cada nova edição.

### A Composição Modular
* **Nó 1 (Compilador de Capítulos):** Consolida os arquivos Markdown das partes no documento mestre.
* **Nó 2 (Fábrica de Diagramas SVG):** Valida dimensões, acessibilidade visual e paleta de cores dos diagramas visuais.
* **Nó 3 (Linter Editorial):** Verifica largura de tabelas, pontuação e conformidade tipográfica contra o manual de estilo.
* **Nó 4 (Motor Typst 0.15):** Compila o PDF final de altíssima resolução gráfica com paginação perfeita em menos de 1 segundo.

```json
{
  "nome": "publicacao-editorial",
  "versao": "2.0",
  "nos": [
    { "id": "consolidar-md", "tarefa": "unir-capitulos", "dependencias": [] },
    { "id": "validar-svg", "tarefa": "checar-diagramas", "dependencias": ["consolidar-md"] },
    { "id": "linter-editorial", "tarefa": "auditar-tabelas", "dependencias": ["validar-svg"] },
    { "id": "compilar-typst", "tarefa": "gerar-pdf-alta-resolucao", "dependencias": ["linter-editorial"] }
  ]
}
```

---

## 33. Receita 7: Esteira de Triagem de Currículos e Recrutamento Ético

### O Problema de Negócios
Departamentos de Recursos Humanos recebem centenas de currículos para cada vaga aberta e precisam realizar a triagem inicial sem incorrer em vieses inconscientes ou discriminações indevidas.

### A Composição Modular
* **Nó 1 (Anonimizador de Dados):** Oculta temporariamente nome, gênero, idade, foto e endereço do candidato.
* **Nó 2 (Avaliador de Competências Técnicas):** Compara a experiência comprovada contra os requisitos exigidos pela vaga.
* **Nó 3 (Gerador de Ranking Justo):** Emite nota fria de compatibilidade baseada exclusivamente em habilidades técnicas comprovadas.
* **Stop-Gate Humano:** O recrutador humano avalia os 5 perfis mais compatíveis antes de agendar as entrevistas presenciais.

```json
{
  "nome": "recrutamento-etico",
  "versao": "2.0",
  "nos": [
    { "id": "anonimizar", "tarefa": "remover-dados-pessoais", "dependencias": [] },
    { "id": "avaliar-skills", "tarefa": "comparar-requisitos", "dependencias": ["anonimizar"] },
    { "id": "gerar-ranking", "tarefa": "calcular-aderencia", "dependencias": ["avaliar-skills"] },
    { "id": "gate-recrutador", "tarefa": "revisao-humana", "dependencias": ["gerar-ranking"] }
  ]
}
```

---

## 34. Receita 8: Esteira de Análise de Sentimento de Clientes

### O Problema de Negócios
Centrais de atendimento ao cliente acumulam milhares de chamados de suporte por dia sem conseguir identificar rapidamente clientes de alto valor com risco iminente de cancelamento (*churn*).

### A Composição Modular
* **Nó 1 (Classificador de Urgência):** Identifica mensagens contendo ameaças de cancelamento, processos judiciais ou falhas operacionais graves.
* **Nó 2 (Extrator de Motivo Raiz):** Isola a causa principal da insatisfação (lentidão do sistema, cobrança indevida ou falha no suporte).
* **Nó 3 (Notificador Instantâneo):** Alerta a gerência de contas no canal dedicado do Slack ou Teams em menos de 30 segundos.

```json
{
  "nome": "analise-sentimento-suporte",
  "versao": "2.0",
  "nos": [
    { "id": "classificar-urgencia", "tarefa": "detectar-risco", "dependencias": [] },
    { "id": "isolar-motivo", "tarefa": "extrair-causa-raiz", "dependencias": ["classificar-urgencia"] },
    { "id": "notificar-equipe", "tarefa": "disparar-alerta-slack", "dependencias": ["isolar-motivo"] }
  ]
}
```

---

## 35. Receita 9: Esteira de Migração de Bases de Dados

### O Problema de Negócios
Migrar informações de um banco de dados relacional legado para um banco em nuvem moderno envolve incompatibilidade de tipos de dados, chaves estrangeiras quebradas e risco severo de corrupção.

### A Composição Modular
* **Nó 1 (Inspetor de Tipos):** Mapeia colunas, tamanhos e tipos de campos de origem e destino, identificando incompatibilidades.
* **Nó 2 (Conversor de Esquema):** Gera scripts de migração com garantia de integridade referencial.
* **Nó 3 (Testador de Carga em Sandbox):** Executa a migração em uma base de testes temporária e valida a contagem de registros linha a linha.
* **Stop-Gate Humano:** O Administrador de Banco de Dados (DBA) autoriza o corte final (*cutover*) para produção.

```json
{
  "nome": "migracao-dados-legado",
  "versao": "2.0",
  "nos": [
    { "id": "inspecionar-tipos", "tarefa": "mapear-colunas", "dependencias": [] },
    { "id": "gerar-scripts", "tarefa": "produzir-migracoes", "dependencias": ["inspecionar-tipos"] },
    { "id": "teste-sandbox", "tarefa": "validar-integridade", "dependencias": ["gerar-scripts"] },
    { "id": "gate-dba", "tarefa": "autorizacao-cutover", "dependencias": ["teste-sandbox"] }
  ]
}
```

---

## 36. Receita 10: O Mega-Loop Autônomo

### O Problema de Negócios
Como integrar todas as operações de uma empresa em uma esteira contínua sem que o sistema se torne um monólito ingovernável e impossível de manter?

### A Composição Modular
Através do operador `compor`, o Mega-Loop conecta esteiras inteiras como se fossem engrenagens de um relógio suíço:
* A esteira de **Especificação** gera a partitura para a esteira de **Desenvolvimento TDD**.
* A esteira de **Desenvolvimento** dispara automaticamente a esteira de **Auditoria de Código**.
* A esteira de **Auditoria** alimenta a esteira de **Geração de Manuais**, que atualiza a documentação corporativa em tempo real.

O resultado é uma fábrica digital autônoma, modular, previsível e com custos rigorosamente sob controle.

```json
{
  "nome": "mega-loop-corporativo",
  "versao": "2.0",
  "nos": [
    { "id": "loop-especificacao", "tarefa": "subloop-especificacao", "dependencias": [] },
    { "id": "loop-desenvolvimento", "tarefa": "subloop-tdd", "dependencias": ["loop-especificacao"] },
    { "id": "loop-auditoria", "tarefa": "subloop-qualidade", "dependencias": ["loop-desenvolvimento"] },
    { "id": "loop-publicacao", "tarefa": "subloop-manuais", "dependencias": ["loop-auditoria"] }
  ]
}
```

### Rastreabilidade Técnica da Parte V

* `internal/process/fabricar.go` — Motor de instanciação das receitas no cofre central.
* `internal/dag/kahn.go` — Validador topológico de precedência entre receitas aninhadas.
* `internal/gate/gate.go` — Mecanismo mecânico de parada e auditoria para Stop-Gates humanos.