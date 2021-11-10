package types

// DateTime is RFC3339 time format
type DateTime string

func (t DateTime) String() string {
	return string(t)
}

// IsZero returns true if DateTime has zero value
func (t DateTime) IsZero() bool {
	return t == ""
}
