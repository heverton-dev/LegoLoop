# Parte IV — Fábrica de Governança e Golden Fixtures

## Da Curadoria ao Software Executável

Assim que o operador aprova a curadoria no Stop-Gate, o pipeline convoca a esteira de fabricação do LegoLoop:


![Diagrama de Arquitetura](/img/diagrama_373f383a.png){width=85%}


### 1. Emissão da Tríade ISO de Governança

A fábrica gera os quatro artefatos canônicos sincronizados:
- **`YYYY-MM-DD_INTENT_<SLUG>.md`:** O propósito de negócio e o valor entregue.
- **`YYYY-MM-DD_SPEC_<SLUG>.md`:** Contratos formais de dados e cenários BDD binários.
- **`YYYY-MM-DD_PLAN_<SLUG>.md`:** Decomposição em tickets verticais com rastreabilidade 1:1.
- **`YYYY-MM-DD_ROUTE_<SLUG>.json`:** O grafo de execução estéril com hashes SHA-256 de integridade.

### 2. A Golden Fixture Sintética (3 Segundos)

Testar pipelines de vídeo não pode depender de renderizar clipes de 10 minutos:
- A fábrica sintetiza um clipe MP4 ultraleve de 3.0 segundos a 60fps.
- Insere uma faixa de áudio de 48kHz contendo um pico vocal para validar o silêncio de ataque (+40ms).
- Permite que a suíte de testes TDD do novo pipeline rode localmente em menos de 2 segundos.
