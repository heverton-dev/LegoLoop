// icons.js — Sistema Canônico de Ícones SVG do LegoLoop Studio (Regra Zero Emojis & Zero innerHTML)
// 100% DOM seguro: cria elementos via document.createElementNS.

const SVG_XMLNS = "http://www.w3.org/2000/svg";

function createSvgElement(tag, attrs) {
  const el = document.createElementNS(SVG_XMLNS, tag);
  for (const [k, v] of Object.entries(attrs || {})) {
    el.setAttribute(k, v);
  }
  return el;
}

function createSvgIconNode(name) {
  const svg = createSvgElement("svg", {
    class: "icon-svg",
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    "stroke-width": "2",
    "stroke-linecap": "round",
    "stroke-linejoin": "round"
  });

  switch (name) {
    case "template":
      svg.append(
        createSvgElement("path", { d: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" }),
        createSvgElement("polyline", { points: "14 2 14 8 20 8" }),
        createSvgElement("line", { x1: "16", y1: "13", x2: "8", y2: "13" }),
        createSvgElement("line", { x1: "16", y1: "17", x2: "8", y2: "17" })
      );
      break;
    case "vault":
      svg.append(
        createSvgElement("rect", { x: "3", y: "3", width: "18", height: "18", rx: "2" }),
        createSvgElement("circle", { cx: "12", cy: "12", r: "3" }),
        createSvgElement("line", { x1: "12", y1: "7", x2: "12", y2: "9" }),
        createSvgElement("line", { x1: "12", y1: "15", x2: "12", y2: "17" })
      );
      break;
    case "camera":
      svg.append(
        createSvgElement("path", { d: "M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z" }),
        createSvgElement("circle", { cx: "12", cy: "13", r: "4" })
      );
      break;
    case "trigger":
      svg.append(
        createSvgElement("polygon", { points: "13 2 3 14 12 14 11 22 21 10 12 10 13 2" })
      );
      break;
    case "search":
      svg.append(
        createSvgElement("circle", { cx: "11", cy: "11", r: "8" }),
        createSvgElement("line", { x1: "21", y1: "21", x2: "16.65", y2: "16.65" })
      );
      break;
    case "shield":
      svg.append(
        createSvgElement("path", { d: "M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" })
      );
      break;
    case "plant":
      svg.append(
        createSvgElement("path", { d: "M7 20h10" }),
        createSvgElement("path", { d: "M12 20v-8" }),
        createSvgElement("path", { d: "M12 12a5 5 0 0 1 5-5 5 5 0 0 0-5-5 5 5 0 0 0-5 5 5 5 0 0 1 5 5z" })
      );
      break;
    case "undo":
      svg.append(
        createSvgElement("path", { d: "M3 7v6h6" }),
        createSvgElement("path", { d: "M21 17a9 9 0 0 0-9-9 9 9 0 0 0-6 2.3L3 13" })
      );
      break;
    case "redo":
      svg.append(
        createSvgElement("path", { d: "M21 7v6h-6" }),
        createSvgElement("path", { d: "M3 17a9 9 0 0 1 9-9 9 9 0 0 1 6 2.3l3 2.7" })
      );
      break;
    case "git":
      svg.append(
        createSvgElement("circle", { cx: "18", cy: "18", r: "3" }),
        createSvgElement("circle", { cx: "6", cy: "6", r: "3" }),
        createSvgElement("path", { d: "M18 15V9a6 6 0 0 0-6-6" }),
        createSvgElement("line", { x1: "6", y1: "9", x2: "6", y2: "21" })
      );
      break;
    case "deploy":
      svg.append(
        createSvgElement("path", { d: "M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" }),
        createSvgElement("line", { x1: "12", y1: "11", x2: "12", y2: "17" }),
        createSvgElement("polyline", { points: "9 14 12 11 15 14" })
      );
      break;
    case "palette":
      svg.append(
        createSvgElement("circle", { cx: "13.5", cy: "6.5", r: ".5", fill: "currentColor" }),
        createSvgElement("circle", { cx: "17.5", cy: "10.5", r: ".5", fill: "currentColor" }),
        createSvgElement("circle", { cx: "8.5", cy: "7.5", r: ".5", fill: "currentColor" }),
        createSvgElement("circle", { cx: "6.5", cy: "12.5", r: ".5", fill: "currentColor" }),
        createSvgElement("path", { d: "M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z" })
      );
      break;
    default: // package
      svg.append(
        createSvgElement("path", { d: "M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" }),
        createSvgElement("polyline", { points: "3.27 6.96 12 12.01 20.73 6.96" }),
        createSvgElement("line", { x1: "12", y1: "22.08", x2: "12", y2: "12" })
      );
      break;
  }
  return svg;
}
