# Capítulo 2: A Arquitetura Visual Three.js & Áudio Gated RMS

A estética distintiva do Cartoon Shorts baseia-se em uma fusão inovadora entre geometria tridimensional e ilustração clássica de quadrinhos.

## 1. O Shader de Tinta (Ink Contract) e Normal/Depth Pass
Ao contrário de animações 3D convencionais com iluminação realista, o motor utiliza um shader customizado em Three.js r160:
* **Normal/Depth Pass:** Em cada quadro, a cena é renderizada primeiro em um buffer que captura apenas as normais das superfícies e a profundidade relativa da câmera.
* **Detecção de Bordas (Sobel/Edge Detection):** Um shader pós-processamento traça linhas de contorno pretas com raio de 2.5px onde quer que haja um salto brusco de profundidade ou uma quina angular.
* **Texturas Pintadas (`paint.js`):** Pranchas de madeira, pedras, grama e azulejos são gerados proceduralmente em Canvas com costuras escuras e baixo contraste interno, garantindo que o traço externo se destaque.

## 2. A Lei do Áudio: Trava Gated RMS a -6 dB
Em vídeos curtos, o maior erro amador é colocar a trilha musical brigando com a voz do narrador.

O pipeline LegoLoop Cartoon Shorts impõe uma **Lei Física de Mixagem**:
* O nível de áudio da voz neural é medido em RMS fechado (*Gated RMS*).
* A trilha musical sintetizada é atenuada de forma precisa para ficar exatamente a **50% do volume da voz (-6.0 dB)**.
* Os efeitos sonoros (SFX) são limitados a picos de no máximo 0.7x o pico da fala.
* O arquivo final é processado pelo filtro `loudnorm` do FFmpeg, cravando **-14 LUFS** de sonoridade e limitador de pico (*true peak*) em **-1 dBFS**.
