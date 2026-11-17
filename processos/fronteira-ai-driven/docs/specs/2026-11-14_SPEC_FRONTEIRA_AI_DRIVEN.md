# SPEC — Arquitetura de Fronteira AI-Driven (2026-11-14)

## Critérios de Aceitação (EARS)
- **REQ-01 (Lazy de Regras):** Quando um nó possuir tags de domínio (ex: `go`, `ui`, `seguranca`), o `Orchestrix` deve injetar apenas as regras correspondentes às tags além do núcleo mínimo.
- **REQ-02 (Roteamento Cognitivo):** Quando `cognitive_complexity` for baixa, o roteador deve selecionar modelo rápido/econômico; quando alta, deve selecionar modelo de fronteira com thinking habilitado.
- **REQ-03 (Mutação e Rollback):** Se o gate binário ou a verificação de mutação falhar na baia, o `runner` deve executar rollback atômico do worktree sem contaminar a branch principal.
- **REQ-04 (Destilação de Memória):** Ao concluir `legoloop fechar` após o Stop-Gate final, o motor deve destilar erros recorrentes em gates mecânicos e limpar a memória episódica do processo.
- **REQ-05 (Segurança de Baia):** Antes de injetar `context_in` externo, o motor deve sanitizar o conteúdo contra prompt injection e restringir ferramentas à allowlist do nó dentro do Path Jail.
- **REQ-06 (Bootstrap Único):** Ao executar `legoloop bootstrap`, o CLI deve instalar/validar binário, skills, `mcp-router`, hooks e ambiente sem configuração manual.
