import React from 'react';
import { interpolate, useCurrentFrame } from 'remotion';
import './composition.css';

export const LegoloopDiagramDesignAnimation: React.FC = () => {
  const frame = useCurrentFrame();
  const opacity = interpolate(frame, [0, 30], [0, 1], { extrapolateRight: 'clamp' });
  const translateY = interpolate(frame, [0, 30], [20, 0], { extrapolateRight: 'clamp' });

  return (
    <div className="diagram-animation-container" style={{ opacity, transform: `translateY(${translateY}px)` }}>
      <header className="diagram-animation-header">
        <span className="diagram-animation-badge">PROCESSO CANÔNICO</span>
        <h1 className="diagram-animation-title">legoloop-diagram-design</h1>
      </header>
      <div className="diagram-animation-stage">
        <p className="diagram-animation-caption">Visualização procedural e determinística LegoLoop.</p>
      </div>
    </div>
  );
};
