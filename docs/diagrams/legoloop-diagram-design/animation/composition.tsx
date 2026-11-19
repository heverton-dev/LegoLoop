```tsx
import React from "react";
import {
  AbsoluteFill,
  Composition,
  Easing,
  interpolate,
  useCurrentFrame,
  useVideoConfig,
} from "remotion";

const FPS = 30;
const WIDTH = 1040;
const HEIGHT = 640;
const FADE_MS = 200;
const STAGGER_MS = 50;
const SLIDE_PX = 24;

const COLORS = {
  paper: "#f5f5f5",
  paper2: "#ececec",
  card: "#ffffff",
  ink: "#111111",
  muted: "#4f5d75",
  soft: "#5a6580",
  ruleSolid: "#bfc0c0",
  skill: "#2563eb",
  script: "#0f172a",
  accent: "#bf4520",
};

const FONT_SANS = "'Geist', system-ui, sans-serif";
const FONT_MONO = "'Geist Mono', ui-monospace, monospace";

const msToFrames = (ms: number): number => Math.round((ms / 1000) * FPS);

type LayerSpec = {
  y: number;
  height: number;
  label: string;
  fill: string;
  stroke: string;
  text: string;
};

type NodeSpec = {
  id: string;
  badge: string;
  badgeWidth: number;
  name: string;
  sub?: string;
  x: number;
  y: number;
  width: number;
  height: number;
  stroke: string;
  strokeWidth: number;
  badgeFill: string;
  subFill: string;
};

type EdgeSpec = {
  d: string;
  stroke: string;
  from: number;
  marker: string;
};

const LAYERS: LayerSpec[] = [
  { y: 24, height: 84, label: "CAMADA 01 · SETUP & GOVERNANÇA DE ARQUIVOS", fill: COLORS.paper2, stroke: COLORS.ruleSolid, text: COLORS.soft },
  { y: 124, height: 108, label: "CAMADA 02 · INGESTÃO TOPOLÓGICA & DESIGN SYSTEM", fill: COLORS.paper2, stroke: COLORS.ruleSolid, text: COLORS.soft },
  { y: 248, height: 114, label: "CAMADA 03 · COMPILAÇÃO EDITORIAL VETORIAL (MOTOR SVG)", fill: "rgba(191,69,32,0.04)", stroke: COLORS.accent, text: COLORS.accent },
  { y: 378, height: 114, label: "CAMADA 04 · CAMADA DE MOVIMENTO & EXPORTAÇÃO", fill: COLORS.paper2, stroke: COLORS.ruleSolid, text: COLORS.soft },
  { y: 508, height: 108, label: "CAMADA 05 · FECHAMENTO CANÔNICO & LIVRO TYPST", fill: "rgba(37,99,235,0.04)", stroke: COLORS.skill, text: COLORS.skill },
];

const NODES: NodeSpec[] = [
  { id: "e01", badge: "e01 SKILL", badgeWidth: 56, name: "inicializar-slug-e-anticorpos", x: 380, y: 52, width: 280, height: 38, stroke: COLORS.ruleSolid, strokeWidth: 1.2, badgeFill: COLORS.skill, subFill: COLORS.soft },
  { id: "e02", badge: "e02 SKILL", badgeWidth: 56, name: "ingestao-e-topologia", sub: "Mapeamento de Peças & Larguras", x: 370, y: 160, width: 300, height: 48, stroke: COLORS.ruleSolid, strokeWidth: 1.2, badgeFill: COLORS.skill, subFill: COLORS.soft },
  { id: "e03", badge: "e03 FOCAL", badgeWidth: 56, name: "geracao-svg-html", sub: "Grid Ortogonal · Zero Transbordamento", x: 340, y: 280, width: 360, height: 52, stroke: COLORS.accent, strokeWidth: 2, badgeFill: COLORS.accent, subFill: COLORS.soft },
  { id: "e04", badge: "e04 SKILL", badgeWidth: 56, name: "micro-animacao", sub: "CSS Keyframes & Remotion Tsx", x: 180, y: 410, width: 280, height: 50, stroke: COLORS.ruleSolid, strokeWidth: 1.2, badgeFill: COLORS.skill, subFill: COLORS.soft },
  { id: "e05", badge: "e05 SCRIPT", badgeWidth: 62, name: "gate-e-exportacao", sub: "export_svg.py > Typst Standalone", x: 580, y: 410, width: 280, height: 50, stroke: COLORS.script, strokeWidth: 1.5, badgeFill: COLORS.script, subFill: COLORS.soft },
  { id: "e06", badge: "e06 SKILL", badgeWidth: 56, name: "fechamento-e-anticorpos", sub: "Relatório Factual & Imunologia", x: 370, y: 540, width: 300, height: 50, stroke: COLORS.skill, strokeWidth: 1.8, badgeFill: COLORS.skill, subFill: COLORS.skill },
];

const EDGES: EdgeSpec[] = [
  { d: "M520,88 L520,154", stroke: COLORS.muted, from: 1, marker: "arr-flow" },
  { d: "M520,208 L520,278", stroke: COLORS.accent, from: 2, marker: "arr-accent" },
  { d: "M420,332 L420,358 L320,358 L320,408", stroke: COLORS.muted, from: 3, marker: "arr-flow" },
  { d: "M620,332 L620,358 L720,358 L720,408", stroke: COLORS.script, from: 4, marker: "arr-flow" },
  { d: "M320,462 L320,488 L470,488 L470,538", stroke: COLORS.muted, from: 5, marker: "arr-flow" },
  { d: "M720,462 L720,488 L570,488 L570,538", stroke: COLORS.skill, from: 5, marker: "arr-typst" },
];

const fadeProgress = (frame: number, startFrame: number): number =>
  interpolate(frame, [startFrame, startFrame + msToFrames(FADE_MS)], [0, 1], {
    easing: Easing.out(Easing.cubic),
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

const nodeStart = (index: number): number =>
  msToFrames(FADE_MS) + msToFrames(STAGGER_MS) * index * 2;

const DiagramLayer: React.FC<{ spec: LayerSpec; index: number }> = ({ spec, index }) => {
  const frame = useCurrentFrame();
  const opacity = fadeProgress(frame, index * msToFrames(STAGGER_MS));
  return (
    <g opacity={opacity}>
      <rect x={24} y={spec.y} width={992} height={spec.height} rx={6} fill={spec.fill} stroke={spec.stroke} strokeWidth={1} />
      <text x={40} y={spec.y + 20} fontFamily={FONT_MONO} fontSize={9} fontWeight={600} letterSpacing="0.14em" fill={spec.text}>
        {spec.label}
      </text>
    </g>
  );
};

const DiagramNode: React.FC<{ spec: NodeSpec; index: number }> = ({ spec, index }) => {
  const frame = useCurrentFrame();
  const start = nodeStart(index);
  const progress = fadeProgress(frame, start);
  const offset = (1 - progress) * SLIDE_PX;
  const hasSub = Boolean(spec.sub);
  const badgeY = spec.id === "e03" ? 10 : 8;
  const badgeX = spec.id === "e03" ? 10 : 8;
  const textX = spec.id === "e03" ? 76 : spec.id === "e05" ? 80 : 74;
  const nameY = hasSub ? (spec.id === "e03" ? 23 : 22) : 23;
  const nameSize = spec.id === "e03" ? 13 : 12;
  const nameWeight = spec.id === "e01" ? 600 : 700;
  return (
    <g transform={`translate(${spec.x}, ${spec.y + offset})`} opacity={progress}>
      <rect x={0} y={0} width={spec.width} height={spec.height} rx={6} fill={COLORS.card} stroke={spec.stroke} strokeWidth={spec.strokeWidth} />
      <rect x={badgeX} y={badgeY} width={spec.badgeWidth} height={22} rx={3} fill={spec.badgeFill} />
      <text x={badgeX + spec.badgeWidth / 2} y={badgeY + 14} fontFamily={FONT_MONO} fontSize={9} fontWeight={700} fill="#ffffff" textAnchor="middle">
        {spec.badge}
      </text>
      <text x={textX} y={nameY} fontFamily={FONT_SANS} fontSize={nameSize} fontWeight={nameWeight} fill={COLORS.ink}>
        {spec.name}
      </text>
      {hasSub ? (
        <text x={textX} y={spec.id === "e03" ? 41 : 38} fontFamily={FONT_MONO} fontSize={10} fill={spec.subFill}>
          {spec.sub}
        </text>
      ) : null}
    </g>
  );
};

const DiagramEdge: React.FC<{ spec: EdgeSpec }> = ({ spec }) => {
  const frame = useCurrentFrame();
  const start = nodeStart(spec.from) - msToFrames(FADE_MS) / 2;
  const opacity = fadeProgress(frame, Math.max(start, 0));
  return (
    <path d={spec.d} fill="none" stroke={spec.stroke} strokeWidth={2} opacity={opacity} markerEnd={`url(#${spec.marker})`} />
  );
};

const Marker: React.FC<{ id: string; fill: string }> = ({ id, fill }) => (
  <marker id={id} markerWidth={8} markerHeight={8} refX={7} refY={4} orient="auto">
    <path d="M0,1 L8,4 L0,7 L1.5,4 Z" fill={fill} />
  </marker>
);

export const LegoloopDiagramComposition: React.FC = () => {
  const { width, height } = useVideoConfig();
  return (
    <AbsoluteFill style={{ backgroundColor: COLORS.paper }}>
      <svg
        width={width}
        height={height}
        viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
        role="img"
        aria-labelledby="comp-title comp-desc"
        xmlns="http://www.w3.org/2000/svg"
      >
        <title id="comp-title">Animação de entrada do fluxo legoloop-diagram-design</title>
        <desc id="comp-desc">
          Seis nós sequenciais entram uma única vez com fade e deslocamento escalonado, seguidos pelos conectores ortogonais.
        </desc>
        <defs>
          <Marker id="arr-flow" fill={COLORS.muted} />
          <Marker id="arr-accent" fill={COLORS.accent} />
          <Marker id="arr-typst" fill={COLORS.skill} />
        </defs>
        <rect width={WIDTH} height={HEIGHT} fill={COLORS.paper} rx={8} />
        {LAYERS.map((layer, i) => (
          <DiagramLayer key={layer.label} spec={layer} index={i} />
        ))}
        {EDGES.map((edge) => (
          <DiagramEdge key={edge.d} spec={edge} />
        ))}
        {NODES.map((node, i) => (
          <DiagramNode key={node.id} spec={node} index={i} />
        ))}
      </svg>
    </AbsoluteFill>
  );
};

export const RemotionRoot: React.FC = () => (
  <Composition
    id="LegoloopDiagramDesign"
    component={LegoloopDiagramComposition}
    durationInFrames={90}
    fps={FPS}
    width={WIDTH}
    height={HEIGHT}
  />
);
```