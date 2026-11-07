# PARTE V — Arquitetura: Padrões e Anti-Padrões

## 24. Os Sete Pecados Capitais do Desenvolvimento com IA

Ao longo de anos acompanhando equipes de tecnologia tentando incorporar inteligência artificial em suas rotinas produtivas, observamos que a esmagadora maioria dos fracassos não decorre de limitações dos modelos matemáticos, mas sim de vícios metodológicos de projeto. Reunimos a seguir os sete anti-padrões mais destrutivos da indústria:

### 1. O Prompt Monolítico da Esperança
Consiste em escrever um prompt de três páginas contendo requisitos de negócios, regras de formatação, tratamento de exceções e exemplos de código, colá-lo no chat e torcer para que a IA acerte tudo de primeira. Esse método ignora que a probabilidade de acerto composto decai exponencialmente a cada nova restrição não testada.

### 2. A Ilusão do Piloto Automático Sem Stop-Gates
Acreditar que um agente autônomo pode alterar bancos de dados de produção ou enviar e-mails para clientes sem nenhum ponto de inspeção humana intermediária. No LegoLoop, decisões irreversíveis exigem compulsoriamente a validação de um stop-gate humano (`SG-*`).

### 3. A Injeção Brutal de Ferramentas (Tool Overload)
Conectar dezenas de servidores MCP na inicialização do assistente, queimando 120.000 tokens antes do primeiro turno. O Meta-MCP Router do LegoLoop elimina esse vício através de descoberta sob demanda (*lazy discovery*).

### 4. A Síndrome do "Código que Compila, mas Está Vazio" (Stubs)
Aceitar códigos gerados por IA contendo funções simuladas com `TODO`, `pass` ou `// implementar depois`. A Lei #1 do LegoLoop estabelece tolerância zero para stubs: o código entregue deve ser real, hermético e funcional.

### 5. O Desenvolvimento Sem Testes Prévios (Anti-TDD)
Pedir à inteligência artificial que implemente a funcionalidade antes de escrever o teste unitário correspondente. Sem um teste automatizado que comprove a falha inicial, é impossível garantir que o modelo não gerou um teste viciado que sempre passa.

### 6. A Poluição Prolixa de Saídas
Permitir que assistentes respondam com longos parágrafos de conversa fiada ("Com certeza! Vou ficar muito feliz em ajudar..."). O protocolo STE-PT impõe comunicação telemétrica enxuta e orientada a evidências.

### 7. A Ignorância da Quebra Linguística
Escrever prompts longos em português sem considerar que o algoritmo de tokenização (BPE) cobra até o triplo por palavras acentuadas. O LegoLoop combate isso aplicando diretrizes de token-economy extrema.

---

## 25. Padrões de Projeto para Esteiras de IA de Alta Resiliência

Para substituir os sete pecados capitais por práticas industriais auditáveis, o LegoLoop adota quatro padrões fundamentais da arquitetura de software moderna adaptados para o universo dos Grandes Modelos de Linguagem:

### Padrão 1: Tubos e Filtros (Pipe-and-Filter)
Cada nó da esteira funciona como um filtro isolado: ele lê dados do cano de entrada, executa sua transformação matemática estrita e despeja o resultado no cano de saída. Nenhum filtro conhece a existência dos filtros anteriores ou posteriores. Isso permite trocar qualquer etapa da esteira sem risco de efeito cascata.

### Padrão 2: Disparo Traçante (Tracer Bullet)
Em vez de tentar construir o sistema inteiro de uma vez, a esteira inicial do LegoLoop implementa um "disparo traçante": uma fatia vertical mínima que percorre todas as camadas do sistema (especificação, teste, código e validação). Assim que o tiro traçante atinge o alvo com sucesso, a esteira é expandida progressivamente.

### Padrão 3: Disjuntor Anti-Loop (Circuit Breaker)
Se um nó entrar em um ciclo vicioso de tentativa e erro (por exemplo, tentando compilar um código quebrado repetidas vezes), o disjuntor do LegoLoop desarma automaticamente após três tentativas infrutíferas, congela a esteira e emite um alerta sonoro e visual no Studio para intervenção humana.

### Padrão 4: Observabilidade Estéril via Sidecar
O monitoramento do consumo de tokens, tempo de execução e cobertura de testes não é executado pelo próprio modelo de IA (o que consumiria ainda mais contexto). Ele é realizado de forma passiva por processos secundários (*sidecars*) do binário em Go, que capturam telemetria do sistema operacional sem gastar um único token de API.

```{=typst}
#painel("Regra da Engenharia Robusta")[
  Um sistema de inteligência artificial de qualidade corporativa não é aquele que nunca erra; é aquele cujas falhas são contidas imediatamente pelo próprio desenho da arquitetura, sem jamais contaminar o restante da operação.
]
```