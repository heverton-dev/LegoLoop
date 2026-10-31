function reorderForLargeContext(models, estTokens, log) {
  if (estTokens <= 100000) return models;

  // 1M models: gemini models (ag/gemini-*, gemini/*)
  const is1MModel = (m) => {
    const lower = m.toLowerCase();
    return (lower.includes("gemini") || lower.includes("flash")) && !lower.startsWith("kr/");
  };
  const isKiro = (m) => m.toLowerCase().startsWith("kr/") || m.toLowerCase().includes("kiro");

  const millionModels = models.filter(is1MModel);
  const otherModels = models.filter(m => !is1MModel(m) && !isKiro(m));
  const kiroModels = models.filter(isKiro);

  // Preferred order: Gemini 3.8/3.6/3.7 Flash first, then other Gemini 1M, then others, and demote Kiro
  millionModels.sort((a, b) => {
    const aFlash = a.includes("flash") ? 1 : 0;
    const bFlash = b.includes("flash") ? 1 : 0;
    return bFlash - aFlash;
  });

  const reordered = [...millionModels, ...otherModels, ...kiroModels];
  if (reordered[0] !== models[0] && log) {
    log.info("COMBO", `large-context (~${estTokens} tokens) → prioritizing 1M models (${reordered[0]})`);
  }
  return reordered;
}

// Test with code-pro models:
const codePro = ["ag/gemini-3.1-pro-low", "ag/claude-opus-4-6-thinking", "gemini/gemini-3.1-pro-preview", "ag/claude-sonnet-4-6", "kr/claude-sonnet-4.5"];

console.log("Normal (<100k):", reorderForLargeContext(codePro, 50000, console));
console.log("Large (>100k):", reorderForLargeContext(codePro, 130000, console));
