```markdown
# SPEC — Sistema de Gestão para Clínicas Veterinárias

## Visão Geral

Sistema web para gestão de agendamentos, prontuários eletrônicos e comunicação automatizada via WhatsApp em clínicas veterinárias. Substituir gestão manual em papel/planilhas, reduzir faltas por esquecimento e centralizar histórico clínico dos animais.

## Escopo

### Funcionalidades Incluídas

**Gestão de Tutores e Pacientes**
- Cadastro de tutores com CPF, nome, telefone (WhatsApp), e-mail
- Cadastro de pets vinculados ao tutor: nome, espécie, raça, sexo, data de nascimento, peso, foto
- Associação de múltiplos animais a um único tutor
- Busca e listagem de tutores e pets

**Agendamento**
- Calendário interativo com visualização semanal/diária
- Criação de consulta: seleção de pet, veterinário, data, horário, tipo de consulta
- Reagendamento e cancelamento com registro de histórico
- Bloqueio de horários indisponíveis
- Painel de controle diário com fila de atendimento e status de confirmação

**Notificações Automatizadas**
- Disparo via WhatsApp Cloud API:
  - 24 horas antes da consulta
  - 2 horas antes da consulta
- Mensagem com dados da consulta, veterinário, horário
- Link para confirmação de presença
- Registro de status de envio e confirmação

**Prontuário Eletrônico**
- Ficha clínica por consulta: anamnese, diagnóstico, prescrição, observações
- Histórico completo de consultas do pet ordenado cronologicamente
- Registro de peso em cada consulta com gráfico de evolução
- Carteira de vacinação digital: vacina, data aplicação, próxima dose, lote
- Upload de exames (PDF, imagens) com data e tipo
- Acesso rápido ao histórico durante atendimento

**Controle de Acesso**
- Perfis: administrador, recepcionista, veterinário, tutor
- Administrador: acesso total
- Recepcionista: agendamentos, cadastros, visualização de prontuários
- Veterinário: prontuários, visualização de agenda própria
- Tutor: visualização de dados dos próprios pets, agendamentos futuros (somente leitura)

### Fora de Escopo

- Telemedicina por videoconferência
- Controle financeiro, fluxo de caixa, emissão de nota fiscal
- E-commerce de produtos/medicamentos
- Aplicativo mobile nativo iOS/Android
- Integração com sistemas de laboratórios externos
- Prescrição eletrônica com assinatura digital certificada ICP-Brasil

## Restrições Técnicas

- Aplicação web responsiva (desktop, tablet, smartphone)
- Compatibilidade: Chrome 90+, Safari 14+, Firefox 88+
- WhatsApp Cloud API oficial (Meta Business Platform) — proibido uso de APIs não oficiais
- Banco de dados PostgreSQL 14+
- Armazenamento em nuvem com criptografia em repouso
- Conformidade LGPD: consentimento, anonimização, exclusão de dados
- Backup automatizado diário
- HTTPS obrigatório em produção

## Modelo de Dados

### Entidades Principais

**Tutor**
- id (PK)
- cpf (único, indexado)
- nome
- telefone_whatsapp (formato: +5511999999999)
- email
- endereco
- data_cadastro
- consentimento_lgpd (booleano)

**Pet**
- id (PK)
- tutor_id (FK → Tutor)
- nome
- especie (enum: Cão, Gato, Ave, Réptil, Roedor, Outro)
- raca
- sexo (enum: Macho, Fêmea)
- data_nascimento
- peso_atual (decimal)
- foto_url
- ativo (booleano)

**Veterinario**
- id (PK)
- usuario_id (FK → Usuario)
- crmv
- nome
- especialidades (array)
- ativo (booleano)

**Agendamento**
- id (PK)
- pet_id (FK → Pet)
- veterinario_id (FK → Veterinario)
- data_hora_inicio (timestamp)
- data_hora_fim (timestamp)
- tipo_consulta (enum: Consulta, Retorno, Vacinação, Cirurgia, Exame)
- status (enum: Agendado, Confirmado, Realizado, Cancelado, Falta)
- observacoes
- criado_por (FK → Usuario)
- data_criacao

**Confirmacao**
- id (PK)
- agendamento_id (FK → Agendamento)
- tipo_lembrete (enum: 24h, 2h)
- enviado_em (timestamp)
- status_envio (enum: Enviado, Falha, Pendente)
- confirmado_em (timestamp, nullable)
- token_confirmacao (único)

**Prontuario**
- id (PK)
- agendamento_id (FK → Agendamento)
- pet_id (FK → Pet)
- veterinario_id (FK → Veterinario)
- data_atendimento (timestamp)
- peso (decimal)
- anamnese (texto)
- diagnostico (texto)
- prescricao (texto)
- observacoes (texto)

**Vacina**
- id (PK)
- pet_id (FK → Pet)
- nome_vacina
- data_aplicacao (date)
- proxima_dose (date, nullable)
- lote
- veterinario_id (FK → Veterinario)

**Exame**
- id (PK)
- pet_id (FK → Pet)
- prontuario_id (FK → Prontuario, nullable)
- tipo_exame
- data_exame (date)
- arquivo_url
- observacoes

**Usuario**
- id (PK)
- email (único)
- senha_hash
- perfil (enum: Administrador, Recepcionista, Veterinario, Tutor)
- ativo (booleano)
- tutor_id (FK → Tutor, nullable)

### Relacionamentos

- Tutor 1:N Pet
- Pet 1:N Agendamento
- Veterinario 1:N Agendamento
- Agendamento 1:1 Prontuario
- Agendamento 1:N Confirmacao
- Pet 1:N Vacina
- Pet 1:N Exame
- Usuario 1:1 Veterinario (quando perfil = Veterinario)
- Usuario 1:1 Tutor (quando perfil = Tutor)

### Índices Obrigatórios

- Tutor.cpf (único)
- Pet.tutor_id
- Agendamento.data_hora_inicio
- Agendamento.veterinario_id
- Agendamento.pet_id
- Confirmacao.agendamento_id
- Confirmacao.token_confirmacao (único)
- Prontuario.pet_id
- Vacina.pet_id
- Exame.pet_id

## Critérios de Aceitação

- [ ] Sistema permite cadastro de tutor com CPF, nome, telefone WhatsApp e vinculação de mínimo 1 pet
- [ ] Recepcionista consegue agendar consulta selecionando pet, veterinário, data e horário em calendário interativo
- [ ] Sistema dispara automaticamente 2 notificações WhatsApp (24h e 2h antes) para cada agendamento confirmado
- [ ] Mensagem de lembrete contém link de confirmação único que atualiza status do agendamento para "Confirmado" quando clicado
- [ ] Veterinário acessa prontuário do pet durante consulta e registra anamnese, diagnóstico, prescrição e peso em até 3 cliques
- [ ] Sistema exibe histórico completo de consultas de um pet ordenado por data decrescente com acesso aos prontuários anteriores
- [ ] Tutor com perfil "Tutor" visualiza somente dados dos próprios pets e não tem acesso a prontuários de outros animais
- [ ] Sistema permite upload de mínimo 1 arquivo de exame (PDF ou imagem) vinculado ao pet e exibe na ficha clínica
- [ ] Painel de controle diário lista todos os agendamentos do dia com status (Agendado/Confirmado/Realizado) e permite filtro por veterinário
- [ ] Sistema funciona responsivamente em navegador mobile (largura 375px) permitindo visualização de agenda e cadastro de consultas
```