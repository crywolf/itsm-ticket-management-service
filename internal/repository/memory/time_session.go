package memory

// TimeSession stored in memory storage
type TimeSession struct {
	ID string

	State string

	Incidents []IncidentInfo

	Work uint

	Travel uint

	TravelBack uint

	TravelDistanceInTravelUnits uint

	CreatedAt string

	CreatedBy string

	UpdatedAt string

	UpdatedBy string
}

// IncidentInfo stored in repository
type IncidentInfo struct {
	ID                 string
	HasSupplierProduct bool
}
