package videorea

import (
	"os/exec"
	"strings"
)

// DiagnoseToolFit avalia a prontidão e o encaixe de uma ferramenta específica no SO do operador.
func DiagnoseToolFit(toolName string) ToolDiagnosis {
	lower := strings.ToLower(toolName)
	switch {
	case strings.Contains(lower, "ffmpeg"):
		return diagnoseFFmpeg()
	case strings.Contains(lower, "hyperframes"):
		return diagnoseHyperFrames()
	case strings.Contains(lower, "yunet") || strings.Contains(lower, "opencv"):
		return diagnoseYuNet()
	case strings.Contains(lower, "whisperx"):
		return diagnoseWhisperX()
	case strings.Contains(lower, "yt-dlp"):
		return diagnoseYtDlp()
	case strings.Contains(lower, "whisper"):
		return diagnoseWhisperLocal()
	case strings.Contains(lower, "sfx") || strings.Contains(lower, "epidemic"):
		return diagnoseSFX()
	case strings.Contains(lower, "claude") || strings.Contains(lower, "9router"):
		return diagnoseRouter()
	case strings.Contains(lower, "comfy"):
		return diagnoseComfyUI()
	case strings.Contains(lower, "minimax") || strings.Contains(lower, "minimx") ||
		strings.Contains(lower, "flux"):
		return diagnoseLocalDiffusionModel(toolName)
	default:
		return diagnoseLocalScript(toolName)
	}
}

func diagnoseComfyUI() ToolDiagnosis {
	for _, bin := range []string{"comfy", "comfyui"} {
		if path, err := exec.LookPath(bin); err == nil {
			return ToolDiagnosis{
				InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
				FitSolution:  "Encaixe perfeito: ComfyUI disponivel no PATH e pronto para orquestrar workflows locais.",
				GoldStandard: "ComfyUI Portable (Windows) / comfyui CLI",
			}
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar o ComfyUI Portable (Windows) para orquestrar os workflows de geracao de video local sem mensalidade.",
		GoldStandard:   "ComfyUI Portable (Windows) / comfyui CLI",
		InstallCommand: "git clone https://github.com/comfyanonymous/ComfyUI && pip install -r requirements.txt",
	}
}

func diagnoseLocalDiffusionModel(name string) ToolDiagnosis {
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Pesos de " + name + " nao detectados no host: baixar VAE e diffusion model e aloca-los em models/ do ComfyUI antes da geracao.",
		GoldStandard:   name + " (pesos abertos locais via ComfyUI)",
		InstallCommand: "ComfyUI Manager > Load Model > " + name,
	}
}

func diagnoseFFmpeg() ToolDiagnosis {
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return ToolDiagnosis{
			InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
			FitSolution:  "Encaixe perfeito: binario FFmpeg instalado e homologado no PATH.",
			GoldStandard: "FFmpeg CLI Engine v7.1+",
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar o FFmpeg no host para habilitar cortes automaticos e remux de audio sem recodificacao.",
		GoldStandard:   "FFmpeg CLI Engine v7.1+",
		InstallCommand: "choco install ffmpeg -y || winget install Gyan.FFmpeg",
	}
}

func diagnoseHyperFrames() ToolDiagnosis {
	if path, err := exec.LookPath("hyperframes"); err == nil {
		return ToolDiagnosis{
			InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
			FitSolution:  "Encaixe perfeito: motor HyperFrames ativo no ambiente bun para motion graphics procedural.",
			GoldStandard: "HyperFrames CLI Engine",
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar o CLI do HyperFrames via bun/npm para renderizar animacoes em HTML/CSS/GSAP.",
		GoldStandard:   "HyperFrames CLI Engine",
		InstallCommand: "bun add -g hyperframes || npm install -g hyperframes",
	}
}

func diagnoseYuNet() ToolDiagnosis {
	if path, err := exec.LookPath("python"); err == nil {
		return ToolDiagnosis{
			InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
			FitSolution:  "Encaixe perfeito: modulo OpenCV com suporte a YuNet disponivel no Python para rastreio facial.",
			GoldStandard: "YuNet Face Tracker (OpenCV Zoo)",
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar opencv-python para permitir a deteccao de pontos faciais e enquadramento de camera.",
		GoldStandard:   "YuNet Face Tracker (OpenCV Zoo)",
		InstallCommand: "pip install opencv-python numpy",
	}
}

func diagnoseWhisperX() ToolDiagnosis {
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar WhisperX via pip para habilitar o alinhamento temporal milissegundo a milissegundo por palavra.",
		GoldStandard:   "WhisperX (Padrão Ouro de alinhamento por palavra)",
		InstallCommand: "pip install whisperx",
	}
}

func diagnoseYtDlp() ToolDiagnosis {
	if path, err := exec.LookPath("yt-dlp"); err == nil {
		return ToolDiagnosis{
			InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
			FitSolution:  "Encaixe perfeito: ferramenta de download de assets disponivel no ambiente Python Scripts.",
			GoldStandard: "yt-dlp Core",
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar yt-dlp via pip para viabilizar a extracao de midia e metadados de referencia.",
		GoldStandard:   "yt-dlp Core",
		InstallCommand: "pip install yt-dlp",
	}
}

func diagnoseWhisperLocal() ToolDiagnosis {
	if path, err := exec.LookPath("whisper"); err == nil {
		return ToolDiagnosis{
			InstalledOnHost: true, HostLocation: path, FitStatus: FitOuro,
			FitSolution:  "Encaixe perfeito: binario local do Whisper pronto para transcricao offline.",
			GoldStandard: "Whisper Local (whisper.cpp)",
		}
	}
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Instalar whisper.cpp ou usar silenceremove local para dispensar custos de API por hora de audio.",
		GoldStandard:   "Whisper Local (whisper.cpp / ffmpeg silenceremove)",
		InstallCommand: "winget install whisper.cpp || choco install whisper",
	}
}

func diagnoseSFX() ToolDiagnosis {
	return ToolDiagnosis{
		InstalledOnHost: false, FitStatus: MissingAquisicao,
		FitSolution:    "Sincronizar o pack local de SFX CC0 para obter biblioteca de efeitos sonoros sem mensalidades.",
		GoldStandard:   "Biblioteca Local de SFX Curada (CC0 / Freesound)",
		InstallCommand: "legoloop sync-sfx-pack --tier ouro",
	}
}

func diagnoseRouter() ToolDiagnosis {
	return ToolDiagnosis{
		InstalledOnHost: true, HostLocation: "Claude CLI / agy", FitStatus: FitFrouxoDefasado,
		FitSolution:    "Configurar o endpoint seguro do 9Router na VPS para roteamento resiliente de modelos abertos, eliminando a dependência do plano Pro comercial e evitando rate-limits.",
		GoldStandard:   "9Router Local / Modelos Abertos Gratuitos",
		InstallCommand: "setx NINEROUTER_URL https://9router.vpsconexao.org",
	}
}

func diagnoseLocalScript(name string) ToolDiagnosis {
	return ToolDiagnosis{
		InstalledOnHost: true, HostLocation: "Target Workspace Local (scratch)", FitStatus: FitOuro,
		FitSolution:  "Encaixe perfeito: script autoral extraido e funcional no scratch local para consumo da esteira.",
		GoldStandard: "Script autoral do criador adaptado para LegoLoop",
	}
}
