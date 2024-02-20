package client

// UserBaseInfo is the base information of a user.
type UserBaseInfo struct {
	// ID is the unique identifier of the user.
	ID string `json:"user_id"`

	// FirstName is the first name of the user.
	FirstName string `json:"fname"`

	// LastName is the last name of the user.
	LastName string `json:"lname"`
}
