package board

// Status representa o estado de um Card no quadro Kanban do GitHub Project.
type Status string

const (
	StatusBacklog  Status = "Backlog"
	StatusReady    Status = "Ready"
	StatusBuilding Status = "Building"
	StatusQA       Status = "QA"
	StatusReview   Status = "Review"
	StatusDone     Status = "Done"
	StatusBlocked  Status = "Blocked"
)

// Card representa uma issue/card mapeada de uma coluna do GitHub Project v2.
type Card struct {
	ID        string   `json:"id"`
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Status    Status   `json:"status"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Branch    string   `json:"branch,omitempty"`
	PRNumber  int      `json:"pr_number,omitempty"`
}

// Wave representa um lote de cards que podem ser executados concorrentemente sem conflito.
type Wave struct {
	Index int    `json:"index"`
	Cards []Card `json:"cards"`
}

// BoardState captura a fotografia instantânea das colunas do projeto.
type BoardState struct {
	TotalCards int               `json:"total_cards"`
	Columns    map[Status][]Card `json:"columns"`
}

// BoardClient define a interface de interação com o GitHub Projects v2.
type BoardClient interface {
	FetchReadyCards() ([]Card, error)
	MoveCard(cardID string, toStatus Status) error
	CheckQuota() (remaining int, resetTime int64, err error)
}
