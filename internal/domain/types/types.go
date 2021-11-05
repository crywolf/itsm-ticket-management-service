package types

// DateTime is RFC3339 time format
type DateTime string

// IsZero returns true if DateTime has zero value
func (t DateTime) IsZero() bool {
	return t == ""
}
