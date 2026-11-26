package render

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

// Sintese propria, deterministica e sem material de terceiros: ruido rosa com semente fixa,
// senoides e um acorde de 4 notas com ciclos inteiros em 20 s para o loop nao ter emenda.
func sinteseSFX(tipo string) []string {
	switch tipo {
	case "troca_cena":
		return []string{"-f", "lavfi", "-i", "anoisesrc=d=0.7:c=pink:r=48000:a=0.4:seed=7",
			"-af", "highpass=f=250,lowpass=f=5000,afade=t=in:d=0.2,afade=t=out:st=0.25:d=0.45,aformat=channel_layouts=stereo"}
	case "movimento_objeto":
		return []string{"-f", "lavfi", "-i", "sine=f=660:d=0.2:r=48000",
			"-af", "afade=t=out:st=0.02:d=0.18,volume=0.8,aformat=channel_layouts=stereo"}
	case "destaque":
		return []string{"-f", "lavfi", "-i", "sine=f=1320:d=0.6:r=48000",
			"-af", "afade=t=out:st=0.05:d=0.55,volume=0.7,aformat=channel_layouts=stereo"}
	}
	return nil
}

func sinteseBGM() []string {
	return []string{
		"-f", "lavfi", "-i", "sine=f=220:d=20:r=48000", "-f", "lavfi", "-i", "sine=f=277.2:d=20:r=48000",
		"-f", "lavfi", "-i", "sine=f=329.6:d=20:r=48000", "-f", "lavfi", "-i", "sine=f=440:d=20:r=48000",
		"-filter_complex", "[0][1][2][3]amix=inputs=4:normalize=0,tremolo=f=0.2:d=0.3,volume=0.25,aformat=channel_layouts=stereo[a]",
		"-map", "[a]",
	}
}

func ordemSFX() []string { return []string{"troca_cena", "movimento_objeto", "destaque"} }

// GarantirAtivosSonoros sintetiza BGM e SFX padrao que faltam; nunca sobrescreve arquivos existentes
// e nao sintetiza o que render.json ja aponta para um arquivo do projeto.
func GarantirAtivosSonoros(ctx context.Context, projDir string, cfg ConfigRender, ex Executor) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if ex == nil {
		ex = execSistema{}
	}
	if err := garantirBGM(ctx, projDir, cfg, ex); err != nil {
		return err
	}
	if cfg.SFXCatalogo != "" {
		return nil
	}
	for _, tipo := range ordemSFX() {
		if err := sintetizarSeFalta(ctx, projDir, CatalogoPadrao()[tipo], sinteseSFX(tipo), ex); err != nil {
			return err
		}
	}
	return nil
}

func garantirBGM(ctx context.Context, projDir string, cfg ConfigRender, ex Executor) error {
	if cfg.BGM == "" {
		return sintetizarSeFalta(ctx, projDir, "audio/bgm.wav", sinteseBGM(), ex)
	}
	caminho, err := resolverSeguro(projDir, cfg.BGM)
	if err != nil {
		return bloqueio("bgm", err)
	}
	if fi, err := os.Stat(caminho); err != nil || fi.IsDir() {
		return bloqueio("bgm", os.ErrNotExist)
	}
	return nil
}

func sintetizarSeFalta(ctx context.Context, projDir, rel string, args []string, ex Executor) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	caminho, err := resolverSeguro(projDir, rel)
	if err != nil {
		return err
	}
	if _, err := os.Stat(caminho); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	tmp := caminho + ".tmp.wav"
	completo := append(append([]string{"-y"}, args...), "-c:a", "pcm_s16le", tmp)
	if err := ex.Executar(ctx, "ffmpeg", completo...); err != nil {
		return errors.Join(err, removerTemp(tmp))
	}
	return os.Rename(tmp, caminho)
}
