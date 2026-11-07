# Parte III — Pilar 2: Whiteboard no Piloto Automático

## 3.1 A Ilusão das Ferramentas de Lousa Branca

Desde o surgimento de softwares como VideoScribe, Doodly e Vyond, vídeos no estilo *whiteboard explainer* (lousa branca com desenho progressivo de traços e personagens) consolidaram-se como uma das formas mais eficazes de retenção cognitiva na educação e no marketing digital.

Essas plataformas proprietárias cobram entre US$ 39 e US$ 49 mensais por um serviço estritamente mecânico: exibir a imagem de uma mão humana que percorre um caminho vetorial pré-gravado, enquanto uma máscara dinâmica revela o desenho por baixo.

Além do custo recorrente elevado, essas soluções sofrem de três limitações graves:
- **Resolução Limitada:** Exportações frequentemente travadas em 1080p ou com artefatos de compressão em gradientes.
- **Acervo Estático e Monótono:** Dependência de bibliotecas fechadas de cliparts que tornam os vídeos idênticos entre diferentes canais.
- **Incapacidade de Integração com Dados:** Impossibilidade de gerar um gráfico vetorial ou fórmula matemática animada a partir de uma API externa em tempo real.

O ecossistema `legoloop-video-master` substitui completamente essas suítes por um **Motor Vetorial Procedural Baseado em SVG Puro e Easing Matemático**.


![Diagrama de Arquitetura](/img/diagrama_e1a5eaea.png){width=85%}


## 3.2 A Matemática do Traçado SVG: `stroke-dasharray`

O princípio fundamental que torna o desenho progressivo infinitamente escalável em 4K nativo baseia-se na propriedade CSS/SVG `stroke-dasharray` e `stroke-dashoffset`.

Para qualquer elemento `<path>` vetorial:
1. Extrai-se o comprimento euclidiano total do caminho:
   $$L = \text{path.getTotalLength()}$$
2. Configura-se o padrão de tracejado com comprimento igual a $L$:
   $$\text{strokeDasharray} = L$$
3. No início da animação ($t=0$), desloca-se o traço em $L$, ocultando o caminho integralmente:
   $$\text{strokeDashoffset}(0) = L$$
4. À medida que o tempo avança de $t=0$ até $t=1$, reduz-se o deslocamento até zero, simulando a caneta desenhando a linha:
   $$\text{strokeDashoffset}(t) = L \cdot (1 - \text{easing}(t))$$

```tsx
// Componente de Traçado Vetorial Automático em Remotion
import React from "react";
import { useCurrentFrame, interpolate } from "remotion";

interface AnimatedPathProps {
  d: string;
  totalLength: number;
  durationInFrames: number;
  strokeColor?: string;
  strokeWidth?: number;
}

export const AnimatedPath: React.FC<AnimatedPathProps> = ({
  d,
  totalLength,
  durationInFrames,
  strokeColor = "#0f172a",
  strokeWidth = 3.5
}) => {
  const frame = useCurrentFrame();

  // Curva de easing humano suave (desenho que desacelera nas curvas complexas)
  const progress = interpolate(frame, [0, durationInFrames], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp"
  });

  const dashoffset = totalLength * (1 - progress);

  return (
    <path
      d={d}
      fill="none"
      stroke={strokeColor}
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeDasharray={totalLength}
      strokeDashoffset={dashoffset}
    />
  );
};
```

## 3.3 Avatares Stickman Paramétricos e Cinemática Direta

Para ilustrar conceitos sem desenhar personagens à mão, o pilar Whiteboard utiliza **Stickmen Paramétricos** modelados via ângulos articulares polares expressos em JSON:

$$\text{Ponto}_B = \text{Ponto}_A + \begin{bmatrix} r \cdot \cos(\theta) \\ r \cdot \sin(\theta) \end{bmatrix}$$

Cada pose emocional (curiosidade, espanto, triunfo ou explicação) é descrita por um conjunto de seis ângulos fundamentais:
- $\theta_1, \theta_2$: Ângulo dos ombros e cotovelos esquerdos.
- $\theta_3, \theta_4$: Ângulo dos ombros e cotovelos direitos.
- $\theta_5, \theta_6$: Inclinação do tronco e cabeça.

```json
{
  "pose": "eureka_triunfo",
  "braco_esquerdo": { "ombro_deg": 135, "cotovelo_deg": 45 },
  "braco_direito": { "ombro_deg": 45, "cotovelo_deg": 135 },
  "tronco_deg": 90,
  "cabeca_inclinacao": 5,
  "expressao": "olhos_abertos_sorriso"
}
```

Ao interpolar entre dicionários de poses no frame rate configurado, o sistema produz animações fluidas e expressivas sem exigir nenhum software gráfico de desktop, gerando arquivos de saída com peso inferior a **50 KB** antes da rasterização final.
