"use strict";

// Janela de virtualização para listas longas (paleta com milhares de peças).
// Todas as linhas têm a mesma altura (rowH): só as visíveis, mais uma folga, entram no DOM.

function virtualWindow(total, rowH, scrollTop, viewH, overscan = 6) {
  if (total <= 0 || rowH <= 0) return { start: 0, end: 0, offset: 0, height: 0 };
  const first = Math.floor(Math.max(0, scrollTop) / rowH);
  const visible = Math.ceil(Math.max(0, viewH) / rowH) + 1;
  const start = Math.min(Math.max(0, first - overscan), total - 1);
  const end = Math.min(total, first + visible + overscan);
  return { start, end, offset: start * rowH, height: total * rowH };
}
