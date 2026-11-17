## 9. Operação Prática, Telemetria & Troubleshooting

## 9.1 Guia de Campo: Execução via Linha de Comando (CLI)

O operador pode disparar o pipeline diretamente pelo terminal através do executável compilado do LegoLoop:

### Disparo Padrão com Variável de Entrada
```bash
legoloop run loops/video-edition/pipeline.dag.json --var input_video="entrada/minha_aula.mp4"
```

### Disparo Customizado com Ajuste Fino de Parâmetros
```bash
legoloop run loops/video-edition/pipeline.dag.json \
  --var input_video="gravacao_bruta.mov" \
  --var pre_roll_ms=50 \
  --var ducking_db="-20dB" \
  --var eye_level_ratio=0.35
```

---

## 9.2 Matriz de Resolução de Falhas (Troubleshooting)

A tabela abaixo mapeia os principais incidentes operacionais encontrados em produção e o protocolo imediato de remediação:

| Incidente Operacional                  | Causa Raiz Provável                       | Ação Corretiva Imediata                       |
| :------------------------------------- | :---------------------------------------- | :-------------------------------------------- |
| **Corte inicial engolindo sílaba plosiva** | Locutor tem ataque muito rápido ou o áudio tem reverberação | Aumente o pre-roll de entrada configurando `--var pre_roll_ms=55`. |
| **Rosto descentralizado durante o zoom**   | Baixa iluminação ou apresentador virou o rosto de perfil | Edite `scratch/yunet_tracking.json` e ative o fallback de centro geométrico. |
| **Trilha sonora abafando a voz do orador** | Compressão sidechain fraca                | Intensifique o ducking passando a variável `--var ducking_db="-22dB"`. |
| **Erro de render no nó Remotion (n09)**    | Node.js sem memória ou pacote NPM ausente | Rode `npm install` dentro de `internal/videomotion` e reinicie o nó. |
| **Timeline FCPXML sem áudio no DaVinci**   | Taxa de amostragem incompatível           | Garanta que o áudio extraído no nó `n01` esteja rigorosamente em 48.000 Hz WAV. |
| **Nota de gravação vazou no corte final**  | O locutor não usou marcadores verbais claros | Edite manualmente o arquivo `scratch/cuts_manifest.json` no Stop-Gate 01. |

---

## 9.3 Telemetria & Eficiência Financeira

Ao término da execução, o motor emite o certificado com o balanço consolidado de recursos:

```
=============================================================
             RELATÓRIO DE TELEMETRIA LEGOLOOP
=============================================================
• Duração do Processamento:       4m 12s
• Tokens Cognitivos Consumidos:   12.400 in / 3.200 out
• Custo de Inteligência (API):    $ 0.041 (R$ 0,24)
• Custo de Infraestrutura Local:  R$ 0,00 (Hardware próprio)
• Economia frente ao SaaS:        94.2% de redução de custo
=============================================================
```

#painel-sucesso("EFICIÊNCIA COMPROVADA", [
  A esteira `video-edition` alcança 94.2% de economia frente a soluções em nuvem proprietárias, consolidando o retorno financeiro imediato da infraestrutura local sem qualquer cobrança por minuto processado.
])
