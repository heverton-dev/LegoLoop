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
	if perfil.Modo == ModoCaveman || perfil.Modo == ModoHibrido {
		b.WriteString(blocoCavemanThinking())
	}
	if perfil.Modo == ModoCaveman {
		b.WriteString(blocoOutputCaveman())
	}
	if perfil.Modo == ModoHibrido {
		b.WriteString(blocoOutputHibrido())
	}
	if perfil.Modo == ModoStept || perfil.Modo == ModoHibrido {
		b.WriteString(blocoCartaoStept(perfil.Pausa))
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

func blocoOutputHibrido() string {
	return `## 3. CONCISÃO DE OUTPUT (Modo Híbrido - Chat Econômico)
- Respostas conversacionais ultra-concisas e diretas (1 a 3 linhas no máximo).
- PROIBIDO re-resumir artefatos, arquivos criados ou diffs de código após a execução.
- Teto estrito de saída: no máximo 300 tokens para respostas normais e 600 tokens para explicações técnicas.
- Use listas curtas ou tabelas em vez de parágrafos prolixos.
- Emita o Cartão de Ação STE-PT antes de qualquer alteração com efeito.

`
}

func diretivaPausa(pausa Pausa) string {
	switch pausa {
	case PausaNunca:
		return "Exiba o cartão a título informativo e prossiga sem aguardar aprovação."
	case PausaSempre:
		return "Pare e aguarde a aprovação do usuário para todas as ações com efeito."
	case PausaDestrutivo:
		fallthrough
	default:
		return "Pare e aguarde a aprovação do usuário apenas em ações destrutivas (deleção, sobrescrita, git reset/clean, drop de dados)."
	}
}

func blocoCartaoStept(pausa Pausa) string {
	var b strings.Builder
	b.WriteString("## 4. CARTÃO DE AÇÃO STE-PT (ASD-STE100)\n")
	b.WriteString("Antes de executar comandos de escrita ou alterações com efeito, emita o cartão:\n\n")
	b.WriteString("```\n")
	b.WriteString("AÇÃO: <Verbo imperativo + Alvo>\n")
	b.WriteString("AVISO: <Impacto de risco ou vazio se seguro>\n")
	b.WriteString("O QUE VAI ACONTECER: <Lista de passos concretos>\n")
	b.WriteString("EFEITO: <Arquivos tocados ou estado final>\n")
	b.WriteString("SE DER ERRADO: <Estratégia de rollback ou recuperação>\n")
	b.WriteString("DECISÃO: [1] seguir  [2] personalizar  [0] cancelar\n")
	b.WriteString("```\n\n")
	fmt.Fprintf(&b, "**Diretiva de Pausa:** %s\n\n", diretivaPausa(pausa))
	return b.String()
}
