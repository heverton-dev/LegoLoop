# 🎯 INTENT: Matriz Canônica de Papéis e Casamento Determinístico de Peças

> **Status:** APPROVED  
> **Data:** 2026-11-10  
> **Autor:** Antigravity  
> **Iniciativa:** Fábrica Determinística LegoLoop  
> **Alvo:** internal/taxonomy, internal/matriz, internal/generator, cmd/legoloop  
> **Sessão:** b92a1402-a673-42dc-b9a5-2b06b573ed56  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** O comando legoloop fabricar busca ferramentas por termos textuais vagos. O sistema pode escolher ferramentas incorretas para os nós do grafo. Não existe catálogo formal de funções operacionais.
* **Quem Sofre:** Desenvolvedores e agentes que operam o comando legoloop fabricar e auditam a Matriz de Peças.
* **O Porquê:** Uma linha de montagem industrial precisa de papéis claros. Cada estação de trabalho requer uma especialidade comprovada. A ferramenta correta deve casar com o crachá do operador.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Peças no cofre não informam a qual função atendem. O resolvedor do processo usa busca aproximada. O risco de alucinação de ferramenta é alto.
  * **Amanhã:** 18 papéis canônicos distribuídos em 7 guildas governam o sistema. Cada peça declara papel primário e papéis secundários. O comando legoloop fabricar seleciona a melhor ferramenta com precisão mecânica. O sistema para imediatamente se encontrar um crachá desconhecido.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Criar o catálogo oficial dos 18 papéis canônicos em 7 guildas em internal/taxonomy e no dicionário JSON oficial.
* [ ] Objetivo 2: Atualizar a estrutura da Matriz de Peças para exigir papel_primario e permitir papeis_secundarios.
* [ ] Objetivo 3: Implementar disjuntor estrito que emite código 1 se encontrar papel fora do catálogo oficial.
* [ ] Objetivo 4: Adicionar filtros --papel e --guilda no comando legoloop matriz ver.
* [ ] Objetivo 5: Implementar o casamento determinístico em 2 níveis no resolvedor do comando legoloop fabricar.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em Alvo.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não criaremos papéis livres ou não homologados fora das 7 guildas oficiais.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| PECA-01 | skill | skills/legoloop-matriz | Consultar catálogo e reputação de peças no cofre. |
| PECA-02 | unit_task | internal/taxonomy | Validar vocabulário e papéis autorizados. |
| PECA-03 | unit_task | internal/matriz | Gerenciar tags de papéis e filtros por crachá. |
| PECA-04 | unit_task | internal/generator | Casar nó do processo com ferramenta compatível. |
| PECA-05 | script | go test ./... | Regressão mecânica e integridade de ponta a ponta. |
| PECA-06 | stop_gate | go build -o bin/legoloop.exe ./cmd/legoloop | Recompilar binário com atualização do PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** manual (via CLI legoloop matriz ver e comando legoloop fabricar).
* **Nível de Risco & Blast Radius:** Baixo (workspace local, focado nos pacotes 	axonomy, matriz e generator).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente determinístico local. Zero chamadas a redes externas para resolução de catálogo.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Execução de legoloop matriz ver --papel engenheiro-frontend filtra peças com sucesso. Execução com papel inválido retorna código 1. Suíte Go 100% verde.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * O catálogo contém exatamente 18 papéis ortogonais em 7 guildas.
  * Peças existentes sem tag receberão migração com seus papéis correspondentes.
* **Riscos e Perguntas em Aberto:** Nenhum. Todas as decisões técnicas fecharam na entrevista socrática.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes passam com zero stubs e zero mocks. | Suíte verde e busca por stubs vazia | go test -count=1 ./... com exit code 0 |
| 2 | Validação de papéis desconhecidos falha de forma mecânica. | Teste unitário de disjuntor com código de erro sentinela | go test ./internal/taxonomy/... verde |
| 3 | Filtro por papel na CLI retorna apenas peças autorizadas. | Execução do comando com flag --papel | legoloop matriz ver --papel <slug> funcional |
| 4 | O comando fabricar prioriza papel primário antes do secundário. | Teste de casamento determinístico no gerador | go test ./internal/generator/... verde |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md](../specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md)
* **Gate de Transição:** Este Intent foi alinhado com 100% de confiança da entrevista socrática.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (TEMPLATE_SPEC.md).*
