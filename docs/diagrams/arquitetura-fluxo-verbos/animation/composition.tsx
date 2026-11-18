import React from 'react';
import { interpolate, spring, useCurrentFrame, useVideoConfig } from 'remotion';

// Composition Remotion para Micro-Animações dos Diagramas LegoLoop
// Segue o Design System Editorial do legoloop-diagram-design

export const DiagramAnimation: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  // Transições em camadas escalonadas
  const layer1Opacity = interpolate(frame, [0, 15], [0, 1], { extrapolateRight: 'clamp' });
  const layer2Opacity = interpolate(frame, [15, 30], [0, 1], { extrapolateRight: 'clamp' });
  const layer3Opacity = interpolate(frame, [30, 45], [0, 1], { extrapolateRight: 'clamp' });

  // Micro-interação de pulso suave no motor 'fabricar'
  const fabricarScale = spring({
    frame: frame - 20,
    fps,
    config: { damping: 12, mass: 0.8, stiffness: 100 },
  });

  return (
    <div
      style={{
        width: '100%',
        height: '100%',
        backgroundColor: '#f5f5f5',
        display: 'flex',
        flexDirection: 'column',
        padding: '32px',
        fontFamily: 'Geist, sans-serif',
      }}
    >
      <div style={{ opacity: layer1Opacity, transition: 'opacity 0.2s ease-out' }}>
        {/* Camada 01: Fundação */}
      </div>
      <div
        style={{
          opacity: layer2Opacity,
          transform: `scale(${Math.max(1, fabricarScale)})`,
          transition: 'opacity 0.2s ease-out',
        }}
      >
        {/* Camada 02: Espinha Dorsal */}
      </div>
      <div style={{ opacity: layer3Opacity, transition: 'opacity 0.2s ease-out' }}>
        {/* Camada 03: Execução e Qualidade */}
      </div>
    </div>
  );
};
