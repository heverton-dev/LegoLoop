# Parte IV — Prática: Operação do Dia a Dia

## 1. Consultar o progresso

```powershell
legoloop sessoes meu-processo --json
```

Saída esperada com duas sessões vivas:

```json
[
  {
    "sessao_id": "abc123",
    "pid": 4321,
    "etapa_atual": "M1.2",
    "proximo_passo": "M1.3"
  },
  {
    "sessao_id": "def456",
    "pid": 4322,
    "etapa_atual": "M2.1",
    "proximo_passo": "M2.2"
  }
]
```

Processo sem sessão viva devolve `[]` e exit 0.

## 2. Antes de executar um nó

1. Calcule a assinatura do erro possível.
2. Chame `brain.GateAntiReincidencia`.
3. `AcaoFluxoNormal`: execute o nó.
4. `AcaoAntidoto`: aplique o antídoto e execute.
5. `AcaoBloqueado`: pare e reporte a causa raiz.

## 3. Depois de um erro novo

1. Registre com `brain.RecordComEscopo`.
2. Escolha o escopo: global para erro de ecossistema, domínio para erro local.
3. Gravar o antídoto transforma bloqueio futuro em `AcaoAntidoto`.

## 4. No Fechamento Canônico

1. `fechar.ZerarSessao` roda no quarteto Lei 11.
2. O envelope volta sem sessões e com `schema_version` correto.
3. Confira com `legoloop sessoes <slug> --json`: deve devolver `[]`.

## 5. Roteiro de recuperação

| Sintoma | Ação |
| :--- | :--- |
| `ErrLockOcupado` | Aguarde o detentor ou confira o `pid` do lock |
| Lock órfão | Nada a fazer: o runtime toma o lock |
| `ErrMemoriaCorrompida` | Confira o backup `.legoloop-bak` |
| `ErrAnticorposCorrompidos` | Confira o backup e regrave os anticorpos |
| `ErrGateAntiReincidencia` | Corrija a causa raiz e registre o antídoto |
