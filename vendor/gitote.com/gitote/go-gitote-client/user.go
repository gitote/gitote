package gitote

import (
	"encoding/json"
	"fmt"
)

// User represents a API user.
type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"login"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

// MarshalJSON implements the json.Marshaler interface for User
func (u User) MarshalJSON() ([]byte, error) {
	// Re-declaring User to avoid recursion
	type shadow User
	return json.Marshal(struct {
		shadow
		// LEGACY [Gitote 1.0]: remove field(s) for backward compatibility
		CompatUserName string `json:"username"`
	}{shadow(u), u.UserName})
}

func (c *Client) GetUserInfo(user string) (*User, error) {
	u := new(User)
	err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s", user), nil, nil, u)
	return u, err
}
