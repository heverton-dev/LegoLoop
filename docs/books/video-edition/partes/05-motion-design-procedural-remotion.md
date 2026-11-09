## 5. Motion Design Procedural com Remotion

## 5.1 O Fim dos Templates Monolíticos do After Effects

A edição comercial tradicional depende fortemente de pacotes de motion graphics pré-renderizados comprados em bancos de mídia (Envato Elements, Motion Array). Esses templates apresentam limitações severas:

* Exigem máquinas potentes para rodar instâncias pesadas do Adobe After Effects;
* Qualquer alteração de uma única linha de texto obriga o editor a abrir um arquivo `.aep` complexo e renderizar novamente em tempo real;
* Dificuldade crônica de versionamento: arquivos binários de projeto não aceitam controle de versão via Git (`git diff`).

---

## 5.2 Animação como Componente React

No LegoLoop, os elementos visuais são codificados como componentes declarativos em **Remotion (React/TypeScript)** ou **HyperFrames**. Cada elemento gráfico (cartões laterais, citações em destaque, números de métricas, infográficos de código) é uma função pura que recebe parâmetros de entrada e renderiza quadros com canal alfa transparente (RGBA):

```tsx
import { interpolate, spring, useCurrentFrame, useVideoConfig } from 'remotion';

export const CartaoDidatico: React.FC<{ titulo: string; itens: string[] }> = ({ titulo, itens }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  // Curva de entrada com frenagem suave (Spring Physics)
  const entrada = spring({
    frame,
    fps,
    config: { damping: 14, mass: 0.8, stiffness: 120 }
  });

  const translateX = interpolate(entrada, [0, 1], [-400, 0]);

  return (
    <div style={{ transform: `translateX(${translateX}px)`, opacity: entrada }} className="cartao-glass-3d">
      <h2 className="titulo-anton">{titulo}</h2>
      <ul>
        {itens.map((it, idx) => (
          <li key={idx}>{it}</li>
        ))}
      </ul>
    </div>
  );
};
```

---

## 5.3 O Contrato de Storyboard Gráfico (`motion_storyboard.json`)

O nó `n05_motion_mapa` analisa a narrativa do vídeo e identifica momentos de alto impacto pedagógico onde um gráfico deve reforçar o aprendizado. Ele gera o manifesto do storyboard gráfico:

```json
{
  "cards": [
    {
      "card_id": "card_003",
      "layout": "side_glass_card",
      "start_ms": 28400,
      "end_ms": 34200,
      "title": "3 PILARES DA SOBERANIA",
      "items": [
        "1. Modelos Locais (Whisper & YuNet)",
        "2. Custo Financeiro R$ 0,00",
        "3. Timeline Universal FCPXML"
      ],
      "icon": "shield-check",
      "target_screen_side": "right"
    }
  ]
}
```

#painel-info("VANTAGENS DO RENDER PROCEDURAL", [
  1. **Automação Total:** O renderizador compila sequências de imagens PNG transparentes ou vídeo ProRes 4444 diretamente pela linha de comando em segundos;
  2. **Consistência de Marca:** Tipografia, pesos visuais, paletas de cores e raios de borda são governados por tokens de design reutilizáveis em CSS/TypeScript;
  3. **Zero Licenças:** O motor roda em qualquer sistema operacional (Linux, macOS, Windows) sem depender de licenças Adobe.
])
