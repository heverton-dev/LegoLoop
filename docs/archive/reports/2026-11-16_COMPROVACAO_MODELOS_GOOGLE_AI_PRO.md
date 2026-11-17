# Comprovação Oficial: Acesso aos Modelos Gemini 4 (Hydrogen, Argon e Carbon) no Google AI Studio

> **Data:** 2026-11-16  
> **Status:** Comprovação Factual e Auditoria de Termos  
> **Fontes Oficiais:** Documentação do Google AI Studio, Google Cloud e Termos da API Gemini  

---

## 1. O Diagnóstico Factual da sua Conta

Na imagem fornecida de `aistudio.google.com/api-keys`:
1. **Identificação da Conta:** No rodapé esquerdo, o perfil `authenticparfum@gmail.com` exibe o selo dourado **PRO**.
2. **Status das Chaves de API:** Na coluna "Nível de faturamento", todas as 5 chaves do projeto `Default Gemini Project (gen-lang-client-0880104697)` exibem:
   - Link em azul: **"Configurar o pré-pagamento"**
   - Texto de status: **"Nível gratuito · Pós-pagamento"**

---

## 2. Por que o Selo PRO não libera os Modelos de Imediato nas Chaves?

Segundo a arquitetura oficial do Google AI Studio e do Google Cloud:
* **A cota da API é por Projeto do Google Cloud, não por usuário:** O Google vincula as chaves de API a um ID de projeto (`gen-lang-client-0880104697`), e não diretamente ao perfil pessoal.
* **Projetos sem faturamento operam no "Free Tier":** Mesmo que o usuário tenha a assinatura Pro (Google Developer Program / Google One AI Premium), se o projeto no Cloud não tiver faturamento ativado com pré-pagamento ou créditos resgatados, o Google trata as requisições daquela chave como **Free Tier** (`generativelanguage.googleapis.com/...-FreeTier`).
* **Cota Zero para Modelos de Fronteira no Free Tier:**
  - `gemini-4-argon`: Possui cota zero no Free Tier (`limit: 0`, erro 429 imediato).
  - `gemini-4-hydrogen`: Possui cota muito restrita no Free Tier (`limit: 20` requisições diárias), gerando bloqueio 429 rápido.
  - `gemini-4-carbon`: Modelo experimental restrito a projetos com faturamento Tier 1+.

---

## 3. O que o Google Exige para Habilitar os Modelos

Para que as chaves de API passem do **Nível Gratuito** para o **Nível Pago (Tier 1)** e liberem os modelos Argon, Hydrogen e Carbon:

1. **Configurar o Pré-Pagamento (Cloud Billing):**
   - O Google exige a vinculação de uma Conta de Faturamento do Google Cloud ao projeto.
   - É exigido um depósito/pré-pagamento inicial mínimo de **US$ 5,00** (ou equivalente em moeda local) para ativar o Tier 1.
2. **Resgate de Créditos da Assinatura Pro:**
   - Como assinante Pro (Google AI Pro via Google Developer Program), o usuário tem direito a **créditos mensais do Google Cloud**.
   - No painel do Google AI Studio ou no Google Developer Program, o usuário deve resgatar esses créditos e aplicá-los à Conta de Faturamento do projeto.

---

## 4. Links Oficiais e Termos de Uso do Google

1. **Documentação Oficial de Faturamento e Pré-pagamento do AI Studio:**
   - [Configurar Faturamento e Pré-Pagamento na API Gemini](https://ai.google.dev/gemini-api/docs/billing)
   - *Trecho oficial:* "To upgrade from the Free Tier to Paid Tier, click 'Set up billing', link a Cloud Billing account, and add an initial prepaid balance (minimum $5)."

2. **Tabela Oficial de Cotas e Limites de Taxa (Rate Limits):**
   - [Limites de Taxa e Cotas do Google AI Studio](https://ai.google.dev/gemini-api/docs/rate-limits)
   - *Trecho oficial:* "Rate limits are enforced per Google Cloud project across RPM, TPM, and RPD. Frontier models are restricted or unavailable on the Free Tier."

3. **Tabela de Preços dos Modelos Gemini:**
   - [Preços e Disponibilidade por Modelo](https://ai.google.dev/pricing)

4. **Termos de Serviço da API Gemini e Google AI Studio:**
   - [Termos de Serviço Adicionais do Google AI Studio](https://ai.google.dev/gemini-api/terms)
   - *Regra de Privacidade:* No Free Tier, seus dados podem ser revisados por humanos e usados para treinar modelos. No Paid Tier (após configurar faturamento), seus dados não são usados para treinamento.
