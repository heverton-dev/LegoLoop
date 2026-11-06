"use strict";

// Lógica e renderização de processos do cofre: listar, refinar, remover, compor e cartão STE-PT.
// Invariantes estritas: zero innerHTML, zero emojis, zero diálogos nativos.

function el(tag, cls, text) {
  const node = document.createElement(tag);
  if (cls) node.className = cls;
  if (text !== undefined && text !== null) node.textContent = text;
  return node;
}

function renderCartaoSTE(cartao) {
  const box = el("div", "cartao-ste");
  if (!cartao) return box;

  const header = el("div", "cartao-ste__header");
  const title = el("h4", "cartao-ste__title", cartao.acao || "Ação de Governança");
  const badgeCls = cartao.irreversivel ? "modal-badge badge--danger" : "modal-badge badge--canon";
  const badgeText = cartao.irreversivel ? "Ação Irreversível" : "Ação Reversível";
  const badge = el("span", badgeCls, badgeText);
  header.append(title, badge);
  box.append(header);

  // Aviso de risco obrigatório antes dos passos
  if (cartao.aviso) {
    const avisoBox = el("div", "cartao-ste__aviso");
    avisoBox.setAttribute("role", "alert");
    const avisoRotulo = el("strong", "cartao-ste__aviso-rotulo", "Aviso de Risco: ");
    const avisoTexto = el("span", "cartao-ste__aviso-texto", cartao.aviso);
    avisoBox.append(avisoRotulo, avisoTexto);
    box.append(avisoBox);
  }

  // Passos que serão executados
  if (Array.isArray(cartao.passos) && cartao.passos.length > 0) {
    const secPassos = el("div", "cartao-ste__secao");
    secPassos.append(el("h5", "cartao-ste__subtitulo", "O que vai acontecer:"));
    const ol = el("ol", "cartao-ste__passos");
    for (const passo of cartao.passos) {
      ol.append(el("li", "cartao-ste__passo-item", passo));
    }
    secPassos.append(ol);
    box.append(secPassos);
  }

  // Efeitos colaterais
  if (Array.isArray(cartao.efeitos) && cartao.efeitos.length > 0) {
    const secEfeitos = el("div", "cartao-ste__secao");
    secEfeitos.append(el("h5", "cartao-ste__subtitulo", "Efeitos colaterais:"));
    const ul = el("ul", "cartao-ste__efeitos");
    for (const efeito of cartao.efeitos) {
      ul.append(el("li", "cartao-ste__efeito-item", efeito));
    }
    secEfeitos.append(ul);
    box.append(secEfeitos);
  }

  // Instrução de reversão / se der errado
  if (cartao.se_der_errado) {
    const secDesfazer = el("div", "cartao-ste__desfazer");
    secDesfazer.append(el("strong", "", "Se der errado: "));
    secDesfazer.append(el("span", "", cartao.se_der_errado));
    box.append(secDesfazer);
  }

  return box;
}

function renderRelatorioCamadas(data) {
  const box = el("div", "relatorio-camadas");
  if (!data) return box;

  // Caso 1: Relatório das 5 camadas da fabricação
  if (data.clarificacao || data.varredura || data.fusao || data.triade || data.gravacao) {
    box.append(renderCamadasFabricacao(data));
    return box;
  }

  // Caso 2: Relatório de Refinar / Compor / Remover (diff e arquivos)
  box.append(renderCamadasMutacao(data));
  return box;
}

function renderCamadasFabricacao(data) {
  const container = el("div", "relatorio-camadas__container");
  const titulo = el("h4", "relatorio-camadas__titulo", "Relatório de Fabricação por Camadas");
  container.append(titulo);

  // Camada 1: Clarificação
  if (data.clarificacao) {
    const c = data.clarificacao;
    const c1 = el("div", "relatorio-camada");
    c1.append(el("h5", "relatorio-camada__nome", "Camada 1: Clarificação"));
    const info = el("p", "relatorio-camada__info", `Confiança: ${c.confianca}% | Turnos: ${c.turnos} | Alvo alcançado: ${c.alcancou ? "Sim" : "Não"}`);
    c1.append(info);
    if (c.espelho && c.espelho.objetivo) {
      c1.append(el("p", "relatorio-camada__detalhe", `Objetivo identificado: ${c.espelho.objetivo}`));
    }
    container.append(c1);
  }

  // Camada 2: Varredura da Matriz
  if (data.varredura) {
    const v = data.varredura;
    const c2 = el("div", "relatorio-camada");
    c2.append(el("h5", "relatorio-camada__nome", "Camada 2: Varredura"));
    const totalPecas = v.pecas_vistas ? Object.values(v.pecas_vistas).reduce((a, b) => a + b, 0) : 0;
    c2.append(el("p", "relatorio-camada__info", `Peças examinadas na matriz: ${totalPecas}`));
    if (Array.isArray(v.cortes) && v.cortes.length > 0) {
      c2.append(el("p", "relatorio-camada__detalhe", `Cortes aplicados: ${v.cortes.join(", ")}`));
    }
    container.append(c2);
  }

  // Camada 3: Fusão e Seleção
  if (data.fusao) {
    const f = data.fusao;
    const c3 = el("div", "relatorio-camada");
    c3.append(el("h5", "relatorio-camada__nome", "Camada 3: Fusão"));
    const numDecisoes = Array.isArray(f.decisoes) ? f.decisoes.length : 0;
    const numLacunas = Array.isArray(f.lacunas) ? f.lacunas.length : 0;
    c3.append(el("p", "relatorio-camada__info", `Decisões de peças: ${numDecisoes} | Lacunas detectadas: ${numLacunas}`));
    container.append(c3);
  }

  // Camada 4: Tríade
  if (data.triade) {
    const t = data.triade;
    const c4 = el("div", "relatorio-camada");
    c4.append(el("h5", "relatorio-camada__nome", "Camada 4: Tríade de Governança"));
    c4.append(el("p", "relatorio-camada__info", `Critérios (ACs): ${t.acs} | Tickets: ${t.tickets} | Coerente: ${t.coerente ? "Sim" : "Não"}`));
    container.append(c4);
  }

  // Camada 5: Gravação
  if (data.gravacao) {
    const g = data.gravacao;
    const c5 = el("div", "relatorio-camada");
    c5.append(el("h5", "relatorio-camada__nome", "Camada 5: Gravação Atômica"));
    const numArqs = Array.isArray(g.arquivos) ? g.arquivos.length : 0;
    const numHarn = Array.isArray(g.harnesses) ? g.harnesses.length : 0;
    c5.append(el("p", "relatorio-camada__info", `Arquivos gravados: ${numArqs} | Vínculos em harnesses: ${numHarn} | Tamanho: ${g.bytes} bytes`));
    container.append(c5);
  }

  return container;
}

function renderCamadasMutacao(data) {
  const container = el("div", "relatorio-camadas__container");
  const titulo = el("h4", "relatorio-camadas__titulo", "Relatório de Mudanças");
  container.append(titulo);

  // Camada 1: Arquivos afetados
  if (Array.isArray(data.arquivos) && data.arquivos.length > 0) {
    const c1 = el("div", "relatorio-camada");
    c1.append(el("h5", "relatorio-camada__nome", "Arquivos no Escopo da Operação:"));
    const ul = el("ul", "relatorio-camada__arquivos");
    for (const arq of data.arquivos) {
      ul.append(el("li", "", arq));
    }
    c1.append(ul);
    container.append(c1);
  }

  // Camada 2: Diff de alterações
  if (data.diff) {
    const c2 = el("div", "relatorio-camada");
    c2.append(el("h5", "relatorio-camada__nome", "Diferenças Planejadas (Diff):"));
    const pre = el("pre", "diff-view", data.diff);
    c2.append(pre);
    container.append(c2);
  }

  // Camada 3: Avisos
  if (Array.isArray(data.avisos) && data.avisos.length > 0) {
    const c3 = el("div", "relatorio-camada");
    c3.append(el("h5", "relatorio-camada__nome", "Avisos de Governança:"));
    const ul = el("ul", "relatorio-camada__avisos");
    for (const aviso of data.avisos) {
      ul.append(el("li", "", aviso));
    }
    c3.append(ul);
    container.append(c3);
  }

  return container;
}

function renderProcessoItem(proc, onAction) {
  const item = el("div", "processo-card");
  item.dataset.nome = proc.nome;

  const header = el("div", "processo-card__header");
  const titleBox = el("div", "processo-card__title-box");
  const title = el("h4", "processo-card__nome", proc.nome);
  const versao = el("span", "badge", proc.versao || "v1");
  titleBox.append(title, versao);

  const actions = el("div", "processo-card__actions");
  const btnRefinar = el("button", "btn btn--sm", "Refinar");
  btnRefinar.type = "button";
  btnRefinar.title = "Refinar etapas ou trocar peças do processo";
  btnRefinar.addEventListener("click", () => onAction("refinar", proc));

  const btnCompor = el("button", "btn btn--sm", "Compor");
  btnCompor.type = "button";
  btnCompor.title = "Compor este processo como sub-etapa de outro";
  btnCompor.addEventListener("click", () => onAction("compor", proc));

  const btnRemover = el("button", "btn btn--sm btn--danger", "Remover");
  btnRemover.type = "button";
  btnRemover.title = "Remover processo do cofre e limpar vínculos";
  btnRemover.addEventListener("click", () => onAction("remover", proc));

  actions.append(btnRefinar, btnCompor, btnRemover);
  header.append(titleBox, actions);
  item.append(header);

  // Meta dados
  const meta = el("div", "processo-card__meta");
  if (proc.descricao || proc.objetivo) {
    meta.append(el("p", "processo-card__desc", proc.descricao || proc.objetivo));
  }
  const tags = el("div", "processo-card__tags");
  if (Array.isArray(proc.etapas)) {
    tags.append(el("span", "processo-tag", `${proc.etapas.length} etapa(s)`));
  }
  if (Array.isArray(proc.harnesses)) {
    tags.append(el("span", "processo-tag", `${proc.harnesses.length} harness(es)`));
  }
  meta.append(tags);
  item.append(meta);

  return item;
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = {
    renderCartaoSTE,
    renderRelatorioCamadas,
    renderProcessoItem,
  };
}
