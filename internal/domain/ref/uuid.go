package ref

// UUID represents UUID of a resource
// swagger:strfmt uuid
type UUID string

func (u UUID) String() string {
	return string(u)
}

// IsZero returns true if UUID has zero value
func (u UUID) IsZero() bool {
	return u == ""
}

// ChannelID represents UUID of a channel
// swagger:strfmt uuid
type ChannelID UUID

func (u ChannelID) String() string {
	return string(u)
}

// IsZero returns true if ChannelID has zero value
func (u ChannelID) IsZero() bool {
	return u == ""
}

// ExternalUserUUID represents UUID of a user (in external service)
// swagger:strfmt uuid
type ExternalUserUUID UUID

func (u ExternalUserUUID) String() string {
	return string(u)
}

// IsZero returns true if ExternalUserUUID has zero value
func (u ExternalUserUUID) IsZero() bool {
	return u == ""
}
