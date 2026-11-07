# PARTE IV — Referência Rápida

## 15. A Tríade de Triggers por Operação

Na coluna Terminal, todo comando é prefixado por `legoloop`; na coluna Chat, toda skill é prefixada por `/legoloop-`.

| Ação desejada | Nível | Terminal | Chat | LegoLoop Studio |
| :--- | :--- | :--- | :--- | :--- |
| **Criar processo** | Macro | `fabricar "<nome>"` | `fabricar` | Botão **+ Novo Processo** |
| **Consultar inventário** | Macro | `listar [--json]` | `listar` | Aba **Processos** |
| **Auditar reputação** | Meso | `matriz [--tier ouro]` | `matriz` | Aba **Matriz** |
| **Encaixar sub-loop** | Meso | `compor <pai> <filho>` | `compor` | Arrastar bloco e traçar aresta |
| **Substituir etapa** | Micro | `refinar "<alvo>"` | `refinar` | Botão direito no nó, **Refinar** |
| **Remover com rollback** | Macro | `remover "<nome>"` | `remover` | Menu do cartão, **Remover** |
| **Ajustar modelo e contexto** | Micro/Nano | Edição do `ROUTE.json` | Pedido livre | Painel de propriedades do nó |
| **Compilar livro em PDF** | Macro | `legopress "<alvo>"` | `legopress` | Botão **Exportar Livro Didático** |

## 16. Síntese

Com a Tríade de Triggers integrada ao padrão BDD e os encadeamentos canônicos mapeados, o LegoLoop deixa de ser apenas uma CLI:

* Em automação headless ou CI/CD, use o **Terminal**.
* Em programação pareada com um modelo de linguagem no editor, use o **Chat do Agente**.
* Para visão panorâmica e edição visual, use o **LegoLoop Studio**.

Em qualquer dos três caminhos, o comportamento contratado pelo BDD (**Dado-Quando-Então**) é idêntico, determinístico e auditável.
