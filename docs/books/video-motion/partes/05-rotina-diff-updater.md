# Parte V — Rotina de Atualização por Diff (LegoPress)

## Governança Viva de Documentação

No ecossistema LegoLoop, nenhuma documentação pode se tornar obsoleta em relação ao código real. Sempre que um arquivo do pipeline (`internal/videomotion/` ou `loops/video-motion/`) for alterado via commit ou refatoração, o motor **LegoPress** sincroniza a obra técnica automaticamente.

## Comandos de Manutenção do Livro

```bash
# Auditar a integridade geométrica do PDF gerado
legoloop legopress auditar docs/books/video-motion/2026-11-07_PIPELINE-MOTION-VIDEO.pdf

# Atualizar capítulos a partir dos diffs recentes do Git
legoloop legopress atualizar docs/books/video-motion/

# Recompilar o livro completo em caso de alterações nas partes
bash docs/books/compilar.sh video-motion
```

A esteira de publicação gera automaticamente o manuscrito consolidado em Markdown, o documento tipográfico final em Typst, o PDF de alta resolução (170x240mm) e as imagens de visualização para cada página na pasta `preview/`.
