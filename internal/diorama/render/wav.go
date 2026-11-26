package render

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var ErrAudioInvalido = errors.New("diorama-render: audio WAV invalido")

// formatoWAV descreve PCM de 16 bits, o unico formato aceito entre provedores.
type formatoWAV struct {
	Canais uint16
	Taxa   uint32
	Bits   uint16
}

func audioInvalido(motivo string) error {
	return fmt.Errorf("%w: %s", ErrAudioInvalido, motivo)
}

func lerFmt(corpo []byte) (formatoWAV, error) {
	if len(corpo) < 16 || binary.LittleEndian.Uint16(corpo[0:2]) != 1 {
		return formatoWAV{}, audioInvalido("fmt ausente ou nao PCM")
	}
	f := formatoWAV{
		Canais: binary.LittleEndian.Uint16(corpo[2:4]),
		Taxa:   binary.LittleEndian.Uint32(corpo[4:8]),
		Bits:   binary.LittleEndian.Uint16(corpo[14:16]),
	}
	if f.Canais == 0 || f.Taxa == 0 || f.Bits != 16 {
		return f, audioInvalido("exige PCM de 16 bits com canais e taxa validos")
	}
	return f, nil
}

// lerWAV extrai formato e amostras, ignorando chunks extras como LIST.
func lerWAV(b []byte) (formatoWAV, []byte, error) {
	var f formatoWAV
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return f, nil, audioInvalido("cabecalho RIFF/WAVE ausente")
	}
	temFmt := false
	for pos := 12; pos+8 <= len(b); {
		id, tam := string(b[pos:pos+4]), int(binary.LittleEndian.Uint32(b[pos+4:pos+8]))
		corpo, fim := pos+8, pos+8+tam
		if fim > len(b) {
			if id != "data" {
				return f, nil, audioInvalido("chunk " + id + " truncado")
			}
			fim = len(b)
		}
		switch id {
		case "fmt ":
			var err error
			if f, err = lerFmt(b[corpo:fim]); err != nil {
				return f, nil, err
			}
			temFmt = true
		case "data":
			if !temFmt {
				return f, nil, audioInvalido("data antes de fmt")
			}
			return f, b[corpo:fim], nil
		}
		pos = fim + tam&1
	}
	return f, nil, audioInvalido("chunk data ausente")
}

// montarWAV grava o cabecalho canonico de 44 bytes seguido das amostras.
func montarWAV(f formatoWAV, dados []byte) []byte {
	bloco := f.Canais * f.Bits / 8
	b := make([]byte, 44, 44+len(dados))
	copy(b[0:], "RIFF")
	binary.LittleEndian.PutUint32(b[4:], uint32(36+len(dados)))
	copy(b[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(b[16:], 16)
	binary.LittleEndian.PutUint16(b[20:], 1)
	binary.LittleEndian.PutUint16(b[22:], f.Canais)
	binary.LittleEndian.PutUint32(b[24:], f.Taxa)
	binary.LittleEndian.PutUint32(b[28:], f.Taxa*uint32(bloco))
	binary.LittleEndian.PutUint16(b[32:], bloco)
	binary.LittleEndian.PutUint16(b[34:], f.Bits)
	copy(b[36:], "data")
	binary.LittleEndian.PutUint32(b[40:], uint32(len(dados)))
	return append(b, dados...)
}

// juntarWAV concatena trechos do mesmo formato em um unico WAV.
func juntarWAV(partes [][]byte) ([]byte, error) {
	if len(partes) == 0 {
		return nil, audioInvalido("nenhum trecho para juntar")
	}
	var base formatoWAV
	var todos []byte
	for i, p := range partes {
		f, dados, err := lerWAV(p)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			base = f
		} else if f != base {
			return nil, audioInvalido(fmt.Sprintf("trecho %d com formato diferente do primeiro", i))
		}
		todos = append(todos, dados...)
	}
	return montarWAV(base, todos), nil
}

// duracaoWAV converte bytes de amostras em segundos.
func duracaoWAV(f formatoWAV, nBytes int) float64 {
	bps := float64(f.Taxa) * float64(f.Canais) * float64(f.Bits/8)
	if bps == 0 {
		return 0
	}
	return float64(nBytes) / bps
}
