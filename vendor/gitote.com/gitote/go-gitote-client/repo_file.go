package gitote

import (
	"fmt"
)

// GetFile downloads a file of repository, ref can be branch/tag/commit.
// e.g.: ref -> master, tree -> macaron.go(no leading slash)
func (c *Client) GetFile(user, repo, ref, tree string) ([]byte, error) {
	return c.getResponse("GET", fmt.Sprintf("/repos/%s/%s/raw/%s/%s", user, repo, ref, tree), nil, nil)
}
