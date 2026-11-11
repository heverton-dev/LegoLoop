package board

// Status representa o estado de um Ticket no quadro Kanban do GitHub Project.
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

// Ticket representa um item mapeado de uma coluna do GitHub Project v2.
type Ticket struct {
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

type Card = Ticket

// Wave representa um lote de tickets que podem ser executados concorrentemente sem conflito.
type Wave struct {
	Index int      `json:"index"`
	Cards []Ticket `json:"tickets"`
}

// BoardState captura a fotografia instantânea das colunas do projeto.
type BoardState struct {
	TotalTickets int                 `json:"total_tickets"`
	Columns      map[Status][]Ticket `json:"columns"`
}

// BoardClient define a interface de interação com o GitHub Projects v2.
type BoardClient interface {
	FetchReadyTickets() ([]Ticket, error)
	MoveTicket(ticketID string, toStatus Status) error
	CheckQuota() (remaining int, resetTime int64, err error)
}
