package memory

// Incident stored in memory storage
type Incident struct {
	ID string

	Number string

	ExternalID string

	ShortDescription string

	Description string

	FieldEngineerID string

	State string

	Timelogs []string

	CreatedAt string

	CreatedBy string

	UpdatedAt string

	UpdatedBy string
}
