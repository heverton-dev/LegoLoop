# Parte IV — Micro: Montagem Determinística e Código

## 4.1 Composição via Remotion e FFmpeg

A montagem final é orquestrada através de um script de timeline programática:

```typescript
import { Composition, Sequence } from 'remotion';
import { SceneCard } from './components/SceneCard';

export const MotionFlowComposition = () => {
  return (
    <Composition
      id="MasterMotion"
      component={SceneSequence}
      durationInFrames={1800} // 60s @ 30fps
      fps={30}
      width={1920}
      height={1080}
    />
  );
};
```

O comando de compilação headless do pipeline executa:
```bash
npx remotion render src/index.ts MasterMotion dist/final_render.mp4 --codec h264 --crf 18
```
