import os

def create_lego_cover_svg(filename):
    palettes = {
        'blue': {
            'top': '#60a5fa', 'top_grad': '#3b82f6', 'front': '#2563eb', 'front_dark': '#1d4ed8',
            'side': '#1e40af', 'side_dark': '#1e3a8a', 'stud_top': '#93c5fd', 'stud_side': '#3b82f6',
            'stud_side_dark': '#1d4ed8', 'highlight': '#bfdbfe', 'shadow': '#0f172a'
        },
        'green': {
            'top': '#34d399', 'top_grad': '#10b981', 'front': '#059669', 'front_dark': '#047857',
            'side': '#065f46', 'side_dark': '#064e3b', 'stud_top': '#6ee7b7', 'stud_side': '#10b981',
            'stud_side_dark': '#047857', 'highlight': '#a7f3d0', 'shadow': '#022c22'
        },
        'amber': {
            'top': '#fbbf24', 'top_grad': '#f59e0b', 'front': '#d97706', 'front_dark': '#b45309',
            'side': '#92400e', 'side_dark': '#78350f', 'stud_top': '#fde68a', 'stud_side': '#f59e0b',
            'stud_side_dark': '#b45309', 'highlight': '#fef3c7', 'shadow': '#451a03'
        },
        'red': {
            'top': '#f87171', 'top_grad': '#ef4444', 'front': '#dc2626', 'front_dark': '#b91c1c',
            'side': '#991b1b', 'side_dark': '#7f1d1d', 'stud_top': '#fca5a5', 'stud_side': '#ef4444',
            'stud_side_dark': '#b91c1c', 'highlight': '#fecaca', 'shadow': '#450a0a'
        },
        'violet': {
            'top': '#a78bfa', 'top_grad': '#8b5cf6', 'front': '#7c3aed', 'front_dark': '#6d28d9',
            'side': '#5b21b6', 'side_dark': '#4c1d95', 'stud_top': '#c4b5fd', 'stud_side': '#8b5cf6',
            'stud_side_dark': '#6d28d9', 'highlight': '#ddd6fe', 'shadow': '#2e1065'
        }
    }

    svg = []
    svg.append('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 600 280" width="100%" height="100%">')
    svg.append('  <defs>')
    svg.append('    <filter id="coverShadow" x="-15%" y="-15%" width="135%" height="145%">')
    svg.append('      <feDropShadow dx="0" dy="16" stdDeviation="14" flood-color="#0f172a" flood-opacity="0.18"/>')
    svg.append('    </filter>')

    # Gradientes para cada cor
    for cname, pal in palettes.items():
        svg.append(f'    <linearGradient id="top_{cname}" x1="0%" y1="0%" x2="100%" y2="100%">')
        svg.append(f'      <stop offset="0%" stop-color="{pal["top"]}"/>')
        svg.append(f'      <stop offset="100%" stop-color="{pal["top_grad"]}"/>')
        svg.append('    </linearGradient>')
        
        svg.append(f'    <linearGradient id="front_{cname}" x1="0%" y1="0%" x2="0%" y2="100%">')
        svg.append(f'      <stop offset="0%" stop-color="{pal["front"]}"/>')
        svg.append(f'      <stop offset="20%" stop-color="{pal["top_grad"]}"/>')
        svg.append(f'      <stop offset="40%" stop-color="{pal["front"]}"/>')
        svg.append(f'      <stop offset="100%" stop-color="{pal["front_dark"]}"/>')
        svg.append('    </linearGradient>')
        
        svg.append(f'    <linearGradient id="side_{cname}" x1="0%" y1="0%" x2="100%" y2="100%">')
        svg.append(f'      <stop offset="0%" stop-color="{pal["side"]}"/>')
        svg.append(f'      <stop offset="100%" stop-color="{pal["side_dark"]}"/>')
        svg.append('    </linearGradient>')
        
        svg.append(f'    <linearGradient id="studCyl_{cname}" x1="0%" y1="0%" x2="100%" y2="0%">')
        svg.append(f'      <stop offset="0%" stop-color="{pal["stud_top"]}"/>')
        svg.append(f'      <stop offset="30%" stop-color="{pal["stud_side"]}"/>')
        svg.append(f'      <stop offset="100%" stop-color="{pal["stud_side_dark"]}"/>')
        svg.append('    </linearGradient>')
        
        svg.append(f'    <linearGradient id="studCap_{cname}" x1="0%" y1="0%" x2="0%" y2="100%">')
        svg.append(f'      <stop offset="0%" stop-color="{pal["highlight"]}"/>')
        svg.append(f'      <stop offset="60%" stop-color="{pal["stud_top"]}"/>')
        svg.append(f'      <stop offset="100%" stop-color="{pal["stud_side"]}"/>')
        svg.append('    </linearGradient>')

    svg.append('  </defs>')
    svg.append('  <g filter="url(#coverShadow)">')

    def render_brick(x, y, w, h, dx, dy, color, studs_n, label=""):
        pal = palettes[color]
        res = []
        res.append(f'    <!-- Bloco {color} ({x},{y}) -->')
        res.append(f'    <g>')
        
        # 1. Face Superior
        res.append(f'      <polygon points="{x},{y} {x+dx},{y+dy} {x+w+dx},{y+dy} {x+w},{y}" fill="url(#top_{color})" stroke="{pal["highlight"]}" stroke-width="0.7"/>')
        # 2. Face Lateral Direita
        res.append(f'      <polygon points="{x+w},{y} {x+w+dx},{y+dy} {x+w+dx},{y+dy+h} {x+w},{y+h}" fill="url(#side_{color})" stroke="{pal["shadow"]}" stroke-width="0.7"/>')
        
        # 3. Studs (pinos)
        rx = 12.5
        ry = 5.2
        h_stud = 9.0
        
        # Interpolação para 2 fileiras se studs_n >= 4
        # Para um bloco 2x2: studs_n = 4 (2x2)
        # Para um bloco 2x4: studs_n = 8 (2x4)
        ncols = studs_n // 2
        for col in range(ncols):
            u = (col + 0.5) / ncols
            # Fileira de trás
            xb = x + dx * 0.3 + u * w
            yb = y + dy * 0.7
            # Fileira da frente
            xf = x + dx * 0.75 + u * w
            yf = y + dy * 0.25
            
            # Stud trás
            res.append(f'      <g transform="translate({xb:.1f}, {yb:.1f})">')
            res.append(f'        <ellipse cx="0" cy="0" rx="{rx}" ry="{ry}" fill="{pal["shadow"]}" opacity="0.3"/>')
            res.append(f'        <path d="M -{rx} -{h_stud} L -{rx} 0 A {rx} {ry} 0 0 0 {rx} 0 L {rx} -{h_stud} Z" fill="url(#studCyl_{color})"/>')
            res.append(f'        <ellipse cx="0" cy="-{h_stud}" rx="{rx}" ry="{ry}" fill="url(#studCap_{color})" stroke="{pal["highlight"]}" stroke-width="0.6"/>')
            res.append(f'        <ellipse cx="0" cy="-{h_stud}" rx="{rx*0.6:.1f}" ry="{ry*0.6:.1f}" fill="none" stroke="{pal["highlight"]}" stroke-width="1.0" opacity="0.9"/>')
            res.append(f'      </g>')
            
            # Stud frente
            res.append(f'      <g transform="translate({xf:.1f}, {yf:.1f})">')
            res.append(f'        <ellipse cx="0" cy="0" rx="{rx}" ry="{ry}" fill="{pal["shadow"]}" opacity="0.3"/>')
            res.append(f'        <path d="M -{rx} -{h_stud} L -{rx} 0 A {rx} {ry} 0 0 0 {rx} 0 L {rx} -{h_stud} Z" fill="url(#studCyl_{color})"/>')
            res.append(f'        <ellipse cx="0" cy="-{h_stud}" rx="{rx}" ry="{ry}" fill="url(#studCap_{color})" stroke="{pal["highlight"]}" stroke-width="0.6"/>')
            res.append(f'        <ellipse cx="0" cy="-{h_stud}" rx="{rx*0.6:.1f}" ry="{ry*0.6:.1f}" fill="none" stroke="{pal["highlight"]}" stroke-width="1.0" opacity="0.9"/>')
            res.append(f'      </g>')
            
        # 4. Face Frontal
        res.append(f'      <rect x="{x}" y="{y}" width="{w}" height="{h}" rx="2.5" fill="url(#front_{color})" stroke="{pal["front_dark"]}" stroke-width="0.8"/>')
        res.append(f'      <line x1="{x+1}" y1="{y+1}" x2="{x+w-1}" y2="{y+1}" stroke="{pal["highlight"]}" stroke-width="1.8" opacity="0.9"/>')
        res.append(f'      <line x1="{x+1}" y1="{y+h-1}" x2="{x+w-1}" y2="{y+h-1}" stroke="{pal["shadow"]}" stroke-width="2.0" opacity="0.7"/>')
        
        if label:
            res.append(f'      <text x="{x + w/2:.1f}" y="{y + h/2 + 6:.1f}" font-family="system-ui, -apple-system, sans-serif" font-size="14" font-weight="900" fill="#ffffff" text-anchor="middle" letter-spacing="2" style="text-shadow: 0 1px 3px rgba(0,0,0,0.5);">{label}</text>')
            
        res.append('    </g>')
        return '\n'.join(res)

    # Composição Modular Tridimensional Harmônica:
    # 4 Blocos em perspectiva formando a base conectada do LegoLoop
    dx = 38
    dy = -26
    
    # Bloco 1 (Base Esquerda): Verde Esmeralda (2x4)
    svg.append(render_brick(x=60, y=140, w=170, h=52, dx=dx, dy=dy, color='green', studs_n=6, label="BUILD"))
    
    # Bloco 2 (Base Direita): Azul Cobalto (2x4)
    svg.append(render_brick(x=270, y=140, w=170, h=52, dx=dx, dy=dy, color='blue', studs_n=6, label="SCALE"))
    
    # Bloco 3 (Degrau Superior Central Encaixado): Carmesim (2x4)
    svg.append(render_brick(x=165, y=75, w=170, h=52, dx=dx, dy=dy, color='red', studs_n=6, label="LEGOLOOP"))
    
    # Bloco 4 (Apoio Lateral Elevado / Conector): Âmbar Dourado (2x2)
    svg.append(render_brick(x=375, y=75, w=100, h=52, dx=dx, dy=dy, color='amber', studs_n=4, label="TDD"))

    svg.append('  </g>')
    svg.append('</svg>')

    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w', encoding='utf-8') as f:
        f.write('\n'.join(svg))
    print(f'Criado {filename}')

if __name__ == '__main__':
    create_lego_cover_svg('docs/books/ativos/lego_cover_blocks.svg')
    create_lego_cover_svg('docs/templates/ativos/lego_cover_blocks.svg')
    create_lego_cover_svg('docs/books/video-edition/ativos/lego_cover_blocks.svg')
