package interview

// MsgConcluir é o pedido padronizado para fechar uma entrevista que estacionou abaixo da meta,
// quando o usuário já respondeu tudo. O motor não inventa restrição: usa só o que foi decidido.
const MsgConcluir = "O usuário já respondeu e aprovou todas as decisões da entrevista. Não há lacuna aberta. " +
	"Emita agora o ESPELHO CANÔNICO final com os seis campos (Resultado Esperado, Quem se Beneficia, Por que Agora, " +
	"Critério de Sucesso, Restrições Inegociáveis, Fora de Escopo) e a CONFIANÇA que você medir. " +
	"Use só o que foi decidido, sem inventar restrições."
