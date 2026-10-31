# Riscos

| Risco | Probabilidade | Impacto | Mitigação |
| --- | --- | --- | --- |
| Bloqueio ou oscilação na API oficial do WhatsApp (Meta Cloud API) | Média | Alto | Implementar fila de retentativas assíncronas com fallback para envio de alertas via e-mail e SMS |
| Não conformidade com a LGPD no tratamento de dados cadastrais e médicos | Baixa | Alto | Aplicar criptografia em repouso e trânsito, controle de acesso baseado em papéis (RBAC) e política de expiração de dados |
| Conflito de agendamento simultâneo para o mesmo veterinário | Média | Médio | Utilizar transações com locks pessimistas no PostgreSQL durante a reserva de horários |
| Perda de histórico clínico por falha de armazenamento ou integridade no banco | Baixa | Crítico | Configurar backups automatizados diários com Point-in-Time Recovery (PITR) e chaves estrangeiras com restrições rígidas |
| Resistência de recepcionistas e veterinários ao abandono do papel | Alta | Médio | Projetar interface simplificada com preenchimento assistido e realizar sessões de treinamento prático antes da migração total |
| Falha no upload e renderização de anexos pesados de exames dos pets | Média | Baixo | Adotar validação de formato e compressão assíncrona de arquivos com armazenamento direto em bucket S3 |