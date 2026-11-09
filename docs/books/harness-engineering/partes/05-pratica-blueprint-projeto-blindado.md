# Capítulo 5: O Blueprint — Como Iniciar Qualquer Projeto Blindado do Zero

Agora que a causa raiz foi exposta, vamos construir o modelo mental prático: **como você deve iniciar qualquer projeto no futuro para que ele já nasça 100% blindado contra a preguiça do LLM?**

Esqueça começar criando arquivos `.cursorrules` ou prompts gigantes no chat. O ponto de partida de um projeto AIDD profissional é a **infraestrutura mecânica de validação**.

---

## Passo 1: O "Guard Rail" Pré-Commit (A Primeira Linha de Defesa)

Antes de escrever a primeira linha de código do produto, configure o seu ambiente para impedir commits ruins. No terminal, crie um script que roda antes de cada alteração:

```bash
#!/usr/bin/env bash
# .git/hooks/pre-commit ou tools/gate.sh

set -e

echo "=== [HARNESS GATE] Verificação Mecânica Inegociável ==="

# 1. Checagem de Zero Stubs (busca mecânica por TODOs ou pass em código de prod)
if grep -rn "TODO: implementar" internal/ src/; then
    echo "ERRO FATAL: Encontrado stub em código de produção!"
    exit 1
fi

# 2. Formatação estrita
if [ -n "$(gofmt -l .)" ]; then
    echo "ERRO FATAL: Arquivos com formatação desalinhada!"
    exit 1
fi

# 3. Testes unitários com exit code auditado
go test -count=1 ./...

echo "=== [HARNESS GATE] APROVADO: O código é real e está testado. ==="
```

Quando esse script está ativo, o modelo não consegue dizer *"Pronto, já implementei"* se ele deixou um mock provisório ou um teste quebrando. O próprio terminal aborta o processo.

---

## Passo 2: O Schema Antes da Função (Type-First Contracts)

Nunca peça para o modelo: *"Crie um serviço de autenticação"*.  
Em vez disso, defina primeiro os tipos e as validações de erro:

```go
type UserRegistrationRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (r UserRegistrationRequest) Validate() error {
    if !strings.Contains(r.Email, "@") {
        return ErrInvalidEmail
    }
    if len(r.Password) < 8 {
        return ErrWeakPassword
    }
    return nil
}
```

Quando o contrato existe e tem um método `.Validate()`, você instrui o agente:
> *"Implemente a lógica de registro. Ela DEVE obrigatoriamente chamar `req.Validate()`. Se falhar, retorne o erro imediatamente sem tocar no banco."*

O modelo não tem espaço para inventar campos ou esquecer validações básicas.

---

## Passo 3: O Padrão de 3 Fases (Intake -> Spec -> Code)

O maior erro dos desenvolvedores que se frustram com IA é pular da ideia direto para o código. Isso é garantia de retrabalho.

O fluxo profissional segue o funil estrito:
1. **Intake (O Porquê):** O que queremos resolver e quem se beneficia (INTENT). Sem falar de código.
2. **Spec (O Quê):** Quais são os contratos de dados, as assinaturas das funções e os critérios binários Given-When-Then (SPEC).
3. **Execução em Tickets Cirúrgicos (O Como):** Fatiar a SPEC em tickets pequenos (máximo 50 linhas de código por ticket), onde cada ticket tem seu próprio teste unitário associado.

---

## Passo 4: O Isolamento de Blast Radius (Workspaces Efêmeros)

Em projetos grandes, nunca deixe o agente trabalhar solto na branch principal (`main`). 

Utilize sempre **Git Worktrees** ou pastas temporárias isoladas. Cada fatia de código é desenvolvida em uma pasta separada. Se o modelo fizer bobagem, você simplesmente deleta a pasta daquela fatia em 1 segundo (`git worktree remove`). O seu código principal nunca é corrompido.
