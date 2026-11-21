# Parte III — Visão Micro: Schemas e Gates

## 1. Schema `memoria_execucao.v1`

```json
{
  "schema_version": "memoria_execucao.v1",
  "processo": "meu-processo",
  "sessoes": {
    "abc123": {
      "processo": "meu-processo",
      "sessao_id": "abc123",
      "pid": 4321,
      "etapa_atual": "M1.2",
      "proximo_passo": "M1.3",
      "atualizado_em": "2026-11-16T22:00:00Z"
    }
  }
}
```

Regras: chave do mapa é o `sessao_id`; `pid` positivo é obrigatório; escrita é
atômica; leitura corrompida gera backup `.legoloop-bak` e sentinela.

## 2. Schema `anticorpos.v1`

```json
{
  "schema_version": "anticorpos.v1",
  "globais": [],
  "dominio": {
    "meu-processo": [
      {
        "pacote": "internal/triagem",
        "causa_raiz": "lock sem PID",
        "assinatura": "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
        "antidoto_aplicado": "limparLockOrfao"
      }
    ]
  }
}
```

Regras: `assinatura` é SHA-256 do erro no formato `sha256:<hex>`; carga legado
é migrada sem perda; banco corrompido vira backup e o sistema segue vazio.

## 3. Gate anti-reincidência

```go
acao, antidoto, err := db.GateAntiReincidencia(pkg, errSig)
```

| Ação                | Quando                           | Resultado                   |
| :------------------ | :------------------------------- | :-------------------------- |
| `AcaoFluxoNormal`     | Erro desconhecido                | Executa o nó                |
| `AcaoAntidoto`        | Erro conhecido com antídoto      | Aplica o antídoto           |
| `AcaoBloqueado`       | Erro conhecido sem antídoto novo | `ErrGateAntiReincidencia`     |

## 4. Lock por arquivo

- O lock guarda o PID do detentor.
- PID morto vira lock órfão e o runtime toma o lock.
- `ErrLockOcupado` significa detentor vivo: aguarde.
