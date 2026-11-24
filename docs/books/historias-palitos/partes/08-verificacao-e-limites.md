# 8. Verificação e limites

Este capítulo registra o que foi testado e o que ainda falta.
Cada teste mostra o comando e o resultado.

## 8.1 Testes executados

| Teste                              | Comando                                    | Resultado                           |
| :--------------------------------- | :----------------------------------------- | :---------------------------------- |
| Dependências                       | `legoloop legopress doctor`                  | pandoc e typst disponíveis (exit 0) |
| Preparação dos cortes first-brands | `python ferramentas/cortes.py cortes/first-brands.json` | 5 cortes preparados (exit 0)        |
| Preparação dos cortes curiosidades | `python ferramentas/cortes.py cortes/curiosidades-malucas.json` | 3 cortes preparados (exit 0)        |
| Render de um corte                 | `cortes.py --render --so 02-invasao-muscular` | 1080 x 1920, 31,7 s, 1 min 20 s     |

## 8.2 Lacunas encontradas

- `cortes.py` diz que confere o gancho. O código não confere o gancho (`ferramentas/cortes.py`).
- O JSON de curiosidades-malucas não tem link do vídeo completo nem status de envio (`cortes/curiosidades-malucas.json`).
- O corte 03-mordida-25cm ainda não foi renderizado.
- A correção de caminhos em `ferramentas/cortes.py` não está commitada.
- O comando `legoloop legopress criar` só imprime uma linha. Não cria a pasta do livro.

## 8.3 Limites conhecidos

- Os cenários foram desenhados em 16:9. No vertical, alguns elementos cortam ou ficam perto da legenda.
- A trilha do vídeo longo serve de base para os cortes. Há um fade curto na emenda.
- Cada corte leva alguns minutos para renderizar. Não cobre publicação real nem geração no Flow.
