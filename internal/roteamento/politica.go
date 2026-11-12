package roteamento

const (
	PoliticaDegrada = "degrada"
	PoliticaFalha   = "falha"
)

func PoliticaDeFalha(classe string) string {
	switch classe {
	case ClasseNano, ClasseMicro:
		return PoliticaDegrada
	}
	return PoliticaFalha
}
