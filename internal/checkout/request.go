package checkout

type ListCheckoutHistoriesPayload struct {
	CustomerID string `schema:"customerId" binding:"omitempty"`
	Limit      int    `schema:"limit" binding:"omitempty"`
	Offset     int    `schema:"offset" binding:"omitempty"`
	CreatedAt  string `schema:"createdAt" binding:"omitempty"`

	CreatedAtSearchType CreatedAtSearchType
}

type CreatedAtSearchType int

const (
	Ascending CreatedAtSearchType = iota
	Descending
	Ignore
)
