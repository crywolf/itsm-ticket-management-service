package memory

// Timelog stored in memory storage
type Timelog struct {
	ID string

	Remote bool

	Start string

	End string

	Work uint

	VisitSummary string

	CreatedAt string

	CreatedBy string

	UpdatedAt string

	UpdatedBy string
}
