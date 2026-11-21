import os

def make_lego_svg(part_num, part_name, color_theme, filename):
    palettes = {
        'green': {
            'top': '#34d399', 'top_grad': '#10b981', 'front': '#059669', 'front_dark': '#047857',
            'side': '#065f46', 'side_dark': '#064e3b', 'stud_top': '#6ee7b7', 'stud_side': '#10b981',
            'stud_side_dark': '#047857', 'highlight': '#a7f3d0', 'shadow': '#022c22', 'accent': '#10b981'
        },
        'blue': {
            'top': '#60a5fa', 'top_grad': '#3b82f6', 'front': '#2563eb', 'front_dark': '#1d4ed8',
            'side': '#1e40af', 'side_dark': '#1e3a8a', 'stud_top': '#93c5fd', 'stud_side': '#3b82f6',
            'stud_side_dark': '#1d4ed8', 'highlight': '#bfdbfe', 'shadow': '#0f172a', 'accent': '#3b82f6'
        },
        'amber': {
            'top': '#fbbf24', 'top_grad': '#f59e0b', 'front': '#d97706', 'front_dark': '#b45309',
            'side': '#92400e', 'side_dark': '#78350f', 'stud_top': '#fde68a', 'stud_side': '#f59e0b',
            'stud_side_dark': '#b45309', 'highlight': '#fef3c7', 'shadow': '#451a03', 'accent': '#f59e0b'
        },
        'violet': {
            'top': '#a78bfa', 'top_grad': '#8b5cf6', 'front': '#7c3aed', 'front_dark': '#6d28d9',
            'side': '#5b21b6', 'side_dark': '#4c1d95', 'stud_top': '#c4b5fd', 'stud_side': '#8b5cf6',
            'stud_side_dark': '#6d28d9', 'highlight': '#ddd6fe', 'shadow': '#2e1065', 'accent': '#8b5cf6'
        }
    }
    pal = palettes.get(color_theme, palettes['green'])
    
    # 8 Studs em 2 fileiras perfeitamente espaçadas (Proporção 2x4 clássica Lego)
    u_vals = [0.155, 0.385, 0.615, 0.845]
    studs_back = []
    studs_front = []
    
    for u in u_vals:
        xb = 110 - 15 + u * 280
        yb = 70
        studs_back.append((round(xb, 1), yb))
        
        xf = 110 - 37 + u * 280
        yf = 86
        studs_front.append((round(xf, 1), yf))
        
    rx = 16.5
    ry = 6.4
    h_stud = 11.5
    
    svg = []
    svg.append('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 460 220" width="100%" height="100%">')
    svg.append('  <defs>')
    # Sombra ampla e suave de contato no chão
    svg.append('    <filter id="dropShad" x="-15%" y="-10%" width="135%" height="145%">')
    svg.append('      <feDropShadow dx="0" dy="16" stdDeviation="12" flood-color="#0f172a" flood-opacity="0.22"/>')
    svg.append('    </filter>')
    
    # Gradiente da Face Superior (Luz Zenital + Brilho Especular plástico)
    svg.append(f'    <linearGradient id="topGrad" x1="0%" y1="0%" x2="100%" y2="100%">')
    svg.append(f'      <stop offset="0%" stop-color="{pal["top"]}"/>')
    svg.append(f'      <stop offset="60%" stop-color="{pal["top_grad"]}"/>')
    svg.append(f'      <stop offset="100%" stop-color="{pal["front"]}"/>')
    svg.append('    </linearGradient>')
    
    # Gradiente da Face Frontal (Luz Difusa com Reflexo Acetinado)
    svg.append(f'    <linearGradient id="frontGrad" x1="0%" y1="0%" x2="0%" y2="100%">')
    svg.append(f'      <stop offset="0%" stop-color="{pal["front"]}"/>')
    svg.append(f'      <stop offset="15%" stop-color="{pal["top_grad"]}"/>')
    svg.append(f'      <stop offset="35%" stop-color="{pal["front"]}"/>')
    svg.append(f'      <stop offset="85%" stop-color="{pal["front_dark"]}"/>')
    svg.append(f'      <stop offset="100%" stop-color="{pal["side"]}"/>')
    svg.append('    </linearGradient>')
    
    # Gradiente da Face Lateral Direita (Sombra Própria)
    svg.append(f'    <linearGradient id="sideGrad" x1="0%" y1="0%" x2="100%" y2="100%">')
    svg.append(f'      <stop offset="0%" stop-color="{pal["side"]}"/>')
    svg.append(f'      <stop offset="100%" stop-color="{pal["side_dark"]}"/>')
    svg.append('    </linearGradient>')
    
    # Gradiente do Cilindro do Stud
    svg.append(f'    <linearGradient id="studCyl" x1="0%" y1="0%" x2="100%" y2="0%">')
    svg.append(f'      <stop offset="0%" stop-color="{pal["stud_top"]}"/>')
    svg.append(f'      <stop offset="25%" stop-color="{pal["stud_side"]}"/>')
    svg.append(f'      <stop offset="75%" stop-color="{pal["stud_side_dark"]}"/>')
    svg.append(f'      <stop offset="100%" stop-color="{pal["side"]}"/>')
    svg.append('    </linearGradient>')
    
    # Gradiente da Tampa do Stud
    svg.append(f'    <linearGradient id="studCap" x1="0%" y1="0%" x2="0%" y2="100%">')
    svg.append(f'      <stop offset="0%" stop-color="{pal["highlight"]}"/>')
    svg.append(f'      <stop offset="50%" stop-color="{pal["stud_top"]}"/>')
    svg.append(f'      <stop offset="100%" stop-color="{pal["stud_side"]}"/>')
    svg.append('    </linearGradient>')
    svg.append('  </defs>')
    
    svg.append('  <g filter="url(#dropShad)">')
    
    # 1. Face Superior do Bloco
    svg.append(f'    <polygon points="60,95 110,60 390,60 340,95" fill="url(#topGrad)" stroke="{pal["highlight"]}" stroke-width="0.75"/>')
    
    # 2. Face Lateral Direita
    svg.append(f'    <polygon points="340,95 390,60 390,135 340,170" fill="url(#sideGrad)" stroke="{pal["shadow"]}" stroke-width="0.75"/>')
    
    # 3. Função para renderizar um Stud autêntico LEGO
    def render_stud(cx, cy):
        lines = []
        lines.append(f'    <g transform="translate({cx}, {cy})">')
        # Sombra de oclusão na base do stud
        lines.append(f'      <ellipse cx="0" cy="0" rx="{rx}" ry="{ry}" fill="{pal["shadow"]}" opacity="0.35"/>')
        # Corpo cilíndrico do stud
        lines.append(f'      <path d="M -{rx} -{h_stud} L -{rx} 0 A {rx} {ry} 0 0 0 {rx} 0 L {rx} -{h_stud} Z" fill="url(#studCyl)"/>')
        # Tampa elíptica superior do stud
        lines.append(f'      <ellipse cx="0" cy="-{h_stud}" rx="{rx}" ry="{ry}" fill="url(#studCap)" stroke="{pal["highlight"]}" stroke-width="0.8"/>')
        # Relevo duplo concêntrico característico LEGO (stud ring)
        lines.append(f'      <ellipse cx="0" cy="-{h_stud}" rx="{round(rx*0.62,1)}" ry="{round(ry*0.62,1)}" fill="none" stroke="{pal["highlight"]}" stroke-width="1.3" opacity="0.9"/>')
        lines.append(f'      <ellipse cx="0.5" cy="-{h_stud-0.4}" rx="{round(rx*0.56,1)}" ry="{round(ry*0.56,1)}" fill="none" stroke="{pal["shadow"]}" stroke-width="0.6" opacity="0.45"/>')
        # Inscrição icônica estilizada de relevo "LEGO" no centro do pino
        lines.append(f'      <text x="0" y="-{round(h_stud-1.5,1)}" font-family="system-ui, -apple-system, sans-serif" font-size="4.8" font-weight="900" font-style="italic" fill="{pal["highlight"]}" text-anchor="middle" letter-spacing="0.4" opacity="0.8">LEGO</text>')
        lines.append('    </g>')
        return '\n'.join(lines)
        
    # Studs da fileira traseira
    for cx, cy in studs_back:
        svg.append(render_stud(cx, cy))
        
    # Studs da fileira frontal
    for cx, cy in studs_front:
        svg.append(render_stud(cx, cy))
        
    # 4. Face Frontal do Bloco (Com cantos sutilmente arredondados para acabamento de injeção plástica)
    svg.append(f'    <rect x="60" y="95" width="280" height="75" rx="3" fill="url(#frontGrad)" stroke="{pal["front_dark"]}" stroke-width="1"/>')
    
    # Chanfro/Bisel de Luz na aresta superior frontal (Highlight plástico de alta precisão)
    svg.append(f'    <line x1="61" y1="96" x2="339" y2="96" stroke="{pal["highlight"]}" stroke-width="2" opacity="0.95"/>')
    # Chanfro de sombra na aresta inferior frontal
    svg.append(f'    <line x1="61" y1="169" x2="339" y2="169" stroke="{pal["shadow"]}" stroke-width="2.2" opacity="0.75"/>')
    # Aresta vertical iluminada à esquerda
    svg.append(f'    <line x1="61" y1="96" x2="61" y2="169" stroke="{pal["highlight"]}" stroke-width="1.5" opacity="0.6"/>')
    # Aresta vertical sombreada à direita
    svg.append(f'    <line x1="339" y1="96" x2="339" y2="169" stroke="{pal["shadow"]}" stroke-width="1.5" opacity="0.5"/>')
    
    # 5. Texto gravado na face frontal
    svg.append(f'    <text x="200" y="142" font-family="system-ui, -apple-system, Segoe UI, Roboto, Helvetica, sans-serif" font-size="26" font-weight="900" fill="#ffffff" text-anchor="middle" letter-spacing="4" style="text-shadow: 0 2px 5px rgba(0,0,0,0.55);">')
    svg.append(f'      {part_name}')
    svg.append(f'    </text>')
    
    # Linha técnica de relevo de injeção plástica abaixo do texto
    svg.append('    <rect x="110" y="152" width="180" height="2" fill="#ffffff" opacity="0.25" rx="1"/>')
    
    svg.append('  </g>')
    svg.append('</svg>')
    
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w', encoding='utf-8') as f:
        f.write('\n'.join(svg))
    print(f'Criado {filename}')

if __name__ == '__main__':
    make_lego_svg('1', 'PARTE 1', 'green', 'docs/books/ativos/lego_brick_part1.svg')
    make_lego_svg('2', 'PARTE 2', 'blue', 'docs/books/ativos/lego_brick_part2.svg')
    make_lego_svg('3', 'PARTE 3', 'amber', 'docs/books/ativos/lego_brick_part3.svg')
    make_lego_svg('4', 'PARTE 4', 'violet', 'docs/books/ativos/lego_brick_part4.svg')
