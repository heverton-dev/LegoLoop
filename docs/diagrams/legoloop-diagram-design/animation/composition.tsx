import React from "react";
import {
  AbsoluteFill,
  Composition,
  Easing,
  interpolate,
  useCurrentFrame,
  useVideoConfig,
} from "remotion";

const TOKENS = {
  paper: "#f5f5f5",
  ink: "#2d3142",
  accent: "#bf4520",
};

const FONT_NAME = "'Geist', system-ui, 'Segoe UI', Arial, sans-serif";
const FONT_TECH = "'Geist Mono', ui-monospace, Consolas, monospace";

const FADE_MS = 200;
const STAGGER_MS = 50;
const SLIDE_PX = 12;

type NodeSpec = {
  id: string;
  x: number;
  y: number;
  tag: string;
  tagWidth: number;
  tagSpacing: number;
  name: string;
  line1: string;
  line2: string;
  focus: boolean;
};

type EdgeSpec = {
  id: string;
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  label: string;
  maskX: number;
  maskY: number;
  maskWidth: number;
  textX: number;
  textY: number;
  anchor: "start" | "middle";
};

const NODES: NodeSpec[] = [
  { id: "e01", x: 40, y: 60, tag: "Q1 ENTRADA", tagWidth: 112, tagSpacing: 0.8, name: "Inicializar slug", line1: "e01 / 5 subdiretórios", line2: "anticorpos.json", focus: false },
  { id: "e02", x: 340, y: 60, tag: "Q2 SELEÇÃO", tagWidth: 112, tagSpacing: 0.8, name: "Selecionar tipo", line1: "e02 / 1 de 50 tipos", line2: "máx. 9 nós, 12 setas", focus: false },
  { id: "e03", x: 640, y: 60, tag: "Q3 GERAÇÃO", tagWidth: 112, tagSpacing: 0.8, name: "Gerar SVG e HTML", line1: "e03 / grid ortogonal", line2: "máscaras opacas", focus: true },
  { id: "e04", x: 640, y: 240, tag: "Q3 GERAÇÃO", tagWidth: 112, tagSpacing: 0.8, name: "Animar em CSS", line1: "e04 / flow-dash", line2: "stroke-dashoffset", focus: false },
  { id: "e05", x: 340, y: 240, tag: "Q4 FECHAMENTO", tagWidth: 112, tagSpacing: 0.4, name: "Exportar e validar", line1: "e05 / SVG e PNG", line2: "gate com exit code 0", focus: true },
  { id: "e06", x: 40, y: 240, tag: "Q4 FECHAMENTO", tagWidth: 112, tagSpacing: 0.4, name: "Fechar e relatar", line1: "e06 / anticorpos.json", line2: "relatório em docs/reports", focus: false },
];

const EDGES: EdgeSpec[] = [
  { id: "a12", x1: 240, y1: 108, x2: 340, y2: 108, label: "MANIFESTO", maskX: 246, maskY: 84, maskWidth: 88, textX: 290, textY: 96, anchor: "middle" },
  { id: "a23", x1: 540, y1: 108, x2: 640, y2: 108, label: "TIPO VISUAL", maskX: 546, maskY: 84, maskWidth: 88, textX: 590, textY: 96, anchor: "middle" },
  { id: "a34", x1: 740, y1: 156, x2: 740, y2: 240, label: "HTML SRC", maskX: 748, maskY: 190, maskWidth: 80, textX: 754, textY: 202, anchor: "start" },
  { id: "a45", x1: 640, y1: 288, x2: 540, y2: 288, label: "ANIMADO", maskX: 546, maskY: 264, maskWidth: 88, textX: 590, textY: 276, anchor: "middle" },
  { id: "a56", x1: 340, y1: 288, x2: 240, y2: 288, label: "GATE EXIT 0", maskX: 246, maskY: 264, maskWidth: 88, textX: 290, textY: 276, anchor: "middle" },
];

const useFadeSlide = (delayFrames: number) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const durationFrames = Math.max(1, Math.round((FADE_MS / 1000) * fps));
  const local = frame - delayFrames;
  const opacity = interpolate(local, [0, durationFrames], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
    easing: Easing.out(Easing.cubic),
  });
  const offset = interpolate(local, [0, durationFrames], [SLIDE_PX, 0], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
    easing: Easing.out(Easing.cubic),
  });
  return { opacity, offset };
};

const Node: React.FC<{ spec: NodeSpec; index: number }> = ({ spec, index }) => {
  const { fps } = useVideoConfig();
  const staggerFrames = Math.max(1, Math.round((STAGGER_MS / 1000) * fps));
  const { opacity, offset } = useFadeSlide(index * staggerFrames);
  const stroke = spec.focus ? TOKENS.accent : TOKENS.ink;
  const fill = spec.focus ? TOKENS.accent : TOKENS.ink;
  const fillOpacity = spec.focus ? 0.08 : 0.04;
  const tagStrokeOpacity = spec.focus ? 1 : 0.72;
  const tagY = spec.y + 14;
  const textBase = spec.y;

  return (
    <g opacity={opacity} transform={`translate(0 ${offset})`}>
      <rect x={spec.x} y={spec.y} width={200} height={96} rx={6} fill={TOKENS.paper} />
      <rect x={spec.x} y={spec.y} width={200} height={96} rx={6} fill={fill} fillOpacity={fillOpacity} stroke={stroke} strokeWidth={1.5} />
      <rect x={spec.x + 16} y={tagY} width={spec.tagWidth} height={16} rx={2} fill="none" stroke={stroke} strokeOpacity={tagStrokeOpacity} />
      <text x={spec.x + 24} y={tagY + 12} fontFamily={FONT_TECH} fontSize={10} letterSpacing={spec.tagSpacing} fill={spec.focus ? TOKENS.accent : TOKENS.ink}>
        {spec.tag}
      </text>
      <text x={spec.x + 16} y={textBase + 52} fontFamily={FONT_NAME} fontSize={14} fontWeight={600} fill={TOKENS.ink}>
        {spec.name}
      </text>
      <text x={spec.x + 16} y={textBase + 70} fontFamily={FONT_TECH} fontSize={10} fill={TOKENS.ink} fillOpacity={0.78}>
        {spec.line1}
      </text>
      <text x={spec.x + 16} y={textBase + 84} fontFamily={FONT_TECH} fontSize={10} fill={TOKENS.ink} fillOpacity={0.78}>
        {spec.line2}
      </text>
    </g>
  );
};

const Edge: React.FC<{ spec: EdgeSpec; index: number }> = ({ spec, index }) => {
  const { fps } = useVideoConfig();
  const staggerFrames = Math.max(1, Math.round((STAGGER_MS / 1000) * fps));
  const base = NODES.length * staggerFrames + Math.round((FADE_MS / 1000) * fps);
  const { opacity } = useFadeSlide(base + index * staggerFrames);

  return (
    <g opacity={opacity}>
      <line x1={spec.x1} y1={spec.y1} x2={spec.x2} y2={spec.y2} stroke={TOKENS.ink} strokeOpacity={0.72} strokeWidth={1.5} markerEnd="url(#legoloop-arrow)" />
      <rect x={spec.maskX} y={spec.maskY} width={spec.maskWidth} height={16} fill={TOKENS.paper} />
      <text x={spec.textX} y={spec.textY} textAnchor={spec.anchor} fontFamily={FONT_TECH} fontSize={10} letterSpacing={0.8} fill={TOKENS.ink}>
        {spec.label}
      </text>
    </g>
  );
};

const Legend: React.FC = () => {
  const { fps } = useVideoConfig();
  const staggerFrames = Math.max(1, Math.round((STAGGER_MS / 1000) * fps));
  const base =
    (NODES.length + EDGES.length) * staggerFrames +
    Math.round((FADE_MS / 1000) * fps);
  const { opacity } = useFadeSlide(base);

  return (
    <g opacity={opacity}>
      <line x1={40} y1={360} x2={840} y2={360} stroke={TOKENS.ink} strokeOpacity={0.24} strokeWidth={1} />
      <rect x={40} y={380} width={24} height={16} rx={6} fill={TOKENS.ink} fillOpacity={0.04} stroke={TOKENS.ink} strokeWidth={1.5} />
      <text x={76} y={392} fontFamily={FONT_TECH} fontSize={10} fill={TOKENS.ink}>Etapa do processo</text>
      <rect x={272} y={380} width={24} height={16} rx={6} fill={TOKENS.accent} fillOpacity={0.08} stroke={TOKENS.accent} strokeWidth={1.5} />
      <text x={308} y={392} fontFamily={FONT_TECH} fontSize={10} fill={TOKENS.ink}>Foco editorial</text>
      <line x1={480} y1={388} x2={512} y2={388} stroke={TOKENS.ink} strokeOpacity={0.72} strokeWidth={1.5} markerEnd="url(#legoloop-arrow)" />
      <text x={528} y={392} fontFamily={FONT_TECH} fontSize={10} fill={TOKENS.ink}>Saída verificada pelo gate</text>
    </g>
  );
};

export const LegoLoopComposition: React.FC = () => {
  return (
    <AbsoluteFill style={{ backgroundColor: TOKENS.paper }}>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 880 420"
        width="100%"
        height="100%"
        role="img"
        aria-labelledby="legoloop-diagram-design-anim-title legoloop-diagram-design-anim-desc"
      >
        <title id="legoloop-diagram-design-anim-title">Fluxo de seis nós do LegoLoop Diagram Design</title>
        <desc id="legoloop-diagram-design-anim-desc">
          Animação de entrada única dos seis nós, das cinco conexões e da legenda do processo.
        </desc>
        <defs>
          <marker id="legoloop-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="8" markerHeight="8" orient="auto-start-reverse">
            <path d="M 0 0 L 10 5 L 0 10 z" fill={TOKENS.ink} fillOpacity={0.72} />
          </marker>
        </defs>
        <rect x={0} y={0} width={880} height={420} fill={TOKENS.paper} />
        {EDGES.map((edge, index) => (
          <Edge key={edge.id} spec={edge} index={index} />
        ))}
        {NODES.map((node, index) => (
          <Node key={node.id} spec={node} index={index} />
        ))}
        <Legend />
      </svg>
    </AbsoluteFill>
  );
};

export const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="LegoLoopDiagramDesign"
      component={LegoLoopComposition}
      durationInFrames={90}
      fps={30}
      width={1760}
      height={840}
    />
  );
};