# Parte III — Meso: Engenharia de Código e Física de Molas

## 3.1 A Linha de Prompt Essencial

No vídeo dissecado, Matt Loui revela que 80% do resultado depende de uma instrução precisa dada ao Claude Code:

#painel-sucesso("A LINHA DE PROMPT MESTRA", [
  *"Create a Remotion component matching this reference. Use the `spring()` function with mass: 0.5, damping: 12, stiffness: 100 for all transitions. Keep dimensions modular via props and ensure a transparent background."*
])

### 3.2 O Poder da Física Spring

Em vez de interpolações lineares duras, o Remotion calcula o movimento através de física real:

```typescript
import { spring, useCurrentFrame, useVideoConfig } from 'remotion';

export const AnimatedCard = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const scale = spring({
    frame,
    fps,
    config: {
      damping: 12,
      mass: 0.5,
      stiffness: 100,
    },
  });

  return (
    <div style={{ transform: `scale(${scale})` }} className="card-glass">
      <h2>Interface Real</h2>
    </div>
  );
};
```
