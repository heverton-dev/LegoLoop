# Parte III — Micro: A Topologia DAG, Stop-Gates & Operação

## 8. Arquitetura do Pipeline DAG & Stop-Gates Humanos

## 8.1 A Estrutura Topológica do Grafo Acíclico

A execução da esteira `video-edition` é governada pelo motor matemático do LegoLoop através de um Grafo Acíclico Dirigido (**DAG**). Isso significa que cada tarefa é isolada, possui dependências explícitas declaradas no campo `depends_on`, e o motor calcula a ordenação topológica via Algoritmo de Kahn, garantindo paralelismo onde for possível e ordem estrita onde houver dependência de dados.

![Topologia do Pipeline DAG de Edição (12 Nós e 3 Stop-Gates Humanos)](img/diagrama_13_guias_dag.svg)

---

## 8.2 Tabela Mestra dos 12 Nós do Pipeline

A tabela a seguir apresenta o contrato de engenharia completo dos 12 nós distribuídos em 11 camadas:

| Camada         | ID do Nó               | Tipo           | Dependências     | Modelo / Runtime | Contrato de Saída (Artefato)        |
| :------------: | :--------------------- | :------------: | :--------------: | :--------------: | :---------------------------------- |
| **C1**             | `n01_iniciar_preparar`     | `script`         | —                | `code-fast`        | `scratch/audio_48k.wav`                |
| **C2**             | `n02_transcricao`         | `task`           | `n01`              | Whisper Local    | `scratch/transcription_words.json`     |
| **C3**             | `n03_decisao_cortes`       | `skill`          | `n02`              | `code-pro`         | `scratch/cuts_manifest.json`           |
| **C4**             | `n04_sg_revisao_corte`      | `stop_gate`       | `n03`              | Humano           | Aprovação Interativa do Corte       |
| **C5**             | `n05_motion_mapa`          | `skill`          | `n04`              | `code-pro`         | `scratch/motion_storyboard.json`       |
| **C6**             | `n06_sg_revisao_motion`     | `stop_gate`       | `n05`              | Humano           | Aprovação Interativa do Motion      |
| **C7**             | `n07_camera_yunet`         | `script`         | `n06`              | OpenCV (CPU)     | `scratch/yunet_tracking.json`          |
| **C7**             | `n08_assets_download`      | `task`           | `n06`              | `code-fast`        | `output/assets_manifest.json`          |
| **C8**             | `n09_motion_render`        | `task`           | `n07`, `n08`         | Remotion (Node)  | `scratch/motion_alpha_bundle/`          |
| **C9**             | `n10_efeitos_sonoros`      | `script`         | `n09`              | FFmpeg Audio     | `scratch/audio_mix_ducked.wav`          |
| **C10**            | `n11_entrega_final`        | `task`           | `n10`              | FFmpeg + FCPXML  | `dist/timeline.fcpxml` + `MP4`          |
| **C11**            | `n12_sg_entrega_final`      | `stop_gate`       | `n11`              | Humano           | Validação e Liberação da Entrega    |

---

## 8.3 O Papel Crítico dos 3 Stop-Gates Humanos

A automação cega é a causa principal da perda de qualidade em produções com IA. O LegoLoop impõe 3 barreiras onde o processo **pausa compulsoriamente** no terminal ou na interface do Studio, exigindo a decisão deliberada do operador:

#painel-alerta("STOP-GATE 01: REVISÃO DE CORTES (n04)", [
  O operador abre o arquivo `scratch/cuts_manifest.json` e visualiza o texto com as marcações de corte e notas de gravação expurgadas. Ele pode remover manualmente qualquer frase que o modelo tenha selecionado incorretamente antes de gerar mídia.
])

#painel-alerta("STOP-GATE 02: REVISÃO DO STORYBOARD (n06)", [
  O operador revisa o manifesto `scratch/motion_storyboard.json`. Ele valida a tipografia, os textos dos cartões e os momentos de entrada, impedindo que a máquina gaste poder computacional renderizando animações com erros de digitação.
])

#painel-alerta("STOP-GATE 03: APROVAÇÃO FINAL DA LINHA DO TEMPO (n12)", [
  O operador abre a timeline gerada `dist/timeline.fcpxml` no DaVinci Resolve ou Premiere e dá o play no arquivo `dist/render_final.mp4`. Aprovada a sincronia e a fluidez, o processo é arquivado como concluído.
])
