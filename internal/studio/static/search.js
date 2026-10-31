"use strict";

// Regra de busca compartilhada (inventário e, depois, a paleta do canvas).
// Sem diferenciar maiúsculas nem acentos; "-", "_", ".", "/", "\" e ":" valem como espaço;
// todas as palavras da consulta precisam aparecer, em qualquer ordem.

function searchFold(s) {
  return String(s || "")
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[-_./\\:\s]+/g, " ")
    .trim();
}

function matchesQuery(haystack, query) {
  const words = searchFold(query).split(" ").filter(Boolean);
  if (words.length === 0) return true;
  const text = searchFold(haystack);
  return words.every((w) => text.includes(w));
}
