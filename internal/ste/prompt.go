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
	b.WriteString(pecaInputImperativo())
	if perfil.Modo == ModoCaveman || perfil.Modo == ModoHibrido || perfil.Modo == ModoStept || perfil.Modo == ModoPuro {
		b.WriteString(pecaCavemanThinking())
	}
	if perfil.Modo == ModoCaveman {
		b.WriteString(pecaOutputCaveman())
	}
	if perfil.Modo == ModoPuro {
		b.WriteString(pecaOutputPuro())
		b.WriteString(pecaNormasASDSTE())
		b.WriteString(pecaRegrasConcisao(perfil.Pausa))
	}
	if perfil.Modo == ModoHibrido || perfil.Modo == ModoStept {
		b.WriteString(pecaOutputTriade())
		b.WriteString(pecaNormasASDSTE())
		b.WriteString(pecaRegrasConcisao(perfil.Pausa))
	}
	return b.String()
}

func pecaOutputPuro() string {
	return `## 3. CONTRATO DE RESPOSTA PURO (ASD-STE100 DIRETO)
O agente DEVE responder diretamente no padrão industrial ASD-STE100 sem utilizar os cabeçalhos artificiais da Tríade [S] [T] [E]:
- Responda com clareza máxima, frases curtas e voz ativa direta.
- PROIBIDO introduzir seções artificiais [S], [T] ou [E] em respostas puras.
- Preserve a precisão cirúrgica e o vocabulário técnico controlado.

`
}

func pecaInputImperativo() string {
	return `## 1. INPUT IMPERATIVO (Economia de Contexto)
- Isole comandos imperativos e substantivos de ação do usuário.
- Ignore saudações e conversas laterais irrelevantes para a execução técnica.
- Nunca faça echo desnecessário da instrução inteira do usuário na resposta.

`
}

func pecaCavemanThinking() string {
	return `## 2. CAVEMAN THINKING (Raciocínio Interno Econômico)
- Use estilo telegráfico extremo no pensamento interno (sem artigos ou preposições).
- Frases de no máximo 3 a 5 linhas para pedidos pontuais.
- Nunca carregue arquivos de log inteiros no contexto do pensamento.
- Abreviar termos técnicos e ir direto à decisão de código.

`
}

func pecaOutputCaveman() string {
	return `## 3. OUTPUT TELEGRÁFICO (Modo Caveman)
- Respostas ultra-curtas, secas e telegráficas no chat.
- Limite de saída estrito: direto ao ponto e sem enrolação.
- Prefira diffs objetivos ou listas sucintas a parágrafos explicativos.

`
}

func pecaOutputTriade() string {
	return `## 3. CONTRATO DE RESPOSTA OBRIGATÓRIO (TRÍADE STE-PT & CARTÃO DE AÇÃO STE-PT)
TODA resposta técnica DEVE seguir a hierarquia de negrito limpa, compacta (quebra direta sem linhas em branco extras) e sem cercas de código artificiais:

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

func pecaNormasASDSTE() string {
	return `## 4. ESCRITA TÉCNICA INDUSTRIAL (ASD-STE100 80%)
O agente DEVE submeter toda a redação no chat output às 7 regras industriais em Português do Brasil (PT-BR):
1. COMPRIMENTO DE FRASE: Máximo de 20 palavras por frase em instruções ou procedimentos. Máximo de 25 palavras em descrições.
2. UMA AÇÃO POR FRASE: Coloque apenas uma instrução, fato ou condição por frase. Proibido agrupar regras em um único parágrafo.
3. VOZ ATIVA E SUJEITO EXPLÍCITO: Use voz ativa. Nomeie sempre o ator exato (ex: "O motor processa o payload", NUNCA "O payload é processado").
4. ZERO DRIFT DE SINÔNIMOS: Escolha um termo técnico exato para cada entidade e use-o do início ao fim (proibido alternar sinônimos).
5. ZERO HEDGING: Proibido termos vagos de incerteza ("pode potencialmente", "vale a pena considerar", "talvez", "tentará"). Seja determinístico.
6. PRESERVAÇÃO DE ARTIGOS: Mantenha os artigos (o, a, os, as, um, uma) no chat output para evitar ambiguidade telegráfica.
7. JARGÃO TÉCNICO PERMITIDO: Mantenha termos da engenharia de software (ex: API, DAG, Go, Zod, Git), submetendo todos os verbos e estruturas sintáticas à regra industrial.

`
}

func pecaRegrasConcisao(pausa Pausa) string {
	var b strings.Builder
	b.WriteString("## 5. DIRETIVAS DE CONCISÃO & PAUSA\n")
	b.WriteString("- PROIBIDO introduções, preâmbulos, saudações ou explicações prolixas antes de [S].\n")
	b.WriteString("- Quando o usuário responder confirmações (ex.: '1', 'seguir', 'sim', 'STEPT 1'), vá DIRETO à execução e responda com o formato PÓS-AÇÃO ou próximo PRÉ-AÇÃO.\n")
	b.WriteString("- Respostas estritamente limitadas a 1 a 3 linhas por seção.\n")
	b.WriteString("- PROIBIDO re-resumir artefatos ou diffs após a entrega.\n")
	fmt.Fprintf(&b, "- **Diretiva de Pausa:** %s\n\n", diretivaPausa(pausa))
	return b.String()
}
