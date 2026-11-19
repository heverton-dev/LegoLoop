#!/usr/bin/env python3
"""Geração Determinística de Animação Remotion (Nó e04).

Emite composition.tsx e composition.css em docs/diagrams/<slug>/animation/
seguindo o Design System LegoPress (Geist Sans/Mono, fundo descansado #f5f5f5).
"""

import sys
import json
from pathlib import Path

def resolver_slug(alvo: str):
    p = Path(alvo)
    if p.is_file():
        if p.name in ("pipeline.dag.json", "manifest.json", "quadrantes.json"):
            return p.parent.name
        if p.parent.name == "src" and p.parent.parent.parent.name == "diagrams":
            return p.parent.parent.name
        return p.stem
    return p.name

def main():
    if len(sys.argv) < 2:
        alvo = "legoloop-diagram-design"
    else:
        alvo = sys.argv[1]

    slug = resolver_slug(alvo)
    anim_dir = Path("docs/diagrams") / slug / "animation"
    anim_dir.mkdir(parents=True, exist_ok=True)

    tsx_file = anim_dir / "composition.tsx"
    css_file = anim_dir / "composition.css"

    pascal_name = "".join(word.capitalize() for word in slug.replace("-", "_").split("_")) + "Animation"

    tsx_content = f"""import React from 'react';
import {{ interpolate, useCurrentFrame }} from 'remotion';
import './composition.css';

export const {pascal_name}: React.FC = () => {{
  const frame = useCurrentFrame();
  const opacity = interpolate(frame, [0, 30], [0, 1], {{ extrapolateRight: 'clamp' }});
  const translateY = interpolate(frame, [0, 30], [20, 0], {{ extrapolateRight: 'clamp' }});

  return (
    <div className="diagram-animation-container" style={{{{ opacity, transform: `translateY(${{translateY}}px)` }}}}>
      <header className="diagram-animation-header">
        <span className="diagram-animation-badge">PROCESSO CANÔNICO</span>
        <h1 className="diagram-animation-title">{slug}</h1>
      </header>
      <div className="diagram-animation-stage">
        <p className="diagram-animation-caption">Visualização procedural e determinística LegoLoop.</p>
      </div>
    </div>
  );
}};
"""

    css_content = """/* Design System LegoPress - Remotion Stylesheet */
.diagram-animation-container {
  width: 100%;
  height: 100%;
  background-color: #f5f5f5;
  color: #1a202c;
  font-family: 'Geist', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  display: flex;
  flex-direction: column;
  padding: 48px;
  box-sizing: border-box;
}

.diagram-animation-header {
  margin-bottom: 24px;
}

.diagram-animation-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  background-color: #0f172a;
  color: #f8fafc;
  padding: 4px 8px;
  border-radius: 4px;
  margin-bottom: 12px;
}

.diagram-animation-title {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.03em;
  margin: 0;
  color: #0f172a;
}

.diagram-animation-stage {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed rgba(79, 93, 117, 0.2);
  border-radius: 8px;
  background-color: #ffffff;
}

.diagram-animation-caption {
  font-size: 14px;
  color: #4f5d75;
  font-family: 'Geist Mono', monospace;
}
"""

    tsx_file.write_text(tsx_content, encoding="utf-8")
    css_file.write_text(css_content, encoding="utf-8")
    print(f"[OK] e04 animacao gerada com sucesso: {tsx_file}")

if __name__ == "__main__":
    main()
