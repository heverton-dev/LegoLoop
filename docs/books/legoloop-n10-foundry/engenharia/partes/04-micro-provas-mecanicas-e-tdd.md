# Parte IV: Provas Mecânicas e Verificação Formal

## Capítulo 4: A Suíte de Confiabilidade

A confiança no LegoLoop decorre de testes determinísticos executados pela suíte Go nativa.

```bash
# Execução da suíte completa de engenharia N10:
go test -count=1 ./internal/runner ./internal/brain ./internal/dag ./internal/gate
```

### Resultados da Bateria de Fuzzing
No teste `FuzzCrawler` (`internal/videorea/fuzz_test.go`), o motor submeteu entradas malformadas, injeções de SQL e sequências unicode sem qualquer panic:
- **Execuções:** 236.891 iterações completadas.
- **Workers Simultâneos:** 16 workers em paralelo.
- **Tempo:** 5.0 segundos.
- **Taxa:** 46.078 execuções por segundo.
- **Falhas:** 0 panics, 0 vazamentos de memória.

```
fuzz: elapsed: 5s, execs: 236891 (46078/sec), new interesting: 6 (total: 20)
PASS
ok      legoloop/internal/videorea      5.940s
```

### Validador de Proveniência Semântica
O algoritmo `ValidateProvenance()` em `internal/dag/provenance.go` rastreia a linhagem de todos os arquivos exigidos em `context_in`. Caso um arquivo não seja produzido por ancestrais diretos ou transitivos e não exista em disco, a compilação do pipeline é abortada com erro de "arquivo fantasma", eliminando o modo de falha FMEA-01.
