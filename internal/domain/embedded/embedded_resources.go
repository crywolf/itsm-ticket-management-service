package embedded

// Resource represents resource that is embedded in other domain object (used in '_embedded' hypermedia field)
type Resource string

func (a Resource) String() string {
	return string(a)
}

// Embedded resources values
const (
	CreatedBy Resource = "CreatedBy"
	UpdatedBy Resource = "UpdatedBy"
)
