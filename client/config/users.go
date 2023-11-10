package config

// UsersBundle defines name for the system users management bundle.
const UsersBundle Bundle = "users"

// Users adds or removes users.
//
// Example payload:
//
//	{
//	 "items": [
//	   {
//	     "username": "test",
//	     "action": "remove"
//	   }
//	 ]
//	}
type Users struct {
	Metadata

	// Users is a list of users to be added or removed.
	Users []User `json:"items,omitempty"`
}

// UserAction defines what to do with a user.
type UserAction string

const (
	UserAdd    UserAction = "add"
	UserRemove UserAction = "remove"
)

// User defines a user to be modified in the system.
type User struct {
	// Username of the user to be added or removed.
	Username string `json:"username"`

	// Action defines what to do with the user.
	Action UserAction `json:"action"`
}
