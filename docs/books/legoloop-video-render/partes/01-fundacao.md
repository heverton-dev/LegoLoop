# Fundacao

O Processo `legoloop-video-render` transforma o roteiro e os ativos em video. Ele junta quatro garantias: a voz bate com o texto (Whisper), o ritmo bate com o alvo (gate de ritmo), o audio bate com o padrao (LUFS medido) e o video bate com o contrato (1080x1920, 30 fps, com voz).

O desenho dos bonecos fica no motor portado, em `motor-local/`, sob licenca MIT. O Go cuida de contratos, voz, mix, gates e orquestracao. O Python so desenha os quadros no Chromium e o ffmpeg monta o MP4.
