package inventario

import "fmt"

// ValidarGenero aceita masculino, feminino e neutro (AC-E07.9).
func ValidarGenero(g string) error {
	switch g {
	case "masculino", "feminino", "neutro":
		return nil
	}
	return fmt.Errorf("%w: %q", ErrGeneroInvalido, g)
}

// ValidarFaixa aceita bebe, crianca, adolescente, adulto e idoso (AC-E07.9).
func ValidarFaixa(f string) error {
	switch f {
	case "bebe", "crianca", "adolescente", "adulto", "idoso":
		return nil
	}
	return fmt.Errorf("%w: %q", ErrFaixaInvalida, f)
}

// ValidarPersonagem exige gênero e faixa etária válidos em item do tipo personagem.
func ValidarPersonagem(it Item) error {
	if it.Tipo != "personagem" {
		return nil
	}
	if err := ValidarGenero(it.Genero); err != nil {
		return fmt.Errorf("personagem %q: %w", it.ID, err)
	}
	if err := ValidarFaixa(it.FaixaEtaria); err != nil {
		return fmt.Errorf("personagem %q: %w", it.ID, err)
	}
	return nil
}

// filtroValido aceita filtro vazio como qualquer valor e valida o resto pelo vocabulário.
func filtroValido(genero, faixa string) error {
	if genero != "" {
		if err := ValidarGenero(genero); err != nil {
			return err
		}
	}
	if faixa != "" {
		return ValidarFaixa(faixa)
	}
	return nil
}

// BuscarPersonagens reaproveita Buscar com tipo personagem e filtra gênero e faixa; vazio vale qualquer valor.
func (s Store) BuscarPersonagens(estilo, genero, faixa string) ([]Item, error) {
	if err := filtroValido(genero, faixa); err != nil {
		return nil, err
	}
	itens, err := s.Buscar(estilo, "personagem", nil)
	if err != nil {
		return nil, err
	}
	out := []Item{}
	for _, it := range itens {
		if (genero == "" || it.Genero == genero) && (faixa == "" || it.FaixaEtaria == faixa) {
			out = append(out, it)
		}
	}
	return out, nil
}
