## 4. Câmera Inteligente YuNet & Linha dos Olhos Estabilizada (35%)

## 4.1 A Ilusão do Zoom Dinâmico Convencional

Em produções modernas para YouTube e redes sociais, o corte de câmera (zoom in em momentos de ênfase e zoom out em introduções) é amplamente utilizado para quebrar a monotonia visual de um plano fixo. No entanto, a maioria dos plugins comerciais aplica reenquadramento cego no centro geométrico do vídeo (50% de altura e largura).

Quando um apresentador se mexe levemente para o lado ou se inclina para frente, o zoom estático corta o topo da cabeça ou posiciona os olhos muito abaixo da linha de leitura do espectador. Pior ainda: se um cartão gráfico animado é inserido na lateral, o apresentador com zoom centralizado colide visualmente com o texto.

---

## 4.2 A Descoberta da Ordem: Por que o Enquadramento é Pós-Motion?

Uma das inovações mais elegantes da metodologia de Lucas Leite reside na **ordem estrita de dependência** dos nós:

#painel-alerta("A ORDEM IMPORTA: MOTION PRIMEIRO, CÂMERA DEPOIS", [
  O cálculo do enquadramento de câmera só é executado no nó `n07_camera_yunet` **após** a geração do mapa de motion no nó `n05_motion_mapa`. Se a cena exige um cartão gráfico na lateral direita, o algoritmo de câmera sabe antecipadamente que o apresentador deve ser deslocado suavemente para o terço esquerdo da tela, garantindo que o texto e a pessoa coexistam sem sobreposição.
])

![Câmera Inteligente YuNet e Motion Procedural sem Colisão](img/diagrama_yunet_eyelevel_remotion.svg)

---

## 4.3 A Regra Áurea dos 35% (Eye-Level Ancorado)

Na cinematografia clássica e na transmissão televisiva profissional, o equilíbrio visual de um plano médio (*medium close-up*) exige que a linha dos olhos do apresentador permaneça ancorada em aproximadamente um terço da altura total do quadro a partir do topo (**35%**):

$$Y_{\text{olhos\_alvo}} = 0.35 \times \text{Altura\_Quadro}$$

Para um vídeo em alta definição (1080p, altura de 1080 pixels), os olhos devem residir estavelmente na coordenada vertical $Y = 378\text{ px}$.

### Rastreamento com OpenCV YuNet
Para alcançar essa precisão sem queimar recursos da GPU com redes neurais pesadas, o pipeline emprega o detector **OpenCV YuNet**, uma rede de detecção facial de altíssima eficiência que roda em CPU local a mais de 120 frames por segundo:

```python
### Trecho representativo do cálculo de enquadramento afim
import cv2

def calcular_enquadramento_camera(frame, face_box, landmarks, largura_alvo=1920, altura_alvo=1080):
    olho_esquerdo = landmarks[0]
    olho_direito = landmarks[1]
    olhos_centro_y = (olho_esquerdo[1] + olho_direito[1]) / 2.0
    
    # Altura desejada: 35% do topo
    y_alvo = altura_alvo * 0.35
    deslocamento_vertical = y_alvo - olhos_centro_y
    
    # Aplica suavização com amortecimento (damping) para evitar tremores
    return aplicar_matriz_afim(frame, deslocamento_vertical, escala=1.15)
```

#painel-sucesso("RESULTADO VISUAL ESTÁVEL E CONFORTÁVEL", [
  Mesmo que o apresentador levante a cadeira ou incline o corpo durante a gravação, a câmera reajusta a escala e a posição de forma suave e ponderada. O corte de câmera ganha a sensação de ter sido operado por um diretor de fotografia experiente, transmitindo autoridade e profissionalismo.
])
