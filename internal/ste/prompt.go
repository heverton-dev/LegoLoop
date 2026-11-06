package ste

import (
	"fmt"
	"strings"
)

// PromptRegra gera o documento markdown de regras com base no perfil.
func PromptRegra(perfil ConfigPerfil) string {
	if perfil.Modo == ModoOff {
		return ""
	}
	var b strings.Builder
	b.WriteString("# 🛡️ REGRAS DE GOVERNANÇA E ECONOMIA STE-PT\n\n")
	b.WriteString(blocoInputImperativo())
	if perfil.Modo == ModoCaveman || perfil.Modo == ModoHibrido || perfil.Modo == ModoStept {
		b.WriteString(blocoCavemanThinking())
	}
	if perfil.Modo == ModoCaveman {
		b.WriteString(blocoOutputCaveman())
	}
	if perfil.Modo == ModoHibrido || perfil.Modo == ModoStept {
		b.WriteString(blocoOutputTriade())
		b.WriteString(blocoRegrasConcisao(perfil.Pausa))
	}
	return b.String()
}

func blocoInputImperativo() string {
	return `## 1. INPUT IMPERATIVO (Economia de Contexto)
- Isole comandos imperativos e substantivos de ação do usuário.
- Ignore saudações e conversas laterais irrelevantes para a execução técnica.
- Nunca faça echo desnecessário da instrução inteira do usuário na resposta.

`
}

func blocoCavemanThinking() string {
	return `## 2. CAVEMAN THINKING (Raciocínio Interno Econômico)
- Use estilo telegráfico extremo no pensamento interno (sem artigos ou preposições).
- Frases de no máximo 3 a 5 linhas para tarefas pontuais.
- Nunca carregue arquivos de log inteiros no contexto do pensamento.
- Abreviar termos técnicos e ir direto à decisão de código.

`
}

func blocoOutputCaveman() string {
	return `## 3. OUTPUT TELEGRÁFICO (Modo Caveman)
- Respostas ultra-curtas, secas e telegráficas no chat.
- Limite de saída estrito: direto ao ponto e sem enrolação.
- Prefira diffs objetivos ou listas sucintas a parágrafos explicativos.

`
}

func blocoOutputTriade() string {
	return `## 3. CONTRATO DE RESPOSTA OBRIGATÓRIO (TRÍADE STE-PT & CARTÃO DE AÇÃO STE-PT)
TODA resposta técnica DEVE seguir a hierarquia de negrito limpa, compacta (quebra direta sem linhas em branco extras) e sem blocos de código artificiais:

### FORMATO PRÉ-AÇÃO (Antes de alterar código ou executar comandos com efeito):
**[S] Situação:** <Cenário factual do problema ou contexto encontrado>
**[T] Teoria:** <Hipótese técnica e plano de ação antes de alterar código>
**[E] Execução:**
· **Ação:** <O que vai ser feito com verbo imperativo e arquivos tocados>
· **Rollback:** <Comando exato de recuperação caso falhe>
· **Decisão:** **[1] seguir** | **[2] personalizar** | **[0] cancelar**

### FORMATO PÓS-AÇÃO (Após executar e validar alterações):
**[S] Situação:** <O que foi executado no sistema>
**[T] Tese:** <A tese técnica comprovada pela alteração>
**[E] Evidência:**
· **Testes:** <Comando de teste executado e resultado real>
· **Build:** <Status de compilação ou verificação real do sistema>
· **Verificação:** <Resultado mensurável final>

`
}

func diretivaPausa(pausa Pausa) string {
	switch pausa {
	case PausaNunca:
		return "Exiba a execução e prossiga sem aguardar aprovação."
	case PausaSempre:
		return "Pare em [E] Execução e aguarde a aprovação do usuário para todas as ações com efeito."
	case PausaDestrutivo:
		fallthrough
	default:
		return "Pare em [E] Execução e aguarde a aprovação do usuário apenas em ações destrutivas (deleção, sobrescrita, git reset/clean, drop de dados)."
	}
}

func blocoRegrasConcisao(pausa Pausa) string {
	var b strings.Builder
	b.WriteString("## 4. DIRETIVAS DE CONCISÃO & PAUSA\n")
	b.WriteString("- PROIBIDO introduções, preâmbulos, saudações ou explicações prolixas antes de [S].\n")
	b.WriteString("- Quando o usuário responder confirmações (ex.: '1', 'seguir', 'sim', 'STEPT 1'), vá DIRETO à execução e responda com o formato PÓS-AÇÃO ou próximo PRÉ-AÇÃO.\n")
	b.WriteString("- Respostas estritamente limitadas a 1 a 3 linhas por bloco.\n")
	b.WriteString("- PROIBIDO re-resumir artefatos ou diffs após a entrega.\n")
	fmt.Fprintf(&b, "- **Diretiva de Pausa:** %s\n\n", diretivaPausa(pausa))
	return b.String()
}
