# Briefing

## Problema
Clínicas veterinárias perdem agendamentos por esquecimento de tutores. Gestão manual em papel ou planilhas descentralizadas causa perda de histórico clínico, atrasos no atendimento e falhas no acompanhamento de vacinas e tratamentos contínuos dos animais.

## Público-alvo
- Recepcionistas e administradores de clínicas veterinárias.
- Médicos veterinários responsáveis pelo atendimento clínico.
- Tutores de animais domésticos atendidos pela clínica.

## Funcionalidades essenciais
- **Agendamento online**: Calendário interativo para marcação, reagendamento e cancelamento de consultas por horário e profissional.
- **Disparo de lembretes via WhatsApp**: Notificações automáticas enviadas 24 horas e 2 horas antes da consulta com links de confirmação.
- **Prontuário eletrônico do pet**: Ficha cadastral do animal contendo histórico de consultas, prescrições, peso, exames anexados e carteira de vacinação digital.
- **Cadastro de clientes e pacientes**: Associação de múltiplos animais ao perfil de um único tutor com dados de contato.
- **Painel de controle**: Visão diária da fila de atendimento e status das confirmações de presença.

## Fora de escopo
- Sistema de telemedicina veterinária por videoconferência.
- Módulo de controle financeiro completo, fluxo de caixa e emissão de notas fiscais.
- Venda online de produtos, rações ou medicamentos (e-commerce).
- Aplicativo mobile nativo para iOS ou Android nesta primeira fase.

## Restrições técnicas
- Aplicação web responsiva compatível com navegadores modernos (Chrome, Safari, Firefox).
- Integração obrigatória com WhatsApp via API oficial (Meta Business Platform / Cloud API) para garantir estabilidade e evitar banimentos de número.
- Armazenamento em nuvem com conformidade à LGPD para proteção dos dados pessoais dos tutores.
- Banco de dados relacional (PostgreSQL) para garantir a integridade referencial entre tutores, pets e prontuários médicos.