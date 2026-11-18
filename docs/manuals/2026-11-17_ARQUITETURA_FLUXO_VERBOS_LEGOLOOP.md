# Arquitetura e Encadeamento Lógico dos Verbos no LegoLoop

> **Data:** 2026-11-17  
> **Status:** Canônico  
> **Padrão:** Vertical Slice Architecture (VSA)  
> **Diagramas:** gerados pela skill `legoloop-diagram-design` (tokens canônicos)

---

## 1. Categorização dos Verbos do LegoLoop

Os comandos do LegoLoop operam em três camadas estruturadas:

1. **Camada de Infraestrutura e Cofre:**
   - `forge`: Inicializa o ecossistema, inspeciona o ambiente e sobe o Studio.
   - `canonize`: Deduplicação mecânica por SHA-256 de arquivos e ativos brutos para o Cofre (`~/.aidd/canonical/`).
   - `doctor`: Diagnóstico de chaves, harnesses e variáveis de ambiente.

2. **Camada da Espinha Dorsal (Os 7 Verbos de Processo):**
   - `listar`: Diagnóstico de inventário, versões e vínculos entre processos.
   - `matriz`: Avaliação de aptidão e reputação de peças via contrato de I/O (`Sai -> Consome`).
   - `fabricar`: Criação de novo processo gerando a cadeia do Sexteto (INTERVIEW, INTENT, SPEC, PLAN, ROUTE, DAG).
   - `refinar`: Substituição atômica de nós e preenchimento de lacunas com preservação de versão.
   - `compor`: Encaixe hierárquico de processos como subprocessos/sub-loops sem ciclos.
   - `fechar`: Homologação canônica determinística em 9 etapas e compilação LegoPress.
   - `remover`: Desvinculação segura com cópia de restauração (`.legoloop-bak`).

3. **Camada de Execução e Qualidade:**
   - `validate`: Verificação mecânica do DAG contra o JSON Schema.
   - `exec`: Motor de execução Orchestrix em contexto estéril.
   - `upgrade`: Elevação de peças ao Padrão Ouro.

---

## 2. Diagrama de Camadas dos Verbos

<div style="font-family:Inter,system-ui,sans-serif;background:#f5f5f5;padding:24px;border-radius:8px;border:1px solid #bfc0c0;">

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 760 480" width="100%" style="display:block;max-width:760px;margin:0 auto;">
  <defs>
    <marker id="arr" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto">
      <path d="M0,0 L8,3.5 L0,7 L2,3.5 Z" fill="#4f5d75"/>
    </marker>
    <marker id="arr-acc" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto">
      <path d="M0,0 L8,3.5 L0,7 L2,3.5 Z" fill="#bf4520"/>
    </marker>
  </defs>

  <!-- BG -->
  <rect width="760" height="480" fill="#f5f5f5" rx="8"/>

  <!-- CAMADA 1: Fundação -->
  <rect x="20" y="20" width="720" height="100" fill="#ececec" rx="6" stroke="#bfc0c0" stroke-width="1"/>
  <text x="36" y="40" font-size="11" fill="#5a6580" font-weight="600" font-family="Inter,system-ui,sans-serif">CAMADA 1 — FUNDAÇÃO E COFRE</text>

  <!-- forge -->
  <rect x="60" y="52" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="130" y="72" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">forge</text>
  <text x="130" y="88" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Inicialização do ecossistema</text>

  <!-- canonize -->
  <rect x="310" y="52" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="380" y="72" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">canonize</text>
  <text x="380" y="88" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Deduplicação SHA-256</text>

  <!-- doctor -->
  <rect x="560" y="52" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="630" y="72" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">doctor</text>
  <text x="630" y="88" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Diagnóstico de ambiente</text>

  <!-- CAMADA 2: Espinha Dorsal -->
  <rect x="20" y="140" width="720" height="160" fill="#ececec" rx="6" stroke="#bf4520" stroke-width="1.5"/>
  <text x="36" y="160" font-size="11" fill="#bf4520" font-weight="600" font-family="Inter,system-ui,sans-serif">CAMADA 2 — ESPINHA DORSAL (OS 7 VERBOS DE PROCESSO)</text>

  <!-- listar -->
  <rect x="40" y="172" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="105" y="192" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">listar</text>
  <text x="105" y="208" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Inventário e vínculos</text>

  <!-- matriz -->
  <rect x="190" y="172" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="255" y="192" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">matriz</text>
  <text x="255" y="208" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Aptidão por IO</text>

  <!-- fabricar -->
  <rect x="340" y="172" width="130" height="48" fill="#2d3142" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="405" y="192" text-anchor="middle" font-size="13" font-weight="700" fill="#f5f5f5" font-family="Inter,system-ui,sans-serif">fabricar</text>
  <text x="405" y="208" text-anchor="middle" font-size="10" fill="#bfc0c0" font-family="Inter,system-ui,sans-serif">Sexteto completo</text>

  <!-- refinar -->
  <rect x="490" y="172" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="555" y="192" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">refinar</text>
  <text x="555" y="208" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Ajuste atômico de nó</text>

  <!-- compor -->
  <rect x="40" y="240" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="105" y="260" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">compor</text>
  <text x="105" y="276" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Subprocesso hierárquico</text>

  <!-- fechar -->
  <rect x="190" y="240" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="255" y="260" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">fechar</text>
  <text x="255" y="276" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">9 etapas canônicas</text>

  <!-- remover -->
  <rect x="340" y="240" width="130" height="48" fill="#fff8f6" rx="4" stroke="#4f5d75" stroke-width="1.5"/>
  <text x="405" y="260" text-anchor="middle" font-size="13" font-weight="700" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">remover</text>
  <text x="405" y="276" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Descarte com rollback</text>

  <!-- CAMADA 3: Execução -->
  <rect x="20" y="320" width="720" height="100" fill="#ececec" rx="6" stroke="#bfc0c0" stroke-width="1"/>
  <text x="36" y="340" font-size="11" fill="#5a6580" font-weight="600" font-family="Inter,system-ui,sans-serif">CAMADA 3 — EXECUÇÃO E QUALIDADE</text>

  <!-- validate -->
  <rect x="80" y="352" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="150" y="372" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">validate</text>
  <text x="150" y="388" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Gate estrutural DAG</text>

  <!-- exec -->
  <rect x="310" y="352" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="380" y="372" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">exec</text>
  <text x="380" y="388" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Motor Orchestrix estéril</text>

  <!-- upgrade -->
  <rect x="540" y="352" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="610" y="372" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">upgrade</text>
  <text x="610" y="388" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Elevação Padrão Ouro</text>

  <!-- LEGENDA -->
  <rect x="20" y="438" width="720" height="32" fill="#ececec" rx="4" stroke="#bfc0c0" stroke-width="1"/>
  <rect x="36" y="451" width="14" height="10" fill="#2d3142" rx="2"/>
  <text x="56" y="460" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Motor principal</text>
  <rect x="160" y="451" width="14" height="10" fill="#fff8f6" rx="2" stroke="#bf4520" stroke-width="1.5"/>
  <text x="180" y="460" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Verbo da espinha dorsal</text>
  <rect x="340" y="451" width="14" height="10" fill="#f5f5f5" rx="2" stroke="#4f5d75" stroke-width="1.5"/>
  <text x="360" y="460" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Verbo auxiliar</text>
  <rect x="470" y="451" width="14" height="10" fill="#f5f5f5" rx="2" stroke="#2d3142" stroke-width="1.5"/>
  <text x="490" y="460" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Fundação / Execução</text>
</svg>

</div>

---

## 3. Diagrama de Encadeamento Lógico (Fluxo do Dia a Dia)

<div style="font-family:Inter,system-ui,sans-serif;background:#f5f5f5;padding:24px;border-radius:8px;border:1px solid #bfc0c0;margin-top:12px;">

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 760 600" width="100%" style="display:block;max-width:760px;margin:0 auto;">
  <defs>
    <marker id="a1" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto">
      <path d="M0,0 L8,3.5 L0,7 L2,3.5 Z" fill="#4f5d75"/>
    </marker>
    <marker id="a2" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto">
      <path d="M0,0 L8,3.5 L0,7 L2,3.5 Z" fill="#bf4520"/>
    </marker>
    <marker id="a3" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto">
      <path d="M0,0 L8,3.5 L0,7 L2,3.5 Z" fill="#2e5aa8"/>
    </marker>
  </defs>

  <rect width="760" height="600" fill="#f5f5f5" rx="8"/>

  <!-- ROW 1: forge -->
  <rect x="310" y="20" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="380" y="40" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">forge</text>
  <text x="380" y="56" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Inicialização</text>

  <!-- forge → listar (vertical) -->
  <line x1="380" y1="68" x2="380" y2="96" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- ROW 2: listar + matriz side by side -->
  <rect x="200" y="100" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="265" y="120" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">listar</text>
  <text x="265" y="136" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Inventário</text>

  <rect x="430" y="100" width="130" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="495" y="120" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">matriz</text>
  <text x="495" y="136" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Aptidão e reputação</text>

  <!-- forge → listar (diagonal via ortogonal) -->
  <polyline points="380,68 380,84 265,84 265,100" fill="none" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>
  <!-- forge → matriz -->
  <polyline points="380,68 380,84 495,84 495,100" fill="none" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- listar → matriz -->
  <line x1="330" y1="124" x2="430" y2="124" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- matriz → fabricar -->
  <line x1="495" y1="148" x2="495" y2="176" stroke="#4f5d75" stroke-width="1.5"/>
  <polyline points="495,176 495,192 405,192 405,208" fill="none" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- ROW 3: fabricar -->
  <rect x="310" y="208" width="190" height="48" fill="#2d3142" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="405" y="228" text-anchor="middle" font-size="13" font-weight="700" fill="#f5f5f5" font-family="Inter,system-ui,sans-serif">fabricar</text>
  <text x="405" y="244" text-anchor="middle" font-size="10" fill="#bfc0c0" font-family="Inter,system-ui,sans-serif">INTERVIEW→INTENT→SPEC→PLAN→ROUTE→DAG</text>

  <!-- fabricar → validate -->
  <line x1="405" y1="256" x2="405" y2="308" stroke="#bf4520" stroke-width="1.5" marker-end="url(#a2)"/>

  <!-- ROW 4: validate -->
  <rect x="310" y="312" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="380" y="332" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">validate</text>
  <text x="380" y="348" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Gate estrutural DAG</text>

  <!-- validate → exec -->
  <line x1="380" y1="360" x2="380" y2="408" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- ROW 5: exec -->
  <rect x="310" y="412" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="380" y="432" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">exec</text>
  <text x="380" y="448" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Motor Orchestrix</text>

  <!-- exec → refinar (lacuna) -->
  <polyline points="310,436 220,436 220,320 200,320" fill="none" stroke="#bf4520" stroke-width="1.5" stroke-dasharray="6,3" marker-end="url(#a2)"/>
  <rect x="60" y="296" width="140" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="1.5"/>
  <text x="130" y="316" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">refinar</text>
  <text x="130" y="332" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Ajuste atômico de nó</text>
  <text x="220" y="380" font-size="9" fill="#bf4520" text-anchor="middle" font-family="Inter,system-ui,sans-serif">Lacuna / Falha</text>

  <!-- refinar → validate -->
  <polyline points="200,320 220,320 220,336 310,336" fill="none" stroke="#bf4520" stroke-width="1.5" marker-end="url(#a2)"/>

  <!-- exec → compor (modular) -->
  <polyline points="450,436 540,436 540,320 560,320" fill="none" stroke="#2e5aa8" stroke-width="1.5" stroke-dasharray="6,3" marker-end="url(#a3)"/>
  <rect x="560" y="296" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#2e5aa8" stroke-width="1.5"/>
  <text x="630" y="316" text-anchor="middle" font-size="13" font-weight="700" fill="#2e5aa8" font-family="Inter,system-ui,sans-serif">compor</text>
  <text x="630" y="332" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Subprocesso hierárquico</text>
  <text x="540" y="380" font-size="9" fill="#2e5aa8" text-anchor="middle" font-family="Inter,system-ui,sans-serif">Composição</text>

  <!-- compor → validate -->
  <polyline points="560,320 540,320 540,336 450,336" fill="none" stroke="#2e5aa8" stroke-width="1.5" marker-end="url(#a3)"/>

  <!-- exec → fechar (sucesso) -->
  <line x1="380" y1="460" x2="380" y2="508" stroke="#4f5d75" stroke-width="1.5" marker-end="url(#a1)"/>

  <!-- ROW 6: fechar -->
  <rect x="280" y="512" width="200" height="48" fill="#fff8f6" rx="4" stroke="#bf4520" stroke-width="2"/>
  <text x="380" y="532" text-anchor="middle" font-size="13" font-weight="700" fill="#bf4520" font-family="Inter,system-ui,sans-serif">fechar</text>
  <text x="380" y="548" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">9 etapas — Manual, Livro, Skill, DAG...</text>

  <!-- fechar → canonize -->
  <polyline points="480,536 560,536 560,76 450,76" fill="none" stroke="#4f5d75" stroke-width="1.5" stroke-dasharray="4,3" marker-end="url(#a1)"/>
  <rect x="430" y="52" width="130" height="48" fill="#f5f5f5" rx="4" stroke="#2d3142" stroke-width="1.5"/>
  <text x="495" y="72" text-anchor="middle" font-size="13" font-weight="700" fill="#2d3142" font-family="Inter,system-ui,sans-serif">canonize</text>
  <text x="495" y="88" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Cofre SHA-256</text>

  <!-- listar → remover (obsoleto) -->
  <polyline points="200,148 180,148 180,500 200,500" fill="none" stroke="#4f5d75" stroke-width="1.5" stroke-dasharray="4,3" marker-end="url(#a1)"/>
  <rect x="200" y="476" width="140" height="48" fill="#f5f5f5" rx="4" stroke="#4f5d75" stroke-width="1.5"/>
  <text x="270" y="496" text-anchor="middle" font-size="13" font-weight="700" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">remover</text>
  <text x="270" y="512" text-anchor="middle" font-size="10" fill="#4f5d75" font-family="Inter,system-ui,sans-serif">Descarte com rollback</text>
  <text x="175" y="330" font-size="9" fill="#4f5d75" text-anchor="end" font-family="Inter,system-ui,sans-serif">Obsoleto</text>

  <!-- LEGENDA -->
  <rect x="20" y="560" width="720" height="32" fill="#ececec" rx="4" stroke="#bfc0c0" stroke-width="1"/>
  <line x1="36" y1="576" x2="56" y2="576" stroke="#4f5d75" stroke-width="1.5"/>
  <path d="M54,572 L62,576 L54,580 L56,576 Z" fill="#4f5d75"/>
  <text x="68" y="580" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Fluxo principal</text>
  <line x1="160" y1="576" x2="180" y2="576" stroke="#bf4520" stroke-width="1.5" stroke-dasharray="5,3"/>
  <path d="M178,572 L186,576 L178,580 L180,576 Z" fill="#bf4520"/>
  <text x="192" y="580" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Lacuna / Falha</text>
  <line x1="300" y1="576" x2="320" y2="576" stroke="#2e5aa8" stroke-width="1.5" stroke-dasharray="5,3"/>
  <path d="M318,572 L326,576 L318,580 L320,576 Z" fill="#2e5aa8"/>
  <text x="332" y="580" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Composição modular</text>
  <line x1="470" y1="576" x2="490" y2="576" stroke="#4f5d75" stroke-width="1.5" stroke-dasharray="3,2"/>
  <path d="M488,572 L496,576 L488,580 L490,576 Z" fill="#4f5d75"/>
  <text x="502" y="580" font-size="10" fill="#2d3142" font-family="Inter,system-ui,sans-serif">Retorno ao cofre</text>
</svg>

</div>

---

## 4. Dinâmica no Dia a Dia

1. **Entrada:** O usuário consulta `listar` e `matriz` para identificar lacunas ou inventário de peças.
2. **Criação:** O verbo `fabricar` sintetiza o processo com seu grafo DAG isolado (VSA).
3. **Validação:** O comando `validate` garante conformidade com o schema antes de qualquer execução.
4. **Iteração:** Casos específicos disparam `refinar` (para trocar nó) ou `compor` (para injetar subprocesso).
5. **Canonização:** O ciclo sempre encerra no verbo `fechar`, que gera o manual, o livro didático e atualiza o cofre.

---

## 5. Tabela de Referência dos Verbos

| Verbo | Camada | Função | Tipo |
|:---|:---:|:---|:---:|
| `forge` | 1 | Inicializa ecossistema e Studio | Fundação |
| `canonize` | 1 | Deduplicação SHA-256 no Cofre | Fundação |
| `doctor` | 1 | Diagnóstico de ambiente e chaves | Fundação |
| `listar` | 2 | Inventário com vínculos e versões | Diagnóstico |
| `matriz` | 2 | Reputação e aptidão de peças por IO | Diagnóstico |
| `fabricar` | 2 | Sexteto INTERVIEW→DAG completo | **Motor principal** |
| `refinar` | 2 | Substituição atômica de nó | Iteração |
| `compor` | 2 | Encaixe de subprocesso hierárquico | Iteração |
| `fechar` | 2 | Homologação canônica em 9 etapas | Encerramento |
| `remover` | 2 | Desvinculação segura com rollback | Encerramento |
| `validate` | 3 | Gate estrutural do DAG vs schema | Qualidade |
| `exec` | 3 | Motor Orchestrix em contexto estéril | Execução |
| `upgrade` | 3 | Elevação de peças ao Padrão Ouro | Qualidade |

> **Fonte:** `internal/cli/cli.go:107-116` · **Skill:** `legoloop-diagram-design` v1.0.0
