import React from 'react';
import { interpolate, useCurrentFrame } from 'remotion';

export const LegoloopDiagramDesignAnimation: React.FC = () => {
  const frame = useCurrentFrame();
  const opacity = interpolate(frame, [0, 30], [0, 1], { extrapolateRight: 'clamp' });

  return (
    <div style={{ width: '100%', height: '100%', background: '#f5f5f5', opacity, padding: 32 }}>
      <h1 style={{ fontFamily: 'Geist, sans-serif' }}>legoloop-diagram-design</h1>
    </div>
  );
};
