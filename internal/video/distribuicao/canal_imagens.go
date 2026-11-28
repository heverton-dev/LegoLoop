package distribuicao

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const scriptCanalImagensPy = `#!/usr/bin/env python3
import json, os, sys
from PIL import Image, ImageDraw, ImageFilter, ImageFont

FONTES = [
    'C:/Windows/Fonts/arialbd.ttf',
    '/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf',
    '/Library/Fonts/Arial Bold.ttf',
]

def fonte(tam):
    for f in FONTES:
        if os.path.exists(f):
            return ImageFont.truetype(f, tam)
    return ImageFont.load_default()

def cobrir(im, w, h):
    r = max(w / im.width, h / im.height)
    im = im.resize((int(im.width * r) + 1, int(im.height * r) + 1), Image.LANCZOS)
    x, y = (im.width - w) // 2, (im.height - h) // 2
    return im.crop((x, y, x + w, y + h))

def achar_base(pasta, prefixo):
    for ext in ('.png', '.jpg', '.jpeg'):
        p = os.path.join(pasta, 'base', prefixo + ext)
        if os.path.exists(p):
            return p
    return ''

def main():
    pasta = os.path.abspath(sys.argv[1])
    T = json.load(open(os.path.join(pasta, 'textos.json'), encoding='utf-8'))
    pal = T['identidade_visual']['paleta']
    ouro = pal.get('ouro', '#FFD24A')
    texto = pal.get('texto', '#FFFFFF')
    fundo = pal.get('fundo', '#0B0E1C')
    saida = os.path.join(pasta, 'imagens')
    os.makedirs(saida, exist_ok=True)
    pp = achar_base(pasta, 'perfil')
    if pp:
        p = cobrir(Image.open(pp).convert('RGB'), 1080, 1080)
        p.save(os.path.join(saida, 'perfil_1080.png'))
        p.resize((800, 800), Image.LANCZOS).save(os.path.join(saida, 'youtube_perfil_800.png'))
        p.resize((150, 150), Image.LANCZOS).save(os.path.join(saida, 'youtube_marca_dagua_150.png'), optimize=True)
        pr = p.resize((98, 98), Image.LANCZOS)
        m = Image.new('L', (98, 98), 0)
        ImageDraw.Draw(m).ellipse((0, 0, 97, 97), fill=255)
        prev = Image.new('RGB', (98, 98), fundo)
        prev.paste(pr, (0, 0), m)
        prev.resize((294, 294), Image.NEAREST).save(os.path.join(saida, 'previa_perfil_98px_ampliada.png'))
    bp = achar_base(pasta, 'banner')
    if bp:
        W, H, SW, SH = 2560, 1440, 1546, 423
        b = cobrir(Image.open(bp).convert('RGB'), W, H)
        cx, cy = W // 2, H // 2
        faixa = Image.new('RGBA', (W, H), (0, 0, 0, 0))
        ImageDraw.Draw(faixa).rounded_rectangle((cx - SW // 2 - 40, cy - SH // 2 - 20, cx + SW // 2 + 40, cy + SH // 2 + 20), 60, fill=(11, 14, 28, 165))
        b = Image.alpha_composite(b.convert('RGBA'), faixa.filter(ImageFilter.GaussianBlur(36))).convert('RGB')
        d = ImageDraw.Draw(b)
        d.text((cx, int(cy - SH * 0.14)), T['nome'].upper(), font=fonte(int(SH * 0.42)), fill=ouro, anchor='mm', stroke_width=9, stroke_fill='#000000')
        d.text((cx, int(cy + SH * 0.28)), T['tagline'].upper(), font=fonte(int(SH * 0.14)), fill=texto, anchor='mm', stroke_width=5, stroke_fill='#000000')
        b.save(os.path.join(saida, 'youtube_banner_2560x1440.png'))
        pv = b.copy()
        ImageDraw.Draw(pv).rectangle((cx - SW // 2, cy - SH // 2, cx + SW // 2, cy + SH // 2), outline='#FF5A5A', width=6)
        pv.resize((1280, 720)).save(os.path.join(saida, 'previa_banner_area_segura.png'))
        b.crop((cx - SW // 2, cy - SH // 2, cx + SW // 2, cy + SH // 2)).save(os.path.join(saida, 'previa_banner_so_area_segura.png'))

if __name__ == '__main__':
    main()
`

// ResultadoImagensCanal resume os arquivos PNG exportados para as redes e se o logo foi sincronizado.
type ResultadoImagensCanal struct {
	Arquivos         []string `json:"arquivos"`
	LogoSincronizado string   `json:"logo_sincronizado,omitempty"`
}

// GerarImagensCanal gera as imagens oficiais e previas do canal a partir de base/perfil e base/banner.
func GerarImagensCanal(ctx context.Context, ex Executor, canalDir, destinoLogoMarca string) (res ResultadoImagensCanal, err error) {
	if err := validarBasesCanal(canalDir); err != nil {
		return res, err
	}
	script, limpar, err := gravarScriptTemp("canal-img-*", "canal_imagens.py", []byte(scriptCanalImagensPy))
	if err != nil {
		return res, err
	}
	defer func() { err = errors.Join(err, limpar()) }()
	if err := ex.Executar(ctx, "python", script, canalDir); err != nil {
		return res, err
	}
	arqs, err := listarPNGsGerados(filepath.Join(canalDir, "imagens"))
	if err != nil {
		return res, err
	}
	res.Arquivos = arqs
	if destinoLogoMarca != "" {
		res.LogoSincronizado, err = copiarMarcaDaguaParaLogo(canalDir, destinoLogoMarca)
	}
	return res, err
}

func validarBasesCanal(canalDir string) error {
	if _, err := os.Stat(filepath.Join(canalDir, "textos.json")); err != nil {
		return fmt.Errorf("%w: textos.json ausente em %s", ErrPedidoInvalido, canalDir)
	}
	for _, base := range []string{"perfil.png", "perfil.jpg", "banner.png", "banner.jpg"} {
		if _, err := os.Stat(filepath.Join(canalDir, "base", base)); err == nil {
			return nil
		}
	}
	return fmt.Errorf("%w: base/perfil ou base/banner ausente em %s", ErrPedidoInvalido, canalDir)
}

func gravarScriptTemp(padrao, nome string, conteudo []byte) (string, func() error, error) {
	dir, err := os.MkdirTemp("", padrao)
	if err != nil {
		return "", nil, err
	}
	caminho := filepath.Join(dir, nome)
	if err := os.WriteFile(caminho, conteudo, 0o600); err != nil {
		return "", nil, errors.Join(err, os.RemoveAll(dir))
	}
	return caminho, func() error { return os.RemoveAll(dir) }, nil
}

func listarPNGsGerados(dir string) ([]string, error) {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: pasta imagens ausente: %v", ErrPedidoInvalido, err)
	}
	var out []string
	for _, e := range entradas {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: nenhuma imagem gerada", ErrPedidoInvalido)
	}
	return out, nil
}

func copiarMarcaDaguaParaLogo(canalDir, destinoLogo string) (string, error) {
	origem := filepath.Join(canalDir, "imagens", "youtube_marca_dagua_150.png")
	dados, err := os.ReadFile(origem)
	if err != nil {
		return "", err
	}
	if err := gravarAtomico(destinoLogo, dados); err != nil {
		return "", err
	}
	return destinoLogo, nil
}
